package main

import (
	"fmt"
	"log"

    "github.com/ercanaziz/notiFY/Betul-Erkoc/Backend/db"
    "github.com/ercanaziz/notiFY/Betul-Erkoc/Backend/models"
    "github.com/ercanaziz/notiFY/Betul-Erkoc/Backend/middleware"

    "github.com/gin-gonic/gin"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	db.InitDB()

	r := router.SetupRouter()

	fmt.Println("🚀 notiFY Sunucusu 8080 portunda hazır...")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Sunucu başlatılamadı: ", err)
	}
}
