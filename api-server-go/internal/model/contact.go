package model

import "time"

type WorkContact struct {
	ID               uint       `gorm:"primaryKey" json:"id"`
	CorpID           uint       `gorm:"column:corp_id;default:0;index" json:"corpId"`
	WxExternalUserID string     `gorm:"column:wx_external_userid;size:100;default:''" json:"wxExternalUserid"`
	Name             string     `gorm:"size:50;default:''" json:"name"`
	NickName         string     `gorm:"column:nick_name;size:50;default:''" json:"nickName"`
	Avatar           string     `gorm:"size:500;default:''" json:"avatar"`
	FollowUpStatus   int        `gorm:"column:follow_up_status;default:0" json:"followUpStatus"`
	Type             int        `gorm:"default:0" json:"type"`
	Gender           int        `gorm:"default:0" json:"gender"`
	Unionid          string     `gorm:"size:100;default:''" json:"unionid"`
	Position         string     `gorm:"size:50;default:''" json:"position"`
	CorpName         string     `gorm:"column:corp_name;size:100;default:''" json:"corpName"`
	CorpFullName     string     `gorm:"column:corp_full_name;size:100;default:''" json:"corpFullName"`
	ExternalProfile  string     `gorm:"column:external_profile;size:2000;default:''" json:"externalProfile"`
	BusinessNo       string     `gorm:"column:business_no;size:50;default:''" json:"businessNo"`
	CreatedAt        time.Time  `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt        time.Time  `gorm:"column:updated_at" json:"updatedAt"`
	DeletedAt        *time.Time `gorm:"column:deleted_at;index" json:"deletedAt"`
}

func (WorkContact) TableName() string { return "mc_work_contact" }

type WorkContactEmployee struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	EmployeeID     uint      `gorm:"column:employee_id;not null;index" json:"employeeId"`
	ContactID      uint      `gorm:"column:contact_id;not null;index" json:"contactId"`
	Remark         string    `gorm:"size:50;default:''" json:"remark"`
	Description    string    `gorm:"size:500;default:''" json:"description"`
	RemarkCorpName string    `gorm:"column:remark_corp_name;size:50;default:''" json:"remarkCorpName"`
	RemarkMobiles  string    `gorm:"column:remark_mobiles;size:200;default:''" json:"remarkMobiles"`
	AddWay         int       `gorm:"column:add_way;default:0" json:"addWay"`
	OperUserID     string    `gorm:"column:oper_userid;size:100;default:''" json:"operUserid"`
	State          string    `gorm:"size:100;default:''" json:"state"`
	CorpID         uint      `gorm:"column:corp_id;default:0;index" json:"corpId"`
	Status         int       `gorm:"default:0" json:"status"`
	CreateTime     time.Time `gorm:"column:create_time" json:"createTime"`
	CreatedAt      time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt      time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

func (WorkContactEmployee) TableName() string { return "mc_work_contact_employee" }

type ContactField struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	Name      string     `gorm:"size:50;not null" json:"name"`
	Label     string     `gorm:"size:50;default:''" json:"label"`
	Type      int        `gorm:"default:0" json:"type"`
	TypeText  string     `gorm:"-" json:"typeText"`
	Options   string     `gorm:"size:2000;default:''" json:"options"`
	Order     int        `gorm:"default:0" json:"order"`
	Status    int        `gorm:"default:1" json:"status"`
	IsSys     int        `gorm:"column:is_sys;default:0" json:"isSys"`
	CreatedAt time.Time  `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt time.Time  `gorm:"column:updated_at" json:"updatedAt"`
	DeletedAt *time.Time `gorm:"column:deleted_at;index" json:"deletedAt"`
}

func (ContactField) TableName() string { return "mc_contact_field" }

type ContactFieldPivot struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	ContactID      uint      `gorm:"column:contact_id;not null;index" json:"contactId"`
	ContactFieldID uint      `gorm:"column:contact_field_id;not null;index" json:"contactFieldId"`
	Value          string    `gorm:"size:500;default:''" json:"value"`
	CreatedAt      time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt      time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

func (ContactFieldPivot) TableName() string { return "mc_contact_field_pivot" }

type ContactEmployeeProcess struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	CorpID           uint      `gorm:"column:corp_id;default:0" json:"corpId"`
	EmployeeID       uint      `gorm:"column:employee_id;default:0" json:"employeeId"`
	ContactID        uint      `gorm:"column:contact_id;default:0" json:"contactId"`
	ContactProcessID uint      `gorm:"column:contact_process_id;default:0" json:"contactProcessId"`
	Content          string    `gorm:"size:2000;default:''" json:"content"`
	FileURL          string    `gorm:"column:file_url;size:500;default:''" json:"fileUrl"`
	CreatedAt        time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt        time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

func (ContactEmployeeProcess) TableName() string { return "mc_contact_employee_process" }

type ContactEmployeeTrack struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	EmployeeID uint      `gorm:"column:employee_id;default:0;index" json:"employeeId"`
	ContactID  uint      `gorm:"column:contact_id;default:0;index" json:"contactId"`
	Event      string    `gorm:"size:50;default:''" json:"event"`
	Content    string    `gorm:"size:2000;default:''" json:"content"`
	CorpID     uint      `gorm:"column:corp_id;default:0" json:"corpId"`
	CreatedAt  time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt  time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

func (ContactEmployeeTrack) TableName() string { return "mc_contact_employee_track" }

type ContactProcess struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	CorpID      uint       `gorm:"column:corp_id;default:0" json:"corpId"`
	Name        string     `gorm:"size:50;default:''" json:"name"`
	Description string     `gorm:"size:500;default:''" json:"description"`
	Order       int        `gorm:"default:0" json:"order"`
	CreatedAt   time.Time  `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt   time.Time  `gorm:"column:updated_at" json:"updatedAt"`
	DeletedAt   *time.Time `gorm:"column:deleted_at;index" json:"deletedAt"`
}

func (ContactProcess) TableName() string { return "mc_contact_process" }

type WorkUnionidExternalUseridMapping struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	CorpID         uint      `gorm:"column:corp_id;default:0" json:"corpId"`
	Unionid        string    `gorm:"size:100;default:''" json:"unionid"`
	Openid         string    `gorm:"size:100;default:''" json:"openid"`
	ExternalUserid string    `gorm:"column:external_userid;size:100;default:''" json:"externalUserid"`
	PendingID      string    `gorm:"column:pending_id;size:100;default:''" json:"pendingId"`
	SubjectType    int       `gorm:"column:subject_type;default:0" json:"subjectType"`
	CreatedAt      time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt      time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

func (WorkUnionidExternalUseridMapping) TableName() string {
	return "mc_work_unionid_external_userid_mapping"
}
