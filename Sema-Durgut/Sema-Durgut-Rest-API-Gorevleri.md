# Sema Durgut - Rest Api Metotları
**API Test Videosu:** Link buraya eklenecek

## 1. Fiyat Geçmişini Getirme
- **Endpoint:** `GET /analysis/history/{id}`
- **Path Parameters:** - `id` (string, required) - Ürünün benzersiz ID'si
- **Açıklama:** Belirli bir ürünün veritabanındaki tüm geçmiş fiyat hareketlerini kronolojik bir liste olarak sunar.
- **Response:** `200 OK` - Ürün fiyat geçmişi başarıyla getirildi

## 2. Fiyat Grafiği Verisi Sağlama
- **Endpoint:** `GET /analysis/chart/{id}`
- **Path Parameters:** - `id` (string, required) - Ürünün benzersiz ID'si
- **Açıklama:** Frontend grafik kütüphaneleri (Chart.js vb.) için optimize edilmiş, zaman serisi tabanlı fiyat verilerini hazırlar.
- **Response:** `200 OK` - Grafik verileri başarıyla oluşturuldu

## 3. En Düşük Fiyat Sorgulama
- **Endpoint:** `GET /analysis/lowest/{id}`
- **Path Parameters:** - `id` (string, required) - Ürünün benzersiz ID'si
- **Açıklama:** Ürünün sisteme kaydedildiği andan itibaren gördüğü tüm zamanların en düşük fiyat bilgisini döner.
- **Response:** `200 OK` - En düşük fiyat bilgisi başarıyla getirildi

## 4. Mağaza Karşılaştırması Yapma
- **Endpoint:** `GET /analysis/compare`
- **Query Parameters:** - `query` (string, required) - Karşılaştırılacak ürün adı
- **Açıklama:** Aynı ürünün farklı mağazalardaki güncel fiyatlarını yan yana getirerek karşılaştırma imkanı sunar.
- **Response:** `200 OK` - Mağaza fiyatları başarıyla karşılaştırıldı

## 5. İndirim Oranı Hesaplama
- **Endpoint:** `GET /analysis/discount/{id}`
- **Path Parameters:** - `id` (string, required) - Ürünün benzersiz ID'si
- **Açıklama:** Ürünün mevcut fiyatını piyasa ortalamasıyla kıyaslar, indirim yüzdesini ve fırsat puanını (Rating) hesaplar.
- **Response:** `200 OK` - İndirim analizi başarıyla tamamlandı

## 6. Renge Göre Filtreleme
- **Endpoint:** `GET /analysis/filter`
- **Query Parameters:** - `name` (string, required) - Aranacak renk (örn: "Siyah")
- **Açıklama:** Kullanıcının takip listesindeki ürünleri belirli renk seçeneklerine göre süzerek listeler.
- **Response:** `200 OK` - Renk filtresine uyan ürünler başarıyla listelendi

## 7. Fiyat Tahmini İzleme (Yapay Zeka)
- **Endpoint:** `GET /analysis/forecast/{id}`
- **Path Parameters:** - `id` (string, required) - Ürünün benzersiz ID'si
- **Açıklama:** Geçmiş fiyat trendlerini analiz ederek ürünün gelecekteki olası fiyat yönünü (Artış/Azalış) ve güven skorunu raporlar.
- **Response:** `200 OK` - Fiyat tahmini ve analiz raporu başarıyla oluşturuldu
