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



type Product struct {
	ID           primitive.ObjectID `json:"id"            bson:"_id,omitempty"`
	UserID       string             `json:"user_id"       bson:"user_id"`
	ProductName  string             `json:"product_name"  bson:"product_name"`
	Brand        string             `json:"brand"         bson:"brand"`
	Color        string             `json:"color"         bson:"color"`
	ProductURL   string             `json:"product_url"   bson:"product_url"`
	CurrentPrice float64            `json:"current_price" bson:"current_price"`
	Category     string             `json:"category"      bson:"category"`
	WatchCount   int                `json:"watch_count"   bson:"watch_count"`
	CreatedAt    time.Time          `json:"created_at"    bson:"created_at"`
}

type Feedback struct {
	ID          primitive.ObjectID `json:"id"           bson:"_id,omitempty"`
	UserID      string             `json:"user_id"      bson:"user_id"`
	Subject     string             `json:"subject"      bson:"subject"`
	Message     string             `json:"message"      bson:"message"`
	Type        string             `json:"type"         bson:"type"` // "bug" | "suggestion" | "other"
	SubmittedAt time.Time          `json:"submitted_at" bson:"submitted_at"`
}

type SubscriptionPlan struct {
	ID          primitive.ObjectID `json:"id"           bson:"_id,omitempty"`
	Name        string             `json:"name"         bson:"name"` // "free" | "basic" | "pro" | "enterprise"
	MaxTracking int                `json:"max_tracking" bson:"max_tracking"`
	Price       float64            `json:"price"        bson:"price"`
	UpdatedAt   time.Time          `json:"updated_at"   bson:"updated_at"`
}
