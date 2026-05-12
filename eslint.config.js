import eslint from '@eslint/js'
import tseslint from 'typescript-eslint'
import pluginVue from 'eslint-plugin-vue'
import globals from 'globals'

export default tseslint.config(
  // --- 共通の無視パターン ---
  {
    ignores: ['**/dist/**', '**/node_modules/**', 'packages/backend/**', '.history/**'],
  },

  // --- JS 推奨ルール ---
  eslint.configs.recommended,

  // --- TypeScript 推奨ルール ---
  ...tseslint.configs.recommended,

  // --- BFF 用設定 ---
  {
    files: ['packages/bff/src/**/*.ts'],
    languageOptions: {
      globals: {
        ...globals.node,
      },
    },
  },

  // --- Frontend 用設定 (Vue) ---
  ...pluginVue.configs['flat/essential'],
  {
    files: ['packages/frontend/src/**/*.vue'],
    languageOptions: {
      globals: {
        ...globals.browser,
      },
      parserOptions: {
        parser: tseslint.parser,
      },
    },
    rules: {
      'vue/multi-word-component-names': 'off',
    },
  },
  {
    files: ['packages/frontend/src/**/*.ts'],
    languageOptions: {
      globals: {
        ...globals.browser,
      },
    },
  },

  // --- 共通ルール上書き ---
  {
    rules: {
      '@typescript-eslint/no-unused-vars': [
        'error',
        { argsIgnorePattern: '^_', varsIgnorePattern: '^_', caughtErrorsIgnorePattern: '^_' },
      ],
      'preserve-caught-error': 'off',
    },
  }
)
