package notification

import (
	"fmt"
	"net/smtp"
)

// EmailService e-posta gönderme işlemlerini yönetecek yapı
type EmailService struct {
	SMTPHost string // Örn: smtp.gmail.com
	SMTPPort string // Örn: 587
	Sender   string // Senin mail adresin (dursoydogukan@gmail.com vb.)
	Password string // Mail şifren (veya Gmail Uygulama Şifresi)
}

// NewEmailService yeni bir e-posta servisi oluşturur
func NewEmailService(host, port, sender, password string) *EmailService {
	return &EmailService{
		SMTPHost: host,
		SMTPPort: port,
		Sender:   sender,
		Password: password,
	}
}

// SendEmail hedef kişiye, belirtilen başlık ve içerikle mail atar
func (e *EmailService) SendPriceAlertEmail(toEmail string, productName string, targetPrice float64) error {
	// 1. Mail sunucusuna giriş yap (Auth)
	auth := smtp.PlainAuth("", e.Sender, e.Password, e.SMTPHost)

	// 2. Mailin konusunu ve içeriğini oluştur
	subject := "🔔 Fiyat Alarmı: " + productName
	body := fmt.Sprintf("Merhaba,\n\nTakip ettiğin '%s' ürününün fiyatı hedeflenen %v TL seviyesine ulaştı!\n\nHemen kontrol et.", productName, targetPrice)

	// 3. Maili Go'nun anlayacağı formata (Byte dizisine) çevir
	msg := []byte("To: " + toEmail + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"Content-Type: text/plain; charset=UTF-8\r\n" +
		"\r\n" +
		body + "\r\n")

	// 4. Maili gönder
	addr := fmt.Sprintf("%s:%s", e.SMTPHost, e.SMTPPort)
	err := smtp.SendMail(addr, auth, e.Sender, []string{toEmail}, msg)
	if err != nil {
		return fmt.Errorf("mail gönderilirken hata oluştu: %v", err)
	}

	fmt.Println("✅ Başarılı: E-posta gönderildi ->", toEmail)
	return nil
}
