package routes

import (
    "github.com/britinogn/bizkeeper/internal/handler"
    "github.com/gin-gonic/gin"
)

func ReorderRoutes(rg *gin.RouterGroup, reorderHandler *handler.ReorderHandler) {
    reorder := rg.Group("/reorder-reminders")
    {
        reorder.GET("", reorderHandler.GetReorderReminders)
    }
}