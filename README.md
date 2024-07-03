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

### API Endpoints

- **POST** /user/create: Creates a random user profile

### Environment Variables

The application uses the following environment variables:

- PORT: The port on which the application runs (default: 3000)
- AWS_ENDPOINT: The AWS endpoint URL (default: http://localhost:4566 for LocalStack)
- AWS_REGION: The AWS region (default: eu-west-2)
- AWS_ACCESS_KEY_ID: AWS access key ID (default: dummy for LocalStack)
- AWS_SECRET_ACCESS_KEY: AWS secret access key (default: dummy for LocalStack)
