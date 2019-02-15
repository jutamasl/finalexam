package main

import (
	"fmt"

	"github.com/jutamasl/finalexam/customer"
	"github.com/jutamasl/finalexam/database"
)

func main() {
	fmt.Println("start conn")
	database.Conn()
	fmt.Println("start create table...")
	customer.CreateTable()
	r := customer.NewRouter()
	r.Run(":2019")
}
