package database

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"

	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDb() {
	dsn := "root:Singha@12@tcp(127.0.0.1:3306)/Library_0706?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err, " Database Connection Failed")
		log.Fatal("connection error: ", err)
	} else {
		fmt.Println("Successfully connected to the database")
	}
	DB = db
}
