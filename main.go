package main

import (
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/katreinhart/askify-api-v1/controller"
)

func main() {

	// get port variable from environment or set to default
	var port string
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	} else {
		port = "8080"
	}

	// set up router
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/", homeHandler)

	// s is a subrouter to handle question routes
	s := r.PathPrefix("/questions").Subrouter()
	s.HandleFunc("/", controller.FetchAllQuestions).Methods("GET")
	s.HandleFunc("/", controller.CreateQuestion).Methods("POST")

	// Logging handler enables standard HTTP logging
	loggedRouter := handlers.LoggingHandler(os.Stdout, r)

	// Recovery handler is a HOF that wraps a router enbling recovery from a panic
	http.ListenAndServe(":"+port, handlers.RecoveryHandler()(loggedRouter))
}

// homeHandler handles the / route
func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{\"message\": \"Hello world\"}"))
}
