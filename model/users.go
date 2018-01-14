package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

// CreateUser checks to see uniqueness of new id, saves the user info if unique, returns error otherwise.
func CreateUser(b []byte) ([]byte, error) {
	var user, dbUser userModel

	err := json.Unmarshal(b, &user)

	if err != nil {
		// handle error
		return []byte(""), err
	}

	db.First(&dbUser, "email = ?", user.Email)
	if dbUser.ID != 0 {
		return []byte("{\"message\": \"User already exists in DB.\"}"), errors.New("User already exists")
	}

	hash, err := hashPassword(user.Password)
	if err != nil {
		// handle internal server error
		return []byte("{\"message\": \"Sorry, something went wrong.\"}"), err
	}

	user.Password = hash

	// save the user in the DB
	db.Save(&user)
	// Get the user back from the database so you have the ID
	db.First(&dbUser, "email = ?", user.Email)

	// create and sign the jwt
	t, err := createAndSignJWT(dbUser)
	// Handle error in JWT creation/signing
	if err != nil {
		fmt.Println(err.Error())
		return []byte("Something went wrong with JWT"), err
	}

	// create transmissin friendly version of user struct and marshal it into JSON
	_user := transformedUser{ID: user.ID, Email: user.Email, Token: t}
	js, err := json.Marshal(_user)

	return js, err
}

// LoginUser handles database side of login, returns JWT token
func LoginUser(b []byte) ([]byte, error) {
	var user, dbUser userModel
	err := json.Unmarshal(b, &user)

	if err != nil {
		// handle internal server error
	}

	db.First(&dbUser, "email = ?", user.Email)

	if dbUser.ID == 0 {
		// handle user not found error
		return []byte("{\"message\": \"Something went wrong with JWT.\"}"), errors.New("not found")
	}

	match := checkPasswordHash(user.Password, dbUser.Password)
	if !match {
		return []byte("{\"message\": \"Check your inputs and try again.\"}"), errors.New("Unauthorized")
	}

	// Create and sign JWT
	t, err := createAndSignJWT(dbUser)
	// Handle error in JWT creation/signing
	if err != nil {
		fmt.Println(err.Error())
		return []byte("Something went wrong with JWT"), err
	}

	// create transmission friendly user struct
	var _user transformedUser
	_user.Email = user.Email
	_user.ID = dbUser.ID
	_user.Token = t

	// marshal user into JSON
	js, err := json.Marshal(_user)

	if err != nil {
		fmt.Println(err.Error())
		return []byte("Error parsing user into JSON"), err
	}

	return js, nil
}

// FetchMyInfo finds the given user in the db and returns info about them
func FetchMyInfo(uid float64) ([]byte, error) {
	var user userModel
	var _user listedUser
	struid := strconv.FormatFloat(uid, 'f', -1, 64)

	db.First(&user, "id = ?", struid)
	if user.ID == 0 {
		return []byte(""), errors.New("User not found")
	}

	_user = listedUser{ID: user.ID, Email: user.Email, Admin: user.Admin}

	js, err := json.Marshal(_user)
	return js, err
}

// user login password helper functions
// from https://gowebexamples.com/password-hashing/
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// my own JWT helper function
func createAndSignJWT(user userModel) (string, error) {
	// jwt stuff
	// create the token
	exp := time.Now().Add(time.Hour * 24).Unix()
	claim := CustomClaims{
		user.ID,
		user.Admin,
		jwt.StandardClaims{
			ExpiresAt: exp,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	secret := []byte(os.Getenv("SECRET"))

	t, err := token.SignedString(secret)

	return t, err
}
