package service

import (
	"encoding/json"
	"fmt"
	"school-go/internal/models"
	"school-go/internal/repository"
)

// OperationLogService 操作日志服务接口
type OperationLogService interface {
	// LogOperation 记录操作日志
	// operationType: 操作类型 (create, read, update, delete)
	// tableName: 表名
	// recordID: 记录ID
	// userID: 用户ID
	// username: 用户名
	// oldData: 旧数据
	// newData: 新数据
	// ipAddress: IP地址
	// userAgent: 用户代理
	LogOperation(operationType, tableName, recordID string, userID int, username, ipAddress, userAgent string, oldData, newData interface{})
	
	// GetAllLogs 获取所有操作日志
	GetAllLogs() ([]models.OperationLog, error)
	
	// GetLogsByUserID 根据用户ID获取操作日志
	GetLogsByUserID(userID int) ([]models.OperationLog, error)
	
	// GetLogsByUsername 根据用户名获取操作日志
	GetLogsByUsername(username string) ([]models.OperationLog, error)
	
	// DeleteLogs 删除操作日志
	DeleteLogs(ids []int) error
}

// operationLogService 操作日志服务实现
type operationLogService struct {
	repo repository.OperationLogRepository
}

// NewOperationLogService 创建操作日志服务实例
func NewOperationLogService() OperationLogService {
	return &operationLogService{
		repo: repository.NewOperationLogRepository(),
	}
}

// LogOperation 记录操作日志
func (s *operationLogService) LogOperation(operationType, tableName, recordID string, userID int, username, ipAddress, userAgent string, oldData, newData interface{}) {
	// 序列化旧数据和新数据
	var oldDataJSON, newDataJSON string
	
	if oldData != nil {
		if data, err := json.Marshal(oldData); err == nil {
			oldDataJSON = string(data)
		}
	}
	
	if newData != nil {
		if data, err := json.Marshal(newData); err == nil {
			newDataJSON = string(data)
		}
	}
	
	// 创建操作日志
	log := &models.OperationLog{
		OperationType: operationType,
		TableName:     tableName,
		RecordID:      recordID,
		UserID:        userID,
		Username:      username,
		OldData:       oldDataJSON,
		NewData:       newDataJSON,
		IPAddress:     ipAddress,
		UserAgent:     userAgent,
	}
	
	// 保存到数据库
	err := s.repo.Create(log)
	if err != nil {
		// 记录错误但不影响主流程
		// 可以考虑使用日志库记录错误
	}
}

// GetAllLogs 获取所有操作日志
func (s *operationLogService) GetAllLogs() ([]models.OperationLog, error) {
	return s.repo.GetAll()
}

// GetLogsByUserID 根据用户ID获取操作日志
func (s *operationLogService) GetLogsByUserID(userID int) ([]models.OperationLog, error) {
	return s.repo.GetByUserID(userID)
}

// GetLogsByUsername 根据用户名获取操作日志
func (s *operationLogService) GetLogsByUsername(username string) ([]models.OperationLog, error) {
	return s.repo.GetByUsername(username)
}

// DeleteLogs 删除操作日志
func (s *operationLogService) DeleteLogs(ids []int) error {
	fmt.Printf("Service层: 开始删除日志，ID列表: %v\n", ids)
	err := s.repo.Delete(ids)
	if err != nil {
		fmt.Printf("Service层: 删除日志失败: %v\n", err)
		return err
	}
	fmt.Printf("Service层: 成功删除 %d 条日志\n", len(ids))
	return nil
}
