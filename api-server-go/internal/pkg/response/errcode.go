package response

// 错误码定义，与原始 PHP 项目 AppErrCode 保持一致
var (
	ErrSuccess      = 0
	ErrServer       = 10000
	ErrParams       = 10001
	ErrAuth         = 10002
	ErrPermission   = 10003
	ErrNotFound     = 10004
	ErrTokenExpired = 10005
	ErrTokenInvalid = 10006
	ErrUserDisabled = 10007
	ErrPassword     = 10008
	ErrPhoneExists  = 10009
	ErrCorpExists   = 10010
	ErrWeChatAPI    = 10011
	ErrFileUpload   = 10012
	ErrDB           = 10013
	ErrRedis        = 10014
	ErrQueue        = 10015
)
