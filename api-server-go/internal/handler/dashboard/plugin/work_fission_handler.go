package plugin

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"mochat-api-server/internal/pkg/response"
	"mochat-api-server/internal/service/plugin"
)

type WorkFissionHandler struct {
	svc *plugin.WorkFissionService
}

func NewWorkFissionHandler(svc *plugin.WorkFissionService) *WorkFissionHandler {
	return &WorkFissionHandler{svc: svc}
}

func (h *WorkFissionHandler) Index(c *gin.Context) {
	corpID, _ := c.Get("corpId")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	offset := (page - 1) * pageSize

	items, total, err := h.svc.List(corpID.(uint), offset, pageSize)
	if err != nil {
		response.Fail(c, response.ErrDB, "获取裂变列表失败")
		return
	}
	response.PageResult(c, items, total, page, pageSize)
}

func (h *WorkFissionHandler) Show(c *gin.Context) {
	response.Success(c, gin.H{})
}

func (h *WorkFissionHandler) Info(c *gin.Context) {
	response.Success(c, gin.H{})
}

func (h *WorkFissionHandler) Invite(c *gin.Context) {
	response.Success(c, gin.H{})
}

func (h *WorkFissionHandler) Store(c *gin.Context) {
	response.SuccessMsg(c, "创建成功")
}

func (h *WorkFissionHandler) Update(c *gin.Context) {
	response.SuccessMsg(c, "更新成功")
}