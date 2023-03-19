package base

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

var MysqlDB *sql.DB

func init() {
	var err error
	//打开sqlite数据库
	MysqlDB, err = sql.Open("mysql", "root:root@tcp(127.0.0.1:13306)/hero_story?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}

	MysqlDB.SetMaxOpenConns(128)
	MysqlDB.SetMaxIdleConns(16)
	MysqlDB.SetConnMaxLifetime(2 * time.Minute)

	if err = MysqlDB.Ping(); err != nil {
		panic(err)
	}

}
