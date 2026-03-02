import { createApp } from 'vue'
import { MotionPlugin } from '@vueuse/motion'
import './style.css'
import App from './App.vue'
import router from './router'

createApp(App).use(router).use(MotionPlugin).mount('#app')
