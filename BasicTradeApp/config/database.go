package config

import (
	"fmt"
	"net"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	// Retrieve environment variables
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := "viaduct.proxy.rlwy.net"
	dbPort := "45408"
	dbName := os.Getenv("DB_NAME")

	// Print environment variables to verify
	fmt.Println("DB_USER:", dbUser)
	fmt.Println("DB_PASS:", dbPass)
	fmt.Println("DB_HOST:", dbHost)
	fmt.Println("DB_PORT:", dbPort)
	fmt.Println("DB_NAME:", dbName)

	// Check DNS resolution
	ips, err := net.LookupIP(dbHost)
	if err != nil {
		panic("DNS lookup failed: " + err.Error())
	}
	fmt.Println("Resolved IP addresses for", dbHost, ":", ips)

	// Create DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPass, dbHost, dbPort, dbName)
	fmt.Println("DSN:", dsn) // Print DSN for verification

	// Open database connection
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database: " + err.Error())
	}
	fmt.Println("Connected to database successfully")

	// Verify current database
	var dbNameCheck string
	db.Raw("SELECT DATABASE()").Scan(&dbNameCheck)
	fmt.Println("Currently connected to database:", dbNameCheck)

	// Check if admins table exists
	var tableExists bool
	db.Raw("SELECT COUNT(*) > 0 FROM information_schema.tables WHERE table_schema = ? AND table_name = 'admins'", dbNameCheck).Scan(&tableExists)
	if !tableExists {
		panic(fmt.Sprintf("Table 'admins' does not exist in database '%s'", dbNameCheck))
	}
	fmt.Println("Table 'admins' exists in database:", dbNameCheck)

	// Test query to verify connection
	var result int
	err = db.Raw("SELECT 1").Scan(&result).Error
	if err != nil {
		panic("Test query failed: " + err.Error())
	}
	fmt.Println("Test query successful, result:", result)

	// Assign db to global variable
	DB = db
}

func main() {
	ConnectDB()
}
