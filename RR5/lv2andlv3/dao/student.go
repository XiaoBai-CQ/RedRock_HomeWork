package dao

import (
	"database/sql"
	"log"
)

func InitDB() (db *sql.DB) {
	//students表中的users
	//密码不给看哼哼
	dsn := "root:********@tcp(127.0.0.1:3306)/students"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect:%v", err)
	}
	return db
}
