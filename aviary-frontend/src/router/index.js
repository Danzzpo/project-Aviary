import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'

// 1. IMPORT VIEWS (HALAMAN PUBLIK)
import HomeView from '../views/HomeView.vue'
import LoginView from '../views/LoginView.vue'
import RegisterView from '../views/RegisterView.vue'

// 2. IMPORT LAYOUT & DASHBOARD VIEWS
import DashboardLayout from '../layouts/DashboardLayout.vue'
import DashboardHome from '../views/dashboard/DashboardHome.vue'
import BirdListView from '../views/dashboard/BirdListView.vue'
import PairingView from '../views/dashboard/PairingView.vue'
import FinanceView from '../views/dashboard/FinanceView.vue'
import SettingsView from '../views/dashboard/SettingsView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    // --- HALAMAN PUBLIK ---
    {
      path: '/',
      name: 'home',
      component: HomeView
    },
    {
      path: '/login',
      name: 'login',
      component: LoginView,
      meta: { guestOnly: true } // Hanya untuk yang BELUM login
    },
    {
      path: '/register',
      name: 'register',
      component: RegisterView,
      meta: { guestOnly: true } // Hanya untuk yang BELUM login
    },

    // --- HALAMAN DASHBOARD (BUTUH LOGIN) ---
    {
      path: '/dashboard',
      component: DashboardLayout, // Parent Layout (Sidebar + Navbar)
      meta: { requiresAuth: true }, // Wajib Login
      children: [
        {
          path: '', // URL: /dashboard
          name: 'dashboard',
          component: DashboardHome
        },
        {
          path: 'birds', // URL: /dashboard/birds
          name: 'bird-list',
          component: BirdListView
        },
        {
          path: 'pairs', // URL: /dashboard/pairing
          name: 'pairing',
          component: PairingView
        },
        {
          path: 'finance', // URL: /dashboard/finance
          name: 'finance',
          component: FinanceView
        },
        {
          path: 'settings', 
          name: 'settings', 
          component: SettingsView 
        }
      ]
    },
    
    // Redirect jika halaman tidak ditemukan (404) ke Home
    {
      path: '/:pathMatch(.*)*',
      redirect: '/'
    }
  ]
})

// --- NAVIGATION GUARD (SATPAM) ---
router.beforeEach((to, from, next) => {
  const authStore = useAuthStore()
  
  // Cek apakah user sudah login (berdasarkan state Pinia)
  const isAuthenticated = authStore.isAuthenticated

  // 1. Jika halaman BUTUH login (requiresAuth) tapi user BELUM login
  if (to.meta.requiresAuth && !isAuthenticated) {
    next({ name: 'login' }) // Tendang ke login
  } 
  // 2. Jika halaman KHUSUS TAMU (guestOnly) tapi user SUDAH login
  else if (to.meta.guestOnly && isAuthenticated) {
    next({ name: 'dashboard' }) // Lempar ke dashboard
  } 
  // 3. Lolos seleksi, silakan masuk
  else {
    next()
  }
})

export default router