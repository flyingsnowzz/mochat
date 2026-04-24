package service

import (
	"mochat-api-server/internal/model"
	"gorm.io/gorm"
)

type WorkContactService struct {
	db *gorm.DB
}

func NewWorkContactService(db *gorm.DB) *WorkContactService {
	return &WorkContactService{db: db}
}

func (s *WorkContactService) List(corpID uint, page, pageSize int, filters map[string]interface{}) ([]model.WorkContact, int64, error) {
	var contacts []model.WorkContact
	var total int64
	query := s.db.Model(&model.WorkContact{}).Where("corp_id = ?", corpID)
	for k, v := range filters {
		query = query.Where(k+" = ?", v)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&contacts).Error; err != nil {
		return nil, 0, err
	}
	return contacts, total, nil
}

func (s *WorkContactService) GetByID(id uint) (*model.WorkContact, error) {
	var contact model.WorkContact
	if err := s.db.First(&contact, id).Error; err != nil {
		return nil, err
	}
	return &contact, nil
}

func (s *WorkContactService) Update(contact *model.WorkContact) error {
	return s.db.Save(contact).Error
}

func (s *WorkContactService) GetByWxExternalUserID(corpID uint, wxID string) (*model.WorkContact, error) {
	var contact model.WorkContact
	if err := s.db.Where("corp_id = ? AND wx_external_userid = ?", corpID, wxID).First(&contact).Error; err != nil {
		return nil, err
	}
	return &contact, nil
}

func (s *WorkContactService) Create(contact *model.WorkContact) error {
	return s.db.Create(contact).Error
}

func (s *WorkContactService) UpdateByID(id uint, updates map[string]interface{}) error {
	return s.db.Model(&model.WorkContact{}).Where("id = ?", id).Updates(updates).Error
}

func (s *WorkContactService) Delete(id uint) error {
	return s.db.Delete(&model.WorkContact{}, id).Error
}

func (s *WorkContactService) ListByOffset(corpID uint, offset, limit int, filters map[string]interface{}) ([]model.WorkContact, int64, error) {
	var contacts []model.WorkContact
	var total int64
	query := s.db.Model(&model.WorkContact{}).Where("corp_id = ?", corpID)
	query.Count(&total)
	if err := query.Offset(offset).Limit(limit).Order("id DESC").Find(&contacts).Error; err != nil {
		return nil, 0, err
	}
	return contacts, total, nil
}

type WorkContactTagService struct {
	db *gorm.DB
}

func NewWorkContactTagService(db *gorm.DB) *WorkContactTagService {
	return &WorkContactTagService{db: db}
}

func (s *WorkContactTagService) List(corpID uint) ([]model.WorkContactTag, error) {
	var tags []model.WorkContactTag
	if err := s.db.Where("corp_id = ?", corpID).Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}

func (s *WorkContactTagService) GetByID(id uint) (*model.WorkContactTag, error) {
	var tag model.WorkContactTag
	if err := s.db.First(&tag, id).Error; err != nil {
		return nil, err
	}
	return &tag, nil
}

func (s *WorkContactTagService) Create(tag *model.WorkContactTag) error {
	return s.db.Create(tag).Error
}

func (s *WorkContactTagService) Update(tag *model.WorkContactTag) error {
	return s.db.Save(tag).Error
}

func (s *WorkContactTagService) Delete(id uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		tx.Where("contact_tag_id = ?", id).Delete(&model.WorkContactTagPivot{})
		return tx.Delete(&model.WorkContactTag{}, id).Error
	})
}

func (s *WorkContactTagService) UpdateByID(id uint, updates map[string]interface{}) error {
	return s.db.Model(&model.WorkContactTag{}).Where("id = ?", id).Updates(updates).Error
}

func (s *WorkContactTagService) ListByOrder(corpID uint) ([]model.WorkContactTag, error) {
	var tags []model.WorkContactTag
	if err := s.db.Where("corp_id = ?", corpID).Order("id ASC").Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}

type ContactFieldService struct {
	db *gorm.DB
}

func NewContactFieldService(db *gorm.DB) *ContactFieldService {
	return &ContactFieldService{db: db}
}

func (s *ContactFieldService) List() ([]model.ContactField, error) {
	var fields []model.ContactField
	if err := s.db.Order("`order` ASC, id ASC").Find(&fields).Error; err != nil {
		return nil, err
	}
	return fields, nil
}

func (s *ContactFieldService) GetByID(id uint) (*model.ContactField, error) {
	var field model.ContactField
	if err := s.db.First(&field, id).Error; err != nil {
		return nil, err
	}
	return &field, nil
}

func (s *ContactFieldService) Create(field *model.ContactField) error {
	return s.db.Create(field).Error
}

func (s *ContactFieldService) Update(field *model.ContactField) error {
	return s.db.Save(field).Error
}

func (s *ContactFieldService) Delete(id uint) error {
	return s.db.Delete(&model.ContactField{}, id).Error
}

func (s *ContactFieldService) ListByStatus(status int, page, pageSize int) ([]model.ContactField, int64, error) {
	var fields []model.ContactField
	var total int64
	query := s.db.Model(&model.ContactField{})
	if status != 2 {
		query = query.Where("status = ?", status)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("`order` DESC, id ASC").Find(&fields).Error; err != nil {
		return nil, 0, err
	}
	return fields, total, nil
}

func (s *ContactFieldService) UpdateStatus(id uint, status int) error {
	return s.db.Model(&model.ContactField{}).Where("id = ?", id).Update("status", status).Error
}

func (s *ContactFieldService) UpdateOrder(id uint, order int) error {
	return s.db.Model(&model.ContactField{}).Where("id = ?", id).Update("order", order).Error
}

func (s *ContactFieldService) DB() *gorm.DB {
	return s.db
}

type ContactFieldPivotService struct {
	db *gorm.DB
}

func NewContactFieldPivotService(db *gorm.DB) *ContactFieldPivotService {
	return &ContactFieldPivotService{db: db}
}

func (s *ContactFieldPivotService) List(contactID uint) ([]model.ContactFieldPivot, error) {
	var pivots []model.ContactFieldPivot
	if err := s.db.Where("contact_id = ?", contactID).Find(&pivots).Error; err != nil {
		return nil, err
	}
	return pivots, nil
}

func (s *ContactFieldPivotService) Update(contactID, fieldID uint, value string) error {
	var pivot model.ContactFieldPivot
	result := s.db.Where("contact_id = ? AND contact_field_id = ?", contactID, fieldID).First(&pivot)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			// 创建新记录
			pivot = model.ContactFieldPivot{
				ContactID:      contactID,
				ContactFieldID: fieldID,
				Value:          value,
			}
			return s.db.Create(&pivot).Error
		}
		return result.Error
	}
	// 更新现有记录
	pivot.Value = value
	return s.db.Save(&pivot).Error
}

func (s *ContactFieldPivotService) BatchUpdate(contactID uint, fields []struct {
	FieldID uint   `json:"fieldId"`
	Value   string `json:"value"`
}) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		for _, field := range fields {
			var pivot model.ContactFieldPivot
			result := tx.Where("contact_id = ? AND contact_field_id = ?", contactID, field.FieldID).First(&pivot)
			if result.Error != nil {
				if result.Error == gorm.ErrRecordNotFound {
					// 创建新记录
					pivot = model.ContactFieldPivot{
						ContactID:      contactID,
						ContactFieldID: field.FieldID,
						Value:          field.Value,
					}
					if err := tx.Create(&pivot).Error; err != nil {
						return err
					}
				} else {
					return result.Error
				}
			} else {
				// 更新现有记录
				pivot.Value = field.Value
				if err := tx.Save(&pivot).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
}

type WorkContactTagGroupService struct {
	db *gorm.DB
}

func NewWorkContactTagGroupService(db *gorm.DB) *WorkContactTagGroupService {
	return &WorkContactTagGroupService{db: db}
}

func (s *WorkContactTagGroupService) List(corpID uint) ([]model.WorkContactTagGroup, error) {
	var groups []model.WorkContactTagGroup
	if err := s.db.Where("corp_id = ?", corpID).Order("`order` ASC, id ASC").Find(&groups).Error; err != nil {
		return nil, err
	}
	return groups, nil
}

func (s *WorkContactTagGroupService) GetByID(id uint) (*model.WorkContactTagGroup, error) {
	var group model.WorkContactTagGroup
	if err := s.db.First(&group, id).Error; err != nil {
		return nil, err
	}
	return &group, nil
}

func (s *WorkContactTagGroupService) Create(group *model.WorkContactTagGroup) error {
	return s.db.Create(group).Error
}

func (s *WorkContactTagGroupService) Update(group *model.WorkContactTagGroup) error {
	return s.db.Save(group).Error
}

func (s *WorkContactTagGroupService) Delete(id uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 先删除关联的标签
		tx.Where("contact_tag_group_id = ?", id).Delete(&model.WorkContactTag{})
		// 删除标签组
		return tx.Delete(&model.WorkContactTagGroup{}, id).Error
	})
}

func (s *WorkContactTagGroupService) UpdateOrder(id uint, order int) error {
	return s.db.Model(&model.WorkContactTagGroup{}).Where("id = ?", id).Update("order", order).Error
}
