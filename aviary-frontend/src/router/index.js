import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'

// Views
import HomeView from '../views/HomeView.vue'
import LoginView from '../views/LoginView.vue'
import RegisterView from '../views/RegisterView.vue'
import DashboardLayout from '../layouts/DashboardLayout.vue'
import DashboardHome from '../views/dashboard/DashboardHome.vue'
import BirdListView from '../views/dashboard/BirdListView.vue'
import PairingView from '../views/dashboard/PairingView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView
    },
    {
      path: '/login',
      name: 'login',
      component: LoginView,
      meta: { guestOnly: true }
       // Hanya boleh diakses kalau BELUM login
    },
    {
      path: '/register',
      name: 'register',
      component: RegisterView,
      meta: { guestOnly: true } // Hanya boleh diakses kalau belum login
    },
    {
      path: '/dashboard',
      component: DashboardLayout, // Pakai Layout yang ada sidebarnya
      meta: { requiresAuth: true }, // WAJIB Login
      children: [
        {
          path: '', // Default: /dashboard
          name: 'dashboard',
          component: DashboardHome
        },
        {
          path: 'birds',  // Jadi URL-nya: /dashboard/birds
          name: 'bird-list',
          component: BirdListView
        },
        {
          path: 'pairs',
          name: 'pairing',
          component: PairingView
        }
      ]
    }
  ]
})

// --- NAVIGATION GUARD (SATPAM) ---
router.beforeEach((to, from, next) => {
  const authStore = useAuthStore()
  const isAuthenticated = authStore.isLoggedIn

  // 1. Cek jika halaman butuh login (Dashboard)
  if (to.meta.requiresAuth && !isAuthenticated) {
    next('/login') // Tendang ke login
  } 
  // 2. Cek jika halaman khusus tamu (Login), tapi user sudah login
  else if (to.meta.guestOnly && isAuthenticated) {
    next('/dashboard') // Lempar ke dashboard
  } 
  else {
    next() // Lanjut
  }
})

export default router