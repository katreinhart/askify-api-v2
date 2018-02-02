package model

import (
	"errors"
	"fmt"
	"os"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"

	// importing postgres dialect for GORM
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// declare DB
var db *gorm.DB

// type declarations for the data model
type (
	// QuestionModel is the DB model for asked questions
	QuestionModel struct {
		gorm.Model
		Question string `json:"question"`
		Answered bool   `json:"answered"`
		UserID   string `json:"userid"`
		FName    string `json:"fname"`
		Cohort   string `json:"cohort"`
	}

	transformedQuestion struct {
		ID       uint   `json:"id"`
		Question string `json:"question"`
		Answered bool   `json:"answered"`
		UserID   int    `json:"userid"`
		FName    string `json:"fname"`
		Cohort   string `json:"cohort"`
	}

	archiveQuestion struct {
		ID       uint            `json:"id"`
		Question string          `json:"question"`
		UserID   int             `json:"userid"`
		FName    string          `json:"fname"`
		Cohort   string          `json:"cohort"`
		Answers  []archiveAnswer `json:"answers"`
	}

	answerModelInput struct {
		Answer string `json:"answer"`
		UserID int    `json:"userid"`
		FName  string `json:"fname"`
		Cohort string `json:"cohort"`
	}

	answerModel struct {
		gorm.Model
		QuestionID int    `json:"questionid"`
		Answer     string `json:"answer"`
		UserID     int    `json:"userid"`
		FName      string `json:"fname"`
		Cohort     string `json:"cohort"`
	}

	transformedAnswer struct {
		ID         uint   `json:"id"`
		QuestionID int    `json:"questionid"`
		Answer     string `json:"answer"`
		UserID     int    `json:"userid"`
		FName      string `json:"fname"`
		Cohort     string `json:"cohort"`
	}

	archiveAnswer struct {
		ID         uint   `json:"id"`
		QuestionID int    `json:"questionid"`
		Answer     string `json:"answer"`
		UserID     int    `json:"uid"`
		FName      string `json:"fname"`
		Cohort     string `json:"cohort"`
	}

	UserModel struct {
		gorm.Model
		Email    string `json:"email"`
		FName    string `json:"fname"`
		Cohort   string `json:"cohort"`
		Password string `json:"password"`
		Admin    bool   `json:"admin"`
	}

	TransformedUser struct {
		ID     uint   `json:"id"`
		Email  string `json:"email"`
		FName  string `json:"fname"`
		Cohort string `json:"cohort"`
		Token  string `json:"token"`
	}

	ListedUser struct {
		ID     uint   `json:"id"`
		FName  string `json:"fname"`
		Cohort string `json:"cohort"`
		Email  string `json:"email"`
		Admin  bool   `json:"admin"`
	}

	// CustomClaims for JWT handling
	CustomClaims struct {
		UID uint `json:"uid"`
		Rol bool `json:"rol"`
		jwt.StandardClaims
	}
)

// Errors
var ErrorUserExists = errors.New("User already exists")

var ErrorNotFound = errors.New("Not found")

var ErrorForbidden = errors.New("Forbidden")

// init function runs at setup; connects to database
func init() {
	_ = godotenv.Load()

	hostname := os.Getenv("HOST")
	dbname := os.Getenv("DBNAME")
	username := os.Getenv("DBUSER")
	password := os.Getenv("PASSWORD")

	dbString := "host=" + hostname + " user=" + username + " dbname=" + dbname + " sslmode=disable password=" + password

	var err error
	db, err = gorm.Open("postgres", dbString)
	if err != nil {
		fmt.Println(err.Error())
		panic("Unable to connect to DB")
	}

	db.AutoMigrate(&QuestionModel{})
	db.AutoMigrate(&answerModel{})
	db.AutoMigrate(&UserModel{})
}
