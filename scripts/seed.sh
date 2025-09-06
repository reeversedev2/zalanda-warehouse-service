#!/bin/bash

echo "🌱 IKEA Warehouse Seed Script"
echo "=============================="

# Check if .env file exists
if [ ! -f .env ]; then
    echo "❌ .env file not found. Please create one with your database credentials."
    echo "Example .env file:"
    echo "DB_USER=your_db_user"
    echo "DB_PASSWORD=your_db_password"
    echo "DB_NAME=your_db_name"
    echo "JWT_SECRET=your-jwt-secret-key"
    exit 1
fi

# Load environment variables
export $(cat .env | grep -v '^#' | xargs)

# Check if required environment variables are set
if [ -z "$DB_USER" ] || [ -z "$DB_PASSWORD" ] || [ -z "$DB_NAME" ]; then
    echo "❌ Missing required environment variables: DB_USER, DB_PASSWORD, DB_NAME"
    exit 1
fi

echo "📊 Database Configuration:"
echo "   Host: db (Docker)"
echo "   User: $DB_USER"
echo "   Database: $DB_NAME"
echo ""

# Run the seed script
echo "🚀 Running seed script..."
go run cmd/seed/main.go

echo ""
echo "🎉 Seed script completed!"
echo ""
echo "📋 Test Credentials:"
echo "   Manager: lars.andersson@ikea.com / IKEA2024!"
echo "   Picker:  mikael.johansson@ikea.com / Picker123!"
echo "   Packer:  anders.pettersson@ikea.com / Packer123!"
echo ""
echo "🔗 API Endpoints to test:"
echo "   POST /api/auth/login"
echo "   GET  /api/protected/profile"
echo "   GET  /api/products"
