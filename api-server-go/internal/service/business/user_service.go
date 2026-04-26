// Package business 提供业务逻辑层服务
// 该目录包含核心业务服务：
// 1. CorpService - 企业服务
// 2. UserService - 用户服务
package business

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"mochat-api-server/internal/model"
)

// UserService 用户服务
// 提供用户的 CRUD 操作和认证功能
// 主要职责：
// 1. 根据 ID 获取用户详情
// 2. 根据手机号获取用户
// 3. 获取用户列表（使用 offset 方式分页）
// 4. 创建用户（密码加密存储）
// 5. 更新用户（密码加密存储）
// 6. 更新用户登录时间
// 7. 验证用户密码
// 8. 删除用户
//
// 依赖：
// - gorm.DB: 数据库连接

type UserService struct {
	db *gorm.DB // 数据库连接
}

// NewUserService 创建用户服务实例
// 参数：db - GORM 数据库连接
// 返回：用户服务实例
func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

// GetByID 根据 ID 获取用户详情
// 查询指定 ID 的用户
// 参数：
//
//	id - 用户 ID
//
// 返回：用户实例和错误信息
func (s *UserService) GetByID(id uint) (*model.User, error) {
	var user model.User
	if err := s.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByPhone 根据手机号获取用户
// 使用手机号查询用户
// 参数：
//
//	phone - 手机号
//
// 返回：用户实例和错误信息
func (s *UserService) GetByPhone(phone string) (*model.User, error) {
	var user model.User
	if err := s.db.Where("phone = ?", phone).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// List 获取用户列表（使用 offset 方式分页）
// 查询用户列表，支持按租户 ID 筛选
// 参数：
//
//	tenantID - 租户 ID，0 表示不限制
//	offset - 偏移量
//	limit - 限制数量
//
// 返回：用户列表、总数和错误信息
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

// Create 创建用户
// 将用户信息保存到数据库，密码使用 bcrypt 加密存储
// 参数：
//
//	user - 用户实例
//
// 返回：错误信息
func (s *UserService) Create(user *model.User) error {
	if user.Password != "" {
		// 使用 bcrypt 对密码进行加密
		hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user.Password = string(hash)
	}
	return s.db.Create(user).Error
}

// Update 更新用户
// 更新用户的指定字段，密码使用 bcrypt 加密存储
// 参数：
//
//	id - 用户 ID
//	updates - 待更新的字段映射
//
// 返回：错误信息
func (s *UserService) Update(id uint, updates map[string]interface{}) error {
	if pwd, ok := updates["password"].(string); ok && pwd != "" {
		// 使用 bcrypt 对密码进行加密
		hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		updates["password"] = string(hash)
	}
	return s.db.Model(&model.User{}).Where("id = ?", id).Updates(updates).Error
}

// UpdateLoginTime 更新用户登录时间
// 更新用户的最后登录时间
// 参数：
//
//	id - 用户 ID
//
// 返回：错误信息
func (s *UserService) UpdateLoginTime(id uint) error {
	now := time.Now()
	return s.db.Model(&model.User{}).Where("id = ?", id).Updates(map[string]interface{}{
		"login_time": now,
	}).Error
}

// VerifyPassword 验证用户密码
// 验证用户的手机号和密码是否匹配，同时检查用户状态
// 参数：
//
//	phone - 手机号
//	password - 密码
//
// 返回：验证成功的用户实例或错误信息
// 可能的错误：用户不存在、密码错误、账户已被禁用
func (s *UserService) VerifyPassword(phone, password string) (*model.User, error) {
	// 获取用户
	user, err := s.GetByPhone(phone)
	if err != nil {
		return nil, err
	}
	
	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("密码错误")
	}
	
	// 检查用户状态
	if user.Status != 1 {
		return nil, errors.New("账户已被禁用")
	}
	return user, nil
}

// Delete 删除用户
// 从数据库中删除指定 ID 的用户
// 参数：
//
//	id - 用户 ID
//
// 返回：错误信息
func (s *UserService) Delete(id uint) error {
	return s.db.Delete(&model.User{}, id).Error
}
