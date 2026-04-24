package plugin

import (
	"gorm.io/gorm"
	"mochat-api-server/internal/model"
)

type WorkTransferService struct {
	db *gorm.DB
}

func NewWorkTransferService(db *gorm.DB) *WorkTransferService {
	return &WorkTransferService{db: db}
}

func (s *WorkTransferService) ListByCorp(corpID uint, offset, limit int) ([]model.WorkTransferLog, int64, error) {
	var items []model.WorkTransferLog
	var total int64
	query := s.db.Model(&model.WorkTransferLog{}).Where("corp_id = ?", corpID)
	query.Count(&total)
	if err := query.Offset(offset).Limit(limit).Order("id DESC").Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

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