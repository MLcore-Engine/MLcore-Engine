// model/role.go
package model

// GetRoleName 返回角色名称
func GetRoleName(role int) string {
	switch role {
	case RoleRoot:
		return "超级管理员"
	case RoleAdmin:
		return "管理员"
	case RoleCommon:
		return "普通用户"
	default:
		return "未知"
	}
}

// IsValidRole 验证角色是否有效
func IsValidRole(role int) bool {
	validRoles := map[int]bool{
		RoleRoot:   true,
		RoleAdmin:  true,
		RoleCommon: true,
	}
	return validRoles[role]
}
