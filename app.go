package main

import ("database/sql"
_"github.com/go-sql-driver/mysql"
"fmt"
)
func main(){
	db,err:=sql.Open("mysql","root:qweasdzxc1@tcp(127.0.0.1:3306)/customer")
   if err !=nil{
		panic(err)
	}
	defer db.Close()
	res,err:=db.Query ("update  users set name='Aski' where name='seee' ")

if err !=nil{
	panic(err)
}
	defer res.Close()
fmt.Print("Mysql is connected and added the values")
}
