/*
 * @Author: Liu Jiarong
 * @Date: 2024-06-28 17:13:04
 * @LastEditors: Liu Jiarong
 * @LastEditTime: 2024-11-20 13:11:34
 * @FilePath: /chatnio/app/vite.config.ts
 * @Description: 
 * 
 * Copyright (c) 2024 by ${git_name_email}, All Rights Reserved. 
 */
import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react-swc'
import path from "path"
import { createHtmlPlugin } from 'vite-plugin-html'
import { createTranslationPlugin } from "./src/translator";

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    react(),
    createHtmlPlugin({
      minify: true,
    }),
    createTranslationPlugin(),
  ],
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src"),
    },
  },
  css: {
    preprocessorOptions: {
      less: {
        javascriptEnabled: true,
      }
    }
  },
  build: {
    manifest: true,
    chunkSizeWarningLimit: 2048,
    rollupOptions: {
      output: {
        entryFileNames: `assets/[name].[hash].js`,
        chunkFileNames: `assets/[name].[hash].js`,
      },
    },
  },
  server: {
    proxy: {
      "/api": {
        target: "http://localhost:8094",
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api/, ""),
        ws: true,
      },
      "/v1": {
        target: "http://localhost:8094",
        changeOrigin: true,
      },
      "/quota": {
        target: "http://localhost:8094",
        changeOrigin: true,
      },
      "/subscription": {
        target: "http://localhost:8094",
        changeOrigin: true,
      }
    }
  }
});
