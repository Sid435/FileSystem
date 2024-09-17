# File System API

This is a File System API built with Go, using PostgreSQL for database storage, Redis for caching, and AWS S3 for file storage. The API provides functionality for user authentication, file uploads, downloads, and deletions.

## Table of Contents

- [File System API](#file-system-api)
  - [Table of Contents](#table-of-contents)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
  - [Configuration](#configuration)
  - [Running the Application](#running-the-application)
  - [API Endpoints](#api-endpoints)
  - [Testing with Postman](#testing-with-postman)
    - [1. User Signup](#1-user-signup)
    - [2. User Login](#2-user-login)
    - [3. Upload File](#3-upload-file)
    - [4. Get Pre-signed URL](#4-get-pre-signed-url)
    - [5. Delete File](#5-delete-file)
  - [Note](#note)

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

## Configuration

1. Create a `.env` file in the root directory with the following content:
   ```
   JWT_SECRET=your_jwt_secret
   AWS_ACCESS_KEY_ID=your_aws_access_key
   AWS_SECRET_ACCESS_KEY=your_aws_secret_key
   AWS_REGION=your_aws_region
   AWS_BUCKET_NAME=your_s3_bucket_name
   PORT=9010
   ```

2. Update the database connection string in `app.go` if necessary.

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
  - DELETE `/files/delete`: Delete a file

## Testing with Postman

### 1. User Signup

- Method: POST
- URL: `http://localhost:9010/auth/signup`
- Body (raw JSON):
  ```json
  {
    "username": "testuser",
    "name": "Test User",
    "age": "25",
    "password": "testpassword"
  }
  ```
- Expected Response: 200 OK with user details

### 2. User Login

- Method: POST
- URL: `http://localhost:9010/auth/login`
- Body (raw JSON):
  ```json
  {
    "username": "testuser",
    "password": "testpassword"
  }
  ```
- Expected Response: 200 OK with JWT token

### 3. Upload File

- Method: POST
- URL: `http://localhost:8080/files/upload`
- Headers:
  - Authorization: Bearer <jwt_token>
- Body (form-data):
  - Key: files
  - Value: Select file(s)
- Expected Response: 200 OK with uploaded file URL(s)

### 4. Get Pre-signed URL

- Method: GET
- URL: `http://localhost:8080/files/get?fileName=example.txt`
- Headers:
  - Authorization: Bearer <jwt_token>
- Expected Response: 200 OK with pre-signed download URL

### 5. Delete File

- Method: DELETE
- URL: `http://localhost:8080/files/delete?fileName=example.txt`
- Headers:
  - Authorization: Bearer <jwt_token>
- Expected Response: 200 OK with success message

Remember to replace `<jwt_token>` with the actual token received from the login endpoint.

## Note

Ensure that your AWS S3 bucket and Redis server are properly configured and accessible. Also, make sure your PostgreSQL database is set up and the connection string in `app.go` is correct.

For any issues or feature requests, please open an issue on the GitHub repository.