package plugin

import (
	"gorm.io/gorm"
	"mochat-api-server/internal/model"
)

// RoomWelcomeService 群聊欢迎语服务
// 提供群聊欢迎语的 CRUD 操作功能
// 主要职责：
// 1. 获取群聊欢迎语列表
// 2. 根据 ID 获取群聊欢迎语详情
// 3. 创建群聊欢迎语
// 4. 更新群聊欢迎语
// 5. 删除群聊欢迎语
//
// 依赖：
// - gorm.DB: 数据库连接

type RoomWelcomeService struct {
	db *gorm.DB // 数据库连接
}

// NewRoomWelcomeService 创建群聊欢迎语服务实例
// 参数：db - GORM 数据库连接
// 返回：群聊欢迎语服务实例
func NewRoomWelcomeService(db *gorm.DB) *RoomWelcomeService {
	return &RoomWelcomeService{db: db}
}

// List 获取群聊欢迎语列表
// 查询指定企业的所有群聊欢迎语
// 参数：
//
//	corpID - 企业 ID
//
// 返回：群聊欢迎语列表和错误信息
func (s *RoomWelcomeService) List(corpID uint) ([]model.RoomWelcomeTemplate, error) {
	var items []model.RoomWelcomeTemplate
	if err := s.db.Where("corp_id = ?", corpID).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// GetByID 根据 ID 获取群聊欢迎语详情
// 查询指定 ID 的群聊欢迎语
// 参数：
//
//	id - 群聊欢迎语 ID
//
// 返回：群聊欢迎语实例和错误信息
func (s *RoomWelcomeService) GetByID(id uint) (*model.RoomWelcomeTemplate, error) {
	var item model.RoomWelcomeTemplate
	if err := s.db.First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

// Create 创建群聊欢迎语
// 将群聊欢迎语信息保存到数据库
// 参数：
//
//	item - 群聊欢迎语实例
//
// 返回：错误信息
func (s *RoomWelcomeService) Create(item *model.RoomWelcomeTemplate) error {
	return s.db.Create(item).Error
}

// Update 更新群聊欢迎语
// 更新数据库中的群聊欢迎语信息，只更新指定字段（msg_text, msg_complex, complex_type）
// 参数：
//
//	id - 群聊欢迎语 ID
//	updates - 待更新的字段映射
//
// 返回：错误信息
func (s *RoomWelcomeService) Update(id uint, updates map[string]interface{}) error {
	// 使用 Select 方法强制更新所有字段，包括零值字段
	return s.db.Model(&model.RoomWelcomeTemplate{}).Where("id = ?", id).Select("msg_text", "msg_complex", "complex_type").Updates(updates).Error
}

// Delete 删除群聊欢迎语
// 从数据库中删除指定 ID 的群聊欢迎语
// 参数：
//
//	id - 群聊欢迎语 ID
//
// 返回：错误信息
func (s *RoomWelcomeService) Delete(id uint) error {
	return s.db.Delete(&model.RoomWelcomeTemplate{}, id).Error
}
