package service

import (
	"mochat-api-server/internal/model"
	"gorm.io/gorm"
)

// WorkRoomService 群聊 Service
// 提供群聊（WorkRoom）的 CRUD 操作功能
// 主要职责：
// 1. 根据 ID 获取群聊详情
// 2. 根据企业微信群 ID 获取群聊
// 3. 获取群聊列表（分页，支持筛选）
// 4. 创建群聊
// 5. 更新群聊
// 6. 删除群聊
// 7. 批量更新群聊分组
//
// 依赖：
// - gorm.DB: 数据库连接

type WorkRoomService struct {
	db *gorm.DB // 数据库连接
}

// NewWorkRoomService 创建群聊 Service 实例
// 参数：db - GORM 数据库连接
// 返回：群聊 Service 实例
func NewWorkRoomService(db *gorm.DB) *WorkRoomService {
	return &WorkRoomService{db: db}
}

// GetByID 根据 ID 获取群聊详情
// 查询指定 ID 的群聊
// 参数：
//
//	id - 群聊 ID
//
// 返回：群聊实例和错误信息
func (s *WorkRoomService) GetByID(id uint) (*model.WorkRoom, error) {
	var room model.WorkRoom
	if err := s.db.First(&room, id).Error; err != nil {
		return nil, err
	}
	return &room, nil
}

// GetByWxChatID 根据企业微信群 ID 获取群聊
// 使用企业微信群 ID 查询群聊
// 参数：
//
//	corpID - 企业 ID
//	wxChatID - 企业微信群 ID
//
// 返回：群聊实例和错误信息
func (s *WorkRoomService) GetByWxChatID(corpID uint, wxChatID string) (*model.WorkRoom, error) {
	var room model.WorkRoom
	if err := s.db.Where("corp_id = ? AND wx_chat_id = ?", corpID, wxChatID).First(&room).Error; err != nil {
		return nil, err
	}
	return &room, nil
}

// List 获取群聊列表（分页）
// 查询指定企业的群聊列表，支持分页和筛选
// 参数：
//
//	corpID - 企业 ID
//	page - 页码
//	pageSize - 每页数量
//	filters - 筛选条件映射
//
// 返回：群聊列表、总数和错误信息
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

// Create 创建群聊
// 将群聊信息保存到数据库
// 参数：
//
//	room - 群聊实例
//
// 返回：错误信息
func (s *WorkRoomService) Create(room *model.WorkRoom) error {
	return s.db.Create(room).Error
}

// Update 更新群聊
// 更新群聊的指定字段
// 参数：
//
//	id - 群聊 ID
//	updates - 待更新的字段映射
//
// 返回：错误信息
func (s *WorkRoomService) Update(id uint, updates map[string]interface{}) error {
	return s.db.Model(&model.WorkRoom{}).Where("id = ?", id).Updates(updates).Error
}

// Delete 删除群聊
// 从数据库中删除指定 ID 的群聊
// 参数：
//
//	id - 群聊 ID
//
// 返回：错误信息
func (s *WorkRoomService) Delete(id uint) error {
	return s.db.Delete(&model.WorkRoom{}, id).Error
}

// BatchUpdateGroup 批量更新群聊分组
// 批量更新指定群聊的分组
// 参数：
//
//	ids - 群聊 ID 列表
//	groupID - 目标分组 ID
//
// 返回：错误信息
func (s *WorkRoomService) BatchUpdateGroup(ids []uint, groupID uint) error {
	return s.db.Model(&model.WorkRoom{}).Where("id IN ?", ids).Update("room_group_id", groupID).Error
}

// WorkContactRoomService 群聊客户关系 Service
// 提供群聊客户关系（WorkContactRoom）的 CRUD 操作功能
// 主要职责：
// 1. 获取群聊客户关系列表（分页）
// 2. 根据 ID 获取群聊客户关系详情
// 3. 根据客户 ID 和群聊 ID 获取群聊客户关系
// 4. 创建群聊客户关系
// 5. 更新群聊客户关系
// 6. 删除群聊客户关系
//
// 依赖：
// - gorm.DB: 数据库连接

type WorkContactRoomService struct {
	db *gorm.DB // 数据库连接
}

// NewWorkContactRoomService 创建群聊客户关系 Service 实例
// 参数：db - GORM 数据库连接
// 返回：群聊客户关系 Service 实例
func NewWorkContactRoomService(db *gorm.DB) *WorkContactRoomService {
	return &WorkContactRoomService{db: db}
}

// List 获取群聊客户关系列表（分页）
// 查询指定群聊的客户关系列表，支持分页
// 参数：
//
//	roomID - 群聊 ID
//	page - 页码
//	pageSize - 每页数量
//
// 返回：群聊客户关系列表、总数和错误信息
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

// GetByID 根据 ID 获取群聊客户关系详情
// 查询指定 ID 的群聊客户关系
// 参数：
//
//	id - 群聊客户关系 ID
//
// 返回：群聊客户关系实例和错误信息
func (s *WorkContactRoomService) GetByID(id uint) (*model.WorkContactRoom, error) {
	var contact model.WorkContactRoom
	if err := s.db.First(&contact, id).Error; err != nil {
		return nil, err
	}
	return &contact, nil
}

// GetByContactAndRoom 根据客户 ID 和群聊 ID 获取群聊客户关系
// 使用客户 ID 和群聊 ID 查询群聊客户关系
// 参数：
//
//	contactID - 客户 ID
//	roomID - 群聊 ID
//
// 返回：群聊客户关系实例和错误信息
func (s *WorkContactRoomService) GetByContactAndRoom(contactID, roomID uint) (*model.WorkContactRoom, error) {
	var contact model.WorkContactRoom
	if err := s.db.Where("contact_id = ? AND room_id = ?", contactID, roomID).First(&contact).Error; err != nil {
		return nil, err
	}
	return &contact, nil
}

// Create 创建群聊客户关系
// 将群聊客户关系信息保存到数据库
// 参数：
//
//	contact - 群聊客户关系实例
//
// 返回：错误信息
func (s *WorkContactRoomService) Create(contact *model.WorkContactRoom) error {
	return s.db.Create(contact).Error
}

// Update 更新群聊客户关系
// 更新数据库中的群聊客户关系信息
// 参数：
//
//	contact - 群聊客户关系实例
//
// 返回：错误信息
func (s *WorkContactRoomService) Update(contact *model.WorkContactRoom) error {
	return s.db.Save(contact).Error
}

// Delete 删除群聊客户关系
// 从数据库中删除指定 ID 的群聊客户关系
// 参数：
//
//	id - 群聊客户关系 ID
//
// 返回：错误信息
func (s *WorkContactRoomService) Delete(id uint) error {
	return s.db.Delete(&model.WorkContactRoom{}, id).Error
}

// WorkRoomGroupService 群聊分组 Service
// 提供群聊分组（WorkRoomGroup）的 CRUD 操作功能
// 主要职责：
// 1. 获取群聊分组列表
// 2. 创建群聊分组
// 3. 更新群聊分组
// 4. 删除群聊分组
//
// 依赖：
// - gorm.DB: 数据库连接

type WorkRoomGroupService struct {
	db *gorm.DB // 数据库连接
}

// NewWorkRoomGroupService 创建群聊分组 Service 实例
// 参数：db - GORM 数据库连接
// 返回：群聊分组 Service 实例
func NewWorkRoomGroupService(db *gorm.DB) *WorkRoomGroupService {
	return &WorkRoomGroupService{db: db}
}

// List 获取群聊分组列表
// 查询指定企业的所有群聊分组
// 参数：
//
//	corpID - 企业 ID
//
// 返回：群聊分组列表和错误信息
func (s *WorkRoomGroupService) List(corpID uint) ([]model.WorkRoomGroup, error) {
	var groups []model.WorkRoomGroup
	if err := s.db.Where("corp_id = ?", corpID).Find(&groups).Error; err != nil {
		return nil, err
	}
	return groups, nil
}

// Create 创建群聊分组
// 将群聊分组信息保存到数据库
// 参数：
//
//	group - 群聊分组实例
//
// 返回：错误信息
func (s *WorkRoomGroupService) Create(group *model.WorkRoomGroup) error {
	return s.db.Create(group).Error
}

// Update 更新群聊分组
// 更新数据库中的群聊分组信息
// 参数：
//
//	group - 群聊分组实例
//
// 返回：错误信息
func (s *WorkRoomGroupService) Update(group *model.WorkRoomGroup) error {
	return s.db.Save(group).Error
}

// Delete 删除群聊分组
// 从数据库中删除指定 ID 的群聊分组
// 参数：
//
//	id - 群聊分组 ID
//
// 返回：错误信息
func (s *WorkRoomGroupService) Delete(id uint) error {
	return s.db.Delete(&model.WorkRoomGroup{}, id).Error
}
