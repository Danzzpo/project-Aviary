<script setup>
import { ref, reactive, onMounted, onUnmounted } from 'vue';
import { useAuthStore } from '../../stores/auth';
import { 
  User, CreditCard, Moon, Sun, Globe, HelpCircle, 
  Shield, Mail, Camera, Loader2 
} from 'lucide-vue-next';

const authStore = useAuthStore();
const activeTab = ref('profile'); 
const isSubmitting = ref(false);
const previewImage = ref(null);
const fileInput = ref(null);

const form = reactive({
  username: '',
  email: '',
  profile_pic: null
});

// Detect Mobile & Short Screen
const isMobile = ref(window.innerWidth < 768);
const isShortScreen = ref(window.innerHeight < 600); // Deteksi layar pendek/zoom

const handleResize = () => {
  isMobile.value = window.innerWidth < 768;
  isShortScreen.value = window.innerHeight < 600;
  
  // Jika mobile ATAU layar pendek (karena zoom), izinkan scroll body
  if (isMobile.value || isShortScreen.value) {
    document.body.style.overflow = '';
  } else {
    document.body.style.overflow = 'hidden';
  }
};

onMounted(() => {
  checkDarkMode();
  if(localStorage.getItem('lang')) {
    language.value = localStorage.getItem('lang');
  }
  
  if (authStore.user) {
    form.username = authStore.user.username;
    form.email = authStore.user.email;
  }

  handleResize();
  window.addEventListener('resize', handleResize);
});

onUnmounted(() => {
  document.body.style.overflow = '';
  window.removeEventListener('resize', handleResize);
});

const getProfileUrl = (path) => {
  if (!path) return null;
  return `http://localhost:8080${path}`;
};

const handleFileChange = (event) => {
  const file = event.target.files[0];
  if (file) {
    form.profile_pic = file;
    previewImage.value = URL.createObjectURL(file);
  }
};

const triggerFileInput = () => {
  fileInput.value.click();
};

const handleUpdateProfile = async () => {
  if (!confirm("Apakah Anda yakin? Profil hanya bisa diubah setiap 30 hari.")) return;

  isSubmitting.value = true;
  
  const formData = new FormData();
  formData.append('username', form.username);
  formData.append('email', form.email);
  if (form.profile_pic) {
    formData.append('profile_pic', form.profile_pic);
  }

  try {
    await authStore.updateProfile(formData);
    alert("Profil berhasil diperbarui!");
    previewImage.value = null; 
  } catch (error) {
    const msg = error.response?.data?.error || "Gagal update profil";
    alert(msg);
  } finally {
    isSubmitting.value = false;
  }
};

const isDarkMode = ref(false);
const checkDarkMode = () => {
  if (localStorage.theme === 'dark' || (!('theme' in localStorage) && window.matchMedia('(prefers-color-scheme: dark)').matches)) {
    isDarkMode.value = true;
    document.documentElement.classList.add('dark');
  } else {
    isDarkMode.value = false;
    document.documentElement.classList.remove('dark');
  }
};
const toggleDarkMode = () => {
  isDarkMode.value = !isDarkMode.value;
  if (isDarkMode.value) {
    document.documentElement.classList.add('dark');
    localStorage.theme = 'dark';
  } else {
    document.documentElement.classList.remove('dark');
    localStorage.theme = 'light';
  }
};
const language = ref('id'); 
const changeLanguage = (lang) => {
  language.value = lang;
  localStorage.setItem('lang', lang);
  alert(`Bahasa diubah ke: ${lang === 'id' ? 'Indonesia' : 'English'} (Simulasi)`);
};

const tabs = [
  { id: 'profile', label: 'Profil', icon: User },
  { id: 'billing', label: 'Langganan', icon: CreditCard },
  { id: 'appearance', label: 'Tampilan', icon: Moon },
  { id: 'language', label: 'Bahasa', icon: Globe },
  { id: 'help', label: 'Bantuan', icon: HelpCircle },
];
</script>

<template>
  <div class="w-full flex flex-col px-1 py-2 h-full min-h-0 overflow-hidden">
    
    <h1 class="text-xl font-bold text-slate-800 mb-2 px-1 flex-shrink-0">Pengaturan</h1>

    <div class="flex flex-col md:grid md:grid-cols-4 gap-4 items-start flex-1 min-h-0">
      
      <div class="md:col-span-1 w-full flex-shrink-0 sticky top-0 bg-slate-50 z-20 md:static pb-1 md:pb-0">
        <div class="flex md:flex-col overflow-x-auto md:overflow-visible gap-2 pb-2 md:pb-0 no-scrollbar">
          <button 
            type="button" 
            v-for="tab in tabs" 
            :key="tab.id"
            @click="activeTab = tab.id"
            class="flex items-center gap-2 px-3 py-2 rounded-xl transition-all duration-200 font-medium text-sm whitespace-nowrap"
            :class="activeTab === tab.id ? 'bg-indigo-600 text-white shadow-sm' : 'bg-white text-slate-500 hover:text-slate-800 border border-slate-100'"
          >
            <component :is="tab.icon" :size="16" />
            {{ tab.label }}
          </button>
        </div>
      </div>

      <div class="md:col-span-3 w-full h-full flex flex-col min-h-0 overflow-hidden">
        
        <div v-if="activeTab === 'profile'" class="bg-white rounded-xl shadow-sm border border-slate-200 flex flex-col h-full overflow-hidden relative">
          
          <div class="p-4 border-b border-slate-100 flex-shrink-0 bg-white z-10">
             <h2 class="text-sm font-bold text-slate-800 flex items-center gap-2">
                <User :size="18" class="text-indigo-600"/> Edit Profil
             </h2>
          </div>
          
          <div class="flex flex-col flex-1 min-h-0 overflow-hidden">
             <form @submit.prevent="handleUpdateProfile" class="flex flex-col h-full">
                
                <div class="flex-1 overflow-y-auto p-4 custom-scrollbar">
                    
                    <div class="mb-4 flex flex-col items-center">
                      <div class="relative group cursor-pointer" @click="triggerFileInput">
                        <div class="w-20 h-20 rounded-full overflow-hidden border-2 border-slate-100 shadow-sm bg-indigo-50 flex items-center justify-center">
                          <img v-if="previewImage" :src="previewImage" class="w-full h-full object-cover" />
                          <img v-else-if="authStore.user?.profile_pic" :src="getProfileUrl(authStore.user.profile_pic)" class="w-full h-full object-cover" />
                          <span v-else class="text-indigo-300 text-2xl font-bold">{{ authStore.user?.username?.charAt(0).toUpperCase() }}</span>
                        </div>
                        <div class="absolute inset-0 bg-black/30 rounded-full flex items-center justify-center opacity-0 group-hover:opacity-100 transition duration-200">
                          <Camera class="text-white" :size="20" />
                        </div>
                      </div>
                      <p class="text-[10px] text-slate-400 mt-1">Ubah Foto</p>
                      <input ref="fileInput" type="file" accept="image/*" class="hidden" @change="handleFileChange" />
                    </div>

                    <div class="space-y-3">
                      <div>
                        <label class="block text-[11px] font-bold text-slate-600 mb-1">Username</label>
                        <input v-model="form.username" type="text" required class="w-full px-3 py-2 bg-slate-50 border border-slate-200 rounded-lg focus:ring-1 focus:ring-indigo-500 outline-none text-slate-800 transition text-sm">
                      </div>
                      
                      <div>
                        <label class="block text-[11px] font-bold text-slate-600 mb-1">Email</label>
                        <input v-model="form.email" type="email" required class="w-full px-3 py-2 bg-slate-50 border border-slate-200 rounded-lg focus:ring-1 focus:ring-indigo-500 outline-none text-slate-800 transition text-sm">
                      </div>

                      <div class="bg-amber-50 border border-amber-100 rounded-lg p-2.5 flex gap-2 items-start">
                        <div class="text-amber-600 mt-0.5 flex-shrink-0"><HelpCircle :size="14"/></div>
                        <p class="text-[10px] text-amber-800 leading-tight">
                          Profil hanya bisa diubah sekali setiap <strong>30 hari</strong>.
                        </p>
                      </div>
                    </div>
                </div>

                <div class="p-4 border-t border-slate-50 bg-white flex-shrink-0 z-20">
                  <button 
                    type="submit" 
                    :disabled="isSubmitting"
                    class="w-full py-2 bg-indigo-600 text-white rounded-lg font-bold hover:bg-indigo-700 transition flex items-center justify-center gap-2 disabled:opacity-50 text-sm shadow-md"
                  >
                    <Loader2 v-if="isSubmitting" class="animate-spin" :size="16" />
                    <span>{{ isSubmitting ? 'Menyimpan...' : 'Simpan' }}</span>
                  </button>
                </div>

             </form>
          </div>
        </div>

        <div v-if="activeTab === 'billing'" class="bg-white rounded-xl shadow-sm border border-slate-200 flex flex-col h-full overflow-hidden">
           <div class="p-4 border-b border-slate-100"><h2 class="text-sm font-bold text-slate-800 flex items-center gap-2"><CreditCard :size="18" class="text-indigo-600"/> Langganan</h2></div>
           <div class="p-4 flex-1 overflow-y-auto custom-scrollbar">
             <div class="bg-gradient-to-r from-indigo-600 to-violet-600 rounded-xl p-5 text-white shadow-lg">
                <p class="text-indigo-100 text-xs uppercase tracking-wider font-bold mb-1">Paket Aktif</p>
                <h3 class="text-xl font-bold">Breeder Pro</h3>
                <p class="text-xs opacity-90 mt-1">Berakhir: 12 Des 2026</p>
             </div>
           </div>
        </div>

        <div v-if="activeTab === 'appearance'" class="bg-white rounded-xl shadow-sm border border-slate-200 flex flex-col h-full overflow-hidden">
           <div class="p-4 border-b border-slate-100"><h2 class="text-sm font-bold text-slate-800 flex items-center gap-2"><Moon :size="18" class="text-indigo-600"/> Tampilan</h2></div>
           <div class="p-4 flex-1 overflow-y-auto custom-scrollbar">
             <div class="flex items-center justify-between p-3 border border-slate-200 rounded-lg bg-slate-50">
              <div class="flex items-center gap-3">
                <div class="p-2 rounded-lg bg-white shadow-sm text-slate-700">
                  <component :is="isDarkMode ? Moon : Sun" :size="20" />
                </div>
                <span class="font-bold text-slate-800 text-sm">Mode Gelap</span>
              </div>
              <button @click="toggleDarkMode" class="w-10 h-6 rounded-full p-1 transition-colors duration-300 focus:outline-none flex items-center" :class="isDarkMode ? 'bg-indigo-600' : 'bg-slate-300'">
                <div class="w-4 h-4 bg-white rounded-full shadow-md transform transition-transform duration-300" :class="isDarkMode ? 'translate-x-4' : 'translate-x-0'"></div>
              </button>
            </div>
           </div>
        </div>

        <div v-if="activeTab === 'language'" class="bg-white rounded-xl shadow-sm border border-slate-200 flex flex-col h-full overflow-hidden">
           <div class="p-4 border-b border-slate-100"><h2 class="text-sm font-bold text-slate-800 flex items-center gap-2"><Globe :size="18" class="text-indigo-600"/> Bahasa</h2></div>
           <div class="p-4 flex-1 overflow-y-auto custom-scrollbar space-y-2">
            <button @click="changeLanguage('id')" class="w-full flex items-center justify-between p-3 border rounded-lg hover:border-indigo-500 transition group bg-white" :class="language === 'id' ? 'border-indigo-600 bg-indigo-50' : 'border-slate-200'">
              <span class="font-bold text-slate-800 text-sm">🇮🇩 Bahasa Indonesia</span>
              <div v-if="language === 'id'" class="w-2 h-2 bg-indigo-600 rounded-full"></div>
            </button>
            <button @click="changeLanguage('en')" class="w-full flex items-center justify-between p-3 border rounded-lg hover:border-indigo-500 transition group bg-white" :class="language === 'en' ? 'border-indigo-600 bg-indigo-50' : 'border-slate-200'">
              <span class="font-bold text-slate-800 text-sm">🇺🇸 English (US)</span>
              <div v-if="language === 'en'" class="w-2 h-2 bg-indigo-600 rounded-full"></div>
            </button>
          </div>
        </div>

        <div v-if="activeTab === 'help'" class="bg-white rounded-xl shadow-sm border border-slate-200 flex flex-col h-full overflow-hidden">
           <div class="p-4 border-b border-slate-100"><h2 class="text-sm font-bold text-slate-800 flex items-center gap-2"><HelpCircle :size="18" class="text-indigo-600"/> Bantuan</h2></div>
           <div class="p-4 flex-1 overflow-y-auto custom-scrollbar">
             <div class="grid grid-cols-2 gap-3 mb-4">
                <div class="p-3 border border-slate-200 rounded-lg hover:bg-slate-50 transition cursor-pointer text-center">
                  <Shield :size="20" class="mx-auto text-blue-600 mb-1"/>
                  <h3 class="font-bold text-slate-800 text-xs">Privasi</h3>
                </div>
                <div class="p-3 border border-slate-200 rounded-lg hover:bg-slate-50 transition cursor-pointer text-center">
                  <Mail :size="20" class="mx-auto text-purple-600 mb-1"/>
                  <h3 class="font-bold text-slate-800 text-xs">Support</h3>
                </div>
              </div>
              <textarea rows="3" class="w-full p-3 border border-slate-200 rounded-lg focus:ring-2 focus:ring-indigo-500 outline-none text-sm resize-none bg-slate-50" placeholder="Tulis kritik dan saran..."></textarea>
              <div class="pt-3">
                 <button class="w-full px-5 py-2 bg-indigo-600 text-white font-bold rounded-lg text-sm hover:bg-indigo-700 transition shadow-md">Kirim</button>
              </div>
           </div>
        </div>

      </div>
    </div>
  </div>
</template>

<style scoped>
.custom-scrollbar {
  -ms-overflow-style: none;
  scrollbar-width: none;
}
.custom-scrollbar::-webkit-scrollbar {
  display: none;
}
.no-scrollbar::-webkit-scrollbar {
  display: none;
}
.no-scrollbar {
  -ms-overflow-style: none;
  scrollbar-width: none;
}
</style>