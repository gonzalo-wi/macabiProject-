package main

import (
	"log"

	mealpersistence "macabi-back/internal/meal/infrastructure/persistence"
	projectpersistence "macabi-back/internal/project/infrastructure/persistence"
	"macabi-back/internal/shared/config"
	"macabi-back/internal/shared/database"
	userpersistence "macabi-back/internal/user/infrastructure/persistence"
	attendancepersistence "macabi-back/internal/attendance/infrastructure/persistence"
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

	if err := mealpersistence.RunMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	if err := projectpersistence.RunMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	if err := attendancepersistence.RunMigrations(db); err != nil {
    	log.Fatalf("Failed to run migrations: %v", err)
	}

	deps := BuildDependencies(db, cfg)

	r := SetupRouter(deps)

	log.Printf("Server starting on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
