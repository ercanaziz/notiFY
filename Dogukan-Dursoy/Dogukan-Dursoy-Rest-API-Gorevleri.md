
https://youtu.be/qsT4Aqke4a4

Doğukan Dursoy - REST API Metotları (Fiyat Takip Sistemi)
1. Fiyat Alarmı Oluşturma
Endpoint: POST /alerts/

Authentication: Bearer Token gerekli

Açıklama: Kullanıcının takip ettiği belirli bir ürün için hedef fiyat belirlemesini sağlar. Fiyat bu seviyeye düştüğünde sistem otomatik tetiklenir.

Response: 201 Created - Fiyat alarmı başarıyla kuruldu.

2. Aktif Alarmları Listeleme
Endpoint: GET /alerts/active

Authentication: Bearer Token gerekli

Açıklama: Sadece giriş yapmış kullanıcıya ait olan ve henüz fiyat hedefine ulaşmamış (tetiklenmemiş) tüm aktif alarmları döndürür.

Response: 200 OK - Aktif alarmlar listesi başarıyla getirildi.

3. Alarm Bilgilerini Güncelleme
Endpoint: PATCH /alerts/{id}

Path Parameters: id (string, required) - Güncellenecek alarmın benzersiz ID'si

Authentication: Bearer Token gerekli

Açıklama: Mevcut bir alarmın hedef fiyatını veya bildirim tercihlerini (email/push) değiştirmek için kullanılır.

Response: 200 OK - Alarm başarıyla güncellendi.

4. Alarm Silme
Endpoint: DELETE /alerts/{id}

Path Parameters: id (string, required) - Silinecek alarmın benzersiz ID'si

Authentication: Bearer Token gerekli

Açıklama: Kullanıcının artık takip etmek istemediği veya iptal etmek istediği fiyat alarmını sistemden kalıcı olarak kaldırır.

Response: 200 OK - Alarm başarıyla silindi.

5. E-posta Bildirimi Gönderme
Endpoint: POST /notify/email

Açıklama: Akıllı takip motoru fiyat düşüşü tespit ettiğinde, SMTP protokolü üzerinden kullanıcıya otomatik bilgilendirme e-postası iletir.

Response: 200 OK - Bilgilendirme e-postası başarıyla gönderildi.

6. Anlık (Push) Bildirim Gönderme
Endpoint: POST /notify/push

Açıklama: Fiyat alarmı tetiklendiğinde kullanıcıya uygulama üzerinden anlık pop-up uyarısı gönderilmesini tetikler.

Response: 200 OK - Anlık bildirim başarıyla iletildi.
