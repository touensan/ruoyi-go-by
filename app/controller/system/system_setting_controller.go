package systemcontroller

import (
	"net/http"
	"ruoyi-go/app/dto"
	"ruoyi-go/app/security"
	"ruoyi-go/app/service"
	"ruoyi-go/framework/response"
	"strings"

	"github.com/gin-gonic/gin"
)

type SystemSettingController struct{}

func (*SystemSettingController) Detail(ctx *gin.Context) {
	site, payment, mail := (&service.SystemSettingService{}).GetSettings()
	response.NewSuccess().SetData("data", dto.SystemSettingResponse{
		Site:      site,
		Payment:   payment,
		Mail:      mail,
		Generated: (&service.SystemSettingService{}).GeneratedCallbackUrls(requestBaseUrl(ctx)),
	}).Json(ctx)
}

func (*SystemSettingController) PublicSite(ctx *gin.Context) {
	site, _, _ := (&service.SystemSettingService{}).GetSettings()
	response.NewSuccess().SetData("data", site).Json(ctx)
}

func (*SystemSettingController) UpdateSite(ctx *gin.Context) {
	var param dto.SiteSetting
	if err := ctx.ShouldBindJSON(&param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}
	if err := (&service.SystemSettingService{}).SaveSiteSetting(param, security.GetAuthUserName(ctx)); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}
	response.NewSuccess().Json(ctx)
}

func (*SystemSettingController) UpdatePayment(ctx *gin.Context) {
	var param dto.PaymentSetting
	if err := ctx.ShouldBindJSON(&param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}
	if err := (&service.SystemSettingService{}).SavePaymentSetting(param, security.GetAuthUserName(ctx)); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}
	response.NewSuccess().Json(ctx)
}

func (*SystemSettingController) UpdateMail(ctx *gin.Context) {
	var param dto.MailSetting
	if err := ctx.ShouldBindJSON(&param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}
	if err := (&service.SystemSettingService{}).SaveMailSetting(param, security.GetAuthUserName(ctx)); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}
	response.NewSuccess().Json(ctx)
}

func (*SystemSettingController) TestPayment(ctx *gin.Context) {
	var param dto.PaymentTestRequest
	if err := ctx.ShouldBindJSON(&param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}
	_, payment, _ := (&service.SystemSettingService{}).GetSettings()
	generated := (&service.SystemSettingService{}).GeneratedCallbackUrls(requestBaseUrl(ctx))
	result := (&service.SystemSettingService{}).TestPayment(payment, param, ctx.ClientIP(), generated.NotifyUrl, generated.ReturnUrl)
	response.NewSuccess().SetData("data", result).Json(ctx)
}

func (*SystemSettingController) TestMail(ctx *gin.Context) {
	var param dto.MailTestRequest
	if err := ctx.ShouldBindJSON(&param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}
	_, _, mail := (&service.SystemSettingService{}).GetSettings()
	result := (&service.SystemSettingService{}).TestMail(mail, param)
	response.NewSuccess().SetData("data", result).Json(ctx)
}

func (*SystemSettingController) PaymentNotify(ctx *gin.Context) {
	if err := ctx.Request.ParseForm(); err != nil {
		ctx.String(http.StatusOK, "fail")
		return
	}
	if err := (&service.SystemSettingService{}).VerifyPaymentCallback(ctx.Request.Form); err != nil {
		ctx.String(http.StatusOK, "fail")
		return
	}
	if strings.ToUpper(ctx.Request.Form.Get("trade_status")) != "TRADE_SUCCESS" {
		ctx.String(http.StatusOK, "fail")
		return
	}
	ctx.String(http.StatusOK, "success")
}

func (*SystemSettingController) PaymentReturn(ctx *gin.Context) {
	if err := ctx.Request.ParseForm(); err != nil {
		ctx.String(http.StatusOK, "支付返回参数无效")
		return
	}
	if err := (&service.SystemSettingService{}).VerifyPaymentCallback(ctx.Request.Form); err != nil {
		ctx.String(http.StatusOK, "支付返回验签失败")
		return
	}
	ctx.String(http.StatusOK, "支付返回验签成功")
}

func requestBaseUrl(ctx *gin.Context) string {
	scheme := strings.TrimSpace(ctx.GetHeader("X-Forwarded-Proto"))
	if scheme == "" {
		scheme = strings.TrimSpace(ctx.GetHeader("X-Forwarded-Scheme"))
	}
	if scheme == "" {
		if ctx.Request.TLS != nil {
			scheme = "https"
		} else {
			scheme = "http"
		}
	}
	host := strings.TrimSpace(ctx.GetHeader("X-Forwarded-Host"))
	if host == "" {
		host = ctx.Request.Host
	}
	return scheme + "://" + host
}
