package main

import (
	"log"

	"github.com/Mayank-Bihani/eta-service/config"
	"github.com/Mayank-Bihani/eta-service/db"
	"github.com/Mayank-Bihani/eta-service/handlers"
	"github.com/Mayank-Bihani/eta-service/router"
	"github.com/Mayank-Bihani/eta-service/services"
)

func main() {
	// Load config
	cfg := config.Load()

	// Init DB connections
	postgres := db.NewPostgres(cfg.DBUrl)
	redis := db.NewRedis(cfg.RedisUrl)

	// Create orders table if not exists
	_, err := postgres.Exec(`
		CREATE TABLE IF NOT EXISTS orders (
			id SERIAL PRIMARY KEY,
			restaurant_id TEXT NOT NULL,
			delivery_lat DOUBLE PRECISION NOT NULL,
			delivery_lng DOUBLE PRECISION NOT NULL,
			item_count INT NOT NULL,
			estimated_eta INT NOT NULL,
			created_at TIMESTAMP DEFAULT NOW()
		)
	`)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	// Wire dependencies
	etaService := services.NewETAService(postgres, redis)
	etaHandler := handlers.NewETAHandler(etaService, postgres, redis)

	r := router.Setup(etaHandler)

	// Start server
	log.Printf("Server starting on port %s", cfg.Port)
	r.Run(":" + cfg.Port)
}
