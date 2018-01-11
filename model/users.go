package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
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

	db.Save(&user)

	// create the token
	exp := time.Now().Add(time.Hour * 24).Unix()
	claim := jwt.StandardClaims{Id: string(dbUser.ID), ExpiresAt: exp}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	secret := []byte(os.Getenv("SECRET"))

	t, err := token.SignedString(secret)

	if err != nil {
		fmt.Println(err.Error())
		return []byte("{\"message\": \"Something went wrong with JWT.\"}"), err
	}

	_user := transformedUser{ID: user.ID, Email: user.Email, Token: t}
	js, err := json.Marshal(_user)

	return js, err
}

// user login password helper functions
// from https://gowebexamples.com/password-hashing/
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
