package history

import (
	"os"

	"github.com/ercanaziz/notiFY/Sema-Durgut/analysis"
	"github.com/ercanaziz/notiFY/Sema-Durgut/db"
	"github.com/gin-gonic/gin"
)

func Start() {
	// 1. Veritabanına (MongoDB Atlas) bağlan
	db.Connect()
	// db.Migrate() sildik! MongoDB'de buna gerek yok.

	// 2. Gin Router (Sunucu) kurulumu
	r := gin.Default()

	// 3. API Uç Noktaları (Route'lar)
	products := r.Group("/products")
	{
		products.GET("/:id/history", analysis.GetPriceHistory)
		products.GET("/:id/chart-data", analysis.GetChartData)
		products.GET("/:id/lowest-price", analysis.GetLowestPrice)
		products.GET("/compare", analysis.CompareStores)
		products.GET("/:id/discount-rate", analysis.GetDiscountRate)
		products.GET("/:id/forecast", analysis.GetForecast)

		// İŞTE BURAYI DÜZELTTİK: "/:id/color-filter" yerine sadece "/filter" yaptık
		products.GET("/filter", analysis.FilterByColor)
	}

	port := os.Getenv("HISTORY_PORT")
	if port == "" {
		port = "8081"
	}
	r.Run(":" + port)
}

/* Küçük ama hayati bir ipucu: Ana projeye kodları geçirirken, kendi test
klasöründeki go.mod ve go.sum dosyalarını ana projeye kopyalama. Ortak
projenin zaten kendine ait bir go.mod dosyası olacaktır. Sen .go
dosyalarını ortak projeye attıktan sonra, ana proje dizininde terminali
açıp go mod tidy yazman yeterlidir. Bu sihirli komut, senin dosyalarında
kullanılan paketleri (gin veya pq gibi) algılar ve ana projeye otomatik
olarak dahil eder. */
