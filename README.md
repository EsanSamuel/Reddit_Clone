# Reddit Clone Backend

This project serves as the backend (server) for a Reddit clone, developed using Go and the Gin web framework. It provides core functionalities for user management and subreddit operations, interacting with a MongoDB database for persistence.

## Description

The backend application handles essential features for a social community platform similar to Reddit. It manages user accounts, including registration, authentication, email verification, and password resets. Furthermore, it facilitates the creation and management of subreddits, allowing users to join communities and assign moderator roles. User profile enhancements, such as avatar uploads, are also supported.

## Features

The server exposes endpoints for the following functionalities:

### User Management
*   **User Registration:** Create new user accounts with email, password, and other profile details.
*   **Email Verification:** Verify user emails using a generated token after registration.
*   **User Login:** Authenticate users with their email and password, generating access and refresh tokens upon successful login.
*   **Password Reset Request:** Initiate a password reset process by generating a reset token for a given email.
*   **Password Reset:** Update a user's password using a valid reset token.
*   **Retrieve All Users:** Fetch a list of all registered users, with support for search (by first name, last name, role), sorting (by creation date), and pagination.
*   **Retrieve Single User:** Get details for a specific user by their ID.
*   **User Avatar Upload:** Allow users to upload an avatar image to their profile.

### Subreddit Management
*   **Create Subreddit:** Establish new subreddits, with the creator automatically assigned as a moderator.
*   **Join Subreddit:** Enable users to join existing subreddits.
*   **Add Moderators:** Grant moderator roles to existing members within a subreddit.

## Installation

To set up and run this project locally, ensure you have a Go environment configured.

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/EsanSamuel/Reddit_Clone.git
    cd Reddit_Clone/server # Navigate to the server directory
    ```

2.  **Install dependencies:**
    The project uses Go modules. The necessary dependencies will be fetched automatically when you build or run the application. Key dependencies include:
    *   `github.com/gin-gonic/gin`: Web framework
    *   `go.mongodb.org/mongo-driver/v2`: MongoDB driver
    *   `golang.org/x/crypto/bcrypt`: For password hashing

    ```bash
    go mod tidy
    ```

3.  **Database Setup:**
    This project uses MongoDB. Ensure you have a MongoDB instance running and accessible. Configure your database connection details (e.g., connection string, database name) in the application's environment or configuration files, which are not provided in these snippets but are implied by the `database` package usage.

4.  **Environment Variables (Implied):**
    The application likely requires environment variables for:
    *   MongoDB connection URI.
    *   Email service credentials for sending verification and other emails.
    *   AWS S3 credentials for avatar uploads (implied by `UploadSingleFileToS3`).
    *   JWT secret keys for token generation.

5.  **Run the application:**
    ```bash
    go run main.go # Assuming main.go is the entry point in the server directory
    ```

## Usage

This is a backend application designed to be consumed by a frontend client or other services via its RESTful API.

### Example API Operations (Conceptual)

While specific routes are not provided, the controllers implement the following actions:

*   **Create a User:** Send a `POST` request to the user creation endpoint with user details (e.g., `email`, `password`, `first_name`, `last_name`).
    ```json
    POST /api/users/register
    Content-Type: application/json

    {
        "email": "test@example.com",
        "password": "securepassword",
        "first_name": "John",
        "last_name": "Doe"
    }
    ```
*   **Login User:** Send a `POST` request to the login endpoint with `email` and `password`.
    ```json
    POST /api/users/login
    Content-Type: application/json

    {
        "email": "test@example.com",
        "password": "securepassword"
    }
    ```
*   **Create a Subreddit:** Send a `POST` request to the subreddit creation endpoint with subreddit details (e.g., `name`, `description`, `creator_id`).
    ```json
    POST /api/subreddits
    Content-Type: application/json

    {
        "name": "MyAwesomeCommunity",
        "description": "A place for awesome people",
        "creator_id": "someUserId123"
    }
    ```

## Folder Structure Explanation

Based on the provided file paths and import statements, the project adheres to a structured backend architecture:

*   `server/controllers/`: Contains handler functions (`gin.HandlerFunc`) for various API endpoints. These files (`subreddit_controller.go`, `user_controller.go`) encapsulate the business logic related to specific domains (subreddits and users).
*   `database/` (inferred): This package (`github.com/EsanSamuel/Reddit_Clone/database`) likely manages the MongoDB connection and provides access to collection objects (`UserCollection`, `SubredditCollection`, `MemberCollection`).
*   `models/` (inferred): This package (`github.com/EsanSamuel/Reddit_Clone/models`) defines the data structures (structs) for various entities like `User`, `SubReddit`, `SubRedditMembers`, `UserLogin`, `UserDTO`, `ForgetPasswordRequestDTO`, and `ForgetPasswordDTO`, which are mapped to database documents and API request/response bodies.
*   `utils/` (inferred): This package (`github.com/EsanSamuel/Reddit_Clone/utils`) provides helper functions for common tasks such as password hashing, token generation, email sending, file validation, and file uploads (e.g., to S3).
*   `jobs/workers/` (inferred): This package (`github.com/EsanSamuel/Reddit_Clone/jobs/workers`) likely handles background tasks, such as sending emails asynchronously (e.g., `SendEmailQueue`).

## Technologies

*   **Language:** Go
*   **Web Framework:** Gin Gonic
*   **Database:** MongoDB
*   **Database Driver:** `go.mongodb.org/mongo-driver/v2`
*   **Password Hashing:** bcrypt
*   **Cloud Storage:** AWS S3 (implied for file uploads)
*   **Background Jobs:** Implied worker/queue system for tasks like email sending.

## License

License information is not specified in the provided project data.