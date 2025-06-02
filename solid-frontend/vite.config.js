import { defineConfig } from 'vite';

export default defineConfig({
  server: {
    port: 5273,
    proxy: {
      '/api': 'http://backend:8100',
    },
  },
});
