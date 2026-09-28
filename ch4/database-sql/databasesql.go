package main

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	ID        int
	Name      sql.NullString
	Email     string
	Age       int
	Birthday  *time.Time
	Salary    Salary
	CreatedAt time.Time
	UpdatedAt string
}

type Salary struct {
	Month int `json:"month"`
	Year  int `json:"year"`
}

func (s *Salary) Scan(src any) error {
	if src == nil {
		return nil
	}

	var buf []byte
	switch v := src.(type) {
	case []byte:
		buf = v
	case string:
		buf = []byte(v)
	default:
		return fmt.Errorf("invalid type: %T", src)
	}

	err := json.Unmarshal(buf, s)
	return err
}

func (s Salary) Value() (driver.Value, error) {
	v, err := json.Marshal(s)
	return string(v), err
}

func CreateUser(db *sql.DB) (int64, error) {
	birthday := time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local)
	user := User{
		Name:     sql.NullString{String: "jianghushinian007", Valid: true},
		Email:    "jianghushinian007@outlook.com",
		Age:      10,
		Birthday: &birthday,
		Salary: Salary{
			Month: 100000,
			Year:  10000000,
		},
	}
	res, err := db.Exec(
		`INSERT INTO user(name, email, age, birthday, salary) VALUES(?, ?, ?, ?, ?)`,
		user.Name, user.Email, user.Age, user.Birthday, user.Salary,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func CreateUsers(db *sql.DB) ([]int64, error) {
	stmt, err := db.Prepare("INSERT INTO user(name, email, age, birthday, salary) VALUES(?, ?, ?, ?, ?)")
	if err != nil {
		return nil, err
	}
	// 注意：预处理对象是需要关闭的
	defer stmt.Close()

	users := []User{
		{
			Name:  sql.NullString{String: "jianghushinian", Valid: true},
			Email: "jianghushinian007@163.com",
		},
		{
			Name:  sql.NullString{String: "江湖十年", Valid: true},
			Email: "jianghushinian@qq.com",
		},
	}

	var ids []int64
	for _, user := range users {
		res, err := stmt.Exec(user.Name, user.Email, user.Age, user.Birthday, user.Salary)
		if err != nil {
			return nil, err
		}
		id, err := res.LastInsertId()
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func GetUsers(db *sql.DB) ([]User, error) {
	rows, err := db.Query("SELECT * FROM user")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Age,
			&user.Birthday, &user.Salary, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, rows.Err()
}

func GetUser(db *sql.DB, id int64) (User, error) {
	var user User
	row := db.QueryRow("SELECT * FROM user WHERE id = ?", id)
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Age,
		&user.Birthday, &user.Salary, &user.CreatedAt, &user.UpdatedAt)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return user, fmt.Errorf("no user with id %d", id)
	case err != nil:
		return user, err
	}
	return user, row.Err()
}

func UpdateUserName(db *sql.DB, id int64, name string) (affected int64, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	res, err := db.ExecContext(ctx, "UPDATE user SET name = ? WHERE id = ?", name, id)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func DeleteUser(db *sql.DB, id int64) (affected int64, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	res, err := db.ExecContext(ctx, "DELETE FROM user WHERE id = ?", id)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func Transaction(db *sql.DB, id int64, name string) error {
	ctx := context.Background()
	// 开启事务
	tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return err
	}
	// 执行 SQL
	_, err = tx.ExecContext(ctx, "UPDATE user SET name = ? WHERE id = ?", name, id)
	if err != nil {
		// 回滚事务
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			log.Printf("update failed: %v, unable to rollback: %v\n", err, rollbackErr)
		}
		return err
	}
	// 提交事务
	return tx.Commit()
}

func HandleUnknownColumns(db *sql.DB, id int64) error {
	rows, _ := db.Query("SELECT * FROM user WHERE id = ?", id)
	defer rows.Close()

	// 获取列名与列类型
	cols, _ := rows.Columns()
	types, _ := rows.ColumnTypes()
	for i, t := range types {
		fmt.Printf("%s: %s\n", cols[i], t.DatabaseTypeName())
	}

	// 使用 sql.RawBytes 动态接收所有列
	vals := make([]sql.RawBytes, len(cols))
	ptrs := make([]any, len(cols))
	for i := range vals {
		ptrs[i] = &vals[i]
	}

	// 扫描并打印结果
	for rows.Next() {
		_ = rows.Scan(ptrs...)
		for i, col := range cols {
			fmt.Printf("%s = %s\n", col, string(vals[i]))
		}
	}
	return rows.Err()
}

func main() {
	dsn := "root:password@tcp(127.0.0.1:3306)/demo?charset=utf8mb4&parseTime=true&loc=Local"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.SetMaxOpenConns(25)                 // 设置最大的并发连接数（in-use + idle）
	db.SetMaxIdleConns(25)                 // 设置最大的空闲连接数（idle）
	db.SetConnMaxLifetime(5 * time.Minute) // 设置连接的最大生命周期

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	user, err := CreateUser(db)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("user: %#v\n", user)

	users, err := CreateUsers(db)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("users: %#v\n", users)

	getUsers, err := GetUsers(db)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("users: %#v\n", getUsers)

	getUser, err := GetUser(db, 1)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("user: %#v\n", getUser)

	affected, err := UpdateUserName(db, 1, "jianghu")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("updated affected: %d\n", affected)

	affected, err = DeleteUser(db, 1)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("deleted affected: %d\n", affected)

	err = Transaction(db, 1, "shinian")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("transaction done")

	err = HandleUnknownColumns(db, 2)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("handle unknown columns")
}
