package service

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"ruoyi-go/app/dto"
	"ruoyi-go/app/model"
	"ruoyi-go/common/uuid"
	"ruoyi-go/config"
	"ruoyi-go/framework/dal"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	ErrPlatformDisabled     = errors.New("platform exchange is disabled")
	ErrPlatformInvalid      = errors.New("platform request is invalid")
	ErrPlatformCodeUsed     = errors.New("platform exchange code is used")
	ErrPlatformClaimed      = errors.New("platform reward is claimed")
	ErrPlatformProvider     = errors.New("platform reward provider is unavailable")
	ErrPlatformPayment      = errors.New("platform payment is unavailable")
	ErrPlatformInsufficient = errors.New("platform points are insufficient")
)

const platformSettingKey = "platform.exchange"

type platformStoredSettings struct {
	Enabled           bool   `json:"enabled"`
	RainyunEnabled    bool   `json:"rainyunEnabled"`
	RainyunAPIBaseURL string `json:"rainyunApiBaseUrl"`
	RainyunInviteURL  string `json:"rainyunInviteUrl"`
	RainyunAPIKey     string `json:"rainyunApiKey"`
}

type platformRainyunRecord struct {
	UserID uint64 `json:"user_id"`
}

type platformRainyunResponse struct {
	Data struct {
		Name         string                  `json:"Name"`
		Records      []platformRainyunRecord `json:"Records"`
		TotalRecords int                     `json:"TotalRecords"`
	} `json:"data"`
	Name         string                  `json:"Name"`
	Records      []platformRainyunRecord `json:"Records"`
	TotalRecords int                     `json:"TotalRecords"`
}

type PlatformService struct {
	db         *gorm.DB
	httpClient *http.Client
	clock      func() time.Time
	random     io.Reader
}

func NewPlatformService() *PlatformService {
	return &PlatformService{
		db: dal.Gorm, httpClient: &http.Client{Timeout: 30 * time.Second},
		clock: time.Now, random: rand.Reader,
	}
}

func (service *PlatformService) account(ctx context.Context, userID int) (model.PlatformPointAccount, error) {
	var account model.PlatformPointAccount
	err := service.db.WithContext(ctx).Where("system_user_id = ?", userID).Take(&account).Error
	if err == nil {
		if account.Status != "ACTIVE" {
			return model.PlatformPointAccount{}, ErrPlatformInvalid
		}
		return account, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return model.PlatformPointAccount{}, err
	}
	account = model.PlatformPointAccount{SystemUserID: userID, Status: "ACTIVE", LockVersion: 1}
	if err := service.db.WithContext(ctx).Create(&account).Error; err != nil {
		if findErr := service.db.WithContext(ctx).Where("system_user_id = ?", userID).Take(&account).Error; findErr != nil {
			return model.PlatformPointAccount{}, err
		}
	}
	return account, nil
}

func (service *PlatformService) Overview(ctx context.Context, userID int) (dto.PlatformOverviewResponse, error) {
	account, err := service.account(ctx, userID)
	if err != nil {
		return dto.PlatformOverviewResponse{}, err
	}
	result := dto.PlatformOverviewResponse{PointsMinor: account.BalanceMinor}
	service.db.WithContext(ctx).Model(&model.PlatformRewardTask{}).Where("status = ?", "ACTIVE").Count(&result.TaskCount)
	service.db.WithContext(ctx).Model(&model.PlatformRewardClaim{}).Where("system_user_id = ?", userID).Count(&result.ClaimCount)
	service.db.WithContext(ctx).Model(&model.PlatformPointRechargeOrder{}).Where("system_user_id = ?", userID).Count(&result.RechargeCount)
	return result, nil
}

func (service *PlatformService) Ledger(ctx context.Context, userID int) ([]dto.PlatformLedgerResponse, error) {
	account, err := service.account(ctx, userID)
	if err != nil {
		return nil, err
	}
	var rows []dto.PlatformLedgerResponse
	err = service.db.WithContext(ctx).Table("platform_point_ledger").
		Where("account_id = ?", account.ID).Order("id DESC").Limit(100).Scan(&rows).Error
	return rows, err
}

func (service *PlatformService) Recharges(ctx context.Context, userID int) ([]dto.PlatformRechargeResponse, error) {
	var orders []model.PlatformPointRechargeOrder
	if err := service.db.WithContext(ctx).Where("system_user_id = ?", userID).
		Order("id DESC").Limit(100).Find(&orders).Error; err != nil {
		return nil, err
	}
	result := make([]dto.PlatformRechargeResponse, 0, len(orders))
	for _, order := range orders {
		result = append(result, platformRechargeResponse(order))
	}
	return result, nil
}

func (service *PlatformService) CreateRecharge(
	ctx context.Context,
	userID int,
	request dto.CreatePlatformRechargeRequest,
	idempotencyKey string,
	clientIP string,
	baseURL string,
) (dto.PlatformRechargeResponse, error) {
	payType := strings.ToLower(strings.TrimSpace(request.PayType))
	if request.Points < 1 || request.Points > 100000 || (payType != "alipay" && payType != "wxpay") ||
		strings.TrimSpace(idempotencyKey) == "" {
		return dto.PlatformRechargeResponse{}, ErrPlatformInvalid
	}
	_, payment, _ := (&SystemSettingService{}).GetSettings()
	if !payment.Enabled || !containsPlatformPayType(payment.EnabledPayTypes, payType) {
		return dto.PlatformRechargeResponse{}, ErrPlatformPayment
	}
	account, err := service.account(ctx, userID)
	if err != nil {
		return dto.PlatformRechargeResponse{}, err
	}
	var order model.PlatformPointRechargeOrder
	err = service.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Where("account_id = ? AND idempotency_key = ?", account.ID, idempotencyKey).Take(&order).Error
		if err == nil {
			if order.Points != request.Points || order.PayType != payType {
				return ErrPlatformInvalid
			}
			return nil
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		publicID, err := uuid.New()
		if err != nil {
			return err
		}
		order = model.PlatformPointRechargeOrder{
			PublicID: publicID, OrderNo: "PT_" + strings.ReplaceAll(publicID, "-", ""),
			AccountID: account.ID, SystemUserID: userID, Points: request.Points,
			AmountMinor: request.Points * 100, PayType: payType, Status: "PENDING",
			IdempotencyKey: idempotencyKey,
		}
		return tx.Create(&order).Error
	})
	if err != nil {
		return dto.PlatformRechargeResponse{}, err
	}
	callbacks := (&SystemSettingService{}).GeneratedCallbackUrls(baseURL)
	payInfo, externalID, err := (&SystemSettingService{}).CreatePointPayment(
		payment, order.OrderNo, order.Points, order.PayType, clientIP,
		callbacks.NotifyUrl, callbacks.ReturnUrl,
	)
	if err != nil {
		return dto.PlatformRechargeResponse{}, ErrPlatformPayment
	}
	if externalID != "" {
		service.db.WithContext(ctx).Model(&order).Update("external_order_id", externalID)
		order.ExternalOrderID = &externalID
	}
	result := platformRechargeResponse(order)
	result.PayInfo = payInfo
	return result, nil
}

func (service *PlatformService) ApplyPaymentCallback(ctx context.Context, values url.Values) error {
	orderNo := strings.TrimSpace(values.Get("out_trade_no"))
	if !strings.HasPrefix(orderNo, "PT_") {
		return nil
	}
	if strings.ToUpper(strings.TrimSpace(values.Get("trade_status"))) != "TRADE_SUCCESS" {
		return ErrPlatformPayment
	}
	return service.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var order model.PlatformPointRechargeOrder
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("order_no = ?", orderNo).Take(&order).Error; err != nil {
			return err
		}
		if order.Status == "PAID" {
			return nil
		}
		if money := strings.TrimSpace(values.Get("money")); money != "" &&
			!platformPaymentAmountMatches(money, order.AmountMinor) {
			return ErrPlatformPayment
		}
		var account model.PlatformPointAccount
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ? AND status = ?", order.AccountID, "ACTIVE").Take(&account).Error; err != nil {
			return err
		}
		credit := order.Points * 100
		if account.BalanceMinor > math.MaxUint64-credit {
			return ErrPlatformInvalid
		}
		newBalance := account.BalanceMinor + credit
		publicID, err := uuid.New()
		if err != nil {
			return err
		}
		ledger := model.PlatformPointLedger{
			PublicID: publicID, AccountID: account.ID, EntryType: "RECHARGE",
			AmountMinor: int64(credit), BalanceAfterMinor: newBalance,
			ReferenceType: "POINT_RECHARGE", ReferenceID: order.PublicID,
			IdempotencyKey: "recharge:" + order.PublicID,
			Description:    fmt.Sprintf("充值到账 %d 积分", order.Points),
		}
		if err := tx.Create(&ledger).Error; err != nil {
			return err
		}
		if err := tx.Model(&account).Updates(map[string]interface{}{
			"balance_minor": newBalance, "lock_version": gorm.Expr("lock_version + 1"),
		}).Error; err != nil {
			return err
		}
		now := service.clock().UTC()
		return tx.Model(&order).Updates(map[string]interface{}{"status": "PAID", "paid_at": now}).Error
	})
}

func (service *PlatformService) settings(ctx context.Context) (platformStoredSettings, error) {
	setting := platformStoredSettings{
		Enabled: true, RainyunAPIBaseURL: "https://api.rainyun.com",
		RainyunInviteURL: "https://www.rainyun.com/",
	}
	var row model.SysSystemSetting
	err := service.db.WithContext(ctx).Where("setting_key = ?", platformSettingKey).Take(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return setting, nil
	}
	if err != nil {
		return setting, err
	}
	if err := json.Unmarshal([]byte(row.SettingValue), &setting); err != nil {
		return setting, err
	}
	setting.RainyunAPIBaseURL = strings.TrimRight(strings.TrimSpace(setting.RainyunAPIBaseURL), "/")
	setting.RainyunInviteURL = strings.TrimSpace(setting.RainyunInviteURL)
	return setting, nil
}

func (service *PlatformService) Settings(ctx context.Context) (dto.PlatformSettings, error) {
	setting, err := service.settings(ctx)
	if err != nil {
		return dto.PlatformSettings{}, err
	}
	return dto.PlatformSettings{
		Enabled: setting.Enabled, RainyunEnabled: setting.RainyunEnabled,
		RainyunAPIBaseURL: setting.RainyunAPIBaseURL, RainyunInviteURL: setting.RainyunInviteURL,
		RainyunKeyConfigured: setting.RainyunAPIKey != "",
	}, nil
}

func (service *PlatformService) SaveSettings(ctx context.Context, request dto.UpdatePlatformSettingsRequest, operator string) (dto.PlatformSettings, error) {
	baseURL := strings.TrimRight(strings.TrimSpace(request.RainyunAPIBaseURL), "/")
	if baseURL == "" {
		baseURL = "https://api.rainyun.com"
	}
	if parsed, err := url.Parse(baseURL); err != nil || parsed.Scheme != "https" || parsed.Host == "" {
		return dto.PlatformSettings{}, ErrPlatformInvalid
	}
	current, err := service.settings(ctx)
	if err != nil {
		return dto.PlatformSettings{}, err
	}
	apiKey := current.RainyunAPIKey
	if request.ClearRainyunKey {
		apiKey = ""
	} else if strings.TrimSpace(request.RainyunAPIKey) != "" {
		apiKey = strings.TrimSpace(request.RainyunAPIKey)
	}
	if request.RainyunEnabled && apiKey == "" {
		return dto.PlatformSettings{}, ErrPlatformProvider
	}
	stored := platformStoredSettings{
		Enabled: request.Enabled, RainyunEnabled: request.RainyunEnabled,
		RainyunAPIBaseURL: baseURL, RainyunInviteURL: strings.TrimSpace(request.RainyunInviteURL),
		RainyunAPIKey: apiKey,
	}
	body, _ := json.Marshal(stored)
	err = service.db.WithContext(ctx).Model(&model.SysSystemSetting{}).
		Where("setting_key = ?", platformSettingKey).Updates(map[string]interface{}{
		"setting_value": string(body), "update_by": operator, "update_time": service.clock().UTC(),
	}).Error
	if err != nil {
		return dto.PlatformSettings{}, err
	}
	return service.Settings(ctx)
}

func (service *PlatformService) Center(ctx context.Context, userID int) (dto.PlatformCenterResponse, error) {
	setting, err := service.settings(ctx)
	if err != nil {
		return dto.PlatformCenterResponse{}, err
	}
	if !setting.Enabled {
		return dto.PlatformCenterResponse{Tasks: []dto.PlatformTaskResponse{}, Claims: []dto.PlatformClaimResponse{}}, nil
	}
	account, err := service.account(ctx, userID)
	if err != nil {
		return dto.PlatformCenterResponse{}, err
	}
	tasks, err := service.tasks(ctx, account.ID, setting, false)
	if err != nil {
		return dto.PlatformCenterResponse{}, err
	}
	claims, _, err := service.claims(ctx, dto.PlatformAdminListRequest{PageRequest: dto.PageRequest{PageNum: 1, PageSize: 100}}, &userID)
	return dto.PlatformCenterResponse{Enabled: true, Tasks: tasks, Claims: claims}, err
}

func (service *PlatformService) Tasks(ctx context.Context) ([]dto.PlatformTaskResponse, error) {
	setting, err := service.settings(ctx)
	if err != nil {
		return nil, err
	}
	return service.tasks(ctx, 0, setting, true)
}

func (service *PlatformService) tasks(ctx context.Context, accountID uint64, setting platformStoredSettings, all bool) ([]dto.PlatformTaskResponse, error) {
	query := service.db.WithContext(ctx).Model(&model.PlatformRewardTask{})
	if !all {
		query = query.Where("status = ?", "ACTIVE")
	}
	var tasks []model.PlatformRewardTask
	if err := query.Order("display_order, id").Find(&tasks).Error; err != nil {
		return nil, err
	}
	statuses := map[uint64]string{}
	if accountID > 0 {
		var claims []model.PlatformRewardClaim
		service.db.WithContext(ctx).Where("account_id = ?", accountID).Find(&claims)
		for _, claim := range claims {
			statuses[claim.TaskID] = claim.Status
		}
	}
	result := make([]dto.PlatformTaskResponse, 0, len(tasks))
	for _, task := range tasks {
		available := setting.Enabled && task.Status == "ACTIVE"
		actionURL := task.ActionURL
		if task.Provider == "RAINYUN" {
			available = available && setting.RainyunEnabled && setting.RainyunAPIKey != ""
			if setting.RainyunInviteURL != "" {
				actionURL = setting.RainyunInviteURL
			}
		}
		result = append(result, dto.PlatformTaskResponse{
			PublicID: task.PublicID, TaskCode: task.TaskCode, Name: task.Name,
			Summary: task.Summary, Provider: task.Provider, VerifyMode: task.VerifyMode,
			RewardPoints: task.RewardPoints, ActionURL: actionURL, Requirements: task.Requirements,
			DisplayOrder: task.DisplayOrder, Status: task.Status, Available: available,
			ClaimStatus: statuses[task.ID], CreatedAt: task.CreatedAt,
		})
	}
	return result, nil
}

func (service *PlatformService) SaveTask(ctx context.Context, publicID string, request dto.SavePlatformTaskRequest) error {
	request.TaskCode = strings.ToUpper(strings.TrimSpace(request.TaskCode))
	request.Provider = strings.ToUpper(strings.TrimSpace(request.Provider))
	request.VerifyMode = strings.ToUpper(strings.TrimSpace(request.VerifyMode))
	request.Status = strings.ToUpper(strings.TrimSpace(request.Status))
	if request.TaskCode == "" || strings.TrimSpace(request.Name) == "" || request.RewardPoints == 0 ||
		(request.Provider != "MANUAL" && request.Provider != "RAINYUN") ||
		(request.VerifyMode != "MANUAL" && request.VerifyMode != "RAINYUN_SUBUSER") ||
		(request.Status != "ACTIVE" && request.Status != "INACTIVE") {
		return ErrPlatformInvalid
	}
	values := map[string]interface{}{
		"task_code": request.TaskCode, "name": strings.TrimSpace(request.Name),
		"summary": strings.TrimSpace(request.Summary), "provider": request.Provider,
		"verify_mode": request.VerifyMode, "reward_points": request.RewardPoints,
		"action_url": strings.TrimSpace(request.ActionURL), "requirements": strings.TrimSpace(request.Requirements),
		"display_order": request.DisplayOrder, "status": request.Status,
	}
	if publicID == "" {
		generated, err := uuid.New()
		if err != nil {
			return err
		}
		values["public_id"] = generated
		return service.db.WithContext(ctx).Model(&model.PlatformRewardTask{}).Create(values).Error
	}
	result := service.db.WithContext(ctx).Model(&model.PlatformRewardTask{}).Where("public_id = ?", publicID).Updates(values)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (service *PlatformService) SubmitClaim(ctx context.Context, userID int, request dto.SubmitPlatformClaimRequest) (dto.PlatformClaimResponse, string, error) {
	setting, err := service.settings(ctx)
	if err != nil {
		return dto.PlatformClaimResponse{}, "", err
	}
	if !setting.Enabled {
		return dto.PlatformClaimResponse{}, "", ErrPlatformDisabled
	}
	var task model.PlatformRewardTask
	if err := service.db.WithContext(ctx).Where("public_id = ? AND status = ?", request.TaskPublicID, "ACTIVE").Take(&task).Error; err != nil {
		return dto.PlatformClaimResponse{}, "", err
	}
	subject := strings.TrimSpace(request.ProviderSubject)
	if subject == "" {
		return dto.PlatformClaimResponse{}, "", ErrPlatformInvalid
	}
	verified := false
	if task.VerifyMode == "RAINYUN_SUBUSER" {
		if !setting.RainyunEnabled || setting.RainyunAPIKey == "" {
			return dto.PlatformClaimResponse{}, "", ErrPlatformProvider
		}
		uid, err := strconv.ParseUint(subject, 10, 64)
		if err != nil || uid == 0 {
			return dto.PlatformClaimResponse{}, "", ErrPlatformInvalid
		}
		verified, err = service.rainyunHasSubuser(ctx, setting, uid)
		if err != nil || !verified {
			return dto.PlatformClaimResponse{}, "", ErrPlatformProvider
		}
	}
	account, err := service.account(ctx, userID)
	if err != nil {
		return dto.PlatformClaimResponse{}, "", err
	}
	var claim model.PlatformRewardClaim
	var plaintext string
	err = service.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		publicID, _ := uuid.New()
		claim = model.PlatformRewardClaim{
			PublicID: publicID, TaskID: task.ID, AccountID: account.ID,
			SystemUserID: userID, ProviderSubject: subject, Status: "PENDING",
		}
		if verified {
			now := service.clock().UTC()
			claim.Status, claim.VerifiedAt, claim.IssuedAt = "ISSUED", &now, &now
		}
		if err := tx.Create(&claim).Error; err != nil {
			return ErrPlatformClaimed
		}
		if verified {
			code, generated, err := service.issueCode(tx, task.RewardPoints, "TASK", &claim.ID, &userID, "完成任务："+task.Name, nil)
			if err != nil {
				return err
			}
			plaintext = generated
			claim.ExchangeCodeID = &code.ID
			return tx.Model(&claim).Update("exchange_code_id", code.ID).Error
		}
		return nil
	})
	if err != nil {
		return dto.PlatformClaimResponse{}, "", err
	}
	row, err := service.claim(ctx, claim.PublicID)
	return row, plaintext, err
}

func (service *PlatformService) ReviewClaim(ctx context.Context, publicID string, request dto.ReviewPlatformClaimRequest) (string, error) {
	var plaintext string
	err := service.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var claim model.PlatformRewardClaim
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("public_id = ?", publicID).Take(&claim).Error; err != nil {
			return err
		}
		if claim.Status != "PENDING" {
			return ErrPlatformClaimed
		}
		if !request.Approved {
			return tx.Model(&claim).Updates(map[string]interface{}{"status": "REJECTED", "review_note": strings.TrimSpace(request.Note)}).Error
		}
		var task model.PlatformRewardTask
		if err := tx.Where("id = ?", claim.TaskID).Take(&task).Error; err != nil {
			return err
		}
		code, generated, err := service.issueCode(tx, task.RewardPoints, "TASK", &claim.ID, &claim.SystemUserID, "完成任务："+task.Name, nil)
		if err != nil {
			return err
		}
		plaintext = generated
		now := service.clock().UTC()
		return tx.Model(&claim).Updates(map[string]interface{}{
			"status": "ISSUED", "exchange_code_id": code.ID, "review_note": strings.TrimSpace(request.Note),
			"verified_at": now, "issued_at": now,
		}).Error
	})
	return plaintext, err
}

func (service *PlatformService) Claims(ctx context.Context, request dto.PlatformAdminListRequest) ([]dto.PlatformClaimResponse, int64, error) {
	return service.claims(ctx, request, nil)
}

func (service *PlatformService) claim(ctx context.Context, publicID string) (dto.PlatformClaimResponse, error) {
	var row dto.PlatformClaimResponse
	err := service.db.WithContext(ctx).Table("platform_reward_claims c").
		Select(`c.public_id, t.public_id AS task_public_id, t.name AS task_name,
			u.user_name AS username, t.provider, c.provider_subject, t.reward_points,
			c.status, e.code_mask, c.review_note, c.verified_at, c.issued_at, c.created_at`).
		Joins("JOIN platform_reward_tasks t ON t.id = c.task_id").
		Joins("JOIN sys_user u ON u.user_id = c.system_user_id").
		Joins("LEFT JOIN platform_exchange_codes e ON e.id = c.exchange_code_id").
		Where("c.public_id = ?", publicID).Take(&row).Error
	return row, err
}

func (service *PlatformService) claims(ctx context.Context, request dto.PlatformAdminListRequest, userID *int) ([]dto.PlatformClaimResponse, int64, error) {
	pageNum, pageSize := normalizePlatformPage(request.PageNum, request.PageSize)
	base := service.db.WithContext(ctx).Table("platform_reward_claims c").
		Joins("JOIN platform_reward_tasks t ON t.id = c.task_id").
		Joins("JOIN sys_user u ON u.user_id = c.system_user_id").
		Joins("LEFT JOIN platform_exchange_codes e ON e.id = c.exchange_code_id")
	if userID != nil {
		base = base.Where("c.system_user_id = ?", *userID)
	}
	if request.Status != "" {
		base = base.Where("c.status = ?", strings.ToUpper(request.Status))
	}
	if request.Username != "" {
		base = base.Where("u.user_name LIKE ?", "%"+request.Username+"%")
	}
	var total int64
	if err := base.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []dto.PlatformClaimResponse
	err := base.Select(`c.public_id, t.public_id AS task_public_id, t.name AS task_name,
			u.user_name AS username, t.provider, c.provider_subject, t.reward_points,
			c.status, e.code_mask, c.review_note, c.verified_at, c.issued_at, c.created_at`).
		Order("c.id DESC").Offset((pageNum - 1) * pageSize).Limit(pageSize).Scan(&rows).Error
	return rows, total, err
}

func (service *PlatformService) GenerateCodes(ctx context.Context, request dto.GeneratePlatformCodesRequest) ([]dto.GeneratedPlatformCodeResponse, error) {
	if request.RewardPoints == 0 || request.Count < 1 || request.Count > 100 {
		return nil, ErrPlatformInvalid
	}
	result := make([]dto.GeneratedPlatformCodeResponse, 0, request.Count)
	err := service.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for i := 0; i < request.Count; i++ {
			code, plaintext, err := service.issueCode(tx, request.RewardPoints, "MANUAL", nil, nil, strings.TrimSpace(request.Description), request.ExpiresAt)
			if err != nil {
				return err
			}
			result = append(result, dto.GeneratedPlatformCodeResponse{PublicID: code.PublicID, Code: plaintext, CodeMask: code.CodeMask, RewardPoints: code.RewardPoints})
		}
		return nil
	})
	return result, err
}

func (service *PlatformService) issueCode(tx *gorm.DB, points uint64, sourceType string, sourceID *uint64, owner *int, description string, expiresAt *time.Time) (model.PlatformExchangeCode, string, error) {
	plaintext, mask, digest, err := service.generateCode()
	if err != nil {
		return model.PlatformExchangeCode{}, "", err
	}
	publicID, _ := uuid.New()
	code := model.PlatformExchangeCode{
		PublicID: publicID, CodeHMAC: digest, CodeMask: mask, RewardPoints: points,
		SourceType: sourceType, SourceID: sourceID, OwnerSystemUserID: owner,
		Status: "UNUSED", Description: description, ExpiresAt: expiresAt,
	}
	err = tx.Create(&code).Error
	return code, plaintext, err
}

func (service *PlatformService) generateCode() (string, string, []byte, error) {
	randomBytes := make([]byte, 15)
	if _, err := io.ReadFull(service.random, randomBytes); err != nil {
		return "", "", nil, err
	}
	encoded := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)
	groups := make([]string, 0, 6)
	for i := 0; i < len(encoded); i += 4 {
		groups = append(groups, encoded[i:i+4])
	}
	plaintext := "R1POINT-" + strings.Join(groups, "-")
	mask := "R1POINT-" + groups[0] + "-****-****-****-****-" + groups[len(groups)-1]
	return plaintext, mask, service.codeDigest(plaintext), nil
}

func (service *PlatformService) codeDigest(value string) []byte {
	key := sha256.Sum256([]byte("R1GO-POINT-CODE-KEY\x00" + config.Data.Token.Secret))
	mac := hmac.New(sha256.New, key[:])
	mac.Write([]byte("R1GO-POINT-CODE-V1\x00"))
	mac.Write([]byte(value))
	return mac.Sum(nil)
}

func normalizePlatformCode(value string) (string, error) {
	normalized := strings.ToUpper(strings.TrimSpace(value))
	if !strings.HasPrefix(normalized, "R1POINT-") || len(normalized) > 80 {
		return "", ErrPlatformInvalid
	}
	for _, character := range normalized {
		if (character < 'A' || character > 'Z') && (character < '2' || character > '7') && (character < '0' || character > '1') && character != '-' {
			return "", ErrPlatformInvalid
		}
	}
	return normalized, nil
}

func (service *PlatformService) Redeem(ctx context.Context, userID int, rawCode string) (dto.RedeemPlatformCodeResponse, error) {
	setting, err := service.settings(ctx)
	if err != nil {
		return dto.RedeemPlatformCodeResponse{}, err
	}
	if !setting.Enabled {
		return dto.RedeemPlatformCodeResponse{}, ErrPlatformDisabled
	}
	codeValue, err := normalizePlatformCode(rawCode)
	if err != nil {
		return dto.RedeemPlatformCodeResponse{}, err
	}
	account, err := service.account(ctx, userID)
	if err != nil {
		return dto.RedeemPlatformCodeResponse{}, err
	}
	var result dto.RedeemPlatformCodeResponse
	err = service.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var code model.PlatformExchangeCode
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("code_hmac = ?", service.codeDigest(codeValue)).Take(&code).Error; err != nil {
			return ErrPlatformInvalid
		}
		if code.Status != "UNUSED" {
			return ErrPlatformCodeUsed
		}
		now := service.clock().UTC()
		if code.ExpiresAt != nil && !code.ExpiresAt.After(now) {
			return ErrPlatformInvalid
		}
		var locked model.PlatformPointAccount
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", account.ID).Take(&locked).Error; err != nil {
			return err
		}
		credit := code.RewardPoints * 100
		if locked.BalanceMinor > math.MaxUint64-credit {
			return ErrPlatformInvalid
		}
		newBalance := locked.BalanceMinor + credit
		publicID, _ := uuid.New()
		ledger := model.PlatformPointLedger{
			PublicID: publicID, AccountID: locked.ID, EntryType: "CODE_REDEMPTION",
			AmountMinor: int64(credit), BalanceAfterMinor: newBalance,
			ReferenceType: "EXCHANGE_CODE", ReferenceID: code.PublicID,
			IdempotencyKey: "exchange:" + code.PublicID,
			Description:    fmt.Sprintf("兑换码到账 %d 积分", code.RewardPoints),
		}
		if err := tx.Create(&ledger).Error; err != nil {
			return err
		}
		if err := tx.Model(&locked).Updates(map[string]interface{}{"balance_minor": newBalance, "lock_version": gorm.Expr("lock_version + 1")}).Error; err != nil {
			return err
		}
		if err := tx.Model(&code).Updates(map[string]interface{}{
			"status": "REDEEMED", "redeemed_account_id": locked.ID,
			"redeemed_system_user_id": userID, "redeemed_at": now, "ledger_id": ledger.ID,
		}).Error; err != nil {
			return err
		}
		result = dto.RedeemPlatformCodeResponse{Points: code.RewardPoints, BalanceAfterMinor: newBalance, CodeMask: code.CodeMask}
		return nil
	})
	return result, err
}

func (service *PlatformService) Accounts(ctx context.Context, request dto.PlatformAdminListRequest) ([]dto.PlatformAccountAdminResponse, int64, error) {
	pageNum, pageSize := normalizePlatformPage(request.PageNum, request.PageSize)
	base := service.db.WithContext(ctx).Table("platform_point_accounts a").Joins("JOIN sys_user u ON u.user_id = a.system_user_id")
	if request.Username != "" {
		base = base.Where("u.user_name LIKE ? OR u.nick_name LIKE ?", "%"+request.Username+"%", "%"+request.Username+"%")
	}
	var total int64
	if err := base.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []dto.PlatformAccountAdminResponse
	err := base.Select("a.system_user_id, u.user_name AS username, u.nick_name AS nickname, a.balance_minor AS points_minor, a.status, a.updated_at").
		Order("a.id DESC").Offset((pageNum - 1) * pageSize).Limit(pageSize).Scan(&rows).Error
	return rows, total, err
}

func (service *PlatformService) Adjust(ctx context.Context, request dto.AdjustPlatformPointsRequest) error {
	if request.SystemUserID <= 0 || request.Points == 0 || strings.TrimSpace(request.Reason) == "" {
		return ErrPlatformInvalid
	}
	account, err := service.account(ctx, request.SystemUserID)
	if err != nil {
		return err
	}
	return service.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var locked model.PlatformPointAccount
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", account.ID).Take(&locked).Error; err != nil {
			return err
		}
		change := request.Points * 100
		var balance uint64
		if change < 0 {
			if locked.BalanceMinor < uint64(-change) {
				return ErrPlatformInsufficient
			}
			balance = locked.BalanceMinor - uint64(-change)
		} else {
			balance = locked.BalanceMinor + uint64(change)
		}
		publicID, _ := uuid.New()
		ledger := model.PlatformPointLedger{
			PublicID: publicID, AccountID: locked.ID, EntryType: "ADJUSTMENT",
			AmountMinor: change, BalanceAfterMinor: balance, ReferenceType: "ADMIN_ADJUST",
			ReferenceID:    fmt.Sprintf("user:%d", request.SystemUserID),
			IdempotencyKey: publicID, Description: strings.TrimSpace(request.Reason),
		}
		if err := tx.Create(&ledger).Error; err != nil {
			return err
		}
		return tx.Model(&locked).Updates(map[string]interface{}{"balance_minor": balance, "lock_version": gorm.Expr("lock_version + 1")}).Error
	})
}

func (service *PlatformService) Codes(ctx context.Context, request dto.PlatformAdminListRequest) ([]dto.PlatformCodeAdminResponse, int64, error) {
	pageNum, pageSize := normalizePlatformPage(request.PageNum, request.PageSize)
	base := service.db.WithContext(ctx).Table("platform_exchange_codes e").
		Joins("LEFT JOIN sys_user owner ON owner.user_id = e.owner_system_user_id").
		Joins("LEFT JOIN sys_user redeemed ON redeemed.user_id = e.redeemed_system_user_id")
	if request.Status != "" {
		base = base.Where("e.status = ?", strings.ToUpper(request.Status))
	}
	var total int64
	if err := base.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []dto.PlatformCodeAdminResponse
	err := base.Select(`e.public_id, e.code_mask, e.reward_points, e.source_type,
			owner.user_name AS owner_username, e.status, redeemed.user_name AS redeemed_by,
			e.description, e.created_at, e.redeemed_at`).
		Order("e.id DESC").Offset((pageNum - 1) * pageSize).Limit(pageSize).Scan(&rows).Error
	return rows, total, err
}

func (service *PlatformService) TestRainyun(ctx context.Context) (map[string]interface{}, error) {
	setting, err := service.settings(ctx)
	if err != nil || !setting.RainyunEnabled || setting.RainyunAPIKey == "" {
		return nil, ErrPlatformProvider
	}
	request, _ := http.NewRequestWithContext(ctx, http.MethodGet, setting.RainyunAPIBaseURL+"/user/", nil)
	request.Header.Set("X-Api-Key", setting.RainyunAPIKey)
	response, err := service.httpClient.Do(request)
	if err != nil {
		return nil, ErrPlatformProvider
	}
	defer response.Body.Close()
	body, _ := io.ReadAll(io.LimitReader(response.Body, 1<<20))
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return nil, ErrPlatformProvider
	}
	var payload platformRainyunResponse
	if json.Unmarshal(body, &payload) != nil {
		return nil, ErrPlatformProvider
	}
	return map[string]interface{}{"connected": true, "account": maskPlatformName(firstPlatformNonEmpty(payload.Data.Name, payload.Name))}, nil
}

func (service *PlatformService) rainyunHasSubuser(ctx context.Context, setting platformStoredSettings, uid uint64) (bool, error) {
	for page := 1; page <= 1000; page++ {
		options := fmt.Sprintf(`{"columnFilters":{},"sort":[{"field":"user_id","type":"asc"}],"page":%d,"perPage":100}`, page)
		request, _ := http.NewRequestWithContext(ctx, http.MethodGet, setting.RainyunAPIBaseURL+"/user/vip/subuser_sale?options="+url.QueryEscape(options), nil)
		request.Header.Set("X-Api-Key", setting.RainyunAPIKey)
		response, err := service.httpClient.Do(request)
		if err != nil {
			return false, ErrPlatformProvider
		}
		body, _ := io.ReadAll(io.LimitReader(response.Body, 4<<20))
		response.Body.Close()
		if response.StatusCode < 200 || response.StatusCode >= 300 {
			return false, ErrPlatformProvider
		}
		var payload platformRainyunResponse
		if json.Unmarshal(body, &payload) != nil {
			return false, ErrPlatformProvider
		}
		records, total := payload.Data.Records, payload.Data.TotalRecords
		if len(records) == 0 {
			records, total = payload.Records, payload.TotalRecords
		}
		for _, record := range records {
			if record.UserID == uid {
				return true, nil
			}
		}
		if page*100 >= total {
			break
		}
	}
	return false, nil
}

func normalizePlatformPage(pageNum, pageSize int) (int, int) {
	if pageNum < 1 {
		pageNum = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return pageNum, pageSize
}

func platformRechargeResponse(order model.PlatformPointRechargeOrder) dto.PlatformRechargeResponse {
	return dto.PlatformRechargeResponse{
		PublicID: order.PublicID, OrderNo: order.OrderNo, Points: order.Points,
		PayType: order.PayType, Status: order.Status, PaidAt: order.PaidAt, CreatedAt: order.CreatedAt,
	}
}

func containsPlatformPayType(values []string, target string) bool {
	for _, value := range values {
		if strings.EqualFold(strings.TrimSpace(value), target) {
			return true
		}
	}
	return false
}

func platformPaymentAmountMatches(value string, expectedMinor uint64) bool {
	amount, err := strconv.ParseFloat(strings.TrimSpace(value), 64)
	amountMinor := math.Round(amount * 100)
	return err == nil && !math.IsNaN(amountMinor) && !math.IsInf(amountMinor, 0) &&
		amountMinor >= 0 && amountMinor <= math.MaxUint64 && uint64(amountMinor) == expectedMinor
}

func firstPlatformNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func maskPlatformName(value string) string {
	runes := []rune(strings.TrimSpace(value))
	if len(runes) == 0 {
		return "-"
	}
	if len(runes) == 1 {
		return "*"
	}
	return string(runes[0]) + "***"
}
