package service

import (
	"ruoyi-go/app/model"
	"ruoyi-go/common/types/constant"
	"ruoyi-go/framework/dal"
)

// EffectiveRoles returns functional roles, including admin -> common inheritance.
// Explicit assignments and data scope remain separate; never persist this union.
func EffectiveRoles(userId int) []model.SysRole {
	roles := make([]model.SysRole, 0)
	if userId <= 0 {
		return roles
	}
	if err := dal.Gorm.Model(model.SysRole{}).
		Joins("JOIN sys_user_role ON sys_user_role.role_id = sys_role.role_id").
		Where("sys_user_role.user_id = ? AND sys_role.status = ?", userId, constant.NORMAL_STATUS).
		Find(&roles).Error; err != nil {
		return nil
	}
	inherit, hasCommon := userId == 1, false
	for _, role := range roles {
		inherit = inherit || role.RoleKey == "admin"
		hasCommon = hasCommon || role.RoleKey == "common"
	}
	if inherit && !hasCommon {
		var common []model.SysRole
		if err := dal.Gorm.Where("role_key = ? AND status = ?", "common", constant.NORMAL_STATUS).Find(&common).Error; err != nil {
			return nil
		}
		roles = append(roles, common...)
	}
	return roles
}

func effectiveRoleIds(userId int) []int {
	ids := make([]int, 0)
	for _, role := range EffectiveRoles(userId) {
		ids = append(ids, role.RoleId)
	}
	return ids
}
