// Package dashboard 提供 Dashboard 相关的 HTTP 处理器
// 该文件包含员工和部门管理的处理器：
// 1. EmployeeHandler - 处理员工的查询、同步等操作
// 2. DepartmentHandler - 处理部门的查询等操作
package organization

import (
	"strconv"

	"mochat-api-server/internal/model"
	"mochat-api-server/internal/pkg/response"
	"mochat-api-server/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// EmployeeHandler 员工管理处理器
// 提供员工的查询、同步等功能
// 主要职责：
// 1. 获取员工列表
// 2. 获取员工搜索条件
// 3. 同步企业微信员工
//
// 依赖服务：
// - WorkEmployeeService: 员工服务
// - WorkAddressSyncService: 通讯录同步服务

type EmployeeHandler struct {
	svc     *service.WorkEmployeeService    // 员工服务
	syncSvc *service.WorkAddressSyncService // 通讯录同步服务
}

// NewEmployeeHandler 创建员工管理处理器实例
// 参数：db - GORM 数据库连接
// 返回：员工管理处理器实例
func NewEmployeeHandler(db *gorm.DB) *EmployeeHandler {
	return &EmployeeHandler{
		svc:     service.NewWorkEmployeeService(db),
		syncSvc: service.NewWorkAddressSyncService(db),
	}
}

// Index 获取员工列表
// 获取企业的员工列表，支持分页和筛选
// 处理流程：
// 1. 获取企业 ID
// 2. 获取分页参数和筛选条件
// 3. 调用服务层获取员工列表
// 4. 构建返回数据
// 5. 返回分页结果
// 参数：
//
//	page - 页码，默认为 1
//	perPage/pageSize - 每页数量，默认为 10
//	status - 状态，默认为 0（全部）
//	name - 员工姓名（用于搜索）
//	contactAuth - 客户联系权限，默认为 "all"（全部）
//
// 返回：包含员工列表、总数、分页信息的响应
func (h *EmployeeHandler) Index(c *gin.Context) {
	// 获取企业 ID
	corpID, exists := c.Get("corpId")
	if !exists {
		response.Fail(c, response.ErrAuth, "未获取到企业信息")
		return
	}

	// 获取分页参数和筛选条件
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("perPage", c.DefaultQuery("pageSize", "10")))
	status, _ := strconv.Atoi(c.DefaultQuery("status", "0"))

	// 调用服务层获取员工列表
	employees, total, err := h.svc.List(corpID.(uint), service.WorkEmployeeListFilter{
		Name:        c.Query("name"),
		Status:      status,
		ContactAuth: c.DefaultQuery("contactAuth", "all"),
	}, page, pageSize)
	if err != nil {
		response.Fail(c, response.ErrDB, "获取员工列表失败")
		return
	}

	// 构建返回数据
	items := make([]gin.H, 0, len(employees))
	for _, employee := range employees {
		items = append(items, gin.H{
			"id":                employee.ID,                                       // 员工 ID
			"name":              employee.Name,                                     // 员工姓名
			"thumbAvatar":       employee.ThumbAvatar,                              // 员工头像
			"status":            employee.Status,                                   // 员工状态
			"statusName":        workEmployeeStatusName(employee.Status),           // 员工状态名称
			"gender":            workEmployeeGenderName(employee.Gender),           // 员工性别
			"contactAuth":       employee.ContactAuth,                              // 客户联系权限
			"contactAuthName":   workEmployeeContactAuthName(employee.ContactAuth), // 客户联系权限名称
			"applyNums":         0,                                                 // 申请数量（暂未实现）
			"addNums":           0,                                                 // 添加数量（暂未实现）
			"messageNums":       0,                                                 // 消息数量（暂未实现）
			"sendMessageNums":   0,                                                 // 发送消息数量（暂未实现）
			"replyMessageRatio": 0,                                                 // 消息回复率（暂未实现）
			"averageReply":      0,                                                 // 平均回复时间（暂未实现）
			"invalidContact":    0,                                                 // 无效客户数量（暂未实现）
		})
	}

	// 返回分页结果
	response.PageResult(c, items, total, page, pageSize)
}

// SearchCondition 获取员工搜索条件
// 获取员工搜索相关的条件数据，包括同步时间、状态列表和客户联系权限列表
// 处理流程：
// 1. 获取企业 ID
// 2. 调用服务层获取最后同步时间
// 3. 构建返回数据，包括同步时间、状态列表和客户联系权限列表
// 4. 返回搜索条件数据
//
// 返回：包含同步时间、状态列表和客户联系权限列表的响应
func (h *EmployeeHandler) SearchCondition(c *gin.Context) {
	// 获取企业 ID
	corpID, _ := c.Get("corpId")
	if corpID == nil {
		corpID = uint(0)
	}

	// 调用服务层获取最后同步时间
	syncTime, err := h.svc.LastSyncTime(corpID.(uint))
	if err != nil {
		response.Fail(c, response.ErrDB, "获取同步信息失败")
		return
	}

	// 构建返回数据
	response.Success(c, gin.H{
		"syncTime": syncTime, // 最后同步时间
		"status": []gin.H{ // 状态列表
			{"id": 1, "name": workEmployeeStatusName(1)},
			{"id": 2, "name": workEmployeeStatusName(2)},
			{"id": 4, "name": workEmployeeStatusName(4)},
			{"id": 5, "name": workEmployeeStatusName(5)},
		},
		"contactAuth": []gin.H{ // 客户联系权限列表
			{"id": 1, "name": workEmployeeContactAuthName(1)},
			{"id": 2, "name": workEmployeeContactAuthName(2)},
		},
	})
}

// SyncEmployee 同步企业微信员工
// 同步企业微信通讯录中的员工信息到系统
// 处理流程：
// 1. 获取企业 ID
// 2. 验证企业 ID 是否有效
// 3. 调用服务层同步企业微信通讯录
// 4. 返回同步结果
//
// 返回：同步成功的消息
func (h *EmployeeHandler) SyncEmployee(c *gin.Context) {
	// 获取企业 ID
	corpID, _ := c.Get("corpId")
	if corpID == nil || corpID.(uint) == 0 {
		response.Fail(c, response.ErrParams, "请先选择企业")
		return
	}

	// 调用服务层同步企业微信通讯录
	if err := h.syncSvc.SyncCorp(corpID.(uint)); err != nil {
		response.Fail(c, response.ErrServer, err.Error())
		return
	}

	// 返回同步结果
	response.SuccessMsg(c, "同步成功")
}

// workEmployeeStatusName 获取员工状态名称
// 根据员工状态码获取对应的状态名称
// 参数：
//
//	status - 员工状态码
//
// 返回：员工状态名称
//
// 状态码对应关系：
// 1 - 已激活
// 2 - 已禁用
// 4 - 未激活
// 5 - 退出企业
func workEmployeeStatusName(status int) string {
	switch status {
	case 1:
		return "已激活"
	case 2:
		return "已禁用"
	case 4:
		return "未激活"
	case 5:
		return "退出企业"
	default:
		return ""
	}
}

// workEmployeeGenderName 获取员工性别名称
// 根据员工性别码获取对应的性别名称
// 参数：
//
//	gender - 员工性别码
//
// 返回：员工性别名称
//
// 性别码对应关系：
// 1 - 男
// 2 - 女
// 其他 - 未定义
func workEmployeeGenderName(gender int) string {
	switch gender {
	case 1:
		return "男"
	case 2:
		return "女"
	default:
		return "未定义"
	}
}

// workEmployeeContactAuthName 获取员工客户联系权限名称
// 根据员工客户联系权限码获取对应的权限名称
// 参数：
//
//	contactAuth - 员工客户联系权限码
//
// 返回：员工客户联系权限名称
//
// 权限码对应关系：
// 1 - 是
// 2 - 否
func workEmployeeContactAuthName(contactAuth int) string {
	switch contactAuth {
	case 1:
		return "是"
	case 2:
		return "否"
	default:
		return ""
	}
}

// DepartmentHandler 部门管理处理器
// 提供部门的查询等功能
// 主要职责：
// 1. 获取部门列表（树形结构）
// 2. 获取部门列表（分页）
// 3. 根据电话选择部门
// 4. 获取部门成员列表
//
// 依赖服务：
// - WorkDepartmentService: 部门服务
// - gorm.DB: 数据库连接

type DepartmentHandler struct {
	db  *gorm.DB                       // 数据库连接
	svc *service.WorkDepartmentService // 部门服务
}

// NewDepartmentHandler 创建部门管理处理器实例
// 参数：db - GORM 数据库连接
// 返回：部门管理处理器实例
func NewDepartmentHandler(db *gorm.DB) *DepartmentHandler {
	return &DepartmentHandler{
		db:  db,
		svc: service.NewWorkDepartmentService(db),
	}
}

// Index 获取部门列表（树形结构）
// 获取企业的部门列表，返回树形结构，并同时返回员工列表
// 处理流程：
// 1. 获取企业 ID
// 2. 调用服务层获取部门列表
// 3. 获取搜索关键词，构建员工查询
// 4. 查询员工列表
// 5. 构建部门树形结构
// 6. 构建员工列表
// 7. 返回部门树形结构和员工列表
// 参数：
//
//	searchKeyWords - 搜索关键词（用于搜索员工）
//
// 返回：包含部门树形结构和员工列表的响应
func (h *DepartmentHandler) Index(c *gin.Context) {
	// 获取企业 ID
	corpID, _ := c.Get("corpId")
	if corpID == nil {
		corpID = uint(0)
	}

	// 调用服务层获取部门列表
	departments, err := h.svc.List(corpID.(uint))
	if err != nil {
		response.Fail(c, response.ErrDB, "获取部门列表失败")
		return
	}

	// 获取搜索关键词，构建员工查询
	searchKeywords := c.Query("searchKeyWords")
	employeeQuery := h.db.Model(&model.WorkEmployee{}).
		Where("corp_id = ? AND deleted_at IS NULL", corpID.(uint)).
		Order("updated_at DESC, id DESC")
	if searchKeywords != "" {
		employeeQuery = employeeQuery.Where("name LIKE ?", "%"+searchKeywords+"%")
	}

	// 查询员工列表
	var employees []model.WorkEmployee
	if err := employeeQuery.Find(&employees).Error; err != nil {
		response.Fail(c, response.ErrDB, "获取员工列表失败")
		return
	}

	// 构建部门树形结构
	departmentNodes := make([]gin.H, 0, len(departments))
	departmentMap := make(map[uint]*gin.H, len(departments))
	for _, department := range departments {
		node := gin.H{
			"id":             department.ID,             // 部门 ID
			"departmentId":   department.ID,             // 部门 ID（兼容字段）
			"name":           department.Name,           // 部门名称
			"parentId":       department.ParentID,       // 父部门 ID
			"wxDepartmentId": department.WxDepartmentID, // 企业微信部门 ID
			"level":          department.Level,          // 部门级别
			"son":            []gin.H{},                 // 子部门列表
		}
		departmentNodes = append(departmentNodes, node)
		departmentMap[department.ID] = &departmentNodes[len(departmentNodes)-1]
	}

	// 构建部门树形结构
	departmentTree := make([]gin.H, 0)
	for i := range departmentNodes {
		node := &departmentNodes[i]
		parentID, _ := (*node)["parentId"].(uint)
		if parentID == 0 {
			departmentTree = append(departmentTree, *node)
			continue
		}
		parentNode, ok := departmentMap[parentID]
		if !ok {
			departmentTree = append(departmentTree, *node)
			continue
		}
		children, _ := (*parentNode)["son"].([]gin.H)
		children = append(children, *node)
		(*parentNode)["son"] = children
	}

	// 构建员工列表
	employeeItems := make([]gin.H, 0, len(employees))
	for _, employee := range employees {
		employeeItems = append(employeeItems, gin.H{
			"id":           employee.ID,               // 员工 ID
			"employeeId":   employee.ID,               // 员工 ID（兼容字段）
			"name":         employee.Name,             // 员工姓名
			"employeeName": employee.Name,             // 员工姓名（兼容字段）
			"wxUserId":     employee.WxUserID,         // 企业微信用户 ID
			"avatar":       employee.Avatar,           // 员工头像
			"thumbAvatar":  employee.ThumbAvatar,      // 员工头像（缩略图）
			"mobile":       employee.Mobile,           // 员工手机号
			"departmentId": employee.MainDepartmentID, // 员工主部门 ID
		})
	}

	// 返回部门树形结构和员工列表
	response.Success(c, gin.H{
		"department": departmentTree, // 部门树形结构
		"employee":   employeeItems,  // 员工列表
	})
}

// PageIndex 获取部门列表（分页）
// 获取企业的部门列表，支持分页
// 处理流程：
// 1. 获取企业 ID
// 2. 获取分页参数
// 3. 调用服务层获取部门列表
// 4. 返回分页结果
// 参数：
//
//	page - 页码，默认为 1
//	pageSize - 每页数量，默认为 20
//
// 返回：包含部门列表、总数、分页信息的响应
func (h *DepartmentHandler) PageIndex(c *gin.Context) {
	// 获取企业 ID
	corpID, _ := c.Get("corpId")
	if corpID == nil {
		corpID = uint(0)
	}

	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	// 调用服务层获取部门列表
	departments, total, err := h.svc.PageList(corpID.(uint), page, pageSize)
	if err != nil {
		response.Fail(c, response.ErrDB, "获取部门列表失败")
		return
	}

	// 返回分页结果
	response.PageResult(c, departments, total, page, pageSize)
}

// SelectByPhone 根据电话选择部门
// 根据电话获取部门信息（暂未实现）
// 处理流程：
// 1. 直接返回空数组
//
// 返回：空数组
func (h *DepartmentHandler) SelectByPhone(c *gin.Context) {
	response.Success(c, []interface{}{})
}

// ShowEmployee 获取部门成员列表
// 获取指定部门的成员列表，支持分页
// 处理流程：
// 1. 获取部门 ID
// 2. 获取分页参数
// 3. 构建查询，关联员工和部门关系表
// 4. 计算总数
// 5. 查询部门成员列表
// 6. 返回分页结果
// 参数：
//
//	departmentId - 部门 ID（查询参数或路径参数）
//	page - 页码，默认为 1
//	perPage/pageSize - 每页数量，默认为 10
//
// 返回：包含部门成员列表、总数、分页信息的响应
func (h *DepartmentHandler) ShowEmployee(c *gin.Context) {
	// 获取部门 ID
	departmentID, _ := strconv.Atoi(c.DefaultQuery("departmentId", c.Param("id")))

	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("perPage", c.DefaultQuery("pageSize", "10")))
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	// 构建查询，关联员工和部门关系表
	query := h.db.Model(&model.WorkEmployeeDepartment{}).
		Joins("JOIN mc_work_employee e ON e.id = mc_work_employee_department.employee_id").
		Where("e.deleted_at IS NULL")
	if departmentID > 0 {
		query = query.Where("mc_work_employee_department.department_id = ?", departmentID)
	}

	// 计算总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		response.Fail(c, response.ErrDB, "获取部门成员失败")
		return
	}

	// 定义返回结构
	type employeeRow struct {
		EmployeeID   uint   `json:"employeeId"`   // 员工 ID
		EmployeeName string `json:"employeeName"` // 员工姓名
		Phone        string `json:"phone"`        // 员工手机号
		RoleName     string `json:"roleName"`     // 员工职位
	}

	// 查询部门成员列表
	var rows []employeeRow
	offset := (page - 1) * pageSize
	if err := query.
		Select("e.id AS employee_id, e.name AS employee_name, e.mobile AS phone, e.position AS role_name").
		Order("mc_work_employee_department.is_leader_in_dept DESC, mc_work_employee_department.`order` ASC, e.id ASC").
		Offset(offset).
		Limit(pageSize).
		Scan(&rows).Error; err != nil {
		response.Fail(c, response.ErrDB, "获取部门成员失败")
		return
	}

	// 返回分页结果
	response.PageResult(c, rows, total, page, pageSize)
}

// DepartmentMemberIndex 获取部门成员列表
// 获取部门成员列表（暂未实现）
// 处理流程：
// 1. 直接返回空列表
//
// 返回：包含空列表的响应
func (h *DepartmentHandler) DepartmentMemberIndex(c *gin.Context) {
	response.Success(c, gin.H{"list": []interface{}{}})
}
