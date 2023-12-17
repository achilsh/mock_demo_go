package sqlitedemo

import (
	// "go-micro.dev/v4/logger"
	"github.com/marmotedu/log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	sql_fake "mock_demo/sql-fake"
)

type UserModel struct {
	//Id   int64  `gorm:"id;primaryKey"`
	Name string `gorm:"name"`
}

func (UserModel) TableName() string {
	return "user"
}
func SqliteDemo() {
	// gorm.// import "gorm.io/driver/mysql"
	// // refer: https://gorm.io/docs/connecting_to_the_database.html#MySQL
	// dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	// db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	// gorm.// import "gorm.io/driver/sqlite"
	// ref: https://gorm.io/docs/connecting_to_the_database.html#SQLite

	db, err := gorm.Open(sqlite.Open("sqlite.db"), &gorm.Config{})
	if err != nil {
		log.Errorf("open sqlite.db fail, e: %v", err)
		return
	}
	defer func() {
		//
	}()
	// 参考：https://juejin.cn/post/7131661977310461965
	// 
	//  AutoMigrate 会创建表、缺失的外键、约束、列和索引。 如果大小、精度、是否为空可以更改，则 AutoMigrate 会改变列的类型。
	//  出于保护您数据的目的，它 不会 删除未使用的列
	if !db.Migrator().HasTable(UserModel{}.TableName()) {
		log.Infof("has not table: %v", UserModel{}.TableName())
		if e := db.Migrator().CreateTable(&UserModel{}); e != nil {
			log.Errorf("create table fail, e: %v", e)
			return
		}
		log.Infof("now create table: %v", UserModel{}.TableName())
	} else {
		log.Infof("has created table: %v", UserModel{}.TableName())

	}

	user := &UserModel{
		Name: "shenzhen",
	}
	if e := db.Create(user).Error; e != nil {
		log.Errorf("insert new item fail, e: %v", e)
		return
	}
}

func InitBusDB() {
	db, err := gorm.Open(sqlite.Open("sqlite.db"), &gorm.Config{})
	if err != nil {
		log.Errorf("open sqlite.db fail, e: %v", err)
		return
	}
	sql_fake.DBHandle = db 
}

func RunDBInProdEnv() {

	if !sql_fake.DBHandle.Migrator().HasTable(UserModel{}.TableName()) {
		log.Infof("has not table: %v", UserModel{}.TableName())
		if e := sql_fake.DBHandle.Migrator().CreateTable(&UserModel{}); e != nil {
			log.Errorf("create table fail, e: %v", e)
			return
		}
		log.Infof("now create table: %v", UserModel{}.TableName())
	} else {
		log.Infof("has created table: %v", UserModel{}.TableName())

	}

	user := &UserModel{
		Name: "beijing",
	}
	if e := sql_fake.DBHandle.Create(user).Error; e != nil {
		log.Errorf("insert new item fail, e: %v", e)
		return
	}
}
