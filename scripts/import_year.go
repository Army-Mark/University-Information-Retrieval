package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "school.db")
	if err != nil {
		fmt.Printf("打开数据库失败: %v\n", err)
		return
	}
	defer db.Close()

	csvFile := "待更新学校清单.csv"
	if _, err := os.Stat(csvFile); os.IsNotExist(err) {
		fmt.Printf("找不到文件: %s\n", csvFile)
		return
	}

	fmt.Println("=" + strings.Repeat("=", 60))
	fmt.Println("建校时间批量导入工具")
	fmt.Println("=" + strings.Repeat("=", 60))
	fmt.Println("\n使用说明:")
	fmt.Println("1. 在CSV文件的'正确建校时间'列填入准确的建校年份")
	fmt.Println("2. 保存CSV文件")
	fmt.Println("3. 运行此脚本导入数据库")
	fmt.Println("\n格式示例:")
	fmt.Println("  正确建校时间: 2001年 或 2001")
	fmt.Println("\n" + "=" + strings.Repeat("=", 60))

	fmt.Printf("\n是否从 '%s' 导入数据? (yes/no): ", csvFile)
	var confirm string
	fmt.Scanln(&confirm)

	if confirm != "yes" {
		fmt.Println("已取消操作")
		return
	}

	updateFromCSV(db, csvFile)
}

func updateFromCSV(db *sql.DB, csvFile string) {
	file, err := os.Open(csvFile)
	if err != nil {
		fmt.Printf("打开CSV文件失败: %v\n", err)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true

	records, err := reader.ReadAll()
	if err != nil {
		fmt.Printf("读取CSV文件失败: %v\n", err)
		return
	}

	if len(records) < 2 {
		fmt.Println("CSV文件格式错误，缺少标题行或数据行")
		return
	}

	// 解析标题行，找到对应的列索引
	header := records[0]
	idIndex := -1
	nameIndex := -1
	yearIndex := -1

	for i, col := range header {
		switch col {
		case "学校ID":
			idIndex = i
		case "学校名称":
			nameIndex = i
		case "正确建校时间":
			yearIndex = i
		}
	}

	if idIndex == -1 || nameIndex == -1 || yearIndex == -1 {
		fmt.Println("CSV文件格式错误，缺少必要的列")
		return
	}

	updatedCount := 0
	skippedCount := 0
	errorCount := 0

	fmt.Printf("\n正在读取CSV文件: %s\n\n", csvFile)

	for i, record := range records[1:] {
		if len(record) <= max(idIndex, nameIndex, yearIndex) {
			fmt.Printf("✗ 错误: 第%d行数据不完整\n", i+2)
			errorCount++
			continue
		}

		schoolID := record[idIndex]
		schoolName := record[nameIndex]
		newYear := record[yearIndex]

		if newYear == "" {
			skippedCount++
			continue
		}

		// 确保年份格式正确
		if !strings.HasSuffix(newYear, "年") {
			newYear += "年"
		}

		result, err := db.Exec(
			"UPDATE universities SET 建校时间 = ? WHERE 学校ID = ?",
			newYear, schoolID,
		)
		if err != nil {
			fmt.Printf("✗ 错误: %s - %v\n", schoolName, err)
			errorCount++
			continue
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			fmt.Printf("✗ 错误: %s - %v\n", schoolName, err)
			errorCount++
			continue
		}

		if rowsAffected > 0 {
			fmt.Printf("✓ 更新: %s -> %s\n", schoolName, newYear)
			updatedCount++
		} else {
			fmt.Printf("✗ 未找到: %s (ID: %s)\n", schoolName, schoolID)
			errorCount++
		}
	}

	fmt.Println("\n" + "=" + strings.Repeat("=", 60))
	fmt.Println("处理完成!")
	fmt.Println("=" + strings.Repeat("=", 60))
	fmt.Printf("成功更新: %d 所\n", updatedCount)
	fmt.Printf("跳过(未填写): %d 所\n", skippedCount)
	fmt.Printf("错误: %d 所\n", errorCount)
}

// 辅助函数
func max(a, b, c int) int {
	if a >= b && a >= c {
		return a
	}
	if b >= a && b >= c {
		return b
	}
	return c
}
