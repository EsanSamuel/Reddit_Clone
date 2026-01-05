# Reddit Clone - Server Controllers

This repository contains the backend controllers for a Reddit clone application, specifically focusing on user authentication, management, and subreddit functionalities. Developed in Go using the Gin web framework and MongoDB for data persistence, these controllers provide the core API endpoints for interacting with user and community data.

## Description

This project comprises the controller logic for managing user accounts and subreddits within a Reddit clone application. It handles operations such as user registration, login, email verification, password reset, and avatar uploads. For subreddits, it supports creation, joining, moderator assignment, and comprehensive searching/retrieval capabilities.

## Features

The following features are implemented through the provided Go controllers:

### User Management

*   **User Registration:** Create new user accounts with password hashing, email verification token generation, and email sending.
*   **Email Verification:** Verify user emails using a generated token.
*   **User Login:** Authenticate users, check email verification status, and generate access and refresh tokens.
*   **Password Reset Request:** Initiate a password reset process by generating a reset token.
*   **Password Reset:** Reset user passwords using a provided token and hash the new password.
*   **Get All Users:** Retrieve a paginated list of all users with search (by first name, last name, role) and sorting options (by creation date).
*   **Get User by ID:** Fetch a single user's details using their unique ID.
*   **Upload Avatar:** Upload a user's profile avatar image to a cloud storage service (implied by S3 utility usage).

### Subreddit Management

*   **Create Subreddit:** Establish new subreddits, automatically assigning the creator as a moderator.
*   **Join Subreddit:** Allow users to become members of a specific subreddit.
*   **Add/Update Moderators:** Assign or promote users to moderator roles within a subreddit.
*   **Get All Subreddits:** Retrieve a paginated list of subreddits with search (by name, description) and sorting options (by creation date).
*   **Get Subreddits User Joined:** List subreddits that a specific user has joined, with sorting (by joined date) and search capabilities.
*   **Get Subreddit by ID:** Fetch details for a single subreddit using its unique ID.

## Installation

To set up the project locally, ensure you have Go and MongoDB installed.

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/EsanSamuel/Reddit_Clone.git
    cd Reddit_Clone/server
    ```
    (Note: The provided file paths indicate the controllers are within a `server` directory)

2.  **Install dependencies:**
    ```bash
    go mod tidy
    ```

3.  **Configure environment variables:**
    *   Database connection string for MongoDB.
    *   API keys/secrets for email sending and S3 (if applicable).
    *   JWT secrets for token generation.
    (Specific environment variables are not provided in the code, but are implied by the functionality.)

## Usage

To run the server and expose the API endpoints:

1.  **Start the MongoDB server.**

2.  **Run the Go application:**
    ```bash
    go run main.go
    ```
    (Assuming `main.go` exists in the `server` directory to bootstrap the application and define routes).

Once the server is running, you can interact with the API endpoints using tools like Postman, Insomnia, or a frontend application.

Example API endpoints (specific routes are not provided, but these are common patterns):

*   `POST /api/users/register` - CreateUser
*   `GET /api/users/verify?token={token}` - VerifyEmail
*   `POST /api/users/login` - Login
*   `GET /api/users` - GetAllUsers
*   `GET /api/users/{userId}` - GetUser
*   `POST /api/subreddits` - CreateSubreddit
*   `GET /api/subreddits` - GetSubReddit
*   `GET /api/subreddits/{id}` - GetSubRedditById

## Folder Structure Explanation

Based on the provided file paths and imports, the `server` directory likely contains the following structure:

```
reddit_clone/
└── server/
    ├── controllers/
    │   ├── subreddit_controller.go  # Handles subreddit-related API logic
    │   └── user_controller.go       # Handles user-related API logic (auth, profile)
    ├── database/                    # Likely contains MongoDB connection and collection setup
    ├── models/                      # Defines data structures (e.g., User, SubReddit, SubRedditMembers)
    ├── utils/                       # Contains utility functions (e.g., password hashing, token generation, email sending, S3 upload)
    ├── jobs/                        # Potentially background job definitions
    │   └── workers/                 # Worker functions for jobs (e.g., SendEmailQueue)
    └── main.go                      # (Inferred) Entry point for the Go application, handles routing
```

## Technologies

*   **Go**: The primary programming language for the backend.
*   **Gin Gonic**: A high-performance HTTP web framework for Go.
*   **MongoDB**: A NoSQL document database used for data storage, accessed via `go.mongodb.org/mongo-driver/v2`.
*   **BSON**: Binary JSON, the serialization format used by MongoDB.
*   **Bcrypt**: A password hashing function used for securing user credentials.
*   **AWS S3 (Implied)**: Used for storing uploaded files, such as user avatars, as suggested by the `UploadSingleFileToS3` utility.

## License

(No license information was provided in the project data.)