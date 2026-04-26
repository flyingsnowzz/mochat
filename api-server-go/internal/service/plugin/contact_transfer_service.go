package plugin

import (
	"mochat-api-server/internal/model"

	"gorm.io/gorm"
)

// WorkTransferService 客户转移服务
// 提供客户转移相关的操作功能
// 主要职责：
// 1. 获取客户转移记录列表（分页）
// 2. 根据 ID 获取客户转移记录详情
// 3. 获取未分配客户列表（分页）
// 4. 根据 ID 获取未分配客户详情
//
// 依赖：
// - gorm.DB: 数据库连接

type WorkTransferService struct {
	db *gorm.DB // 数据库连接
}

// NewWorkTransferService 创建客户转移服务实例
// 参数：db - GORM 数据库连接
// 返回：客户转移服务实例
func NewWorkTransferService(db *gorm.DB) *WorkTransferService {
	return &WorkTransferService{db: db}
}

// ListTransferLogs 获取客户转移记录列表（分页）
// 查询指定企业的客户转移记录列表，支持分页
// 参数：
//
//	corpID - 企业 ID
//	offset - 偏移量
//	limit - 限制数量
//
// 返回：客户转移记录列表、总数和错误信息
func (s *WorkTransferService) ListTransferLogs(corpID uint, offset, limit int) ([]model.WorkTransferLog, int64, error) {
	var items []model.WorkTransferLog
	var total int64
	query := s.db.Model(&model.WorkTransferLog{}).Where("corp_id = ?", corpID)
	query.Count(&total)
	if err := query.Offset(offset).Limit(limit).Order("id DESC").Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

// GetTransferLogByID 根据 ID 获取客户转移记录详情
// 查询指定 ID 的客户转移记录
// 参数：
//
//	id - 客户转移记录 ID
//
// 返回：客户转移记录实例和错误信息
func (s *WorkTransferService) GetTransferLogByID(id uint) (*model.WorkTransferLog, error) {
	var item model.WorkTransferLog
	if err := s.db.First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

// ListUnassigned 获取未分配客户列表（分页）
// 查询指定企业的未分配客户列表，支持分页
// 参数：
//
//	corpID - 企业 ID
//	offset - 偏移量
//	limit - 限制数量
//
// 返回：未分配客户列表、总数和错误信息
func (s *WorkTransferService) ListUnassigned(corpID uint, offset, limit int) ([]model.WorkUnassigned, int64, error) {
	var items []model.WorkUnassigned
	var total int64
	query := s.db.Model(&model.WorkUnassigned{}).Where("corp_id = ?", corpID)
	query.Count(&total)
	if err := query.Offset(offset).Limit(limit).Order("id DESC").Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

// GetUnassignedByID 根据 ID 获取未分配客户详情
// 查询指定 ID 的未分配客户
// 参数：
//
//	id - 未分配客户 ID
//
// 返回：未分配客户实例和错误信息
func (s *WorkTransferService) GetUnassignedByID(id uint) (*model.WorkUnassigned, error) {
	var item model.WorkUnassigned
	if err := s.db.First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}
