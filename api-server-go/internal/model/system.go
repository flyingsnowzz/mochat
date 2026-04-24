package model

import "time"

type BusinessLog struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	BusinessID  uint      `gorm:"column:business_id;default:0" json:"businessId"`
	Params      string    `gorm:"size:2000;default:''" json:"params"`
	Event       string    `gorm:"size:50;default:''" json:"event"`
	OperationID uint      `gorm:"column:operation_id;default:0" json:"operationId"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

func (BusinessLog) TableName() string { return "mc_business_log" }

type CorpDayData struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	CorpID         uint      `gorm:"column:corp_id;default:0;index" json:"corpId"`
	AddContactNum  int       `gorm:"column:add_contact_num;default:0" json:"addContactNum"`
	AddRoomNum     int       `gorm:"column:add_room_num;default:0" json:"addRoomNum"`
	AddIntoRoomNum int       `gorm:"column:add_into_room_num;default:0" json:"addIntoRoomNum"`
	LossContactNum int       `gorm:"column:loss_contact_num;default:0" json:"lossContactNum"`
	QuitRoomNum    int       `gorm:"column:quit_room_num;default:0" json:"quitRoomNum"`
	Date           string    `gorm:"size:20;default:''" json:"date"`
	CreatedAt      time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt      time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

func (CorpDayData) TableName() string { return "mc_corp_day_data" }

type OfficialAccount struct {
	ID                uint       `gorm:"primaryKey" json:"id"`
	AppType           int        `gorm:"column:app_type;default:0" json:"appType"`
	Appid             string     `gorm:"size:50;default:''" json:"appid"`
	AuthorizedStatus  int        `gorm:"column:authorized_status;default:0" json:"authorizedStatus"`
	AuthorizerAppid   string     `gorm:"column:authorizer_appid;size:50;default:''" json:"authorizerAppid"`
	AuthorizationCode string     `gorm:"column:authorization_code;size:200;default:''" json:"authorizationCode"`
	PreAuthCode       string     `gorm:"column:pre_auth_code;size:200;default:''" json:"preAuthCode"`
	HeadImg           string     `gorm:"column:head_img;size:500;default:''" json:"headImg"`
	Avatar            string     `gorm:"size:500;default:''" json:"avatar"`
	BusinessInfo      string     `gorm:"column:business_info;size:2000;default:''" json:"businessInfo"`
	Modules           string     `gorm:"size:2000;default:''" json:"modules"`
	Nickname          string     `gorm:"size:50;default:''" json:"nickname"`
	ServiceTypeInfo   string     `gorm:"column:service_type_info;size:200;default:''" json:"serviceTypeInfo"`
	VerifyTypeInfo    string     `gorm:"column:verify_type_info;size:200;default:''" json:"verifyTypeInfo"`
	OriginalID        string     `gorm:"column:original_id;size:50;default:''" json:"originalId"`
	FuncInfo          string     `gorm:"column:func_info;size:2000;default:''" json:"funcInfo"`
	PrincipalName     string     `gorm:"column:principal_name;size:100;default:''" json:"principalName"`
	Alias             string     `gorm:"size:50;default:''" json:"alias"`
	QrcodeURL         string     `gorm:"column:qrcode_url;size:500;default:''" json:"qrcodeUrl"`
	LocalQrcodeURL    string     `gorm:"column:local_qrcode_url;size:500;default:''" json:"localQrcodeUrl"`
	CallbackSuffix    string     `gorm:"column:callback_suffix;size:50;default:''" json:"callbackSuffix"`
	CallbackVerified  int        `gorm:"column:callback_verified;default:0" json:"callbackVerified"`
	UserName          string     `gorm:"column:user_name;size:50;default:''" json:"userName"`
	EncodingAesKey    string     `gorm:"column:encoding_aes_key;size:100;default:''" json:"encodingAesKey"`
	NotifyURL         string     `gorm:"column:notify_url;size:500;default:''" json:"notifyUrl"`
	Secret            string     `gorm:"size:100;default:''" json:"secret"`
	Token             string     `gorm:"size:100;default:''" json:"token"`
	TenantID          uint       `gorm:"column:tenant_id;default:0" json:"tenantId"`
	CorpID            uint       `gorm:"column:corp_id;default:0" json:"corpId"`
	CreateUserID      uint       `gorm:"column:create_user_id;default:0" json:"createUserId"`
	CreatedAt         time.Time  `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt         time.Time  `gorm:"column:updated_at" json:"updatedAt"`
	DeletedAt         *time.Time `gorm:"column:deleted_at;index" json:"deletedAt"`
}

func (OfficialAccount) TableName() string { return "mc_official_account" }

type OfficialAccountSet struct {
	ID                uint      `gorm:"primaryKey" json:"id"`
	OfficialAccountID uint      `gorm:"column:official_account_id;default:0" json:"officialAccountId"`
	Type              int       `gorm:"default:0" json:"type"`
	TenantID          uint      `gorm:"column:tenant_id;default:0" json:"tenantId"`
	CorpID            uint      `gorm:"column:corp_id;default:0" json:"corpId"`
	CreateUserID      uint      `gorm:"column:create_user_id;default:0" json:"createUserId"`
	CreatedAt         time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt         time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

func (OfficialAccountSet) TableName() string { return "mc_official_account_set" }

type SystemConfig struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	Name        string     `gorm:"size:50;default:''" json:"name"`
	Remark      string     `gorm:"size:500;default:''" json:"remark"`
	Description string     `gorm:"size:500;default:''" json:"description"`
	TenantID    uint       `gorm:"column:tenant_id;default:0" json:"tenantId"`
	CorpID      uint       `gorm:"column:corp_id;default:0" json:"corpId"`
	CreatorID   uint       `gorm:"column:creator_id;default:0" json:"creatorId"`
	CreatedAt   time.Time  `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt   time.Time  `gorm:"column:updated_at" json:"updatedAt"`
	DeletedAt   *time.Time `gorm:"column:deleted_at;index" json:"deletedAt"`
}

func (SystemConfig) TableName() string { return "mc_system_config" }

type SystemConfigValue struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Type        int       `gorm:"default:0" json:"type"`
	TargetID    uint      `gorm:"column:target_id;default:0" json:"targetId"`
	ConfigID    uint      `gorm:"column:config_id;default:0" json:"configId"`
	Value       string    `gorm:"size:2000;default:''" json:"value"`
	Description string    `gorm:"size:500;default:''" json:"description"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

func (SystemConfigValue) TableName() string { return "mc_system_config_value" }

type SysLog struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	URLPath     string    `gorm:"column:url_path;size:255;default:''" json:"urlPath"`
	Method      string    `gorm:"size:10;default:''" json:"method"`
	Query       string    `gorm:"size:2000;default:''" json:"query"`
	Body        string    `gorm:"size:2000;default:''" json:"body"`
	MenuID      uint      `gorm:"column:menu_id;default:0" json:"menuId"`
	MenuName    string    `gorm:"column:menu_name;size:50;default:''" json:"menuName"`
	OperateID   uint      `gorm:"column:operate_id;default:0" json:"operateId"`
	OperateName string    `gorm:"column:operate_name;size:50;default:''" json:"operateName"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

func (SysLog) TableName() string { return "mc_sys_log" }

type Plugin struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	CorpID    uint       `gorm:"column:corp_id;default:0" json:"corpId"`
	Name      string     `gorm:"size:50;default:''" json:"name"`
	Version   string     `gorm:"size:20;default:''" json:"version"`
	Content   string     `gorm:"size:2000;default:''" json:"content"`
	Status    int        `gorm:"default:1" json:"status"`
	CreatedAt time.Time  `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt time.Time  `gorm:"column:updated_at" json:"updatedAt"`
	DeletedAt *time.Time `gorm:"column:deleted_at;index" json:"deletedAt"`
}

func (Plugin) TableName() string { return "mc_plugin" }
