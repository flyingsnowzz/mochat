package dashboard

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"mochat-api-server/internal/model"
	"mochat-api-server/internal/pkg/response"
	"mochat-api-server/internal/service"
)

type EmployeeHandler struct {
	svc     *service.WorkEmployeeService
	syncSvc *service.WorkAddressSyncService
}

func NewEmployeeHandler(db *gorm.DB) *EmployeeHandler {
	return &EmployeeHandler{
		svc:     service.NewWorkEmployeeService(db),
		syncSvc: service.NewWorkAddressSyncService(db),
	}
}

func (h *EmployeeHandler) Index(c *gin.Context) {
	corpID, _ := c.Get("corpId")
	if corpID == nil {
		corpID = uint(0)
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("perPage", c.DefaultQuery("pageSize", "10")))
	status, _ := strconv.Atoi(c.DefaultQuery("status", "0"))

	employees, total, err := h.svc.List(corpID.(uint), service.WorkEmployeeListFilter{
		Name:        c.Query("name"),
		Status:      status,
		ContactAuth: c.DefaultQuery("contactAuth", "all"),
	}, page, pageSize)
	if err != nil {
		response.Fail(c, response.ErrDB, "获取员工列表失败")
		return
	}
	items := make([]gin.H, 0, len(employees))
	for _, employee := range employees {
		items = append(items, gin.H{
			"id":                employee.ID,
			"name":              employee.Name,
			"thumbAvatar":       employee.ThumbAvatar,
			"status":            employee.Status,
			"statusName":        workEmployeeStatusName(employee.Status),
			"gender":            workEmployeeGenderName(employee.Gender),
			"contactAuth":       employee.ContactAuth,
			"contactAuthName":   workEmployeeContactAuthName(employee.ContactAuth),
			"applyNums":         0,
			"addNums":           0,
			"messageNums":       0,
			"sendMessageNums":   0,
			"replyMessageRatio": 0,
			"averageReply":      0,
			"invalidContact":    0,
		})
	}
	response.PageResult(c, items, total, page, pageSize)
}

func (h *EmployeeHandler) SearchCondition(c *gin.Context) {
	corpID, _ := c.Get("corpId")
	if corpID == nil {
		corpID = uint(0)
	}

	syncTime, err := h.svc.LastSyncTime(corpID.(uint))
	if err != nil {
		response.Fail(c, response.ErrDB, "获取同步信息失败")
		return
	}

	response.Success(c, gin.H{
		"syncTime": syncTime,
		"status": []gin.H{
			{"id": 1, "name": workEmployeeStatusName(1)},
			{"id": 2, "name": workEmployeeStatusName(2)},
			{"id": 4, "name": workEmployeeStatusName(4)},
			{"id": 5, "name": workEmployeeStatusName(5)},
		},
		"contactAuth": []gin.H{
			{"id": 1, "name": workEmployeeContactAuthName(1)},
			{"id": 2, "name": workEmployeeContactAuthName(2)},
		},
	})
}

func (h *EmployeeHandler) SyncEmployee(c *gin.Context) {
	corpID, _ := c.Get("corpId")
	if corpID == nil || corpID.(uint) == 0 {
		response.Fail(c, response.ErrParams, "请先选择企业")
		return
	}

	if err := h.syncSvc.SyncCorp(corpID.(uint)); err != nil {
		response.Fail(c, response.ErrServer, err.Error())
		return
	}
	response.SuccessMsg(c, "同步成功")
}

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

type DepartmentHandler struct {
	db  *gorm.DB
	svc *service.WorkDepartmentService
}

func NewDepartmentHandler(db *gorm.DB) *DepartmentHandler {
	return &DepartmentHandler{
		db:  db,
		svc: service.NewWorkDepartmentService(db),
	}
}

func (h *DepartmentHandler) Index(c *gin.Context) {
	corpID, _ := c.Get("corpId")
	if corpID == nil {
		corpID = uint(0)
	}
	departments, err := h.svc.List(corpID.(uint))
	if err != nil {
		response.Fail(c, response.ErrDB, "获取部门列表失败")
		return
	}

	searchKeywords := c.Query("searchKeyWords")
	employeeQuery := h.db.Model(&model.WorkEmployee{}).
		Where("corp_id = ? AND deleted_at IS NULL", corpID.(uint)).
		Order("updated_at DESC, id DESC")
	if searchKeywords != "" {
		employeeQuery = employeeQuery.Where("name LIKE ?", "%"+searchKeywords+"%")
	}

	var employees []model.WorkEmployee
	if err := employeeQuery.Find(&employees).Error; err != nil {
		response.Fail(c, response.ErrDB, "获取员工列表失败")
		return
	}

	departmentNodes := make([]gin.H, 0, len(departments))
	departmentMap := make(map[uint]*gin.H, len(departments))
	for _, department := range departments {
		node := gin.H{
			"id":           department.ID,
			"departmentId": department.ID,
			"name":         department.Name,
			"parentId":     department.ParentID,
			"wxDepartmentId": department.WxDepartmentID,
			"level":        department.Level,
			"son":          []gin.H{},
		}
		departmentNodes = append(departmentNodes, node)
		departmentMap[department.ID] = &departmentNodes[len(departmentNodes)-1]
	}

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

	employeeItems := make([]gin.H, 0, len(employees))
	for _, employee := range employees {
		employeeItems = append(employeeItems, gin.H{
			"id":           employee.ID,
			"employeeId":   employee.ID,
			"name":         employee.Name,
			"employeeName": employee.Name,
			"wxUserId":     employee.WxUserID,
			"avatar":       employee.Avatar,
			"thumbAvatar":  employee.ThumbAvatar,
			"mobile":       employee.Mobile,
			"departmentId": employee.MainDepartmentID,
		})
	}

	response.Success(c, gin.H{
		"department": departmentTree,
		"employee":   employeeItems,
	})
}

func (h *DepartmentHandler) PageIndex(c *gin.Context) {
	corpID, _ := c.Get("corpId")
	if corpID == nil {
		corpID = uint(0)
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	departments, total, err := h.svc.PageList(corpID.(uint), page, pageSize)
	if err != nil {
		response.Fail(c, response.ErrDB, "获取部门列表失败")
		return
	}
	response.PageResult(c, departments, total, page, pageSize)
}

func (h *DepartmentHandler) SelectByPhone(c *gin.Context) {
	response.Success(c, []interface{}{})
}

func (h *DepartmentHandler) ShowEmployee(c *gin.Context) {
	departmentID, _ := strconv.Atoi(c.DefaultQuery("departmentId", c.Param("id")))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("perPage", c.DefaultQuery("pageSize", "10")))
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	query := h.db.Model(&model.WorkEmployeeDepartment{}).
		Joins("JOIN mc_work_employee e ON e.id = mc_work_employee_department.employee_id").
		Where("e.deleted_at IS NULL")
	if departmentID > 0 {
		query = query.Where("mc_work_employee_department.department_id = ?", departmentID)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		response.Fail(c, response.ErrDB, "获取部门成员失败")
		return
	}

	type employeeRow struct {
		EmployeeID   uint   `json:"employeeId"`
		EmployeeName string `json:"employeeName"`
		Phone        string `json:"phone"`
		RoleName     string `json:"roleName"`
	}

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

	response.PageResult(c, rows, total, page, pageSize)
}

func (h *DepartmentHandler) DepartmentMemberIndex(c *gin.Context) {
	response.Success(c, gin.H{"list": []interface{}{}})
}
