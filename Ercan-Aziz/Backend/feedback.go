package feedback

import (
	"fmt"
	"os"

	"github.com/ercanaziz/notiFY/Ercan-Aziz/Backend/db"
	"github.com/ercanaziz/notiFY/Ercan-Aziz/Backend/router"
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
