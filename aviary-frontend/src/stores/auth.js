import { defineStore } from 'pinia';
import Api from '../api';

export const useAuthStore = defineStore('auth', {
  state: () => ({
    user: null,
    isAuthenticated: false 
  }),

  actions: {
    // 1. REGISTER (Fungsi yang tadi hilang)
    async register(formData) {
      try {
        // Kirim data ke Backend
        await Api.post('/auth/register', formData);
        // Jika sukses, return true
        return true;
      } catch (err) {
        // Lempar error agar bisa ditangkap di RegisterView (untuk alert)
        throw err;
      }
    },

    // 2. LOGIN
    async login(credentials) {
      try {
        await Api.post('/auth/login', credentials);
        // Token otomatis diurus browser (Cookie), kita cukup set state jadi true
        this.isAuthenticated = true;
        return true;
      } catch (err) {
        throw err;
      }
    },

    // 3. LOGOUT
    async logout() {
      try {
        await Api.post('/auth/logout');
      } catch (e) {
        console.error("Logout error", e);
      } finally {
        this.isAuthenticated = false;
        this.user = null;
        window.location.href = '/login';
      }
    }
  }
});