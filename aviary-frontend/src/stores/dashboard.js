import { defineStore } from 'pinia';
import Api from '../api';

export const useDashboardStore = defineStore('dashboard', {
  state: () => ({
    stats: {
      total_birds: 0,
      active_pairs: 0,
      incubating_eggs: 0,
      total_available: 0,
      total_sold: 0,
      total_deceased: 0
    },
    isLoading: false
  }),

  actions: {
    async fetchStats() {
      this.isLoading = true;
      try {
        const response = await Api.get('/dashboard/stats');
        
        // UPDATE PENTING:
        // Karena backend kirim struct langsung (tanpa bungkus "data" lagi di JSON-nya),
        // Maka Axios akan menyimpannya di response.data
        this.stats = response.data; 

      } catch (error) {
        console.error("Gagal load stats:", error);
      } finally {
        this.isLoading = false;
      }
    }
  }
});