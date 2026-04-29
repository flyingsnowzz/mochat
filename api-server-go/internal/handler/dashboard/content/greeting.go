package content

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"mochat-api-server/internal/model"
	"mochat-api-server/internal/pkg/response"
	"mochat-api-server/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GreetingHandler 欢迎语管理处理器
// 处理客户欢迎语的增删改查操作
type GreetingHandler struct {
	db  *gorm.DB
	svc *service.GreetingService
}

// NewGreetingHandler 创建欢迎语管理处理器实例
func NewGreetingHandler(db *gorm.DB) *GreetingHandler {
	return &GreetingHandler{
		db:  db,
		svc: service.NewGreetingService(db),
	}
}

// Index 获取欢迎语列表
// 支持分页查询企业的欢迎语列表，返回包含员工信息和素材内容
// 处理流程：
// 1. 获取企业 ID 和分页参数
// 2. 调用服务层获取欢迎语列表
// 3. 构建返回数据，包括：
//   - 欢迎语列表项
//   - 是否存在通用欢迎语
//   - 已配置欢迎语的员工列表
//   - 分页信息
//
// 参数：
//
//	page - 页码，默认为 1
//	perPage/pageSize - 每页数量，默认为 10
//
// 返回：包含欢迎语列表、是否有通用欢迎语、已配置的员工列表、分页信息的响应
func (h *GreetingHandler) Index(c *gin.Context) {
	// 获取企业 ID
	corpID, exists := c.Get("corpId")
	if !exists || corpID.(uint) == 0 {
		response.Fail(c, response.ErrParams, "未获取到企业信息")
		return
	}

	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("perPage", c.DefaultQuery("pageSize", "10")))

	// 调用服务层获取欢迎语列表
	greetings, total, err := h.svc.List(corpID.(uint), page, pageSize)
	if err != nil {
		response.Fail(c, response.ErrDB, "获取欢迎语列表失败")
		return
	}

	// 构建返回数据
	items := make([]gin.H, 0, len(greetings))
	hadGeneral := 0                         // 是否存在通用欢迎语（1-存在，0-不存在）
	hadEmployees := make([]uint, 0)         // 已配置欢迎语的员工 ID 列表
	seenEmployee := make(map[uint]struct{}) // 用于去重员工 ID

	// 处理每个欢迎语
	for _, greeting := range greetings {
		// 构建欢迎语列表项
		item := h.buildGreetingListItem(greeting)
		items = append(items, item)

		// 检查是否为通用欢迎语
		if greeting.RangeType == 1 {
			hadGeneral = 1
		}

		// 收集已配置欢迎语的员工 ID
		for _, employee := range greetingEmployeeIDs(greeting.Employees) {
			if _, ok := seenEmployee[employee]; ok {
				continue // 跳过已处理的员工
			}
			seenEmployee[employee] = struct{}{}
			hadEmployees = append(hadEmployees, employee)
		}
	}

	// 返回成功响应
	response.Success(c, gin.H{
		"list":         items,        // 欢迎语列表
		"hadGeneral":   hadGeneral,   // 是否存在通用欢迎语
		"hadEmployees": hadEmployees, // 已配置欢迎语的员工列表
		"page": gin.H{
			"perPage":     pageSize,                       // 每页数量
			"total":       total,                          // 总记录数
			"totalPage":   calcTotalPage(total, pageSize), // 总页数
			"currentPage": page,                           // 当前页码
		},
		"currentPage": page,     // 当前页码（兼容旧接口）
		"pageSize":    pageSize, // 每页数量（兼容旧接口）
		"total":       total,    // 总记录数（兼容旧接口）
	})
}

// Show 获取欢迎语详情
// 根据 ID 获取指定欢迎语的详细信息
// 参数：id/greetingId - 欢迎语 ID
// 返回：欢迎语详情信息
func (h *GreetingHandler) Show(c *gin.Context) {
	id := getGreetingUintID(c, "id", "greetingId")
	greeting, err := h.svc.GetByID(uint(id))
	if err != nil {
		response.Fail(c, response.ErrNotFound, "欢迎语不存在")
		return
	}
	response.Success(c, h.buildGreetingDetailItem(*greeting))
}

// Store 创建新欢迎语
// 创建新的客户欢迎语配置
// 参数（JSON Body）：
//
//	rangeType - 范围类型（1-全部，2-指定）
//	employees - 员工 ID 列表
//	type - 欢迎语类型
//	words - 欢迎语文本
//	mediumId - 素材 ID
//
// 返回：创建的欢迎语详情
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

	corpID, exists := c.Get("corpId")
	if !exists || corpID.(uint) == 0 {
		response.Fail(c, response.ErrParams, "请先选择企业")
		return
	}

	// 验证 rangeType 参数
	if req.RangeType != 1 && req.RangeType != 2 {
		response.Fail(c, response.ErrParams, "适用成员类型 值必须在列表内：[1,2]")
		return
	}

	// 处理欢迎语类型，将逗号分隔的类型转换为 -type- 格式
	typeArr := strings.Split(req.Type, ",")
	greetingType := ""
	if len(typeArr) == 1 {
		greetingType = fmt.Sprintf("-%s-", strings.TrimSpace(typeArr[0]))
	} else {
		cleanTypes := make([]string, 0, len(typeArr))
		for _, t := range typeArr {
			cleanTypes = append(cleanTypes, strings.TrimSpace(t))
		}
		greetingType = fmt.Sprintf("-%s-", strings.Join(cleanTypes, "-"))
	}

	// 处理员工列表
	employeesJSON := "[]"
	if req.Employees != "" {
		employeeIDs := strings.Split(req.Employees, ",")
		cleanIDs := make([]uint, 0, len(employeeIDs))
		for _, idStr := range employeeIDs {
			idStr = strings.TrimSpace(idStr)
			if idStr == "" {
				continue
			}
			if id, err := strconv.ParseUint(idStr, 10, 32); err == nil {
				cleanIDs = append(cleanIDs, uint(id))
			}
		}
		if len(cleanIDs) > 0 {
			employeesJSONBytes, _ := json.Marshal(cleanIDs)
			employeesJSON = string(employeesJSONBytes)
		}
	}

	greeting := model.Greeting{
		CorpID:    corpID.(uint),
		Type:      greetingType,
		Words:     req.Words,
		MediumID:  req.MediumID,
		RangeType: req.RangeType,
		Employees: employeesJSON,
	}

	if err := h.svc.Create(&greeting); err != nil {
		response.Fail(c, response.ErrDB, "创建欢迎语失败")
		return
	}

	response.Success(c, h.buildGreetingDetailItem(greeting))
}

// Update 更新欢迎语
// 更新指定欢迎语的配置信息
// 参数（JSON Body）：
//
//	greetingId - 欢迎语 ID
//	rangeType - 范围类型
//	employees - 员工 ID 列表
//	type - 欢迎语类型
//	words - 欢迎语文本
//	mediumId - 素材 ID
//
// 返回：更新后的欢迎语详情
func (h *GreetingHandler) Update(c *gin.Context) {
	id := getGreetingUintID(c, "id", "greetingId")
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
			response.Fail(c, response.ErrParams, "欢迎语ID 必填")
			return
		}
	}

	// 验证 rangeType 参数
	if req.RangeType != 1 && req.RangeType != 2 {
		response.Fail(c, response.ErrParams, "适用成员类型 值必须在列表内：[1,2]")
		return
	}

	greeting, err := h.svc.GetByID(uint(id))
	if err != nil {
		response.Fail(c, response.ErrNotFound, "欢迎语不存在")
		return
	}

	// 处理欢迎语类型
	typeArr := strings.Split(req.Type, ",")
	greetingType := ""
	if len(typeArr) == 1 {
		greetingType = fmt.Sprintf("-%s-", strings.TrimSpace(typeArr[0]))
	} else {
		cleanTypes := make([]string, 0, len(typeArr))
		for _, t := range typeArr {
			cleanTypes = append(cleanTypes, strings.TrimSpace(t))
		}
		greetingType = fmt.Sprintf("-%s-", strings.Join(cleanTypes, "-"))
	}

	// 处理员工列表
	employeesJSON := "[]"
	if req.Employees != "" {
		employeeIDs := strings.Split(req.Employees, ",")
		cleanIDs := make([]uint, 0, len(employeeIDs))
		for _, idStr := range employeeIDs {
			idStr = strings.TrimSpace(idStr)
			if idStr == "" {
				continue
			}
			if id, err := strconv.ParseUint(idStr, 10, 32); err == nil {
				cleanIDs = append(cleanIDs, uint(id))
			}
		}
		if len(cleanIDs) > 0 {
			employeesJSONBytes, _ := json.Marshal(cleanIDs)
			employeesJSON = string(employeesJSONBytes)
		}
	}

	greeting.RangeType = req.RangeType
	greeting.Employees = employeesJSON
	greeting.Type = greetingType
	greeting.Words = req.Words
	greeting.MediumID = req.MediumID

	if err := h.svc.Update(greeting); err != nil {
		response.Fail(c, response.ErrDB, "更新欢迎语失败")
		return
	}

	response.Success(c, h.buildGreetingDetailItem(*greeting))
}

// Destroy 删除欢迎语
// 根据 ID 删除指定的欢迎语
// 参数：id/greetingId - 欢迎语 ID
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
		response.Fail(c, response.ErrParams, "欢迎语ID 必填")
		return
	}

	// 获取欢迎语信息，验证企业归属
	greeting, err := h.svc.GetByID(uint(id))
	if err != nil {
		response.Fail(c, response.ErrNotFound, "欢迎语不存在，不可操作")
		return
	}

	corpID, exists := c.Get("corpId")
	if !exists || greeting.CorpID != corpID.(uint) {
		response.Fail(c, response.ErrParams, "此欢迎语不归属当前企业，不可操作")
		return
	}

	if err := h.svc.Delete(uint(id)); err != nil {
		response.Fail(c, response.ErrDB, "删除欢迎语失败")
		return
	}

	response.SuccessMsg(c, "删除成功")
}

// buildGreetingListItem 构建欢迎语列表项
// 将模型数据转换为前端需要的格式
// 参数：greeting - 欢迎语模型
// 返回：包含列表项数据
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
		mediumService := service.NewMediumService(h.db)
		if medium, err := mediumService.GetByID(greeting.MediumID); err == nil {
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

// buildGreetingDetailItem 构建欢迎语详情项
// 用于 Show、Store、Update 方法返回详情数据
// 参数：greeting - 欢迎语模型
// 返回：包含详情数据的 gin.H
func (h *GreetingHandler) buildGreetingDetailItem(greeting model.Greeting) gin.H {
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
		mediumService := service.NewMediumService(h.db)
		if medium, err := mediumService.GetByID(greeting.MediumID); err == nil {
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

// loadGreetingEmployees 加载欢迎语关联的员工信息
// 根据员工 ID 列表批量加载员工数据
// 参数：ids - 员工 ID 列表
// 返回：员工 ID 到员工模型的映射
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

// greetingEmployeeJSON 将员工数据转换为 JSON 字符串
// 参数：raw - 原始员工数据（JSON 字符串或逗号分隔字符串）
// 返回：标准 JSON 数组字符串
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

// greetingEmployeeIDs 解析员工 ID 列表
// 支持 JSON 数组和逗号分隔两种格式
// 参数：raw - 原始员工数据
// 返回：员工 ID 列表
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

// greetingTypeText 转换欢迎语类型为文本
// 参数：greetingType - 欢迎语类型字符串（格式如 "-1-6-" 或 "1,6"）
// 返回：类型名称文本，如 "文本+小程序"
func greetingTypeText(greetingType string) string {
	greetingType = strings.TrimSpace(greetingType)
	if greetingType == "" {
		return ""
	}

	// 优先按连字符分割，兼容 PHP 的 -1-6- 格式
	if strings.Contains(greetingType, "-") {
		greetingType = strings.Trim(greetingType, "-")
	} else {
		// 按逗号分割，兼容 1,6 格式
		greetingType = strings.ReplaceAll(greetingType, ",", "-")
	}

	parts := strings.Split(greetingType, "-")
	names := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		switch part {
		case "1":
			names = append(names, "文本")
		case "2":
			names = append(names, "图片")
		case "3":
			names = append(names, "图文")
		case "4":
			names = append(names, "音频")
		case "5":
			names = append(names, "视频")
		case "6":
			names = append(names, "小程序")
		case "7":
			names = append(names, "文件")
		}
	}
	return strings.Join(names, "+")
}

// greetingRangeTypeText 转换范围类型为文本
// 参数：rangeType - 范围类型（1-全部，2-指定）
// 返回：范围类型文本
func greetingRangeTypeText(rangeType int) string {
	if rangeType == 2 {
		return "指定成员"
	}
	return "全部成员"
}

// calcTotalPage 计算总页数
// 参数：
//
//	total - 总记录数
//	pageSize - 每页数量
//
// 返回：总页数
func calcTotalPage(total int64, pageSize int) int {
	if pageSize <= 0 {
		return 0
	}
	if total == 0 {
		return 0
	}
	return int((total + int64(pageSize) - 1) / int64(pageSize))
}

// getGreetingUintID 获取欢迎语 ID
// 从路径参数或查询参数中获取 ID，支持多个参数名
// 参数：
//
//	c - Gin 上下文
//	keys - 可能的参数名列表
//
// 返回：ID 值
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
