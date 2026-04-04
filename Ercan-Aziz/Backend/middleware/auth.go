package middleware

import "github.com/gin-gonic/gin"

func MockAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetHeader("X-User-Id")
		if userID == "" {
			c.JSON(401, gin.H{"error": "Yetkisiz işlem! Lütfen Header kısmına X-User-Id ekleyin."})
			c.Abort()
			return
		}
		c.Set("user_id", userID)
		c.Next()
	}
}

func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		adminKey := c.GetHeader("X-Admin-Key")
		if adminKey != "secret-admin-key" {
			c.JSON(401, gin.H{"error": "Yetkisiz işlem! Geçerli bir X-Admin-Key gereklidir."})
			c.Abort()
			return
		}
		c.Next()
	}
}
