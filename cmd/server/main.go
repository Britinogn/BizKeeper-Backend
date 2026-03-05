package main

import (
	"context"
	"log"

	"github.com/britinogn/bizkeeper/config"
	"github.com/britinogn/bizkeeper/internal/db"
	"github.com/britinogn/bizkeeper/internal/handler"
	"github.com/britinogn/bizkeeper/internal/repository"
	"github.com/britinogn/bizkeeper/internal/routes"
	"github.com/britinogn/bizkeeper/internal/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type handlers struct {
	auth *handler.AuthHandler
	purchase *handler.PurchaseHandler	
}

func initHandlers(database *gorm.DB) *handlers {
	// Repositories
	userRepo := repository.NewUserRepository(database)
	purchaseRepo := repository.NewPurchaseRepository(database)

	// Services
	authService := services.NewAuthService(userRepo)
	purchaseService := services.NewPurchaseService(purchaseRepo)

	// Handlers
	return &handlers{
		auth: handler.NewAuthHandler(authService),
		purchase: handler.NewPurchaseHandler(purchaseService),
	}
}

func main() {
	db.Init()
	cfg := config.Load()
	database, err := db.ConnectPostgres(context.Background(), cfg)
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}
	defer db.Close()
	log.Println("✓ Database connected successfully")

	h := initHandlers(database)

	r := gin.Default()
	routes.SetupRoutes(r, h.auth, h.purchase)

	for _, route := range r.Routes() {
		log.Println(route.Method, route.Path)
	}

	log.Println("Server running on port 8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}