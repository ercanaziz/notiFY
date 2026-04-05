package handlers

import (
	"context"
	"net/http"
	"time"

	//"nisanur-sutcu/db"
	//"nisanur-sutcu/models"
	
	"github.com/ercanaziz/notiFY/nisanur-sutcu/backend/db"
	"github.com/ercanaziz/notiFY/nisanur-sutcu/backend/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// 1. Trending Fonksiyonu
func GetTrending(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	findOptions := options.Find().SetSort(bson.D{{Key: "watch_count", Value: -1}}).SetLimit(10)
	cursor, err := db.WatchlistCol.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Veriler alınamadı"})
		return
	}
	var results []models.WatchlistItem = []models.WatchlistItem{}
	cursor.All(ctx, &results)
	c.JSON(http.StatusOK, results)
}

// 2. Add To Watchlist Fonksiyonu
func AddToWatchlist(c *gin.Context) {
	var item models.WatchlistItem
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(400, gin.H{"error": "Geçersiz veri"})
		return
	}
	item.ID = primitive.NewObjectID()
	item.UserID = c.MustGet("user_id").(string)
	item.CreatedAt = time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	db.WatchlistCol.InsertOne(ctx, item)
	c.JSON(201, gin.H{"message": "Eklendi", "id": item.ID})
}

// 3. Delete Fonksiyonu
func DeleteFromWatchlist(c *gin.Context) {
	idParam := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(400, gin.H{"error": "Geçersiz ID formatı!"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"_id": objectID}

	res, err := db.WatchlistCol.DeleteOne(ctx, filter)
	if err != nil {
		c.JSON(500, gin.H{"error": "Veritabanı hatası oluştu"})
		return
	}

	if res.DeletedCount == 0 {
		c.JSON(404, gin.H{"message": "Silinecek ürün bulunamadı. ID doğru mu?"})
		return
	}
	c.JSON(200, gin.H{"message": "Ürün başarıyla silindi"})
}

// 4. Arama Fonksiyonu (Kullanıcıya özel arama yapar)
func SearchProducts(c *gin.Context) {
	query := c.Query("q")
	userID := c.MustGet("user_id").(string)

	filter := bson.M{
		"user_id":      userID,
		"product_name": bson.M{"$regex": query, "$options": "i"},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, _ := db.WatchlistCol.Find(ctx, filter)
	var results []models.WatchlistItem = []models.WatchlistItem{}
	cursor.All(ctx, &results)
	c.JSON(200, results)
}

// 5. Ürün Detayı (Tıklanınca watch_count değerini 1 artırır)
func GetProductDetail(c *gin.Context) {
	idParam := c.Param("id")
	objectID, _ := primitive.ObjectIDFromHex(idParam)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"_id": objectID}
	update := bson.M{"$inc": bson.M{"watch_count": 1}}

	var updatedItem models.WatchlistItem
	err := db.WatchlistCol.FindOneAndUpdate(
		ctx,
		filter,
		update,
		nil,
	).Decode(&updatedItem)

	if err != nil {
		c.JSON(404, gin.H{"error": "Ürün bulunamadı"})
		return
	}
	c.JSON(200, updatedItem)
}

// 6. Kategori Listeleme
func GetCategories(c *gin.Context) {
	userID := c.MustGet("user_id").(string)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"user_id": userID}
	categories, err := db.WatchlistCol.Distinct(ctx, "category", filter)

	if err != nil {
		c.JSON(500, gin.H{"error": "Kategoriler getirilemedi!"})
		return
	}
	c.JSON(200, gin.H{"categories": categories})
}

// 7. Takip Listesini Görüntüle (Kullanıcıya özel)
func GetWatchlist(c *gin.Context) {
	userID := c.MustGet("user_id").(string)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Sadece bu kullanıcıya ait olan ürünleri getir
	cursor, err := db.WatchlistCol.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Liste alınamadı"})
		return
	}

	var results []models.WatchlistItem = []models.WatchlistItem{}
	cursor.All(ctx, &results)
	c.JSON(http.StatusOK, results)
}
