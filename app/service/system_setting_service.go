package service

import (
	"bytes"
	"crypto"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"mime"
	"net"
	"net/http"
	"net/smtp"
	"net/url"
	"ruoyi-go/app/dto"
	"ruoyi-go/app/model"
	"ruoyi-go/framework/dal"
	"sort"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm/clause"
)

const (
	systemSettingKeySite    = "site"
	systemSettingKeyPayment = "payment"
	systemSettingKeyMail    = "mail"

	paymentTestAmount         = "1.00"
	paymentTestOrderName      = "RuoYi-Go BY支付功能测试"
	paymentTestDevice         = "pc"
	paymentTestV2Method       = "web"
	paymentRequestTimeoutSec  = 15
	mailRequestTimeoutSeconds = 15
)

type SystemSettingService struct{}

func EnsureSystemSettingBaseline() error {
	if err := dal.Gorm.AutoMigrate(&model.SysSystemSetting{}); err != nil {
		return err
	}

	service := &SystemSettingService{}
	if err := service.ensureDefaultSetting(systemSettingKeySite, "site", service.defaultSiteSetting(), "站点配置", "system"); err != nil {
		return err
	}
	if err := service.ensureDefaultSetting(systemSettingKeyPayment, "payment", service.defaultPaymentSetting(), "支付配置", "system"); err != nil {
		return err
	}
	if err := service.ensureDefaultSetting(systemSettingKeyMail, "mail", service.defaultMailSetting(), "邮箱配置", "system"); err != nil {
		return err
	}
	site, payment, mail := service.GetSettings()
	if err := service.SaveSiteSetting(site, "system"); err != nil {
		return err
	}
	if err := service.SavePaymentSetting(payment, "system"); err != nil {
		return err
	}
	if err := service.SaveMailSetting(mail, "system"); err != nil {
		return err
	}

	return ensureSystemSettingMenu()
}

func (s *SystemSettingService) GetSettings() (dto.SiteSetting, dto.PaymentSetting, dto.MailSetting) {
	site := s.defaultSiteSetting()
	payment := s.defaultPaymentSetting()
	mail := s.defaultMailSetting()

	s.loadSetting(systemSettingKeySite, &site)
	s.loadSetting(systemSettingKeyPayment, &payment)
	s.loadSetting(systemSettingKeyMail, &mail)

	return normalizeSiteSetting(site), normalizePaymentSetting(payment), normalizeMailSetting(mail)
}

func (s *SystemSettingService) SaveSiteSetting(setting dto.SiteSetting, userName string) error {
	setting = normalizeSiteSetting(setting)
	return s.saveSetting(systemSettingKeySite, "site", setting, "站点配置", userName)
}

func (s *SystemSettingService) SavePaymentSetting(setting dto.PaymentSetting, userName string) error {
	setting = normalizePaymentSetting(setting)
	return s.saveSetting(systemSettingKeyPayment, "payment", setting, "支付配置", userName)
}

func (s *SystemSettingService) SaveMailSetting(setting dto.MailSetting, userName string) error {
	setting = normalizeMailSetting(setting)
	return s.saveSetting(systemSettingKeyMail, "mail", setting, "邮箱配置", userName)
}

func (s *SystemSettingService) GeneratedCallbackUrls(baseUrl string) dto.SystemSettingGeneratedUrls {
	baseUrl = strings.TrimRight(baseUrl, "/")
	return dto.SystemSettingGeneratedUrls{
		NotifyUrl: baseUrl + "/api/system-config/payment/notify",
		ReturnUrl: baseUrl + "/api/system-config/payment/return",
	}
}

func (s *SystemSettingService) TestPayment(setting dto.PaymentSetting, req dto.PaymentTestRequest, clientIP, generatedNotifyUrl, generatedReturnUrl string) dto.PaymentTestResponse {
	setting = normalizePaymentSetting(setting)

	outTradeNo := "R1GO" + time.Now().Format("20060102150405")
	payType := firstEnabledPayType(setting.EnabledPayTypes)
	notifyUrl := firstNonEmpty(req.NotifyUrl, setting.NotifyUrl, generatedNotifyUrl)
	returnUrl := firstNonEmpty(req.ReturnUrl, setting.ReturnUrl, generatedReturnUrl)

	result := dto.PaymentTestResponse{
		Action:     "function",
		Version:    setting.EpayVersion,
		OutTradeNo: outTradeNo,
		PayType:    payType,
	}

	if payType == "" {
		result.Steps = append(result.Steps, errorStep("配置检查", "请至少启用支付宝支付或微信支付", 0))
		return result
	}
	if strings.TrimSpace(setting.GatewayUrl) == "" {
		result.Steps = append(result.Steps, errorStep("配置检查", "请先填写支付网关地址", 0))
		return result
	}
	if strings.TrimSpace(setting.MerchantId) == "" {
		result.Steps = append(result.Steps, errorStep("配置检查", "请先填写商户ID", 0))
		return result
	}
	if setting.EpayVersion == "v1" && strings.TrimSpace(setting.MerchantKey) == "" {
		result.Steps = append(result.Steps, errorStep("配置检查", "V1测试需要填写商户密钥", 0))
		return result
	}
	if setting.EpayVersion == "v2" && strings.TrimSpace(setting.MerchantPrivateKey) == "" {
		result.Steps = append(result.Steps, errorStep("配置检查", "V2测试需要填写商户私钥", 0))
		return result
	}

	if setting.EpayVersion == "v2" {
		result = s.testEpayV2(setting, result, epayTestContext{
			OutTradeNo: outTradeNo,
			Amount:     paymentTestAmount,
			PayType:    payType,
			Device:     paymentTestDevice,
			Method:     paymentTestV2Method,
			NotifyUrl:  notifyUrl,
			ReturnUrl:  returnUrl,
			OrderName:  paymentTestOrderName,
			ClientIP:   clientIP,
		})
	} else {
		result = s.testEpayV1(setting, result, epayTestContext{
			OutTradeNo: outTradeNo,
			Amount:     paymentTestAmount,
			PayType:    payType,
			Device:     paymentTestDevice,
			NotifyUrl:  notifyUrl,
			ReturnUrl:  returnUrl,
			OrderName:  paymentTestOrderName,
			ClientIP:   clientIP,
		})
	}

	result.Success = len(result.Steps) > 0
	for _, step := range result.Steps {
		if step.Status == "error" {
			result.Success = false
			break
		}
	}

	return result
}

func (s *SystemSettingService) TestMail(setting dto.MailSetting, req dto.MailTestRequest) dto.MailTestResponse {
	setting = normalizeMailSetting(setting)
	start := time.Now()
	to := strings.TrimSpace(firstNonEmpty(req.To, setting.TestRecipient))
	subject := firstNonEmpty(req.Subject, "RuoYi-Go BY邮箱测试")
	body := firstNonEmpty(req.Body, "这是一封来自 RuoYi-Go BY 系统配置的 SMTP 测试邮件。")

	if to == "" {
		return dto.MailTestResponse{Success: false, Message: "请填写测试收件人", Server: mailServerAddr(setting), ElapsedMs: time.Since(start).Milliseconds()}
	}
	if strings.TrimSpace(setting.Host) == "" || setting.Port <= 0 {
		return dto.MailTestResponse{Success: false, Message: "请先填写SMTP服务器和端口", Server: mailServerAddr(setting), ElapsedMs: time.Since(start).Milliseconds()}
	}
	if strings.TrimSpace(setting.Username) == "" || strings.TrimSpace(setting.Password) == "" {
		return dto.MailTestResponse{Success: false, Message: "请先填写SMTP账号和授权码/密码", Server: mailServerAddr(setting), ElapsedMs: time.Since(start).Milliseconds()}
	}

	if err := sendSMTPMail(setting, to, subject, body); err != nil {
		return dto.MailTestResponse{Success: false, Message: err.Error(), Server: mailServerAddr(setting), ElapsedMs: time.Since(start).Milliseconds()}
	}

	return dto.MailTestResponse{Success: true, Message: "邮件已发送", Server: mailServerAddr(setting), ElapsedMs: time.Since(start).Milliseconds()}
}

func (s *SystemSettingService) VerifyPaymentCallback(values url.Values) error {
	_, payment, _ := s.GetSettings()
	params := valuesToStringMap(values)
	if payment.EpayVersion == "v2" {
		if strings.TrimSpace(payment.PlatformPublicKey) == "" {
			return errors.New("未配置平台公钥，无法验签")
		}
		return verifyRSASign(params, payment.PlatformPublicKey)
	}
	expected := signMD5(params, payment.MerchantKey)
	actual := strings.ToLower(strings.TrimSpace(params["sign"]))
	if actual == "" || expected != actual {
		return errors.New("MD5签名验证失败")
	}
	return nil
}

func (s *SystemSettingService) ensureDefaultSetting(key, group string, value interface{}, remark, userName string) error {
	var count int64
	if err := dal.Gorm.Model(&model.SysSystemSetting{}).Where("setting_key = ?", key).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}
	return s.saveSetting(key, group, value, remark, userName)
}

func (s *SystemSettingService) loadSetting(key string, dest interface{}) {
	var setting model.SysSystemSetting
	if err := dal.Gorm.Model(&model.SysSystemSetting{}).Where("setting_key = ?", key).Take(&setting).Error; err != nil {
		return
	}
	if strings.TrimSpace(setting.SettingValue) == "" {
		return
	}
	_ = json.Unmarshal([]byte(setting.SettingValue), dest)
}

func (s *SystemSettingService) saveSetting(key, group string, value interface{}, remark, userName string) error {
	body, err := json.Marshal(value)
	if err != nil {
		return err
	}

	var count int64
	if err := dal.Gorm.Model(&model.SysSystemSetting{}).Where("setting_key = ?", key).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return dal.Gorm.Model(&model.SysSystemSetting{}).Create(&model.SysSystemSetting{
			SettingKey:   key,
			SettingGroup: group,
			SettingValue: string(body),
			Remark:       remark,
			CreateBy:     userName,
		}).Error
	}

	return dal.Gorm.Model(&model.SysSystemSetting{}).
		Where("setting_key = ?", key).
		Updates(map[string]interface{}{
			"setting_group": group,
			"setting_value": string(body),
			"remark":        remark,
			"update_by":     userName,
			"update_time":   time.Now(),
		}).Error
}

func (s *SystemSettingService) defaultSiteSetting() dto.SiteSetting {
	return dto.SiteSetting{
		Title:                "RuoYi-Go BY",
		Logo:                 "",
		Favicon:              "/admin/favicon.ico",
		Description:          "RuoYi-Go BY 开源项目",
		Keywords:             "RuoYi-Go BY,Go,Gin,Vue3,RuoYi",
		FrontendHeadCode:     "",
		SiteUrl:              "",
		IcpNo:                "",
		PublicSecurityNo:     "",
		CustomerServiceEmail: "",
		Copyright:            "Copyright © 2026 RuoYi-Go BY. All Rights Reserved.",
		DefaultLanguage:      "zh-CN",
		EnableSeo:            true,
	}
}

func (s *SystemSettingService) defaultPaymentSetting() dto.PaymentSetting {
	return dto.PaymentSetting{
		Enabled:         false,
		Provider:        "epay",
		EpayVersion:     "v1",
		EnabledPayTypes: []string{"alipay"},
	}
}

func (s *SystemSettingService) defaultMailSetting() dto.MailSetting {
	return dto.MailSetting{
		Enabled:    false,
		Provider:   "custom",
		Port:       465,
		Encryption: "ssl",
	}
}

func normalizeSiteSetting(setting dto.SiteSetting) dto.SiteSetting {
	setting.Title = strings.TrimSpace(firstNonEmpty(setting.Title, "RuoYi-Go BY"))
	setting.DefaultLanguage = strings.TrimSpace(firstNonEmpty(setting.DefaultLanguage, "zh-CN"))
	setting.Favicon = strings.TrimSpace(setting.Favicon)
	setting.Logo = strings.TrimSpace(setting.Logo)
	setting.SiteUrl = strings.TrimSpace(setting.SiteUrl)
	setting.CustomerServiceEmail = strings.TrimSpace(setting.CustomerServiceEmail)
	return setting
}

func normalizePaymentSetting(setting dto.PaymentSetting) dto.PaymentSetting {
	setting.Provider = "epay"
	setting.EpayVersion = strings.ToLower(strings.TrimSpace(firstNonEmpty(setting.EpayVersion, "v1")))
	if setting.EpayVersion != "v2" {
		setting.EpayVersion = "v1"
	}
	setting.GatewayUrl = strings.TrimRight(strings.TrimSpace(setting.GatewayUrl), "/")
	setting.MerchantId = strings.TrimSpace(setting.MerchantId)
	setting.NotifyUrl = strings.TrimSpace(setting.NotifyUrl)
	setting.ReturnUrl = strings.TrimSpace(setting.ReturnUrl)
	setting.EnabledPayTypes = normalizeEnabledPayTypes(setting.EnabledPayTypes)
	return setting
}

func normalizeMailSetting(setting dto.MailSetting) dto.MailSetting {
	setting.Provider = strings.ToLower(strings.TrimSpace(firstNonEmpty(setting.Provider, "custom")))
	if setting.Provider == "qq" {
		if strings.TrimSpace(setting.Host) == "" {
			setting.Host = "smtp.qq.com"
		}
		if setting.Port == 0 {
			setting.Port = 465
		}
		if strings.TrimSpace(setting.Encryption) == "" {
			setting.Encryption = "ssl"
		}
	} else if setting.Provider == "gmail" {
		if strings.TrimSpace(setting.Host) == "" {
			setting.Host = "smtp.gmail.com"
		}
		if setting.Port == 0 {
			setting.Port = 587
		}
		if strings.TrimSpace(setting.Encryption) == "" {
			setting.Encryption = "starttls"
		}
	}
	setting.Host = strings.TrimSpace(setting.Host)
	setting.Username = strings.TrimSpace(setting.Username)
	setting.FromEmail = strings.TrimSpace(firstNonEmpty(setting.FromEmail, setting.Username))
	setting.FromName = strings.TrimSpace(firstNonEmpty(setting.FromName, "RuoYi-Go BY"))
	setting.Encryption = strings.ToLower(strings.TrimSpace(firstNonEmpty(setting.Encryption, "ssl")))
	if setting.Encryption != "ssl" && setting.Encryption != "starttls" && setting.Encryption != "none" {
		setting.Encryption = "ssl"
	}
	if setting.Port <= 0 {
		if setting.Encryption == "starttls" {
			setting.Port = 587
		} else {
			setting.Port = 465
		}
	}
	return setting
}

func normalizeEnabledPayTypes(values []string) []string {
	allowed := map[string]bool{"alipay": true, "wxpay": true}
	seen := make(map[string]bool, len(values))
	result := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.ToLower(strings.TrimSpace(value))
		if !allowed[value] || seen[value] {
			continue
		}
		seen[value] = true
		result = append(result, value)
	}
	if len(result) == 0 {
		return []string{"alipay"}
	}
	return result
}

func firstEnabledPayType(values []string) string {
	for _, preferred := range []string{"alipay", "wxpay"} {
		for _, value := range values {
			if value == preferred {
				return value
			}
		}
	}
	return ""
}

type epayTestContext struct {
	OutTradeNo string
	Amount     string
	PayType    string
	Device     string
	Method     string
	NotifyUrl  string
	ReturnUrl  string
	OrderName  string
	ClientIP   string
}

type epayHTTPResponse struct {
	StatusCode int
	Body       string
	Data       map[string]interface{}
}

func (s *SystemSettingService) testEpayV1(setting dto.PaymentSetting, result dto.PaymentTestResponse, ctx epayTestContext) dto.PaymentTestResponse {
	result.Steps = append(result.Steps, s.epayV1Merchant(setting))
	createStep, tradeNo, payInfo := s.epayV1Create(setting, ctx)
	result.Steps = append(result.Steps, createStep)
	if tradeNo != "" {
		result.TradeNo = tradeNo
	}
	if payInfo != "" {
		result.PayInfo = payInfo
	}
	result.Steps = append(result.Steps, s.epayV1Query(setting, ctx.OutTradeNo, result.TradeNo))
	return result
}

func (s *SystemSettingService) testEpayV2(setting dto.PaymentSetting, result dto.PaymentTestResponse, ctx epayTestContext) dto.PaymentTestResponse {
	result.Steps = append(result.Steps, s.epayV2Merchant(setting))
	createStep, tradeNo, payInfo := s.epayV2Create(setting, ctx)
	result.Steps = append(result.Steps, createStep)
	if tradeNo != "" {
		result.TradeNo = tradeNo
	}
	if payInfo != "" {
		result.PayInfo = payInfo
	}
	result.Steps = append(result.Steps, s.epayV2Query(setting, ctx.OutTradeNo, result.TradeNo))
	result.Steps = append(result.Steps, s.epayV2Close(setting, ctx.OutTradeNo, result.TradeNo))
	return result
}

func (s *SystemSettingService) epayV1Merchant(setting dto.PaymentSetting) dto.SystemConfigTestStep {
	params := map[string]string{"act": "query", "pid": setting.MerchantId, "key": setting.MerchantKey}
	step, resp, err := requestEpay(setting, "GET", epayEndpoint(setting.GatewayUrl, "/api.php"), params)
	step.Name = "查询商户"
	if err != nil {
		step.Status = "error"
		step.Message = err.Error()
		return step
	}
	step.Response = maskPayload(resp.Data)
	if !codeEquals(resp.Data, "1") {
		step.Status = "error"
		step.Message = responseMessage(resp.Data, "商户查询失败")
		return step
	}
	step.Status = "success"
	step.Message = "商户查询成功"
	return step
}

func (s *SystemSettingService) epayV1Create(setting dto.PaymentSetting, ctx epayTestContext) (dto.SystemConfigTestStep, string, string) {
	params := map[string]string{
		"pid":          setting.MerchantId,
		"type":         ctx.PayType,
		"out_trade_no": ctx.OutTradeNo,
		"notify_url":   ctx.NotifyUrl,
		"return_url":   ctx.ReturnUrl,
		"name":         ctx.OrderName,
		"money":        ctx.Amount,
		"clientip":     firstNonEmpty(ctx.ClientIP, "127.0.0.1"),
		"device":       ctx.Device,
		"param":        "system-config-test",
	}
	params["sign"] = signMD5(params, setting.MerchantKey)
	params["sign_type"] = "MD5"
	step, resp, err := requestEpay(setting, "POST", epayEndpoint(setting.GatewayUrl, "/mapi.php"), params)
	step.Name = "创建测试订单"
	if err != nil {
		step.Status = "error"
		step.Message = err.Error()
		return step, "", ""
	}
	step.Response = maskPayload(resp.Data)
	tradeNo := getMapString(resp.Data, "trade_no")
	payInfo := firstNonEmpty(getMapString(resp.Data, "payurl"), getMapString(resp.Data, "qrcode"), getMapString(resp.Data, "urlscheme"))
	if !codeEquals(resp.Data, "1") {
		step.Status = "error"
		step.Message = responseMessage(resp.Data, "创建订单失败")
		return step, tradeNo, payInfo
	}
	step.Status = "success"
	step.Message = "测试订单已创建"
	return step, tradeNo, payInfo
}

func (s *SystemSettingService) epayV1Query(setting dto.PaymentSetting, outTradeNo, tradeNo string) dto.SystemConfigTestStep {
	if outTradeNo == "" && tradeNo == "" {
		return errorStep("查询订单", "请提供商户订单号或平台订单号", 0)
	}
	params := map[string]string{"act": "order", "pid": setting.MerchantId, "key": setting.MerchantKey}
	if tradeNo != "" {
		params["trade_no"] = tradeNo
	} else {
		params["out_trade_no"] = outTradeNo
	}
	step, resp, err := requestEpay(setting, "GET", epayEndpoint(setting.GatewayUrl, "/api.php"), params)
	step.Name = "查询订单"
	if err != nil {
		step.Status = "error"
		step.Message = err.Error()
		return step
	}
	step.Response = maskPayload(resp.Data)
	if !codeEquals(resp.Data, "1") {
		step.Status = "error"
		step.Message = responseMessage(resp.Data, "订单查询失败")
		return step
	}
	step.Status = "success"
	step.Message = "订单查询成功"
	return step
}

func (s *SystemSettingService) epayV2Merchant(setting dto.PaymentSetting) dto.SystemConfigTestStep {
	params := map[string]string{
		"pid":       setting.MerchantId,
		"timestamp": strconv.FormatInt(time.Now().Unix(), 10),
	}
	return s.requestEpayV2(setting, "查询商户", epayEndpoint(setting.GatewayUrl, "/api/merchant/info"), params, "0")
}

func (s *SystemSettingService) epayV2Create(setting dto.PaymentSetting, ctx epayTestContext) (dto.SystemConfigTestStep, string, string) {
	params := map[string]string{
		"pid":          setting.MerchantId,
		"method":       ctx.Method,
		"device":       ctx.Device,
		"type":         ctx.PayType,
		"out_trade_no": ctx.OutTradeNo,
		"notify_url":   ctx.NotifyUrl,
		"return_url":   ctx.ReturnUrl,
		"name":         ctx.OrderName,
		"money":        ctx.Amount,
		"clientip":     firstNonEmpty(ctx.ClientIP, "127.0.0.1"),
		"param":        "system-config-test",
		"timestamp":    strconv.FormatInt(time.Now().Unix(), 10),
	}
	step := s.requestEpayV2(setting, "创建测试订单", epayEndpoint(setting.GatewayUrl, "/api/pay/create"), params, "0")
	data, _ := step.Response.(map[string]interface{})
	return step, getMapString(data, "trade_no"), getMapString(data, "pay_info")
}

func (s *SystemSettingService) epayV2Query(setting dto.PaymentSetting, outTradeNo, tradeNo string) dto.SystemConfigTestStep {
	if outTradeNo == "" && tradeNo == "" {
		return errorStep("查询订单", "请提供商户订单号或平台订单号", 0)
	}
	params := map[string]string{
		"pid":       setting.MerchantId,
		"timestamp": strconv.FormatInt(time.Now().Unix(), 10),
	}
	if tradeNo != "" {
		params["trade_no"] = tradeNo
	} else {
		params["out_trade_no"] = outTradeNo
	}
	return s.requestEpayV2(setting, "查询订单", epayEndpoint(setting.GatewayUrl, "/api/pay/query"), params, "0")
}

func (s *SystemSettingService) epayV2Close(setting dto.PaymentSetting, outTradeNo, tradeNo string) dto.SystemConfigTestStep {
	if outTradeNo == "" && tradeNo == "" {
		return errorStep("关闭订单", "请提供商户订单号或平台订单号", 0)
	}
	params := map[string]string{
		"pid":       setting.MerchantId,
		"timestamp": strconv.FormatInt(time.Now().Unix(), 10),
	}
	if tradeNo != "" {
		params["trade_no"] = tradeNo
	} else {
		params["out_trade_no"] = outTradeNo
	}
	return s.requestEpayV2(setting, "关闭订单", epayEndpoint(setting.GatewayUrl, "/api/pay/close"), params, "0")
}

func (s *SystemSettingService) requestEpayV2(setting dto.PaymentSetting, name, endpoint string, params map[string]string, successCode string) dto.SystemConfigTestStep {
	signed, err := signRSAParams(params, setting.MerchantPrivateKey)
	if err != nil {
		return errorStep(name, err.Error(), 0)
	}
	step, resp, err := requestEpay(setting, "POST", endpoint, signed)
	step.Name = name
	if err != nil {
		step.Status = "error"
		step.Message = err.Error()
		return step
	}
	step.Response = maskPayload(resp.Data)
	if codeEquals(resp.Data, successCode) && strings.TrimSpace(setting.PlatformPublicKey) != "" && getMapString(resp.Data, "sign") != "" {
		if err := verifyRSASign(mapInterfaceToString(resp.Data), setting.PlatformPublicKey); err != nil {
			step.Status = "error"
			step.Message = "平台返回验签失败：" + err.Error()
			return step
		}
	}
	if !codeEquals(resp.Data, successCode) {
		step.Status = "error"
		step.Message = responseMessage(resp.Data, name+"失败")
		return step
	}
	step.Status = "success"
	step.Message = name + "成功"
	return step
}

func requestEpay(setting dto.PaymentSetting, method, endpoint string, params map[string]string) (dto.SystemConfigTestStep, epayHTTPResponse, error) {
	start := time.Now()
	step := dto.SystemConfigTestStep{
		Method:  method,
		Url:     endpoint,
		Request: maskPayload(params),
	}
	client := &http.Client{Timeout: time.Duration(paymentRequestTimeoutSec) * time.Second}
	form := url.Values{}
	for key, value := range params {
		if value != "" {
			form.Set(key, value)
		}
	}

	var req *http.Request
	var err error
	if method == "GET" {
		requestUrl := endpoint
		if strings.Contains(requestUrl, "?") {
			requestUrl += "&" + form.Encode()
		} else {
			requestUrl += "?" + form.Encode()
		}
		step.Url = requestUrl
		req, err = http.NewRequest(http.MethodGet, requestUrl, nil)
	} else {
		req, err = http.NewRequest(http.MethodPost, endpoint, strings.NewReader(form.Encode()))
		if err == nil {
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		}
	}
	if err != nil {
		step.ElapsedMs = time.Since(start).Milliseconds()
		return step, epayHTTPResponse{}, err
	}
	req.Header.Set("Accept", "application/json")

	httpResp, err := client.Do(req)
	step.ElapsedMs = time.Since(start).Milliseconds()
	if err != nil {
		return step, epayHTTPResponse{}, err
	}
	defer httpResp.Body.Close()

	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return step, epayHTTPResponse{}, err
	}

	resp := epayHTTPResponse{StatusCode: httpResp.StatusCode, Body: string(body), Data: map[string]interface{}{}}
	decoder := json.NewDecoder(bytes.NewReader(body))
	decoder.UseNumber()
	if err := decoder.Decode(&resp.Data); err != nil {
		resp.Data = map[string]interface{}{"raw": strings.TrimSpace(string(body))}
		return step, resp, fmt.Errorf("上游返回非JSON内容，HTTP状态：%d", httpResp.StatusCode)
	}
	if httpResp.StatusCode < 200 || httpResp.StatusCode >= 300 {
		return step, resp, fmt.Errorf("上游HTTP状态异常：%d", httpResp.StatusCode)
	}
	return step, resp, nil
}

func signMD5(params map[string]string, key string) string {
	base := signatureBase(params)
	sum := md5.Sum([]byte(base + key))
	return hex.EncodeToString(sum[:])
}

func signRSAParams(params map[string]string, privateKeyPEM string) (map[string]string, error) {
	privateKey, err := parseRSAPrivateKey(privateKeyPEM)
	if err != nil {
		return nil, err
	}
	signed := make(map[string]string, len(params)+2)
	for key, value := range params {
		signed[key] = value
	}
	base := signatureBase(signed)
	digest := sha256.Sum256([]byte(base))
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, digest[:])
	if err != nil {
		return nil, err
	}
	signed["sign"] = base64.StdEncoding.EncodeToString(signature)
	signed["sign_type"] = "RSA"
	return signed, nil
}

func verifyRSASign(params map[string]string, publicKeyPEM string) error {
	signatureText := strings.TrimSpace(params["sign"])
	if signatureText == "" {
		return errors.New("缺少sign")
	}
	publicKey, err := parseRSAPublicKey(publicKeyPEM)
	if err != nil {
		return err
	}
	signature, err := base64.StdEncoding.DecodeString(signatureText)
	if err != nil {
		return err
	}
	base := signatureBase(params)
	digest := sha256.Sum256([]byte(base))
	return rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, digest[:], signature)
}

func signatureBase(params map[string]string) string {
	keys := make([]string, 0, len(params))
	for key, value := range params {
		if key == "sign" || key == "sign_type" || strings.TrimSpace(value) == "" {
			continue
		}
		keys = append(keys, key)
	}
	sort.Strings(keys)
	pairs := make([]string, 0, len(keys))
	for _, key := range keys {
		pairs = append(pairs, key+"="+params[key])
	}
	return strings.Join(pairs, "&")
}

func parseRSAPrivateKey(privateKeyPEM string) (*rsa.PrivateKey, error) {
	privateKeyPEM = normalizePEM(privateKeyPEM)
	block, _ := pem.Decode([]byte(privateKeyPEM))
	if block == nil {
		return nil, errors.New("商户私钥PEM格式无效")
	}
	if key, err := x509.ParsePKCS1PrivateKey(block.Bytes); err == nil {
		return key, nil
	}
	parsed, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	key, ok := parsed.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("商户私钥不是RSA私钥")
	}
	return key, nil
}

func parseRSAPublicKey(publicKeyPEM string) (*rsa.PublicKey, error) {
	publicKeyPEM = normalizePEM(publicKeyPEM)
	block, _ := pem.Decode([]byte(publicKeyPEM))
	if block == nil {
		return nil, errors.New("平台公钥PEM格式无效")
	}
	if key, err := x509.ParsePKIXPublicKey(block.Bytes); err == nil {
		if rsaKey, ok := key.(*rsa.PublicKey); ok {
			return rsaKey, nil
		}
	}
	key, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return key, nil
}

func normalizePEM(value string) string {
	value = strings.TrimSpace(value)
	value = strings.ReplaceAll(value, `\n`, "\n")
	return value
}

func sendSMTPMail(setting dto.MailSetting, to, subject, body string) error {
	addr := mailServerAddr(setting)
	timeout := time.Duration(mailRequestTimeoutSeconds) * time.Second

	var client *smtp.Client
	if setting.Encryption == "ssl" {
		conn, err := tls.DialWithDialer(&net.Dialer{Timeout: timeout}, "tcp", addr, &tls.Config{
			ServerName: setting.Host,
			MinVersion: tls.VersionTLS12,
		})
		if err != nil {
			return err
		}
		c, err := smtp.NewClient(conn, setting.Host)
		if err != nil {
			_ = conn.Close()
			return err
		}
		client = c
	} else {
		conn, err := net.DialTimeout("tcp", addr, timeout)
		if err != nil {
			return err
		}
		c, err := smtp.NewClient(conn, setting.Host)
		if err != nil {
			_ = conn.Close()
			return err
		}
		client = c
		if setting.Encryption == "starttls" {
			if err := client.StartTLS(&tls.Config{ServerName: setting.Host, MinVersion: tls.VersionTLS12}); err != nil {
				_ = client.Close()
				return err
			}
		}
	}
	defer client.Close()

	if setting.Username != "" || setting.Password != "" {
		if err := client.Auth(smtp.PlainAuth("", setting.Username, setting.Password, setting.Host)); err != nil {
			return err
		}
	}

	from := firstNonEmpty(setting.FromEmail, setting.Username)
	if err := client.Mail(from); err != nil {
		return err
	}
	if err := client.Rcpt(to); err != nil {
		return err
	}
	writer, err := client.Data()
	if err != nil {
		return err
	}
	if _, err := writer.Write([]byte(buildMailMessage(setting.FromName, from, to, subject, body))); err != nil {
		_ = writer.Close()
		return err
	}
	if err := writer.Close(); err != nil {
		return err
	}
	return client.Quit()
}

func buildMailMessage(fromName, from, to, subject, body string) string {
	encodedFromName := mime.QEncoding.Encode("UTF-8", firstNonEmpty(fromName, from))
	encodedSubject := mime.QEncoding.Encode("UTF-8", subject)
	headers := []string{
		"From: " + encodedFromName + " <" + from + ">",
		"To: " + to,
		"Subject: " + encodedSubject,
		"MIME-Version: 1.0",
		`Content-Type: text/plain; charset="UTF-8"`,
		"Content-Transfer-Encoding: 8bit",
	}
	return strings.Join(headers, "\r\n") + "\r\n\r\n" + body
}

func mailServerAddr(setting dto.MailSetting) string {
	return net.JoinHostPort(setting.Host, strconv.Itoa(setting.Port))
}

func ensureSystemSettingMenu() error {
	isFrame := 1
	isCache := 0
	menus := []model.SysMenu{
		{MenuId: 107, MenuName: "系统配置", ParentId: 1, OrderNum: 0, Path: "settings", Component: "system/settings/index", RouteName: "Settings", IsFrame: &isFrame, IsCache: &isCache, MenuType: "C", Visible: "0", Perms: "system:setting:list", Icon: "system", Status: "0", CreateBy: "system", Remark: "系统配置菜单"},
		{MenuId: 1035, MenuName: "系统配置查询", ParentId: 107, OrderNum: 1, Path: "#", Component: "", RouteName: "", IsFrame: &isFrame, IsCache: &isCache, MenuType: "F", Visible: "0", Perms: "system:setting:query", Icon: "#", Status: "0", CreateBy: "system"},
		{MenuId: 1036, MenuName: "系统配置修改", ParentId: 107, OrderNum: 2, Path: "#", Component: "", RouteName: "", IsFrame: &isFrame, IsCache: &isCache, MenuType: "F", Visible: "0", Perms: "system:setting:edit", Icon: "#", Status: "0", CreateBy: "system"},
		{MenuId: 1037, MenuName: "支付配置测试", ParentId: 107, OrderNum: 3, Path: "#", Component: "", RouteName: "", IsFrame: &isFrame, IsCache: &isCache, MenuType: "F", Visible: "0", Perms: "system:setting:pay:test", Icon: "#", Status: "0", CreateBy: "system"},
		{MenuId: 1038, MenuName: "邮箱配置测试", ParentId: 107, OrderNum: 4, Path: "#", Component: "", RouteName: "", IsFrame: &isFrame, IsCache: &isCache, MenuType: "F", Visible: "0", Perms: "system:setting:mail:test", Icon: "#", Status: "0", CreateBy: "system"},
	}

	for _, menu := range menus {
		var count int64
		if err := dal.Gorm.Model(&model.SysMenu{}).Where("menu_id = ?", menu.MenuId).Count(&count).Error; err != nil {
			return err
		}
		if count == 0 {
			if err := dal.Gorm.Model(&model.SysMenu{}).Create(&menu).Error; err != nil {
				return err
			}
			continue
		}
		if err := dal.Gorm.Model(&model.SysMenu{}).Where("menu_id = ?", menu.MenuId).Updates(map[string]interface{}{
			"menu_name":   menu.MenuName,
			"parent_id":   menu.ParentId,
			"order_num":   menu.OrderNum,
			"path":        menu.Path,
			"component":   menu.Component,
			"route_name":  menu.RouteName,
			"is_frame":    1,
			"is_cache":    0,
			"menu_type":   menu.MenuType,
			"visible":     menu.Visible,
			"perms":       menu.Perms,
			"icon":        menu.Icon,
			"status":      menu.Status,
			"remark":      menu.Remark,
			"update_by":   "system",
			"update_time": time.Now(),
		}).Error; err != nil {
			return err
		}
	}

	roleMenus := []model.SysRoleMenu{
		{RoleId: 2, MenuId: 107},
		{RoleId: 2, MenuId: 1035},
		{RoleId: 2, MenuId: 1036},
		{RoleId: 2, MenuId: 1037},
		{RoleId: 2, MenuId: 1038},
	}
	return dal.Gorm.Clauses(clause.OnConflict{DoNothing: true}).Create(&roleMenus).Error
}

func epayEndpoint(baseUrl, path string) string {
	return strings.TrimRight(baseUrl, "/") + path
}

func valuesToStringMap(values url.Values) map[string]string {
	params := make(map[string]string, len(values))
	for key := range values {
		params[key] = values.Get(key)
	}
	return params
}

func mapInterfaceToString(values map[string]interface{}) map[string]string {
	params := make(map[string]string, len(values))
	for key, value := range values {
		switch v := value.(type) {
		case json.Number:
			params[key] = v.String()
		case string:
			params[key] = v
		case nil:
			params[key] = ""
		default:
			params[key] = fmt.Sprint(v)
		}
	}
	return params
}

func maskPayload(value interface{}) interface{} {
	switch typed := value.(type) {
	case map[string]string:
		masked := make(map[string]string, len(typed))
		for key, value := range typed {
			if isSecretKey(key) {
				masked[key] = maskSecret(value)
			} else {
				masked[key] = value
			}
		}
		return masked
	case map[string]interface{}:
		masked := make(map[string]interface{}, len(typed))
		for key, value := range typed {
			if isSecretKey(key) {
				masked[key] = maskSecret(fmt.Sprint(value))
			} else {
				masked[key] = value
			}
		}
		return masked
	default:
		return value
	}
}

func isSecretKey(key string) bool {
	key = strings.ToLower(key)
	return key == "key" || key == "merchant_key" || key == "password" || strings.Contains(key, "private")
}

func maskSecret(value string) string {
	if value == "" {
		return ""
	}
	if len(value) <= 8 {
		return "******"
	}
	return value[:4] + "******" + value[len(value)-4:]
}

func codeEquals(data map[string]interface{}, expected string) bool {
	return getMapString(data, "code") == expected
}

func getMapString(data map[string]interface{}, key string) string {
	if data == nil {
		return ""
	}
	value, ok := data[key]
	if !ok || value == nil {
		return ""
	}
	switch typed := value.(type) {
	case json.Number:
		return typed.String()
	case string:
		return typed
	default:
		return fmt.Sprint(typed)
	}
}

func responseMessage(data map[string]interface{}, fallback string) string {
	if msg := getMapString(data, "msg"); msg != "" {
		return msg
	}
	return fallback
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func errorStep(name, message string, elapsedMs int64) dto.SystemConfigTestStep {
	return dto.SystemConfigTestStep{Name: name, Status: "error", Message: message, ElapsedMs: elapsedMs}
}
