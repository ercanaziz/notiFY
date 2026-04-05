package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	backendDB "github.com/ercanaziz/notiFY/Ercan-Aziz/Backend/db"
	backendRouter "github.com/ercanaziz/notiFY/Ercan-Aziz/Backend/router"
	historyDB "github.com/ercanaziz/notiFY/Sema-Durgut/db"
	history "github.com/ercanaziz/notiFY/Sema-Durgut"
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
