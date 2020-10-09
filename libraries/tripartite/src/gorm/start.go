package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"log"
)

// gorm.Model 官方定义的通用模型
//type Model struct {
//	ID        uint `gorm:"primary_key"`
//	CreatedAt time.Time
//	UpdatedAt time.Time
//	DeletedAt *time.Time `sql:"index"`
//}

type Product struct {
	gorm.Model // 官方定义的通用模型
	Code       string
	Price      uint
}

func main() {
	db, err := gorm.Open("sqlite3", "gorm_test.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// AutoMigrate run auto migration for given models,
	// will only add missing fields, won't delete/change current data
	db.AutoMigrate(&Product{})

	// 插入一条数据
	db.Create(&Product{Code: "SN2020001", Price: 122})
	db.Create(&Product{Code: "SN2020002", Price: 123})

	// 查询一条数据
	var product Product

	// find product with id 1
	db.First(&product, 1)
	//fmt.Printf("%+v\n", product)

	// find product with code SN2020001
	db.First(&product, "code = ?", "SN2020001")
	//fmt.Printf("%+v\n", product)

	// 查询多条数据，标记已删除的默认无法获取
	var products []Product
	db.Find(&products, "code = ?", "SN2020")
	//fmt.Printf("%+v\n", products)

	// Update - update product's price to 2000
	db.Model(&product).Update("Price", 2000)
	// product price 字段同步更新
	fmt.Printf("%+v\n", product)

	// Delete - delete product
	db.Delete(&product)
	// product DeletedAt 字段不会更新
	fmt.Printf("%+v\n", product)
}
