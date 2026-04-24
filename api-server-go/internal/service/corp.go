package service

import (
	"gorm.io/gorm"
	"mochat-api-server/internal/model"
)

type CorpService struct {
	db *gorm.DB
}

func NewCorpService(db *gorm.DB) *CorpService {
	return &CorpService{db: db}
}

func (s *CorpService) List(tenantID uint, corpName string, page, pageSize int) ([]model.Corp, int64, error) {
	var corps []model.Corp
	var total int64
	query := s.db.Model(&model.Corp{})
	if tenantID > 0 {
		query = query.Where("tenant_id = ?", tenantID)
	}
	if corpName != "" {
		query = query.Where("name LIKE ?", "%"+corpName+"%")
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * pageSize
	if err := query.Order("id DESC").Offset(offset).Limit(pageSize).Find(&corps).Error; err != nil {
		return nil, 0, err
	}
	return corps, total, nil
}

func (s *CorpService) GetByID(id uint) (*model.Corp, error) {
	var corp model.Corp
	if err := s.db.First(&corp, id).Error; err != nil {
		return nil, err
	}
	return &corp, nil
}

func (s *CorpService) Create(corp *model.Corp) error {
	return s.db.Create(corp).Error
}

func (s *CorpService) Update(corp *model.Corp) error {
	return s.db.Save(corp).Error
}

func (s *CorpService) Select(tenantID uint) ([]model.Corp, error) {
	var corps []model.Corp
	query := s.db.Model(&model.Corp{})
	if tenantID > 0 {
		query = query.Where("tenant_id = ?", tenantID)
	}
	if err := query.Find(&corps).Error; err != nil {
		return nil, err
	}
	return corps, nil
}
