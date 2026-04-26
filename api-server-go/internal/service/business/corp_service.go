// Package business 提供业务逻辑层服务
// 该目录包含核心业务服务：
// 1. CorpService - 企业服务
// 2. UserService - 用户服务
package business

import (
	"gorm.io/gorm"
	"mochat-api-server/internal/model"
)

// CorpService 企业服务
// 提供企业的 CRUD 操作功能
// 主要职责：
// 1. 根据 ID 获取企业详情
// 2. 根据企业微信 CorpID 获取企业
// 3. 获取企业列表（使用 offset 方式分页）
// 4. 创建企业
// 5. 更新企业
// 6. 删除企业
//
// 依赖：
// - gorm.DB: 数据库连接

type CorpService struct {
	db *gorm.DB // 数据库连接
}

// NewCorpService 创建企业服务实例
// 参数：db - GORM 数据库连接
// 返回：企业服务实例
func NewCorpService(db *gorm.DB) *CorpService {
	return &CorpService{db: db}
}

// GetByID 根据 ID 获取企业详情
// 查询指定 ID 的企业
// 参数：
//
//	id - 企业 ID
//
// 返回：企业实例和错误信息
func (s *CorpService) GetByID(id uint) (*model.Corp, error) {
	var corp model.Corp
	if err := s.db.First(&corp, id).Error; err != nil {
		return nil, err
	}
	return &corp, nil
}

// GetByWxCorpid 根据企业微信 CorpID 获取企业
// 使用企业微信 CorpID 查询企业
// 参数：
//
//	wxCorpid - 企业微信 CorpID
//
// 返回：企业实例和错误信息
func (s *CorpService) GetByWxCorpid(wxCorpid string) (*model.Corp, error) {
	var corp model.Corp
	if err := s.db.Where("wx_corpid = ?", wxCorpid).First(&corp).Error; err != nil {
		return nil, err
	}
	return &corp, nil
}

// List 获取企业列表（使用 offset 方式分页）
// 查询企业列表，支持按租户 ID 筛选
// 参数：
//
//	tenantID - 租户 ID，0 表示不限制
//	offset - 偏移量
//	limit - 限制数量
//
// 返回：企业列表、总数和错误信息
func (s *CorpService) List(tenantID uint, offset, limit int) ([]model.Corp, int64, error) {
	var corps []model.Corp
	var total int64
	query := s.db.Model(&model.Corp{})
	if tenantID > 0 {
		query = query.Where("tenant_id = ?", tenantID)
	}
	query.Count(&total)
	if err := query.Offset(offset).Limit(limit).Order("id DESC").Find(&corps).Error; err != nil {
		return nil, 0, err
	}
	return corps, total, nil
}

// Create 创建企业
// 将企业信息保存到数据库
// 参数：
//
//	corp - 企业实例
//
// 返回：错误信息
func (s *CorpService) Create(corp *model.Corp) error {
	return s.db.Create(corp).Error
}

// Update 更新企业
// 更新企业的指定字段
// 参数：
//
//	id - 企业 ID
//	updates - 待更新的字段映射
//
// 返回：错误信息
func (s *CorpService) Update(id uint, updates map[string]interface{}) error {
	return s.db.Model(&model.Corp{}).Where("id = ?", id).Updates(updates).Error
}

// Delete 删除企业
// 从数据库中删除指定 ID 的企业
// 参数：
//
//	id - 企业 ID
//
// 返回：错误信息
func (s *CorpService) Delete(id uint) error {
	return s.db.Delete(&model.Corp{}, id).Error
}
