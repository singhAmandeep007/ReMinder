# Go API Server

## Overview
This project is a scalable and maintainable API server built with Go, featuring JWT authentication and a database-agnostic architecture. It supports both SQL (SQLite) and NoSQL (Firestore,MongoDB) databases, allowing for flexible data storage options.

## Project Structure
```
go-api-server
├── internal
    ├── cmd
│   └── server
│       └── main.go          # Entry point of the application
│   ├── internal
├── pkg
│   └── logger               # Logging utility
│       └── logger.go       # Structured logging
├── .env                     # Environment variables
├── .env.example             # Example environment variables
├── .gitignore               # Git ignore file
├── go.mod                   # Module dependencies
├── go.work                  # Workspace setup
└── README.md                # Project documentation
```

## Getting Started

### Prerequisites
- Go (version 1.16 or higher)
- SQLite (if using SQLite)
- MongoDB (if using MongoDB)

### Installation
1. Clone the repository:
   ```
   git clone https://github.com/yourusername/go-api-server.git
   cd go-api-server
   ```

2. Create a `.env` file based on the `.env.example` file and configure your environment variables.

3. Install dependencies:
   ```
   go mod tidy
   ```

## Contributing
Contributions are welcome! Please open an issue or submit a pull request for any enhancements or bug fixes.

## License
This project is licensed under the MIT License. See the LICENSE file for details.


## How to run

1. Install Go: https://golang.org/doc/install
2. Install Nodemon: `npm install -g nodemon`
3. Clone the repository: `git clone`
4. Change to the project directory: `cd backend/gin-server`
5. Install dependencies: `go mod tidy`
6. Run the server: `go run cmd/server/main.go` or with hot reloading: `nodemon`
