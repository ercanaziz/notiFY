package router

import (
	"context"
	"net/http"

	
    "github.com/ercanaziz/notiFY/Betul-Erkoc/Backend/handlers"
    "github.com/ercanaziz/notiFY/Betul-Erkoc/Backend/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Public rotalar
	r.POST("/api/auth/register", handlers.Register)
	r.POST("/api/auth/login", handlers.Login)

	// Korumalı rotalar
	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.PUT("/users/profile", func(c *gin.Context) {
			withUserID(c, handlers.UpdateProfile)
		})

		protected.DELETE("/users/profile", func(c *gin.Context) {
			withUserID(c, handlers.DeleteProfile)
		})

		protected.PATCH("/users/password", func(c *gin.Context) {
			withUserID(c, handlers.ChangePassword)
		})

		protected.POST("/auth/logout", func(c *gin.Context) {
			withUserID(c, handlers.Logout)
		})
	}

	return r
}

// withUserID - Gin context'indeki userID'yi standart http.Request context'ine taşır
func withUserID(c *gin.Context, handler func(http.ResponseWriter, *http.Request)) {
	userID, _ := c.Get("userID")
	ctx := context.WithValue(c.Request.Context(), "userID", userID)
	handler(c.Writer, c.Request.WithContext(ctx))
}
