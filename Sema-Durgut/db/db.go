package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Global değişkenlerimiz (mongo.Database kullanıyoruz)
var MongoClient *mongo.Client
var DB *mongo.Database

func Connect() {
	//  MONGO_URI kullanacağız
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		log.Fatal("MONGO_URI ortam değişkeni bulunamadı!")
	}

	// Bağlantı için 10 saniyelik bir zaman aşımı süresi belirliyoruz
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Atlas'a bağlan
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("MongoDB bağlantı hatası: ", err)
	}

	// Bağlantının gerçekten kurulup kurulmadığını test et (Ping at)
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("MongoDB Atlas'a ulaşılamıyor: ", err)
	}

	MongoClient = client

	// Ekran görüntüsündeki veritabanı adını tam buraya yazıyoruz
	DB = client.Database("notiFY_DB")

	fmt.Println("MongoDB Atlas bağlantısı başarıyla kuruldu!")
}
