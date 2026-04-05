package auth

import (
    "context"
    "net/http"
    "time"
    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
)

type AuthHandler struct {
    UserCollection *mongo.Collection
}

func NewAuthHandler(col *mongo.Collection) *AuthHandler {
    return &AuthHandler{UserCollection: col}
}

func (h *AuthHandler) Register(c *gin.Context) {
    var input RegisterInput
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Veri formatı hatalı"})
        return
    }

    newUser := bson.M{
        "email":     input.Email,
        "password":  input.Password,
        "firstName": input.FirstName,
        "lastName":  input.LastName,
        "isEmailVerified": false,
        "createdOn": time.Now(),
    }

    _, err := h.UserCollection.InsertOne(context.Background(), newUser)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Kayıt başarısız"})
        return
    }
    c.JSON(http.StatusCreated, gin.H{"message": "Başarıyla kayıt oldun!"})
}

func (h *AuthHandler) Login(c *gin.Context) {
    var input LoginInput
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Giriş bilgileri eksik"})
        return
    }

    var foundUser bson.M
    err := h.UserCollection.FindOne(context.Background(), bson.M{
        "email":    input.Email,
        "password": input.Password,
    }).Decode(&foundUser)

    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "E-posta veya şifre hatalı"})
        return
    }

    // Token oluşturma
    claims := jwt.MapClaims{
        "id":  foundUser["_id"].(primitive.ObjectID).Hex(),
        "exp": time.Now().Add(time.Hour * 24).Unix(),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, _ := token.SignedString(mySigningKey)

    c.JSON(http.StatusOK, gin.H{"token": tokenString})
}