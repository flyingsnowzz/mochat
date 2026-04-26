package plugin

import (
	"gorm.io/gorm"
	"mochat-api-server/internal/model"
)

// RoomMessageBatchSendService 群聊消息批量发送服务
// 提供群聊消息批量发送记录的查询功能
// 主要职责：
// 1. 获取群聊消息批量发送记录列表（分页）
// 2. 根据 ID 获取群聊消息批量发送记录详情
//
// 依赖：
// - gorm.DB: 数据库连接
//
// 注意：该服务目前只提供查询功能，批量发送的具体实现可能在其他模块

type RoomMessageBatchSendService struct {
	db *gorm.DB // 数据库连接
}

// NewRoomMessageBatchSendService 创建群聊消息批量发送服务实例
// 参数：db - GORM 数据库连接
// 返回：群聊消息批量发送服务实例
func NewRoomMessageBatchSendService(db *gorm.DB) *RoomMessageBatchSendService {
	return &RoomMessageBatchSendService{db: db}
}

// List 获取群聊消息批量发送记录列表（分页）
// 查询指定企业的群聊消息批量发送记录列表，支持分页
// 参数：
//
//	corpID - 企业 ID
//	offset - 偏移量
//	limit - 限制数量
//
// 返回：群聊消息批量发送记录列表、总数和错误信息
func (s *RoomMessageBatchSendService) List(corpID uint, offset, limit int) ([]model.RoomMessageBatchSend, int64, error) {
	var items []model.RoomMessageBatchSend
	var total int64
	query := s.db.Model(&model.RoomMessageBatchSend{}).Where("corp_id = ?", corpID)
	query.Count(&total)
	if err := query.Offset(offset).Limit(limit).Order("id DESC").Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

// GetByID 根据 ID 获取群聊消息批量发送记录详情
// 查询指定 ID 的群聊消息批量发送记录
// 参数：
//
//	id - 群聊消息批量发送记录 ID
//
// 返回：群聊消息批量发送记录实例和错误信息
func (s *RoomMessageBatchSendService) GetByID(id uint) (*model.RoomMessageBatchSend, error) {
	var item model.RoomMessageBatchSend
	if err := s.db.First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}
