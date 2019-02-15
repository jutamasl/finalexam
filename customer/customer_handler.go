package customer

import (
	"fmt"
	"log"
	"net/http"

	//"database/sql"
	"github.com/gin-gonic/gin"

	"strconv"

	"github.com/jutamasl/finalexam/database"
)

var customer1 = []Customer{}

func CreateTable() {
	fmt.Println("start func create table")
	ctb := ` 
	CREATE TABLE IF NOT EXISTS customers (
		id SERIAL PRIMARY KEY,
		name TEXT,
		email TEXT,
		status TEXT
	);`
	_, err := database.Conn().Exec(ctb)
	if err != nil {
		fmt.Println("cannot crate table", err)
	}
	fmt.Println("create table success")
}

func authMiddleware(c *gin.Context) { // HL
	log.Println("start middleware")
	authKey := c.GetHeader("Authorization") // HL
	if authKey != "token2019" {
		c.JSON(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized)) // HL
		c.Abort()                                                                 // HL
		return
	}

	c.Next()
	log.Println("end middleware")
}

func NewRouter() *gin.Engine {
	r := gin.Default()
	r.Use(authMiddleware)
	v1 := r.Group("")
	v1.GET("/customers", getCustomerHandler)
	v1.POST("/customers", createCustomerHandler)
	v1.GET("/customers/:id", getCustomerByIDHandler)
	v1.PUT("/customers/:id", updateCustomerHandler)
	v1.DELETE("/customers/:id", deleteCustomerHandler)
	return r
}

func getCustomerHandler(c *gin.Context) {
	stmt, err := database.Conn().Prepare("select id, name,email,status from customers")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	rows, err := stmt.Query()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	var customer1 = []Customer{}
	fmt.Printf("c1 %v \n", customer1)
	for rows.Next() {
		t := Customer{}
		err := rows.Scan(&t.ID, &t.Name, &t.Email, &t.Status)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		customer1 = append(customer1, t)
		fmt.Printf("c1 %v \n", customer1)
	}
	c.JSON(http.StatusOK, customer1)
}

var id = 0

func createCustomerHandler(c *gin.Context) {
	var item Customer

	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//fmt.Printf("ti = %v, st = %v \n", item.Title, item.Status)
	//row := database.InsertCust(item.Name, item.Email, item.Status)
	row := database.Conn().QueryRow("insert into customers (name, email, status) values ($1, $2, $3) returning id", item.Name, item.Email, item.Status)
	var id int
	var err error
	err = row.Scan(&id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"scan message": err.Error()})
		return
	}
	item.ID = id
	c.JSON(http.StatusCreated, item)
}

func getCustomerByIDHandler(c *gin.Context) {
	stmt, err := database.Conn().Prepare("select id, name, email,status from customers where id=$1;")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	rid, _ := strconv.Atoi(c.Param("id"))
	row := stmt.QueryRow(rid)
	t := Customer{}
	err = row.Scan(&t.ID, &t.Name, &t.Email, &t.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, t)
}
func updateCustomerHandler(c *gin.Context) {
	var item Customer
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rid, _ := strconv.Atoi(c.Param("id"))
	//fmt.Println("update", rid)
	stmt, err := database.Conn().Prepare("update customers set name=$2, email=$3, status=$4 where id=$1;")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"prep message": err.Error()})
		return
	}

	if _, err := stmt.Exec(rid, item.Name, item.Email, item.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"exec message": err.Error()})
		return
	}
	fmt.Printf("id %v \n", rid)
	item.ID = rid
	c.JSON(http.StatusOK, item)
}

func deleteCustomerHandler(c *gin.Context) {
	rid, _ := strconv.Atoi(c.Param("id"))
	fmt.Println("input", rid)
	stmt, err := database.Conn().Prepare("delete from customers where id=$1;")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"prep message": err.Error()})
		return
	}

	if _, err := stmt.Exec(rid); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"exec message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "customer deleted"})
}
