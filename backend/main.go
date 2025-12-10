package main

import (
	"log"
	"permit-app/database"
	"permit-app/helper"
	"permit-app/model"
	"permit-app/repo/notificationRepository"
	"permit-app/repo/permitRepository"
	"permit-app/repo/userRepository"
	"permit-app/routes"
	"permit-app/scheduler"
	"permit-app/service/notificationService"
	"strconv"
	"time"
)

func main() {
	db, err := database.ConnectDB()
	if err != nil {
		log.Fatal("Error connect to database", err)
		panic(err)
	}

	// Auto migrate notification table with silent logger to avoid constraint warnings
	sqlDB, _ := db.DB()
	if sqlDB != nil {
		// Check if notifications table exists
		var exists bool
		db.Raw("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'notifications')").Scan(&exists)
		
		if !exists {
			// Only migrate if table doesn't exist
			if err := db.AutoMigrate(&model.Notification{}); err != nil {
				log.Printf("Warning: Failed to auto migrate notification table: %v", err)
			} else {
				log.Println("Notifications table created successfully")
			}
		} else {
			log.Println("Notifications table already exists, skipping migration")
		}
	}

	// Initialize notification scheduler
	notificationRepo := notificationRepository.NewNotificationRepository(db)
	permitRepo := permitRepository.NewPermitRepository(db)
	userRepo := userRepository.NewUserRepository(db)
	notificationSvc := notificationService.NewNotificationService(notificationRepo, permitRepo, userRepo)
	
	notificationScheduler := scheduler.NewScheduler(notificationSvc)
	
	// Start scheduler based on mode
	schedulerMode := helper.GetEnv("SCHEDULER_MODE")
	if schedulerMode == "testing" {
		intervalMinutes := helper.GetEnv("SCHEDULER_INTERVAL_MINUTES")
		if intervalMinutes == "" {
			intervalMinutes = "5" // default 5 minutes
		}
		if minutes, err := strconv.Atoi(intervalMinutes); err == nil {
			notificationScheduler.StartWithInterval(time.Duration(minutes) * time.Minute)
		} else {
			log.Printf("Invalid SCHEDULER_INTERVAL_MINUTES, using default 5 minutes")
			notificationScheduler.StartWithInterval(5 * time.Minute)
		}
	} else {
		// production mode - run daily at scheduled time
		notificationScheduler.Start()
	}
	defer notificationScheduler.Stop()

	log.Println("Notification scheduler started successfully")

	app := routes.NewRoute(db)

	apiPort := helper.GetEnv("PORT")
	log.Fatal(app.Run(":" + apiPort))
}
