<script setup>
import { RouterLink } from 'vue-router';
import { Bird, ShieldCheck, TrendingUp, ChevronRight } from 'lucide-vue-next';
</script>

<template>
  <div class="min-h-screen bg-slate-50 font-sans text-slate-800">
    
    <nav class="container mx-auto px-6 py-6 flex justify-between items-center">
      <div class="flex items-center gap-2">
        <div class="bg-emerald-600 p-2 rounded-lg text-white">
          <Bird :size="24" />
        </div>
        <span class="text-xl font-bold tracking-tight text-slate-900">Aviary<span class="text-emerald-600">Pro</span></span>
      </div>
      <div class="space-x-4">
        <RouterLink to="/login" class="px-5 py-2.5 text-sm font-medium text-emerald-600 hover:text-emerald-700 transition">
          Masuk
        </RouterLink>
        <RouterLink to="/register" class="px-5 py-2.5 text-sm font-medium bg-emerald-600 text-white rounded-full hover:bg-emerald-700 shadow-lg shadow-emerald-200 transition">
          Daftar Gratis
        </RouterLink>
      </div>
    </nav>

    <header class="container mx-auto px-6 py-16 md:py-24 text-center max-w-4xl">
      <div class="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-emerald-50 text-emerald-700 text-xs font-semibold uppercase tracking-wider mb-6 border border-emerald-100">
        <span class="w-2 h-2 rounded-full bg-emerald-500 animate-pulse"></span>
        Sistem Manajemen Kandang No. #1
      </div>
      <h1 class="text-4xl md:text-6xl font-extrabold text-slate-900 leading-tight mb-6">
        Kelola Peternakan Burung <br class="hidden md:block" />
        <span class="text-transparent bg-clip-text bg-gradient-to-r from-emerald-600 to-teal-500">Lebih Cerdas & Efisien</span>
      </h1>
      <p class="text-lg text-slate-600 mb-10 max-w-2xl mx-auto leading-relaxed">
        Pantau stok, pelacakan garis keturunan (pedigree), dan manajemen keuangan dalam satu aplikasi terintegrasi. Tingkatkan produktivitas breeding Anda sekarang.
      </p>
      
      <div class="flex flex-col sm:flex-row gap-4 justify-center">
        <RouterLink to="/login" class="inline-flex items-center justify-center gap-2 px-8 py-4 text-base font-bold text-white bg-emerald-600 rounded-xl hover:bg-emerald-700 transition shadow-xl shadow-emerald-200">
          Mulai Sekarang
          <ChevronRight :size="20" />
        </RouterLink>
        <button class="inline-flex items-center justify-center px-8 py-4 text-base font-bold text-slate-700 bg-white border border-slate-200 rounded-xl hover:bg-slate-50 transition">
          Pelajari Fitur
        </button>
      </div>
    </header>

    <section class="container mx-auto px-6 py-12">
      <div class="grid md:grid-cols-3 gap-8">
        <div class="bg-white p-8 rounded-2xl shadow-sm border border-slate-100 hover:shadow-md transition">
          <div class="w-12 h-12 bg-blue-50 text-blue-600 rounded-xl flex items-center justify-center mb-6">
            <Bird :size="24" />
          </div>
          <h3 class="text-xl font-bold mb-3">Manajemen Stok & Silsilah</h3>
          <p class="text-slate-500 leading-relaxed">
            Data induk jantan/betina tersimpan rapi dengan fitur pelacakan pedigree otomatis hingga ke buyut.
          </p>
        </div>
        <div class="bg-white p-8 rounded-2xl shadow-sm border border-slate-100 hover:shadow-md transition">
          <div class="w-12 h-12 bg-emerald-50 text-emerald-600 rounded-xl flex items-center justify-center mb-6">
            <ShieldCheck :size="24" />
          </div>
          <h3 class="text-xl font-bold mb-3">Keamanan Data Terjamin</h3>
          <p class="text-slate-500 leading-relaxed">
            Sistem kami menggunakan enkripsi modern dan backup otomatis. Data kandang Anda aman 24/7.
          </p>
        </div>
        <div class="bg-white p-8 rounded-2xl shadow-sm border border-slate-100 hover:shadow-md transition">
          <div class="w-12 h-12 bg-orange-50 text-orange-600 rounded-xl flex items-center justify-center mb-6">
            <TrendingUp :size="24" />
          </div>
          <h3 class="text-xl font-bold mb-3">Laporan Keuangan</h3>
          <p class="text-slate-500 leading-relaxed">
            Pantau profitabilitas ternak Anda dengan grafik pemasukan dan pengeluaran yang real-time.
          </p>
        </div>
      </div>
    </section>
  </div>
</template>