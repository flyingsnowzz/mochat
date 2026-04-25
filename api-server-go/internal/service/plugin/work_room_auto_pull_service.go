package plugin

import (
	"gorm.io/gorm"
	"mochat-api-server/internal/model"
)

type WorkRoomAutoPullService struct {
	db *gorm.DB
}

func NewWorkRoomAutoPullService(db *gorm.DB) *WorkRoomAutoPullService {
	return &WorkRoomAutoPullService{db: db}
}

func (s *WorkRoomAutoPullService) List(corpID uint, qrcodeName string, offset, limit int) ([]model.WorkRoomAutoPull, int64, error) {
	var items []model.WorkRoomAutoPull
	var total int64
	query := s.db.Model(&model.WorkRoomAutoPull{}).Where("corp_id = ?", corpID)
	
	if qrcodeName != "" {
		query = query.Where("qrcode_name LIKE ?", "%"+qrcodeName+"%")
	}
	
	query.Count(&total)
	if err := query.Offset(offset).Limit(limit).Order("id DESC").Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (s *WorkRoomAutoPullService) GetByID(id uint) (*model.WorkRoomAutoPull, error) {
	var item model.WorkRoomAutoPull
	if err := s.db.First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (s *WorkRoomAutoPullService) Create(item *model.WorkRoomAutoPull) error {
	return s.db.Create(item).Error
}

func (s *WorkRoomAutoPullService) Update(id uint, updates map[string]interface{}) error {
	return s.db.Model(&model.WorkRoomAutoPull{}).Where("id = ?", id).Updates(updates).Error
}

func (s *WorkRoomAutoPullService) Delete(id uint) error {
	return s.db.Delete(&model.WorkRoomAutoPull{}, id).Error
}
