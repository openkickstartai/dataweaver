package main

import (
	"log"
	"net/http"
	"os"

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
		v1.GET("/workflows/:id", apiHandler.GetWorkflow)
		v1.POST("/transform", apiHandler.TransformData)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	log.Printf("DataWeaver server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}