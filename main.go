package main

import (
	"os"

	backendDB "github.com/ercanaziz/notiFY/Ercan-Aziz/Backend/db"
	backendRouter "github.com/ercanaziz/notiFY/Ercan-Aziz/Backend/router"
	history "github.com/ercanaziz/notiFY/Sema-Durgut"
	historyDB "github.com/ercanaziz/notiFY/Sema-Durgut/db"
	"github.com/gin-gonic/gin"
)

func main() {

	backendDB.Connect()
	historyDB.Connect()

	r := gin.Default()
	backendRouter.RegisterRoutes(r)
	history.RegisterRoutes(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
