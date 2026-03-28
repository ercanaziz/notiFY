package main

import (
	"fmt"
	"net/http"
	"time"
)

// User - Döküman
type User struct {
	ID              string    `json:"_id"`
	Email           string    `json:"email"`
	FirstName       string    `json:"firstName"`
	LastName        string    `json:"lastName"`
	IsEmailVerified bool      `json:"isEmailVerified"`
	CreatedOn       time.Time `json:"createdOn"`
}

// RegisterInput - Yeni Kullanıcı Kaydı için beklenen veri
type RegisterInput struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

// LoginInput - Kullanıcı Girişi için beklenen veri
type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// ProfileUpdateInput - Profil Güncelleme için beklenen veri
type ProfileUpdateInput struct {
	FirstName                string                   `json:"firstName"`
	LastName                 string                   `json:"lastName"`
	CommunicationPreferences CommunicationPreferences `json:"communicationPreferences"`
}

type CommunicationPreferences struct {
	Newsletter       bool `json:"newsletter"`
	SmsNotifications bool `json:"smsNotifications"`
}

// Handlers
func registerUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	fmt.Fprintln(w, "201: Kullanıcı başarıyla oluşturuldu")
}

func loginUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	fmt.Fprintln(w, "200: Giriş başarılı, token üretildi")
}

func profileHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPut:
		fmt.Fprintln(w, "200: Profil başarıyla güncellendi")
	case http.MethodDelete:
		fmt.Fprintln(w, "204: Kullanıcı hesabı kalıcı olarak silindi")
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func changePassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	fmt.Fprintln(w, "200: Şifre başarıyla güncellendi")
}

func logoutUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "200: Başarıyla çıkış yapıldı")
}

// yönlendirme ve sunucu başlatma
func main() {
	// API Rotaları 
	http.HandleFunc("/api/auth/register", registerUser)
	http.HandleFunc("/api/auth/login", loginUser)
	http.HandleFunc("/api/users/profile", profileHandler)
	http.HandleFunc("/api/users/password", changePassword)
	http.HandleFunc("/api/auth/logout", logoutUser)

	fmt.Println("notiFY API Sunucusu 8080 portunda çalışıyor...")
	http.ListenAndServe(":8080", nil)
}