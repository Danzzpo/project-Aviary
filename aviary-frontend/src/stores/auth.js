import { defineStore } from 'pinia';
import Api from '../api';

export const useAuthStore = defineStore('auth', {
  state: () => ({
    token: localStorage.getItem('token') || '', // Ambil dari storage jika di-refresh
    user: null,
  }),
  
  getters: {
    isLoggedIn: (state) => !!state.token, // Mengembalikan true jika token ada
  },

  actions: {
    async login(email, password) {
      try {
        const response = await Api.post('/auth/login', { email, password });
        
        // Simpan data dari response backend
        this.token = response.data.token;
        this.user = response.data.user;

        // Simpan ke LocalStorage agar tidak hilang saat refresh
        localStorage.setItem('token', this.token);
        
        return true; // Login sukses
      } catch (error) {
        throw error.response.data.error || 'Login Gagal';
      }
    },
    async register(formData) {
      try {
        // Kirim data: username, email, password
        await Api.post('/auth/register', formData);
        return true; // Sukses
      } catch (error) {
        throw error.response?.data?.error || 'Registrasi Gagal';
      }
    },
    logout() {
      this.token = '';
      this.user = null;
      localStorage.removeItem('token');
      // Redirect logic nanti diurus di component
    }
  }
});