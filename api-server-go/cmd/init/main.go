package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/bcrypt"
)

var rootCmd = &cobra.Command{
	Use:   "init",
	Short: "MoChat API Server 初始化工具",
	Long:  "用于初始化 MoChat API Server 配置和数据库",
	Run: func(cmd *cobra.Command, args []string) {
		// 实现初始化逻辑
		fmt.Println("开始初始化 MoChat API Server...")

		// 创建 .env 文件
		createEnvFile()

		// 设置域名
		setDomain()

		// 初始化数据库
		mysqlInit()

		// 执行 SQL 文件
		mysqlDataInit()

		// 初始化 Redis
		redisInit()

		// 注册管理员
		registerAdmin()

		// 初始化租户
		initTenant()

		fmt.Println("\n项目初始化完成！")
		fmt.Println("请使用设置的管理员账号登录系统")
	},
}

// 读取用户输入
func readInput(prompt string, defaultVal string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s (默认: %s): ", prompt, defaultVal)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input == "" {
		return defaultVal
	}
	return input
}

// 读取密码
func readPassword(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s: ", prompt)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

// 创建 .env 文件
func createEnvFile() {
	envPath := ".env"
	if _, err := os.Stat(envPath); os.IsNotExist(err) {
		examplePath := "config/config.example.yaml"
		if _, err := os.Stat(examplePath); err == nil {
			exampleContent, err := ioutil.ReadFile(examplePath)
			if err == nil {
				err = ioutil.WriteFile(envPath, exampleContent, 0644)
				if err == nil {
					fmt.Println("已创建 .env 文件")
				}
			}
		}
	}
}

// 设置域名
func setDomain() {
	dashboardBaseUrl := readInput("输入商户后台地址", "http://scrm.mo.chat")
	apiBaseUrl := readInput("输入后端接口地址", "http://api.mo.chat")
	sidebarBaseUrl := readInput("输入聊天侧边栏地址", "http://sidebar.mo.chat")
	operationBaseUrl := readInput("输入运营工具地址", "http://op.mo.chat")

	// 更新配置文件
	updateConfigFile(map[string]string{
		"wechat.dashboard_base_url": dashboardBaseUrl,
		"wechat.api_base_url":       apiBaseUrl,
		"wechat.sidebar_base_url":   sidebarBaseUrl,
		"wechat.operation_base_url": operationBaseUrl,
	})
}

// 初始化数据库
func mysqlInit() {
	host := readInput("输入 MySQL 主机地址", "127.0.0.1")
	port := readInput("输入 MySQL 端口", "3306")
	database := readInput("输入 MySQL 数据库名", "mochat")
	username := readInput("输入 MySQL 用户名", "root")
	password := readPassword("输入 MySQL 密码")

	// 更新配置文件
	updateConfigFile(map[string]string{
		"db.host":     host,
		"db.port":     port,
		"db.database": database,
		"db.username": username,
		"db.password": password,
	})

	// 测试数据库连接
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		username, password, host, port, database)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Printf("数据库连接失败: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		fmt.Printf("数据库连接失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("数据库连接成功")
}

// 执行 SQL 文件
func mysqlDataInit() {
	sqlFile := "../storage/install/mochat.sql"
	if _, err := os.Stat(sqlFile); os.IsNotExist(err) {
		sqlFile = "storage/install/mochat.sql"
	}

	if _, err := os.Stat(sqlFile); os.IsNotExist(err) {
		fmt.Println("SQL 文件不存在，请手动执行 SQL 文件")
		return
	}

	// 读取配置
	host := getConfigValue("db.host", "127.0.0.1")
	port := getConfigValue("db.port", "3306")
	database := getConfigValue("db.database", "mochat")
	username := getConfigValue("db.username", "root")
	password := getConfigValue("db.password", "")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		username, password, host, port, database)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Printf("数据库连接失败: %v\n", err)
		return
	}
	defer db.Close()

	// 读取 SQL 文件
	sqlContent, err := ioutil.ReadFile(sqlFile)
	if err != nil {
		fmt.Printf("读取 SQL 文件失败: %v\n", err)
		return
	}

	// 执行 SQL 语句
	sqlStatements := strings.Split(string(sqlContent), ";")
	for _, stmt := range sqlStatements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}
		_, err := db.Exec(stmt)
		if err != nil {
			fmt.Printf("执行 SQL 语句失败: %v\n", err)
			// 继续执行其他语句
		}
	}

	fmt.Println("SQL 文件执行完成")
}

// 初始化 Redis
func redisInit() {
	host := readInput("输入 Redis 主机地址", "127.0.0.1")
	port := readInput("输入 Redis 端口", "6379")
	auth := readPassword("输入 Redis 密码")
	db := readInput("输入 Redis 数据库编号", "0")

	// 更新配置文件
	updateConfigFile(map[string]string{
		"redis.host": host,
		"redis.port": port,
		"redis.auth": auth,
		"redis.db":   db,
	})

	fmt.Println("Redis 配置完成")
}

// 注册管理员
func registerAdmin() {
	// 读取配置
	host := getConfigValue("db.host", "127.0.0.1")
	port := getConfigValue("db.port", "3306")
	database := getConfigValue("db.database", "mochat")
	username := getConfigValue("db.username", "root")
	password := getConfigValue("db.password", "")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		username, password, host, port, database)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Printf("数据库连接失败: %v\n", err)
		return
	}
	defer db.Close()

	// 检查是否已有管理员
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM mc_user WHERE isSuperAdmin = 1").Scan(&count)
	if err != nil {
		fmt.Printf("查询管理员失败: %v\n", err)
		return
	}

	if count > 0 {
		fmt.Println("已存在管理员账号")
		return
	}

	// 输入管理员信息
	phone := readInput("输入管理员手机号", "13800138000")
	adminPassword := readPassword("输入管理员密码")

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
	if err != nil {
		fmt.Printf("密码加密失败: %v\n", err)
		return
	}

	// 插入管理员账号
	_, err = db.Exec(
		"INSERT INTO mc_user (phone, password, status, isSuperAdmin, created_at, updated_at) VALUES (?, ?, 1, 1, NOW(), NOW())",
		phone, string(hashedPassword),
	)
	if err != nil {
		fmt.Printf("创建管理员失败: %v\n", err)
		return
	}

	fmt.Println("管理员账号创建成功")
}

// 初始化租户
func initTenant() {
	// 读取配置
	host := getConfigValue("db.host", "127.0.0.1")
	port := getConfigValue("db.port", "3306")
	database := getConfigValue("db.database", "mochat")
	username := getConfigValue("db.username", "root")
	password := getConfigValue("db.password", "")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		username, password, host, port, database)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Printf("数据库连接失败: %v\n", err)
		return
	}
	defer db.Close()

	// 输入服务器 IP
	ips := readInput("输入服务器 IP (多个 IP 用逗号分隔)", "127.0.0.1")

	// 插入租户信息
	_, err = db.Exec(
		"INSERT INTO mc_tenant (name, status, server_ips, created_at, updated_at) VALUES ('默认租户', 1, ?, NOW(), NOW())",
		fmt.Sprintf("[%s]", strings.Replace(ips, ",", ", ", -1)),
	)
	if err != nil {
		fmt.Printf("创建租户失败: %v\n", err)
		return
	}

	fmt.Println("租户信息创建成功")
}

// 更新配置文件
func updateConfigFile(keyValues map[string]string) {
	configPath := ".env"
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		configPath = "config/config.example.yaml"
	}

	content, err := ioutil.ReadFile(configPath)
	if err != nil {
		fmt.Printf("读取配置文件失败: %v\n", err)
		return
	}

	contentStr := string(content)
	for key, value := range keyValues {
		// 简单的配置更新逻辑
		// 注意：这是一个简化版，实际项目中可能需要更复杂的 YAML 解析
		keyPath := strings.Replace(key, ".", "_", -1)
		oldLine := fmt.Sprintf("%s: %s", keyPath, getConfigValue(key, ""))
		newLine := fmt.Sprintf("%s: %s", keyPath, value)
		contentStr = strings.Replace(contentStr, oldLine, newLine, -1)
	}

	err = ioutil.WriteFile(configPath, []byte(contentStr), 0644)
	if err != nil {
		fmt.Printf("更新配置文件失败: %v\n", err)
	}
}

// 获取配置值
func getConfigValue(key, defaultValue string) string {
	configPath := ".env"
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		configPath = "config/config.example.yaml"
	}

	content, err := ioutil.ReadFile(configPath)
	if err != nil {
		return defaultValue
	}

	// 简单的配置读取逻辑
	// 注意：这是一个简化版，实际项目中可能需要更复杂的 YAML 解析
	keyPath := strings.Replace(key, ".", "_", -1)
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, keyPath+":") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				return strings.TrimSpace(parts[1])
			}
		}
	}

	return defaultValue
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
