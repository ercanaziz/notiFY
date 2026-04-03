import axios from 'axios';

const API_URL = 'https://notify-n.onrender.com';
const TOKEN = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NzUzMzQyNzUsImlkIjoiNjlkMDIxYmQ1YTc2NjMwNTQ4OWY3MDYxIn0.HmwjayETu9vSJOThebXR0i-e7mtaGGsOI6Q4QhFLgxA";

const api = axios.create({
    baseURL: API_URL,
    headers: {
        // Başına Bearer ve boşluğu burada ekliyoruz
        'Authorization': `Bearer ${TOKEN}`
    }
});

export const productService = {
getTrending: () => api.get('/products/trending'),
getWatchlist: () => api.get('/watchlist'),
searchProducts: (query) => api.get(`/products/search?q=${query}`),
addToWatchlist: (item) => api.post('/watchlist/add', item),
removeFromWatchlist: (id) => api.delete(`/watchlist/${id}`),
getCategories: () => api.get('/products/category'),
};