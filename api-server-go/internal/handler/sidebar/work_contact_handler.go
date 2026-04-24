package sidebar

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"mochat-api-server/internal/model"
	"mochat-api-server/internal/pkg/response"
	"mochat-api-server/internal/service"
)

type WorkContactHandler struct {
	db         *gorm.DB
	contactSvc *service.WorkContactService
}

func NewWorkContactHandler(db *gorm.DB) *WorkContactHandler {
	return &WorkContactHandler{
		db:         db,
		contactSvc: service.NewWorkContactService(db),
	}
}

func (h *WorkContactHandler) Show(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	contact, err := h.contactSvc.GetByID(uint(id))
	if err != nil {
		response.Fail(c, response.ErrNotFound, "客户不存在")
		return
	}
	response.Success(c, contact)
}

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
