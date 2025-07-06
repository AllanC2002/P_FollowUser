# FollowUser Microservice

## Project Overview

P_FollowUser is a Go-based microservice designed to handle user-following functionality. It allows users to follow and unfollow other users within an application. The service uses the Gin web framework for routing, GORM as an ORM for MySQL database interactions, and JWT for securing its API endpoint.

## Folder Structure

The project is organized into the following main directories:

*   **`.github/`**: Contains GitHub Actions workflows.
    *   `workflows/publish.yml`: Defines a CI/CD pipeline for building the application's Docker image, publishing it to Docker Hub, and deploying it to an AWS EC2 instance on pushes/PRs to the `qa` branch or when new version tags are pushed.
*   **`connection/`**: Manages the database connection.
    *   `mysql.go`: Contains the logic to establish a connection to a MySQL database using credentials and connection details provided via environment variables.
*   **`functions/`**: Houses the core business logic of the service.
    *   `follow.go`: Implements the `FollowUser` function, which includes validating user profiles, checking existing follow relationships, and creating or updating follow records in the database.
*   **`models/`**: Defines the data structures (GORM models) that map to database tables.
    *   `models.go`: Contains the `Profile` and `Followers` structs, representing user profiles and their follow relationships.
*   **`tests/`**: Includes test scripts for the application.
    *   `test.py`: A Python script used for integration testing of the `/follow` endpoint. It demonstrates how to obtain a JWT from a separate authentication service and then use that token to make an authenticated request to this microservice.
*   **`dockerfile`**: Instructions for building the Docker image for the application.
*   **`go.mod`**: Go modules file, defining the project's module path and dependencies.
*   **`go.sum`**: Contains the checksums of direct and indirect dependencies.
*   **`main.go`**: The main entry point for the application. It initializes the database connection, sets up the Gin HTTP router, defines the API endpoints, and handles JWT-based authentication.
*   **`.gitignore`**: Specifies intentionally untracked files that Git should ignore.

## Backend Design Pattern

The microservice employs a **service-oriented approach** within its monolithic structure.
*   **Routing & Request Handling:** `main.go` uses the Gin framework to define routes and handle incoming HTTP requests. It acts as the primary interface for the service.
*   **Business Logic:** The `functions/` directory encapsulates the core business logic. For instance, `functions/follow.go` contains the logic for the user-following feature. This promotes separation of concerns, making the `main.go` file leaner and focused on request/response handling and middleware.
*   **Data Modeling:** The `models/` directory defines GORM models, which represent the database schema and facilitate database operations.
*   **Database Interaction:** The `connection/` package is responsible for establishing and managing the database connection.

This structure is common in Go web services, where HTTP handlers in `main.go` (or a dedicated handlers package) call specific functions from service packages (like `functions/`) to perform operations.

## Communication Architecture

*   **RESTful API:** The service exposes its functionality via a RESTful HTTP API. Currently, it has one primary endpoint for user-following.
*   **Synchronous Communication:** API requests are handled synchronously. The client sends a request and waits for a response.
*   **Database:** The service interacts with a MySQL database for data persistence, managed by GORM.
*   **Authentication:** Endpoint security is provided by **JSON Web Tokens (JWT)**.
    *   Clients must include a JWT in the `Authorization` header with the `Bearer` scheme (e.g., `Authorization: Bearer <your_jwt_token>`).
    *   The JWT is expected to contain a `user_id` claim, which identifies the user performing the action (the follower).

## Folder Pattern

The project utilizes a **component-based** or **feature-based** folder pattern:
*   `connection/`: Groups all database connection logic.
*   `functions/`: Groups business logic units or service functions.
*   `models/`: Groups data model definitions.
*   `tests/`: Groups test files.

This organization helps in maintaining a clear separation of concerns and makes it easier to locate code related to specific functionalities.

## Endpoint Instructions

### Follow a User

*   **Endpoint:** `POST /follow`
*   **Purpose:** Allows an authenticated user (follower) to follow another user (the one being followed).
*   **Authentication:**
    *   Required: Yes
    *   Method: JWT Bearer Token.
    *   Header: `Authorization: Bearer <token>`
    *   The JWT must contain a `user_id` claim (float64) representing the ID of the user initiating the follow request.
*   **Request Body (JSON):**
    ```json
    {
      "id_following": <integer> // The ID of the user to be followed
    }
    ```
*   **Responses:**
    *   **`201 Created`**: Successfully followed the user.
        ```json
        {
          "message": "Followed successfully"
        }
        ```
    *   **`200 OK`**: If the user is already following the target or if a previously "unfollowed" status is re-activated.
        ```json
        // Option 1
        {
          "message": "Already following"
        }
        // Option 2
        {
          "message": "Follow re-activated"
        }
        ```
    *   **`400 Bad Request`**:
        *   If the `id_follower` (from token) and `id_following` (from request body) are the same.
            ```json
            {
              "error": "You cannot follow yourself"
            }
            ```
        *   If the request body is not valid JSON.
            ```json
            {
              "error": "Invalid JSON"
            }
            ```
    *   **`401 Unauthorized`**:
        *   If the `Authorization` header is missing or not in "Bearer <token>" format.
            ```json
            {
              "error": "Token missing or invalid"
            }
            ```
        *   If the provided JWT is invalid, expired, or malformed.
            ```json
            {
              "error": "Invalid or expired token"
            }
            ```
        *   If the token claims cannot be parsed correctly.
            ```json
            {
              "error": "Invalid token claims"
            }
            ```
        *   If the `user_id` claim is missing or not a number in the token.
            ```json
            {
              "error": "user_id not found in token"
            }
            ```
    *   **`404 Not Found`**:
        *   If the follower's profile (derived from `user_id` in token) is not found or inactive.
            ```json
            {
              "error": "Follower profile not found or inactive"
            }
            ```
        *   If the target user's profile (`id_following`) is not found or inactive.
            ```json
            {
              "error": "Following profile not found or inactive"
            }
            ```
    *   **`500 Internal Server Error`**: If there's an issue with the database or any other unexpected server-side error. The error message will provide more details.
        ```json
        {
          "error": "<specific database or server error message>"
        }
        ```
