package main

import (
	"fmt"
	"log"
	"permit-app/helper"
)

func main2() {
	fmt.Println("Testing Email Configuration...")
	fmt.Println("====================================")

	config := helper.GetEmailConfig()
	fmt.Printf("SMTP Host: %s\n", config.Host)
	fmt.Printf("SMTP Port: %d\n", config.Port)
	fmt.Printf("Username: %s\n", config.Username)
	fmt.Printf("From Address: %s\n", config.From)
	fmt.Printf("From Name: %s\n", config.FromName)
	fmt.Println("====================================")

	// Test send email
	testEmail := []string{"ferdianrafli125@gmail.com"} // Ganti dengan email tujuan test
	subject := "Test Email - Permit Management System"
	body := `
<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background-color: #4f46e5; color: white; padding: 20px; text-align: center; }
        .content { background-color: #f9f9f9; padding: 20px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Test Email</h1>
        </div>
        <div class="content">
            <p>This is a test email from Permit Management System.</p>
            <p>If you receive this email, the email configuration is working correctly.</p>
        </div>
    </div>
</body>
</html>
`

	fmt.Println("\nSending test email...")
	err := helper.SendEmail(testEmail, subject, body)
	if err != nil {
		log.Fatalf("Failed to send email: %v", err)
	}

	fmt.Println("✓ Test email sent successfully!")
	fmt.Println("Please check your inbox.")
}
