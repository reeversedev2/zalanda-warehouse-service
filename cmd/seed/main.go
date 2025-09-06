package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/reeversedev2/zalanda-warehouse-service/pkg/handlers"
	"github.com/reeversedev2/zalanda-warehouse-service/pkg/models"
	"github.com/reeversedev2/zalanda-warehouse-service/pkg/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	// Command line flags
	var (
		force     = flag.Bool("force", false, "Force seeding even if data exists")
		clear     = flag.Bool("clear", false, "Clear existing data before seeding")
		env       = flag.String("env", "development", "Environment (development, staging, production)")
		dryRun    = flag.Bool("dry-run", false, "Show what would be seeded without actually doing it")
		help      = flag.Bool("help", false, "Show help")
	)
	flag.Parse()

	if *help {
		showHelp()
		return
	}

	fmt.Printf("🌱 IKEA Warehouse Seed Script (Environment: %s)\n", strings.ToUpper(*env))
	fmt.Println(strings.Repeat("=", 50))

	if *dryRun {
		fmt.Println("🔍 DRY RUN MODE - No changes will be made")
		showSeedPlan()
		return
	}

	// Environment validation
	if *env == "production" && !*force {
		fmt.Println("⚠️  WARNING: You are about to seed PRODUCTION database!")
		fmt.Print("Type 'yes' to continue: ")
		var confirmation string
		fmt.Scanln(&confirmation)
		if strings.ToLower(confirmation) != "yes" {
			fmt.Println("❌ Seeding cancelled.")
			return
		}
	}

	// Connect to database
	db := connectToDatabase(*env)
	if db == nil {
		log.Fatal("❌ Failed to connect to database")
	}

	// Check if data already exists
	if !*force && !*clear {
		if hasExistingData(db) {
			fmt.Println("⚠️  Database already contains data. Use --force to overwrite or --clear to clear first.")
			return
		}
	}

	startTime := time.Now()

	// Clear existing data if requested
	if *clear {
		fmt.Println("🧹 Clearing existing data...")
		clearDatabase(db)
	}

	// Seed data
	fmt.Println("🏢 Seeding companies...")
	companies := seedCompanies(db)

	fmt.Println("👥 Seeding warehouse employees...")
	users := seedUsers(db)

	fmt.Println("📦 Seeding IKEA products...")
	products := seedProducts(db, companies)

	duration := time.Since(startTime)

	fmt.Println("✅ Seed script completed successfully!")
	fmt.Printf("📊 Summary:\n")
	fmt.Printf("   - Companies: %d\n", len(companies))
	fmt.Printf("   - Users: %d\n", len(users))
	fmt.Printf("   - Products: %d\n", len(products))
	fmt.Printf("   - Duration: %v\n", duration.Round(time.Millisecond))
}

func connectToDatabase(env string) *gorm.DB {
	// Determine database host based on environment
	host := getEnv("DB_HOST", "localhost")
	if env == "production" {
		host = getEnv("DB_HOST", "db") // Use Docker service name in production
	}

	// SSL mode based on environment
	sslMode := "disable"
	if env == "production" {
		sslMode = getEnv("DB_SSLMODE", "require")
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=UTC",
		host,
		getEnv("DB_USER", "warehouse_user"),
		getEnv("DB_PASSWORD", "secure_warehouse_password"),
		getEnv("DB_NAME", "ikea_warehouse_db"),
		getEnv("DB_PORT", "5432"),
		sslMode,
	)

	fmt.Printf("🔗 Connecting to database: %s@%s:%s/%s\n", 
		getEnv("DB_USER", "warehouse_user"), 
		host, 
		getEnv("DB_PORT", "5432"), 
		getEnv("DB_NAME", "ikea_warehouse_db"))

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // Silent in production
	})

	if err != nil {
		fmt.Printf("❌ Failed to connect to database: %v\n", err)
		return nil
	}

	fmt.Println("✅ Database connected successfully")

	// Run migrations
	fmt.Println("🔄 Running database migrations...")
	if err := db.AutoMigrate(&models.Product{}, &models.Company{}, &models.User{}); err != nil {
		fmt.Printf("❌ Migration failed: %v\n", err)
		return nil
	}

	return db
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func showHelp() {
	fmt.Println("🌱 IKEA Warehouse Seed Script")
	fmt.Println("Usage: go run cmd/seed/main.go [options]")
	fmt.Println("\nOptions:")
	fmt.Println("  --force       Force seeding even if data exists")
	fmt.Println("  --clear       Clear existing data before seeding")
	fmt.Println("  --env         Environment (development, staging, production)")
	fmt.Println("  --dry-run     Show what would be seeded without actually doing it")
	fmt.Println("  --help        Show this help message")
	fmt.Println("\nExamples:")
	fmt.Println("  go run cmd/seed/main.go --env=development")
	fmt.Println("  go run cmd/seed/main.go --env=production --force")
	fmt.Println("  go run cmd/seed/main.go --dry-run")
	fmt.Println("  go run cmd/seed/main.go --clear --env=staging")
}

func hasExistingData(db *gorm.DB) bool {
	var count int64
	db.Model(&models.User{}).Count(&count)
	return count > 0
}

func clearDatabase(db *gorm.DB) {
	// Disable foreign key checks temporarily
	db.Exec("SET session_replication_role = replica;")
	
	// Clear data in correct order (respecting foreign keys)
	db.Exec("DELETE FROM products")
	db.Exec("DELETE FROM users")
	db.Exec("DELETE FROM companies")
	
	// Re-enable foreign key checks
	db.Exec("SET session_replication_role = DEFAULT;")
	
	fmt.Println("✅ Database cleared successfully")
}

func showSeedPlan() {
	fmt.Println("📋 Seed Plan:")
	fmt.Println("   - 3 Companies (IKEA Sweden AB, IKEA Supply Chain, IKEA Logistics)")
	fmt.Println("   - 12 Users (3 Managers, 4 Pickers, 5 Packers)")
	fmt.Println("   - 33+ IKEA Products (Furniture, Kitchen, Storage, etc.)")
	fmt.Println("\n🔐 Default Credentials:")
	fmt.Println("   Manager: lars.andersson@ikea.com / IKEA2024!")
	fmt.Println("   Picker:  mikael.johansson@ikea.com / Picker123!")
	fmt.Println("   Packer:  anders.pettersson@ikea.com / Packer123!")
}

func seedCompanies(db *gorm.DB) []models.Company {
	companies := []models.Company{
		{
			Name:  "IKEA Sweden AB",
			Image: "https://images.unsplash.com/photo-1586023492125-27b2c045efd7?w=150&h=150&fit=crop",
		},
		{
			Name:  "IKEA Supply Chain",
			Image: "https://images.unsplash.com/photo-1566576912321-d58ddd7a6088?w=150&h=150&fit=crop",
		},
		{
			Name:  "IKEA Logistics",
			Image: "https://images.unsplash.com/photo-1586528116311-ad8dd3c8310d?w=150&h=150&fit=crop",
		},
	}

	for i := range companies {
		if err := db.Create(&companies[i]).Error; err != nil {
			log.Printf("Error creating company %s: %v", companies[i].Name, err)
		}
	}

	return companies
}

func seedUsers(db *gorm.DB) []models.User {
	// Hash passwords
	managerPassword, _ := handlers.HashPassword("IKEA2024!")
	pickerPassword, _ := handlers.HashPassword("Picker123!")
	packerPassword, _ := handlers.HashPassword("Packer123!")

	users := []models.User{
		// Managers
		{
			Email:      "lars.andersson@ikea.com",
			Password:   managerPassword,
			Role:       utils.Manager.GetRole(),
			Name:       "Lars Andersson",
			EmployeeID: "IKEA-MGR-001",
			IsActive:   true,
		},
		{
			Email:      "anna.lindberg@ikea.com",
			Password:   managerPassword,
			Role:       utils.Manager.GetRole(),
			Name:       "Anna Lindberg",
			EmployeeID: "IKEA-MGR-002",
			IsActive:   true,
		},
		{
			Email:      "erik.svensson@ikea.com",
			Password:   managerPassword,
			Role:       utils.Manager.GetRole(),
			Name:       "Erik Svensson",
			EmployeeID: "IKEA-MGR-003",
			IsActive:   true,
		},

		// Pickers
		{
			Email:      "mikael.johansson@ikea.com",
			Password:   pickerPassword,
			Role:       utils.Picker.GetRole(),
			Name:       "Mikael Johansson",
			EmployeeID: "IKEA-PCK-001",
			IsActive:   true,
		},
		{
			Email:      "sofia.gustafsson@ikea.com",
			Password:   pickerPassword,
			Role:       utils.Picker.GetRole(),
			Name:       "Sofia Gustafsson",
			EmployeeID: "IKEA-PCK-002",
			IsActive:   true,
		},
		{
			Email:      "olof.nilsson@ikea.com",
			Password:   pickerPassword,
			Role:       utils.Picker.GetRole(),
			Name:       "Olof Nilsson",
			EmployeeID: "IKEA-PCK-003",
			IsActive:   true,
		},
		{
			Email:      "emma.eriksson@ikea.com",
			Password:   pickerPassword,
			Role:       utils.Picker.GetRole(),
			Name:       "Emma Eriksson",
			EmployeeID: "IKEA-PCK-004",
			IsActive:   true,
		},

		// Packers
		{
			Email:      "anders.pettersson@ikea.com",
			Password:   packerPassword,
			Role:       utils.Packer.GetRole(),
			Name:       "Anders Pettersson",
			EmployeeID: "IKEA-PKR-001",
			IsActive:   true,
		},
		{
			Email:      "lisa.hansson@ikea.com",
			Password:   packerPassword,
			Role:       utils.Packer.GetRole(),
			Name:       "Lisa Hansson",
			EmployeeID: "IKEA-PKR-002",
			IsActive:   true,
		},
		{
			Email:      "marcus.olsson@ikea.com",
			Password:   packerPassword,
			Role:       utils.Packer.GetRole(),
			Name:       "Marcus Olsson",
			EmployeeID: "IKEA-PKR-003",
			IsActive:   true,
		},
		{
			Email:      "jenny.karlsson@ikea.com",
			Password:   packerPassword,
			Role:       utils.Packer.GetRole(),
			Name:       "Jenny Karlsson",
			EmployeeID: "IKEA-PKR-004",
			IsActive:   true,
		},
		{
			Email:      "daniel.larsson@ikea.com",
			Password:   packerPassword,
			Role:       utils.Packer.GetRole(),
			Name:       "Daniel Larsson",
			EmployeeID: "IKEA-PKR-005",
			IsActive:   true,
		},
	}

	for i := range users {
		if err := db.Create(&users[i]).Error; err != nil {
			log.Printf("Error creating user %s: %v", users[i].Name, err)
		}
	}

	return users
}

func seedProducts(db *gorm.DB, companies []models.Company) []models.Product {
	// IKEA product categories and items
	products := []models.Product{
		// Furniture - Living Room
		{
			Name:      "BILLY Bookcase",
			CompanyID: int(companies[0].ID),
			Price:     79.99,
			Category:  "Furniture - Living Room",
			Expire:    "2025-12-31",
			Status:    "in_stock",
			Image:     "https://images.unsplash.com/photo-1586023492125-27b2c045efd7?w=300&h=300&fit=crop",
		},
		{
			Name:      "EKTORP Sofa",
			CompanyID: int(companies[0].ID),
			Price:     399.99,
			Category:  "Furniture - Living Room",
			Expire:    "2025-12-31",
			Status:    "in_stock",
			Image:     "https://images.unsplash.com/photo-1586023492125-27b2c045efd7?w=300&h=300&fit=crop",
		},
		{
			Name:      "LACK Coffee Table",
			CompanyID: int(companies[0].ID),
			Price:     29.99,
			Category:  "Furniture - Living Room",
			Expire:    "2025-12-31",
			Status:    "in_stock",
			Image:     "https://images.unsplash.com/photo-1586023492125-27b2c045efd7?w=300&h=300&fit=crop",
		},
		{
			Name:      "POÄNG Armchair",
			CompanyID: int(companies[0].ID),
			Price:     149.99,
			Category:  "Furniture - Living Room",
			Expire:    "2025-12-31",
			Status:    "in_stock",
			Image:     "https://images.unsplash.com/photo-1586023492125-27b2c045efd7?w=300&h=300&fit=crop",
		},

		// Furniture - Bedroom
		{
			Name:      "MALM Bed Frame",
			CompanyID: int(companies[0].ID),
			Price:     199.99,
			Category:  "Furniture - Bedroom",
			Expire:    "2025-12-31",
			Status:    "in_stock",
			Image:     "https://images.unsplash.com/photo-1586023492125-27b2c045efd7?w=300&h=300&fit=crop",
		},
		{
			Name:      "HEMNES Dresser",
			CompanyID: int(companies[0].ID),
			Price:     249.99,
			Category:  "Furniture - Bedroom",
			Expire:    "2025-12-31",
			Status:    "in_stock",
			Image:     "https://images.unsplash.com/photo-1586023492125-27b2c045efd7?w=300&h=300&fit=crop",
		},
		{
			Name:      "NORDLI Bed Frame",
			CompanyID: int(companies[0].ID),
			Price:     299.99,
			Category:  "Furniture - Bedroom",
			Expire:    "2025-12-31",
			Status:    "in_stock",
			Image:     "https://images.unsplash.com/photo-1586023492125-27b2c045efd7?w=300&h=300&fit=crop",
		},

		// Kitchen & Dining
		{
			Name:      "KALLAX Shelf Unit",
			CompanyID: int(companies[0].ID),
			Price:     89.99,
			Category:  "Kitchen & Dining",
			Expire:    "2025-12-31",
			Status:    "in_stock",
			Image:     "https://images.unsplash.com/photo-1586023492125-27b2c045efd7?w=300&h=300&fit=crop",
		},
		{
			Name:      "FÖRHÖJA Kitchen Cart",
			CompanyID: int(companies[0].ID),
			Price:     179.99,
			Category:  "Kitchen & Dining",
			Expire:    "2025-12-31",
			Status:    "in_stock",
			Image:     "https://images.unsplash.com/photo-1586023492125-27b2c045efd7?w=300&h=300&fit=crop",
		},
		{
			Name:      "INGO Dining Table",
			CompanyID: int(companies[0].ID),
			Price:     129.99,
			Category:  "Kitchen & Dining",
			Expire:    "2025-12-31",
			Status:    "in_stock",
			Image:     "https://images.unsplash.com/photo-1586023492125-27b2c045efd7?w=300&h=300&fit=crop",
		},

		// Storage & Organization
		{
			Name:      "KALLAX Storage Unit",
			CompanyID: int(companies[0].ID),
			Price:     69.99,
			Category:  "Storage & Organization",
			Expire:    "2025-12-31",
			Status:    "in_stock",
			Image:     "https://images.unsplash.com/photo-1586023492125-27b2c045efd7?w=300&h=300&fit=crop",
		},
		{
			Name:      "TROFAST Storage Box",
			CompanyID: int(companies[0].ID),
			Price:     12.99,
			Category:  "Storage & Organization",
			Expire:    "2025-12-31",
			Status:    "in_stock",
			Image:     "https://images.unsplash.com/photo-1586023492125-27b2c045efd7?w=300&h=300&fit=crop",
		},
		{
			Name:      "SKUBB Storage Box",
			CompanyID: int(companies[0].ID),
			Price:     7.99,
			Category:  "Storage & Organization",
			Expire:    "2025-12-31",
			Status:    "in_stock",
			Image:     "https://images.unsplash.com/photo-1586023492125-27b2c045efd7?w=300&h=300&fit=crop",
		},

		// Lighting
		{
			Name:      "FADO Table Lamp",
			CompanyID: int(companies[0].ID),
			Price:     24.99,
			Category:  "Lighting",
			Expire:    "2025-12-31",
			Status:    "in_stock",
			Image:     "https://images.unsplash.com/photo-1586023492125-27b2c045efd7?w=300&h=300&fit=crop",
		},
		{
			Name:      "HEKTAR Floor Lamp",
			CompanyID: int(companies[0].ID),
			Price:     39.99,
			Category:  "Lighting",
			Expire:    "2025-12-31",
			Status:    "in_stock",
			Image:     "https://images.unsplash.com/photo-1586023492125-27b2c045efd7?w=300&h=300&fit=crop",
		},
		{
			Name:      "RANARP Work Lamp",
			CompanyID: int(companies[0].ID),
			Price:     49.99,
			Category:  "Lighting",
			Expire:    "2025-12-31",
			Status:    "in_stock",
			Image:     "https://images.unsplash.com/photo-1586023492125-27b2c045efd7?w=300&h=300&fit=crop",
		},

		// Textiles & Rugs
		{
			Name:      "VINDUM Rug",
			CompanyID: int(companies[0].ID),
			Price:     79.99,
			Category:  "Textiles & Rugs",
			Expire:    "2025-12-31",
			Status:    "in_stock",
			Image:     "https://images.unsplash.com/photo-1586023492125-27b2c045efd7?w=300&h=300&fit=crop",
		},
		{
			Name:      "LUDDE Cushion Cover",
			CompanyID: int(companies[0].ID),
			Price:     9.99,
			Category:  "Textiles & Rugs",
			Expire:    "2025-12-31",
			Status:    "in_stock",
			Image:     "https://images.unsplash.com/photo-1586023492125-27b2c045efd7?w=300&h=300&fit=crop",
		},
		{
			Name:      "STRANDMON Wing Chair",
			CompanyID: int(companies[0].ID),
			Price:     199.99,
			Category:  "Textiles & Rugs",
			Expire:    "2025-12-31",
			Status:    "in_stock",
			Image:     "https://images.unsplash.com/photo-1586023492125-27b2c045efd7?w=300&h=300&fit=crop",
		},

		// Bathroom
		{
			Name:      "GODMORGON Bathroom Cabinet",
			CompanyID: int(companies[0].ID),
			Price:     149.99,
			Category:  "Bathroom",
			Expire:    "2025-12-31",
			Status:    "in_stock",
			Image:     "https://images.unsplash.com/photo-1586023492125-27b2c045efd7?w=300&h=300&fit=crop",
		},
		{
			Name:      "RINNIG Shower Curtain",
			CompanyID: int(companies[0].ID),
			Price:     12.99,
			Category:  "Bathroom",
			Expire:    "2025-12-31",
			Status:    "in_stock",
			Image:     "https://images.unsplash.com/photo-1586023492125-27b2c045efd7?w=300&h=300&fit=crop",
		},

		// Office & Workspace
		{
			Name:      "BEKANT Desk",
			CompanyID: int(companies[0].ID),
			Price:     179.99,
			Category:  "Office & Workspace",
			Expire:    "2025-12-31",
			Status:    "in_stock",
			Image:     "https://images.unsplash.com/photo-1586023492125-27b2c045efd7?w=300&h=300&fit=crop",
		},
		{
			Name:      "MARKUS Office Chair",
			CompanyID: int(companies[0].ID),
			Price:     199.99,
			Category:  "Office & Workspace",
			Expire:    "2025-12-31",
			Status:    "in_stock",
			Image:     "https://images.unsplash.com/photo-1586023492125-27b2c045efd7?w=300&h=300&fit=crop",
		},
		{
			Name:      "ALEX Drawer Unit",
			CompanyID: int(companies[0].ID),
			Price:     89.99,
			Category:  "Office & Workspace",
			Expire:    "2025-12-31",
			Status:    "in_stock",
			Image:     "https://images.unsplash.com/photo-1586023492125-27b2c045efd7?w=300&h=300&fit=crop",
		},

		// Outdoor
		{
			Name:      "APPLARÖ Outdoor Table",
			CompanyID: int(companies[0].ID),
			Price:     199.99,
			Category:  "Outdoor",
			Expire:    "2025-12-31",
			Status:    "in_stock",
			Image:     "https://images.unsplash.com/photo-1586023492125-27b2c045efd7?w=300&h=300&fit=crop",
		},
		{
			Name:      "FÄRLÖV Outdoor Chair",
			CompanyID: int(companies[0].ID),
			Price:     79.99,
			Category:  "Outdoor",
			Expire:    "2025-12-31",
			Status:    "in_stock",
			Image:     "https://images.unsplash.com/photo-1586023492125-27b2c045efd7?w=300&h=300&fit=crop",
		},

		// Kids & Baby
		{
			Name:      "STUVA Storage Combination",
			CompanyID: int(companies[0].ID),
			Price:     149.99,
			Category:  "Kids & Baby",
			Expire:    "2025-12-31",
			Status:    "in_stock",
			Image:     "https://images.unsplash.com/photo-1586023492125-27b2c045efd7?w=300&h=300&fit=crop",
		},
		{
			Name:      "SMÅSTAD Wardrobe",
			CompanyID: int(companies[0].ID),
			Price:     199.99,
			Category:  "Kids & Baby",
			Expire:    "2025-12-31",
			Status:    "in_stock",
			Image:     "https://images.unsplash.com/photo-1586023492125-27b2c045efd7?w=300&h=300&fit=crop",
		},

		// Some products with different statuses for testing
		{
			Name:      "MALM Bed Frame - White",
			CompanyID: int(companies[0].ID),
			Price:     199.99,
			Category:  "Furniture - Bedroom",
			Expire:    "2025-12-31",
			Status:    "low_stock",
			Image:     "https://images.unsplash.com/photo-1586023492125-27b2c045efd7?w=300&h=300&fit=crop",
		},
		{
			Name:      "BILLY Bookcase - Oak",
			CompanyID: int(companies[0].ID),
			Price:     89.99,
			Category:  "Furniture - Living Room",
			Expire:    "2025-12-31",
			Status:    "out_of_stock",
			Image:     "https://images.unsplash.com/photo-1586023492125-27b2c045efd7?w=300&h=300&fit=crop",
		},
		{
			Name:      "POÄNG Armchair - Black",
			CompanyID: int(companies[0].ID),
			Price:     149.99,
			Category:  "Furniture - Living Room",
			Expire:    "2025-12-31",
			Status:    "discontinued",
			Image:     "https://images.unsplash.com/photo-1586023492125-27b2c045efd7?w=300&h=300&fit=crop",
		},
		{
			Name:      "EKTORP Sofa - Gray",
			CompanyID: int(companies[0].ID),
			Price:     399.99,
			Category:  "Furniture - Living Room",
			Expire:    "2025-12-31",
			Status:    "in_stock",
			Image:     "https://images.unsplash.com/photo-1586023492125-27b2c045efd7?w=300&h=300&fit=crop",
		},
		{
			Name:      "LACK Coffee Table - Black",
			CompanyID: int(companies[0].ID),
			Price:     29.99,
			Category:  "Furniture - Living Room",
			Expire:    "2025-12-31",
			Status:    "in_stock",
			Image:     "https://images.unsplash.com/photo-1586023492125-27b2c045efd7?w=300&h=300&fit=crop",
		},
	}

	// Create products in batches
	for i := range products {
		if err := db.Create(&products[i]).Error; err != nil {
			log.Printf("Error creating product %s: %v", products[i].Name, err)
		}
	}

	fmt.Printf("✅ Created %d IKEA products\n", len(products))
	return products
}
