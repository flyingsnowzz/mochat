package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

// Load 加载配置文件并返回 *Config 实例
// 配置加载优先级：环境变量 > config/config.yaml > 默认值
// 环境变量前缀为 MOCHAT_，使用下划线分隔嵌套层级，例如 MOCHAT_DB_HOST
func Load() (*Config, error) {
	v := viper.New()

	// 设置默认值
	setDefaults(v)

	// 配置文件路径
	configPath := os.Getenv("MOCHAT_CONFIG_PATH")
	if configPath == "" {
		configPath = "config"
	}

	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(configPath)

	// 支持环境变量覆盖，前缀 MOCHAT_
	v.SetEnvPrefix("MOCHAT")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// 读取配置文件
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("read config file error: %w", err)
		}
		// 配置文件不存在时不报错，允许完全通过环境变量配置
	}

	// 解析配置到结构体
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config error: %w", err)
	}

	return &cfg, nil
}

// setDefaults 设置配置默认值
func setDefaults(v *viper.Viper) {
	// App
	v.SetDefault("app.name", "mochat")
	v.SetDefault("app.env", "production")

	// DB
	v.SetDefault("db.driver", "mysql")
	v.SetDefault("db.host", "127.0.0.1")
	v.SetDefault("db.port", 3306)
	v.SetDefault("db.database", "mochat")
	v.SetDefault("db.username", "root")
	v.SetDefault("db.password", "")
	v.SetDefault("db.charset", "utf8mb4")
	v.SetDefault("db.collation", "utf8mb4_unicode_ci")
	v.SetDefault("db.prefix", "mc_")
	v.SetDefault("db.max_idle_time", 60)
	v.SetDefault("db.min_conns", 1)
	v.SetDefault("db.max_conns", 10)
	v.SetDefault("db.conn_timeout", 10)
	v.SetDefault("db.wait_timeout", 3)

	// Redis
	v.SetDefault("redis.host", "127.0.0.1")
	v.SetDefault("redis.auth", "")
	v.SetDefault("redis.port", 6379)
	v.SetDefault("redis.db", 0)
	v.SetDefault("redis.min_conns", 1)
	v.SetDefault("redis.max_conns", 10)
	v.SetDefault("redis.conn_timeout", 10)
	v.SetDefault("redis.wait_timeout", 3)
	v.SetDefault("redis.max_idle_time", 60)

	// JWT
	v.SetDefault("jwt.dashboard_secret", "3S6ybWbSy&23fFeq8")
	v.SetDefault("jwt.dashboard_prefix", "mc_jwt_")
	v.SetDefault("jwt.sidebar_secret", "Br3LXhp&Ysha1zRDh")

	// File
	v.SetDefault("file.driver", "local")

	// Server
	v.SetDefault("server.port", 9501)
	v.SetDefault("server.mode", "release")
	v.SetDefault("server.read_timeout", 30)
	v.SetDefault("server.write_timeout", 30)
	v.SetDefault("server.max_coroutine", 100000)

	// Log
	v.SetDefault("log.level", "info")
	v.SetDefault("log.output", "stdout")
}
