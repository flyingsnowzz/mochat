import { defineConfig } from 'vitest/config'
import path from 'path'

export default defineConfig({
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src')
    }
  },
  test: {
    environment: 'node',
    globals: true,
    include: ['src/api/__tests__/**/*.test.ts'],
    exclude: ['node_modules', 'dist'],
    coverage: {
      reporter: ['text', 'json', 'html'],
      include: ['src/api/**/*.ts'],
      exclude: ['src/api/__tests__/**']
    }
  }
})
