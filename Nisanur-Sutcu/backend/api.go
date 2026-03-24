package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type WatchlistItem struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID       string             `json:"user_id" bson:"user_id"`
	ProductName  string             `json:"product_name" bson:"product_name"`
	ProductURL   string             `json:"product_url" bson:"product_url"`
	CurrentPrice string             `json:"current_price" bson:"current_price"`
	Category     string             `json:"category" bson:"category"`
	WatchCount   int                `json:"watch_count" bson:"watch_count"`
}

var collection *mongo.Collection

func init() {
	clientOptions := options.Client().ApplyURI("mongodb+srv://nisa:stcnokta@test.ek07wik.mongodb.net/?appName=test")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	if err = client.Ping(context.TODO(), nil); err != nil {
		log.Fatal("MongoDB bağlantısı başarısız!", err)
	}
	fmt.Println("✅ MongoDB Bağlantısı Başarılı!")
	collection = client.Database("notiFY_DB").Collection("watchlist")

}

func main() {
	r := gin.Default()

	// 1. ÜRÜN ARAMA
	r.GET("/products/search", func(c *gin.Context) {
		query := c.Query("q")
		if query == "" {
			c.JSON(400, gin.H{"error": "Arama terimi boş olamaz!"})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		filter := bson.M{"product_name": bson.M{"$regex": query, "$options": "i"}}
		cursor, err := collection.Find(ctx, filter)
		if err != nil {
			c.JSON(500, gin.H{"error": "Arama sırasında hata oluştu!"})
			return
		}

		var results []WatchlistItem
		if err = cursor.All(ctx, &results); err != nil {
			c.JSON(500, gin.H{"error": "Sonuçlar okunamadı!"})
			return
		}
		c.JSON(200, results)
	})

	// 2. TAKİP LİSTESİNE EKLEME
	r.POST("/watchlist/add", func(c *gin.Context) {
		var item WatchlistItem
		if err := c.ShouldBindJSON(&item); err != nil {
			c.JSON(400, gin.H{"error": "Veri formatı hatalı!"})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		item.ID = primitive.NewObjectID()
		item.WatchCount = 1

		_, err := collection.InsertOne(ctx, item)
		if err != nil {
			c.JSON(500, gin.H{"error": "Ürün eklenemedi!"})
			return
		}
		c.JSON(201, gin.H{"message": "Ürün takip listesine eklendi!", "id": item.ID})
	})

	// 3. TAKİP LİSTESİNİ GÖRÜNTÜLEME
	r.GET("/watchlist", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		cursor, err := collection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(500, gin.H{"error": "Liste getirilemedi!"})
			return
		}

		var results []WatchlistItem
		if err = cursor.All(ctx, &results); err != nil {
			c.JSON(500, gin.H{"error": "Sonuçlar okunamadı!"})
			return
		}
		c.JSON(200, results)
	})

	// 4. TAKİP LİSTESİNDEN ÇIKARMA
	r.DELETE("/watchlist/:id", func(c *gin.Context) {
		idParam := c.Param("id")
		objectID, err := primitive.ObjectIDFromHex(idParam)
		if err != nil {
			c.JSON(400, gin.H{"error": "Geçersiz ID formatı!"})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		result, err := collection.DeleteOne(ctx, bson.M{"_id": objectID})
		if err != nil {
			c.JSON(500, gin.H{"error": "Silme işlemi başarısız!"})
			return
		}
		if result.DeletedCount == 0 {
			c.JSON(404, gin.H{"error": "Ürün bulunamadı!"})
			return
		}
		c.JSON(200, gin.H{"message": "Ürün takip listesinden çıkarıldı!"})
	})

	// 5. KATEGORİ LİSTELEME
	r.GET("/products/category", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		categories, err := collection.Distinct(ctx, "category", bson.M{})
		if err != nil {
			c.JSON(500, gin.H{"error": "Kategoriler getirilemedi!"})
			return
		}
		c.JSON(200, gin.H{"categories": categories})
	})

	// 6. POPÜLER ÜRÜNLER (TRENDING)
	r.GET("/products/trending", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		opts := options.Find().
			SetSort(bson.D{{Key: "watch_count", Value: -1}}).
			SetLimit(10)

		cursor, err := collection.Find(ctx, bson.M{}, opts)
		if err != nil {
			c.JSON(500, gin.H{"error": "Trending ürünler getirilemedi!"})
			return
		}

		var results []WatchlistItem
		if err = cursor.All(ctx, &results); err != nil {
			c.JSON(500, gin.H{"error": "Sonuçlar okunamadı!"})
			return
		}
		c.JSON(200, results)
	})

	r.Run(":8080")
}
