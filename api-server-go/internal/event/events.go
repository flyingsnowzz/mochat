package event

type BaseEvent struct {
	CorpID uint `json:"corpId"`
}

func (e BaseEvent) Name() string { return "base" }

type AddContactEvent struct {
	BaseEvent
	ExternalUserID string `json:"externalUserId"`
	UserID         string `json:"userId"`
	WelcomeCode    string `json:"welcomeCode"`
}

func (e AddContactEvent) Name() string { return "contact.add" }

type DeleteContactEvent struct {
	BaseEvent
	ExternalUserID string `json:"externalUserId"`
	UserID         string `json:"userId"`
}

func (e DeleteContactEvent) Name() string { return "contact.delete" }

type UpdateContactEvent struct {
	BaseEvent
	ExternalUserID string `json:"externalUserId"`
	UserID         string `json:"userId"`
}

func (e UpdateContactEvent) Name() string { return "contact.update" }

type DeleteFollowEmployeeEvent struct {
	BaseEvent
	ExternalUserID string `json:"externalUserId"`
	UserID         string `json:"userId"`
}

func (e DeleteFollowEmployeeEvent) Name() string { return "contact.follow_delete" }

type CreateTagEvent struct {
	BaseEvent
	TagID   string `json:"tagId"`
	TagName string `json:"tagName"`
}

func (e CreateTagEvent) Name() string { return "tag.create" }

type UpdateTagEvent struct {
	BaseEvent
	TagID   string `json:"tagId"`
	TagName string `json:"tagName"`
}

func (e UpdateTagEvent) Name() string { return "tag.update" }

type DeleteTagEvent struct {
	BaseEvent
	TagID string `json:"tagId"`
}

func (e DeleteTagEvent) Name() string { return "tag.delete" }

type CreateDepartmentEvent struct {
	BaseEvent
	DepartmentID int `json:"departmentId"`
}

func (e CreateDepartmentEvent) Name() string { return "department.create" }

type UpdateDepartmentEvent struct {
	BaseEvent
	DepartmentID int `json:"departmentId"`
}

func (e UpdateDepartmentEvent) Name() string { return "department.update" }

type DeleteDepartmentEvent struct {
	BaseEvent
	DepartmentID int `json:"departmentId"`
}

func (e DeleteDepartmentEvent) Name() string { return "department.delete" }

type CreateEmployeeEvent struct {
	BaseEvent
	UserID string `json:"userId"`
}

func (e CreateEmployeeEvent) Name() string { return "employee.create" }

type UpdateEmployeeEvent struct {
	BaseEvent
	UserID string `json:"userId"`
}

func (e UpdateEmployeeEvent) Name() string { return "employee.update" }

type DeleteEmployeeEvent struct {
	BaseEvent
	UserID string `json:"userId"`
}

func (e DeleteEmployeeEvent) Name() string { return "employee.delete" }

type CreateRoomEvent struct {
	BaseEvent
	ChatID string `json:"chatId"`
}

func (e CreateRoomEvent) Name() string { return "room.create" }

type UpdateRoomEvent struct {
	BaseEvent
	ChatID string `json:"chatId"`
}

func (e UpdateRoomEvent) Name() string { return "room.update" }

type DismissRoomEvent struct {
	BaseEvent
	ChatID string `json:"chatId"`
}

func (e DismissRoomEvent) Name() string { return "room.dismiss" }
