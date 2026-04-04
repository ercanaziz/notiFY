package handlers

import (
	"context"
	"time"

	"notify-api/db"
	"notify-api/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetCategories(c *gin.Context) {
	categoryName := c.Query("name")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if categoryName == "" {
		categories, err := db.ProductCollection.Distinct(ctx, "category", bson.M{})
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

	filter := bson.M{"category": bson.M{"$regex": categoryName, "$options": "i"}}
	cursor, err := db.ProductCollection.Find(ctx, filter)
	if err != nil {
		c.JSON(500, gin.H{"error": "Ürünler getirilemedi!"})
		return
	}

	var results []models.Product
	if err = cursor.All(ctx, &results); err != nil {
		c.JSON(500, gin.H{"error": "Sonuçlar okunamadı!"})
		return
	}
	if results == nil {
		results = []models.Product{}
	}
	c.JSON(200, gin.H{"products": results, "total": len(results)})
}

func GetProducts(c *gin.Context) {
	brand  := c.Query("brand")
	sortBy := c.Query("sort_by")
	order  := c.Query("order")

	if sortBy != "" && sortBy != "price" && sortBy != "date" {
		c.JSON(400, gin.H{"error": "sort_by yalnızca 'price' veya 'date' olabilir!"})
		return
	}
	if order != "" && order != "asc" && order != "desc" {
		c.JSON(400, gin.H{"error": "order yalnızca 'asc' veya 'desc' olabilir!"})
		return
	}

	filter := bson.M{}
	if brand != "" {
		filter["brand"] = bson.M{"$regex": brand, "$options": "i"}
	}

	sortValue := 1
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

	cursor, err := db.ProductCollection.Find(ctx, filter, findOpts)
	if err != nil {
		c.JSON(500, gin.H{"error": "Ürünler getirilemedi!"})
		return
	}

	var results []models.Product
	if err = cursor.All(ctx, &results); err != nil {
		c.JSON(500, gin.H{"error": "Sonuçlar okunamadı!"})
		return
	}
	if results == nil {
		results = []models.Product{}
	}
	c.JSON(200, gin.H{"products": results, "total": len(results)})
}
