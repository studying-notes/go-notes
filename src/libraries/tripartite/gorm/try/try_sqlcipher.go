package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	//_ "github.com/jinzhu/gorm/dialects/sqlite"
	_ "github.com/fujiawei-dev/go-sqlcipher"
	"log"
)

type Product struct {
	gorm.Model // 官方定义的通用模型
	Code       string
	Price      uint
}

type Producer struct {
	gorm.Model // 官方定义的通用模型
	Age        int
	Code       string
	//Price      uint
	ProductID int64
}

func main() {
	db, err := gorm.Open("sqlite3", "gorm_test.db?_key=LGaMU5RdVImAL9CN")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.LogMode(true)

	// AutoMigrate run auto migration for given models,
	// will only add missing fields, won't delete/change current data
	db.AutoMigrate(&Product{}, &Producer{})

	// 插入一条数据
	db.Create(&Product{Code: "SN2020001", Price: 122})
	//db.Create(&Product{Code: "SN2020002", Price: 123})
	db.Create(&Producer{Code: "SN2020002", ProductID: 1, Age: 100})

	// 查询一条数据
	//var product = Product{Code: "SN2020001"}

	//product.ID = 1

	// Update - update product's price to 2000
	//db.Model(&product).
	//	Where("code = ?", "SN2020001").
	//	Or("code = ?", "SN2020002").
	//	Where("price > ?", 1).
	//	Update("Price", 2000)

	var product struct {
		Code  string
		Price uint
		Age   int
	}

	db.Model(&Product{}).Update("price", "1")

	//db.Table("products AS p").
	//	Joins("LEFT JOIN producers as q ON p.id=q.product_id").
	//	Select("*").
	//	Limit(1).
	//	Scan(&product)

	db.Model(&Product{}).Scan(&product)
	//db.Model(&Product{}).Find(&product)
	//db.Model(&Product{}).Find(&product)

	// product price 字段同步更新
	fmt.Printf("%+v\n", product)
}
