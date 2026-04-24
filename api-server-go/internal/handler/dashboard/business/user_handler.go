package dashboard

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"mochat-api-server/internal/middleware"
	"mochat-api-server/internal/model"
	"mochat-api-server/internal/pkg/response"
	"mochat-api-server/internal/service/business"
)

type UserHandler struct {
	svc     *business.UserService
	corpSvc *business.CorpService
	jwtCfg  response.JWTConfig
}

func NewUserHandler(svc *business.UserService, corpSvc *business.CorpService, jwtCfg response.JWTConfig) *UserHandler {
	return &UserHandler{svc: svc, corpSvc: corpSvc, jwtCfg: jwtCfg}
}

func (h *UserHandler) Auth(c *gin.Context) {
	var req struct {
		Phone    string `json:"phone" binding:"required"`
		Password string `json:"password" binding:"required"`
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

func (h *UserHandler) LoginShow(c *gin.Context) {
	response.Success(c, gin.H{
		"copyright": "MoChat",
		"logo":      "/static/image/logo.png",
	})
}

func (h *UserHandler) Logout(c *gin.Context) {
	response.SuccessMsg(c, "退出成功")
}

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

type StoreUserRequest struct {
	Phone      string `json:"phone" binding:"required"`
	Password   string `json:"password" binding:"required"`
	Name       string `json:"name"`
	Gender     int    `json:"gender"`
	Department string `json:"department"`
	Position   string `json:"position"`
	TenantID   uint   `json:"tenantId"`
}

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

func (h *UserHandler) PasswordUpdate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, response.ErrParams, "参数错误")
		return
	}
	var req struct {
		Password string `json:"password" binding:"required"`
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

func (h *UserHandler) StatusUpdate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, response.ErrParams, "参数错误")
		return
	}
	var req struct {
		Status int `json:"status" binding:"required"`
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