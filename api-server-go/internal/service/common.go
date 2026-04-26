// Package service 提供业务逻辑层服务
// 该目录包含各种业务服务：
// 1. BusinessLogService - 业务日志服务
// 2. CorpDayDataService - 企业日数据服务
// 3. ChatToolService - 聊天工具服务
// 4. WorkAgentService - 企业微信应用服务
// 5. GreetingService - 欢迎语服务
// 6. TenantService - 租户服务
package service

import (
	"gorm.io/gorm"
	"mochat-api-server/internal/model"
)

// BusinessLogService 业务日志服务
// 提供业务操作的日志记录功能
// 主要职责：
// 1. 创建业务日志
//
// 依赖：
// - gorm.DB: 数据库连接

type BusinessLogService struct {
	db *gorm.DB // 数据库连接
}

// NewBusinessLogService 创建业务日志服务实例
// 参数：db - GORM 数据库连接
// 返回：业务日志服务实例
func NewBusinessLogService(db *gorm.DB) *BusinessLogService {
	return &BusinessLogService{db: db}
}

// Create 创建业务日志
// 将业务日志记录到数据库
// 参数：
//
//	log - 业务日志实例
//
// 返回：错误信息
func (s *BusinessLogService) Create(log *model.BusinessLog) error {
	return s.db.Create(log).Error
}

// CorpDayDataService 企业日数据服务
// 提供企业每日数据查询功能
// 主要职责：
// 1. 查询企业在指定日期范围内的每日数据
//
// 依赖：
// - gorm.DB: 数据库连接

type CorpDayDataService struct {
	db *gorm.DB // 数据库连接
}

// NewCorpDayDataService 创建企业日数据服务实例
// 参数：db - GORM 数据库连接
// 返回：企业日数据服务实例
func NewCorpDayDataService(db *gorm.DB) *CorpDayDataService {
	return &CorpDayDataService{db: db}
}

// List 获取企业每日数据列表
// 查询企业在指定日期范围内的每日数据
// 参数：
//
//	corpID - 企业 ID
//	startDate - 开始日期（格式：YYYY-MM-DD），空字符串表示不限制
//	endDate - 结束日期（格式：YYYY-MM-DD），空字符串表示不限制
//
// 返回：企业每日数据列表和错误信息
func (s *CorpDayDataService) List(corpID uint, startDate, endDate string) ([]model.CorpDayData, error) {
	var data []model.CorpDayData
	query := s.db.Where("corp_id = ?", corpID)
	if startDate != "" {
		query = query.Where("date >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("date <= ?", endDate)
	}
	if err := query.Order("date ASC").Find(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

// ChatToolService 聊天工具服务
// 提供聊天工具的查询功能
// 主要职责：
// 1. 查询聊天工具列表
//
// 依赖：
// - gorm.DB: 数据库连接

type ChatToolService struct {
	db *gorm.DB // 数据库连接
}

// NewChatToolService 创建聊天工具服务实例
// 参数：db - GORM 数据库连接
// 返回：聊天工具服务实例
func NewChatToolService(db *gorm.DB) *ChatToolService {
	return &ChatToolService{db: db}
}

// List 获取聊天工具列表
// 查询所有聊天工具
// 返回：聊天工具列表和错误信息
func (s *ChatToolService) List() ([]model.ChatTool, error) {
	var tools []model.ChatTool
	if err := s.db.Find(&tools).Error; err != nil {
		return nil, err
	}
	return tools, nil
}

// WorkAgentService 企业微信应用服务
// 提供企业微信应用的创建和查询功能
// 主要职责：
// 1. 创建企业微信应用
// 2. 根据企业 ID 获取应用列表
//
// 依赖：
// - gorm.DB: 数据库连接

type WorkAgentService struct {
	db *gorm.DB // 数据库连接
}

// NewWorkAgentService 创建企业微信应用服务实例
// 参数：db - GORM 数据库连接
// 返回：企业微信应用服务实例
func NewWorkAgentService(db *gorm.DB) *WorkAgentService {
	return &WorkAgentService{db: db}
}

// Create 创建企业微信应用
// 将企业微信应用信息保存到数据库
// 参数：
//
//	agent - 企业微信应用实例
//
// 返回：错误信息
func (s *WorkAgentService) Create(agent *model.WorkAgent) error {
	return s.db.Create(agent).Error
}

// GetByCorpID 根据企业 ID 获取应用列表
// 查询指定企业的所有企业微信应用
// 参数：
//
//	corpID - 企业 ID
//
// 返回：企业微信应用列表和错误信息
func (s *WorkAgentService) GetByCorpID(corpID uint) ([]model.WorkAgent, error) {
	var agents []model.WorkAgent
	if err := s.db.Where("corp_id = ?", corpID).Find(&agents).Error; err != nil {
		return nil, err
	}
	return agents, nil
}

// GreetingService 欢迎语服务
// 提供欢迎语的 CRUD 操作功能
// 主要职责：
// 1. 获取欢迎语列表（分页）
// 2. 根据 ID 获取欢迎语详情
// 3. 创建欢迎语
// 4. 更新欢迎语
// 5. 删除欢迎语
//
// 依赖：
// - gorm.DB: 数据库连接

type GreetingService struct {
	db *gorm.DB // 数据库连接
}

// NewGreetingService 创建欢迎语服务实例
// 参数：db - GORM 数据库连接
// 返回：欢迎语服务实例
func NewGreetingService(db *gorm.DB) *GreetingService {
	return &GreetingService{db: db}
}

// List 获取欢迎语列表（分页）
// 查询指定企业的欢迎语列表，支持分页
// 参数：
//
//	corpID - 企业 ID
//	page - 页码
//	pageSize - 每页数量
//
// 返回：欢迎语列表、总数和错误信息
func (s *GreetingService) List(corpID uint, page, pageSize int) ([]model.Greeting, int64, error) {
	var greetings []model.Greeting
	var total int64
	query := s.db.Model(&model.Greeting{}).Where("corp_id = ?", corpID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Find(&greetings).Error; err != nil {
		return nil, 0, err
	}
	return greetings, total, nil
}

// GetByID 根据 ID 获取欢迎语详情
// 查询指定 ID 的欢迎语
// 参数：
//
//	id - 欢迎语 ID
//
// 返回：欢迎语实例和错误信息
func (s *GreetingService) GetByID(id uint) (*model.Greeting, error) {
	var greeting model.Greeting
	if err := s.db.First(&greeting, id).Error; err != nil {
		return nil, err
	}
	return &greeting, nil
}

// Create 创建欢迎语
// 将欢迎语保存到数据库
// 参数：
//
//	greeting - 欢迎语实例
//
// 返回：错误信息
func (s *GreetingService) Create(greeting *model.Greeting) error {
	return s.db.Create(greeting).Error
}

// Update 更新欢迎语
// 更新数据库中的欢迎语信息
// 参数：
//
//	greeting - 欢迎语实例
//
// 返回：错误信息
func (s *GreetingService) Update(greeting *model.Greeting) error {
	return s.db.Save(greeting).Error
}

// Delete 删除欢迎语
// 从数据库中删除指定 ID 的欢迎语
// 参数：
//
//	id - 欢迎语 ID
//
// 返回：错误信息
func (s *GreetingService) Delete(id uint) error {
	return s.db.Delete(&model.Greeting{}, id).Error
}

// TenantService 租户服务
// 提供租户的创建功能
// 主要职责：
// 1. 创建租户
//
// 依赖：
// - gorm.DB: 数据库连接

type TenantService struct {
	db *gorm.DB // 数据库连接
}

// NewTenantService 创建租户服务实例
// 参数：db - GORM 数据库连接
// 返回：租户服务实例
func NewTenantService(db *gorm.DB) *TenantService {
	return &TenantService{db: db}
}

// Create 创建租户
// 将租户信息保存到数据库
// 参数：
//
//	tenant - 租户实例
//
// 返回：错误信息
func (s *TenantService) Create(tenant *model.Tenant) error {
	return s.db.Create(tenant).Error
}
