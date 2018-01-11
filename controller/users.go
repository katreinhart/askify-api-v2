package controller

import (
	"bytes"
	"net/http"

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
