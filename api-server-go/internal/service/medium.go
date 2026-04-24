package service

import (
	"encoding/json"
	"fmt"

	"gorm.io/gorm"
	"mochat-api-server/internal/model"
)

type MediumService struct {
	db *gorm.DB
}

func NewMediumService(db *gorm.DB) *MediumService {
	return &MediumService{db: db}
}

func (s *MediumService) List(corpID uint, page, pageSize int, mediumType int, groupID uint) ([]model.Medium, int64, error) {
	var media []model.Medium
	var total int64
	query := s.db.Model(&model.Medium{}).Where("corp_id = ?", corpID)
	if mediumType > 0 {
		query = query.Where("type = ?", mediumType)
	}
	if groupID > 0 {
		query = query.Where("medium_group_id = ?", groupID)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&media).Error; err != nil {
		return nil, 0, err
	}
	return media, total, nil
}

func (s *MediumService) ListWithSearch(corpID uint, page, pageSize int, mediumType int, groupID uint, searchStr string) ([]model.Medium, int64, error) {
	var media []model.Medium
	var total int64
	query := s.db.Model(&model.Medium{}).Where("corp_id = ?", corpID)
	if mediumType > 0 {
		query = query.Where("type = ?", mediumType)
	}
	if groupID > 0 {
		query = query.Where("medium_group_id = ?", groupID)
	}
	if searchStr != "" {
		query = query.Where("content LIKE ?", "%"+searchStr+"%")
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&media).Error; err != nil {
		return nil, 0, err
	}
	return media, total, nil
}

func (s *MediumService) GetByID(id uint) (*model.Medium, error) {
	var medium model.Medium
	if err := s.db.First(&medium, id).Error; err != nil {
		return nil, err
	}
	return &medium, nil
}

func (s *MediumService) Create(medium *model.Medium) error {
	return s.db.Create(medium).Error
}

func (s *MediumService) Update(medium *model.Medium) error {
	return s.db.Save(medium).Error
}

func (s *MediumService) Delete(id uint) error {
	return s.db.Delete(&model.Medium{}, id).Error
}

func ParseMediumContent(raw string) map[string]interface{} {
	if raw == "" {
		return map[string]interface{}{}
	}
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(raw), &result); err != nil {
		return map[string]interface{}{}
	}
	return result
}

func MarshalMediumContent(content interface{}) (string, error) {
	if content == nil {
		return "{}", nil
	}
	data, err := json.Marshal(content)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func MediumTypeName(mediumType int) string {
	switch mediumType {
	case 1:
		return "文本"
	case 2:
		return "图片"
	case 3:
		return "图文"
	case 4:
		return "音频"
	case 5:
		return "视频"
	case 6:
		return "小程序"
	case 7:
		return "文件"
	default:
		return fmt.Sprintf("%d", mediumType)
	}
}

type MediumGroupService struct {
	db *gorm.DB
}

func NewMediumGroupService(db *gorm.DB) *MediumGroupService {
	return &MediumGroupService{db: db}
}

func (s *MediumGroupService) List(corpID uint) ([]model.MediumGroup, error) {
	var groups []model.MediumGroup
	if err := s.db.Where("corp_id = ?", corpID).Order("`order` ASC").Find(&groups).Error; err != nil {
		return nil, err
	}
	return groups, nil
}

func (s *MediumGroupService) Create(group *model.MediumGroup) error {
	return s.db.Create(group).Error
}

func (s *MediumGroupService) Update(group *model.MediumGroup) error {
	return s.db.Save(group).Error
}

func (s *MediumGroupService) Delete(id uint) error {
	return s.db.Delete(&model.MediumGroup{}, id).Error
}
