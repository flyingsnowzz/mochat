package service

import (
	"mochat-api-server/internal/model"
	"gorm.io/gorm"
)

type WorkRoomService struct {
	db *gorm.DB
}

func NewWorkRoomService(db *gorm.DB) *WorkRoomService {
	return &WorkRoomService{db: db}
}

func (s *WorkRoomService) GetByID(id uint) (*model.WorkRoom, error) {
	var room model.WorkRoom
	if err := s.db.First(&room, id).Error; err != nil {
		return nil, err
	}
	return &room, nil
}

type WorkContactRoomService struct {
	db *gorm.DB
}

func NewWorkContactRoomService(db *gorm.DB) *WorkContactRoomService {
	return &WorkContactRoomService{db: db}
}

func (s *WorkContactRoomService) List(roomID uint, page, pageSize int) ([]model.WorkContactRoom, int64, error) {
	var contacts []model.WorkContactRoom
	var total int64
	query := s.db.Model(&model.WorkContactRoom{}).Where("room_id = ?", roomID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&contacts).Error; err != nil {
		return nil, 0, err
	}
	return contacts, total, nil
}

func (s *WorkContactRoomService) GetByID(id uint) (*model.WorkContactRoom, error) {
	var contact model.WorkContactRoom
	if err := s.db.First(&contact, id).Error; err != nil {
		return nil, err
	}
	return &contact, nil
}

func (s *WorkContactRoomService) GetByContactAndRoom(contactID, roomID uint) (*model.WorkContactRoom, error) {
	var contact model.WorkContactRoom
	if err := s.db.Where("contact_id = ? AND room_id = ?", contactID, roomID).First(&contact).Error; err != nil {
		return nil, err
	}
	return &contact, nil
}

func (s *WorkContactRoomService) Create(contact *model.WorkContactRoom) error {
	return s.db.Create(contact).Error
}

func (s *WorkContactRoomService) Update(contact *model.WorkContactRoom) error {
	return s.db.Save(contact).Error
}

func (s *WorkContactRoomService) Delete(id uint) error {
	return s.db.Delete(&model.WorkContactRoom{}, id).Error
}

func (s *WorkRoomService) GetByWxChatID(corpID uint, wxChatID string) (*model.WorkRoom, error) {
	var room model.WorkRoom
	if err := s.db.Where("corp_id = ? AND wx_chat_id = ?", corpID, wxChatID).First(&room).Error; err != nil {
		return nil, err
	}
	return &room, nil
}

func (s *WorkRoomService) List(corpID uint, page, pageSize int, filters map[string]interface{}) ([]model.WorkRoom, int64, error) {
	var rooms []model.WorkRoom
	var total int64
	query := s.db.Model(&model.WorkRoom{}).Where("corp_id = ?", corpID)
	for k, v := range filters {
		query = query.Where(k+" = ?", v)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&rooms).Error; err != nil {
		return nil, 0, err
	}
	return rooms, total, nil
}

func (s *WorkRoomService) Create(room *model.WorkRoom) error {
	return s.db.Create(room).Error
}

func (s *WorkRoomService) Update(id uint, updates map[string]interface{}) error {
	return s.db.Model(&model.WorkRoom{}).Where("id = ?", id).Updates(updates).Error
}

func (s *WorkRoomService) Delete(id uint) error {
	return s.db.Delete(&model.WorkRoom{}, id).Error
}

func (s *WorkRoomService) BatchUpdateGroup(ids []uint, groupID uint) error {
	return s.db.Model(&model.WorkRoom{}).Where("id IN ?", ids).Update("room_group_id", groupID).Error
}

type WorkRoomGroupService struct {
	db *gorm.DB
}

func NewWorkRoomGroupService(db *gorm.DB) *WorkRoomGroupService {
	return &WorkRoomGroupService{db: db}
}

func (s *WorkRoomGroupService) List(corpID uint) ([]model.WorkRoomGroup, error) {
	var groups []model.WorkRoomGroup
	if err := s.db.Where("corp_id = ?", corpID).Find(&groups).Error; err != nil {
		return nil, err
	}
	return groups, nil
}

func (s *WorkRoomGroupService) Create(group *model.WorkRoomGroup) error {
	return s.db.Create(group).Error
}

func (s *WorkRoomGroupService) Update(group *model.WorkRoomGroup) error {
	return s.db.Save(group).Error
}

func (s *WorkRoomGroupService) Delete(id uint) error {
	return s.db.Delete(&model.WorkRoomGroup{}, id).Error
}
