//siteden veri cekme DENEME AMACLI

package main

import (
	"fmt"     //Ekrana yazı yazdırmak (çıktı vermek) için kullanılır.
	"strings" //Yazıları kesmek, birleştirmek, büyütmek/küçültmek veya boşlukları silmek için kullanılır

	"github.com/gocolly/colly/v2" //Senin adına bir web sitesine gider, HTML kodlarını indirir ve istediğin parçaları (fiyat, başlık vb.) bulup sana getirir.
	//Colly: Siteye bağlanır ve HTML sayfasını bir paket gibi alır.
	//Strings: O paketin içindeki düzensiz metinleri (fiyatlardaki boşluklar gibi) ayıklar ve temizler.
	//Fmt: Sonuçları terminal ekranına düzenli bir şekilde basar.
)

func main() {
	c := colly.NewCollector(
		colly.AllowedDomains("www.eyyo.com.tr"), //güvenlilk. sadece o sitede kal
		colly.UserAgent("Mozilla/5.0"),          //kimlik gizleme. tarayıcı gibi davranarak engellenmeyi önler
		colly.CacheDir("./eyyo_cache"),          //hızlandırma. aynı sayfayı tekrar ziyaret ederken önbellekten hızlıca getirir
	)

	// Her bir ürün kartını bulduğunda çalışacak olan kısım
	// Detay sayfası olduğu için .productItem yerine #product-right kullandık
	c.OnHTML("#product-right", func(e *colly.HTMLElement) {
		// Ürün ismini ve fiyatını "cımbızla" çekiyoruz
		name := e.ChildText("#product-title")
		// e.ChildText fonksiyonu etiketin içindeki yazıyı (650,00) çeker
		price := e.ChildText("span.product-price-not-vat")

		// Çıktıyı temizleyelim (boşlukları silmek için)
		name = strings.TrimSpace(name)
		price = strings.TrimSpace(price)

		// Değişkenler burada tanımlandığı için kontrol ve yazdırma burada olmalı
		if name != "" {
			fmt.Println("---------------------------")
			fmt.Printf("👗 Ürün Adı: %s\n", name)
			fmt.Printf("💰 Ürün Fiyatı: %s TL\n", price)
			fmt.Println("---------------------------")
		}
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Eyyo taranıyor:", r.URL)
	})

	err := c.Visit("https://www.eyyo.com.tr/siyah-fermuarli-esofman-takim-ate-4643")
	if err != nil {
		fmt.Println("Hata:", err)
	}
}
