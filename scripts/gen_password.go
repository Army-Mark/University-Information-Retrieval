package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// 生成密码的 bcrypt 哈希值
	password := "admin123"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("生成哈希值失败:", err)
		return
	}

	fmt.Println("密码:", password)
	fmt.Println("哈希值:", string(hashedPassword))
}
