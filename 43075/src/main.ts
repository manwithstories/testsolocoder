import { createApp } from 'vue'
import { createPinia } from 'pinia'
import './style.css'
import App from './App.vue'
import router from './router'
import { useSettingsStore } from './stores/settings'
import { useCategoriesStore } from './stores/categories'
import { useRecordsStore } from './stores/records'

const app = createApp(App)
const pinia = createPinia()

app.use(pinia)
app.use(router)

app.mount('#app')

const settingsStore = useSettingsStore()
const categoriesStore = useCategoriesStore()
const recordsStore = useRecordsStore()

settingsStore.loadSettings()
categoriesStore.loadCategories()
recordsStore.loadRecords()
