package sqlite_fake

import (
	"testing"

	sql_fake "mock_demo/sql-fake"
	sqlitedemo "mock_demo/sql-fake/sqlite_demo"
)

func TestRunQuery(t *testing.T) {
	init_test_db(t)

	h := &sqlitedemo.UserModel{
		Name: "hanzhou",
	}
	if e := sql_fake.DBHandle.Create(h).Error; e != nil {
		t.Errorf("insert new fail, e: %v", e)
		return
	}

	r := []*sqlitedemo.UserModel{}
	if e := sql_fake.DBHandle.Find(&r, "name='hanzhou'").Error; e != nil {
		t.Logf("seach cond fail, e: %v", e)
		return
	}

	t.Logf("%v", len(r))
	for _, v := range r {
		t.Logf("%#v", v)
	}
}
