# API Tasarımı - OpenAPI Specification Örneği

**OpenAPI Spesifikasyon Dosyası:** [lamine.yaml](lamine.yaml)

Bu doküman, HERMES Ekibi tarafından geliştirilen "Fiyat Takip ve Akıllı Alarm Sistemi" projesi için OpenAPI Specification (OAS) 3.0 standardına göre hazırlanmış örnek bir API tasarımını içermektedir.

## OpenAPI Specification

```yaml
openapi: 3.0.3
info:
  title: Fiyat Takip ve Akıllı Bildirim API'si
  version: 1.0.0
  description: >
    Tüketicilere ve e-ticaret profesyonellerine, piyasadaki anlık fiyat değişimlerini 
    kaçırmadan en doğru zamanda aksiyon alma imkanı sunar. Redis ve Kafka/RabbitMQ 
    tabanlı asenkron mimari üzerinden fiyat dalgalanmalarını analiz eder ve kullanıcıların 
    belirlediği limitlerin altına inen fiyatları "satın alma sinyali" olarak iletir.
  contact:
    name: HERMES

servers:
  - url: https://api.fiyattakip.com
    description: Production Sunucusu
  - url: https://staging.fiyattakip.com
    description: Test Sunucusu (Staging)
  - url: http://localhost:8080
    description: Local Geliştirme (Go/Fiber veya Gin)

tags:
  - name: Auth & Users
    description: Kullanıcı doğrulama ve profil işlemleri (Betül Erkoç)
  - name: Products & Watchlist
    description: Ürün arama, filtreleme ve takip listesi yönetimi (Nisanur Sütçü)
  - name: Alerts & Notifications
    description: Fiyat alarmı oluşturma, yönetimi ve anlık bildirimler (Doğukan Dursoy)
  - name: Price Analytics
    description: Fiyat geçmişi, grafikleri ve karşılaştırmalar (Sema Durgut)
  - name: Core Lists & Plans
    description: Filtreleme, sıralama, abonelik ve geri bildirim (Ercan Aziz)

security:
  - BearerAuth: []

paths:
  # ==========================================
  # AUTH & USERS (Betül Erkoç)
  # ==========================================
  /auth/register:
    post:
      tags: [Auth & Users]
      summary: Yeni Kullanıcı Kaydı
      operationId: registerUser
      description: Kullanıcı bilgilerini sisteme kaydeder ve hesap oluşturur.
      security: [] 
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/RegisterInput'
      responses:
        "201": { description: Kullanıcı başarıyla oluşturuldu }
        "400": 
          description: Geçersiz veri girişi
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'
  
  /auth/login:
    post:
      tags: [Auth & Users]
      summary: Kullanıcı Girişi
      operationId: loginUser
      description: Kimlik bilgilerini doğrular ve JWT erişim token'ı üretir.
      security: []
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/LoginInput'
      responses:
        "200": 
          description: Giriş başarılı, token üretildi
          content:
            application/json:
              schema:
                type: object
                properties:
                  token:
                    type: string
                    example: "eyJhbGciOiJIUzI1Ni..."
  /auth/logout:
    post:
      tags: [Auth & Users]
      summary: Oturum Kapatma
      operationId: logoutUser
      security:
        - BearerAuth: []
      responses:
        "200": { description: Başarıyla çıkış yapıldı }   

  /user/profile:
    put:
      tags: [Auth & Users]
      summary: Profil Bilgilerini Güncelleme
      operationId: updateProfile
      security:
        - BearerAuth: []
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/ProfileUpdateInput'
      responses:
        "200": { description: Profil başarıyla güncellendi }
        "401": { $ref: '#/components/responses/Unauthorized' }
    delete:
      tags: [Auth & Users]
      summary: Hesap Silme
      operationId: deleteAccount
      security:
        - BearerAuth: []
      responses:
        "204": { description: Kullanıcı hesabı kalıcı olarak silindi }
        "401": { $ref: '#/components/responses/Unauthorized' }

  /user/password:
    patch:
      tags: [Auth & Users]
      summary: Şifre Değiştirme
      operationId: changePassword
      security:
        - BearerAuth: []
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/PasswordUpdateInput'
      responses:
        "200": { description: Şifre başarıyla güncellendi }
        "400": { description: Geçersiz veri }
        "401": { description: Yetkisiz erişim }

  # ==========================================
  # PRODUCTS & WATCHLIST (Nisanur Sütçü & Ercan Aziz)
  # ==========================================
  
  /products:
    get:
      summary: Ürünleri listele, filtrele ve sırala
      description: Marka filtresi ile fiyat/tarih sıralama seçeneklerini bir arada sunar.
      tags:
        - Products
      parameters:
        - name: brand
          in: query
          description: Markaya göre filtrele.
          schema:
            type: string
        - name: sortBy
          in: query
          description: Sıralama kriteri.
          schema:
            type: string
            enum: [price, createdAt]
        - name: order
          in: query
          description: Sıralama yönü.
          schema:
            type: string
            enum: [asc, desc]
            default: asc
      responses:
        '200':
          description: Filtrelenmiş ve sıralanmış ürün listesi.
          content:
            application/json:
              schema:
                type: array
                items:
                  $ref: '#/components/schemas/Product'

  /watchlist:
    get:
      tags:
        - Products
      summary: Ürün listesi
      description: Kullanıcının takip ettiği ürünleri listeler
      operationId: listProducts
      security:
        - bearerAuth: []
      responses:
        "200":
          description: Ürünler başarıyla listelendi
          content:
            application/json:
              schema:
                type: array
                items:
                  $ref: "#/components/schemas/Product"
  /watchlist/add:
    post:
      tags:
        - Products
      summary: Ürün ekleme
      description: Bir ürünü takip listesine ekler
      operationId: addProduct
      security:
        - bearerAuth: []
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: "#/components/schemas/ProductInput"
      responses:
        "201":
          description: Ürün başarıyla eklendi
        "400":
          $ref: "#/components/responses/BadRequest"

  /products/search:
    get:
      tags:
        - Products
      summary: Ürün arama
      description: Anahtar kelime ile ürün arama
      operationId: searchProducts
      parameters:
        - name: query
          in: query
          required: true
          description: Aranacak ürün adı veya anahtar kelime
          schema:
            type: string
      responses:
        "200":
          description: Arama sonuçları
          content:
            application/json:
              schema:
                type: array
                items:
                  $ref: "#/components/schemas/Product"


  /products/category:
    get:
      tags:
        - Categories
      summary: Kategori listesi
      description: Ürün kategorilerini listeler
      operationId: listCategories
      responses:
        "200":
          description: Kategoriler başarıyla listelendi
          content:
            application/json:
              schema:
                type: array
                items:
                  $ref: "#/components/schemas/Category"
  /products/categories:
    get:
      summary: Kategoriye göre ürünleri getir
      tags:
        - Products
      parameters:
        - name: categoryName
          in: query
          required: true
          schema:
            type: string
          example: "Elektronik"
      responses:
        '200':
          description: Kategorize edilmiş ürün listesi.
          content:
            application/json:
              schema:
                type: array
                items:
                  $ref: '#/components/schemas/Product'
        '404':
          $ref: '#/components/responses/NotFound'

  /products/trending:
    get:
      tags:
        - Products
      summary: Favori ürünler
      description: En sık kontrol edilen ürünleri listeler
      operationId: trendProducts
      responses:
        "200":
          description: Trend ürünler listelendi
          content:
            application/json:
              schema:
                type: array
                items:
                  $ref: "#/components/schemas/Product"

  /watchlists/{id}:
    delete:
      tags:
        - Products
      summary: Ürün silme
      description: Takip edilen ürünü listeden kaldırır
      operationId: deleteProduct
      security:
        - bearerAuth: []
      parameters:
        - name: id
          in: path
          required: true
          description: Silinecek ürün ID
          schema:
            type: string
      responses:
        "204":
          description: Ürün başarıyla silindi
        "404":
          $ref: "#/components/responses/NotFound"

  # ==========================================
  # ALERTS & NOTIFICATIONS (Doğukan Dursoy)
  # ==========================================
  /alerts:
    get:
      tags: [Alerts & Notifications]
      summary: Aktif Alarmları Listeleme
      operationId: listActiveAlerts
      responses:
        "200":
          description: Kullanıcının aktif fiyat alarmları listelendi
          content:
            application/json:
              schema:
                type: array
                items:
                  $ref: '#/components/schemas/Alert'
        "401":
          description: Yetkisiz erişim (Kullanıcı giriş yapmamış)
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'
    post:
      tags: [Alerts & Notifications]
      summary: Fiyat Alarmı Oluşturma
      description: Belirli bir ürün için hedef fiyat belirlenir. Kafka/RabbitMQ kuyruğuna dinleyici olarak eklenir.
      operationId: createPriceAlert
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/AlertInput'
      responses:
        "201":
          description: Alarm başarıyla kuruldu
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Alert'
        "400":
          description: Geçersiz veya eksik veri gönderildi (örn. productId eksik)
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'
        "401":
          description: Yetkisiz erişim

  /alerts/{alertId}:
    put:
      tags: [Alerts & Notifications]
      summary: Alarm Durumu Güncelleme
      description: Alarmı aktif/pasif duruma getirme veya hedef fiyatı güncelleme.
      operationId: updateAlertStatus
      parameters:
        - name: alertId
          in: path
          required: true
          schema: { type: string }
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              properties:
                isActive: { type: boolean }
                targetPrice: { type: number, format: float }
      responses:
        "200": 
          description: Alarm güncellendi
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Alert'
        "400":
          description: Geçersiz veri formatı
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'
        "404":
          description: Güncellenmek istenen alarm bulunamadı
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'
    delete:
      tags: [Alerts & Notifications]
      summary: Alarm Silme
      operationId: deleteAlert
      parameters:
        - name: alertId
          in: path
          required: true
          schema: { type: string }
      responses:
        "204": 
          description: Alarm başarıyla silindi (İçerik dönmez)
        "404":
          description: Silinmek istenen alarm zaten yok
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'

  /notify/email:
    post:
      tags: [Alerts & Notifications]
      summary: E-posta Bildirimi Gönderme
      description: Asenkron worker'lar tarafından fiyat düştüğünde tetiklenen iç servis uç noktası.
      operationId: sendEmailNotification
      responses:
        "202": 
          description: E-posta başarıyla kuyruğa eklendi
        "500":
          description: Kuyruk servisi (RabbitMQ/Kafka) bağlantı hatası
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'

  /notify/push:
    post:
      tags: [Alerts & Notifications]
      summary: Anlık Bildirim (Push) Gönderme
      operationId: sendPushNotification
      responses:
        "202": 
          description: Push bildirim başarıyla kuyruğa eklendi
        "500":
          description: Sunucu içi hata
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'

  # ==========================================
  # PRICE ANALYTICS (Sema Durgut)
  # ==========================================
  /products/{id}/history:
    get:
      tags: 
        - products
      summary: Ürün Fiyat Geçmişi Listeleme
      description: Seçili ürünün geçmiş tarihlerdeki fiyat kayıtlarını liste olarak sunar.
      operationId: getPriceHistory
      security:
        - bearerAuth: []
      parameters:
        - $ref: '#/components/parameters/ProductIdParam'
      responses:
        '200':
          description: Fiyat geçmişi başarıyla getirildi.
          content:
            application/json:
              schema:
                type: array
                items:
                  $ref: '#/components/schemas/PriceHistory'
        '401':
          $ref: '#/components/responses/Unauthorized'
        '404':
          $ref: '#/components/responses/NotFound'

  /products/{id}/lowest-price:
    get:
      tags: 
        - products
      summary: Tüm Zamanların En Düşük Fiyatı
      description: Ürünün tüm zamanlardaki en düşük fiyat bilgisini kullanıcıya gösterir.
      operationId: getLowestPrice
      parameters:
        - $ref: '#/components/parameters/ProductIdParam'
      responses:
        '200':
          description: Başarılı
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/PriceInfo'
        '404':
          $ref: '#/components/responses/NotFound'

  /products/compare:
    get:
      tags: 
        - products
      summary: Mağazalar Arası Fiyat Karşılaştırması
      description: Aynı ürünün farklı e-ticaret sitelerindeki fiyatlarını yan yana getirir.
      operationId: compareStores
      parameters:
        - name: query
          in: query
          description: Karşılaştırılacak ürünün adı, modeli veya SKU kodu.
          required: true
          schema:
            type: string
      responses:
        '200':
          description: Karşılaştırma listesi başarıyla oluşturuldu.
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ComparisonList'
        '400':
          $ref: '#/components/responses/BadRequest'

  /products/{id}/discount-rate:
    get:
      tags: 
        - products
      summary: İndirim Oranı Analizi
      description: Ürünün piyasa ortalamasına ve geçmiş fiyatlarına göre kârlılık oranını hesaplar.
      operationId: calculateDiscountRate
      parameters:
        - $ref: '#/components/parameters/ProductIdParam'
      responses:
        '200':
          description: Analiz tamamlandı.
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/DiscountAnalysis'

  /products/{id}/forecast:
    get:
      tags: 
        - products
      summary: Yapay Zeka Destekli Fiyat Tahmini
      description: Mevcut trendlere göre fiyatın gelecekteki olası yönünü raporlar (AI Destekli).
      operationId: getPriceForecast
      parameters:
        - $ref: '#/components/parameters/ProductIdParam'
      responses:
        '200':
          description: Tahmin raporu başarıyla oluşturuldu.
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ForecastData'
        '500':
          description: Tahmin motoru hatası veya yetersiz veri.
  /products/filter/color:
    get:
      tags:
        - products
      summary: Renk Bazlı Ürün Filtreleme
      description: Kullanıcıların ürünleri belirli renk seçeneklerine göre süzerek listelemesi için kullanılır.
      operationId: getProductsByColor
      parameters:
        - name: color
          in: query
          description: Filtrelenecek renk (örn. Siyah, Kırmızı, Space Gray)
          required: true
          schema:
            type: string
      responses:
        '200':
          description: Filtrelenmiş liste döndürüldü.
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ProductList'
  /products/{productId}/chart-data:
    get:
      tags: 
        - products
      summary: Grafik İçin Zaman Serisi Verisi
      description: Grafik çizimi için zaman serisi tabanlı fiyat verilerini hazırlar.
      operationId: getChartData
      parameters:
        - $ref: '#/components/parameters/ProductIdParam'
      responses:
        '200':
          description: Grafik verisi başarıyla hazırlandı.
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ChartData'
        '404':
          $ref: '#/components/responses/NotFound'

  # ==========================================
  # SYSTEM & MISC (Ercan Aziz)
  # ==========================================
  /support/feedback:
    post:
      summary: Kullanıcı geri bildirimi oluştur
      description: Teknik sorunları veya önerileri yönetime iletir.
      tags:
        - Support
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/Feedback'
      responses:
        '201':
          description: Geri bildirim başarıyla alındı.
        '400':
          $ref: '#/components/responses/BadRequest'

  /subscriptions:
    put:
      summary: Abonelik planı limitlerini güncelle
      description: Kullanıcıların takip edebileceği ürün limitlerini belirler (Admin yetkisi gerektirir).
      tags:
        - Admin
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/SubscriptionPlan'
      responses:
        '200':
          description: Plan başarıyla güncellendi.
        '403':
          description: Yetkisiz erişim.




components:
  securitySchemes:
    BearerAuth:
      type: http
      scheme: bearer
      bearerFormat: JWT
  
  parameters:
    ProductIdParam:
      name: productId
      in: path
      description: İşlem yapılacak ürünün benzersiz ID değeri.
      required: true
      schema:
        type: string
        format: uuid
      

  schemas:
    Feedback:
      type: object
      required:
        - userId
        - subject
        - message
      properties:
        userId:
          type: integer
        subject:
          type: string
        message:
          type: string

    SubscriptionPlan:
      type: object
      required:
        - planType
        - productLimit
      properties:
        planType:
          type: string
          enum: [Basic, Premium, Enterprise]
        productLimit:
          type: integer
          example: 50

    PriceHistory:
      type: object
      properties:
        price:
          type: number
          format: float
          example: 1250.50
        currency:
          type: string
          example: "TRY"
        recordedAt:
          type: string
          format: date-time
          example: "2024-03-01T14:30:00Z"
        storeName:
          type: string
          example: "Amazon"

    ChartData:
      type: object
      properties:
        labels:
          type: array
          items:
            type: string
            format: date
          example: ["2024-01-01", "2024-02-01", "2024-03-01"]
        datasets:
          type: array
          items:
            type: object
            properties:
              label:
                type: string
                example: "Fiyat Trendi"
              data:
                type: array
                items:
                  type: number
                example: [1100, 1300, 1250]

    PriceInfo:
      type: object
      properties:
        lowestPrice:
          type: number
          example: 899.99
        highestPrice:
          type: number
          example: 1450.00
        currentPrice:
          type: number
          example: 1050.00
        attainedAt:
          type: string
          format: date-time
          example: "2023-11-15T09:00:00Z"
        storeUrl:
          type: string
          format: uri
          example: "https://example.com/product/123"

    ComparisonList:
      type: object
      properties:
        productName:
          type: string
          example: "iPhone 15 Pro 128GB"
        lastUpdated:
          type: string
          format: date-time
        offers:
          type: array
          items:
            type: object
            properties:
              storeName:
                type: string
                example: "Hepsiburada"
              price:
                type: number
                example: 72000.00
              inStock:
                type: boolean
                example: true
              url:
                type: string
                format: uri

    DiscountAnalysis:
      type: object
      properties:
        currentPrice:
          type: number
          example: 450.00
        averagePrice:
          type: number
          example: 600.00
        discountPercentage:
          type: integer
          example: 25
        rating:
          type: string
          enum: [VERY_GOOD, GOOD, NEUTRAL, BAD]
          description: Fiyatın ne kadar avantajlı olduğunu belirten etiket.
          example: "VERY_GOOD"

    ProductList:
      type: array
      items:
        type: object
        properties:
          id:
            type: string
          name:
            type: string
          color:
            type: string
            example: "Uzay Grisi"
          price:
            type: number
          imageUrl:
            type: string

    ForecastData:
      type: object
      properties:
        productId:
          type: string
        prediction:
          type: string
          enum: [UP, DOWN, STABLE]
          example: "DOWN"
        confidenceScore:
          type: number
          format: float
          description: Tahminin güvenilirlik oranı (0-1 arası).
          example: 0.85
        estimatedPriceNextMonth:
          type: number
          example: 980.00
        aiComment:
          type: string
          example: "Stok verileri ve kampanya dönemleri yaklaştığı için fiyatın %5 düşmesi bekleniyor."
    Product:
      type: object
      required:
        - _id
        - name
        - price
      properties:
        _id:
          type: string
          example: "prod_123"
        name:
          type: string
          example: "Örnek Ürün"
        price:
          type: number
          example: 1999.99
        url:
          type: string
          example: "https://example.com/product"
        brand:
          type: string
        category:
          type: string
        createdAt:
          type: string
          format: date-time
        description:
          type: string

    ProductInput:
      type: object
      required:
        - name
        - url
      properties:
        name:
          type: string
        url:
          type: string
        categoryId:
          type: string

    Category:
      type: object
      required:
        - _id
        - name
      properties:
        _id:
          type: string
          example: "cat_001"
        name:
          type: string
          example: "Elektronik"
    RegisterInput:
      type: object
      required:
        - email
        - password
        - firstName
        - lastName
      properties:
        email:
          type: string
          format: email
          description: Sisteme kayıt için kullanılacak e-posta adresi.
          example: "betul.erkoc@example.com"
        password:
          type: string
          format: password
          description: En az 8 karakterli, güvenli şifre.
          example: "Suleyman123!"
        firstName:
          type: string
          description: Kullanıcının adı.
          example: "Betül"
        lastName:
          type: string
          description: Kullanıcının soyadı.
          example: "Erkoç"

    LoginInput:
      type: object
      required:
        - email
        - password
      properties:
        email:
          type: string
          format: email
          example: "betul.erkoc@example.com"
        password:
          type: string
          format: password
          example: "Suleyman123!"

    ProfileUpdateInput:
      type: object
      properties:
        firstName:
          type: string
          description: Güncellenecek ad bilgisi.
          example: "Betül"
        lastName:
          type: string
          description: Güncellenecek soyad bilgisi.
          example: "Erkoç"
        communicationPreferences:
          type: object
          description: Kullanıcının bildirim ve iletişim tercihleri.
          properties:
            newsletter:
              type: boolean
              example: true
            smsNotifications:
              type: boolean
              example: false

    PasswordUpdateInput:
      type: object
      required:
        - oldPassword
        - newPassword
      properties:
        oldPassword:
          type: string
          format: password
          description: Mevcut şifre doğrulaması.
        newPassword:
          type: string
          format: password
          description: Belirlenecek yeni güçlü şifre.
          example: "YeniGucluSifre2026!"

    User:
      type: object
      properties:
        _id:
          type: string
          example: "user_65432"
          description: Kullanıcının benzersiz sistem kimliği.
        email:
          type: string
          format: email
          example: "betul.erkoc@example.com"
        firstName:
          type: string
          example: "Betül"
        lastName:
          type: string
          example: "Erkoç"
        isEmailVerified:
          type: boolean
          example: true
        createdOn:
          type: string
          format: date-time
          description: Hesabın oluşturulma zamanı.
    Alert:
      type: object
      properties:
        _id:
          type: string
          description: Alarmın benzersiz ID'si
          example: "alrt_64a7f9"
        productId:
          type: string
          description: Takip edilen ürünün ID'si
          example: "prod_9872"
        targetPrice:
          type: number
          format: float
          description: Beklenen düşüş fiyatı
          example: 1250.50
        currentPriceAtCreation:
          type: number
          format: float
          description: Alarm kurulduğu andaki mevcut fiyat
          example: 1500.00
        isActive:
          type: boolean
          description: Alarm şu an devrede mi?
          example: true
        createdOn:
          type: string
          format: date-time
          description: Alarmın kurulma zamanı
          example: "2026-03-08T15:00:00Z"

   
    AlertInput:
      type: object
      required:
        - productId
        - targetPrice
      properties:
        productId:
          type: string
          description: Hangi ürün takip edilecek?
          example: "prod_9872"
        targetPrice:
          type: number
          format: float
          description: Fiyat kaça düşerse haber verelim?
          example: 1250.50

    Error:
      type: object
      properties:
        message:
          type: string
          example: "Yetkisiz erişim veya geçersiz veri girişi."
        code:
          type: integer
          example: 401

  responses:

    Unauthorized:

      description: Yetkisiz erişim. Geçerli bir Bearer Token gerekli.
   

    BadRequest:
      description: Geçersiz istek
      content:
        application/json:
          schema:
            $ref: "#/components/schemas/Error"

    NotFound:
      description: Kaynak bulunamadı
      content:
        application/json:
          schema:
            $ref: "#/components/schemas/Error"