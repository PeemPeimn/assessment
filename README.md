# Assessment

This is a public repository for KBTG Go Software Engineering Bootcamp's post-test assessment.

Created by: 

Peemwaruch Intamool, intamoon.peem@gmail.com


## Setting up the environment
Please check these requirements before running the program.
* Set your `PORT` to `:2565`.
* Set your `DATABASE_URL` to the URL of your Postgres database.
* To run the integration tests, make sure your machine can run docker-compose.

## How to run the program
* To run the main program: `DATABASE_URL=postgres://... PORT=:2565 go run server.go`
* To run unit tests: `go test -v ./...`
* To run the integration tests: `docker-compose -f docker-compose.test.yml up --build --abort-on-container-exit --exit-code-from it_tests` 

## Notes

* `Echo` library is used to implement APIs.
* Authorization is handled by the middleware. Your request will be accept if and only if your request's header has the `Authorization` field with value `November 10, 2009` (as stated in the Postman tests.)
* Each user story is created in its own branch. You can check with `git log --graph` afther cloning this project.
* Expenses routes' logic is implemented in the `expenses` folder.
* `db.go` contains code used to handle database connections.
* `handler.go` contains all route handling functions.
* `handler_it_test.go` consists of integration tests for each handler function and other files that end with `_test.go` are unit tests code.