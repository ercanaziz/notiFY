package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ProductCollection *mongo.Collection
var FeedbackCollection *mongo.Collection
var PlanCollection *mongo.Collection

func Connect() {
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		log.Fatal("🚨 HATA: MONGO_URI bulunamadı!")
	}
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	if err = client.Ping(context.TODO(), nil); err != nil {
		log.Fatal("MongoDB bağlantısı başarısız!", err)
	}
	fmt.Println("✅ MongoDB Bağlantısı Başarılı!")

	priceDB := client.Database("priceTracker_DB")
	FeedbackCollection = priceDB.Collection("feedback")
	PlanCollection = priceDB.Collection("subscription_plans")

	notifyDB := client.Database("notiFY_DB")
	ProductCollection = notifyDB.Collection("watchlist")
}
