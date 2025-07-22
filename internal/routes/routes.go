package routes

import (
	"OppoCalypse/internal/handlers"
	"OppoCalypse/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	// Public routes
	router.GET("/login", handlers.ShowLoginPage)
	router.POST("/login", handlers.Login)

	// Authenticated routes
	auth := router.Group("/")
	auth.Use(middleware.AuthRequired())
	{
		auth.GET("/", handlers.GetTransactions)
		auth.GET("/transactions/new", handlers.NewTransactionForm)
		auth.POST("/transactions", handlers.CreateTransaction)
		auth.GET("/transactions/:id/edit", handlers.EditTransactionForm)
		auth.POST("/transactions/:id", handlers.UpdateTransaction)
		auth.POST("/transactions/:id/delete", handlers.DeleteTransaction)
		auth.POST("/logout", handlers.Logout)
	}
}
