package controller

import (
	"net/http"

	"github.com/katreinhart/askify-api-v2/model"
)

func handleErrorAndRespond(js []byte, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")

	// handle the error cases
	if err != nil {
		if err == model.ErrorNotFound {
			w.WriteHeader(http.StatusNotFound)
		} else if err == model.ErrorUserExists {
			w.WriteHeader(http.StatusFound)
		} else if err == model.ErrorForbidden {
			w.WriteHeader(http.StatusUnauthorized)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		// in any case, send back the message created in the model
		w.Write(js)
		return
	}

	// Handle the success case
	w.WriteHeader(http.StatusOK)
	w.Write(js)
}
