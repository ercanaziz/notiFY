## 1. Kayıt Olma 
**API Metodu:** `POST /auth/register` 

**Açıklama:** Kullanıcıların e-posta ve şifre ile sisteme yeni hesap tanımlamasını sağlar.

## 2. Giriş Yapma 
**API Metodu:** `POST /auth/login` 

**Açıklama:** Kayıtlı kullanıcıların kimlik doğrulamasını yaparak sistem erişim izni almasını sağlar.

## 3. Profil Güncelleme 
**API Metodu:** `PUT /user/profile`  

**Açıklama:** Kullanıcının ad, soyad ve iletişim tercihlerini düzenlemesine olanak tanır.

## 4. Şifre Değiştirme 
**API Metodu:** `PATCH /user/password`  

**Açıklama:** Mevcut kullanıcının güvenliği için şifresini yenilemesini sağlar.

## 5. Hesap Silme 
**API Metodu:** `DELETE /user/account`  

**Açıklama:** Kullanıcının tüm hesap verilerini ve tercihlerini sistemden kalıcı olarak kaldırır.

## 6. Çıkış Yapma 
**API Metodu:** `POST /auth/logout`  

**Açıklama:** Aktif kullanıcı oturumunu sonlandırarak güvenli çıkış yapılmasını sağlar.

