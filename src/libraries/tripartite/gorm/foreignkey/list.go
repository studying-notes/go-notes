package main

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/satori/go.uuid"
)

type Email struct {
	ID      uuid.UUID `gorm:"primary_key;type:char(36);"`
	Address string
	Uid     uuid.UUID
}

type User struct {
	ID     uuid.UUID `gorm:"primary_key;type:char(36);"`
	Name   string
	Emails []Email `gorm:"ForeignKey:Uid"`
}

func main() {
	// Connect
	db, err := gorm.Open("sqlite3", "ForeignKey.db")
	if err != nil {
		fmt.Println("failed to connect database")
	}
	defer db.Close()

	// Migrate
	db.AutoMigrate(
		&User{},
		&Email{},
	)

	/*
	   id, err := uuid.FromString("0ff7a01f-9ab6-4041-b372-4375ce3d4065")
	   // Create user
	   db.Debug().Create(
	       &User{
	           ID:id,
	           Name: "admin",
	           Emails: []Email{
	               {
	                   ID:uuid.Must(uuid.NewV4()),
	                   Address: "admin@admin.com",
	                   Uid: id,
	               },
	               {
	                   ID:uuid.Must(uuid.NewV4()),
	                   Address: "test@admin.com",
	                   Uid: id,
	               },
	           },
	       },
	   )

	   /*
	   // Get user id
	   id, err := uuid.FromString("0ff7a01f-9ab6-4041-b372-4375ce3d4065")
	   if err != nil {
	       fmt.Println("parse uuid failed")
	   }

	   // Create emails
	   emails := []Email {
	       {
	           ID:uuid.Must(uuid.NewV4()),
	           Address: "admin@admin.com",
	           Uid: id,
	       },
	       {
	           ID:uuid.Must(uuid.NewV4()),
	           Address: "test@admin.com",
	           Uid: id,
	       },
	   }
	   fmt.Println(emails)
	   for _, email := range emails{
	       db.Debug().Create(email)
	   }
	*/

	var user User
	err = db.Debug().Where(&User{Name: "admin"}).First(&user).Error
	if err != nil {
		fmt.Println("Query failed")
	}
	fmt.Println(user)

	var emails []Email
	err = db.Debug().Model(&user).Related(&emails, "Emails").Error
	if err != nil {
		fmt.Println("Query failed")
	}
	fmt.Println(emails)

	user.Emails = emails
	fmt.Println(user)

	/*
	   var email Email
	   err = db.Debug().Where(&Email{Address: "admin@admin.com"}).First(&email).Error
	   if err != nil {
	       fmt.Println("Query failed")
	   }
	   fmt.Println(email)
	*/
}
