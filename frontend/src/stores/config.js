import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '../api'

export const useConfigStore = defineStore('config', () => {
    const siteName = ref('短链接服务')
    const allowRegistration = ref(true)
    const allowGuestToCreateUrl = ref(false)
    const loaded = ref(false)

    async function loadConfig() {
        if (loaded.value) return

        try {
            const res = await api.get('/config')
            if (res.code === 200) {
                siteName.value = res.data.site_name
                allowRegistration.value = res.data.allow_registration ?? true
                allowGuestToCreateUrl.value = res.data.allow_guest_to_create_url ?? false
                loaded.value = true
            }
        } catch (error) {
            console.error('Failed to load config:', error)
        }
    }

    return {
        siteName,
        allowRegistration,
        allowGuestToCreateUrl,
        loaded,
        loadConfig
    }
})
