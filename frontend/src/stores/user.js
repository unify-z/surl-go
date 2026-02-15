import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useUserStore = defineStore('user', () => {
    const userId = ref(localStorage.getItem('user_id') || '')
    const username = ref(localStorage.getItem('username') || '')
    const isAdmin = ref(localStorage.getItem('is_admin') === 'true')
    const token = ref(localStorage.getItem('token') || '')

    const isLoggedIn = computed(() => !!token.value)

    function setUser(data) {
        userId.value = data.user_id
        username.value = data.username
        isAdmin.value = data.is_admin

        localStorage.setItem('user_id', data.user_id)
        localStorage.setItem('username', data.username)
        localStorage.setItem('is_admin', data.is_admin)
    }

    function setToken(newToken) {
        token.value = newToken
        localStorage.setItem('token', newToken)
    }

    function logout() {
        userId.value = ''
        username.value = ''
        isAdmin.value = false
        token.value = ''

        localStorage.removeItem('user_id')
        localStorage.removeItem('username')
        localStorage.removeItem('is_admin')
        localStorage.removeItem('token')
    }

    return {
        userId,
        username,
        isAdmin,
        token,
        isLoggedIn,
        setUser,
        setToken,
        logout
    }
})
