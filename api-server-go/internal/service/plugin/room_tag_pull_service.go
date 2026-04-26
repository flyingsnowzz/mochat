package plugin

import (
	"mochat-api-server/internal/model"

	"gorm.io/gorm"
)

// RoomTagPullService 群聊标签拉取服务
// 提供群聊标签拉取记录的查询功能
// 主要职责：
// 1. 获取群聊标签拉取记录列表（分页）
// 2. 根据 ID 获取群聊标签拉取记录详情
//
// 依赖：
// - gorm.DB: 数据库连接
//
// 注意：该服务目前只提供查询功能，拉取的具体实现可能在其他模块

type RoomTagPullService struct {
	db *gorm.DB // 数据库连接
}

// NewRoomTagPullService 创建群聊标签拉取服务实例
// 参数：db - GORM 数据库连接
// 返回：群聊标签拉取服务实例
func NewRoomTagPullService(db *gorm.DB) *RoomTagPullService {
	return &RoomTagPullService{db: db}
}

// List 获取群聊标签拉取记录列表（分页）
// 查询指定企业的群聊标签拉取记录列表，支持分页和名称搜索
// 参数：
//
//	corpID - 企业 ID
//	name - 名称（模糊搜索），空字符串表示不限制
//	offset - 偏移量
//	limit - 限制数量
//
// 返回：群聊标签拉取记录列表、总数和错误信息
func (s *RoomTagPullService) List(corpID uint, name string, offset, limit int) ([]model.RoomTagPull, int64, error) {
	var items []model.RoomTagPull
	var total int64
	query := s.db.Model(&model.RoomTagPull{}).Where("corp_id = ?", corpID)
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	query.Count(&total)
	if err := query.Offset(offset).Limit(limit).Order("id DESC").Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

// GetByID 根据 ID 获取群聊标签拉取记录详情
// 查询指定 ID 的群聊标签拉取记录
// 参数：
//
//	id - 群聊标签拉取记录 ID
//
// 返回：群聊标签拉取记录实例和错误信息
func (s *RoomTagPullService) GetByID(id uint) (*model.RoomTagPull, error) {
	var item model.RoomTagPull
	if err := s.db.First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

// Create 创建群聊标签拉取记录
// 将群聊标签拉取记录信息保存到数据库
// 参数：
//
//	item - 群聊标签拉取记录实例
//
// 返回：错误信息
func (s *RoomTagPullService) Create(item *model.RoomTagPull) error {
	return s.db.Create(item).Error
}
