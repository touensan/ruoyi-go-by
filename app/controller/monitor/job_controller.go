package monitorcontroller

import (
	"github.com/gin-gonic/gin"
	"ruoyi-go/framework/response"
)

type JobController struct{}
type JobLogController struct{}

// The native mainline has no task store or scheduler. Do not return synthetic results.
func (*JobController) List(ctx *gin.Context) {
	response.Unsupported(ctx, "定时任务与任务日志")
}

func (*JobController) Detail(ctx *gin.Context) {
	response.Unsupported(ctx, "定时任务与任务日志")
}

func (*JobController) Create(ctx *gin.Context) {
	response.Unsupported(ctx, "定时任务与任务日志")
}

func (*JobController) Update(ctx *gin.Context) {
	response.Unsupported(ctx, "定时任务与任务日志")
}

func (*JobController) Remove(ctx *gin.Context) {
	response.Unsupported(ctx, "定时任务与任务日志")
}

func (*JobController) ChangeStatus(ctx *gin.Context) {
	response.Unsupported(ctx, "定时任务与任务日志")
}

func (*JobController) Run(ctx *gin.Context) { response.Unsupported(ctx, "定时任务与任务日志") }

func (*JobController) Export(ctx *gin.Context) {
	response.Unsupported(ctx, "定时任务与任务日志")
}

func (*JobLogController) List(ctx *gin.Context) {
	response.Unsupported(ctx, "定时任务与任务日志")
}

func (*JobLogController) Detail(ctx *gin.Context) {
	response.Unsupported(ctx, "定时任务与任务日志")
}

func (*JobLogController) Remove(ctx *gin.Context) {
	response.Unsupported(ctx, "定时任务与任务日志")
}

func (*JobLogController) Clean(ctx *gin.Context) {
	response.Unsupported(ctx, "定时任务与任务日志")
}

func (*JobLogController) Export(ctx *gin.Context) {
	response.Unsupported(ctx, "定时任务与任务日志")
}
