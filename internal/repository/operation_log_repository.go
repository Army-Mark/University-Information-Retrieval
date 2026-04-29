package repository

import (
	"database/sql"
	"fmt"
	"school-go/internal/models"
)

// OperationLogRepository 操作日志仓库接口
type OperationLogRepository interface {
	// GetAll 获取所有操作日志
	GetAll() ([]models.OperationLog, error)
	
	// GetByUserID 根据用户ID获取操作日志
	GetByUserID(userID int) ([]models.OperationLog, error)
	
	// GetByUsername 根据用户名获取操作日志
	GetByUsername(username string) ([]models.OperationLog, error)
	
	// GetByID 根据ID获取操作日志
	GetByID(id int) (models.OperationLog, error)
	
	// Create 创建操作日志
	Create(log *models.OperationLog) error
	
	// Delete 删除操作日志
	Delete(ids []int) error
}

// operationLogRepository 操作日志仓库实现
type operationLogRepository struct {
	db *sql.DB
}

// NewOperationLogRepository 创建操作日志仓库实例
func NewOperationLogRepository() OperationLogRepository {
	return &operationLogRepository{
		db: GetDB(),
	}
}

// GetAll 获取所有操作日志
func (r *operationLogRepository) GetAll() ([]models.OperationLog, error) {
	query := `
	SELECT id, operation_type, table_name, record_id, user_id, username, 
	       operation_time, old_data, new_data, ip_address, user_agent 
	FROM operation_logs 
	ORDER BY operation_time DESC
	`
	
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var logs []models.OperationLog
	for rows.Next() {
		var log models.OperationLog
		err := rows.Scan(
			&log.ID,
			&log.OperationType,
			&log.TableName,
			&log.RecordID,
			&log.UserID,
			&log.Username,
			&log.OperationTime,
			&log.OldData,
			&log.NewData,
			&log.IPAddress,
			&log.UserAgent,
		)
		if err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}
	
	return logs, nil
}

// GetByID 根据ID获取操作日志
func (r *operationLogRepository) GetByID(id int) (models.OperationLog, error) {
	query := `
	SELECT id, operation_type, table_name, record_id, user_id, username, 
	       operation_time, old_data, new_data, ip_address, user_agent 
	FROM operation_logs 
	WHERE id = ?
	`
	
	var log models.OperationLog
	err := r.db.QueryRow(query, id).Scan(
		&log.ID,
		&log.OperationType,
		&log.TableName,
		&log.RecordID,
		&log.UserID,
		&log.Username,
		&log.OperationTime,
		&log.OldData,
		&log.NewData,
		&log.IPAddress,
		&log.UserAgent,
	)
	
	if err != nil {
		return models.OperationLog{}, err
	}
	
	return log, nil
}

// Create 创建操作日志
func (r *operationLogRepository) Create(log *models.OperationLog) error {
	query := `
	INSERT INTO operation_logs (
		operation_type, table_name, record_id, user_id, username, 
		old_data, new_data, ip_address, user_agent
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	
	_, err := r.db.Exec(
		query,
		log.OperationType,
		log.TableName,
		log.RecordID,
		log.UserID,
		log.Username,
		log.OldData,
		log.NewData,
		log.IPAddress,
		log.UserAgent,
	)
	
	return err
}

// GetByUserID 根据用户ID获取操作日志
func (r *operationLogRepository) GetByUserID(userID int) ([]models.OperationLog, error) {
	query := `
	SELECT id, operation_type, table_name, record_id, user_id, username, 
	       operation_time, old_data, new_data, ip_address, user_agent 
	FROM operation_logs 
	WHERE user_id = ?
	ORDER BY operation_time DESC
	`
	
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var logs []models.OperationLog
	for rows.Next() {
		var log models.OperationLog
		err := rows.Scan(
			&log.ID,
			&log.OperationType,
			&log.TableName,
			&log.RecordID,
			&log.UserID,
			&log.Username,
			&log.OperationTime,
			&log.OldData,
			&log.NewData,
			&log.IPAddress,
			&log.UserAgent,
		)
		if err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}
	
	return logs, nil
}

// GetByUsername 根据用户名获取操作日志
func (r *operationLogRepository) GetByUsername(username string) ([]models.OperationLog, error) {
	query := `
	SELECT id, operation_type, table_name, record_id, user_id, username, 
	       operation_time, old_data, new_data, ip_address, user_agent 
	FROM operation_logs 
	WHERE username = ?
	ORDER BY operation_time DESC
	`
	
	rows, err := r.db.Query(query, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var logs []models.OperationLog
	for rows.Next() {
		var log models.OperationLog
		err := rows.Scan(
			&log.ID,
			&log.OperationType,
			&log.TableName,
			&log.RecordID,
			&log.UserID,
			&log.Username,
			&log.OperationTime,
			&log.OldData,
			&log.NewData,
			&log.IPAddress,
			&log.UserAgent,
		)
		if err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}
	
	return logs, nil
}

// Delete 删除操作日志
func (r *operationLogRepository) Delete(ids []int) error {
	if len(ids) == 0 {
		return nil
	}

	// 构建SQL语句
	query := `DELETE FROM operation_logs WHERE id IN (`
	args := make([]interface{}, len(ids))

	for i, id := range ids {
		if i > 0 {
			query += `, `
		}
		query += `?`
		args[i] = id
	}
	query += `)`

	fmt.Printf("Repository层: 执行删除SQL: %s, 参数: %v\n", query, args)
	result, err := r.db.Exec(query, args...)
	if err != nil {
		fmt.Printf("Repository层: 删除失败: %v\n", err)
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	fmt.Printf("Repository层: 影响的行数: %d\n", rowsAffected)
	return nil
}
