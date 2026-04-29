package service

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"school-go/internal/models"
)

// SettingsService 设置服务接口
// 定义了用户设置相关的业务逻辑操作
type SettingsService interface {
	// GetSettings 获取用户设置
	// username: 用户名
	// 返回用户设置和可能的错误
	GetSettings(username string) (models.UserSettings, error)
	
	// SaveSettings 保存用户设置
	// settings: 用户设置
	// 返回可能的错误
	SaveSettings(settings models.UserSettings) error
	
	// AddToSearchHistory 添加到搜索历史
	// username: 用户名
	// keyword: 搜索关键词
	// 返回可能的错误
	AddToSearchHistory(username string, keyword string) error
	
	// AddToFavorites 添加到收藏
	// username: 用户名
	// schoolID: 学校ID
	// 返回可能的错误
	AddToFavorites(username string, schoolID string) error
	
	// RemoveFromFavorites 从收藏中移除
	// username: 用户名
	// schoolID: 学校ID
	// 返回可能的错误
	RemoveFromFavorites(username string, schoolID string) error
}

// settingsService 设置服务实现
// 封装了用户设置相关的业务逻辑
type settingsService struct {
	// settingsDir 设置文件目录
	settingsDir string
	// mutex 并发安全锁
	mutex sync.RWMutex
}

// NewSettingsService 创建设置服务实例
// 返回 SettingsService 接口实现
func NewSettingsService() SettingsService {
	settingsDir := filepath.Join("data", "settings")
	os.MkdirAll(settingsDir, 0755)
	
	return &settingsService{
		settingsDir: settingsDir,
	}
}

// GetSettings 获取用户设置
func (s *settingsService) GetSettings(username string) (models.UserSettings, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	
	filePath := s.getSettingsFilePath(username)
	
	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// 返回默认设置
		return models.GetDefaultSettings(username), nil
	}
	
	// 读取文件
	data, err := os.ReadFile(filePath)
	if err != nil {
		return models.UserSettings{}, err
	}
	
	// 解析JSON
	var settings models.UserSettings
	if err := json.Unmarshal(data, &settings); err != nil {
		return models.UserSettings{}, err
	}
	
	return settings, nil
}

// SaveSettings 保存用户设置
func (s *settingsService) SaveSettings(settings models.UserSettings) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	
	filePath := s.getSettingsFilePath(settings.Username)
	
	// 序列化JSON
	data, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return err
	}
	
	// 写入文件
	return os.WriteFile(filePath, data, 0644)
}

// AddToSearchHistory 添加到搜索历史
func (s *settingsService) AddToSearchHistory(username string, keyword string) error {
	// 获取当前设置
	settings, err := s.GetSettings(username)
	if err != nil {
		return err
	}
	
	// 检查是否已存在
	for _, item := range settings.SearchHistory {
		if item == keyword {
			return nil
		}
	}
	
	// 添加到历史记录
	settings.SearchHistory = append(settings.SearchHistory, keyword)
	
	// 限制历史记录长度
	if len(settings.SearchHistory) > 10 {
		settings.SearchHistory = settings.SearchHistory[len(settings.SearchHistory)-10:]
	}
	
	// 保存设置
	return s.SaveSettings(settings)
}

// AddToFavorites 添加到收藏
func (s *settingsService) AddToFavorites(username string, schoolID string) error {
	// 获取当前设置
	settings, err := s.GetSettings(username)
	if err != nil {
		return err
	}
	
	// 检查是否已存在
	for _, item := range settings.FavoriteSchools {
		if item == schoolID {
			return nil
		}
	}
	
	// 添加到收藏
	settings.FavoriteSchools = append(settings.FavoriteSchools, schoolID)
	
	// 保存设置
	return s.SaveSettings(settings)
}

// RemoveFromFavorites 从收藏中移除
func (s *settingsService) RemoveFromFavorites(username string, schoolID string) error {
	// 获取当前设置
	settings, err := s.GetSettings(username)
	if err != nil {
		return err
	}
	
	// 移除收藏
	newFavorites := []string{}
	for _, item := range settings.FavoriteSchools {
		if item != schoolID {
			newFavorites = append(newFavorites, item)
		}
	}
	settings.FavoriteSchools = newFavorites
	
	// 保存设置
	return s.SaveSettings(settings)
}

// getSettingsFilePath 获取设置文件路径
func (s *settingsService) getSettingsFilePath(username string) string {
	return filepath.Join(s.settingsDir, username+"_settings.json")
}
