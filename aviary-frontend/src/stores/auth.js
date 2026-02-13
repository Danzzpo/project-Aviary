import { defineStore } from 'pinia';
import Api from '../api';

export const useAuthStore = defineStore('auth', {
  state: () => ({
    // 1. Saat refresh, ambil data dari localStorage. Jika tidak ada, kembalikan null/false
    user: JSON.parse(localStorage.getItem('user_data')) || null,
    isAuthenticated: localStorage.getItem('is_logged_in') === 'true' 
  }),

  actions: {
    // --- REGISTER ---
    async register(formData) {
      try {
        await Api.post('/auth/register', formData);
        return true;
      } catch (err) {
        throw err;
      }
    },

    // --- LOGIN ---
    async login(credentials) {
      try {
        const response = await Api.post('/auth/login', credentials);
        
        if (response.data && response.data.user) {
          // Update State Pinia
          this.user = response.data.user;
          this.isAuthenticated = true;

          // 2. SIMPAN KE LOCAL STORAGE AGAR TAHAN REFRESH
          localStorage.setItem('user_data', JSON.stringify(response.data.user));
          localStorage.setItem('is_logged_in', 'true');
        }
        
        return true;
      } catch (err) {
        throw err;
      }
    },

    // --- LOGOUT ---
    async logout() {
      try {
        await Api.post('/auth/logout');
      } catch (e) {
        console.error("Logout error", e);
      } finally {
        // Hapus State Pinia
        this.isAuthenticated = false;
        this.user = null;

        // 3. HAPUS DARI LOCAL STORAGE SAAT LOGOUT
        localStorage.removeItem('user_data');
        localStorage.removeItem('is_logged_in');
        
        window.location.href = '/login';
      }
    },

    // --- UPDATE PROFILE ---
    async updateProfile(formData) {
      try {
        const response = await Api.put('/auth/profile', formData);
        
        if (response.data && response.data.user) {
          // Update State Pinia
          this.user = response.data.user;
          
          // 4. UPDATE DATA DI LOCAL STORAGE JUGA AGAR FOTO/NAMA BARU TIDAK HILANG SAAT REFRESH
          localStorage.setItem('user_data', JSON.stringify(response.data.user));
        }
        
        return response.data;
      } catch (error) {
        throw error;
      }
    }
  }
});