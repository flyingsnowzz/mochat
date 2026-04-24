package dashboard

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"mochat-api-server/internal/pkg/response"
	"mochat-api-server/internal/service"
)

type ContactFieldPivotHandler struct {
	pivotSvc *service.ContactFieldPivotService
}

func NewContactFieldPivotHandler(db *gorm.DB) *ContactFieldPivotHandler {
	return &ContactFieldPivotHandler{
		pivotSvc: service.NewContactFieldPivotService(db),
	}
}

func (h *ContactFieldPivotHandler) Index(c *gin.Context) {
	contactIDStr := c.Query("contactId")
	if contactIDStr == "" {
		response.Fail(c, response.ErrParams, "缺少 contactId 参数")
		return
	}

	contactID, err := strconv.ParseUint(contactIDStr, 10, 32)
	if err != nil {
		response.Fail(c, response.ErrParams, "contactId 参数错误")
		return
	}

	pivots, err := h.pivotSvc.List(uint(contactID))
	if err != nil {
		response.Fail(c, response.ErrDB, "获取字段值列表失败")
		return
	}

	response.Success(c, gin.H{"list": pivots})
}

func (h *ContactFieldPivotHandler) Update(c *gin.Context) {
	contactID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Fail(c, response.ErrParams, "参数错误")
		return
	}

	var req struct {
		Fields []struct {
			FieldID uint   `json:"fieldId"`
			Value   string `json:"value"`
		}
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ErrParams, "参数错误")
		return
	}

	if err := h.pivotSvc.BatchUpdate(uint(contactID), req.Fields); err != nil {
		response.Fail(c, response.ErrDB, "更新字段值失败")
		return
	}

	response.SuccessMsg(c, "更新成功")
}
