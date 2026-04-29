package api

import (
	"github.com/gin-gonic/gin"
	"school-go/internal/api/handlers"
	"school-go/internal/api/middleware"
	"school-go/internal/pkg/errors"
	"time"
)

// RegisterRoutes 注册所有路由
// r: Gin 引擎实例
// rootDir: 应用根目录
func RegisterRoutes(r *gin.Engine, rootDir string) {
	// 静态文件
	// 注册 /logo 路径，指向 static/logo 目录
	r.Static("/logo", "static/logo")
	// 注册 /static 路径，指向 static 目录
	r.Static("/static", "static")

	// 为静态资源添加缓存控制中间件
	r.Use(func(c *gin.Context) {
		// 检查是否是静态资源请求
		if len(c.Request.URL.Path) > 7 && (c.Request.URL.Path[:7] == "/static" || c.Request.URL.Path[:5] == "/logo") {
			// 设置缓存控制头，缓存1天
			c.Header("Cache-Control", "public, max-age=86400")
			// 设置过期时间，1天后
			c.Header("Expires", time.Now().Add(24*time.Hour).Format(time.RFC1123))
		}
		c.Next()
	})

	// 模板渲染
	// 加载 templates 目录下的所有 HTML 模板
	r.LoadHTMLGlob("templates/*.html")

	// 注册错误处理中间件
	r.Use(errors.ErrorMiddleware())

	// 处理器
	// 创建学校处理器实例
	universityHandler := handlers.NewUniversityHandler()
	// 创建用户处理器实例
	userHandler := handlers.NewUserHandler()
	// 创建设置处理器实例
	settingsHandler := handlers.NewSettingsHandler()
	// 创建操作日志处理器实例
	operationLogHandler := handlers.NewOperationLogHandler()

	// 公开路由
	// 首页
	r.GET("/", handlers.Index)
	// 搜索学校
	r.GET("/search", universityHandler.Search)
	// 学校详情
	r.GET("/university/:id", universityHandler.GetUniversity)
	// 获取滚动数据
	r.GET("/api/scrolling_data", universityHandler.GetScrollingData)
	// 更新滚动位置
	r.POST("/api/scrolling_position", universityHandler.UpdateScrollingPosition)

	// 认证相关
	// 用户登录
	r.POST("/login", userHandler.Login)
	// 用户登出
	r.POST("/logout", userHandler.Logout)
	// 检查登录状态
	r.GET("/check_login", userHandler.CheckLogin)

	// 需要认证的路由
	// 创建需要认证的路由组
	authGroup := r.Group("/")
	// 使用认证中间件
	authGroup.Use(middleware.AuthRequired())
	{
		// 用户账户管理
		// 普通用户和管理员都可以访问账户管理页面
		authGroup.GET("/account", userHandler.Account)
		authGroup.GET("/get_accounts", userHandler.GetAccounts)
		// 只有管理员可以添加和删除账户
			authGroup.POST("/add_account", middleware.RoleRequired("admin"), userHandler.AddAccount)
			// 所有登录用户都可以更新自己的账户
			authGroup.POST("/update_account", userHandler.UpdateAccount)
			authGroup.POST("/delete_account", middleware.RoleRequired("admin"), userHandler.DeleteAccount)

		// 学校管理
		// 普通用户和管理员都可以添加和编辑院校信息
		authGroup.GET("/add_school", universityHandler.AddSchoolPage)
		authGroup.POST("/upload_logo", universityHandler.UploadLogo)
		authGroup.POST("/add_school", universityHandler.AddSchool)
		authGroup.GET("/edit/:id", universityHandler.EditUniversity)
		authGroup.POST("/save", universityHandler.SaveUniversity)
		// 仅管理员可删除院校
		authGroup.POST("/delete_school", middleware.RoleRequired("admin"), universityHandler.DeleteSchool)

		// 操作日志 - 所有登录用户可访问
			authGroup.GET("/operation_logs", operationLogHandler.OperationLogsPage)
			authGroup.GET("/get_operation_logs", operationLogHandler.GetOperationLogs)
			// 只有管理员可以删除操作日志
			authGroup.POST("/delete_operation_logs", middleware.RoleRequired("admin"), operationLogHandler.DeleteOperationLogs)

		// 用户设置 - 所有登录用户可访问
		authGroup.GET("/api/settings", settingsHandler.GetSettings)
		authGroup.POST("/api/settings", settingsHandler.SaveSettings)
		authGroup.GET("/api/favorites", settingsHandler.GetFavorites)
		authGroup.POST("/api/favorites", settingsHandler.AddToFavorites)
		authGroup.DELETE("/api/favorites/:id", settingsHandler.RemoveFromFavorites)
	}
}
