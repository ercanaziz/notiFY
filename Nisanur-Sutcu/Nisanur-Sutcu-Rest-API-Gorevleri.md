# Nisanur Sütcü - REST API Metotları
**API Test Videosu:** Link buraya eklenecek

## 1. Ürün Arama
- **Endpoint:** `GET /products/search`
- **Query Parameters:** - `q` (string, required) - Aranacak ürün adı
- **Authentication:** Bearer Token gerekli
- **Açıklama:** Kullanıcının kendi takip listesi içinde anahtar kelime ile regex tabanlı arama yapmasını sağlar.
- **Response:** `200 OK` - Eşleşen ürünlerin listesi başarıyla getirildi

## 2. Takip Listesine Ekleme
- **Endpoint:** `POST /watchlist/add`
- **Request Body:** ```json
  {
    "product_name": "iPhone 15 Pro",
    "brand": "Apple",
    "current_price": 75000.0,
    "category": "Elektronik",
    "product_url": "[https://example.com/urun](https://example.com/urun)"
  }
- **Authentication:** Bearer Token gerekli
- **Açıklama:** Yeni bir ürünü, isteği atan kullanıcının user_id bilgisiyle MongoDB'ye kaydeder.
- **Response:** 201 Created - Ürün başarıyla takip listesine eklendi

## 3. Takip Listesini Listeleme
- **Endpoint:** GET /watchlist
- **Authentication:** Bearer Token gerekli
- **Açıklama:** Sadece giriş yapmış kullanıcıya ait olan tüm takip listesini döndürür.
- **Response:** 200 OK - Kullanıcının takip listesi başarıyla getirildi

## 4. Takip Listesinden Çıkarma
- **Endpoint:** DELETE /watchlist/{id}
- **Path Parameters:** - id (string, required) - Silinecek ürünün benzersiz ID'si
- **Authentication:** Bearer Token gerekli
- **Açıklama:** Kullanıcının artık takip etmek istemediği bir ürünü veritabanından siler.
- **Response:** 200 OK - Ürün başarıyla listeden temizlendi

## 5. Kategori Listeleme
- **Endpoint:** GET /products/category
- **Authentication:** Bearer Token gerekli
- **Açıklama:** Kullanıcının takip listesindeki ürünlere ait benzersiz kategori isimlerini listeler.
- **Response:** 200 OK - Kategoriler başarıyla getirildi

## 6. Popüler Ürünleri Görüntüleme
- **Endpoint:** GET /products/trending
- **Açıklama:** Veritabanındaki tüm ürünler arasında watch_count değeri en yüksek olan ilk 10 ürünü büyükten küçüğe sıralayarak getirir.
- **Response:** 200 OK - En popüler ürünler başarıyla listelendi
