// Package plugin 提供插件相关的业务服务
// 该目录包含各种插件功能的服务：
// 1. ChannelCodeService - 渠道码服务
// 2. ContactMessageBatchSendService - 客户消息批量发送服务
// 3. ContactTransferService - 客户转账服务
// 4. RoomMessageBatchSendService - 群聊消息批量发送服务
// 5. RoomTagPullService - 群聊标签拉取服务
// 6. RoomWelcomeService - 群聊欢迎语服务
// 7. WorkFissionService - 任务宝服务
// 8. WorkRoomAutoPullService - 自动拉群服务
package plugin

import (
	"gorm.io/gorm"
	"mochat-api-server/internal/model"
)

// ChannelCodeService 渠道码服务
// 提供渠道码的 CRUD 操作功能
// 主要职责：
// 1. 获取渠道码列表（分页）
// 2. 根据 ID 获取渠道码详情
// 3. 创建渠道码
// 4. 更新渠道码
// 5. 删除渠道码
// 6. 更新渠道码状态
//
// 依赖：
// - gorm.DB: 数据库连接

type ChannelCodeService struct {
	db *gorm.DB // 数据库连接
}

// NewChannelCodeService 创建渠道码服务实例
// 参数：db - GORM 数据库连接
// 返回：渠道码服务实例
func NewChannelCodeService(db *gorm.DB) *ChannelCodeService {
	return &ChannelCodeService{db: db}
}

// List 获取渠道码列表（分页）
// 查询指定企业的渠道码列表，支持分页和名称搜索
// 参数：
//
//	corpID - 企业 ID
//	qrcodeName - 渠道码名称（模糊搜索），空字符串表示不限制
//	offset - 偏移量
//	limit - 限制数量
//
// 返回：渠道码列表、总数和错误信息
func (s *ChannelCodeService) List(corpID uint, qrcodeName string, offset, limit int) ([]model.ChannelCode, int64, error) {
	var items []model.ChannelCode
	var total int64
	query := s.db.Model(&model.ChannelCode{}).Where("corp_id = ?", corpID)
	if qrcodeName != "" {
		query = query.Where("qrcode_name LIKE ?", "%"+qrcodeName+"%")
	}
	query.Count(&total)
	if err := query.Offset(offset).Limit(limit).Order("id DESC").Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

// GetByID 根据 ID 获取渠道码详情
// 查询指定 ID 的渠道码
// 参数：
//
//	id - 渠道码 ID
//
// 返回：渠道码实例和错误信息
func (s *ChannelCodeService) GetByID(id uint) (*model.ChannelCode, error) {
	var item model.ChannelCode
	if err := s.db.First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

// Create 创建渠道码
// 将渠道码信息保存到数据库
// 参数：
//
//	item - 渠道码实例
//
// 返回：错误信息
func (s *ChannelCodeService) Create(item *model.ChannelCode) error {
	return s.db.Create(item).Error
}

// Update 更新渠道码
// 更新数据库中的渠道码信息
// 参数：
//
//	item - 渠道码实例
//
// 返回：错误信息
func (s *ChannelCodeService) Update(item *model.ChannelCode) error {
	return s.db.Save(item).Error
}

// UpdateByID 根据 ID 更新渠道码
// 使用映射更新渠道码的指定字段
// 参数：
//
//	id - 渠道码 ID
//	updates - 待更新的字段映射
//
// 返回：错误信息
func (s *ChannelCodeService) UpdateByID(id uint, updates map[string]interface{}) error {
	return s.db.Model(&model.ChannelCode{}).Where("id = ?", id).Updates(updates).Error
}

// Delete 删除渠道码
// 从数据库中删除指定 ID 的渠道码
// 参数：
//
//	id - 渠道码 ID
//
// 返回：错误信息
func (s *ChannelCodeService) Delete(id uint) error {
	return s.db.Delete(&model.ChannelCode{}, id).Error
}

// UpdateStatus 更新渠道码状态
// 更新渠道码的启用/禁用状态
// 参数：
//
//	id - 渠道码 ID
//	status - 状态
//
// 返回：错误信息
func (s *ChannelCodeService) UpdateStatus(id uint, status int) error {
	return s.db.Model(&model.ChannelCode{}).Where("id = ?", id).Update("status", status).Error
}
