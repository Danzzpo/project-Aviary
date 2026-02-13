<script setup>
import { onMounted, ref, reactive, computed } from 'vue';
import { useBirdStore } from '../../stores/bird';
import { Plus, Search, Bird, X, Loader2, Save, Edit2, Info, ArrowRight, ArrowUpRight } from 'lucide-vue-next';
import SearchableSelect from '../../components/SearchableSelect.vue';

const birdStore = useBirdStore();

// --- STATE ---
const showModal = ref(false);       
const showDetailModal = ref(false); 
const isSubmitting = ref(false);
const isEditing = ref(false);
const editId = ref(null);
const searchQuery = ref('');        

const activeBird = ref(null); 

// Form Data
const form = reactive({
  ring_number: '', 
  species: 'Lovebird', 
  mutation: '',
  gender: 'UNKNOWN', 
  status: 'AVAILABLE', 
  sire_id: null, 
  dam_id: null
});

onMounted(() => {
  birdStore.fetchBirds();
});

// --- FITUR SEARCH & LIMIT 5 DATA ---
const filteredBirds = computed(() => {
  let result = birdStore.birds;

  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase();
    result = result.filter(bird => 
      (bird.ring_number && bird.ring_number.toLowerCase().includes(query)) ||
      (bird.species && bird.species.toLowerCase().includes(query)) ||
      (bird.mutation && bird.mutation.toLowerCase().includes(query))
    );
  }

  return result.slice(0, 5); 
});

// --- DATA UTILS ---
const sireOptions = computed(() => {
  return birdStore.birds
    .filter(b => b.gender === 'M' && (b.id || b.ID) !== editId.value)
    .map(b => ({ id: b.id || b.ID, label: b.ring_number || b.RingNumber, subLabel: b.mutation || b.Mutation }));
});

const damOptions = computed(() => {
  return birdStore.birds
    .filter(b => b.gender === 'F' && (b.id || b.ID) !== editId.value)
    .map(b => ({ id: b.id || b.ID, label: b.ring_number || b.RingNumber, subLabel: b.mutation || b.Mutation }));
});

// --- LOGIC TRACKING (SILSILAH) ---
const activeChildren = computed(() => {
  if (!activeBird.value) return [];
  const parentId = activeBird.value.id || activeBird.value.ID;
  return birdStore.birds.filter(b => 
    b.sire_id === parentId || b.dam_id === parentId || b.SireID === parentId || b.DamID === parentId
  );
});

const openDetailModal = (birdId) => {
  const target = birdStore.birds.find(b => (b.id || b.ID) === birdId);
  if (target) {
    activeBird.value = target;
    showDetailModal.value = true;
  } else {
    alert("Data burung tidak ditemukan.");
  }
};

const navigateToBird = (birdId) => {
  const target = birdStore.birds.find(b => (b.id || b.ID) === birdId);
  if (target) {
    activeBird.value = target; 
  }
};

// --- CRUD ACTIONS ---
const openAddModal = () => {
  isEditing.value = false;
  editId.value = null;
  resetForm();
  showModal.value = true;
};

const openEditModal = (bird) => {
  isEditing.value = true;
  // Amankan ID (Antisipasi Golang huruf besar/kecil)
  editId.value = bird.id || bird.ID; 
  
  Object.assign(form, {
    ring_number: bird.ring_number || bird.RingNumber || '', 
    species: bird.species || bird.Species || 'Lovebird', 
    mutation: bird.mutation || bird.Mutation || '',
    gender: bird.gender || bird.Gender || 'UNKNOWN', 
    status: bird.status || bird.Status || 'AVAILABLE', 
    sire_id: bird.sire_id || bird.SireID || null,
    dam_id: bird.dam_id || bird.DamID || null
  });
  showModal.value = true;
};

const resetForm = () => {
  Object.assign(form, {
    ring_number: '', species: 'Lovebird', mutation: '',
    gender: 'UNKNOWN', status: 'AVAILABLE', sire_id: null, dam_id: null
  });
};

const handleSubmit = async () => {
  isSubmitting.value = true;
  
  // PERBAIKAN 1: Format Payload yang sangat aman untuk Golang
  // Jika form.sire_id kosong, kita kirim angka 0 (Golang GORM mengerti 0 berarti tidak ada relasi untuk uint)
  const payload = {
    ring_number: form.ring_number,
    species: form.species,
    mutation: form.mutation,
    gender: form.gender,
    status: form.status,
    sire_id: Number(form.sire_id) || 0,
    dam_id: Number(form.dam_id) || 0,
  };

  try {
    if (isEditing.value) {
      if (!editId.value) throw new Error("ID Burung tidak valid saat mencoba edit!");
      await birdStore.updateBird(editId.value, payload);
    } else {
      await birdStore.addBird(payload);
    }
    
    showModal.value = false;
    
    // Auto refresh halaman agar tabel langsung up-to-date (opsional tapi aman)
    await birdStore.fetchBirds(); 
    
    if(activeBird.value && isEditing.value) {
      navigateToBird(activeBird.value.id || activeBird.value.ID);
    }
  } catch (error) {
    console.error("DEBUG API ERROR:", error);
    // Tampilkan pesan error spesifik dari Backend jika ada
    const msg = error.response?.data?.error || error.response?.data?.message || error.message || 'Terjadi kesalahan sistem.';
    alert('GAGAL MENYIMPAN:\n' + msg);
  } finally {
    isSubmitting.value = false;
  }
};

const getStatusColor = (status) => {
  switch (status) {
    case 'AVAILABLE': return 'bg-emerald-100 text-emerald-700 ring-1 ring-emerald-600/20';
    case 'PAIRED': return 'bg-blue-100 text-blue-700 ring-1 ring-blue-600/20';
    case 'SOLD': return 'bg-purple-100 text-purple-700 ring-1 ring-purple-600/20';
    case 'DEAD': return 'bg-slate-100 text-slate-700 ring-1 ring-slate-600/20';
    default: return 'bg-slate-50 text-slate-600';
  }
};

const inputClass = "w-full px-4 py-2.5 bg-white border border-slate-200 rounded-lg focus:ring-2 focus:ring-emerald-500 focus:border-emerald-500 outline-none transition text-sm shadow-sm";
</script>

<template>
  <div class="h-full w-full flex flex-col p-1">
    
    <div class="flex flex-col md:flex-row justify-between items-start md:items-center mb-6 gap-4 flex-shrink-0">
      <div>
        <h1 class="text-2xl font-bold text-slate-800 tracking-tight">Stok Burung</h1>
        <p class="text-slate-500 text-sm">Menampilkan maksimal 5 data terbaru.</p>
      </div>

      <div class="flex flex-col sm:flex-row gap-3 w-full md:w-auto">
        <div class="relative group w-full md:w-64">
          <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none text-slate-400 group-focus-within:text-emerald-500">
            <Search :size="18" />
          </div>
          <input 
            v-model="searchQuery"
            type="text" 
            placeholder="Cari Ring / Spesies..." 
            class="block w-full pl-10 pr-3 py-2.5 border border-slate-200 rounded-xl leading-5 bg-white placeholder-slate-400 focus:outline-none focus:bg-white focus:ring-2 focus:ring-emerald-500 focus:border-emerald-500 transition shadow-sm sm:text-sm"
          />
        </div>

        <button 
          @click="openAddModal"
          class="flex items-center justify-center gap-2 bg-emerald-600 hover:bg-emerald-700 text-white px-5 py-2.5 rounded-xl font-medium transition shadow-lg shadow-emerald-200 active:scale-95 whitespace-nowrap"
        >
          <Plus :size="20" />
          <span class="hidden sm:inline">Tambah</span>
        </button>
      </div>
    </div>

    <div class="bg-white rounded-2xl shadow-sm border border-slate-200 flex flex-col flex-1 overflow-hidden">
      
      <div v-if="birdStore.isLoading" class="flex-1 flex flex-col items-center justify-center text-slate-400 gap-3">
        <Loader2 class="animate-spin" :size="32" />
        <p class="text-sm font-medium">Mengambil data...</p>
      </div>
      
      <div v-else class="w-full flex-1 overflow-y-auto no-scroll-visual">
        <table class="w-full text-left border-collapse">
          <thead class="bg-slate-50 text-slate-600 text-xs uppercase tracking-wider font-bold sticky top-0 z-10 shadow-sm border-b border-slate-200">
            <tr>
              <th class="p-4 w-[20%]">Identitas</th>
              <th class="p-4 w-[20%]">Genetik</th>
              <th class="p-4 w-[25%]">Indukan</th>
              <th class="p-4 w-[15%]">Gender</th>
              <th class="p-4 w-[10%]">Status</th>
              <th class="p-4 w-[10%] text-right">Aksi</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100">
            <tr v-for="bird in filteredBirds" :key="bird.id || bird.ID" class="hover:bg-slate-50 transition group">
              
              <td class="p-4 align-top">
                <div class="font-bold text-slate-800 text-base">{{ bird.ring_number || bird.RingNumber }}</div>
              </td>

              <td class="p-4 align-top">
                <div class="font-medium text-slate-700">{{ bird.species || bird.Species }}</div>
                <div class="text-[10px] text-slate-500 inline-block">
                  {{ bird.mutation || bird.Mutation || '-' }}
                </div>
              </td>

              <td class="p-4 align-top">
                <div class="flex gap-2">
                  <div v-if="bird.sire || bird.Sire" @click="openDetailModal(bird.sire_id || bird.SireID)" class="cursor-pointer hover:text-blue-600 text-xs font-bold text-slate-500 bg-slate-100 px-2 py-1 rounded">
                    B: {{ bird.sire?.ring_number || bird.Sire?.RingNumber || 'Lihat' }}
                  </div>
                  <div v-if="bird.dam || bird.Dam" @click="openDetailModal(bird.dam_id || bird.DamID)" class="cursor-pointer hover:text-pink-600 text-xs font-bold text-slate-500 bg-slate-100 px-2 py-1 rounded">
                    I: {{ bird.dam?.ring_number || bird.Dam?.RingNumber || 'Lihat' }}
                  </div>
                </div>
              </td>

              <td class="p-4 align-top">
                  <span v-if="(bird.gender || bird.Gender) === 'M'" class="text-blue-600 font-bold text-xs">Jantan</span>
                  <span v-else-if="(bird.gender || bird.Gender) === 'F'" class="text-pink-600 font-bold text-xs">Betina</span>
                  <span v-else class="text-slate-400 text-xs">?</span>
              </td>

              <td class="p-4 align-top">
                <span :class="`px-2 py-0.5 rounded text-[10px] font-bold uppercase ${getStatusColor(bird.status || bird.Status)}`">
                  {{ bird.status || bird.Status }}
                </span>
              </td>

              <td class="p-4 align-top text-right">
                <div class="flex justify-end gap-1">
                  <button @click="openDetailModal(bird.id || bird.ID)" class="p-1.5 text-slate-400 hover:text-blue-600 hover:bg-blue-50 rounded transition">
                    <Info :size="16" />
                  </button>
                  <button @click="openEditModal(bird)" class="p-1.5 text-slate-400 hover:text-emerald-600 hover:bg-emerald-50 rounded transition">
                    <Edit2 :size="16" />
                  </button>
                </div>
              </td>
            </tr>
            
            <tr v-if="filteredBirds.length === 0">
              <td colspan="6" class="p-8 text-center text-slate-500 text-sm">
                Data tidak ditemukan.
              </td>
            </tr>

            <tr v-if="filteredBirds.length > 0">
               <td colspan="6" class="p-3 text-center text-xs text-slate-400 bg-slate-50/50 italic">
                  Menampilkan maksimal 5 data. Gunakan pencarian untuk menemukan burung lainnya.
               </td>
            </tr>

          </tbody>
        </table>
      </div>
    </div>

    <div v-if="showDetailModal && activeBird" class="fixed inset-0 z-50 flex items-center justify-center p-4">
      <div class="absolute inset-0 bg-slate-900/60 backdrop-blur-sm" @click="showDetailModal = false"></div>
      
      <div class="bg-white rounded-2xl shadow-2xl w-full max-w-3xl relative z-10 overflow-hidden flex flex-col max-h-[85vh] animate-in fade-in zoom-in-95 duration-200">
        <div class="bg-slate-800 text-white px-6 py-4 flex justify-between items-center flex-shrink-0">
          <div class="flex items-center gap-3">
             <Bird class="text-emerald-400" :size="24" />
             <div>
               <h3 class="font-bold text-lg">{{ activeBird.ring_number || activeBird.RingNumber }}</h3>
               <p class="text-slate-400 text-xs">{{ activeBird.species || activeBird.Species }}</p>
             </div>
          </div>
          <button @click="showDetailModal = false" class="text-slate-400 hover:text-white transition"><X :size="24" /></button>
        </div>
        
        <div class="overflow-y-auto p-6 bg-slate-50 no-scroll-visual">
           <div class="grid grid-cols-2 gap-4 mb-6">
              <div class="bg-white p-3 rounded-xl border border-slate-200 text-center">
                 <span class="text-[10px] font-bold text-blue-600 block mb-1">BAPAK</span>
                 <span class="font-bold text-slate-800">{{ (activeBird.sire || activeBird.Sire) ? (activeBird.sire?.ring_number || activeBird.Sire?.RingNumber) : '-' }}</span>
              </div>
              <div class="bg-white p-3 rounded-xl border border-slate-200 text-center">
                 <span class="text-[10px] font-bold text-pink-600 block mb-1">IBU</span>
                 <span class="font-bold text-slate-800">{{ (activeBird.dam || activeBird.Dam) ? (activeBird.dam?.ring_number || activeBird.Dam?.RingNumber) : '-' }}</span>
              </div>
           </div>
           
           <h4 class="text-xs font-bold text-slate-500 uppercase mb-3">Keturunan</h4>
           <div v-if="activeChildren.length > 0" class="space-y-2">
              <div v-for="child in activeChildren" :key="child.id || child.ID" @click="navigateToBird(child.id || child.ID)" class="bg-white p-3 border border-slate-200 hover:border-emerald-300 cursor-pointer rounded-xl flex justify-between items-center transition">
                 <span class="font-bold text-slate-700">{{ child.ring_number || child.RingNumber }}</span>
                 <ArrowRight :size="16" class="text-slate-300" />
              </div>
           </div>
           <div v-else class="text-center text-slate-400 text-sm py-6 border-2 border-dashed border-slate-200 rounded-xl">Belum ada keturunan</div>
        </div>
      </div>
    </div>

    <div v-if="showModal" class="fixed inset-0 z-50 flex items-center justify-center p-4">
      <div class="absolute inset-0 bg-slate-900/60 backdrop-blur-sm" @click="showModal = false"></div>
      
      <div class="bg-white rounded-2xl shadow-2xl w-full max-w-xl relative z-10 overflow-hidden flex flex-col max-h-[90vh] animate-in fade-in zoom-in-95 duration-200">
        
        <div class="bg-white px-6 py-4 border-b border-slate-100 flex justify-between items-center flex-shrink-0">
          <h3 class="font-bold text-lg text-slate-800">{{ isEditing ? 'Edit Data Burung' : 'Tambah Burung Baru' }}</h3>
          <button @click="showModal = false" class="text-slate-400 hover:text-slate-600"><X :size="20" /></button>
        </div>

        <div class="overflow-y-auto p-6 flex-1 no-scroll-visual">
          <form @submit.prevent="handleSubmit" id="birdForm" class="space-y-5">
             
             <div>
                <label class="block text-sm font-medium text-slate-700 mb-1.5">No. Ring</label>
                <input v-model="form.ring_number" type="text" required :class="inputClass" placeholder="Contoh: BNR-001" />
             </div>
             
             <div class="grid grid-cols-2 gap-4">
                <div>
                   <label class="block text-sm font-medium text-slate-700 mb-1.5">Spesies</label>
                   <select v-model="form.species" :class="inputClass">
                     <option value="Lovebird">Lovebird</option>
                     <option value="Murai Batu">Murai Batu</option>
                     <option value="Kenari">Kenari</option>
                   </select>
                </div>
                <div>
                   <label class="block text-sm font-medium text-slate-700 mb-1.5">Mutasi / Warna</label>
                   <input v-model="form.mutation" type="text" :class="inputClass" placeholder="Contoh: Biola Blue" />
                </div>
             </div>
             
             <div>
                <label class="block text-sm font-medium text-slate-700 mb-2">Jenis Kelamin</label>
                <div class="flex gap-3">
                   <label class="flex-1 flex items-center gap-2 p-2.5 border border-slate-200 rounded-lg cursor-pointer hover:bg-slate-50 transition">
                     <input type="radio" v-model="form.gender" value="M" class="accent-blue-600 w-4 h-4"> <span class="text-sm font-medium">Jantan</span>
                   </label>
                   <label class="flex-1 flex items-center gap-2 p-2.5 border border-slate-200 rounded-lg cursor-pointer hover:bg-slate-50 transition">
                     <input type="radio" v-model="form.gender" value="F" class="accent-pink-600 w-4 h-4"> <span class="text-sm font-medium">Betina</span>
                   </label>
                   <label class="flex-1 flex items-center gap-2 p-2.5 border border-slate-200 rounded-lg cursor-pointer hover:bg-slate-50 transition">
                     <input type="radio" v-model="form.gender" value="UNKNOWN" class="accent-slate-600 w-4 h-4"> <span class="text-sm font-medium">?</span>
                   </label>
                </div>
             </div>
             
             <div>
                 <label class="block text-sm font-medium text-slate-700 mb-1.5">Status</label>
                 <select v-model="form.status" :class="inputClass">
                    <option value="AVAILABLE">Tersedia (Available)</option>
                    <option value="PAIRED">Dijodohkan (Paired)</option>
                    <option value="SOLD">Terjual (Sold)</option>
                    <option value="DEAD">Mati (Dead)</option>
                 </select>
             </div>
             
             <div class="grid grid-cols-2 gap-4 border-t border-slate-100 pt-4 mt-2">
                <SearchableSelect label="Pilih Bapak" placeholder="Cari Jantan..." :options="sireOptions" v-model="form.sire_id" />
                <SearchableSelect label="Pilih Ibu" placeholder="Cari Betina..." :options="damOptions" v-model="form.dam_id" />
             </div>

          </form>
        </div>

        <div class="p-4 border-t border-slate-100 bg-slate-50 flex justify-end gap-3 flex-shrink-0">
           <button type="button" @click="showModal = false" class="px-5 py-2.5 text-slate-600 font-medium hover:bg-slate-200 rounded-xl transition text-sm">Batal</button>
           <button type="submit" form="birdForm" :disabled="isSubmitting" class="flex items-center gap-2 bg-emerald-600 text-white px-6 py-2.5 rounded-xl font-bold hover:bg-emerald-700 transition disabled:opacity-70 shadow-md shadow-emerald-200 text-sm">
             <Loader2 v-if="isSubmitting" class="animate-spin" :size="16" />
             <Save v-else :size="16" /> Simpan Data
           </button>
        </div>

      </div>
    </div>

  </div>
</template>

<style>
/* === PERBAIKAN 2: MENGHAPUS SEMUA SCROLLBAR SECARA GLOBAL (CLEAN LOOK) === */
/* CSS ini akan "membunuh" scrollbar bawaan browser (yang ada di foto Anda) di SELURUH HALAMAN ini */
html, body {
  -ms-overflow-style: none !important; /* IE and Edge */
  scrollbar-width: none !important; /* Firefox */
}

::-webkit-scrollbar {
  display: none !important; /* Chrome, Safari and Opera */
  width: 0 !important;
  height: 0 !important;
  background: transparent !important;
}

/* Class tambahan khusus untuk div internal */
.no-scroll-visual {
  -ms-overflow-style: none;
  scrollbar-width: none;
}
.no-scroll-visual::-webkit-scrollbar {
  display: none;
}
</style>