package main

import sqlitedemo "mock_demo/sql-fake/sqlite_demo"

func main() {
	sqlitedemo.InitBusDB()
	sqlitedemo.RunDBInProdEnv()
	

	// sqlitedemo.SqliteDemo()
}