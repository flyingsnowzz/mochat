package plugin

import (
	"encoding/json"
	"strings"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"mochat-api-server/internal/model"
	"mochat-api-server/internal/pkg/response"
	"mochat-api-server/internal/service/plugin"
)

type ChannelCodeHandler struct {
	db       *gorm.DB
	svc      *plugin.ChannelCodeService
	groupSvc *plugin.ChannelCodeGroupService
}

type channelCodeGroupItem struct {
	GroupID uint   `json:"groupId"`
	Name    string `json:"name"`
}

func NewChannelCodeHandler(db *gorm.DB, svc *plugin.ChannelCodeService, groupSvc *plugin.ChannelCodeGroupService) *ChannelCodeHandler {
	return &ChannelCodeHandler{db: db, svc: svc, groupSvc: groupSvc}
}

func (h *ChannelCodeHandler) Index(c *gin.Context) {
	corpIDValue, ok := getCorpID(c)
	if !ok {
		response.Fail(c, response.ErrParams, "请先选择企业")
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", c.DefaultQuery("perPage", "20")))
	groupID, _ := strconv.Atoi(c.DefaultQuery("groupId", "0"))
	codeType, _ := strconv.Atoi(c.DefaultQuery("type", "0"))
	name := strings.TrimSpace(c.Query("name"))
	offset := (page - 1) * pageSize

	var (
		codes []model.ChannelCode
		total int64
	)
	query := h.db.Model(&model.ChannelCode{}).Where("corp_id = ? AND deleted_at IS NULL", corpIDValue)
	if groupID > 0 {
		query = query.Where("group_id = ?", groupID)
	}
	if codeType > 0 {
		query = query.Where("type = ?", codeType)
	}
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if err := query.Count(&total).Error; err != nil {
		response.Fail(c, response.ErrDB, "获取渠道码列表失败")
		return
	}
	if err := query.Order("id DESC").Offset(offset).Limit(pageSize).Find(&codes).Error; err != nil {
		response.Fail(c, response.ErrDB, "获取渠道码列表失败")
		return
	}

	list, err := h.buildChannelCodeList(corpIDValue, codes)
	if err != nil {
		response.Fail(c, response.ErrDB, "获取渠道码列表失败")
		return
	}
	response.PageResult(c, list, total, page, pageSize)
}

func (h *ChannelCodeHandler) Show(c *gin.Context) {
	corpIDValue, ok := getCorpID(c)
	if !ok {
		response.Fail(c, response.ErrParams, "请先选择企业")
		return
	}
	id := channelCodeIDFromRequest(c)
	if id == 0 {
		response.Fail(c, response.ErrParams, "参数错误")
		return
	}

	code, err := h.svc.GetByID(id)
	if err != nil || code.CorpID != corpIDValue {
		response.Fail(c, response.ErrDB, "获取渠道码详情失败")
		return
	}

	tagGroups, selectedTags, err := h.buildChannelCodeTagGroups(corpIDValue, code.Tags)
	if err != nil {
		response.Fail(c, response.ErrDB, "获取渠道码详情失败")
		return
	}

	baseInfo := gin.H{
		"groupId":       code.GroupID,
		"name":          code.Name,
		"autoAddFriend": code.AutoAddFriend,
		"tags":          tagGroups,
		"selectedTags":  selectedTags,
	}

	response.Success(c, gin.H{
		"baseInfo":         baseInfo,
		"drainageEmployee": h.normalizeChannelCodeDrainageEmployee(corpIDValue, code.DrainageEmployee),
		"welcomeMessage":   normalizeChannelCodeWelcomeMessage(code.WelcomeMessage),
	})
}

func (h *ChannelCodeHandler) Store(c *gin.Context) {
	corpIDValue, ok := getCorpID(c)
	if !ok {
		response.Fail(c, response.ErrParams, "请先选择企业")
		return
	}

	var req struct {
		BaseInfo struct {
			GroupID       uint        `json:"groupId"`
			Name          string      `json:"name" binding:"required"`
			AutoAddFriend int         `json:"autoAddFriend"`
			Tags          interface{} `json:"tags"`
		} `json:"baseInfo" binding:"required"`
		DrainageEmployee map[string]interface{} `json:"drainageEmployee" binding:"required"`
		WelcomeMessage   map[string]interface{} `json:"welcomeMessage" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ErrParams, "参数错误")
		return
	}
	tags, err := h.contactTags(corpIDValue)
	if err != nil {
		response.Fail(c, response.ErrDB, "创建渠道码失败")
		return
	}

	code := &model.ChannelCode{
		CorpID:           corpIDValue,
		GroupID:          req.BaseInfo.GroupID,
		Name:             strings.TrimSpace(req.BaseInfo.Name),
		QrcodeURL:        "https://dummyimage.com/240x240/1677ff/ffffff&text=QR",
		WxConfigID:       "",
		AutoAddFriend:    req.BaseInfo.AutoAddFriend,
		Tags:             normalizeChannelCodeTagJSON(req.BaseInfo.Tags, tags),
		Type:             uintInt(req.DrainageEmployee["type"]),
		DrainageEmployee: mustJSONString(req.DrainageEmployee, "{}"),
		WelcomeMessage:   mustJSONString(req.WelcomeMessage, "{}"),
	}
	if code.Name == "" {
		response.Fail(c, response.ErrParams, "参数错误")
		return
	}
	if err := h.svc.Create(code); err != nil {
		response.Fail(c, response.ErrDB, "创建渠道码失败")
		return
	}

	response.Success(c, gin.H{"id": code.ID})
}

func (h *ChannelCodeHandler) Update(c *gin.Context) {
	corpIDValue, ok := getCorpID(c)
	if !ok {
		response.Fail(c, response.ErrParams, "请先选择企业")
		return
	}

	var req struct {
		ChannelCodeID any `json:"channelCodeId" binding:"required"`
		BaseInfo      struct {
			GroupID       uint        `json:"groupId"`
			Name          string      `json:"name" binding:"required"`
			AutoAddFriend int         `json:"autoAddFriend"`
			Tags          interface{} `json:"tags"`
		} `json:"baseInfo" binding:"required"`
		DrainageEmployee map[string]interface{} `json:"drainageEmployee" binding:"required"`
		WelcomeMessage   map[string]interface{} `json:"welcomeMessage" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ErrParams, "参数错误")
		return
	}
	channelCodeID := uintFromAny(req.ChannelCodeID)
	if channelCodeID == 0 {
		response.Fail(c, response.ErrParams, "参数错误")
		return
	}

	code, err := h.svc.GetByID(channelCodeID)
	if err != nil || code.CorpID != corpIDValue {
		response.Fail(c, response.ErrDB, "修改渠道码失败")
		return
	}
	tags, err := h.contactTags(corpIDValue)
	if err != nil {
		response.Fail(c, response.ErrDB, "修改渠道码失败")
		return
	}

	updates := map[string]interface{}{
		"group_id":          req.BaseInfo.GroupID,
		"name":              strings.TrimSpace(req.BaseInfo.Name),
		"auto_add_friend":   req.BaseInfo.AutoAddFriend,
		"tags":              normalizeChannelCodeTagJSON(req.BaseInfo.Tags, tags),
		"type":              uintInt(req.DrainageEmployee["type"]),
		"drainage_employee": mustJSONString(req.DrainageEmployee, "{}"),
		"welcome_message":   mustJSONString(req.WelcomeMessage, "{}"),
	}
	if err := h.svc.Update(channelCodeID, updates); err != nil {
		response.Fail(c, response.ErrDB, "修改渠道码失败")
		return
	}
	response.SuccessMsg(c, "更新成功")
}

func (h *ChannelCodeHandler) Contact(c *gin.Context) {
	corpIDValue, ok := getCorpID(c)
	if !ok {
		response.Fail(c, response.ErrParams, "请先选择企业")
		return
	}

	channelCodeID := channelCodeIDFromRequest(c)
	if channelCodeID == 0 {
		response.Fail(c, response.ErrParams, "参数错误")
		return
	}
	code, err := h.svc.GetByID(channelCodeID)
	if err != nil || code.CorpID != corpIDValue {
		response.Fail(c, response.ErrDB, "获取扫码客户失败")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", c.DefaultQuery("perPage", "20")))
	offset := (page - 1) * pageSize

	type row struct {
		ContactID   uint   `json:"contactId"`
		Name        string `json:"name"`
		Employee    string `json:"employee"`
		CreateTime  string `json:"createTime"`
	}

	var total int64
	baseQuery := h.db.Table("mc_work_contact_employee AS wce").
		Joins("JOIN mc_work_contact AS wc ON wc.id = wce.contact_id").
		Joins("LEFT JOIN mc_work_employee AS we ON we.id = wce.employee_id").
		Where("wce.corp_id = ?", corpIDValue)

	if err := baseQuery.Count(&total).Error; err != nil {
		response.Fail(c, response.ErrDB, "获取扫码客户失败")
		return
	}

	var rows []row
	if err := baseQuery.
		Select("wce.contact_id AS contact_id, wc.name AS name, we.name AS employee, DATE_FORMAT(wce.create_time, '%Y-%m-%d %H:%i:%s') AS create_time").
		Order("wce.contact_id DESC").
		Offset(offset).
		Limit(pageSize).
		Scan(&rows).Error; err != nil {
		response.Fail(c, response.ErrDB, "获取扫码客户失败")
		return
	}

	list := make([]gin.H, 0, len(rows))
	for _, item := range rows {
		list = append(list, gin.H{
			"contactId":  item.ContactID,
			"name":       item.Name,
			"employees":  item.Employee,
			"createTime": item.CreateTime,
		})
	}

	response.PageResult(c, list, total, page, pageSize)
}

func (h *ChannelCodeHandler) StatisticsIndex(c *gin.Context) {
	corpIDValue, ok := getCorpID(c)
	if !ok {
		response.Fail(c, response.ErrParams, "请先选择企业")
		return
	}
	channelCodeID := channelCodeIDFromRequest(c)
	if channelCodeID == 0 {
		response.Fail(c, response.ErrParams, "参数错误")
		return
	}
	code, err := h.svc.GetByID(channelCodeID)
	if err != nil || code.CorpID != corpIDValue {
		response.Fail(c, response.ErrDB, "获取渠道码统计失败")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", c.DefaultQuery("perPage", "20")))
	typeValue, _ := strconv.Atoi(c.DefaultQuery("type", "2"))

	rows, err := h.channelCodeStatisticRows(corpIDValue, typeValue, c.Query("startTime"), c.Query("endTime"))
	if err != nil {
		response.Fail(c, response.ErrDB, "获取渠道码统计失败")
		return
	}

	total := int64(len(rows))
	start := (page - 1) * pageSize
	if start > len(rows) {
		start = len(rows)
	}
	end := start + pageSize
	if end > len(rows) {
		end = len(rows)
	}
	response.PageResult(c, rows[start:end], total, page, pageSize)
}

func (h *ChannelCodeHandler) Statistics(c *gin.Context) {
	corpIDValue, ok := getCorpID(c)
	if !ok {
		response.Fail(c, response.ErrParams, "请先选择企业")
		return
	}
	channelCodeID := channelCodeIDFromRequest(c)
	if channelCodeID == 0 {
		response.Fail(c, response.ErrParams, "参数错误")
		return
	}
	code, err := h.svc.GetByID(channelCodeID)
	if err != nil || code.CorpID != corpIDValue {
		response.Fail(c, response.ErrDB, "获取渠道码统计失败")
		return
	}

	typeValue, _ := strconv.Atoi(c.DefaultQuery("type", "2"))
	rows, err := h.channelCodeStatisticRows(corpIDValue, typeValue, c.Query("startTime"), c.Query("endTime"))
	if err != nil {
		response.Fail(c, response.ErrDB, "获取渠道码统计失败")
		return
	}

	addToday, defriendToday, deleteToday, netToday, err := h.channelCodeTodaySummary(corpIDValue)
	if err != nil {
		response.Fail(c, response.ErrDB, "获取渠道码统计失败")
		return
	}

	addTotal := 0
	defriendTotal := 0
	deleteTotal := 0
	netTotal := 0
	for _, item := range rows {
		addTotal += item["addNumRange"].(int)
		defriendTotal += item["defriendNumRange"].(int)
		deleteTotal += item["deleteNumRange"].(int)
		netTotal += item["netNumRange"].(int)
	}

	response.Success(c, gin.H{
		"addNum":          addToday,
		"defriendNum":     defriendToday,
		"deleteNum":       deleteToday,
		"netNum":          netToday,
		"addNumLong":      addTotal,
		"defriendNumLong": defriendTotal,
		"deleteNumLong":   deleteTotal,
		"netNumLong":      netTotal,
		"list":            rows,
	})
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

func channelCodeIDFromRequest(c *gin.Context) uint {
	if id, err := strconv.Atoi(c.Param("id")); err == nil && id > 0 {
		return uint(id)
	}
	if id, err := strconv.Atoi(c.Query("channelCodeId")); err == nil && id > 0 {
		return uint(id)
	}
	return 0
}

func (h *ChannelCodeHandler) buildChannelCodeList(corpID uint, codes []model.ChannelCode) ([]gin.H, error) {
	groupMap, err := h.channelCodeGroupMap(corpID)
	if err != nil {
		return nil, err
	}
	tags, err := h.contactTags(corpID)
	if err != nil {
		return nil, err
	}
	tagNameMap := make(map[uint]string, len(tags))
	for _, tag := range tags {
		tagNameMap[tag.ID] = tag.Name
	}

	list := make([]gin.H, 0, len(codes))
	for _, code := range codes {
		list = append(list, gin.H{
			"channelCodeId":  code.ID,
			"groupId":        code.GroupID,
			"groupName":      groupMap[code.GroupID],
			"name":           code.Name,
			"qrcodeUrl":      code.QrcodeURL,
			"type":           channelCodeTypeText(code.Type),
			"contactNum":     0,
			"tags":           tagNames(code.Tags, tags, tagNameMap),
			"autoAddFriend":  channelCodeAutoAddFriendText(code.AutoAddFriend),
			"drainageConfig": parseChannelCodeJSONObject(code.DrainageEmployee),
		})
	}
	return list, nil
}

func (h *ChannelCodeHandler) buildChannelCodeTagGroups(corpID uint, rawTags string) ([]gin.H, []uint, error) {
	var groups []model.WorkContactTagGroup
	if err := h.db.Where("corp_id = ? AND deleted_at IS NULL", corpID).Order("id ASC").Find(&groups).Error; err != nil {
		return nil, nil, err
	}
	tags, err := h.contactTags(corpID)
	if err != nil {
		return nil, nil, err
	}
	selectedTagIDs := normalizeChannelCodeTagIDs(rawTags, tags)

	selectedSet := make(map[uint]struct{}, len(selectedTagIDs))
	for _, id := range selectedTagIDs {
		selectedSet[id] = struct{}{}
	}

	groupTags := make(map[uint][]gin.H)
	for _, tag := range tags {
		_, selected := selectedSet[tag.ID]
		groupTags[tag.ContactTagGroupID] = append(groupTags[tag.ContactTagGroupID], gin.H{
			"tagId":      tag.ID,
			"tagName":    tag.Name,
			"isSelected": selectedFlag(selected),
		})
	}

	result := make([]gin.H, 0, len(groups))
	for _, group := range groups {
		items := groupTags[group.ID]
		if items == nil {
			items = []gin.H{}
		}
		result = append(result, gin.H{
			"groupId":   group.ID,
			"groupName": group.GroupName,
			"list":      items,
		})
	}
	return result, selectedTagIDs, nil
}

func (h *ChannelCodeHandler) channelCodeGroupMap(corpID uint) (map[uint]string, error) {
	var groups []model.ChannelCodeGroup
	if err := h.db.Where("corp_id = ? AND deleted_at IS NULL", corpID).Find(&groups).Error; err != nil {
		return nil, err
	}
	result := map[uint]string{0: "未分组"}
	for _, group := range groups {
		result[group.ID] = group.Name
	}
	return result, nil
}

func (h *ChannelCodeHandler) contactTags(corpID uint) ([]model.WorkContactTag, error) {
	var tags []model.WorkContactTag
	if err := h.db.Where("corp_id = ? AND deleted_at IS NULL", corpID).Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}

func (h *ChannelCodeHandler) employeeNameMap(corpID uint) (map[uint]string, error) {
	var employees []model.WorkEmployee
	if err := h.db.Where("corp_id = ? AND deleted_at IS NULL", corpID).Find(&employees).Error; err != nil {
		return nil, err
	}
	result := make(map[uint]string, len(employees))
	for _, employee := range employees {
		result[employee.ID] = employee.Name
	}
	return result, nil
}

func tagNames(raw string, tags []model.WorkContactTag, nameMap map[uint]string) []string {
	ids := normalizeChannelCodeTagIDs(raw, tags)
	result := make([]string, 0, len(ids))
	for _, id := range ids {
		if name := nameMap[id]; name != "" {
			result = append(result, name)
		}
	}
	return result
}

func (h *ChannelCodeHandler) normalizeChannelCodeDrainageEmployee(corpID uint, raw string) gin.H {
	result := parseChannelCodeJSONObject(raw)
	if len(result) == 0 {
		return gin.H{
			"type": 1,
			"employees": []gin.H{},
			"specialPeriod": gin.H{
				"status": 2,
				"detail": []gin.H{},
			},
			"addMax": gin.H{
				"status":          2,
				"employees":       []gin.H{},
				"spareEmployeeIds": []uint{},
			},
		}
	}

	employeeMap, err := h.employeeNameMap(corpID)
	if err != nil {
		employeeMap = map[uint]string{}
	}

	result["type"] = uintInt(result["type"])
	result["employees"] = normalizeChannelCodeEmployeeRules(result["employees"], employeeMap, uintInt(result["type"]))
	result["specialPeriod"] = normalizeChannelCodeSpecialPeriod(result["specialPeriod"], employeeMap, uintInt(result["type"]))
	result["addMax"] = normalizeChannelCodeAddMax(result["addMax"])
	return result
}

func parseChannelCodeJSONObject(raw string) map[string]interface{} {
	result := map[string]interface{}{}
	if strings.TrimSpace(raw) == "" {
		return result
	}
	_ = json.Unmarshal([]byte(raw), &result)
	return result
}

func mustJSONString(v interface{}, fallback string) string {
	if v == nil {
		return fallback
	}
	data, err := json.Marshal(v)
	if err != nil || string(data) == "null" {
		return fallback
	}
	return string(data)
}

func parseChannelCodeUintSlice(raw string) []uint {
	if strings.TrimSpace(raw) == "" {
		return nil
	}
	var ids []uint
	if err := json.Unmarshal([]byte(raw), &ids); err == nil {
		return ids
	}
	var mixed []interface{}
	if err := json.Unmarshal([]byte(raw), &mixed); err == nil {
		return interfaceToUintSlice(mixed)
	}
	return nil
}

func normalizeChannelCodeTagJSON(raw interface{}, tags []model.WorkContactTag) string {
	return mustJSONString(normalizeChannelCodeTagIDs(raw, tags), "[]")
}

func normalizeChannelCodeTagIDs(raw interface{}, tags []model.WorkContactTag) []uint {
	ids := channelCodeTagIDsFromAny(raw)
	if len(ids) == 0 {
		return []uint{}
	}
	byID := make(map[uint]struct{}, len(tags))
	byWxID := make(map[string]uint, len(tags))
	for _, tag := range tags {
		byID[tag.ID] = struct{}{}
		if strings.TrimSpace(tag.WxContactTagID) != "" {
			byWxID[strings.TrimSpace(tag.WxContactTagID)] = tag.ID
		}
	}

	result := make([]uint, 0, len(ids))
	seen := make(map[uint]struct{}, len(ids))
	for _, id := range ids {
		resolved := id
		if _, ok := byID[id]; !ok {
			if mapped, ok := byWxID[strconv.FormatUint(uint64(id), 10)]; ok {
				resolved = mapped
			} else {
				continue
			}
		}
		if _, ok := seen[resolved]; ok {
			continue
		}
		seen[resolved] = struct{}{}
		result = append(result, resolved)
	}
	return result
}

func channelCodeTagIDsFromAny(raw interface{}) []uint {
	switch val := raw.(type) {
	case nil:
		return []uint{}
	case string:
		return parseChannelCodeUintSlice(val)
	case []uint:
		return val
	case []int:
		result := make([]uint, 0, len(val))
		for _, item := range val {
			if item > 0 {
				result = append(result, uint(item))
			}
		}
		return result
	case []string:
		result := make([]uint, 0, len(val))
		for _, item := range val {
			if id := uintFromAny(item); id > 0 {
				result = append(result, id)
			}
		}
		return result
	case []interface{}:
		return interfaceToUintSlice(val)
	default:
		if id := uintFromAny(val); id > 0 {
			return []uint{id}
		}
	}
	return []uint{}
}

func selectedFlag(selected bool) int {
	if selected {
		return 1
	}
	return 2
}

func uintInt(v interface{}) int {
	switch val := v.(type) {
	case float64:
		return int(val)
	case int:
		return val
	case string:
		n, _ := strconv.Atoi(val)
		return n
	default:
		return 0
	}
}

func channelCodeTypeText(v int) string {
	switch v {
	case 1:
		return "单人"
	case 2:
		return "多人"
	default:
		return ""
	}
}

func channelCodeAutoAddFriendText(v int) string {
	if v == 1 {
		return "开启"
	}
	if v == 2 {
		return "关闭"
	}
	return ""
}

func normalizeChannelCodeEmployeeRules(raw interface{}, employeeMap map[uint]string, memberType int) []gin.H {
	list, ok := raw.([]interface{})
	if !ok {
		return []gin.H{}
	}
	result := make([]gin.H, 0, len(list))
	for _, item := range list {
		rule, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		week := uintInt(rule["week"])
		if week == 0 && rule["week"] == nil {
			week = uintInt(rule["weekDay"])
		}
		slots := normalizeChannelCodeTimeSlots(rule["timeSlot"], employeeMap, memberType)
		result = append(result, gin.H{
			"week":     week,
			"weekDay":  week,
			"timeSlot": slots,
		})
	}
	return result
}

func normalizeChannelCodeTimeSlots(raw interface{}, employeeMap map[uint]string, memberType int) []gin.H {
	list, ok := raw.([]interface{})
	if !ok {
		return []gin.H{}
	}
	result := make([]gin.H, 0, len(list))
	for _, item := range list {
		slot, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		employeeIDs := interfaceToUintSlice(slot["employeeId"])
		selectMembers := interfaceToStringSlice(slot["selectMembers"])
		if len(selectMembers) == 0 {
			for _, id := range employeeIDs {
				if name := employeeMap[id]; name != "" {
					selectMembers = append(selectMembers, name)
				}
			}
		}

		normalized := gin.H{
			"startTime":        stringValue(slot["startTime"]),
			"endTime":          stringValue(slot["endTime"]),
			"employeeId":       employeeIDs,
			"departmentId":     interfaceToUintSlice(slot["departmentId"]),
			"selectMembers":    selectMembers,
			"selectDepartment": interfaceToStringSlice(slot["selectDepartment"]),
		}
		if memberType == 1 {
			label := ""
			key := ""
			if len(selectMembers) > 0 {
				label = selectMembers[0]
			}
			if len(employeeIDs) > 0 {
				key = strconv.Itoa(int(employeeIDs[0]))
			}
			normalized["employeeSelect"] = gin.H{
				"label": label,
				"key":   key,
			}
		}
		result = append(result, normalized)
	}
	return result
}

func normalizeChannelCodeSpecialPeriod(raw interface{}, employeeMap map[uint]string, memberType int) gin.H {
	result := gin.H{
		"status": 2,
		"detail": []gin.H{},
	}
	data, ok := raw.(map[string]interface{})
	if !ok {
		return result
	}
	result["status"] = uintInt(data["status"])
	list, ok := data["detail"].([]interface{})
	if !ok {
		return result
	}
	detail := make([]gin.H, 0, len(list))
	for _, item := range list {
		row, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		detail = append(detail, gin.H{
			"startDate":  stringValue(row["startDate"]),
			"endDate":    stringValue(row["endDate"]),
			"dataString": interfaceToStringSlice(row["dataString"]),
			"timeSlot":   normalizeChannelCodeTimeSlots(row["timeSlot"], employeeMap, memberType),
		})
	}
	result["detail"] = detail
	return result
}

func normalizeChannelCodeAddMax(raw interface{}) gin.H {
	result := gin.H{
		"status":           2,
		"employees":        []gin.H{},
		"spareEmployeeIds": []uint{},
	}
	data, ok := raw.(map[string]interface{})
	if !ok {
		return result
	}
	result["status"] = uintInt(data["status"])
	if employees, ok := data["employees"].([]interface{}); ok {
		rows := make([]gin.H, 0, len(employees))
		for _, item := range employees {
			if row, ok := item.(map[string]interface{}); ok {
				rows = append(rows, gin.H{
					"employeeId":   uintInt(row["employeeId"]),
					"employeeName": stringValue(row["employeeName"]),
					"max":          stringValue(row["max"]),
				})
			}
		}
		result["employees"] = rows
	}
	result["spareEmployeeIds"] = interfaceToUintSlice(data["spareEmployeeIds"])
	return result
}

func normalizeChannelCodeWelcomeMessage(raw string) gin.H {
	data := parseChannelCodeJSONObject(raw)
	result := gin.H{
		"scanCodePush": 1,
		"messageDetail": []gin.H{
			{
				"type":           1,
				"welcomeContent": "",
				"mediumId":       "",
				"content":        gin.H{},
			},
			{
				"type":   2,
				"status": 2,
				"detail": []gin.H{},
			},
			{
				"type":   3,
				"status": 2,
				"detail": []gin.H{},
			},
		},
	}
	if len(data) == 0 {
		return result
	}

	result["scanCodePush"] = uintInt(data["scanCodePush"])
	rawDetails, ok := data["messageDetail"].([]interface{})
	if !ok {
		return result
	}

	details := result["messageDetail"].([]gin.H)
	if len(rawDetails) > 0 {
		details[0] = normalizeChannelCodeGeneralMessage(rawDetails[0])
	}
	if len(rawDetails) > 1 {
		details[1] = normalizeChannelCodeTimedMessage(rawDetails[1], 2)
	}
	if len(rawDetails) > 2 {
		details[2] = normalizeChannelCodeTimedMessage(rawDetails[2], 3)
	}
	result["messageDetail"] = details
	return result
}

func normalizeChannelCodeGeneralMessage(raw interface{}) gin.H {
	result := gin.H{
		"type":           1,
		"welcomeContent": "",
		"mediumId":       "",
		"content":        gin.H{},
	}
	data, ok := raw.(map[string]interface{})
	if !ok {
		return result
	}
	result["type"] = uintInt(data["type"])
	result["welcomeContent"] = stringValue(data["welcomeContent"])
	result["mediumId"] = stringValue(data["mediumId"])
	if content, ok := data["content"].(map[string]interface{}); ok {
		result["content"] = content
	}
	return result
}

func normalizeChannelCodeTimedMessage(raw interface{}, messageType int) gin.H {
	result := gin.H{
		"type":   messageType,
		"status": 2,
		"detail": []gin.H{},
	}
	data, ok := raw.(map[string]interface{})
	if !ok {
		return result
	}
	result["type"] = uintInt(data["type"])
	if uintInt(data["type"]) == 0 {
		result["type"] = messageType
	}
	result["status"] = uintInt(data["status"])

	rawDetail, ok := data["detail"].([]interface{})
	if !ok {
		return result
	}
	rows := make([]gin.H, 0, len(rawDetail))
	for _, item := range rawDetail {
		row, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		normalized := gin.H{
			"key":         stringValue(row["key"]),
			"chooseCycle": interfaceToIntSlice(row["chooseCycle"]),
			"startDate":   stringValue(row["startDate"]),
			"endDate":     stringValue(row["endDate"]),
			"dataString":  interfaceToStringSlice(row["dataString"]),
			"timeSlot":    normalizeChannelCodeWelcomeSlots(row["timeSlot"]),
		}
		rows = append(rows, normalized)
	}
	result["detail"] = rows
	return result
}

func normalizeChannelCodeWelcomeSlots(raw interface{}) []gin.H {
	list, ok := raw.([]interface{})
	if !ok {
		return []gin.H{}
	}
	result := make([]gin.H, 0, len(list))
	for _, item := range list {
		slot, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		content := gin.H{}
		if val, ok := slot["content"].(map[string]interface{}); ok {
			content = val
		}
		result = append(result, gin.H{
			"welcomeContent": stringValue(slot["welcomeContent"]),
			"mediumId":       stringValue(slot["mediumId"]),
			"startTime":      stringValue(slot["startTime"]),
			"endTime":        stringValue(slot["endTime"]),
			"content":        content,
		})
	}
	return result
}

func interfaceToUintSlice(v interface{}) []uint {
	switch val := v.(type) {
	case []uint:
		return val
	case []interface{}:
		result := make([]uint, 0, len(val))
		for _, item := range val {
			n := uintInt(item)
			if n > 0 {
				result = append(result, uint(n))
			}
		}
		return result
	case float64, int, string:
		n := uintInt(val)
		if n > 0 {
			return []uint{uint(n)}
		}
	}
	return []uint{}
}

func interfaceToStringSlice(v interface{}) []string {
	switch val := v.(type) {
	case []string:
		return val
	case []interface{}:
		result := make([]string, 0, len(val))
		for _, item := range val {
			s := stringValue(item)
			if s != "" {
				result = append(result, s)
			}
		}
		return result
	case string:
		if val == "" {
			return []string{}
		}
		return []string{val}
	}
	return []string{}
}

func interfaceToIntSlice(v interface{}) []int {
	switch val := v.(type) {
	case []int:
		return val
	case []interface{}:
		result := make([]int, 0, len(val))
		for _, item := range val {
			result = append(result, uintInt(item))
		}
		return result
	case float64, int, string:
		n := uintInt(val)
		if n > 0 || val == 0 {
			return []int{n}
		}
	}
	return []int{}
}

func stringValue(v interface{}) string {
	if s, ok := v.(string); ok {
		return s
	}
	if n, ok := v.(float64); ok {
		return strconv.FormatFloat(n, 'f', 0, 64)
	}
	return ""
}

func uintFromAny(v interface{}) uint {
	switch val := v.(type) {
	case uint:
		return val
	case int:
		if val > 0 {
			return uint(val)
		}
	case int64:
		if val > 0 {
			return uint(val)
		}
	case float64:
		if val > 0 {
			return uint(val)
		}
	case string:
		n, _ := strconv.Atoi(strings.TrimSpace(val))
		if n > 0 {
			return uint(n)
		}
	}
	return 0
}

func (h *ChannelCodeHandler) channelCodeStatisticRows(corpID uint, typeValue int, startTime, endTime string) ([]gin.H, error) {
	periodFormat, parseLayout := statisticPeriodFormat(typeValue)
	if periodFormat == "" {
		periodFormat = "%Y-%m-%d"
		parseLayout = "2006-01-02"
	}

	query := h.db.Table("mc_work_contact_employee").
		Select(
			"DATE_FORMAT(create_time, ?) AS period, "+
				"SUM(CASE WHEN status = 1 THEN 1 ELSE 0 END) AS add_num, "+
				"SUM(CASE WHEN status = 3 THEN 1 ELSE 0 END) AS defriend_num, "+
				"SUM(CASE WHEN status = 2 THEN 1 ELSE 0 END) AS delete_num",
			periodFormat,
		).
		Where("corp_id = ?", corpID)

	if startTime != "" {
		query = query.Where("DATE(create_time) >= ?", startTime)
	}
	if endTime != "" {
		query = query.Where("DATE(create_time) <= ?", endTime)
	}

	type rawRow struct {
		Period       string `json:"period"`
		AddNum       int    `json:"addNum"`
		DefriendNum  int    `json:"defriendNum"`
		DeleteNum    int    `json:"deleteNum"`
	}
	var rawRows []rawRow
	if err := query.Group("period").Order("period ASC").Scan(&rawRows).Error; err != nil {
		return nil, err
	}

	rows := make([]gin.H, 0, len(rawRows))
	for _, row := range rawRows {
		rows = append(rows, gin.H{
			"time":             formatStatisticTime(row.Period, parseLayout, typeValue),
			"addNumRange":      row.AddNum,
			"defriendNumRange": row.DefriendNum,
			"deleteNumRange":   row.DeleteNum,
			"netNumRange":      row.AddNum - row.DefriendNum - row.DeleteNum,
		})
	}
	return rows, nil
}

func (h *ChannelCodeHandler) channelCodeTodaySummary(corpID uint) (int, int, int, int, error) {
	type summary struct {
		AddNum      int `json:"addNum"`
		DefriendNum int `json:"defriendNum"`
		DeleteNum   int `json:"deleteNum"`
	}
	var item summary
	err := h.db.Table("mc_work_contact_employee").
		Select(
			"SUM(CASE WHEN status = 1 THEN 1 ELSE 0 END) AS add_num, "+
				"SUM(CASE WHEN status = 3 THEN 1 ELSE 0 END) AS defriend_num, "+
				"SUM(CASE WHEN status = 2 THEN 1 ELSE 0 END) AS delete_num",
		).
		Where("corp_id = ? AND DATE(create_time) = CURDATE()", corpID).
		Scan(&item).Error
	if err != nil {
		return 0, 0, 0, 0, err
	}
	return item.AddNum, item.DefriendNum, item.DeleteNum, item.AddNum - item.DefriendNum - item.DeleteNum, nil
}

func statisticPeriodFormat(typeValue int) (string, string) {
	switch typeValue {
	case 3:
		return "%Y-%m", "2006-01"
	case 2:
		return "%x-W%v", "2006-W01"
	default:
		return "%Y-%m-%d", "2006-01-02"
	}
}

func formatStatisticTime(period, layout string, typeValue int) string {
	if period == "" {
		return ""
	}
	if typeValue == 2 {
		return period
	}
	t, err := time.Parse(layout, period)
	if err != nil {
		return period
	}
	if typeValue == 3 {
		return t.Format("2006-01")
	}
	return t.Format("2006-01-02")
}
