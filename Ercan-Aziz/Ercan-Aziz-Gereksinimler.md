1. **Kullanıcı Geri Bildirimi Alma**

    **API Metodu**:`POST /support/feedback`

    **Açıklama**: Kullanıcıların yaşadığı teknik sorunları yönetime iletmesini sağlar.

2. **Ürün Çeşidine (Kategori) Göre Listeleme**

    **API Metodu**: `GET /products/categories`

    **Açıklama**: Ürünleri belirli bir kategoriye göre filtrelemek için kullanılır.

3. **Markaya Göre Listeleme**

    **API Metodu**: `GET /products`

    **Açıklama**: Belirli bir markanın ürünlerini getirmek için kullanılır.

4. **Fiyata Göre Sıralama**

    **API Metodu**:`GET /products`

    **Açıklama**: Ürünleri fiyata göre artan veya azalan şekilde sıralamak için kullanılır.

5. **Tarihe Göre Sıralama**

    **API Metodu**: `GET /products`

    **Açıklama**: Ürünlerin eklenme tarihine göre eskiden yeniye veya yeniden eskiye sıralanmasıdır.

6.  **Abonelik Planı Belirleme**

    **API Metodu**: `PUT /admin/subscription-plans`

    **Açıklama**: Kullanıcıların kaç ürün takip edebileceğine dair limitleri yönetir.