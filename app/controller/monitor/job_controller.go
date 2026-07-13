package monitorcontroller

import (
	"ruoyi-go/common/excel"
	"ruoyi-go/framework/response"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type JobController struct{}
type JobLogController struct{}

type jobItem struct {
	JobId          int    `json:"jobId" excel:"name:任务编号;"`
	JobName        string `json:"jobName" excel:"name:任务名称;"`
	JobGroup       string `json:"jobGroup" excel:"name:任务组名;"`
	InvokeTarget   string `json:"invokeTarget" excel:"name:调用目标字符串;"`
	CronExpression string `json:"cronExpression" excel:"name:执行表达式;"`
	MisfirePolicy  string `json:"misfirePolicy" excel:"name:计划策略;"`
	Concurrent     string `json:"concurrent" excel:"name:并发执行;"`
	Status         string `json:"status" excel:"name:状态;replace:0_正常,1_暂停;"`
	CreateTime     string `json:"createTime" excel:"name:创建时间;"`
}

type jobLogItem struct {
	JobLogId      int    `json:"jobLogId" excel:"name:任务日志编号;"`
	JobName       string `json:"jobName" excel:"name:任务名称;"`
	JobGroup      string `json:"jobGroup" excel:"name:任务组名;"`
	InvokeTarget  string `json:"invokeTarget" excel:"name:调用目标字符串;"`
	JobMessage    string `json:"jobMessage" excel:"name:日志信息;"`
	ExceptionInfo string `json:"exceptionInfo" excel:"name:异常信息;"`
	Status        string `json:"status" excel:"name:状态;replace:0_正常,1_失败;"`
	CreateTime    string `json:"createTime" excel:"name:创建时间;"`
}

func (*JobController) List(ctx *gin.Context) {
	response.NewSuccess().SetPageData([]jobItem{}, 0).Json(ctx)
}

func (*JobController) Detail(ctx *gin.Context) {
	jobId, _ := strconv.Atoi(ctx.Param("jobId"))
	response.NewSuccess().SetData("data", jobItem{
		JobId:          jobId,
		JobName:        "暂无定时任务",
		JobGroup:       "DEFAULT",
		InvokeTarget:   "",
		CronExpression: "",
		MisfirePolicy:  "1",
		Concurrent:     "1",
		Status:         "1",
	}).Json(ctx)
}

func (*JobController) Create(ctx *gin.Context) {
	response.NewSuccess().Json(ctx)
}

func (*JobController) Update(ctx *gin.Context) {
	response.NewSuccess().Json(ctx)
}

func (*JobController) Remove(ctx *gin.Context) {
	response.NewSuccess().Json(ctx)
}

func (*JobController) ChangeStatus(ctx *gin.Context) {
	response.NewSuccess().Json(ctx)
}

func (*JobController) Run(ctx *gin.Context) {
	response.NewSuccess().SetMsg("当前 Go 基线未启用定时任务调度器，已接受执行请求").Json(ctx)
}

func (*JobController) Export(ctx *gin.Context) {
	file, err := excel.NormalDynamicExport("Sheet1", "", "", false, false, []jobItem{}, nil)
	if err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}
	excel.DownLoadExcel("job_"+time.Now().Format("20060102150405"), ctx.Writer, file)
}

func (*JobLogController) List(ctx *gin.Context) {
	response.NewSuccess().SetPageData([]jobLogItem{}, 0).Json(ctx)
}

func (*JobLogController) Detail(ctx *gin.Context) {
	jobLogId, _ := strconv.Atoi(ctx.Param("jobLogId"))
	response.NewSuccess().SetData("data", jobLogItem{JobLogId: jobLogId, Status: "0"}).Json(ctx)
}

func (*JobLogController) Remove(ctx *gin.Context) {
	response.NewSuccess().Json(ctx)
}

func (*JobLogController) Clean(ctx *gin.Context) {
	response.NewSuccess().Json(ctx)
}

func (*JobLogController) Export(ctx *gin.Context) {
	file, err := excel.NormalDynamicExport("Sheet1", "", "", false, false, []jobLogItem{}, nil)
	if err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}
	excel.DownLoadExcel("jobLog_"+time.Now().Format("20060102150405"), ctx.Writer, file)
}
