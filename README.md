# Banking-system

A simple banking system built with Go, featuring user authentication, account management, and RESTful APIs.

## Features
- User registration and login with JWT authentication
- Account creation, retrieval, update, and deletion
- Secure password handling
- PostgreSQL database integration
- RESTful API endpoints

## Project Structure
```
Banking-system/
├── Bank/
│   ├── Authantication/    # JWT and user login logic
│   ├── config/            # Configuration loader
│   ├── controllers/       # API server and handlers
│   ├── db/                # Database connection and storage logic
│   ├── models/            # Data models
│   ├── go.mod, go.sum     # Go dependencies
│   └── main.go            # Application entry point
└── README.md
```

## Getting Started

### Prerequisites
- Go 1.18 or higher
- PostgreSQL database

### Setup
1. Clone the repository:
   ```sh
   git clone <repo-url>
   cd Banking-system/Bank
   ```
2. Configure your environment variables in a `.env` file:
   ```env
   POSTGRES_URL=postgres://user:password@localhost:5432/dbname?sslmode=disable
   JWT_KEY=your_jwt_secret
   ```
3. Install dependencies:
   ```sh
   go mod tidy
   ```
4. Run the application:
   ```sh
   go run main.go
   ```
   Or build and run:
   ```sh
   go build -o gobank
   ./gobank
   ```

## API Endpoints
- `POST /register` - Register a new user
- `POST /login` - User login
- `POST /accounts` - Create account
- `GET /accounts/{id}` - Get account by ID
- `DELETE /accounts/{id}` - Delete account
- `POST /transfer` - Transfer between accounts

## License
MIT