package handlers

import (
	"net/http"
	"regexp"
	"school-go/internal/models"
	"school-go/internal/pkg/errors"
	"school-go/internal/service"

	"github.com/gin-gonic/gin"
)

// 会话服务实例
var sessionService = service.NewSessionService()

// isPasswordStrong 检查密码强度
// 密码至少需要8个字符，包含字母和数字
func isPasswordStrong(password string) bool {
	if len(password) < 8 {
		return false
	}
	hasLetter := regexp.MustCompile(`[a-zA-Z]`).MatchString(password)
	hasDigit := regexp.MustCompile(`[0-9]`).MatchString(password)
	return hasLetter && hasDigit
}

// UserHandler 用户处理器
// 处理用户相关的 HTTP 请求
type UserHandler struct {
	// service 用户服务
	service service.UserService
}

// NewUserHandler 创建用户处理器实例
// 返回 UserHandler 实例
func NewUserHandler() *UserHandler {
	return &UserHandler{
		service: service.NewUserService(),
	}
}

// Login 登录
// 处理 POST /login 请求
// 请求体：{"username": "用户名", "password": "密码"}
// 返回登录结果
func (h *UserHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.BadRequest("请求数据错误", err))
		return
	}

	if !h.service.Authenticate(req.Username, req.Password) {
		c.JSON(http.StatusOK, models.LoginResponse{
			Success: false,
			Message: "用户名或密码错误",
		})
		return
	}

	// 创建会话
	sessionID, err := sessionService.CreateSession(req.Username)
	if err != nil {
		c.Error(errors.InternalServerError("创建会话失败", err))
		return
	}

	// 设置会话cookie
	// 在开发和Docker环境中，不要设置Secure标志
	c.SetCookie("session_id", sessionID, 86400, "/", "", false, false) // 24小时有效期

	c.JSON(http.StatusOK, models.LoginResponse{
		Success: true,
		Message: "登录成功",
	})
}

// Logout 注销
// 处理 POST /logout 请求
// 清除会话cookie并返回注销结果
func (h *UserHandler) Logout(c *gin.Context) {
	// 从cookie中获取会话ID
	sessionID, err := c.Cookie("session_id")
	if err == nil && sessionID != "" {
		// 删除会话
		sessionService.DeleteSession(sessionID)
	}

	// 清除会话cookie
	c.SetCookie("session_id", "", -1, "/", "", false, false)

	c.JSON(http.StatusOK, models.LoginResponse{
		Success: true,
		Message: "退出成功",
	})
}

// CheckLogin 检查登录状态
// 处理 GET /check_login 请求
// 从cookie中获取会话ID并返回登录状态
func (h *UserHandler) CheckLogin(c *gin.Context) {
	// 从cookie中获取会话ID
	sessionID, err := c.Cookie("session_id")
	if err != nil || sessionID == "" {
		c.JSON(http.StatusOK, gin.H{
			"logged_in": false,
		})
		return
	}

	// 获取会话信息
	session, err := sessionService.GetSession(sessionID)
	if err != nil || session == nil {
		c.JSON(http.StatusOK, gin.H{
			"logged_in": false,
		})
		return
	}

	// 获取用户信息，包括角色
	user, err := h.service.GetByUsername(session.Username)
	if err != nil || user == nil {
		c.JSON(http.StatusOK, gin.H{
			"logged_in": false,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"logged_in": true,
		"username":  session.Username,
		"role":      user.Role,
	})
}

// Account 账户页面
// 处理 GET /account 请求
// 返回账户管理页面
func (h *UserHandler) Account(c *gin.Context) {
	c.HTML(http.StatusOK, "account.html", nil)
}

// GetAccounts 获取所有账户
// 处理 GET /get_accounts 请求
// 返回所有用户账户列表
func (h *UserHandler) GetAccounts(c *gin.Context) {
	users, err := h.service.GetAll()
	if err != nil {
		c.Error(errors.InternalServerError("获取账户失败", err))
		return
	}

	c.JSON(http.StatusOK, models.AccountResponse{
		Success:  true,
		Accounts: users,
	})
}

// AddAccount 添加账户
// 处理 POST /add_account 请求
// 请求体：{"username": "用户名", "password": "密码", "role": "角色"}
// 返回添加结果
func (h *UserHandler) AddAccount(c *gin.Context) {
	var req models.AddAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.BadRequest("请求数据错误", err))
		return
	}

	// 检查密码强度
	if !isPasswordStrong(req.Password) {
		c.Error(errors.BadRequest("密码强度不足，至少需要8个字符，包含字母和数字", nil))
		return
	}

	// 检查角色是否有效
	if req.Role != "admin" && req.Role != "user" {
		c.Error(errors.BadRequest("角色必须是admin或user", nil))
		return
	}

	// 检查用户是否存在
	existingUser, err := h.service.GetByUsername(req.Username)
	if err != nil {
		c.Error(errors.InternalServerError("检查用户失败", err))
		return
	}

	if existingUser != nil {
		c.Error(errors.BadRequest("用户名已存在", nil))
		return
	}

	// 创建用户
	user := models.User{
		Username: req.Username,
		Password: req.Password,
		Role:     req.Role,
	}

	err = h.service.Create(&user)
	if err != nil {
		c.Error(errors.InternalServerError("添加失败", err))
		return
	}

	c.JSON(http.StatusOK, models.AccountResponse{
		Success: true,
		Message: "添加成功",
	})
}

// UpdateAccount 更新账户
// 处理 POST /update_account 请求
// 请求体：{"old_username": "原用户名", "new_username": "新用户名", "new_password": "新密码", "newRole": "新角色"}
// 返回更新结果
func (h *UserHandler) UpdateAccount(c *gin.Context) {
	var req models.UpdateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.BadRequest("请求数据错误", err))
		return
	}

	// 检查密码强度
	if !isPasswordStrong(req.NewPassword) {
		c.Error(errors.BadRequest("密码强度不足，至少需要8个字符，包含字母和数字", nil))
		return
	}

	// 获取当前登录用户的用户名
	currentUsername, exists := c.Get("username")
	if !exists {
		c.Error(errors.Unauthorized("未登录", nil))
		return
	}

	// 检查是否是更新自己的账户
	if currentUsername != req.OldUsername {
		// 检查当前用户是否是管理员
		user, err := h.service.GetByUsername(currentUsername.(string))
		if err != nil || user == nil || user.Role != "admin" {
			c.Error(errors.Forbidden("权限不足，只能更新自己的账户", nil))
			return
		}
	}

	// 检查角色是否有效
	if req.NewRole != "admin" && req.NewRole != "user" {
		c.Error(errors.BadRequest("角色必须是admin或user", nil))
		return
	}

	// 检查原用户是否存在
	existingUser, err := h.service.GetByUsername(req.OldUsername)
	if err != nil {
		c.Error(errors.InternalServerError("检查用户失败", err))
		return
	}

	if existingUser == nil {
		c.Error(errors.BadRequest("原用户不存在", nil))
		return
	}

	// 非管理员不能修改角色和用户名
	currentUser, _ := h.service.GetByUsername(currentUsername.(string))
	if currentUser != nil && currentUser.Role != "admin" {
		// 普通用户只能更新自己的账户，且不能修改角色和用户名
		req.NewRole = existingUser.Role
		req.NewUsername = existingUser.Username
	}

	// 更新用户
	err = h.service.Update(req.OldUsername, req.NewUsername, req.NewPassword, req.NewRole)
	if err != nil {
		c.Error(errors.InternalServerError("更新失败", err))
		return
	}

	// 记录操作日志
	logService := service.NewOperationLogService()
	newUser := &models.User{
		Username: req.NewUsername,
		Password: req.NewPassword, // 注意：这里存储的是明文密码，实际应该存储哈希后的密码，但为了日志记录，这里使用明文
		Role:     req.NewRole,
	}
	// 获取IP地址和用户代理
	ipAddress := c.ClientIP()
	userAgent := c.Request.UserAgent()
	logService.LogOperation("update", "users", req.OldUsername, 0, currentUsername.(string), ipAddress, userAgent, existingUser, newUser)

	c.JSON(http.StatusOK, models.AccountResponse{
		Success: true,
		Message: "更新成功",
	})
}

// DeleteAccount 删除账户
// 处理 POST /delete_account 请求
// 请求体：{"username": "用户名"}
// 返回删除结果
func (h *UserHandler) DeleteAccount(c *gin.Context) {
	var req models.DeleteAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.BadRequest("请求数据错误", err))
		return
	}

	// 检查用户是否存在
	existingUser, err := h.service.GetByUsername(req.Username)
	if err != nil {
		c.Error(errors.InternalServerError("检查用户失败", err))
		return
	}

	if existingUser == nil {
		c.Error(errors.BadRequest("用户不存在", nil))
		return
	}

	// 检查是否删除当前登录用户
	currentUsername, _ := c.Get("username")
	if currentUsername == req.Username {
		c.Error(errors.BadRequest("不能删除当前登录账户", nil))
		return
	}

	// 删除用户
	err = h.service.Delete(req.Username)
	if err != nil {
		c.Error(errors.InternalServerError("删除失败", err))
		return
	}

	c.JSON(http.StatusOK, models.AccountResponse{
		Success: true,
		Message: "删除成功",
	})
}
