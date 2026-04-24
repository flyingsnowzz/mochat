package dashboard

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"mochat-api-server/internal/model"
	"mochat-api-server/internal/pkg/response"
	"mochat-api-server/internal/service"
)

type MenuHandler struct {
	svc *service.RbacMenuService
}

type menuIndexItem struct {
	MenuID      uint            `json:"menuId"`
	MenuPath    string          `json:"menuPath"`
	Name        string          `json:"name"`
	Level       int             `json:"level"`
	LevelName   string          `json:"levelName"`
	ParentID    uint            `json:"parentId"`
	Icon        string          `json:"icon"`
	Status      string          `json:"status"`
	OperateName string          `json:"operateName"`
	UpdatedAt   string          `json:"updatedAt"`
	Children    []menuIndexItem `json:"children"`
}

func NewMenuHandler(db *gorm.DB) *MenuHandler {
	return &MenuHandler{svc: service.NewRbacMenuService(db)}
}

func (h *MenuHandler) Index(c *gin.Context) {
	name := strings.TrimSpace(c.Query("name"))
	menus, err := h.svc.List(name)
	if err != nil {
		response.Fail(c, response.ErrDB, "获取菜单列表失败")
		return
	}
	items := buildMenuIndexTree(menus)
	response.PageResult(c, items, int64(len(items)), 1, len(items))
}

func (h *MenuHandler) Select(c *gin.Context) {
	menus, err := h.svc.Select()
	if err != nil {
		response.Fail(c, response.ErrDB, "获取菜单选择列表失败")
		return
	}
	response.Success(c, buildMenuSelectTree(menus))
}

func (h *MenuHandler) Show(c *gin.Context) {
	id := getMenuID(c)
	menu, err := h.svc.GetByID(uint(id))
	if err != nil {
		response.Fail(c, response.ErrNotFound, "菜单不存在")
		return
	}
	response.Success(c, buildMenuShowData(menu))
}

func (h *MenuHandler) IconIndex(c *gin.Context) {
	response.Success(c, []string{})
}

func (h *MenuHandler) Store(c *gin.Context) {
	var menu model.RbacMenu
	if err := c.ShouldBindJSON(&menu); err != nil {
		response.Fail(c, response.ErrParams, "参数错误")
		return
	}
	if err := h.svc.Create(&menu); err != nil {
		response.Fail(c, response.ErrDB, "创建菜单失败")
		return
	}
	response.Success(c, menu)
}

func (h *MenuHandler) Update(c *gin.Context) {
	id := getMenuID(c)
	menu, err := h.svc.GetByID(uint(id))
	if err != nil {
		response.Fail(c, response.ErrNotFound, "菜单不存在")
		return
	}
	var req struct {
		MenuID         uint   `json:"menuId"`
		Name           string `json:"name"`
		Icon           string `json:"icon"`
		LinkURL        string `json:"linkUrl"`
		LinkType       int    `json:"linkType"`
		IsPageMenu     int    `json:"isPageMenu"`
		DataPermission int    `json:"dataPermission"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ErrParams, "参数错误")
		return
	}
	menu.Name = req.Name
	menu.Icon = req.Icon
	menu.LinkURL = req.LinkURL
	if req.LinkType > 0 {
		menu.LinkType = req.LinkType
	}
	if req.IsPageMenu > 0 {
		menu.IsPageMenu = req.IsPageMenu
	}
	if req.DataPermission > 0 {
		menu.DataPermission = req.DataPermission
	}
	if err := h.svc.Update(menu); err != nil {
		response.Fail(c, response.ErrDB, "更新菜单失败")
		return
	}
	response.Success(c, buildMenuShowData(menu))
}

func (h *MenuHandler) Destroy(c *gin.Context) {
	id := getMenuID(c)
	if err := h.svc.Delete(uint(id)); err != nil {
		response.Fail(c, response.ErrDB, "删除菜单失败")
		return
	}
	response.SuccessMsg(c, "删除成功")
}

func (h *MenuHandler) StatusUpdate(c *gin.Context) {
	var req struct {
		MenuID uint `json:"menuId" binding:"required"`
		Status int  `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ErrParams, "参数错误")
		return
	}
	menu, err := h.svc.GetByID(req.MenuID)
	if err != nil {
		response.Fail(c, response.ErrNotFound, "菜单不存在")
		return
	}
	menu.Status = req.Status
	if err := h.svc.Update(menu); err != nil {
		response.Fail(c, response.ErrDB, "更新状态失败")
		return
	}
	response.SuccessMsg(c, "状态更新成功")
}

func buildMenuIndexTree(menus []model.RbacMenu) []menuIndexItem {
	childrenMap := make(map[uint][]model.RbacMenu)
	for _, menu := range menus {
		childrenMap[menu.ParentID] = append(childrenMap[menu.ParentID], menu)
	}
	return buildMenuIndexChildren(childrenMap, 0, "")
}

func buildMenuIndexChildren(childrenMap map[uint][]model.RbacMenu, parentID uint, prefix string) []menuIndexItem {
	children := childrenMap[parentID]
	items := make([]menuIndexItem, 0, len(children))
	for idx, menu := range children {
		path := strconv.Itoa(idx + 1)
		if prefix != "" {
			path = prefix + "-" + path
		}
		updatedAt := ""
		if !menu.UpdatedAt.IsZero() {
			updatedAt = menu.UpdatedAt.Format("2006-01-02 15:04:05")
		}
		item := menuIndexItem{
			MenuID:      menu.ID,
			MenuPath:    path,
			Name:        menu.Name,
			Level:       menu.Level,
			LevelName:   strconv.Itoa(menu.Level),
			ParentID:    menu.ParentID,
			Icon:        menu.Icon,
			Status:      strconv.Itoa(menu.Status),
			OperateName: menu.OperateName,
			UpdatedAt:   updatedAt,
		}
		item.Children = buildMenuIndexChildren(childrenMap, menu.ID, path)
		items = append(items, item)
	}
	return items
}

func buildMenuSelectTree(menus []model.RbacMenu) []gin.H {
	childrenMap := make(map[uint][]model.RbacMenu)
	for _, menu := range menus {
		childrenMap[menu.ParentID] = append(childrenMap[menu.ParentID], menu)
	}
	var build func(parentID uint) []gin.H
	build = func(parentID uint) []gin.H {
		children := childrenMap[parentID]
		items := make([]gin.H, 0, len(children))
		for _, menu := range children {
			items = append(items, gin.H{
				"menuId":         menu.ID,
				"name":           menu.Name,
				"level":          menu.Level,
				"dataPermission": menu.DataPermission,
				"parentId":       menu.ParentID,
				"children":       build(menu.ID),
			})
		}
		return items
	}
	return build(0)
}

func buildMenuShowData(menu *model.RbacMenu) gin.H {
	return gin.H{
		"menuId":         menu.ID,
		"firstMenuId":    firstParentFromPath(menu.Path, 1),
		"secondMenuId":   firstParentFromPath(menu.Path, 2),
		"thirdMenuId":    firstParentFromPath(menu.Path, 3),
		"fourthMenuId":   firstParentFromPath(menu.Path, 4),
		"level":          menu.Level,
		"name":           menu.Name,
		"icon":           menu.Icon,
		"isPageMenu":     menu.IsPageMenu,
		"linkUrl":        menu.LinkURL,
		"linkType":       menu.LinkType,
		"dataPermission": menu.DataPermission,
		"status":         menu.Status,
	}
}

func firstParentFromPath(path string, index int) uint {
	if path == "" {
		return 0
	}
	parts := strings.Split(path, "#-#")
	if len(parts) <= index {
		return 0
	}
	value := strings.Trim(parts[index], "#")
	id, _ := strconv.ParseUint(value, 10, 32)
	return uint(id)
}

func getMenuID(c *gin.Context) uint64 {
	if id, err := strconv.ParseUint(c.Param("id"), 10, 32); err == nil && id > 0 {
		return id
	}
	if id, err := strconv.ParseUint(c.Query("menuId"), 10, 32); err == nil && id > 0 {
		return id
	}
	var req struct {
		MenuID uint `json:"menuId"`
	}
	if err := c.ShouldBindJSON(&req); err == nil && req.MenuID > 0 {
		return uint64(req.MenuID)
	}
	return 0
}
