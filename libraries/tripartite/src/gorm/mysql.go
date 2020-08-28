package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"time"
)

// 定义一个数据模型
// 列名是字段名的蛇形小写，即下划线形式
type User struct {
	Id       uint   `gorm:"AUTO_INCREMENT"`
	Name     string `gorm:"size:50"`
	Age      int    `gorm:"size:3"`
	Birthday *time.Time
	Email    string `gorm:"type:varchar(50);unique_index"`
	PassWord string `gorm:"type:varchar(25)"`
}

func main() {
	db, err := gorm.Open("mysql", "root:root@(127.0.0.1:3306)/grom_test?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 默认的表名都是结构体名称的复数形式，User 结构体默认创建的表为 users
	// db.SingularTable(true) 可以取消表名的复数形式，使表名和结构体名称一致
	db.AutoMigrate(&User{})

	// 添加唯一索引
	db.Model(&User{}).AddUniqueIndex("name_email", "id", "name", "email")

	// 插入记录
	//db.Create(&User{Name: "xlj", Age: 18, Email: "xlj@xlj.org"})

	gm := db.Create(&User{Name: "xll", Age: 18, Email: "xlj@xlj.org"})
	if err := gm.Error; err != nil {// 处理错误
		fmt.Println(err)
	}

	// 查看插入后的全部元素
	var users []User
	db.Find(&users)
	fmt.Printf("%+v\n", users)

	var user User
	// 查询一条记录
	db.First(&user, "name = ?", "xlj")
	fmt.Printf("%+v\n", user)

	// 更新记录
	db.Model(&user).Update("name", "xxj")
	fmt.Printf("%+v\n", user)

	// 删除记录
	db.Delete(&user)
}
