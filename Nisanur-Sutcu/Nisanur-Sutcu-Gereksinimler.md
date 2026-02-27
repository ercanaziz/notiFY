1. **Ürün Arama**

    **API Metodu**:`GET /products/search`

    **Açıklama**: Veritabanındaki ürünler arasında anahtar kelime ile arama yapılmasını sağlar.

2. **Takip Listesine Ekleme**

    **API Metodu**: `POST /watchlist/add`

    **Açıklama**: Bir ürünün kullanıcının kişisel takip listesine kaydedilmesini sağlar.

3. **Takip Listesini Listeleme**

    **API Metodu**: `GET /watchlist`

    **Açıklama**: Kullanıcının takip ettiği tüm ürünleri tek bir ekranda listeler.

4. **Takip Listesinden Çıkarma**

    **API Metodu**:`DELETE /watchlist/{id}`

    **Açıklama**: Artık takip edilmek istenmeyen bir ürünü listeden temizler.

5. **Kategori Listeleme**

    **API Metodu**: `GET /categories`

    **Açıklama**: Ürünlerin gruplandığı (Teknoloji, Kitap vb.) kategorileri görüntüler.

6.  **Popüler Ürünleri Görüntüleme**

    **API Metodu**: `GET /products/trending`

    **Açıklama**: En çok takip edilen ürünlerin hızlıca ana sayfada listelenmesini sağlar.