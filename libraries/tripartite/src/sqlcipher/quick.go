package main

import (
	"fmt"
	_ "github.com/fujiawei-dev/go-sqlcipher"
	//"database/sql"
	sql "github.com/jmoiron/sqlx"
	"log"
)

func main() {
	db, err := sql.Open("sqlite3", "users.db?_key=LGaMU5RdVImAL9CN")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 必须写入一个数据，否则仅仅只是创建了一个未加密的空白文件
	c := "CREATE TABLE IF NOT EXISTS `users` (`id` INTEGER PRIMARY KEY, `name` char, `password` chart, UNIQUE(`name`));"
	_, err = db.Exec(c)
	if err != nil {
		fmt.Println(err)
		return
	}
	d := "INSERT INTO `users` (name, password) values('py', 314159);"
	_, err = db.Exec(d)
	if err != nil {
		fmt.Println(err)
	}

	e := "select name, password from users;"
	rows, err := db.Query(e)
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()

	var (
		name     string
		password string
	)
	for rows.Next() {
		_ = rows.Scan(&name, &password)
		fmt.Print("{\"name\":\"" + name + "\", \"password\": \"" + password + "\"}")
	}
	rows.Close()
}
