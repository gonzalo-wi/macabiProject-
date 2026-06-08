package main

import (
	"log"

	apideps "macabi-back/cmd/api/deps"
	"macabi-back/internal/shared/config"
	"macabi-back/internal/shared/database"
	userpersistence "macabi-back/internal/user/infrastructure/persistence"
)

func main() {
	cfg := config.Load()
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Invalid configuration: %v", err)
	}

	db := database.NewPostgresConnection(cfg.DSN())

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get underlying DB: %v", err)
	}
	defer sqlDB.Close()

	if err := userpersistence.RunMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	d := apideps.Build(db, cfg)

	scheduler := StartScheduler(d.SendEventReminders)
	defer scheduler.Stop()

	r := SetupRouter(d)

	log.Printf("Server starting on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
