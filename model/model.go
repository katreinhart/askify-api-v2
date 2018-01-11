package model

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"

	// importing postgres dialect for GORM
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// declare DB
var db *gorm.DB

// type declarations for the data model
type (
	questionModel struct {
		gorm.Model
		Question string `json:"question"`
		Answered bool   `json:"answered"`
		UserID   int    `json:"userid"`
	}

	transformedQuestion struct {
		ID         uint   `json:"id"`
		Question   string `json:"question"`
		Answered   bool   `json:"answered"`
		AnsweredBy string `json:"answered_by"`
	}
)

// init function runs at setup; connects to database
func init() {
	_ = godotenv.Load()

	hostname := os.Getenv("HOST")
	dbname := os.Getenv("DBNAME")
	username := os.Getenv("USER")
	password := os.Getenv("PASSWORD")

	dbString := "host=" + hostname + " user=" + username + " dbname=" + dbname + " sslmode=disable password=" + password

	var err error
	db, err = gorm.Open("postgres", dbString)
	if err != nil {
		fmt.Println(err.Error())
		panic("Unable to connect to DB")
	}

	db.AutoMigrate(&questionModel{})
}
