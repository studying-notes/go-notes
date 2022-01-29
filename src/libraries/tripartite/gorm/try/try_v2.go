/*
 * @Date: 2021.07.23 15:20
 * @Description: Omit
 * @LastEditors: Rustle Karl
 * @LastEditTime: 2021.07.23 15:20
 */

package main

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Journal struct {
	gorm.Model
	Content string
	Subject int8
	Images  string
}

func main() {
	db, err := gorm.Open(sqlite.Open("journal.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&Journal{})
}
