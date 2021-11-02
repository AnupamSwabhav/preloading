package main

import (
	"log"
	"preloading/test/address"
	"preloading/test/controller"
	"preloading/test/repository"
	"preloading/test/service"
	"preloading/test/student"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB
var err error

func main() {

	db, err = gorm.Open("mysql", "root:password@tcp(127.0.0.1:3306)/httpcon?charset=utf8&parseTime=True")

	if err != nil {
		log.Println("Connection Failed to Open")
	} else {
		log.Println("Connection Established")
	}
	db.AutoMigrate(&student.Student{})
	db.AutoMigrate(&address.Address{})
	db.Model(&address.Address{}).AddForeignKey("student_id", "students(id)", "CASCADE", "CASCADE")
	repo := repository.NewRepository()
	ser := service.NewStudentService(db, repo)
	con := controller.New(ser)
	con.HandleRequest()

}
