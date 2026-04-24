package plugin

import (
	"gorm.io/gorm"
	"mochat-api-server/internal/model"
)

type RoomTagPullService struct {
	db *gorm.DB
}

func NewRoomTagPullService(db *gorm.DB) *RoomTagPullService {
	return &RoomTagPullService{db: db}
}

func (s *RoomTagPullService) List(corpID uint, offset, limit int) ([]model.RoomTagPull, int64, error) {
	var items []model.RoomTagPull
	var total int64
	query := s.db.Model(&model.RoomTagPull{}).Where("corp_id = ?", corpID)
	query.Count(&total)
	if err := query.Offset(offset).Limit(limit).Order("id DESC").Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (s *RoomTagPullService) GetByID(id uint) (*model.RoomTagPull, error) {
	var item model.RoomTagPull
	if err := s.db.First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (s *RoomTagPullService) Create(item *model.RoomTagPull) error {
	return s.db.Create(item).Error
}