import { defineStore } from 'pinia';
import Api from '../api';

export const useFinanceStore = defineStore('finance', {
  state: () => ({
    transactions: [],
    summary: {
      income: 0,
      expense: 0,
      balance: 0
    },
    isLoading: false
  }),

  actions: {
    async fetchTransactions() {
      this.isLoading = true;
      try {
        const response = await Api.get('/finance');
        this.transactions = response.data.data;
        // Sekalian ambil summary biar update
        this.fetchSummary();
      } catch (err) {
        console.error(err);
      } finally {
        this.isLoading = false;
      }
    },

    async fetchSummary() {
      try {
        const response = await Api.get('/finance/summary');
        this.summary = response.data;
      } catch (err) {
        console.error(err);
      }
    },

    async addTransaction(form) {
      try {
        await Api.post('/finance', form);
        await this.fetchTransactions(); // Refresh list
      } catch (err) {
        throw err;
      }
    },

    async deleteTransaction(id) {
      try {
        await Api.delete(`/finance/${id}`);
        await this.fetchTransactions();
      } catch (err) {
        throw err;
      }
    }
  }
});