package mysqlx

import (
	"database/sql"
	"fmt"
	"myinfo"
	_ "mysql"
	"time"
)

//var db *sql.DB //数据库

//**************************************数据库查询函数*******************************
/*定义一个函数类*/
func MySql_Init(sqluser, sqlpass, sqlserver, sqlport, sqlbase string) (db *sql.DB, ret int) {
	ret = 1
	if db != nil {
		return db, 0
	}
	var err error
	//db, err = sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/test?charset=utf8&allowOldPasswords=1")
	//db, err = sql.Open("mysql", "root:root!@#@tcp(47.96.82.26:82)/square?charset=utf8&allowOldPasswords=1")
	db, err = sql.Open("mysql", sqluser+":"+sqlpass+"@tcp("+sqlserver+":"+sqlport+")/"+sqlbase+"?charset=utf8&allowOldPasswords=1")
	defer db.Close() //关闭

	if err != nil {
		myinfo.SaveLogEx("MySql_Init:执行失败\r\n")
		return nil, ret
	} else {
		myinfo.SaveLogEx("MySql_Init:初始化成功\r\n")
	}
	db.SetConnMaxLifetime(time.Second * 25)
	db.SetMaxIdleConns(100)
	db.SetMaxOpenConns(100)
	bRet := show_tables(db)
	if bRet {
		return db, 0
	}
	return db, 1
}

func show_tables(db *sql.DB) bool {
	if db == nil {
		return false
	}
	rows, err := db.Query("show tables")
	defer rows.Close()

	if err != nil {
		fmt.Println("Query:Error")
		return false
	}

	var Num int
	Num = 0
	for rows.Next() {
		Num++
	}
	return (Num > 0)
}
