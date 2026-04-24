package model

import "time"

type WorkContactTag struct {
	ID                uint       `gorm:"primaryKey" json:"id"`
	WxContactTagID    string     `gorm:"column:wx_contact_tag_id;size:100;default:''" json:"wxContactTagId"`
	CorpID            uint       `gorm:"column:corp_id;default:0;index" json:"corpId"`
	Name              string     `gorm:"size:50;default:''" json:"name"`
	Order             int        `gorm:"default:0" json:"order"`
	ContactTagGroupID uint       `gorm:"column:contact_tag_group_id;default:0;index" json:"contactTagGroupId"`
	CreatedAt         time.Time  `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt         time.Time  `gorm:"column:updated_at" json:"updatedAt"`
	DeletedAt         *time.Time `gorm:"column:deleted_at;index" json:"deletedAt"`
}

func (WorkContactTag) TableName() string { return "mc_work_contact_tag" }

type WorkContactTagGroup struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	WxGroupID string     `gorm:"column:wx_group_id;size:100;default:''" json:"wxGroupId"`
	CorpID    uint       `gorm:"column:corp_id;default:0;index" json:"corpId"`
	GroupName string     `gorm:"column:group_name;size:50;default:''" json:"groupName"`
	Order     int        `gorm:"default:0" json:"order"`
	CreatedAt time.Time  `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt time.Time  `gorm:"column:updated_at" json:"updatedAt"`
	DeletedAt *time.Time `gorm:"column:deleted_at;index" json:"deletedAt"`
}

func (WorkContactTagGroup) TableName() string { return "mc_work_contact_tag_group" }

type WorkContactTagPivot struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	ContactID    uint      `gorm:"column:contact_id;not null;index" json:"contactId"`
	EmployeeID   uint      `gorm:"column:employee_id;not null;index" json:"employeeId"`
	ContactTagID uint      `gorm:"column:contact_tag_id;not null;index" json:"contactTagId"`
	Type         int       `gorm:"default:0" json:"type"`
	CreatedAt    time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt    time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

func (WorkContactTagPivot) TableName() string { return "mc_work_contact_tag_pivot" }
