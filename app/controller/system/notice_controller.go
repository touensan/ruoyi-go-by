package systemcontroller

import (
	"ruoyi-go/framework/response"

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
		SetData("supported", false).
		SetMsg("当前 Go 主线尚未实现公告持久化，顶部通知为空").
		Json(ctx)
}

func (*NoticeController) List(ctx *gin.Context) {
	response.Unsupported(ctx, "公告持久化与已读记录")
}

func (*NoticeController) Detail(ctx *gin.Context) {
	response.Unsupported(ctx, "公告持久化与已读记录")
}

func (*NoticeController) Create(ctx *gin.Context) {
	response.Unsupported(ctx, "公告持久化与已读记录")
}

func (*NoticeController) Update(ctx *gin.Context) {
	response.Unsupported(ctx, "公告持久化与已读记录")
}

func (*NoticeController) Remove(ctx *gin.Context) {
	response.Unsupported(ctx, "公告持久化与已读记录")
}

func (*NoticeController) MarkRead(ctx *gin.Context) {
	response.Unsupported(ctx, "公告持久化与已读记录")
}

func (*NoticeController) MarkReadAll(ctx *gin.Context) {
	response.Unsupported(ctx, "公告持久化与已读记录")
}

func (*NoticeController) ReadUsers(ctx *gin.Context) {
	response.Unsupported(ctx, "公告持久化与已读记录")
}
