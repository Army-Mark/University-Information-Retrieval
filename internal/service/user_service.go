package service

import (
	"golang.org/x/crypto/bcrypt"
	"school-go/internal/models"
	"school-go/internal/repository"
)

// UserService 用户服务接口
// 定义了用户相关的业务逻辑操作
// 作为 repository 层和 handler 层之间的桥梁
type UserService interface {
	// GetAll 获取所有用户
	// 返回用户列表和可能的错误
	GetAll() ([]models.User, error)

	// GetByUsername 根据用户名获取用户
	// username: 用户名
	// 返回用户信息和可能的错误
	GetByUsername(username string) (*models.User, error)

	// Create 创建用户
	// user: 用户信息
	// 返回可能的错误
	Create(user *models.User) error

	// Update 更新用户
	// oldUsername: 原用户名
	// newUsername: 新用户名
	// newPassword: 新密码
	// newRole: 新角色
	// 返回可能的错误
	Update(oldUsername, newUsername, newPassword, newRole string) error

	// Delete 删除用户
	// username: 用户名
	// 返回可能的错误
	Delete(username string) error

	// SaveAll 保存所有用户
	// users: 用户列表
	// 返回可能的错误
	SaveAll(users []models.User) error

	// Authenticate 验证用户身份
	// username: 用户名
	// password: 密码
	// 返回验证结果
	Authenticate(username, password string) bool
}

// userService 用户服务实现
// 封装了用户相关的业务逻辑
type userService struct {
	// repo 用户数据仓库
	repo repository.UserRepository
}

// NewUserService 创建用户服务实例
// 返回 UserService 接口实现
func NewUserService() UserService {
	return &userService{
		repo: repository.NewUserRepository(),
	}
}

// GetAll 获取所有用户
func (s *userService) GetAll() ([]models.User, error) {
	return s.repo.GetAll()
}

// GetByUsername 根据用户名获取用户
func (s *userService) GetByUsername(username string) (*models.User, error) {
	return s.repo.GetByUsername(username)
}

// Create 创建用户
// 对密码进行bcrypt哈希处理后保存
func (s *userService) Create(user *models.User) error {
	// 对密码进行哈希处理
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	
	err = s.repo.Create(user)
	if err == nil {
		// 记录操作日志
		logService := NewOperationLogService()
		logService.LogOperation("create", "users", user.Username, 0, "system", "", "", nil, user)
	}
	return err
}

// Update 更新用户
// 对新密码进行bcrypt哈希处理后保存
func (s *userService) Update(oldUsername, newUsername, newPassword, newRole string) error {
	// 获取旧数据
	oldUser, err := s.repo.GetByUsername(oldUsername)
	if err != nil {
		return err
	}
	
	// 对新密码进行哈希处理
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	
	err = s.repo.Update(oldUsername, newUsername, string(hashedPassword), newRole)
	if err == nil {
		// 记录操作日志
		logService := NewOperationLogService()
		newUser := &models.User{
			Username: newUsername,
			Password: string(hashedPassword),
			Role:     newRole,
		}
		logService.LogOperation("update", "users", oldUsername, 0, "system", "", "", oldUser, newUser)
	}
	return err
}

// Delete 删除用户
func (s *userService) Delete(username string) error {
	// 获取旧数据
	oldUser, err := s.repo.GetByUsername(username)
	if err != nil {
		return err
	}
	
	err = s.repo.Delete(username)
	if err == nil {
		// 记录操作日志
		logService := NewOperationLogService()
		logService.LogOperation("delete", "users", username, 0, "system", "", "", oldUser, nil)
	}
	return err
}

// SaveAll 保存所有用户
// 对所有用户的密码进行bcrypt哈希处理后保存
func (s *userService) SaveAll(users []models.User) error {
	// 对所有用户的密码进行哈希处理
	for i := range users {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(users[i].Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		users[i].Password = string(hashedPassword)
	}
	return s.repo.SaveAll(users)
}

// Authenticate 验证用户身份
// 使用bcrypt验证密码是否正确
func (s *userService) Authenticate(username, password string) bool {
	user, err := s.repo.GetByUsername(username)
	if err != nil || user == nil {
		return false
	}
	// 使用bcrypt验证密码
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return err == nil
}
