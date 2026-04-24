package service

import (
	"gorm.io/gorm"
	"mochat-api-server/internal/model"
)

// MediumGroupServiceAlt 素材分组服务(别名，用于兼容)
type MediumGroupServiceAlt struct {
	db *gorm.DB
}

func NewMediumGroupServiceAlt(db *gorm.DB) *MediumGroupServiceAlt {
	return &MediumGroupServiceAlt{db: db}
}

func (s *MediumGroupServiceAlt) List(corpID uint) ([]model.MediumGroup, error) {
	var groups []model.MediumGroup
	if err := s.db.Where("corp_id = ?", corpID).Find(&groups).Error; err != nil {
		return nil, err
	}
	return groups, nil
}

func (s *MediumGroupServiceAlt) Create(group *model.MediumGroup) error {
	return s.db.Create(group).Error
}

func (s *MediumGroupServiceAlt) Update(group *model.MediumGroup) error {
	return s.db.Save(group).Error
}

func (s *MediumGroupServiceAlt) Delete(id uint) error {
	return s.db.Delete(&model.MediumGroup{}, id).Error
}
