package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5" // Yeni eklenen kütüphane
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Model Yapısı
type WatchlistItem struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID       string             `json:"user_id" bson:"user_id"`
	ProductName  string             `json:"product_name" bson:"product_name"`
	Brand        string             `json:"brand" bson:"brand"`
	Color        string             `json:"color" bson:"color"`
	ProductURL   string             `json:"product_url" bson:"product_url"`
	CurrentPrice float64            `json:"current_price" bson:"current_price"`
	Category     string             `json:"category" bson:"category"`
	WatchCount   int                `json:"watch_count" bson:"watch_count"`
	CreatedAt    time.Time          `json:"created_at" bson:"created_at"`
}

var collection *mongo.Collection
var mySigningKey = []byte("notify_token") // Betül'ün gönderdiği anahtar

func init() {
	_ = godotenv.Load()

	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		log.Fatal("🚨 HATA: MONGO_URI bulunamadı!")
	}

	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal("❌ MongoDB bağlantı hatası: ", err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal("❌ MongoDB'ye ulaşılamıyor: ", err)
	}

	fmt.Println("✅ MongoDB Bağlantısı Başarılı!")
	collection = client.Database("notiFY_DB").Collection("watchlist")
}

// 🔐 GERÇEK JWT AUTH MIDDLEWARE
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, gin.H{"error": "Yetkisiz işlem! Token bulunamadı."})
			c.Abort()
			return
		}

		// "Bearer <token>" formatını ayıkla
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// İmza yöntemini kontrol et
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("beklenmeyen imza yöntemi")
			}
			return mySigningKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(401, gin.H{"error": "Geçersiz veya süresi dolmuş token!"})
			c.Abort()
			return
		}

		// Token içindeki 'id'yi al ve Context'e kaydet
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			userID := claims["id"].(string)
			c.Set("user_id", userID)
			c.Next()
		} else {
			c.JSON(401, gin.H{"error": "Token verisi okunamadı!"})
			c.Abort()
		}
	}
}

func main() {
	r := gin.Default()
	r.Use(cors.Default())

	// --- PUBLIC (HERKESE AÇIK) ROTALAR ---

	// Popüler Ürünler
	r.GET("/products/trending", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		opts := options.Find().SetSort(bson.D{{Key: "watch_count", Value: -1}}).SetLimit(10)
		cursor, err := collection.Find(ctx, bson.M{}, opts)
		if err != nil {
			c.JSON(500, gin.H{"error": "Veriler getirilemedi"})
			return
		}
		var results []WatchlistItem
		cursor.All(ctx, &results)
		c.JSON(200, results)
	})

	// Ürün Detayı (Tıklanma artırır)
	r.GET("/products/detail/:id", func(c *gin.Context) {
		idParam := c.Param("id")
		objectID, _ := primitive.ObjectIDFromHex(idParam)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		filter := bson.M{"_id": objectID}
		update := bson.M{"$inc": bson.M{"watch_count": 1}}
		var updatedItem WatchlistItem
		err := collection.FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&updatedItem)

		if err != nil {
			c.JSON(404, gin.H{"error": "Ürün bulunamadı"})
			return
		}
		c.JSON(200, updatedItem)
	})

	// --- PRIVATE (KORUMALI) ROTALAR ---
	authorized := r.Group("/")
	authorized.Use(AuthMiddleware())
	{
		// Arama (Kişiye Özel)
		authorized.GET("/products/search", func(c *gin.Context) {
			query := c.Query("q")
			userID := c.MustGet("user_id").(string)
			filter := bson.M{
				"user_id":      userID,
				"product_name": bson.M{"$regex": query, "$options": "i"},
			}
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			cursor, _ := collection.Find(ctx, filter)
			var results []WatchlistItem
			cursor.All(ctx, &results)
			c.JSON(200, results)
		})

		// Takip Listesi Görüntüle
		authorized.GET("/watchlist", func(c *gin.Context) {
			userID := c.MustGet("user_id").(string)
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			cursor, _ := collection.Find(ctx, bson.M{"user_id": userID})
			var results []WatchlistItem = []WatchlistItem{}
			cursor.All(ctx, &results)
			c.JSON(200, results)
		})

		// Takip Listesine Ekle
		authorized.POST("/watchlist/add", func(c *gin.Context) {
			var item WatchlistItem
			if err := c.ShouldBindJSON(&item); err != nil {
				c.JSON(400, gin.H{"error": "Geçersiz veri"})
				return
			}
			item.ID = primitive.NewObjectID()
			item.UserID = c.MustGet("user_id").(string)
			item.WatchCount = 1
			item.CreatedAt = time.Now()

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			collection.InsertOne(ctx, item)
			c.JSON(201, gin.H{"message": "Eklendi", "id": item.ID})
		})

		// Takip Listesinden Sil
		authorized.DELETE("/watchlist/:id", func(c *gin.Context) {
			idParam := c.Param("id")
			objectID, _ := primitive.ObjectIDFromHex(idParam)
			userID := c.MustGet("user_id").(string)

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			filter := bson.M{"_id": objectID, "user_id": userID}
			res, _ := collection.DeleteOne(ctx, filter)
			if res.DeletedCount == 0 {
				c.JSON(404, gin.H{"error": "Bulunamadı veya yetki yok"})
				return
			}
			c.JSON(200, gin.H{"message": "Silindi"})
		})

		// KATEGORİ LİSTELEME

		authorized.GET("/products/category", func(c *gin.Context) {
			userID := c.MustGet("user_id").(string)
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			filter := bson.M{"user_id": userID}
			categories, err := collection.Distinct(ctx, "category", filter)

			if err != nil {
				c.JSON(500, gin.H{"error": "Kategoriler getirilemedi!"})
				return
			}

			if categories == nil {
				categories = []interface{}{}
			}
			c.JSON(200, gin.H{"categories": categories})
		})
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run("0.0.0.0:" + port)
}
