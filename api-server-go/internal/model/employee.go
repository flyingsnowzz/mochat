package model

import "time"

type WorkEmployee struct {
	ID                 uint       `gorm:"primaryKey" json:"id"`
	WxUserID           string     `gorm:"column:wx_user_id;size:100;default:''" json:"wxUserId"`
	CorpID             uint       `gorm:"column:corp_id;default:0;index" json:"corpId"`
	Name               string     `gorm:"size:50;default:''" json:"name"`
	Mobile             string     `gorm:"size:20;default:''" json:"mobile"`
	Position           string     `gorm:"size:50;default:''" json:"position"`
	Gender             int        `gorm:"default:0" json:"gender"`
	Email              string     `gorm:"size:100;default:''" json:"email"`
	Avatar             string     `gorm:"size:500;default:''" json:"avatar"`
	ThumbAvatar        string     `gorm:"column:thumb_avatar;size:500;default:''" json:"thumbAvatar"`
	Telephone          string     `gorm:"size:20;default:''" json:"telephone"`
	Alias              string     `gorm:"size:50;default:''" json:"alias"`
	Extattr            string     `gorm:"size:2000;default:''" json:"extattr"`
	Status             int        `gorm:"default:1" json:"status"`
	QrCode             string     `gorm:"column:qr_code;size:500;default:''" json:"qrCode"`
	ExternalProfile    string     `gorm:"column:external_profile;size:2000;default:''" json:"externalProfile"`
	ExternalPosition   string     `gorm:"column:external_position;size:50;default:''" json:"externalPosition"`
	Address            string     `gorm:"size:200;default:''" json:"address"`
	OpenUserID         string     `gorm:"column:open_user_id;size:100;default:''" json:"openUserId"`
	WxMainDepartmentID int        `gorm:"column:wx_main_department_id;default:0" json:"wxMainDepartmentId"`
	MainDepartmentID   uint       `gorm:"column:main_department_id;default:0" json:"mainDepartmentId"`
	LogUserID          uint       `gorm:"column:log_user_id;default:0" json:"logUserId"`
	ContactAuth        int        `gorm:"column:contact_auth;default:0" json:"contactAuth"`
	AuditStatus        int        `gorm:"column:audit_status;default:0" json:"auditStatus"`
	CreatedAt          time.Time  `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt          time.Time  `gorm:"column:updated_at" json:"updatedAt"`
	DeletedAt          *time.Time `gorm:"column:deleted_at;index" json:"deletedAt"`
}

func (WorkEmployee) TableName() string { return "mc_work_employee" }

type WorkDepartment struct {
	ID             uint       `gorm:"primaryKey" json:"id"`
	WxDepartmentID int        `gorm:"column:wx_department_id;default:0" json:"wxDepartmentId"`
	CorpID         uint       `gorm:"column:corp_id;default:0;index" json:"corpId"`
	Name           string     `gorm:"size:100;default:''" json:"name"`
	ParentID       uint       `gorm:"column:parent_id;default:0" json:"parentId"`
	WxParentID     int        `gorm:"column:wx_parentid;default:0" json:"wxParentId"`
	Order          int        `gorm:"default:0" json:"order"`
	Level          int        `gorm:"default:0" json:"level"`
	Path           string     `gorm:"size:500;default:''" json:"path"`
	CreatedAt      time.Time  `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt      time.Time  `gorm:"column:updated_at" json:"updatedAt"`
	DeletedAt      *time.Time `gorm:"column:deleted_at;index" json:"deletedAt"`
}

func (WorkDepartment) TableName() string { return "mc_work_department" }

type WorkEmployeeDepartment struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	EmployeeID     uint      `gorm:"column:employee_id;not null;index" json:"employeeId"`
	DepartmentID   uint      `gorm:"column:department_id;not null;index" json:"departmentId"`
	IsLeaderInDept int       `gorm:"column:is_leader_in_dept;default:0" json:"isLeaderInDept"`
	Order          int       `gorm:"default:0" json:"order"`
	CreatedAt      time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt      time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

func (WorkEmployeeDepartment) TableName() string { return "mc_work_employee_department" }

type WorkEmployeeTag struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	WxTagID   int       `gorm:"column:wx_tagid;default:0" json:"wxTagid"`
	CorpID    uint      `gorm:"column:corp_id;default:0;index" json:"corpId"`
	TagName   string    `gorm:"column:tag_name;size:50;default:''" json:"tagName"`
	CreatedAt time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

func (WorkEmployeeTag) TableName() string { return "mc_work_employee_tag" }

type WorkEmployeeTagPivot struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	EmployeeID uint      `gorm:"column:employee_id;not null;index" json:"employeeId"`
	TagID      uint      `gorm:"column:tag_id;not null;index" json:"tagId"`
	CreatedAt  time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt  time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

func (WorkEmployeeTagPivot) TableName() string { return "mc_work_employee_tag_pivot" }

type WorkEmployeeStatistic struct {
	ID                  uint    `gorm:"primaryKey" json:"id"`
	CorpID              uint    `gorm:"column:corp_id;default:0;index" json:"corpId"`
	EmployeeID          uint    `gorm:"column:employee_id;default:0;index" json:"employeeId"`
	NewApplyCnt         int     `gorm:"column:new_apply_cnt;default:0" json:"newApplyCnt"`
	NewContactCnt       int     `gorm:"column:new_contact_cnt;default:0" json:"newContactCnt"`
	ChatCnt             int     `gorm:"column:chat_cnt;default:0" json:"chatCnt"`
	MessageCnt          int     `gorm:"column:message_cnt;default:0" json:"messageCnt"`
	ReplyPercentage     float64 `gorm:"column:reply_percentage;default:0" json:"replyPercentage"`
	AvgReplyTime        int     `gorm:"column:avg_reply_time;default:0" json:"avgReplyTime"`
	NegativeFeedbackCnt int     `gorm:"column:negative_feedback_cnt;default:0" json:"negativeFeedbackCnt"`
	SynTime             string  `gorm:"column:syn_time;size:50;default:''" json:"synTime"`
	CreatedAt           time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt           time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

func (WorkEmployeeStatistic) TableName() string { return "mc_work_employee_statistic" }
