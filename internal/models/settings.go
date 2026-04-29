package models

// UserSettings 用户个性化设置模型
// 存储用户的偏好设置
type UserSettings struct {
	// Username 用户名
	Username string `json:"username"`
	// Theme 主题设置 (light/dark)
	Theme string `json:"theme"`
	// Language 语言设置
	Language string `json:"language"`
	// DefaultView 默认视图设置
	DefaultView string `json:"default_view"`
	// Notifications 通知设置
	Notifications bool `json:"notifications"`
	// SearchHistory 搜索历史
	SearchHistory []string `json:"search_history"`
	// FavoriteSchools 收藏的学校
	FavoriteSchools []string `json:"favorite_schools"`
}

// GetDefaultSettings 获取默认设置
// 返回默认的用户设置
func GetDefaultSettings(username string) UserSettings {
	return UserSettings{
		Username:        username,
		Theme:           "light",
		Language:        "zh-CN",
		DefaultView:     "list",
		Notifications:   true,
		SearchHistory:   []string{},
		FavoriteSchools: []string{},
	}
}
