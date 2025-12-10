package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	ID       int64  `gorm:"primaryKey"`
	Username string
	Email    string
	Password string
	IsActive bool
}

func main() {
	// Load .env
	godotenv.Load("../.env")

	// Build DSN
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	// Connect
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	fmt.Println("✅ Connected to database successfully")

	// Generate new hash for "password123"
	newPassword := "password123"
	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), 10)
	if err != nil {
		log.Fatal("Failed to generate hash:", err)
	}

	fmt.Printf("\n🔐 New hash for password '%s':\n%s\n\n", newPassword, string(hash))

	// Update all users
	result := db.Model(&User{}).Where("1=1").Update("password", string(hash))
	if result.Error != nil {
		log.Fatal("Failed to update passwords:", result.Error)
	}

	fmt.Printf("✅ Successfully updated %d user passwords\n", result.RowsAffected)

	// Verify admin user
	var admin User
	db.Where("username = ?", "admin").First(&admin)

	// Test password
	err = bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(newPassword))
	if err == nil {
		fmt.Println("\n✅ Password verification successful for admin user!")
		fmt.Println("   You can now login with:")
		fmt.Printf("   Username: %s\n", admin.Username)
		fmt.Printf("   Email: %s\n", admin.Email)
		fmt.Printf("   Password: %s\n", newPassword)
	} else {
		fmt.Println("\n❌ Password verification failed:", err)
	}
}
