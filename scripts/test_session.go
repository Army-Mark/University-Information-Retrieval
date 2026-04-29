package main

import (
	"fmt"
	"school-go/internal/service"
)

func main() {
	// 创会话服务实例
	sessionService1 := service.NewSessionService()
	sessionService2 := service.NewSessionService()

	fmt.Println("sessionService1 == sessionService2:", sessionService1 == sessionService2)

	// 创建会话
	sessionID, err := sessionService1.CreateSession("admin")
	if err != nil {
		fmt.Println("创建会话失败:", err)
		return
	}

	fmt.Println("会话ID:", sessionID)

	// 从不同的会话服务实例获取会话
	session1, err := sessionService1.GetSession(sessionID)
	if err != nil {
		fmt.Println("获取会话1失败:", err)
		return
	}

	session2, err := sessionService2.GetSession(sessionID)
	if err != nil {
		fmt.Println("获取会话2失败:", err)
		return
	}

	fmt.Println("会话1:", session1)
	fmt.Println("会话2:", session2)
}
