package analysis

import (
	"github.com/gin-gonic/gin"
	"math"
	"net/http"
	"notiFY/models" // Modellerin bulunduğu paket eklendi
)

// GET /products/:id/history
func GetPriceHistory(c *gin.Context) {
	id := c.Param("id")
	history, err := GetHistoryFromDB(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Veri çekilemedi"})
		return
	}
	c.JSON(http.StatusOK, history)
}

// GET /products/:id/lowest-price
func GetLowestPrice(c *gin.Context) {
	id := c.Param("id")
	lowest, _ := GetLowestPriceFromDB(id)

	// Gerçek senaryoda currentPrice'ı products tablosundan çekersin
	c.JSON(http.StatusOK, models.PriceInfo{
		LowestPrice:  lowest,
		CurrentPrice: lowest + 50, // Örnek veri
	})
}

// GET /products/compare?query=iphone
func CompareStores(c *gin.Context) {
	query := c.Query("query")
	offers, _ := GetStoreComparisons(query)

	c.JSON(http.StatusOK, models.ComparisonList{
		ProductName: query,
		Offers:      offers,
	})
}

// GET /products/:id/discount-rate
func GetDiscountRate(c *gin.Context) {
	id := c.Param("id")
	avgPrice, _ := GetAveragePriceFromDB(id)
	currentPrice := 900.0 // Örnek: Normalde db.DB.QueryRow("SELECT current_price...") ile çekilir.

	// İndirim Yüzdesi Hesaplama Mantığı
	var discount int
	rating := "NEUTRAL"
	if currentPrice < avgPrice {
		discount = int(((avgPrice - currentPrice) / avgPrice) * 100)
		if discount > 15 {
			rating = "VERY_GOOD"
		} else {
			rating = "GOOD"
		}
	}

	c.JSON(http.StatusOK, models.DiscountAnalysis{
		CurrentPrice:       currentPrice,
		AveragePrice:       math.Round(avgPrice*100) / 100,
		DiscountPercentage: discount,
		Rating:             rating,
	})
}

// GET /products/:id/chart-data
func GetChartData(c *gin.Context) {
	id := c.Param("id")

	// 1. Repository'den veriyi iste (Hata yönetimiyle beraber)
	// Not: GetChartDataFromDB fonksiyonu içindeki 'B Planı' sayesinde veri yoksa güncel fiyatı dönecek
	chartData, err := GetChartDataFromDB(id)
	if err != nil {
		// Eğer veritabanı bağlantısı veya ID hatası varsa 500 döner
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Grafik verisi hazirlanamadi: " + err.Error()})
		return
	}

	// 2. Eğer veri bossa (hiçbir kayıt bulunamadıysa)
	if len(chartData.Labels) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"message":  "Bu urun icin henuz fiyat verisi bulunmamaktadir.",
			"labels":   []string{},
			"datasets": []models.Dataset{},
		})
		return
	}

	// 3. Her şey yolundaysa 200 OK ile veriyi bas
	c.JSON(http.StatusOK, chartData)
}

// GET /products/:id/forecast
func GetForecast(c *gin.Context) {
	id := c.Param("id")

	pred, conf, estPrice := GetPriceForecast(id)

	c.JSON(200, gin.H{
		"productId":               id,
		"prediction":              pred,
		"confidenceScore":         conf,
		"estimatedPriceNextMonth": estPrice,
		"aiMethod":                "Mean Reversion Algorithm", // Havalı bir isim ekledik :)
	})
}

// get /products/filter?name=red
func FilterByColor(c *gin.Context) {
	color := c.Query("name") // ?name=red gibi
	products, _ := GetProductsByColor(color)
	c.JSON(200, products)
}
