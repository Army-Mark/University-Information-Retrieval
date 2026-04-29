package service

import (
	"testing"
	"school-go/internal/models"
)

// TestUserService_Create 测试创建用户功能
func TestUserService_Create(t *testing.T) {
	service := NewUserService()
	
	// 创建测试用户
	user := &models.User{
		Username: "testuser",
		Password: "Test1234", // 符合密码强度要求
		Role:     "user",
	}
	
	err := service.Create(user)
	if err != nil {
		t.Errorf("创建用户失败: %v", err)
	}
	
	// 验证用户是否创建成功
	createdUser, err := service.GetByUsername("testuser")
	if err != nil || createdUser == nil {
		t.Errorf("获取创建的用户失败: %v", err)
	}
	
	// 验证密码是否被哈希处理
	if createdUser.Password == "Test1234" {
		t.Error("密码未被哈希处理")
	}
	
	// 清理测试数据
	service.Delete("testuser")
}

// TestUserService_Authenticate 测试用户认证功能
func TestUserService_Authenticate(t *testing.T) {
	service := NewUserService()
	
	// 创建测试用户
	user := &models.User{
		Username: "authtest",
		Password: "Test1234",
		Role:     "user",
	}
	service.Create(user)
	
	// 测试正确的密码
	if !service.Authenticate("authtest", "Test1234") {
		t.Error("正确密码认证失败")
	}
	
	// 测试错误的密码
	if service.Authenticate("authtest", "WrongPassword") {
		t.Error("错误密码认证成功")
	}
	
	// 测试不存在的用户
	if service.Authenticate("nonexistent", "Test1234") {
		t.Error("不存在用户认证成功")
	}
	
	// 清理测试数据
	service.Delete("authtest")
}

// TestUserService_Update 测试更新用户功能
func TestUserService_Update(t *testing.T) {
	service := NewUserService()
	
	// 创建测试用户
	user := &models.User{
		Username: "updateuser",
		Password: "Test1234",
		Role:     "user",
	}
	service.Create(user)
	
	// 更新用户
	err := service.Update("updateuser", "updateduser", "NewPass123", "admin")
	if err != nil {
		t.Errorf("更新用户失败: %v", err)
	}
	
	// 验证更新是否成功
	updatedUser, err := service.GetByUsername("updateduser")
	if err != nil || updatedUser == nil {
		t.Errorf("获取更新后的用户失败: %v", err)
	}
	
	if updatedUser.Role != "admin" {
		t.Error("用户角色未更新")
	}
	
	// 验证新密码是否被哈希处理
	if updatedUser.Password == "NewPass123" {
		t.Error("新密码未被哈希处理")
	}
	
	// 清理测试数据
	service.Delete("updateduser")
}

// TestUserService_Delete 测试删除用户功能
func TestUserService_Delete(t *testing.T) {
	service := NewUserService()
	
	// 创建测试用户
	user := &models.User{
		Username: "deleteuser",
		Password: "Test1234",
		Role:     "user",
	}
	service.Create(user)
	
	// 验证用户存在
	_, err := service.GetByUsername("deleteuser")
	if err != nil {
		t.Errorf("用户创建失败: %v", err)
	}
	
	// 删除用户
	err = service.Delete("deleteuser")
	if err != nil {
		t.Errorf("删除用户失败: %v", err)
	}
	
	// 验证用户已删除
	deletedUser, err := service.GetByUsername("deleteuser")
	if deletedUser != nil {
		t.Error("用户未被删除")
	}
}

// TestUserService_GetAll 测试获取所有用户功能
func TestUserService_GetAll(t *testing.T) {
	service := NewUserService()
	
	// 获取初始用户数量
	initialUsers, err := service.GetAll()
	if err != nil {
		t.Errorf("获取初始用户失败: %v", err)
	}
	initialCount := len(initialUsers)
	
	// 创建测试用户
	user := &models.User{
		Username: "getalltest",
		Password: "Test1234",
		Role:     "user",
	}
	service.Create(user)
	
	// 验证用户数量增加
	updatedUsers, err := service.GetAll()
	if err != nil {
		t.Errorf("获取更新后的用户失败: %v", err)
	}
	if len(updatedUsers) != initialCount+1 {
		t.Errorf("用户数量未增加，期望: %d, 实际: %d", initialCount+1, len(updatedUsers))
	}
	
	// 清理测试数据
	service.Delete("getalltest")
}
