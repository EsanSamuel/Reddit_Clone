# 🚀 Reddit Clone Server 🌐

Welcome to the server-side repository for the Reddit Clone project! This powerful backend is built with Go and Gin Gonic, designed to support a robust social media platform akin to Reddit. It provides comprehensive functionalities for user management, subreddit creation and moderation, dynamic post handling (including media uploads), an advanced comment system, and cutting-edge AI-powered features for content summarization and semantic search.

This server is the backbone, ensuring all client-side interactions are seamless, secure, and enriched with intelligent features. Dive in to explore the core services that bring our Reddit clone to life! ✨

---

## 📝 Description

This repository contains the Go-based server for a Reddit clone, meticulously engineered to handle all core functionalities of a modern social content platform. It leverages MongoDB for data persistence, Gin Gonic for efficient routing, and integrates with AI services for intelligent content processing.

The server manages a wide array of features, from user authentication and profile management to complex interactions involving posts, comments, and subreddits. It's designed for scalability and performance, incorporating background job processing for tasks like AI embedding generation and email notifications. Secure file uploads to AWS S3 are also a key part of its media handling capabilities.

---

## ✨ Features

Our Reddit Clone server comes packed with a comprehensive set of features to power a dynamic and interactive community platform.

### 🧠 AI-Powered Content Intelligence

*   **Thread Summarization (Post & Comments):** 🗣️💬
    *   Automatically generates concise summaries of posts and their associated comment threads.
    *   Provides key points of the post, main opinions/arguments from comments, and identifies disagreements or recurring ideas.
    *   Ensures a neutral and informative tone for summarization.
*   **Semantic Search within Posts & Comments:** 🔍💡
    *   Allows users to perform natural language queries within a specific post's title and its comments.
    *   Utilizes AI embeddings and cosine similarity to find the most relevant content chunks.
    *   Processes content dynamically by chunking text, generating embeddings, scoring relevance, and then using an AI assistant to formulate an answer based on highly relevant information.

### 💬 Comment Management

*   **Create Comments:** ✍️
    *   Users can create new comments on posts, supporting a threaded discussion model.
    *   Automatically updates comment counts on the parent post and parent comment (for replies).
*   **Retrieve Post Comments:** 📄
    *   Fetches all comments associated with a specific post.
    *   Includes robust search capabilities by comment `type` or `content`.
    *   Supports sorting comments by `created_at` (ascending or descending).
    *   Implements pagination for efficient data loading.
*   **Retrieve Child Comments (Replies):** 🌲
    *   Enables fetching all replies to a specific parent comment.
    *   Offers search, sort, and pagination options similar to post comments.
*   **Get Comment by ID:** 🆔
    *   Retrieves a single comment using its unique identifier.

### 📰 Post Management

*   **Create Posts:** ➕
    *   Allows users to create new posts, supporting various content types.
    *   Handles multipart form data for file uploads (e.g., images, videos) to AWS S3.
    *   Automatically assigns a unique post ID and sets creation/update timestamps.
    *   Triggers a background job to generate AI embeddings for the new post, enhancing search and summarization capabilities.
    *   Increments the `posts_count` for the associated subreddit.
*   **Retrieve All Posts:** 📚
    *   Fetches a list of all posts across the platform.
    *   Includes powerful search functionality by post `title` or `content`.
    *   Supports sorting posts by `created_at` (ascending or descending).
    *   Provides pagination for managing large datasets.
*   **Retrieve Subreddit Posts:** 🏘️
    *   Gets all posts belonging to a specific subreddit.
    *   Offers search, sort, and pagination specific to subreddit content.
*   **Retrieve Posts by Tag:** 🏷️
    *   Aggregates and returns posts categorized by their tags, showing post counts per tag.
*   **Get Post by ID:** 🔍
    *   Retrieves a single post by its unique identifier, including its AI embeddings.
*   **Upvote Posts:** 👍
    *   Enables users to upvote posts.
    *   Ensures a user cannot upvote or downvote the same post multiple times.
    *   Increments the `up_vote` count for the post.
*   **Downvote Posts:** 👎
    *   Allows users to downvote posts.
    *   Ensures a user cannot upvote or downvote the same post multiple times.
    *   Decrements the `down_vote` count for the post.

### 🏘️ Subreddit Management

*   **Create Subreddits:** 🏡
    *   Users can create new subreddits with unique IDs, names, and descriptions.
    *   The creator is automatically assigned as the first 'MODERATOR' and added to the member list.
    *   Increments the `members_count` for the newly created subreddit.
*   **Join Subreddits:** 🤝
    *   Allows users to join existing subreddits.
    *   Prevents users from joining the same subreddit multiple times.
    *   Increments the `members_count` for the joined subreddit.
*   **Add Moderators:** 👑
    *   Enables designated users to add or promote other members to 'MODERATOR' roles within a subreddit.
    *   Updates the member's role if they are already a member, or adds them as a new moderator.
*   **Retrieve All Subreddits:** 🌍
    *   Fetches a list of all subreddits on the platform.
    *   Supports searching subreddits by `name` or `description`.
    *   Offers sorting by `created_at` and pagination.
*   **Retrieve User's Joined Subreddits:** 👤🏠
    *   Retrieves all subreddits a specific user has joined.
    *   Includes search and sort capabilities for the joined subreddits list.
*   **Get Subreddit by ID:** 📍
    *   Retrieves a single subreddit's details using its unique identifier.
*   **Leave Subreddit:** 🚪
    *   Allows a user to leave a subreddit, removing their membership record.

### 👤 User Management & Authentication

*   **Create User Account:** 📝
    *   Facilitates new user registration with password hashing (bcrypt) for security.
    *   Generates a unique `UserId` and sets creation/update timestamps.
    *   Prevents duplicate registrations using the same email address.
    *   Sends an email verification token to the user upon registration.
*   **Email Verification:** ✅📧
    *   Verifies a user's email address using a unique token.
    *   Updates the `email_verified` status and removes the verification token.
    *   Queues a follow-up email notification job upon successful verification.
*   **User Login:** 🔑
    *   Authenticates users by email and password, verifying email status and comparing hashed passwords.
    *   Generates secure JWT `access_token` and `refresh_token`.
    *   Sets these tokens as HttpOnly, Secure, and SameSiteNoneMode cookies for enhanced security.
    *   Returns a `UserDTO` containing essential user information and tokens.
*   **User Logout:** 🚶‍♀️
    *   Allows users to securely log out.
    *   Invalidates access and refresh tokens by clearing them from the database and expiring their respective cookies.
*   **Password Reset Request:** ❓
    *   Initiates a password reset process by generating a `reset_token`.
    *   Updates the user's record with the generated reset token.
*   **Password Reset:** 🔄
    *   Enables users to reset their password using a valid reset token.
    *   Hashes the new password and updates the user's password in the database.
    *   Removes the used reset token for security.
*   **Retrieve All Users:** 👥
    *   Fetches a list of all registered users.
    *   Supports searching by `first_name`, `last_name`, or `role`.
    *   Allows sorting by `created_at` and provides pagination.
*   **Get User by ID:** 🧍
    *   Retrieves a single user's profile details using their unique `userId`.
*   **Upload Avatar:** 🖼️
    *   Enables users to upload a profile avatar.
    *   Validates the uploaded file to ensure it's an image.
    *   Uploads the image to AWS S3 and updates the user's `avatar` URL in their profile.

---

## 🛠️ Installation

To get the Reddit Clone Server up and running locally, follow these steps.

### Prerequisites

Before you begin, ensure you have the following installed on your system:

*   **Go:** Version 1.18 or higher. You can download it from [golang.org](https://golang.org/dl/).
*   **MongoDB:** A running MongoDB instance (local or remote). You can download MongoDB Community Server from [mongodb.com](https://www.mongodb.com/try/download/community).
*   **Git:** For cloning the repository.

### Setup Steps

1.  **Clone the Repository:**
    Start by cloning the server repository to your local machine:

    ```bash
    git clone https://github.com/EsanSamuel/Reddit_Clone.git
    cd Reddit_Clone/server
    ```

2.  **Install Go Modules:**
    Navigate to the `server` directory and install the necessary Go modules:

    ```bash
    go mod tidy
    ```

3.  **Configure Environment Variables:**
    The server requires several environment variables for database connection, AI service integration, and AWS S3 configuration. Create a `.env` file in the `server` directory based on a template (if available, otherwise infer from the code).
    *Example `.env` (adjust values as per your setup):*

    ```ini
    MONGO_URI="mongodb://localhost:27017"
    MONGO_DB_NAME="reddit_clone_db"
    JWT_SECRET_KEY="your_jwt_secret"
    REFRESH_JWT_SECRET_KEY="your_refresh_jwt_secret"

    # AI Service Configuration (e.g., API keys, endpoints)
    AI_API_KEY="your_ai_service_api_key"
    AI_BASE_URL="https://api.example.com/ai" # Placeholder, actual URL will vary

    # AWS S3 Configuration
    AWS_REGION="your-aws-region"
    AWS_ACCESS_KEY_ID="your_aws_access_key_id"
    AWS_SECRET_ACCESS_KEY="your_aws_secret_access_key"
    AWS_S3_BUCKET_NAME="your-s3-bucket-name"

    # Email Service Configuration (for verification and notifications)
    # E.g., for SendGrid, Mailgun, or similar SMTP setup
    EMAIL_HOST="smtp.example.com"
    EMAIL_PORT="587"
    EMAIL_USERNAME="your_email_username"
    EMAIL_PASSWORD="your_email_password"
    EMAIL_FROM="no-reply@redditclone.com"
    ```

    **Note:** The exact environment variables will depend on the `config` and `utils` packages, which were not provided. The above is an educated guess based on common practices for such a project.

4.  **Run the MongoDB Database:**
    Ensure your MongoDB instance is running. For a local setup, you might start it via:

    ```bash
    # On Linux/macOS
    sudo systemctl start mongod # Or brew services start mongodb-community
    # On Windows
    "C:\Program Files\MongoDB\Server\X.X\bin\mongod.exe" --dbpath "C:\data\db"
    ```

5.  **Start the Server:**
    Once all dependencies are installed and environment variables are configured, you can start the server:

    ```bash
    go run main.go # Assuming your main entry point is main.go in the server directory
    ```

    The server should now be running, typically on `http://localhost:8080` (or whatever port is configured in `config`).

---

## 🚀 Usage

The Reddit Clone Server exposes a RESTful API that can be consumed by any client application. Below are examples of how to interact with some of the core functionalities.

**Base URL:** `http://localhost:8080/api/v1` (or your configured server address)

### 👤 User Endpoints

*   **Register a New User:**
    `POST /api/v1/users/signup`
    ```json
    {
        "firstName": "John",
        "lastName": "Doe",
        "email": "john.doe@example.com",
        "password": "StrongPassword123",
        "role": "USER"
    }
    ```

*   **Login User:**
    `POST /api/v1/users/login`
    ```json
    {
        "email": "john.doe@example.com",
        "password": "StrongPassword123"
    }
    ```

*   **Upload User Avatar:**
    `POST /api/v1/users/:userId/avatar`
    *   **Content-Type:** `multipart/form-data`
    *   **Form Field:** `avatar` (file)

### 🏘️ Subreddit Endpoints

*   **Create a New Subreddit:**
    `POST /api/v1/subreddits`
    ```json
    {
        "name": "MyCoolSubreddit",
        "description": "A place for cool things!",
        "creatorId": "user_id_of_creator"
    }
    ```

*   **Get All Subreddits (with Search and Pagination):**
    `GET /api/v1/subreddits?search=cool&page=1&sort=desc`

*   **Join a Subreddit:**
    `POST /api/v1/subreddits/join`
    ```json
    {
        "subredditId": "subreddit_id_to_join",
        "userId": "user_id_joining"
    }
    ```

### 📰 Post Endpoints

*   **Create a New Post (with optional files):**
    `POST /api/v1/posts`
    *   **Content-Type:** `multipart/form-data`
    *   **Form Fields:**
        *   `title`: "My New Post"
        *   `content`: "This is the content of my post."
        *   `type`: "text" (or "image", "link")
        *   `subredditId`: "id_of_subreddit"
        *   `userId`: "id_of_user"
        *   `tags`: ["tag1", "tag2"] (as `tags[0]=tag1&tags[1]=tag2` or similar array encoding)
        *   `files`: (multiple file inputs for images/videos)

*   **Get Posts for a Specific Subreddit:**
    `GET /api/v1/subreddits/:subreddit_id/posts?page=1&sort=desc`

*   **Upvote a Post:**
    `POST /api/v1/posts/upvote`
    ```json
    {
        "postId": "post_id_to_upvote",
        "userId": "user_id_upvoting"
    }
    ```

### 💬 Comment Endpoints

*   **Create a New Comment:**
    `POST /api/v1/comments`
    ```json
    {
        "postId": "post_id_to_comment_on",
        "userId": "user_id_commenting",
        "content": "Great post!",
        "parentId": "optional_parent_comment_id_for_replies"
    }
    ```

*   **Get Comments for a Post:**
    `GET /api/v1/posts/:post_id/comments?search=great&page=1`

### 🧠 AI Endpoints

*   **Summarize Post Threads:**
    `GET /api/v1/ai/threads-summary/:post_id`
    *   Returns a concise summary of the post and its comments.

*   **Semantic Search within Post Details:**
    `GET /api/v1/ai/post-search/:postId?query=What are the main arguments against this idea?`
    *   Returns an AI-generated answer based on the most semantically relevant content chunks from the post and its comments.

---

## 🗂️ Folder Structure Explanation

The server's codebase is organized to promote modularity, readability, and maintainability, following common Go project layouts.

```
reddit_clone/server/
├── controllers/
│   ├── ai_controller.go        # Handles AI-powered features like summarization and semantic search.
│   ├── comment_controller.go   # Manages comment creation, retrieval, and moderation logic.
│   ├── post_controller.go      # Manages post creation, retrieval, voting, and file uploads.
│   ├── subreddit_controller.go # Manages subreddit creation, membership, and moderation.
│   └── user_controller.go      # Handles user authentication, profile management, and account actions.
├── config/                     # (Implied) Contains application configuration, e.g., AI service setup, database credentials.
├── database/                   # (Implied) Manages database connection and collection interactions (e.g., MongoDB).
├── helper/                     # (Implied) Provides utility functions like text chunking and cosine similarity for AI features.
├── jobs/                       # (Implied) Directory for background jobs.
│   └── workers/                # (Implied) Contains worker functions for asynchronous tasks (e.g., AI embeddings, email sending).
├── models/                     # (Implied) Defines data structures (structs) for various entities like User, Post, Comment, Subreddit.
├── utils/                      # (Implied) Provides common utilities like password hashing, token generation, email services, and S3 file uploads.
└── main.go                     # (Implied) The main entry point of the Go application, setting up routes and starting the server.
```

**Explanation of Key Directories:**

*   **`controllers/`**: This directory houses the business logic for handling incoming HTTP requests. Each `*_controller.go` file is responsible for a specific domain (AI, comments, posts, subreddits, users) and contains functions that map to API endpoints. They interact with the `models` and `database` layers to perform operations.
*   **`config/`**: (Inferred from `ai_controller.go` imports) This likely contains configurations for external services, database connections, AI service credentials, and other application-wide settings.
*   **`database/`**: (Inferred from all controllers' imports) This package is responsible for establishing and managing the connection to the MongoDB database and providing interfaces to interact with various collections (e.g., `PostCollection`, `CommentCollection`, `UserCollection`).
*   **`helper/`**: (Inferred from `ai_controller.go` imports) This package provides general utility functions that support specific features, such as `ChunkText` for breaking down large text into smaller parts and `CosineSimilarity` for comparing vector embeddings, crucial for semantic search.
*   **`jobs/workers/`**: (Inferred from `post_controller.go` and `user_controller.go` imports) This package is designed to handle background tasks asynchronously. Examples include generating AI embeddings for new posts or sending verification/welcome emails, preventing these long-running operations from blocking the main request-response cycle.
*   **`models/`**: (Inferred from all controllers' imports) This package defines the data structures (Go structs) that represent the various entities in the application (e.g., `User`, `Post`, `Comment`, `SubReddit`, `PostUpvote`, `SubRedditMembers`). These structs are often mapped to MongoDB documents.
*   **`utils/`**: (Inferred from `post_controller.go` and `user_controller.go` imports) This package contains common, reusable utility functions that don't belong to any specific domain. This includes functions for password hashing (`HashPassword`), JWT token generation and management (`GenerateTokens`, `UpdateTokens`), email handling (`SendVerificationEmail`), and file upload to cloud storage like AWS S3 (`UploadFiles`, `UploadSingleFileToS3`, `IsFileImage`).

---

## 💻 Technologies

This server is built using a modern and robust stack:

*   **Go (Golang)**: The primary programming language, chosen for its performance, concurrency features, and strong type system. 🚀
*   **Gin Gonic**: A high-performance HTTP web framework for Go, used for building the RESTful API endpoints. ⚡
*   **MongoDB**: A NoSQL document database, providing a flexible and scalable data store for all application data. 💾
*   **AI/ML Integration**: Utilizes external AI services for advanced functionalities like:
    *   **Content Summarization**: Condensing long texts into brief, informative summaries. 📝
    *   **Text Embeddings**: Converting text into numerical vectors for semantic understanding. 📊
    *   **Cosine Similarity**: Measuring the similarity between text embeddings for search and relevance. 📐
*   **AWS S3**: Cloud storage service used for efficiently storing user-uploaded files, such as post media and avatars. ☁️
*   **Bcrypt**: A robust password hashing function employed for securely storing user passwords. 🔒
*   **JSON**: The standard data interchange format used for communication between the server and clients. 📨
*   **Goroutines & `context`**: Go's native concurrency mechanisms are leveraged for efficient handling of multiple requests and managing request lifecycles. 🔄
*   **Regular Expressions**: Used for advanced search functionalities within various resource endpoints. 🔍

---