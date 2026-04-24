package dashboard

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"mochat-api-server/internal/model"
	"mochat-api-server/internal/pkg/response"
	"mochat-api-server/internal/service/business"
)

type CorpHandler struct {
	svc *business.CorpService
}

func NewCorpHandler(svc *business.CorpService) *CorpHandler {
	return &CorpHandler{svc: svc}
}

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

func (h *CorpHandler) Select(c *gin.Context) {
	tenantID, _ := c.Get("tenantId")
	corps, _, err := h.svc.List(tenantID.(uint), 0, 1000)
	if err != nil {
		response.Fail(c, response.ErrDB, "获取企业列表失败")
		return
	}
	response.Success(c, corps)
}

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

type StoreCorpRequest struct {
	Name           string `json:"name" binding:"required"`
	WxCorpid      string `json:"wxCorpid"`
	EmployeeSecret string `json:"employeeSecret"`
	ContactSecret  string `json:"contactSecret"`
	Token          string `json:"token"`
	EncodingAesKey string `json:"encodingAesKey"`
	TenantID       uint   `json:"tenantId"`
}

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

type BindCorpRequest struct {
	WxCorpid       string `json:"wxCorpid" binding:"required"`
	EmployeeSecret string `json:"employeeSecret" binding:"required"`
	ContactSecret  string `json:"contactSecret" binding:"required"`
	Token          string `json:"token" binding:"required"`
	EncodingAesKey string `json:"encodingAesKey" binding:"required"`
}

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

func (h *CorpHandler) WeWorkCallback(c *gin.Context) {
	c.String(http.StatusOK, "success")
}