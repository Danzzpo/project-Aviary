import { defineStore } from 'pinia';
import Api from '../api';

export const useBirdStore = defineStore('bird', {
  state: () => ({
    birds: [],
    isLoading: false,
    error: null,
  }),

  actions: {
    // --- FUNGSI AMBIL DATA (YANG SUDAH ADA) ---
    async fetchBirds() {
      this.isLoading = true;
      try {
        const response = await Api.get('/birds');
        this.birds = response.data.data;
      } catch (err) {
        console.error("Gagal ambil data:", err);
      } finally {
        this.isLoading = false;
      }
    },

    // --- FUNGSI BARU: TAMBAH DATA ---
    async addBird(formData) {
      try {
        // Kirim data ke Backend (POST)
        await Api.post('/birds', formData);
        
        // Jika sukses, ambil ulang data terbaru biar tabel update otomatis
        await this.fetchBirds(); 
        return true; // Beri sinyal sukses ke komponen
      } catch (err) {
        console.error("Gagal tambah burung:", err);
        throw err; // Lempar error biar bisa ditangkap komponen
      }
    }
  }
});