## Table of Contents

- [Overview](#overview)
- [Golang API Application Structure](#golang-api-application-structure)
- [Routes and Specifications](#routes-and-specifications)
- [E2E Tesing](#e2e-tesing)
- [Database Schema diagram](#database-schema-diagram)
- [Redis Revoke Refresh Token](#redis-revoke-refresh-token)


# Overview

This repository houses a Golang-based API application designed for managing organizations. The application includes features such as token management, CRUD operations for organizations, user invitations, and integration with MongoDB using Docker.

# Golang API Application Structure

The project structure is designed to assist you in getting started quickly. You can modify it as needed for your specific requirements.

- **cmd/**: Contains the main application file.
  - **main.go**: The entry point of the application.

- **pkg/**: Core logic of the application divided into different packages.
  - **api/**: API handling components.
    - **handlers/**: API route handlers.
    - **middleware/**: Middleware functions.
    - **routes/**: Route definitions.
  - **controllers/**: Business logic for each route.
  - **database/**: Database-related code.
    - **mongodb/**
      - **models/**: Data models.
      - **repository/**: Database operations.
  - **utils/**: Utility functions.
  - **app.go**: Application initialization and setup.

- **docker/**: Docker-related files.
  - **Dockerfile**: Instructions for building the application image.

- **docker-compose.yaml**: Configuration for Docker Compose.

- **config/**: Configuration files for the application.
  - **app-config.yaml**: General application settings.
  - **database-config.yaml**: Database connection details.

- **tests/**: Directory for tests.
  - **e2e/**: End-to-End tests.
  - **unit/**: Unit tests.

- **.gitignore**: Specifies files and directories to be ignored by Git.

## Getting Started

To begin working with the application, follow the instructions in the project documentation. Feel free to adjust the project structure as needed based on your preferences and evolving project requirements.

# Routes and Specifications

### Signup Endpoint:
```
Request Shema: POST /signup
Request Body:
JSON {
    "name": "string",
    "email": "string",
    "password": "string",
}
Response Schema:
JSON {
    "message": "string"
}
```

### Signin Endpoint:
```
Request Shema: POST /signin
Request Body:
JSON {
    "email": "string",
    "password": "string",
}
Response Schema:
JSON {
    "message": "string",
    "access_token": "string",
    "refresh_token": "string",
}
```
### Refresh Token Endpoint:
```
Request Shema: POST /refresh-token
Request Body:
JSON {
    "refresh_token": "string",
}
Response Schema:
JSON {
    "message": "string",
    "access_token": "string",
    "refresh_token": "string",
}
```

### Create Organization Endpoint:
```
Request Shema: POST /organization
Authorization: Bearer [Token]
Request Body:
JSON {
    "name": "string",
    "description": "string",
}
Response Schema:
JSON {
    "organization_id": "string",
}
```

### Read Organization Endpoint:
```
Request Shema: GET /organization/{organization_id}
Authorization: Bearer [Token]
Response Schema:
JSON {
    "organization_id": "string",
    "name": "string",
    "description": "string",
    "organization_members": [
        {
            "name": "string",
            "email": "string",
            "access_level": "string",
        },
        ...
    ],
}
```

### Read All Organizations Endpoint:
```
Request Shema: GET /organization
Authorization: Bearer [Token]
Response Schema:
JSON [
    {
        "organization_id": "string",
        "name": "string",
        "description": "string",
        "organization_members": [
            {
                "name": "string",
                "email": "string",
                "access_level": "string",
            },
            ...
        ],
    },
    ...
]
```

### Update Organization Endpoint:
```
Request Shema: PUT /organization/{organization_id}
Authorization: Bearer [Token]
Request Body:
JSON {
    "name": "string",
    "description": "string",
}
Response Schema:
JSON {
    "organization_id": "string",
    "name": "string",
    "description": "string",
}
```

### Delete Organization Endpoint:
```
Request Shema: DELETE /organization/{organization_id}
Authorization: Bearer [Token]
Response Schema:
JSON {
    "message": "string",
}
```

### Invite User to Organization Endpoint:
```
Request Shema: POST /organization/{organization_id}/invite
Authorization: Bearer [Token]
Request Body:
JSON {
    "user_email": "string",
}
Response Schema:
JSON {
    "message": "string",
}
```

### Refresh Token Revocation Integrated With Redis:
``` 
Request Shema: POST /revoke-refresh-token/
Authorization: Bearer [Token]
Request Body:
JSON {
"refresh_token": "string",
}
Response Schema:
JSON {
"message": "string",
}
```

# E2E Tesing

E2E testing is crucial for ensuring the reliability of user authentication. We're currently focusing on testing user signup and signin processes. Additional tests will be added to cover more scenarios and functionalities, ensuring the overall robustness of our system.

To Run Test:
```
>> cd test/e2e
>> go test
```
similar output should be found:
```
Pass
ok organization_management/tests/e2e    1.547s
```

# Database Schema diagram

```
Users Collection:
{
    "_id": ObjectId,
    "username": string,
    "email": string,
    "password": string,
}

Organizations Collection:
{
    "_id": ObjectId,
    "name": string,
    "description": string,
    "organization_members": [
        {
            "name": string,
            "email": string,
            "access_level": string
        },
        ...
    ],
}
```

# Redis Revoke Refresh Token

Redis was integrated into the project to handle refresh token revocation, improving security and access management. This involved setting up a Redis client, creating a token repository, adding a revoke refresh token endpoint, integrating the controller with the token repository, and implementing response handling. This integration 