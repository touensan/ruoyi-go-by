package monitorcontroller

import (
	"context"
	"fmt"
	"net"
	"os"
	"runtime"
	"ruoyi-go/app/token"
	rediskey "ruoyi-go/common/types/redis-key"
	"ruoyi-go/config"
	"ruoyi-go/framework/dal"
	"ruoyi-go/framework/response"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

type ServerController struct{}
type CacheController struct{}
type OnlineController struct{}

var serverStartTime = time.Now()

func (*ServerController) Info(ctx *gin.Context) {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	total, free := systemMemory()
	used := total - free
	memUsage := percent(float64(used), float64(total))
	heapTotal := float64(mem.Sys) / 1024 / 1024
	heapUsed := float64(mem.Alloc) / 1024 / 1024

	response.NewSuccess().SetData("data", gin.H{
		"cpu": gin.H{
			"cpuNum": runtime.NumCPU(),
			"used":   0,
			"sys":    0,
			"free":   100,
		},
		"mem": gin.H{
			"total": roundGB(total),
			"used":  roundGB(used),
			"free":  roundGB(free),
			"usage": memUsage,
		},
		"jvm": gin.H{
			"name":      "Go Runtime",
			"version":   runtime.Version(),
			"home":      runtime.GOROOT(),
			"total":     fmt.Sprintf("%.2f", heapTotal),
			"used":      fmt.Sprintf("%.2f", heapUsed),
			"free":      fmt.Sprintf("%.2f", maxFloat(heapTotal-heapUsed, 0)),
			"usage":     percent(heapUsed, heapTotal),
			"startTime": serverStartTime.Format("2006-01-02 15:04:05"),
			"runTime":   time.Since(serverStartTime).Round(time.Second).String(),
			"inputArgs": strings.Join(os.Args, " "),
		},
		"sys": gin.H{
			"computerName": hostname(),
			"computerIp":   firstIP(),
			"osName":       runtime.GOOS,
			"osArch":       runtime.GOARCH,
			"userDir":      workingDir(),
		},
		"sysFiles": diskInfo(),
	}).Json(ctx)
}

func (*CacheController) Info(ctx *gin.Context) {
	info := redisInfo(ctx.Request.Context())
	dbSize, _ := dal.Redis.DBSize(ctx.Request.Context()).Result()
	response.NewSuccess().SetData("data", gin.H{
		"info":         info,
		"dbSize":       dbSize,
		"commandStats": []gin.H{},
	}).Json(ctx)
}

func (*CacheController) Names(ctx *gin.Context) {
	names := []gin.H{
		{"cacheName": rediskey.SysConfigKey, "remark": "参数配置"},
		{"cacheName": rediskey.SysDictKey, "remark": "字典数据"},
		{"cacheName": rediskey.UserTokenKey, "remark": "登录用户"},
		{"cacheName": rediskey.CaptchaCodeKey, "remark": "验证码"},
		{"cacheName": rediskey.LoginPasswordErrorKey, "remark": "登录错误次数"},
	}
	response.NewSuccess().SetData("data", names).Json(ctx)
}

func (*CacheController) Keys(ctx *gin.Context) {
	cacheName := ctx.Param("cacheName")
	keys, err := scanRedisKeys(ctx.Request.Context(), cacheName+"*")
	if err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}
	response.NewSuccess().SetData("data", keys).Json(ctx)
}

func (*CacheController) Value(ctx *gin.Context) {
	cacheName := ctx.Param("cacheName")
	cacheKey := ctx.Param("cacheKey")
	value, err := dal.Redis.Get(ctx.Request.Context(), cacheKey).Result()
	if err != nil {
		value = ""
	}
	response.NewSuccess().SetData("data", gin.H{
		"cacheName":  cacheName,
		"cacheKey":   cacheKey,
		"cacheValue": value,
		"remark":     "",
	}).Json(ctx)
}

func (*CacheController) ClearName(ctx *gin.Context) {
	keys, err := scanRedisKeys(ctx.Request.Context(), ctx.Param("cacheName")+"*")
	if err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}
	if len(keys) > 0 {
		if err = dal.Redis.Del(ctx.Request.Context(), keys...).Err(); err != nil {
			response.NewError().SetMsg(err.Error()).Json(ctx)
			return
		}
	}
	response.NewSuccess().Json(ctx)
}

func (*CacheController) ClearKey(ctx *gin.Context) {
	if err := dal.Redis.Del(ctx.Request.Context(), ctx.Param("cacheKey")).Err(); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}
	response.NewSuccess().Json(ctx)
}

func (*CacheController) ClearAll(ctx *gin.Context) {
	if err := dal.Redis.FlushDB(ctx.Request.Context()).Err(); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}
	response.NewSuccess().Json(ctx)
}

func (*OnlineController) List(ctx *gin.Context) {
	keys, err := scanRedisKeys(ctx.Request.Context(), rediskey.UserTokenKey+"*")
	if err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}
	rows := make([]gin.H, 0, len(keys))
	for _, key := range keys {
		var auth token.UserTokenResponse
		if err := dal.Redis.Get(ctx.Request.Context(), key).Scan(&auth); err != nil {
			continue
		}
		rows = append(rows, gin.H{
			"tokenId":       strings.TrimPrefix(key, rediskey.UserTokenKey),
			"deptName":      auth.DeptName,
			"userName":      auth.UserName,
			"ipaddr":        "",
			"loginLocation": "",
			"browser":       "",
			"os":            "",
			"loginTime":     auth.ExpireTime.Format("2006-01-02 15:04:05"),
		})
	}
	response.NewSuccess().SetPageData(rows, len(rows)).Json(ctx)
}

func (*OnlineController) ForceLogout(ctx *gin.Context) {
	if err := dal.Redis.Del(ctx.Request.Context(), rediskey.UserTokenKey+ctx.Param("tokenId")).Err(); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}
	response.NewSuccess().Json(ctx)
}

func scanRedisKeys(ctx context.Context, pattern string) ([]string, error) {
	var cursor uint64
	keys := make([]string, 0)
	for {
		batch, next, err := dal.Redis.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			return nil, err
		}
		keys = append(keys, batch...)
		cursor = next
		if cursor == 0 {
			return keys, nil
		}
	}
}

func redisInfo(ctx context.Context) map[string]string {
	infoText, err := dal.Redis.Info(ctx).Result()
	if err != nil {
		return map[string]string{}
	}
	info := make(map[string]string)
	for _, line := range strings.Split(infoText, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, value, ok := strings.Cut(line, ":")
		if ok {
			info[key] = value
		}
	}
	return info
}

func systemMemory() (uint64, uint64) {
	data, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		return 0, 0
	}
	values := make(map[string]uint64)
	for _, line := range strings.Split(string(data), "\n") {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		value, _ := strconv.ParseUint(fields[1], 10, 64)
		values[strings.TrimSuffix(fields[0], ":")] = value * 1024
	}
	total := values["MemTotal"]
	free := values["MemAvailable"]
	if free == 0 {
		free = values["MemFree"]
	}
	return total, free
}

func diskInfo() []gin.H {
	paths := []string{"/"}
	if config.Data.Ruoyi.UploadPath != "" {
		paths = append(paths, strings.TrimRight(config.Data.Ruoyi.UploadPath, "/"))
	}
	items := make([]gin.H, 0, len(paths))
	seen := make(map[string]bool)
	for _, path := range paths {
		if path == "" || seen[path] {
			continue
		}
		seen[path] = true
		var stat syscall.Statfs_t
		if err := syscall.Statfs(path, &stat); err != nil {
			continue
		}
		total := stat.Blocks * uint64(stat.Bsize)
		free := stat.Bavail * uint64(stat.Bsize)
		used := total - free
		items = append(items, gin.H{
			"dirName":     path,
			"sysTypeName": "linux",
			"typeName":    "local",
			"total":       byteSize(total),
			"free":        byteSize(free),
			"used":        byteSize(used),
			"usage":       percent(float64(used), float64(total)),
		})
	}
	return items
}

func hostname() string {
	name, _ := os.Hostname()
	return name
}

func workingDir() string {
	dir, _ := os.Getwd()
	return dir
}

func firstIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, addr := range addrs {
		ipNet, ok := addr.(*net.IPNet)
		if !ok || ipNet.IP.IsLoopback() {
			continue
		}
		if ip := ipNet.IP.To4(); ip != nil {
			return ip.String()
		}
	}
	return "127.0.0.1"
}

func roundGB(value uint64) string {
	return fmt.Sprintf("%.2f", float64(value)/1024/1024/1024)
}

func byteSize(value uint64) string {
	gb := float64(value) / 1024 / 1024 / 1024
	if gb >= 1 {
		return fmt.Sprintf("%.2f GB", gb)
	}
	return fmt.Sprintf("%.2f MB", float64(value)/1024/1024)
}

func percent(used, total float64) float64 {
	if total <= 0 {
		return 0
	}
	value, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", used/total*100), 64)
	return value
}

func maxFloat(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}
