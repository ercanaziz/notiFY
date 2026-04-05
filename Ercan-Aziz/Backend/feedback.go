package feedback

import (
	"fmt"
	"os"

	"Ercan-Aziz/Backend/db"
	"Ercan-Aziz/Backend/router"
)

func Start() {
	db.Connect()

	r := router.Setup()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Println("🚀 Sunucu port " + port + " üzerinde çalışıyor...")
	r.Run(":" + port)
}
