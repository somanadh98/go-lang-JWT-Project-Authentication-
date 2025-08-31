Project: Go-Based RESTful Authentication API
A secure and scalable RESTful API server built with Go (Golang) to handle user authentication and authorization using JSON Web Tokens (JWT). This project serves as a foundational microservice for applications requiring robust user management.
Key Features:
Secure User Authentication:
User Signup: Endpoint for new user registration with server-side validation. Passwords are never stored in plaintext; they are securely hashed using the bcrypt algorithm.
User Login: Endpoint for authenticating users with email and password. It verifies credentials against the stored hash and prevents timing attacks.
JWT-Based Session Management:
Token Generation: Upon successful login, the API generates both a short-lived access token and a long-lived refresh token.
Stateless Authorization: The stateless nature of JWTs allows for easy scaling of backend services, as no session state needs to be stored on the server.
Token Refresh: A mechanism to use the refresh token to obtain a new access token without requiring the user to log in again.
Role-Based Access Control (RBAC):
The API implements middleware to protect specific routes based on user roles (e.g., USER vs. ADMIN).
Demonstrated with an admin-only endpoint (/users) to fetch a list of all registered users.
RESTful API Endpoints:
POST /users/signup: Register a new user.
POST /users/login: Authenticate a user and receive JWTs.
GET /users: (Admin only) Retrieve a paginated list of all users.
GET /users/:user_id: (Protected) Retrieve profile information for a specific user.
Database Integration:
Utilizes MongoDB as the data store for user information.
Efficiently queries the database to check for existing users, retrieve user data, and update tokens.
Technical Stack:
Language: Go (Golang)
Web Framework: Gin (for high-performance HTTP routing and middleware)
Database: MongoDB (with the official go.mongodb.org/mongo-driver)
Authentication: JSON Web Tokens (JWT)
Password Hashing: Bcrypt (golang.org/x/crypto/bcrypt)
Validation: go-playground/validator/v10 for struct validation.
