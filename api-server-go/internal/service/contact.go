package service

import (
	"mochat-api-server/internal/model"
	"gorm.io/gorm"
)

// WorkContactService 客户 Service
// 提供客户（WorkContact）的 CRUD 操作功能
// 主要职责：
// 1. 获取客户列表（分页）
// 2. 根据 ID 获取客户详情
// 3. 更新客户信息
// 4. 根据企业微信外部用户 ID 获取客户
// 5. 创建客户
// 6. 根据 ID 更新客户指定字段
// 7. 删除客户
// 8. 使用 offset 方式获取客户列表
//
// 依赖：
// - gorm.DB: 数据库连接

type WorkContactService struct {
	db *gorm.DB // 数据库连接
}

// NewWorkContactService 创建客户 Service 实例
// 参数：db - GORM 数据库连接
// 返回：客户 Service 实例
func NewWorkContactService(db *gorm.DB) *WorkContactService {
	return &WorkContactService{db: db}
}

// List 获取客户列表（分页）
// 查询指定企业的客户列表，支持分页和筛选
// 参数：
//
//	corpID - 企业 ID
//	page - 页码
//	pageSize - 每页数量
//	filters - 筛选条件映射
//
// 返回：客户列表、总数和错误信息
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

// GetByID 根据 ID 获取客户详情
// 查询指定 ID 的客户
// 参数：
//
//	id - 客户 ID
//
// 返回：客户实例和错误信息
func (s *WorkContactService) GetByID(id uint) (*model.WorkContact, error) {
	var contact model.WorkContact
	if err := s.db.First(&contact, id).Error; err != nil {
		return nil, err
	}
	return &contact, nil
}

// Update 更新客户信息
// 更新数据库中的客户信息
// 参数：
//
//	contact - 客户实例
//
// 返回：错误信息
func (s *WorkContactService) Update(contact *model.WorkContact) error {
	return s.db.Save(contact).Error
}

// GetByWxExternalUserID 根据企业微信外部用户 ID 获取客户
// 使用企业微信外部用户 ID 查询客户
// 参数：
//
//	corpID - 企业 ID
//	wxID - 企业微信外部用户 ID
//
// 返回：客户实例和错误信息
func (s *WorkContactService) GetByWxExternalUserID(corpID uint, wxID string) (*model.WorkContact, error) {
	var contact model.WorkContact
	if err := s.db.Where("corp_id = ? AND wx_external_userid = ?", corpID, wxID).First(&contact).Error; err != nil {
		return nil, err
	}
	return &contact, nil
}

// Create 创建客户
// 将客户信息保存到数据库
// 参数：
//
//	contact - 客户实例
//
// 返回：错误信息
func (s *WorkContactService) Create(contact *model.WorkContact) error {
	return s.db.Create(contact).Error
}

// UpdateByID 根据 ID 更新客户指定字段
// 更新客户的指定字段
// 参数：
//
//	id - 客户 ID
//	updates - 待更新的字段映射
//
// 返回：错误信息
func (s *WorkContactService) UpdateByID(id uint, updates map[string]interface{}) error {
	return s.db.Model(&model.WorkContact{}).Where("id = ?", id).Updates(updates).Error
}

// Delete 删除客户
// 从数据库中删除指定 ID 的客户
// 参数：
//
//	id - 客户 ID
//
// 返回：错误信息
func (s *WorkContactService) Delete(id uint) error {
	return s.db.Delete(&model.WorkContact{}, id).Error
}

// ListByOffset 使用 offset 方式获取客户列表
// 使用 offset 和 limit 方式查询客户列表，不支持筛选
// 参数：
//
//	corpID - 企业 ID
//	offset - 偏移量
//	limit - 限制数量
//	filters - 筛选条件映射（暂未使用）
//
// 返回：客户列表、总数和错误信息
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

// WorkContactTagService 客户标签 Service
// 提供客户标签的 CRUD 操作功能
// 主要职责：
// 1. 获取客户标签列表
// 2. 根据 ID 获取客户标签详情
// 3. 创建客户标签
// 4. 更新客户标签
// 5. 删除客户标签（同时删除关联的客户标签中间表记录）
// 6. 根据 ID 更新客户标签指定字段
// 7. 按顺序获取客户标签列表
//
// 依赖：
// - gorm.DB: 数据库连接

type WorkContactTagService struct {
	db *gorm.DB // 数据库连接
}

// NewWorkContactTagService 创建客户标签 Service 实例
// 参数：db - GORM 数据库连接
// 返回：客户标签 Service 实例
func NewWorkContactTagService(db *gorm.DB) *WorkContactTagService {
	return &WorkContactTagService{db: db}
}

// List 获取客户标签列表
// 查询指定企业的所有客户标签
// 参数：
//
//	corpID - 企业 ID
//
// 返回：客户标签列表和错误信息
func (s *WorkContactTagService) List(corpID uint) ([]model.WorkContactTag, error) {
	var tags []model.WorkContactTag
	if err := s.db.Where("corp_id = ?", corpID).Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}

// GetByID 根据 ID 获取客户标签详情
// 查询指定 ID 的客户标签
// 参数：
//
//	id - 客户标签 ID
//
// 返回：客户标签实例和错误信息
func (s *WorkContactTagService) GetByID(id uint) (*model.WorkContactTag, error) {
	var tag model.WorkContactTag
	if err := s.db.First(&tag, id).Error; err != nil {
		return nil, err
	}
	return &tag, nil
}

// Create 创建客户标签
// 将客户标签信息保存到数据库
// 参数：
//
//	tag - 客户标签实例
//
// 返回：错误信息
func (s *WorkContactTagService) Create(tag *model.WorkContactTag) error {
	return s.db.Create(tag).Error
}

// Update 更新客户标签
// 更新数据库中的客户标签信息
// 参数：
//
//	tag - 客户标签实例
//
// 返回：错误信息
func (s *WorkContactTagService) Update(tag *model.WorkContactTag) error {
	return s.db.Save(tag).Error
}

// Delete 删除客户标签
// 删除指定 ID 的客户标签，同时删除关联的客户标签中间表记录
// 使用事务确保数据一致性
// 参数：
//
//	id - 客户标签 ID
//
// 返回：错误信息
func (s *WorkContactTagService) Delete(id uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 删除关联的客户标签中间表记录
		tx.Where("contact_tag_id = ?", id).Delete(&model.WorkContactTagPivot{})
		// 删除客户标签
		return tx.Delete(&model.WorkContactTag{}, id).Error
	})
}

// UpdateByID 根据 ID 更新客户标签指定字段
// 更新客户标签的指定字段
// 参数：
//
//	id - 客户标签 ID
//	updates - 待更新的字段映射
//
// 返回：错误信息
func (s *WorkContactTagService) UpdateByID(id uint, updates map[string]interface{}) error {
	return s.db.Model(&model.WorkContactTag{}).Where("id = ?", id).Updates(updates).Error
}

// ListByOrder 按顺序获取客户标签列表
// 查询指定企业的所有客户标签，按 ID 升序排列
// 参数：
//
//	corpID - 企业 ID
//
// 返回：客户标签列表和错误信息
func (s *WorkContactTagService) ListByOrder(corpID uint) ([]model.WorkContactTag, error) {
	var tags []model.WorkContactTag
	if err := s.db.Where("corp_id = ?", corpID).Order("id ASC").Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}

// ContactFieldService 客户字段 Service
// 提供客户字段的 CRUD 操作功能
// 主要职责：
// 1. 获取客户字段列表
// 2. 根据 ID 获取客户字段详情
// 3. 创建客户字段
// 4. 更新客户字段
// 5. 删除客户字段
// 6. 按状态获取客户字段列表（分页）
// 7. 更新客户字段状态
// 8. 更新客户字段排序
// 9. 获取数据库连接
//
// 依赖：
// - gorm.DB: 数据库连接

type ContactFieldService struct {
	db *gorm.DB // 数据库连接
}

// NewContactFieldService 创建客户字段 Service 实例
// 参数：db - GORM 数据库连接
// 返回：客户字段 Service 实例
func NewContactFieldService(db *gorm.DB) *ContactFieldService {
	return &ContactFieldService{db: db}
}

// List 获取客户字段列表
// 查询所有客户字段，按排序字段和 ID 升序排列
// 返回：客户字段列表和错误信息
func (s *ContactFieldService) List() ([]model.ContactField, error) {
	var fields []model.ContactField
	if err := s.db.Order("`order` ASC, id ASC").Find(&fields).Error; err != nil {
		return nil, err
	}
	return fields, nil
}

// GetByID 根据 ID 获取客户字段详情
// 查询指定 ID 的客户字段
// 参数：
//
//	id - 客户字段 ID
//
// 返回：客户字段实例和错误信息
func (s *ContactFieldService) GetByID(id uint) (*model.ContactField, error) {
	var field model.ContactField
	if err := s.db.First(&field, id).Error; err != nil {
		return nil, err
	}
	return &field, nil
}

// Create 创建客户字段
// 将客户字段信息保存到数据库
// 参数：
//
//	field - 客户字段实例
//
// 返回：错误信息
func (s *ContactFieldService) Create(field *model.ContactField) error {
	return s.db.Create(field).Error
}

// Update 更新客户字段
// 更新数据库中的客户字段信息
// 参数：
//
//	field - 客户字段实例
//
// 返回：错误信息
func (s *ContactFieldService) Update(field *model.ContactField) error {
	return s.db.Save(field).Error
}

// Delete 删除客户字段
// 从数据库中删除指定 ID 的客户字段
// 参数：
//
//	id - 客户字段 ID
//
// 返回：错误信息
func (s *ContactFieldService) Delete(id uint) error {
	return s.db.Delete(&model.ContactField{}, id).Error
}

// ListByStatus 按状态获取客户字段列表（分页）
// 查询指定状态的客户字段列表，支持分页
// 参数：
//
//	status - 状态（0-禁用，1-启用，2-全部）
//	page - 页码
//	pageSize - 每页数量
//
// 返回：客户字段列表、总数和错误信息
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

// UpdateStatus 更新客户字段状态
// 更新客户的字段状态
// 参数：
//
//	id - 客户字段 ID
//	status - 状态
//
// 返回：错误信息
func (s *ContactFieldService) UpdateStatus(id uint, status int) error {
	return s.db.Model(&model.ContactField{}).Where("id = ?", id).Update("status", status).Error
}

// UpdateOrder 更新客户字段排序
// 更新客户的字段排序
// 参数：
//
//	id - 客户字段 ID
//	order - 排序值
//
// 返回：错误信息
func (s *ContactFieldService) UpdateOrder(id uint, order int) error {
	return s.db.Model(&model.ContactField{}).Where("id = ?", id).Update("order", order).Error
}

// DB 获取数据库连接
// 返回数据库连接实例
// 返回：GORM 数据库连接
func (s *ContactFieldService) DB() *gorm.DB {
	return s.db
}

// ContactFieldPivotService 客户字段值 Service
// 提供客户字段值的查询和更新功能
// 主要职责：
// 1. 获取客户字段值列表
// 2. 更新客户字段值
// 3. 批量更新客户字段值
//
// 依赖：
// - gorm.DB: 数据库连接

type ContactFieldPivotService struct {
	db *gorm.DB // 数据库连接
}

// NewContactFieldPivotService 创建客户字段值 Service 实例
// 参数：db - GORM 数据库连接
// 返回：客户字段值 Service 实例
func NewContactFieldPivotService(db *gorm.DB) *ContactFieldPivotService {
	return &ContactFieldPivotService{db: db}
}

// List 获取客户字段值列表
// 查询指定客户的所有字段值
// 参数：
//
//	contactID - 客户 ID
//
// 返回：客户字段值列表和错误信息
func (s *ContactFieldPivotService) List(contactID uint) ([]model.ContactFieldPivot, error) {
	var pivots []model.ContactFieldPivot
	if err := s.db.Where("contact_id = ?", contactID).Find(&pivots).Error; err != nil {
		return nil, err
	}
	return pivots, nil
}

// Update 更新客户字段值
// 更新指定客户的指定字段值，如果记录不存在则创建
// 参数：
//
//	contactID - 客户 ID
//	fieldID - 字段 ID
//	value - 字段值
//
// 返回：错误信息
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

// BatchUpdate 批量更新客户字段值
// 批量更新指定客户的多个字段值，使用事务确保数据一致性
// 参数：
//
//	contactID - 客户 ID
//	fields - 字段值列表，包含 FieldID 和 Value
//
// 返回：错误信息
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

// WorkContactTagGroupService 客户标签组 Service
// 提供客户标签组的 CRUD 操作功能
// 主要职责：
// 1. 获取客户标签组列表
// 2. 根据 ID 获取客户标签组详情
// 3. 创建客户标签组
// 4. 更新客户标签组
// 5. 删除客户标签组（同时删除关联的标签）
// 6. 更新客户标签组排序
//
// 依赖：
// - gorm.DB: 数据库连接

type WorkContactTagGroupService struct {
	db *gorm.DB // 数据库连接
}

// NewWorkContactTagGroupService 创建客户标签组 Service 实例
// 参数：db - GORM 数据库连接
// 返回：客户标签组 Service 实例
func NewWorkContactTagGroupService(db *gorm.DB) *WorkContactTagGroupService {
	return &WorkContactTagGroupService{db: db}
}

// List 获取客户标签组列表
// 查询指定企业的所有客户标签组，按排序字段和 ID 升序排列
// 参数：
//
//	corpID - 企业 ID
//
// 返回：客户标签组列表和错误信息
func (s *WorkContactTagGroupService) List(corpID uint) ([]model.WorkContactTagGroup, error) {
	var groups []model.WorkContactTagGroup
	if err := s.db.Where("corp_id = ?", corpID).Order("`order` ASC, id ASC").Find(&groups).Error; err != nil {
		return nil, err
	}
	return groups, nil
}

// GetByID 根据 ID 获取客户标签组详情
// 查询指定 ID 的客户标签组
// 参数：
//
//	id - 客户标签组 ID
//
// 返回：客户标签组实例和错误信息
func (s *WorkContactTagGroupService) GetByID(id uint) (*model.WorkContactTagGroup, error) {
	var group model.WorkContactTagGroup
	if err := s.db.First(&group, id).Error; err != nil {
		return nil, err
	}
	return &group, nil
}

// Create 创建客户标签组
// 将客户标签组信息保存到数据库
// 参数：
//
//	group - 客户标签组实例
//
// 返回：错误信息
func (s *WorkContactTagGroupService) Create(group *model.WorkContactTagGroup) error {
	return s.db.Create(group).Error
}

// Update 更新客户标签组
// 更新数据库中的客户标签组信息
// 参数：
//
//	group - 客户标签组实例
//
// 返回：错误信息
func (s *WorkContactTagGroupService) Update(group *model.WorkContactTagGroup) error {
	return s.db.Save(group).Error
}

// Delete 删除客户标签组
// 删除指定 ID 的客户标签组，同时删除关联的标签
// 使用事务确保数据一致性
// 参数：
//
//	id - 客户标签组 ID
//
// 返回：错误信息
func (s *WorkContactTagGroupService) Delete(id uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 先删除关联的标签
		tx.Where("contact_tag_group_id = ?", id).Delete(&model.WorkContactTag{})
		// 删除标签组
		return tx.Delete(&model.WorkContactTagGroup{}, id).Error
	})
}

// UpdateOrder 更新客户标签组排序
// 更新客户标签组的排序
// 参数：
//
//	id - 客户标签组 ID
//	order - 排序值
//
// 返回：错误信息
func (s *WorkContactTagGroupService) UpdateOrder(id uint, order int) error {
	return s.db.Model(&model.WorkContactTagGroup{}).Where("id = ?", id).Update("order", order).Error
}
