package alarm

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AlarmRepository struct {
	Collection *mongo.Collection
}

// Parametreyi *mongo.Collection yapıyoruz ki main.go bize tam hedefi göstersin
func NewAlarmRepository(col *mongo.Collection) *AlarmRepository {
	return &AlarmRepository{
		Collection: col,
	}
}

// Veritabanına yeni alarm kaydeder
func (r *AlarmRepository) CreateAlert(alert Alert) (Alert, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := r.Collection.InsertOne(ctx, alert)
	if err == nil {
		alert.ID = res.InsertedID.(primitive.ObjectID)
	}
	return alert, err
}

// Kullanıcının aktif alarmlarını getirir
func (r *AlarmRepository) GetUserAlerts(userID string) ([]Alert, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var alerts []Alert
	filter := bson.M{"user_id": userID, "is_active": true}

	cursor, err := r.Collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	err = cursor.All(ctx, &alerts)
	return alerts, err
}
func (r *AlarmRepository) DeleteAlert(alertID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 1. ID'yi ObjectID'ye çevir
	objID, err := primitive.ObjectIDFromHex(alertID)
	if err != nil {
		fmt.Printf("❌ ID Çevrim Hatası: %s (Format yanlış olabilir)\n", alertID)
		return err
	}

	// 2. Filtre: SADECE bu ID'ye sahip olanı bul
	filter := bson.M{"_id": objID}

	// 3. Güncelleme: is_active'i FALSE yap
	update := bson.M{"$set": bson.M{"is_active": false}}

	// 🔍 DEBUG: Hangi koleksiyonda neyi arıyoruz?
	fmt.Printf("🔍 Silme Denemesi -> Tablo: %s | ID: %s\n", r.Collection.Name(), alertID)

	result, err := r.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		fmt.Printf("❌ DB Hatası: %v\n", err)
		return err
	}

	// 4. SONUÇ ANALİZİ
	if result.MatchedCount == 0 {
		fmt.Printf("⚠️ BULUNAMADI! Veritabanında '%s' ID'li bir kayıt yok. Postman'deki ID'yi kontrol et!\n", alertID)
		return fmt.Errorf("kayit bulunamadi")
	}

	fmt.Printf("🗑️ BAŞARI! Kayıt pasife çekildi. Matched: %d, Modified: %d\n", result.MatchedCount, result.ModifiedCount)
	return nil
}

func (r *AlarmRepository) GetActiveAlerts() ([]Alert, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var alerts []Alert
	// Sadece aktif olanları getir
	filter := bson.M{"is_active": true}
	cursor, err := r.Collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	err = cursor.All(ctx, &alerts)
	return alerts, err
}

func (r *AlarmRepository) DeactivateAlert(id interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var objID primitive.ObjectID
	var err error

	// ID'nin tipine göre işlem yapıyoruz
	switch v := id.(type) {
	case string:
		objID, err = primitive.ObjectIDFromHex(v)
		if err != nil {
			return err
		}
	case primitive.ObjectID:
		objID = v
	default:
		return fmt.Errorf("geçersiz ID tipi")
	}

	filter := bson.M{"_id": objID}
	update := bson.M{"$set": bson.M{"is_active": false}}

	_, err = r.Collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *AlarmRepository) UpdateAlertPrice(id string, newPrice float64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objID, _ := primitive.ObjectIDFromHex(id)

	// Sadece target_price alanını güncelliyoruz
	update := bson.M{"$set": bson.M{"target_price": newPrice}}

	_, err := r.Collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	return err
}

// GetProductPrice, notiFY_DB veritabanındaki 'products' tablosundan güncel fiyatı çeker
func (r *AlarmRepository) GetProductPrice(productName string) (float64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client := r.Collection.Database().Client()
	productsCollection := client.Database("notiFY_DB").Collection("watchlist")
	var item struct {
		CurrentPrice float64 `bson:"current_price"`
	}

	// KRİTİK DEĞİŞİKLİK: Birebir eşleşme yerine, büyük/küçük harf
	// umursamayan (Options: "i") bir Regex araması yapıyoruz!
	filter := bson.M{
		"product_name": primitive.Regex{
			Pattern: "^" + productName + "$",
			Options: "i",
		},
	}

	err := productsCollection.FindOne(ctx, filter).Decode(&item)
	if err != nil {
		return 0, err
	}

	return item.CurrentPrice, nil
}
