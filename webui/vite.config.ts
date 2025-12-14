import { defineConfig, loadEnv } from 'vite'
import { svelte } from '@sveltejs/vite-plugin-svelte'
import { VitePWA } from 'vite-plugin-pwa'

// https://vite.dev/config/
export default defineConfig(({ mode }) => {
  // Load env file based on `mode` in the current working directory.
  const env = loadEnv(mode, process.cwd(), '')
  
  return {
    plugins: [
      svelte(),
      VitePWA({
        registerType: 'autoUpdate',
        includeAssets: ['favicon.ico', 'apple-touch-icon.png'],
        manifest: {
          name: 'Docker Simple Panel',
          short_name: 'Docker Panel',
          description: 'A simple web application for managing Docker containers',
          theme_color: '#0db7ed',
          background_color: '#ffffff',
          display: 'standalone',
          scope: '/',
          start_url: '/',
          icons: [
            {
              src: 'icon-192x192.png',
              sizes: '192x192',
              type: 'image/png',
              purpose: 'any'
            },
            {
              src: 'icon-192x192.png',
              sizes: '192x192',
              type: 'image/png',
              purpose: 'maskable'
            },
            {
              src: 'icon-512x512.png',
              sizes: '512x512',
              type: 'image/png',
              purpose: 'any'
            },
            {
              src: 'icon-512x512.png',
              sizes: '512x512',
              type: 'image/png',
              purpose: 'maskable'
            }
          ]
        },
        workbox: {
          // Cache API responses
          runtimeCaching: [
            {
              urlPattern: /^https?:\/\/.*\/api\/.*/i,
              handler: 'StaleWhileRevalidate',
              options: {
                cacheName: 'api-cache',
                expiration: {
                  maxEntries: 50,
                  maxAgeSeconds: 30 * 60 // 30 minutes
                },
                cacheableResponse: {
                  statuses: [0, 200]
                }
              }
            }
          ],
          // Don't cache these paths
          navigateFallbackDenylist: [/^\/api/]
        }
      })
    ],
    server: {
      proxy: {
        // Proxy API requests to backend server in development
        '/api': {
          target: env.VITE_API_PROXY_TARGET || 'http://localhost:3000',
          changeOrigin: true,
          secure: false
        }
      }
    }
  }
})
