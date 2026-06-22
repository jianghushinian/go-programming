package main

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/reflectx"
)

type User struct {
	ID        int
	Name      sql.NullString `json:"username"`
	Email     string
	Age       int
	Birthday  time.Time
	Salary    Salary
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
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

func ConnDB(driver, dsn string) (*sqlx.DB, error) {
	// 1. 使用 sqlx.Open 连接数据库
	db, err := sqlx.Open(driver, dsn)

	// 2. 使用 sqlx.Open 变体方法连接数据库，如果出现错误直接 panic
	db = sqlx.MustOpen(driver, dsn)

	// 3. 如果已经有了 *sql.DB 对象，可以使用 sqlx.NewDb 连接数据库，得到 *sqlx.DB 对象
	sqlDB, err := sql.Open(driver, dsn)
	db = sqlx.NewDb(sqlDB, driver)

	// 4. 使用 sqlx.Connect 连接数据库，等价于 sqlx.Open + db.Ping
	db, err = sqlx.Connect(driver, dsn)

	// 5. 使用 sqlx.Connect 变体方法连接数据库，如果出现错误直接 panic
	db = sqlx.MustConnect(driver, dsn)

	return db, err
}

func MustCreateUser(db *sqlx.DB, user User) (int64, error) {
	res := db.MustExec(
		`INSERT INTO user(name, email, age, birthday, salary) VALUES(?, ?, ?, ?, ?)`,
		user.Name, user.Email, user.Age, user.Birthday, user.Salary,
	)
	return res.LastInsertId()
}

func QueryxUsers(db *sqlx.DB) ([]User, error) {
	var us []User
	rows, _ := db.Queryx("SELECT * FROM user")
	defer rows.Close()

	for rows.Next() {
		var u User
		// sqlx 提供了便捷方法 StructScan 可以将查询结果直接扫描到结构体
		_ = rows.StructScan(&u)
		us = append(us, u)
	}
	return us, nil
}

func QueryRowxUser(db *sqlx.DB, id int) (User, error) {
	var u User
	err := db.QueryRowx("SELECT * FROM user WHERE id = ?", id).StructScan(&u)
	return u, err
}

func GetUser(db *sqlx.DB, id int) (User, error) {
	var u User
	err := db.Get(&u, "SELECT * FROM user WHERE id = ?", id)
	return u, err
}

func SelectUsers(db *sqlx.DB) ([]User, error) {
	var us []User
	err := db.Select(&us, "SELECT * FROM user")
	return us, err
}

func SqlxIn(db *sqlx.DB, ids []int64) ([]User, error) {
	query, args, _ := sqlx.In("SELECT * FROM user WHERE id IN (?)", ids)
	query = db.Rebind(query)

	var us []User
	err := db.Select(&us, query, args...)
	return us, err
}

func NamedExec(db *sqlx.DB) (sql.Result, error) {
	// 使用 map
	m := map[string]any{
		"email": "jianghushinian007@outlook.com",
		"age":   30,
	}
	return db.NamedExec(`UPDATE user SET age = :age WHERE email = :email`, m)
}

func NamedQuery(db *sqlx.DB) ([]User, error) {
	// 使用结构体
	u := User{
		Email: "jianghushinian007@outlook.com",
		Age:   30,
	}
	rows, _ := db.NamedQuery("SELECT * FROM user WHERE email = :email OR age = :age", u)
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		_ = rows.StructScan(&user)
		users = append(users, user)
	}
	return users, nil
}

func MustTransaction(db *sqlx.DB) error {
	tx := db.MustBegin()
	tx.MustExec("UPDATE user SET age = 25 WHERE id = ?", 1)
	return tx.Commit()
}

func Transaction(db *sqlx.DB, id int64, name string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	res, err := tx.Exec("UPDATE user SET name = ? WHERE id = ?", name, id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	fmt.Printf("rowsAffected: %d\n", rowsAffected)
	return tx.Commit()
}

func PreparexGetUser(db *sqlx.DB) (User, error) {
	stmt, _ := db.Preparex(`SELECT * FROM user WHERE id = ?`)

	var u User
	err := stmt.Get(&u, 1)
	return u, err
}

func GetUserMissingAge(db *sqlx.DB, id int) (User, error) {
	var user struct {
		ID    int
		Name  string
		Email string
		// 缺少 Age 字段
	}
	// err := db.Get(&user, "SELECT id, name, email, age FROM user WHERE id = ?", id)
	err := db.Unsafe().Get(&user, "SELECT id, name, email, age FROM user WHERE id = ?", id)
	if err != nil {
		return User{}, err
	}
	return User{
		ID:    user.ID,
		Name:  sql.NullString{String: user.Name},
		Email: user.Email,
	}, nil
}

func Scan2Map(db *sqlx.DB) ([]map[string]any, error) {
	rows, _ := db.Queryx("SELECT * FROM user")
	defer rows.Close()

	var result []map[string]any
	for rows.Next() {
		m := map[string]any{}
		_ = rows.MapScan(m)
		result = append(result, m)
	}
	return result, nil
}

func Scan2Slice(db *sqlx.DB) ([][]any, error) {
	rows, _ := db.Queryx("SELECT * FROM user")
	defer rows.Close()

	var result [][]any
	for rows.Next() {
		cols, _ := rows.SliceScan()
		result = append(result, cols)
	}
	return result, nil
}

func MapperFuncUseToUpper(db *sqlx.DB) (User, error) {
	copyDB := sqlx.NewDb(db.DB, db.DriverName())
	copyDB.MapperFunc(strings.ToUpper)

	var user User
	err := copyDB.Get(&user, "SELECT id as ID, name as NAME, email as EMAIL FROM user WHERE id = ?", 1)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func MapperFuncUseJsonTag(db *sqlx.DB) (User, error) {
	copyDB := sqlx.NewDb(db.DB, db.DriverName())
	copyDB.Mapper = reflectx.NewMapperFunc("json", strings.ToLower)

	var user User
	err := copyDB.Get(&user, "SELECT id, name as username, email FROM user WHERE id = ?", 1)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func main() {
	dsn := "root:password@tcp(127.0.0.1:3306)/demo?charset=utf8mb4&parseTime=true&loc=Local"
	db, err := ConnDB("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	user := User{
		Name:     sql.NullString{String: "jianghushinian007", Valid: true},
		Email:    "jianghushinian007@outlook.com",
		Age:      10,
		Birthday: time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local),
		Salary: Salary{
			Month: 100000,
			Year:  10000000,
		},
	}
	createUser, err := MustCreateUser(db, user)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("user: %#v\n", createUser)

	queryRowxUser, err := QueryRowxUser(db, 1)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("user: %#v\n", queryRowxUser)

	getUser, err := GetUser(db, 1)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("user: %#v\n", getUser)

	users, err := SelectUsers(db)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("users: %#v\n", users)

	in, err := SqlxIn(db, []int64{1, 2, 3})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("in: %#v\n", in)

	exec, err := NamedExec(db)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("exec: %#v\n", exec)

	query, err := NamedQuery(db)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("query: %#v\n", query)

	err = MustTransaction(db)
	if err != nil {
		log.Fatal(err)
	}

	err = Transaction(db, 1, "jianghushinian")
	if err != nil {
		log.Fatal(err)
	}

	preparexGetUser, err := PreparexGetUser(db)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("preparexGetUser: %#v\n", preparexGetUser)

	u, err := GetUserMissingAge(db, 1)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("u: %#v\n", u)

	scan2Map, err := Scan2Map(db)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("scan2Map: %#v\n", scan2Map)

	scan2Slice, err := Scan2Slice(db)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("scan2Slice: %#v\n", scan2Slice)

	upper, err := MapperFuncUseToUpper(db)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("upper: %#v\n", upper)

	tag, err := MapperFuncUseJsonTag(db)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("tag: %#v\n", tag)
}
