package model

import (
	"time"

	"gorm.io/gorm"
)

// RoomWelcomeTemplate 入群欢迎语模板
type RoomWelcomeTemplate struct {
	ID                uint           `gorm:"primaryKey" json:"id"`
	CorpID            string         `gorm:"column:corp_id" json:"corpId"`
	MsgText           string         `gorm:"column:msg_text" json:"msgText"`
	ComplexType       string         `gorm:"column:complex_type" json:"complexType"`
	MsgComplex        string         `gorm:"column:msg_complex;type:json" json:"msgComplex"`
	ComplexTemplateID string         `gorm:"column:complex_template_id" json:"complexTemplateId"`
	CreateUserID      int            `gorm:"column:create_user_id" json:"createUserId"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`
	CreatedAt         time.Time      `json:"createdAt"`
	UpdatedAt         time.Time      `json:"updatedAt"`
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
	QrcodeURL    string         `gorm:"column:qrcode_url" json:"qrcodeUrl"`
	WxConfigID   string         `gorm:"column:wx_config_id" json:"wxConfigId"`
	IsVerified   uint           `gorm:"column:is_verified" json:"isVerified"`
	LeadingWords string         `gorm:"column:leading_words;type:text" json:"leadingWords"`
	Tags         string         `gorm:"column:tags;type:json" json:"tags"`
	Employees    string         `gorm:"column:employees;type:json" json:"employees"`
	Rooms        string         `gorm:"column:rooms;type:json" json:"rooms"`
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
	ID               uint           `gorm:"primaryKey" json:"id"`
	CorpID           uint           `gorm:"column:corp_id" json:"corpId"`
	GroupID          uint           `gorm:"column:group_id" json:"groupId"`
	Name             string         `gorm:"column:name" json:"name"`
	QrcodeURL        string         `gorm:"column:qrcode_url" json:"qrcodeUrl"`
	WxConfigID       string         `gorm:"column:wx_config_id" json:"wxConfigId"`
	AutoAddFriend    int            `gorm:"column:auto_add_friend" json:"autoAddFriend"`
	Tags             string         `gorm:"column:tags" json:"tags"`
	Type             uint           `gorm:"column:type" json:"type"`
	DrainageEmployee string         `gorm:"column:drainage_employee" json:"drainageEmployee"`
	WelcomeMessage   string         `gorm:"column:welcome_message" json:"welcomeMessage"`
	CreatedAt        time.Time      `json:"createdAt"`
	UpdatedAt        time.Time      `json:"updatedAt"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
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
	ID                 uint           `gorm:"primaryKey" json:"id"`
	CorpID             uint           `gorm:"column:corp_id" json:"corpId"`
	UserID             uint           `gorm:"column:user_id" json:"userId"`
	UserName           string         `gorm:"column:user_name" json:"userName"`
	EmployeeIDs        string         `gorm:"column:employee_ids;type:json" json:"employeeIds"`
	FilterParams       string         `gorm:"column:filter_params;type:json" json:"filterParams"`
	FilterParamsDetail string         `gorm:"column:filter_params_detail;type:json" json:"filterParamsDetail"`
	Content            string         `gorm:"column:content;type:json" json:"content"`
	SendWay            int            `gorm:"column:send_way" json:"sendWay"`
	DefiniteTime       *time.Time     `gorm:"column:definite_time" json:"definiteTime"`
	SendTime           *time.Time     `gorm:"column:send_time" json:"sendTime"`
	SendEmployeeTotal  uint           `gorm:"column:send_employee_total" json:"sendEmployeeTotal"`
	SendContactTotal   uint           `gorm:"column:send_contact_total" json:"sendContactTotal"`
	SendTotal          uint           `gorm:"column:send_total" json:"sendTotal"`
	NotSendTotal       uint           `gorm:"column:not_send_total" json:"notSendTotal"`
	ReceivedTotal      uint           `gorm:"column:received_total" json:"receivedTotal"`
	NotReceivedTotal   uint           `gorm:"column:not_received_total" json:"notReceivedTotal"`
	ReceiveLimitTotal  uint           `gorm:"column:receive_limit_total" json:"receiveLimitTotal"`
	NotFriendTotal     uint           `gorm:"column:not_friend_total" json:"notFriendTotal"`
	SendStatus         int            `gorm:"column:send_status" json:"sendStatus"`
	CreatedAt          time.Time      `json:"createdAt"`
	UpdatedAt          time.Time      `json:"updatedAt"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (ContactMessageBatchSend) TableName() string {
	return "mc_contact_message_batch_send"
}

// WorkTransferLog 客户迁移日志
type WorkTransferLog struct {
	ID                 uint      `gorm:"primaryKey" json:"id"`
	CorpID             uint      `gorm:"column:corp_id" json:"corpId"`
	Status             int       `gorm:"column:status" json:"status"`
	Type               int       `gorm:"column:type" json:"type"`
	Name               string    `gorm:"column:name" json:"name"`
	ContactID          string    `gorm:"column:contact_id" json:"contactId"`
	HandoverEmployeeID string    `gorm:"column:handover_employee_id" json:"handoverEmployeeId"`
	TakeoverEmployeeID string    `gorm:"column:takeover_employee_id" json:"takeoverEmployeeId"`
	State              int       `gorm:"column:state" json:"state"`
	CreatedAt          time.Time `json:"createdAt"`
	UpdatedAt          time.Time `json:"updatedAt"`
}

// TableName 指定表名
func (WorkTransferLog) TableName() string {
	return "mc_work_transfer_log"
}

// WorkUnassigned 未分配客户
type WorkUnassigned struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	CorpID         uint           `gorm:"column:corp_id" json:"corpId"`
	HandoverUserID string         `gorm:"column:handover_userid" json:"handoverUserid"`
	ExternalUserID string         `gorm:"column:external_userid" json:"externalUserid"`
	DimissionTime  int            `gorm:"column:dimission_time" json:"dimissionTime"`
	CreatedAt      time.Time      `json:"createdAt"`
	UpdatedAt      time.Time      `json:"updatedAt"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (WorkUnassigned) TableName() string {
	return "mc_work_unassigned"
}

// RoomMessageBatchSend 客户群群发
type RoomMessageBatchSend struct {
	ID                uint           `gorm:"primaryKey" json:"id"`
	CorpID            uint           `gorm:"column:corp_id" json:"corpId"`
	UserID            uint           `gorm:"column:user_id" json:"userId"`
	UserName          string         `gorm:"column:user_name" json:"userName"`
	EmployeeIDs       string         `gorm:"column:employee_ids;type:json" json:"employeeIds"`
	BatchTitle        string         `gorm:"column:batch_title" json:"batchTitle"`
	Content           string         `gorm:"column:content;type:json" json:"content"`
	SendWay           int            `gorm:"column:send_way" json:"sendWay"`
	DefiniteTime      *time.Time     `gorm:"column:definite_time" json:"definiteTime"`
	SendTime          *time.Time     `gorm:"column:send_time" json:"sendTime"`
	SendRoomTotal     uint           `gorm:"column:send_room_total" json:"sendRoomTotal"`
	SendEmployeeTotal uint           `gorm:"column:send_employee_total" json:"sendEmployeeTotal"`
	SendTotal         uint           `gorm:"column:send_total" json:"sendTotal"`
	NotSendTotal      uint           `gorm:"column:not_send_total" json:"notSendTotal"`
	ReceivedTotal     uint           `gorm:"column:received_total" json:"receivedTotal"`
	NotReceivedTotal  uint           `gorm:"column:not_received_total" json:"notReceivedTotal"`
	SendStatus        int            `gorm:"column:send_status" json:"sendStatus"`
	CreatedAt         time.Time      `json:"createdAt"`
	UpdatedAt         time.Time      `json:"updatedAt"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (RoomMessageBatchSend) TableName() string {
	return "mc_room_message_batch_send"
}

// WorkFission 企微任务宝
type WorkFission struct {
	ID                    uint           `gorm:"primaryKey" json:"id"`
	CorpID                uint           `gorm:"column:corp_id" json:"corpId"`
	ActiveName            string         `gorm:"column:active_name" json:"activeName"`
	ServiceEmployees      string         `gorm:"column:service_employees;type:json" json:"serviceEmployees"`
	AutoPass              int            `gorm:"column:auto_pass" json:"autoPass"`
	AutoAddTag            int            `gorm:"column:auto_add_tag" json:"autoAddTag"`
	ContactTags           string         `gorm:"column:contact_tags;type:json" json:"contactTags"`
	EndTime               *time.Time     `gorm:"column:end_time" json:"endTime"`
	QrCodeInvalid         int            `gorm:"column:qr_code_invalid" json:"qrCodeInvalid"`
	Tasks                 string         `gorm:"column:tasks;type:json" json:"tasks"`
	NewFriend             int            `gorm:"column:new_friend" json:"newFriend"`
	DeleteInvalid         int            `gorm:"column:delete_invalid" json:"deleteInvalid"`
	ReceivePrize          int            `gorm:"column:receive_prize" json:"receivePrize"`
	ReceivePrizeEmployees string         `gorm:"column:receive_prize_employees;type:json" json:"receivePrizeEmployees"`
	ReceiveLinks          string         `gorm:"column:receive_links;type:json" json:"receiveLinks"`
	ReceiveQrcode         string         `gorm:"column:receive_qrcode;type:json" json:"receiveQrcode"`
	CreateUserID          int            `gorm:"column:create_user_id" json:"createUserId"`
	DeletedAt             gorm.DeletedAt `gorm:"index" json:"-"`
	CreatedAt             time.Time      `json:"createdAt"`
	UpdatedAt             time.Time      `json:"updatedAt"`
}

// TableName 指定表名
func (WorkFission) TableName() string {
	return "mc_work_fission"
}
