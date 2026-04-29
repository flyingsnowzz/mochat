package model

import "time"

type WorkAgent struct {
	ID                 uint      `gorm:"primaryKey" json:"id"`
	CorpID             uint      `gorm:"column:corp_id;default:0;index" json:"corpId"`
	WxAgentID          int       `gorm:"column:wx_agent_id;default:0" json:"wxAgentId"`
	WxSecret           string    `gorm:"column:wx_secret;size:100;default:''" json:"wxSecret"`
	Name               string    `gorm:"size:50;default:''" json:"name"`
	SquareLogoURL      string    `gorm:"column:square_logo_url;size:500;default:''" json:"squareLogoUrl"`
	Description        string    `gorm:"size:500;default:''" json:"description"`
	Close              int       `gorm:"default:0" json:"close"`
	RedirectDomain     string    `gorm:"column:redirect_domain;size:200;default:''" json:"redirectDomain"`
	ReportLocationFlag int       `gorm:"column:report_location_flag;default:0" json:"reportLocationFlag"`
	IsReportEnter      int       `gorm:"column:is_reportenter;default:0" json:"isReportenter"`
	HomeURL            string    `gorm:"column:home_url;size:500;default:''" json:"homeUrl"`
	CreatedAt          time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt          time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

func (WorkAgent) TableName() string { return "mc_work_agent" }

type ChatTool struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	PageName  string    `gorm:"column:page_name;size:50;default:''" json:"pageName"`
	PageFlag  string    `gorm:"column:page_flag;size:50;default:''" json:"pageFlag"`
	Status    int       `gorm:"default:1" json:"status"`
	CreatedAt time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

func (ChatTool) TableName() string { return "mc_chat_tool" }

type Medium struct {
	ID             uint       `gorm:"primaryKey" json:"id"`
	MediaID        string     `gorm:"column:media_id;size:255;default:''" json:"mediaId"`
	LastUploadTime uint       `gorm:"column:last_upload_time;default:0" json:"lastUploadTime"`
	Type           int        `gorm:"default:1" json:"type"`
	IsSync         int        `gorm:"column:is_sync;default:1" json:"isSync"`
	Content        string     `gorm:"type:json" json:"content"`
	CorpID         uint       `gorm:"column:corp_id;default:0;index" json:"corpId"`
	MediumGroupID  uint       `gorm:"column:medium_group_id;default:0;index" json:"mediumGroupId"`
	UserID         uint       `gorm:"column:user_id;default:0" json:"userId"`
	UserName       string     `gorm:"column:user_name;size:255;default:''" json:"userName"`
	CreatedAt      time.Time  `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt      time.Time  `gorm:"column:updated_at" json:"updatedAt"`
	DeletedAt      *time.Time `gorm:"column:deleted_at;index" json:"deletedAt"`
}

func (Medium) TableName() string { return "mc_medium" }

type MediumGroup struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	CorpID    uint       `gorm:"column:corp_id;default:0;index" json:"corpId"`
	Name      string     `gorm:"size:50;default:''" json:"name"`
	Order     int        `gorm:"default:0" json:"order"`
	CreatedAt time.Time  `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt time.Time  `gorm:"column:updated_at" json:"updatedAt"`
	DeletedAt *time.Time `gorm:"column:deleted_at;index" json:"deletedAt"`
}

func (MediumGroup) TableName() string { return "mc_medium_group" }

type Greeting struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	CorpID    uint       `gorm:"column:corp_id;default:0;index" json:"corpId"`
	Type      string     `gorm:"size:255;default:''" json:"type"`
	Words     string     `gorm:"type:text" json:"words"`
	MediumID  uint       `gorm:"column:medium_id;default:0" json:"mediumId"`
	RangeType int        `gorm:"column:range_type;default:1" json:"rangeType"`
	Employees string     `gorm:"type:json" json:"employees"`
	CreatedAt time.Time  `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt time.Time  `gorm:"column:updated_at" json:"updatedAt"`
	DeletedAt *time.Time `gorm:"column:deleted_at;index" json:"deletedAt"`
}

func (Greeting) TableName() string { return "mc_greeting" }
