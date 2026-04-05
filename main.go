package main

import (
	"os"

	login "github.com/ercanaziz/notiFY/Betul-Erkoc/Backend"
	alert "github.com/ercanaziz/notiFY/Dogukan-Dursoy/dogukan_dursoy/cmd/server"
	product "github.com/ercanaziz/notiFY/Nisanur-Sutcu/backend"
	productDB "github.com/ercanaziz/notiFY/Nisanur-Sutcu/backend/db"

	backendDB "github.com/ercanaziz/notiFY/Ercan-Aziz/Backend/db"
	backendRouter "github.com/ercanaziz/notiFY/Ercan-Aziz/Backend/router"
	history "github.com/ercanaziz/notiFY/Sema-Durgut"
	historyDB "github.com/ercanaziz/notiFY/Sema-Durgut/db"
	"github.com/gin-gonic/gin"
)

func main() {

	backendDB.Connect()
	historyDB.Connect()
	productDB.Connect()
	go login.Start()
	go alert.Start()

	r := gin.Default()
	backendRouter.RegisterRoutes(r)
	history.RegisterRoutes(r)
	product.RegisterRoutes(r)
	login.RegisterRoutes(r)
	alert.RegisterRoutes(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
