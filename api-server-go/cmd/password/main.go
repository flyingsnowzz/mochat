package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// 要加密的密码
	password := "123456"
	
	// 使用与初始化代码相同的加密方法
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Printf("密码加密失败: %v\n", err)
		return
	}
	
	fmt.Println("原始密码:", password)
	fmt.Println("加密后密码:", string(hashedPassword))
	fmt.Println("\n请将以下 SQL 语句执行到数据库中:")
	fmt.Println("INSERT INTO mc_user (phone, password, status, isSuperAdmin, created_at, updated_at) VALUES ('13800138000', '" + string(hashedPassword) + "', 1, 1, NOW(), NOW())")
}
