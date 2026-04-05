package main

import (
	"github.com/gin-gonic/gin"
	"notiFY/Sema-Durgut/analysis"
	"notiFY/Sema-Durgut/db"
)

func main() {
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

	// 4. Sunucuyu Başlat
	r.Run(":8080") // localhost:8080'de dinliyor
}

/* Küçük ama hayati bir ipucu: Ana projeye kodları geçirirken, kendi test
klasöründeki go.mod ve go.sum dosyalarını ana projeye kopyalama. Ortak
projenin zaten kendine ait bir go.mod dosyası olacaktır. Sen .go
dosyalarını ortak projeye attıktan sonra, ana proje dizininde terminali
açıp go mod tidy yazman yeterlidir. Bu sihirli komut, senin dosyalarında
kullanılan paketleri (gin veya pq gibi) algılar ve ana projeye otomatik
olarak dahil eder. */
