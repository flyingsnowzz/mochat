package rbac

import (
	"gorm.io/gorm"
	"mochat-api-server/internal/model"
)

// UserCan 检查用户是否拥有指定权限
func UserCan(db *gorm.DB, userID uint, permission string) bool {
	var count int64
	db.Raw(`
		SELECT COUNT(*) FROM mc_rbac_user_role ur
		JOIN mc_rbac_role_menu rm ON ur.role_id = rm.role_id
		JOIN mc_rbac_menu m ON rm.menu_id = m.id
		WHERE ur.user_id = ? AND m.link_url = ? AND m.status = 1
	`, userID, permission).Count(&count)
	return count > 0
}

// UserRoles 获取用户的所有有效角色
func UserRoles(db *gorm.DB, userID uint) []model.RbacRole {
	var roles []model.RbacRole
	db.Raw(`
		SELECT r.* FROM mc_rbac_role r
		JOIN mc_rbac_user_role ur ON r.id = ur.role_id
		WHERE ur.user_id = ? AND r.status = 1
	`, userID).Scan(&roles)
	return roles
}

// UserPermissions 获取用户的所有权限标识列表
func UserPermissions(db *gorm.DB, userID uint) []string {
	var permissions []string
	db.Raw(`
		SELECT DISTINCT m.link_url FROM mc_rbac_menu m
		JOIN mc_rbac_role_menu rm ON m.id = rm.menu_id
		JOIN mc_rbac_user_role ur ON rm.role_id = ur.role_id
		WHERE ur.user_id = ? AND m.status = 1 AND m.link_url != ''
	`, userID).Scan(&permissions)
	return permissions
}

// PermissionsCreate 为角色创建权限关联（先删后插，事务保护）
func PermissionsCreate(db *gorm.DB, roleID uint, menuIDs []uint) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("role_id = ?", roleID).Delete(&model.RbacRoleMenu{}).Error; err != nil {
			return err
		}
		for _, menuID := range menuIDs {
			roleMenu := model.RbacRoleMenu{RoleID: roleID, MenuID: menuID}
			if err := tx.Create(&roleMenu).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// PermissionsToRole 获取角色关联的菜单ID列表
func PermissionsToRole(db *gorm.DB, roleID uint) []uint {
	var menuIDs []uint
	db.Model(&model.RbacRoleMenu{}).Where("role_id = ?", roleID).Pluck("menu_id", &menuIDs)
	return menuIDs
}
