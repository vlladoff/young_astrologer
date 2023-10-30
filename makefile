API_IMAGE = young_astrologer_api
WORKER_IMAGE = young_astrologer_worker

.PHONY: all compose-up compose-down clean help

all: compose-up

compose-up:
	@echo "Starting the services using Docker Compose..."
	docker-compose up -d

compose-down:
	@echo "Stopping and removing Docker Compose services..."
	docker-compose down

clean:
	@echo "Cleaning up binary files..."
	rm -f $(API_IMAGE) $(WORKER_IMAGE)

help:
	@echo "Available targets:"
	@echo "  all        : Build and run services using Docker Compose"
	@echo "  compose-up : Build and run services using Docker Compose"
	@echo "  compose-down: Stop and remove Docker Compose services"
	@echo "  clean      : Clean up binary files"