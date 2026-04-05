package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv" // BU SATIR EKSİKTİ: .env okumak için gerekli
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client
var DB *mongo.Database

func Connect() {
	// .env dosyasını yükle (Bu kısım eklenmeli)
	// Dosya yolunu bir üst klasörde olduğu için "../.env" veya direkt ".env" denemelisin
	err := godotenv.Load(".env")
	if err != nil {
		// Eğer .env bulunamazsa en azından log basalım ama durmayalım (belki sistemde tanımlıdır)
		fmt.Println(".env dosyası yüklenemedi, sistem değişkenlerine bakılıyor...")
	}

	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		log.Fatal("MONGO_URI ortam değişkeni bulunamadı!")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("MongoDB bağlantı hatası: ", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("MongoDB Atlas'a ulaşılamıyor: ", err)
	}

	MongoClient = client
	DB = client.Database("notiFY_DB")

	fmt.Println("MongoDB Atlas bağlantısı başarıyla kuruldu!")
}
