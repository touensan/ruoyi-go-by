package controller

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"ruoyi-go/common/upload"
	"ruoyi-go/config"
	"ruoyi-go/framework/response"
	"strings"

	"github.com/gin-gonic/gin"
)

type CommonController struct{}

func (*CommonController) Upload(ctx *gin.Context) {
	fileHeader, err := ctx.FormFile("file")
	if err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	result, err := upload.New(upload.SetRandomName(true)).SetFile(&upload.File{
		FileName:    fileHeader.Filename,
		FileSize:    int(fileHeader.Size),
		FileType:    fileHeader.Header.Get("Content-Type"),
		FileHeader:  fileHeader.Header,
		FileContent: content,
	}).Save()
	if err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	response.NewSuccess().
		SetData("fileName", result.FileName).
		SetData("newFileName", result.FileName).
		SetData("originalFilename", result.OriginalName).
		SetData("url", result.Url).
		Json(ctx)
}

func (*CommonController) Download(ctx *gin.Context) {
	fileName := filepath.Clean(ctx.Query("fileName"))
	if fileName == "." || strings.Contains(fileName, "..") {
		response.NewError().SetMsg("文件名无效").Json(ctx)
		return
	}
	serveDownload(ctx, filepath.Join(config.Data.Ruoyi.UploadPath, fileName), filepath.Base(fileName))
}

func (*CommonController) DownloadResource(ctx *gin.Context) {
	resource := filepath.Clean(ctx.Query("resource"))
	if resource == "." || strings.Contains(resource, "..") {
		response.NewError().SetMsg("资源路径无效").Json(ctx)
		return
	}
	serveDownload(ctx, resource, filepath.Base(resource))
}

func (*CommonController) Druid(ctx *gin.Context) {
	ctx.Data(http.StatusOK, "text/html; charset=utf-8", []byte(`<!doctype html><html><head><meta charset="utf-8"><title>RuoYi-Go BY 数据源监控</title></head><body style="font-family:sans-serif;padding:24px;"><h3>RuoYi-Go BY 数据源监控</h3><p>当前后端使用 GORM 连接池，未集成 Java Druid 控制台。</p></body></html>`))
}

func (*CommonController) Swagger(ctx *gin.Context) {
	ctx.Data(http.StatusOK, "text/html; charset=utf-8", []byte(`<!doctype html><html><head><meta charset="utf-8"><title>RuoYi-Go BY API 文档</title></head><body style="font-family:sans-serif;padding:24px;"><h3>RuoYi-Go BY API 文档</h3><p>当前 ruoyi-go 基线未集成 Swagger UI，接口以后端路由为准。</p></body></html>`))
}

func serveDownload(ctx *gin.Context, path string, downloadName string) {
	if _, err := os.Stat(path); err != nil {
		response.NewError().SetMsg("文件不存在").Json(ctx)
		return
	}
	ctx.Header("download-filename", downloadName)
	ctx.FileAttachment(path, downloadName)
}
