package db

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var WatchlistCol *mongo.Collection

func Connect() {
	uri := os.Getenv("MONGO_URI")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal("🚨 MongoDB Bağlantı Hatası:", err)
	}
	// Koleksiyonu global değişkene atıyoruz
	WatchlistCol = client.Database("notiFY_DB").Collection("watchlist")
}
