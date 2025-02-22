# Gin Server

## How to run

1. Install Go: https://golang.org/doc/install
2. Install Nodemon: `npm install -g nodemon`
3. Clone the repository: `git clone`
4. Change to the project directory: `cd backend/gin-server`
5. Install dependencies: `go mod tidy`
6. Run the server: `go run cmd/server/main.go` or with hot reloading: `nodemon`


## Project structure

- `cmd/server/`: Contains the main.go file, the entry point of your server application. It's responsible for wiring up all the components (configuration, database, handlers, middleware) and starting the Gin server.
- `internal/`: This is where most of your application's business logic resides. It's considered internal because it's not meant to be imported by external packages.
  - `internal/db/`: Manages database connections.
  - `internal/handlers/`: Contains Gin handlers that receive HTTP requests, handle encoding/decoding, call services for business logic, and return HTTP responses. We'll version these handlers.
  - `internal/middleware/`: Implements Gin middlewares for cross-cutting concerns like authentication, logging, and request ID tracing.
  - `internal/models/`: Defines data models (structs) that represent your application entities (User). These should be database-agnostic as much as possible.
  - `internal/repositories/`: Contains database access code. Repositories are responsible for CRUD operations on the database and should be independent of the HTTP layer.
  - `internal/services/`: Contains the core business logic of your application. Services should be independent of HTTP and database details, making them easily testable and reusable.
  - `internal/config/`: Handles loading and managing application configuration from environment variables or config files.
  - `internal/router/`: Configures the Gin router, registering routes, middleware, and handlers.
  - `internal/logger/`: Initializes the application logger, allowing you to log messages with different severity levels.
  - `internal/utils/`: Contains utility functions that don't fit elsewhere.
- `tests/`: Contains tests for your application (unit tests, integration tests, end-to-end tests).
- `.env`: Stores environment-specific configuration variables (e.g., database connection strings, JWT secrets).
- `go.mod`, `go.sum`: Go modules files for dependency management.
- `README.md`: Project documentation.
