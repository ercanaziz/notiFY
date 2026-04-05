package login

import (
	"fmt"
	"log"
	"github.com/gin-gonic/gin"
    
    "github.com/ercanaziz/notiFY/Betul-Erkoc/Backend/router"
    "github.com/ercanaziz/notiFY/Betul-Erkoc/Backend/db"
    
)

func Start() {
	db.InitDB()

	r := router.SetupRouter()

	fmt.Println("🚀 notiFY Sunucusu 8080 portunda hazır...")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Sunucu başlatılamadı: ", err)
	}
}


func RegisterRoutes(r *gin.Engine) {
    router.SetupRouter(r)  // eğer SetupRouter *gin.Engine alacak şekilde değiştirirsen
}