package main

import (
	"database/sql"
	//_ "github.com/mattn/go-sqlite3"
	_ "github.com/fujiawei-dev/go-sqlcipher"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	_ = os.Remove("./zhengjiangui2.db")

	db, err := sql.Open("sqlite3", "./zhengjiangui2.db?_key=LGaMU5RdVImAL9CN")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlStmt, _ := ioutil.ReadFile("./w_all.sql")
	_, err = db.Exec(string(sqlStmt))
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}
}
