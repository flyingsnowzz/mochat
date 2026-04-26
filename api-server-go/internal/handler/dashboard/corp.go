// Package dashboard 提供 Dashboard 相关的 HTTP 处理器
// 该文件包含企业管理的处理器：
// CorpHandler - 处理企业的增删改查、绑定等操作
package dashboard

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"mochat-api-server/internal/model"
	"mochat-api-server/internal/pkg/response"
	mochatRedis "mochat-api-server/internal/redis"
	"mochat-api-server/internal/service"
)

// CorpHandler 企业管理处理器
// 提供企业的增删改查、绑定等功能
// 主要职责：
// 1. 获取企业列表
// 2. 获取企业选择列表
// 3. 获取企业详情
// 4. 创建企业
// 5. 更新企业
// 6. 绑定企业
// 7. 处理企业微信回调
//
// 依赖服务：
// - CorpService: 企业服务
// - gorm.DB: 数据库连接
// - Redis: 用于缓存用户绑定的企业信息

type CorpHandler struct {
	svc *service.CorpService // 企业服务
	db  *gorm.DB            // 数据库连接
}

// corpIndexItem 企业列表项
// 用于返回企业列表数据

type corpIndexItem struct {
	CorpID           uint   `json:"corpId"`           // 企业 ID
	CorpName         string `json:"corpName"`         // 企业名称
	WxCorpID         string `json:"wxCorpId"`         // 企业微信 ID
	EmployeeSecret   string `json:"employeeSecret"`   // 员工权限密钥
	ContactSecret    string `json:"contactSecret"`    // 客户联系权限密钥
	EventCallback    string `json:"eventCallback"`    // 事件回调地址
	Token            string `json:"token"`            // 令牌
	EncodingAesKey   string `json:"encodingAesKey"`   // 加密密钥
	CreatedAt        string `json:"createdAt"`        // 创建时间
	ChatApplyStatus  int    `json:"chatApplyStatus"`  // 聊天申请状态
	ChatStatus       int    `json:"chatStatus"`       // 聊天状态
	MessageCreatedAt string `json:"messageCreatedAt"` // 消息创建时间
}

// NewCorpHandler 创建企业管理处理器实例
// 参数：db - GORM 数据库连接
// 返回：企业管理处理器实例
func NewCorpHandler(db *gorm.DB) *CorpHandler {
	return &CorpHandler{svc: service.NewCorpService(db), db: db}
}

// Index 获取企业列表
// 获取企业列表，支持分页和企业名称搜索
// 处理流程：
// 1. 获取租户 ID、页码、每页数量和企业名称参数
// 2. 调用服务层获取企业列表
// 3. 构建返回数据
// 4. 返回分页结果
// 参数：
//
//	page - 页码，默认为 1
//	perPage/pageSize - 每页数量，默认为 20
//	corpName - 企业名称（用于搜索）
//
// 返回：包含企业列表、总数、分页信息的响应
func (h *CorpHandler) Index(c *gin.Context) {
	// 获取租户 ID
	tenantID, _ := c.Get("tenantId")
	
	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("perPage", c.DefaultQuery("pageSize", "20")))
	
	// 获取企业名称搜索参数
	corpName := strings.TrimSpace(c.Query("corpName"))

	// 调用服务层获取企业列表
	corps, total, err := h.svc.List(tenantID.(uint), corpName, page, pageSize)
	if err != nil {
		response.Fail(c, response.ErrDB, "获取企业列表失败")
		return
	}

	// 构建返回数据
	items := make([]corpIndexItem, 0, len(corps))
	for _, corp := range corps {
		createdAt := ""
		if !corp.CreatedAt.IsZero() {
			createdAt = corp.CreatedAt.Format("2006-01-02 15:04:05")
		}
		items = append(items, corpIndexItem{
			CorpID:           corp.ID,
			CorpName:         corp.Name,
			WxCorpID:         corp.WxCorpid,
			EmployeeSecret:   corp.EmployeeSecret,
			ContactSecret:    corp.ContactSecret,
			EventCallback:    corp.EventCallback,
			Token:            corp.Token,
			EncodingAesKey:   corp.EncodingAesKey,
			CreatedAt:        createdAt,
			ChatApplyStatus:  0, // 暂未实现
			ChatStatus:       0, // 暂未实现
			MessageCreatedAt: createdAt,
		})
	}
	
	// 返回分页结果
	response.PageResult(c, items, total, page, pageSize)
}

// Select 获取企业选择列表
// 获取用户可选择的企业列表，超级管理员可看到所有企业，普通用户只能看到自己所属的企业
// 处理流程：
// 1. 获取租户 ID 和用户 ID
// 2. 查询用户信息，判断是否为超级管理员
// 3. 如果是超级管理员，返回所有企业
// 4. 如果是普通用户，返回用户所属的企业
//
// 返回：企业选择列表
func (h *CorpHandler) Select(c *gin.Context) {
	// 获取租户 ID 和用户 ID
	tenantID, _ := c.Get("tenantId")
	userID, _ := c.Get("userId")

	// 查询用户信息，判断是否为超级管理员
	var user model.User
	if err := h.db.Select("id", "isSuperAdmin").First(&user, userID.(uint)).Error; err != nil {
		response.Fail(c, response.ErrDB, "获取企业选择列表失败")
		return
	}

	// 定义企业选项结构体
	type corpOption struct {
		CorpID   uint   `json:"corpId"`   // 企业 ID
		CorpName string `json:"corpName"` // 企业名称
	}
	var data []corpOption

	// 超级管理员可看到所有企业
	if user.IsSuperAdmin == 1 {
		var corps []model.Corp
		if err := h.db.Where("tenant_id = ?", tenantID.(uint)).Order("id ASC").Find(&corps).Error; err != nil {
			response.Fail(c, response.ErrDB, "获取企业选择列表失败")
			return
		}
		for _, corp := range corps {
			data = append(data, corpOption{CorpID: corp.ID, CorpName: corp.Name})
		}
		response.Success(c, data)
		return
	}

	// 普通用户只能看到自己所属的企业
	if err := h.db.Raw(`
		SELECT DISTINCT c.id AS corp_id, c.name AS corp_name
		FROM mc_corp c
		JOIN mc_work_employee we ON we.corp_id = c.id
		WHERE c.tenant_id = ? AND we.log_user_id = ? AND c.deleted_at IS NULL AND we.deleted_at IS NULL
		ORDER BY c.id ASC
	`, tenantID.(uint), user.ID).Scan(&data).Error; err != nil {
		response.Fail(c, response.ErrDB, "获取企业选择列表失败")
		return
	}

	// 返回企业选择列表
	response.Success(c, data)
}

// Show 获取企业详情
// 根据企业 ID 获取其详细信息
// 处理流程：
// 1. 获取企业 ID
// 2. 调用服务层获取企业详情
// 3. 验证企业是否属于当前租户
// 4. 返回企业详情信息
// 参数：
//
//	id/corpId - 企业 ID（路径参数或查询参数）
//
// 返回：企业详情信息
func (h *CorpHandler) Show(c *gin.Context) {
	// 获取企业 ID
	id := getCorpIDParam(c)
	
	// 调用服务层获取企业详情
	corp, err := h.svc.GetByID(uint(id))
	if err != nil {
		response.Fail(c, response.ErrNotFound, "企业不存在")
		return
	}
	
	// 验证企业是否属于当前租户
	tenantID, _ := c.Get("tenantId")
	if corp.TenantID != tenantID.(uint) {
		response.Fail(c, response.ErrNotFound, "企业不存在")
		return
	}
	
	// 返回企业详情信息
	response.Success(c, gin.H{
		"corpId":         corp.ID,         // 企业 ID
		"corpName":       corp.Name,       // 企业名称
		"wxCorpId":       corp.WxCorpid,   // 企业微信 ID
		"employeeSecret": corp.EmployeeSecret, // 员工权限密钥
		"contactSecret":  corp.ContactSecret,  // 客户联系权限密钥
		"eventCallback":  corp.EventCallback,  // 事件回调地址
		"token":          corp.Token,          // 令牌
		"encodingAesKey": corp.EncodingAesKey, // 加密密钥
		"tenantId":       corp.TenantID,       // 租户 ID
	})
}

// Store 创建企业
// 创建新的企业
// 处理流程：
// 1. 绑定请求参数
// 2. 获取租户 ID
// 3. 创建企业对象并设置属性
// 4. 调用服务层创建企业
// 5. 返回创建的企业信息
// 请求体（JSON）：
//
//	corpName - 企业名称（必填）
//	wxCorpId - 企业微信 ID（必填）
//	employeeSecret - 员工权限密钥（必填）
//	contactSecret - 客户联系权限密钥（必填）
//	eventCallback - 事件回调地址
//	token - 令牌
//	encodingAesKey - 加密密钥
//
// 返回：创建的企业详情信息
func (h *CorpHandler) Store(c *gin.Context) {
	// 绑定请求参数
	var req struct {
		CorpName       string `json:"corpName" binding:"required"`       // 企业名称
		WxCorpID       string `json:"wxCorpId" binding:"required"`       // 企业微信 ID
		EmployeeSecret string `json:"employeeSecret" binding:"required"` // 员工权限密钥
		ContactSecret  string `json:"contactSecret" binding:"required"`  // 客户联系权限密钥
		EventCallback  string `json:"eventCallback"`                    // 事件回调地址
		Token          string `json:"token"`                            // 令牌
		EncodingAesKey string `json:"encodingAesKey"`                   // 加密密钥
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ErrParams, "参数错误")
		return
	}
	
	// 获取租户 ID
	tenantID, _ := c.Get("tenantId")
	
	// 创建企业对象并设置属性
	corp := model.Corp{
		Name:           strings.TrimSpace(req.CorpName),
		WxCorpid:       strings.TrimSpace(req.WxCorpID),
		EmployeeSecret: strings.TrimSpace(req.EmployeeSecret),
		ContactSecret:  strings.TrimSpace(req.ContactSecret),
		EventCallback:  strings.TrimSpace(req.EventCallback),
		Token:          strings.TrimSpace(req.Token),
		EncodingAesKey: strings.TrimSpace(req.EncodingAesKey),
		TenantID:       tenantID.(uint),
	}
	
	// 调用服务层创建企业
	if err := h.svc.Create(&corp); err != nil {
		response.Fail(c, response.ErrDB, "创建企业失败")
		return
	}
	
	// 返回创建的企业信息
	response.Success(c, gin.H{
		"corpId":         corp.ID,         // 企业 ID
		"corpName":       corp.Name,       // 企业名称
		"wxCorpId":       corp.WxCorpid,   // 企业微信 ID
		"employeeSecret": corp.EmployeeSecret, // 员工权限密钥
		"contactSecret":  corp.ContactSecret,  // 客户联系权限密钥
		"eventCallback":  corp.EventCallback,  // 事件回调地址
		"token":          corp.Token,          // 令牌
		"encodingAesKey": corp.EncodingAesKey, // 加密密钥
		"tenantId":       corp.TenantID,       // 租户 ID
	})
}

// Update 更新企业
// 更新指定企业的信息
// 处理流程：
// 1. 绑定请求参数
// 2. 获取企业 ID（从请求体或路径参数）
// 3. 调用服务层获取企业详情
// 4. 验证企业是否属于当前租户
// 5. 更新企业信息
// 6. 调用服务层更新企业
// 7. 返回更新后的企业信息
// 请求体（JSON）：
//
//	corpId - 企业 ID（可选，若不提供则从路径参数获取）
//	corpName - 企业名称（必填）
//	wxCorpId - 企业微信 ID（必填）
//	employeeSecret - 员工权限密钥（必填）
//	contactSecret - 客户联系权限密钥（必填）
//	eventCallback - 事件回调地址
//	token - 令牌
//	encodingAesKey - 加密密钥
//
// 返回：更新后的企业详情信息
func (h *CorpHandler) Update(c *gin.Context) {
	// 绑定请求参数
	var req struct {
		CorpID         uint   `json:"corpId"`                         // 企业 ID
		CorpName       string `json:"corpName" binding:"required"`       // 企业名称
		WxCorpID       string `json:"wxCorpId" binding:"required"`       // 企业微信 ID
		EmployeeSecret string `json:"employeeSecret" binding:"required"` // 员工权限密钥
		ContactSecret  string `json:"contactSecret" binding:"required"`  // 客户联系权限密钥
		EventCallback  string `json:"eventCallback"`                    // 事件回调地址
		Token          string `json:"token"`                            // 令牌
		EncodingAesKey string `json:"encodingAesKey"`                   // 加密密钥
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ErrParams, "参数错误")
		return
	}
	
	// 获取企业 ID（从请求体或路径参数）
	id := req.CorpID
	if id == 0 {
		id = uint(getCorpIDParam(c))
	}
	
	// 调用服务层获取企业详情
	corp, err := h.svc.GetByID(id)
	if err != nil {
		response.Fail(c, response.ErrNotFound, "企业不存在")
		return
	}
	
	// 验证企业是否属于当前租户
	tenantID, _ := c.Get("tenantId")
	if corp.TenantID != tenantID.(uint) {
		response.Fail(c, response.ErrNotFound, "企业不存在")
		return
	}
	
	// 更新企业信息
	corp.Name = strings.TrimSpace(req.CorpName)
	corp.WxCorpid = strings.TrimSpace(req.WxCorpID)
	corp.EmployeeSecret = strings.TrimSpace(req.EmployeeSecret)
	corp.ContactSecret = strings.TrimSpace(req.ContactSecret)
	corp.EventCallback = strings.TrimSpace(req.EventCallback)
	corp.Token = strings.TrimSpace(req.Token)
	corp.EncodingAesKey = strings.TrimSpace(req.EncodingAesKey)
	corp.TenantID = tenantID.(uint)
	
	// 调用服务层更新企业
	if err := h.svc.Update(corp); err != nil {
		response.Fail(c, response.ErrDB, "更新企业失败")
		return
	}
	
	// 返回更新后的企业信息
	response.Success(c, gin.H{
		"corpId":         corp.ID,         // 企业 ID
		"corpName":       corp.Name,       // 企业名称
		"wxCorpId":       corp.WxCorpid,   // 企业微信 ID
		"employeeSecret": corp.EmployeeSecret, // 员工权限密钥
		"contactSecret":  corp.ContactSecret,  // 客户联系权限密钥
		"eventCallback":  corp.EventCallback,  // 事件回调地址
		"token":          corp.Token,          // 令牌
		"encodingAesKey": corp.EncodingAesKey, // 加密密钥
		"tenantId":       corp.TenantID,       // 租户 ID
	})
}

// Bind 绑定企业
// 绑定用户与企业的关系，将用户绑定到指定企业
// 处理流程：
// 1. 绑定请求参数
// 2. 获取用户 ID 和租户 ID
// 3. 验证企业是否存在且属于当前租户
// 4. 查询用户信息
// 5. 对于非超级管理员，验证用户是否属于该企业
// 6. 将用户绑定的企业信息缓存到 Redis
// 7. 返回绑定结果
// 请求体（JSON）：
//
//	corpId - 企业 ID（必填）
//
// 返回：绑定成功的消息
func (h *CorpHandler) Bind(c *gin.Context) {
	// 绑定请求参数
	var req struct {
		CorpID uint `json:"corpId" binding:"required"` // 企业 ID
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ErrParams, "参数错误")
		return
	}

	// 获取用户 ID 和租户 ID
	userID, _ := c.Get("userId")
	tenantID, _ := c.Get("tenantId")

	// 验证企业是否存在且属于当前租户
	corp, err := h.svc.GetByID(req.CorpID)
	if err != nil || corp.TenantID != tenantID.(uint) {
		response.Fail(c, response.ErrNotFound, "企业不存在")
		return
	}

	// 查询用户信息
	var user model.User
	if err := h.db.Select("id", "isSuperAdmin").First(&user, userID.(uint)).Error; err != nil {
		response.Fail(c, response.ErrDB, "绑定企业失败")
		return
	}

	// 对于非超级管理员，验证用户是否属于该企业
	employeeID := uint(0)
	if user.IsSuperAdmin == 0 {
		var employee model.WorkEmployee
		if err := h.db.Select("id").Where("corp_id = ? AND log_user_id = ?", req.CorpID, user.ID).First(&employee).Error; err != nil {
			response.Fail(c, response.ErrParams, "当前用户不归属该企业，不可操作")
			return
		}
		employeeID = employee.ID
	}

	// 将用户绑定的企业信息缓存到 Redis
	if mochatRedis.RDB == nil {
		response.Fail(c, response.ErrServer, "Redis 未初始化")
		return
	}
	cacheKey := fmt.Sprintf("mc:user.%d", user.ID)
	cacheVal := fmt.Sprintf("%d-%d", req.CorpID, employeeID)
	if err := mochatRedis.RDB.Set(context.Background(), cacheKey, cacheVal, 0).Err(); err != nil {
		response.Fail(c, response.ErrRedis, "绑定企业失败")
		return
	}

	// 返回绑定结果
	response.Success(c, nil)
}

// WeWorkCallback 企业微信回调处理
// 处理企业微信的回调请求，主要用于验证回调 URL
// 处理流程：
// 1. 获取 echostr 参数
// 2. 如果存在 echostr 参数，直接返回（用于验证回调 URL）
// 3. 否则返回 200 状态码
// 参数：
//
//	echostr - 企业微信回调验证参数
//
// 返回：如果有 echostr 参数则返回该参数，否则返回 200 状态码
func (h *CorpHandler) WeWorkCallback(c *gin.Context) {
	// 获取 echostr 参数
	echoStr := c.Query("echostr")
	
	// 如果存在 echostr 参数，直接返回（用于验证回调 URL）
	if echoStr != "" {
		c.String(200, echoStr)
		return
	}
	
	// 否则返回 200 状态码
	c.Status(200)
}

// getCorpIDParam 获取企业 ID 参数
// 从路径参数或查询参数中获取企业 ID
// 处理流程：
// 1. 尝试从路径参数中获取 id
// 2. 如果路径参数中没有，则尝试从查询参数中获取 corpId
// 3. 如果都没有或解析失败，返回 0
// 参数：
//
//	c - Gin 上下文
//
// 返回：企业 ID 或 0
func getCorpIDParam(c *gin.Context) uint64 {
	// 尝试从路径参数中获取 id
	if id, err := strconv.ParseUint(c.Param("id"), 10, 32); err == nil && id > 0 {
		return id
	}
	// 如果路径参数中没有，则尝试从查询参数中获取 corpId
	if id, err := strconv.ParseUint(c.Query("corpId"), 10, 32); err == nil && id > 0 {
		return id
	}
	// 如果都没有或解析失败，返回 0
	return 0
}
