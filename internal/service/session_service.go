package service

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"
	"school-go/internal/models"
)

// SessionService 会话服务接口
// 定义了会话管理相关的业务逻辑操作
type SessionService interface {
	// CreateSession 创建会话
	// 生成随机会话ID并存储会话信息
	// username: 用户名
	// 返回会话ID和可能的错误
	CreateSession(username string) (string, error)
	
	// GetSession 获取会话
	// 检查会话是否存在且未过期
	// sessionID: 会话ID
	// 返回会话信息和可能的错误
	GetSession(sessionID string) (*models.Session, error)
	
	// DeleteSession 删除会话
	// 从会话存储中移除指定的会话
	// sessionID: 会话ID
	DeleteSession(sessionID string)
}

// sessionService 会话服务实现
// 封装了会话管理相关的业务逻辑
type sessionService struct {
	// sessions 会话存储
	sessions map[string]*models.Session
	// mutex 并发安全锁
	mutex sync.RWMutex
}

// 全局会话服务实例
var globalSessionService SessionService

// 确保全局会话服务实例只初始化一次
var sessionServiceOnce sync.Once

// GetSessionService 获取全局会话服务实例
// 返回 SessionService 接口实现
func GetSessionService() SessionService {
	sessionServiceOnce.Do(func() {
		globalSessionService = &sessionService{
			sessions: make(map[string]*models.Session),
		}
	})
	return globalSessionService
}

// NewSessionService 创建会话服务实例
// 返回 SessionService 接口实现
func NewSessionService() SessionService {
	return GetSessionService()
}

// generateSessionID 生成随机会话ID
// 使用32字节随机数据生成会话ID
// 返回会话ID和可能的错误
func generateSessionID() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// CreateSession 创建会话
// 生成随机会话ID并存储会话信息
// 会话有效期为24小时
func (s *sessionService) CreateSession(username string) (string, error) {
	sessionID, err := generateSessionID()
	if err != nil {
		return "", err
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.sessions[sessionID] = &models.Session{
		SessionID: sessionID,
		Username:  username,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(24 * time.Hour), // 会话有效期24小时
	}

	return sessionID, nil
}

// GetSession 获取会话
// 检查会话是否存在且未过期
// 如果会话过期则自动删除
func (s *sessionService) GetSession(sessionID string) (*models.Session, error) {
	s.mutex.RLock()
	session, exists := s.sessions[sessionID]
	if !exists {
		s.mutex.RUnlock()
		return nil, nil
	}

	// 检查会话是否过期
	if time.Now().After(session.ExpiresAt) {
		s.mutex.RUnlock()
		s.mutex.Lock()
		delete(s.sessions, sessionID)
		s.mutex.Unlock()
		return nil, nil
	}

	sessionCopy := *session
	s.mutex.RUnlock()

	return &sessionCopy, nil
}

// DeleteSession 删除会话
// 从会话存储中移除指定的会话
func (s *sessionService) DeleteSession(sessionID string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	delete(s.sessions, sessionID)
}
