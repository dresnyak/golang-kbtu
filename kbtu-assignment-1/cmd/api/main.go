package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"kbtu-assignment-1/internal/handlers"
	"kbtu-assignment-1/internal/middleware"
	"kbtu-assignment-1/internal/models"
)

func main() {
	store := models.NewTaskStore()
	taskHandler := &handlers.TaskHandler{Store: store}

	mux := http.NewServeMux()
	mux.Handle("/tasks", taskHandler)

	handler := middleware.Logging(middleware.Auth(mux))

	srv := &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	go func() {
		log.Println("Server starting on :8080")
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Println("Server stopped")
}
