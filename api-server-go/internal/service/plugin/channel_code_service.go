package plugin

import (
	"gorm.io/gorm"
	"mochat-api-server/internal/model"
)

type ChannelCodeService struct {
	db *gorm.DB
}

func NewChannelCodeService(db *gorm.DB) *ChannelCodeService {
	return &ChannelCodeService{db: db}
}

func (s *ChannelCodeService) List(corpID uint, groupID uint, offset, limit int) ([]model.ChannelCode, int64, error) {
	var codes []model.ChannelCode
	var total int64
	query := s.db.Model(&model.ChannelCode{}).Where("corp_id = ?", corpID)
	if groupID > 0 {
		query = query.Where("group_id = ?", groupID)
	}
	query.Count(&total)
	if err := query.Offset(offset).Limit(limit).Order("id DESC").Find(&codes).Error; err != nil {
		return nil, 0, err
	}
	return codes, total, nil
}

func (s *ChannelCodeService) GetByID(id uint) (*model.ChannelCode, error) {
	var code model.ChannelCode
	if err := s.db.First(&code, id).Error; err != nil {
		return nil, err
	}
	return &code, nil
}

func (s *ChannelCodeService) Create(code *model.ChannelCode) error {
	return s.db.Create(code).Error
}

func (s *ChannelCodeService) Update(id uint, updates map[string]interface{}) error {
	return s.db.Model(&model.ChannelCode{}).Where("id = ?", id).Updates(updates).Error
}

func (s *ChannelCodeService) Delete(id uint) error {
	return s.db.Delete(&model.ChannelCode{}, id).Error
}

type ChannelCodeGroupService struct {
	db *gorm.DB
}

func NewChannelCodeGroupService(db *gorm.DB) *ChannelCodeGroupService {
	return &ChannelCodeGroupService{db: db}
}

func (s *ChannelCodeGroupService) List(corpID uint) ([]model.ChannelCodeGroup, error) {
	var groups []model.ChannelCodeGroup
	if err := s.db.Where("corp_id = ?", corpID).Find(&groups).Error; err != nil {
		return nil, err
	}
	return groups, nil
}

func (s *ChannelCodeGroupService) Create(group *model.ChannelCodeGroup) error {
	return s.db.Create(group).Error
}

func (s *ChannelCodeGroupService) Update(id uint, updates map[string]interface{}) error {
	return s.db.Model(&model.ChannelCodeGroup{}).Where("id = ?", id).Updates(updates).Error
}

func (s *ChannelCodeGroupService) Delete(id uint) error {
	return s.db.Delete(&model.ChannelCodeGroup{}, id).Error
}