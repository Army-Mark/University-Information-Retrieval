package main

import (
	"fmt"
	"os"

	"school-go/internal/config"
	"school-go/internal/repository"
)

func main() {
	// 设置应用根目录
	rootDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("获取工作目录失败: %v\n", err)
		os.Exit(1)
	}
	repository.SetRootDir(rootDir)

	// 加载配置
	cfg := config.Load()

	fmt.Println("=== 初始化数据库 ===")
	fmt.Printf("数据库路径: %s\n", cfg.DBPath)

	// 初始化数据库
	if err := repository.InitDB(); err != nil {
		fmt.Printf("数据库初始化失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("数据库初始化成功")

	// 获取数据库连接
	db := repository.GetDB()

	// 检查用户表
	fmt.Println("\n=== 检查用户表 ===")
	var exists bool
	err = db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='users'").Scan(&exists)
	if err != nil {
		fmt.Printf("检查用户表失败: %v\n", err)
	} else {
		fmt.Printf("用户表存在: %v\n", exists)
	}

	// 检查管理员用户
	fmt.Println("\n=== 检查管理员用户 ===")
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM users WHERE username = 'admin'").Scan(&count)
	if err != nil {
		fmt.Printf("检查管理员用户失败: %v\n", err)
	} else {
		fmt.Printf("管理员用户数量: %d\n", count)
		if count > 0 {
			var username, role string
			err = db.QueryRow("SELECT username, role FROM users WHERE username = 'admin'").Scan(&username, &role)
			if err == nil {
				fmt.Printf("管理员用户名: %s, 角色: %s\n", username, role)
			}
		}
	}

	fmt.Println("\n=== 数据库初始化完成 ===")
}