JWT Authentication Server
A Go-based JWT authentication server with MongoDB integration.

Features
User registration and login
JWT token generation and validation
MongoDB database integration
Role-based access control
Password hashing with bcrypt
RESTful API endpoints
Setup
Install Dependencies

go mod tidy
Environment Variables Create a .env file in the root directory with the following variables:

MONGODB_URL=mongodb://localhost:27017
DATABASE_NAME=jwt_auth_db
SECRET_KEY=your-secret-key-here-make-it-long-and-secure
PORT=8000
MongoDB Make sure MongoDB is running on your system or update the MONGODB_URL to point to your MongoDB instance.

Run the Server

go run main.go
API Endpoints
Authentication
POST /users/signup - User registration
POST /users/login - User login
Protected Routes (require authentication)
GET /users - Get all users (Admin only)
GET /users/:user_id - Get specific user
GET /api-1 - Test protected endpoint
GET /api-2 - Test protected endpoint
Project Structure
├── controllers/     # HTTP request handlers
├── database/        # Database connection and utilities
├── helpers/         # Helper functions for auth and tokens
├── middleware/      # Authentication middleware
├── models/          # Data models
├── routes/          # Route definitions
└── main.go         # Application entry point
Dependencies
gin-gonic/gin - Web framework
go.mongodb.org/mongo-driver - MongoDB driver
github.com/dgrijalva/jwt-go - JWT implementation
golang.org/x/crypto/bcrypt - Password hashing
github.com/joho/godotenv - Environment variable loading
github.com/go-playground/validator/v10 - Data validation
