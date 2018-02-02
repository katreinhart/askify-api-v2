package model

import (
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

// CreateUser checks to see uniqueness of new id, saves the user info if unique, returns error otherwise.
func CreateUser(u UserModel) (TransformedUser, error) {

	// declare data structures to be used
	var dbUser UserModel

	// See if the user exists in the database. If so, return an error (no duplicates allowed)
	db.First(&dbUser, "email = ?", u.Email)
	if dbUser.ID != 0 {
		return TransformedUser{}, ErrorUserExists
	}

	// Hash the user's password using bcrypt helper function and handle any error
	hash, err := hashPassword(u.Password)
	if err != nil {
		return TransformedUser{}, err
	}

	// Overwrite the user's password with the hashed version (no plaintext storage of passwords)
	u.Password = hash

	// save the user in the DB
	db.Save(&u)

	// Get the user back from the database so you have the correct ID
	db.First(&dbUser, "email = ?", u.Email)

	// create and sign the JWT
	t, err := createAndSignJWT(dbUser)

	// Handle error in JWT creation/signing
	if err != nil {
		return TransformedUser{}, err
	}

	// create transformed version of user structure, marshal it into JSON and return
	_user := TransformedUser{ID: u.ID, Email: u.Email, FName: u.FName, Cohort: u.Cohort, Token: t}
	return _user, nil
}

// LoginUser handles database side of login, returns JWT token
func LoginUser(u UserModel) (TransformedUser, error) {

	// Declare data types and unmarshal JSON into user struct
	var dbUser UserModel

	db.First(&dbUser, "email = ?", u.Email)

	// handle user not found error
	if dbUser.ID == 0 {
		return TransformedUser{}, ErrorNotFound
	}

	// See if password matches the hashed password from the database
	match := checkPasswordHash(u.Password, dbUser.Password)
	if !match {
		return TransformedUser{}, ErrorForbidden
	}

	// Create and sign JWT; handle any error
	t, err := createAndSignJWT(dbUser)
	if err != nil {
		return TransformedUser{}, err
	}

	// create transmission friendly user struct
	_user := TransformedUser{ID: dbUser.ID, Email: dbUser.Email, FName: dbUser.FName, Cohort: dbUser.Cohort, Token: t}
	return _user, nil
}

// FetchMyInfo finds the given user in the db and returns info about them
func FetchMyInfo(uid string) (ListedUser, error) {

	// Declare data types
	var user UserModel
	var _user ListedUser

	// Find the user in the DB
	db.First(&user, "id = ?", uid)
	if user.ID == 0 {
		return ListedUser{}, ErrorNotFound
	}

	// Transform user into a listedUser format with ID, email, fname, admin status, cohort
	_user = ListedUser{ID: user.ID, FName: user.FName, Cohort: user.Cohort, Email: user.Email, Admin: user.Admin}

	return _user, nil
}

// --------------------- Helper Functions ---------------------
// user login password helper functions
// from https://gowebexamples.com/password-hashing/
func hashPassword(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(b), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// JWT helper function
func createAndSignJWT(user UserModel) (string, error) {

	// create the expiration time, build claim, create and sign token, and return.
	e := time.Now().Add(time.Hour * 24).Unix()
	c := CustomClaims{
		user.ID,
		user.Admin,
		jwt.StandardClaims{
			ExpiresAt: e,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	secret := []byte(os.Getenv("SECRET"))
	t, err := token.SignedString(secret)
	fmt.Println(t)
	return t, err
}
