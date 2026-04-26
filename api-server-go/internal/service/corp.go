package service

import (
	"gorm.io/gorm"
	"mochat-api-server/internal/model"
)

// CorpService 企业 Service
// 提供企业（Corp）的 CRUD 操作功能
// 主要职责：
// 1. 获取企业列表（分页，支持按名称筛选）
// 2. 根据 ID 获取企业详情
// 3. 创建企业
// 4. 更新企业
// 5. 获取企业选择列表
//
// 依赖：
// - gorm.DB: 数据库连接

type CorpService struct {
	db *gorm.DB // 数据库连接
}

// NewCorpService 创建企业 Service 实例
// 参数：db - GORM 数据库连接
// 返回：企业 Service 实例
func NewCorpService(db *gorm.DB) *CorpService {
	return &CorpService{db: db}
}

// List 获取企业列表（分页）
// 查询企业列表，支持按租户 ID 和企业名称筛选
// 参数：
//
//	tenantID - 租户 ID，0 表示不限制
//	corpName - 企业名称（模糊搜索），空字符串表示不限制
//	page - 页码
//	pageSize - 每页数量
//
// 返回：企业列表、总数和错误信息
func (s *CorpService) List(tenantID uint, corpName string, page, pageSize int) ([]model.Corp, int64, error) {
	var corps []model.Corp
	var total int64
	query := s.db.Model(&model.Corp{})
	if tenantID > 0 {
		query = query.Where("tenant_id = ?", tenantID)
	}
	if corpName != "" {
		query = query.Where("name LIKE ?", "%"+corpName+"%")
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * pageSize
	if err := query.Order("id DESC").Offset(offset).Limit(pageSize).Find(&corps).Error; err != nil {
		return nil, 0, err
	}
	return corps, total, nil
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
// 更新数据库中的企业信息
// 参数：
//
//	corp - 企业实例
//
// 返回：错误信息
func (s *CorpService) Update(corp *model.Corp) error {
	return s.db.Save(corp).Error
}

// Select 获取企业选择列表
// 查询企业列表，用于下拉选择
// 参数：
//
//	tenantID - 租户 ID，0 表示不限制
//
// 返回：企业列表和错误信息
func (s *CorpService) Select(tenantID uint) ([]model.Corp, error) {
	var corps []model.Corp
	query := s.db.Model(&model.Corp{})
	if tenantID > 0 {
		query = query.Where("tenant_id = ?", tenantID)
	}
	if err := query.Find(&corps).Error; err != nil {
		return nil, err
	}
	return corps, nil
}
