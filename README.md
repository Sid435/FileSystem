# File System Service

This project is a file storage and management system built using Go, Gin framework, PostgreSQL, Redis, and AWS S3. It allows users to securely upload files, retrieve pre-signed URLs for downloading files, and manage file metadata. The system uses JWT-based authentication and caches file metadata with Redis for better performance.

## Features

- **User Authentication**: Signup and login using secure password hashing (bcrypt) and JWT tokens.
- **File Upload**: Upload multiple files to AWS S3.
- **Pre-Signed URLs**: Retrieve time-limited pre-signed URLs to download files.
- **File Metadata**: Stores file metadata such as filename, size, and content type in PostgreSQL.
- **Caching**: Uses Redis to cache file metadata for faster access.

## Table of Contents

- [File System Service](#file-system-service)
  - [Features](#features)
  - [Table of Contents](#table-of-contents)
  - [Installation](#installation)
  - [Environment Variables](#environment-variables)
  - [Usage](#usage)
    - [1. Signup](#1-signup)
    - [2. Login](#2-login)
    - [3. Upload Files](#3-upload-files)
    - [4. Get Pre-Signed URLs](#4-get-pre-signed-urls)
  - [Docker](#docker)
    - [Docker Services](#docker-services)

## Installation

1. Clone the repository:

    ```bash
    git clone https://github.com/sid/FileSystem.git
    cd FileSystem
    ```

2. Install Go dependencies:

    ```bash
    go mod download
    ```

3. Set up PostgreSQL and Redis using Docker:

    ```bash
    docker-compose up -d
    ```

4. Build the Go application:

    ```bash
    go build -o main .
    ```

5. Run the application:

    ```bash
    ./main
    ```

## Environment Variables

Create a `.env` file in the root directory with the following variables:

```bash
# AWS Configuration
AWS_ACCESS_KEY_ID=your_aws_access_key_id
AWS_SECRET_ACCESS_KEY=your_aws_secret_access_key
AWS_REGION=your_aws_region
AWS_BUCKET_NAME=your_s3_bucket_name

# Database Configuration
DB_HOST=db
DB_PORT=5432
DB_USER=siddharth
DB_PASSWORD=siddharth
DB_NAME=file_system

# Redis Configuration
REDIS_HOST=redis
REDIS_PORT=6379
```

## Usage

### 1. Signup

- **Endpoint**: `/auth/signup`
- **Method**: `POST`
- **Description**: Registers a new user to the system.
- **Request Body**:

    ```json
    {
      "username": "your_username",
      "password": "your_password"
    }
    ```

- **Response**:

    ```json
    {
      "ID": 1,
      "Username": "your_username"
    }
    ```

- **Example**:

    ```bash
    curl -X POST http://localhost:8080/auth/signup \
    -H "Content-Type: application/json" \
    -d '{"username": "testuser", "password": "password123"}'
    ```

### 2. Login

- **Endpoint**: `/auth/login`
- **Method**: `POST`
- **Description**: Authenticates a user and returns a JWT token.
- **Request Body**:

    ```json
    {
      "username": "your_username",
      "password": "your_password"
    }
    ```

- **Response**:

    ```json
    {
      "token": "your_jwt_token"
    }
    ```

- **Example**:

    ```bash
    curl -X POST http://localhost:8080/auth/login \
    -H "Content-Type: application/json" \
    -d '{"username": "testuser", "password": "password123"}'
    ```

### 3. Upload Files

- **Endpoint**: `/files/upload`
- **Method**: `POST`
- **Description**: Uploads multiple files to AWS S3. Requires JWT authentication.
- **Headers**: 
  - `Authorization`: Bearer `your_jwt_token`
  - `Content-Type`: multipart/form-data
- **Form Data**:
  - `files`: One or more files to be uploaded.

- **Response**:

    ```json
    {
      "urls": [
        "https://your_bucket.s3.amazonaws.com/testuser/file1.png",
        "https://your_bucket.s3.amazonaws.com/testuser/file2.jpg"
      ]
    }
    ```

- **Example**:

    ```bash
    curl -X POST http://localhost:8080/files/upload \
    -H "Authorization: Bearer your_jwt_token" \
    -F "files=@path_to_file1" \
    -F "files=@path_to_file2"
    ```

### 4. Get Pre-Signed URLs

- **Endpoint**: `/files/get`
- **Method**: `GET`
- **Description**: Generates a pre-signed URL for a file stored in S3. Requires JWT authentication.
- **Headers**: 
  - `Authorization`: Bearer `your_jwt_token`
- **Query Params**:
  - `fileName`: The name of the file to generate the pre-signed URL for.

- **Response**:

    ```json
    {
      "download_url": "https://your_bucket.s3.amazonaws.com/testuser/file1.png?X-Amz-Signature=..."
    }
    ```

- **Example**:

    ```bash
    curl -X GET http://localhost:8080/files/get?fileName=file1.png \
    -H "Authorization: Bearer your_jwt_token"
    ```

## Docker

The project includes a `docker-compose.yml` file for easy setup of PostgreSQL, Redis, and the application. To start the services:

```bash
docker-compose up --build
```

This will set up:

- **PostgreSQL** for file metadata storage.
- **Redis** for caching.
- **AWS S3** as the file storage backend (you need valid AWS credentials).

### Docker Services

- **App**: Runs the Go application and exposes it on port `8080`.
- **PostgreSQL**: Database service exposed on port `5432`.
- **Redis**: Caching service exposed on port `6379`.


---
