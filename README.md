# ✨ Reddit Clone API - User Service ✨

Welcome to the backend API for a Reddit Clone project, meticulously crafted with Go! This service is dedicated to managing all aspects of user authentication, authorization, and profile management, providing a robust foundation for a dynamic social platform. 🚀

## 📝 Description

This Go-based API module, specifically the `controllers` package, is the core of user interaction. It provides a comprehensive suite of endpoints for users to register, log in, manage their credentials, and interact with their profiles. From secure password hashing to email verification and token-based authentication, this service ensures a smooth and secure user experience. It leverages the power of the Gin web framework for efficient routing and MongoDB for persistent data storage. 🛡️

## 🌟 Features

Our user service comes packed with essential functionalities to empower user engagement:

*   **🔑 User Registration (`CreateUser`):**
    *   Enables new users to sign up securely.
    *   Automatically hashes passwords using `bcrypt` for top-tier security.
    *   Assigns a unique `UserId` using MongoDB's `ObjectID`.
    *   Initiates an email verification flow by generating and sending a `verificationToken`.
    *   Prevents duplicate user registrations based on email.
    *   Timestamps user creation and update for auditing.

*   **✅ Email Verification (`VerifyEmail`):**
    *   Allows users to verify their email address using a unique token.
    *   Updates the user's `email_verified` status to `true` upon successful verification.
    *   Removes the `verification_token` after successful verification.
    *   Integrates with a `workers` queue to potentially send post-verification emails or perform other background tasks.

*   **🚪 User Login (`Login`):**
    *   Authenticates users by verifying their email and securely comparing passwords using `bcrypt`.
    *   Ensures that only verified email accounts can log in.
    *   Upon successful login, generates secure `token` and `refreshToken` for session management.
    *   Provides user details including `UserId`, `FirstName`, `LastName`, `Email`, `Role`, and the generated tokens.

*   **🔒 Password Reset Request (`ResetPasswordRequest`):**
    *   Enables users to request a password reset by providing their email.
    *   Generates a unique `reset_token` and updates the user's record in the database.
    *   *Note: Emailing the reset token is implied but not directly shown in this snippet for `ResetPasswordRequest`.*

*   **🔄 Password Reset Confirmation (`ResetPassword`):**
    *   Allows users to set a new password using a valid `reset_token`.
    *   Hashes the new password with `bcrypt` before updating.
    *   Clears the `reset_token` after the password has been successfully updated.

*   **👥 User Listing with Advanced Filtering (`GetAllUsers`):**
    *   Fetches a list of all registered users.
    *   **Search Functionality:** Supports searching users by `first_name`, `last_name`, or `role` using case-insensitive regular expressions.
    *   **Sorting Options:** Allows sorting results by `created_at` in ascending (`asc`) or descending (`desc`) order.
    *   **Pagination:** Implements pagination to retrieve users in manageable chunks, supporting `page` and `perPage` parameters.

*   **👤 Single User Retrieval (`GetUser`):**
    *   Retrieves detailed information for a specific user based on their `userId`.
    *   Ensures robust error handling for missing `userId` parameters.

*   **⬆️ File Uploads (`UploadFiles`):**
    *   Handles multi-part form data to facilitate the upload of multiple files.
    *   Utilizes a dedicated `utils.UploadFiles` function to process and store the uploaded files.
    *   Returns the URLs or identifiers of the uploaded files.

## ⚙️ Installation

To get this powerful user service up and running locally, follow these steps.

### Prerequisites

*   **Go:** Ensure you have Go installed (version 1.18 or higher is recommended).
    *   [Download Go](https://golang.org/dl/)
*   **MongoDB:** A running MongoDB instance is required for data persistence. You can run it locally or use a cloud service like MongoDB Atlas.

### Steps

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/EsanSamuel/Reddit_Clone.git
    cd Reddit_Clone
    ```

2.  **Install Go Modules:**
    Navigate to the project root and install the necessary Go modules.
    ```bash
    go mod tidy
    ```
    This command will download all the dependencies listed in `go.mod` (which include `gin-gonic`, `mongo-driver`, `bcrypt`, etc.).

3.  **MongoDB Setup:**
    Ensure your MongoDB instance is accessible. You will need to configure the connection string, likely through environment variables or a configuration file, which this `controllers` package then uses via `database.UserCollection`.

## 🚀 Usage

Once installed, you can run the application and interact with its API endpoints.

### Running the Application

While the `controllers` package defines the handlers, you'll need a `main.go` file (not provided in the snippet) to set up the Gin router and start the HTTP server. A typical `main.go` might look something like this (for illustrative purposes):

```go
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/EsanSamuel/Reddit_Clone/controllers" // Assuming this path
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// Example of registering routes
	router.POST("/users", controllers.CreateUser())
	router.GET("/users/verify", controllers.VerifyEmail())
	router.POST("/login", controllers.Login())
	router.POST("/password-reset-request", controllers.ResetPasswordRequest())
	router.POST("/password-reset", controllers.ResetPassword())
	router.GET("/users", controllers.GetAllUsers())
	router.GET("/users/:userId", controllers.GetUser())
	router.POST("/upload", controllers.UploadFiles()) // Example, might need middleware for auth

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if not specified
	}

	fmt.Printf("Server starting on port %s\n", port)
	log.Fatal(router.Run(":" + port))
}
```

Then, you can run your application:

```bash
go run main.go
```

### API Endpoints Examples (using cURL)

Here are some examples of how to interact with the API:

#### 1. Create a New User (Register)

```bash
curl -X POST http://localhost:8080/users \
     -H "Content-Type: application/json" \
     -d '{
           "first_name": "John",
           "last_name": "Doe",
           "email": "john.doe@example.com",
           "password": "securePassword123",
           "role": "user"
         }'
```

#### 2. Verify Email

After creating a user, you would receive a verification email containing a token.
```bash
curl -X GET "http://localhost:8080/users/verify?token=YOUR_VERIFICATION_TOKEN"
```

#### 3. User Login

```bash
curl -X POST http://localhost:8080/login \
     -H "Content-Type: application/json" \
     -d '{
           "email": "john.doe@example.com",
           "password": "securePassword123"
         }'
```

#### 4. Get All Users (with search, sort, and pagination)

```bash
# Get all users, sorted by creation date descending, 2nd page, 5 items per page
curl -X GET "http://localhost:8080/users?search=john&sort=desc&page=2&perPage=5"
```

#### 5. Get a Single User

```bash
curl -X GET "http://localhost:8080/users/YOUR_USER_ID"
```

#### 6. Upload Files

```bash
curl -X POST http://localhost:8080/upload \
     -H "Content-Type: multipart/form-data" \
     -F "files=@/path/to/your/image1.jpg" \
     -F "files=@/path/to/your/document.pdf"
```

## 📂 Folder Structure Explanation

The project is thoughtfully organized to promote modularity and maintainability. Based on the imports and the context of the `controllers` package, here's an overview of the likely structure:

```
Reddit_Clone/
├── controllers/          # 🎯 This directory contains the HTTP handlers (the code you provided).
│   └── user_controller.go  # Example: Handles all user-related API logic.
├── database/             # 💾 Manages database connections and provides access to collections (e.g., UserCollection).
│   └── mongo.go            # Example: MongoDB connection setup.
├── jobs/                 # ⚙️ Directory for background jobs or asynchronous tasks.
│   └── workers/          # 👷‍♂️ Contains worker functions, like the email queue mentioned for post-verification tasks.
│       └── email_worker.go # Example: Handles sending emails asynchronously.
├── models/               # 📦 Defines data structures (structs) for various entities (e.g., User, UserLogin, ForgetPasswordDTO).
│   └── user_model.go       # Example: User data model.
├── utils/                # 🛠️ Houses utility functions such as password hashing, token generation, email sending, and file uploads.
│   └── auth_utils.go       # Example: Token generation, password hashing.
│   └── email_utils.go      # Example: Email sending functions.
│   └── file_utils.go       # Example: File upload logic.
└── main.go               # 🚀 The entry point of the application, where routes are defined and the server starts.
└── go.mod                # Go module definition file.
└── go.sum                # Go module checksums.
```

## 💻 Technologies

This service is built upon a robust stack of modern technologies:

*   **Go (Golang):** The primary programming language, chosen for its performance, concurrency, and strong type safety. 💖
*   **Gin Web Framework:** A high-performance HTTP web framework for Go, used for building powerful and efficient APIs. ⚡
*   **MongoDB:** A popular NoSQL database, providing a flexible and scalable solution for data storage. 📊
*   **go.mongodb.org/mongo-driver:** The official Go driver for MongoDB, facilitating seamless interaction with the database. 🐘
*   **golang.org/x/crypto/bcrypt:** A robust library for secure password hashing, ensuring user credentials are well-protected. 🔒
*   **context (Go Standard Library):** Used for managing request-scoped values, cancellation signals, and deadlines across API operations. ⏳
*   **regexp (Go Standard Library):** Utilized for advanced search functionalities, enabling pattern matching for user queries. 🔍
*   **time (Go Standard Library):** For handling timestamps and managing time-related operations, like user creation/update times. ⏰
*   **net/http (Go Standard Library):** Fundamental for handling HTTP requests and responses. 🌐
*   **fmt, strconv, strings (Go Standard Library):** Essential utilities for formatting, string conversions, and manipulations. ✍️
*   **JWT (JSON Web Tokens):** Implied by `utils.GenerateTokens` and `utils.UpdateTokens` for secure, stateless authentication. 🔑
*   **Email Service Integration:** Implied by `utils.SendVerificationEmail` and `workers.SendEmailQueue`, indicating the ability to dispatch emails. 📧

## 📄 License

The license information for this project is not provided in the current snippet. Please refer to the project's root directory for the full `LICENSE` file.
