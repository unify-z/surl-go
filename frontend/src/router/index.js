import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '../stores/user'
import { useConfigStore } from '../stores/config'

const routes = [
    {
        path: '/',
        redirect: '/dashboard'
    },
    {
        path: '/login',
        name: 'Login',
        component: () => import('../views/Login.vue'),
        meta: { requiresGuest: true }
    },
    {
        path: '/dashboard',
        name: 'Dashboard',
        component: () => import('../views/Dashboard.vue'),
        meta: { requiresAuth: true, allowGuest: true }
    },
    {
        path: '/admin',
        name: 'Admin',
        component: () => import('../views/Admin.vue'),
        meta: { requiresAuth: true, requiresAdmin: true }
    }
]

const router = createRouter({
    history: createWebHistory(),
    routes
})

router.beforeEach(async (to, from, next) => {
    const userStore = useUserStore()
    const configStore = useConfigStore()

    if (!configStore.loaded) {
        await configStore.loadConfig()
    }

    const { isLoggedIn, isAdmin } = userStore

    if (to.meta.requiresAuth) {
        if (isLoggedIn) {
            if (to.meta.requiresAdmin && !isAdmin) {
                return next('/dashboard')
            }
            return next()
        } else {
            if (to.meta.allowGuest && configStore.allowGuestToCreateUrl) {
                return next()
            }
            return next('/login')
        }
    }

    if (to.meta.requiresGuest && isLoggedIn) {
        return next('/dashboard')
    }

    next()
})

export default router