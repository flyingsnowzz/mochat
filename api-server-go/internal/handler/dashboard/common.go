package dashboard

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"mochat-api-server/internal/model"
	"mochat-api-server/internal/pkg/response"
	"mochat-api-server/internal/pkg/storage"
	"mochat-api-server/internal/service"
)

type CommonHandler struct{}

func NewCommonHandler() *CommonHandler {
	return &CommonHandler{}
}

func (h *CommonHandler) Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.Fail(c, response.ErrParams, "请上传文件")
		return
	}

	filename := storage.GenerateFilename(file.Filename)
	src, err := file.Open()
	if err != nil {
		response.Fail(c, response.ErrFileUpload, "读取文件失败")
		return
	}
	defer src.Close()

	path, err := storage.DefaultStorage.Upload(src, filename)
	if err != nil {
		response.Fail(c, response.ErrFileUpload, "上传文件失败")
		return
	}

	url := storage.DefaultStorage.GetURL(path)
	response.Success(c, gin.H{
		"url":  url,
		"path": path,
		"name": file.Filename,
		"size": file.Size,
	})
}

func (h *CommonHandler) UploadFile(c *gin.Context) {
	h.Upload(c)
}

type IndexHandler struct {
	dataSvc *service.CorpDayDataService
}

func NewIndexHandler(db *gorm.DB) *IndexHandler {
	return &IndexHandler{dataSvc: service.NewCorpDayDataService(db)}
}

func (h *IndexHandler) Index(c *gin.Context) {
	response.Success(c, gin.H{})
}

func (h *IndexHandler) LineChat(c *gin.Context) {
	response.Success(c, gin.H{})
}

type ChatToolHandler struct {
	svc *service.ChatToolService
}

func NewChatToolHandler(db *gorm.DB) *ChatToolHandler {
	return &ChatToolHandler{svc: service.NewChatToolService(db)}
}

func (h *ChatToolHandler) Index(c *gin.Context) {
	tools, err := h.svc.List()
	if err != nil {
		response.Fail(c, response.ErrDB, "获取侧边工具栏列表失败")
		return
	}
	response.Success(c, tools)
}

type AgentHandler struct {
	svc *service.WorkAgentService
}

func NewAgentHandler(db *gorm.DB) *AgentHandler {
	return &AgentHandler{svc: service.NewWorkAgentService(db)}
}

func (h *AgentHandler) Store(c *gin.Context) {
	var agent model.WorkAgent
	if err := c.ShouldBindJSON(&agent); err != nil {
		response.Fail(c, response.ErrParams, "参数错误")
		return
	}
	if err := h.svc.Create(&agent); err != nil {
		response.Fail(c, response.ErrDB, "创建应用失败")
		return
	}
	response.Success(c, agent)
}

func (h *AgentHandler) TxtVerifyShow(c *gin.Context) {
	response.Success(c, gin.H{})
}

func (h *AgentHandler) TxtVerifyUpload(c *gin.Context) {
	response.SuccessMsg(c, "上传成功")
}

type GreetingHandler struct {
	db  *gorm.DB
	svc *service.GreetingService
}

func NewGreetingHandler(db *gorm.DB) *GreetingHandler {
	return &GreetingHandler{db: db, svc: service.NewGreetingService(db)}
}

func (h *GreetingHandler) Index(c *gin.Context) {
	corpID, _ := c.Get("corpId")
	if corpID == nil {
		corpID = uint(0)
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("perPage", c.DefaultQuery("pageSize", "10")))

	greetings, total, err := h.svc.List(corpID.(uint), page, pageSize)
	if err != nil {
		response.Fail(c, response.ErrDB, "获取欢迎语列表失败")
		return
	}
	items := make([]gin.H, 0, len(greetings))
	hadGeneral := 0
	hadEmployees := make([]uint, 0)
	seenEmployee := make(map[uint]struct{})
	for _, greeting := range greetings {
		item := h.buildGreetingListItem(greeting)
		items = append(items, item)
		if greeting.RangeType == 1 {
			hadGeneral = 1
		}
		for _, employee := range greetingEmployeeIDs(greeting.Employees) {
			if _, ok := seenEmployee[employee]; ok {
				continue
			}
			seenEmployee[employee] = struct{}{}
			hadEmployees = append(hadEmployees, employee)
		}
	}
	response.Success(c, gin.H{
		"list":         items,
		"hadGeneral":   hadGeneral,
		"hadEmployees": hadEmployees,
		"page": gin.H{
			"perPage":     pageSize,
			"total":       total,
			"totalPage":   calcTotalPage(total, pageSize),
			"currentPage": page,
		},
		"currentPage": page,
		"pageSize":    pageSize,
		"total":       total,
	})
}

func (h *GreetingHandler) Show(c *gin.Context) {
	id := getGreetingUintID(c, "id", "greetingId")
	greeting, err := h.svc.GetByID(uint(id))
	if err != nil {
		response.Fail(c, response.ErrNotFound, "欢迎语不存在")
		return
	}
	response.Success(c, h.buildGreetingDetailItem(*greeting))
}

func (h *GreetingHandler) Store(c *gin.Context) {
	var req struct {
		RangeType int    `json:"rangeType" binding:"required"`
		Employees string `json:"employees"`
		Type      string `json:"type" binding:"required"`
		Words     string `json:"words"`
		MediumID  uint   `json:"mediumId"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ErrParams, "参数错误")
		return
	}
	corpID, _ := c.Get("corpId")
	if corpID == nil || corpID.(uint) == 0 {
		response.Fail(c, response.ErrParams, "请先选择企业")
		return
	}
	greeting := model.Greeting{
		CorpID:    corpID.(uint),
		Type:      strings.TrimSpace(req.Type),
		Words:     req.Words,
		MediumID:  req.MediumID,
		RangeType: req.RangeType,
		Employees: greetingEmployeeJSON(req.Employees),
	}
	if err := h.svc.Create(&greeting); err != nil {
		response.Fail(c, response.ErrDB, "创建欢迎语失败")
		return
	}
	response.Success(c, h.buildGreetingDetailItem(greeting))
}

func (h *GreetingHandler) Update(c *gin.Context) {
	id := getGreetingUintID(c, "id", "greetingId")
	var greeting *model.Greeting
	var err error

	var req struct {
		GreetingID uint   `json:"greetingId"`
		RangeType  int    `json:"rangeType" binding:"required"`
		Employees  string `json:"employees"`
		Type       string `json:"type" binding:"required"`
		Words      string `json:"words"`
		MediumID   uint   `json:"mediumId"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ErrParams, "参数错误")
		return
	}
	if id == 0 {
		id = uint64(req.GreetingID)
		if id == 0 {
			response.Fail(c, response.ErrParams, "参数错误")
			return
		}
	}
	greeting, err = h.svc.GetByID(uint(id))
	if err != nil {
		response.Fail(c, response.ErrNotFound, "欢迎语不存在")
		return
	}
	greeting.RangeType = req.RangeType
	greeting.Employees = greetingEmployeeJSON(req.Employees)
	greeting.Type = strings.TrimSpace(req.Type)
	greeting.Words = req.Words
	greeting.MediumID = req.MediumID
	if err := h.svc.Update(greeting); err != nil {
		response.Fail(c, response.ErrDB, "更新欢迎语失败")
		return
	}
	response.Success(c, h.buildGreetingDetailItem(*greeting))
}

func (h *GreetingHandler) Destroy(c *gin.Context) {
	id := getGreetingUintID(c, "id", "greetingId")
	if id == 0 {
		var req struct {
			GreetingID uint `json:"greetingId"`
		}
		if err := c.ShouldBindJSON(&req); err == nil {
			id = uint64(req.GreetingID)
		}
	}
	if id == 0 {
		response.Fail(c, response.ErrParams, "参数错误")
		return
	}
	if err := h.svc.Delete(uint(id)); err != nil {
		response.Fail(c, response.ErrDB, "删除欢迎语失败")
		return
	}
	response.SuccessMsg(c, "删除成功")
}

func (h *GreetingHandler) buildGreetingListItem(greeting model.Greeting) gin.H {
	employeeIDs := greetingEmployeeIDs(greeting.Employees)
	employeeMap := h.loadGreetingEmployees(employeeIDs)
	employeeNames := make([]string, 0, len(employeeIDs))
	for _, employeeID := range employeeIDs {
		employee, ok := employeeMap[employeeID]
		if !ok {
			continue
		}
		employeeNames = append(employeeNames, employee.Name)
	}

	mediumContent := gin.H{}
	typeText := greetingTypeText(greeting.Type)
	if greeting.MediumID > 0 {
		if medium, err := service.NewMediumService(h.db).GetByID(greeting.MediumID); err == nil {
			content := service.ParseMediumContent(medium.Content)
			appendMediumFullPaths(content)
			mediumContent = gin.H(content)
			if greeting.Words == "" {
				typeText = service.MediumTypeName(medium.Type)
			}
		}
	}

	return gin.H{
		"greetingId":    greeting.ID,
		"rangeType":     greeting.RangeType,
		"rangeTypeText": greetingRangeTypeText(greeting.RangeType),
		"employees":     employeeNames,
		"words":         greeting.Words,
		"mediumId":      greeting.MediumID,
		"mediumContent": mediumContent,
		"type":          greeting.Type,
		"typeText":      typeText,
		"createdAt":     greeting.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

func (h *GreetingHandler) buildGreetingDetailItem(greeting model.Greeting) gin.H {
	employeeIDs := greetingEmployeeIDs(greeting.Employees)
	employeeMap := h.loadGreetingEmployees(employeeIDs)
	employees := make([]gin.H, 0, len(employeeIDs))
	for _, employeeID := range employeeIDs {
		employee, ok := employeeMap[employeeID]
		if !ok {
			continue
		}
		employees = append(employees, gin.H{
			"employeeId":   employee.ID,
			"employeeName": employee.Name,
			"name":         employee.Name,
			"id":           employee.ID,
		})
	}

	mediumContent := gin.H{}
	typeText := greetingTypeText(greeting.Type)
	if greeting.MediumID > 0 {
		if medium, err := service.NewMediumService(h.db).GetByID(greeting.MediumID); err == nil {
			content := service.ParseMediumContent(medium.Content)
			appendMediumFullPaths(content)
			mediumContent = gin.H(content)
			if greeting.Words == "" {
				typeText = service.MediumTypeName(medium.Type)
			}
		}
	}

	return gin.H{
		"greetingId":    greeting.ID,
		"rangeType":     greeting.RangeType,
		"rangeTypeText": greetingRangeTypeText(greeting.RangeType),
		"employees":     employees,
		"words":         greeting.Words,
		"mediumId":      greeting.MediumID,
		"mediumContent": mediumContent,
		"type":          greeting.Type,
		"typeText":      typeText,
		"createdAt":     greeting.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

func (h *GreetingHandler) loadGreetingEmployees(ids []uint) map[uint]model.WorkEmployee {
	result := make(map[uint]model.WorkEmployee)
	if len(ids) == 0 {
		return result
	}
	var employees []model.WorkEmployee
	if err := h.db.Where("id IN ?", ids).Find(&employees).Error; err != nil {
		return result
	}
	for _, employee := range employees {
		result[employee.ID] = employee
	}
	return result
}

func greetingEmployeeJSON(raw string) string {
	ids := greetingEmployeeIDs(raw)
	if len(ids) == 0 {
		return "[]"
	}
	data, err := json.Marshal(ids)
	if err != nil {
		return "[]"
	}
	return string(data)
}

func greetingEmployeeIDs(raw string) []uint {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return []uint{}
	}
	if strings.HasPrefix(raw, "[") {
		var ids []uint
		if err := json.Unmarshal([]byte(raw), &ids); err == nil {
			return ids
		}
	}

	parts := strings.Split(raw, ",")
	result := make([]uint, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		id, err := strconv.ParseUint(part, 10, 32)
		if err != nil {
			continue
		}
		result = append(result, uint(id))
	}
	return result
}

func greetingTypeText(greetingType string) string {
	parts := strings.Split(greetingType, ",")
	names := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		switch part {
		case "1":
			names = append(names, "文本")
		case "2":
			names = append(names, "图片")
		case "3":
			names = append(names, "图文")
		case "6":
			names = append(names, "小程序")
		}
	}
	return strings.Join(names, "+")
}

func greetingRangeTypeText(rangeType int) string {
	if rangeType == 2 {
		return "指定成员"
	}
	return "全部成员"
}

func calcTotalPage(total int64, pageSize int) int {
	if pageSize <= 0 {
		return 0
	}
	if total == 0 {
		return 0
	}
	return int((total + int64(pageSize) - 1) / int64(pageSize))
}

func getGreetingUintID(c *gin.Context, keys ...string) uint64 {
	for _, key := range keys {
		if val := c.Param(key); val != "" {
			id, _ := strconv.ParseUint(val, 10, 32)
			return id
		}
		if val := c.Query(key); val != "" {
			id, _ := strconv.ParseUint(val, 10, 32)
			return id
		}
	}
	return 0
}
