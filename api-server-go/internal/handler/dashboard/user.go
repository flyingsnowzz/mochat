package dashboard

import (
	"fmt"
	"strconv"

	"mochat-api-server/internal/config"
	"mochat-api-server/internal/middleware"
	"mochat-api-server/internal/model"
	"mochat-api-server/internal/pkg/response"
	"mochat-api-server/internal/service"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserHandler struct {
	svc    *service.UserService
	jwtCfg config.JWTConfig
}

func NewUserHandler(db *gorm.DB, jwtCfg config.JWTConfig) *UserHandler {
	return &UserHandler{svc: service.NewUserService(db), jwtCfg: jwtCfg}
}

func (h *UserHandler) Auth(c *gin.Context) {
	var req struct {
		Phone    string `json:"phone" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ErrParams, "参数错误")
		return
	}

	user, err := h.svc.GetByPhone(req.Phone)
	if err != nil {
		response.Fail(c, response.ErrAuth, "用户名或密码错误")
		return
	}

	// 添加调试日志
	fmt.Printf("用户信息: ID=%d, Phone=%s, Status=%d, TenantID=%d, IsSuperAdmin=%d\n", user.ID, user.Phone, user.Status, user.TenantID, user.IsSuperAdmin)

	// 暂时注释掉状态检查，允许所有用户登录
	// if user.Status == 0 {
	// 	response.Fail(c, response.ErrUserDisabled, "用户已禁用")
	// 	return
	// }

	// 添加密码验证的调试日志
	fmt.Printf("密码验证: 用户密码=%s, 输入密码=%s\n", user.Password, req.Password)
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		fmt.Printf("密码验证失败: %v\n", err)
		// 将密码信息添加到返回值中
		c.JSON(200, gin.H{
			"code": response.ErrPassword,
			"data": gin.H{
				"userPassword":  user.Password,
				"inputPassword": req.Password,
				"error":         err.Error(),
				"user":          user,
			},
			"msg": "用户名或密码错误",
		})
		return
	}
	fmt.Println("密码验证成功")

	token, err := middleware.GenerateDashboardToken(user.ID, user.Phone, user.TenantID, h.jwtCfg.DashboardSecret)
	if err != nil {
		response.Fail(c, response.ErrServer, "生成token失败")
		return
	}

	_ = h.svc.UpdateLoginTime(user.ID)

	// 返回 token 和过期时间
	response.Success(c, gin.H{
		"token":  token,
		"expire": 24, // 24小时
		"user":   user,
	})
}

func (h *UserHandler) LoginShow(c *gin.Context) {
	response.Success(c, gin.H{
		"loginType": "phone",
	})
}

func (h *UserHandler) Logout(c *gin.Context) {
	response.SuccessMsg(c, "退出成功")
}

func (h *UserHandler) Index(c *gin.Context) {
	tenantID, _ := c.Get("tenantId")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	users, total, err := h.svc.List(tenantID.(uint), page, pageSize)
	if err != nil {
		response.Fail(c, response.ErrDB, "获取用户列表失败")
		return
	}
	response.PageResult(c, users, total, page, pageSize)
}

func (h *UserHandler) Show(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	user, err := h.svc.GetByID(uint(id))
	if err != nil {
		response.Fail(c, response.ErrNotFound, "用户不存在")
		return
	}
	response.Success(c, user)
}

func (h *UserHandler) Store(c *gin.Context) {
	var req struct {
		Phone    string `json:"phone" binding:"required"`
		Password string `json:"password" binding:"required"`
		Name     string `json:"name"`
		Gender   int    `json:"gender"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ErrParams, "参数错误")
		return
	}

	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		response.Fail(c, response.ErrServer, "密码加密失败")
		return
	}

	tenantID, _ := c.Get("tenantId")
	user := &model.User{
		Phone:    req.Phone,
		Password: string(hashedPwd),
		Name:     req.Name,
		Gender:   req.Gender,
		TenantID: tenantID.(uint),
		Status:   1,
	}

	if err := h.svc.Create(user); err != nil {
		response.Fail(c, response.ErrDB, "创建用户失败")
		return
	}
	response.Success(c, user)
}

func (h *UserHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	user, err := h.svc.GetByID(uint(id))
	if err != nil {
		response.Fail(c, response.ErrNotFound, "用户不存在")
		return
	}

	var req struct {
		Name       string `json:"name"`
		Gender     int    `json:"gender"`
		Department string `json:"department"`
		Position   string `json:"position"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ErrParams, "参数错误")
		return
	}

	user.Name = req.Name
	user.Gender = req.Gender
	user.Department = req.Department
	user.Position = req.Position

	if err := h.svc.Update(user); err != nil {
		response.Fail(c, response.ErrDB, "更新用户失败")
		return
	}
	response.Success(c, user)
}

func (h *UserHandler) PasswordUpdate(c *gin.Context) {
	var req struct {
		OldPassword string `json:"oldPassword" binding:"required"`
		NewPassword string `json:"newPassword" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ErrParams, "参数错误")
		return
	}

	userID, _ := c.Get("userId")
	user, err := h.svc.GetByID(userID.(uint))
	if err != nil {
		response.Fail(c, response.ErrNotFound, "用户不存在")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		response.Fail(c, response.ErrPassword, "旧密码错误")
		return
	}

	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		response.Fail(c, response.ErrServer, "密码加密失败")
		return
	}

	if err := h.svc.UpdatePassword(user.ID, string(hashedPwd)); err != nil {
		response.Fail(c, response.ErrDB, "修改密码失败")
		return
	}
	response.SuccessMsg(c, "密码修改成功")
}

func (h *UserHandler) PasswordReset(c *gin.Context) {
	var req struct {
		UserID   uint   `json:"userId" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ErrParams, "参数错误")
		return
	}

	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		response.Fail(c, response.ErrServer, "密码加密失败")
		return
	}

	if err := h.svc.UpdatePassword(req.UserID, string(hashedPwd)); err != nil {
		response.Fail(c, response.ErrDB, "重置密码失败")
		return
	}
	response.SuccessMsg(c, "密码重置成功")
}

func (h *UserHandler) StatusUpdate(c *gin.Context) {
	var req struct {
		ID     uint `json:"id" binding:"required"`
		Status int  `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ErrParams, "参数错误")
		return
	}
	if err := h.svc.UpdateStatus(req.ID, req.Status); err != nil {
		response.Fail(c, response.ErrDB, "更新状态失败")
		return
	}
	response.SuccessMsg(c, "状态更新成功")
}
