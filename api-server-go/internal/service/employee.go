package service

import (
	"gorm.io/gorm"
	"mochat-api-server/internal/model"
)

// WorkEmployeeService 员工 Service
// 提供员工（WorkEmployee）的 CRUD 操作功能
// 主要职责：
// 1. 根据 ID 获取员工详情
// 2. 根据企业微信用户 ID 获取员工
// 3. 获取员工列表（分页，支持多条件筛选）
// 4. 获取员工最后同步时间
// 5. 创建员工
// 6. 更新员工
// 7. 删除员工
//
// 依赖：
// - gorm.DB: 数据库连接

type WorkEmployeeService struct {
	db *gorm.DB // 数据库连接
}

// WorkEmployeeListFilter 员工列表筛选条件
// 用于筛选员工列表的条件
type WorkEmployeeListFilter struct {
	Name        string // 员工姓名（模糊搜索）
	Status      int    // 员工状态
	ContactAuth string // 客户联系权限（"all" 表示全部）
}

// NewWorkEmployeeService 创建员工 Service 实例
// 参数：db - GORM 数据库连接
// 返回：员工 Service 实例
func NewWorkEmployeeService(db *gorm.DB) *WorkEmployeeService {
	return &WorkEmployeeService{db: db}
}

// GetByID 根据 ID 获取员工详情
// 查询指定 ID 的员工
// 参数：
//
//	id - 员工 ID
//
// 返回：员工实例和错误信息
func (s *WorkEmployeeService) GetByID(id uint) (*model.WorkEmployee, error) {
	var emp model.WorkEmployee
	if err := s.db.First(&emp, id).Error; err != nil {
		return nil, err
	}
	return &emp, nil
}

// GetByWxUserID 根据企业微信用户 ID 获取员工
// 使用企业微信用户 ID 查询员工
// 参数：
//
//	corpID - 企业 ID
//	wxUserID - 企业微信用户 ID
//
// 返回：员工实例和错误信息
func (s *WorkEmployeeService) GetByWxUserID(corpID uint, wxUserID string) (*model.WorkEmployee, error) {
	var emp model.WorkEmployee
	if err := s.db.Where("corp_id = ? AND wx_user_id = ?", corpID, wxUserID).First(&emp).Error; err != nil {
		return nil, err
	}
	return &emp, nil
}

// List 获取员工列表（分页）
// 查询指定企业的员工列表，支持分页和多条件筛选
// 参数：
//
//	corpID - 企业 ID
//	filter - 筛选条件（姓名、状态、客户联系权限）
//	page - 页码
//	pageSize - 每页数量
//
// 返回：员工列表、总数和错误信息
func (s *WorkEmployeeService) List(corpID uint, filter WorkEmployeeListFilter, page, pageSize int) ([]model.WorkEmployee, int64, error) {
	var employees []model.WorkEmployee
	var total int64
	query := s.db.Model(&model.WorkEmployee{}).Where("corp_id = ?", corpID)
	if filter.Name != "" {
		query = query.Where("name LIKE ?", "%"+filter.Name+"%")
	}
	if filter.Status > 0 {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.ContactAuth != "" && filter.ContactAuth != "all" {
		query = query.Where("contact_auth = ?", filter.ContactAuth)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * pageSize
	if err := query.Order("updated_at DESC, id DESC").Offset(offset).Limit(pageSize).Find(&employees).Error; err != nil {
		return nil, 0, err
	}
	return employees, total, nil
}

// LastSyncTime 获取员工最后同步时间
// 查询员工的最后同步时间
// 参数：
//
//	corpID - 企业 ID
//
// 返回：最后同步时间和错误信息（如果未找到同步记录，返回空字符串）
func (s *WorkEmployeeService) LastSyncTime(corpID uint) (string, error) {
	if corpID == 0 {
		return "", nil
	}

	var item model.WorkUpdateTime
	err := s.db.Where("corp_id = ? AND type = ?", corpID, workUpdateTimeTypeEmployee).
		Order("id DESC").
		First(&item).Error
	if err == nil {
		return item.LastUpdateTime, nil
	}
	if err == gorm.ErrRecordNotFound {
		return "", nil
	}
	return "", err
}

// Create 创建员工
// 将员工信息保存到数据库
// 参数：
//
//	emp - 员工实例
//
// 返回：错误信息
func (s *WorkEmployeeService) Create(emp *model.WorkEmployee) error {
	return s.db.Create(emp).Error
}

// Update 更新员工
// 更新员工的指定字段
// 参数：
//
//	id - 员工 ID
//	updates - 待更新的字段映射
//
// 返回：错误信息
func (s *WorkEmployeeService) Update(id uint, updates map[string]interface{}) error {
	return s.db.Model(&model.WorkEmployee{}).Where("id = ?", id).Updates(updates).Error
}

// Delete 删除员工
// 从数据库中删除指定 ID 的员工
// 参数：
//
//	id - 员工 ID
//
// 返回：错误信息
func (s *WorkEmployeeService) Delete(id uint) error {
	return s.db.Delete(&model.WorkEmployee{}, id).Error
}

// WorkDepartmentService 部门 Service
// 提供部门（WorkDepartment）的 CRUD 操作功能
// 主要职责：
// 1. 获取部门列表
// 2. 获取部门列表（分页）
// 3. 根据 ID 获取部门详情
// 4. 创建部门
// 5. 更新部门
// 6. 删除部门
//
// 依赖：
// - gorm.DB: 数据库连接

type WorkDepartmentService struct {
	db *gorm.DB // 数据库连接
}

// NewWorkDepartmentService 创建部门 Service 实例
// 参数：db - GORM 数据库连接
// 返回：部门 Service 实例
func NewWorkDepartmentService(db *gorm.DB) *WorkDepartmentService {
	return &WorkDepartmentService{db: db}
}

// List 获取部门列表
// 查询指定企业的所有部门，按排序字段升序排列
// 参数：
//
//	corpID - 企业 ID
//
// 返回：部门列表和错误信息
func (s *WorkDepartmentService) List(corpID uint) ([]model.WorkDepartment, error) {
	var departments []model.WorkDepartment
	if err := s.db.Where("corp_id = ?", corpID).Order("`order` ASC").Find(&departments).Error; err != nil {
		return nil, err
	}
	return departments, nil
}

// PageList 获取部门列表（分页）
// 查询指定企业的部门列表，支持分页
// 参数：
//
//	corpID - 企业 ID
//	page - 页码
//	pageSize - 每页数量
//
// 返回：部门列表、总数和错误信息
func (s *WorkDepartmentService) PageList(corpID uint, page, pageSize int) ([]model.WorkDepartment, int64, error) {
	var departments []model.WorkDepartment
	var total int64
	query := s.db.Model(&model.WorkDepartment{}).Where("corp_id = ?", corpID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Find(&departments).Error; err != nil {
		return nil, 0, err
	}
	return departments, total, nil
}

// GetByID 根据 ID 获取部门详情
// 查询指定 ID 的部门
// 参数：
//
//	id - 部门 ID
//
// 返回：部门实例和错误信息
func (s *WorkDepartmentService) GetByID(id uint) (*model.WorkDepartment, error) {
	var dept model.WorkDepartment
	if err := s.db.First(&dept, id).Error; err != nil {
		return nil, err
	}
	return &dept, nil
}

// Create 创建部门
// 将部门信息保存到数据库
// 参数：
//
//	dept - 部门实例
//
// 返回：错误信息
func (s *WorkDepartmentService) Create(dept *model.WorkDepartment) error {
	return s.db.Create(dept).Error
}

// Update 更新部门
// 更新部门的指定字段
// 参数：
//
//	id - 部门 ID
//	updates - 待更新的字段映射
//
// 返回：错误信息
func (s *WorkDepartmentService) Update(id uint, updates map[string]interface{}) error {
	return s.db.Model(&model.WorkDepartment{}).Where("id = ?", id).Updates(updates).Error
}

// Delete 删除部门
// 从数据库中删除指定 ID 的部门
// 参数：
//
//	id - 部门 ID
//
// 返回：错误信息
func (s *WorkDepartmentService) Delete(id uint) error {
	return s.db.Delete(&model.WorkDepartment{}, id).Error
}
