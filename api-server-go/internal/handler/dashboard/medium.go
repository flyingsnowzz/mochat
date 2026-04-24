package dashboard

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"mochat-api-server/internal/model"
	"mochat-api-server/internal/pkg/response"
	"mochat-api-server/internal/service"
)

type MediumHandler struct {
	db       *gorm.DB
	svc      *service.MediumService
	groupSvc *service.MediumGroupService
}

func NewMediumHandler(db *gorm.DB) *MediumHandler {
	return &MediumHandler{
		db:       db,
		svc:      service.NewMediumService(db),
		groupSvc: service.NewMediumGroupService(db),
	}
}

func (h *MediumHandler) Index(c *gin.Context) {
	corpID, _ := c.Get("corpId")
	if corpID == nil {
		corpID = uint(0)
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("perPage", c.DefaultQuery("pageSize", "10")))
	mediumType, _ := strconv.Atoi(c.DefaultQuery("type", "0"))
	groupID, _ := strconv.ParseUint(c.DefaultQuery("mediumGroupId", c.DefaultQuery("groupId", "0")), 10, 32)
	searchStr := c.Query("searchStr")

	media, total, err := h.svc.ListWithSearch(corpID.(uint), page, pageSize, mediumType, uint(groupID), searchStr)
	if err != nil {
		response.Fail(c, response.ErrDB, "获取素材列表失败")
		return
	}
	groupNames := h.loadMediumGroupNames(corpID.(uint))
	items := make([]gin.H, 0, len(media))
	for _, medium := range media {
		items = append(items, buildMediumResponseItem(medium, groupNames[medium.MediumGroupID]))
	}
	response.PageResult(c, items, total, page, pageSize)
}

func (h *MediumHandler) Show(c *gin.Context) {
	id := getUintID(c, "id")
	medium, err := h.svc.GetByID(uint(id))
	if err != nil {
		response.Fail(c, response.ErrNotFound, "素材不存在")
		return
	}

	corpID, _ := c.Get("corpId")
	groupNames := h.loadMediumGroupNames(asUint(corpID))
	response.Success(c, buildMediumResponseItem(*medium, groupNames[medium.MediumGroupID]))
}

func (h *MediumHandler) Store(c *gin.Context) {
	var req struct {
		Type          int                    `json:"type" binding:"required"`
		IsSync        int                    `json:"isSync"`
		Content       map[string]interface{} `json:"content" binding:"required"`
		MediumGroupID uint                   `json:"mediumGroupId"`
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

	content, err := service.MarshalMediumContent(req.Content)
	if err != nil {
		response.Fail(c, response.ErrParams, "参数错误")
		return
	}
	if req.IsSync == 0 {
		req.IsSync = 1
	}

	medium := model.Medium{
		Type:          req.Type,
		IsSync:        req.IsSync,
		Content:       content,
		CorpID:        corpID.(uint),
		MediumGroupID: req.MediumGroupID,
		UserID:        asUint(userID),
		UserName:      h.loadUserName(asUint(userID)),
	}
	if err := h.svc.Create(&medium); err != nil {
		response.Fail(c, response.ErrDB, "创建素材失败")
		return
	}

	groupNames := h.loadMediumGroupNames(corpID.(uint))
	response.Success(c, buildMediumResponseItem(medium, groupNames[medium.MediumGroupID]))
}

func (h *MediumHandler) Update(c *gin.Context) {
	id := getUintID(c, "id")
	medium, err := h.svc.GetByID(uint(id))
	if err != nil {
		response.Fail(c, response.ErrNotFound, "素材不存在")
		return
	}

	var req struct {
		ID            uint                   `json:"id"`
		Type          int                    `json:"type"`
		IsSync        int                    `json:"isSync"`
		Content       map[string]interface{} `json:"content"`
		MediumGroupID uint                   `json:"mediumGroupId"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ErrParams, "参数错误")
		return
	}
	if req.Type > 0 {
		medium.Type = req.Type
	}
	if req.IsSync > 0 {
		medium.IsSync = req.IsSync
	}
	if req.Content != nil {
		content, err := service.MarshalMediumContent(req.Content)
		if err != nil {
			response.Fail(c, response.ErrParams, "参数错误")
			return
		}
		medium.Content = content
	}
	medium.MediumGroupID = req.MediumGroupID
	if err := h.svc.Update(medium); err != nil {
		response.Fail(c, response.ErrDB, "更新素材失败")
		return
	}

	corpID, _ := c.Get("corpId")
	groupNames := h.loadMediumGroupNames(asUint(corpID))
	response.Success(c, buildMediumResponseItem(*medium, groupNames[medium.MediumGroupID]))
}

func (h *MediumHandler) Destroy(c *gin.Context) {
	id := getUintID(c, "id")
	if err := h.svc.Delete(uint(id)); err != nil {
		response.Fail(c, response.ErrDB, "删除素材失败")
		return
	}
	response.SuccessMsg(c, "删除成功")
}

func (h *MediumHandler) GroupUpdate(c *gin.Context) {
	var req struct {
		ID      uint `json:"id" binding:"required"`
		GroupID uint `json:"mediumGroupId"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ErrParams, "参数错误")
		return
	}
	medium, err := h.svc.GetByID(req.ID)
	if err != nil {
		response.Fail(c, response.ErrNotFound, "素材不存在")
		return
	}
	medium.MediumGroupID = req.GroupID
	if err := h.svc.Update(medium); err != nil {
		response.Fail(c, response.ErrDB, "移动素材失败")
		return
	}
	response.SuccessMsg(c, "移动成功")
}

func (h *MediumHandler) loadMediumGroupNames(corpID uint) map[uint]string {
	result := map[uint]string{0: "未分组"}
	if corpID == 0 {
		return result
	}
	groups, err := h.groupSvc.List(corpID)
	if err != nil {
		return result
	}
	for _, group := range groups {
		result[group.ID] = group.Name
	}
	return result
}

func (h *MediumHandler) loadUserName(userID uint) string {
	if userID == 0 {
		return ""
	}
	userSvc := service.NewUserService(h.db)
	user, err := userSvc.GetByID(userID)
	if err != nil {
		return ""
	}
	return user.Name
}

func buildMediumResponseItem(medium model.Medium, groupName string) gin.H {
	content := service.ParseMediumContent(medium.Content)
	appendMediumFullPaths(content)
	title := mediumTitle(medium.Type, content)
	return gin.H{
		"id":              medium.ID,
		"title":           title,
		"type":            service.MediumTypeName(medium.Type),
		"typeValue":       medium.Type,
		"content":         content,
		"mediumGroupId":   medium.MediumGroupID,
		"mediumGroupName": groupName,
		"mediaId":         medium.MediaID,
		"userId":          medium.UserID,
		"userName":        medium.UserName,
		"isSync":          medium.IsSync,
		"createdAt":       medium.CreatedAt.Format("2006-01-02 15:04:05"),
		"updatedAt":       medium.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func appendMediumFullPaths(content map[string]interface{}) {
	pairs := map[string]string{
		"imagePath": "imageFullPath",
		"videoPath": "videoFullPath",
		"voicePath": "voiceFullPath",
		"filePath":  "fileFullPath",
	}
	for rawKey, fullKey := range pairs {
		if path, ok := content[rawKey].(string); ok && path != "" {
			content[fullKey] = fmt.Sprintf("http://localhost:9501/uploads/%s", path)
		}
	}
}

func mediumTitle(mediumType int, content map[string]interface{}) string {
	switch mediumType {
	case 1:
		return stringValue(content["title"])
	case 2:
		return stringValue(content["imageName"])
	case 3, 6:
		return stringValue(content["title"])
	case 4:
		return stringValue(content["voiceName"])
	case 5:
		return stringValue(content["videoName"])
	case 7:
		return stringValue(content["fileName"])
	default:
		return ""
	}
}

func stringValue(value interface{}) string {
	if text, ok := value.(string); ok {
		return text
	}
	return ""
}

func getUintID(c *gin.Context, key string) uint64 {
	if val := c.Param(key); val != "" {
		id, _ := strconv.ParseUint(val, 10, 32)
		return id
	}
	if val := c.Query(key); val != "" {
		id, _ := strconv.ParseUint(val, 10, 32)
		return id
	}
	if val := c.Query("mediumId"); val != "" {
		id, _ := strconv.ParseUint(val, 10, 32)
		return id
	}
	if val := c.DefaultPostForm(key, ""); val != "" {
		id, _ := strconv.ParseUint(val, 10, 32)
		return id
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

type MediumGroupHandler struct {
	svc *service.MediumGroupService
}

func NewMediumGroupHandler(db *gorm.DB) *MediumGroupHandler {
	return &MediumGroupHandler{svc: service.NewMediumGroupService(db)}
}

func (h *MediumGroupHandler) Index(c *gin.Context) {
	corpID, _ := c.Get("corpId")
	if corpID == nil {
		corpID = uint(0)
	}
	groups, err := h.svc.List(corpID.(uint))
	if err != nil {
		response.Fail(c, response.ErrDB, "获取素材分组列表失败")
		return
	}
	response.Success(c, groups)
}

func (h *MediumGroupHandler) Store(c *gin.Context) {
	var group model.MediumGroup
	if err := c.ShouldBindJSON(&group); err != nil {
		response.Fail(c, response.ErrParams, "参数错误")
		return
	}
	if err := h.svc.Create(&group); err != nil {
		response.Fail(c, response.ErrDB, "创建素材分组失败")
		return
	}
	response.Success(c, group)
}

func (h *MediumGroupHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	group := &model.MediumGroup{}
	group.ID = uint(id)
	if err := c.ShouldBindJSON(group); err != nil {
		response.Fail(c, response.ErrParams, "参数错误")
		return
	}
	if err := h.svc.Update(group); err != nil {
		response.Fail(c, response.ErrDB, "更新素材分组失败")
		return
	}
	response.Success(c, group)
}

func (h *MediumGroupHandler) Destroy(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if err := h.svc.Delete(uint(id)); err != nil {
		response.Fail(c, response.ErrDB, "删除素材分组失败")
		return
	}
	response.SuccessMsg(c, "删除成功")
}
