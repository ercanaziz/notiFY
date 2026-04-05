package router

import (
	"github.com/ercanaziz/notiFY/Ercan-Aziz/Backend/handlers"
	"github.com/ercanaziz/notiFY/Ercan-Aziz/Backend/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Setup() *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:    []string{"Origin", "Content-Type", "X-User-Id", "X-Admin-Key"},
	}))

	// Ürünler
	r.GET("/products", handlers.GetProducts)
	r.GET("/products/categories", handlers.GetCategories)

	// Support
	r.POST("/support/feedback", middleware.MockAuthMiddleware(), handlers.PostFeedback)

	// Admin
	r.GET("/admin/subscription-plans", middleware.AdminAuthMiddleware(), handlers.GetPlans)
	r.PUT("/admin/subscription-plans", middleware.AdminAuthMiddleware(), handlers.UpdatePlan)

	return r
}
