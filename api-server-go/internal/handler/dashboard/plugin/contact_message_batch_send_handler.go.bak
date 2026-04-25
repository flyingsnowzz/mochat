package plugin

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"mochat-api-server/internal/model"
	"mochat-api-server/internal/pkg/response"
	"mochat-api-server/internal/service"
	pluginService "mochat-api-server/internal/service/plugin"
)

type ContactMessageBatchSendHandler struct {
	db  *gorm.DB
	svc *pluginService.ContactMessageBatchSendService
}

func NewContactMessageBatchSendHandler(db *gorm.DB, svc *pluginService.ContactMessageBatchSendService) *ContactMessageBatchSendHandler {
	return &ContactMessageBatchSendHandler{db: db, svc: svc}
}

func (h *ContactMessageBatchSendHandler) Index(c *gin.Context) {
	corpID, _ := c.Get("corpId")
	if corpID == nil {
		corpID = uint(0)
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("perPage", "10"))
	if perPage <= 0 {
		perPage = 10
	}
	offset := (page - 1) * perPage

	items, total, err := h.svc.List(corpID.(uint), offset, perPage)
	if err != nil {
		response.Fail(c, response.ErrDB, "获取客户群发列表失败")
		return
	}

	list := make([]gin.H, 0, len(items))
	for _, item := range items {
		list = append(list, h.buildBatchListItem(item))
	}
	response.PageResult(c, list, total, page, perPage)
}

func (h *ContactMessageBatchSendHandler) Store(c *gin.Context) {
	var req struct {
		EmployeeIDs  []uint `json:"employeeIds" binding:"required"`
		FilterParams string `json:"filterParams"`
		Content      string `json:"content" binding:"required"`
		SendWay      string `json:"sendWay"`
		DefiniteTime string `json:"definiteTime"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ErrParams, "参数错误")
		return
	}

	corpID, _ := c.Get("corpId")
	userID, _ := c.Get("userId")
	if corpID == nil || corpID.(uint) == 0 {
		response.Fail(c, response.ErrParams, "请先选择企业")
		return
	}

	sendWay, _ := strconv.Atoi(req.SendWay)
	if sendWay == 0 {
		sendWay = 1
	}

	filterParams := strings.TrimSpace(req.FilterParams)
	if filterParams == "" {
		filterParams = "{}"
	}

	item := model.ContactMessageBatchSend{
		CorpID:             corpID.(uint),
		UserID:             asUint(userID),
		UserName:           h.loadUserName(asUint(userID)),
		EmployeeIDs:        marshalUintJSON(req.EmployeeIDs),
		FilterParams:       filterParams,
		FilterParamsDetail: buildFilterParamsDetail(filterParams),
		Content:            req.Content,
		SendWay:            sendWay,
		SendStatus:         0,
		NotSendTotal:       len(req.EmployeeIDs),
	}
	if sendWay == 1 {
		now := time.Now()
		item.DefiniteTime = &now
	} else if definiteTime, err := time.ParseInLocation("2006-01-02 15:04:05", req.DefiniteTime, time.Local); err == nil {
		item.DefiniteTime = &definiteTime
	}

	if err := h.svc.Create(&item); err != nil {
		response.Fail(c, response.ErrDB, "创建客户群发失败")
		return
	}
	response.Success(c, gin.H{"id": item.ID})
}

func (h *ContactMessageBatchSendHandler) Show(c *gin.Context) {
	id := getBatchID(c)
	item, err := h.svc.GetByID(uint(id))
	if err != nil {
		response.Fail(c, response.ErrNotFound, "客户群发不存在")
		return
	}
	response.Success(c, h.buildBatchDetailItem(*item))
}

func (h *ContactMessageBatchSendHandler) Destroy(c *gin.Context) {
	id := getBatchID(c)
	if id == 0 {
		response.Fail(c, response.ErrParams, "参数错误")
		return
	}
	if err := h.db.Delete(&model.ContactMessageBatchSend{}, id).Error; err != nil {
		response.Fail(c, response.ErrDB, "删除客户群发失败")
		return
	}
	response.SuccessMsg(c, "删除成功")
}

func (h *ContactMessageBatchSendHandler) MessageShow(c *gin.Context) {
	id := getBatchID(c)
	item, err := h.svc.GetByID(uint(id))
	if err != nil {
		response.Fail(c, response.ErrNotFound, "客户群发不存在")
		return
	}
	response.Success(c, parseJSONList(item.Content))
}

func (h *ContactMessageBatchSendHandler) Remind(c *gin.Context) {
	response.SuccessMsg(c, "提醒成功")
}

func (h *ContactMessageBatchSendHandler) EmployeeSendIndex(c *gin.Context) {
	id := getBatchID(c)
	item, err := h.svc.GetByID(uint(id))
	if err != nil {
		response.Fail(c, response.ErrNotFound, "客户群发不存在")
		return
	}

	sendStatus, _ := strconv.Atoi(c.DefaultQuery("sendStatus", "0"))
	keywords := strings.TrimSpace(c.DefaultQuery("keyWords", ""))
	if keywords == "" {
		keywords = strings.TrimSpace(c.DefaultQuery("keywords", ""))
	}

	employees := h.batchEmployees(item.EmployeeIDs)
	list := make([]gin.H, 0, len(employees))
	for _, employee := range employees {
		status := 0
		if item.SendStatus == 1 {
			status = 1
		}
		if sendStatus != status {
			continue
		}
		if keywords != "" && !strings.Contains(employee.Name, keywords) {
			continue
		}
		list = append(list, gin.H{
			"employeeId":       employee.ID,
			"employeeName":     employee.Name,
			"employeeAvatar":   employee.Avatar,
			"status":           status,
			"sendContactTotal": item.SendContactTotal,
			"sendTime":         formatTimePtr(item.SendTime),
		})
	}
	response.Success(c, gin.H{"list": list})
}

func (h *ContactMessageBatchSendHandler) ContactReceiveIndex(c *gin.Context) {
	id := getBatchID(c)
	_, err := h.svc.GetByID(uint(id))
	if err != nil {
		response.Fail(c, response.ErrNotFound, "客户群发不存在")
		return
	}

	sendStatus, _ := strconv.Atoi(c.DefaultQuery("sendStatus", "1"))
	keywords := strings.TrimSpace(c.DefaultQuery("keyWords", ""))
	var rows []model.ContactMessageBatchSendResult
	query := h.db.Model(&model.ContactMessageBatchSendResult{}).Where("batch_id = ?", id)
	if sendStatus >= 0 {
		query = query.Where("status = ?", sendStatus)
	}
	if err := query.Order("id DESC").Find(&rows).Error; err != nil {
		response.Fail(c, response.ErrDB, "获取客户详情失败")
		return
	}

	contactIDs := make([]uint, 0, len(rows))
	for _, row := range rows {
		if row.ContactID > 0 {
			contactIDs = append(contactIDs, row.ContactID)
		}
	}
	contactMap := h.loadContacts(contactIDs)

	list := make([]gin.H, 0, len(rows))
	for _, row := range rows {
		contact := contactMap[row.ContactID]
		name := contact.Name
		if name == "" {
			name = row.ExternalUserID
		}
		if keywords != "" && !strings.Contains(name, keywords) {
			continue
		}
		list = append(list, gin.H{
			"contactId":     row.ContactID,
			"contactName":   name,
			"contactAvatar": contact.Avatar,
			"status":        row.Status,
			"sendTime":      formatTimePtr(row.SendTime),
		})
	}
	response.Success(c, gin.H{"list": list})
}

func (h *ContactMessageBatchSendHandler) buildBatchListItem(item model.ContactMessageBatchSend) gin.H {
	return gin.H{
		"id":               item.ID,
		"content":          parseJSONList(item.Content),
		"sendWay":          item.SendWay,
		"definiteTime":     formatTimePtr(item.DefiniteTime),
		"sendTotal":        item.SendTotal,
		"receivedTotal":    item.ReceivedTotal,
		"notSendTotal":     item.NotSendTotal,
		"notReceivedTotal": item.NotReceivedTotal,
		"createdAt":        item.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

func (h *ContactMessageBatchSendHandler) buildBatchDetailItem(item model.ContactMessageBatchSend) gin.H {
	filterParams := parseJSONObject(item.FilterParams)
	filterDetail := parseJSONObject(item.FilterParamsDetail)
	if len(filterDetail) == 0 {
		filterDetail = gin.H{
			"rooms": []interface{}{},
			"tags":  []interface{}{},
		}
	}
	return gin.H{
		"id":                 item.ID,
		"creator":            item.UserName,
		"createdAt":          item.CreatedAt.Format("2006-01-02 15:04:05"),
		"content":            parseJSONList(item.Content),
		"filterParams":       filterParams,
		"filterParamsDetail": filterDetail,
		"sendTotal":          item.SendTotal,
		"receivedTotal":      item.ReceivedTotal,
		"notSendTotal":       item.NotSendTotal,
		"notReceivedTotal":   item.NotReceivedTotal,
		"receiveLimitTotal":  item.ReceiveLimitTotal,
		"notFriendTotal":     item.NotFriendTotal,
	}
}

func (h *ContactMessageBatchSendHandler) loadUserName(userID uint) string {
	if userID == 0 {
		return ""
	}
	user, err := service.NewUserService(h.db).GetByID(userID)
	if err != nil {
		return ""
	}
	return user.Name
}

func (h *ContactMessageBatchSendHandler) batchEmployees(raw string) []model.WorkEmployee {
	ids := parseUintSlice(raw)
	if len(ids) == 0 {
		return []model.WorkEmployee{}
	}
	var employees []model.WorkEmployee
	if err := h.db.Where("id IN ?", ids).Find(&employees).Error; err != nil {
		return []model.WorkEmployee{}
	}
	return employees
}

func (h *ContactMessageBatchSendHandler) loadContacts(ids []uint) map[uint]model.WorkContact {
	result := make(map[uint]model.WorkContact)
	if len(ids) == 0 {
		return result
	}
	var contacts []model.WorkContact
	if err := h.db.Where("id IN ?", ids).Find(&contacts).Error; err != nil {
		return result
	}
	for _, contact := range contacts {
		result[contact.ID] = contact
	}
	return result
}

func parseJSONList(raw string) []gin.H {
	if strings.TrimSpace(raw) == "" {
		return []gin.H{}
	}
	var result []gin.H
	if err := json.Unmarshal([]byte(raw), &result); err != nil {
		return []gin.H{}
	}
	return result
}

func parseJSONObject(raw string) gin.H {
	if strings.TrimSpace(raw) == "" || raw == "{}" {
		return gin.H{}
	}
	var result gin.H
	if err := json.Unmarshal([]byte(raw), &result); err != nil {
		return gin.H{}
	}
	return result
}

func buildFilterParamsDetail(filterParams string) string {
	params := parseJSONObject(filterParams)
	if len(params) == 0 {
		return `{"rooms":[],"tags":[]}`
	}
	if _, ok := params["rooms"]; !ok {
		params["rooms"] = []interface{}{}
	}
	if _, ok := params["tags"]; !ok {
		params["tags"] = []interface{}{}
	}
	if _, ok := params["excludeContacts"]; !ok {
		params["excludeContacts"] = []interface{}{}
	}
	data, err := json.Marshal(params)
	if err != nil {
		return `{"rooms":[],"tags":[]}`
	}
	return string(data)
}

func parseUintSlice(raw string) []uint {
	if strings.TrimSpace(raw) == "" {
		return []uint{}
	}
	var ids []uint
	if err := json.Unmarshal([]byte(raw), &ids); err == nil {
		return ids
	}
	return []uint{}
}

func marshalUintJSON(ids []uint) string {
	if len(ids) == 0 {
		return "[]"
	}
	data, err := json.Marshal(ids)
	if err != nil {
		return "[]"
	}
	return string(data)
}

func formatTimePtr(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}

func getBatchID(c *gin.Context) uint64 {
	if val := c.Query("batchId"); val != "" {
		id, _ := strconv.ParseUint(val, 10, 32)
		return id
	}
	if val := c.Param("id"); val != "" {
		id, _ := strconv.ParseUint(val, 10, 32)
		return id
	}
	var body struct {
		BatchID uint `json:"batchId"`
	}
	if err := c.ShouldBindJSON(&body); err == nil {
		return uint64(body.BatchID)
	}
	return 0
}

func asUint(value interface{}) uint {
	if value == nil {
		return 0
	}
	if v, ok := value.(uint); ok {
		return v
	}
	return 0
}
