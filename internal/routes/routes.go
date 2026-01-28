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
	auth.GET("/export", handlers.ShowExportPage)
	auth.POST("/export", handlers.ExportTransactions)
	auth.GET("/budgets", handlers.ListBudgets)
	auth.GET("/budgets/new", handlers.NewBudgetForm)
	auth.POST("/budgets", handlers.CreateBudget)
	auth.GET("/budgets/:id/edit", handlers.EditBudgetForm)
	auth.POST("/budgets/:id", handlers.UpdateBudget)
	auth.POST("/budgets/:id/delete", handlers.DeleteBudget)
	auth.GET("/graphs", handlers.ShowGraphs)
		auth.POST("/logout", handlers.Logout)
	}

	// Admin routes
	admin := router.Group("/admin")
	admin.Use(middleware.AuthRequired())
	admin.Use(middleware.AdminRequired())
	{
		admin.GET("/users", handlers.ListUsers)
		admin.GET("/users/new", handlers.NewUserForm)
		admin.POST("/users", handlers.CreateUser)
		admin.GET("/users/:id/edit", handlers.EditUserForm)
		admin.POST("/users/:id", handlers.UpdateUser)
		admin.POST("/users/:id/delete", handlers.DeleteUser)
		admin.POST("/users/:id/reset-pin", handlers.ResetUserPIN)
	}
}
