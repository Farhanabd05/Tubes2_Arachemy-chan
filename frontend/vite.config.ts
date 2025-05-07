// vite.config.ts
import { defineConfig } from 'vite';

export default defineConfig({
  server: {
    proxy: {
      '/search': 'http://localhost:8080',
    }
  }
});
