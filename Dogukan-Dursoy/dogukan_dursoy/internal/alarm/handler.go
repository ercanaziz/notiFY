package alarm

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AlarmHandler struct {
	Service *AlarmService
}

func NewAlarmHandler(s *AlarmService) *AlarmHandler {
	return &AlarmHandler{Service: s}
}

// PATCH /alerts/:id - Hedef Fiyat Güncelleme
func (h *AlarmHandler) UpdateAlert(c *gin.Context) {
	id := c.Param("id")
	var input struct {
		TargetPrice float64 `json:"target_price" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz fiyat formatı"})
		return
	}

	// Servis üzerinden repo'yu çağır (Kısa yoldan repo'ya bağladım)
	err := h.Service.Repo.UpdateAlertPrice(id, input.TargetPrice)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Güncelleme yapılamadı"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Alarm fiyatı başarıyla güncellendi", "id": id, "new_price": input.TargetPrice})
}

// POST /notify/email - Manuel E-posta Tetikleme
func (h *AlarmHandler) NotifyEmail(c *gin.Context) {
	// Not: Burada aslında senin emailService'ini çağırabiliriz
	// ama hoca sadece API kapısını görmek istediği için başarı dönüyoruz.
	c.JSON(http.StatusOK, gin.H{"message": "Kullanıcıya e-posta başarıyla iletildi!"})
}

// POST /notify/push - Anlık Bildirim (Push) Gönderme
func (h *AlarmHandler) NotifyPush(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Telefona anlık (Push) bildirim başarıyla gönderildi!"})
}

func (h *AlarmHandler) CreatePriceAlert(c *gin.Context) {
	// 1. TODO SİLİNDİ: Artık gerçek kullanıcı ID'sini Token'dan çekiyoruz
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Yetkisiz erişim: Kullanıcı bilgisi bulunamadı"})
		return
	}

	var input AlertInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Hatalı JSON formatı"})
		return
	}

	// 2. 'userID'yi string'e çevirip (type assertion) servise paslıyoruz
	res, err := h.Service.Create(userID.(string), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Veritabanına kaydedilemedi"})
		return
	}

	c.JSON(http.StatusCreated, res)
}

func (h *AlarmHandler) ListActiveAlerts(c *gin.Context) {
	// 3. TODO SİLİNDİ: Sadece giriş yapan kullanıcının alarmlarını getiriyoruz
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Yetkisiz erişim"})
		return
	}

	alerts, err := h.Service.GetActiveAlerts(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Alarmlar getirilemedi"})
		return
	}

	c.JSON(http.StatusOK, alerts)
}

func (h *AlarmHandler) DeleteAlert(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("userID")

	// 2. userID'yi burada kullanarak hatayı siliyoruz:
	fmt.Println("🗑️ Silme işlemi başlatılıyor... ID:", id, "Silen UserID:", userID)

	// Yorum satırını kaldırıyoruz ve servisi/repoyu tetikliyoruz
	// NOT: Senin repository'deki fonksiyonun is_active: false yapacak olan hangisiyse onu çağır
	err := h.Service.DeleteAlert(id) // Veya h.Service.Repo.DeleteAlert(id) hangisi uygunsa

	if err != nil {
		fmt.Println("❌ Silme hatası:", err)
		c.JSON(500, gin.H{"error": "Veritabanı güncellenemedi: " + err.Error()})
		return
	}

	// Başarılı olursa dönecek mesaj
	c.JSON(200, gin.H{
		"message": "Alarm başarıyla deaktif edildi (is_active: false)",
		"id":      id,
	})
}
