import { useEffect, useState } from 'react'
import { productService } from './productService'
import './index.css'

function App() {
  const [trending, setTrending] = useState([])
  const [watchlist, setWatchlist] = useState([])
  const [categories, setCategories] = useState([])
  const [searchTerm, setSearchTerm] = useState('')
  const [searchResults, setSearchResults] = useState([])
  const [isListOpen, setIsListOpen] = useState(true)

  useEffect(() => {
    loadData()
  }, [])
  const loadData = async () => {
    try {
      console.log("Veri çekme işlemi başladı...");

      const [trendRes, watchRes, catRes] = await Promise.all([
        productService.getTrending(),
        productService.getWatchlist(),
        productService.getCategories()
      ]);

      console.log("Gelen Trend Verisi:", trendRes.data);
      console.log("Gelen Liste Verisi:", watchRes.data);

      setTrending(Array.isArray(trendRes.data) ? trendRes.data : []);
      setWatchlist(Array.isArray(watchRes.data) ? watchRes.data : []);
      setCategories(catRes.data.categories || []);

    } catch (err) {
      console.error("Veri çekilirken hata oluştu:", err);
      if (err.response) {
        console.log("Hata Kodu:", err.response.status);
        console.log("Hata Mesajı:", err.response.data);
      }
    }
  };
 //ürün arama authorized.GET("/products/search")
  const handleSearch = async () => {
    if (!searchTerm) { setSearchResults([]); return; }
    try {
      const res = await productService.searchProducts(searchTerm)
      setSearchResults(Array.isArray(res.data) ? res.data : [])
    } catch (e) { console.error(e) }
  }

  //takip listesine ekleme authorized.POST("/watchlist/add")
  const addToWatchlist = async (product) => {
    try {
      await productService.addToWatchlist(product)
      loadData() // Listeyi tazele
      alert("Takip listesine eklendi!")
    } catch (err) { alert("Eklenemedi, ürün zaten listede olabilir.") }
  }

  //takip listesinden çıkarma authorized.DELETE("/watchlist/:id")
  const removeFromWatchlist = async (id) => {
    try {
      await productService.removeFromWatchlist(id)
      loadData()
    } catch (e) { console.error(e) }
  }


  return (
    <div className="app-container">
      <header className="header">
        <h1 className="logo" onClick={() => window.location.reload()} style={{ cursor: 'pointer' }}>notiFY<span>.</span></h1>
        <div className="search-box">
          <input
            placeholder="Ürün ara..."
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
            onKeyDown={(e) => e.key === 'Enter' && handleSearch()}
          />
          <button onClick={handleSearch}>Ara</button>
        </div>
      </header>

      <div className="main-grid">
        <aside className="sidebar">
          <h3>Kategoriler</h3>
          {categories.length > 0 ? categories.map((cat, i) => (
            <div key={i} className="category-item" onClick={() => { setSearchTerm(cat); handleSearch(); }} style={{ cursor: 'pointer' }}>
              {cat}
            </div>
          )) : <p>Kategori yok</p>}
        </aside>

        <section className="content">
          <h3>{searchResults.length > 0 ? "Sonuçlar" : "Popüler Ürünler"}</h3>
          <div className="product-grid">
            {(searchResults.length > 0 ? searchResults : trending).map((item, index) => (
              <div key={item.id || index} className="product-card">
                <strong>{item.product_name}</strong>
                <p>{item.current_price} TL</p>
                <button className="add-btn" onClick={() => addToWatchlist(item)}>Takip Et</button>
              </div>
            ))}
          </div>
        </section>

        <aside className="watchlist-panel">
          <h3 onClick={() => setIsListOpen(!isListOpen)} style={{ cursor: 'pointer' }}>
            Takip Listem {isListOpen ? '▼' : '▶'}
          </h3>
          {isListOpen && (
            <div className="watchlist-content">
              {watchlist.length === 0 ? <p style={{ fontSize: '12px', color: '#888' }}>Listeniz boş.</p> :
                watchlist.map(item => (
                  <div key={item.id} className="watchlist-item">
                    <span>{item.product_name}</span>
                    <button className="del-btn" onClick={() => removeFromWatchlist(item.id)}>✕</button>
                  </div>
                ))
              }
            </div>
          )}
        </aside>
      </div>
    </div>
  )
}

export default App