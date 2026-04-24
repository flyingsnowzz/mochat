package plugin

import (
	"gorm.io/gorm"
	"mochat-api-server/internal/model"
)

type RoomMessageBatchSendService struct {
	db *gorm.DB
}

func NewRoomMessageBatchSendService(db *gorm.DB) *RoomMessageBatchSendService {
	return &RoomMessageBatchSendService{db: db}
}

func (s *RoomMessageBatchSendService) List(corpID uint, offset, limit int) ([]model.RoomMessageBatchSend, int64, error) {
	var items []model.RoomMessageBatchSend
	var total int64
	query := s.db.Model(&model.RoomMessageBatchSend{}).Where("corp_id = ?", corpID)
	query.Count(&total)
	if err := query.Offset(offset).Limit(limit).Order("id DESC").Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (s *RoomMessageBatchSendService) GetByID(id uint) (*model.RoomMessageBatchSend, error) {
	var item model.RoomMessageBatchSend
	if err := s.db.First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (s *RoomMessageBatchSendService) Create(item *model.RoomMessageBatchSend) error {
	return s.db.Create(item).Error
}

func (s *RoomMessageBatchSendService) Update(id uint, updates map[string]interface{}) error {
	return s.db.Model(&model.RoomMessageBatchSend{}).Where("id = ?", id).Updates(updates).Error
}