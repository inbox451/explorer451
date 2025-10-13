// https://nuxt.com/docs/api/configuration/nuxt-config
import tailwindcss from '@tailwindcss/vite'

// =========================
// Nuxt Configuration
// =========================
export default defineNuxtConfig({
  // General
  compatibilityDate: '2025-05-15',
  devtools: { enabled: true },
  ssr: false,

  // CSS
  css: ['~/assets/css/main.css'],

  // Modules
  modules: ['@nuxt/eslint', '@nuxt/ui', 'shadcn-nuxt', '@pinia/nuxt'],

  // Runtime Config
  runtimeConfig: {
    public: {
      apiEndpoint: process.env.API_BASE_URL || 'http://localhost:8080/api'
    }
  },

  // Vite
  vite: {
    plugins: [
      tailwindcss(),
    ],
  },

  // Typescript
  typescript: {
    // typeCheck: true
  },

  // UI: shadcn
  shadcn: {
    /**
     * Prefix for all the imported component
     */
    prefix: '',
    /**
     * Directory that the component lives in.
     * @default "./components/ui"
     */
    componentDir: './components/ui',
  },
})