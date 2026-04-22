package db

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewPostgres(dbUrl string) *sqlx.DB {
	db, err := sqlx.Connect("postgres", dbUrl)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}

	// Connection pool settings for high concurrency
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)

	log.Println("PostgreSQL connected successfully")
	return db
}