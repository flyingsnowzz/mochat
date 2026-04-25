package dashboard

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"mochat-api-server/internal/model"
	"mochat-api-server/internal/pkg/response"
	"mochat-api-server/internal/service"
)

type ContactHandler struct {
	db         *gorm.DB
	contactSvc *service.WorkContactService
	tagSvc     *service.WorkContactTagService
	fieldSvc   *service.ContactFieldService
}

func NewContactHandler(db *gorm.DB) *ContactHandler {
	return &ContactHandler{
		db:         db,
		contactSvc: service.NewWorkContactService(db),
		tagSvc:     service.NewWorkContactTagService(db),
		fieldSvc:   service.NewContactFieldService(db),
	}
}

func (h *ContactHandler) Index(c *gin.Context) {
	corpID, _ := c.Get("corpId")
	if corpID == nil {
		corpID = uint(0)
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", c.DefaultQuery("perPage", "20")))

	contacts, total, err := h.contactSvc.List(corpID.(uint), page, pageSize, nil)
	if err != nil {
		response.Fail(c, response.ErrDB, "获取客户列表失败")
		return
	}
	list, err := h.buildContactList(corpID.(uint), contacts)
	if err != nil {
		response.Fail(c, response.ErrDB, "获取客户列表失败")
		return
	}
	response.PageResult(c, list, total, page, pageSize)
}

func (h *ContactHandler) Show(c *gin.Context) {
	id, ok := parseContactID(c)
	if !ok {
		response.Fail(c, response.ErrParams, "参数错误")
		return
	}
	contact, err := h.contactSvc.GetByID(uint(id))
	if err != nil {
		response.Fail(c, response.ErrNotFound, "客户不存在")
		return
	}
	response.Success(c, contact)
}

func (h *ContactHandler) Update(c *gin.Context) {
	id, ok := parseContactID(c)
	if !ok {
		response.Fail(c, response.ErrParams, "参数错误")
		return
	}
	contact, err := h.contactSvc.GetByID(uint(id))
	if err != nil {
		response.Fail(c, response.ErrNotFound, "客户不存在")
		return
	}
	if err := c.ShouldBindJSON(contact); err != nil {
		response.Fail(c, response.ErrParams, "参数错误")
		return
	}
	if err := h.contactSvc.Update(contact); err != nil {
		response.Fail(c, response.ErrDB, "更新客户失败")
		return
	}
	response.Success(c, contact)
}

func (h *ContactHandler) SynContact(c *gin.Context) {
	corpID, _ := c.Get("corpId")
	if corpID == nil {
		response.Fail(c, response.ErrParams, "未获取到企业信息")
		return
	}

	// 异步提交同步任务
	go func() {
		// 这里可以实现与企业微信的同步逻辑
		// 1. 获取企业微信客户端
		// 2. 拉取客户列表
		// 3. 同步到数据库
	}()

	response.SuccessMsg(c, "同步任务已提交")
}

func (h *ContactHandler) Track(c *gin.Context) {
	id, ok := parseContactID(c)
	if !ok {
		response.Fail(c, response.ErrParams, "参数错误")
		return
	}
	
	// 获取客户轨迹记录
	var tracks []model.ContactEmployeeTrack
	result := h.db.Where("contact_id = ?", uint(id)).Order("created_at DESC").Find(&tracks)
	if result.Error != nil {
		response.Fail(c, response.ErrDB, "获取客户轨迹失败")
		return
	}

	response.Success(c, tracks)
}

func (h *ContactHandler) LossContact(c *gin.Context) {
	corpID, _ := c.Get("corpId")
	if corpID == nil {
		corpID = uint(0)
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	// 查询已删除的客户
	var contacts []model.WorkContact
	var total int64

	query := h.db.Model(&model.WorkContact{}).Where("corp_id = ? AND deleted_at IS NOT NULL", corpID)
	if err := query.Count(&total).Error; err != nil {
		response.Fail(c, response.ErrDB, "获取流失客户列表失败")
		return
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("deleted_at DESC").Find(&contacts).Error; err != nil {
		response.Fail(c, response.ErrDB, "获取流失客户列表失败")
		return
	}

	response.PageResult(c, contacts, total, page, pageSize)
}

func (h *ContactHandler) Source(c *gin.Context) {
	corpID, _ := c.Get("corpId")
	if corpID == nil {
		corpID = uint(0)
	}

	// 统计客户来源
	type SourceStat struct {
		AddWay    int   `json:"addWay"`
		Count     int64 `json:"count"`
		AddWayStr string `json:"addWayStr"`
	}

	var stats []SourceStat

	// 定义添加方式映射
	addWayMap := map[int]string{
		0: "未知",
		1: "扫描二维码",
		2: "搜索手机号",
		3: "名片分享",
		4: "群聊",
		5: "手机通讯录",
		6: "微信好友",
		7: "来自微信的添加",
		8: "安装第三方应用",
		9: "搜索邮箱",
		10: "企业微信内部成员共享",
		11: "管理员/负责人分配",
	}

	// 查询各来源的客户数量
	rows, err := h.db.Model(&model.WorkContactEmployee{}).Select("add_way, count(*) as count").Where("corp_id = ?", corpID).Group("add_way").Rows()
	if err != nil {
		response.Fail(c, response.ErrDB, "获取客户来源统计失败")
		return
	}
	defer rows.Close()

	for rows.Next() {
		var stat SourceStat
		if err := rows.Scan(&stat.AddWay, &stat.Count); err != nil {
			response.Fail(c, response.ErrDB, "获取客户来源统计失败")
			return
		}
		stat.AddWayStr = addWayMap[stat.AddWay]
		stats = append(stats, stat)
	}

	response.Success(c, stats)
}

func (h *ContactHandler) BatchLabeling(c *gin.Context) {
	var req struct {
		ContactIDs []uint `json:"contactIds"`
		TagIDs     []uint `json:"tagIds"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ErrParams, "参数错误")
		return
	}

	// 批量打标签
	err := h.db.Transaction(func(tx *gorm.DB) error {
		// 先删除已有的标签关联
		for _, contactID := range req.ContactIDs {
			if err := tx.Where("contact_id = ?", contactID).Delete(&model.WorkContactTagPivot{}).Error; err != nil {
				return err
			}
		}

		// 添加新的标签关联
		for _, contactID := range req.ContactIDs {
			for _, tagID := range req.TagIDs {
				pivot := &model.WorkContactTagPivot{
					ContactID:    contactID,
					ContactTagID: tagID,
				}
				if err := tx.Create(pivot).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})

	if err != nil {
		response.Fail(c, response.ErrDB, "批量打标签失败")
		return
	}

	response.SuccessMsg(c, "批量打标签成功")
}

func (h *ContactHandler) UpdateByID(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ErrParams, err.Error())
		return
	}
	if err := h.contactSvc.UpdateByID(uint(id), req); err != nil {
		response.Fail(c, response.ErrDB, "更新客户失败")
		return
	}
	response.SuccessMsg(c, "更新成功")
}

func (h *ContactHandler) GetByWxExternalUserID(c *gin.Context) {
	corpID, _ := c.Get("corpId")
	wxID := c.Query("wxExternalUserid")
	if wxID == "" {
		response.Fail(c, response.ErrParams, "wxExternalUserid参数不能为空")
		return
	}
	contact, err := h.contactSvc.GetByWxExternalUserID(corpID.(uint), wxID)
	if err != nil {
		response.Fail(c, response.ErrNotFound, "客户不存在")
		return
	}
	response.Success(c, contact)
}

func (h *ContactHandler) buildContactList(corpID uint, contacts []model.WorkContact) ([]gin.H, error) {
	if len(contacts) == 0 {
		return []gin.H{}, nil
	}

	contactIDs := make([]uint, 0, len(contacts))
	for _, contact := range contacts {
		contactIDs = append(contactIDs, contact.ID)
	}

	employeeMap, err := h.workEmployeeNameMap(corpID)
	if err != nil {
		return nil, err
	}
	relations, err := h.latestContactEmployeeMap(corpID, contactIDs)
	if err != nil {
		return nil, err
	}
	roomNames, err := h.contactRoomNames(contactIDs)
	if err != nil {
		return nil, err
	}
	tagNames, err := h.contactTagNames(contactIDs)
	if err != nil {
		return nil, err
	}

	list := make([]gin.H, 0, len(contacts))
	for _, contact := range contacts {
		row := relations[contact.ID]
		employeeName := ""
		employeeID := uint(0)
		addWay := 0
		createTime := ""
		if row != nil {
			employeeID = row.EmployeeID
			addWay = row.AddWay
			createTime = row.CreateTime.Format("2006-01-02 15:04:05")
			employeeName = employeeMap[row.EmployeeID]
		}

		list = append(list, gin.H{
			"id":          contact.ID,
			"name":        firstNonEmpty(strings.TrimSpace(contact.Name), strings.TrimSpace(contact.NickName)),
			"avatar":      contact.Avatar,
			"employeeName": employeeName,
			"roomName":    defaultStringSlice(roomNames[contact.ID]),
			"addWayText":  contactAddWayText(addWay),
			"tag":         defaultStringSlice(tagNames[contact.ID]),
			"createTime":  createTime,
			"contactId":   contact.ID,
			"employeeId":  employeeID,
			"isContact":   2,
			"businessNo":  contact.BusinessNo,
		})
	}
	return list, nil
}

func (h *ContactHandler) latestContactEmployeeMap(corpID uint, contactIDs []uint) (map[uint]*model.WorkContactEmployee, error) {
	var rows []model.WorkContactEmployee
	if err := h.db.Where("corp_id = ? AND contact_id IN ?", corpID, contactIDs).
		Order("create_time DESC, id DESC").
		Find(&rows).Error; err != nil {
		return nil, err
	}
	result := make(map[uint]*model.WorkContactEmployee, len(contactIDs))
	for i := range rows {
		if _, ok := result[rows[i].ContactID]; ok {
			continue
		}
		result[rows[i].ContactID] = &rows[i]
	}
	return result, nil
}

func (h *ContactHandler) workEmployeeNameMap(corpID uint) (map[uint]string, error) {
	var rows []model.WorkEmployee
	if err := h.db.Where("corp_id = ? AND deleted_at IS NULL", corpID).Find(&rows).Error; err != nil {
		return nil, err
	}
	result := make(map[uint]string, len(rows))
	for _, row := range rows {
		result[row.ID] = row.Name
	}
	return result, nil
}

func (h *ContactHandler) contactRoomNames(contactIDs []uint) (map[uint][]string, error) {
	var rows []struct {
		ContactID uint   `json:"contactId"`
		RoomName  string `json:"roomName"`
	}
	if err := h.db.Table("mc_work_contact_room AS wcr").
		Select("wcr.contact_id AS contact_id, wr.name AS room_name").
		Joins("LEFT JOIN mc_work_room AS wr ON wr.id = wcr.room_id").
		Where("wcr.contact_id IN ?", contactIDs).
		Where("wr.deleted_at IS NULL").
		Scan(&rows).Error; err != nil {
		return nil, err
	}
	result := make(map[uint][]string, len(contactIDs))
	for _, row := range rows {
		if strings.TrimSpace(row.RoomName) == "" {
			continue
		}
		result[row.ContactID] = appendIfMissing(result[row.ContactID], row.RoomName)
	}
	return result, nil
}

func (h *ContactHandler) contactTagNames(contactIDs []uint) (map[uint][]string, error) {
	var rows []struct {
		ContactID uint   `json:"contactId"`
		TagName   string `json:"tagName"`
	}
	if err := h.db.Table("mc_work_contact_tag_pivot AS pivot").
		Select("pivot.contact_id AS contact_id, tag.name AS tag_name").
		Joins("LEFT JOIN mc_work_contact_tag AS tag ON tag.id = pivot.contact_tag_id").
		Where("pivot.contact_id IN ?", contactIDs).
		Scan(&rows).Error; err != nil {
		return nil, err
	}
	result := make(map[uint][]string, len(contactIDs))
	for _, row := range rows {
		if strings.TrimSpace(row.TagName) == "" {
			continue
		}
		result[row.ContactID] = appendIfMissing(result[row.ContactID], row.TagName)
	}
	return result, nil
}

func defaultStringSlice(items []string) []string {
	if items == nil {
		return []string{}
	}
	return items
}

func appendIfMissing(items []string, value string) []string {
	for _, item := range items {
		if item == value {
			return items
		}
	}
	return append(items, value)
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}

func contactAddWayText(addWay int) string {
	switch addWay {
	case 1:
		return "扫描二维码"
	case 2:
		return "搜索手机号"
	case 3:
		return "名片分享"
	case 4:
		return "群聊"
	case 5:
		return "手机通讯录"
	case 6:
		return "微信好友"
	case 7:
		return "来自微信的添加"
	case 8:
		return "安装第三方应用"
	case 9:
		return "搜索邮箱"
	case 10:
		return "企业微信内部成员共享"
	case 11:
		return "管理员/负责人分配"
	default:
		return "未知"
	}
}

func splitContactFieldOptions(raw string) []string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return []string{}
	}
	if strings.HasPrefix(raw, "[") {
		var arr []string
		if err := json.Unmarshal([]byte(raw), &arr); err == nil {
			return arr
		}
	}
	parts := strings.FieldsFunc(raw, func(r rune) bool {
		return r == ',' || r == '\n' || r == '\r'
	})
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		result = append(result, part)
	}
	return result
}

func parseContactID(c *gin.Context) (uint64, bool) {
	for _, raw := range []string{c.Param("id"), c.Query("contactId"), c.Query("id")} {
		raw = strings.TrimSpace(raw)
		if raw == "" || raw == "undefined" || raw == "null" {
			continue
		}
		id, err := strconv.ParseUint(raw, 10, 32)
		if err == nil && id > 0 {
			return id, true
		}
	}
	return 0, false
}

func parseIDFromParamOrBody(c *gin.Context) (uint, bool) {
	if raw := strings.TrimSpace(c.Param("id")); raw != "" {
		if id, err := strconv.ParseUint(raw, 10, 32); err == nil && id > 0 {
			return uint(id), true
		}
	}
	var req struct {
		ID uint `json:"id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		return 0, false
	}
	if req.ID == 0 {
		return 0, false
	}
	return req.ID, true
}

func uintFromAny(v interface{}) uint {
	switch val := v.(type) {
	case float64:
		if val < 0 {
			return 0
		}
		return uint(val)
	case float32:
		if val < 0 {
			return 0
		}
		return uint(val)
	case int:
		if val < 0 {
			return 0
		}
		return uint(val)
	case int64:
		if val < 0 {
			return 0
		}
		return uint(val)
	case uint:
		return val
	case uint64:
		return uint(val)
	case string:
		val = strings.TrimSpace(val)
		if val == "" {
			return 0
		}
		n, err := strconv.ParseUint(val, 10, 32)
		if err != nil {
			return 0
		}
		return uint(n)
	default:
		return 0
	}
}

func intFromAny(v interface{}) int {
	switch val := v.(type) {
	case float64:
		return int(val)
	case float32:
		return int(val)
	case int:
		return val
	case int64:
		return int(val)
	case uint:
		return int(val)
	case uint64:
		return int(val)
	case string:
		val = strings.TrimSpace(val)
		if val == "" {
			return 0
		}
		n, err := strconv.Atoi(val)
		if err != nil {
			return 0
		}
		return n
	default:
		return 0
	}
}

func stringFromAny(v interface{}) string {
	switch val := v.(type) {
	case string:
		return val
	case float64:
		return strconv.FormatFloat(val, 'f', -1, 64)
	case float32:
		return strconv.FormatFloat(float64(val), 'f', -1, 32)
	case int:
		return strconv.Itoa(val)
	case int64:
		return strconv.FormatInt(val, 10)
	case uint:
		return strconv.FormatUint(uint64(val), 10)
	case uint64:
		return strconv.FormatUint(val, 10)
	default:
		return ""
	}
}

func normalizeContactFieldOptions(v interface{}) string {
	switch val := v.(type) {
	case nil:
		return "[]"
	case string:
		s := strings.TrimSpace(val)
		if s == "" || s == "[]" {
			return "[]"
		}
		if strings.HasPrefix(s, "[") {
			var arr []string
			if err := json.Unmarshal([]byte(s), &arr); err == nil {
				return marshalContactFieldOptions(arr)
			}
		}
		return marshalContactFieldOptions(splitContactFieldOptions(s))
	case []string:
		return marshalContactFieldOptions(val)
	case []interface{}:
		items := make([]string, 0, len(val))
		for _, item := range val {
			text := strings.TrimSpace(stringFromAny(item))
			if text == "" {
				continue
			}
			items = append(items, text)
		}
		return marshalContactFieldOptions(items)
	default:
		text := strings.TrimSpace(stringFromAny(val))
		if text == "" {
			return "[]"
		}
		return marshalContactFieldOptions([]string{text})
	}
}

func marshalContactFieldOptions(items []string) string {
	normalized := make([]string, 0, len(items))
	for _, item := range items {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		normalized = append(normalized, item)
	}
	if len(normalized) == 0 {
		return "[]"
	}
	data, err := json.Marshal(normalized)
	if err != nil {
		return "[]"
	}
	return string(data)
}

func (h *ContactHandler) ListByOffset(c *gin.Context) {
	corpID, _ := c.Get("corpId")
	if corpID == nil {
		corpID = uint(0)
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	offset := (page - 1) * pageSize

	contacts, total, err := h.contactSvc.ListByOffset(corpID.(uint), offset, pageSize, nil)
	if err != nil {
		response.Fail(c, response.ErrDB, "获取客户列表失败")
		return
	}
	response.PageResult(c, contacts, total, page, pageSize)
}

type ContactFieldHandler struct {
	svc       *service.ContactFieldService
	pivotSvc  *service.ContactFieldPivotService
}

func NewContactFieldHandler(db *gorm.DB) *ContactFieldHandler {
	return &ContactFieldHandler{
		svc:       service.NewContactFieldService(db),
		pivotSvc:  service.NewContactFieldPivotService(db),
	}
}

func (h *ContactFieldHandler) Index(c *gin.Context) {
	status, _ := strconv.Atoi(c.DefaultQuery("status", "2"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", c.DefaultQuery("perPage", "10")))

	fields, total, err := h.svc.ListByStatus(status, page, pageSize)
	if err != nil {
		response.Fail(c, response.ErrDB, "获取字段列表失败")
		return
	}

	// 处理 typeText
	for i := range fields {
		fields[i].TypeText = getFieldTypeText(fields[i].Type)
	}

	response.PageResult(c, fields, total, page, pageSize)
}

func (h *ContactFieldHandler) Show(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	field, err := h.svc.GetByID(uint(id))
	if err != nil {
		response.Fail(c, response.ErrNotFound, "字段不存在")
		return
	}

	// 处理 typeText
	field.TypeText = getFieldTypeText(field.Type)

	response.Success(c, field)
}

func (h *ContactFieldHandler) Store(c *gin.Context) {
	var field model.ContactField
	if err := c.ShouldBindJSON(&field); err != nil {
		response.Fail(c, response.ErrParams, "参数错误")
		return
	}
	if err := h.svc.Create(&field); err != nil {
		response.Fail(c, response.ErrDB, "创建字段失败")
		return
	}

	// 处理 typeText
	field.TypeText = getFieldTypeText(field.Type)

	response.Success(c, field)
}

func (h *ContactFieldHandler) Update(c *gin.Context) {
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ErrParams, "参数错误")
		return
	}
	id := uint(0)
	if raw, ok := req["id"]; ok {
		id = uintFromAny(raw)
	}
	if id == 0 {
		if parsedID, err := strconv.ParseUint(c.Param("id"), 10, 32); err == nil {
			id = uint(parsedID)
		}
	}
	if id == 0 {
		response.Fail(c, response.ErrParams, "参数错误")
		return
	}

	field, err := h.svc.GetByID(id)
	if err != nil {
		response.Fail(c, response.ErrNotFound, "字段不存在")
		return
	}
	if raw, ok := req["name"]; ok {
		field.Name = strings.TrimSpace(stringFromAny(raw))
	}
	if raw, ok := req["label"]; ok {
		field.Label = strings.TrimSpace(stringFromAny(raw))
		// old frontend often omits `name`, keep name aligned with label if empty.
		if strings.TrimSpace(field.Name) == "" {
			field.Name = field.Label
		}
	}
	if raw, ok := req["type"]; ok {
		field.Type = intFromAny(raw)
	}
	if raw, ok := req["options"]; ok {
		field.Options = normalizeContactFieldOptions(raw)
	}
	if raw, ok := req["order"]; ok {
		field.Order = intFromAny(raw)
	}
	if raw, ok := req["status"]; ok {
		field.Status = intFromAny(raw)
	}
	if raw, ok := req["isSys"]; ok {
		field.IsSys = intFromAny(raw)
	}
	if err := h.svc.Update(field); err != nil {
		response.Fail(c, response.ErrDB, "更新字段失败")
		return
	}

	// 处理 typeText
	field.TypeText = getFieldTypeText(field.Type)

	response.Success(c, field)
}

func (h *ContactFieldHandler) Destroy(c *gin.Context) {
	id, ok := parseIDFromParamOrBody(c)
	if !ok {
		response.Fail(c, response.ErrParams, "参数错误")
		return
	}
	if err := h.svc.Delete(id); err != nil {
		response.Fail(c, response.ErrDB, "删除字段失败")
		return
	}
	response.SuccessMsg(c, "删除成功")
}

func (h *ContactFieldHandler) StatusUpdate(c *gin.Context) {
	var req struct {
		ID     uint `json:"id"`
		Status int `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ErrParams, "参数错误")
		return
	}
	if req.ID == 0 {
		if id, err := strconv.ParseUint(c.Param("id"), 10, 32); err == nil {
			req.ID = uint(id)
		}
	}
	if req.ID == 0 {
		response.Fail(c, response.ErrParams, "参数错误")
		return
	}

	if err := h.svc.UpdateStatus(req.ID, req.Status); err != nil {
		response.Fail(c, response.ErrDB, "更新状态失败")
		return
	}

	response.SuccessMsg(c, "更新成功")
}

func (h *ContactFieldHandler) Portrait(c *gin.Context) {
	idText := c.Param("id")
	if idText != "" {
		id, _ := strconv.ParseUint(idText, 10, 32)
		fieldPivots, err := h.pivotSvc.List(uint(id))
		if err != nil {
			response.Fail(c, response.ErrDB, "获取客户画像失败")
			return
		}
		response.Success(c, fieldPivots)
		return
	}

	fields, err := h.svc.List()
	if err != nil {
		response.Fail(c, response.ErrDB, "获取客户画像字段失败")
		return
	}

	list := make([]gin.H, 0, len(fields))
	for _, field := range fields {
		if field.Status != 1 {
			continue
		}
		list = append(list, gin.H{
			"fieldId":  field.ID,
			"id":       field.ID,
			"name":     firstNonEmpty(strings.TrimSpace(field.Name), strings.TrimSpace(field.Label)),
			"label":    field.Label,
			"type":     field.Type,
			"typeText": getFieldTypeText(field.Type),
			"options":  splitContactFieldOptions(field.Options),
		})
	}

	response.Success(c, list)
}

func (h *ContactFieldHandler) BatchUpdate(c *gin.Context) {
	var req struct {
		Fields []struct {
			ID    uint `json:"id"`
			Order int  `json:"order"`
		} `json:"fields"`
		Update []struct {
			ID    uint `json:"id"`
			Order int  `json:"order"`
		} `json:"update"`
		Destroy []uint `json:"destroy"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ErrParams, "参数错误")
		return
	}

	fields := req.Fields
	if len(fields) == 0 {
		fields = req.Update
	}
	for _, item := range fields {
		if err := h.svc.UpdateOrder(item.ID, item.Order); err != nil {
			response.Fail(c, response.ErrDB, "批量更新失败")
			return
		}
	}
	for _, id := range req.Destroy {
		if id == 0 {
			continue
		}
		if err := h.svc.Delete(id); err != nil {
			response.Fail(c, response.ErrDB, "批量更新失败")
			return
		}
	}

	response.SuccessMsg(c, "批量更新成功")
}

type ContactTagHandler struct {
	svc *service.WorkContactTagService
}

func NewContactTagHandler(db *gorm.DB) *ContactTagHandler {
	return &ContactTagHandler{svc: service.NewWorkContactTagService(db)}
}

func (h *ContactTagHandler) Index(c *gin.Context) {
	corpID, _ := c.Get("corpId")
	if corpID == nil {
		corpID = uint(0)
	}
	tags, err := h.svc.List(corpID.(uint))
	if err != nil {
		response.Fail(c, response.ErrDB, "获取标签列表失败")
		return
	}
	response.Success(c, tags)
}

func (h *ContactTagHandler) AllTag(c *gin.Context) {
	h.Index(c)
}

func (h *ContactTagHandler) Store(c *gin.Context) {
	var tag model.WorkContactTag
	if err := c.ShouldBindJSON(&tag); err != nil {
		response.Fail(c, response.ErrParams, "参数错误")
		return
	}
	if err := h.svc.Create(&tag); err != nil {
		response.Fail(c, response.ErrDB, "创建标签失败")
		return
	}
	response.Success(c, tag)
}

func (h *ContactTagHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	tag, err := h.svc.GetByID(uint(id))
	if err != nil {
		response.Fail(c, response.ErrNotFound, "标签不存在")
		return
	}
	if err := c.ShouldBindJSON(tag); err != nil {
		response.Fail(c, response.ErrParams, "参数错误")
		return
	}
	if err := h.svc.Update(tag); err != nil {
		response.Fail(c, response.ErrDB, "更新标签失败")
		return
	}
	response.Success(c, tag)
}

func (h *ContactTagHandler) Destroy(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if err := h.svc.Delete(uint(id)); err != nil {
		response.Fail(c, response.ErrDB, "删除标签失败")
		return
	}
	response.SuccessMsg(c, "删除成功")
}

func (h *ContactTagHandler) Detail(c *gin.Context)            { h.Index(c) }
func (h *ContactTagHandler) ContactTagList(c *gin.Context)    { h.Index(c) }
func (h *ContactTagHandler) Move(c *gin.Context)              { response.SuccessMsg(c, "移动成功") }
func (h *ContactTagHandler) SynContactTag(c *gin.Context)     { response.SuccessMsg(c, "同步任务已提交") }

func (h *ContactTagHandler) UpdateByID(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ErrParams, err.Error())
		return
	}
	if err := h.svc.UpdateByID(uint(id), req); err != nil {
		response.Fail(c, response.ErrDB, "更新标签失败")
		return
	}
	response.SuccessMsg(c, "更新成功")
}

func (h *ContactTagHandler) ListByOrder(c *gin.Context) {
	corpID, _ := c.Get("corpId")
	if corpID == nil {
		corpID = uint(0)
	}
	tags, err := h.svc.ListByOrder(corpID.(uint))
	if err != nil {
		response.Fail(c, response.ErrDB, "获取标签列表失败")
		return
	}
	response.Success(c, tags)
}

// 辅助函数：获取字段类型文本
func getFieldTypeText(fieldType int) string {
	switch fieldType {
	case 1:
		return "单行文本"
	case 2:
		return "多行文本"
	case 3:
		return "数字"
	case 4:
		return "单选"
	case 5:
		return "多选"
	case 6:
		return "日期"
	case 7:
		return "附件"
	default:
		return "未知"
	}
}
