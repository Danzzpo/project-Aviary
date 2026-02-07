<script setup>
import { onMounted, ref, reactive, computed } from 'vue';
import { useFinanceStore } from '../../stores/finance';
import { Wallet, TrendingUp, TrendingDown, Plus, Trash2, ArrowUpRight, ArrowDownRight } from 'lucide-vue-next';

const financeStore = useFinanceStore();
const showModal = ref(false);
const isSubmitting = ref(false);

const form = reactive({
  type: 'EXPENSE', // Default Pengeluaran
  category: 'PAKAN',
  amount: '',
  date: new Date().toISOString().split('T')[0],
  description: ''
});

onMounted(() => {
  financeStore.fetchTransactions();
});

// Format Rupiah
const formatRupiah = (number) => {
  return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(number);
};

const handleSubmit = async () => {
  if(!form.amount || form.amount <= 0) return alert("Jumlah uang harus diisi");
  
  isSubmitting.value = true;
  try {
    // Konversi amount ke float
    await financeStore.addTransaction({
      ...form,
      amount: parseFloat(form.amount)
    });
    showModal.value = false;
    form.amount = ''; form.description = ''; // Reset
  } catch (e) {
    alert("Gagal simpan");
  } finally {
    isSubmitting.value = false;
  }
};

const handleDelete = async (id) => {
  if(confirm("Hapus catatan ini?")) {
    await financeStore.deleteTransaction(id);
  }
};
</script>

<template>
  <div class="pb-10">
    <div class="flex justify-between items-center mb-6">
      <div>
        <h1 class="text-2xl font-bold text-slate-800">Keuangan</h1>
        <p class="text-slate-500">Catat untung rugi peternakan Anda.</p>
      </div>
      <button @click="showModal = true" class="bg-indigo-600 hover:bg-indigo-700 text-white px-5 py-2.5 rounded-xl font-bold flex items-center gap-2 shadow-lg shadow-indigo-200 transition">
        <Plus :size="20"/> Catat Transaksi
      </button>
    </div>

    <div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
      <div class="bg-white p-6 rounded-2xl border border-slate-200 shadow-sm">
        <div class="flex items-center gap-4 mb-2">
          <div class="p-3 bg-blue-50 text-blue-600 rounded-xl"><Wallet :size="24"/></div>
          <span class="text-sm font-bold text-slate-400 uppercase">Saldo Bersih</span>
        </div>
        <h2 class="text-3xl font-bold text-slate-800">{{ formatRupiah(financeStore.summary.balance) }}</h2>
      </div>

      <div class="bg-white p-6 rounded-2xl border border-slate-200 shadow-sm">
        <div class="flex items-center gap-4 mb-2">
          <div class="p-3 bg-emerald-50 text-emerald-600 rounded-xl"><TrendingUp :size="24"/></div>
          <span class="text-sm font-bold text-slate-400 uppercase">Pemasukan</span>
        </div>
        <h2 class="text-3xl font-bold text-emerald-600">+ {{ formatRupiah(financeStore.summary.income) }}</h2>
      </div>

      <div class="bg-white p-6 rounded-2xl border border-slate-200 shadow-sm">
        <div class="flex items-center gap-4 mb-2">
          <div class="p-3 bg-red-50 text-red-600 rounded-xl"><TrendingDown :size="24"/></div>
          <span class="text-sm font-bold text-slate-400 uppercase">Pengeluaran</span>
        </div>
        <h2 class="text-3xl font-bold text-red-600">- {{ formatRupiah(financeStore.summary.expense) }}</h2>
      </div>
    </div>

    <div class="bg-white rounded-2xl border border-slate-200 shadow-sm overflow-hidden">
      <div class="p-5 border-b border-slate-100 font-bold text-slate-700">Riwayat Transaksi</div>
      
      <div v-if="financeStore.transactions.length === 0" class="p-10 text-center text-slate-400">
        Belum ada data keuangan.
      </div>

      <table v-else class="w-full text-left border-collapse">
        <thead class="bg-slate-50 text-slate-500 text-xs uppercase font-bold">
          <tr>
            <th class="p-4">Tanggal</th>
            <th class="p-4">Kategori</th>
            <th class="p-4">Deskripsi</th>
            <th class="p-4 text-right">Jumlah</th>
            <th class="p-4 text-center">Aksi</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-slate-100 text-sm">
          <tr v-for="trx in financeStore.transactions" :key="trx.id" class="hover:bg-slate-50 transition">
            <td class="p-4 text-slate-600">{{ new Date(trx.date).toLocaleDateString('id-ID') }}</td>
            <td class="p-4 font-bold text-slate-700">{{ trx.category }}</td>
            <td class="p-4 text-slate-500">{{ trx.description || '-' }}</td>
            <td class="p-4 text-right font-bold" :class="trx.type === 'INCOME' ? 'text-emerald-600' : 'text-red-600'">
              {{ trx.type === 'INCOME' ? '+' : '-' }} {{ formatRupiah(trx.amount) }}
            </td>
            <td class="p-4 text-center">
              <button @click="handleDelete(trx.id)" class="text-slate-300 hover:text-red-500 transition"><Trash2 :size="16"/></button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <div v-if="showModal" class="fixed inset-0 z-50 flex items-center justify-center p-4">
      <div class="absolute inset-0 bg-slate-900/60 backdrop-blur-sm" @click="showModal = false"></div>
      <div class="bg-white rounded-2xl shadow-2xl w-full max-w-md relative z-10 p-6 animate-in fade-in zoom-in-95">
        <h3 class="font-bold text-lg mb-4">Catat Keuangan</h3>
        <form @submit.prevent="handleSubmit" class="space-y-4">
          
          <div class="flex bg-slate-100 p-1 rounded-lg">
            <button type="button" @click="form.type = 'INCOME'" class="flex-1 py-2 text-sm font-bold rounded-md transition" :class="form.type === 'INCOME' ? 'bg-white text-emerald-600 shadow-sm' : 'text-slate-500 hover:text-slate-700'">Pemasukan</button>
            <button type="button" @click="form.type = 'EXPENSE'" class="flex-1 py-2 text-sm font-bold rounded-md transition" :class="form.type === 'EXPENSE' ? 'bg-white text-red-600 shadow-sm' : 'text-slate-500 hover:text-slate-700'">Pengeluaran</button>
          </div>

          <div>
            <label class="block text-sm font-bold text-slate-700 mb-1">Kategori</label>
            <select v-model="form.category" class="w-full px-4 py-2 border border-slate-200 rounded-lg outline-none focus:ring-2 focus:ring-indigo-500 bg-white">
              <option v-if="form.type === 'INCOME'" value="JUAL_BURUNG">Penjualan Burung</option>
              <option v-if="form.type === 'INCOME'" value="LAINNYA">Pemasukan Lain</option>
              
              <option v-if="form.type === 'EXPENSE'" value="PAKAN">Pakan (Milet/Jagung)</option>
              <option v-if="form.type === 'EXPENSE'" value="VITAMIN">Obat & Vitamin</option>
              <option v-if="form.type === 'EXPENSE'" value="PERLENGKAPAN">Kandang & Glodok</option>
              <option v-if="form.type === 'EXPENSE'" value="OPERASIONAL">Listrik & Air</option>
              <option v-if="form.type === 'EXPENSE'" value="LAINNYA">Lain-lain</option>
            </select>
          </div>

          <div>
            <label class="block text-sm font-bold text-slate-700 mb-1">Jumlah (Rp)</label>
            <input v-model="form.amount" type="number" min="0" required class="w-full px-4 py-2 border border-slate-200 rounded-lg outline-none focus:ring-2 focus:ring-indigo-500" placeholder="0">
          </div>

          <div>
            <label class="block text-sm font-bold text-slate-700 mb-1">Tanggal</label>
            <input v-model="form.date" type="date" required class="w-full px-4 py-2 border border-slate-200 rounded-lg outline-none focus:ring-2 focus:ring-indigo-500">
          </div>

          <div>
            <label class="block text-sm font-bold text-slate-700 mb-1">Keterangan (Opsional)</label>
            <textarea v-model="form.description" rows="2" class="w-full px-4 py-2 border border-slate-200 rounded-lg outline-none focus:ring-2 focus:ring-indigo-500" placeholder="Contoh: Beli milet 10kg"></textarea>
          </div>

          <div class="flex justify-end gap-2 pt-2">
            <button type="button" @click="showModal = false" class="text-slate-500 hover:bg-slate-100 px-4 py-2 rounded-lg font-medium">Batal</button>
            <button type="submit" :disabled="isSubmitting" class="bg-indigo-600 text-white px-6 py-2 rounded-lg font-bold hover:bg-indigo-700 disabled:opacity-50">Simpan</button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>