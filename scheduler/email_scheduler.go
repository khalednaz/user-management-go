package scheduler

import (
	"fmt"
	"sync"
	"time"

	"go-user-app/utils"
)

var (
	mu       sync.Mutex
	running  = false
	stopChan chan struct{}
	emailNow string
)

func StartEmailScheduler(email string) {
	mu.Lock()
	defer mu.Unlock()

	// Stop existing scheduler if running
	if running {
		fmt.Println("Stopping previous email scheduler for:", emailNow)
		stopChan <- struct{}{}
		running = false
	}

	emailNow = email
	stopChan = make(chan struct{})
	running = true

	fmt.Println("Starting email scheduler for:", email)

	go func() {
		ticker := time.NewTicker(60 * time.Second) // Adjust interval as needed
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				err := utils.SendEmail(emailNow, "Reminder", "This is a scheduled email.")
				if err != nil {
					fmt.Printf("Failed to send email to %s: %v\n", emailNow, err)
				} else {
					fmt.Printf("Email successfully sent to %s\n", emailNow)
				}
			case <-stopChan:
				fmt.Println("Email scheduler stopped for:", emailNow)
				return
			}
		}
	}()
}
