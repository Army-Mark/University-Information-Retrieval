package models

// OperationLog 操作日志模型
type OperationLog struct {
	ID            int    `json:"id"`
	OperationType string `json:"operation_type"`
	TableName     string `json:"table_name"`
	RecordID      string `json:"record_id"`
	UserID        int    `json:"user_id"`
	Username      string `json:"username"`
	OperationTime string `json:"operation_time"`
	OldData       string `json:"old_data"`
	NewData       string `json:"new_data"`
	IPAddress     string `json:"ip_address"`
	UserAgent     string `json:"user_agent"`
}
