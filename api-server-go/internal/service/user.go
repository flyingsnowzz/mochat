package service

import (
	"mochat-api-server/internal/model"
	"gorm.io/gorm"
)

// UserService 用户 Service
// 提供用户（User）的 CRUD 操作功能
// 主要职责：
// 1. 根据手机号获取用户
// 2. 根据 ID 获取用户详情
// 3. 获取用户列表（分页）
// 4. 创建用户
// 5. 更新用户
// 6. 更新用户状态
// 7. 更新用户密码
// 8. 更新用户登录时间
//
// 依赖：
// - gorm.DB: 数据库连接

type UserService struct {
	db *gorm.DB // 数据库连接
}

// NewUserService 创建用户 Service 实例
// 参数：db - GORM 数据库连接
// 返回：用户 Service 实例
func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
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

// List 获取用户列表（分页）
// 查询用户列表，支持按租户 ID 筛选
// 参数：
//
//	tenantID - 租户 ID，0 表示不限制
//	page - 页码
//	pageSize - 每页数量
//
// 返回：用户列表、总数和错误信息
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

// Create 创建用户
// 将用户信息保存到数据库
// 参数：
//
//	user - 用户实例
//
// 返回：错误信息
func (s *UserService) Create(user *model.User) error {
	return s.db.Create(user).Error
}

// Update 更新用户
// 更新数据库中的用户信息
// 参数：
//
//	user - 用户实例
//
// 返回：错误信息
func (s *UserService) Update(user *model.User) error {
	return s.db.Save(user).Error
}

// UpdateStatus 更新用户状态
// 更新用户的账户状态
// 参数：
//
//	id - 用户 ID
//	status - 状态
//
// 返回：错误信息
func (s *UserService) UpdateStatus(id uint, status int) error {
	return s.db.Model(&model.User{}).Where("id = ?", id).Update("status", status).Error
}

// UpdatePassword 更新用户密码
// 更新用户的登录密码
// 参数：
//
//	id - 用户 ID
//	password - 新密码
//
// 返回：错误信息
func (s *UserService) UpdatePassword(id uint, password string) error {
	return s.db.Model(&model.User{}).Where("id = ?", id).Update("password", password).Error
}

// UpdateLoginTime 更新用户登录时间
// 更新用户的最后登录时间
// 参数：
//
//	id - 用户 ID
//
// 返回：错误信息
func (s *UserService) UpdateLoginTime(id uint) error {
	return s.db.Model(&model.User{}).Where("id = ?", id).Update("login_time", gorm.Expr("NOW()")).Error
}
