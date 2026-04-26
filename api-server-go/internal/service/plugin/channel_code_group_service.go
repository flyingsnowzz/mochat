package plugin

import (
	"mochat-api-server/internal/model"

	"gorm.io/gorm"
)

// ChannelCodeGroupService 渠道码分组服务
// 提供渠道码分组的 CRUD 操作功能
// 主要职责：
// 1. 获取渠道码分组列表（根据企业 ID）
// 2. 根据 ID 获取渠道码分组详情
// 3. 创建渠道码分组
// 4. 更新渠道码分组
// 5. 删除渠道码分组
// 6. 根据分组名称获取渠道码分组
// 7. 根据分组名称列表获取渠道码分组列表
//
// 依赖：
// - gorm.DB: 数据库连接

type ChannelCodeGroupService struct {
	db *gorm.DB // 数据库连接
}

// NewChannelCodeGroupService 创建渠道码分组服务实例
// 参数：db - GORM 数据库连接
// 返回：渠道码分组服务实例
func NewChannelCodeGroupService(db *gorm.DB) *ChannelCodeGroupService {
	return &ChannelCodeGroupService{db: db}
}

// List 获取渠道码分组列表
// 根据企业 ID 获取渠道码分组列表
// 参数：
//
//	corpID - 企业 ID
//	columns - 查询字段，默认为全部字段
//
// 返回：渠道码分组列表和错误信息
func (s *ChannelCodeGroupService) List(corpID uint, columns ...string) ([]model.ChannelCodeGroup, error) {
	var groups []model.ChannelCodeGroup
	query := s.db.Model(&model.ChannelCodeGroup{}).Where("corp_id = ?", corpID)
	if len(columns) > 0 {
		query = query.Select(columns)
	}
	if err := query.Find(&groups).Error; err != nil {
		return nil, err
	}
	return groups, nil
}

// GetByID 根据 ID 获取渠道码分组详情
// 查询指定 ID 的渠道码分组
// 参数：
//
//	id - 渠道码分组 ID
//	columns - 查询字段，默认为全部字段
//
// 返回：渠道码分组实例和错误信息
func (s *ChannelCodeGroupService) GetByID(id uint, columns ...string) (*model.ChannelCodeGroup, error) {
	var group model.ChannelCodeGroup
	query := s.db.Model(&model.ChannelCodeGroup{}).Where("id = ?", id)
	if len(columns) > 0 {
		query = query.Select(columns)
	}
	if err := query.First(&group).Error; err != nil {
		return nil, err
	}
	return &group, nil
}

// Create 创建渠道码分组
// 将渠道码分组信息保存到数据库
// 参数：
//
//	group - 渠道码分组实例
//
// 返回：错误信息
func (s *ChannelCodeGroupService) Create(group *model.ChannelCodeGroup) error {
	return s.db.Create(group).Error
}

// Update 更新渠道码分组
// 更新指定 ID 的渠道码分组信息
// 参数：
//
//	id - 渠道码分组 ID
//	updates - 待更新的字段映射
//
// 返回：错误信息
func (s *ChannelCodeGroupService) Update(id uint, updates map[string]interface{}) error {
	return s.db.Model(&model.ChannelCodeGroup{}).Where("id = ?", id).Updates(updates).Error
}

// Delete 删除渠道码分组
// 从数据库中删除指定 ID 的渠道码分组
// 参数：
//
//	id - 渠道码分组 ID
//
// 返回：错误信息
func (s *ChannelCodeGroupService) Delete(id uint) error {
	return s.db.Delete(&model.ChannelCodeGroup{}, id).Error
}

// GetByName 根据分组名称获取渠道码分组
// 使用分组名称查询渠道码分组
// 参数：
//
//	name - 分组名称
//	columns - 查询字段，默认为全部字段
//
// 返回：渠道码分组列表和错误信息
func (s *ChannelCodeGroupService) GetByName(name string, columns ...string) ([]model.ChannelCodeGroup, error) {
	var groups []model.ChannelCodeGroup
	query := s.db.Model(&model.ChannelCodeGroup{}).Where("name = ?", name)
	if len(columns) > 0 {
		query = query.Select(columns)
	}
	if err := query.Find(&groups).Error; err != nil {
		return nil, err
	}
	return groups, nil
}

// GetByNames 根据分组名称列表获取渠道码分组
// 使用分组名称列表查询渠道码分组
// 参数：
//
//	names - 分组名称列表
//	columns - 查询字段，默认为全部字段
//
// 返回：渠道码分组列表和错误信息
func (s *ChannelCodeGroupService) GetByNames(names []string, columns ...string) ([]model.ChannelCodeGroup, error) {
	if len(names) == 0 {
		return []model.ChannelCodeGroup{}, nil
	}
	var groups []model.ChannelCodeGroup
	query := s.db.Model(&model.ChannelCodeGroup{}).Where("name IN ?", names)
	if len(columns) > 0 {
		query = query.Select(columns)
	}
	if err := query.Find(&groups).Error; err != nil {
		return nil, err
	}
	return groups, nil
}
