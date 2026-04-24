package config

// Config 全局配置结构体，对应原始 PHP 项目的 .env.example 配置项
type Config struct {
	App    AppConfig
	DB     DBConfig
	Redis  RedisConfig
	JWT    JWTConfig
	File   FileConfig
	WeChat WeChatConfig
	Proxy  ProxyConfig
	Server ServerConfig
	Log    LogConfig
}

// AppConfig 应用基础配置
type AppConfig struct {
	Name string `mapstructure:"name" yaml:"name"`
	Env  string `mapstructure:"env" yaml:"env"` // dev, production
}

// DBConfig 数据库配置
type DBConfig struct {
	Driver      string  `mapstructure:"driver" yaml:"driver"`
	Host        string  `mapstructure:"host" yaml:"host"`
	Port        int     `mapstructure:"port" yaml:"port"`
	Database    string  `mapstructure:"database" yaml:"database"`
	Username    string  `mapstructure:"username" yaml:"username"`
	Password    string  `mapstructure:"password" yaml:"password"`
	Charset     string  `mapstructure:"charset" yaml:"charset"`
	Collation   string  `mapstructure:"collation" yaml:"collation"`
	Prefix      string  `mapstructure:"prefix" yaml:"prefix"` // mc_
	MaxIdleTime int     `mapstructure:"max_idle_time" yaml:"max_idle_time"`
	MinConns    int     `mapstructure:"min_conns" yaml:"min_conns"`
	MaxConns    int     `mapstructure:"max_conns" yaml:"max_conns"`
	ConnTimeout float64 `mapstructure:"conn_timeout" yaml:"conn_timeout"`
	WaitTimeout float64 `mapstructure:"wait_timeout" yaml:"wait_timeout"`
}

// RedisConfig Redis 配置
type RedisConfig struct {
	Host        string  `mapstructure:"host" yaml:"host"`
	Auth        string  `mapstructure:"auth" yaml:"auth"`
	Port        int     `mapstructure:"port" yaml:"port"`
	DB          int     `mapstructure:"db" yaml:"db"`
	MinConns    int     `mapstructure:"min_conns" yaml:"min_conns"`
	MaxConns    int     `mapstructure:"max_conns" yaml:"max_conns"`
	ConnTimeout float64 `mapstructure:"conn_timeout" yaml:"conn_timeout"`
	WaitTimeout float64 `mapstructure:"wait_timeout" yaml:"wait_timeout"`
	MaxIdleTime float64 `mapstructure:"max_idle_time" yaml:"max_idle_time"`
}

// JWTConfig JWT 认证配置
type JWTConfig struct {
	DashboardSecret string `mapstructure:"dashboard_secret" yaml:"dashboard_secret"`
	DashboardPrefix string `mapstructure:"dashboard_prefix" yaml:"dashboard_prefix"`
	SidebarSecret   string `mapstructure:"sidebar_secret" yaml:"sidebar_secret"`
}

// FileConfig 文件存储配置
type FileConfig struct {
	Driver string      `mapstructure:"driver" yaml:"driver"` // local, oss, cos, s3, minio, qiniu
	OSS    OSSConfig   `mapstructure:"oss" yaml:"oss"`
	COS    COSConfig   `mapstructure:"cos" yaml:"cos"`
	S3     S3Config    `mapstructure:"s3" yaml:"s3"`
	Qiniu  QiniuConfig `mapstructure:"qiniu" yaml:"qiniu"`
}

// OSSConfig 阿里云 OSS 配置
type OSSConfig struct {
	AccessID     string `mapstructure:"access_id" yaml:"access_id"`
	AccessSecret string `mapstructure:"access_secret" yaml:"access_secret"`
	Bucket       string `mapstructure:"bucket" yaml:"bucket"`
	Endpoint     string `mapstructure:"endpoint" yaml:"endpoint"`
}

// COSConfig 腾讯云 COS 配置
type COSConfig struct {
	AppID     string `mapstructure:"app_id" yaml:"app_id"`
	SecretID  string `mapstructure:"secret_id" yaml:"secret_id"`
	SecretKey string `mapstructure:"secret_key" yaml:"secret_key"`
	Bucket    string `mapstructure:"bucket" yaml:"bucket"`
	Region    string `mapstructure:"region" yaml:"region"`
}

// S3Config AWS S3 / MinIO 配置
type S3Config struct {
	AccessKey string `mapstructure:"access_key" yaml:"access_key"`
	SecretKey string `mapstructure:"secret_key" yaml:"secret_key"`
	Bucket    string `mapstructure:"bucket" yaml:"bucket"`
	Region    string `mapstructure:"region" yaml:"region"`
	Endpoint  string `mapstructure:"endpoint" yaml:"endpoint"`
}

// QiniuConfig 七牛云配置
type QiniuConfig struct {
	AccessKey string `mapstructure:"access_key" yaml:"access_key"`
	SecretKey string `mapstructure:"secret_key" yaml:"secret_key"`
	Bucket    string `mapstructure:"bucket" yaml:"bucket"`
	Domain    string `mapstructure:"domain" yaml:"domain"`
}

// WeChatConfig 微信相关配置
type WeChatConfig struct {
	OpenPlatform     OpenPlatformConfig `mapstructure:"open_platform" yaml:"open_platform"`
	APIBaseURL       string             `mapstructure:"api_base_url" yaml:"api_base_url"`
	DashboardBaseURL string             `mapstructure:"dashboard_base_url" yaml:"dashboard_base_url"`
	SidebarBaseURL   string             `mapstructure:"sidebar_base_url" yaml:"sidebar_base_url"`
	OperationBaseURL string             `mapstructure:"operation_base_url" yaml:"operation_base_url"`
}

// OpenPlatformConfig 微信开放平台第三方平台配置
type OpenPlatformConfig struct {
	AppID  string `mapstructure:"app_id" yaml:"app_id"`
	Secret string `mapstructure:"secret" yaml:"secret"`
	Token  string `mapstructure:"token" yaml:"token"`
	AESKey string `mapstructure:"aes_key" yaml:"aes_key"`
}

// ProxyConfig 代理配置
type ProxyConfig struct {
	HTTP  string `mapstructure:"http" yaml:"http"`
	HTTPS string `mapstructure:"https" yaml:"https"`
}

// ServerConfig HTTP 服务器配置
type ServerConfig struct {
	Port         int    `mapstructure:"port" yaml:"port"`
	Mode         string `mapstructure:"mode" yaml:"mode"` // debug, release
	ReadTimeout  int    `mapstructure:"read_timeout" yaml:"read_timeout"`
	WriteTimeout int    `mapstructure:"write_timeout" yaml:"write_timeout"`
	MaxCoroutine int    `mapstructure:"max_coroutine" yaml:"max_coroutine"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level  string `mapstructure:"level" yaml:"level"`
	Output string `mapstructure:"output" yaml:"output"`
}
