package analysis

import (
	"context"
	"github.com/ercanaziz/notiFY/Sema-Durgut/db" // Veritabanı bağlantısı
	"github.com/ercanaziz/notiFY/models"         // Modellerin bulunduğu paket eklendi
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// 1. GetHistoryFromDB: Belirli bir ürünün geçmiş fiyat hareketlerini getirir
func GetHistoryFromDB(productID string) ([]models.PriceHistory, error) {
	var history []models.PriceHistory
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Geçmiş fiyat verileri için koleksiyon
	collection := db.DB.Collection("price_histories")

	filter := bson.M{"product_id": productID}
	opts := options.Find().SetSort(bson.D{{Key: "recorded_at", Value: -1}})

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &history); err != nil {
		return nil, err
	}
	return history, nil
}

// 2. GetLowestPriceFromDB: Bir ürünün bugüne kadar gördüğü en dip fiyatı bulur
func GetLowestPriceFromDB(productID string) (float64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := db.DB.Collection("price_histories")
	filter := bson.M{"product_id": productID}

	// En düşük fiyatı (1) bulmak için sıralıyoruz
	opts := options.FindOne().SetSort(bson.D{{Key: "price", Value: 1}})

	var result models.PriceHistory
	err := collection.FindOne(ctx, filter, opts).Decode(&result)
	if err != nil {
		return 0, err
	}
	return result.Price, nil
}

// 3. GetStoreComparisons: Girilen kelimeye göre mağazalardaki güncel teklifleri karşılaştırır
func GetStoreComparisons(query string) ([]models.Offer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Arkadaşlarının kullandığı koleksiyon ismi: watchlist
	productsCollection := db.DB.Collection("watchlist")

	// Ekrandaki 'product_name' alanına göre Regex araması
	productFilter := bson.M{"product_name": bson.M{"$regex": primitive.Regex{Pattern: query, Options: "i"}}}

	cursor, err := productsCollection.Find(ctx, productFilter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var offers []models.Offer
	for cursor.Next(ctx) {
		var p struct {
			Brand string  `bson:"brand"`
			Price float64 `bson:"current_price"`
			URL   string  `bson:"product_url"`
		}
		if err := cursor.Decode(&p); err == nil {
			offers = append(offers, models.Offer{
				StoreName: p.Brand,
				Price:     p.Price,
				URL:       p.URL,
			})
		}
	}
	return offers, nil
}

// 4. GetAveragePriceFromDB: İndirim analizi için ortalama fiyat hesaplar
func GetAveragePriceFromDB(productID string) (float64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := db.DB.Collection("price_histories")

	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: bson.M{"product_id": productID}}},
		bson.D{{Key: "$group", Value: bson.M{
			"_id":      nil,
			"avgPrice": bson.M{"$avg": "$price"},
		}}},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)

	var result []struct {
		AvgPrice float64 `bson:"avgPrice"`
	}
	if err = cursor.All(ctx, &result); err != nil || len(result) == 0 {
		return 0, nil
	}

	return result[0].AvgPrice, nil
}

// 5. GetProductsByColor: Belirli bir renge sahip ürünleri filtreler
func GetProductsByColor(color string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := db.DB.Collection("watchlist")

	// Ekrandaki 'color' alanına göre filtre
	filter := bson.M{"color": bson.M{"$regex": primitive.Regex{Pattern: color, Options: "i"}}}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var productNames []string
	for cursor.Next(ctx) {
		var p struct {
			Name string `bson:"product_name"`
		}
		if err := cursor.Decode(&p); err == nil {
			productNames = append(productNames, p.Name)
		}
	}
	return productNames, nil
}

// 6. GetPriceForecast: Ortalama ve güncel fiyatı kıyaslayarak tahmin üretir
func GetPriceForecast(productID string) (string, float64, float64) {
	// 1. Ortalama fiyatı çek
	avgPrice, _ := GetAveragePriceFromDB(productID)

	// 2. Güncel fiyat (Gerçek veride bu watchlist'ten çekilmeli)
	currentPrice := 1450.0

	prediction := "STABLE"
	confidence := 0.70

	// Eğer veri yoksa tahmini stabil tut
	if avgPrice == 0 {
		return "STABLE", 0.50, currentPrice
	}

	// Basit Algoritma: Ortalamadan sapma kontrolü
	if currentPrice > avgPrice*1.10 {
		prediction = "DOWN"
		confidence = 0.85
	} else if currentPrice < avgPrice*0.90 {
		prediction = "UP"
		confidence = 0.80
	}

	estimatedNextMonth := (currentPrice + avgPrice) / 2
	return prediction, confidence, estimatedNextMonth
}

// GetChartDataFromDB: Grafik için veri hazırlar, geçmiş yoksa watchlist'ten güncel fiyatı çeker.
func GetChartDataFromDB(productID string) (models.ChartData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var labels []string
	var prices []float64

	// 1. ADIM: Önce geçmiş fiyatlara bak (price_histories)
	historyCol := db.DB.Collection("price_histories")
	filter := bson.M{"product_id": productID}
	cursor, err := historyCol.Find(ctx, filter)

	if err == nil {
		defer cursor.Close(ctx)
		for cursor.Next(ctx) {
			var h models.PriceHistory
			if err := cursor.Decode(&h); err == nil {
				labels = append(labels, h.RecordedAt.Format("02/01"))
				prices = append(prices, h.Price)
			}
		}
	}

	// 2. ADIM: EĞER GEÇMİŞ YOKSA (B Planı - Ekran görüntüsündeki veriyi kullan)
	if len(prices) == 0 {
		watchlistCol := db.DB.Collection("watchlist")

		// Gelen ID'yi MongoDB ObjectID'sine çevir
		objID, err := primitive.ObjectIDFromHex(productID)
		if err == nil {
			var p struct {
				CurrentPrice float64 `bson:"current_price"`
			}
			// Watchlist koleksiyonundan o anki fiyatı çek
			err := watchlistCol.FindOne(ctx, bson.M{"_id": objID}).Decode(&p)

			if err == nil {
				labels = append(labels, "Güncel")
				prices = append(prices, p.CurrentPrice) // İşte o 55000.5 buraya geliyor!
			}
		}
	}

	return models.ChartData{
		Labels:   labels,
		Datasets: []models.Dataset{{Label: "Ürün Fiyatı", Data: prices}},
	}, nil
}
