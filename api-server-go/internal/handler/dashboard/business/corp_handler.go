package dashboard

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"mochat-api-server/internal/model"
	"mochat-api-server/internal/pkg/response"
	"mochat-api-server/internal/service/business"
)

// CorpHandler 企业管理处理器
// 处理企业相关的HTTP请求，包括企业列表、详情、创建、更新和绑定等操作

type CorpHandler struct {
	svc *business.CorpService // 企业服务实例
}

// NewCorpHandler 创建企业处理器实例
// 参数:
//   - svc: 企业服务实例
// 返回值:
//   - *CorpHandler: 企业处理器实例

func NewCorpHandler(svc *business.CorpService) *CorpHandler {
	return &CorpHandler{svc: svc}
}

// Index 获取企业列表
// 请求方法: GET
// 请求路径: /dashboard/corp/index
// 请求参数:
//   - page: 页码，默认1
//   - pageSize: 每页数量，默认20
// 响应:
//   - 成功: 企业列表数据，包含分页信息
//   - 失败: 错误信息

func (h *CorpHandler) Index(c *gin.Context) {
	tenantID, _ := c.Get("tenantId")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	corps, total, err := h.svc.List(tenantID.(uint), offset, pageSize)
	if err != nil {
		response.Fail(c, response.ErrDB, "获取企业列表失败")
		return
	}
	response.PageResult(c, corps, total, page, pageSize)
}

// Select 获取企业选择列表
// 请求方法: GET
// 请求路径: /dashboard/corp/select
// 响应:
//   - 成功: 企业列表数据
//   - 失败: 错误信息

func (h *CorpHandler) Select(c *gin.Context) {
	tenantID, _ := c.Get("tenantId")
	corps, _, err := h.svc.List(tenantID.(uint), 0, 1000)
	if err != nil {
		response.Fail(c, response.ErrDB, "获取企业列表失败")
		return
	}
	response.Success(c, corps)
}

// Show 获取企业详情
// 请求方法: GET
// 请求路径: /dashboard/corp/show/:id
// 请求参数:
//   - id: 企业ID
// 响应:
//   - 成功: 企业详情数据
//   - 失败: 错误信息

func (h *CorpHandler) Show(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, response.ErrParams, "参数错误")
		return
	}
	corp, err := h.svc.GetByID(uint(id))
	if err != nil {
		response.Fail(c, response.ErrNotFound, "企业不存在")
		return
	}
	response.Success(c, corp)
}

// StoreCorpRequest 创建企业请求结构
// 包含创建企业所需的字段

type StoreCorpRequest struct {
	Name           string `json:"name" binding:"required"`           // 企业名称
	WxCorpid      string `json:"wxCorpid"`                           // 企业微信CorpID
	EmployeeSecret string `json:"employeeSecret"`                    // 企业微信员工权限密钥
	ContactSecret  string `json:"contactSecret"`                     // 企业微信通讯录权限密钥
	Token          string `json:"token"`                             // 企业微信回调token
	EncodingAesKey string `json:"encodingAesKey"`                    // 企业微信消息加密密钥
	TenantID       uint   `json:"tenantId"`                         // 租户ID
}

// Store 创建企业
// 请求方法: POST
// 请求路径: /dashboard/corp/store
// 请求体:
//   - StoreCorpRequest 结构
// 响应:
//   - 成功: 创建的企业数据
//   - 失败: 错误信息

func (h *CorpHandler) Store(c *gin.Context) {
	var req StoreCorpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ErrParams, err.Error())
		return
	}
	corp := &model.Corp{
		Name:           req.Name,
		WxCorpid:       req.WxCorpid,
		EmployeeSecret: req.EmployeeSecret,
		ContactSecret:  req.ContactSecret,
		Token:          req.Token,
		EncodingAesKey:  req.EncodingAesKey,
		TenantID:       req.TenantID,
	}
	if err := h.svc.Create(corp); err != nil {
		response.Fail(c, response.ErrDB, "创建企业失败")
		return
	}
	response.Success(c, corp)
}

// Update 更新企业
// 请求方法: PUT
// 请求路径: /dashboard/corp/update/:id
// 请求参数:
//   - id: 企业ID
// 请求体:
//   - 要更新的字段和值
// 响应:
//   - 成功: 更新成功消息
//   - 失败: 错误信息

func (h *CorpHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, response.ErrParams, "参数错误")
		return
	}
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ErrParams, err.Error())
		return
	}
	if err := h.svc.Update(uint(id), req); err != nil {
		response.Fail(c, response.ErrDB, "更新企业失败")
		return
	}
	response.SuccessMsg(c, "更新成功")
}

// BindCorpRequest 绑定企业微信请求结构
// 包含绑定企业微信所需的字段

type BindCorpRequest struct {
	WxCorpid       string `json:"wxCorpid" binding:"required"`       // 企业微信CorpID
	EmployeeSecret string `json:"employeeSecret" binding:"required"` // 企业微信员工权限密钥
	ContactSecret  string `json:"contactSecret" binding:"required"`  // 企业微信通讯录权限密钥
	Token          string `json:"token" binding:"required"`          // 企业微信回调token
	EncodingAesKey string `json:"encodingAesKey" binding:"required"` // 企业微信消息加密密钥
}

// Bind 绑定企业微信
// 请求方法: POST
// 请求路径: /dashboard/corp/bind/:id
// 请求参数:
//   - id: 企业ID
// 请求体:
//   - BindCorpRequest 结构
// 响应:
//   - 成功: 绑定成功消息
//   - 失败: 错误信息

func (h *CorpHandler) Bind(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, response.ErrParams, "参数错误")
		return
	}
	var req BindCorpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ErrParams, err.Error())
		return
	}
	updates := map[string]interface{}{
		"wx_corpid":        req.WxCorpid,
		"employee_secret":  req.EmployeeSecret,
		"contact_secret":   req.ContactSecret,
		"token":            req.Token,
		"encoding_aes_key":  req.EncodingAesKey,
	}
	if err := h.svc.Update(uint(id), updates); err != nil {
		response.Fail(c, response.ErrDB, "绑定失败")
		return
	}
	response.SuccessMsg(c, "绑定成功")
}

// WeWorkCallback 企业微信回调
// 请求方法: GET/POST
// 请求路径: /dashboard/corp/weWorkCallback
// 响应:
//   - 成功: "success"

func (h *CorpHandler) WeWorkCallback(c *gin.Context) {
	c.String(http.StatusOK, "success")
}