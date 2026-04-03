import axios from 'axios';


const API_URL = 'https://notify-n.onrender.com';


const TOKEN = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NzUzMzQyNzUsImlkIjoiNjlkMDIxYmQ1YTc2NjMwNTQ4OWY3MDYxIn0.HmwjayETu9vSJOThebXR0i-e7mtaGGsOI6Q4QhFLgxA";

const api = axios.create({
    baseURL: API_URL,
    headers: {
        'Authorization': `Bearer ${TOKEN}`
    }
});

export const productService = {
    // 6. Popüler Ürünler
    getTrending: () => api.get('/products/trending'),

    // 3. Takip Listesini Listeleme
    getWatchlist: () => api.get('/watchlist'),

    // 1. Ürün Arama (Kategoriler için de bunu kullanıyoruz)
    searchProducts: (query) => api.get(`/products/search?q=${query}`),

    // 2. Takip Listesine Ekleme
    addToWatchlist: (item) => api.post('/watchlist/add', item),

    // 4. Takip Listesinden Çıkarma
    removeFromWatchlist: (id) => api.delete(`/watchlist/${id}`),

    // 5. Kategori Listeleme
    getCategories: () => api.get('/products/category'),
};