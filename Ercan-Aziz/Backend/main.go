package main

import (
	"fmt"
	"os"

	"notify-api/db"
	"notify-api/router"
)

func main() {
	db.Connect()

	r := router.Setup()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Println("🚀 Sunucu port " + port + " üzerinde çalışıyor...")
	r.Run(":" + port)
}
