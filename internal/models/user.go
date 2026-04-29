package models

// User 用户模型
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"` // 角色：admin, user
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// AccountResponse 账户响应
type AccountResponse struct {
	Success  bool   `json:"success"`
	Message  string `json:"message"`
	Accounts []User `json:"accounts,omitempty"`
	LoggedIn bool   `json:"logged_in,omitempty"`
	Username string `json:"username,omitempty"`
}

// AddAccountRequest 添加账户请求
type AddAccountRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role" binding:"required"` // 角色：admin, user
}

// UpdateAccountRequest 更新账户请求
type UpdateAccountRequest struct {
	OldUsername string `json:"oldUsername" binding:"required"`
	NewUsername string `json:"newUsername" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required"`
	NewRole     string `json:"newRole" binding:"required"` // 角色：admin, user
}

// DeleteAccountRequest 删除账户请求
type DeleteAccountRequest struct {
	Username string `json:"username" binding:"required"`
}
