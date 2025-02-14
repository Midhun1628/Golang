package main

import (
	"fmt"
	"log"
	"database/sql"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	 "github.com/gin-gonic/gin"
)


type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func main() {

	db, err := sql.Open("mysql", "root:qweasdzxc1@tcp(127.0.0.1:3306)/customer")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		rows, err := db.Query("SELECT id, name, age FROM users")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var users []User
		for rows.Next() {
			var user User
			if err := rows.Scan(&user.ID, &user.Name, &user.Age); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			users = append(users, user)
		}

		
		c.JSON(http.StatusOK, users)
	})

	fmt.Println("Server running on 3000")
	r.Run(":3000")
}
