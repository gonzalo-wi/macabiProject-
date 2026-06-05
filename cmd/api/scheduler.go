package main

import (
	"context"
	"log"

	"github.com/robfig/cron/v3"

	eventusecases "macabi-back/internal/event/application/usecases"
)

func StartScheduler(sendReminders *eventusecases.SendEventReminders) *cron.Cron {
	c := cron.New()

	c.AddFunc("0 9 * * *", func() {
		log.Println("scheduler: running event reminder job")
		if err := sendReminders.Execute(context.Background()); err != nil {
			log.Printf("scheduler: event reminder job failed: %v", err)
		}
	})

	c.Start()
	log.Println("scheduler: started (event reminders at 09:00 UTC daily)")
	return c
}
