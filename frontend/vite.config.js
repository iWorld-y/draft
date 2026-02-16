import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig({
  plugins: [react()],
  server: {
    proxy: {
      '/articles': 'http://localhost:8124',
      '/helloworld': 'http://localhost:8124', // Just in case
    }
  }
})
