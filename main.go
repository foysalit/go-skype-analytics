package main

import (
    "github.com/jinzhu/gorm"
    _ "github.com/mattn/go-sqlite3"
    "fmt"
    "os"
    "time"
)

type Contact struct {
	SkypeName 	string  `gorm:"column:skypename"`
	FullName 	string  `gorm:"column:fullname"`
}

type Message struct {
	To 			string  `gorm:"column:dialog_partner"`
	From 		string  `gorm:"column:author"`
	Timestamp 	int64  `gorm:"column:timestamp"`
	Content		string `gorm:"column:body_xml"`
	Type 		int `gorm:"column:chatmsg_type"`
}

func GetDbName () string {
	skypePath := "/home/foysal/.Skype/"
	skypePath += os.Args[1]
	skypePath += "/main.db"
	return skypePath
}

func GetMessages(db gorm.DB) {
	var messages []Message
	partner := os.Args[2]

	if (!HasContact(db, partner)) {
		fmt.Printf("You don't have a contact with username %q", partner)
		return
	}

	db.Table("messages").Where("dialog_partner = ?", partner).Or("author = ?", partner).Where("chatmsg_type = ?", 3).Order("timestamp desc").Limit(30).Find(&messages)
	
	if (len(messages) <= 0) {
		fmt.Println("No records found!")
	} else {
		for _, message := range messages {
			fmt.Println(message.From, message.Content, time.Unix(message.Timestamp, 0))
		}
	}
}

func GetContacts(db gorm.DB) {
	var contacts []Contact
	db.Debug().Table("Contacts").Find(&contacts)
	
	if (len(contacts) <= 0) {
		fmt.Println("No records found!")
	} else {
		for _,contact := range contacts {
			fmt.Println(contact.SkypeName)
		}
	}
}

func HasContact(db gorm.DB, username string) bool {
	count := 0
	db.Table("contacts").Where("skypename = ?", username).Count(&count)
	return (count > 0)
}

func main() {
	if (len(os.Args) < 3) {
		fmt.Println("First parameter must be your username and the Second parameter must be the other user!")
		return
	}

	db, err := gorm.Open("sqlite3", GetDbName())

	if (err != nil) {
		fmt.Println("error")
		return
	}

	db.DB()

	db.DB().Ping()
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)
	db.LogMode(true)

	GetMessages(db)
	// GetContacts(db)
}