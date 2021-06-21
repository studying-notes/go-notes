package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
)

func main() {
	_ = os.Remove("./foo.db")

	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 创建数据表
	sqlStmt := `
		CREATE TABLE IF NOT EXISTS userinfo
		(
			uid        INTEGER PRIMARY KEY AUTOINCREMENT,
			username   VARCHAR(64) NULL,
			department VARCHAR(64) NULL,
			created    DATE        NULL
		);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}

	// 事务插入
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare("INSERT INTO userinfo(username, department, created) values(?,?,?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	for i := 0; i < 100; i++ {
		_, err = stmt.Exec("Zhu", "Yue", fmt.Sprintf("VOL.%03d", i))
		if err != nil {
			log.Fatal(err)
		}
	}
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}

	// 查询整个表
	rows, err := db.Query("select username, department from userinfo")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var username string
		var department string
		err = rows.Scan(&username, &department)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(username, department)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	// 查询指定条件
	stmt, err = db.Prepare("select created from userinfo where uid = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	var created string
	err = stmt.QueryRow("3").Scan(&created)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(created)

	// 删除
	//_, err = db.Exec("delete from userinfo")
	//if err != nil {
	//	log.Fatal(err)
	//}

	// 插入
	_, err = db.Exec("INSERT INTO userinfo(username, department, created) values('foo', 'bar', 'baz'), ('foo2', 'bar2', 'baz2')")
	if err != nil {
		log.Fatal(err)
	}
}
