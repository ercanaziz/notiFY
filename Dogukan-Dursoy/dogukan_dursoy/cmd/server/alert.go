package alert

import (
	"context"
	"fmt"
	"log"
	"time"

	// Kendi paketlerin (Yeni modül ismine göre güncellendi)
	"github.com/ercanaziz/notiFY/Dogukan-Dursoy/dogukan_dursoy/auth"
	"github.com/ercanaziz/notiFY/Dogukan-Dursoy/dogukan_dursoy/internal/alarm"
	"github.com/ercanaziz/notiFY/Dogukan-Dursoy/dogukan_dursoy/internal/notification"

	// Dış kütüphaneler
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

func Start() {
	uri := "mongodb+srv://teamnotify:notiFY@test.ek07wik.mongodb.net/?appName=test"

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal("❌ MongoDB'ye bağlanılamadı: ", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("❌ MongoDB'ye ping atılamadı: ", err)
	}

	Client = client

	fmt.Println("✅ MongoDB bağlantısı başarılı!")

	dbAlarms := Client.Database("price_tracker_db")
	dbUsers := Client.Database("notiFY_DB")

	repo := alarm.NewAlarmRepository(dbAlarms.Collection("alerts"))

	emailSvc := notification.NewEmailService(
		"smtp.gmail.com",
		"587",
		"notifycmp@gmail.com",
		"vctxxnhuknyvnjno",
	)

	go func() {
		fmt.Println("🤖 Akıllı Fiyat Takip Motoru başlatıldı...")
		for {
			activeAlerts, err := repo.GetActiveAlerts()
			if err != nil {
				fmt.Println("❌ Alarmlar çekilemedi:", err)
			} else {
				for _, alert := range activeAlerts {
					currentMarketPrice, err := repo.GetProductPrice(alert.ProductID)
					if err != nil {
						continue
					}

					if currentMarketPrice <= alert.TargetPrice {
						fmt.Printf("🔥 %s tetiklendi! Kullanıcı aranıyor: %s\n", alert.ProductID, alert.UserID)

						var foundUser bson.M
						userObjID, _ := primitive.ObjectIDFromHex(alert.UserID)

						userCtx, userCancel := context.WithTimeout(context.Background(), 5*time.Second)
						err := dbUsers.Collection("users").FindOne(userCtx, bson.M{"_id": userObjID}).Decode(&foundUser)
						userCancel()

						targetEmail := "dursoydogukan@gmail.com"
						if err == nil && foundUser["email"] != nil {
							targetEmail = foundUser["email"].(string)
						}

						fmt.Printf("📧 Mail gönderiliyor: %s -> %s\n", alert.ProductID, targetEmail)

						err = emailSvc.SendPriceAlertEmail(targetEmail, alert.ProductID, alert.TargetPrice)
						if err == nil {
							repo.DeactivateAlert(alert.ID)
							fmt.Println("✅ Bildirim gönderildi ve alarm susturuldu.")
						}
					}
				}
			}
			time.Sleep(1 * time.Minute)
		}
	}()
}
func RegisterRoutes(r *gin.Engine) {
	dbAlarms := Client.Database("price_tracker_db")
	dbUsers := Client.Database("notiFY_DB")

	repo := alarm.NewAlarmRepository(dbAlarms.Collection("alerts"))
	svc := alarm.NewAlarmService(repo)
	hdl := alarm.NewAlarmHandler(svc)
	authHdl := auth.NewAuthHandler(dbUsers.Collection("users"))

	r.POST("/register", authHdl.Register)
	r.POST("/login", authHdl.Login)

	protected := r.Group("/alerts")
	protected.Use(auth.AuthMiddleware())
	{
		protected.POST("", hdl.CreatePriceAlert)
		protected.GET("/active", hdl.ListActiveAlerts)
		protected.DELETE("/:id", hdl.DeleteAlert)
		protected.PATCH("/:id", hdl.UpdateAlert)
	}

	r.POST("/notify/email", hdl.NotifyEmail)
	r.POST("/notify/push", hdl.NotifyPush)
}
