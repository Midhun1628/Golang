package main

import "fmt"

type Personal struct{
	name string
	age int
}

var Per1 Personal
var Per2 Personal

func main(){
	Per1.name="Midhun"
	Per1.age=12
	Per2.name="Midhunn"
	Per2.age=13

fmt.Print("Name : " ,Per1.name)
fmt.Print(" Age : " ,Per1.age)
fmt.Print(" Name :" ,Per2.name)
fmt.Print("age : " ,Per2.age)


}

