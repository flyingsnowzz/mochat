package dashboard

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"mochat-api-server/internal/model"
	"mochat-api-server/internal/pkg/response"
	"mochat-api-server/internal/service"
)

type WorkContactTagGroupHandler struct {
	groupSvc *service.WorkContactTagGroupService
}

func NewWorkContactTagGroupHandler(db *gorm.DB) *WorkContactTagGroupHandler {
	return &WorkContactTagGroupHandler{
		groupSvc: service.NewWorkContactTagGroupService(db),
	}
}

func (h *WorkContactTagGroupHandler) Index(c *gin.Context) {
	corpID, _ := c.Get("corpId")
	if corpID == nil {
		response.Fail(c, response.ErrParams, "未获取到企业信息")
		return
	}

	groups, err := h.groupSvc.List(corpID.(uint))
	if err != nil {
		response.Fail(c, response.ErrDB, "获取标签组列表失败")
		return
	}

	response.Success(c, gin.H{"list": groups})
}

func (h *WorkContactTagGroupHandler) Detail(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	group, err := h.groupSvc.GetByID(uint(id))
	if err != nil {
		response.Fail(c, response.ErrNotFound, "标签组不存在")
		return
	}

	response.Success(c, group)
}

func (h *WorkContactTagGroupHandler) Store(c *gin.Context) {
	corpID, _ := c.Get("corpId")
	if corpID == nil {
		response.Fail(c, response.ErrParams, "未获取到企业信息")
		return
	}

	var group model.WorkContactTagGroup
	if err := c.ShouldBindJSON(&group); err != nil {
		response.Fail(c, response.ErrParams, "参数错误")
		return
	}

	group.CorpID = corpID.(uint)
	if err := h.groupSvc.Create(&group); err != nil {
		response.Fail(c, response.ErrDB, "创建标签组失败")
		return
	}

	response.Success(c, group)
}

func (h *WorkContactTagGroupHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	group, err := h.groupSvc.GetByID(uint(id))
	if err != nil {
		response.Fail(c, response.ErrNotFound, "标签组不存在")
		return
	}

	if err := c.ShouldBindJSON(group); err != nil {
		response.Fail(c, response.ErrParams, "参数错误")
		return
	}

	if err := h.groupSvc.Update(group); err != nil {
		response.Fail(c, response.ErrDB, "更新标签组失败")
		return
	}

	response.Success(c, group)
}

func (h *WorkContactTagGroupHandler) Destroy(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if err := h.groupSvc.Delete(uint(id)); err != nil {
		response.Fail(c, response.ErrDB, "删除标签组失败")
		return
	}

	response.SuccessMsg(c, "删除成功")
}
