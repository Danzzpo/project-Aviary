import { createApp } from 'vue'
import { createPinia } from 'pinia' 
import App from './App.vue'
import router from './router' // Asumsi router sudah ada

import './assets/main.css'

const app = createApp(App)

app.use(createPinia()) // <--- Pakai ini
app.use(router)

app.mount('#app')