package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"nudgepay/internal/api"
	"nudgepay/internal/config"
	"nudgepay/internal/db"
	"nudgepay/internal/services"
)

func main() {
	cfg := config.Load()
	database, err := db.New(cfg.DBPath)
	if err != nil {
		log.Fatalf("db error: %v", err)
	}
	defer database.Close()

	app := api.NewApp(database, cfg)

	if cfg.WorkerEnabled {
		go runWorker(database)
	}

	go func() {
		if err := app.Listen(cfg.Addr); err != nil {
			log.Printf("server stopped: %v", err)
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Printf("shutdown error: %v", err)
	}
}

func runWorker(database *sql.DB) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		sendDueForAllOrgs(database)
		<-ticker.C
	}
}

func sendDueForAllOrgs(database *sql.DB) {
	rows, err := database.Query(`SELECT id FROM organizations`)
	if err != nil {
		log.Printf("worker org query error: %v", err)
		return
	}
	defer rows.Close()

	now := time.Now().UTC()
	for rows.Next() {
		var orgID string
		if err := rows.Scan(&orgID); err != nil {
			log.Printf("worker org scan error: %v", err)
			continue
		}
		if _, err := services.SendDueReminders(database, orgID, now); err != nil {
			log.Printf("worker send error: %v", err)
		}
	}
}
