<script setup>
import { onMounted, ref, reactive, computed } from 'vue';
import { usePairStore } from '../../stores/pair';
import { useBirdStore } from '../../stores/bird';
import { useProductionStore } from '../../stores/production';
import { Heart, Plus, X, Loader2, Egg, Unlink, Trash2 } from 'lucide-vue-next'; // <--- Tambah Trash2
import SearchableSelect from '../../components/SearchableSelect.vue';

const pairStore = usePairStore();
const birdStore = useBirdStore();
const productionStore = useProductionStore();

const showAddModal = ref(false); 
const showEggModal = ref(false); 
const isSubmitting = ref(false);
const activePair = ref(null);    

// Default form date hari ini
const form = reactive({
  cage_name: '',
  sire_id: null,
  dam_id: null,
  date: new Date().toISOString().split('T')[0]
});

const eggFormDate = ref(new Date().toISOString().split('T')[0]);

// Load data saat halaman dibuka
onMounted(() => {
  pairStore.fetchPairs();
  birdStore.fetchBirds();
});

// Filter Burung untuk Dropdown (Hanya yang AVAILABLE)
const availableSires = computed(() => {
  return birdStore.birds
    .filter(b => b.gender === 'M' && b.status === 'AVAILABLE')
    .map(b => ({ id: b.id, label: b.ring_number, subLabel: b.mutation }));
});

const availableDams = computed(() => {
  return birdStore.birds
    .filter(b => b.gender === 'F' && b.status === 'AVAILABLE')
    .map(b => ({ id: b.id, label: b.ring_number, subLabel: b.mutation }));
});

// --- ACTIONS ---

const handleSubmitPair = async () => {
  isSubmitting.value = true;
  try {
    await pairStore.createPair(form);
    showAddModal.value = false;
    form.cage_name = ''; form.sire_id = null; form.dam_id = null;
  } catch (error) {
    alert("Gagal: " + (error.response?.data?.error || "Terjadi kesalahan"));
  } finally {
    isSubmitting.value = false;
  }
};

const handleDisband = async (pairId) => {
  if (confirm("Yakin ingin membubarkan pasangan ini?")) {
    await pairStore.disbandPair(pairId);
  }
};

const openEggModal = async (pair) => {
  activePair.value = pair;
  showEggModal.value = true;
  productionStore.activeProduction = null; 
  await productionStore.fetchProduction(pair.id);
};

const handleAddEgg = async () => {
  if (!eggFormDate.value) return alert("Pilih tanggal dulu");
  isSubmitting.value = true;
  try {
    await productionStore.addEgg(activePair.value.id, eggFormDate.value);
  } catch (error) {
    alert("Gagal tambah telur");
  } finally {
    isSubmitting.value = false;
  }
};

const handleChangeEggStatus = async (eggId, status) => {
  await productionStore.updateEggStatus(eggId, activePair.value.id, status);
};

// --- FITUR BARU: HAPUS TELUR ---
const handleDeleteEgg = async (eggId) => {
  if (confirm("Hapus telur ini secara permanen? Data tidak bisa dikembalikan.")) {
    try {
      await productionStore.deleteEgg(eggId, activePair.value.id);
    } catch (error) {
      alert("Gagal menghapus telur");
    }
  }
};

// Helper untuk mengambil Nama Burung
const getRingNumber = (bird, backupId) => {
  if (!bird) return `(ID: ${backupId} Data Hilang)`; 
  return bird.ring_number || bird.RingNumber || `(ID: ${bird.id} No Name)`;
};

const getMutation = (bird) => {
  if (!bird) return '-';
  return bird.mutation || bird.Mutation || '-';
};

// CSS Classes
const labelClass = "block text-sm font-medium text-slate-700 mb-1";
const inputClass = "w-full px-4 py-2 border border-slate-200 rounded-lg focus:ring-2 focus:ring-emerald-500 outline-none transition";
</script>

<template>
  <div>
    <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center mb-8 gap-4">
      <div>
        <h1 class="text-2xl font-bold text-slate-800">Penjodohan</h1>
        <p class="text-slate-500 text-sm">Kelola pasangan dan produksi telur.</p>
      </div>
      <button 
        @click="showAddModal = true"
        class="flex items-center gap-2 bg-pink-600 hover:bg-pink-700 text-white px-5 py-2.5 rounded-xl font-medium transition shadow-lg shadow-pink-200"
      >
        <Plus :size="20" />
        Pasangan Baru
      </button>
    </div>

    <div v-if="pairStore.isLoading" class="text-center py-12 text-slate-400">
      <Loader2 class="animate-spin mx-auto mb-2" :size="32" /> Memuat data...
    </div>

    <div v-else-if="pairStore.pairs.length === 0" class="text-center py-16 bg-white rounded-2xl border border-dashed border-slate-300">
      <Heart class="text-pink-300 mx-auto mb-4" :size="48" />
      <h3 class="text-lg font-medium text-slate-900">Belum ada pasangan</h3>
      <p class="text-slate-500">Mulai jodohkan burung untuk melihat data di sini.</p>
    </div>

    <div v-else class="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-6">
      <div v-for="pair in pairStore.pairs" :key="pair.id" class="bg-white rounded-2xl shadow-sm border border-slate-200 overflow-hidden hover:shadow-md transition flex flex-col">
        
        <div class="bg-slate-50 px-5 py-3 border-b border-slate-100 flex justify-between items-center">
          <div class="flex items-center gap-2">
            <span class="bg-slate-800 text-white text-xs font-bold px-2 py-1 rounded">
              {{ pair.cage_name }}
            </span>
            <span class="text-xs text-slate-500">{{ new Date(pair.pairing_date).toLocaleDateString('id-ID') }}</span>
          </div>
          <span class="flex h-2 w-2 relative">
            <span class="relative inline-flex rounded-full h-2 w-2 bg-emerald-500"></span>
          </span>
        </div>

        <div class="p-5 flex items-center justify-between relative flex-1">
          
          <div class="text-center w-1/2 pr-2">
            <div :class="pair.sire ? 'text-blue-600' : 'text-red-500'" class="font-bold text-sm truncate">
              {{ getRingNumber(pair.sire, pair.sire_id) }}
            </div>
            <div class="text-xs text-slate-400 truncate">
              {{ getMutation(pair.sire) }}
            </div>
          </div>
          
          <Heart class="text-pink-300 shrink-0" :size="20" />

          <div class="text-center w-1/2 pl-2">
            <div :class="pair.dam ? 'text-pink-600' : 'text-red-500'" class="font-bold text-sm truncate">
              {{ getRingNumber(pair.dam, pair.dam_id) }}
            </div>
            <div class="text-xs text-slate-400 truncate">
              {{ getMutation(pair.dam) }}
            </div>
          </div>

        </div>

        <div class="p-3 bg-slate-50 border-t border-slate-100 grid grid-cols-2 gap-2">
          <button 
            @click="openEggModal(pair)"
            class="flex items-center justify-center gap-2 bg-white border border-slate-200 hover:border-orange-400 hover:text-orange-600 text-slate-600 py-2 rounded-lg text-sm font-medium transition"
          >
            <Egg :size="16" /> Telur
          </button>
          <button 
            @click="handleDisband(pair.id)"
            class="flex items-center justify-center gap-2 bg-white border border-slate-200 hover:border-red-400 hover:text-red-600 text-slate-600 py-2 rounded-lg text-sm font-medium transition"
          >
            <Unlink :size="16" /> Bubar
          </button>
        </div>
      </div>
    </div>

    <div v-if="showAddModal" class="fixed inset-0 z-50 flex items-center justify-center p-4">
      <div class="absolute inset-0 bg-slate-900/60 backdrop-blur-sm" @click="showAddModal = false"></div>
      <div class="bg-white rounded-2xl shadow-2xl w-full max-w-lg relative z-10 p-6 animate-in fade-in zoom-in-95">
        <h3 class="font-bold text-lg mb-4">Pasangan Baru</h3>
        <form @submit.prevent="handleSubmitPair" class="space-y-4">
          <div>
            <label :class="labelClass">Kandang</label>
            <input v-model="form.cage_name" type="text" required :class="inputClass" placeholder="No. Kandang" />
          </div>
          <div>
            <label :class="labelClass" class="text-blue-600">Pejantan</label>
            <SearchableSelect :options="availableSires" v-model="form.sire_id" placeholder="Cari Jantan..." />
          </div>
          <div>
            <label :class="labelClass" class="text-pink-600">Betina</label>
            <SearchableSelect :options="availableDams" v-model="form.dam_id" placeholder="Cari Betina..." />
          </div>
          <div>
            <label :class="labelClass">Tanggal</label>
            <input v-model="form.date" type="date" :class="inputClass" />
          </div>
          <div class="flex justify-end gap-2 pt-4">
            <button type="button" @click="showAddModal = false" class="text-slate-600 px-4 py-2 hover:bg-slate-50 rounded-lg">Batal</button>
            <button type="submit" :disabled="isSubmitting" class="bg-emerald-600 text-white px-4 py-2 rounded-lg font-bold hover:bg-emerald-700 disabled:opacity-50">Simpan</button>
          </div>
        </form>
      </div>
    </div>

    <div v-if="showEggModal && activePair" class="fixed inset-0 z-50 flex items-center justify-center p-4">
      <div class="absolute inset-0 bg-slate-900/60 backdrop-blur-sm" @click="showEggModal = false"></div>
      
      <div class="bg-white rounded-2xl shadow-2xl w-full max-w-2xl relative z-10 flex flex-col max-h-[90vh] animate-in fade-in zoom-in-95">
        <div class="px-6 py-4 border-b border-slate-100 flex justify-between items-center bg-slate-50 rounded-t-2xl">
          <div>
            <h3 class="font-bold text-lg text-slate-800 flex items-center gap-2">
              <Egg class="text-orange-500" /> Telur: {{ activePair.cage_name }}
            </h3>
            <p class="text-xs text-slate-500">
               {{ getRingNumber(activePair.sire, activePair.sire_id) }} x {{ getRingNumber(activePair.dam, activePair.dam_id) }}
            </p>
          </div>
          <button @click="showEggModal = false"><X class="text-slate-400" /></button>
        </div>

        <div class="p-6 overflow-y-auto bg-slate-50/50 flex-1">
          <div class="bg-white p-4 rounded-xl border border-slate-200 shadow-sm mb-6 flex items-end gap-3">
            <div class="flex-1">
              <label class="text-xs font-bold text-slate-500 mb-1 block">Tanggal</label>
              <input v-model="eggFormDate" type="date" :class="inputClass" />
            </div>
            <button @click="handleAddEgg" :disabled="isSubmitting" class="bg-orange-500 hover:bg-orange-600 text-white px-4 py-2.5 rounded-lg font-bold text-sm flex items-center gap-2">
              <Plus :size="18" /> Tambah
            </button>
          </div>

          <div v-if="productionStore.isLoading" class="text-center py-8"><Loader2 class="animate-spin mx-auto text-slate-400" /></div>
          
          <div v-else-if="!productionStore.activeProduction || !productionStore.activeProduction.eggs || productionStore.activeProduction.eggs.length === 0" class="text-center py-10 border-2 border-dashed border-slate-200 rounded-xl">
             <p class="text-slate-400 text-sm">Belum ada telur.</p>
          </div>

          <div v-else class="space-y-3">
            <div v-for="egg in productionStore.activeProduction.eggs" :key="egg.id" class="bg-white p-4 rounded-xl border border-slate-200 shadow-sm flex flex-col sm:flex-row items-center justify-between gap-4">
              
              <div class="flex items-center gap-4">
                <div class="w-10 h-10 rounded-full bg-orange-100 text-orange-600 flex items-center justify-center font-bold text-lg">{{ egg.egg_order }}</div>
                <div>
                  <div class="text-sm font-bold text-slate-700">{{ new Date(egg.laid_date).toLocaleDateString('id-ID') }}</div>
                  <div class="text-xs text-slate-400">Est. Netas: {{ new Date(egg.est_hatch_date).toLocaleDateString('id-ID') }}</div>
                </div>
              </div>
              
              <div class="flex items-center gap-2">
                <select 
                  :value="egg.status" 
                  @change="handleChangeEggStatus(egg.id, $event.target.value)" 
                  class="py-1.5 pl-3 pr-8 rounded-lg text-xs font-bold border outline-none cursor-pointer transition"
                  :class="{
                    'text-slate-600 border-slate-200 bg-white': egg.status === 'PENDING',
                    'text-emerald-600 border-emerald-200 bg-emerald-50': egg.status === 'HATCHED',
                    'text-red-600 border-red-200 bg-red-50': ['INFERTILE', 'DIS', 'BROKEN'].includes(egg.status)
                  }"
                >
                  <option value="PENDING">PENDING</option>
                  <option value="FERTILE">FERTIL</option>
                  <option value="INFERTILE">ZONK</option>
                  <option value="DIS">DIS</option>
                  <option value="BROKEN">PECAH</option>
                  <option value="HATCHED">MENETAS</option>
                </select>

                <button 
                  @click="handleDeleteEgg(egg.id)"
                  class="p-2 text-slate-400 hover:text-red-600 hover:bg-red-50 rounded-lg transition"
                  title="Hapus Telur Permanen"
                >
                  <Trash2 :size="16" />
                </button>
              </div>

            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>