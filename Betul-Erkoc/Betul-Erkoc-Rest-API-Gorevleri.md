# Betül Erkoç — REST API Görevleri

## 1. Kayıt Olma
**API Metodu:** `POST /auth/register`

**Açıklama:** Kullanıcıların e-posta ve şifre ile sisteme yeni hesap tanımlamasını sağlar.

**Request Body:**
```json
{
  "email": "betul.erkoc@example.com",
  "password": "Suleyman123!",
  "firstName": "Betül",
  "lastName": "Erkoç"
}
```

**Yanıtlar:**
- `201` → Kullanıcı başarıyla oluşturuldu
- `400` → Geçersiz veri girişi

---

## 2. Giriş Yapma
**API Metodu:** `POST /auth/login`

**Açıklama:** Kayıtlı kullanıcıların kimlik doğrulamasını yaparak sistem erişim izni almasını sağlar.

**Request Body:**
```json
{
  "email": "betul.erkoc@example.com",
  "password": "Suleyman123!"
}
```

**Yanıtlar:**
- `200` → Giriş başarılı, JWT token döner
- `401` → E-posta veya şifre hatalı

---

## 3. Profil Güncelleme
**API Metodu:** `PUT /user/profile`

**Açıklama:** Kullanıcının ad, soyad ve iletişim tercihlerini düzenlemesine olanak tanır.

**Yetkilendirme:** `Bearer Token` gereklidir.

**Request Body:**
```json
{
  "firstName": "Betül",
  "lastName": "Erkoç",
  "communicationPreferences": {
    "newsletter": true,
    "smsNotifications": false
  }
}
```

**Yanıtlar:**
- `200` → Profil başarıyla güncellendi
- `401` → Yetkisiz erişim

---

## 4. Şifre Değiştirme
**API Metodu:** `PATCH /user/password`

**Açıklama:** Mevcut kullanıcının güvenliği için şifresini yenilemesini sağlar.

**Yetkilendirme:** `Bearer Token` gereklidir.

**Request Body:**
```json
{
  "oldPassword": "Suleyman123!",
  "newPassword": "YeniGucluSifre2026!"
}
```

**Yanıtlar:**
- `200` → Şifre başarıyla güncellendi
- `400` → Geçersiz veri
- `401` → Yetkisiz erişim

---

## 5. Hesap Silme
**API Metodu:** `DELETE /user/profile`

**Açıklama:** Kullanıcının tüm hesap verilerini ve tercihlerini sistemden kalıcı olarak kaldırır.

**Yetkilendirme:** `Bearer Token` gereklidir.

**Yanıtlar:**
- `204` → Kullanıcı hesabı kalıcı olarak silindi
- `401` → Yetkisiz erişim

---

## 6. Çıkış Yapma
**API Metodu:** `POST /auth/logout`

**Açıklama:** Aktif kullanıcı oturumunu sonlandırarak güvenli çıkış yapılmasını sağlar.

**Yetkilendirme:** `Bearer Token` gereklidir.

**Yanıtlar:**
- `200` → Başarıyla çıkış yapıldı
- `401` → Yetkisiz erişim
# Betül Erkoç — REST API Görevleri

## 1. Kayıt Olma
**API Metodu:** `POST /auth/register`

**Açıklama:** Kullanıcıların e-posta ve şifre ile sisteme yeni hesap tanımlamasını sağlar.

**Request Body:**
```json
{
  "email": "betul.erkoc@example.com",
  "password": "Suleyman123!",
  "firstName": "Betül",
  "lastName": "Erkoç"
}
```

**Yanıtlar:**
- `201` → Kullanıcı başarıyla oluşturuldu
- `400` → Geçersiz veri girişi

---

## 2. Giriş Yapma
**API Metodu:** `POST /auth/login`

**Açıklama:** Kayıtlı kullanıcıların kimlik doğrulamasını yaparak sistem erişim izni almasını sağlar.

**Request Body:**
```json
{
  "email": "betul.erkoc@example.com",
  "password": "Suleyman123!"
}
```

**Yanıtlar:**
- `200` → Giriş başarılı, JWT token döner
- `401` → E-posta veya şifre hatalı

---

## 3. Profil Güncelleme
**API Metodu:** `PUT /user/profile`

**Açıklama:** Kullanıcının ad, soyad ve iletişim tercihlerini düzenlemesine olanak tanır.

**Yetkilendirme:** `Bearer Token` gereklidir.

**Request Body:**
```json
{
  "firstName": "Betül",
  "lastName": "Erkoç",
  "communicationPreferences": {
    "newsletter": true,
    "smsNotifications": false
  }
}
```

**Yanıtlar:**
- `200` → Profil başarıyla güncellendi
- `401` → Yetkisiz erişim

---

## 4. Şifre Değiştirme
**API Metodu:** `PATCH /user/password`

**Açıklama:** Mevcut kullanıcının güvenliği için şifresini yenilemesini sağlar.

**Yetkilendirme:** `Bearer Token` gereklidir.

**Request Body:**
```json
{
  "oldPassword": "Suleyman123!",
  "newPassword": "YeniGucluSifre2026!"
}
```

**Yanıtlar:**
- `200` → Şifre başarıyla güncellendi
- `400` → Geçersiz veri
- `401` → Yetkisiz erişim

---

## 5. Hesap Silme
**API Metodu:** `DELETE /user/profile`

**Açıklama:** Kullanıcının tüm hesap verilerini ve tercihlerini sistemden kalıcı olarak kaldırır.

**Yetkilendirme:** `Bearer Token` gereklidir.

**Yanıtlar:**
- `204` → Kullanıcı hesabı kalıcı olarak silindi
- `401` → Yetkisiz erişim

---

## 6. Çıkış Yapma
**API Metodu:** `POST /auth/logout`

**Açıklama:** Aktif kullanıcı oturumunu sonlandırarak güvenli çıkış yapılmasını sağlar.

**Yetkilendirme:** `Bearer Token` gereklidir.

**Yanıtlar:**
- `200` → Başarıyla çıkış yapıldı
- `401` → Yetkisiz erişim
