package system

import (
	"strconv"
	"strings"

	"mochat-api-server/internal/model"
	"mochat-api-server/internal/pkg/response"
	"mochat-api-server/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RoleHandler struct {
	svc *service.RbacRoleService
	db  *gorm.DB
}

type permissionMenu struct {
	ID             uint             `json:"id"`
	MenuID         uint             `json:"menuId"`
	ParentID       uint             `json:"parentId"`
	Name           string           `json:"name"`
	Level          int              `json:"level"`
	Icon           string           `json:"icon"`
	LinkType       int              `json:"linkType"`
	LinkURL        string           `json:"linkUrl"`
	DataPermission int              `json:"dataPermission"`
	IsPageMenu     int              `json:"isPageMenu"`
	Children       []permissionMenu `json:"children"`
}

func NewRoleHandler(db *gorm.DB) *RoleHandler {
	return &RoleHandler{svc: service.NewRbacRoleService(db), db: db}
}

func (h *RoleHandler) Index(c *gin.Context) {
	// 获取租户ID，默认为0
	tenantID := uint(0)
	if id, exists := c.Get("tenantId"); exists {
		if tid, ok := id.(uint); ok {
			tenantID = tid
		}
	}

	// 获取分页参数，支持perPage和pageSize两种参数名
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("perPage", "10"))
	if pageSize <= 0 {
		pageSize, _ = strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	}

	roles, total, err := h.svc.List(tenantID, page, pageSize)
	if err != nil {
		response.Fail(c, response.ErrDB, "获取角色列表失败")
		return
	}
	response.PageResult(c, roles, total, page, pageSize)
}

func (h *RoleHandler) Select(c *gin.Context) {
	// 获取租户ID，默认为0
	tenantID := uint(0)
	if id, exists := c.Get("tenantId"); exists {
		if tid, ok := id.(uint); ok {
			tenantID = tid
		}
	}
	roles, err := h.svc.Select(tenantID)
	if err != nil {
		response.Fail(c, response.ErrDB, "获取角色选择列表失败")
		return
	}
	response.Success(c, roles)
}

func (h *RoleHandler) Show(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	role, err := h.svc.GetByID(uint(id))
	if err != nil {
		response.Fail(c, response.ErrNotFound, "角色不存在")
		return
	}
	response.Success(c, role)
}

func (h *RoleHandler) Store(c *gin.Context) {
	var role model.RbacRole
	if err := c.ShouldBindJSON(&role); err != nil {
		response.Fail(c, response.ErrParams, "参数错误")
		return
	}
	if err := h.svc.Create(&role); err != nil {
		response.Fail(c, response.ErrDB, "创建角色失败")
		return
	}
	response.Success(c, role)
}

func (h *RoleHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	role, err := h.svc.GetByID(uint(id))
	if err != nil {
		response.Fail(c, response.ErrNotFound, "角色不存在")
		return
	}
	if err := c.ShouldBindJSON(role); err != nil {
		response.Fail(c, response.ErrParams, "参数错误")
		return
	}
	if err := h.svc.Update(role); err != nil {
		response.Fail(c, response.ErrDB, "更新角色失败")
		return
	}
	response.Success(c, role)
}

func (h *RoleHandler) Destroy(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if err := h.svc.Delete(uint(id)); err != nil {
		response.Fail(c, response.ErrDB, "删除角色失败")
		return
	}
	response.SuccessMsg(c, "删除成功")
}

func (h *RoleHandler) StatusUpdate(c *gin.Context) {
	var req struct {
		ID     uint `json:"id" binding:"required"`
		Status int  `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ErrParams, "参数错误")
		return
	}
	role, err := h.svc.GetByID(req.ID)
	if err != nil {
		response.Fail(c, response.ErrNotFound, "角色不存在")
		return
	}
	role.Status = req.Status
	if err := h.svc.Update(role); err != nil {
		response.Fail(c, response.ErrDB, "更新状态失败")
		return
	}
	response.SuccessMsg(c, "状态更新成功")
}

func (h *RoleHandler) PermissionShow(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	menuIDs, err := h.svc.GetMenus(uint(id))
	if err != nil {
		response.Fail(c, response.ErrDB, "获取权限失败")
		return
	}
	response.Success(c, gin.H{"menuIds": menuIDs})
}

func (h *RoleHandler) PermissionStore(c *gin.Context) {
	var req struct {
		RoleID  uint   `json:"roleId" binding:"required"`
		MenuIDs []uint `json:"menuIds"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ErrParams, "参数错误")
		return
	}
	if err := h.svc.SetMenus(req.RoleID, req.MenuIDs); err != nil {
		response.Fail(c, response.ErrDB, "保存权限失败")
		return
	}
	response.SuccessMsg(c, "权限保存成功")
}

func (h *RoleHandler) PermissionByUser(c *gin.Context) {
	userID, _ := c.Get("userId")
	var user model.User
	if err := h.db.Select("id", "isSuperAdmin").First(&user, userID.(uint)).Error; err != nil {
		response.Fail(c, response.ErrDB, "获取权限失败")
		return
	}

	menus, err := h.permissionMenusForUser(user)
	if err != nil {
		response.Fail(c, response.ErrDB, "获取权限失败")
		return
	}

	response.Success(c, h.buildPermissionTree(menus))
}

func (h *RoleHandler) ShowEmployee(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var userRoles []model.RbacUserRole
	h.db.Where("role_id = ?", id).Find(&userRoles)
	response.Success(c, userRoles)
}

func (h *RoleHandler) PermissionShowBySvc(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	menus, err := h.svc.GetMenus(uint(id))
	if err != nil {
		response.Fail(c, response.ErrDB, "获取权限失败")
		return
	}
	response.Success(c, menus)
}

func (h *RoleHandler) PermissionStoreBySvc(c *gin.Context) {
	var req struct {
		RoleID  uint   `json:"roleId" binding:"required"`
		MenuIDs []uint `json:"menuIds"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ErrParams, err.Error())
		return
	}
	if err := h.svc.SetMenus(req.RoleID, req.MenuIDs); err != nil {
		response.Fail(c, response.ErrDB, "保存权限失败")
		return
	}
	response.SuccessMsg(c, "权限保存成功")
}

func (h *RoleHandler) ListByOffset(c *gin.Context) {
	tenantID, _ := c.Get("tenantId")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	offset := (page - 1) * pageSize

	roles, total, err := h.svc.ListByOffset(tenantID.(uint), offset, pageSize)
	if err != nil {
		response.Fail(c, response.ErrDB, "获取角色列表失败")
		return
	}
	response.PageResult(c, roles, total, page, pageSize)
}

func (h *RoleHandler) permissionMenusForUser(user model.User) ([]model.RbacMenu, error) {
	var menus []model.RbacMenu

	query := h.db.Model(&model.RbacMenu{}).
		Where("status = 1 AND is_page_menu = 1").
		Order("sort ASC, id ASC")

	if user.IsSuperAdmin == 1 {
		if err := query.Find(&menus).Error; err != nil {
			return nil, err
		}
		return menus, nil
	}

	if err := query.
		Joins("JOIN mc_rbac_role_menu rm ON rm.menu_id = mc_rbac_menu.id").
		Joins("JOIN mc_rbac_user_role ur ON ur.role_id = rm.role_id").
		Where("ur.user_id = ?", user.ID).
		Distinct().
		Find(&menus).Error; err != nil {
		return nil, err
	}

	return menus, nil
}

func (h *RoleHandler) buildPermissionTree(menus []model.RbacMenu) []permissionMenu {
	childrenByParent := make(map[uint][]model.RbacMenu, len(menus))
	menuByID := make(map[uint]model.RbacMenu, len(menus))
	for _, menu := range menus {
		childrenByParent[menu.ParentID] = append(childrenByParent[menu.ParentID], menu)
		menuByID[menu.ID] = menu
	}

	rootParentIDs := []uint{0}
	for _, menu := range menus {
		if menu.ParentID != 0 {
			if _, exists := menuByID[menu.ParentID]; !exists {
				rootParentIDs = append(rootParentIDs, menu.ParentID)
			}
		}
	}

	tree := make([]permissionMenu, 0)
	visitedRoots := make(map[uint]bool)
	for _, rootParentID := range rootParentIDs {
		nodes := h.buildPermissionChildren(childrenByParent, rootParentID, map[uint]bool{})
		for _, node := range nodes {
			if visitedRoots[node.ID] {
				continue
			}
			visitedRoots[node.ID] = true
			tree = append(tree, node)
		}
	}

	return tree
}

func (h *RoleHandler) buildPermissionChildren(childrenByParent map[uint][]model.RbacMenu, parentID uint, path map[uint]bool) []permissionMenu {
	menus := childrenByParent[parentID]
	tree := make([]permissionMenu, 0, len(menus))
	for _, menu := range menus {
		if menu.ID == 0 || menu.ID == parentID || path[menu.ID] {
			continue
		}

		nextPath := make(map[uint]bool, len(path)+1)
		for id := range path {
			nextPath[id] = true
		}
		nextPath[menu.ID] = true

		node := permissionMenu{
			ID:             menu.ID,
			MenuID:         menu.ID,
			ParentID:       menu.ParentID,
			Name:           menu.Name,
			Level:          menu.Level,
			Icon:           menu.Icon,
			LinkType:       menu.LinkType,
			LinkURL:        strings.Replace(menu.LinkURL, "/dashboard", "", 1),
			DataPermission: menu.DataPermission,
			IsPageMenu:     menu.IsPageMenu,
			Children:       h.buildPermissionChildren(childrenByParent, menu.ID, nextPath),
		}
		tree = append(tree, node)
	}
	return tree
}
