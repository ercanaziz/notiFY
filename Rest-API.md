**REST API Adresi:** https://notify-2hde.onrender.com

# REST API Görev Dağılımı
Nisanur Sütcü
- **Backend Mimari Tasarımı:** Monolitik (tek dosya) yapıdaki Go kodunu; handlers, models ve db katmanlarına ayırarak Modüler Mimari (MVC yapısına uygun) geçişini sağladım.

- **Veritabanı Entegrasyonu:** MongoDB sürücüsü kullanılarak veritabanı bağlantı katmanının oluşturulması ve koleksiyon yönetimini gerçekleştirdim.

- **CRUD Operasyonları:** Kullanıcıya özel takip listesi yönetimi için "Ekleme, Listeleme, Silme ve Arama" fonksiyonlarını geliştirdim.

- **Gelişmiş Sorgulama Logic'leri:** * Regex Search: Ürün isimleri üzerinden büyük/küçük harf duyarsız, esnek arama algoritmasını kurguladım.

- **Aggregation & Sorting:** Popüler ürünlerin tespiti için watch_count değerine göre büyükten küçüğe sıralama (Sort) ve Limit operasyonlarını yönettim.

- **Distinct Query:** Takip listesindeki verilerden benzersiz kategori isimlerini çeken Distinct sorgusunu yazdım.

- **Güvenlik ve Yetkilendirme:** API rotalarını korumak için JWT (JSON Web Token) tabanlı AuthMiddleware entegrasyonunu yaptım ve tüm isteklere user_id bazlı filtreleme ekledim.

- **Hata Yönetimi (Error Handling):** Geçersiz ID formatları, bulunamayan kayıtlar (404) ve veritabanı hataları için özelleştirilmiş JSON hata mesajları ve HTTP statü kodlarını yapılandırdım.

- # Betül Erkoç

**Gereksinim Analizi ve Kullanıcı Hikayelerinin Oluşturulması:** Sisteme kayıt olma, giriş yapma, profil güncelleme, şifre değiştirme, hesap silme ve çıkış yapma olmak üzere 6 temel kullanıcı gereksinimini belirleyerek her biri için API metodunu, açıklamasını ve kapsamını tanımladım.

**API Endpoint Tasarımı:** RESTful mimari prensiplerine uygun olarak `POST /auth/register`, `POST /auth/login`, `PUT /user/profile`, `PATCH /user/password`, `DELETE /user/profile` ve `POST /auth/logout` endpoint'lerini tasarlayarak her birinin HTTP metodunu, yolunu ve işlevsel sorumluluğunu belirledim.

**OpenAPI (YAML) Şema Dokümantasyonu:** Tüm endpoint'lerin istek gövdelerini (request body), başarılı ve hatalı yanıt kodlarını (200, 201, 204, 400, 401) ve operasyon kimliklerini (operationId) kapsayan eksiksiz bir OpenAPI spesifikasyonu oluşturdum.

**Veri Modelleme ve Şema Tasarımı:** `RegisterInput`, `LoginInput`, `ProfileUpdateInput`, `PasswordUpdateInput`, `User` ve `Error` olmak üzere 6 adet bileşen şeması (component schema) tasarladım; her alan için veri tipi, format, zorunluluk durumu, açıklama ve örnek değerleri tanımladım.

**Güvenlik Şeması Tanımı:** JWT tabanlı kimlik doğrulama için `BearerAuth` güvenlik şemasını `components/securitySchemes` altında tanımladım; kimlik doğrulama gerektirmeyen endpoint'leri (`security: []`) ile korumalı endpoint'leri ayrıştırarak yetkilendirme katmanını yapılandırdım.

**İletişim Tercihleri Veri Yapısı:** Profil güncelleme şeması kapsamında kullanıcıya ait `newsletter` ve `smsNotifications` alanlarını barındıran iç içe `communicationPreferences` nesnesini modelledim.

**Hata Yanıtı Standardizasyonu:** Tüm hatalı isteklere tutarlı bir yapı sunmak amacıyla `message` ve `code` alanlarından oluşan merkezi `Error` şemasını tasarlayarak ilgili endpoint'lerde `$ref` ile referans gösterdim.



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


Nisanur Sütcü - REST API Metotları

## 1. Ürün Arama
- **Endpoint:** `GET /products/search`
- **Query Parameters:** - `q` (string, required) - Aranacak ürün adı
- **Authentication:** Bearer Token gerekli
- **Açıklama:** Kullanıcının kendi takip listesi içinde anahtar kelime ile regex tabanlı arama yapmasını sağlar.
- **Response:** `200 OK` - Eşleşen ürünlerin listesi başarıyla getirildi

## 2. Takip Listesine Ekleme
- **Endpoint:** `POST /watchlist/add`
- **Request Body:**
   ```json
  {
    "product_name": "iPhone 15 Pro",
    "brand": "Apple",
    "current_price": 75000.0,
    "category": "Elektronik",
    "product_url": "[https://example.com/urun](https://example.com/urun)"
  }
   ```
- **Authentication:** Bearer Token gerekli
- **Açıklama:** Yeni bir ürünü, isteği atan kullanıcının user_id bilgisiyle MongoDB'ye kaydeder.
- **Response:** `201 Created` - Ürün başarıyla takip listesine eklendi

## 3. Takip Listesini Listeleme
- **Endpoint:** `GET /watchlist`
- **Authentication:** Bearer Token gerekli
- **Açıklama:** Sadece giriş yapmış kullanıcıya ait olan tüm takip listesini döndürür.
- **Response:** `200 OK` - Kullanıcının takip listesi başarıyla getirildi

## 4. Takip Listesinden Çıkarma
- **Endpoint:** `DELETE /watchlist/{id}`
- **Path Parameters:** - `id` (string, required) - Silinecek ürünün benzersiz ID'si
- **Authentication:** Bearer Token gerekli
- **Açıklama:** Kullanıcının artık takip etmek istemediği bir ürünü veritabanından siler.
- **Response:** `200 OK` - Ürün başarıyla listeden temizlendi

## 5. Kategori Listeleme
- **Endpoint:** `GET /products/category`
- **Authentication:** Bearer Token gerekli
- **Açıklama:** Kullanıcının takip listesindeki ürünlere ait benzersiz kategori isimlerini listeler.
- **Response:** `200 OK` - Kategoriler başarıyla getirildi

## 6. Popüler Ürünleri Görüntüleme
- **Endpoint:** `GET /products/trending`
- **Açıklama:** Veritabanındaki tüm ürünler arasında watch_count değeri en yüksek olan ilk 10 ürünü büyükten küçüğe sıralayarak getirir.
- **Response:** `200 OK` - En popüler ürünler başarıyla listelendi


## Grup Üyelerinin REST API Metotları

1. [Ercan Aziz'in REST API Metotları](./Ercan-Aziz/Ercan-Aziz-Rest-API-Gorevleri.md)
2. [Sema Durgut'un REST API Metotları](./Sema-Durgut/Sema-Durgut-Rest-API-Gorevleri.md)
3. [Dogukan Dursoy'un REST API Metotları](./Dogukan-Dursoy/Dogukan-Dursoy-Rest-API-Gorevleri.md)
4. [Betül Erkoç'un REST API Metotları](./Betul-Erkoc/Betul-Erkoc-Rest-API-Gorevleri.md)
5. [Nisanur Sütçü'nün REST API Metotları](./Nisanur-Sutcu/Nisanur-Sutcu-Rest-API-Gorevleri.md)
