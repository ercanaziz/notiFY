package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client
var UserCollection *mongo.Collection

func InitDB() {
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		uri = "mongodb+srv://teamnotify:notiFY@test.ek07wik.mongodb.net/?appName=test"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("Bağlantı ayarı hatası:", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("MongoDB'ye bağlanılamadı (IP izni sorun olabilir):", err)
	}

	MongoClient = client
	UserCollection = client.Database("notiFY_DB").Collection("users")

	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true),
	}

	_, err = UserCollection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		fmt.Println("Index oluşturulurken bir not: Muhtemelen zaten var veya bağlantı kısıtlı.")
	}

	fmt.Println("🚀 MongoDB veritabanına ve Index yapısına başarıyla bağlanıldı!")
}
