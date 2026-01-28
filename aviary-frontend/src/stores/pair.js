import { defineStore } from 'pinia';
import Api from '../api';
import { useBirdStore } from './bird'; // Kita butuh update stok burung juga

export const usePairStore = defineStore('pair', {
  state: () => ({
    pairs: [],
    isLoading: false,
  }),

  actions: {
    async fetchPairs() {
      this.isLoading = true;
      try {
        const response = await Api.get('/pairs');
        this.pairs = response.data.data;
      } catch (err) {
        console.error(err);
      } finally {
        this.isLoading = false;
      }
    },

    async createPair(formData) {
      try {
        await Api.post('/pairs', formData);
        await this.fetchPairs(); // Refresh list pasangan
        
        // Refresh juga list burung agar status AVAILABLE -> PAIRED terupdate di cache
        const birdStore = useBirdStore();
        await birdStore.fetchBirds();
        
        return true;
      } catch (err) {
        throw err;
      }
    },

    async disbandPair(pairID) {
      try {
        await Api.put(`/pairs/${pairID}/disband`);
        await this.fetchPairs();
        const birdStore = useBirdStore();
        await birdStore.fetchBirds();
      } catch (err) {
        throw err;
      }
    }
  }
});