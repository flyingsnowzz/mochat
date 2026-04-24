package plugin

import (
	"gorm.io/gorm"
	"mochat-api-server/internal/model"
)

type ContactMessageBatchSendService struct {
	db *gorm.DB
}

func NewContactMessageBatchSendService(db *gorm.DB) *ContactMessageBatchSendService {
	return &ContactMessageBatchSendService{db: db}
}

func (s *ContactMessageBatchSendService) List(corpID uint, offset, limit int) ([]model.ContactMessageBatchSend, int64, error) {
	var items []model.ContactMessageBatchSend
	var total int64
	query := s.db.Model(&model.ContactMessageBatchSend{}).Where("corp_id = ?", corpID)
	query.Count(&total)
	if err := query.Offset(offset).Limit(limit).Order("id DESC").Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (s *ContactMessageBatchSendService) GetByID(id uint) (*model.ContactMessageBatchSend, error) {
	var item model.ContactMessageBatchSend
	if err := s.db.First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (s *ContactMessageBatchSendService) Create(item *model.ContactMessageBatchSend) error {
	return s.db.Create(item).Error
}

func (s *ContactMessageBatchSendService) Update(id uint, updates map[string]interface{}) error {
	return s.db.Model(&model.ContactMessageBatchSend{}).Where("id = ?", id).Updates(updates).Error
}