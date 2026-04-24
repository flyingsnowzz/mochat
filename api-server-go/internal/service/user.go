package service

import (
	"mochat-api-server/internal/model"
	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

func (s *UserService) GetByPhone(phone string) (*model.User, error) {
	var user model.User
	if err := s.db.Where("phone = ?", phone).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserService) GetByID(id uint) (*model.User, error) {
	var user model.User
	if err := s.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserService) List(tenantID uint, page, pageSize int) ([]model.User, int64, error) {
	var users []model.User
	var total int64
	query := s.db.Model(&model.User{})
	if tenantID > 0 {
		query = query.Where("tenant_id = ?", tenantID)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

func (s *UserService) Create(user *model.User) error {
	return s.db.Create(user).Error
}

func (s *UserService) Update(user *model.User) error {
	return s.db.Save(user).Error
}

func (s *UserService) UpdateStatus(id uint, status int) error {
	return s.db.Model(&model.User{}).Where("id = ?", id).Update("status", status).Error
}

func (s *UserService) UpdatePassword(id uint, password string) error {
	return s.db.Model(&model.User{}).Where("id = ?", id).Update("password", password).Error
}

func (s *UserService) UpdateLoginTime(id uint) error {
	return s.db.Model(&model.User{}).Where("id = ?", id).Update("login_time", gorm.Expr("NOW()")).Error
}
