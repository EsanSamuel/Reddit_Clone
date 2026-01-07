# 🚀 Reddit Clone Server Backend

Welcome to the backend server for a Reddit-style social media platform! This Go-based API powers core functionalities such as managing subreddits, posts, comments, user authentication, and even leverages AI for intelligent content summarization.

---

## 📝 Description

This project serves as the robust backend for a Reddit-like application, built with Go and the Gin web framework. It handles all data persistence with MongoDB, offering a comprehensive suite of features for community management, user interaction, content creation, and sophisticated content moderation through AI-powered summaries. The server is designed for scalability and maintainability, ensuring a smooth and responsive experience for users.

---

## ✨ Features

This server provides a rich set of functionalities, covering various aspects of a social media platform:

### 🤖 AI-Powered Content Summarization
*   **Threads Summary**: Generate a concise and informative summary of a post and its associated comments using an AI assistant. This includes a post summary, comment thread summary (main opinions, arguments, insights), tone analysis, and optional highlights of disagreements or recurring ideas.

### 💬 Comment Management
*   **Create Comments**: Users can create new comments on posts.
*   **Nested Comments**: Supports replying to existing comments, creating threaded discussions.
*   **Retrieve Comments**:
    *   Fetch all comments for a specific post.
    *   Fetch child comments for a specific parent comment.
    *   Retrieve a single comment by its ID.
*   **Search Comments**: Filter comments by `type` or `content` using a case-insensitive regex search.
*   **Sort Comments**: Order comments by `created_at` in ascending or descending order.
*   **Paginate Comments**: Retrieve comments in chunks for efficient loading.
*   **Comment Count Updates**: Automatically increments comment counts for posts and parent comments upon new comment creation.

### 📰 Post Management
*   **Create Posts**: Users can create new posts, including support for multipart form data uploads (e.g., images, files).
*   **Retrieve Posts**:
    *   Fetch all posts.
    *   Fetch posts belonging to a specific subreddit.
    *   Fetch posts categorized by tags, providing counts for each tag.
    *   Retrieve a single post by its ID.
*   **Search Posts**: Filter posts by `title` or `content` using a case-insensitive regex search.
*   **Sort Posts**: Order posts by `created_at` in ascending or descending order.
*   **Paginate Posts**: Retrieve posts in chunks for efficient loading.
*   **Upvote Posts**: Allow users to upvote posts, incrementing the post's `up_vote` count. Includes checks to prevent duplicate upvotes or upvoting after downvoting.
*   **Downvote Posts**: Allow users to downvote posts, decrementing the post's `down_vote` count. Includes checks to prevent duplicate downvotes or downvoting after upvoting.
*   **Post Count Updates**: Automatically increments post counts for subreddits upon new post creation.

### 🏛️ Subreddit Management
*   **Create Subreddits**: Users can create new communities (subreddits). The creator is automatically added as a `MODERATOR`.
*   **Join Subreddits**: Users can join existing subreddits.
*   **Add Moderators**: Subreddit creators can promote members to moderators.
*   **Leave Subreddits**: Users can leave subreddits they have joined.
*   **Retrieve Subreddits**:
    *   Fetch all subreddits.
    *   Fetch subreddits joined by a specific user.
    *   Retrieve a single subreddit by its ID.
*   **Search Subreddits**: Filter subreddits by `name` or `description` using a case-insensitive regex search.
*   **Sort Subreddits**: Order subreddits by `created_at` or `joined_at` in ascending or descending order.
*   **Paginate Subreddits**: Retrieve subreddits in chunks for efficient loading.
*   **Member Count Updates**: Automatically increments member counts for subreddits upon new member joining.

### 👤 User Authentication & Management
*   **User Registration**: Create new user accounts with password hashing.
*   **Email Verification**: New users receive a verification email to activate their account.
*   **User Login**: Authenticate users, generate access and refresh tokens, and set HTTP-only cookies.
*   **User Logout**: Invalidate user sessions and clear authentication cookies.
*   **Password Reset**:
    *   Request a password reset via email.
    *   Reset password using a generated token.
*   **Retrieve Users**:
    *   Fetch all user accounts.
    *   Retrieve a single user by their ID.
*   **Search Users**: Filter users by `first_name`, `last_name`, or `role` using a case-insensitive regex search.
*   **Sort Users**: Order users by `created_at` in ascending or descending order.
*   **Paginate Users**: Retrieve users in chunks for efficient loading.
*   **Avatar Upload**: Allow users to upload a profile avatar, with image validation.

---

## 🛠️ Installation

To get this server up and running on your local machine, follow these steps:

1.  **Clone the Repository**:
    ```bash
    git clone https://github.com/EsanSamuel/Reddit_Clone.git
    cd Reddit_Clone/server
    ```

2.  **Install Dependencies**:
    The project uses Go modules. Ensure you have Go installed (version 1.18 or higher recommended).
    ```bash
    go mod tidy
    ```

3.  **Configure Environment Variables**:
    This server relies on environment variables for sensitive information like database connection strings, AI service API keys, and S3 credentials. While specific variables aren't detailed here, you'll typically need to set up:
    *   MongoDB connection URI.
    *   AI service API key (e.g., `GEMINI_API_KEY`).
    *   JWT secret keys for token generation.
    *   AWS S3 credentials (access key, secret key, region) for file uploads.
    *   Email service credentials for sending verification/reset emails.

    Create a `.env` file in the `server` directory or set them directly in your environment.

---

## 🏃 Usage

Once installed and configured, you can run the server:

1.  **Run the Server**:
    Assuming your main application entry point is `main.go` in the `server` directory:
    ```bash
    go run main.go
    ```
    The server will typically start on a port like `8080` (or as configured).

2.  **API Endpoints**:
    The server exposes various API endpoints accessible via HTTP requests. Examples include:
    *   `POST /api/v1/users/signup` to create a new user.
    *   `POST /api/v1/users/login` for user authentication.
    *   `POST /api/v1/subreddits` to create a new subreddit.
    *   `GET /api/v1/subreddits` to fetch all subreddits (with search/sort/pagination).
    *   `POST /api/v1/posts` to create a new post.
    *   `GET /api/v1/posts/{post_id}/comments` to get comments for a post (with search/sort/pagination).
    *   `GET /api/v1/ai/summary/{post_id}` to get an AI summary of a post's threads.

    You can interact with these endpoints using tools like Postman, Insomnia, or directly from your frontend application.

---

## 📁 Folder Structure Explanation

The provided files give insight into the `server/controllers` directory. A typical Go backend project structure would likely look something like this:

```
reddit_clone/
├── server/
│   ├── main.go               # Entry point for the application
│   ├── config/               # Configuration files (e.g., database, AI, environment)
│   │   ├── config.go         # Example: AI client initialization
│   ├── database/             # Database connection and collection initialization
│   │   ├── database.go       # Example: MongoDB client setup
│   ├── models/               # Data structures (structs) for various entities (User, Post, Comment, Subreddit, etc.)
│   │   ├── user.go
│   │   ├── post.go
│   │   ├── comment.go
│   │   ├── subreddit.go
│   │   └── ...
│   ├── controllers/          # API handlers (functions that process HTTP requests)
│   │   ├── ai_controller.go        # Handles AI-related requests (e.g., summarization)
│   │   ├── comment_controller.go   # Manages comment creation, retrieval, etc.
│   │   ├── post_controller.go      # Handles post creation, retrieval, voting, etc.
│   │   ├── subreddit_controller.go # Manages subreddit creation, joining, members, etc.
│   │   └── user_controller.go      # Handles user authentication, profile management, etc.
│   ├── routes/               # Defines the API endpoints and maps them to controllers
│   │   └── routes.go
│   ├── middleware/           # Custom middleware for authentication, logging, etc.
│   │   └── auth_middleware.go
│   ├── utils/                # Utility functions (e.g., password hashing, token generation, file uploads, email sending)
│   │   ├── helpers.go
│   │   ├── aws_s3.go
│   │   ├── email.go
│   │   └── ...
│   ├── jobs/                 # Background jobs or worker processes
│   │   └── workers/          # Worker functions (e.g., sending emails)
│   │       └── email_worker.go
│   ├── go.mod                # Go module definition
│   └── go.sum                # Go module checksums
└── (frontend/client/ etc.)   # Other parts of the monorepo, if applicable
```

---

## 💻 Technologies

The project primarily uses the following technologies:

*   **Go**: The main programming language for the backend.
*   **Gin Gonic**: A high-performance HTTP web framework for Go.
*   **MongoDB**: A NoSQL document database used for data storage.
*   **Go MongoDB Driver**: Official MongoDB driver for Go.
*   **BSON**: Binary JSON format used by MongoDB.
*   **Context Package**: For managing request-scoped values, cancellation, and deadlines.
*   **bcrypt**: For secure password hashing.
*   **JWT (JSON Web Tokens)**: Likely used for user authentication and authorization (implied by token generation/refresh).
*   **AWS S3**: Utilized for storing user-uploaded files (e.g., post files, user avatars).
*   **AI Integration**: An external AI service (implied by `config.Ai`) is used for generating content summaries.
*   **`net/http`**: Standard Go package for HTTP client and server functionality.
*   **`encoding/json`**: Standard Go package for JSON encoding and decoding.
*   **`regexp`**: Standard Go package for regular expression matching, used in search functionalities.

---

## 📜 License

The license information was not provided with the project data. Please refer to the project repository for details on licensing.