package errors

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// ErrorMiddleware 错误处理中间件
func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// 处理错误
		if len(c.Errors) > 0 {
			for _, e := range c.Errors {
				// 记录错误日志
				log.Printf("Error: %v", e.Err)

				// 处理 AppError 类型的错误
				if appErr, ok := e.Err.(*AppError); ok {
					c.JSON(appErr.Code, gin.H{
						"success": false,
						"message": appErr.Message,
					})
					return
				}

				// 处理 TypedError 类型的错误
				if typedErr, ok := e.Err.(*TypedError); ok {
					c.JSON(typedErr.Code, gin.H{
						"success": false,
						"message": typedErr.Message,
						"type":    typedErr.Type,
					})
					return
				}

				// 处理其他类型的错误
				c.JSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"message": "内部服务器错误",
				})
				return
			}
		}
	}
}
