import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue({
    refTransform: true,
    reactivityTransform: true
  }),],

  server:{

    proxy:{
      '/user':{
        target: "http://localhost:3000/"
      }
    }
  },

  base: './',
})
