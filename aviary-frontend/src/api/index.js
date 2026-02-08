import axios from 'axios';

const Api = axios.create({
    baseURL: 'http://localhost:8080/api', 
    withCredentials: true, 
});

// INTERCEPTOR RESPONSE (RAHASIA AGAR TIDAK LOGOUT SENDIRI)
// Fungsi ini otomatis jalan SETELAH menerima jawaban dari backend
Api.interceptors.response.use(
  (response) => {
    return response; // Jika sukses (200 OK), lanjut.
  },
  async (error) => {
    const originalRequest = error.config;

    // Cek: Apakah errornya 401 (Unauthorized) DAN belum pernah dicoba ulang?
    if (error.response && error.response.status === 401 && !originalRequest._retry) {
      
      originalRequest._retry = true;

      try {
        // Minta Backend perbarui cookie (Browser otomatis bawa refresh_token)
        await Api.post('/auth/refresh');

        // Jika sukses refresh, panggil ulang request awal yang tadi gagal
        return Api(originalRequest);

      } catch (refreshError) {
        // Jika Refresh gagal (berarti sesi 7 hari habis)
        console.error("Sesi habis, login ulang.");
        // Redirect paksa ke login
        window.location.href = '/login'; 
        return Promise.reject(refreshError);
      }
    }

    return Promise.reject(error);
  }
);

export default Api;