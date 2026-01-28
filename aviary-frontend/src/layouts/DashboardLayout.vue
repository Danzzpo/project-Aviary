<script setup>
import { ref } from 'vue';
import { RouterLink, RouterView, useRoute, useRouter } from 'vue-router';
import { useAuthStore } from '../stores/auth';
import { 
  LayoutDashboard, 
  Bird, 
  Heart, 
  DollarSign, 
  User, 
  LogOut, 
  Menu, 
  X 
} from 'lucide-vue-next';

const authStore = useAuthStore();
const router = useRouter();
const route = useRoute();
const isSidebarOpen = ref(false);

const menuItems = [
  { name: 'Dashboard', path: '/dashboard', icon: LayoutDashboard },
  { name: 'Stok Burung', path: '/dashboard/birds', icon: Bird },
  { name: 'Penjodohan', path: '/dashboard/pairs', icon: Heart },
  { name: 'Keuangan', path: '/dashboard/finance', icon: DollarSign },
  { name: 'Profil', path: '/dashboard/profile', icon: User },
];

const handleLogout = () => {
  authStore.logout();
  router.push('/login');
};
</script>

<template>
  <div class="min-h-screen bg-slate-50 flex">
    
    <aside 
      class="fixed inset-y-0 left-0 z-50 w-64 bg-slate-900 text-white transition-transform duration-300 transform lg:translate-x-0 lg:static lg:inset-0"
      :class="isSidebarOpen ? 'translate-x-0' : '-translate-x-full'"
    >
      <div class="h-16 flex items-center px-6 border-b border-slate-800">
        <Bird class="text-emerald-500 mr-2" :size="24" />
        <span class="text-xl font-bold tracking-wide">Aviary<span class="text-emerald-500">Pro</span></span>
      </div>

      <nav class="p-4 space-y-2 mt-4">
        <RouterLink 
          v-for="item in menuItems" 
          :key="item.path"
          :to="item.path"
          class="flex items-center gap-3 px-4 py-3 rounded-xl transition-colors"
          :class="route.path === item.path ? 'bg-emerald-600 text-white' : 'text-slate-400 hover:bg-slate-800 hover:text-white'"
        >
          <component :is="item.icon" :size="20" />
          <span class="font-medium">{{ item.name }}</span>
        </RouterLink>

        <button 
          @click="handleLogout"
          class="w-full flex items-center gap-3 px-4 py-3 rounded-xl text-red-400 hover:bg-red-500/10 hover:text-red-300 transition mt-8"
        >
          <LogOut :size="20" />
          <span class="font-medium">Keluar</span>
        </button>
      </nav>
    </aside>

    <div class="flex-1 flex flex-col min-w-0 overflow-hidden">
      
      <header class="bg-white border-b border-slate-200 h-16 flex items-center justify-between px-4 sm:px-6">
        <button @click="isSidebarOpen = !isSidebarOpen" class="lg:hidden p-2 text-slate-600">
          <Menu v-if="!isSidebarOpen" />
          <X v-else />
        </button>

        <div class="flex items-center gap-4 ml-auto">
          <div class="text-right hidden sm:block">
            <p class="text-sm font-bold text-slate-800">{{ authStore.user?.username || 'Peternak' }}</p>
            <p class="text-xs text-slate-500">{{ authStore.user?.email }}</p>
          </div>
          <div class="w-10 h-10 bg-emerald-100 rounded-full flex items-center justify-center text-emerald-700 font-bold border border-emerald-200">
            {{ authStore.user?.username?.charAt(0).toUpperCase() || 'U' }}
          </div>
        </div>
      </header>

      <main class="flex-1 overflow-auto p-4 sm:p-6 md:p-8">
        <RouterView />
      </main>

    </div>

    <div 
      v-if="isSidebarOpen" 
      @click="isSidebarOpen = false"
      class="fixed inset-0 bg-black/50 z-40 lg:hidden"
    ></div>

  </div>
</template>