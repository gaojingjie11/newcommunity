import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { fileURLToPath, URL } from 'node:url'

// https://vite.dev/config/
export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    }
  },
  server: {
    port: 5173,
    proxy: {
      '/api': {
          // target: 'http://42.193.104.173:8082', // 连接服务器 Nginx 端口
        target: 'http://127.0.0.1:8080', // 连接本地后端 (dev.yaml 配置为 8080)
        changeOrigin: true
      }
    }
  }
})
