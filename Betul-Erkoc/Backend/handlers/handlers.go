package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

    "github.com/ercanaziz/notiFY/Betul-Erkoc/Backend/db"
    "github.com/ercanaziz/notiFY/Betul-Erkoc/Backend/middleware"
    "github.com/ercanaziz/notiFY/Betul-Erkoc/Backend/models"

    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive" )

// Register - Yeni kullanıcı kaydı
func Register(c *gin.Context) {
	var input models.RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Veri formatı hatalı"})
		return
	}

	newUser := bson.M{
		"email":           input.Email,
		"password":        input.Password,
		"firstName":       input.FirstName,
		"lastName":        input.LastName,
		"isEmailVerified": false,
		"createdOn":       time.Now(),
	}

	_, err := db.UserCollection.InsertOne(context.Background(), newUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Kayıt başarısız"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Başarıyla kayıt oldun!"})
}

// Login - Kullanıcı girişi ve JWT token üretimi
func Login(c *gin.Context) {
	var input models.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Giriş bilgileri eksik"})
		return
	}

	var foundUser bson.M
	err := db.UserCollection.FindOne(context.Background(), bson.M{
		"email":    input.Email,
		"password": input.Password,
	}).Decode(&foundUser)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "E-posta veya şifre hatalı"})
		return
	}

	claims := jwt.MapClaims{
		"id":  foundUser["_id"].(primitive.ObjectID).Hex(),
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(middleware.MySigningKey)

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

// UpdateProfile - Profil güncelleme (PUT)
func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Yetkisiz"})
		return
	}
	objID, _ := primitive.ObjectIDFromHex(userID)

	var input models.ProfileUpdateInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Geçersiz veri", http.StatusBadRequest)
		return
	}

	update := bson.M{
		"$set": bson.M{
			"firstName":                input.FirstName,
			"lastName":                 input.LastName,
			"communicationPreferences": input.CommunicationPreferences,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db.UserCollection.UpdateOne(ctx, bson.M{"_id": objID}, update)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Profil güncellendi"})
}

// DeleteProfile - Hesap silme (DELETE)
func DeleteProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Yetkisiz"})
		return
	}
	objID, _ := primitive.ObjectIDFromHex(userID)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := db.UserCollection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		http.Error(w, "Silme hatası", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ChangePassword - Şifre değiştirme (PATCH)
func ChangePassword(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	objID, _ := primitive.ObjectIDFromHex(userID)

	var input models.PasswordUpdateInput
	json.NewDecoder(r.Body).Decode(&input)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{"$set": bson.M{"password": input.NewPassword}}
	db.UserCollection.UpdateOne(ctx, bson.M{"_id": objID}, update)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Şifre güncellendi"})
}

// Logout - Çıkış yap (POST)
func Logout(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(string)

	w.Header().Set("Content-Type", "application/json")

	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Zaten giriş yapılmamış veya yetkisiz erişim"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Başarıyla çıkış yapıldı. Güle güle, ID: " + userID,
	})
}
