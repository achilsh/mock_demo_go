package sqlite_fake

import (
	"fmt"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	sql_fake "mock_demo/sql-fake"
	sqlitedemo "mock_demo/sql-fake/sqlite_demo"
)

//在单元测试的时候封装一个接口，该接口就是创建一个新的的db handle 初始化全局句柄。
//在单元测试开始时调用该接口即可。参考： https://juejin.cn/post/7131661977310461965
func init_test_db(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		fmt.Println("open fail.e: ", err)
		return
	}
	sql_fake.DBHandle = db

	if !sql_fake.DBHandle.Migrator().HasTable(sqlitedemo.UserModel{}.TableName()) {
		t.Logf("has not table: %v", sqlitedemo.UserModel{}.TableName())
		if e := sql_fake.DBHandle.Migrator().CreateTable(&sqlitedemo.UserModel{}); e != nil {
			t.Errorf("create table fail, e: %v", e)
			return
		}
	}
}
func TestSqliteFake(t *testing.T) {
	init_test_db(t)
	//把 db name 文件 替换成  "file::memory:?cache=shared" ；这样操作数据库就是内存数据库。
	//就不会对数据文件产生污染，同时也可以验证我们的sql语句是否正确。

	u := sqlitedemo.UserModel{
		Name: "nanchang",
	}
	e := sql_fake.DBHandle.Create(&u).Error
	if e != nil {
		t.Errorf("create item name: shenzhen fail, e: %v", e)
		return
	}

	r := []*sqlitedemo.UserModel{}
	e = sql_fake.DBHandle.Find(&r).Error
	if e != nil {
		t.Errorf("query all user model fail, e: %v", e)
		return
	}
	for _, v := range r {
		if v == nil {
			continue
		}
		t.Logf("data: %v", *v)
	}

}
