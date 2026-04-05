package feedback

import (
	"fmt"
	"os"

	"notify-api/db"
	"notify-api/router"
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
