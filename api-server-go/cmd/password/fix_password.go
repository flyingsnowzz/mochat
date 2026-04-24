package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// 要加密的密码
	password := "123456"
	
	// 生成完整的密码哈希
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Printf("密码加密失败: %v\n", err)
		return
	}
	
	fmt.Println("原始密码:", password)
	fmt.Println("完整的加密后密码:", string(hashedPassword))
	fmt.Println("\n请将以下 SQL 语句执行到数据库中:")
	fmt.Println("UPDATE mc_user SET password = '" + string(hashedPassword) + "' WHERE phone = '13800138000';")
}
