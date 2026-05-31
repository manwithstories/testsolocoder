/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        coffee: {
          50: '#fdf8f3',
          100: '#f7e9d7',
          200: '#ecd1a8',
          300: '#deb573',
          400: '#cf9647',
          500: '#c07a2a',
          600: '#a85f22',
          700: '#8b461f',
          800: '#6b341c',
          900: '#4a2316',
        },
      },
    },
  },
  plugins: [],
}
