import axios from 'axios';

const Api = axios.create({
  baseURL: 'http://localhost:8080/api', // Alamat Backend Go
  headers: {
    'Content-Type': 'application/json',
  },
});

// INTERCEPTOR (PENTING!)
// Fungsi ini otomatis jalan SEBELUM request dikirim ke backend
Api.interceptors.request.use(
  (config) => {
    // Ambil token dari LocalStorage (kita akan simpan di sini nanti)
    const token = localStorage.getItem('token');
    
    if (token) {
      // Tempelkan token ke Header: "Authorization: Bearer <token>"
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

export default Api;