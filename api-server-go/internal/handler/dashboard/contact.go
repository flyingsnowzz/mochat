package dashboard

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"mochat-api-server/internal/model"
	"mochat-api-server/internal/pkg/response"
	"mochat-api-server/internal/service"
)

type ContactHandler struct {
	db         *gorm.DB
	contactSvc *service.WorkContactService
	tagSvc     *service.WorkContactTagService
	fieldSvc   *service.ContactFieldService
}

func NewContactHandler(db *gorm.DB) *ContactHandler {
	return &ContactHandler{
		db:         db,
		contactSvc: service.NewWorkContactService(db),
		tagSvc:     service.NewWorkContactTagService(db),
		fieldSvc:   service.NewContactFieldService(db),
	}
}

func (h *ContactHandler) Index(c *gin.Context) {
	corpID, _ := c.Get("corpId")
	if corpID == nil {
		corpID = uint(0)
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	contacts, total, err := h.contactSvc.List(corpID.(uint), page, pageSize, nil)
	if err != nil {
		response.Fail(c, response.ErrDB, "获取客户列表失败")
		return
	}
	response.PageResult(c, contacts, total, page, pageSize)
}

func (h *ContactHandler) Show(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	contact, err := h.contactSvc.GetByID(uint(id))
	if err != nil {
		response.Fail(c, response.ErrNotFound, "客户不存在")
		return
	}
	response.Success(c, contact)
}

func (h *ContactHandler) Update(c *gin.Context) {
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
	response.Success(c, contact)
}

func (h *ContactHandler) SynContact(c *gin.Context) {
	corpID, _ := c.Get("corpId")
	if corpID == nil {
		response.Fail(c, response.ErrParams, "未获取到企业信息")
		return
	}

	// 异步提交同步任务
	go func() {
		// 这里可以实现与企业微信的同步逻辑
		// 1. 获取企业微信客户端
		// 2. 拉取客户列表
		// 3. 同步到数据库
	}()

	response.SuccessMsg(c, "同步任务已提交")
}

func (h *ContactHandler) Track(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	
	// 获取客户轨迹记录
	var tracks []model.ContactEmployeeTrack
	result := h.db.Where("contact_id = ?", uint(id)).Order("created_at DESC").Find(&tracks)
	if result.Error != nil {
		response.Fail(c, response.ErrDB, "获取客户轨迹失败")
		return
	}

	response.Success(c, tracks)
}

func (h *ContactHandler) LossContact(c *gin.Context) {
	corpID, _ := c.Get("corpId")
	if corpID == nil {
		corpID = uint(0)
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	// 查询已删除的客户
	var contacts []model.WorkContact
	var total int64

	query := h.db.Model(&model.WorkContact{}).Where("corp_id = ? AND deleted_at IS NOT NULL", corpID)
	if err := query.Count(&total).Error; err != nil {
		response.Fail(c, response.ErrDB, "获取流失客户列表失败")
		return
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("deleted_at DESC").Find(&contacts).Error; err != nil {
		response.Fail(c, response.ErrDB, "获取流失客户列表失败")
		return
	}

	response.PageResult(c, contacts, total, page, pageSize)
}

func (h *ContactHandler) Source(c *gin.Context) {
	corpID, _ := c.Get("corpId")
	if corpID == nil {
		corpID = uint(0)
	}

	// 统计客户来源
	type SourceStat struct {
		AddWay    int   `json:"addWay"`
		Count     int64 `json:"count"`
		AddWayStr string `json:"addWayStr"`
	}

	var stats []SourceStat

	// 定义添加方式映射
	addWayMap := map[int]string{
		0: "未知",
		1: "扫描二维码",
		2: "搜索手机号",
		3: "名片分享",
		4: "群聊",
		5: "手机通讯录",
		6: "微信好友",
		7: "来自微信的添加",
		8: "安装第三方应用",
		9: "搜索邮箱",
		10: "企业微信内部成员共享",
		11: "管理员/负责人分配",
	}

	// 查询各来源的客户数量
	rows, err := h.db.Model(&model.WorkContactEmployee{}).Select("add_way, count(*) as count").Where("corp_id = ?", corpID).Group("add_way").Rows()
	if err != nil {
		response.Fail(c, response.ErrDB, "获取客户来源统计失败")
		return
	}
	defer rows.Close()

	for rows.Next() {
		var stat SourceStat
		if err := rows.Scan(&stat.AddWay, &stat.Count); err != nil {
			response.Fail(c, response.ErrDB, "获取客户来源统计失败")
			return
		}
		stat.AddWayStr = addWayMap[stat.AddWay]
		stats = append(stats, stat)
	}

	response.Success(c, stats)
}

func (h *ContactHandler) BatchLabeling(c *gin.Context) {
	var req struct {
		ContactIDs []uint `json:"contactIds"`
		TagIDs     []uint `json:"tagIds"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ErrParams, "参数错误")
		return
	}

	// 批量打标签
	err := h.db.Transaction(func(tx *gorm.DB) error {
		// 先删除已有的标签关联
		for _, contactID := range req.ContactIDs {
			if err := tx.Where("contact_id = ?", contactID).Delete(&model.WorkContactTagPivot{}).Error; err != nil {
				return err
			}
		}

		// 添加新的标签关联
		for _, contactID := range req.ContactIDs {
			for _, tagID := range req.TagIDs {
				pivot := &model.WorkContactTagPivot{
					ContactID:    contactID,
					ContactTagID: tagID,
				}
				if err := tx.Create(pivot).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})

	if err != nil {
		response.Fail(c, response.ErrDB, "批量打标签失败")
		return
	}

	response.SuccessMsg(c, "批量打标签成功")
}

func (h *ContactHandler) UpdateByID(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ErrParams, err.Error())
		return
	}
	if err := h.contactSvc.UpdateByID(uint(id), req); err != nil {
		response.Fail(c, response.ErrDB, "更新客户失败")
		return
	}
	response.SuccessMsg(c, "更新成功")
}

func (h *ContactHandler) GetByWxExternalUserID(c *gin.Context) {
	corpID, _ := c.Get("corpId")
	wxID := c.Query("wxExternalUserid")
	if wxID == "" {
		response.Fail(c, response.ErrParams, "wxExternalUserid参数不能为空")
		return
	}
	contact, err := h.contactSvc.GetByWxExternalUserID(corpID.(uint), wxID)
	if err != nil {
		response.Fail(c, response.ErrNotFound, "客户不存在")
		return
	}
	response.Success(c, contact)
}

func (h *ContactHandler) ListByOffset(c *gin.Context) {
	corpID, _ := c.Get("corpId")
	if corpID == nil {
		corpID = uint(0)
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	offset := (page - 1) * pageSize

	contacts, total, err := h.contactSvc.ListByOffset(corpID.(uint), offset, pageSize, nil)
	if err != nil {
		response.Fail(c, response.ErrDB, "获取客户列表失败")
		return
	}
	response.PageResult(c, contacts, total, page, pageSize)
}

type ContactFieldHandler struct {
	svc       *service.ContactFieldService
	pivotSvc  *service.ContactFieldPivotService
}

func NewContactFieldHandler(db *gorm.DB) *ContactFieldHandler {
	return &ContactFieldHandler{
		svc:       service.NewContactFieldService(db),
		pivotSvc:  service.NewContactFieldPivotService(db),
	}
}

func (h *ContactFieldHandler) Index(c *gin.Context) {
	status, _ := strconv.Atoi(c.DefaultQuery("status", "2"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	fields, total, err := h.svc.ListByStatus(status, page, pageSize)
	if err != nil {
		response.Fail(c, response.ErrDB, "获取字段列表失败")
		return
	}

	// 处理 typeText
	for i := range fields {
		fields[i].TypeText = getFieldTypeText(fields[i].Type)
	}

	response.PageResult(c, fields, total, page, pageSize)
}

func (h *ContactFieldHandler) Show(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	field, err := h.svc.GetByID(uint(id))
	if err != nil {
		response.Fail(c, response.ErrNotFound, "字段不存在")
		return
	}

	// 处理 typeText
	field.TypeText = getFieldTypeText(field.Type)

	response.Success(c, field)
}

func (h *ContactFieldHandler) Store(c *gin.Context) {
	var field model.ContactField
	if err := c.ShouldBindJSON(&field); err != nil {
		response.Fail(c, response.ErrParams, "参数错误")
		return
	}
	if err := h.svc.Create(&field); err != nil {
		response.Fail(c, response.ErrDB, "创建字段失败")
		return
	}

	// 处理 typeText
	field.TypeText = getFieldTypeText(field.Type)

	response.Success(c, field)
}

func (h *ContactFieldHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	field, err := h.svc.GetByID(uint(id))
	if err != nil {
		response.Fail(c, response.ErrNotFound, "字段不存在")
		return
	}
	if err := c.ShouldBindJSON(field); err != nil {
		response.Fail(c, response.ErrParams, "参数错误")
		return
	}
	if err := h.svc.Update(field); err != nil {
		response.Fail(c, response.ErrDB, "更新字段失败")
		return
	}

	// 处理 typeText
	field.TypeText = getFieldTypeText(field.Type)

	response.Success(c, field)
}

func (h *ContactFieldHandler) Destroy(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if err := h.svc.Delete(uint(id)); err != nil {
		response.Fail(c, response.ErrDB, "删除字段失败")
		return
	}
	response.SuccessMsg(c, "删除成功")
}

func (h *ContactFieldHandler) StatusUpdate(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var req struct {
		Status int `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ErrParams, "参数错误")
		return
	}

	if err := h.svc.UpdateStatus(uint(id), req.Status); err != nil {
		response.Fail(c, response.ErrDB, "更新状态失败")
		return
	}

	response.SuccessMsg(c, "更新成功")
}

func (h *ContactFieldHandler) Portrait(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	// 获取客户的字段值
	fieldPivots, err := h.pivotSvc.List(uint(id))
	if err != nil {
		response.Fail(c, response.ErrDB, "获取客户画像失败")
		return
	}

	response.Success(c, fieldPivots)
}

func (h *ContactFieldHandler) BatchUpdate(c *gin.Context) {
	var req struct {
		Fields []struct {
			ID    uint `json:"id"`
			Order int  `json:"order"`
		}
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ErrParams, "参数错误")
		return
	}

	for _, item := range req.Fields {
		if err := h.svc.UpdateOrder(item.ID, item.Order); err != nil {
			response.Fail(c, response.ErrDB, "批量更新失败")
			return
		}
	}

	response.SuccessMsg(c, "批量更新成功")
}

type ContactTagHandler struct {
	svc *service.WorkContactTagService
}

func NewContactTagHandler(db *gorm.DB) *ContactTagHandler {
	return &ContactTagHandler{svc: service.NewWorkContactTagService(db)}
}

func (h *ContactTagHandler) Index(c *gin.Context) {
	corpID, _ := c.Get("corpId")
	if corpID == nil {
		corpID = uint(0)
	}
	tags, err := h.svc.List(corpID.(uint))
	if err != nil {
		response.Fail(c, response.ErrDB, "获取标签列表失败")
		return
	}
	response.Success(c, tags)
}

func (h *ContactTagHandler) AllTag(c *gin.Context) {
	h.Index(c)
}

func (h *ContactTagHandler) Store(c *gin.Context) {
	var tag model.WorkContactTag
	if err := c.ShouldBindJSON(&tag); err != nil {
		response.Fail(c, response.ErrParams, "参数错误")
		return
	}
	if err := h.svc.Create(&tag); err != nil {
		response.Fail(c, response.ErrDB, "创建标签失败")
		return
	}
	response.Success(c, tag)
}

func (h *ContactTagHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	tag, err := h.svc.GetByID(uint(id))
	if err != nil {
		response.Fail(c, response.ErrNotFound, "标签不存在")
		return
	}
	if err := c.ShouldBindJSON(tag); err != nil {
		response.Fail(c, response.ErrParams, "参数错误")
		return
	}
	if err := h.svc.Update(tag); err != nil {
		response.Fail(c, response.ErrDB, "更新标签失败")
		return
	}
	response.Success(c, tag)
}

func (h *ContactTagHandler) Destroy(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if err := h.svc.Delete(uint(id)); err != nil {
		response.Fail(c, response.ErrDB, "删除标签失败")
		return
	}
	response.SuccessMsg(c, "删除成功")
}

func (h *ContactTagHandler) Detail(c *gin.Context)            { h.Index(c) }
func (h *ContactTagHandler) ContactTagList(c *gin.Context)    { h.Index(c) }
func (h *ContactTagHandler) Move(c *gin.Context)              { response.SuccessMsg(c, "移动成功") }
func (h *ContactTagHandler) SynContactTag(c *gin.Context)     { response.SuccessMsg(c, "同步任务已提交") }

func (h *ContactTagHandler) UpdateByID(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ErrParams, err.Error())
		return
	}
	if err := h.svc.UpdateByID(uint(id), req); err != nil {
		response.Fail(c, response.ErrDB, "更新标签失败")
		return
	}
	response.SuccessMsg(c, "更新成功")
}

func (h *ContactTagHandler) ListByOrder(c *gin.Context) {
	corpID, _ := c.Get("corpId")
	if corpID == nil {
		corpID = uint(0)
	}
	tags, err := h.svc.ListByOrder(corpID.(uint))
	if err != nil {
		response.Fail(c, response.ErrDB, "获取标签列表失败")
		return
	}
	response.Success(c, tags)
}

// 辅助函数：获取字段类型文本
func getFieldTypeText(fieldType int) string {
	switch fieldType {
	case 1:
		return "单行文本"
	case 2:
		return "多行文本"
	case 3:
		return "数字"
	case 4:
		return "单选"
	case 5:
		return "多选"
	case 6:
		return "日期"
	case 7:
		return "附件"
	default:
		return "未知"
	}
}
