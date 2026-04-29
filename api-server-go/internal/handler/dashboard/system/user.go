package system

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"mochat-api-server/internal/middleware"
	"mochat-api-server/internal/model"
	"mochat-api-server/internal/pkg/logger"
	"mochat-api-server/internal/pkg/response"
	mochatRedis "mochat-api-server/internal/redis"
	"mochat-api-server/internal/service/business"

	"github.com/gin-gonic/gin"
)

// UserHandler 用户管理处理器
// 处理用户相关的HTTP请求，包括用户认证、登录、登出、列表、详情、创建、更新、密码修改和状态更新等操作

type UserHandler struct {
	svc     *business.UserService // 用户服务实例
	corpSvc *business.CorpService // 企业服务实例
	jwtCfg  response.JWTConfig    // JWT配置
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

	// 登录后自动写入用户-企业绑定缓存
	h.ensureCorpBinding(user.ID)

	response.Success(c, gin.H{
		"token":    token,
		"userInfo": user,
	})
}

// ensureCorpBinding 确保 Redis 中存在 mc:user.{userId} 缓存
// 如果缓存不存在，从 mc_work_employee 表中查询用户的第一个企业并写入
func (h *UserHandler) ensureCorpBinding(userID uint) {
	fmt.Printf("[ensureCorpBinding] start userID=%d\n", userID)
	if mochatRedis.RDB == nil {
		fmt.Println("[ensureCorpBinding] RDB is nil, skip")
		return
	}
	cacheKey := "mc:user." + strconv.Itoa(int(userID))
	ctx := context.Background()

	// 如果缓存已存在，不覆盖
	if cached, err := mochatRedis.RDB.Get(ctx, cacheKey).Result(); err == nil && cached != "" {
		fmt.Printf("[ensureCorpBinding] cache already exists: %s\n", cached)
		return
	}
	fmt.Println("[ensureCorpBinding] cache not found, querying DB")

	// 查询该用户关联的第一个员工
	var employee model.WorkEmployee
	if err := h.corpSvc.DB().Where("log_user_id = ?", userID).First(&employee).Error; err != nil {
		fmt.Printf("[ensureCorpBinding] DB query failed: %v\n", err)
		return
	}
	fmt.Printf("[ensureCorpBinding] found employee: id=%d corpId=%d\n", employee.ID, employee.CorpID)

	if employee.CorpID == 0 || employee.ID == 0 {
		fmt.Println("[ensureCorpBinding] corpID or employeeID is 0, skip")
		return
	}

	// 写入缓存 mc:user.{userId} = {corpId}-{employeeId}
	cacheValue := strconv.Itoa(int(employee.CorpID)) + "-" + strconv.Itoa(int(employee.ID))
	err := mochatRedis.RDB.Set(ctx, cacheKey, cacheValue, 0).Err()
	if err != nil {
		fmt.Printf("[ensureCorpBinding] Redis set failed: %v\n", err)
		return
	}
	fmt.Printf("[ensureCorpBinding] Redis set success: %s = %s\n", cacheKey, cacheValue)
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

// CorpBind 绑定（选择）企业
// 请求方法: POST
// 请求路径: /corp/bind
// 请求体:
//   - corpId: 企业ID
// 响应:
//   - 成功: 绑定成功消息
//   - 失败: 错误信息

func (h *UserHandler) CorpBind(c *gin.Context) {
	userID, exists := c.Get("userId")
	if !exists {
		response.Fail(c, response.ErrAuth, "未登录")
		return
	}

	var req struct {
		CorpID uint `json:"corpId" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ErrParams, "参数错误")
		return
	}

	// 验证企业是否存在且属于当前租户
	tenantID, _ := c.Get("tenantId")
	corp, err := h.corpSvc.GetByID(req.CorpID)
	if err != nil || corp.TenantID != tenantID.(uint) {
		response.Fail(c, response.ErrNotFound, "企业不存在或无权限")
		return
	}

	// 保存用户选择的企业到 Redis
	// 格式: mc:user.{userID} = "{corpId}-{employeeId}"
	if mochatRedis.RDB != nil {
		cacheKey := "mc:user." + strconv.Itoa(int(userID.(uint)))
		cached, err := mochatRedis.RDB.Get(context.Background(), cacheKey).Result()
		logger.Sugar.Infof("CorpBind cacheKey=%s, cached=%s", cacheKey, cached)
		if err == nil && cached != "" {
			parts := strings.SplitN(cached, "-", 2)
			if len(parts) == 2 {
				// 保留原有的 employeeId，只更新 corpId
				logger.Sugar.Infof("CorpBind update corpId to %d, employeeId=%s", req.CorpID, parts[1])
				err = mochatRedis.RDB.Set(context.Background(), cacheKey, strconv.Itoa(int(req.CorpID))+"-"+parts[1], 24*3600).Err()
			} else {
				// 只有 corpId，没有 employeeId
				logger.Sugar.Infof("CorpBind set corpId to %d", req.CorpID)
				err = mochatRedis.RDB.Set(context.Background(), cacheKey, strconv.Itoa(int(req.CorpID)), 24*3600).Err()
			}
		} else {
			// 没有缓存，只设置 corpId
			logger.Sugar.Infof("CorpBind set corpId to %d", req.CorpID)
			err = mochatRedis.RDB.Set(context.Background(), cacheKey, strconv.Itoa(int(req.CorpID)), 24*3600).Err()
		}
		logger.Sugar.Infof("CorpBind Redis set success: %s = %s", cacheKey, strconv.Itoa(int(req.CorpID)))
		if err != nil {
			response.Fail(c, response.ErrServer, "设置企业失败")
			return
		}
	}

	response.SuccessMsg(c, "选择企业成功")
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
	if err != nil || id == 0 {
		userId := c.Query("userId")
		if userId == "" {
			userId = c.Query("id")
		}
		if userId != "" {
			id, _ = strconv.ParseUint(userId, 10, 64)
		}
	}
	if id == 0 {
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
	Phone      string `json:"phone" binding:"required"`    // 手机号
	Password   string `json:"password" binding:"required"` // 密码
	Name       string `json:"name"`                        // 姓名
	Gender     int    `json:"gender"`                      // 性别
	Department string `json:"department"`                  // 部门
	Position   string `json:"position"`                    // 职位
	TenantID   uint   `json:"tenantId"`                    // 租户ID
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
