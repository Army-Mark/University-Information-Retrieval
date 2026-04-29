package middleware

import (
	"github.com/gin-gonic/gin"
	"school-go/internal/service"
)

// 会话服务实例
var sessionService = service.NewSessionService()

// AuthRequired 认证中间件
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从cookie中获取会话ID
		sessionID, err := c.Cookie("session_id")
		
		// 如果cookie中没有，尝试从请求头中获取
		if err != nil || sessionID == "" {
			sessionID = c.Request.Header.Get("X-Session-ID")
		}
		
		if sessionID == "" {
			c.JSON(401, gin.H{"success": false, "message": "请先登录"})
			c.Abort()
			return
		}

		// 获取会话信息
		session, err := sessionService.GetSession(sessionID)
		if err != nil || session == nil {
			c.JSON(401, gin.H{"success": false, "message": "会话已过期，请重新登录"})
			c.Abort()
			return
		}

		// 获取用户信息，包括角色
		userService := service.NewUserService()
		user, err := userService.GetByUsername(session.Username)
		if err != nil || user == nil {
			c.JSON(401, gin.H{"success": false, "message": "用户不存在"})
			c.Abort()
			return
		}

		// 将用户名和角色设置到上下文中，方便后续处理
		c.Set("username", session.Username)
		c.Set("role", user.Role)
		c.Next()
	}
}

// RoleRequired 角色验证中间件
func RoleRequired(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从上下文中获取角色
		role, exists := c.Get("role")
		if !exists {
			c.JSON(401, gin.H{"success": false, "message": "请先登录"})
			c.Abort()
			return
		}

		// 检查角色是否在允许的角色列表中
		roleStr := role.(string)
		allowed := false
		for _, r := range roles {
			if r == roleStr {
				allowed = true
				break
			}
		}

		if !allowed {
			// 检查是否是 AJAX 请求
			if c.GetHeader("X-Requested-With") == "XMLHttpRequest" || c.GetHeader("Content-Type") == "application/json" {
				// 对于 AJAX 请求，返回 JSON 响应
				c.JSON(403, gin.H{"success": false, "message": "权限不足"})
			} else {
				// 对于普通请求，返回 HTML 页面
				c.HTML(403, "error.html", gin.H{
					"title": "权限不足",
					"message": "您没有足够的权限访问此页面",
				})
			}
			c.Abort()
			return
		}

		c.Next()
	}
}
