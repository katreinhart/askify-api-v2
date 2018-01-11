package main

import (
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/katreinhart/askify-api-v2/controller"
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
	s.HandleFunc("/{id}", controller.FetchSingleQuestion).Methods("GET")
	s.HandleFunc("/{id}", controller.UpdateQuestion).Methods("PUT")

	// nested answer routes
	s.HandleFunc("/{id}/answers", controller.FetchQuestionAnswers).Methods("GET")
	s.HandleFunc("/{id}/answers", controller.CreateAnswer).Methods("POST")
	s.HandleFunc("/{id}/answers/{aid}", controller.FetchSingleAnswer).Methods("GET")

	// u is another subrouter to handle users routes
	u := r.PathPrefix("/users").Subrouter()
	u.HandleFunc("/register", controller.CreateUser).Methods("POST")
	u.HandleFunc("/login", controller.LoginUser).Methods("POST")

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
