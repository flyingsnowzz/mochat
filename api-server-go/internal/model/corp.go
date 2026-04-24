package model

import "time"

type Corp struct {
	ID             uint       `gorm:"column:id;primaryKey" json:"id"`
	Name           string     `gorm:"column:name;size:255;not null" json:"name"`
	WxCorpid       string     `gorm:"column:wx_corpid;size:50;default:''" json:"wxCorpid"`
	SocialCode     string     `gorm:"column:social_code;size:50;default:''" json:"socialCode"`
	EmployeeSecret string     `gorm:"column:employee_secret;size:100;default:''" json:"employeeSecret"`
	EventCallback  string     `gorm:"column:event_callback;size:255;default:''" json:"eventCallback"`
	ContactSecret  string     `gorm:"column:contact_secret;size:100;default:''" json:"contactSecret"`
	Token          string     `gorm:"size:100;default:''" json:"token"`
	EncodingAesKey string     `gorm:"column:encoding_aes_key;size:100;default:''" json:"encodingAesKey"`
	TenantID       uint       `gorm:"column:tenant_id;default:0" json:"tenantId"`
	CreatedAt      time.Time  `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt      time.Time  `gorm:"column:updated_at" json:"updatedAt"`
	DeletedAt      *time.Time `gorm:"column:deleted_at;index" json:"deletedAt"`
}

func (Corp) TableName() string { return "mc_corp" }

type Tenant struct {
	ID              uint       `gorm:"column:id;primaryKey" json:"id"`
	Name            string     `gorm:"column:name;size:255;not null" json:"name"`
	Status          int        `gorm:"default:1" json:"status"`
	Logo            string     `gorm:"size:255;default:''" json:"logo"`
	LoginBackground string     `gorm:"column:login_background;size:255;default:''" json:"loginBackground"`
	URL             string     `gorm:"size:255;default:''" json:"url"`
	Copyright       string     `gorm:"size:255;default:''" json:"copyright"`
	ServerIps       string     `gorm:"column:server_ips;size:500;default:''" json:"serverIps"`
	CreatedAt       time.Time  `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt       time.Time  `gorm:"column:updated_at" json:"updatedAt"`
	DeletedAt       *time.Time `gorm:"column:deleted_at;index" json:"deletedAt"`
}

func (Tenant) TableName() string { return "mc_tenant" }

type WorkUpdateTime struct {
	ID             uint      `gorm:"column:id;primaryKey" json:"id"`
	CorpID         uint      `gorm:"column:corp_id;default:0" json:"corpId"`
	Type           int       `gorm:"default:0" json:"type"`
	LastUpdateTime string    `gorm:"column:last_update_time;size:50;default:''" json:"lastUpdateTime"`
	ErrorMsg       string    `gorm:"column:error_msg;size:500;default:''" json:"errorMsg"`
	CreatedAt      time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt      time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

func (WorkUpdateTime) TableName() string { return "mc_work_update_time" }
