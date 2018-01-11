# Askify 2.0 API 
a gSchool Question Queue App
written in Golang by [katreinhart](kat.reinhart@gmail.com)

## Askify is a question queue built for students at Galvanize.  
This API is my independent study project for WDI Q4 and is implemented in Go using the Gorilla toolkit and no web framework (because Go doesn't need a web framework!)

### Getting started
1. You will need [Go](http://golang.org) installed to run this project. 
1. Fork and clone the repository to your local machine into your ```$GOPATH``` (usually /Users/(your user name)/go/).
1. Navigate into the project folder and run ```go install``` to install dependencies. 
1. You will need a local postgres server running - [here](https://launchschool.com/blog/how-to-install-postgresql-on-a-mac) is a good overview on how to do that.
1. Create a local postgres database called askify_v2_dev. 
1. Change the name of ```.env.sample``` to ```.env```
1. ```go run main.go``` will run the server locally, and you can test API endpoints from Postman or similar. 

### Routes Implemented 
* ```/questions/``` GET, POST
* ```/questions/{id}``` GET, PUT
* ```/questions/{id}/answers``` GET, POST
* ```/questions/{id}/answers/{aid}``` GET

### Work In Progress
This is a work in progress (January 2018). 