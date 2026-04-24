package model

import "time"

type WorkRoom struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	CorpID      uint       `gorm:"column:corp_id;default:0;index" json:"corpId"`
	WxChatID    string     `gorm:"column:wx_chat_id;size:100;default:''" json:"wxChatId"`
	Name        string     `gorm:"size:50;default:''" json:"name"`
	OwnerID     uint       `gorm:"column:owner_id;default:0" json:"ownerId"`
	Notice      string     `gorm:"size:500;default:''" json:"notice"`
	Status      int        `gorm:"default:0" json:"status"`
	CreateTime  time.Time  `gorm:"column:create_time" json:"createTime"`
	RoomMax     int        `gorm:"column:room_max;default:0" json:"roomMax"`
	RoomGroupID uint       `gorm:"column:room_group_id;default:0;index" json:"roomGroupId"`
	CreatedAt   time.Time  `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt   time.Time  `gorm:"column:updated_at" json:"updatedAt"`
	DeletedAt   *time.Time `gorm:"column:deleted_at;index" json:"deletedAt"`
}

func (WorkRoom) TableName() string { return "mc_work_room" }

type WorkRoomGroup struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	CorpID    uint       `gorm:"column:corp_id;default:0;index" json:"corpId"`
	Name      string     `gorm:"size:50;default:''" json:"name"`
	CreatedAt time.Time  `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt time.Time  `gorm:"column:updated_at" json:"updatedAt"`
	DeletedAt *time.Time `gorm:"column:deleted_at;index" json:"deletedAt"`
}

func (WorkRoomGroup) TableName() string { return "mc_work_room_group" }

type WorkContactRoom struct {
	ID         uint       `gorm:"primaryKey" json:"id"`
	WxUserID   string     `gorm:"column:wx_user_id;size:100;default:''" json:"wxUserId"`
	ContactID  uint       `gorm:"column:contact_id;default:0;index" json:"contactId"`
	EmployeeID uint       `gorm:"column:employee_id;default:0" json:"employeeId"`
	Unionid    string     `gorm:"size:100;default:''" json:"unionid"`
	RoomID     uint       `gorm:"column:room_id;default:0;index" json:"roomId"`
	JoinScene  int        `gorm:"column:join_scene;default:0" json:"joinScene"`
	Type       int        `gorm:"default:0" json:"type"`
	Status     int        `gorm:"default:0" json:"status"`
	JoinTime   time.Time  `gorm:"column:join_time" json:"joinTime"`
	OutTime    *time.Time `gorm:"column:out_time" json:"outTime"`
	CreatedAt  time.Time  `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt  time.Time  `gorm:"column:updated_at" json:"updatedAt"`
}

func (WorkContactRoom) TableName() string { return "mc_work_contact_room" }
