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
