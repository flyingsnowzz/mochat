package sidebar

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"mochat-api-server/internal/model"
	"mochat-api-server/internal/pkg/response"
	"mochat-api-server/internal/service"
)

// WorkContactHandler 侧边栏客户处理器
// 处理侧边栏相关的客户操作，包括客户详情、客户详细信息、客户更新和客户轨迹等

type WorkContactHandler struct {
	db         *gorm.DB                // 数据库连接
	contactSvc *service.WorkContactService // 客户服务实例
}

// NewWorkContactHandler 创建客户处理器实例
// 参数:
//   - db: 数据库连接
// 返回值:
//   - *WorkContactHandler: 客户处理器实例

func NewWorkContactHandler(db *gorm.DB) *WorkContactHandler {
	return &WorkContactHandler{
		db:         db,
		contactSvc: service.NewWorkContactService(db),
	}
}

// Show 获取客户详情
// 请求方法: GET
// 请求路径: /sidebar/workContact/show/:id
// 请求参数:
//   - id: 客户ID
// 响应:
//   - 成功: 客户详情数据
//   - 失败: 错误信息

func (h *WorkContactHandler) Show(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	contact, err := h.contactSvc.GetByID(uint(id))
	if err != nil {
		response.Fail(c, response.ErrNotFound, "客户不存在")
		return
	}
	response.Success(c, contact)
}

// Detail 获取客户详细信息
// 请求方法: GET
// 请求路径: /sidebar/workContact/detail/:id
// 请求参数:
//   - id: 客户ID
// 响应:
//   - 成功: 包含客户详情、标签和轨迹的对象
//   - 失败: 错误信息

func (h *WorkContactHandler) Detail(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	contact, err := h.contactSvc.GetByID(uint(id))
	if err != nil {
		response.Fail(c, response.ErrNotFound, "客户不存在")
		return
	}

	// 获取客户标签
	var tags []model.WorkContactTag
	h.db.Model(&model.WorkContactTag{}).Joins("JOIN work_contact_tag_pivot ON work_contact_tag.id = work_contact_tag_pivot.contact_tag_id").Where("work_contact_tag_pivot.contact_id = ?", uint(id)).Find(&tags)

	// 获取客户轨迹
	var tracks []model.ContactEmployeeTrack
	h.db.Where("contact_id = ?", uint(id)).Order("created_at DESC").Limit(20).Find(&tracks)

	response.Success(c, gin.H{
		"contact": contact,
		"tags":    tags,
		"tracks":  tracks,
	})
}

// Update 更新客户信息
// 请求方法: PUT
// 请求路径: /sidebar/workContact/update/:id
// 请求参数:
//   - id: 客户ID
// 请求体:
//   - 客户信息
// 响应:
//   - 成功: 更新成功消息
//   - 失败: 错误信息

func (h *WorkContactHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	contact, err := h.contactSvc.GetByID(uint(id))
	if err != nil {
		response.Fail(c, response.ErrNotFound, "客户不存在")
		return
	}

	if err := c.ShouldBindJSON(contact); err != nil {
		response.Fail(c, response.ErrParams, "参数错误")
		return
	}

	if err := h.contactSvc.Update(contact); err != nil {
		response.Fail(c, response.ErrDB, "更新客户失败")
		return
	}

	response.SuccessMsg(c, "更新成功")
}

// Track 获取客户轨迹
// 请求方法: GET
// 请求路径: /sidebar/workContact/track/:id
// 请求参数:
//   - id: 客户ID
// 响应:
//   - 成功: 包含客户轨迹列表的对象
//   - 失败: 错误信息

func (h *WorkContactHandler) Track(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	// 获取客户轨迹记录
	var tracks []model.ContactEmployeeTrack
	result := h.db.Where("contact_id = ?", uint(id)).Order("created_at DESC").Find(&tracks)
	if result.Error != nil {
		response.Fail(c, response.ErrDB, "获取客户轨迹失败")
		return
	}

	response.Success(c, gin.H{"list": tracks})
}
