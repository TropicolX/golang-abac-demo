# Golang ABAC Demo

This repository demonstrates the implementation of Attribute-Based Access Control (ABAC) in a Golang-based document management system using Permify.

## Table of Contents

- [Introduction](#introduction)
- [Features](#features)
- [Setup](#setup)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
  - [Running the Application](#running-the-application)
- [API Endpoints](#api-endpoints)

## Introduction

The Golang ABAC Demo is an internal document management system that illustrates how ABAC can provide granular and dynamic access control based on various attributes such as user roles, department, document classification, and more. The backend is built with Golang, and Permify is used to manage ABAC policies and evaluations.

## Features

- User authentication and authorization
- Document upload, view, edit, and delete functionalities
- Granular access control using ABAC with Permify
- Middleware for access checks
- Logging of all requests

## Setup

### Prerequisites

- Go 1.16+
- Docker

### Installation

1. **Clone the repository**

   ```sh
   git clone https://github.com/TropicolX/golang-abac-demo.git
   cd golang-abac-demo
   ```

2. **Set up Permify**

   Pull and run the Permify Docker container:

   ```sh
   docker pull permify/permify
   docker run -d -p 3476:3476 --name permify permify/permify
   ```

3. **Install Go dependencies**

   ```sh
   go mod tidy
   ```

### Running the Application

1. **Start the server**

   ```sh
   go run cmd/server/main.go
   ```

2. **Verify the setup**

   Access the API endpoints using a tool like Postman or cURL.

## API Endpoints

- **POST /login**: User login
- **POST /api/documents**: Upload a document
- **GET /api/documents/{id}**: View a document
- **PUT /api/documents/{id}**: Edit a document
- **DELETE /api/documents/{id}**: Delete a document

### Example Requests

**Login**

```sh
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{
        "username": "user",
        "password": "password"
      }'
```

**Upload Document**

```sh
curl -X POST http://localhost:8080/api/documents \
  -H "Authorization: Bearer <your-token>" \
  -H "Content-Type: application/json" \
  -d '{
        "title": "Sample Document",
        "content": "This is a sample document.",
        "classification": "internal",
        "department": "IT"
      }'
```

**View Document**

```sh
curl -X GET http://localhost:8080/api/documents/<document-id> \
  -H "Authorization: Bearer <your-token>"
```

**Edit Document**

```sh
curl -X PUT http://localhost:8080/api/documents/<document-id> \
  -H "Authorization: Bearer <your-token>" \
  -H "Content-Type: application/json" \
  -d '{
        "title": "Updated Document Title",
        "content": "Updated content."
      }'
```

**Delete Document**

```sh
curl -X DELETE http://localhost:8080/api/documents/<document-id> \
  -H "Authorization: Bearer <your-token>"
```
