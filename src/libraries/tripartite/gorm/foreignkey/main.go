/*
 * @Date: 2021.12.25 16:56
 * @Description: Omit
 * @LastEditors: Rustle Karl
 * @LastEditTime: 2021.12.25 16:56
 */

package main

import (
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type User struct {
	gorm.Model
	Name         string
	CompanyRefer int
	Company      Company `gorm:"foreignKey:CompanyRefer"`
}

type Company struct {
	ID   int
	Name string
}

func main() {
	db, err := gorm.Open("sqlite3", "main.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// AutoMigrate run auto migration for given models,
	// will only add missing fields, won't delete/change current data
	db.AutoMigrate(&User{}, &Company{})

	var users []User
	db.Debug().Preload("Company").Find(&users)
}
