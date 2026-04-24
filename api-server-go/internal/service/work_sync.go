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

const workUpdateTimeTypeEmployee = 1

type WorkAddressSyncService struct {
	db *gorm.DB
}

func NewWorkAddressSyncService(db *gorm.DB) *WorkAddressSyncService {
	return &WorkAddressSyncService{db: db}
}

func (s *WorkAddressSyncService) SyncCorp(corpID uint) error {
	var corp model.Corp
	if err := s.db.First(&corp, corpID).Error; err != nil {
		return err
	}
	if corp.WxCorpid == "" || corp.EmployeeSecret == "" {
		return errors.New("当前企业未配置企业微信通讯录同步凭证")
	}

	wechat.RefreshApp(corp.ID)
	addressClient := wechat.GetWorkApp(&corp).GetAddressList()

	departments, err := addressClient.DepartmentList()
	if err != nil {
		return fmt.Errorf("获取企业微信部门失败: %w", err)
	}

	followUsers := s.loadFollowUsers(&corp)
	userDetails, err := s.loadUsers(addressClient, departments)
	if err != nil {
		return err
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		deptMap, err := s.syncDepartments(tx, corp.ID, departments)
		if err != nil {
			return err
		}
		employeeIDMap, err := s.syncEmployees(tx, &corp, deptMap, userDetails, followUsers)
		if err != nil {
			return err
		}
		if err := s.replaceEmployeeDepartments(tx, deptMap, userDetails, employeeIDMap); err != nil {
			return err
		}
		if err := s.touchWorkUpdateTime(tx, corp.ID, workUpdateTimeTypeEmployee); err != nil {
			return err
		}
		return nil
	})
}

func (s *WorkAddressSyncService) loadFollowUsers(corp *model.Corp) map[string]struct{} {
	result := make(map[string]struct{})
	if corp.ContactSecret == "" {
		return result
	}

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

func (s *WorkAddressSyncService) loadUsers(client *addresslist.Client, departments []*addresslist.Department) ([]*addresslist.UserGetResponse, error) {
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

func (s *WorkAddressSyncService) syncDepartments(tx *gorm.DB, corpID uint, remoteDepartments []*addresslist.Department) (map[int]model.WorkDepartment, error) {
	existing := make([]model.WorkDepartment, 0)
	if err := tx.Where("corp_id = ?", corpID).Find(&existing).Error; err != nil {
		return nil, err
	}

	existingByWxID := make(map[int]model.WorkDepartment, len(existing))
	for _, department := range existing {
		existingByWxID[department.WxDepartmentID] = department
	}

	now := time.Now()
	for _, remote := range remoteDepartments {
		if remote == nil {
			continue
		}

		item, ok := existingByWxID[remote.ID]
		if ok {
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

	departments := make([]model.WorkDepartment, 0)
	if err := tx.Where("corp_id = ?", corpID).Order("id ASC").Find(&departments).Error; err != nil {
		return nil, err
	}

	deptMap := make(map[int]model.WorkDepartment, len(departments))
	for _, department := range departments {
		deptMap[department.WxDepartmentID] = department
	}

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

func (s *WorkAddressSyncService) syncEmployees(tx *gorm.DB, corp *model.Corp, deptMap map[int]model.WorkDepartment, users []*addresslist.UserGetResponse, followUsers map[string]struct{}) (map[string]uint, error) {
	wxUserIDs := make([]string, 0, len(users))
	for _, user := range users {
		if user == nil || user.UserID == "" {
			continue
		}
		wxUserIDs = append(wxUserIDs, user.UserID)
	}

	existing := make([]model.WorkEmployee, 0)
	if len(wxUserIDs) > 0 {
		if err := tx.Where("corp_id = ? AND wx_user_id IN ?", corp.ID, wxUserIDs).Find(&existing).Error; err != nil {
			return nil, err
		}
	}

	existingByWxID := make(map[string]model.WorkEmployee, len(existing))
	for _, employee := range existing {
		existingByWxID[employee.WxUserID] = employee
	}

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

		if existingEmployee, ok := existingByWxID[user.UserID]; ok {
			employee.ID = existingEmployee.ID
			employee.CreatedAt = existingEmployee.CreatedAt
			if err := tx.Save(&employee).Error; err != nil {
				return nil, err
			}
			employeeIDMap[user.UserID] = employee.ID
			continue
		}

		employee.CreatedAt = now
		if err := tx.Create(&employee).Error; err != nil {
			return nil, err
		}
		employeeIDMap[user.UserID] = employee.ID
	}

	return employeeIDMap, nil
}

func (s *WorkAddressSyncService) replaceEmployeeDepartments(tx *gorm.DB, deptMap map[int]model.WorkDepartment, users []*addresslist.UserGetResponse, employeeIDMap map[string]uint) error {
	if len(employeeIDMap) == 0 {
		return nil
	}

	employeeIDs := make([]uint, 0, len(employeeIDMap))
	seenEmployeeIDs := make(map[uint]struct{}, len(employeeIDMap))
	for _, id := range employeeIDMap {
		if _, ok := seenEmployeeIDs[id]; ok {
			continue
		}
		seenEmployeeIDs[id] = struct{}{}
		employeeIDs = append(employeeIDs, id)
	}

	if err := tx.Where("employee_id IN ?", employeeIDs).Delete(&model.WorkEmployeeDepartment{}).Error; err != nil {
		return err
	}

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

func (s *WorkAddressSyncService) touchWorkUpdateTime(tx *gorm.DB, corpID uint, syncType int) error {
	var item model.WorkUpdateTime
	now := time.Now().Format("2006-01-02 15:04:05")
	err := tx.Where("corp_id = ? AND type = ?", corpID, syncType).First(&item).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
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
	return tx.Model(&model.WorkUpdateTime{}).Where("id = ?", item.ID).Updates(map[string]interface{}{
		"last_update_time": now,
		"error_msg":        "",
		"updated_at":       time.Now(),
	}).Error
}

func (s *WorkAddressSyncService) loadTenantUserIDsByPhone(tx *gorm.DB, tenantID uint, users []*addresslist.UserGetResponse) (map[string]uint, error) {
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

	userModels := make([]model.User, 0)
	if err := tx.Select("id", "phone").Where("tenant_id = ? AND phone IN ?", tenantID, phones).Find(&userModels).Error; err != nil {
		return nil, err
	}
	for _, user := range userModels {
		result[user.Phone] = user.ID
	}
	return result, nil
}

func buildDepartmentRelation(wxDepartmentID int, deptMap map[int]model.WorkDepartment) (uint, int, string) {
	department, ok := deptMap[wxDepartmentID]
	if !ok {
		return 0, 0, ""
	}

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

	parentID := uint(0)
	if department.WxParentID > 0 {
		if parent, ok := deptMap[department.WxParentID]; ok {
			parentID = parent.ID
		}
	}

	path := ""
	for idx, id := range ids {
		if idx > 0 {
			path += "-"
		}
		path += "#" + strconv.FormatUint(uint64(id), 10) + "#"
	}
	return parentID, len(ids), path
}

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

func intAt(values []int, index int) int {
	if index < 0 || index >= len(values) {
		return 0
	}
	return values[index]
}

func departmentIDByWX(wxDepartmentID int, deptMap map[int]model.WorkDepartment) uint {
	if department, ok := deptMap[wxDepartmentID]; ok {
		return department.ID
	}
	return 0
}

func boolToContactAuth(yes bool) int {
	if yes {
		return 1
	}
	return 2
}
