package handlers

import (
	"fmt"
	"net/http"
	"school-go/internal/models"
	"school-go/internal/service"

	"github.com/gin-gonic/gin"
)

// OperationLogHandler 操作日志处理器
type OperationLogHandler struct {
	logService service.OperationLogService
}

// NewOperationLogHandler 创建操作日志处理器实例
func NewOperationLogHandler() *OperationLogHandler {
	return &OperationLogHandler{
		logService: service.NewOperationLogService(),
	}
}

// OperationLogsPage 操作日志页面
func (h *OperationLogHandler) OperationLogsPage(c *gin.Context) {
	c.HTML(http.StatusOK, "operation_logs.html", nil)
}

// GetOperationLogs 获取操作日志
func (h *OperationLogHandler) GetOperationLogs(c *gin.Context) {
	// 获取当前用户信息
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "未登录",
		})
		return
	}

	role, exists := c.Get("role")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "用户信息错误",
		})
		return
	}

	var logs []models.OperationLog
	var err error

	// 根据角色获取操作日志
	if role.(string) == "admin" {
		// 管理员可以查看所有操作日志
		logs, err = h.logService.GetAllLogs()
	} else {
		// 普通用户只能查看自己的操作日志
		logs, err = h.logService.GetLogsByUsername(username.(string))
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "获取操作日志失败",
		})
		return
	}

	// 确保即使没有日志，也返回一个空数组，而不是 null
	if logs == nil {
		logs = []models.OperationLog{}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"logs":    logs,
	})
}

// DeleteOperationLogs 删除操作日志
func (h *OperationLogHandler) DeleteOperationLogs(c *gin.Context) {
	// 获取当前用户角色
	role, exists := c.Get("role")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "用户信息错误",
		})
		return
	}

	// 只有管理员可以删除操作日志
	if role.(string) != "admin" {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "权限不足，只有管理员可以删除操作日志",
		})
		return
	}

	var req struct {
		LogIDs []int `json:"log_ids"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求数据错误",
		})
		return
	}

	if len(req.LogIDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请选择要删除的日志",
		})
		return
	}

	// 记录要删除的日志ID
	fmt.Printf("删除操作日志，ID列表: %v\n", req.LogIDs)

	err := h.logService.DeleteLogs(req.LogIDs)
	if err != nil {
		fmt.Printf("删除操作日志失败: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "删除操作日志失败",
		})
		return
	}

	fmt.Printf("成功删除 %d 条操作日志\n", len(req.LogIDs))
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "删除成功",
	})
}
