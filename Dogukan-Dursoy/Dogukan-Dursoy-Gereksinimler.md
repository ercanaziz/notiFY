1. **Fiyat Alarmı Oluşturma**

    **API Metodu**:`POST /alerts`

    **Açıklama**: Belirli bir ürün fiyatı düştüğünde tetiklenecek hedef fiyat alarmı kurar.

2. **Aktif Alarmları Listeleme**

    **API Metodu**: `GET /alerts/active`

    **Açıklama**: Kullanıcının henüz tetiklenmemiş aktif tüm alarmlarını gösterir.

3. **Alarm Durumu Güncelleme**

    **API Metodu**: `PATCH /alerts/{id}`

    **Açıklama**: Kurulu olan alarmın hedef fiyatını veya bildirim tipini değiştirir.

4. **Alarm Silme**

    **API Metodu**:`DELETE /alerts/{id}`

    **Açıklama**: İptal edilmek istenen fiyat alarmını sistemden kaldırır.

5. **E-posta Bildirimi Gönderme**

    **API Metodu**: `POST /notify/email`

    **Açıklama**: Fiyat düşüşü anında kullanıcıya otomatik bilgilendirme e-postası iletir.

6.  **Anlık Bildirim Gönderme**

    **API Metodu**: `POST /notify/push`

    **Açıklama**: Uygulama açıkken kullanıcıya anlık pop-up uyarısı gönderilmesini sağlar.