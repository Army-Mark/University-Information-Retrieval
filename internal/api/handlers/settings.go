package handlers

import (
	"net/http"
	"school-go/internal/service"

	"github.com/gin-gonic/gin"
)

// SettingsHandler 设置处理器
// 处理用户个性化设置相关的HTTP请求
type SettingsHandler struct {
	// settingsService 设置服务
	settingsService service.SettingsService
}

// NewSettingsHandler 创建设置处理器实例
// 返回 SettingsHandler 实例
func NewSettingsHandler() *SettingsHandler {
	return &SettingsHandler{
		settingsService: service.NewSettingsService(),
	}
}

// GetSettings 获取用户设置
// 处理 GET /api/settings 请求
func (h *SettingsHandler) GetSettings(c *gin.Context) {
	// 从上下文获取用户名
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "未登录",
		})
		return
	}

	// 获取用户设置
	settings, err := h.settingsService.GetSettings(username.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "获取设置失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"settings": settings,
	})
}

// SaveSettings 保存用户设置
// 处理 POST /api/settings 请求
// 请求体: {
//   "theme": "light",
//   "language": "zh-CN",
//   "default_view": "list",
//   "notifications": true
// }
func (h *SettingsHandler) SaveSettings(c *gin.Context) {
	// 从上下文获取用户名
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "未登录",
		})
		return
	}

	// 解析请求体
	var req struct {
		Theme         string `json:"theme"`
		Language      string `json:"language"`
		DefaultView   string `json:"default_view"`
		Notifications bool   `json:"notifications"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求数据错误: " + err.Error(),
		})
		return
	}

	// 获取当前设置
	settings, err := h.settingsService.GetSettings(username.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "获取设置失败: " + err.Error(),
		})
		return
	}

	// 更新设置
	settings.Theme = req.Theme
	settings.Language = req.Language
	settings.DefaultView = req.DefaultView
	settings.Notifications = req.Notifications

	// 保存设置
	if err := h.settingsService.SaveSettings(settings); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "保存设置失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "设置保存成功",
		"settings": settings,
	})
}

// AddToFavorites 添加到收藏
// 处理 POST /api/favorites 请求
// 请求体: {"school_id": "123"}
func (h *SettingsHandler) AddToFavorites(c *gin.Context) {
	// 从上下文获取用户名
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "未登录",
		})
		return
	}

	// 解析请求体
	var req struct {
		SchoolID string `json:"school_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求数据错误: " + err.Error(),
		})
		return
	}

	// 添加到收藏
	if err := h.settingsService.AddToFavorites(username.(string), req.SchoolID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "添加收藏失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "添加收藏成功",
	})
}

// RemoveFromFavorites 从收藏中移除
// 处理 DELETE /api/favorites/:id 请求
func (h *SettingsHandler) RemoveFromFavorites(c *gin.Context) {
	// 从上下文获取用户名
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "未登录",
		})
		return
	}

	// 获取学校ID
	schoolID := c.Param("id")

	// 从收藏中移除
	if err := h.settingsService.RemoveFromFavorites(username.(string), schoolID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "移除收藏失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "移除收藏成功",
	})
}

// GetFavorites 获取用户收藏
// 处理 GET /api/favorites 请求
func (h *SettingsHandler) GetFavorites(c *gin.Context) {
	// 从上下文获取用户名
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "未登录",
		})
		return
	}

	// 获取用户设置
	settings, err := h.settingsService.GetSettings(username.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "获取设置失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"favorites": settings.FavoriteSchools,
	})
}
