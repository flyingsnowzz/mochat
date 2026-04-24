package model

import "time"

type RbacMenu struct {
	ID             uint       `gorm:"column:id;primaryKey" json:"id"`
	ParentID       uint       `gorm:"column:parent_id;default:0" json:"parentId"`
	Name           string     `gorm:"column:name;size:50;not null" json:"name"`
	Level          int        `gorm:"column:level;default:0" json:"level"`
	Path           string     `gorm:"column:path;size:255;default:''" json:"path"`
	Icon           string     `gorm:"column:icon;size:100;default:''" json:"icon"`
	Status         int        `gorm:"column:status;default:1" json:"status"`
	LinkType       int        `gorm:"column:link_type;default:0" json:"linkType"`
	IsPageMenu     int        `gorm:"column:is_page_menu;default:1" json:"isPageMenu"`
	LinkURL        string     `gorm:"column:link_url;size:255;default:''" json:"linkUrl"`
	DataPermission int        `gorm:"column:data_permission;default:0" json:"dataPermission"`
	OperateID      uint       `gorm:"column:operate_id;default:0" json:"operateId"`
	OperateName    string     `gorm:"column:operate_name;size:50;default:''" json:"operateName"`
	Sort           int        `gorm:"column:sort;default:0" json:"sort"`
	CreatedAt      time.Time  `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt      time.Time  `gorm:"column:updated_at" json:"updatedAt"`
	DeletedAt      *time.Time `gorm:"column:deleted_at;index" json:"deletedAt"`
}

func (RbacMenu) TableName() string { return "mc_rbac_menu" }

type RbacRole struct {
	ID             uint       `gorm:"column:id;primaryKey" json:"id"`
	TenantID       uint       `gorm:"column:tenant_id;default:0" json:"tenantId"`
	Name           string     `gorm:"column:name;size:50;not null" json:"name"`
	Remarks        string     `gorm:"column:remarks;size:500;default:''" json:"remarks"`
	Status         int        `gorm:"column:status;default:1" json:"status"`
	OperateID      uint       `gorm:"column:operate_id;default:0" json:"operateId"`
	OperateName    string     `gorm:"column:operate_name;size:50;default:''" json:"operateName"`
	DataPermission int        `gorm:"column:data_permission;default:0" json:"dataPermission"`
	CreatedAt      time.Time  `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt      time.Time  `gorm:"column:updated_at" json:"updatedAt"`
	DeletedAt      *time.Time `gorm:"column:deleted_at;index" json:"deletedAt"`
}

func (RbacRole) TableName() string { return "mc_rbac_role" }

type RbacRoleMenu struct {
	ID        uint      `gorm:"column:id;primaryKey" json:"id"`
	RoleID    uint      `gorm:"column:role_id;not null;index" json:"roleId"`
	MenuID    uint      `gorm:"column:menu_id;not null;index" json:"menuId"`
	CreatedAt time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

func (RbacRoleMenu) TableName() string { return "mc_rbac_role_menu" }

type RbacUserRole struct {
	ID        uint      `gorm:"column:id;primaryKey" json:"id"`
	UserID    uint      `gorm:"column:user_id;not null;index" json:"userId"`
	RoleID    uint      `gorm:"column:role_id;not null;index" json:"roleId"`
	CreatedAt time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

func (RbacUserRole) TableName() string { return "mc_rbac_user_role" }

type RbacUserDepartment struct {
	ID               uint      `gorm:"column:id;primaryKey" json:"id"`
	CorpID           uint      `gorm:"column:corp_id;default:0" json:"corpId"`
	UserID           uint      `gorm:"column:user_id;default:0" json:"userId"`
	WorkDepartmentID uint      `gorm:"column:work_department_id;default:0" json:"workDepartmentId"`
	CreatedAt        time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt        time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

func (RbacUserDepartment) TableName() string { return "mc_rbac_user_department" }
