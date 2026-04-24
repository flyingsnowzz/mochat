package dashboard

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
	svc *service.WorkDepartmentService
}

func NewDepartmentHandler(db *gorm.DB) *DepartmentHandler {
	return &DepartmentHandler{svc: service.NewWorkDepartmentService(db)}
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
	response.Success(c, departments)
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
	response.Success(c, []interface{}{})
}

func (h *DepartmentHandler) DepartmentMemberIndex(c *gin.Context) {
	response.Success(c, gin.H{"list": []interface{}{}})
}
