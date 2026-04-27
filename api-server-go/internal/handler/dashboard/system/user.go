package system

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"mochat-api-server/internal/middleware"
	"mochat-api-server/internal/model"
	"mochat-api-server/internal/pkg/response"
	"mochat-api-server/internal/service/business"
)

// UserHandler 用户管理处理器
// 处理用户相关的HTTP请求，包括用户认证、登录、登出、列表、详情、创建、更新、密码修改和状态更新等操作

type UserHandler struct {
	svc     *business.UserService    // 用户服务实例
	corpSvc *business.CorpService    // 企业服务实例
	jwtCfg  response.JWTConfig       // JWT配置
}

// NewUserHandler 创建用户处理器实例
// 参数:
//   - svc: 用户服务实例
//   - corpSvc: 企业服务实例
//   - jwtCfg: JWT配置
// 返回值:
//   - *UserHandler: 用户处理器实例

func NewUserHandler(svc *business.UserService, corpSvc *business.CorpService, jwtCfg response.JWTConfig) *UserHandler {
	return &UserHandler{svc: svc, corpSvc: corpSvc, jwtCfg: jwtCfg}
}

// Auth 用户认证
// 请求方法: POST
// 请求路径: /dashboard/user/auth
// 请求体:
//   - phone: 手机号
//   - password: 密码
// 响应:
//   - 成功: 包含token和用户信息的对象
//   - 失败: 错误信息

func (h *UserHandler) Auth(c *gin.Context) {
	var req struct {
		Phone    string `json:"phone" binding:"required"`    // 手机号
		Password string `json:"password" binding:"required"` // 密码
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ErrParams, err.Error())
		return
	}
	user, err := h.svc.VerifyPassword(req.Phone, req.Password)
	if err != nil {
		response.Fail(c, response.ErrPassword, err.Error())
		return
	}
	token, err := middleware.GenerateDashboardToken(user.ID, user.Phone, user.TenantID, h.jwtCfg.DashboardSecret)
	if err != nil {
		response.Fail(c, response.ErrServer, "生成令牌失败")
		return
	}
	h.svc.UpdateLoginTime(user.ID)
	response.Success(c, gin.H{
		"token":    token,
		"userInfo": user,
	})
}

// LoginShow 登录页面信息
// 请求方法: GET
// 请求路径: /dashboard/user/loginShow
// 响应:
//   - 成功: 包含版权和logo信息的对象

func (h *UserHandler) LoginShow(c *gin.Context) {
	response.Success(c, gin.H{
		"copyright": "MoChat",
		"logo":      "/static/image/logo.png",
	})
}

// Logout 登出
// 请求方法: POST
// 请求路径: /dashboard/user/logout
// 响应:
//   - 成功: 退出成功消息

func (h *UserHandler) Logout(c *gin.Context) {
	response.SuccessMsg(c, "退出成功")
}

// Index 获取用户列表
// 请求方法: GET
// 请求路径: /dashboard/user/index
// 请求参数:
//   - page: 页码，默认1
//   - pageSize: 每页数量，默认20
// 响应:
//   - 成功: 用户列表数据，包含分页信息
//   - 失败: 错误信息

func (h *UserHandler) Index(c *gin.Context) {
	tenantID, _ := c.Get("tenantId")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	offset := (page - 1) * pageSize

	users, total, err := h.svc.List(tenantID.(uint), offset, pageSize)
	if err != nil {
		response.Fail(c, response.ErrDB, "获取用户列表失败")
		return
	}
	response.PageResult(c, users, total, page, pageSize)
}

// Show 获取用户详情
// 请求方法: GET
// 请求路径: /dashboard/user/show/:id
// 请求参数:
//   - id: 用户ID
// 响应:
//   - 成功: 用户详情数据
//   - 失败: 错误信息

func (h *UserHandler) Show(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, response.ErrParams, "参数错误")
		return
	}
	user, err := h.svc.GetByID(uint(id))
	if err != nil {
		response.Fail(c, response.ErrNotFound, "用户不存在")
		return
	}
	response.Success(c, user)
}

// StoreUserRequest 创建用户请求结构
// 包含创建用户所需的字段

type StoreUserRequest struct {
	Phone      string `json:"phone" binding:"required"` // 手机号
	Password   string `json:"password" binding:"required"` // 密码
	Name       string `json:"name"` // 姓名
	Gender     int    `json:"gender"` // 性别
	Department string `json:"department"` // 部门
	Position   string `json:"position"` // 职位
	TenantID   uint   `json:"tenantId"` // 租户ID
}

// Store 创建用户
// 请求方法: POST
// 请求路径: /dashboard/user/store
// 请求体:
//   - StoreUserRequest 结构
// 响应:
//   - 成功: 创建的用户数据
//   - 失败: 错误信息

func (h *UserHandler) Store(c *gin.Context) {
	var req StoreUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ErrParams, err.Error())
		return
	}
	user := &model.User{
		Phone:      req.Phone,
		Password:   req.Password,
		Name:       req.Name,
		Gender:     req.Gender,
		Department: req.Department,
		Position:   req.Position,
		TenantID:   req.TenantID,
		Status:     1,
	}
	if err := h.svc.Create(user); err != nil {
		response.Fail(c, response.ErrDB, "创建用户失败")
		return
	}
	response.Success(c, user)
}

// Update 更新用户
// 请求方法: PUT
// 请求路径: /dashboard/user/update/:id
// 请求参数:
//   - id: 用户ID
// 请求体:
//   - 要更新的字段和值
// 响应:
//   - 成功: 更新成功消息
//   - 失败: 错误信息

func (h *UserHandler) Update(c *gin.Context) {
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
		response.Fail(c, response.ErrDB, "更新用户失败")
		return
	}
	response.SuccessMsg(c, "更新成功")
}

// PasswordUpdate 修改密码
// 请求方法: PUT
// 请求路径: /dashboard/user/passwordUpdate/:id
// 请求参数:
//   - id: 用户ID
// 请求体:
//   - password: 新密码
// 响应:
//   - 成功: 修改成功消息
//   - 失败: 错误信息

func (h *UserHandler) PasswordUpdate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, response.ErrParams, "参数错误")
		return
	}
	var req struct {
		Password string `json:"password" binding:"required"` // 新密码
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ErrParams, err.Error())
		return
	}
	if err := h.svc.Update(uint(id), map[string]interface{}{"password": req.Password}); err != nil {
		response.Fail(c, response.ErrDB, "修改密码失败")
		return
	}
	response.SuccessMsg(c, "修改成功")
}

// PasswordReset 重置密码
// 请求方法: PUT
// 请求路径: /dashboard/user/passwordReset/:id
// 请求参数:
//   - id: 用户ID
// 响应:
//   - 成功: 重置成功消息，新密码为123456
//   - 失败: 错误信息

func (h *UserHandler) PasswordReset(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, response.ErrParams, "参数错误")
		return
	}
	if err := h.svc.Update(uint(id), map[string]interface{}{"password": "123456"}); err != nil {
		response.Fail(c, response.ErrDB, "重置密码失败")
		return
	}
	response.SuccessMsg(c, "重置成功，新密码为123456")
}

// StatusUpdate 更新用户状态
// 请求方法: PUT
// 请求路径: /dashboard/user/statusUpdate/:id
// 请求参数:
//   - id: 用户ID
// 请求体:
//   - status: 状态，1为启用，0为禁用
// 响应:
//   - 成功: 更新成功消息
//   - 失败: 错误信息

func (h *UserHandler) StatusUpdate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, response.ErrParams, "参数错误")
		return
	}
	var req struct {
		Status int `json:"status" binding:"required"` // 状态
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ErrParams, err.Error())
		return
	}
	if err := h.svc.Update(uint(id), map[string]interface{}{"status": req.Status}); err != nil {
		response.Fail(c, response.ErrDB, "更新状态失败")
		return
	}
	response.SuccessMsg(c, "更新成功")
}