/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        garden: {
          50: '#f0fdf4',
          100: '#dcfce7',
          200: '#bbf7d0',
          300: '#86efac',
          400: '#4ade80',
          500: '#22c55e',
          600: '#16a34a',
          700: '#15803d',
          800: '#166534',
          900: '#14532d',
        },
        earth: {
          50: '#fdf8f6',
          100: '#f5ebe4',
          200: '#e8d5c4',
          300: '#d4b89e',
          400: '#c19a74',
          500: '#a67c52',
          600: '#8b6340',
          700: '#6f4c32',
          800: '#573b28',
          900: '#402a1d',
        },
      },
    },
  },
  plugins: [],
}
