package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type User struct {
	gorm.Model
	UUID         string         `gorm:"unique"`
	Name         string         `gorm:"column:name"`
	Email        *string        `gorm:"column:email"`
	Age          uint8          `gorm:"column:age"`
	Birthday     *time.Time     `gorm:"column:birthday"`
	MemberNumber sql.NullString `gorm:"column:member_number"`
	ActivatedAt  sql.NullTime   `gorm:"column:activated_at"`
}

func (u *User) TableName() string {
	return "user"
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.UUID = uuid.New().String()
	if u.Name == "admin" {
		return errors.New("invalid name")
	}
	return nil
}

func ConnectMySQL(host, port, user, pass, dbname string) (*gorm.DB, error) {
	format := "%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := fmt.Sprintf(format, user, pass, host, port, dbname)

	slowLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: 3 * time.Millisecond, // 超过 3ms 即记录日志
			LogLevel:      logger.Warn,
		},
	)

	return gorm.Open(mysql.Open(dsn), &gorm.Config{
		// DryRun: true,
		// Logger: logger.Default.LogMode(logger.Info),
		Logger: slowLogger,
	})
}

func SetConnect(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	sqlDB.SetMaxOpenConns(100)                 // 设置数据库的最大打开连接数
	sqlDB.SetMaxIdleConns(100)                 // 设置最大空闲连接数
	sqlDB.SetConnMaxLifetime(10 * time.Second) // 设置空闲连接最大存活时间
	return nil
}

func CreateUser(db *gorm.DB) {
	now := time.Now()
	email := "u1@jianghushinian.com"
	user := User{Name: "user1", Email: &email, Age: 18, Birthday: &now}

	result := db.Create(&user)

	fmt.Printf("user: %+v\n", user) // user.ID 会被自动填充
	fmt.Printf("affected rows: %d\n", result.RowsAffected)
	fmt.Printf("error: %v\n", result.Error)
}

func CreateUsers(db *gorm.DB) {
	now := time.Now()
	email2 := "u2@jianghushinian.com"
	email3 := "u3@jianghushinian.com"
	users := []User{
		{Name: "user2", Email: &email2, Age: 19, Birthday: &now},
		{Name: "user3", Email: &email3, Age: 20, Birthday: &now},
	}

	result := db.Create(&users)

	fmt.Printf("users: %+v\n", users) // user.ID 会被自动填充
	fmt.Printf("affected rows: %d\n", result.RowsAffected)
	fmt.Printf("error: %v\n", result.Error)
}

func FirstLastUser(db *gorm.DB) (*User, *User) {
	var first, last User
	db.First(&first) // 查询第一条
	db.Last(&last)   // 查询最后一条
	return &first, &last
}

func WhereFindUsers(db *gorm.DB) []User {
	var users []User
	db.Where("name != ?", "unknown").Find(&users)
	return users
}

func SelectOrderLimitOffsetUsers(db *gorm.DB) []User {
	var users []User
	db.Select("name", "age").
		Where("age >= ?", 18).
		Order("id desc").
		Limit(1).Offset(2).
		Find(&users)
	return users
}

func CountUsers(db *gorm.DB) uint {
	var count int64
	db.Model(&User{}).Where("age >= ?", 18).Count(&count)
	return uint(count)
}

func GroupHaving(db *gorm.DB) []float64 {
	subQuery := db.Select("AVG(age)").Where("name LIKE ?", "user%").Table("user")
	var results []float64
	db.Model(&User{}).
		Select("AVG(age) as avgage").
		Group("name").
		Having("AVG(age) > (?)", subQuery).
		Find(&results)
	return results
}

func SaveUser(db *gorm.DB) {
	var user User
	db.First(&user)
	user.Name = "江湖十年"
	user.Age = 30
	db.Save(&user)
}

func UpdateUser(db *gorm.DB) {
	var user User
	db.First(&user)
	// 更新单个字段
	db.Model(&user).Update("name", "江湖")
	// 更新多个字段
	db.Model(&user).Updates(User{Name: "十年", Age: 0})
}

func UpdatesMapUser(db *gorm.DB) {
	var user User
	db.First(&user)
	db.Model(&user).Updates(map[string]any{"name": "十年", "age": 0})
}

func UpdateExprUser(db *gorm.DB) {
	var user User
	db.First(&user)
	db.Model(&user).Update("age", gorm.Expr("age + ?", 1))
}

func DeleteUser(db *gorm.DB) {
	var user User
	db.First(&user)
	db.Where("name = ?", "十年").Delete(&user)

	var deletedUser User
	db.Unscoped().Where("name = ?", "十年").First(&deletedUser)
	fmt.Printf("deletedUser: %+v\n", deletedUser)

	db.Unscoped().Where("name = ?", "十年").Delete(&user)

	var deletedUser2 User
	db.Unscoped().Where("name = ?", "十年").First(&deletedUser2)
	fmt.Printf("deletedUser2: %+v\n", deletedUser2)
}

type Post struct {
	gorm.Model
	Title    string     `gorm:"column:title"`
	Content  string     `gorm:"column:content"`
	Comments []*Comment `gorm:"foreignKey:PostID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;references:ID"`
	Tags     []*Tag     `gorm:"many2many:posts_tags"`
}

type Comment struct {
	gorm.Model
	Content string `gorm:"column:content"`
	PostID  uint   `gorm:"column:post_id"`
	Post    *Post
}

type Tag struct {
	gorm.Model
	Name string  `gorm:"column:name"`
	Post []*Post `gorm:"many2many:posts_tags"`
}

func CreatePost(db *gorm.DB) {
	post := Post{
		Title:   "post1",
		Content: "content1",
		Comments: []*Comment{
			{Content: "comment1"},
			{Content: "comment2"},
		},
		Tags: []*Tag{
			{Name: "tag1"},
			{Name: "tag2"},
		},
	}
	db.Create(&post)
}

func AssociationComments(db *gorm.DB) {
	var (
		post     Post
		comments []*Comment
	)

	post.ID = 1
	db.Model(&post).Association("Comments").Find(&comments)
	fmt.Printf("comments: %+v\n", comments)
}

func PreloadTags(db *gorm.DB) {
	var post Post
	db.Preload("Comments").Preload("Tags").First(&post)
	fmt.Printf("post: %+v\n", post)
}

func Joins(db *gorm.DB) {
	type PostComment struct {
		Title   string
		Comment string
	}

	var result PostComment
	db.Model(&Post{}).
		Select("posts.title, comments.content AS comment").
		Joins("LEFT JOIN comments ON comments.post_id = posts.id").
		Where("posts.id = ?", 1).
		Scan(&result)
	fmt.Printf("result: %+v\n", result)
}

func ReplaceComments(db *gorm.DB) {
	var post Post
	db.First(&post)

	newComment := Comment{Content: "Updated Comment"}
	err := db.Model(&post).Association("Comments").Replace([]*Comment{&newComment})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("post: %+v\n", post)
}

func DeletePost(db *gorm.DB) {
	var post Post
	db.Preload("Comments").Preload("Tags").First(&post)
	// db.Where("id = ?", 1).Delete(&post)
	err := db.Model(&post).Association("Comments").Delete(post.Comments)
	if err != nil {
		log.Fatal(err)
	}
}

func Transaction(db *gorm.DB) {
	err := db.Transaction(func(tx *gorm.DB) error {
		post := Post{
			Title: "Hello World",
		}
		if err := tx.Create(&post).Error; err != nil {
			return err
		}
		comment := Comment{
			Content: "你好，世界！",
			PostID:  post.ID,
		}
		return tx.Create(&comment).Error
	})
	if err != nil {
		log.Fatal(err)
	}
}

func BeginTransaction(db *gorm.DB) error {
	tx := db.Begin()    // 开启事务
	defer tx.Rollback() // 回滚事务

	post := Post{
		Title: "Hello World",
	}
	if err := tx.Create(&post).Error; err != nil {
		return err
	}
	comment := Comment{
		Content: "你好，世界！",
		PostID:  post.ID,
	}
	if err := tx.Create(&comment).Error; err != nil {
		return err
	}

	return tx.Commit().Error // 提交事务
}

func RawSQL(db *gorm.DB) {
	type UserResult struct {
		ID   uint
		Name string
		Age  uint8
	}

	var userRes UserResult
	db.Raw(`SELECT id, name, age FROM user WHERE id = ?`, 3).Scan(&userRes)
	fmt.Printf("affected rows: %d\n", db.RowsAffected)
	fmt.Printf("error: %+v\n", db.Error)
	fmt.Printf("userRes: %+v\n", userRes)
}

func ExprSQL(db *gorm.DB) {
	var sumage int
	db.Raw(`SELECT SUM(age) as sumage FROM user WHERE member_number ?`, gorm.Expr("IS NULL")).Scan(&sumage)
	fmt.Printf("sumage: %d\n", sumage)
}

func ExecSQL(db *gorm.DB) {
	// 批量更新
	db.Exec("UPDATE user SET age = ? WHERE id IN ?", 18, []int64{1, 2})
	// 使用表达式更新
	db.Exec(`UPDATE user SET age = ? WHERE name = ?`, gorm.Expr("age * ? + ?", 1, 2), "Jianghu")
	// 删除表
	db.Exec("DROP TABLE user")
	// 创建表
	db.Exec(`
CREATE TABLE user (
  id            int(11) NOT NULL AUTO_INCREMENT,
  name          varchar(50)           DEFAULT '' COMMENT '用户名',
  email         varchar(255) NOT NULL DEFAULT '' COMMENT '邮箱',
  age           tinyint(4) NOT NULL DEFAULT '0' COMMENT '年龄',
  birthday      datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '生日',
  member_number varchar(50) COMMENT '成员编号',
  activated_at  datetime COMMENT '激活时间',
  created_at    datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at    datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at    datetime,
  PRIMARY KEY (id),
  UNIQUE KEY u_email (email),
  INDEX idx_deleted_at(deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='用户表';
`)
	// 插入
	db.Exec(`
INSERT INTO user (id, name, age) VALUES (?, ?, ?)
`, 1, "Jianghu", 18)
}

func Debug(db *gorm.DB) {
	u := User{}
	db.Debug().First(&u)
	fmt.Printf("u: %+v\n", u)
}

func DryRun(db *gorm.DB) {
	var user User
	stmt := db.Session(&gorm.Session{DryRun: true}).First(&user, 1).Statement
	fmt.Printf("user: %+v\n", user)
	fmt.Printf("stmt: %+v\n", stmt.SQL.String())
	fmt.Printf("stmt: %+v\n", stmt.Vars)
}

func main() {
	db, err := ConnectMySQL("127.0.0.1", "3306", "root", "password", "demo")
	if err != nil {
		log.Fatal(err)
	}

	err = SetConnect(db)
	if err != nil {
		log.Fatal(err)
	}

	CreateUser(db)
	CreateUsers(db)

	firstUser, lastUser := FirstLastUser(db)
	fmt.Printf("firstUser: %+v, lastUser: %+v\n", firstUser, lastUser)

	findUsers := WhereFindUsers(db)
	fmt.Printf("findUsers: %+v\n", findUsers)

	users := SelectOrderLimitOffsetUsers(db)
	fmt.Printf("users: %+v\n", users)

	countUsers := CountUsers(db)
	fmt.Printf("countUsers: %+v\n", countUsers)

	groupHaving := GroupHaving(db)
	fmt.Printf("groupHaving: %+v\n", groupHaving)

	SaveUser(db)
	UpdateUser(db)
	UpdatesMapUser(db)
	UpdateExprUser(db)

	DeleteUser(db)

	CreatePost(db)

	AssociationComments(db)

	PreloadTags(db)

	Joins(db)

	ReplaceComments(db)

	DeletePost(db)

	Transaction(db)

	err = BeginTransaction(db)
	if err != nil {
		log.Fatal(err)
	}

	RawSQL(db)
	ExprSQL(db)
	ExecSQL(db)

	Debug(db)
	DryRun(db)
}
