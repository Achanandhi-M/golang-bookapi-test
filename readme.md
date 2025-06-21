```markdown
# Book Tracker API

A RESTful API built with Go for tracking books you're reading. Manage book details like title, author, and reading progress, with data stored in a PostgreSQL database. The project is containerized using Docker and includes unit, integration, and API tests to ensure reliability.

## Features
- **Create a book**: Add a new book with title, author, and progress (POST `/books`).
- **Retrieve all books**: List all books (GET `/books`).
- **Update a book**: Modify a book's details by ID (PUT `/books/{id}`).
- **Delete a book**: Remove a book by ID (DELETE `/books/{id}`).
- **Validation**: Ensures non-empty title, author, and non-negative progress.
- **High Test Coverage**: ~86% coverage with unit, integration, and API tests.

## Tech Stack
- **Language**: Go (Golang)
- **Database**: PostgreSQL
- **HTTP Router**: [gorilla/mux](https://github.com/gorilla/mux)
- **Database Library**: [jmoiron/sqlx](https://github.com/jmoiron/sqlx)
- **Containerization**: Docker & Docker Compose
- **Testing**: Go's `testing` package with `httptest`

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

- **`cmd/server/main.go`**: Application entry point, sets up the HTTP server.
- **`internal/`**: Private packages for database, handlers, models, and repository.
- **`tests/`**: Unit, integration, and API tests for robust coverage.
- **`.env`**: Environment variables for database configuration.
- **`Dockerfile` & `docker-compose.yml`**: Containerizes the app and database.

## Prerequisites
- **Docker**: Install [Docker](https://docs.docker.com/get-docker/) and [Docker Compose](https://docs.docker.com/compose/install/).
- **Git**: To clone the repository.
- **curl** (optional): For manual API testing.
- **Go** (optional): Only needed if running without Docker.

## Setup Instructions
1. **Clone the Repository**:
   ```bash
   git clone https://github.com/<your-username>/book-tracker.git
   cd book-tracker
   ```

2. **Create a `.env` File**:
   Create a `.env` file in the project root with the following content:
   ```plaintext
   POSTGRES_HOST=db
   POSTGRES_PORT=5432
   POSTGRES_USER=bookuser
   POSTGRES_PASSWORD=bookpassword
   POSTGRES_DB=bookdb
   ```
   - These variables configure the PostgreSQL connection.
   - Ensure `POSTGRES_HOST=db` matches the service name in `docker-compose.yml`.

3. **Build and Start Containers**:
   ```bash
   docker-compose up -d
   ```
   - Starts two containers:
     - `book-tracker-app-1`: Go API on `http://localhost:8080`.
     - `book-tracker-db-1`: PostgreSQL database.
   - Verify containers are running:
     ```bash
     docker ps
     ```

## Running the Application
The API is accessible at `http://localhost:8080` after starting the containers. Test the API with `curl` or a tool like Postman.

### Example `curl` Commands
1. **Create a Book**:
   ```bash
   curl -X POST http://localhost:8080/books \
     -H "Content-Type: application/json" \
     -d '{"title":"The Hobbit","author":"J.R.R. Tolkien","progress":50}'
   ```
   **Expected**: `201 Created`, returns book with ID.

2. **Retrieve All Books**:
   ```bash
   curl -X GET http://localhost:8080/books
   ```
   **Expected**: `200 OK`, returns array of books.

3. **Update a Book** (replace `1` with actual book ID):
   ```bash
   curl -X PUT http://localhost:8080/books/1 \
     -H "Content-Type: application/json" \
     -d '{"title":"The Hobbit Updated","author":"J.R.R. Tolkien","progress":75}'
   ```
   **Expected**: `200 OK`, returns updated book.

4. **Delete a Book** (replace `1` with actual book ID):
   ```bash
   curl -X DELETE http://localhost:8080/books/1
   ```
   **Expected**: `204 No Content`.

5. **Error Case: Invalid Input**:
   ```bash
   curl -X POST http://localhost:8080/books \
     -H "Content-Type: application/json" \
     -d '{}'
   ```
   **Expected**: `400 Bad Request`, error message.

For a full test script, see `test_api_curl_commands.sh` in the repository (run with `bash test_api_curl_commands.sh`).

## Running Tests
The project includes unit, integration, and API tests, achieving ~86% coverage.

1. **Access the App Container**:
   ```bash
   docker exec -it book-tracker-app-1 sh
   ```

2. **Run All Tests**:
   Inside the container, execute:
   ```bash
   go test -v -coverpkg=./internal/db,./internal/handlers,./internal/repository ./tests/... -coverprofile=coverage.out
   ```
   - `-v`: Verbose output, shows test details.
   - `-coverpkg`: Measures coverage for specified packages.
   - `-coverprofile`: Saves coverage data to `coverage.out`.
   - **Expected Output**: Test results (PASS/FAIL) and coverage (e.g., 86.6%).

3. **View Coverage Report**:
   ```bash
   go tool cover -func=coverage.out
   ```
   Shows per-function coverage.

4. **Generate HTML Coverage Report**:
   ```bash
   go tool cover -html=coverage.out -o coverage.html
   ```
   Copy the HTML file to your host machine:
   ```bash
   exit  # Exit container
   docker cp book-tracker-app-1:/app/coverage.html .
   ```
   Open `coverage.html` in a browser to visualize coverage.

## Troubleshooting
- **API Not Responding**:
  - Check container logs: `docker logs book-tracker-app-1`.
  - Ensure database is running: `docker ps`.
  - Verify `.env` credentials match `docker-compose.yml`.
- **Test Failures**:
  - Run tests with `-v` for detailed output: `go test -v ./tests/...`.
  - Check database connection: `psql -h localhost -p 5432 -U bookuser -d bookdb`.
  - Clear Go cache: `go clean -cache`.
- **Port Conflicts**:
  - If port `8080` is in use, change the port mapping in `docker-compose.yml` (e.g., `8081:8080`).
- **Database Issues**:
  - Truncate table to reset:
    ```bash
    docker exec -it book-tracker-app-1 psql -h db -U bookuser -d bookdb -c "TRUNCATE TABLE books RESTART IDENTITY;"
    ```

## Contributing
1. Fork the repository.
2. Create a branch: `git checkout -b feature/your-feature`.
3. Commit changes: `git commit -m "Add your feature"`.
4. Push to your fork: `git push origin feature/your-feature`.
5. Open a pull request.


## Acknowledgments
- Built with inspiration from Go's simplicity and clean architecture principles.
- Thanks to [gorilla/mux](https://github.com/gorilla/mux) and [sqlx](https://github.com/jmoiron/sqlx) for robust libraries.
```

