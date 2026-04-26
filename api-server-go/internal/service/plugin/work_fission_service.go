package plugin

import (
	"gorm.io/gorm"
	"mochat-api-server/internal/model"
)

// WorkFissionService 任务宝服务
// 提供任务宝活动的 CRUD 操作功能
// 主要职责：
// 1. 获取任务宝活动列表（分页）
// 2. 根据 ID 获取任务宝活动详情
// 3. 创建任务宝活动
// 4. 更新任务宝活动
//
// 依赖：
// - gorm.DB: 数据库连接
//
// 注意：任务宝是一种营销插件，用于通过邀请好友获得奖励

type WorkFissionService struct {
	db *gorm.DB // 数据库连接
}

// NewWorkFissionService 创建任务宝服务实例
// 参数：db - GORM 数据库连接
// 返回：任务宝服务实例
func NewWorkFissionService(db *gorm.DB) *WorkFissionService {
	return &WorkFissionService{db: db}
}

// List 获取任务宝活动列表（分页）
// 查询指定企业的任务宝活动列表，支持分页
// 参数：
//
//	corpID - 企业 ID
//	offset - 偏移量
//	limit - 限制数量
//
// 返回：任务宝活动列表、总数和错误信息
func (s *WorkFissionService) List(corpID uint, offset, limit int) ([]model.WorkFission, int64, error) {
	var items []model.WorkFission
	var total int64
	query := s.db.Model(&model.WorkFission{}).Where("corp_id = ?", corpID)
	query.Count(&total)
	if err := query.Offset(offset).Limit(limit).Order("id DESC").Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

// GetByID 根据 ID 获取任务宝活动详情
// 查询指定 ID 的任务宝活动
// 参数：
//
//	id - 任务宝活动 ID
//
// 返回：任务宝活动实例和错误信息
func (s *WorkFissionService) GetByID(id uint) (*model.WorkFission, error) {
	var item model.WorkFission
	if err := s.db.First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

// Create 创建任务宝活动
// 将任务宝活动信息保存到数据库
// 参数：
//
//	item - 任务宝活动实例
//
// 返回：错误信息
func (s *WorkFissionService) Create(item *model.WorkFission) error {
	return s.db.Create(item).Error
}

// Update 更新任务宝活动
// 更新数据库中的任务宝活动信息
// 参数：
//
//	id - 任务宝活动 ID
//	updates - 待更新的字段映射
//
// 返回：错误信息
func (s *WorkFissionService) Update(id uint, updates map[string]interface{}) error {
	return s.db.Model(&model.WorkFission{}).Where("id = ?", id).Updates(updates).Error
}
