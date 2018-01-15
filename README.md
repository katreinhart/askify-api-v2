# Askify 2.0 API 
a gSchool Question Queue API
written in Go by [katreinhart](kat.reinhart@gmail.com)

## Askify is a question queue built for students at Galvanize.  
This API is my independent study project for WDI Q4 and is implemented in Go using the Gorilla toolkit for parsing URL tokens and the Negroni framework for handling middleware. 

### Getting started
1. You will need [Go](http://golang.org) installed to run this project. 
1. Fork and clone the repository to your local machine into your ```$GOPATH``` (usually /Users/(your user name)/go/).
1. Navigate into the project folder and run ```go install``` to install dependencies. 
1. You will need a local postgres server running - [here](https://launchschool.com/blog/how-to-install-postgresql-on-a-mac) is a good overview on how to do that on MacOS.
1. Create a local postgres database called askify_v2_dev. 
1. Change the name of ```.env.sample``` to ```.env``` and edit the file with your credentials or secret key as necessary
1. ```go run main.go``` will run the server locally, and you can test API endpoints from Postman or similar. 
1. ```go build``` will compile the project into an executable binary. 

### API Routes Implemented 
(all API routes begin with ```/api```)
* ```/questions/``` GET, POST
* ```/questions/{id}``` GET, PUT
* ```/questions/{id}/answers``` GET, POST
* ```/questions/{id}/answers/{aid}``` GET
* ```/questions/open``` GET <-- returns all unanswered questions
* ```/user ``` GET - returns information about the user based on contents of JWT
* ```/users/{id}/questions``` returns all quetions a user has asked
* ```/archive``` returns all answered question with nested answers

### Authentication & Authorization
* Authentication routes begin with ```/auth```
* Routes are ```/register``` and ```/login```
* Both auth routes require an "email" and "password" field.
* Auth register route also takes "fname" and "cohort" fields
* Auth routes return user information and a JWT token.
* API routes begin with /api/ and are JWT-token protected.
* Use the JWT token returned from register or login as an authorization bearer token to access these routes. 

#### Backlog
* Ensure update permissions only to admins and user who owns question/answer
* Does deeply nested archive route return user names? 
* Implement Delete functionality? 

* ~~Use ```https://github.com/auth0/go-jwt-middleware``` go-jwt-middleware to protect routes.~~
* ~~Update user registration to take fname and cohort ID~~
* ~~Implement ```/users/{id}/questions``` GET route to see all users' questions~~
* ~~Implement ```/archive``` GET route to get all answered questions with nested answers.~~
* ~~POST to /questions should not require "userid" field - that information should be parsed from token.~~

### Work In Progress
This is a work in progress as of January 2018.
