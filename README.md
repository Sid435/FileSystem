# File System Management Application

## Table of Contents
- [File System Management Application](#file-system-management-application)
  - [Table of Contents](#table-of-contents)
  - [Overview](#overview)
  - [Features](#features)
  - [Prerequisites](#prerequisites)
  - [Setup and Installation](#setup-and-installation)
  - [API Endpoints](#api-endpoints)
    - [Authentication](#authentication)
    - [File Operations](#file-operations)
  - [Testing](#testing)
    - [Postman Testing Guide](#postman-testing-guide)
    - [Unit Testing](#unit-testing)
  - [Deployment](#deployment)

## Overview

The File System Management Application is a robust, Go-based solution designed to efficiently manage file uploads and metadata. It leverages PostgreSQL for data persistence, AWS S3 for secure file storage, and Redis for performance optimization through caching. This application provides a scalable architecture suitable for handling high traffic and large volumes of data, while ensuring secure file access through pre-signed URLs and JWT-based authentication.

## Features

- **Secure File Upload**: Directly upload files to AWS S3 with proper access controls.
- **Metadata Management**: Efficiently store and retrieve file metadata using PostgreSQL.
- **Performance Optimization**: Utilize Redis caching to reduce database load and improve response times.
- **Secure File Access**: Generate pre-signed URLs for temporary, secure file downloads.
- **JWT Authentication**: Protect API endpoints with JSON Web Token based authentication.
- **Scalable Architecture**: Designed to handle high traffic and large volumes of data.
- **Docker Support**: Easy setup and deployment using Docker and Docker Compose.

## Prerequisites

Before setting up the application, ensure you have the following installed:

- [Go](https://golang.org/doc/install) (version 1.16 or later)
- [Docker](https://docs.docker.com/get-docker/)
- [Docker Compose](https://docs.docker.com/compose/install/)
- [AWS Account](https://aws.amazon.com/) with S3 and IAM access
- [PostgreSQL](https://www.postgresql.org/download/) (if running without Docker)
- [Redis](https://redis.io/download) (if running without Docker)

## Setup and Installation

1. **Clone the Repository**

   ```bash
   git clone https://github.com/your-username/file-system-management.git
   cd file-system-management
   ```

2. **Configure Environment Variables**

   Create a `.env` file in the root directory:

   ```ini
   AWS_ACCESS_KEY_ID=your_aws_access_key_id
   AWS_SECRET_ACCESS_KEY=your_aws_secret_access_key
   AWS_REGION=your_aws_region
   DB_HOST=db
   DB_PORT=5432
   DB_USER=your_db_username
   DB_PASSWORD=your_db_password
   DB_NAME=file_system
   REDIS_HOST=redis
   REDIS_PORT=6379
   JWT_SECRET=your_jwt_secret_key
   ```

   Replace placeholder values with your actual credentials and configuration.

3. **Build and Run with Docker Compose**

   ```bash
   docker-compose up --build
   ```

   This command builds the Docker images and starts all necessary services.

4. **Run Migrations (if applicable)**

   If your application uses database migrations:

   ```bash
   make migrate
   ```

5. **Start the Application**

   ```bash
   make run
   ```

   The application will start and listen on the configured ports.

## API Endpoints

### Authentication
- **Sign Up**: `POST /auth/signup`
  - Create a new user account.
  - Body: `username`, `password`

- **Login**: `POST /auth/login`
  - Authenticate and receive a JWT token.
  - Body: `username`, `password`

### File Operations
- **Upload File**: `POST /files/upload`
  - Upload one or more files to S3.
  - Headers: `Authorization: Bearer <token>`
  - Body: `multipart/form-data` with file(s)

- **Get File URL**: `GET /files/get`
  - Retrieve a pre-signed URL for file download.
  - Headers: `Authorization: Bearer <token>`
  - Query Parameters: `fileName`

## Testing

### Postman Testing Guide

We provide a Postman collection for comprehensive API testing. Follow these steps to test the application:

1. **Setup**
   - Import `File_System_Management.postman_collection.json` into Postman.
   - Create a Postman environment and add a variable named `jwt_token`.

2. **Test Sequence**
   a. **Authentication**
      - Run "Sign Up" to create a new user.
      - Run "Login" to obtain a JWT token (automatically sets `jwt_token`).

   b. **File Operations**
      - Use "Upload File" to upload a file (select file in form-data body).
      - Use "Get File URL" to retrieve a pre-signed URL (update `fileName` as needed).

   c. **Error Handling**
      - Run "Missing Authorization" (no JWT token provided).
      - Run "Invalid File Name" (non-existent file).

3. **Verify Responses**
   - Check for appropriate status codes (200 OK, 201 Created).
   - Verify meaningful error messages in error responses.
   - Confirm upload responses include file URLs.
   - Ensure file URL retrieval returns valid pre-signed URLs.

4. **Additional Tests**
   - Test caching by requesting the same file metadata multiple times.
   - Verify pre-signed URLs by accessing files through a web browser.

Note: Replace `/path/to/your/file.txt` in the "Upload File" request with an actual file path when testing.

### Unit Testing

Run the application's unit tests with:

```bash
make test
```

This executes all unit tests and provides a coverage report.

## Deployment

For production deployment, consider:

1. Using a reverse proxy (e.g., Nginx) for SSL termination and load balancing.
2. Setting up monitoring and logging (e.g., Prometheus, Grafana, ELK stack).
3. Implementing proper error handling and recovery mechanisms.
4. Regularly backing up your database and S3 bucket.
5. Implementing a CI/CD pipeline for automated testing and deployment.
