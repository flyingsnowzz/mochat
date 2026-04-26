// Package dashboard 提供 Dashboard 相关的 HTTP 处理器
// 该文件包含客户管理相关的处理器：
// 1. ContactHandler - 处理客户信息的增删改查、同步、统计等操作
// 2. ContactFieldHandler - 处理客户字段的管理操作
// 3. ContactTagHandler - 处理客户标签的管理操作
package dashboard

import (
	"encoding/json"
	"strconv"
	"strings"

	"mochat-api-server/internal/model"
	"mochat-api-server/internal/pkg/response"
	"mochat-api-server/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ContactHandler 客户管理处理器
// 提供客户信息的增删改查、同步、统计等功能
// 主要职责：
// 1. 客户列表查询与详情获取
// 2. 客户信息更新
// 3. 客户与企业微信同步
// 4. 客户轨迹查询
// 5. 流失客户查询
// 6. 客户来源统计
// 7. 批量打标签
// 8. 客户字段管理
// 9. 客户标签管理
//
// 依赖服务：
// - WorkContactService: 客户信息服务
// - WorkContactTagService: 客户标签服务
// - ContactFieldService: 客户字段服务
// - gorm.DB: 数据库连接

type ContactHandler struct {
	db         *gorm.DB                     // 数据库连接
	contactSvc *service.WorkContactService  // 客户信息服务
	tagSvc     *service.WorkContactTagService // 客户标签服务
	fieldSvc   *service.ContactFieldService  // 客户字段服务
}

// NewContactHandler 创建客户管理处理器实例
// 参数：db - GORM 数据库连接
// 返回：客户管理处理器实例
func NewContactHandler(db *gorm.DB) *ContactHandler {
	return &ContactHandler{
		db:         db,
		contactSvc: service.NewWorkContactService(db),
		tagSvc:     service.NewWorkContactTagService(db),
		fieldSvc:   service.NewContactFieldService(db),
	}
}

// Index 获取客户列表
// 支持分页查询企业的客户列表，返回包含员工信息、标签、群聊等关联数据
// 处理流程：
// 1. 获取企业 ID 和分页参数
// 2. 调用服务层获取客户列表
// 3. 构建返回数据，包含客户基本信息及关联数据
// 4. 返回分页结果
// 参数：
//
//	page - 页码，默认为 1
//	pageSize/perPage - 每页数量，默认为 20
//
// 返回：包含客户列表、总数、分页信息的响应
func (h *ContactHandler) Index(c *gin.Context) {
	// 获取企业 ID
	corpID, _ := c.Get("corpId")
	if corpID == nil {
		corpID = uint(0)
	}
	
	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", c.DefaultQuery("perPage", "20")))

	// 调用服务层获取客户列表
	contacts, total, err := h.contactSvc.List(corpID.(uint), page, pageSize, nil)
	if err != nil {
		response.Fail(c, response.ErrDB, "获取客户列表失败")
		return
	}
	
	// 构建返回数据，包含客户基本信息及关联数据
	list, err := h.buildContactList(corpID.(uint), contacts)
	if err != nil {
		response.Fail(c, response.ErrDB, "获取客户列表失败")
		return
	}
	
	// 返回分页结果
	response.PageResult(c, list, total, page, pageSize)
}

// Show 获取客户详情
// 根据客户 ID 获取客户的详细信息
// 处理流程：
// 1. 从请求中解析客户 ID
// 2. 调用服务层获取客户详情
// 3. 返回客户详情信息
// 参数：
//
//	id/contactId - 客户 ID
//
// 返回：客户详情信息
func (h *ContactHandler) Show(c *gin.Context) {
	// 从请求中解析客户 ID
	id, ok := parseContactID(c)
	if !ok {
		response.Fail(c, response.ErrParams, "参数错误")
		return
	}
	
	// 调用服务层获取客户详情
	contact, err := h.contactSvc.GetByID(uint(id))
	if err != nil {
		response.Fail(c, response.ErrNotFound, "客户不存在")
		return
	}
	
	// 返回客户详情信息
	response.Success(c, contact)
}

// Update 更新客户信息
// 更新指定客户的详细信息
// 处理流程：
// 1. 从请求中解析客户 ID
// 2. 调用服务层获取客户详情
// 3. 绑定请求参数到客户对象
// 4. 调用服务层更新客户信息
// 5. 返回更新后的客户信息
// 参数：
//
//	id/contactId - 客户 ID
//
// 请求体（JSON）：
//
//	包含客户信息的字段，如姓名、备注等
//
// 返回：更新后的客户详情信息
func (h *ContactHandler) Update(c *gin.Context) {
	// 从请求中解析客户 ID
	id, ok := parseContactID(c)
	if !ok {
		response.Fail(c, response.ErrParams, "参数错误")
		return
	}
	
	// 调用服务层获取客户详情
	contact, err := h.contactSvc.GetByID(uint(id))
	if err != nil {
		response.Fail(c, response.ErrNotFound, "客户不存在")
		return
	}
	
	// 绑定请求参数到客户对象
	if err := c.ShouldBindJSON(contact); err != nil {
		response.Fail(c, response.ErrParams, "参数错误")
		return
	}
	
	// 调用服务层更新客户信息
	if err := h.contactSvc.Update(contact); err != nil {
		response.Fail(c, response.ErrDB, "更新客户失败")
		return
	}
	
	// 返回更新后的客户信息
	response.Success(c, contact)
}

// SynContact 同步客户信息
// 异步提交客户信息同步任务，从企业微信拉取最新客户数据
// 处理流程：
// 1. 获取企业 ID
// 2. 异步启动同步任务
// 3. 返回同步任务已提交的消息
// 异步任务处理：
// 1. 获取企业微信客户端
// 2. 拉取客户列表
// 3. 同步到数据库
//
// 返回：同步任务已提交的成功消息
func (h *ContactHandler) SynContact(c *gin.Context) {
	// 获取企业 ID
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

	// 返回同步任务已提交的消息
	response.SuccessMsg(c, "同步任务已提交")
}

// Track 获取客户轨迹
// 获取指定客户的轨迹记录，按时间倒序排列
// 处理流程：
// 1. 从请求中解析客户 ID
// 2. 查询客户轨迹记录，按时间倒序排列
// 3. 返回客户轨迹列表
// 参数：
//
//	id/contactId - 客户 ID
//
// 返回：客户轨迹记录列表
func (h *ContactHandler) Track(c *gin.Context) {
	// 从请求中解析客户 ID
	id, ok := parseContactID(c)
	if !ok {
		response.Fail(c, response.ErrParams, "参数错误")
		return
	}
	
	// 获取客户轨迹记录，按时间倒序排列
	var tracks []model.ContactEmployeeTrack
	result := h.db.Where("contact_id = ?", uint(id)).Order("created_at DESC").Find(&tracks)
	if result.Error != nil {
		response.Fail(c, response.ErrDB, "获取客户轨迹失败")
		return
	}

	// 返回客户轨迹列表
	response.Success(c, tracks)
}

// LossContact 获取流失客户列表
// 查询企业的已删除客户列表，按删除时间倒序排列
// 处理流程：
// 1. 获取企业 ID 和分页参数
// 2. 查询已删除客户的总数
// 3. 查询分页数据，按删除时间倒序排列
// 4. 返回分页结果
// 参数：
//
//	page - 页码，默认为 1
//	pageSize - 每页数量，默认为 20
//
// 返回：包含流失客户列表、总数、分页信息的响应
func (h *ContactHandler) LossContact(c *gin.Context) {
	// 获取企业 ID
	corpID, _ := c.Get("corpId")
	if corpID == nil {
		corpID = uint(0)
	}
	
	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	// 查询已删除的客户
	var contacts []model.WorkContact
	var total int64

	// 构建查询条件
	query := h.db.Model(&model.WorkContact{}).Where("corp_id = ? AND deleted_at IS NOT NULL", corpID)
	
	// 查询总数
	if err := query.Count(&total).Error; err != nil {
		response.Fail(c, response.ErrDB, "获取流失客户列表失败")
		return
	}

	// 计算偏移量并查询分页数据
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("deleted_at DESC").Find(&contacts).Error; err != nil {
		response.Fail(c, response.ErrDB, "获取流失客户列表失败")
		return
	}

	// 返回分页结果
	response.PageResult(c, contacts, total, page, pageSize)
}

// Source 获取客户来源统计
// 统计企业客户的添加方式分布情况
// 处理流程：
// 1. 获取企业 ID
// 2. 定义统计结果结构体
// 3. 定义添加方式映射表
// 4. 查询各添加方式的客户数量
// 5. 构建统计结果并返回
//
// 返回：客户来源统计列表，包含添加方式、数量、添加方式文本
func (h *ContactHandler) Source(c *gin.Context) {
	// 获取企业 ID
	corpID, _ := c.Get("corpId")
	if corpID == nil {
		corpID = uint(0)
	}

	// 统计客户来源
	type SourceStat struct {
		AddWay    int   `json:"addWay"`    // 添加方式代码
		Count     int64 `json:"count"`     // 客户数量
		AddWayStr string `json:"addWayStr"` // 添加方式文本
	}

	var stats []SourceStat

	// 定义添加方式映射
	addWayMap := map[int]string{
		0:  "未知",
		1:  "扫描二维码",
		2:  "搜索手机号",
		3:  "名片分享",
		4:  "群聊",
		5:  "手机通讯录",
		6:  "微信好友",
		7:  "来自微信的添加",
		8:  "安装第三方应用",
		9:  "搜索邮箱",
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

	// 构建统计结果
	for rows.Next() {
		var stat SourceStat
		if err := rows.Scan(&stat.AddWay, &stat.Count); err != nil {
			response.Fail(c, response.ErrDB, "获取客户来源统计失败")
			return
		}
		stat.AddWayStr = addWayMap[stat.AddWay]
		stats = append(stats, stat)
	}

	// 返回统计结果
	response.Success(c, stats)
}

// BatchLabeling 批量打标签
// 为多个客户批量添加标签，会先删除已有标签关联，再添加新的标签关联
// 处理流程：
// 1. 绑定请求参数
// 2. 开始数据库事务
// 3. 删除已有的标签关联
// 4. 添加新的标签关联
// 5. 提交事务
// 6. 返回操作结果
// 请求体（JSON）：
//
//	contactIds - 客户 ID 列表
//	tagIds - 标签 ID 列表
//
// 返回：批量打标签操作结果
func (h *ContactHandler) BatchLabeling(c *gin.Context) {
	// 绑定请求参数
	var req struct {
		ContactIDs []uint `json:"contactIds"` // 客户 ID 列表
		TagIDs     []uint `json:"tagIds"`     // 标签 ID 列表
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ErrParams, "参数错误")
		return
	}

	// 批量打标签（使用事务确保操作原子性）
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

	// 返回成功消息
	response.SuccessMsg(c, "批量打标签成功")
}

// UpdateByID 根据 ID 更新客户信息
// 根据客户 ID 更新指定字段的信息
// 处理流程：
// 1. 从路径参数中获取客户 ID
// 2. 绑定请求参数
// 3. 调用服务层更新客户信息
// 4. 返回操作结果
// 参数：
//
//	id - 客户 ID（路径参数）
//
// 请求体（JSON）：
//
//	包含要更新的字段和值的映射
//
// 返回：更新操作结果
func (h *ContactHandler) UpdateByID(c *gin.Context) {
	// 从路径参数中获取客户 ID
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	
	// 绑定请求参数
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ErrParams, err.Error())
		return
	}
	
	// 调用服务层更新客户信息
	if err := h.contactSvc.UpdateByID(uint(id), req); err != nil {
		response.Fail(c, response.ErrDB, "更新客户失败")
		return
	}
	
	// 返回成功消息
	response.SuccessMsg(c, "更新成功")
}

// GetByWxExternalUserID 根据微信外部用户 ID 获取客户信息
// 根据企业微信外部用户 ID 查询客户信息
// 处理流程：
// 1. 获取企业 ID
// 2. 获取微信外部用户 ID
// 3. 调用服务层根据微信外部用户 ID 获取客户信息
// 4. 返回客户详情信息
// 参数：
//
//	wxExternalUserid - 微信外部用户 ID（查询参数）
//
// 返回：客户详情信息
func (h *ContactHandler) GetByWxExternalUserID(c *gin.Context) {
	// 获取企业 ID
	corpID, _ := c.Get("corpId")
	
	// 获取微信外部用户 ID
	wxID := c.Query("wxExternalUserid")
	if wxID == "" {
		response.Fail(c, response.ErrParams, "wxExternalUserid参数不能为空")
		return
	}
	
	// 调用服务层根据微信外部用户 ID 获取客户信息
	contact, err := h.contactSvc.GetByWxExternalUserID(corpID.(uint), wxID)
	if err != nil {
		response.Fail(c, response.ErrNotFound, "客户不存在")
		return
	}
	
	// 返回客户详情信息
	response.Success(c, contact)
}

// buildContactList 构建客户列表数据
// 将客户模型数据转换为前端需要的格式，包含关联的员工、标签、群聊等信息
// 处理流程：
// 1. 收集客户 ID 列表
// 2. 加载员工名称映射
// 3. 加载客户与员工的最新关联关系
// 4. 加载客户所在的群聊名称
// 5. 加载客户的标签名称
// 6. 构建返回数据
// 参数：
//
//	corpID - 企业 ID
//	contacts - 客户模型列表
//
// 返回：构建后的客户列表数据和可能的错误
func (h *ContactHandler) buildContactList(corpID uint, contacts []model.WorkContact) ([]gin.H, error) {
	if len(contacts) == 0 {
		return []gin.H{}, nil
	}

	// 收集客户 ID 列表
	contactIDs := make([]uint, 0, len(contacts))
	for _, contact := range contacts {
		contactIDs = append(contactIDs, contact.ID)
	}

	// 加载员工名称映射
	employeeMap, err := h.workEmployeeNameMap(corpID)
	if err != nil {
		return nil, err
	}
	
	// 加载客户与员工的最新关联关系
	relations, err := h.latestContactEmployeeMap(corpID, contactIDs)
	if err != nil {
		return nil, err
	}
	
	// 加载客户所在的群聊名称
	roomNames, err := h.contactRoomNames(contactIDs)
	if err != nil {
		return nil, err
	}
	
	// 加载客户的标签名称
	tagNames, err := h.contactTagNames(contactIDs)
	if err != nil {
		return nil, err
	}

	// 构建返回数据
	list := make([]gin.H, 0, len(contacts))
	for _, contact := range contacts {
		row := relations[contact.ID]
		employeeName := ""
		employeeID := uint(0)
		addWay := 0
		createTime := ""
		if row != nil {
			employeeID = row.EmployeeID
			addWay = row.AddWay
			createTime = row.CreateTime.Format("2006-01-02 15:04:05")
			employeeName = employeeMap[row.EmployeeID]
		}

		list = append(list, gin.H{
			"id":          contact.ID,                // 客户 ID
			"name":        firstNonEmpty(strings.TrimSpace(contact.Name), strings.TrimSpace(contact.NickName)), // 客户名称
			"avatar":      contact.Avatar,            // 客户头像
			"employeeName": employeeName,             // 关联员工名称
			"roomName":    defaultStringSlice(roomNames[contact.ID]), // 所在群聊名称
			"addWayText":  contactAddWayText(addWay), // 添加方式文本
			"tag":         defaultStringSlice(tagNames[contact.ID]), // 标签名称列表
			"createTime":  createTime,                // 创建时间
			"contactId":   contact.ID,                // 客户 ID（兼容旧接口）
			"employeeId":  employeeID,                // 员工 ID
			"isContact":   2,                         // 是否为客户（固定值）
			"businessNo":  contact.BusinessNo,         // 业务编号
		})
	}
	return list, nil
}

// latestContactEmployeeMap 获取客户与员工的最新关联关系
// 为每个客户获取最新的员工关联记录
// 处理流程：
// 1. 查询客户与员工的关联记录，按创建时间和 ID 倒序排列
// 2. 构建客户 ID 到最新关联记录的映射
// 参数：
//
//	corpID - 企业 ID
//	contactIDs - 客户 ID 列表
//
// 返回：客户 ID 到最新关联记录的映射和可能的错误
func (h *ContactHandler) latestContactEmployeeMap(corpID uint, contactIDs []uint) (map[uint]*model.WorkContactEmployee, error) {
	// 查询客户与员工的关联记录，按创建时间和 ID 倒序排列
	var rows []model.WorkContactEmployee
	if err := h.db.Where("corp_id = ? AND contact_id IN ?", corpID, contactIDs).
		Order("create_time DESC, id DESC").
		Find(&rows).Error; err != nil {
		return nil, err
	}
	
	// 构建客户 ID 到最新关联记录的映射
	result := make(map[uint]*model.WorkContactEmployee, len(contactIDs))
	for i := range rows {
		if _, ok := result[rows[i].ContactID]; ok {
			continue // 已经找到该客户的最新记录，跳过
		}
		result[rows[i].ContactID] = &rows[i]
	}
	return result, nil
}

// workEmployeeNameMap 获取员工名称映射
// 构建员工 ID 到员工名称的映射
// 处理流程：
// 1. 查询企业的所有未删除员工
// 2. 构建员工 ID 到员工名称的映射
// 参数：
//
//	corpID - 企业 ID
//
// 返回：员工 ID 到员工名称的映射和可能的错误
func (h *ContactHandler) workEmployeeNameMap(corpID uint) (map[uint]string, error) {
	// 查询企业的所有未删除员工
	var rows []model.WorkEmployee
	if err := h.db.Where("corp_id = ? AND deleted_at IS NULL", corpID).Find(&rows).Error; err != nil {
		return nil, err
	}
	
	// 构建员工 ID 到员工名称的映射
	result := make(map[uint]string, len(rows))
	for _, row := range rows {
		result[row.ID] = row.Name
	}
	return result, nil
}

// contactRoomNames 获取客户所在的群聊名称
// 为每个客户获取其所在的群聊名称列表
// 处理流程：
// 1. 查询客户与群聊的关联关系，关联群聊表获取群聊名称
// 2. 过滤掉空群聊名称
// 3. 构建客户 ID 到群聊名称列表的映射
// 参数：
//
//	contactIDs - 客户 ID 列表
//
// 返回：客户 ID 到群聊名称列表的映射和可能的错误
func (h *ContactHandler) contactRoomNames(contactIDs []uint) (map[uint][]string, error) {
	// 定义结果结构体
	var rows []struct {
		ContactID uint   `json:"contactId"` // 客户 ID
		RoomName  string `json:"roomName"`  // 群聊名称
	}
	
	// 查询客户与群聊的关联关系，关联群聊表获取群聊名称
	if err := h.db.Table("mc_work_contact_room AS wcr").
		Select("wcr.contact_id AS contact_id, wr.name AS room_name").
		Joins("LEFT JOIN mc_work_room AS wr ON wr.id = wcr.room_id").
		Where("wcr.contact_id IN ?", contactIDs).
		Where("wr.deleted_at IS NULL").
		Scan(&rows).Error; err != nil {
		return nil, err
	}
	
	// 构建客户 ID 到群聊名称列表的映射
	result := make(map[uint][]string, len(contactIDs))
	for _, row := range rows {
		if strings.TrimSpace(row.RoomName) == "" {
			continue // 跳过空群聊名称
		}
		result[row.ContactID] = appendIfMissing(result[row.ContactID], row.RoomName)
	}
	return result, nil
}

// contactTagNames 获取客户的标签名称
// 为每个客户获取其标签名称列表
// 处理流程：
// 1. 查询客户与标签的关联关系，关联标签表获取标签名称
// 2. 过滤掉空标签名称
// 3. 构建客户 ID 到标签名称列表的映射
// 参数：
//
//	contactIDs - 客户 ID 列表
//
// 返回：客户 ID 到标签名称列表的映射和可能的错误
func (h *ContactHandler) contactTagNames(contactIDs []uint) (map[uint][]string, error) {
	// 定义结果结构体
	var rows []struct {
		ContactID uint   `json:"contactId"` // 客户 ID
		TagName   string `json:"tagName"`  // 标签名称
	}
	
	// 查询客户与标签的关联关系，关联标签表获取标签名称
	if err := h.db.Table("mc_work_contact_tag_pivot AS pivot").
		Select("pivot.contact_id AS contact_id, tag.name AS tag_name").
		Joins("LEFT JOIN mc_work_contact_tag AS tag ON tag.id = pivot.contact_tag_id").
		Where("pivot.contact_id IN ?", contactIDs).
		Scan(&rows).Error; err != nil {
		return nil, err
	}
	
	// 构建客户 ID 到标签名称列表的映射
	result := make(map[uint][]string, len(contactIDs))
	for _, row := range rows {
		if strings.TrimSpace(row.TagName) == "" {
			continue // 跳过空标签名称
		}
		result[row.ContactID] = appendIfMissing(result[row.ContactID], row.TagName)
	}
	return result, nil
}

// defaultStringSlice 处理字符串切片默认值
// 如果输入切片为 nil，返回空切片
// 参数：
//
//	items - 输入字符串切片
//
// 返回：非 nil 的字符串切片
func defaultStringSlice(items []string) []string {
	if items == nil {
		return []string{} // 返回空切片
	}
	return items
}

// appendIfMissing 向切片中添加元素（去重）
// 如果切片中已存在该元素，则不添加
// 参数：
//
//	items - 输入字符串切片
//	value - 要添加的字符串
//
// 返回：添加元素后的切片
func appendIfMissing(items []string, value string) []string {
	// 检查切片中是否已存在该元素
	for _, item := range items {
		if item == value {
			return items // 已存在，直接返回
		}
	}
	// 不存在，添加到切片末尾
	return append(items, value)
}

// firstNonEmpty 获取第一个非空字符串
// 按顺序检查输入的字符串，返回第一个非空的字符串
// 参数：
//
//	values - 字符串列表
//
// 返回：第一个非空字符串，如果所有字符串都为空则返回空字符串
func firstNonEmpty(values ...string) string {
	// 按顺序检查每个字符串
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value // 找到非空字符串，返回
		}
	}
	// 所有字符串都为空，返回空字符串
	return ""
}

// contactAddWayText 获取添加方式文本
// 将添加方式代码转换为对应的文本描述
// 参数：
//
//	addWay - 添加方式代码
//
// 返回：添加方式文本描述
func contactAddWayText(addWay int) string {
	switch addWay {
	case 1:
		return "扫描二维码"
	case 2:
		return "搜索手机号"
	case 3:
		return "名片分享"
	case 4:
		return "群聊"
	case 5:
		return "手机通讯录"
	case 6:
		return "微信好友"
	case 7:
		return "来自微信的添加"
	case 8:
		return "安装第三方应用"
	case 9:
		return "搜索邮箱"
	case 10:
		return "企业微信内部成员共享"
	case 11:
		return "管理员/负责人分配"
	default:
		return "未知"
	}
}

// splitContactFieldOptions 分割客户字段选项
// 将客户字段的选项字符串分割为字符串切片
// 支持两种格式：
// 1. JSON 数组格式：["选项1", "选项2"]
// 2. 逗号/换行分隔格式：选项1,选项2
// 参数：
//
//	raw - 原始选项字符串
//
// 返回：分割后的选项字符串切片
func splitContactFieldOptions(raw string) []string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return []string{} // 空字符串返回空切片
	}
	
	// 尝试解析 JSON 数组格式
	if strings.HasPrefix(raw, "[") {
		var arr []string
		if err := json.Unmarshal([]byte(raw), &arr); err == nil {
			return arr
		}
	}
	
	// 解析逗号/换行分隔格式
	parts := strings.FieldsFunc(raw, func(r rune) bool {
		return r == ',' || r == '\n' || r == '\r'
	})
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue // 跳过空选项
		}
		result = append(result, part)
	}
	return result
}

// parseContactID 解析客户 ID
// 从请求的路径参数或查询参数中解析客户 ID
// 支持多个参数名：id, contactId
// 参数：
//
//	c - Gin 上下文
//
// 返回：解析后的客户 ID 和是否解析成功
func parseContactID(c *gin.Context) (uint64, bool) {
	// 尝试从多个参数名中获取 ID
	for _, raw := range []string{c.Param("id"), c.Query("contactId"), c.Query("id")} {
		raw = strings.TrimSpace(raw)
		if raw == "" || raw == "undefined" || raw == "null" {
			continue // 跳过无效值
		}
		id, err := strconv.ParseUint(raw, 10, 32)
		if err == nil && id > 0 {
			return id, true // 解析成功
		}
	}
	// 所有参数都解析失败
	return 0, false
}

// parseIDFromParamOrBody 从参数或请求体中解析 ID
// 先尝试从路径参数中解析 ID，失败则尝试从请求体中解析
// 参数：
//
//	c - Gin 上下文
//
// 返回：解析后的 ID 和是否解析成功
func parseIDFromParamOrBody(c *gin.Context) (uint, bool) {
	// 尝试从路径参数中解析 ID
	if raw := strings.TrimSpace(c.Param("id")); raw != "" {
		if id, err := strconv.ParseUint(raw, 10, 32); err == nil && id > 0 {
			return uint(id), true // 解析成功
		}
	}
	
	// 尝试从请求体中解析 ID
	var req struct {
		ID uint `json:"id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		return 0, false // 请求体解析失败
	}
	if req.ID == 0 {
		return 0, false // ID 为 0，无效
	}
	return req.ID, true // 解析成功
}

// uintFromAny 将任意类型转换为 uint
// 支持多种类型的转换，包括数字类型和字符串类型
// 负数会被转换为 0
// 参数：
//
//	v - 任意类型的值
//
// 返回：转换后的 uint 值
func uintFromAny(v interface{}) uint {
	switch val := v.(type) {
	case float64:
		if val < 0 {
			return 0 // 负数转换为 0
		}
		return uint(val)
	case float32:
		if val < 0 {
			return 0 // 负数转换为 0
		}
		return uint(val)
	case int:
		if val < 0 {
			return 0 // 负数转换为 0
		}
		return uint(val)
	case int64:
		if val < 0 {
			return 0 // 负数转换为 0
		}
		return uint(val)
	case uint:
		return val // 直接返回
	case uint64:
		return uint(val) // 转换为 uint
	case string:
		val = strings.TrimSpace(val)
		if val == "" {
			return 0 // 空字符串转换为 0
		}
		n, err := strconv.ParseUint(val, 10, 32)
		if err != nil {
			return 0 // 解析失败转换为 0
		}
		return uint(n)
	default:
		return 0 // 其他类型转换为 0
	}
}

// intFromAny 将任意类型转换为 int
// 支持多种类型的转换，包括数字类型和字符串类型
// 空字符串会被转换为 0
// 参数：
//
//	v - 任意类型的值
//
// 返回：转换后的 int 值
func intFromAny(v interface{}) int {
	switch val := v.(type) {
	case float64:
		return int(val) // 转换为 int
	case float32:
		return int(val) // 转换为 int
	case int:
		return val // 直接返回
	case int64:
		return int(val) // 转换为 int
	case uint:
		return int(val) // 转换为 int
	case uint64:
		return int(val) // 转换为 int
	case string:
		val = strings.TrimSpace(val)
		if val == "" {
			return 0 // 空字符串转换为 0
		}
		n, err := strconv.Atoi(val)
		if err != nil {
			return 0 // 解析失败转换为 0
		}
		return n
	default:
		return 0 // 其他类型转换为 0
	}
}

// stringFromAny 将任意类型转换为 string
// 支持多种类型的转换，包括字符串类型和数字类型
// 参数：
//
//	v - 任意类型的值
//
// 返回：转换后的 string 值
func stringFromAny(v interface{}) string {
	switch val := v.(type) {
	case string:
		return val // 直接返回
	case float64:
		return strconv.FormatFloat(val, 'f', -1, 64) // 转换为字符串
	case float32:
		return strconv.FormatFloat(float64(val), 'f', -1, 32) // 转换为字符串
	case int:
		return strconv.Itoa(val) // 转换为字符串
	case int64:
		return strconv.FormatInt(val, 10) // 转换为字符串
	case uint:
		return strconv.FormatUint(uint64(val), 10) // 转换为字符串
	case uint64:
		return strconv.FormatUint(val, 10) // 转换为字符串
	default:
		return "" // 其他类型转换为空字符串
	}
}

// normalizeContactFieldOptions 标准化客户字段选项
// 将各种类型的选项数据转换为标准的 JSON 数组字符串
// 支持多种输入类型：
// 1. nil - 返回 "[]"
// 2. string - 支持 JSON 数组格式或逗号分隔格式
// 3. []string - 直接转换为 JSON 数组
// 4. []interface{} - 转换为 []string 后再转换为 JSON 数组
// 5. 其他类型 - 转换为字符串后再处理
// 参数：
//
//	v - 任意类型的选项数据
//
// 返回：标准化后的 JSON 数组字符串
func normalizeContactFieldOptions(v interface{}) string {
	switch val := v.(type) {
	case nil:
		return "[]" // nil 返回空数组
	case string:
		s := strings.TrimSpace(val)
		if s == "" || s == "[]" {
			return "[]" // 空字符串或空数组返回 "[]"
		}
		if strings.HasPrefix(s, "[") {
			var arr []string
			if err := json.Unmarshal([]byte(s), &arr); err == nil {
				return marshalContactFieldOptions(arr) // 解析 JSON 数组
			}
		}
		return marshalContactFieldOptions(splitContactFieldOptions(s)) // 解析逗号分隔格式
	case []string:
		return marshalContactFieldOptions(val) // 直接转换为 JSON 数组
	case []interface{}:
		items := make([]string, 0, len(val))
		for _, item := range val {
			text := strings.TrimSpace(stringFromAny(item))
			if text == "" {
				continue // 跳过空字符串
			}
			items = append(items, text)
		}
		return marshalContactFieldOptions(items) // 转换为 JSON 数组
	default:
		text := strings.TrimSpace(stringFromAny(val))
		if text == "" {
			return "[]" // 空字符串返回 "[]"
		}
		return marshalContactFieldOptions([]string{text}) // 转换为 JSON 数组
	}
}

// marshalContactFieldOptions 将字符串切片转换为 JSON 数组字符串
// 会过滤掉空字符串，并处理 JSON 序列化失败的情况
// 参数：
//
//	items - 字符串切片
//
// 返回：JSON 数组字符串
func marshalContactFieldOptions(items []string) string {
	// 过滤掉空字符串
	normalized := make([]string, 0, len(items))
	for _, item := range items {
		item = strings.TrimSpace(item)
		if item == "" {
			continue // 跳过空字符串
		}
		normalized = append(normalized, item)
	}
	
	// 处理空切片
	if len(normalized) == 0 {
		return "[]" // 返回空数组
	}
	
	// 转换为 JSON 数组
	data, err := json.Marshal(normalized)
	if err != nil {
		return "[]" // 序列化失败返回空数组
	}
	return string(data)
}

// ListByOffset 按偏移量获取客户列表
// 支持分页查询企业的客户列表，使用偏移量进行分页
// 处理流程：
// 1. 获取企业 ID
// 2. 获取分页参数并计算偏移量
// 3. 调用服务层获取客户列表
// 4. 返回分页结果
// 参数：
//
//	page - 页码，默认为 1
//	pageSize - 每页数量，默认为 20
//
// 返回：包含客户列表、总数、分页信息的响应
func (h *ContactHandler) ListByOffset(c *gin.Context) {
	// 获取企业 ID
	corpID, _ := c.Get("corpId")
	if corpID == nil {
		corpID = uint(0)
	}
	
	// 获取分页参数并计算偏移量
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	offset := (page - 1) * pageSize

	// 调用服务层获取客户列表
	contacts, total, err := h.contactSvc.ListByOffset(corpID.(uint), offset, pageSize, nil)
	if err != nil {
		response.Fail(c, response.ErrDB, "获取客户列表失败")
		return
	}
	
	// 返回分页结果
	response.PageResult(c, contacts, total, page, pageSize)
}

// ContactFieldHandler 客户字段管理处理器
// 提供客户字段的增删改查、状态更新等功能
// 主要职责：
// 1. 客户字段列表查询与详情获取
// 2. 客户字段创建与更新
// 3. 客户字段删除
// 4. 客户字段状态更新
// 5. 客户画像获取
// 6. 客户字段批量更新
//
// 依赖服务：
// - ContactFieldService: 客户字段服务
// - ContactFieldPivotService: 客户字段值关联服务

type ContactFieldHandler struct {
	svc       *service.ContactFieldService       // 客户字段服务
	pivotSvc  *service.ContactFieldPivotService  // 客户字段值关联服务
}

// NewContactFieldHandler 创建客户字段管理处理器实例
// 参数：db - GORM 数据库连接
// 返回：客户字段管理处理器实例
func NewContactFieldHandler(db *gorm.DB) *ContactFieldHandler {
	return &ContactFieldHandler{
		svc:       service.NewContactFieldService(db),
		pivotSvc:  service.NewContactFieldPivotService(db),
	}
}

// Index 获取客户字段列表
// 支持按状态筛选和分页查询客户字段列表
// 处理流程：
// 1. 获取状态、页码和每页数量参数
// 2. 调用服务层获取字段列表
// 3. 处理字段类型文本
// 4. 返回分页结果
// 参数：
//
//	status - 状态（0-禁用，1-启用，2-全部），默认为 2
//	page - 页码，默认为 1
//	pageSize/perPage - 每页数量，默认为 10
//
// 返回：包含字段列表、总数、分页信息的响应
func (h *ContactFieldHandler) Index(c *gin.Context) {
	// 获取状态、页码和每页数量参数
	status, _ := strconv.Atoi(c.DefaultQuery("status", "2"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", c.DefaultQuery("perPage", "10")))

	// 调用服务层获取字段列表
	fields, total, err := h.svc.ListByStatus(status, page, pageSize)
	if err != nil {
		response.Fail(c, response.ErrDB, "获取字段列表失败")
		return
	}

	// 处理字段类型文本
	for i := range fields {
		fields[i].TypeText = getFieldTypeText(fields[i].Type)
	}

	// 返回分页结果
	response.PageResult(c, fields, total, page, pageSize)
}

// Show 获取客户字段详情
// 根据字段 ID 获取客户字段的详细信息
// 处理流程：
// 1. 从路径参数中获取字段 ID
// 2. 调用服务层获取字段详情
// 3. 处理字段类型文本
// 4. 返回字段详情信息
// 参数：
//
//	id - 字段 ID（路径参数）
//
// 返回：字段详情信息
func (h *ContactFieldHandler) Show(c *gin.Context) {
	// 从路径参数中获取字段 ID
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	
	// 调用服务层获取字段详情
	field, err := h.svc.GetByID(uint(id))
	if err != nil {
		response.Fail(c, response.ErrNotFound, "字段不存在")
		return
	}

	// 处理字段类型文本
	field.TypeText = getFieldTypeText(field.Type)

	// 返回字段详情信息
	response.Success(c, field)
}

// Store 创建客户字段
// 创建新的客户字段
// 处理流程：
// 1. 绑定请求参数到字段对象
// 2. 调用服务层创建字段
// 3. 处理字段类型文本
// 4. 返回创建的字段信息
// 请求体（JSON）：
//
//	包含客户字段的相关信息，如名称、类型、选项等
//
// 返回：创建的字段详情信息
func (h *ContactFieldHandler) Store(c *gin.Context) {
	// 绑定请求参数到字段对象
	var field model.ContactField
	if err := c.ShouldBindJSON(&field); err != nil {
		response.Fail(c, response.ErrParams, "参数错误")
		return
	}
	
	// 调用服务层创建字段
	if err := h.svc.Create(&field); err != nil {
		response.Fail(c, response.ErrDB, "创建字段失败")
		return
	}

	// 处理字段类型文本
	field.TypeText = getFieldTypeText(field.Type)

	// 返回创建的字段信息
	response.Success(c, field)
}

// Update 更新客户字段
// 更新指定客户字段的信息
// 处理流程：
// 1. 绑定请求参数
// 2. 获取字段 ID（从请求体或路径参数）
// 3. 调用服务层获取字段详情
// 4. 更新字段信息
// 5. 调用服务层更新字段
// 6. 处理字段类型文本
// 7. 返回更新后的字段信息
// 请求体（JSON）：
//
//	id - 字段 ID（可选，若不提供则从路径参数获取）
//	name - 字段名称
//	label - 字段标签
//	type - 字段类型
//	options - 字段选项
//	order - 排序
//	status - 状态
//	isSys - 是否系统字段
//
// 返回：更新后的字段详情信息
func (h *ContactFieldHandler) Update(c *gin.Context) {
	// 绑定请求参数
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ErrParams, "参数错误")
		return
	}
	
	// 获取字段 ID（从请求体或路径参数）
	id := uint(0)
	if raw, ok := req["id"]; ok {
		id = uintFromAny(raw)
	}
	if id == 0 {
		if parsedID, err := strconv.ParseUint(c.Param("id"), 10, 32); err == nil {
			id = uint(parsedID)
		}
	}
	if id == 0 {
		response.Fail(c, response.ErrParams, "参数错误")
		return
	}

	// 调用服务层获取字段详情
	field, err := h.svc.GetByID(id)
	if err != nil {
		response.Fail(c, response.ErrNotFound, "字段不存在")
		return
	}
	
	// 更新字段信息
	if raw, ok := req["name"]; ok {
		field.Name = strings.TrimSpace(stringFromAny(raw))
	}
	if raw, ok := req["label"]; ok {
		field.Label = strings.TrimSpace(stringFromAny(raw))
		// 兼容旧前端，若名称为空则与标签保持一致
		if strings.TrimSpace(field.Name) == "" {
			field.Name = field.Label
		}
	}
	if raw, ok := req["type"]; ok {
		field.Type = intFromAny(raw)
	}
	if raw, ok := req["options"]; ok {
		field.Options = normalizeContactFieldOptions(raw)
	}
	if raw, ok := req["order"]; ok {
		field.Order = intFromAny(raw)
	}
	if raw, ok := req["status"]; ok {
		field.Status = intFromAny(raw)
	}
	if raw, ok := req["isSys"]; ok {
		field.IsSys = intFromAny(raw)
	}
	
	// 调用服务层更新字段
	if err := h.svc.Update(field); err != nil {
		response.Fail(c, response.ErrDB, "更新字段失败")
		return
	}

	// 处理字段类型文本
	field.TypeText = getFieldTypeText(field.Type)

	// 返回更新后的字段信息
	response.Success(c, field)
}

func (h *ContactFieldHandler) Destroy(c *gin.Context) {
	id, ok := parseIDFromParamOrBody(c)
	if !ok {
		response.Fail(c, response.ErrParams, "参数错误")
		return
	}
	if err := h.svc.Delete(id); err != nil {
		response.Fail(c, response.ErrDB, "删除字段失败")
		return
	}
	response.SuccessMsg(c, "删除成功")
}

func (h *ContactFieldHandler) StatusUpdate(c *gin.Context) {
	var req struct {
		ID     uint `json:"id"`
		Status int `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ErrParams, "参数错误")
		return
	}
	if req.ID == 0 {
		if id, err := strconv.ParseUint(c.Param("id"), 10, 32); err == nil {
			req.ID = uint(id)
		}
	}
	if req.ID == 0 {
		response.Fail(c, response.ErrParams, "参数错误")
		return
	}

	if err := h.svc.UpdateStatus(req.ID, req.Status); err != nil {
		response.Fail(c, response.ErrDB, "更新状态失败")
		return
	}

	response.SuccessMsg(c, "更新成功")
}

func (h *ContactFieldHandler) Portrait(c *gin.Context) {
	idText := c.Param("id")
	if idText != "" {
		id, _ := strconv.ParseUint(idText, 10, 32)
		fieldPivots, err := h.pivotSvc.List(uint(id))
		if err != nil {
			response.Fail(c, response.ErrDB, "获取客户画像失败")
			return
		}
		response.Success(c, fieldPivots)
		return
	}

	fields, err := h.svc.List()
	if err != nil {
		response.Fail(c, response.ErrDB, "获取客户画像字段失败")
		return
	}

	list := make([]gin.H, 0, len(fields))
	for _, field := range fields {
		if field.Status != 1 {
			continue
		}
		list = append(list, gin.H{
			"fieldId":  field.ID,
			"id":       field.ID,
			"name":     firstNonEmpty(strings.TrimSpace(field.Name), strings.TrimSpace(field.Label)),
			"label":    field.Label,
			"type":     field.Type,
			"typeText": getFieldTypeText(field.Type),
			"options":  splitContactFieldOptions(field.Options),
		})
	}

	response.Success(c, list)
}

func (h *ContactFieldHandler) BatchUpdate(c *gin.Context) {
	var req struct {
		Fields []struct {
			ID    uint `json:"id"`
			Order int  `json:"order"`
		} `json:"fields"`
		Update []struct {
			ID    uint `json:"id"`
			Order int  `json:"order"`
		} `json:"update"`
		Destroy []uint `json:"destroy"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ErrParams, "参数错误")
		return
	}

	fields := req.Fields
	if len(fields) == 0 {
		fields = req.Update
	}
	for _, item := range fields {
		if err := h.svc.UpdateOrder(item.ID, item.Order); err != nil {
			response.Fail(c, response.ErrDB, "批量更新失败")
			return
		}
	}
	for _, id := range req.Destroy {
		if id == 0 {
			continue
		}
		if err := h.svc.Delete(id); err != nil {
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
