package systemcontroller

import (
	"ruoyi-go/framework/response"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type NoticeController struct{}

type noticeItem struct {
	NoticeId      int    `json:"noticeId"`
	NoticeTitle   string `json:"noticeTitle"`
	NoticeType    string `json:"noticeType"`
	NoticeContent string `json:"noticeContent"`
	Status        string `json:"status"`
	CreateTime    string `json:"createTime"`
	IsRead        bool   `json:"isRead"`
}

type noticeReadUser struct {
	UserId     int    `json:"userId"`
	UserName   string `json:"userName"`
	NickName   string `json:"nickName"`
	ReadTime   string `json:"readTime"`
	NoticeId   int    `json:"noticeId"`
	NoticeType string `json:"noticeType"`
}

// 当前 ruoyi-go 基线没有公告表；这些接口用于兼容 RuoYi-Vue3-ts 顶部通知组件。
func (*NoticeController) ListTop(ctx *gin.Context) {
	response.NewSuccess().
		SetData("data", []noticeItem{}).
		SetData("unreadCount", 0).
		Json(ctx)
}

func (*NoticeController) List(ctx *gin.Context) {
	response.NewSuccess().SetPageData([]noticeItem{}, 0).Json(ctx)
}

func (*NoticeController) Detail(ctx *gin.Context) {
	noticeId, _ := strconv.Atoi(ctx.Param("noticeId"))
	response.NewSuccess().SetData("data", noticeItem{
		NoticeId:      noticeId,
		NoticeTitle:   "暂无公告",
		NoticeType:    "2",
		NoticeContent: "",
		Status:        "0",
		IsRead:        true,
	}).Json(ctx)
}

func (*NoticeController) Create(ctx *gin.Context) {
	response.NewSuccess().Json(ctx)
}

func (*NoticeController) Update(ctx *gin.Context) {
	response.NewSuccess().Json(ctx)
}

func (*NoticeController) Remove(ctx *gin.Context) {
	response.NewSuccess().Json(ctx)
}

func (*NoticeController) MarkRead(ctx *gin.Context) {
	response.NewSuccess().Json(ctx)
}

func (*NoticeController) MarkReadAll(ctx *gin.Context) {
	if strings.TrimSpace(ctx.Query("ids")) == "" {
		response.NewSuccess().Json(ctx)
		return
	}
	response.NewSuccess().Json(ctx)
}

func (*NoticeController) ReadUsers(ctx *gin.Context) {
	response.NewSuccess().SetPageData([]noticeReadUser{}, 0).Json(ctx)
}
