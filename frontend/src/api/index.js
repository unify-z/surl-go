import axios from 'axios'
import { useUserStore } from '../stores/user'

const api = axios.create({
    baseURL: '/api',
    timeout: 10000
})

api.interceptors.request.use(
    config => {
        const userStore = useUserStore()
        if (userStore.token) {
            config.headers.Authorization = `Bearer ${userStore.token}`
        }
        return config
    },
    error => Promise.reject(error)
)

api.interceptors.response.use(
    response => response.data,
    error => {

        return Promise.reject(error)
    }
)

export default api
