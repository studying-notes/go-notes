package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
)

type Model struct {
	gorm.Model  // 官方定义的通用模型
	Email       string
	Password    string
	PhoneNumber string
	UserName    string
	FirstName   string
	LastName    string
	UUID        string
	Year        string
	Amount      float64
}

func (Model) TableName() string {
	return "users"
}

type SubModel struct {
	gorm.Model
	Code   string
	Price  uint
	UserID uint
}

func (SubModel) TableName() string {
	return "books"
}

func main() {
	//port := "3306"
	port := "4567"
	db, err := gorm.Open("mysql", fmt.Sprintf(
		"root:root@(127.0.0.1:%s)/gorm_test?", port)+
		"charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	db.LogMode(true)

	//db.AutoMigrate(&Model{})
	//db.AutoMigrate(&SubModel{})

	var user Model
	gm := db.First(&user)

	fmt.Println(gm.RowsAffected)
	fmt.Println(gm.Error)
}
