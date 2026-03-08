package routes

import (
	"github.com/britinogn/bizkeeper/internal/handler"
	"github.com/britinogn/bizkeeper/internal/middleware"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(rg *gin.RouterGroup, authHandler *handler.AuthHandler, authLimiter *middleware.RateLimiter) {
	auth := rg.Group("/auth")
	auth.Use(authLimiter.Middleware())
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}
}

func UserRoutes(rg *gin.RouterGroup, authHandler *handler.AuthHandler, apiLimiter *middleware.RateLimiter) {
	user := rg.Group("/settings")
	user.Use(apiLimiter.Middleware())
	{
		user.PUT("/update", authHandler.UpdateUser)
		user.DELETE("/delete", authHandler.DeleteUser)
		user.GET("/profile", authHandler.GetProfile)
	}
}