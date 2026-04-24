package service

import (
	"gorm.io/gorm"
	"mochat-api-server/internal/model"
)

type RbacMenuService struct {
	db *gorm.DB
}

func NewRbacMenuService(db *gorm.DB) *RbacMenuService {
	return &RbacMenuService{db: db}
}

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

func (s *RbacMenuService) GetByID(id uint) (*model.RbacMenu, error) {
	var menu model.RbacMenu
	if err := s.db.First(&menu, id).Error; err != nil {
		return nil, err
	}
	return &menu, nil
}

func (s *RbacMenuService) Create(menu *model.RbacMenu) error {
	return s.db.Create(menu).Error
}

func (s *RbacMenuService) Update(menu *model.RbacMenu) error {
	return s.db.Save(menu).Error
}

func (s *RbacMenuService) Delete(id uint) error {
	return s.db.Delete(&model.RbacMenu{}, id).Error
}

func (s *RbacMenuService) Select() ([]model.RbacMenu, error) {
	var menus []model.RbacMenu
	if err := s.db.Where("status = 1 AND is_page_menu = 1").Order("sort ASC").Find(&menus).Error; err != nil {
		return nil, err
	}
	return menus, nil
}

type RbacRoleService struct {
	db *gorm.DB
}

func NewRbacRoleService(db *gorm.DB) *RbacRoleService {
	return &RbacRoleService{db: db}
}

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

func (s *RbacRoleService) GetByID(id uint) (*model.RbacRole, error) {
	var role model.RbacRole
	if err := s.db.First(&role, id).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (s *RbacRoleService) Create(role *model.RbacRole) error {
	return s.db.Create(role).Error
}

func (s *RbacRoleService) Update(role *model.RbacRole) error {
	return s.db.Save(role).Error
}

func (s *RbacRoleService) Delete(id uint) error {
	return s.db.Delete(&model.RbacRole{}, id).Error
}

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

func (s *RbacRoleService) UpdateByID(id uint, updates map[string]interface{}) error {
	return s.db.Model(&model.RbacRole{}).Where("id = ?", id).Updates(updates).Error
}

func (s *RbacRoleService) DeleteWithRelations(id uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		tx.Where("role_id = ?", id).Delete(&model.RbacRoleMenu{})
		tx.Where("role_id = ?", id).Delete(&model.RbacUserRole{})
		return tx.Delete(&model.RbacRole{}, id).Error
	})
}

func (s *RbacRoleService) SetMenus(roleID uint, menuIDs []uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		tx.Where("role_id = ?", roleID).Delete(&model.RbacRoleMenu{})
		for _, menuID := range menuIDs {
			if err := tx.Create(&model.RbacRoleMenu{RoleID: roleID, MenuID: menuID}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *RbacRoleService) GetMenus(roleID uint) ([]uint, error) {
	var menuIDs []uint
	if err := s.db.Model(&model.RbacRoleMenu{}).Where("role_id = ?", roleID).Pluck("menu_id", &menuIDs).Error; err != nil {
		return nil, err
	}
	return menuIDs, nil
}

func (s *RbacMenuService) ListByIDs(ids []uint) ([]model.RbacMenu, error) {
	var menus []model.RbacMenu
	if err := s.db.Where("id IN ? AND status = 1", ids).Find(&menus).Error; err != nil {
		return nil, err
	}
	return menus, nil
}

func (s *RbacMenuService) UpdateByID(id uint, updates map[string]interface{}) error {
	return s.db.Model(&model.RbacMenu{}).Where("id = ?", id).Updates(updates).Error
}
