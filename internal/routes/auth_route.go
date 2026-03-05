package routes

import (
	"github.com/britinogn/bizkeeper/internal/handler"
	// "github.com/britinogn/bizkeeper/internal/middleware"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(rg *gin.RouterGroup, authHandler *handler.AuthHandler) {
	auth := rg.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}
}


func UserRoutes(rg *gin.RouterGroup, authHandler *handler.AuthHandler) {
	user := rg.Group("/settings")
	// user.Use(middleware.AuthMiddleware())
	{
		user.PUT("/update", authHandler.UpdateUser)
		user.DELETE("/delete", authHandler.DeleteUser)
		user.GET("/profile", authHandler.GetProfile)
	}

}