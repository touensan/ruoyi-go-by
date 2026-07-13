package model

import "ruoyi-go/framework/datetime"

type SysSystemSetting struct {
	SettingKey   string `gorm:"primaryKey;size:100"`
	SettingGroup string `gorm:"size:50;index"`
	SettingValue string `gorm:"type:longtext"`
	Remark       string
	CreateBy     string
	CreateTime   datetime.Datetime `gorm:"autoCreateTime"`
	UpdateBy     string
	UpdateTime   datetime.Datetime `gorm:"autoUpdateTime"`
}

func (SysSystemSetting) TableName() string {
	return "sys_system_setting"
}
