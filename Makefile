# Variable
DC=docker compose
DC_FILE=docker-compose.yml
DB_USER=postgres
DB_NAME=user_db
DB_PASSWORD=postgres

.PHONY: help build down start stop restart logs status db-connect db-reset db-status db-backup db-restore

# Default target
.DEFAULT_GOAL := help

# Colors for output
RESET=\033[0m
BOLD=\033[1m
GREEN=\033[32m
YELLOW=\033[33m
RED=\033[31m
MAGENTA=\033[35m

# Docker compose commands
build: ## Build & start the containers
	@echo "$(GREEN)Starting build & run containers...$(RESET)"
	$(DC) -f $(DC_FILE) up -d
	@echo "\n$(GREEN)Container is running...$(RESET)"
	@$(DC) -f $(DC_FILE) ps
	@echo "\n$(GREEN)Access pgAdmin:$(RESET) http://localhost:5050"
	@echo "$(MAGENTA)Email:$(RESET) admin@admin.com"
	@echo "$(MAGENTA)Password:$(RESET) admin"
	@echo "\n$(MAGENTA)Username:$(RESET) $(DB_USER)"
	@echo "$(MAGENTA)Password:$(RESET) $(DB_PASSWORD)"
	@echo "$(MAGENTA)Database:$(RESET) $(DB_NAME)"

down: ## Stop & remove the containers (include volumes)
	@echo "$(YELLOW)Stopping & removing containers...$(RESET)"
	$(DC) -f $(DC_FILE) down -v

start: ## Start the containers
	@echo "$(GREEN)Starting containers...$(RESET)"
	$(DC) -f $(DC_FILE) start

stop: ## Stop the containers
	@echo "$(YELLOW)Stopping containers...$(RESET)"
	$(DC) -f $(DC_FILE) stop

restart: ## Restart the containers
	@echo "$(YELLOW)Restarting containers...$(RESET)"
	$(DC) -f $(DC_FILE) restart

status: ## Show the status of the containers
	@echo "$(GREEN)Container status...$(RESET)"
	$(DC) -f $(DC_FILE) ps

logs: ## Show the logs of the containers
	@echo "$(GREEN)Container logs...$(RESET)"
	$(DC) -f $(DC_FILE) logs -f

db-connect: ## Connect to the database
	@echo "$(GREEN)Connecting to the database...$(RESET)"
	@docker exec -it $(DB_NAME) psql -U $(DB_USER) -d $(DB_NAME)

db-reset: ## Reset the database
	@echo "$(YELLOW)Resetting the database...$(RESET)"
	@docker exec -it $(DB_NAME) psql -U $(DB_USER) -d $(DB_NAME) -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"

db-status: ## Show the status of the database
	@echo "$(GREEN)Database status...$(RESET)"
	@docker exec -it $(DB_NAME) psql -U $(DB_USER) -d $(DB_NAME) -c "\l"

db-backup: ## Backup the database
	@echo "$(GREEN)Backing up the database...$(RESET)"
	@docker exec -it $(DB_NAME) pg_dump -U $(DB_USER) -d $(DB_NAME) > backup.sql

db-restore: ## Restore the database
	@echo "$(GREEN)Restoring the database...$(RESET)"
	@docker exec -i $(DB_NAME) psql -U $(DB_USER) -d $(DB_NAME) < backup.sql


# Help command
help: ## Show this help
	@echo "$(BOLD)Available commands:$(RESET)"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[32m%-15s\033[0m %s\n", $$1, $$2}'