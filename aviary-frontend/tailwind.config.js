/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {},
  },
  // TAMBAHKAN INI:
  plugins: [
    require('tailwind-scrollbar')({ nocompatible: true }),
  ],
}