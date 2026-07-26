package model

import "time"

type PlatformPointAccount struct {
	ID           uint64    `gorm:"primaryKey;autoIncrement"`
	SystemUserID int       `gorm:"not null;uniqueIndex"`
	BalanceMinor uint64    `gorm:"not null"`
	Status       string    `gorm:"size:16;not null"`
	LockVersion  uint64    `gorm:"not null"`
	CreatedAt    time.Time `gorm:"not null"`
	UpdatedAt    time.Time `gorm:"not null"`
}

func (PlatformPointAccount) TableName() string { return "platform_point_accounts" }

type PlatformPointLedger struct {
	ID                uint64    `gorm:"primaryKey;autoIncrement"`
	PublicID          string    `gorm:"size:36;not null;uniqueIndex"`
	AccountID         uint64    `gorm:"not null;index"`
	EntryType         string    `gorm:"size:24;not null"`
	AmountMinor       int64     `gorm:"not null"`
	BalanceAfterMinor uint64    `gorm:"not null"`
	ReferenceType     string    `gorm:"size:32;not null"`
	ReferenceID       string    `gorm:"size:64;not null"`
	IdempotencyKey    string    `gorm:"size:64;not null"`
	Description       string    `gorm:"size:255;not null"`
	CreatedAt         time.Time `gorm:"not null"`
}

func (PlatformPointLedger) TableName() string { return "platform_point_ledger" }

type PlatformPointRechargeOrder struct {
	ID              uint64  `gorm:"primaryKey;autoIncrement"`
	PublicID        string  `gorm:"size:36;not null;uniqueIndex"`
	OrderNo         string  `gorm:"size:64;not null;uniqueIndex"`
	AccountID       uint64  `gorm:"not null;index"`
	SystemUserID    int     `gorm:"not null;index"`
	Points          uint64  `gorm:"not null"`
	AmountMinor     uint64  `gorm:"not null"`
	PayType         string  `gorm:"size:16;not null"`
	Status          string  `gorm:"size:16;not null;index"`
	IdempotencyKey  string  `gorm:"size:64;not null"`
	ExternalOrderID *string `gorm:"size:191"`
	PaidAt          *time.Time
	CreatedAt       time.Time `gorm:"not null"`
	UpdatedAt       time.Time `gorm:"not null"`
}

func (PlatformPointRechargeOrder) TableName() string { return "platform_point_recharge_orders" }

type PlatformRewardTask struct {
	ID           uint64    `gorm:"primaryKey;autoIncrement"`
	PublicID     string    `gorm:"size:36;not null;uniqueIndex"`
	TaskCode     string    `gorm:"size:64;not null;uniqueIndex"`
	Name         string    `gorm:"size:120;not null"`
	Summary      string    `gorm:"size:500;not null"`
	Provider     string    `gorm:"size:32;not null"`
	VerifyMode   string    `gorm:"size:32;not null"`
	RewardPoints uint64    `gorm:"not null"`
	ActionURL    string    `gorm:"size:500;not null"`
	Requirements string    `gorm:"size:1000;not null"`
	DisplayOrder int       `gorm:"not null"`
	Status       string    `gorm:"size:16;not null;index"`
	CreatedAt    time.Time `gorm:"not null"`
	UpdatedAt    time.Time `gorm:"not null"`
}

func (PlatformRewardTask) TableName() string { return "platform_reward_tasks" }

type PlatformRewardClaim struct {
	ID              uint64  `gorm:"primaryKey;autoIncrement"`
	PublicID        string  `gorm:"size:36;not null;uniqueIndex"`
	TaskID          uint64  `gorm:"not null;index"`
	AccountID       uint64  `gorm:"not null;index"`
	SystemUserID    int     `gorm:"not null;index"`
	ProviderSubject string  `gorm:"size:128;not null"`
	Status          string  `gorm:"size:16;not null;index"`
	ExchangeCodeID  *uint64 `gorm:"index"`
	ReviewNote      string  `gorm:"size:500;not null"`
	VerifiedAt      *time.Time
	IssuedAt        *time.Time
	CreatedAt       time.Time `gorm:"not null"`
	UpdatedAt       time.Time `gorm:"not null"`
}

func (PlatformRewardClaim) TableName() string { return "platform_reward_claims" }

type PlatformExchangeCode struct {
	ID                   uint64  `gorm:"primaryKey;autoIncrement"`
	PublicID             string  `gorm:"size:36;not null;uniqueIndex"`
	CodeHMAC             []byte  `gorm:"column:code_hmac;type:binary(32);not null;uniqueIndex"`
	CodeMask             string  `gorm:"size:40;not null"`
	RewardPoints         uint64  `gorm:"not null"`
	SourceType           string  `gorm:"size:32;not null"`
	SourceID             *uint64 `gorm:"index"`
	OwnerSystemUserID    *int    `gorm:"index"`
	Status               string  `gorm:"size:16;not null;index"`
	Description          string  `gorm:"size:255;not null"`
	ExpiresAt            *time.Time
	RedeemedAccountID    *uint64 `gorm:"index"`
	RedeemedSystemUserID *int    `gorm:"index"`
	RedeemedAt           *time.Time
	LedgerID             *uint64
	CreatedAt            time.Time `gorm:"not null"`
	UpdatedAt            time.Time `gorm:"not null"`
}

func (PlatformExchangeCode) TableName() string { return "platform_exchange_codes" }
