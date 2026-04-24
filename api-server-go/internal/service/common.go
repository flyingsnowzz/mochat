package service

import (
	"gorm.io/gorm"
	"mochat-api-server/internal/model"
)

type BusinessLogService struct {
	db *gorm.DB
}

func NewBusinessLogService(db *gorm.DB) *BusinessLogService {
	return &BusinessLogService{db: db}
}

func (s *BusinessLogService) Create(log *model.BusinessLog) error {
	return s.db.Create(log).Error
}

type CorpDayDataService struct {
	db *gorm.DB
}

func NewCorpDayDataService(db *gorm.DB) *CorpDayDataService {
	return &CorpDayDataService{db: db}
}

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

type ChatToolService struct {
	db *gorm.DB
}

func NewChatToolService(db *gorm.DB) *ChatToolService {
	return &ChatToolService{db: db}
}

func (s *ChatToolService) List() ([]model.ChatTool, error) {
	var tools []model.ChatTool
	if err := s.db.Find(&tools).Error; err != nil {
		return nil, err
	}
	return tools, nil
}

type WorkAgentService struct {
	db *gorm.DB
}

func NewWorkAgentService(db *gorm.DB) *WorkAgentService {
	return &WorkAgentService{db: db}
}

func (s *WorkAgentService) Create(agent *model.WorkAgent) error {
	return s.db.Create(agent).Error
}

func (s *WorkAgentService) GetByCorpID(corpID uint) ([]model.WorkAgent, error) {
	var agents []model.WorkAgent
	if err := s.db.Where("corp_id = ?", corpID).Find(&agents).Error; err != nil {
		return nil, err
	}
	return agents, nil
}

type GreetingService struct {
	db *gorm.DB
}

func NewGreetingService(db *gorm.DB) *GreetingService {
	return &GreetingService{db: db}
}

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

func (s *GreetingService) GetByID(id uint) (*model.Greeting, error) {
	var greeting model.Greeting
	if err := s.db.First(&greeting, id).Error; err != nil {
		return nil, err
	}
	return &greeting, nil
}

func (s *GreetingService) Create(greeting *model.Greeting) error {
	return s.db.Create(greeting).Error
}

func (s *GreetingService) Update(greeting *model.Greeting) error {
	return s.db.Save(greeting).Error
}

func (s *GreetingService) Delete(id uint) error {
	return s.db.Delete(&model.Greeting{}, id).Error
}

type TenantService struct {
	db *gorm.DB
}

func NewTenantService(db *gorm.DB) *TenantService {
	return &TenantService{db: db}
}

func (s *TenantService) Create(tenant *model.Tenant) error {
	return s.db.Create(tenant).Error
}
