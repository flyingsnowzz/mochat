package model

import (
	"time"

	"gorm.io/gorm"
)

// RoomWelcomeTemplate 入群欢迎语模板
type RoomWelcomeTemplate struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	CorpID      uint           `gorm:"column:corp_id" json:"corpId"`
	MsgText     string         `gorm:"column:msg_text" json:"msgText"`
	MsgComplex  string         `gorm:"column:msg_complex" json:"msgComplex"`
	ComplexType int            `gorm:"column:complex_type" json:"complexType"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (RoomWelcomeTemplate) TableName() string {
	return "mc_room_welcome_template"
}

// WorkRoomAutoPull 自动拉群
type WorkRoomAutoPull struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	CorpID       uint           `gorm:"column:corp_id" json:"corpId"`
	QrcodeName   string         `gorm:"column:qrcode_name" json:"qrcodeName"`
	IsVerified   int            `gorm:"column:is_verified" json:"isVerified"`
	LeadingWords string         `gorm:"column:leading_words" json:"leadingWords"`
	Employees    string         `gorm:"column:employees" json:"employees"`
	Tags         string         `gorm:"column:tags" json:"tags"`
	Rooms        string         `gorm:"column:rooms" json:"rooms"`
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (WorkRoomAutoPull) TableName() string {
	return "mc_work_room_auto_pull"
}

// RoomTagPull 标签建群
type RoomTagPull struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	Name          string         `gorm:"column:name" json:"name"`
	Employees     string         `gorm:"column:employees" json:"employees"`
	ChooseContact string         `gorm:"column:choose_contact" json:"chooseContact"`
	Guide         string         `gorm:"column:guide" json:"guide"`
	Rooms         string         `gorm:"column:rooms" json:"rooms"`
	FilterContact int            `gorm:"column:filter_contact" json:"filterContact"`
	ContactNum    int            `gorm:"column:contact_num" json:"contactNum"`
	WxTid         string         `gorm:"column:wx_tid" json:"wxTid"`
	TenantID      int            `gorm:"column:tenant_id" json:"tenantId"`
	CorpID        int            `gorm:"column:corp_id" json:"corpId"`
	CreateUserID  int            `gorm:"column:create_user_id" json:"createUserId"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
	CreatedAt     time.Time      `json:"createdAt"`
	UpdatedAt     time.Time      `json:"updatedAt"`
}

// TableName 指定表名
func (RoomTagPull) TableName() string {
	return "mc_room_tag_pull"
}

// ChannelCode 渠道码
type ChannelCode struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	CorpID          uint           `gorm:"column:corp_id" json:"corpId"`
	GroupID         uint           `gorm:"column:group_id" json:"groupId"`
	Name            string         `gorm:"column:name" json:"name"`
	QrcodeURL       string         `gorm:"column:qrcode_url" json:"qrcodeUrl"`
	WxConfigID      string         `gorm:"column:wx_config_id" json:"wxConfigId"`
	AutoAddFriend   int            `gorm:"column:auto_add_friend" json:"autoAddFriend"`
	Tags            string         `gorm:"column:tags" json:"tags"`
	Type            uint           `gorm:"column:type" json:"type"`
	DrainageEmployee string `gorm:"column:drainage_employee" json:"drainageEmployee"`
	WelcomeMessage  string         `gorm:"column:welcome_message" json:"welcomeMessage"`
	CreatedAt       time.Time      `json:"createdAt"`
	UpdatedAt       time.Time      `json:"updatedAt"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (ChannelCode) TableName() string {
	return "mc_channel_code"
}

// ChannelCodeGroup 渠道码分组
type ChannelCodeGroup struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CorpID    uint           `gorm:"column:corp_id" json:"corpId"`
	Name      string         `gorm:"column:name" json:"name"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (ChannelCodeGroup) TableName() string {
	return "mc_channel_code_group"
}

// ContactMessageBatchSend 客户群发
type ContactMessageBatchSend struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CorpID    uint           `gorm:"column:corp_id" json:"corpId"`
	Name      string         `gorm:"column:name" json:"name"`
	Status    int            `gorm:"column:status" json:"status"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (ContactMessageBatchSend) TableName() string {
	return "mc_contact_message_batch_send"
}

// WorkTransferLog 客户迁移日志
type WorkTransferLog struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CorpID    uint           `gorm:"column:corp_id" json:"corpId"`
	Status    int            `gorm:"column:status" json:"status"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (WorkTransferLog) TableName() string {
	return "mc_work_transfer_log"
}

// WorkUnassigned 未分配客户
type WorkUnassigned struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CorpID    uint           `gorm:"column:corp_id" json:"corpId"`
	Status    int            `gorm:"column:status" json:"status"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (WorkUnassigned) TableName() string {
	return "mc_work_unassigned"
}

// RoomMessageBatchSend 客户群群发
type RoomMessageBatchSend struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CorpID    uint           `gorm:"column:corp_id" json:"corpId"`
	Name      string         `gorm:"column:name" json:"name"`
	Status    int            `gorm:"column:status" json:"status"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (RoomMessageBatchSend) TableName() string {
	return "mc_room_message_batch_send"
}

// WorkFission 企微任务宝
type WorkFission struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CorpID    uint           `gorm:"column:corp_id" json:"corpId"`
	Name      string         `gorm:"column:name" json:"name"`
	Status    int            `gorm:"column:status" json:"status"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (WorkFission) TableName() string {
	return "mc_work_fission"
}
