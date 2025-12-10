package scheduler

import (
	"fmt"
	"log"
	"os"
	"permit-app/service/notificationService"
	"strconv"
	"time"
)

type Scheduler struct {
	notificationService notificationService.NotificationService
	ticker              *time.Ticker
	done                chan bool
}

func NewScheduler(notificationService notificationService.NotificationService) *Scheduler {
	return &Scheduler{
		notificationService: notificationService,
		done:                make(chan bool),
	}
}

// Start memulai scheduler untuk check permit expiry setiap hari pada jam 8 pagi
func (s *Scheduler) Start() {
	// Jalankan pertama kali saat start
	go func() {
		log.Println("Notification Scheduler: Initial check for expiring permits")
		if err := s.notificationService.CheckAndSendExpiryNotifications(); err != nil {
			log.Printf("Error in initial notification check: %v", err)
		}
	}()

	// Schedule untuk check setiap hari jam 8 pagi
	go func() {
		for {
			now := time.Now()
			// Hitung waktu ke jam berikutnya berdasarkan konfigurasi
			next := time.Date(now.Year(), now.Month(), now.Day(), s.GetScheduledHour(), s.GetScheduledMinute(), 0, 0, now.Location())
			if now.After(next) {
				next = next.Add(24 * time.Hour)
			}
			
			duration := time.Until(next)
			log.Printf("Notification Scheduler: Next run scheduled at %v (in %v)", next.Format("2006-01-02 15:04:05"), duration)
			
			timer := time.NewTimer(duration)
			
			select {
			case <-timer.C:
				log.Println("Notification Scheduler: Running scheduled check for expiring permits")
				if err := s.notificationService.CheckAndSendExpiryNotifications(); err != nil {
					log.Printf("Error in scheduled notification check: %v", err)
				}
			case <-s.done:
				timer.Stop()
				log.Println("Notification Scheduler: Stopped")
				return
			}
		}
	}()

	log.Println("Notification Scheduler: Started successfully")
}

// StartWithInterval untuk testing - check setiap interval tertentu
func (s *Scheduler) StartWithInterval(interval time.Duration) {
	s.ticker = time.NewTicker(interval)
	
	// Jalankan pertama kali
	go func() {
		log.Println("Notification Scheduler: Initial check for expiring permits")
		if err := s.notificationService.CheckAndSendExpiryNotifications(); err != nil {
			log.Printf("Error in initial notification check: %v", err)
		}
	}()
	
	go func() {
		for {
			select {
			case <-s.ticker.C:
				log.Println("Notification Scheduler: Running periodic check for expiring permits")
				if err := s.notificationService.CheckAndSendExpiryNotifications(); err != nil {
					log.Printf("Error in notification check: %v", err)
				}
			case <-s.done:
				s.ticker.Stop()
				log.Println("Notification Scheduler: Stopped")
				return
			}
		}
	}()
	
	fmt.Printf("Notification Scheduler: Started with %v interval\n", interval)
}

// Stop menghentikan scheduler
func (s *Scheduler) Stop() {
	s.done <- true
}

// GetScheduledHour returns the hour when scheduler should run (default: 8)
func (s *Scheduler) GetScheduledHour() int {
	hour := 8 // default
	if hourStr := getEnv("SCHEDULER_HOUR", "8"); hourStr != "" {
		if h, err := strconv.Atoi(hourStr); err == nil && h >= 0 && h <= 23 {
			hour = h
		}
	}
	return hour
}

// GetScheduledMinute returns the minute when scheduler should run (default: 0)
func (s *Scheduler) GetScheduledMinute() int {
	minute := 0 // default
	if minuteStr := getEnv("SCHEDULER_MINUTE", "0"); minuteStr != "" {
		if m, err := strconv.Atoi(minuteStr); err == nil && m >= 0 && m <= 59 {
			minute = m
		}
	}
	return minute
}

// getEnv gets environment variable with default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
