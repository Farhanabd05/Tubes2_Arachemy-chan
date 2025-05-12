import { defineConfig } from 'vite';
import dns from 'node:dns';

dns.setDefaultResultOrder('verbatim');

export default defineConfig({
  server: {
    host: '0.0.0.0',
    port: 5173,
    strictPort: true,
  },
});