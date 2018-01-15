package controller

import (
	"bytes"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/katreinhart/askify-api-v2/model"
)

// CreateUser handles creation of new users via POST to /users/register
func CreateUser(w http.ResponseWriter, r *http.Request) {
	// Read in http request body
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	b := []byte(buf.String())

	js, err := model.CreateUser(b)

	handleErrorAndRespond(js, err, w)
}

// LoginUser handles login of existing user via POST to /users/login
func LoginUser(w http.ResponseWriter, r *http.Request) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	b := []byte(buf.String())

	js, err := model.LoginUser(b)

	handleErrorAndRespond(js, err, w)
}

// FetchUserInfo gets the information about the current user (by token) and returns it in JSON format.
func FetchUserInfo(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user")
	tok := user.(*jwt.Token)

	// no token present, so this is an unauthorized request.
	if tok == nil {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("{\"message\": \"Unauthorized\"}"))
		return
	}

	// get claims from token
	claims := tok.Claims.(jwt.MapClaims)
	// parse uid out of claims.
	uid, ok := claims["uid"].(float64)

	// Error parsing uid from token.
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("{\"message\": \"Could not parse token.\"}"))
	}

	// get the user from the model
	js, err := model.FetchMyInfo(uid)

	handleErrorAndRespond(js, err, w)
}

// FetchUserQuestions handles user question route
func FetchUserQuestions(w http.ResponseWriter, r *http.Request) {
	// get the URL parameter from the http request
	vars := mux.Vars(r)
	id := vars["id"]

	// fetch the question and an error if there is one
	js, err := model.FetchUserQuestions(id)

	handleErrorAndRespond(js, err, w)
}
