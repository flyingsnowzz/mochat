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

type CorpHandler struct {
	svc *service.CorpService
	db  *gorm.DB
}

type corpIndexItem struct {
	CorpID           uint   `json:"corpId"`
	CorpName         string `json:"corpName"`
	WxCorpID         string `json:"wxCorpId"`
	EmployeeSecret   string `json:"employeeSecret"`
	ContactSecret    string `json:"contactSecret"`
	EventCallback    string `json:"eventCallback"`
	Token            string `json:"token"`
	EncodingAesKey   string `json:"encodingAesKey"`
	CreatedAt        string `json:"createdAt"`
	ChatApplyStatus  int    `json:"chatApplyStatus"`
	ChatStatus       int    `json:"chatStatus"`
	MessageCreatedAt string `json:"messageCreatedAt"`
}

func NewCorpHandler(db *gorm.DB) *CorpHandler {
	return &CorpHandler{svc: service.NewCorpService(db), db: db}
}

func (h *CorpHandler) Index(c *gin.Context) {
	tenantID, _ := c.Get("tenantId")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("perPage", c.DefaultQuery("pageSize", "20")))
	corpName := strings.TrimSpace(c.Query("corpName"))

	corps, total, err := h.svc.List(tenantID.(uint), corpName, page, pageSize)
	if err != nil {
		response.Fail(c, response.ErrDB, "获取企业列表失败")
		return
	}

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
			ChatApplyStatus:  0,
			ChatStatus:       0,
			MessageCreatedAt: createdAt,
		})
	}
	response.PageResult(c, items, total, page, pageSize)
}

func (h *CorpHandler) Select(c *gin.Context) {
	tenantID, _ := c.Get("tenantId")
	userID, _ := c.Get("userId")

	var user model.User
	if err := h.db.Select("id", "isSuperAdmin").First(&user, userID.(uint)).Error; err != nil {
		response.Fail(c, response.ErrDB, "获取企业选择列表失败")
		return
	}

	type corpOption struct {
		CorpID   uint   `json:"corpId"`
		CorpName string `json:"corpName"`
	}
	var data []corpOption

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

	response.Success(c, data)
}

func (h *CorpHandler) Show(c *gin.Context) {
	id := getCorpIDParam(c)
	corp, err := h.svc.GetByID(uint(id))
	if err != nil {
		response.Fail(c, response.ErrNotFound, "企业不存在")
		return
	}
	tenantID, _ := c.Get("tenantId")
	if corp.TenantID != tenantID.(uint) {
		response.Fail(c, response.ErrNotFound, "企业不存在")
		return
	}
	response.Success(c, gin.H{
		"corpId":         corp.ID,
		"corpName":       corp.Name,
		"wxCorpId":       corp.WxCorpid,
		"employeeSecret": corp.EmployeeSecret,
		"contactSecret":  corp.ContactSecret,
		"eventCallback":  corp.EventCallback,
		"token":          corp.Token,
		"encodingAesKey": corp.EncodingAesKey,
		"tenantId":       corp.TenantID,
	})
}

func (h *CorpHandler) Store(c *gin.Context) {
	var req struct {
		CorpName       string `json:"corpName" binding:"required"`
		WxCorpID       string `json:"wxCorpId" binding:"required"`
		EmployeeSecret string `json:"employeeSecret" binding:"required"`
		ContactSecret  string `json:"contactSecret" binding:"required"`
		EventCallback  string `json:"eventCallback"`
		Token          string `json:"token"`
		EncodingAesKey string `json:"encodingAesKey"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ErrParams, "参数错误")
		return
	}
	tenantID, _ := c.Get("tenantId")
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
	if err := h.svc.Create(&corp); err != nil {
		response.Fail(c, response.ErrDB, "创建企业失败")
		return
	}
	response.Success(c, gin.H{
		"corpId":         corp.ID,
		"corpName":       corp.Name,
		"wxCorpId":       corp.WxCorpid,
		"employeeSecret": corp.EmployeeSecret,
		"contactSecret":  corp.ContactSecret,
		"eventCallback":  corp.EventCallback,
		"token":          corp.Token,
		"encodingAesKey": corp.EncodingAesKey,
		"tenantId":       corp.TenantID,
	})
}

func (h *CorpHandler) Update(c *gin.Context) {
	var req struct {
		CorpID         uint   `json:"corpId"`
		CorpName       string `json:"corpName" binding:"required"`
		WxCorpID       string `json:"wxCorpId" binding:"required"`
		EmployeeSecret string `json:"employeeSecret" binding:"required"`
		ContactSecret  string `json:"contactSecret" binding:"required"`
		EventCallback  string `json:"eventCallback"`
		Token          string `json:"token"`
		EncodingAesKey string `json:"encodingAesKey"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ErrParams, "参数错误")
		return
	}
	id := req.CorpID
	if id == 0 {
		id = uint(getCorpIDParam(c))
	}
	corp, err := h.svc.GetByID(id)
	if err != nil {
		response.Fail(c, response.ErrNotFound, "企业不存在")
		return
	}
	tenantID, _ := c.Get("tenantId")
	if corp.TenantID != tenantID.(uint) {
		response.Fail(c, response.ErrNotFound, "企业不存在")
		return
	}
	corp.Name = strings.TrimSpace(req.CorpName)
	corp.WxCorpid = strings.TrimSpace(req.WxCorpID)
	corp.EmployeeSecret = strings.TrimSpace(req.EmployeeSecret)
	corp.ContactSecret = strings.TrimSpace(req.ContactSecret)
	corp.EventCallback = strings.TrimSpace(req.EventCallback)
	corp.Token = strings.TrimSpace(req.Token)
	corp.EncodingAesKey = strings.TrimSpace(req.EncodingAesKey)
	corp.TenantID = tenantID.(uint)
	if err := h.svc.Update(corp); err != nil {
		response.Fail(c, response.ErrDB, "更新企业失败")
		return
	}
	response.Success(c, gin.H{
		"corpId":         corp.ID,
		"corpName":       corp.Name,
		"wxCorpId":       corp.WxCorpid,
		"employeeSecret": corp.EmployeeSecret,
		"contactSecret":  corp.ContactSecret,
		"eventCallback":  corp.EventCallback,
		"token":          corp.Token,
		"encodingAesKey": corp.EncodingAesKey,
		"tenantId":       corp.TenantID,
	})
}

func (h *CorpHandler) Bind(c *gin.Context) {
	var req struct {
		CorpID uint `json:"corpId" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ErrParams, "参数错误")
		return
	}

	userID, _ := c.Get("userId")
	tenantID, _ := c.Get("tenantId")

	corp, err := h.svc.GetByID(req.CorpID)
	if err != nil || corp.TenantID != tenantID.(uint) {
		response.Fail(c, response.ErrNotFound, "企业不存在")
		return
	}

	var user model.User
	if err := h.db.Select("id", "isSuperAdmin").First(&user, userID.(uint)).Error; err != nil {
		response.Fail(c, response.ErrDB, "绑定企业失败")
		return
	}

	employeeID := uint(0)
	if user.IsSuperAdmin == 0 {
		var employee model.WorkEmployee
		if err := h.db.Select("id").Where("corp_id = ? AND log_user_id = ?", req.CorpID, user.ID).First(&employee).Error; err != nil {
			response.Fail(c, response.ErrParams, "当前用户不归属该企业，不可操作")
			return
		}
		employeeID = employee.ID
	}

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

	response.Success(c, nil)
}

func (h *CorpHandler) WeWorkCallback(c *gin.Context) {
	echoStr := c.Query("echostr")
	if echoStr != "" {
		c.String(200, echoStr)
		return
	}
	c.Status(200)
}

func getCorpIDParam(c *gin.Context) uint64 {
	if id, err := strconv.ParseUint(c.Param("id"), 10, 32); err == nil && id > 0 {
		return id
	}
	if id, err := strconv.ParseUint(c.Query("corpId"), 10, 32); err == nil && id > 0 {
		return id
	}
	return 0
}
