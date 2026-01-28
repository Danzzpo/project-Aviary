<script setup>
import { ref, reactive } from 'vue';
import { useRouter, RouterLink } from 'vue-router';
import { useAuthStore } from '../stores/auth';
import { User, Mail, Lock, Loader2, AlertCircle, ArrowLeft, CheckCircle } from 'lucide-vue-next';

const router = useRouter();
const authStore = useAuthStore();

const form = reactive({
  username: '',
  email: '',
  password: ''
});

const isLoading = ref(false);
const errorMessage = ref('');
const successMessage = ref('');

const handleRegister = async () => {
  isLoading.value = true;
  errorMessage.value = '';
  successMessage.value = '';

  try {
    await authStore.register(form);
    successMessage.value = 'Registrasi berhasil! Mengalihkan ke halaman login...';
    
    // Tunggu 2 detik lalu pindah ke login
    setTimeout(() => {
      router.push('/login');
    }, 2000);
    
  } catch (error) {
    errorMessage.value = error;
  } finally {
    isLoading.value = false;
  }
};
</script>

<template>
  <div class="min-h-screen flex items-center justify-center bg-slate-50 p-4 relative overflow-hidden">
    
    <div class="absolute top-0 left-0 w-full h-full overflow-hidden z-0">
      <div class="absolute -top-[20%] -right-[10%] w-[500px] h-[500px] bg-emerald-200 rounded-full mix-blend-multiply filter blur-3xl opacity-30 animate-blob"></div>
      <div class="absolute -bottom-[20%] -left-[10%] w-[500px] h-[500px] bg-teal-200 rounded-full mix-blend-multiply filter blur-3xl opacity-30 animate-blob animation-delay-2000"></div>
    </div>

    <div class="w-full max-w-md bg-white/80 backdrop-blur-xl border border-white/20 shadow-2xl rounded-2xl p-8 z-10 animate-in fade-in zoom-in duration-300">
      
      <div class="mb-6 text-center">
        <RouterLink to="/" class="inline-flex items-center text-sm text-slate-400 hover:text-emerald-600 mb-6 transition">
          <ArrowLeft :size="16" class="mr-1" /> Kembali ke Home
        </RouterLink>
        <h2 class="text-3xl font-bold text-slate-800 mb-2">Buat Akun Baru</h2>
        <p class="text-slate-500">Mulai kelola peternakan Anda hari ini.</p>
      </div>

      <div v-if="errorMessage" class="mb-6 bg-red-50 border border-red-200 text-red-600 px-4 py-3 rounded-lg text-sm flex items-center gap-2">
        <AlertCircle :size="18" /> {{ errorMessage }}
      </div>
      
      <div v-if="successMessage" class="mb-6 bg-emerald-50 border border-emerald-200 text-emerald-600 px-4 py-3 rounded-lg text-sm flex items-center gap-2">
        <CheckCircle :size="18" /> {{ successMessage }}
      </div>

      <form @submit.prevent="handleRegister" class="space-y-4">
        
        <div>
          <label class="block text-sm font-medium text-slate-700 mb-1.5">Username</label>
          <div class="relative">
            <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none text-slate-400">
              <User :size="18" />
            </div>
            <input 
              v-model="form.username" 
              type="text" 
              required
              class="block w-full pl-10 pr-3 py-2.5 bg-white border border-slate-200 rounded-xl text-slate-900 placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-emerald-500 focus:border-transparent transition sm:text-sm"
              placeholder="Nama Peternakan / User"
            />
          </div>
        </div>

        <div>
          <label class="block text-sm font-medium text-slate-700 mb-1.5">Email Address</label>
          <div class="relative">
            <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none text-slate-400">
              <Mail :size="18" />
            </div>
            <input 
              v-model="form.email" 
              type="email" 
              required
              class="block w-full pl-10 pr-3 py-2.5 bg-white border border-slate-200 rounded-xl text-slate-900 placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-emerald-500 focus:border-transparent transition sm:text-sm"
              placeholder="nama@email.com"
            />
          </div>
        </div>

        <div>
          <label class="block text-sm font-medium text-slate-700 mb-1.5">Password</label>
          <div class="relative">
            <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none text-slate-400">
              <Lock :size="18" />
            </div>
            <input 
              v-model="form.password" 
              type="password" 
              required
              minlength="6"
              class="block w-full pl-10 pr-3 py-2.5 bg-white border border-slate-200 rounded-xl text-slate-900 placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-emerald-500 focus:border-transparent transition sm:text-sm"
              placeholder="Minimal 6 karakter"
            />
          </div>
        </div>

        <button 
          type="submit" 
          :disabled="isLoading"
          class="w-full flex justify-center py-3 px-4 border border-transparent rounded-xl shadow-lg shadow-emerald-200 text-sm font-bold text-white bg-emerald-600 hover:bg-emerald-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-emerald-500 transition disabled:opacity-70 disabled:cursor-not-allowed mt-6"
        >
          <span v-if="isLoading" class="flex items-center gap-2">
            <Loader2 class="animate-spin" :size="18" /> Mendaftarkan...
          </span>
          <span v-else>Daftar Sekarang</span>
        </button>

      </form>

      <div class="mt-8 text-center text-sm text-slate-500">
        Sudah punya akun? 
        <RouterLink to="/login" class="font-semibold text-emerald-600 hover:text-emerald-500 transition">
          Masuk di sini
        </RouterLink>
      </div>

    </div>
  </div>
</template>

<style scoped>
@keyframes blob {
  0% { transform: translate(0px, 0px) scale(1); }
  33% { transform: translate(30px, -50px) scale(1.1); }
  66% { transform: translate(-20px, 20px) scale(0.9); }
  100% { transform: translate(0px, 0px) scale(1); }
}
.animate-blob {
  animation: blob 7s infinite;
}
.animation-delay-2000 {
  animation-delay: 2s;
}
</style>