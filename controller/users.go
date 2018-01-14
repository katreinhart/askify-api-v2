package controller

import (
	"bytes"
	"fmt"
	"net/http"
	"os"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/katreinhart/askify-api-v2/model"
)

// CreateUser handles creation of new users via POST to /users/register
func CreateUser(w http.ResponseWriter, r *http.Request) {
	// Read in http request body
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	b := []byte(buf.String())

	js, err := model.CreateUser(b)

	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(js)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(js)
}

// LoginUser handles login of existing user via POST to /users/login
func LoginUser(w http.ResponseWriter, r *http.Request) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	b := []byte(buf.String())

	js, err := model.LoginUser(b)

	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(js)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(js)
}

func FetchUserInfo(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user")
	tok := user.(*jwt.Token)

	if tok == nil {
		panic("token not found")
	}

	fmt.Fprintf(os.Stderr, "Header:\n%v\n", tok.Header)
	fmt.Fprintf(os.Stderr, "Claims:\n%v\n", tok.Claims)
	claims := tok.Claims.(jwt.MapClaims)
	fmt.Fprintf(os.Stderr, "UID:\n%v\n", claims["uid"])
	uid, ok := claims["uid"].(float64)

	if !ok {
		panic("what the fuck")
	}
	js, err := model.FetchMyInfo(uid)

	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		if err.Error() == "Not found" {
			w.WriteHeader(http.StatusNotFound)
			w.Write(js)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Something went wrong"))
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Write(js)
}
