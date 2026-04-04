package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ---------------------------------------------------------
// 1. ORTAK VERİTABANI MODELİ (Sema, Ercan ve Nisa'nın Birleşimi)
// ---------------------------------------------------------
type Watchlist struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID       string             `json:"userId" bson:"user_id"`
	ProductName  string             `json:"productName" bson:"product_name"`
	Brand        string             `json:"brand" bson:"brand"`
	Color        string             `json:"color" bson:"color"`
	ProductURL   string             `json:"productUrl" bson:"product_url"`
	CurrentPrice float64            `json:"currentPrice" bson:"current_price"`
	Category     string             `json:"category" bson:"category"`
	WatchCount   int                `json:"watchCount" bson:"watch_count"`
	CreatedAt    time.Time          `json:"createdAt" bson:"created_at"`
}

// ---------------------------------------------------------
// 2. KULLANICI VE AUTH MODELLERİ (Betül)
// ---------------------------------------------------------
type User struct {
	ID              string    `json:"_id"`
	Email           string    `json:"email"`
	FirstName       string    `json:"firstName"`
	LastName        string    `json:"lastName"`
	IsEmailVerified bool      `json:"isEmailVerified"`
	CreatedOn       time.Time `json:"createdOn"`
}

type RegisterInput struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ProfileUpdateInput struct {
	FirstName                string                   `json:"firstName"`
	LastName                 string                   `json:"lastName"`
	CommunicationPreferences CommunicationPreferences `json:"communicationPreferences"`
}

type CommunicationPreferences struct {
	Newsletter       bool `json:"newsletter"`
	SmsNotifications bool `json:"smsNotifications"`
}

type PasswordUpdateInput struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

// ---------------------------------------------------------
// 3. ALARM VE BİLDİRİM MODELLERİ (Doğukan)
// ---------------------------------------------------------
type Alert struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID      string             `bson:"user_id" json:"user_id"`
	ProductID   string             `bson:"product_id" json:"product_id"`
	TargetPrice float64            `bson:"target_price" json:"target_price"`
	IsActive    bool               `bson:"is_active" json:"is_active"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
}

type AlertInput struct {
	ProductID   string  `json:"product_id" binding:"required"`
	TargetPrice float64 `json:"target_price" binding:"required"`
}

// ---------------------------------------------------------
// 4. SİSTEM VE ÜYELİK MODELLERİ (Ercan)
// ---------------------------------------------------------
type Feedback struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID      string             `json:"user_id" bson:"user_id"`
	Subject     string             `json:"subject" bson:"subject"`
	Message     string             `json:"message" bson:"message"`
	Type        string             `json:"type" bson:"type"`
	SubmittedAt time.Time          `json:"submitted_at" bson:"submitted_at"`
}

type SubscriptionPlan struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name"`
	MaxTracking int                `json:"max_tracking" bson:"max_tracking"`
	Price       float64            `json:"price" bson:"price"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}

// ---------------------------------------------------------
// 5. ANALİZ VE GRAFİK MODELLERİ (Sema)
// ---------------------------------------------------------
type PriceHistory struct {
	Price      float64   `json:"price" bson:"price"`
	Currency   string    `json:"currency" bson:"currency"`
	RecordedAt time.Time `json:"recordedAt" bson:"recorded_at"`
	StoreName  string    `json:"storeName" bson:"store_name"`
}

type ChartData struct {
	Labels   []string  `json:"labels"`
	Datasets []Dataset `json:"datasets"`
}

type Dataset struct {
	Label string    `json:"label"`
	Data  []float64 `json:"data"`
}

type PriceInfo struct {
	LowestPrice  float64 `json:"lowestPrice"`
	CurrentPrice float64 `json:"currentPrice"`
}

type ComparisonList struct {
	ProductName string  `json:"productName"`
	Offers      []Offer `json:"offers"`
}

type Offer struct {
	StoreName string  `json:"storeName"`
	Price     float64 `json:"price"`
	URL       string  `json:"url"`
}

type DiscountAnalysis struct {
	CurrentPrice       float64 `json:"currentPrice"`
	AveragePrice       float64 `json:"averagePrice"`
	DiscountPercentage int     `json:"discountPercentage"`
	Rating             string  `json:"rating"`
}

type ForecastData struct {
	ProductID               string  `json:"productId"`
	Prediction              string  `json:"prediction"`
	ConfidenceScore         float64 `json:"confidenceScore"`
	EstimatedPriceNextMonth float64 `json:"estimatedPriceNextMonth"`
}
