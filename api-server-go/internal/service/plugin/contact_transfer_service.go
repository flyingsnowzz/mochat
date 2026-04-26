package plugin

import (
	"gorm.io/gorm"
	"mochat-api-server/internal/model"
)

// ContactTransferService 客户转账服务
// 提供客户转账记录的查询功能
// 主要职责：
// 1. 获取客户转账记录列表（分页）
// 2. 根据 ID 获取客户转账记录详情
//
// 依赖：
// - gorm.DB: 数据库连接
//
// 注意：该服务目前只提供查询功能，转账的具体实现可能在其他模块

type ContactTransferService struct {
	db *gorm.DB // 数据库连接
}

// NewContactTransferService 创建客户转账服务实例
// 参数：db - GORM 数据库连接
// 返回：客户转账服务实例
func NewContactTransferService(db *gorm.DB) *ContactTransferService {
	return &ContactTransferService{db: db}
}

// List 获取客户转账记录列表（分页）
// 查询指定企业的客户转账记录列表，支持分页
// 参数：
//
//	corpID - 企业 ID
//	offset - 偏移量
//	limit - 限制数量
//
// 返回：客户转账记录列表、总数和错误信息
func (s *ContactTransferService) List(corpID uint, offset, limit int) ([]model.ContactTransfer, int64, error) {
	var items []model.ContactTransfer
	var total int64
	query := s.db.Model(&model.ContactTransfer{}).Where("corp_id = ?", corpID)
	query.Count(&total)
	if err := query.Offset(offset).Limit(limit).Order("id DESC").Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

// GetByID 根据 ID 获取客户转账记录详情
// 查询指定 ID 的客户转账记录
// 参数：
//
//	id - 客户转账记录 ID
//
// 返回：客户转账记录实例和错误信息
func (s *ContactTransferService) GetByID(id uint) (*model.ContactTransfer, error) {
	var item model.ContactTransfer
	if err := s.db.First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}
