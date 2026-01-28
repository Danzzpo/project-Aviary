import { defineStore } from 'pinia';
import Api from '../api';

export const useProductionStore = defineStore('production', {
  state: () => ({
    activeProduction: null, // Menyimpan data clutch/periode bertelur aktif
    isLoading: false,
  }),

  actions: {
    // 1. Ambil Data Produksi & Telur
    async fetchProduction(pairId) {
      this.isLoading = true;
      try {
        const response = await Api.get(`/pairs/${pairId}/production`);
        // Backend mengirim format: { data: { ... } }
        this.activeProduction = response.data.data; 
      } catch (err) {
        console.error("Gagal load produksi:", err);
        this.activeProduction = null;
      } finally {
        this.isLoading = false;
      }
    },

    // 2. Tambah Telur Baru
    async addEgg(pairId, laidDate) {
      try {
        // Kirim tanggal bertelur ke backend
        await Api.post(`/pairs/${pairId}/eggs`, { laid_date: laidDate });
        // Refresh data agar telur baru langsung muncul di layar
        await this.fetchProduction(pairId); 
        return true;
      } catch (err) {
        throw err; // Lempar error agar bisa ditangkap di Vue Component (untuk alert)
      }
    },

    // 3. Update Status Telur (Menetas, Zonk, Dis, dll)
    async updateEggStatus(eggId, pairId, newStatus) {
      try {
        await Api.put(`/eggs/${eggId}/status`, { status: newStatus });
        await this.fetchProduction(pairId); // Refresh data
      } catch (err) {
        console.error("Gagal update status:", err);
        throw err;
      }
    },

    // 4. Hapus Telur (FITUR BARU)
    // Digunakan untuk menghapus telur zonk/rusak agar database bersih
    async deleteEgg(eggId, pairId) {
      try {
        await Api.delete(`/eggs/${eggId}`);
        await this.fetchProduction(pairId); // Refresh data agar telur hilang dari layar
      } catch (err) {
        console.error("Gagal menghapus telur:", err);
        throw err;
      }
    }
  }
});