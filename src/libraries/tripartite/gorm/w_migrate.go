package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"time"
)

type ItemProduct struct {
	ID         int64     `gorm:"column:id;primary_key;auto_increment;" json:"id"`
	Name       string    `gorm:"column:name;type:varchar(255)"         json:"name"`
	IsDelete   int       `gorm:"column:is_delete;type:int(11)"         json:"isDelete"`
	UnitID     int64     `gorm:"column:unit_id;type:int(11)"         json:"unitID"`
	CreateTime time.Time `gorm:"column:create_time;type:datetime;not null;default:CURRENT_TIMESTAMP" json:"createTime"`
	UpdateTime time.Time `gorm:"column:update_time;type:datetime;not null;default:CURRENT_TIMESTAMP" json:"updateTime"`
}

func (ItemProduct) TableName() string {
	return "item_product"
}

func main() {
	db, err := gorm.Open("mysql", "wing:wing@tcp(118.178.86.183:14010)/"+
		"asset_cabinet?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.SingularTable(true)
	db.AutoMigrate(&ItemProduct{})
}
