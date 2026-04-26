package plugin

import (
	"gorm.io/gorm"
	"mochat-api-server/internal/model"
)

// WorkRoomAutoPullService 自动拉群服务
// 提供自动拉群配置的 CRUD 操作功能
// 主要职责：
// 1. 获取自动拉群配置列表（分页）
// 2. 根据 ID 获取自动拉群配置详情
// 3. 创建自动拉群配置
// 4. 更新自动拉群配置
// 5. 删除自动拉群配置
//
// 依赖：
// - gorm.DB: 数据库连接
//
// 注意：自动拉群是一种插件功能，用于自动将客户拉入群聊

type WorkRoomAutoPullService struct {
	db *gorm.DB // 数据库连接
}

// NewWorkRoomAutoPullService 创建自动拉群服务实例
// 参数：db - GORM 数据库连接
// 返回：自动拉群服务实例
func NewWorkRoomAutoPullService(db *gorm.DB) *WorkRoomAutoPullService {
	return &WorkRoomAutoPullService{db: db}
}

// List 获取自动拉群配置列表（分页）
// 查询指定企业的自动拉群配置列表，支持分页和名称搜索
// 参数：
//
//	corpID - 企业 ID
//	qrcodeName - 群码名称（模糊搜索），空字符串表示不限制
//	offset - 偏移量
//	limit - 限制数量
//
// 返回：自动拉群配置列表、总数和错误信息
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

// GetByID 根据 ID 获取自动拉群配置详情
// 查询指定 ID 的自动拉群配置
// 参数：
//
//	id - 自动拉群配置 ID
//
// 返回：自动拉群配置实例和错误信息
func (s *WorkRoomAutoPullService) GetByID(id uint) (*model.WorkRoomAutoPull, error) {
	var item model.WorkRoomAutoPull
	if err := s.db.First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

// Create 创建自动拉群配置
// 将自动拉群配置信息保存到数据库
// 参数：
//
//	item - 自动拉群配置实例
//
// 返回：错误信息
func (s *WorkRoomAutoPullService) Create(item *model.WorkRoomAutoPull) error {
	return s.db.Create(item).Error
}

// Update 更新自动拉群配置
// 更新数据库中的自动拉群配置信息
// 参数：
//
//	id - 自动拉群配置 ID
//	updates - 待更新的字段映射
//
// 返回：错误信息
func (s *WorkRoomAutoPullService) Update(id uint, updates map[string]interface{}) error {
	return s.db.Model(&model.WorkRoomAutoPull{}).Where("id = ?", id).Updates(updates).Error
}

// Delete 删除自动拉群配置
// 从数据库中删除指定 ID 的自动拉群配置
// 参数：
//
//	id - 自动拉群配置 ID
//
// 返回：错误信息
func (s *WorkRoomAutoPullService) Delete(id uint) error {
	return s.db.Delete(&model.WorkRoomAutoPull{}, id).Error
}
