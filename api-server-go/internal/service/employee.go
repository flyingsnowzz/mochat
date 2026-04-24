package service

import (
	"gorm.io/gorm"
	"mochat-api-server/internal/model"
)

type WorkEmployeeService struct {
	db *gorm.DB
}

type WorkEmployeeListFilter struct {
	Name        string
	Status      int
	ContactAuth string
}

func NewWorkEmployeeService(db *gorm.DB) *WorkEmployeeService {
	return &WorkEmployeeService{db: db}
}

func (s *WorkEmployeeService) GetByID(id uint) (*model.WorkEmployee, error) {
	var emp model.WorkEmployee
	if err := s.db.First(&emp, id).Error; err != nil {
		return nil, err
	}
	return &emp, nil
}

func (s *WorkEmployeeService) GetByWxUserID(corpID uint, wxUserID string) (*model.WorkEmployee, error) {
	var emp model.WorkEmployee
	if err := s.db.Where("corp_id = ? AND wx_user_id = ?", corpID, wxUserID).First(&emp).Error; err != nil {
		return nil, err
	}
	return &emp, nil
}

func (s *WorkEmployeeService) List(corpID uint, filter WorkEmployeeListFilter, page, pageSize int) ([]model.WorkEmployee, int64, error) {
	var employees []model.WorkEmployee
	var total int64
	query := s.db.Model(&model.WorkEmployee{}).Where("corp_id = ?", corpID)
	if filter.Name != "" {
		query = query.Where("name LIKE ?", "%"+filter.Name+"%")
	}
	if filter.Status > 0 {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.ContactAuth != "" && filter.ContactAuth != "all" {
		query = query.Where("contact_auth = ?", filter.ContactAuth)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * pageSize
	if err := query.Order("updated_at DESC, id DESC").Offset(offset).Limit(pageSize).Find(&employees).Error; err != nil {
		return nil, 0, err
	}
	return employees, total, nil
}

func (s *WorkEmployeeService) LastSyncTime(corpID uint) (string, error) {
	if corpID == 0 {
		return "", nil
	}

	var item model.WorkUpdateTime
	err := s.db.Where("corp_id = ? AND type = ?", corpID, workUpdateTimeTypeEmployee).
		Order("id DESC").
		First(&item).Error
	if err == nil {
		return item.LastUpdateTime, nil
	}
	if err == gorm.ErrRecordNotFound {
		return "", nil
	}
	return "", err
}

func (s *WorkEmployeeService) Create(emp *model.WorkEmployee) error {
	return s.db.Create(emp).Error
}

func (s *WorkEmployeeService) Update(id uint, updates map[string]interface{}) error {
	return s.db.Model(&model.WorkEmployee{}).Where("id = ?", id).Updates(updates).Error
}

func (s *WorkEmployeeService) Delete(id uint) error {
	return s.db.Delete(&model.WorkEmployee{}, id).Error
}

type WorkDepartmentService struct {
	db *gorm.DB
}

func NewWorkDepartmentService(db *gorm.DB) *WorkDepartmentService {
	return &WorkDepartmentService{db: db}
}

func (s *WorkDepartmentService) List(corpID uint) ([]model.WorkDepartment, error) {
	var departments []model.WorkDepartment
	if err := s.db.Where("corp_id = ?", corpID).Order("`order` ASC").Find(&departments).Error; err != nil {
		return nil, err
	}
	return departments, nil
}

func (s *WorkDepartmentService) PageList(corpID uint, page, pageSize int) ([]model.WorkDepartment, int64, error) {
	var departments []model.WorkDepartment
	var total int64
	query := s.db.Model(&model.WorkDepartment{}).Where("corp_id = ?", corpID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Find(&departments).Error; err != nil {
		return nil, 0, err
	}
	return departments, total, nil
}

func (s *WorkDepartmentService) GetByID(id uint) (*model.WorkDepartment, error) {
	var dept model.WorkDepartment
	if err := s.db.First(&dept, id).Error; err != nil {
		return nil, err
	}
	return &dept, nil
}

func (s *WorkDepartmentService) Create(dept *model.WorkDepartment) error {
	return s.db.Create(dept).Error
}

func (s *WorkDepartmentService) Update(id uint, updates map[string]interface{}) error {
	return s.db.Model(&model.WorkDepartment{}).Where("id = ?", id).Updates(updates).Error
}

func (s *WorkDepartmentService) Delete(id uint) error {
	return s.db.Delete(&model.WorkDepartment{}, id).Error
}
