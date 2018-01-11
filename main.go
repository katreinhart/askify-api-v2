package main

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
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

	http.ListenAndServe(":"+port, r)

}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{\"message\": \"Hello world\"}"))
}
