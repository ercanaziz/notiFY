package main

import (
	"fmt"
	"os"
	"strings"

	"nisanur-sutcu/db"       // Kendi modül adın
	"nisanur-sutcu/handlers" // Yazdığın fonksiyonlar

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

var mySigningKey = []byte("notify_token")

// 🔐 JWT AUTH MIDDLEWARE
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, gin.H{"error": "Yetkisiz işlem! Token bulunamadı."})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return mySigningKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(401, gin.H{"error": "Geçersiz token!"})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			userID := claims["id"].(string)
			c.Set("user_id", userID)
			c.Next()
		}
	}
}

func main() {
	_ = godotenv.Load()
	db.Connect()

	r := gin.Default()

	// CORS Ayarları
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// --- ROTALAR ---
	r.GET("/products/trending", handlers.GetTrending)
	r.GET("/products/detail/:id", handlers.GetProductDetail)

	authorized := r.Group("/")
	authorized.Use(AuthMiddleware())
	{
		authorized.GET("/products/search", handlers.SearchProducts)
		authorized.GET("/watchlist", handlers.GetWatchlist)
		authorized.POST("/watchlist/add", handlers.AddToWatchlist)
		authorized.DELETE("/watchlist/:id", handlers.DeleteFromWatchlist)
		authorized.GET("/products/category", handlers.GetCategories)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
	}
	fmt.Println("🚀 Sunucu hazır, Postman'i açabilirsin!")
	r.Run(":" + port)
}
