<script setup>
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '../stores/auth';
import { Mail, Lock, Loader2, AlertCircle, ArrowLeft } from 'lucide-vue-next';

const email = ref('');
const password = ref('');
const isLoading = ref(false);
const errorMessage = ref('');
const authStore = useAuthStore();
const router = useRouter();

const handleLogin = async () => {
  isLoading.value = true;
  errorMessage.value = '';

  try {
    await authStore.login(email.value, password.value);
    // Jika sukses, arahkan ke dashboard (kita buat nanti)
    router.push('/dashboard'); 
  } catch (error) {
    errorMessage.value = typeof error === 'string' ? error : 'Login gagal, periksa email/password';
  } finally {
    isLoading.value = false;
  }
};
</script>

<template>
  <div class="min-h-screen flex items-center justify-center bg-slate-50 p-4 relative overflow-hidden">
    
    <div class="absolute top-0 left-0 w-full h-full overflow-hidden z-0">
      <div class="absolute -top-[20%] -left-[10%] w-[500px] h-[500px] bg-emerald-200 rounded-full mix-blend-multiply filter blur-3xl opacity-30 animate-blob"></div>
      <div class="absolute -bottom-[20%] -right-[10%] w-[500px] h-[500px] bg-teal-200 rounded-full mix-blend-multiply filter blur-3xl opacity-30 animate-blob animation-delay-2000"></div>
    </div>

    <div class="w-full max-w-md bg-white/80 backdrop-blur-xl border border-white/20 shadow-2xl rounded-2xl p-8 z-10">
      
      <div class="mb-8 text-center">
        <RouterLink to="/" class="inline-flex items-center text-sm text-slate-400 hover:text-emerald-600 mb-6 transition">
          <ArrowLeft :size="16" class="mr-1" /> Kembali ke Home
        </RouterLink>
        <h2 class="text-3xl font-bold text-slate-800 mb-2">Selamat Datang</h2>
        <p class="text-slate-500">Masuk untuk mengelola peternakan Anda</p>
      </div>

      <div v-if="errorMessage" class="mb-6 bg-red-50 border border-red-200 text-red-600 px-4 py-3 rounded-lg text-sm flex items-center gap-2">
        <AlertCircle :size="18" />
        {{ errorMessage }}
      </div>

      <form @submit.prevent="handleLogin" class="space-y-5">
        
        <div>
          <label class="block text-sm font-medium text-slate-700 mb-1.5">Email Address</label>
          <div class="relative">
            <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none text-slate-400">
              <Mail :size="18" />
            </div>
            <input 
              v-model="email" 
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
              v-model="password" 
              type="password" 
              required
              class="block w-full pl-10 pr-3 py-2.5 bg-white border border-slate-200 rounded-xl text-slate-900 placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-emerald-500 focus:border-transparent transition sm:text-sm"
              placeholder="••••••••"
            />
          </div>
        </div>

        <button 
          type="submit" 
          :disabled="isLoading"
          class="w-full flex justify-center py-3 px-4 border border-transparent rounded-xl shadow-sm text-sm font-bold text-white bg-emerald-600 hover:bg-emerald-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-emerald-500 transition disabled:opacity-70 disabled:cursor-not-allowed"
        >
          <span v-if="isLoading" class="flex items-center gap-2">
            <Loader2 class="animate-spin" :size="18" /> Memproses...
          </span>
          <span v-else>Masuk ke Dashboard</span>
        </button>

      </form>

      <div class="mt-8 text-center text-sm text-slate-500">
        Belum punya akun? 
        <RouterLink to="/register" class="font-semibold text-emerald-600 hover:text-emerald-500 transition">
          Daftar Sekarang
        </RouterLink>
      </div>

    </div>
  </div>
</template>

<style scoped>
/* Animasi blob background */
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