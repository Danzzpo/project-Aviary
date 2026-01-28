<script setup>
import { onMounted, ref, reactive, computed } from 'vue';
import { useBirdStore } from '../../stores/bird';
import { Plus, Search, Filter, Bird, X, Loader2, Save, Edit2, Info, ArrowRight, ArrowUpRight } from 'lucide-vue-next';
import SearchableSelect from '../../components/SearchableSelect.vue';

const birdStore = useBirdStore();

// --- STATE ---
const showModal = ref(false);       // Modal Tambah/Edit
const showDetailModal = ref(false); // Modal Detail (Tracking)
const isSubmitting = ref(false);
const isEditing = ref(false);
const editId = ref(null);
const searchQuery = ref('');        // State untuk Search

// State untuk Detail Modal
const activeBird = ref(null); // Burung yang sedang dilihat detailnya

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

// --- 1. FITUR SEARCH (FILTERING) ---
const filteredBirds = computed(() => {
  if (!searchQuery.value) return birdStore.birds;
  
  const query = searchQuery.value.toLowerCase();
  return birdStore.birds.filter(bird => 
    bird.ring_number.toLowerCase().includes(query) ||
    bird.species.toLowerCase().includes(query) ||
    (bird.mutation && bird.mutation.toLowerCase().includes(query))
  );
});

// --- 2. DATA UTILS ---
const sireOptions = computed(() => {
  return birdStore.birds
    .filter(b => b.gender === 'M' && b.id !== editId.value)
    .map(b => ({ id: b.id, label: b.ring_number, subLabel: b.mutation }));
});

const damOptions = computed(() => {
  return birdStore.birds
    .filter(b => b.gender === 'F' && b.id !== editId.value)
    .map(b => ({ id: b.id, label: b.ring_number, subLabel: b.mutation }));
});

// --- 3. LOGIC TRACKING (SILSILAH) ---

// Cari anak-anak dari burung yang sedang aktif (ActiveBird)
const activeChildren = computed(() => {
  if (!activeBird.value) return [];
  return birdStore.birds.filter(b => 
    b.sire_id === activeBird.value.id || 
    b.dam_id === activeBird.value.id
  );
});

// Buka Modal Detail
const openDetailModal = (birdId) => {
  // Cari data burung lengkap dari store
  const target = birdStore.birds.find(b => b.id === birdId);
  if (target) {
    activeBird.value = target;
    showDetailModal.value = true;
  } else {
    alert("Data burung tidak ditemukan.");
  }
};

// Navigasi di dalam Modal (Ganti view ke burung lain)
const navigateToBird = (birdId) => {
  const target = birdStore.birds.find(b => b.id === birdId);
  if (target) {
    activeBird.value = target; // Otomatis reaktif mengupdate tampilan modal
  }
};

// --- 4. CRUD ACTIONS ---
const openAddModal = () => {
  isEditing.value = false;
  editId.value = null;
  resetForm();
  showModal.value = true;
};

const openEditModal = (bird) => {
  isEditing.value = true;
  editId.value = bird.id;
  Object.assign(form, {
    ring_number: bird.ring_number,
    species: bird.species,
    mutation: bird.mutation,
    gender: bird.gender,
    status: bird.status,
    sire_id: bird.sire_id,
    dam_id: bird.dam_id
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
  try {
    if (isEditing.value) {
      await birdStore.updateBird(editId.value, form);
    } else {
      await birdStore.addBird(form);
    }
    showModal.value = false;
    // Jika sedang buka detail, refresh detailnya juga (optional)
    if(activeBird.value && isEditing.value) navigateToBird(activeBird.value.id);
  } catch (error) {
    const msg = error.response?.data?.error || 'Terjadi kesalahan sistem.';
    alert('GAGAL: ' + msg);
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
const radioClass = "flex items-center gap-3 cursor-pointer border border-slate-200 p-3 rounded-lg flex-1 hover:bg-slate-50 transition text-sm font-medium text-slate-700 shadow-sm relative overflow-hidden";
</script>

<template>
  <div class="h-[calc(100vh-100px)] flex flex-col">
    
    <div class="flex flex-col md:flex-row justify-between items-start md:items-center mb-6 gap-4 flex-shrink-0">
      <div>
        <h1 class="text-2xl font-bold text-slate-800 tracking-tight">Stok Burung</h1>
        <p class="text-slate-500 text-sm">Database lengkap indukan dan anakan farm Anda.</p>
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
          <span class="hidden sm:inline">Tambah Burung</span>
          <span class="sm:hidden">Tambah</span>
        </button>
      </div>
    </div>

    <div class="bg-white rounded-2xl shadow-sm border border-slate-200 overflow-hidden flex flex-col flex-1 min-h-0">
      
      <div v-if="birdStore.isLoading" class="flex-1 flex flex-col items-center justify-center text-slate-400 gap-3">
        <Loader2 class="animate-spin" :size="32" />
        <p class="text-sm font-medium">Mengambil data kandang...</p>
      </div>
      
      <div v-else class="flex-1 overflow-auto scrollbar-thin scrollbar-thumb-slate-300 scrollbar-track-slate-50">
        <table class="w-full text-left border-collapse">
          <thead class="bg-slate-50 text-slate-600 text-xs uppercase tracking-wider font-bold sticky top-0 z-10 shadow-sm">
            <tr>
              <th class="p-4 border-b border-slate-200 w-[20%]">Identitas</th>
              <th class="p-4 border-b border-slate-200 w-[20%]">Genetik</th>
              <th class="p-4 border-b border-slate-200 w-[25%]">Indukan (Tracking)</th>
              <th class="p-4 border-b border-slate-200 w-[15%]">Gender</th>
              <th class="p-4 border-b border-slate-200 w-[10%]">Status</th>
              <th class="p-4 border-b border-slate-200 w-[10%] text-right">Aksi</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100">
            <tr v-for="bird in filteredBirds" :key="bird.id" class="hover:bg-slate-50 transition group">
              
              <td class="p-4 align-top">
                <div class="font-bold text-slate-800 text-base">{{ bird.ring_number }}</div>
                <div class="text-[10px] text-slate-400 font-mono mt-1">ID: #{{ bird.id }}</div>
              </td>

              <td class="p-4 align-top">
                <div class="font-medium text-slate-700">{{ bird.species }}</div>
                <div class="text-xs text-slate-500 bg-slate-100 px-2 py-0.5 rounded-full inline-block mt-1 border border-slate-200">
                  {{ bird.mutation || 'Original' }}
                </div>
              </td>

              <td class="p-4 align-top">
                <div class="flex flex-col gap-2">
                  <div v-if="bird.sire" 
                       @click="openDetailModal(bird.sire.id)"
                       class="flex items-center gap-2 text-xs cursor-pointer hover:bg-blue-50 p-1.5 rounded-lg transition border border-transparent hover:border-blue-100 group/sire w-fit">
                    <span class="w-5 h-5 flex items-center justify-center bg-blue-100 text-blue-700 font-bold rounded text-[10px]">B</span>
                    <span class="font-medium text-slate-600 group-hover/sire:text-blue-600 group-hover/sire:underline">
                      {{ bird.sire.ring_number }}
                    </span>
                    <ArrowUpRight :size="12" class="text-slate-300 group-hover/sire:text-blue-400" />
                  </div>
                  
                  <div v-if="bird.dam" 
                       @click="openDetailModal(bird.dam.id)"
                       class="flex items-center gap-2 text-xs cursor-pointer hover:bg-pink-50 p-1.5 rounded-lg transition border border-transparent hover:border-pink-100 group/dam w-fit">
                    <span class="w-5 h-5 flex items-center justify-center bg-pink-100 text-pink-700 font-bold rounded text-[10px]">I</span>
                    <span class="font-medium text-slate-600 group-hover/dam:text-pink-600 group-hover/dam:underline">
                      {{ bird.dam.ring_number }}
                    </span>
                    <ArrowUpRight :size="12" class="text-slate-300 group-hover/dam:text-pink-400" />
                  </div>
                </div>
              </td>

              <td class="p-4 align-top">
                 <span v-if="bird.gender === 'M'" class="inline-flex items-center gap-1.5 text-blue-700 bg-blue-50 border border-blue-200 px-2.5 py-1 rounded-md text-xs font-bold shadow-sm">
                   <span class="w-1.5 h-1.5 rounded-full bg-blue-500"></span> Jantan
                 </span>
                 <span v-else-if="bird.gender === 'F'" class="inline-flex items-center gap-1.5 text-pink-700 bg-pink-50 border border-pink-200 px-2.5 py-1 rounded-md text-xs font-bold shadow-sm">
                   <span class="w-1.5 h-1.5 rounded-full bg-pink-500"></span> Betina
                 </span>
                 <span v-else class="text-slate-500 bg-slate-100 border border-slate-200 px-2.5 py-1 rounded-md text-xs font-medium">Unknown</span>
              </td>

              <td class="p-4 align-top">
                <span :class="`px-3 py-1 rounded-full text-xs font-bold uppercase tracking-wider ${getStatusColor(bird.status)}`">
                  {{ bird.status }}
                </span>
              </td>

              <td class="p-4 align-top text-right">
                <div class="flex justify-end gap-2">
                  <button 
                    @click="openDetailModal(bird.id)"
                    class="p-2 bg-white border border-slate-200 text-slate-400 hover:text-blue-600 hover:border-blue-500 hover:shadow-md rounded-lg transition"
                    title="Lihat Detail & Silsilah"
                  >
                    <Info :size="16" />
                  </button>
                  <button 
                    @click="openEditModal(bird)"
                    class="p-2 bg-white border border-slate-200 text-slate-400 hover:text-emerald-600 hover:border-emerald-500 hover:shadow-md rounded-lg transition"
                    title="Edit Data"
                  >
                    <Edit2 :size="16" />
                  </button>
                </div>
              </td>
            </tr>
            
            <tr v-if="filteredBirds.length === 0">
              <td colspan="6" class="p-12 text-center">
                <div class="inline-flex bg-slate-50 p-4 rounded-full mb-4">
                  <Search class="text-slate-300" :size="32" />
                </div>
                <h3 class="text-slate-900 font-medium">Data tidak ditemukan</h3>
                <p class="text-slate-500 text-sm mt-1">Coba kata kunci lain atau tambah data baru.</p>
              </td>
            </tr>

          </tbody>
        </table>
      </div>
    </div>

    <div v-if="showDetailModal && activeBird" class="fixed inset-0 z-50 flex items-center justify-center p-4">
      <div class="absolute inset-0 bg-slate-900/60 backdrop-blur-sm" @click="showDetailModal = false"></div>
      
      <div class="bg-white rounded-2xl shadow-2xl w-full max-w-3xl relative z-10 overflow-hidden flex flex-col max-h-[85vh] animate-in fade-in zoom-in-95 duration-200">
        
        <div class="bg-slate-800 text-white px-6 py-4 flex justify-between items-center">
          <div class="flex items-center gap-3">
            <div class="p-2 bg-white/10 rounded-lg">
              <Bird class="text-emerald-400" :size="24" />
            </div>
            <div>
              <h3 class="font-bold text-lg">{{ activeBird.ring_number }}</h3>
              <p class="text-slate-400 text-xs flex gap-2">
                {{ activeBird.species }} 
                <span class="text-slate-600">|</span> 
                {{ activeBird.mutation || 'Original' }}
              </p>
            </div>
          </div>
          <button @click="showDetailModal = false" class="text-slate-400 hover:text-white transition">
            <X :size="24" />
          </button>
        </div>

        <div class="overflow-y-auto p-6 bg-slate-50">
          
          <div class="mb-8">
            <h4 class="text-xs font-bold text-slate-500 uppercase tracking-wider mb-3">Orang Tua (Parents)</h4>
            <div class="grid grid-cols-2 gap-4">
              <div 
                @click="activeBird.sire ? navigateToBird(activeBird.sire.id) : null"
                :class="activeBird.sire ? 'bg-white hover:border-blue-400 cursor-pointer' : 'bg-slate-100 border-dashed'"
                class="border border-slate-200 rounded-xl p-4 transition relative group"
              >
                <span class="absolute top-2 right-2 text-[10px] font-bold bg-blue-100 text-blue-700 px-2 py-0.5 rounded">SIRE (Bapak)</span>
                <div v-if="activeBird.sire">
                  <p class="font-bold text-slate-800 text-lg">{{ activeBird.sire.ring_number }}</p>
                  <p class="text-sm text-slate-500">{{ activeBird.sire.mutation }}</p>
                  <p class="text-xs text-blue-600 mt-2 flex items-center gap-1 group-hover:underline">
                    Lihat Profil <ArrowUpRight :size="12" />
                  </p>
                </div>
                <div v-else class="text-slate-400 text-sm py-2">Tidak ada data Bapak</div>
              </div>

              <div 
                @click="activeBird.dam ? navigateToBird(activeBird.dam.id) : null"
                :class="activeBird.dam ? 'bg-white hover:border-pink-400 cursor-pointer' : 'bg-slate-100 border-dashed'"
                class="border border-slate-200 rounded-xl p-4 transition relative group"
              >
                <span class="absolute top-2 right-2 text-[10px] font-bold bg-pink-100 text-pink-700 px-2 py-0.5 rounded">DAM (Ibu)</span>
                <div v-if="activeBird.dam">
                  <p class="font-bold text-slate-800 text-lg">{{ activeBird.dam.ring_number }}</p>
                  <p class="text-sm text-slate-500">{{ activeBird.dam.mutation }}</p>
                  <p class="text-xs text-pink-600 mt-2 flex items-center gap-1 group-hover:underline">
                    Lihat Profil <ArrowUpRight :size="12" />
                  </p>
                </div>
                <div v-else class="text-slate-400 text-sm py-2">Tidak ada data Ibu</div>
              </div>
            </div>
          </div>

          <div class="bg-white border border-slate-200 rounded-xl p-6 mb-8 shadow-sm">
             <h4 class="text-xs font-bold text-slate-400 uppercase tracking-wider mb-4 border-b pb-2">Detail Burung Ini</h4>
             <div class="grid grid-cols-2 md:grid-cols-4 gap-6">
               <div>
                 <span class="text-xs text-slate-500 block">Status</span>
                 <span :class="`px-2 py-0.5 rounded text-xs font-bold ${getStatusColor(activeBird.status)}`">{{ activeBird.status }}</span>
               </div>
               <div>
                 <span class="text-xs text-slate-500 block">Gender</span>
                 <span v-if="activeBird.gender === 'M'" class="text-blue-600 font-bold">Jantan</span>
                 <span v-else-if="activeBird.gender === 'F'" class="text-pink-600 font-bold">Betina</span>
                 <span v-else class="text-slate-500">Unknown</span>
               </div>
               <div>
                 <span class="text-xs text-slate-500 block">Umur / DOB</span>
                 <span class="text-slate-800 font-medium">-</span>
               </div>
               <div>
                 <span class="text-xs text-slate-500 block">Catatan</span>
                 <span class="text-slate-800 font-medium truncate">{{ activeBird.notes || '-' }}</span>
               </div>
             </div>
          </div>

          <div>
            <h4 class="text-xs font-bold text-slate-500 uppercase tracking-wider mb-3">Keturunan (Offspring)</h4>
            
            <div v-if="activeChildren.length > 0" class="space-y-2">
              <div 
                v-for="child in activeChildren" 
                :key="child.id"
                @click="navigateToBird(child.id)"
                class="bg-white border border-slate-200 hover:border-emerald-400 p-3 rounded-lg flex justify-between items-center cursor-pointer transition group"
              >
                <div class="flex items-center gap-3">
                   <div class="w-8 h-8 rounded-full bg-slate-100 flex items-center justify-center text-slate-500 font-bold text-xs">
                     {{ child.gender === 'M' ? '♂' : (child.gender === 'F' ? '♀' : '?') }}
                   </div>
                   <div>
                     <p class="font-bold text-slate-800 text-sm group-hover:text-emerald-700">{{ child.ring_number }}</p>
                     <p class="text-xs text-slate-500">{{ child.mutation }}</p>
                   </div>
                </div>
                <ArrowRight :size="16" class="text-slate-300 group-hover:text-emerald-500" />
              </div>
            </div>

            <div v-else class="text-center py-6 border-2 border-dashed border-slate-200 rounded-xl">
              <p class="text-slate-400 text-sm">Belum ada data keturunan yang tercatat.</p>
            </div>
          </div>

        </div>
        
        <div class="p-4 bg-white border-t border-slate-200 flex justify-end">
           <button @click="openEditModal(activeBird)" class="flex items-center gap-2 text-slate-600 hover:text-emerald-600 px-4 py-2 rounded-lg font-medium transition hover:bg-slate-50">
             <Edit2 :size="16" /> Edit Burung Ini
           </button>
        </div>
      </div>
    </div>

    <div v-if="showModal" class="fixed inset-0 z-50 flex items-center justify-center p-4">
      <div class="absolute inset-0 bg-slate-900/60 backdrop-blur-sm" @click="showModal = false"></div>
      
      <div class="bg-white rounded-2xl shadow-2xl w-full max-w-2xl relative z-10 overflow-hidden flex flex-col max-h-[90vh] animate-in fade-in zoom-in-95 duration-200">
        <div class="bg-white px-6 py-4 border-b border-slate-100 flex justify-between items-center z-20">
          <div>
            <h3 class="font-bold text-lg text-slate-800">
              {{ isEditing ? 'Edit Data Burung' : 'Tambah Burung Baru' }}
            </h3>
            <p class="text-xs text-slate-500 mt-0.5">Pastikan data yang dimasukkan valid.</p>
          </div>
          <button @click="showModal = false" class="bg-slate-50 p-2 rounded-full text-slate-400 hover:bg-red-50 hover:text-red-500 transition">
            <X :size="20" />
          </button>
        </div>

        <div class="overflow-y-auto p-6 scrollbar-thin scrollbar-thumb-slate-300 scrollbar-track-transparent">
          <form @submit.prevent="handleSubmit" class="space-y-8">
            <div class="space-y-4">
              <div class="flex items-center gap-2 mb-2">
                <div class="w-8 h-8 rounded-full bg-emerald-100 text-emerald-600 flex items-center justify-center font-bold text-xs">1</div>
                <h4 class="text-sm font-bold text-slate-700 uppercase tracking-wide">Identitas</h4>
              </div>
              <div class="grid grid-cols-1 md:grid-cols-2 gap-5 pl-10">
                <div>
                  <label class="block text-sm font-medium text-slate-700 mb-1.5">No. Ring</label>
                  <input v-model="form.ring_number" type="text" required :class="inputClass" />
                </div>
                <div>
                  <label class="block text-sm font-medium text-slate-700 mb-1.5">Status</label>
                  <select v-model="form.status" :class="inputClass" class="bg-white">
                    <option value="AVAILABLE">Available</option>
                    <option value="PAIRED">Paired</option>
                    <option value="SOLD">Sold</option>
                    <option value="DEAD">Dead</option>
                  </select>
                </div>
              </div>
              <div class="pl-10">
                 <label class="block text-sm font-medium text-slate-700 mb-2">Jenis Kelamin</label>
                 <div class="flex gap-4">
                  <label :class="radioClass" class="has-[:checked]:border-blue-500 has-[:checked]:bg-blue-50">
                    <input type="radio" v-model="form.gender" value="M" class="accent-blue-600 w-4 h-4"> <span>Jantan</span>
                  </label>
                  <label :class="radioClass" class="has-[:checked]:border-pink-500 has-[:checked]:bg-pink-50">
                    <input type="radio" v-model="form.gender" value="F" class="accent-pink-600 w-4 h-4"> <span>Betina</span>
                  </label>
                  <label :class="radioClass" class="has-[:checked]:border-slate-500 has-[:checked]:bg-slate-100">
                    <input type="radio" v-model="form.gender" value="UNKNOWN" class="accent-slate-600 w-4 h-4"> <span>Unknown</span>
                  </label>
                </div>
              </div>
            </div>
            
            <hr class="border-slate-100" />

            <div class="space-y-4">
              <div class="flex items-center gap-2 mb-2">
                 <div class="w-8 h-8 rounded-full bg-purple-100 text-purple-600 flex items-center justify-center font-bold text-xs">2</div>
                 <h4 class="text-sm font-bold text-slate-700 uppercase tracking-wide">Genetik</h4>
              </div>
              <div class="grid grid-cols-1 md:grid-cols-2 gap-5 pl-10">
                <div>
                  <label class="block text-sm font-medium text-slate-700 mb-1.5">Spesies</label>
                  <select v-model="form.species" :class="inputClass">
                    <option value="Lovebird">Lovebird</option>
                    <option value="Murai Batu">Murai Batu</option>
                    <option value="Kenari">Kenari</option>
                  </select>
                </div>
                <div>
                  <label class="block text-sm font-medium text-slate-700 mb-1.5">Mutasi</label>
                  <input v-model="form.mutation" type="text" :class="inputClass" />
                </div>
              </div>
            </div>

            <hr class="border-slate-100" />

            <div class="space-y-4">
               <div class="flex items-center gap-2 mb-2">
                 <div class="w-8 h-8 rounded-full bg-orange-100 text-orange-600 flex items-center justify-center font-bold text-xs">3</div>
                 <h4 class="text-sm font-bold text-slate-700 uppercase tracking-wide">Silsilah</h4>
              </div>
              <div class="grid grid-cols-1 md:grid-cols-2 gap-5 pl-10">
                <SearchableSelect label="Pilih Bapak" placeholder="Cari Bapak..." :options="sireOptions" v-model="form.sire_id" />
                <SearchableSelect label="Pilih Ibu" placeholder="Cari Ibu..." :options="damOptions" v-model="form.dam_id" />
              </div>
            </div>

            <div class="pt-4 flex justify-end gap-3 sticky bottom-0 bg-white border-t border-slate-100 mt-4 py-4 z-20">
              <button type="button" @click="showModal = false" class="px-6 py-2.5 text-slate-600 font-medium hover:bg-slate-50 rounded-xl transition">Batal</button>
              <button type="submit" :disabled="isSubmitting" class="flex items-center gap-2 bg-emerald-600 text-white px-8 py-2.5 rounded-xl font-bold hover:bg-emerald-700 transition disabled:opacity-70">
                <Loader2 v-if="isSubmitting" class="animate-spin" :size="20" />
                <Save v-else :size="20" /> Simpan
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>

  </div>
</template>