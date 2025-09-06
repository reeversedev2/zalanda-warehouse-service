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

# seed database with IKEA warehouse data
seed:
	@echo "🌱 Seeding IKEA warehouse database..."
	@./scripts/seed.sh

# run seed script directly (without shell script)
seed-direct:
	@echo "🌱 Running seed script directly..."
	@go run cmd/seed/main.go

# seed for development environment
seed-dev:
	@echo "🌱 Seeding development database..."
	@go run cmd/seed/main.go --env=development

# seed for production environment (with confirmation)
seed-prod:
	@echo "🌱 Seeding production database..."
	@go run cmd/seed/main.go --env=production --force

# dry run - show what would be seeded
seed-dry-run:
	@echo "🔍 Dry run - showing seed plan..."
	@go run cmd/seed/main.go --dry-run

# clear and seed database
seed-clear:
	@echo "🧹 Clearing and seeding database..."
	@go run cmd/seed/main.go --clear --env=development

# show seed script help
seed-help:
	@echo "📖 Seed script help..."
	@go run cmd/seed/main.go --help