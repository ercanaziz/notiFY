package main

import (
	"fmt"
	"log"
    
    "github.com/ercanaziz/notiFY/Betul-Erkoc/Backend/router"
    "github.com/ercanaziz/notiFY/Betul-Erkoc/Backend/db"
    
)

func main() {
	db.InitDB()

	r := router.SetupRouter()

	fmt.Println("🚀 notiFY Sunucusu 8080 portunda hazır...")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Sunucu başlatılamadı: ", err)
	}
}
