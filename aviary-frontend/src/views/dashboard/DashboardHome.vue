<script setup>
import { onMounted } from 'vue';
import { useDashboardStore } from '../../stores/dashboard';
import { 
  Bird, 
  Heart, 
  Egg, 
  CheckCircle, 
  DollarSign, 
  AlertOctagon,
  ArrowRight
} from 'lucide-vue-next';

const dashboardStore = useDashboardStore();

onMounted(() => {
  dashboardStore.fetchStats();
});
</script>

<template>
  <div class="pb-10">
    <div class="flex flex-col md:flex-row justify-between items-start md:items-center mb-8 gap-4">
      <div>
        <h1 class="text-2xl font-bold text-slate-800">Dashboard Overview</h1>
        <p class="text-slate-500">Pantau perkembangan ternak Anda secara real-time.</p>
      </div>
      <router-link to="/dashboard/pairs" class="flex items-center gap-2 text-sm font-bold text-blue-600 hover:text-blue-700 bg-blue-50 px-4 py-2 rounded-lg transition">
        Ke Penjodohan <ArrowRight :size="16" />
      </router-link>
    </div>

    <div v-if="dashboardStore.isLoading" class="grid grid-cols-1 md:grid-cols-3 gap-6 animate-pulse">
      <div class="h-32 bg-slate-200 rounded-2xl"></div>
      <div class="h-32 bg-slate-200 rounded-2xl"></div>
      <div class="h-32 bg-slate-200 rounded-2xl"></div>
    </div>

    <div v-else>
      <h3 class="text-sm font-bold text-slate-400 uppercase tracking-wider mb-4">Produksi & Populasi</h3>
      <div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-10">
        
        <div class="bg-white p-6 rounded-2xl shadow-sm border border-slate-100 flex items-center gap-5 hover:shadow-md transition">
          <div class="w-14 h-14 rounded-2xl bg-blue-50 text-blue-600 flex items-center justify-center">
            <Bird :size="28" />
          </div>
          <div>
            <p class="text-slate-500 text-sm font-medium mb-1">Total Populasi</p>
            <h2 class="text-3xl font-bold text-slate-800">{{ dashboardStore.stats.total_birds }} <span class="text-sm font-normal text-slate-400">Ekor</span></h2>
          </div>
        </div>

        <div class="bg-white p-6 rounded-2xl shadow-sm border border-slate-100 flex items-center gap-5 hover:shadow-md transition">
          <div class="w-14 h-14 rounded-2xl bg-pink-50 text-pink-600 flex items-center justify-center">
            <Heart :size="28" />
          </div>
          <div>
            <p class="text-slate-500 text-sm font-medium mb-1">Pasangan Aktif</p>
            <h2 class="text-3xl font-bold text-slate-800">{{ dashboardStore.stats.active_pairs }} <span class="text-sm font-normal text-slate-400">Pasang</span></h2>
          </div>
        </div>

        <div class="bg-white p-6 rounded-2xl shadow-sm border border-slate-100 flex items-center gap-5 hover:shadow-md transition">
          <div class="w-14 h-14 rounded-2xl bg-orange-50 text-orange-600 flex items-center justify-center">
            <Egg :size="28" />
          </div>
          <div>
            <p class="text-slate-500 text-sm font-medium mb-1">Sedang Mengeram</p>
            <h2 class="text-3xl font-bold text-slate-800">{{ dashboardStore.stats.incubating_eggs }} <span class="text-sm font-normal text-slate-400">Butir</span></h2>
          </div>
        </div>
      </div>

      <h3 class="text-sm font-bold text-slate-400 uppercase tracking-wider mb-4">Status Inventaris</h3>
      <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-6">
        
        <div class="bg-white p-5 rounded-xl border border-slate-200 flex flex-col justify-between h-full">
          <div class="flex justify-between items-start mb-4">
            <div class="p-2 bg-emerald-100 text-emerald-600 rounded-lg"><CheckCircle :size="20"/></div>
            <span class="text-xs font-bold bg-emerald-50 text-emerald-600 px-2 py-1 rounded">READY</span>
          </div>
          <div>
            <h4 class="text-2xl font-bold text-slate-700">{{ dashboardStore.stats.total_available }}</h4>
            <p class="text-xs text-slate-500 font-medium">Burung Siap Jodoh</p>
          </div>
        </div>

        <div class="bg-white p-5 rounded-xl border border-slate-200 flex flex-col justify-between h-full">
          <div class="flex justify-between items-start mb-4">
            <div class="p-2 bg-teal-100 text-teal-600 rounded-lg"><DollarSign :size="20"/></div>
            <span class="text-xs font-bold bg-teal-50 text-teal-600 px-2 py-1 rounded">SOLD</span>
          </div>
          <div>
            <h4 class="text-2xl font-bold text-slate-700">{{ dashboardStore.stats.total_sold }}</h4>
            <p class="text-xs text-slate-500 font-medium">Total Terjual</p>
          </div>
        </div>

        <div class="bg-white p-5 rounded-xl border border-slate-200 flex flex-col justify-between h-full">
          <div class="flex justify-between items-start mb-4">
            <div class="p-2 bg-slate-100 text-slate-600 rounded-lg"><AlertOctagon :size="20"/></div>
            <span class="text-xs font-bold bg-slate-100 text-slate-500 px-2 py-1 rounded">RIP</span>
          </div>
          <div>
            <h4 class="text-2xl font-bold text-slate-700">{{ dashboardStore.stats.total_deceased }}</h4>
            <p class="text-xs text-slate-500 font-medium">Burung Mati</p>
          </div>
        </div>

      </div>
    </div>
  </div>
</template>