# Book Tracker API

## A RESTful API built with Go for tracking books you're reading. You can manage book details like title, author, and reading progress, with all data stored in a PostgreSQL database. The project is containerized using Docker and includes unit, integration, and API tests to ensure reliability.

### Features

* Create a book: Add a new book with title, author, and progress (POST `/books`)
* Retrieve all books: List all books (GET `/books`)
* Update a book: Modify a book's details by ID (PUT `/books/{id}`)
* Delete a book: Remove a book by ID (DELETE `/books/{id}`)
* Validation: Ensures non-empty title, author, and non-negative progress
* High Test Coverage: Approximately 86% coverage with unit, integration, and API tests

### Tech Stack

* Language: Go (Golang)
* Database: PostgreSQL
* HTTP Router: gorilla/mux
* Database Library: sqlx
* Containerization: Docker & Docker Compose
* Testing: Go’s testing package with httptest

## Project Structure

```
book-tracker/
├── Dockerfile
├── docker-compose.yml
├── go.mod
├── go.sum
├── .env
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── db/
│   │   └── db.go
│   ├── handlers/
│   │   └── handlers.go
│   ├── models/
│   │   └── book.go
│   └── repository/
│       └── book_repository.go
└── tests/
    ├── api/
    │   └── api_test.go
    ├── integration/
    │   └── repository_test.go
    └── unit/
        └── handlers_test.go
```

## Explanation of Structure:

* `cmd/server/main.go`: The entry point that sets up the HTTP server
* `internal/`: Contains private application logic including handlers, models, db code, and repository pattern
* `tests/`: Includes unit, integration, and API tests to ensure proper behavior
* `.env`: Stores environment variables for configuring the database
* `Dockerfile` & `docker-compose.yml`: Used for containerizing and orchestrating the app and database

## Prerequisites

To run this project, make sure you have:

* Docker: [https://docs.docker.com/get-docker/](https://docs.docker.com/get-docker/)
* Docker Compose: [https://docs.docker.com/compose/install/](https://docs.docker.com/compose/install/)
* Git: For cloning the repository
* curl (optional): For manually testing the API
* Go (optional): Only required if not using Docker

## Setup Instructions

Step 1: Clone the Repository

```bash
git clone https://github.com/<your-username>/book-tracker.git
cd book-tracker
```

Step 2: Create a `.env` File

```env
POSTGRES_HOST=db
POSTGRES_PORT=5432
POSTGRES_USER=bookuser
POSTGRES_PASSWORD=bookpassword
POSTGRES_DB=bookdb
```

Make sure the value of `POSTGRES_HOST` matches the service name (`db`) in `docker-compose.yml`.

Step 3: Build and Start the Containers

```bash
docker-compose up -d
```

This will start:

* `book-tracker-app-1`: Go API server accessible at `http://localhost:8080`
* `book-tracker-db-1`: PostgreSQL instance

Verify that containers are running:

```bash
docker ps
```

## Running the Application

Once the containers are up, the API is available at `http://localhost:8080`. You can test it using curl or Postman.

Example curl Commands

Create a Book

```bash
curl -X POST http://localhost:8080/books \
  -H "Content-Type: application/json" \
  -d '{"title":"The Hobbit","author":"J.R.R. Tolkien","progress":50}'
```

Expected: HTTP 201 Created with the created book object

Retrieve All Books

```bash
curl -X GET http://localhost:8080/books
```

Expected: HTTP 200 OK with a list of books

Update a Book (Replace `1` with actual ID)

```bash
curl -X PUT http://localhost:8080/books/1 \
  -H "Content-Type: application/json" \
  -d '{"title":"The Hobbit Updated","author":"J.R.R. Tolkien","progress":75}'
```

Expected: HTTP 200 OK with updated book data

Delete a Book (Replace `1` with actual ID)

```bash
curl -X DELETE http://localhost:8080/books/1
```

Expected: HTTP 204 No Content

Invalid Input Example

```bash
curl -X POST http://localhost:8080/books \
  -H "Content-Type: application/json" \
  -d '{}'
```

Expected: HTTP 400 Bad Request with error message

Note: A script `test_api_curl_commands.sh` is available for testing all endpoints at once.

## Running Tests

The app includes all types of tests and reaches approximately 86% coverage.

Step 1: Access the App Container

```bash
docker exec -it book-tracker-app-1 sh
```

Step 2: Run All Tests

```bash
go test -v -coverpkg=./internal/db,./internal/handlers,./internal/repository ./tests/... -coverprofile=coverage.out
```

* `-v` shows verbose output
* `-coverpkg` tracks coverage across multiple packages
* `-coverprofile` stores coverage data

Step 3: View Function-Level Coverage

```bash
go tool cover -func=coverage.out
```

Step 4: Generate and View HTML Coverage Report

```bash
go tool cover -html=coverage.out -o coverage.html
exit
docker cp book-tracker-app-1:/app/coverage.html .
```

Open `coverage.html` in your browser to explore the coverage visually.

## Troubleshooting

API Not Responding

* Check logs: `docker logs book-tracker-app-1`
* Verify PostgreSQL is running: `docker ps`
* Confirm `.env` matches `docker-compose.yml`

Test Failures

* Run with verbose mode: `go test -v ./tests/...`
* Ensure DB is accessible: `psql -h localhost -p 5432 -U bookuser -d bookdb`
* Clean cache: `go clean -cache`

Port Conflicts

* If `8080` is busy, change it in `docker-compose.yml` to another port like `8081:8080`

Resetting Database

To wipe the `books` table:

```bash
docker exec -it book-tracker-app-1 psql -h db -U bookuser -d bookdb -c "TRUNCATE TABLE books RESTART IDENTITY;"
```

## Contributing

1. Fork this repository
2. Create a feature branch: `git checkout -b feature/your-feature`
3. Commit your work: `git commit -m "Add your feature"`
4. Push it: `git push origin feature/your-feature`
5. Open a pull request

## Acknowledgments

This project was inspired by Go's clean and minimal architecture. Thanks to open-source contributors of libraries like gorilla/mux and sqlx for making development easy and robust.