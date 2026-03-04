package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/britinogn/bizkeeper/internal/db"
)

func main() {
	// Load env
	db.Init()

	// Connect to database
	ctx := context.Background()
	_, err := db.ConnectPostgres(ctx)
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}
	defer db.Close()

	// Setup Gin
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "BizKeeper is running"})
	})

	// Start server
	log.Println("Server running on port 8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}