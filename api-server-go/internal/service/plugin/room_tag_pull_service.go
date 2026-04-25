package plugin

import (
	"gorm.io/gorm"

	"mochat-api-server/internal/model"
)

type RoomTagPullService struct {
	db *gorm.DB
}

func NewRoomTagPullService(db *gorm.DB) *RoomTagPullService {
	return &RoomTagPullService{
		db: db,
	}
}

// List 获取标签建群列表
func (s *RoomTagPullService) List(name string, page, perPage int) (map[string]interface{}, error) {
	var items []model.RoomTagPull
	var total int64

	// 构建查询
	query := s.db.Model(&model.RoomTagPull{})
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// 分页查询
	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Find(&items).Error; err != nil {
		return nil, err
	}

	// 构建响应
	result := map[string]interface{}{
		"list":  items,
		"total": total,
		"page": map[string]int{
			"currentPage": page,
			"perPage":     perPage,
			"total":       int(total),
			"totalPage":   (int(total) + perPage - 1) / perPage,
		},
		"pageSize": perPage,
	}

	return result, nil
}

// GetByID 根据ID获取标签建群详情
func (s *RoomTagPullService) GetByID(id uint) (model.RoomTagPull, error) {
	var item model.RoomTagPull
	if err := s.db.First(&item, id).Error; err != nil {
		return item, err
	}
	return item, nil
}

// Create 创建标签建群
func (s *RoomTagPullService) Create(item *model.RoomTagPull) error {
	return s.db.Create(item).Error
}
