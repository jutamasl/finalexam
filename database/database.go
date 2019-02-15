package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var db *sql.DB

func Conn() *sql.DB {
	if db != nil {
		fmt.Println("olddddd")
		return db
	}
	var err error
	db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	fmt.Println("open connection")
	if err != nil {
		log.Fatal("cannot connect db", err)
	}
	fmt.Println("connect db..")
	return db
}
