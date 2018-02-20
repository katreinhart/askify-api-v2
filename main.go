package main

import (
	"net/http"
	"os"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/codegangsta/negroni"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/katreinhart/askify-api-v2/controller"
	"github.com/rs/cors"
)

func main() {

	// get port variable from environment or set to default
	var port string
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	} else {
		port = "8080"
	}

	// CORS middleware setup
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedHeaders:   []string{"Accept-Encoding", "Accept-Language", "Authorization", "Content-Type"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "OPTIONS", "DELETE"},
		AllowCredentials: true,
	})

	// set up router
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/", homeHandler)

	// s is a subrouter to handle question routes
	api := r.PathPrefix("/api").Subrouter()

	// This isn't particularly RESTful but it works
	api.HandleFunc("/queue", controller.FetchQueue).Methods("GET")
	api.HandleFunc("/archive", controller.FetchArchive).Methods("GET")

	// get user info based on token
	api.HandleFunc("/user", controller.FetchUserInfo).Methods("GET")
	api.HandleFunc("/user/{id}/questions", controller.FetchUserQuestions).Methods("GET")

	// questions routes
	api.HandleFunc("/questions/", controller.FetchAllQuestions).Methods("GET")
	api.HandleFunc("/questions/", controller.CreateQuestion).Methods("POST")
	api.HandleFunc("/questions/{id}", controller.FetchSingleQuestion).Methods("GET")
	api.HandleFunc("/questions/{id}", controller.UpdateQuestion).Methods("PUT")

	// nested answer routes
	api.HandleFunc("/questions/{id}/answers", controller.FetchQuestionAnswers).Methods("GET")
	api.HandleFunc("/questions/{id}/answers", controller.CreateAnswer).Methods("POST")
	api.HandleFunc("/questions/{id}/answers/{aid}", controller.FetchSingleAnswer).Methods("GET")
	api.HandleFunc("/questions/{id}/answers/{aid}", controller.UpdateAnswer).Methods("PUT")

	// handle update cohorts
	api.HandleFunc("/cohorts", controller.AddCohort).Methods("POST")

	// u is another subrouter to handle auth routes
	u := r.PathPrefix("/auth").Subrouter()
	u.HandleFunc("/register", controller.CreateUser).Methods("POST")
	u.HandleFunc("/login", controller.LoginUser).Methods("POST")

	// cohort routes for fetching & adding to list of active cohorts
	u.HandleFunc("/cohorts", controller.FetchCohortList).Methods("GET")

	// JWT Middleware handles authorization configuration
	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("SECRET")), nil
		},
		SigningMethod: jwt.SigningMethodHS256,
	})

	// muxRouter uses Negroni handles the middleware for authorization
	muxRouter := http.NewServeMux()
	muxRouter.Handle("/", r)
	muxRouter.Handle("/api/", negroni.New(
		negroni.HandlerFunc(jwtMiddleware.HandlerWithNext),
		negroni.Wrap(api),
	))

	// Negroni handles the middleware chaining with next
	n := negroni.Classic()

	// Use CORS
	n.Use(c)

	// handle routes with the muxRouter
	n.UseHandler(muxRouter)

	// listen and serve!
	http.ListenAndServe(":"+port, handlers.RecoveryHandler()(n))
}

// homeHandler handles the / route
func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{\"message\": \"Hello world\"}"))
}
