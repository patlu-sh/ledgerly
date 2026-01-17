package routes

import (
	"ledgerly/handlers"
	"ledgerly/middleware"
	"ledgerly/models"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	h := handlers.NewHandler()

	// Rate Limiter (100 req/s, burst 200)
	r.Use(middleware.DefaultRateLimiter().Middleware())

	// Swagger UI
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Auth Routes
	auth := r.Group("/auth")
	{
		auth.POST("/login", h.Login)
	}

	// Protected Routes
	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware())

	// Petty Cash Routes
	pc := protected.Group("/petty-cash")
	{
		// POST has custom RBAC in handler to allow employees to create debit transactions
		pc.POST("", h.CreatePettyCashTransaction) 
		pc.GET("", middleware.PermissionMiddleware(models.PermissionPettyCashViewList), h.ListPettyCashTransactions)
		pc.GET("/balance", middleware.PermissionMiddleware(models.PermissionPettyCashViewBalance), h.GetPettyCashBalance)
	}

	// Expense Routes
	ex := protected.Group("/expenses")
	{
		ex.POST("", middleware.PermissionMiddleware(models.PermissionExpensesCreate), h.CreateExpense)
		ex.GET("", middleware.PermissionMiddleware(models.PermissionExpensesViewOwn), h.ListExpenses)
	}

	// Reporting Routes
	rp := protected.Group("/reports")
	rp.Use(middleware.PermissionMiddleware(models.PermissionReportsView))
	{
		rp.GET("/expenses-summary", h.GetExpenseSummary)
		rp.GET("/petty-cash-summary", h.GetPettyCashSummary)
	}

	return r
}