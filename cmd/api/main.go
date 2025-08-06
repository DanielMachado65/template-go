package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"template-go.com/internal/db"
	"template-go.com/internal/server"
)

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func main() {
	mongoURI := getenv("MONGO_URI", "mongodb://localhost:27017")
	dbName := getenv("DB_NAME", "appdb")
	port := getenv("PORT", "8080")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to Mongo
	client, database, err := db.Connect(ctx, mongoURI, dbName)
	if err != nil {
		log.Fatalf("failed to connect to mongo: %v", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = client.Disconnect(ctx)
	}()

	// Router
	r := server.NewRouter(database)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	// Start server
	go func() {
		log.Printf("listening on :%s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutting down...")

	ctxShut, cancelShut := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShut()
	if err := srv.Shutdown(ctxShut); err != nil {
		log.Printf("server shutdown error: %v", err)
	}
	log.Println("bye")
}
