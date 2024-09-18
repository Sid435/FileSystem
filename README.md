# File System API

This is a File System API built with Go, using PostgreSQL for database storage, Redis for caching, and AWS S3 for file storage. The API provides functionality for user authentication, file uploads, downloads, and deletions.

## Table of Contents

- [File System API](#file-system-api)
  - [Table of Contents](#table-of-contents)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
  - [Running the Application](#running-the-application)
  - [API Endpoints](#api-endpoints)
  - [Testing with Postman](#testing-with-postman)
    - [Using the Postman Collection](#using-the-postman-collection)
      - [Authentication](#authentication)
      - [File Operations](#file-operations)
      - [Error Handling](#error-handling)
    - [Note](#note)
  - [Note](#note-1)

## Prerequisites

- Go 1.16+
- PostgreSQL
- Redis
- AWS S3 account and credentials
- Postman (for testing)

## Installation

1. Clone the repository:
   ```
   git clone https://github.com/yourusername/FileSystem.git
   cd FileSystem
   ```

2. Install dependencies:
   ```
   go mod tidy
   ```
## Running the Application

1. Start the application:
   ```
   go run main.go
   ```

2. The application will start two servers:
   - File operations server on port 8080
   - Authentication server on port 9010 (or the port specified in the `PORT` environment variable)

## API Endpoints

- Authentication:
  - POST `/auth/signup`: User registration
  - POST `/auth/login`: User login

- File Operations:
  - POST `/files/upload`: Upload a file
  - GET `/files/get`: Get a pre-signed URL for file download
  - DELETE `/files/delete`: Delete a file (not included in the current Postman collection)

## Testing with Postman

We've provided a Postman collection to help you test the API endpoints. Follow these steps to use it:

1. Import the `File_System_Management.postman_collection.json` file into Postman.
2. The collection is organized into folders: Auth, File Operations, and Error Handling.

### Using the Postman Collection

#### Authentication

1. Sign Up
   - Request: POST `http://localhost:9010/auth/signup`
   - Body: x-www-form-urlencoded
     - username: testuser
     - password: testpassword
   - Expected Response: 200 OK with user details

2. Login
   - Request: POST `http://localhost:9010/auth/login`
   - Body: x-www-form-urlencoded
     - username: testuser
     - password: testpassword
   - Expected Response: 200 OK with JWT token
   - After successful login, copy the JWT token and update the `jwt_token` variable in the collection variables.

#### File Operations

3. Upload File
   - Request: POST `http://localhost:8080/files/upload`
   - Headers: 
     - Authorization: Bearer {{jwt_token}}
   - Body: form-data
     - Key: files
     - Value: Select file(s)
   - Expected Response: 200 OK with uploaded file URL(s)

4. Get File URL
   - Request: GET `http://localhost:9010/files/get?fileName=example.txt`
   - Headers:
     - Authorization: Bearer {{jwt_token}}
   - Expected Response: 200 OK with pre-signed download URL

#### Error Handling

5. Missing Authorization
   - Request: GET `http://localhost:8080/files/get?fileName=example.txt`
   - No Authorization header
   - Expected Response: 401 Unauthorized

6. Invalid File Name
   - Request: GET `http://localhost:8080/files/get?fileName=nonexistent.txt`
   - Headers:
     - Authorization: Bearer {{jwt_token}}
   - Expected Response: 404 Not Found

### Note

- The Postman collection uses environment variables. Make sure to set the `jwt_token` variable after login.
- The collection doesn't include a test for the file deletion endpoint. You may want to add this to your collection.
- Some URLs in the collection use port 9010 while others use 8080. Ensure your servers are running on these ports or update the collection accordingly.

## Note

Ensure that your AWS S3 bucket and Redis server are properly configured and accessible. Also, make sure your PostgreSQL database is set up and the connection string in `app.go` is correct.

For any issues or feature requests, please open an issue on the GitHub repository.