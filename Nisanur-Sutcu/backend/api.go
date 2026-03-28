package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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

func init() {
	// 1. Yerel geliştirme için .env dosyasını yükle.
	_ = godotenv.Load()

	// 2. MONGO_URI'yi işletim sistemi değişkenlerinden (Environment Variables) çek.

	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		log.Fatal("🚨 HATA: MONGO_URI bulunamadı! Lütfen .env dosyanızı veya Render ayarlarınızı kontrol edin.")
	}

	// 3. MongoDB bağlantı seçeneklerini ayarla
	clientOptions := options.Client().ApplyURI(uri)

	// 4. Bağlantıyı başlat
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal("❌ MongoDB bağlantı hatası: ", err)
	}

	// 5. Bağlantının gerçekten kurulup kurulmadığını test et (Ping)
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal("❌ MongoDB'ye ulaşılamıyor (Ping başarısız): ", err)
	}

	fmt.Println("✅ MongoDB Bağlantısı Başarılı!")

	// 6. Koleksiyonu tanımla
	collection = client.Database("notiFY_DB").Collection("watchlist")
}

// 🛡️ SAHTE GİRİŞ (MOCK AUTH) MİDDLEWARE'İ
func MockAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Postman'de Headers kısmına eklediğimiz X-User-Id değerini okur
		userID := c.GetHeader("X-User-Id")

		if userID == "" {
			c.JSON(401, gin.H{"error": "Yetkisiz işlem! Lütfen Header kısmına X-User-Id ekleyin."})
			c.Abort() // İsteği burada keser, API'ye ulaşmasını engeller
			return
		}

		// ID'yi API içinde kullanabilmek için Gin Context'ine kaydediyoruz
		c.Set("user_id", userID)
		c.Next()
	}
}

func main() {
	r := gin.Default()
	r.Use(cors.Default())

	// 1. ÜRÜN ARAMA

	r.GET("/products/search", MockAuthMiddleware(), func(c *gin.Context) {
		query := c.Query("q")
		if query == "" {
			c.JSON(400, gin.H{"error": "Arama terimi boş olamaz!"})
			return
		}

		// İstek yapan kullanıcıyı middleware'den alıyoruz
		userID := c.MustGet("user_id").(string)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// SADECE bu kullanıcının listesinde, bu kelimeyi içeren ürünleri ara
		filter := bson.M{
			"user_id":      userID,
			"product_name": bson.M{"$regex": query, "$options": "i"},
		}

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

		if results == nil {
			results = []WatchlistItem{}
		}
		c.JSON(200, results)
	})

	// 5. KATEGORİ LİSTELEME (Kişiye Özel)
	r.GET("/products/category", MockAuthMiddleware(), func(c *gin.Context) {
		// İstek yapan kullanıcıyı alıyoruz
		userID := c.MustGet("user_id").(string)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// SADECE bu kullanıcıya ait ürünlerdeki benzersiz kategorileri bul
		filter := bson.M{"user_id": userID}
		categories, err := collection.Distinct(ctx, "category", filter)

		if err != nil {
			c.JSON(500, gin.H{"error": "Kategoriler getirilemedi!"})
			return
		}

		if categories == nil {
			categories = []interface{}{} // Boş liste dönsün
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

	// ÜRÜN DETAYI
	r.GET("/products/detail/:id", func(c *gin.Context) {
		idParam := c.Param("id")
		objectID, err := primitive.ObjectIDFromHex(idParam)
		if err != nil {
			c.JSON(400, gin.H{"error": "Geçersiz ID formatı!"})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		filter := bson.M{"_id": objectID}
		update := bson.M{"$inc": bson.M{"watch_count": 1}}
		opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

		var updatedItem WatchlistItem
		err = collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updatedItem)

		if err != nil {
			c.JSON(404, gin.H{"error": "Ürün bulunamadı!"})
			return
		}

		c.JSON(200, updatedItem)
	})

	// --------------------------------------------------------
	// 🔒 KORUMALI (PRIVATE) API'LER (Sadece giriş yapanlar)
	// MockAuthMiddleware kullandığımıza dikkat et!
	// --------------------------------------------------------

	// 2. TAKİP LİSTESİNE EKLEME
	r.POST("/watchlist/add", MockAuthMiddleware(), func(c *gin.Context) {
		var item WatchlistItem
		if err := c.ShouldBindJSON(&item); err != nil {
			c.JSON(400, gin.H{"error": "Veri formatı hatalı!"})
			return
		}

		// Middleware'den kullanıcının ID'sini çekiyoruz
		userID := c.MustGet("user_id").(string)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// TEMİZLENMİŞ VE GÜNCEL HALİ
		item.ID = primitive.NewObjectID()
		item.UserID = userID
		item.WatchCount = 1
		item.CreatedAt = time.Now()

		_, err := collection.InsertOne(ctx, item)
		if err != nil {
			c.JSON(500, gin.H{"error": "Ürün eklenemedi!"})
			return
		}
		c.JSON(201, gin.H{"message": "Ürün takip listesine eklendi!", "id": item.ID})
	})

	// 3. TAKİP LİSTESİNİ GÖRÜNTÜLEME
	r.GET("/watchlist", MockAuthMiddleware(), func(c *gin.Context) {
		// Kimin listesini getireceğimizi öğreniyoruz
		userID := c.MustGet("user_id").(string)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// SADECE bu kullanıcıya ait olanları getir!
		filter := bson.M{"user_id": userID}
		cursor, err := collection.Find(ctx, filter)
		if err != nil {
			c.JSON(500, gin.H{"error": "Liste getirilemedi!"})
			return
		}

		var results []WatchlistItem
		if err = cursor.All(ctx, &results); err != nil {
			c.JSON(500, gin.H{"error": "Sonuçlar okunamadı!"})
			return
		}

		// Eğer liste boşsa null yerine boş dizi dönelim (Frontend için daha iyi)
		if results == nil {
			results = []WatchlistItem{}
		}
		c.JSON(200, results)
	})

	// 4. TAKİP LİSTESİNDEN ÇIKARMA
	r.DELETE("/watchlist/:id", MockAuthMiddleware(), func(c *gin.Context) {
		idParam := c.Param("id")
		objectID, err := primitive.ObjectIDFromHex(idParam)
		if err != nil {
			c.JSON(400, gin.H{"error": "Geçersiz ID formatı!"})
			return
		}

		// Silmek isteyen kullanıcı kim?
		userID := c.MustGet("user_id").(string)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Güvenlik: Silinecek ürün hem verilen ID'ye sahip olmalı, HEM DE bu kullanıcının olmalı!
		filter := bson.M{"_id": objectID, "user_id": userID}
		result, err := collection.DeleteOne(ctx, filter)
		if err != nil {
			c.JSON(500, gin.H{"error": "Silme işlemi başarısız!"})
			return
		}

		if result.DeletedCount == 0 {
			c.JSON(404, gin.H{"error": "Ürün bulunamadı veya bu ürünü silme yetkiniz yok!"})
			return
		}
		c.JSON(200, gin.H{"message": "Ürün takip listesinden çıkarıldı!"})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run("0.0.0.0:" + port)
}
