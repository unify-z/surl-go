import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'
import './style.css'
import { useConfigStore } from './stores/config'

const app = createApp(App)
const pinia = createPinia()

app.use(pinia)
app.use(router)

// Pre‑load configuration as early as possible
const configStore = useConfigStore()
configStore.loadConfig()

app.mount('#app')
