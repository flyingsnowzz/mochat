package plugin

import (
	"gorm.io/gorm"
	"mochat-api-server/internal/model"
)

type WorkFissionService struct {
	db *gorm.DB
}

func NewWorkFissionService(db *gorm.DB) *WorkFissionService {
	return &WorkFissionService{db: db}
}

func (s *WorkFissionService) List(corpID uint, offset, limit int) ([]model.WorkFission, int64, error) {
	var items []model.WorkFission
	var total int64
	query := s.db.Model(&model.WorkFission{}).Where("corp_id = ?", corpID)
	query.Count(&total)
	if err := query.Offset(offset).Limit(limit).Order("id DESC").Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (s *WorkFissionService) GetByID(id uint) (*model.WorkFission, error) {
	var item model.WorkFission
	if err := s.db.First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (s *WorkFissionService) Create(item *model.WorkFission) error {
	return s.db.Create(item).Error
}

func (s *WorkFissionService) Update(id uint, updates map[string]interface{}) error {
	return s.db.Model(&model.WorkFission{}).Where("id = ?", id).Updates(updates).Error
}