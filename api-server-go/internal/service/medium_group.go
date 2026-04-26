package service

import (
	"gorm.io/gorm"
	"mochat-api-server/internal/model"
)

// MediumGroupServiceAlt 素材分组服务(别名，用于兼容)
// 提供素材分组（MediumGroup）的 CRUD 操作功能
// 主要职责：
// 1. 获取素材分组列表
// 2. 创建素材分组
// 3. 更新素材分组
// 4. 删除素材分组
//
// 注意：此服务为 MediumGroupService 的别名版本，用于兼容旧代码
// 新代码应使用 MediumGroupService
//
// 依赖：
// - gorm.DB: 数据库连接

type MediumGroupServiceAlt struct {
	db *gorm.DB // 数据库连接
}

// NewMediumGroupServiceAlt 创建素材分组服务(别名)实例
// 参数：db - GORM 数据库连接
// 返回：素材分组服务实例
func NewMediumGroupServiceAlt(db *gorm.DB) *MediumGroupServiceAlt {
	return &MediumGroupServiceAlt{db: db}
}

// List 获取素材分组列表
// 查询指定企业的所有素材分组
// 参数：
//
//	corpID - 企业 ID
//
// 返回：素材分组列表和错误信息
func (s *MediumGroupServiceAlt) List(corpID uint) ([]model.MediumGroup, error) {
	var groups []model.MediumGroup
	if err := s.db.Where("corp_id = ?", corpID).Find(&groups).Error; err != nil {
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
func (s *MediumGroupServiceAlt) Create(group *model.MediumGroup) error {
	return s.db.Create(group).Error
}

// Update 更新素材分组
// 更新数据库中的素材分组信息
// 参数：
//
//	group - 素材分组实例
//
// 返回：错误信息
func (s *MediumGroupServiceAlt) Update(group *model.MediumGroup) error {
	return s.db.Save(group).Error
}

// Delete 删除素材分组
// 从数据库中删除指定 ID 的素材分组
// 参数：
//
//	id - 素材分组 ID
//
// 返回：错误信息
func (s *MediumGroupServiceAlt) Delete(id uint) error {
	return s.db.Delete(&model.MediumGroup{}, id).Error
}
