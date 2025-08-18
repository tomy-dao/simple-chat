import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import tailwindcss from '@tailwindcss/vite'
import path from 'path'
import { fileURLToPath } from 'url'
import { loadEnv } from 'vite'

const __dirname = path.dirname(fileURLToPath(import.meta.url))

// https://vite.dev/config/
export default defineConfig((mode) => {
  // eslint-disable-next-line no-undef
  const env = loadEnv(mode, process.cwd(), '')
  
  return {
    define: Object.fromEntries(
      Object.entries(env).map(([key, val]) => [
        `import.meta.env.${key}`, JSON.stringify(val)
      ])
    ),
    plugins: [react(), tailwindcss()],
    resolve: {
      alias: {
        "@": path.resolve(__dirname, "./src"),
      },
    },
    build: {
      assetsInlineLimit: 0 // Disable inlining for stricter CSP
    },
    server: {
      headers: {
        'Content-Security-Policy': `
          default-src 'self';
          script-src 'self' 'unsafe-eval';
          style-src 'self' 'unsafe-inline';
          connect-src 'self' ws: wss: ${env.VITE_SOCKET_BASE_URL};
        `.replace(/\s{2,}/g, ' ').trim()
      }
    }
  }
})
