import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig({
  plugins: [react()],
  server: {
    host: '0.0.0.0',
    proxy: {
      '/articles': 'http://backend-dev:8000',
      '/helloworld': 'http://backend-dev:8000',
    }
  }
})
