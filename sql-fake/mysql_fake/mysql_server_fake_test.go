package mysqlfake

// This is an example of how to implement a MySQL server.
// After running the example, you may connect to it using the following:
// > mysql --host=localhost --port=3306 --user=root mydb --execute="SELECT * FROM mytable;"
// The included MySQL client is used in this example, however any MySQL-compatible client will work.

import (
	"context"
	"fmt"
	"testing"
	"time"

	sqle "github.com/dolthub/go-mysql-server"
	"github.com/dolthub/go-mysql-server/memory"
	"github.com/dolthub/go-mysql-server/server"
	"github.com/dolthub/go-mysql-server/sql"
	"github.com/dolthub/go-mysql-server/sql/information_schema"
	"github.com/dolthub/go-mysql-server/sql/types"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 内存数据库-mysql
// 主要是用来创建数据库，表，如果要查询时，还需插入一定的数据。
func createTestDataBase() *memory.Database {
	fmt.Println("run create database for test.")

	dbName := "demo_fake"
	tabName := "test_fake_demo_tab"

	db := memory.NewDatabase(dbName)
	tab := memory.NewTable(tabName, sql.NewPrimaryKeySchema(
		sql.Schema{
			{Name: "name", Type: types.Text, Nullable: false, Source: tabName},
			{Name: "email", Type: types.Text, Nullable: false, Source: tabName},
			{Name: "created_at", Type: types.Int64, Nullable: false, Source: tabName},
		}), db.GetForeignKeyCollection())

	//{Name: "phone_numbers", Type: types.JSON, Nullable: false, Source: tabName},

	db.AddTable(tabName, tab)

	ctx := sql.NewContext(context.Background())
	tab.Insert(ctx, sql.NewRow("sim", "aaa@126.com", time.Now().UnixMilli()))  //types.MustJSON(`[11111,22222]`),
	tab.Insert(ctx, sql.NewRow("alen", "bbb@126.com", time.Now().UnixMilli())) //types.MustJSON(`[4444,3333]`),

	return db
}

func runMysqlServer(ctx context.Context, t *testing.T) {
	engine := sqle.NewDefault(
		sql.NewDatabaseProvider(
			createTestDataBase(),
			information_schema.NewInformationSchemaDatabase(),
		))

	cfg := server.Config{
		Protocol: "tcp",
		Address:  fmt.Sprintf("%s:%d", "127.0.0.1", 3306),
	}

	s, e := server.NewDefaultServer(cfg, engine)
	if e != nil {
		t.Logf("create fail, e: %v", e)
		return
	}

	if e := s.Start(); e != nil {
		t.Logf("run mysql server fail, e: %v", e)
		return
	}

	t.Logf("start mysql fake server ok, then you can run cmd:  mysql -h 127.0.0.1  demo_fake --execute=\"SELECT * FROM test_fake_demo_tab\"   ")

	<-ctx.Done()
	t.Logf("now stop mysql server....")
	s.Close()

}
func runMysqlClient(cancel context.CancelFunc, t *testing.T) {
	defer cancel()
	//
	dsn := "root:123@tcp(127.0.0.1:3306)/demo_fake?charset=utf8mb4&loc=Local" //&parseTime=false
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Errorf("connect local mysql server fail, e: %v\n", err)
		return
	}

	type DemoModel struct {
		Name      string `gorm:"name"`
		Email     string `gorm:"email"`
		CreatedAt int64  `gorm:"created_at"`
		//PhoneNumbers []string  `gorm:"phone_numbers"`
	}
	var r []*DemoModel

	e := db.Table("test_fake_demo_tab").Find(&r).Error
	if e != nil {
		t.Errorf("query fail, e: %v", e)
		return
	}

	for _, v := range r {
		if v == nil {
			continue
		}
		t.Logf("item: %+v", v)
	}
}

func TestMysqlSrvFake(t *testing.T) {
	ctx, cf := context.WithCancel(context.Background())

	go runMysqlServer(ctx, t)

	time.Sleep(1 * time.Second)
	_ = cf
	runMysqlClient(cf, t)
}
