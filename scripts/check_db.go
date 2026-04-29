package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "school.db")
	if err != nil {
		fmt.Printf("打开数据库失败: %v\n", err)
		return
	}
	defer db.Close()

	fmt.Println("=== 表结构 ===")
	rows, err := db.Query("PRAGMA table_info(universities)")
	if err != nil {
		fmt.Printf("查询表结构失败: %v\n", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var cid int
		var name string
		var typ string
		var notnull int
		var dflt_value interface{}
		var pk int
		err := rows.Scan(&cid, &name, &typ, &notnull, &dflt_value, &pk)
		if err != nil {
			fmt.Printf("扫描行失败: %v\n", err)
			continue
		}
		fmt.Printf("  %s (%s)\n", name, typ)
	}

	fmt.Println("\n=== 学校列表（前20条）===")
	rows, err = db.Query("SELECT 学校ID, 学校名称, 占地面积 FROM universities LIMIT 20")
	if err != nil {
		fmt.Printf("查询学校列表失败: %v\n", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id, name, area string
		err := rows.Scan(&id, &name, &area)
		if err != nil {
			fmt.Printf("扫描行失败: %v\n", err)
			continue
		}
		if area == "" {
			area = "未填写"
		}
		fmt.Printf("  %s | %s | 占地面积: %s\n", id, name, area)
	}

	fmt.Println("\n=== 统计信息 ===")
	var total int
	err = db.QueryRow("SELECT COUNT(*) FROM universities").Scan(&total)
	if err != nil {
		fmt.Printf("查询总学校数失败: %v\n", err)
		return
	}

	var hasArea int
	err = db.QueryRow("SELECT COUNT(*) FROM universities WHERE 占地面积 IS NOT NULL AND 占地面积 != ''").Scan(&hasArea)
	if err != nil {
		fmt.Printf("查询已填写占地面积失败: %v\n", err)
		return
	}

	fmt.Printf("  总学校数: %d\n", total)
	fmt.Printf("  已填写占地面积: %d\n", hasArea)
	fmt.Printf("  未填写占地面积: %d\n", total-hasArea)
}
