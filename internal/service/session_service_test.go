package service

import (
	"testing"
	"time"
)

// TestSessionService_CreateSession 测试创建会话功能
func TestSessionService_CreateSession(t *testing.T) {
	service := NewSessionService()
	
	// 创建会话
	sessionID, err := service.CreateSession("testuser")
	if err != nil {
		t.Errorf("创建会话失败: %v", err)
	}
	
	if sessionID == "" {
		t.Error("会话ID为空")
	}
	
	// 验证会话是否创建成功
	session, err := service.GetSession(sessionID)
	if err != nil || session == nil {
		t.Errorf("获取会话失败: %v", err)
	}
	
	if session.Username != "testuser" {
		t.Errorf("会话用户名不匹配，期望: testuser, 实际: %s", session.Username)
	}
}

// TestSessionService_GetSession 测试获取会话功能
func TestSessionService_GetSession(t *testing.T) {
	service := NewSessionService()
	
	// 创建会话
	sessionID, _ := service.CreateSession("testuser")
	
	// 测试获取存在的会话
	session, err := service.GetSession(sessionID)
	if err != nil || session == nil {
		t.Errorf("获取会话失败: %v", err)
	}
	
	// 测试获取不存在的会话
	nonexistentSession, err := service.GetSession("nonexistent")
	if nonexistentSession != nil {
		t.Error("获取不存在的会话返回非空")
	}
}

// TestSessionService_DeleteSession 测试删除会话功能
func TestSessionService_DeleteSession(t *testing.T) {
	service := NewSessionService()
	
	// 创建会话
	sessionID, _ := service.CreateSession("testuser")
	
	// 验证会话存在
	session, _ := service.GetSession(sessionID)
	if session == nil {
		t.Error("会话创建失败")
	}
	
	// 删除会话
	service.DeleteSession(sessionID)
	
	// 验证会话已删除
	deletedSession, _ := service.GetSession(sessionID)
	if deletedSession != nil {
		t.Error("会话未被删除")
	}
}

// TestSessionService_SessionExpiration 测试会话过期功能
func TestSessionService_SessionExpiration(t *testing.T) {
	service := NewSessionService()
	
	// 创建会话
	sessionID, _ := service.CreateSession("testuser")
	
	// 验证会话存在
	session, _ := service.GetSession(sessionID)
	if session == nil {
		t.Error("会话创建失败")
	}
	
	// 模拟会话过期（这里我们无法直接修改时间，所以这个测试可能无法完全覆盖过期场景）
	// 但我们可以验证会话创建时设置了过期时间
	if session.ExpiresAt.Before(time.Now()) {
		t.Error("会话过期时间设置错误")
	}
}
