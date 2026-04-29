package main

import (
	"fmt"
	"school-go/internal/repository"
	"school-go/internal/service"
)

func main() {
	// 初始化数据库
	err := repository.InitDB()
	if err != nil {
		fmt.Println("数据库初始化失败:", err)
		return
	}

	// 创作用户服务
	userService := service.NewUserService()

	// 测试获取用户
	user, err := userService.GetByUsername("admin")
	if err != nil {
		fmt.Println("获取用户失败:", err)
		return
	}

	if user == nil {
		fmt.Println("用户不存在")
		return
	}

	fmt.Println("用户信息:")
	fmt.Println("用户名:", user.Username)
	fmt.Println("密码:", user.Password)
	fmt.Println("角色:", user.Role)

	// 测试认证
	password := "admin123"
	authenticated := userService.Authenticate("admin", password)
	fmt.Println("认证结果:", authenticated)
}
