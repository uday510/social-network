package main

import (
	"log"

	"github.com/uday510/go-crud-app/internal/db"
	"github.com/uday510/go-crud-app/internal/env"
	store2 "github.com/uday510/go-crud-app/internal/store"
)

func main() {
	log.Println("Starting application...")

	addr := env.GetString("DB_ADDR", "postgres://user:password@localhost/social?sslmode=disable")
	log.Println("Connecting to DB...")

	conn, err := db.New(addr, 3, 3, "15m")
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	log.Println("Connected to DB")

	store := store2.NewStorage(conn)
	log.Println("Initialized storage layer")

	log.Println("Seeding database...")
	db.Seed(store, conn)

	log.Println("Database seeding complete.")
}
