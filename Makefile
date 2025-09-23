all: build test

build:
	@echo "Building..."
	@go build -o main cmd/notification/main.go

run:
	@echo "Start container..."
	> logs/app.log
	@go run cmd/notification/main.go

docker-build:
	@echo "Start build service in docker..."
	docker build -t norification-service .

docker-run:
	@echo "Start container service..."
	docker run -p 8080:8080 --name norification-service -d norification-service

docker-compose-build:
	@echo "Start docker-compose build..."
	docker compose --build

docker-compose-run:
	@if docker compose up --build 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose up --build; \
	fi

docker-compose-down:
	@if docker compose down 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose down; \
	fi

test:
	@echo "Testing..."
	@go test ./... -v

itest:
	@echo "Running integration tests..."
	@go test ./internal/database -v

clean:
	@echo "Cleaning..."
	@rm -f main

clean-logs:
	@echo "Cleaning..."
	> ./logs/app.log

watch:
	@if command -v air > /dev/null; then \
            air; \
            echo "Watching...";\
        else \
            read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
            if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
                go install github.com/air-verse/air@latest; \
                air; \
                echo "Watching...";\
            else \
                echo "You chose not to install air. Exiting..."; \
                exit 1; \
            fi; \
        fi

.PHONY: all build run test clean watch docker-compose-run docker-compose-down itest clean-logs docker-build docker-run
