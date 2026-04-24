package business

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

func (s *CorpService) GetByID(id uint) (*model.Corp, error) {
	var corp model.Corp
	if err := s.db.First(&corp, id).Error; err != nil {
		return nil, err
	}
	return &corp, nil
}

func (s *CorpService) GetByWxCorpid(wxCorpid string) (*model.Corp, error) {
	var corp model.Corp
	if err := s.db.Where("wx_corpid = ?", wxCorpid).First(&corp).Error; err != nil {
		return nil, err
	}
	return &corp, nil
}

func (s *CorpService) List(tenantID uint, offset, limit int) ([]model.Corp, int64, error) {
	var corps []model.Corp
	var total int64
	query := s.db.Model(&model.Corp{})
	if tenantID > 0 {
		query = query.Where("tenant_id = ?", tenantID)
	}
	query.Count(&total)
	if err := query.Offset(offset).Limit(limit).Order("id DESC").Find(&corps).Error; err != nil {
		return nil, 0, err
	}
	return corps, total, nil
}

func (s *CorpService) Create(corp *model.Corp) error {
	return s.db.Create(corp).Error
}

func (s *CorpService) Update(id uint, updates map[string]interface{}) error {
	return s.db.Model(&model.Corp{}).Where("id = ?", id).Updates(updates).Error
}

func (s *CorpService) Delete(id uint) error {
	return s.db.Delete(&model.Corp{}, id).Error
}