package repository

import (
	"database/sql"
	"log"
	"os"
	"sync"
	"time"

	_ "modernc.org/sqlite"
)

// 全局变量定义
var (
	// db 数据库连接实例
	db *sql.DB
	// rootDir 应用根目录
	rootDir string
	// initOnce 确保数据库初始化只执行一次
	initOnce sync.Once
	// initErr 初始化错误
	initErr error
)

// SetRootDir 设置应用根目录
// dir: 应用根目录路径
func SetRootDir(dir string) {
	rootDir = dir
}

// InitDB 初始化数据库连接
// 确保数据库连接只初始化一次，设置连接池参数，并创建必要的索引
// 返回可能的错误
func InitDB() error {
	initOnce.Do(func() {
		// 从环境变量获取数据库路径，默认使用相对路径
		dbPath := os.Getenv("DB_PATH")
		if dbPath == "" {
			dbPath = "school.db"
		}

		// 打开数据库连接
		db, initErr = sql.Open("sqlite", dbPath)
		if initErr != nil {
			return
		}

		// 设置连接池参数
		// 最大打开连接数
		db.SetMaxOpenConns(25)
		// 最大空闲连接数
		db.SetMaxIdleConns(10)
		// 连接最大生命周期
		db.SetConnMaxLifetime(time.Hour)
		// 连接最大空闲时间
		db.SetConnMaxIdleTime(30 * time.Minute)

		// 测试数据库连接
		initErr = db.Ping()
		if initErr != nil {
			return
		}

		// 创建数据库表
		createTables()

		// 创建数据库索引
		createIndexes()
	})
	return initErr
}

// createTables 创建数据库表
// 确保必要的表结构存在
func createTables() {
	// 检查用户表是否存在
	var exists bool
	err := db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='users'").Scan(&exists)
	if err != nil {
		return
	}

	if !exists {
		// 创建新的用户表，包含 role 列
		userTable := `
		CREATE TABLE users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL,
			role TEXT NOT NULL DEFAULT 'user'
		);
		`
		if _, err := db.Exec(userTable); err != nil {
			return
		}
	} else {
		// 检查表结构，看看是否有 role 列
		var hasRoleColumn bool
		err := db.QueryRow("SELECT COUNT(*) FROM pragma_table_info('users') WHERE name='role'").Scan(&hasRoleColumn)
		if err != nil || !hasRoleColumn {
			// 如果没有 role 列，需要重新创建表
			// 1. 创建临时表
			if _, err := db.Exec("CREATE TABLE users_temp AS SELECT * FROM users"); err != nil {
				return
			}
			// 2. 删除原表
			if _, err := db.Exec("DROP TABLE users"); err != nil {
				return
			}
			// 3. 创建新表
			userTable := `
			CREATE TABLE users (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				username TEXT UNIQUE NOT NULL,
				password TEXT NOT NULL,
				role TEXT NOT NULL DEFAULT 'user'
			);
			`
			if _, err := db.Exec(userTable); err != nil {
				return
			}
			// 4. 从临时表导入数据
			if _, err := db.Exec("INSERT INTO users (id, username, password, role) SELECT id, username, password, 'user' FROM users_temp"); err != nil {
				return
			}
			// 5. 删除临时表
			if _, err := db.Exec("DROP TABLE users_temp"); err != nil {
				return
			}
		}
	}

	// 创建操作记录表
	operationLogTable := `
	CREATE TABLE IF NOT EXISTS operation_logs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		operation_type TEXT NOT NULL,
		table_name TEXT NOT NULL,
		record_id TEXT NOT NULL,
		user_id INTEGER,
		username TEXT NOT NULL,
		operation_time DATETIME DEFAULT CURRENT_TIMESTAMP,
		old_data TEXT,
		new_data TEXT,
		ip_address TEXT,
		user_agent TEXT
	);
	`
	if _, err := db.Exec(operationLogTable); err != nil {
		log.Printf("Error creating operation_logs table: %v", err)
	}

	// 创建默认管理员用户（如果不存在）
	// 先检查是否存在
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM users WHERE username = 'admin'").Scan(&count)
	if err != nil {
		log.Printf("Error checking admin user: %v", err)
	} else {
		log.Printf("Admin user count: %d", count)
		if count == 0 {
			// 密码使用 bcrypt 哈希值，对应密码 admin123
			hashedPassword := "$2a$10$99yjfkfaDT.k3Iot5rqvHe.uier9SJA4yQ3uM6u9TOUjAmDBaFzZu"
			insertAdmin := `
			INSERT INTO users (username, password, role) 
			VALUES ('admin', ?, 'admin');
			`
			_, err := db.Exec(insertAdmin, hashedPassword)
			if err != nil {
				log.Printf("Error inserting admin user: %v", err)
			} else {
				log.Println("Admin user created successfully")
			}
		} else {
			// 更新现有管理员用户的密码和角色
			hashedPassword := "$2a$10$99yjfkfaDT.k3Iot5rqvHe.uier9SJA4yQ3uM6u9TOUjAmDBaFzZu"
			updateAdmin := `
			UPDATE users SET password = ?, role = 'admin' WHERE username = 'admin';
			`
			_, err := db.Exec(updateAdmin, hashedPassword)
			if err != nil {
				log.Printf("Error updating admin user: %v", err)
			} else {
				log.Println("Admin user updated successfully")
			}
		}
	}
}

// createIndexes 创建数据库索引
// 为常用查询字段创建索引，提高查询性能
func createIndexes() {
	// 定义需要创建的索引
	indexes := []string{
		// 为学校ID创建索引，加速ID查询
		"CREATE INDEX IF NOT EXISTS idx_school_id ON universities(学校ID)",
		// 为学校名称创建索引，加速名称搜索
		"CREATE INDEX IF NOT EXISTS idx_school_name ON universities(学校名称)",
		// 为用户名创建索引，加速用户查询
		"CREATE INDEX IF NOT EXISTS idx_username ON users(username)",
	}

	// 执行索引创建
	for _, idx := range indexes {
		if _, err := db.Exec(idx); err != nil {
			return
		}
	}
}

// GetDB 获取数据库连接
// 如果数据库连接未初始化，则自动初始化
// 返回数据库连接实例
func GetDB() *sql.DB {
	if db == nil {
		InitDB()
	}
	return db
}

// CloseDB 关闭数据库连接
// 返回可能的错误
func CloseDB() error {
	if db != nil {
		return db.Close()
	}
	return nil
}

// GetDBStats 获取数据库连接池状态
// 返回数据库连接池的统计信息
func GetDBStats() sql.DBStats {
	if db == nil {
		InitDB()
	}
	return db.Stats()
}
