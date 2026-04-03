import { useEffect, useState } from 'react'
import { productService } from './productService'
import './index.css'

function App() {
  const [trending, setTrending] = useState([])
  const [watchlist, setWatchlist] = useState([])
  const [categories, setCategories] = useState([])
  const [searchTerm, setSearchTerm] = useState('')
  const [searchResults, setSearchResults] = useState([])

  // LİSTEYİ GÖSTERİP GİZLEMEK İÇİN YENİ DURUM
  const [isListOpen, setIsListOpen] = useState(false)

  useEffect(() => {
    loadData()
  }, [])

  const loadData = async () => {
    try {
      const t = await productService.getTrending(); setTrending(t.data || [])
      const w = await productService.getWatchlist(); setWatchlist(w.data || [])
      const c = await productService.getCategories(); setCategories(c.data.categories || [])
    } catch (err) { console.log("Veri çekilemedi, token bekliyoruz...") }
  }

  const handleSearch = async () => {
    if (!searchTerm) { setSearchResults([]); return; }
    const res = await productService.searchProducts(searchTerm)
    setSearchResults(res.data || [])
  }

  const removeFromWatchlist = async (id) => {
    await productService.removeFromWatchlist(id)
    setWatchlist(watchlist.filter(item => item.id !== id))
  }

  return (
    <div className="app-container">
      <header className="header">
        <h1 className="logo">notiFY<span>.</span></h1>
        <div className="search-box">
          <input
            placeholder="iPhone, Kitap, Ayakkabı..."
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
            onKeyDown={(e) => e.key === 'Enter' && handleSearch()}
          />
          <button onClick={handleSearch}>Ara</button>
        </div>
      </header>

      <div className="main-grid">
        <aside className="sidebar">
          <h3>📂 Kategoriler</h3>
          {categories.map((cat, i) => <div key={i} className="category-item">{cat}</div>)}
        </aside>

        <section className="content">
          <h3>{searchResults.length > 0 ? "🔎 Sonuçlar" : "🔥 Trendler"}</h3>
          <div className="product-grid">
            {(searchResults.length > 0 ? searchResults : trending).map(item => (
              <div key={item.id} className="product-card">
                <strong>{item.product_name}</strong>
                <p>{item.current_price} TL</p>
                <button className="add-btn" onClick={() => alert("Eklendi!")}>Takip Et</button>
              </div>
            ))}
          </div>
        </section>

        {/* --- TIKLANABİLİR TAKİP LİSTESİ --- */}
        <aside className="watchlist-panel">
          <h3
            onClick={() => setIsListOpen(!isListOpen)}
            style={{ cursor: 'pointer', display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}
          >
            📋 Takip Listem {isListOpen ? '▼' : '▶'}
          </h3>

          {isListOpen && (
            <div className="watchlist-content">
              {watchlist.length === 0 ? <p style={{ fontSize: '12px', color: '#888' }}>Listeniz henüz boş.</p> :
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