package handlers

import (
	"context"
	"time"

	"notify-api/db"
	"notify-api/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetPlans(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := db.PlanCollection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(500, gin.H{"error": "Planlar getirilemedi!"})
		return
	}

	var plans []models.SubscriptionPlan
	if err = cursor.All(ctx, &plans); err != nil {
		c.JSON(500, gin.H{"error": "Sonuçlar okunamadı!"})
		return
	}
	if plans == nil {
		plans = []models.SubscriptionPlan{}
	}
	c.JSON(200, gin.H{"plans": plans, "total": len(plans)})
}

func UpdatePlan(c *gin.Context) {
	var req models.SubscriptionPlan
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Veri formatı hatalı! " + err.Error()})
		return
	}

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

	var updatedPlan models.SubscriptionPlan
	err := db.PlanCollection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updatedPlan)
	if err != nil {
		c.JSON(500, gin.H{"error": "Plan güncellenemedi!"})
		return
	}

	c.JSON(200, gin.H{
		"message": "Plan başarıyla güncellendi!",
		"plan":    updatedPlan,
	})
}
