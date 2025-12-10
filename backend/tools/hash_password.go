package main

import (
	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run hash_password.go <password>")
		os.Exit(1)
	}

	password := os.Args[1]
	
	// Generate hash with cost 10 (same as in the application)
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		fmt.Println("Error generating hash:", err)
		os.Exit(1)
	}

	fmt.Println("\n✅ Password hashed successfully!")
	fmt.Printf("Password: %s\n", password)
	fmt.Printf("Hash: %s\n\n", string(hash))
	
	// Verify the hash works
	err = bcrypt.CompareHashAndPassword(hash, []byte(password))
	if err == nil {
		fmt.Println("✅ Verification successful - hash is valid")
	} else {
		fmt.Println("❌ Verification failed - hash is invalid")
	}
	
	fmt.Println("\nSQL Update Statement:")
	fmt.Printf("UPDATE users SET password = '%s' WHERE username = 'admin';\n", string(hash))
}
