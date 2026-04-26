package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/silenceper/wechat/v2/work/addresslist"
	"gorm.io/gorm"

	"mochat-api-server/internal/model"
	"mochat-api-server/internal/pkg/wechat"
)

// 同步类型常量
const workUpdateTimeTypeEmployee = 1 // 员工同步类型

// WorkAddressSyncService 企业微信通讯录同步服务
// 提供企业微信通讯录的同步功能，包括部门和员工的同步
// 主要职责：
// 1. 同步企业的部门信息
// 2. 同步企业的员工信息
// 3. 管理员工与部门的关联关系
// 4. 记录同步时间
//
// 依赖：
// - gorm.DB: 数据库连接

type WorkAddressSyncService struct {
	db *gorm.DB // 数据库连接
}

// NewWorkAddressSyncService 创建企业微信通讯录同步服务实例
// 参数：db - GORM 数据库连接
// 返回：企业微信通讯录同步服务实例
func NewWorkAddressSyncService(db *gorm.DB) *WorkAddressSyncService {
	return &WorkAddressSyncService{db: db}
}

// SyncCorp 同步企业通讯录
// 从企业微信同步部门、员工和员工部门关联信息
// 处理流程：
// 1. 获取企业信息
// 2. 验证企业是否配置了通讯录同步凭证
// 3. 刷新企业微信应用
// 4. 获取企业微信部门列表
// 5. 加载跟进用户列表
// 6. 加载用户详情
// 7. 同步部门信息
// 8. 同步员工信息
// 9. 替换员工部门关联
// 10. 更新同步时间
// 参数：
//
//	corpID - 企业 ID
//
// 返回：错误信息
func (s *WorkAddressSyncService) SyncCorp(corpID uint) error {
	// 获取企业信息
	var corp model.Corp
	if err := s.db.First(&corp, corpID).Error; err != nil {
		return err
	}
	
	// 验证企业是否配置了通讯录同步凭证
	if corp.WxCorpid == "" || corp.EmployeeSecret == "" {
		return errors.New("当前企业未配置企业微信通讯录同步凭证")
	}

	// 刷新企业微信应用
	wechat.RefreshApp(corp.ID)
	addressClient := wechat.GetWorkApp(&corp).GetAddressList()

	// 获取企业微信部门列表
	departments, err := addressClient.DepartmentList()
	if err != nil {
		return fmt.Errorf("获取企业微信部门失败: %w", err)
	}

	// 加载跟进用户列表
	followUsers := s.loadFollowUsers(&corp)
	
	// 加载用户详情
	userDetails, err := s.loadUsers(addressClient, departments)
	if err != nil {
		return err
	}

	// 在事务中执行同步操作
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 同步部门信息
		deptMap, err := s.syncDepartments(tx, corp.ID, departments)
		if err != nil {
			return err
		}
		// 同步员工信息
		employeeIDMap, err := s.syncEmployees(tx, &corp, deptMap, userDetails, followUsers)
		if err != nil {
			return err
		}
		// 替换员工部门关联
		if err := s.replaceEmployeeDepartments(tx, deptMap, userDetails, employeeIDMap); err != nil {
			return err
		}
		// 更新同步时间
		if err := s.touchWorkUpdateTime(tx, corp.ID, workUpdateTimeTypeEmployee); err != nil {
			return err
		}
		return nil
	})
}

// loadFollowUsers 加载跟进用户列表
// 从企业微信获取配置了客户联系功能的员工 ID 列表
// 参数：
//
//	corp - 企业实例
//
// 返回：用户 ID 到空结构体的映射
func (s *WorkAddressSyncService) loadFollowUsers(corp *model.Corp) map[string]struct{} {
	result := make(map[string]struct{})
	// 如果企业未配置客户联系密钥，直接返回空结果
	if corp.ContactSecret == "" {
		return result
	}

	// 获取跟进用户列表
	users, err := wechat.GetContactApp(corp).GetExternalContact().GetFollowUserList()
	if err != nil {
		return result
	}
	for _, userID := range users {
		if userID == "" {
			continue
		}
		result[userID] = struct{}{}
	}
	return result
}

// loadUsers 加载用户详情
// 从企业微信获取所有部门的用户详情
// 处理流程：
// 1. 遍历所有部门，获取每个部门的用户 ID 列表
// 2. 根据用户 ID 获取每个用户的详情
// 参数：
//
//	client - 企业微信通讯录客户端
//	departments - 部门列表
//
// 返回：用户详情列表和错误信息
func (s *WorkAddressSyncService) loadUsers(client *addresslist.Client, departments []*addresslist.Department) ([]*addresslist.UserGetResponse, error) {
	// 收集所有用户 ID
	userIDs := make(map[string]struct{})
	for _, department := range departments {
		if department == nil {
			continue
		}
		users, err := client.UserSimpleList(department.ID)
		if err != nil {
			return nil, fmt.Errorf("读取部门 %d 成员失败: %w", department.ID, err)
		}
		for _, user := range users {
			if user == nil || user.UserID == "" {
				continue
			}
			userIDs[user.UserID] = struct{}{}
		}
	}

	// 获取每个用户的详情
	details := make([]*addresslist.UserGetResponse, 0, len(userIDs))
	for userID := range userIDs {
		detail, err := client.UserGet(userID)
		if err != nil {
			return nil, fmt.Errorf("读取成员 %s 详情失败: %w", userID, err)
		}
		details = append(details, detail)
	}
	return details, nil
}

// syncDepartments 同步部门信息
// 将企业微信的部门信息同步到本地数据库
// 处理流程：
// 1. 获取现有的部门列表
// 2. 按企业微信部门 ID 建立映射
// 3. 遍历远程部门，更新或创建部门
// 4. 重新计算部门的父子关系和层级路径
// 参数：
//
//	tx - 数据库事务
//	corpID - 企业 ID
//	remoteDepartments - 远程部门列表
//
// 返回：企业微信部门 ID 到本地部门的映射和错误信息
func (s *WorkAddressSyncService) syncDepartments(tx *gorm.DB, corpID uint, remoteDepartments []*addresslist.Department) (map[int]model.WorkDepartment, error) {
	// 获取现有的部门列表
	existing := make([]model.WorkDepartment, 0)
	if err := tx.Where("corp_id = ?", corpID).Find(&existing).Error; err != nil {
		return nil, err
	}

	// 按企业微信部门 ID 建立映射
	existingByWxID := make(map[int]model.WorkDepartment, len(existing))
	for _, department := range existing {
		existingByWxID[department.WxDepartmentID] = department
	}

	now := time.Now()
	// 遍历远程部门，更新或创建
	for _, remote := range remoteDepartments {
		if remote == nil {
			continue
		}

		item, ok := existingByWxID[remote.ID]
		if ok {
			// 更新现有部门
			item.Name = remote.Name
			item.WxParentID = remote.ParentID
			item.Order = remote.Order
			item.UpdatedAt = now
			if err := tx.Save(&item).Error; err != nil {
				return nil, err
			}
			existingByWxID[remote.ID] = item
			continue
		}

		// 创建新部门
		item = model.WorkDepartment{
			WxDepartmentID: remote.ID,
			CorpID:         corpID,
			Name:           remote.Name,
			WxParentID:     remote.ParentID,
			Order:          remote.Order,
			CreatedAt:      now,
			UpdatedAt:      now,
		}
		if err := tx.Create(&item).Error; err != nil {
			return nil, err
		}
		existingByWxID[remote.ID] = item
	}

	// 重新获取所有部门
	departments := make([]model.WorkDepartment, 0)
	if err := tx.Where("corp_id = ?", corpID).Order("id ASC").Find(&departments).Error; err != nil {
		return nil, err
	}

	// 重新建立映射
	deptMap := make(map[int]model.WorkDepartment, len(departments))
	for _, department := range departments {
		deptMap[department.WxDepartmentID] = department
	}

	// 计算并更新部门的父子关系和层级路径
	for _, department := range departments {
		parentID, level, path := buildDepartmentRelation(department.WxDepartmentID, deptMap)
		if err := tx.Model(&model.WorkDepartment{}).
			Where("id = ?", department.ID).
			Updates(map[string]interface{}{
				"parent_id":  parentID,
				"level":      level,
				"path":       path,
				"updated_at": now,
			}).Error; err != nil {
			return nil, err
		}
	}

	// 重新获取所有部门并建立最终映射
	refreshed := make([]model.WorkDepartment, 0)
	if err := tx.Where("corp_id = ?", corpID).Find(&refreshed).Error; err != nil {
		return nil, err
	}
	deptMap = make(map[int]model.WorkDepartment, len(refreshed))
	for _, department := range refreshed {
		deptMap[department.WxDepartmentID] = department
	}
	return deptMap, nil
}

// syncEmployees 同步员工信息
// 将企业微信的员工信息同步到本地数据库
// 处理流程：
// 1. 收集所有用户 ID
// 2. 获取现有的员工列表
// 3. 加载租户用户 ID（按手机号关联）
// 4. 遍历远程用户，更新或创建员工
// 参数：
//
//	tx - 数据库事务
//	corp - 企业实例
//	deptMap - 企业微信部门 ID 到本地部门的映射
//	users - 远程用户列表
//	followUsers - 跟进用户映射
//
// 返回：企业微信用户 ID 到本地员工 ID 的映射和错误信息
func (s *WorkAddressSyncService) syncEmployees(tx *gorm.DB, corp *model.Corp, deptMap map[int]model.WorkDepartment, users []*addresslist.UserGetResponse, followUsers map[string]struct{}) (map[string]uint, error) {
	// 收集所有用户 ID
	wxUserIDs := make([]string, 0, len(users))
	for _, user := range users {
		if user == nil || user.UserID == "" {
			continue
		}
		wxUserIDs = append(wxUserIDs, user.UserID)
	}

	// 获取现有的员工列表
	existing := make([]model.WorkEmployee, 0)
	if len(wxUserIDs) > 0 {
		if err := tx.Where("corp_id = ? AND wx_user_id IN ?", corp.ID, wxUserIDs).Find(&existing).Error; err != nil {
			return nil, err
		}
	}

	// 按企业微信用户 ID 建立映射
	existingByWxID := make(map[string]model.WorkEmployee, len(existing))
	for _, employee := range existing {
		existingByWxID[employee.WxUserID] = employee
	}

	// 加载租户用户 ID（按手机号关联）
	userIDsByPhone, err := s.loadTenantUserIDsByPhone(tx, corp.TenantID, users)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	employeeIDMap := make(map[string]uint, len(users))
	for _, user := range users {
		if user == nil || user.UserID == "" {
			continue
		}

		// 判断是否有客户联系权限
		_, hasFollowAuth := followUsers[user.UserID]
		employee := model.WorkEmployee{
			WxUserID:           user.UserID,
			CorpID:             corp.ID,
			Name:               user.Name,
			Mobile:             user.Mobile,
			Position:           user.Position,
			Gender:             parseInt(user.Gender),
			Email:              user.Email,
			Avatar:             user.Avatar,
			ThumbAvatar:        user.ThumbAvatar,
			Telephone:          user.Telephone,
			Alias:              user.Alias,
			Extattr:            marshalJSON(user.Extattr),
			Status:             user.Status,
			QrCode:             user.QrCode,
			ExternalProfile:    marshalJSON(user.ExternalProfile),
			ExternalPosition:   user.ExternalPosition,
			Address:            user.Address,
			OpenUserID:         user.OpenUserid,
			WxMainDepartmentID: user.MainDepartment,
			MainDepartmentID:   departmentIDByWX(user.MainDepartment, deptMap),
			LogUserID:          userIDsByPhone[user.Mobile],
			ContactAuth:        boolToContactAuth(hasFollowAuth),
			UpdatedAt:          now,
		}

		// 如果已存在，更新
		if existingEmployee, ok := existingByWxID[user.UserID]; ok {
			employee.ID = existingEmployee.ID
			employee.CreatedAt = existingEmployee.CreatedAt
			if err := tx.Save(&employee).Error; err != nil {
				return nil, err
			}
			employeeIDMap[user.UserID] = employee.ID
			continue
		}

		// 创建新员工
		employee.CreatedAt = now
		if err := tx.Create(&employee).Error; err != nil {
			return nil, err
		}
		employeeIDMap[user.UserID] = employee.ID
	}

	return employeeIDMap, nil
}

// replaceEmployeeDepartments 替换员工部门关联
// 更新员工的部门关联关系
// 处理流程：
// 1. 删除现有的员工部门关联
// 2. 创建新的员工部门关联
// 参数：
//
//	tx - 数据库事务
//	deptMap - 企业微信部门 ID 到本地部门的映射
//	users - 远程用户列表
//	employeeIDMap - 企业微信用户 ID 到本地员工 ID 的映射
//
// 返回：错误信息
func (s *WorkAddressSyncService) replaceEmployeeDepartments(tx *gorm.DB, deptMap map[int]model.WorkDepartment, users []*addresslist.UserGetResponse, employeeIDMap map[string]uint) error {
	if len(employeeIDMap) == 0 {
		return nil
	}

	// 收集员工 ID 列表
	employeeIDs := make([]uint, 0, len(employeeIDMap))
	seenEmployeeIDs := make(map[uint]struct{}, len(employeeIDMap))
	for _, id := range employeeIDMap {
		if _, ok := seenEmployeeIDs[id]; ok {
			continue
		}
		seenEmployeeIDs[id] = struct{}{}
		employeeIDs = append(employeeIDs, id)
	}

	// 删除现有的员工部门关联
	if err := tx.Where("employee_id IN ?", employeeIDs).Delete(&model.WorkEmployeeDepartment{}).Error; err != nil {
		return err
	}

	// 创建新的员工部门关联
	relations := make([]model.WorkEmployeeDepartment, 0)
	now := time.Now()
	for _, user := range users {
		if user == nil || user.UserID == "" {
			continue
		}

		employeeID := employeeIDMap[user.UserID]
		if employeeID == 0 {
			continue
		}

		// 遍历用户所属的部门
		for idx, wxDepartmentID := range user.Department {
			department, ok := deptMap[wxDepartmentID]
			if !ok {
				continue
			}
			relations = append(relations, model.WorkEmployeeDepartment{
				EmployeeID:     employeeID,
				DepartmentID:   department.ID,
				IsLeaderInDept: intAt(user.IsLeaderInDept, idx),
				Order:          intAt(user.Order, idx),
				CreatedAt:      now,
				UpdatedAt:      now,
			})
		}
	}

	if len(relations) == 0 {
		return nil
	}
	return tx.Create(&relations).Error
}

// touchWorkUpdateTime 更新同步时间
// 记录或更新同步时间
// 参数：
//
//	tx - 数据库事务
//	corpID - 企业 ID
//	syncType - 同步类型
//
// 返回：错误信息
func (s *WorkAddressSyncService) touchWorkUpdateTime(tx *gorm.DB, corpID uint, syncType int) error {
	var item model.WorkUpdateTime
	now := time.Now().Format("2006-01-02 15:04:05")
	err := tx.Where("corp_id = ? AND type = ?", corpID, syncType).First(&item).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 创建新记录
		return tx.Create(&model.WorkUpdateTime{
			CorpID:         corpID,
			Type:           syncType,
			LastUpdateTime: now,
			ErrorMsg:       "",
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}).Error
	}
	if err != nil {
		return err
	}
	// 更新现有记录
	return tx.Model(&model.WorkUpdateTime{}).Where("id = ?", item.ID).Updates(map[string]interface{}{
		"last_update_time": now,
		"error_msg":        "",
		"updated_at":       time.Now(),
	}).Error
}

// loadTenantUserIDsByPhone 按手机号加载租户用户 ID
// 建立手机号到用户 ID 的映射，用于关联员工和用户
// 参数：
//
//	tx - 数据库事务
//	tenantID - 租户 ID
//	users - 用户列表
//
// 返回：手机号到用户 ID 的映射和错误信息
func (s *WorkAddressSyncService) loadTenantUserIDsByPhone(tx *gorm.DB, tenantID uint, users []*addresslist.UserGetResponse) (map[string]uint, error) {
	// 收集所有手机号
	phones := make([]string, 0)
	seen := make(map[string]struct{})
	for _, user := range users {
		if user == nil || user.Mobile == "" {
			continue
		}
		if _, ok := seen[user.Mobile]; ok {
			continue
		}
		seen[user.Mobile] = struct{}{}
		phones = append(phones, user.Mobile)
	}

	result := make(map[string]uint, len(phones))
	if len(phones) == 0 {
		return result, nil
	}

	// 查询用户
	userModels := make([]model.User, 0)
	if err := tx.Select("id", "phone").Where("tenant_id = ? AND phone IN ?", tenantID, phones).Find(&userModels).Error; err != nil {
		return nil, err
	}
	for _, user := range userModels {
		result[user.Phone] = user.ID
	}
	return result, nil
}

// buildDepartmentRelation 构建部门关系
// 根据企业微信部门 ID 计算父部门 ID、层级和路径
// 参数：
//
//	wxDepartmentID - 企业微信部门 ID
//	deptMap - 企业微信部门 ID 到本地部门的映射
//
// 返回：父部门 ID、层级和路径
func buildDepartmentRelation(wxDepartmentID int, deptMap map[int]model.WorkDepartment) (uint, int, string) {
	department, ok := deptMap[wxDepartmentID]
	if !ok {
		return 0, 0, ""
	}

	// 构建部门 ID 链
	ids := make([]uint, 0, 6)
	current := department
	for {
		ids = append([]uint{current.ID}, ids...)
		parent, hasParent := deptMap[current.WxParentID]
		if current.WxParentID == 0 || !hasParent {
			break
		}
		current = parent
	}

	// 计算父部门 ID
	parentID := uint(0)
	if department.WxParentID > 0 {
		if parent, ok := deptMap[department.WxParentID]; ok {
			parentID = parent.ID
		}
	}

	// 构建路径
	path := ""
	for idx, id := range ids {
		if idx > 0 {
			path += "-"
		}
		path += "#" + strconv.FormatUint(uint64(id), 10) + "#"
	}
	return parentID, len(ids), path
}

// marshalJSON 序列化 JSON
// 将对象序列化为 JSON 字符串
// 参数：
//
//	v - 待序列化的对象
//
// 返回：JSON 字符串
func marshalJSON(v interface{}) string {
	if v == nil {
		return "[]"
	}
	data, err := json.Marshal(v)
	if err != nil {
		return "[]"
	}
	if string(data) == "null" {
		return "[]"
	}
	return string(data)
}

// parseInt 解析整数
// 将字符串解析为整数
// 参数：
//
//	value - 字符串值
//
// 返回：整数值
func parseInt(value string) int {
	if value == "" {
		return 0
	}
	num, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}
	return num
}

// intAt 安全获取整数数组成员
// 从整数数组中安全获取指定索引的值
// 参数：
//
//	values - 整数数组
//	index - 索引
//
// 返回：指定索引的值，如果索引无效返回 0
func intAt(values []int, index int) int {
	if index < 0 || index >= len(values) {
		return 0
	}
	return values[index]
}

// departmentIDByWX 根据企业微信部门 ID 获取本地部门 ID
// 参数：
//
//	wxDepartmentID - 企业微信部门 ID
//	deptMap - 企业微信部门 ID 到本地部门的映射
//
// 返回：本地部门 ID，如果不存在返回 0
func departmentIDByWX(wxDepartmentID int, deptMap map[int]model.WorkDepartment) uint {
	if department, ok := deptMap[wxDepartmentID]; ok {
		return department.ID
	}
	return 0
}

// boolToContactAuth 将布尔值转换为客户联系权限
// 参数：
//
//	yes - 是否有客户联系权限
//
// 返回：1（有权限）或 2（无权限）
func boolToContactAuth(yes bool) int {
	if yes {
		return 1
	}
	return 2
}
