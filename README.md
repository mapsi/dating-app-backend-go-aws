# Dating App Backend

This project is a backend service for a dating app, built with Go and using DynamoDB (emulated with LocalStack) for data storage.

## Features

- Create random user profiles
- Store user data in DynamoDB
- Structured logging with slog
- Docker and Docker Compose setup for easy deployment
- Makefile for common operations

## Setup

1. Ensure you have Go, Docker, and Docker Compose installed on your system.
2. Clone this repository.
3. Run `go mod download` to download dependencies.

## Running the Application

You can use the provided Makefile to run the application:

```bash
# Start the application and create a user
make

# Just start the Docker Compose setup
make up

# Create a user
make create-user

# View logs
make logs

# Stop the application
make down

# Clean up Docker resources
make clean
```

## API Endpoints

- **POST** `/user/create`: Creates a random user profile

### Login Endpoint

To use the login endpoint, send a `POST` request to `/login` with the following JSON body. A JWT token will be returned.

```json
{
  "email": "user@example.com",
  "password": "us3rP4ss0rd"
}
```

Example

```
curl -X POST -H "Content-Type: application/json" -d '{"email":"user@example.com","password":"us3rP4ss0rd"}' http://localhost:3000/login
```

### Discover Endpoint

To use the discover endpoint, send an authenticated GET request to `/discover`. The endpoint will return a list of potential matches, excluding the current user and users that have already been swiped on.

Example:

```
curl -X GET -H "Authorization: Bearer <your_jwt_token>" http://localhost:3000/discover
```

Response format:

```json
{
  "results": [
    {
      "id": "01F8Z6ARNVT4VQ3HTBD7BTHVF9",
      "name": "John Doe",
      "gender": "male",
      "age": 30
    },
    {
      "id": "01F8Z6ARNVT4VQ3HTBD7BTHVG9",
      "name": "Jane Smith",
      "gender": "female",
      "age": 28
    },
    ...
  ]
}
```

### Swipe Endpoint

To use the swipe endpoint, send an authenticated POST request to `/swipe` with the following JSON body:

```
{
   "swipedId": "user_id_of_the_swiped_profile",
   "preference": "YES" or "NO"
}
```

The server will respond with a result indicating whether there was a match:

```json
{
  "result": {
    "matched": true,
    "matchID": "01F8Z6ARNVT4VQ3HTBD7BTHVF9"
  }
}
```

Note:

- "NO" represents a dislike, while "YES" represents a like.
- The matchID filed is only include if `matched` is true.

Example:

```
curl -X POST http://localhost:3000/swipe \
     -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
     -H "Content-Type: application/json" \
     -d '{
           "swipedId": "01F8Z6ARNVT4VQ3HTBD7BTHVF9",
           "preference": "YES"
         }'
```

## Environment Variables

The application uses the following environment variables:

- PORT: The port on which the application runs (default: 3000)
- AWS_ENDPOINT: The AWS endpoint URL (default: http://localhost:4566 for LocalStack)
- AWS_REGION: The AWS region (default: eu-west-2)
- AWS_ACCESS_KEY_ID: AWS access key ID (default: dummy for LocalStack)
- AWS_SECRET_ACCESS_KEY: AWS secret access key (default: dummy for LocalStack)
- JWT_SECRET: A secret key used for signing and verifying JWT tokens

## Authentication

The application uses JWT (JSON Web Tokens) for authentication. Protected routes require a valid JWT token to be included in the Authorization header of the request.

### Authenticating Requests

To authenticate a request, include the JWT token in the Authorization header like this: `Authorization: Bearer <your_jwt_token>`

### Protected Routes

The following routes are protected and require authentication:

- **GET** `/discover`: Fetches profiles of potential matches

