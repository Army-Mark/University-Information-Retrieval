package api

import (
	"net/http"
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
	logs, err := h.logService.GetAllLogs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "获取操作日志失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"logs":    logs,
	})
}
