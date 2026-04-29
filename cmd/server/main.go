package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"school-go/internal/api"
	"school-go/internal/config"
	"school-go/internal/repository"
)

// main 应用程序入口函数
// 负责初始化配置、数据库连接、设置路由并启动服务器
func main() {
	// 获取应用根目录
	// 尝试获取可执行文件所在目录的绝对路径
	rootDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Printf("Warning: Failed to get root directory: %v", err)
		// 如果获取失败，使用当前工作目录
		rootDir, _ = os.Getwd()
	}
	log.Printf("Application root directory: %s", rootDir)

	// 加载应用配置
	// 从环境变量或配置文件中加载配置信息
	cfg := config.Load()

	// 设置应用根目录
	// 用于后续文件路径的解析
	repository.SetRootDir(rootDir)

	// 初始化数据库连接
	// 连接到 SQLite 数据库并创建必要的索引
	if err := repository.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// 设置 Gin 框架运行模式
	// 调试模式下会输出详细的日志信息
	if cfg.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// 创建 Gin 引擎实例
	// 默认包含日志和恢复中间件
	r := gin.Default()

	// 注册所有路由
	// 包括公开路由和需要认证的路由
	api.RegisterRoutes(r, rootDir)

	// 启动服务器
	// 从环境变量获取端口和主机配置，默认使用 0.0.0.0:5000
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	host := os.Getenv("HOST")
	if host == "" {
		host = "0.0.0.0"
	}

	log.Printf("Server starting on %s:%s...", host, port)
	// 启动 HTTP 服务器
	if err := r.Run(host + ":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
