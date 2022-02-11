package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
)

type PacketPerson struct {
	Id       int    `grom:"id" json:"id"`
	IsDelete int    `grom:"is_delete" json:"is_delete"`
	PersonID int    `grom:"ryid"` //人员ID
	PacketID string `grom:"zbh"`  //组包号
}

func (model PacketPerson) TableName() string {
	return "packet_person"
}

type Doc struct {
	Id       int    `grom:"id" json:"id"`
	IsDelete int    `grom:"is_delete" json:"is_delete"`
	PersonId int    `grom:"person_id" json:"person_id"`
	Zjbs     string `grom:"zjbs" json:"zjbs"`
	Zjhm     string `grom:"zjhm" json:"zjhm"`
	Zjlx     string `grom:"zjlx" json:"zjlx"`
}

func (model Doc) TableName() string {
	return "doc"
}

type Person struct {
	Bbzt     string `grom:"bbzt" json:"bbzt"`
	Dwdm     string `grom:"dwdm" json:"dwdm"`
	Id       int    `grom:"id" json:"id"`
	IsDelete int    `grom:"is_delete" json:"is_delete"`
	Rybs     string `grom:"rybs" json:"rybs"`
}

func (model Person) TableName() string {
	return "person"
}

func main() {
	db, err := gorm.Open("mysql", "wing:wing@(118.178.86.183:14010)/zjg_6s?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	db.LogMode(true)

	var packets []Person
	gm := db.Where("is_delete = ?", 0).Find(&packets)
	fmt.Println(gm.Error)
	fmt.Println(gm.RowsAffected)
	fmt.Println(packets)
}
