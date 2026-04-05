package alarm

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Veritabanına kaydedilecek ana şema
type Alert struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID      string             `bson:"user_id" json:"user_id"`       // İleride Betül'ün JWT token'ından gelecek
	ProductID   string             `bson:"product_id" json:"product_id"` // Sema'nın ürün tablosuyla eşleşecek ID
	TargetPrice float64            `bson:"target_price" json:"target_price"`
	IsActive    bool               `bson:"is_active" json:"is_active"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
}

// Postman'den / API'den gelecek olan veri
type AlertInput struct {
	ProductID   string  `json:"product_id" binding:"required"`
	TargetPrice float64 `json:"target_price" binding:"required"`
}