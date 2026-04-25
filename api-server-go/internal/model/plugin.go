package model

import "time"

// Channel Code Plugin
type ChannelCode struct {
	ID               uint       `gorm:"column:id;primaryKey" json:"id"`
	CorpID           uint       `gorm:"column:corp_id;default:0;index" json:"corpId"`
	GroupID          uint       `gorm:"column:group_id;default:0;index" json:"groupId"`
	Name             string     `gorm:"column:name;size:255;default:''" json:"name"`
	QrcodeURL        string     `gorm:"column:qrcode_url;size:500;default:''" json:"qrcodeUrl"`
	WxConfigID       string     `gorm:"column:wx_config_id;size:100;default:''" json:"wxConfigId"`
	AutoAddFriend    int        `gorm:"column:auto_add_friend;default:0" json:"autoAddFriend"`
	Tags             string     `gorm:"column:tags;type:json" json:"tags"`
	Type             int        `gorm:"default:0" json:"type"`
	DrainageEmployee string     `gorm:"column:drainage_employee;type:json" json:"drainageEmployee"`
	WelcomeMessage   string     `gorm:"column:welcome_message;type:json" json:"welcomeMessage"`
	CreatedAt        time.Time  `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt        time.Time  `gorm:"column:updated_at" json:"updatedAt"`
	DeletedAt        *time.Time `gorm:"column:deleted_at;index" json:"deletedAt"`
}

func (ChannelCode) TableName() string { return "mc_channel_code" }

type ChannelCodeGroup struct {
	ID        uint       `gorm:"column:id;primaryKey" json:"id"`
	CorpID    uint       `gorm:"column:corp_id;default:0;index" json:"corpId"`
	Name      string     `gorm:"column:name;size:255;default:''" json:"name"`
	CreatedAt time.Time  `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt time.Time  `gorm:"column:updated_at" json:"updatedAt"`
	DeletedAt *time.Time `gorm:"column:deleted_at;index" json:"deletedAt"`
}

func (ChannelCodeGroup) TableName() string { return "mc_channel_code_group" }

// Contact Message Batch Send Plugin
type ContactMessageBatchSend struct {
	ID                 uint       `gorm:"primaryKey" json:"id"`
	CorpID             uint       `gorm:"column:corp_id;default:0;index" json:"corpId"`
	UserID             uint       `gorm:"column:user_id;default:0" json:"userId"`
	UserName           string     `gorm:"column:user_name;size:50;default:''" json:"userName"`
	EmployeeIDs        string     `gorm:"column:employee_ids;size:2000;default:''" json:"employeeIds"`
	FilterParams       string     `gorm:"column:filter_params;size:2000;default:''" json:"filterParams"`
	FilterParamsDetail string     `gorm:"column:filter_params_detail;size:2000;default:''" json:"filterParamsDetail"`
	Content            string     `gorm:"size:2000;default:''" json:"content"`
	SendWay            int        `gorm:"column:send_way;default:0" json:"sendWay"`
	DefiniteTime       *time.Time `gorm:"column:definite_time" json:"definiteTime"`
	SendTime           *time.Time `gorm:"column:send_time" json:"sendTime"`
	SendEmployeeTotal  int        `gorm:"column:send_employee_total;default:0" json:"sendEmployeeTotal"`
	SendContactTotal   int        `gorm:"column:send_contact_total;default:0" json:"sendContactTotal"`
	SendTotal          int        `gorm:"column:send_total;default:0" json:"sendTotal"`
	NotSendTotal       int        `gorm:"column:not_send_total;default:0" json:"notSendTotal"`
	ReceivedTotal      int        `gorm:"column:received_total;default:0" json:"receivedTotal"`
	NotReceivedTotal   int        `gorm:"column:not_received_total;default:0" json:"notReceivedTotal"`
	ReceiveLimitTotal  int        `gorm:"column:receive_limit_total;default:0" json:"receiveLimitTotal"`
	NotFriendTotal     int        `gorm:"column:not_friend_total;default:0" json:"notFriendTotal"`
	SendStatus         int        `gorm:"column:send_status;default:0" json:"sendStatus"`
	CreatedAt          time.Time  `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt          time.Time  `gorm:"column:updated_at" json:"updatedAt"`
	DeletedAt          *time.Time `gorm:"column:deleted_at;index" json:"deletedAt"`
}

func (ContactMessageBatchSend) TableName() string { return "mc_contact_message_batch_send" }

type ContactMessageBatchSendEmployee struct {
	ID               uint       `gorm:"primaryKey" json:"id"`
	BatchID          uint       `gorm:"column:batch_id;not null;index" json:"batchId"`
	EmployeeID       uint       `gorm:"column:employee_id;default:0;index" json:"employeeId"`
	WxUserID         string     `gorm:"column:wx_user_id;size:100;default:''" json:"wxUserId"`
	SendContactTotal int        `gorm:"column:send_contact_total;default:0" json:"sendContactTotal"`
	ErrCode          int        `gorm:"column:err_code;default:0" json:"errCode"`
	ErrMsg           string     `gorm:"column:err_msg;size:200;default:''" json:"errMsg"`
	MsgID            string     `gorm:"column:msg_id;size:100;default:''" json:"msgId"`
	SendTime         *time.Time `gorm:"column:send_time" json:"sendTime"`
	LastSyncTime     *time.Time `gorm:"column:last_sync_time" json:"lastSyncTime"`
	Status           int        `gorm:"default:0" json:"status"`
	ReceiveStatus    int        `gorm:"column:receive_status;default:0" json:"receiveStatus"`
	CreatedAt        time.Time  `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt        time.Time  `gorm:"column:updated_at" json:"updatedAt"`
}

func (ContactMessageBatchSendEmployee) TableName() string {
	return "mc_contact_message_batch_send_employee"
}

type ContactMessageBatchSendResult struct {
	ID             uint       `gorm:"primaryKey" json:"id"`
	BatchID        uint       `gorm:"column:batch_id;not null;index" json:"batchId"`
	EmployeeID     uint       `gorm:"column:employee_id;default:0" json:"employeeId"`
	ContactID      uint       `gorm:"column:contact_id;default:0" json:"contactId"`
	ExternalUserID string     `gorm:"column:external_user_id;size:100;default:''" json:"externalUserId"`
	UserID         string     `gorm:"column:user_id;size:100;default:''" json:"userId"`
	Status         int        `gorm:"default:0" json:"status"`
	SendTime       *time.Time `gorm:"column:send_time" json:"sendTime"`
	CreatedAt      time.Time  `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt      time.Time  `gorm:"column:updated_at" json:"updatedAt"`
}

func (ContactMessageBatchSendResult) TableName() string {
	return "mc_contact_message_batch_send_result"
}

// Room Message Batch Send Plugin
type RoomMessageBatchSend struct {
	ID                uint       `gorm:"primaryKey" json:"id"`
	CorpID            uint       `gorm:"column:corp_id;default:0;index" json:"corpId"`
	UserID            uint       `gorm:"column:user_id;default:0" json:"userId"`
	UserName          string     `gorm:"column:user_name;size:50;default:''" json:"userName"`
	EmployeeIDs       string     `gorm:"column:employee_ids;size:2000;default:''" json:"employeeIds"`
	BatchTitle        string     `gorm:"column:batch_title;size:100;default:''" json:"batchTitle"`
	Content           string     `gorm:"size:2000;default:''" json:"content"`
	SendWay           int        `gorm:"column:send_way;default:0" json:"sendWay"`
	DefiniteTime      *time.Time `gorm:"column:definite_time" json:"definiteTime"`
	SendTime          *time.Time `gorm:"column:send_time" json:"sendTime"`
	SendRoomTotal     int        `gorm:"column:send_room_total;default:0" json:"sendRoomTotal"`
	SendEmployeeTotal int        `gorm:"column:send_employee_total;default:0" json:"sendEmployeeTotal"`
	SendTotal         int        `gorm:"column:send_total;default:0" json:"sendTotal"`
	NotSendTotal      int        `gorm:"column:not_send_total;default:0" json:"notSendTotal"`
	ReceivedTotal     int        `gorm:"column:received_total;default:0" json:"receivedTotal"`
	NotReceivedTotal  int        `gorm:"column:not_received_total;default:0" json:"notReceivedTotal"`
	SendStatus        int        `gorm:"column:send_status;default:0" json:"sendStatus"`
	CreatedAt         time.Time  `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt         time.Time  `gorm:"column:updated_at" json:"updatedAt"`
	DeletedAt         *time.Time `gorm:"column:deleted_at;index" json:"deletedAt"`
}

func (RoomMessageBatchSend) TableName() string { return "mc_room_message_batch_send" }

type RoomMessageBatchSendEmployee struct {
	ID            uint       `gorm:"primaryKey" json:"id"`
	BatchID       uint       `gorm:"column:batch_id;not null;index" json:"batchId"`
	EmployeeID    uint       `gorm:"column:employee_id;default:0;index" json:"employeeId"`
	WxUserID      string     `gorm:"column:wx_user_id;size:100;default:''" json:"wxUserId"`
	SendRoomTotal int        `gorm:"column:send_room_total;default:0" json:"sendRoomTotal"`
	ErrCode       int        `gorm:"column:err_code;default:0" json:"errCode"`
	ErrMsg        string     `gorm:"column:err_msg;size:200;default:''" json:"errMsg"`
	MsgID         string     `gorm:"column:msg_id;size:100;default:''" json:"msgId"`
	SendTime      *time.Time `gorm:"column:send_time" json:"sendTime"`
	LastSyncTime  *time.Time `gorm:"column:last_sync_time" json:"lastSyncTime"`
	Status        int        `gorm:"default:0" json:"status"`
	ReceiveStatus int        `gorm:"column:receive_status;default:0" json:"receiveStatus"`
	CreatedAt     time.Time  `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt     time.Time  `gorm:"column:updated_at" json:"updatedAt"`
}

func (RoomMessageBatchSendEmployee) TableName() string {
	return "mc_room_message_batch_send_employee"
}

type RoomMessageBatchSendResult struct {
	ID              uint       `gorm:"primaryKey" json:"id"`
	BatchID         uint       `gorm:"column:batch_id;not null;index" json:"batchId"`
	EmployeeID      uint       `gorm:"column:employee_id;default:0" json:"employeeId"`
	RoomID          uint       `gorm:"column:room_id;default:0" json:"roomId"`
	RoomName        string     `gorm:"column:room_name;size:50;default:''" json:"roomName"`
	RoomEmployeeNum int        `gorm:"column:room_employee_num;default:0" json:"roomEmployeeNum"`
	RoomCreateTime  *time.Time `gorm:"column:room_create_time" json:"roomCreateTime"`
	ChatID          string     `gorm:"column:chat_id;size:100;default:''" json:"chatId"`
	UserID          string     `gorm:"column:user_id;size:100;default:''" json:"userId"`
	Status          int        `gorm:"default:0" json:"status"`
	SendTime        *time.Time `gorm:"column:send_time" json:"sendTime"`
	CreatedAt       time.Time  `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt       time.Time  `gorm:"column:updated_at" json:"updatedAt"`
}

func (RoomMessageBatchSendResult) TableName() string { return "mc_room_message_batch_send_result" }

// Contact Transfer Plugin
type WorkTransferLog struct {
	ID                 uint       `gorm:"primaryKey" json:"id"`
	CorpID             uint       `gorm:"column:corp_id;default:0;index" json:"corpId"`
	Status             int        `gorm:"default:0" json:"status"`
	Type               int        `gorm:"default:0" json:"type"`
	Name               string     `gorm:"size:50;default:''" json:"name"`
	ContactID          uint       `gorm:"column:contact_id;default:0" json:"contactId"`
	HandoverEmployeeID uint       `gorm:"column:handover_employee_id;default:0" json:"handoverEmployeeId"`
	TakeoverEmployeeID uint       `gorm:"column:takeover_employee_id;default:0" json:"takeoverEmployeeId"`
	State              string     `gorm:"size:50;default:''" json:"state"`
	CreatedAt          time.Time  `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt          time.Time  `gorm:"column:updated_at" json:"updatedAt"`
	DeletedAt          *time.Time `gorm:"column:deleted_at;index" json:"deletedAt"`
}

func (WorkTransferLog) TableName() string { return "mc_work_transfer_log" }

type WorkUnassigned struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	CorpID         uint      `gorm:"column:corp_id;default:0;index" json:"corpId"`
	HandoverUserid string    `gorm:"column:handover_userid;size:100;default:''" json:"handoverUserid"`
	ExternalUserid string    `gorm:"column:external_userid;size:100;default:''" json:"externalUserid"`
	DimissionTime  time.Time `gorm:"column:dimission_time" json:"dimissionTime"`
	CreatedAt      time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt      time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

func (WorkUnassigned) TableName() string { return "mc_work_unassigned" }

// Room Auto Pull Plugin
type WorkRoomAutoPull struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	CorpID       uint       `gorm:"column:corp_id;default:0;index" json:"corpId"`
	QrcodeName   string     `gorm:"column:qrcode_name;size:50;default:''" json:"qrcodeName"`
	QrcodeURL    string     `gorm:"column:qrcode_url;size:500;default:''" json:"qrcodeUrl"`
	WxConfigID   string     `gorm:"column:wx_config_id;size:100;default:''" json:"wxConfigId"`
	IsVerified   int        `gorm:"column:is_verified;default:0" json:"isVerified"`
	LeadingWords string     `gorm:"column:leading_words;size:500;default:''" json:"leadingWords"`
	Tags         string     `gorm:"size:2000;default:''" json:"tags"`
	Employees    string     `gorm:"size:2000;default:''" json:"employees"`
	Rooms        string     `gorm:"size:2000;default:''" json:"rooms"`
	CreatedAt    time.Time  `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt    time.Time  `gorm:"column:updated_at" json:"updatedAt"`
	DeletedAt    *time.Time `gorm:"column:deleted_at;index" json:"deletedAt"`
}

func (WorkRoomAutoPull) TableName() string { return "mc_work_room_auto_pull" }

// Room Tag Pull Plugin
type RoomTagPull struct {
	ID            uint       `gorm:"primaryKey" json:"id"`
	Name          string     `gorm:"size:50;default:''" json:"name"`
	Employees     string     `gorm:"size:2000;default:''" json:"employees"`
	ChooseContact string     `gorm:"column:choose_contact;size:2000;default:''" json:"chooseContact"`
	Guide         string     `gorm:"size:2000;default:''" json:"guide"`
	Rooms         string     `gorm:"size:2000;default:''" json:"rooms"`
	FilterContact string     `gorm:"column:filter_contact;size:2000;default:''" json:"filterContact"`
	ContactNum    int        `gorm:"column:contact_num;default:0" json:"contactNum"`
	WxTid         string     `gorm:"column:wx_tid;size:100;default:''" json:"wxTid"`
	TenantID      uint       `gorm:"column:tenant_id;default:0" json:"tenantId"`
	CorpID        uint       `gorm:"column:corp_id;default:0;index" json:"corpId"`
	CreateUserID  uint       `gorm:"column:create_user_id;default:0" json:"createUserId"`
	CreatedAt     time.Time  `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt     time.Time  `gorm:"column:updated_at" json:"updatedAt"`
	DeletedAt     *time.Time `gorm:"column:deleted_at;index" json:"deletedAt"`
}

func (RoomTagPull) TableName() string { return "mc_room_tag_pull" }

type RoomTagPullContact struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	RoomTagPullID    uint      `gorm:"column:room_tag_pull_id;default:0;index" json:"roomTagPullId"`
	ContactID        uint      `gorm:"column:contact_id;default:0" json:"contactId"`
	WxExternalUserid string    `gorm:"column:wx_external_userid;size:100;default:''" json:"wxExternalUserid"`
	ContactName      string    `gorm:"column:contact_name;size:50;default:''" json:"contactName"`
	EmployeeID       uint      `gorm:"column:employee_id;default:0" json:"employeeId"`
	WxUserID         string    `gorm:"column:wx_user_id;size:100;default:''" json:"wxUserId"`
	SendStatus       int       `gorm:"column:send_status;default:0" json:"sendStatus"`
	IsJoinRoom       int       `gorm:"column:is_join_room;default:0" json:"isJoinRoom"`
	RoomID           uint      `gorm:"column:room_id;default:0" json:"roomId"`
	CreatedAt        time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt        time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

func (RoomTagPullContact) TableName() string { return "mc_room_tag_pull_contact" }

// Room Welcome Plugin
type RoomWelcomeTemplate struct {
	ID                uint       `gorm:"primaryKey" json:"id"`
	CorpID            uint       `gorm:"column:corp_id;default:0;index" json:"corpId"`
	MsgText           string     `gorm:"column:msg_text;size:2000;default:''" json:"msgText"`
	ComplexType       int        `gorm:"column:complex_type;default:0" json:"complexType"`
	MsgComplex        string     `gorm:"column:msg_complex;size:2000;default:''" json:"msgComplex"`
	ComplexTemplateID string     `gorm:"column:complex_template_id;size:100;default:''" json:"complexTemplateId"`
	CreateUserID      uint       `gorm:"column:create_user_id;default:0" json:"createUserId"`
	CreatedAt         time.Time  `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt         time.Time  `gorm:"column:updated_at" json:"updatedAt"`
	DeletedAt         *time.Time `gorm:"column:deleted_at;index" json:"deletedAt"`
}

func (RoomWelcomeTemplate) TableName() string { return "mc_room_welcome_template" }

// Work Fission Plugin
type WorkFission struct {
	ID                    uint       `gorm:"primaryKey" json:"id"`
	CorpID                uint       `gorm:"column:corp_id;default:0;index" json:"corpId"`
	ActiveName            string     `gorm:"column:active_name;size:100;default:''" json:"activeName"`
	ServiceEmployees      string     `gorm:"column:service_employees;size:2000;default:''" json:"serviceEmployees"`
	AutoPass              int        `gorm:"column:auto_pass;default:0" json:"autoPass"`
	AutoAddTag            int        `gorm:"column:auto_add_tag;default:0" json:"autoAddTag"`
	ContactTags           string     `gorm:"column:contact_tags;size:2000;default:''" json:"contactTags"`
	EndTime               *time.Time `gorm:"column:end_time" json:"endTime"`
	QrCodeInvalid         int        `gorm:"column:qr_code_invalid;default:0" json:"qrCodeInvalid"`
	Tasks                 string     `gorm:"size:2000;default:''" json:"tasks"`
	NewFriend             int        `gorm:"column:new_friend;default:0" json:"newFriend"`
	DeleteInvalid         int        `gorm:"column:delete_invalid;default:0" json:"deleteInvalid"`
	ReceivePrize          int        `gorm:"column:receive_prize;default:0" json:"receivePrize"`
	ReceivePrizeEmployees string     `gorm:"column:receive_prize_employees;size:2000;default:''" json:"receivePrizeEmployees"`
	ReceiveLinks          string     `gorm:"column:receive_links;size:2000;default:''" json:"receiveLinks"`
	ReceiveQrcode         string     `gorm:"column:receive_qrcode;size:500;default:''" json:"receiveQrcode"`
	CreateUserID          uint       `gorm:"column:create_user_id;default:0" json:"createUserId"`
	CreatedAt             time.Time  `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt             time.Time  `gorm:"column:updated_at" json:"updatedAt"`
	DeletedAt             *time.Time `gorm:"column:deleted_at;index" json:"deletedAt"`
}

func (WorkFission) TableName() string { return "mc_work_fission" }

type WorkFissionContact struct {
	ID                        uint      `gorm:"primaryKey" json:"id"`
	FissionID                 uint      `gorm:"column:fission_id;default:0;index" json:"fissionId"`
	UnionID                   string    `gorm:"column:union_id;size:100;default:''" json:"unionId"`
	Nickname                  string    `gorm:"size:50;default:''" json:"nickname"`
	Avatar                    string    `gorm:"size:500;default:''" json:"avatar"`
	ContactSuperiorUserParent uint      `gorm:"column:contact_superior_user_parent;default:0" json:"contactSuperiorUserParent"`
	Level                     int       `gorm:"default:0" json:"level"`
	Employee                  string    `gorm:"size:2000;default:''" json:"employee"`
	InviteCount               int       `gorm:"column:invite_count;default:0" json:"inviteCount"`
	Loss                      int       `gorm:"default:0" json:"loss"`
	Status                    int       `gorm:"default:0" json:"status"`
	ReceiveLevel              int       `gorm:"column:receive_level;default:0" json:"receiveLevel"`
	IsNew                     int       `gorm:"column:is_new;default:0" json:"isNew"`
	ExternalUserID            string    `gorm:"column:external_user_id;size:100;default:''" json:"externalUserId"`
	QrcodeID                  uint      `gorm:"column:qrcode_id;default:0" json:"qrcodeId"`
	QrcodeURL                 string    `gorm:"column:qrcode_url;size:500;default:''" json:"qrcodeUrl"`
	CreatedAt                 time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt                 time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

func (WorkFissionContact) TableName() string { return "mc_work_fission_contact" }

type WorkFissionInvite struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	FissionID uint      `gorm:"column:fission_id;default:0;index" json:"fissionId"`
	Type      int       `gorm:"default:0" json:"type"`
	Text      string    `gorm:"size:500;default:''" json:"text"`
	LinkTitle string    `gorm:"column:link_title;size:100;default:''" json:"linkTitle"`
	LinkDesc  string    `gorm:"column:link_desc;size:200;default:''" json:"linkDesc"`
	LinkPic   string    `gorm:"column:link_pic;size:500;default:''" json:"linkPic"`
	WxLinkPic string    `gorm:"column:wx_link_pic;size:500;default:''" json:"wxLinkPic"`
	CreatedAt time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

func (WorkFissionInvite) TableName() string { return "mc_work_fission_invite" }

type WorkFissionPoster struct {
	ID                uint      `gorm:"primaryKey" json:"id"`
	FissionID         uint      `gorm:"column:fission_id;default:0;index" json:"fissionId"`
	PosterType        int       `gorm:"column:poster_type;default:0" json:"posterType"`
	CoverPic          string    `gorm:"column:cover_pic;size:500;default:''" json:"coverPic"`
	WxCoverPic        string    `gorm:"column:wx_cover_pic;size:500;default:''" json:"wxCoverPic"`
	FowardText        string    `gorm:"column:foward_text;size:500;default:''" json:"fowardText"`
	AvatarShow        int       `gorm:"column:avatar_show;default:0" json:"avatarShow"`
	NicknameShow      int       `gorm:"column:nickname_show;default:0" json:"nicknameShow"`
	NicknameColor     string    `gorm:"column:nickname_color;size:20;default:''" json:"nicknameColor"`
	CardCorpImageName string    `gorm:"column:card_corp_image_name;size:50;default:''" json:"cardCorpImageName"`
	CardCorpName      string    `gorm:"column:card_corp_name;size:50;default:''" json:"cardCorpName"`
	CardCorpLogo      string    `gorm:"column:card_corp_logo;size:500;default:''" json:"cardCorpLogo"`
	QrcodeW           int       `gorm:"column:qrcode_w;default:0" json:"qrcodeW"`
	QrcodeH           int       `gorm:"column:qrcode_h;default:0" json:"qrcodeH"`
	QrcodeX           int       `gorm:"column:qrcode_x;default:0" json:"qrcodeX"`
	QrcodeY           int       `gorm:"column:qrcode_y;default:0" json:"qrcodeY"`
	QrcodeID          uint      `gorm:"column:qrcode_id;default:0" json:"qrcodeId"`
	QrcodeURL         string    `gorm:"column:qrcode_url;size:500;default:''" json:"qrcodeUrl"`
	CreatedAt         time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt         time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

func (WorkFissionPoster) TableName() string { return "mc_work_fission_poster" }

type WorkFissionPush struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	FissionID      uint      `gorm:"column:fission_id;default:0;index" json:"fissionId"`
	PushEmployee   int       `gorm:"column:push_employee;default:0" json:"pushEmployee"`
	PushContact    int       `gorm:"column:push_contact;default:0" json:"pushContact"`
	MsgText        string    `gorm:"column:msg_text;size:500;default:''" json:"msgText"`
	MsgComplex     string    `gorm:"column:msg_complex;size:2000;default:''" json:"msgComplex"`
	MsgComplexType int       `gorm:"column:msg_complex_type;default:0" json:"msgComplexType"`
	CreatedAt      time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt      time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

func (WorkFissionPush) TableName() string { return "mc_work_fission_push" }

type WorkFissionWelcome struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	FissionID    uint      `gorm:"column:fission_id;default:0;index" json:"fissionId"`
	MsgText      string    `gorm:"column:msg_text;size:500;default:''" json:"msgText"`
	LinkTitle    string    `gorm:"column:link_title;size:100;default:''" json:"linkTitle"`
	LinkDesc     string    `gorm:"column:link_desc;size:200;default:''" json:"linkDesc"`
	LinkCoverURL string    `gorm:"column:link_cover_url;size:500;default:''" json:"linkCoverUrl"`
	LinkWxURL    string    `gorm:"column:link_wx_url;size:500;default:''" json:"linkWxUrl"`
	CreatedAt    time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt    time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

func (WorkFissionWelcome) TableName() string { return "mc_work_fission_welcome" }
