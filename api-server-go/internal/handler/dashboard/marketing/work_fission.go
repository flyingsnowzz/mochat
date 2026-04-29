package marketing

import (
	"mochat-api-server/internal/pkg/response"
	"mochat-api-server/internal/pkg/utils"
	"mochat-api-server/internal/service/plugin"
	"strconv"

	"github.com/gin-gonic/gin"
)

type WorkFissionHandler struct {
	service *plugin.WorkFissionService
}

func NewWorkFissionHandler(service *plugin.WorkFissionService) *WorkFissionHandler {
	return &WorkFissionHandler{service: service}
}

// ============ Dashboard 端 ============

func (h *WorkFissionHandler) Index(c *gin.Context) {
	corpID, exists := c.Get("corpId")
	if !exists {
		response.Fail(c, response.ErrParams, "请先选择企业")
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("perPage", c.DefaultQuery("pageSize", "10")))
	offset := (page - 1) * pageSize

	items, total, err := h.service.List(corpID.(uint), offset, pageSize)
	if err != nil {
		response.Fail(c, response.ErrDB, "获取任务宝列表失败")
		return
	}

	list := make([]gin.H, 0, len(items))
	for _, item := range items {
		status := utils.CalcFissionStatus(item.EndTime)
		list = append(list, gin.H{
			"id":                    item.ID,
			"corpId":                item.CorpID,
			"activeName":            item.ActiveName,
			"status":                status,
			"serviceEmployees":      utils.ParseJSONArray(item.ServiceEmployees),
			"autoPass":              item.AutoPass,
			"autoAddTag":            item.AutoAddTag,
			"contactTags":           utils.ParseJSONArray(item.ContactTags),
			"endTime":               utils.FormatTimePtr(item.EndTime),
			"qrCodeInvalid":         item.QrCodeInvalid,
			"tasks":                 utils.ParseJSONArray(item.Tasks),
			"newFriend":             item.NewFriend,
			"deleteInvalid":         item.DeleteInvalid,
			"receivePrize":          item.ReceivePrize,
			"receivePrizeEmployees": utils.ParseJSONArray(item.ReceivePrizeEmployees),
			"receiveLinks":          utils.ParseJSONArray(item.ReceiveLinks),
			"receiveQrcode":         utils.ParseJSONArray(item.ReceiveQrcode),
			"createUserId":          item.CreateUserID,
			"createdAt":             item.CreatedAt.Format("2006-01-02 15:04:05"),
			"updatedAt":             item.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	response.Success(c, gin.H{
		"list":  list,
		"total": total,
		"page":  page,
		"size":  pageSize,
	})
}

func (h *WorkFissionHandler) Show(c *gin.Context) {
	response.Success(c, nil)
}

func (h *WorkFissionHandler) Store(c *gin.Context) {
	response.SuccessMsg(c, "创建成功")
}

func (h *WorkFissionHandler) Update(c *gin.Context) {
	response.SuccessMsg(c, "更新成功")
}

func (h *WorkFissionHandler) Destroy(c *gin.Context) {
	response.SuccessMsg(c, "删除成功")
}

func (h *WorkFissionHandler) Statistics(c *gin.Context) {
	response.Success(c, gin.H{"list": []interface{}{}})
}

func (h *WorkFissionHandler) Info(c *gin.Context) {
	response.Success(c, nil)
}

func (h *WorkFissionHandler) Invite(c *gin.Context) {
	response.Success(c, gin.H{"list": []interface{}{}})
}

func (h *WorkFissionHandler) InviteData(c *gin.Context) {
	response.Success(c, nil)
}

func (h *WorkFissionHandler) InviteDetail(c *gin.Context) {
	response.Success(c, gin.H{"list": []interface{}{}})
}

func (h *WorkFissionHandler) ChooseContact(c *gin.Context) {
	response.Success(c, gin.H{"list": []interface{}{}})
}

// ============ Operation 端 ============

func (h *WorkFissionHandler) Auth(c *gin.Context) {
	response.Success(c, nil)
}

func (h *WorkFissionHandler) InviteFriends(c *gin.Context) {
	response.Success(c, nil)
}

func (h *WorkFissionHandler) TaskData(c *gin.Context) {
	response.Success(c, nil)
}

func (h *WorkFissionHandler) Receive(c *gin.Context) {
	response.Success(c, nil)
}

func (h *WorkFissionHandler) Poster(c *gin.Context) {
	response.Success(c, nil)
}

func (h *WorkFissionHandler) OpenUserInfo(c *gin.Context) {
	response.Success(c, nil)
}
