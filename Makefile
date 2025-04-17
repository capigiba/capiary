# Variables
CMD_DIR := cmd/server
SETUP_DIR := cmd/setup
APP_NAME := capiary
SCRIPT_DIR := scripts
INITIALIZE_DIR := internal/initialize
CMD_MIGRATION := cmd/migration

# Default target
all: build

# Run the application
run:
	go run $(CMD_DIR)/main.go

setup:
	go run $(SETUP_DIR)/main.go

# Generate wire dependencies
wire:
	go install github.com/google/wire/cmd/wire@latest
	cd $(INITIALIZE_DIR) && wire

# Build the application
build:
	cd $(CMD_DIR) && go build -o $(APP_NAME)

# Clean the generated binaries
clean:
	rm -f $(CMD_DIR)/$(APP_NAME)

# Swagger
swag:
	go install github.com/swaggo/swag/cmd/swag@latest
	swag init -g $(CMD_DIR)/main.go -o ./docs

# Run all steps from the script
run-all:
	bash $(SCRIPT_DIR)/run-all.sh

migration:
	go run $(CMD_MIGRATION)/main.go
	

# Help
help:
	@echo "Makefile for $(APP_NAME)"
	@echo
	@echo "Usage:"
	@echo "  make run         Run the application"
	@echo "  make swag		  Run the swagger"
	@echo "  make build       Build the application"
	@echo "  make clean       Clean the generated binaries"
	@echo "  make help        Show this help message"