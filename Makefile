#build
build:
	@echo "Building server..."
	@docker compose build
# start server
start:
	@echo "Starting server..."
	@docker compose up --remove-orphans
# stop server
stop:
	@echo "Stopping server..."
	@docker compose down