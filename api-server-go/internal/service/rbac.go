package service

import (
	"gorm.io/gorm"
	"mochat-api-server/internal/model"
)

// RbacMenuService 菜单 Service
// 提供菜单（RbacMenu）的 CRUD 操作功能
// 主要职责：
// 1. 获取菜单列表（支持按名称搜索）
// 2. 根据 ID 获取菜单详情
// 3. 创建菜单
// 4. 更新菜单
// 5. 删除菜单
// 6. 获取菜单选择列表（只获取状态为启用且为页面菜单的）
// 7. 根据 ID 列表获取菜单
// 8. 根据 ID 更新菜单指定字段
//
// 依赖：
// - gorm.DB: 数据库连接

type RbacMenuService struct {
	db *gorm.DB // 数据库连接
}

// NewRbacMenuService 创建菜单 Service 实例
// 参数：db - GORM 数据库连接
// 返回：菜单 Service 实例
func NewRbacMenuService(db *gorm.DB) *RbacMenuService {
	return &RbacMenuService{db: db}
}

// List 获取菜单列表
// 查询菜单列表，支持按名称搜索，按排序字段和 ID 升序排列
// 参数：
//
//	name - 菜单名称（模糊搜索），空字符串表示不限制
//
// 返回：菜单列表和错误信息
func (s *RbacMenuService) List(name string) ([]model.RbacMenu, error) {
	var menus []model.RbacMenu
	query := s.db.Model(&model.RbacMenu{})
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if err := query.Order("sort ASC, id ASC").Find(&menus).Error; err != nil {
		return nil, err
	}
	return menus, nil
}

// GetByID 根据 ID 获取菜单详情
// 查询指定 ID 的菜单
// 参数：
//
//	id - 菜单 ID
//
// 返回：菜单实例和错误信息
func (s *RbacMenuService) GetByID(id uint) (*model.RbacMenu, error) {
	var menu model.RbacMenu
	if err := s.db.First(&menu, id).Error; err != nil {
		return nil, err
	}
	return &menu, nil
}

// Create 创建菜单
// 将菜单信息保存到数据库
// 参数：
//
//	menu - 菜单实例
//
// 返回：错误信息
func (s *RbacMenuService) Create(menu *model.RbacMenu) error {
	return s.db.Create(menu).Error
}

// Update 更新菜单
// 更新数据库中的菜单信息
// 参数：
//
//	menu - 菜单实例
//
// 返回：错误信息
func (s *RbacMenuService) Update(menu *model.RbacMenu) error {
	return s.db.Save(menu).Error
}

// Delete 删除菜单
// 从数据库中删除指定 ID 的菜单
// 参数：
//
//	id - 菜单 ID
//
// 返回：错误信息
func (s *RbacMenuService) Delete(id uint) error {
	return s.db.Delete(&model.RbacMenu{}, id).Error
}

// Select 获取菜单选择列表
// 查询状态为启用且为页面菜单的菜单列表，用于下拉选择
// 返回：菜单列表和错误信息
func (s *RbacMenuService) Select() ([]model.RbacMenu, error) {
	var menus []model.RbacMenu
	if err := s.db.Where("status = 1 AND is_page_menu = 1").Order("sort ASC").Find(&menus).Error; err != nil {
		return nil, err
	}
	return menus, nil
}

// ListByIDs 根据 ID 列表获取菜单
// 查询指定 ID 列表的菜单，只返回状态为启用的菜单
// 参数：
//
//	ids - 菜单 ID 列表
//
// 返回：菜单列表和错误信息
func (s *RbacMenuService) ListByIDs(ids []uint) ([]model.RbacMenu, error) {
	var menus []model.RbacMenu
	if err := s.db.Where("id IN ? AND status = 1", ids).Find(&menus).Error; err != nil {
		return nil, err
	}
	return menus, nil
}

// UpdateByID 根据 ID 更新菜单指定字段
// 更新菜单的指定字段
// 参数：
//
//	id - 菜单 ID
//	updates - 待更新的字段映射
//
// 返回：错误信息
func (s *RbacMenuService) UpdateByID(id uint, updates map[string]interface{}) error {
	return s.db.Model(&model.RbacMenu{}).Where("id = ?", id).Updates(updates).Error
}

// RbacRoleService 角色 Service
// 提供角色（RbacRole）的 CRUD 操作功能
// 主要职责：
// 1. 获取角色列表（分页）
// 2. 根据 ID 获取角色详情
// 3. 创建角色
// 4. 更新角色
// 5. 删除角色（同时删除关联的菜单和用户角色）
// 6. 获取角色选择列表
// 7. 使用 offset 方式获取角色列表
// 8. 根据 ID 更新角色指定字段
// 9. 设置角色的菜单权限
// 10. 获取角色的菜单 ID 列表
//
// 依赖：
// - gorm.DB: 数据库连接

type RbacRoleService struct {
	db *gorm.DB // 数据库连接
}

// NewRbacRoleService 创建角色 Service 实例
// 参数：db - GORM 数据库连接
// 返回：角色 Service 实例
func NewRbacRoleService(db *gorm.DB) *RbacRoleService {
	return &RbacRoleService{db: db}
}

// List 获取角色列表（分页）
// 查询角色列表，支持按租户 ID 筛选
// 参数：
//
//	tenantID - 租户 ID，0 表示不限制
//	page - 页码
//	pageSize - 每页数量
//
// 返回：角色列表、总数和错误信息
func (s *RbacRoleService) List(tenantID uint, page, pageSize int) ([]model.RbacRole, int64, error) {
	var roles []model.RbacRole
	var total int64
	query := s.db.Model(&model.RbacRole{})
	if tenantID > 0 {
		query = query.Where("tenant_id = ?", tenantID)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Find(&roles).Error; err != nil {
		return nil, 0, err
	}
	return roles, total, nil
}

// GetByID 根据 ID 获取角色详情
// 查询指定 ID 的角色
// 参数：
//
//	id - 角色 ID
//
// 返回：角色实例和错误信息
func (s *RbacRoleService) GetByID(id uint) (*model.RbacRole, error) {
	var role model.RbacRole
	if err := s.db.First(&role, id).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

// Create 创建角色
// 将角色信息保存到数据库
// 参数：
//
//	role - 角色实例
//
// 返回：错误信息
func (s *RbacRoleService) Create(role *model.RbacRole) error {
	return s.db.Create(role).Error
}

// Update 更新角色
// 更新数据库中的角色信息
// 参数：
//
//	role - 角色实例
//
// 返回：错误信息
func (s *RbacRoleService) Update(role *model.RbacRole) error {
	return s.db.Save(role).Error
}

// Delete 删除角色
// 从数据库中删除指定 ID 的角色
// 参数：
//
//	id - 角色 ID
//
// 返回：错误信息
func (s *RbacRoleService) Delete(id uint) error {
	return s.db.Delete(&model.RbacRole{}, id).Error
}

// Select 获取角色选择列表
// 查询状态为启用的角色列表，用于下拉选择
// 参数：
//
//	tenantID - 租户 ID，0 表示不限制
//
// 返回：角色列表和错误信息
func (s *RbacRoleService) Select(tenantID uint) ([]model.RbacRole, error) {
	var roles []model.RbacRole
	query := s.db.Model(&model.RbacRole{}).Where("status = 1")
	if tenantID > 0 {
		query = query.Where("tenant_id = ?", tenantID)
	}
	if err := query.Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

// ListByOffset 使用 offset 方式获取角色列表
// 使用 offset 和 limit 方式查询角色列表
// 参数：
//
//	tenantID - 租户 ID，0 表示不限制
//	offset - 偏移量
//	limit - 限制数量
//
// 返回：角色列表、总数和错误信息
func (s *RbacRoleService) ListByOffset(tenantID uint, offset, limit int) ([]model.RbacRole, int64, error) {
	var roles []model.RbacRole
	var total int64
	query := s.db.Model(&model.RbacRole{})
	if tenantID > 0 {
		query = query.Where("tenant_id = ?", tenantID)
	}
	query.Count(&total)
	if err := query.Offset(offset).Limit(limit).Order("id DESC").Find(&roles).Error; err != nil {
		return nil, 0, err
	}
	return roles, total, nil
}

// UpdateByID 根据 ID 更新角色指定字段
// 更新角色的指定字段
// 参数：
//
//	id - 角色 ID
//	updates - 待更新的字段映射
//
// 返回：错误信息
func (s *RbacRoleService) UpdateByID(id uint, updates map[string]interface{}) error {
	return s.db.Model(&model.RbacRole{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteWithRelations 删除角色（同时删除关联数据）
// 删除指定 ID 的角色，同时删除关联的菜单和用户角色
// 使用事务确保数据一致性
// 参数：
//
//	id - 角色 ID
//
// 返回：错误信息
func (s *RbacRoleService) DeleteWithRelations(id uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 删除角色菜单关联
		tx.Where("role_id = ?", id).Delete(&model.RbacRoleMenu{})
		// 删除用户角色关联
		tx.Where("role_id = ?", id).Delete(&model.RbacUserRole{})
		// 删除角色
		return tx.Delete(&model.RbacRole{}, id).Error
	})
}

// SetMenus 设置角色的菜单权限
// 设置指定角色的菜单权限，先删除旧关联再创建新关联
// 使用事务确保数据一致性
// 参数：
//
//	roleID - 角色 ID
//	menuIDs - 菜单 ID 列表
//
// 返回：错误信息
func (s *RbacRoleService) SetMenus(roleID uint, menuIDs []uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 删除旧的角色菜单关联
		tx.Where("role_id = ?", roleID).Delete(&model.RbacRoleMenu{})
		// 创建新的角色菜单关联
		for _, menuID := range menuIDs {
			if err := tx.Create(&model.RbacRoleMenu{RoleID: roleID, MenuID: menuID}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// GetMenus 获取角色的菜单 ID 列表
// 查询指定角色关联的所有菜单 ID
// 参数：
//
//	roleID - 角色 ID
//
// 返回：菜单 ID 列表和错误信息
func (s *RbacRoleService) GetMenus(roleID uint) ([]uint, error) {
	var menuIDs []uint
	if err := s.db.Model(&model.RbacRoleMenu{}).Where("role_id = ?", roleID).Pluck("menu_id", &menuIDs).Error; err != nil {
		return nil, err
	}
	return menuIDs, nil
}
