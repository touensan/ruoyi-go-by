package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"ruoyi-go/app/router"
	"ruoyi-go/app/service"
	"ruoyi-go/config"
	"ruoyi-go/framework/dal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func main() {

	// dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := config.Data.Mysql.Username + ":" + config.Data.Mysql.Password + "@tcp(" + config.Data.Mysql.Host + ":" + strconv.Itoa(config.Data.Mysql.Port) + ")/" + config.Data.Mysql.Database + "?charset=" + config.Data.Mysql.Charset + "&parseTime=True&loc=Local"

	// 初始化数据访问层
	dal.InitDal(&dal.Config{
		GomrConfig: &dal.GomrConfig{
			Dialector: mysql.New(mysql.Config{
				DSN:                       dsn,
				DefaultStringSize:         191,  // 兼容 MySQL 5.6/5.7 utf8mb4 索引长度。
				DontSupportRenameIndex:    true, // 兼容 MySQL 5.7 及以下。
				DontSupportRenameColumn:   true, // 兼容 MySQL 5.7 及以下。
				SkipInitializeWithVersion: false,
			}),
			Opts: &gorm.Config{
				SkipDefaultTransaction: true, // 跳过默认事务
				NamingStrategy: schema.NamingStrategy{
					SingularTable: true,
				},
				Logger: logger.New(log.Default(), logger.Config{
					// LogLevel: logger.Silent, // 不打印日志
					LogLevel:                  logger.Error, // 打印错误日志
					IgnoreRecordNotFoundError: true,
				}),
			},
			MaxOpenConns: config.Data.Mysql.MaxOpenConns,
			MaxIdleConns: config.Data.Mysql.MaxIdleConns,
		},
		RedisConfig: &dal.RedisConfig{
			Host:     config.Data.Redis.Host,
			Port:     config.Data.Redis.Port,
			Database: config.Data.Redis.Database,
			Password: config.Data.Redis.Password,
		},
	})

	// Production releases must migrate explicitly, before either version starts.
	if os.Getenv("RUOYI_AUTO_SETUP") != "false" {
		if err := service.EnsureSystemSettingBaseline(); err != nil {
			log.Fatalf("初始化系统配置失败: %v", err)
		}
	}

	// 设置模式
	listenHost := ""
	if port := os.Getenv("RUOYI_PORT"); port != "" {
		parsed, err := strconv.Atoi(port)
		if err != nil || parsed < 1 || parsed > 65535 {
			log.Fatal("RUOYI_PORT 无效")
		}
		config.Data.Server.Port = parsed
		config.Data.Server.Host = "127.0.0.1"
		listenHost = "127.0.0.1"
	}
	gin.SetMode(config.Data.Server.Mode)

	// 初始化gin
	server := gin.New()

	// 使用恢复中间件
	server.Use(gin.Recovery())
	server.GET("/internal/ready", func(ctx *gin.Context) {
		host, _, err := net.SplitHostPort(ctx.Request.RemoteAddr)
		if err != nil || net.ParseIP(host) == nil || !net.ParseIP(host).IsLoopback() {
			ctx.Status(http.StatusNotFound)
			return
		}
		check, cancel := context.WithTimeout(ctx.Request.Context(), 2*time.Second)
		defer cancel()
		db, err := dal.Gorm.DB()
		if err == nil {
			err = db.PingContext(check)
		}
		if err == nil {
			err = dal.Redis.Ping(check).Err()
		}
		if err != nil {
			ctx.Status(http.StatusServiceUnavailable)
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"status": "ok", "version": config.Data.Ruoyi.Version})
	})

	// 托管 RuoYi-Vue3-ts 管理后台。前端使用 history 路由，刷新 /admin/* 时回退到 index.html。
	server.Any("/", func(ctx *gin.Context) {
		ctx.Redirect(http.StatusFound, "/admin/")
	})
	server.Static("/admin", "web/admin")
	server.NoRoute(func(ctx *gin.Context) {
		if strings.HasPrefix(ctx.Request.URL.Path, "/admin") {
			ctx.File("web/admin/index.html")
			return
		}
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": "not found"})
	})

	// 设置上传文件目录。配置值作为磁盘目录使用，URL 路由补齐为绝对路径。
	uploadRoot := strings.Trim(config.Data.Ruoyi.UploadPath, "/")
	if uploadRoot == "" {
		uploadRoot = "runtime/uploads"
	}
	server.Static("/"+uploadRoot, uploadRoot)

	// 注册路由
	router.Register(server)

	httpServer := &http.Server{Addr: net.JoinHostPort(listenHost, strconv.Itoa(config.Data.Server.Port)), Handler: server}
	stopped := make(chan os.Signal, 1)
	signal.Notify(stopped, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-stopped
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()
		if err := httpServer.Shutdown(ctx); err != nil {
			log.Printf("优雅退出失败: %v", err)
		}
	}()
	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
