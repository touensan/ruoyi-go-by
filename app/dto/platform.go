package dto

import "time"

type PlatformOverviewResponse struct {
	PointsMinor   uint64 `json:"pointsMinor"`
	TaskCount     int64  `json:"taskCount"`
	ClaimCount    int64  `json:"claimCount"`
	RechargeCount int64  `json:"rechargeCount"`
}

type PlatformLedgerResponse struct {
	PublicID          string    `json:"publicId"`
	EntryType         string    `json:"entryType"`
	AmountMinor       int64     `json:"amountMinor"`
	BalanceAfterMinor uint64    `json:"balanceAfterMinor"`
	Description       string    `json:"description"`
	CreatedAt         time.Time `json:"createdAt"`
}

type CreatePlatformRechargeRequest struct {
	Points  uint64 `json:"points"`
	PayType string `json:"payType"`
}

type PlatformRechargeResponse struct {
	PublicID  string     `json:"publicId"`
	OrderNo   string     `json:"orderNo"`
	Points    uint64     `json:"points"`
	PayType   string     `json:"payType"`
	Status    string     `json:"status"`
	PayInfo   string     `json:"payInfo,omitempty"`
	PaidAt    *time.Time `json:"paidAt,omitempty"`
	CreatedAt time.Time  `json:"createdAt"`
}

type PlatformSettings struct {
	Enabled              bool   `json:"enabled"`
	RainyunEnabled       bool   `json:"rainyunEnabled"`
	RainyunAPIBaseURL    string `json:"rainyunApiBaseUrl"`
	RainyunInviteURL     string `json:"rainyunInviteUrl"`
	RainyunAPIKey        string `json:"rainyunApiKey,omitempty"`
	RainyunKeyConfigured bool   `json:"rainyunKeyConfigured"`
}

type UpdatePlatformSettingsRequest struct {
	Enabled           bool   `json:"enabled"`
	RainyunEnabled    bool   `json:"rainyunEnabled"`
	RainyunAPIBaseURL string `json:"rainyunApiBaseUrl"`
	RainyunInviteURL  string `json:"rainyunInviteUrl"`
	RainyunAPIKey     string `json:"rainyunApiKey"`
	ClearRainyunKey   bool   `json:"clearRainyunKey"`
}

type PlatformTaskResponse struct {
	PublicID     string    `json:"publicId"`
	TaskCode     string    `json:"taskCode"`
	Name         string    `json:"name"`
	Summary      string    `json:"summary"`
	Provider     string    `json:"provider"`
	VerifyMode   string    `json:"verifyMode"`
	RewardPoints uint64    `json:"rewardPoints"`
	ActionURL    string    `json:"actionUrl"`
	Requirements string    `json:"requirements"`
	DisplayOrder int       `json:"displayOrder"`
	Status       string    `json:"status"`
	Available    bool      `json:"available"`
	ClaimStatus  string    `json:"claimStatus,omitempty"`
	CreatedAt    time.Time `json:"createdAt"`
}

type SavePlatformTaskRequest struct {
	TaskCode     string `json:"taskCode"`
	Name         string `json:"name"`
	Summary      string `json:"summary"`
	Provider     string `json:"provider"`
	VerifyMode   string `json:"verifyMode"`
	RewardPoints uint64 `json:"rewardPoints"`
	ActionURL    string `json:"actionUrl"`
	Requirements string `json:"requirements"`
	DisplayOrder int    `json:"displayOrder"`
	Status       string `json:"status"`
}

type SubmitPlatformClaimRequest struct {
	TaskPublicID    string `json:"taskPublicId"`
	ProviderSubject string `json:"providerSubject"`
}

type PlatformClaimResponse struct {
	PublicID        string     `json:"publicId"`
	TaskPublicID    string     `json:"taskPublicId"`
	TaskName        string     `json:"taskName"`
	Username        string     `json:"username,omitempty"`
	Provider        string     `json:"provider"`
	ProviderSubject string     `json:"providerSubject"`
	RewardPoints    uint64     `json:"rewardPoints"`
	Status          string     `json:"status"`
	CodeMask        string     `json:"codeMask,omitempty"`
	ReviewNote      string     `json:"reviewNote,omitempty"`
	VerifiedAt      *time.Time `json:"verifiedAt,omitempty"`
	IssuedAt        *time.Time `json:"issuedAt,omitempty"`
	CreatedAt       time.Time  `json:"createdAt"`
}

type PlatformCenterResponse struct {
	Enabled bool                    `json:"enabled"`
	Tasks   []PlatformTaskResponse  `json:"tasks"`
	Claims  []PlatformClaimResponse `json:"claims"`
}

type RedeemPlatformCodeRequest struct {
	Code string `json:"code"`
}

type RedeemPlatformCodeResponse struct {
	Points            uint64 `json:"points"`
	BalanceAfterMinor uint64 `json:"balanceAfterMinor"`
	CodeMask          string `json:"codeMask"`
}

type GeneratePlatformCodesRequest struct {
	RewardPoints uint64     `json:"rewardPoints"`
	Count        int        `json:"count"`
	Description  string     `json:"description"`
	ExpiresAt    *time.Time `json:"expiresAt"`
}

type GeneratedPlatformCodeResponse struct {
	PublicID     string `json:"publicId"`
	Code         string `json:"code"`
	CodeMask     string `json:"codeMask"`
	RewardPoints uint64 `json:"rewardPoints"`
}

type PlatformCodeAdminResponse struct {
	PublicID      string     `json:"publicId"`
	CodeMask      string     `json:"codeMask"`
	RewardPoints  uint64     `json:"rewardPoints"`
	SourceType    string     `json:"sourceType"`
	OwnerUsername string     `json:"ownerUsername,omitempty"`
	Status        string     `json:"status"`
	RedeemedBy    string     `json:"redeemedBy,omitempty"`
	Description   string     `json:"description"`
	CreatedAt     time.Time  `json:"createdAt"`
	RedeemedAt    *time.Time `json:"redeemedAt,omitempty"`
}

type PlatformAdminListRequest struct {
	PageRequest
	Username string `query:"username" form:"username"`
	Status   string `query:"status" form:"status"`
}

type PlatformAccountAdminResponse struct {
	SystemUserID int       `json:"systemUserId"`
	Username     string    `json:"username"`
	Nickname     string    `json:"nickname"`
	PointsMinor  uint64    `json:"pointsMinor"`
	Status       string    `json:"status"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type AdjustPlatformPointsRequest struct {
	SystemUserID int    `json:"systemUserId"`
	Points       int64  `json:"points"`
	Reason       string `json:"reason"`
}

type ReviewPlatformClaimRequest struct {
	Approved bool   `json:"approved"`
	Note     string `json:"note"`
}
