package models

import (
	"time"
)

// Session 会话模型
type Session struct {
	SessionID string    `json:"session_id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}
