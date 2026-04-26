package service

import (
	"encoding/json"
	"fmt"

	"gorm.io/gorm"
	"mochat-api-server/internal/model"
)

// MediumService 素材 Service
// 提供素材（Medium）的 CRUD 操作功能
// 主要职责：
// 1. 获取素材列表（分页，支持按类型和分组筛选）
// 2. 获取素材列表（分页，支持搜索）
// 3. 根据 ID 获取素材详情
// 4. 创建素材
// 5. 更新素材
// 6. 删除素材
//
// 依赖：
// - gorm.DB: 数据库连接

type MediumService struct {
	db *gorm.DB // 数据库连接
}

// NewMediumService 创建素材 Service 实例
// 参数：db - GORM 数据库连接
// 返回：素材 Service 实例
func NewMediumService(db *gorm.DB) *MediumService {
	return &MediumService{db: db}
}

// List 获取素材列表（分页）
// 查询指定企业的素材列表，支持按类型和分组筛选
// 参数：
//
//	corpID - 企业 ID
//	page - 页码
//	pageSize - 每页数量
//	mediumType - 素材类型，0 表示全部
//	groupID - 素材分组 ID，0 表示全部
//
// 返回：素材列表、总数和错误信息
func (s *MediumService) List(corpID uint, page, pageSize int, mediumType int, groupID uint) ([]model.Medium, int64, error) {
	var media []model.Medium
	var total int64
	query := s.db.Model(&model.Medium{}).Where("corp_id = ?", corpID)
	if mediumType > 0 {
		query = query.Where("type = ?", mediumType)
	}
	if groupID > 0 {
		query = query.Where("medium_group_id = ?", groupID)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&media).Error; err != nil {
		return nil, 0, err
	}
	return media, total, nil
}

// ListWithSearch 获取素材列表（分页，支持搜索）
// 查询指定企业的素材列表，支持按类型、分组筛选和内容搜索
// 参数：
//
//	corpID - 企业 ID
//	page - 页码
//	pageSize - 每页数量
//	mediumType - 素材类型，0 表示全部
//	groupID - 素材分组 ID，0 表示全部
//	searchStr - 搜索关键词（搜索 content 字段）
//
// 返回：素材列表、总数和错误信息
func (s *MediumService) ListWithSearch(corpID uint, page, pageSize int, mediumType int, groupID uint, searchStr string) ([]model.Medium, int64, error) {
	var media []model.Medium
	var total int64
	query := s.db.Model(&model.Medium{}).Where("corp_id = ?", corpID)
	if mediumType > 0 {
		query = query.Where("type = ?", mediumType)
	}
	if groupID > 0 {
		query = query.Where("medium_group_id = ?", groupID)
	}
	if searchStr != "" {
		query = query.Where("content LIKE ?", "%"+searchStr+"%")
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&media).Error; err != nil {
		return nil, 0, err
	}
	return media, total, nil
}

// GetByID 根据 ID 获取素材详情
// 查询指定 ID 的素材
// 参数：
//
//	id - 素材 ID
//
// 返回：素材实例和错误信息
func (s *MediumService) GetByID(id uint) (*model.Medium, error) {
	var medium model.Medium
	if err := s.db.First(&medium, id).Error; err != nil {
		return nil, err
	}
	return &medium, nil
}

// Create 创建素材
// 将素材信息保存到数据库
// 参数：
//
//	medium - 素材实例
//
// 返回：错误信息
func (s *MediumService) Create(medium *model.Medium) error {
	return s.db.Create(medium).Error
}

// Update 更新素材
// 更新数据库中的素材信息
// 参数：
//
//	medium - 素材实例
//
// 返回：错误信息
func (s *MediumService) Update(medium *model.Medium) error {
	return s.db.Save(medium).Error
}

// Delete 删除素材
// 从数据库中删除指定 ID 的素材
// 参数：
//
//	id - 素材 ID
//
// 返回：错误信息
func (s *MediumService) Delete(id uint) error {
	return s.db.Delete(&model.Medium{}, id).Error
}

// ParseMediumContent 解析素材内容
// 将素材的 JSON 内容字符串解析为 map
// 参数：
//
//	raw - 素材内容 JSON 字符串
//
// 返回：解析后的 map，如果解析失败则返回空 map
func ParseMediumContent(raw string) map[string]interface{} {
	if raw == "" {
		return map[string]interface{}{}
	}
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(raw), &result); err != nil {
		return map[string]interface{}{}
	}
	return result
}

// MarshalMediumContent 序列化素材内容
// 将 map 序列化为 JSON 字符串
// 参数：
//
//	content - 素材内容（map 或其他可序列化类型）
//
// 返回：JSON 字符串和错误信息
func MarshalMediumContent(content interface{}) (string, error) {
	if content == nil {
		return "{}", nil
	}
	data, err := json.Marshal(content)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// MediumTypeName 获取素材类型名称
// 根据素材类型值返回对应的中文名称
// 参数：
//
//	mediumType - 素材类型值
//
// 返回：素材类型中文名称
//
// 类型对应关系：
// 1 - 文本
// 2 - 图片
// 3 - 图文
// 4 - 音频
// 5 - 视频
// 6 - 小程序
// 7 - 文件
func MediumTypeName(mediumType int) string {
	switch mediumType {
	case 1:
		return "文本"
	case 2:
		return "图片"
	case 3:
		return "图文"
	case 4:
		return "音频"
	case 5:
		return "视频"
	case 6:
		return "小程序"
	case 7:
		return "文件"
	default:
		return fmt.Sprintf("%d", mediumType)
	}
}

// MediumGroupService 素材分组 Service
// 提供素材分组（MediumGroup）的 CRUD 操作功能
// 主要职责：
// 1. 获取素材分组列表
// 2. 创建素材分组
// 3. 更新素材分组
// 4. 删除素材分组
//
// 依赖：
// - gorm.DB: 数据库连接

type MediumGroupService struct {
	db *gorm.DB // 数据库连接
}

// NewMediumGroupService 创建素材分组 Service 实例
// 参数：db - GORM 数据库连接
// 返回：素材分组 Service 实例
func NewMediumGroupService(db *gorm.DB) *MediumGroupService {
	return &MediumGroupService{db: db}
}

// List 获取素材分组列表
// 查询指定企业的所有素材分组，按排序字段升序排列
// 参数：
//
//	corpID - 企业 ID
//
// 返回：素材分组列表和错误信息
func (s *MediumGroupService) List(corpID uint) ([]model.MediumGroup, error) {
	var groups []model.MediumGroup
	if err := s.db.Where("corp_id = ?", corpID).Order("`order` ASC").Find(&groups).Error; err != nil {
		return nil, err
	}
	return groups, nil
}

// Create 创建素材分组
// 将素材分组信息保存到数据库
// 参数：
//
//	group - 素材分组实例
//
// 返回：错误信息
func (s *MediumGroupService) Create(group *model.MediumGroup) error {
	return s.db.Create(group).Error
}

// Update 更新素材分组
// 更新数据库中的素材分组信息
// 参数：
//
//	group - 素材分组实例
//
// 返回：错误信息
func (s *MediumGroupService) Update(group *model.MediumGroup) error {
	return s.db.Save(group).Error
}

// Delete 删除素材分组
// 从数据库中删除指定 ID 的素材分组
// 参数：
//
//	id - 素材分组 ID
//
// 返回：错误信息
func (s *MediumGroupService) Delete(id uint) error {
	return s.db.Delete(&model.MediumGroup{}, id).Error
}
