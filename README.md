# 🚀 Reddit Clone Backend API

Welcome to the backend API for the Reddit Clone project! This repository contains the server-side logic and API endpoints developed in Go, powering a robust social media platform akin to Reddit. It focuses on managing posts, comments, subreddits, user interactions, and even includes an AI-powered content summarization feature.

This `README.md` provides a comprehensive overview of the project's features, setup, and usage, generated directly from the provided source code.

## 📖 Description

This project serves as the backend infrastructure for a Reddit-like application, built with Go and leveraging MongoDB for data persistence. It provides a full suite of functionalities for users to create and interact with posts, manage subreddits, leave comments, and engage with content through upvoting and downvoting. A standout feature is the integration of an AI assistant to summarize post threads, enhancing user experience by providing quick insights into discussions. The API is designed to be scalable, offering robust data handling, search, sorting, and pagination capabilities across various resources.

## ✨ Features

The backend API offers a rich set of features, categorized by their respective domains:

### 🤖 AI-Powered Features

*   **Threads Summary**: Automatically summarizes a post and its associated comments using an AI assistant. This provides a concise overview of the discussion, including key post points, main opinions from comments, overall tone, and highlights disagreements or recurring ideas.

    ```go
    // Example AI summary prompt structure
    prompt := fmt.Sprintf(`You are an AI assistant. I will provide you with a post and its associated comments. Summarize the content...`)
    summary, err := config.Ai(prompt)
    ```

### 💬 Comment Management

*   **Create Comment**: Allows users to add new comments to a post. Supports nested comments by linking to a `ParentID`. Increments `comment_count` on the associated post and parent comment.
*   **Get Post Comments**: Retrieves all comments for a specific post.
    *   Supports **search functionality** on comment `type` and `content`.
    *   Allows **sorting** comments by `created_at` in ascending (`asc`) or descending (`desc`) order.
    *   Implements **pagination** to fetch comments in chunks (default 9 per page).
*   **Get Parent Comments**: Fetches replies to a specific parent comment.
    *   Supports **search functionality** on comment `type` and `content`.
    *   Allows **sorting** comments by `created_at` in ascending (`asc`) or descending (`desc`) order.
    *   Implements **pagination** to fetch comments in chunks (default 9 per page).
*   **Get Comment by ID**: Retrieves a single comment using its unique identifier.

### 📝 Post Management

*   **Create Post**: Enables users to create new posts.
    *   Supports different content types, including **multipart form data** for file uploads (e.g., images).
    *   Automatically increments the `posts_count` of the associated subreddit upon successful creation.
*   **Get All Posts**: Fetches a list of all posts.
    *   Supports **search functionality** on post `title` and `content`.
    *   Allows **sorting** posts by `created_at` in ascending (`asc`) or descending (`desc`) order.
    *   Implements **pagination** to fetch posts in chunks (default 9 per page).
*   **Get Subreddit Posts**: Retrieves posts belonging to a specific subreddit.
    *   Supports **search functionality** on post `title` and `content`.
    *   Allows **sorting** posts by `created_at` in ascending (`asc`) or descending (`desc`) order.
    *   Implements **pagination** to fetch posts in chunks (default 9 per page).
*   **Get Tag Posts**: Organizes and retrieves posts grouped by their associated tags, showing post counts per tag.
*   **Get Post by ID**: Fetches a single post using its unique identifier.
*   **Upvote Post**: Allows a user to upvote a post. Prevents duplicate upvotes or downvotes from the same user. Increments `up_vote` count on the post.
*   **Downvote Post**: Allows a user to downvote a post. Prevents duplicate downvotes or upvotes from the same user. Decrements `down_vote` count on the post.

### 🏛️ Subreddit Management

*   **Create Subreddit**: Enables users to create new subreddits.
    *   Automatically assigns the creator as a `MODERATOR` and adds them to the subreddit's members, incrementing the `members_count`.
*   **Join Subreddit**: Allows a user to join an existing subreddit. Prevents duplicate joins. Increments `members_count`.
*   **Add Moderators**: Grants moderator status to a user within a subreddit. If the user is already a member, their role is updated; otherwise, they are added as a new moderator.
*   **Get All Subreddits**: Retrieves a list of all available subreddits.
    *   Supports **search functionality** on subreddit `name` and `description`.
    *   Allows **sorting** subreddits by `created_at` in ascending (`asc`) or descending (`desc`) order.
    *   Implements **pagination** to fetch subreddits in chunks (default 9 per page).
*   **Get User Joined Subreddits**: Fetches subreddits that a specific user has joined.
    *   Supports **sorting** by `joined_at`.
    *   Supports **search functionality** on subreddit `name` and `description` within the joined subreddits.
*   **Get Subreddit by ID**: Retrieves a single subreddit using its unique identifier.
*   **Leave Subreddit**: Allows a user to leave a subreddit.

### 👤 User Management & Authentication

*   **Create User**: Registers a new user with hashed password storage.
    *   Generates and sends a **verification email** with a unique token.
    *   Ensures unique email addresses.
*   **Verify Email**: Verifies a user's email address using a provided token. Marks the user as `email_verified` and potentially queues a welcome email.
*   **Login**: Authenticates a user based on email and password.
    *   Checks if the email is verified.
    *   Generates and updates **authentication tokens** (access and refresh tokens).
*   **Reset Password Request**: Initiates a password reset process by generating a reset token and associating it with the user's email.
*   **Reset Password**: Allows a user to set a new password using a valid reset token.
*   **Get All Users**: Retrieves a list of all registered users.
    *   Supports **search functionality** on user `first_name`, `last_name`, and `role`.
    *   Allows **sorting** users by `created_at` in ascending (`asc`) or descending (`desc`) order.
    *   Implements **pagination** to fetch users in chunks (default 9 per page).
*   **Get User by ID**: Fetches a single user's details using their unique identifier.
*   **Upload Avatar**: Allows a user to upload an avatar image. Validates the file type to ensure it's an image and uploads it to S3. Updates the user's avatar URL.

## 🛠️ Installation

To set up and run this Go backend service locally, follow these steps.

### Prerequisites

*   **Go**: Ensure you have Go installed (version 1.18 or higher is recommended).
    *   [Download and Install Go](https://golang.org/doc/install)
*   **MongoDB**: This project uses MongoDB as its database. You will need a running MongoDB instance.
    *   [Install MongoDB Community Edition](https://docs.mongodb.com/manual/installation/)
    *   Alternatively, use a cloud-hosted MongoDB Atlas instance.

### Setup Steps

1.  **Clone the Repository**:
    ```bash
    git clone https://github.com/EsanSamuel/Reddit_Clone.git
    cd Reddit_Clone/server # Navigate to the server directory
    ```
    *(Note: The provided file paths indicate the Go code is within a `server` directory)*

2.  **Install Dependencies**:
    The project uses Go modules for dependency management. Navigate to the `server` directory (or wherever your `go.mod` file resides) and install the necessary packages.
    ```bash
    go mod tidy
    ```

3.  **Configuration**:
    Based on the source code, the application requires configuration for the database, a potential AI service, and an S3 bucket for file uploads. You will need to set up environment variables or a configuration file.
    *   **Database**: A MongoDB connection URI is essential.
        ```bash
        # Example environment variable (adjust as per your actual configuration)
        export MONGODB_URI="mongodb://localhost:27017/reddit_clone_db"
        ```
    *   **AI Service**: The `ai_controller.go` suggests an AI service is integrated via `config.Ai()`. This likely requires API keys or endpoint URLs.
        ```bash
        # Example environment variable (adjust as per your actual configuration)
        export AI_API_KEY="your_ai_service_api_key"
        export AI_ENDPOINT="your_ai_service_endpoint"
        ```
    *   **AWS S3**: For file uploads (e.g., user avatars, post files), an S3 bucket is used. This requires AWS credentials and bucket details.
        ```bash
        # Example environment variables (adjust as per your actual configuration)
        export AWS_ACCESS_KEY_ID="YOUR_AWS_ACCESS_KEY_ID"
        export AWS_SECRET_ACCESS_KEY="YOUR_AWS_SECRET_ACCESS_KEY"
        export AWS_REGION="your-aws-region"
        export AWS_S3_BUCKET_NAME="your-s3-bucket-name"
        ```
    *   **Email Service**: For user verification emails.
        ```bash
        # Example environment variables (adjust as per your actual configuration)
        export EMAIL_SERVICE_API_KEY="your_email_service_api_key"
        export SENDER_EMAIL="your_sender_email@example.com"
        ```
    *   **JWT Secrets**: For token generation.
        ```bash
        export JWT_SECRET="super_secret_jwt_key"
        export REFRESH_SECRET="super_secret_refresh_key"
        ```
    *(Note: Specific variable names and required values are inferred, you might need to consult other parts of the project's codebase or documentation for exact details.)*

4.  **Run the Application**:
    Once dependencies are installed and configuration is set, you can run the server.
    ```bash
    go run main.go # Assuming your entry point is main.go in the server directory
    ```
    The API should now be running, typically on `http://localhost:8080` (or as configured).

## 💡 Usage

This section outlines how to interact with the API endpoints. Since explicit routes are not provided, common RESTful API patterns are inferred from the controller names and functionalities.

### API Endpoints Overview

#### 🤖 AI Summary

*   **GET `/posts/{post_id}/summary`**
    *   Summarizes a post and its comments using AI.

#### 💬 Comment Endpoints

*   **POST `/comments`**
    *   **Description**: Creates a new comment.
    *   **Body**: `{"post_id": "...", "content": "...", "user_id": "...", "parent_id": "..." (optional)}`
*   **GET `/posts/{post_id}/comments`**
    *   **Description**: Retrieves comments for a specific post.
    *   **Query Params**: `search={query}`, `sort=asc|desc`, `page={number}`
*   **GET `/comments/{parent_id}/replies`**
    *   **Description**: Retrieves replies to a specific parent comment.
    *   **Query Params**: `search={query}`, `sort=asc|desc`, `page={number}`
*   **GET `/comments/{id}`**
    *   **Description**: Retrieves a single comment by ID.

#### 📝 Post Endpoints

*   **POST `/posts`**
    *   **Description**: Creates a new post. Supports `application/json` or `multipart/form-data` for file uploads.
    *   **Body (JSON)**: `{"subreddit_id": "...", "title": "...", "content": "...", "type": "text", "tags": ["tag1", "tag2"], "user_id": "..."}`
    *   **Body (Multipart)**: `subreddit_id`, `title`, `content`, `type`, `tags` (as form fields), `files` (as file input)
*   **GET `/posts`**
    *   **Description**: Retrieves all posts.
    *   **Query Params**: `search={query}`, `sort=asc|desc`, `page={number}`
*   **GET `/subreddits/{subreddit_id}/posts`**
    *   **Description**: Retrieves posts within a specific subreddit.
    *   **Query Params**: `search={query}`, `sort=asc|desc`, `page={number}`
*   **GET `/posts/tags`**
    *   **Description**: Retrieves posts grouped by tags.
*   **GET `/posts/{id}`**
    *   **Description**: Retrieves a single post by ID.
*   **POST `/posts/{post_id}/upvote`**
    *   **Description**: Upvotes a post.
    *   **Body**: `{"user_id": "..."}`
*   **POST `/posts/{post_id}/downvote`**
    *   **Description**: Downvotes a post.
    *   **Body**: `{"user_id": "..."}`

#### 🏛️ Subreddit Endpoints

*   **POST `/subreddits`**
    *   **Description**: Creates a new subreddit.
    *   **Body**: `{"name": "...", "description": "...", "creator_id": "..."}`
*   **POST `/subreddits/join`**
    *   **Description**: Allows a user to join a subreddit.
    *   **Body**: `{"user_id": "...", "subreddit_id": "..."}`
*   **POST `/subreddits/moderators`**
    *   **Description**: Adds a user as a moderator to a subreddit.
    *   **Body**: `{"user_id": "...", "subreddit_id": "..."}`
*   **GET `/subreddits`**
    *   **Description**: Retrieves all subreddits.
    *   **Query Params**: `search={query}`, `sort=asc|desc`, `page={number}`
*   **GET `/users/{user_id}/subreddits/joined`**
    *   **Description**: Retrieves subreddits a user has joined.
    *   **Query Params**: `search={query}`, `sort=asc|desc`
*   **GET `/subreddits/{id}`**
    *   **Description**: Retrieves a single subreddit by ID.
*   **DELETE `/users/{user_id}/subreddits/leave`**
    *   **Description**: Allows a user to leave a subreddit.

#### 👤 User Endpoints

*   **POST `/users/register`**
    *   **Description**: Registers a new user.
    *   **Body**: `{"first_name": "...", "last_name": "...", "email": "...", "password": "..."}`
*   **GET `/users/verify-email?token={token}`**
    *   **Description**: Verifies a user's email with a token.
*   **POST `/users/login`**
    *   **Description**: Authenticates a user.
    *   **Body**: `{"email": "...", "password": "..."}`
*   **POST `/users/forgot-password-request`**
    *   **Description**: Initiates a password reset.
    *   **Body**: `{"email": "..."}`
*   **POST `/users/reset-password?token={token}`**
    *   **Description**: Resets a user's password.
    *   **Body**: `{"password": "..."}`
*   **GET `/users`**
    *   **Description**: Retrieves all users.
    *   **Query Params**: `search={query}`, `sort=asc|desc`, `page={number}`
*   **GET `/users/{userId}`**
    *   **Description**: Retrieves a single user by ID.
*   **POST `/users/{userId}/avatar`**
    *   **Description**: Uploads a user's avatar image.
    *   **Body**: `avatar` (multipart file input)

## 📁 Folder Structure Explanation

The provided file paths indicate a clear separation of concerns within the `server` directory:

```
reddit_clone/
└── server/
    ├── controllers/
    │   ├── ai_controller.go          # Handles AI-powered summarization logic
    │   ├── comment_controller.go     # Manages comment-related API endpoints
    │   ├── post_controller.go        # Manages post-related API endpoints
    │   ├── subreddit_controller.go   # Manages subreddit-related API endpoints
    │   └── user_controller.go        # Manages user authentication and profile API endpoints
    ├── database/                     # (Inferred) Contains database connection and collection setup
    ├── models/                       # (Inferred) Defines data structures (structs) for the application
    ├── utils/                        # (Inferred) Provides utility functions like file uploads, password hashing, token generation, email sending
    ├── config/                       # (Inferred) Manages application configuration, including AI service integration
    ├── jobs/                         # (Inferred) Contains background job workers (e.g., email queue)
    └── main.go                       # (Inferred) The main application entry point (not provided, but standard for Go apps)
```

*   **`controllers/`**: This directory holds the Gin Gonic handler functions that define the API endpoints. Each `_controller.go` file groups related API logic, making the codebase organized and maintainable.
*   **`database/`**: (Inferred) Based on `database.PostCollection.FindOne`, this package is responsible for initializing the database connection (MongoDB) and providing access to various collections.
*   **`models/`**: (Inferred) Based on `var post models.Post`, this package likely contains the Go structs that define the data models for entities like `Post`, `Comment`, `User`, `Subreddit`, etc., often mapping directly to database schema.
*   **`utils/`**: (Inferred) The code snippets show calls to `utils.HashPassword`, `utils.GenerateTokens`, `utils.SendVerificationEmail`, `utils.UploadFiles`, `utils.IsFileImage`, `utils.UploadSingleFileToS3`. This package contains helper functions that perform common tasks across the application.
*   **`config/`**: (Inferred) The `config.Ai(prompt)` call suggests this package manages application configurations, including credentials and client setup for external services like the AI assistant.
*   **`jobs/`**: (Inferred) `workers.SendEmailQueue` implies a background job processing system, likely for tasks such as sending emails asynchronously.

## 💻 Technologies

This backend API is built using a modern and efficient technology stack:

*   **Go (Golang)**: The primary programming language, chosen for its performance, concurrency features, and robust standard library.
*   **Gin Gonic**: A high-performance HTTP web framework for Go, used for building the RESTful API endpoints.
*   **MongoDB**: A NoSQL document database, utilized for flexible and scalable data storage.
    *   `go.mongodb.org/mongo-driver/v2/bson`: The official MongoDB Go Driver's BSON package for encoding/decoding BSON documents.
*   **BSON (`go.mongodb.org/mongo-driver/v2/bson`)**: Binary JSON format used by MongoDB.
*   **`context`**: Go's standard library package for managing request-scoped values, deadlines, and cancellation signals.
*   **`time`**: Go's standard library for time-related operations (e.g., `time.Now()`, `time.Second`).
*   **`regexp`**: Go's standard library for regular expression matching, used for robust search functionalities.
*   **`encoding/json`**: Go's standard library for JSON encoding and decoding, used for API responses and AI prompts.
*   **`fmt`**: Go's standard library for formatted I/O, used for string formatting (e.g., AI prompt creation) and debugging.
*   **`net/http`**: Go's standard library for HTTP client and server implementations.
*   **`strconv`**: Go's standard library for string conversions, particularly for converting query parameters to integers (e.g., page numbers).
*   **`strings`**: Go's standard library for string manipulation (e.g., `TrimSpace`, `HasPrefix`).
*   **`golang.org/x/crypto/bcrypt`**: Used for secure password hashing and comparison.
*   **AWS S3 (inferred via `utils.UploadSingleFileToS3`)**: Cloud storage service for handling file uploads (e.g., user avatars, post media).
*   **AI Service (inferred via `config.Ai`)**: An external AI service integrated for content summarization.

## 📄 License

The license information for this project is not specified in the provided source code snippets. Please refer to the main project repository or contact the author for licensing details.