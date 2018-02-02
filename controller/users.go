package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

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

	var u model.UserModel

	err := json.Unmarshal(b, &u)
	if err != nil {
		handleErrorAndRespond(nil, err, w)
		return
	}

	// Create user in model
	_u, err := model.CreateUser(u)

	if err != nil {
		handleErrorAndRespond(nil, err, w)
		return
	}

	js, err := json.Marshal(_u)

	handleErrorAndRespond(js, err, w)
}

// LoginUser handles login of existing user via POST to /users/login
func LoginUser(w http.ResponseWriter, r *http.Request) {

	// Read in http request body
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	b := []byte(buf.String())

	var u model.UserModel
	err := json.Unmarshal(b, &u)

	if err != nil {
		handleErrorAndRespond(nil, err, w)
		return
	}

	var _u model.TransformedUser
	_u, err = model.LoginUser(u)

	if err != nil {
		handleErrorAndRespond(nil, err, w)
		return
	}

	js, err := json.Marshal(_u)
	handleErrorAndRespond(js, err, w)
}

// FetchUserInfo gets the information about the current user (by token) and returns it in JSON format.
func FetchUserInfo(w http.ResponseWriter, r *http.Request) {
	uid, err := GetUIDFromBearerToken(r)

	if err != nil {
		handleErrorAndRespond(nil, err, w)
		return
	}
	// get the user from the model

	_user, err := model.FetchMyInfo(uid)

	if err != nil {
		handleErrorAndRespond(nil, err, w)
	}
	js, err := json.Marshal(_user)

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

// GetUIDFromBearerToken does what it says on the tin
func GetUIDFromBearerToken(r *http.Request) (string, error) {
	user := r.Context().Value("user")
	tok := user.(*jwt.Token)
	var err error

	// no token present, so this is an unauthorized request.
	if tok == nil {
		err = errors.New("Forbidden")
	}

	// get claims from token
	claims := tok.Claims.(jwt.MapClaims)

	// parse uid out of claims.
	uid, ok := claims["uid"].(float64)

	// Error parsing uid from token.
	if !ok {
		err = errors.New("Forbidden")
	}

	// UID parsed from token is of type float64; we need it as a string.
	struid := strconv.FormatFloat(uid, 'f', -1, 64)

	return struid, err
}
