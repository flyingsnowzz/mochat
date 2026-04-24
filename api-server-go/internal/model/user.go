package model

import "time"

type User struct {
	ID           uint       `gorm:"column:id;primaryKey" json:"id"`
	Phone        string     `gorm:"column:phone;size:11;not null" json:"phone"`
	Password     string     `gorm:"column:password;size:255;not null" json:"-"`
	Name         string     `gorm:"column:name;size:255;default:''" json:"name"`
	Gender       int        `gorm:"column:gender;default:0" json:"gender"`
	Department   string     `gorm:"column:department;size:255;default:''" json:"department"`
	Position     string     `gorm:"column:position;size:255;default:''" json:"position"`
	LoginTime    *time.Time `gorm:"column:login_time" json:"loginTime"`
	Status       int        `gorm:"column:status;default:1" json:"status"`
	TenantID     uint       `gorm:"column:tenant_id;default:1" json:"tenantId"`
	IsSuperAdmin int        `gorm:"column:isSuperAdmin;default:0" json:"isSuperAdmin"`
	CreatedAt    time.Time  `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt    time.Time  `gorm:"column:updated_at" json:"updatedAt"`
	DeletedAt    *time.Time `gorm:"column:deleted_at;index" json:"deletedAt"`
}

func (User) TableName() string { return "mc_user" }
