package controller_test

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http/httptest"
	monitor "ruoyi-go/app/controller/monitor"
	system "ruoyi-go/app/controller/system"
	tool "ruoyi-go/app/controller/tool"
	"strings"
	"testing"
)

func TestUnsupportedCapabilities(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cases := []struct {
		name    string
		handler gin.HandlerFunc
	}{
		{"JobController/List", (&monitor.JobController{}).List},
		{"JobController/Detail", (&monitor.JobController{}).Detail},
		{"JobController/Create", (&monitor.JobController{}).Create},
		{"JobController/Update", (&monitor.JobController{}).Update},
		{"JobController/Remove", (&monitor.JobController{}).Remove},
		{"JobController/ChangeStatus", (&monitor.JobController{}).ChangeStatus},
		{"JobController/Run", (&monitor.JobController{}).Run},
		{"JobController/Export", (&monitor.JobController{}).Export},
		{"JobLogController/List", (&monitor.JobLogController{}).List},
		{"JobLogController/Detail", (&monitor.JobLogController{}).Detail},
		{"JobLogController/Remove", (&monitor.JobLogController{}).Remove},
		{"JobLogController/Clean", (&monitor.JobLogController{}).Clean},
		{"JobLogController/Export", (&monitor.JobLogController{}).Export},
		{"GenController/Update", (&tool.GenController{}).Update},
		{"GenController/ImportTable", (&tool.GenController{}).ImportTable},
		{"GenController/CreateTable", (&tool.GenController{}).CreateTable},
		{"GenController/Preview", (&tool.GenController{}).Preview},
		{"GenController/Remove", (&tool.GenController{}).Remove},
		{"GenController/GenCode", (&tool.GenController{}).GenCode},
		{"GenController/SynchDb", (&tool.GenController{}).SynchDb},
		{"GenController/BatchGenCode", (&tool.GenController{}).BatchGenCode},
		{"NoticeController/List", (&system.NoticeController{}).List},
		{"NoticeController/Detail", (&system.NoticeController{}).Detail},
		{"NoticeController/Create", (&system.NoticeController{}).Create},
		{"NoticeController/Update", (&system.NoticeController{}).Update},
		{"NoticeController/Remove", (&system.NoticeController{}).Remove},
		{"NoticeController/MarkRead", (&system.NoticeController{}).MarkRead},
		{"NoticeController/MarkReadAll", (&system.NoticeController{}).MarkReadAll},
		{"NoticeController/ReadUsers", (&system.NoticeController{}).ReadUsers},
	}
	for _, item := range cases {
		t.Run(item.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			ctx.Request = httptest.NewRequest("POST", "/", nil)
			item.handler(ctx)
			if w.Code != 501 {
				t.Fatalf("status=%d body=%s", w.Code, w.Body.String())
			}
			var body struct {
				Code      int    `json:"code"`
				Msg       string `json:"msg"`
				Supported *bool  `json:"supported"`
			}
			if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
				t.Fatal(err)
			}
			if body.Code != 501 || body.Supported == nil || *body.Supported || !strings.Contains(body.Msg, "未执行") {
				t.Fatalf("misleading response: %s", w.Body.String())
			}
			if strings.Contains(w.Header().Get("Content-Type"), "zip") {
				t.Fatal("placeholder archive must not be offered as generated code")
			}
		})
	}
}

func TestNoticeHeaderDeclaresUnsupported(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	(&system.NoticeController{}).ListTop(ctx)
	var body struct {
		Supported   *bool `json:"supported"`
		Data        []any `json:"data"`
		UnreadCount int   `json:"unreadCount"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if w.Code != 200 || body.Supported == nil || *body.Supported || len(body.Data) != 0 || body.UnreadCount != 0 {
		t.Fatalf("invalid empty header compatibility response: %s", w.Body.String())
	}
}
