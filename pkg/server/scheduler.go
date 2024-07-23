package server

import (
	"log"
	"time"

	"github.com/abyan-dev/auth/pkg/model"
	"github.com/go-co-op/gocron/v2"
	"gorm.io/gorm"
)

func CreateCleanupScheduler(db *gorm.DB) gocron.Scheduler {
	if db == nil {
		log.Fatal("Database connection is required for cleanup scheduler")
	}

	s, err := gocron.NewScheduler()
	if err != nil {
		log.Fatalf("failed to create scheduler: %v", err)
	}

	_, err = s.NewJob(
		gocron.DurationJob(time.Minute*10), gocron.NewTask(
			func() {
				err := db.Where("verified = ?", false).Delete(&model.User{}).Error
				if err != nil {
					log.Printf("Error cleaning up users: %v", err)
				}
			},
		),
	)

	if err != nil {
		log.Fatalf("failed to schedule job: %v", err)
	}

	return s
}
