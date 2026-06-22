//go:build ignore
// +build ignore

package main

import (
	"database/sql"
	"fmt"

	// 1. 注册驱动
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// 2. 建立连接
	db, _ := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/demo?charset=utf8mb4&parseTime=true&loc=Local")
	defer db.Close()

	// 3. 执行 SQL
	rows, _ := db.Query("SELECT id, name FROM user")
	defer rows.Close()

	// 4. 解析数据
	for rows.Next() {
		var (
			id   int
			name string
		)
		_ = rows.Scan(&id, &name)
		fmt.Printf("id: %d, name: %s\n", id, name)
	}
}
