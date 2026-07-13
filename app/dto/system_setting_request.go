package dto

type SiteSetting struct {
	Title                string `json:"title"`
	Logo                 string `json:"logo"`
	Favicon              string `json:"favicon"`
	Description          string `json:"description"`
	Keywords             string `json:"keywords"`
	FrontendHeadCode     string `json:"frontendHeadCode"`
	SiteUrl              string `json:"siteUrl"`
	IcpNo                string `json:"icpNo"`
	PublicSecurityNo     string `json:"publicSecurityNo"`
	CustomerServiceEmail string `json:"customerServiceEmail"`
	Copyright            string `json:"copyright"`
	DefaultLanguage      string `json:"defaultLanguage"`
	EnableSeo            bool   `json:"enableSeo"`
}

type PaymentSetting struct {
	Enabled            bool     `json:"enabled"`
	Provider           string   `json:"provider"`
	EpayVersion        string   `json:"epayVersion"`
	GatewayUrl         string   `json:"gatewayUrl"`
	MerchantId         string   `json:"merchantId"`
	MerchantKey        string   `json:"merchantKey"`
	MerchantPrivateKey string   `json:"merchantPrivateKey"`
	PlatformPublicKey  string   `json:"platformPublicKey"`
	EnabledPayTypes    []string `json:"enabledPayTypes"`
	NotifyUrl          string   `json:"notifyUrl"`
	ReturnUrl          string   `json:"returnUrl"`
}

type MailSetting struct {
	Enabled       bool   `json:"enabled"`
	Provider      string `json:"provider"`
	Host          string `json:"host"`
	Port          int    `json:"port"`
	Username      string `json:"username"`
	Password      string `json:"password"`
	FromEmail     string `json:"fromEmail"`
	FromName      string `json:"fromName"`
	Encryption    string `json:"encryption"`
	TestRecipient string `json:"testRecipient"`
}

type SystemSettingGeneratedUrls struct {
	NotifyUrl string `json:"notifyUrl"`
	ReturnUrl string `json:"returnUrl"`
}

type SystemSettingResponse struct {
	Site      SiteSetting                `json:"site"`
	Payment   PaymentSetting             `json:"payment"`
	Mail      MailSetting                `json:"mail"`
	Generated SystemSettingGeneratedUrls `json:"generated"`
}

type PaymentTestRequest struct {
	Action    string `json:"action"`
	NotifyUrl string `json:"notifyUrl"`
	ReturnUrl string `json:"returnUrl"`
}

type SystemConfigTestStep struct {
	Name      string      `json:"name"`
	Status    string      `json:"status"`
	Message   string      `json:"message"`
	Method    string      `json:"method,omitempty"`
	Url       string      `json:"url,omitempty"`
	Request   interface{} `json:"request,omitempty"`
	Response  interface{} `json:"response,omitempty"`
	ElapsedMs int64       `json:"elapsedMs"`
}

type PaymentTestResponse struct {
	Success    bool                   `json:"success"`
	Action     string                 `json:"action"`
	Version    string                 `json:"version"`
	OutTradeNo string                 `json:"outTradeNo"`
	TradeNo    string                 `json:"tradeNo"`
	PayType    string                 `json:"payType"`
	PayInfo    string                 `json:"payInfo"`
	Steps      []SystemConfigTestStep `json:"steps"`
}

type MailTestRequest struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

type MailTestResponse struct {
	Success   bool   `json:"success"`
	Message   string `json:"message"`
	Server    string `json:"server"`
	ElapsedMs int64  `json:"elapsedMs"`
}
