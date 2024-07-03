# Makefile

# Variables
COMPOSE_FILE := docker-compose.yml
API_URL := http://localhost:3000

# Default target
all: up create-user

# Start Docker Compose in detached mode
up:
	docker-compose -f $(COMPOSE_FILE) up -d

# Stop Docker Compose
down:
	docker-compose -f $(COMPOSE_FILE) down

# Send a POST request to create a user
create-user:
	@echo "Creating a new user..."
	@curl -X POST $(API_URL)/user/create

# Restart the Docker Compose setup
restart: down up

# Show logs
logs:
	docker-compose -f $(COMPOSE_FILE) logs -f

# Clean up Docker resources
clean:
	docker-compose -f $(COMPOSE_FILE) down -v --remove-orphans

# Help target
help:
	@echo "Available targets:"
	@echo "  up           - Start Docker Compose in detached mode"
	@echo "  down         - Stop Docker Compose"
	@echo "  create-user  - Send a POST request to create a user"
	@echo "  restart      - Restart the Docker Compose setup"
	@echo "  logs         - Show logs"
	@echo "  clean        - Clean up Docker resources"
	@echo "  all          - Run 'up' and 'create-user' targets"

.PHONY: all up down create-user restart logs clean help