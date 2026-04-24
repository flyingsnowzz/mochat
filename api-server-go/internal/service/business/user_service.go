package business

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"mochat-api-server/internal/model"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

func (s *UserService) GetByID(id uint) (*model.User, error) {
	var user model.User
	if err := s.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserService) GetByPhone(phone string) (*model.User, error) {
	var user model.User
	if err := s.db.Where("phone = ?", phone).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserService) List(tenantID uint, offset, limit int) ([]model.User, int64, error) {
	var users []model.User
	var total int64
	query := s.db.Model(&model.User{})
	if tenantID > 0 {
		query = query.Where("tenant_id = ?", tenantID)
	}
	query.Count(&total)
	if err := query.Offset(offset).Limit(limit).Order("id DESC").Find(&users).Error; err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

func (s *UserService) Create(user *model.User) error {
	if user.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user.Password = string(hash)
	}
	return s.db.Create(user).Error
}

func (s *UserService) Update(id uint, updates map[string]interface{}) error {
	if pwd, ok := updates["password"].(string); ok && pwd != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		updates["password"] = string(hash)
	}
	return s.db.Model(&model.User{}).Where("id = ?", id).Updates(updates).Error
}

func (s *UserService) UpdateLoginTime(id uint) error {
	now := time.Now()
	return s.db.Model(&model.User{}).Where("id = ?", id).Updates(map[string]interface{}{
		"login_time": now,
	}).Error
}

func (s *UserService) VerifyPassword(phone, password string) (*model.User, error) {
	user, err := s.GetByPhone(phone)
	if err != nil {
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("密码错误")
	}
	if user.Status != 1 {
		return nil, errors.New("账户已被禁用")
	}
	return user, nil
}

func (s *UserService) Delete(id uint) error {
	return s.db.Delete(&model.User{}, id).Error
}