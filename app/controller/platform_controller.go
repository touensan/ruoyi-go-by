package controller

import (
	"errors"
	"net/http"
	"strings"

	"ruoyi-go/app/dto"
	"ruoyi-go/app/security"
	"ruoyi-go/app/service"
	"ruoyi-go/framework/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PlatformController struct{}

func writePlatformError(ctx *gin.Context, err error) {
	status, message := http.StatusInternalServerError, "平台积分服务暂时不可用"
	switch {
	case errors.Is(err, service.ErrPlatformDisabled):
		status, message = http.StatusForbidden, "兑换中心暂未开放"
	case errors.Is(err, service.ErrPlatformInvalid):
		status, message = http.StatusBadRequest, "请求参数或兑换码无效"
	case errors.Is(err, service.ErrPlatformCodeUsed):
		status, message = http.StatusConflict, "兑换码已经使用"
	case errors.Is(err, service.ErrPlatformClaimed):
		status, message = http.StatusConflict, "该任务或核验账号已经领取过奖励"
	case errors.Is(err, service.ErrPlatformProvider):
		status, message = http.StatusServiceUnavailable, "雨云任务暂不可核验"
	case errors.Is(err, service.ErrPlatformPayment):
		status, message = http.StatusServiceUnavailable, "支付服务暂不可用"
	case errors.Is(err, service.ErrPlatformInsufficient):
		status, message = http.StatusPaymentRequired, "积分不足"
	case errors.Is(err, gorm.ErrRecordNotFound):
		status, message = http.StatusNotFound, "记录不存在"
	}
	response.NewError().SetStatus(status).SetMsg(message).Json(ctx)
}

func (*PlatformController) Overview(ctx *gin.Context) {
	data, err := service.NewPlatformService().Overview(ctx.Request.Context(), security.GetAuthUserId(ctx))
	if err != nil {
		writePlatformError(ctx, err)
		return
	}
	response.NewSuccess().SetData("data", data).Json(ctx)
}

func (*PlatformController) Ledger(ctx *gin.Context) {
	data, err := service.NewPlatformService().Ledger(ctx.Request.Context(), security.GetAuthUserId(ctx))
	if err != nil {
		writePlatformError(ctx, err)
		return
	}
	response.NewSuccess().SetData("data", data).Json(ctx)
}

func (*PlatformController) Recharges(ctx *gin.Context) {
	data, err := service.NewPlatformService().Recharges(ctx.Request.Context(), security.GetAuthUserId(ctx))
	if err != nil {
		writePlatformError(ctx, err)
		return
	}
	response.NewSuccess().SetData("data", data).Json(ctx)
}

func (*PlatformController) CreateRecharge(ctx *gin.Context) {
	var request dto.CreatePlatformRechargeRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		writePlatformError(ctx, service.ErrPlatformInvalid)
		return
	}
	idempotencyKey := strings.TrimSpace(ctx.GetHeader("Idempotency-Key"))
	data, err := service.NewPlatformService().CreateRecharge(
		ctx.Request.Context(), security.GetAuthUserId(ctx), request, idempotencyKey,
		ctx.ClientIP(), platformBaseURL(ctx),
	)
	if err != nil {
		writePlatformError(ctx, err)
		return
	}
	ctx.Header("Cache-Control", "no-store")
	response.NewSuccess().SetStatus(http.StatusCreated).SetData("data", data).Json(ctx)
}

func (*PlatformController) Center(ctx *gin.Context) {
	data, err := service.NewPlatformService().Center(ctx.Request.Context(), security.GetAuthUserId(ctx))
	if err != nil {
		writePlatformError(ctx, err)
		return
	}
	response.NewSuccess().SetData("data", data).Json(ctx)
}

func (*PlatformController) SubmitClaim(ctx *gin.Context) {
	var request dto.SubmitPlatformClaimRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		writePlatformError(ctx, service.ErrPlatformInvalid)
		return
	}
	claim, code, err := service.NewPlatformService().SubmitClaim(ctx.Request.Context(), security.GetAuthUserId(ctx), request)
	if err != nil {
		writePlatformError(ctx, err)
		return
	}
	ctx.Header("Cache-Control", "no-store")
	response.NewSuccess().SetStatus(http.StatusCreated).SetData("data", map[string]interface{}{"claim": claim, "code": code}).Json(ctx)
}

func (*PlatformController) Redeem(ctx *gin.Context) {
	var request dto.RedeemPlatformCodeRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		writePlatformError(ctx, service.ErrPlatformInvalid)
		return
	}
	data, err := service.NewPlatformService().Redeem(ctx.Request.Context(), security.GetAuthUserId(ctx), request.Code)
	if err != nil {
		writePlatformError(ctx, err)
		return
	}
	response.NewSuccess().SetData("data", data).Json(ctx)
}

type PlatformAdminController struct{}

func (*PlatformAdminController) Settings(ctx *gin.Context) {
	data, err := service.NewPlatformService().Settings(ctx.Request.Context())
	if err != nil {
		writePlatformError(ctx, err)
		return
	}
	ctx.Header("Cache-Control", "no-store")
	response.NewSuccess().SetData("data", data).Json(ctx)
}

func (*PlatformAdminController) SaveSettings(ctx *gin.Context) {
	var request dto.UpdatePlatformSettingsRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		writePlatformError(ctx, service.ErrPlatformInvalid)
		return
	}
	data, err := service.NewPlatformService().SaveSettings(ctx.Request.Context(), request, security.GetAuthUserName(ctx))
	if err != nil {
		writePlatformError(ctx, err)
		return
	}
	ctx.Header("Cache-Control", "no-store")
	response.NewSuccess().SetData("data", data).Json(ctx)
}

func (*PlatformAdminController) TestRainyun(ctx *gin.Context) {
	data, err := service.NewPlatformService().TestRainyun(ctx.Request.Context())
	if err != nil {
		writePlatformError(ctx, err)
		return
	}
	response.NewSuccess().SetData("data", data).Json(ctx)
}

func (*PlatformAdminController) Tasks(ctx *gin.Context) {
	data, err := service.NewPlatformService().Tasks(ctx.Request.Context())
	if err != nil {
		writePlatformError(ctx, err)
		return
	}
	response.NewSuccess().SetData("data", data).Json(ctx)
}

func (*PlatformAdminController) CreateTask(ctx *gin.Context) {
	var request dto.SavePlatformTaskRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		writePlatformError(ctx, service.ErrPlatformInvalid)
		return
	}
	if err := service.NewPlatformService().SaveTask(ctx.Request.Context(), "", request); err != nil {
		writePlatformError(ctx, err)
		return
	}
	response.NewSuccess().SetStatus(http.StatusCreated).Json(ctx)
}

func (*PlatformAdminController) UpdateTask(ctx *gin.Context) {
	var request dto.SavePlatformTaskRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		writePlatformError(ctx, service.ErrPlatformInvalid)
		return
	}
	if err := service.NewPlatformService().SaveTask(ctx.Request.Context(), ctx.Param("publicId"), request); err != nil {
		writePlatformError(ctx, err)
		return
	}
	response.NewSuccess().Json(ctx)
}

func (*PlatformAdminController) Claims(ctx *gin.Context) {
	var request dto.PlatformAdminListRequest
	if err := ctx.ShouldBindQuery(&request); err != nil {
		writePlatformError(ctx, service.ErrPlatformInvalid)
		return
	}
	rows, total, err := service.NewPlatformService().Claims(ctx.Request.Context(), request)
	if err != nil {
		writePlatformError(ctx, err)
		return
	}
	response.NewSuccess().SetPageData(rows, int(total)).Json(ctx)
}

func (*PlatformAdminController) ReviewClaim(ctx *gin.Context) {
	var request dto.ReviewPlatformClaimRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		writePlatformError(ctx, service.ErrPlatformInvalid)
		return
	}
	code, err := service.NewPlatformService().ReviewClaim(ctx.Request.Context(), ctx.Param("publicId"), request)
	if err != nil {
		writePlatformError(ctx, err)
		return
	}
	ctx.Header("Cache-Control", "no-store")
	response.NewSuccess().SetData("data", map[string]string{"code": code}).Json(ctx)
}

func (*PlatformAdminController) Codes(ctx *gin.Context) {
	var request dto.PlatformAdminListRequest
	if err := ctx.ShouldBindQuery(&request); err != nil {
		writePlatformError(ctx, service.ErrPlatformInvalid)
		return
	}
	rows, total, err := service.NewPlatformService().Codes(ctx.Request.Context(), request)
	if err != nil {
		writePlatformError(ctx, err)
		return
	}
	response.NewSuccess().SetPageData(rows, int(total)).Json(ctx)
}

func (*PlatformAdminController) GenerateCodes(ctx *gin.Context) {
	var request dto.GeneratePlatformCodesRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		writePlatformError(ctx, service.ErrPlatformInvalid)
		return
	}
	data, err := service.NewPlatformService().GenerateCodes(ctx.Request.Context(), request)
	if err != nil {
		writePlatformError(ctx, err)
		return
	}
	ctx.Header("Cache-Control", "no-store")
	response.NewSuccess().SetStatus(http.StatusCreated).SetData("data", data).Json(ctx)
}

func (*PlatformAdminController) Accounts(ctx *gin.Context) {
	var request dto.PlatformAdminListRequest
	if err := ctx.ShouldBindQuery(&request); err != nil {
		writePlatformError(ctx, service.ErrPlatformInvalid)
		return
	}
	rows, total, err := service.NewPlatformService().Accounts(ctx.Request.Context(), request)
	if err != nil {
		writePlatformError(ctx, err)
		return
	}
	response.NewSuccess().SetPageData(rows, int(total)).Json(ctx)
}

func (*PlatformAdminController) Adjust(ctx *gin.Context) {
	var request dto.AdjustPlatformPointsRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		writePlatformError(ctx, service.ErrPlatformInvalid)
		return
	}
	if err := service.NewPlatformService().Adjust(ctx.Request.Context(), request); err != nil {
		writePlatformError(ctx, err)
		return
	}
	response.NewSuccess().Json(ctx)
}

func platformBaseURL(ctx *gin.Context) string {
	scheme := strings.TrimSpace(ctx.GetHeader("X-Forwarded-Proto"))
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
