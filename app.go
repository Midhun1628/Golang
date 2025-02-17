package main

import(
   "net/http"
   "database/sql"
   "fmt"
   "encoding/json"
   "log"
_ "github.com/go-sql-driver/mysql"

)
var db *sql.DB 

type User struct{
    ID int      `json:"id"`
    NAME string `json:"name"`
    AGE  int    `json:"age"`
}


func Display(w http.ResponseWriter,r *http.Request){

    var err error

    result,err:=db.Query("Select id,name,age from users")
    
    if err  !=nil{
        log.Fatal("Failed to fetch data",err)    
        http.Error(w,"Internal server error",http.StatusInternalServerError)

    }
     defer result.Close()

     var users []User
    
    for result.Next(){
        var user User
        err:=result.Scan(&user.ID,&user.NAME,&user.AGE)
        if err!=nil{
            log.Fatal("Failed to scan data",err)      
            http.Error(w,"Failed to scan data",http.StatusInternalServerError)
        }
        users=append(users,user)
        }
        w.Header().Set("Content-Type","application/json")
        json.NewEncoder(w).Encode(users)
    
    }
    func main(){ 

var err error

        db,err=sql.Open("mysql","root:qweasdzxc1@tcp(0.0.0.0:3306)/customer")
    
        if err !=nil{
            panic(err)
        }
        
        err= db.Ping()
        
        if err!=nil{
            log.Fatal("Connecting to db failed",err)
        }

        fmt.Println("Database connected")        
        defer db.Close()
        
        http.HandleFunc("/",Display)
        fmt.Println("Sever is running in 3000 port")
        http.ListenAndServe(":3000",nil)
        
    }