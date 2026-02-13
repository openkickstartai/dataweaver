package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/dataweaver/internal/api"
	"github.com/dataweaver/internal/config"
	"github.com/dataweaver/internal/database"
)

func main() {
	cfg := config.Load()
	
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	router := gin.Default()
	apiHandler := api.NewHandler(db)
	
	v1 := router.Group("/api/v1")
	{
		v1.GET("/health", apiHandler.Health)
		v1.POST("/schemas/detect", apiHandler.DetectSchema)
		v1.POST("/workflows", apiHandler.CreateWorkflow)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("DataWeaver server starting on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	log.Printf("Received signal %s, shutting down gracefully...", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Println("Server exited cleanly")

	
	log.Printf("DataWeaver server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}