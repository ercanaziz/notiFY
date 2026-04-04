package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)



func main() {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:    []string{"Origin", "Content-Type", "X-User-Id", "X-Admin-Key"},
	}))

	// =========================================================================
	// 1. KULLANICI GERİ BİLDİRİMİ ALMA
	// POST /support/feedback
	// Header: X-User-Id zorunlu
	// Body: { "subject": "...", "message": "...", "type": "bug|suggestion|other" }
	// =========================================================================
	r.POST("/support/feedback", MockAuthMiddleware(), func(c *gin.Context) {
		userID := c.MustGet("user_id").(string)

		var feedback Feedback
		if err := c.ShouldBindJSON(&feedback); err != nil {
			c.JSON(400, gin.H{"error": "Veri formatı hatalı! " + err.Error()})
			return
		}

		// Alan doğrulamaları
		if feedback.Subject == "" || feedback.Message == "" {
			c.JSON(400, gin.H{"error": "Subject ve message alanları boş olamaz!"})
			return
		}
		if feedback.Type != "bug" && feedback.Type != "suggestion" && feedback.Type != "other" {
			c.JSON(400, gin.H{"error": "type alanı yalnızca 'bug', 'suggestion' veya 'other' olabilir!"})
			return
		}

		feedback.ID = primitive.NewObjectID()
		feedback.UserID = userID
		feedback.SubmittedAt = time.Now()

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		_, err := feedbackCollection.InsertOne(ctx, feedback)
		if err != nil {
			c.JSON(500, gin.H{"error": "Geri bildirim kaydedilemedi!"})
			return
		}

		c.JSON(201, gin.H{
			"message":      "Geri bildiriminiz alındı. En kısa sürede değerlendirilerek dönüş yapılacaktır.",
			"feedback_id":  feedback.ID,
			"submitted_at": feedback.SubmittedAt,
		})
	})

	// =========================================================================
	// 2. ÜRÜN ÇEŞİDİNE (KATEGORİ) GÖRE LİSTELEME
	// GET /products/categories              → Tüm benzersiz kategorileri listeler
	// GET /products/categories?name=Telefon → O kategorideki ürünleri getirir
	// =========================================================================
	r.GET("/products/categories", func(c *gin.Context) {
		categoryName := c.Query("name")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Kategori adı verilmemişse → benzersiz tüm kategorileri döndür
		if categoryName == "" {
			categories, err := productCollection.Distinct(ctx, "category", bson.M{})
			if err != nil {
				c.JSON(500, gin.H{"error": "Kategoriler getirilemedi!"})
				return
			}
			if categories == nil {
				categories = []interface{}{}
			}
			c.JSON(200, gin.H{"categories": categories})
			return
		}

		// Kategori adı verilmişse → o kategorideki ürünleri getir
		filter := bson.M{"category": bson.M{"$regex": categoryName, "$options": "i"}}
		cursor, err := productCollection.Find(ctx, filter)
		if err != nil {
			c.JSON(500, gin.H{"error": "Ürünler getirilemedi!"})
			return
		}

		var results []Product
		if err = cursor.All(ctx, &results); err != nil {
			c.JSON(500, gin.H{"error": "Sonuçlar okunamadı!"})
			return
		}
		if results == nil {
			results = []Product{}
		}
		c.JSON(200, gin.H{"products": results, "total": len(results)})
	})

	// =========================================================================
	// 3. MARKAYA GÖRE LİSTELEME
	// 4. FİYATA GÖRE SIRALAMA
	// 5. TARİHE GÖRE SIRALAMA
	//
	// GET /products
	// Query Parametreleri:
	//   ?brand=Apple              → Markaya göre filtrele
	//   ?sort_by=price&order=asc  → Fiyata göre artan sırala
	//   ?sort_by=price&order=desc → Fiyata göre azalan sırala
	//   ?sort_by=date&order=asc   → Eskiden yeniye sırala
	//   ?sort_by=date&order=desc  → Yeniden eskiye sırala
	//
	// Hepsi birlikte kullanılabilir:
	//   ?brand=Apple&sort_by=price&order=desc
	// =========================================================================
	r.GET("/products", func(c *gin.Context) {
		brand := c.Query("brand")
		sortBy := c.Query("sort_by") // "price" veya "date"
		order := c.Query("order")    // "asc" veya "desc"

		// Parametre doğrulamaları
		if sortBy != "" && sortBy != "price" && sortBy != "date" {
			c.JSON(400, gin.H{"error": "sort_by yalnızca 'price' veya 'date' olabilir!"})
			return
		}
		if order != "" && order != "asc" && order != "desc" {
			c.JSON(400, gin.H{"error": "order yalnızca 'asc' veya 'desc' olabilir!"})
			return
		}

		// ── Filtre oluştur ───────────────────────────────────────────────────
		filter := bson.M{}
		if brand != "" {
			// Büyük/küçük harf duyarsız arama
			filter["brand"] = bson.M{"$regex": brand, "$options": "i"}
		}

		// ── Sıralama oluştur ─────────────────────────────────────────────────
		sortValue := 1 // varsayılan: artan (asc)
		if order == "desc" {
			sortValue = -1
		}

		findOpts := options.Find()
		switch sortBy {
		case "price":
			findOpts.SetSort(bson.D{{Key: "current_price", Value: sortValue}})
		case "date":
			findOpts.SetSort(bson.D{{Key: "created_at", Value: sortValue}})
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		cursor, err := productCollection.Find(ctx, filter, findOpts)
		if err != nil {
			c.JSON(500, gin.H{"error": "Ürünler getirilemedi!"})
			return
		}

		var results []Product
		if err = cursor.All(ctx, &results); err != nil {
			c.JSON(500, gin.H{"error": "Sonuçlar okunamadı!"})
			return
		}
		if results == nil {
			results = []Product{}
		}
		c.JSON(200, gin.H{"products": results, "total": len(results)})
	})

	// =========================================================================
	// 6. ABONELİK PLANI BELİRLEME
	// PUT  /admin/subscription-plans → Belirli planın limitini güncelle
	// GET  /admin/subscription-plans → Tüm planları listele
	// Header: X-Admin-Key zorunlu
	// =========================================================================

	// Tüm planları listele
	r.GET("/admin/subscription-plans", AdminAuthMiddleware(), func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		cursor, err := planCollection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(500, gin.H{"error": "Planlar getirilemedi!"})
			return
		}

		var plans []SubscriptionPlan
		if err = cursor.All(ctx, &plans); err != nil {
			c.JSON(500, gin.H{"error": "Sonuçlar okunamadı!"})
			return
		}
		if plans == nil {
			plans = []SubscriptionPlan{}
		}
		c.JSON(200, gin.H{"plans": plans, "total": len(plans)})
	})

	// Plan limitini güncelle (yoksa oluştur)
	r.PUT("/admin/subscription-plans", AdminAuthMiddleware(), func(c *gin.Context) {
		var req SubscriptionPlan
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": "Veri formatı hatalı! " + err.Error()})
			return
		}

		// Alan doğrulamaları
		validPlans := map[string]bool{"free": true, "basic": true, "pro": true, "enterprise": true}
		if !validPlans[req.Name] {
			c.JSON(400, gin.H{"error": "name yalnızca 'free', 'basic', 'pro' veya 'enterprise' olabilir!"})
			return
		}
		if req.MaxTracking <= 0 {
			c.JSON(400, gin.H{"error": "max_tracking en az 1 olmalıdır!"})
			return
		}
		if req.Price < 0 {
			c.JSON(400, gin.H{"error": "Fiyat negatif olamaz!"})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Planı güncelle; yoksa yeni oluştur (upsert)
		filter := bson.M{"name": req.Name}
		update := bson.M{
			"$set": bson.M{
				"max_tracking": req.MaxTracking,
				"price":        req.Price,
				"updated_at":   time.Now(),
			},
			"$setOnInsert": bson.M{
				"_id": primitive.NewObjectID(),
			},
		}
		opts := options.FindOneAndUpdate().
			SetUpsert(true).
			SetReturnDocument(options.After)

		var updatedPlan SubscriptionPlan
		err := planCollection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updatedPlan)
		if err != nil {
			c.JSON(500, gin.H{"error": "Plan güncellenemedi!"})
			return
		}

		c.JSON(200, gin.H{
			"message": "Plan başarıyla güncellendi!",
			"plan":    updatedPlan,
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Println("🚀 Sunucu port " + port + " üzerinde çalışıyor...")
	r.Run(":" + port)
}