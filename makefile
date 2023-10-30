API_IMAGE = young_astrologer_api
WORKER_IMAGE = young_astrologer_worker

.PHONY: all api worker clean compose-up compose-down

all: api worker

# Build the API application
api:
	@echo "Building the API application..."
	cd cmd/api && go build -o ../$(API_IMAGE)

# Build the worker application
worker:
	@echo "Building the worker application..."
	cd cmd/worker && go build -o ../$(WORKER_IMAGE)

# Build and run both API and worker using Docker Compose
compose-up: api worker
	@echo "Starting the services using Docker Compose..."
	docker-compose up -d

# Stop and remove Docker Compose services
compose-down:
	@echo "Stopping and removing Docker Compose services..."
	docker-compose down

# Clean up binary files
clean:
	@echo "Cleaning up binary files..."
	rm -f $(API_IMAGE) $(WORKER_IMAGE)

# Help target to display available targets
help:
	@echo "Available targets:"
	@echo "  all        : Build both API and worker"
	@echo "  api        : Build the API application"
	@echo "  worker     : Build the worker application"
	@echo "  compose-up : Build and run services using Docker Compose"
	@echo "  compose-down: Stop and remove Docker Compose services"
	@echo "  clean      : Clean up binary files"