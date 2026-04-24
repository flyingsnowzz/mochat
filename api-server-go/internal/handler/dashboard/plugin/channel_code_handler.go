package plugin

import (
	"strings"
	"strconv"

	"github.com/gin-gonic/gin"
	"mochat-api-server/internal/model"
	"mochat-api-server/internal/pkg/response"
	"mochat-api-server/internal/service/plugin"
)

type ChannelCodeHandler struct {
	svc      *plugin.ChannelCodeService
	groupSvc *plugin.ChannelCodeGroupService
}

type channelCodeGroupItem struct {
	GroupID uint   `json:"groupId"`
	Name    string `json:"name"`
}

func NewChannelCodeHandler(svc *plugin.ChannelCodeService, groupSvc *plugin.ChannelCodeGroupService) *ChannelCodeHandler {
	return &ChannelCodeHandler{svc: svc, groupSvc: groupSvc}
}

func (h *ChannelCodeHandler) Index(c *gin.Context) {
	corpIDValue, ok := getCorpID(c)
	if !ok {
		response.Fail(c, response.ErrParams, "请先选择企业")
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	groupID, _ := strconv.Atoi(c.DefaultQuery("groupId", "0"))
	offset := (page - 1) * pageSize

	codes, total, err := h.svc.List(corpIDValue, uint(groupID), offset, pageSize)
	if err != nil {
		response.Fail(c, response.ErrDB, "获取渠道码列表失败")
		return
	}
	response.PageResult(c, codes, total, page, pageSize)
}

func (h *ChannelCodeHandler) Show(c *gin.Context) {
	response.Success(c, gin.H{})
}

func (h *ChannelCodeHandler) Store(c *gin.Context) {
	response.SuccessMsg(c, "创建成功")
}

func (h *ChannelCodeHandler) Update(c *gin.Context) {
	response.SuccessMsg(c, "更新成功")
}

func (h *ChannelCodeHandler) GroupIndex(c *gin.Context) {
	corpIDValue, ok := getCorpID(c)
	if !ok {
		response.Fail(c, response.ErrParams, "请先选择企业")
		return
	}

	groups, err := h.groupSvc.List(corpIDValue)
	if err != nil {
		response.Fail(c, response.ErrDB, "获取渠道码分组失败")
		return
	}

	data := make([]channelCodeGroupItem, 0, len(groups)+1)
	data = append(data, channelCodeGroupItem{GroupID: 0, Name: "未分组"})
	for _, group := range groups {
		data = append(data, channelCodeGroupItem{GroupID: group.ID, Name: group.Name})
	}
	response.Success(c, data)
}

func (h *ChannelCodeHandler) GroupStore(c *gin.Context) {
	corpIDValue, ok := getCorpID(c)
	if !ok {
		response.Fail(c, response.ErrParams, "请先选择企业")
		return
	}

	var req struct {
		Name any `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ErrParams, "参数错误")
		return
	}

	names := normalizeGroupNames(req.Name)
	if len(names) == 0 {
		response.Fail(c, response.ErrParams, "参数错误")
		return
	}

	for _, name := range names {
		group := &model.ChannelCodeGroup{
			CorpID: corpIDValue,
			Name:   name,
		}
		if err := h.groupSvc.Create(group); err != nil {
			response.Fail(c, response.ErrDB, "创建分组失败")
			return
		}
	}

	response.Success(c, gin.H{})
}

func (h *ChannelCodeHandler) GroupUpdate(c *gin.Context) {
	var req struct {
		GroupID uint   `json:"groupId" binding:"required"`
		Name    string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ErrParams, "参数错误")
		return
	}

	if err := h.groupSvc.Update(req.GroupID, map[string]interface{}{"name": strings.TrimSpace(req.Name)}); err != nil {
		response.Fail(c, response.ErrDB, "修改分组失败")
		return
	}

	response.Success(c, gin.H{})
}

func (h *ChannelCodeHandler) GroupMove(c *gin.Context) {
	var req struct {
		ChannelCodeID uint `json:"channelCodeId" binding:"required"`
		GroupID       uint `json:"groupId"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ErrParams, "参数错误")
		return
	}

	if err := h.svc.Update(req.ChannelCodeID, map[string]interface{}{"group_id": req.GroupID}); err != nil {
		response.Fail(c, response.ErrDB, "移动分组失败")
		return
	}

	response.Success(c, gin.H{})
}

func getCorpID(c *gin.Context) (uint, bool) {
	corpID, exists := c.Get("corpId")
	if !exists {
		return 0, false
	}
	corpIDValue, ok := corpID.(uint)
	if !ok || corpIDValue == 0 {
		return 0, false
	}
	return corpIDValue, true
}

func normalizeGroupNames(raw any) []string {
	var names []string
	switch value := raw.(type) {
	case string:
		name := strings.TrimSpace(value)
		if name != "" {
			names = append(names, name)
		}
	case []any:
		for _, item := range value {
			if str, ok := item.(string); ok {
				name := strings.TrimSpace(str)
				if name != "" {
					names = append(names, name)
				}
			}
		}
	case []string:
		for _, item := range value {
			name := strings.TrimSpace(item)
			if name != "" {
				names = append(names, name)
			}
		}
	}
	return names
}
