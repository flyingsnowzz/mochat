package plugin

import (
	"gorm.io/gorm"
	"mochat-api-server/internal/model"
)

type RoomWelcomeService struct {
	db *gorm.DB
}

func NewRoomWelcomeService(db *gorm.DB) *RoomWelcomeService {
	return &RoomWelcomeService{db: db}
}

func (s *RoomWelcomeService) List(corpID uint) ([]model.RoomWelcomeTemplate, error) {
	var items []model.RoomWelcomeTemplate
	if err := s.db.Where("corp_id = ?", corpID).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (s *RoomWelcomeService) GetByID(id uint) (*model.RoomWelcomeTemplate, error) {
	var item model.RoomWelcomeTemplate
	if err := s.db.First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (s *RoomWelcomeService) Create(item *model.RoomWelcomeTemplate) error {
	return s.db.Create(item).Error
}

func (s *RoomWelcomeService) Update(id uint, updates map[string]interface{}) error {
	return s.db.Model(&model.RoomWelcomeTemplate{}).Where("id = ?", id).Updates(updates).Error
}