package main

import (
	"database/sql"
	"fmt"
	"net/http"
_ "github.com/go-sql-driver/mysql"
)


func display(w http.ResponseWriter,r * http.Request){

	
	fmt.Fprint(w," data is sucessfully displayed on terminal")
}


func main() {
	
	// create a database object which can be used
	// to connect with database.
	db, err := sql.Open("mysql", "root:qweasdzxc1@tcp(0.0.0.0:3306)/customer")
	
if err !=nil{
	panic(err)
}

err = db.Ping()

if err!=nil{
	panic(err)
}

result,err:=db.Query("select * from users")

if err !=nil{
	panic(err)	
}

for result.Next(){

	
	var id int 
	var name string
	var age  int

	err= result.Scan(&id,&name,&age)

	if err !=nil{
		panic(err)
	}

fmt.Printf("id: %d ,name: %s,age: %d \n",id,name,age )
	
	
}

http.HandleFunc("/",display)
fmt.Print("Server is running on 3000")
http.ListenAndServe(":3000",nil)


}