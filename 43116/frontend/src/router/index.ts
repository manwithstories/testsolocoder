import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router'
import { useUserStore } from '@/stores/user'

const routes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/Login.vue'),
    meta: { requiresAuth: false }
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('@/views/Register.vue'),
    meta: { requiresAuth: false }
  },
  {
    path: '/',
    component: () => import('@/layouts/DefaultLayout.vue'),
    meta: { requiresAuth: true },
    children: [
      {
        path: '',
        name: 'Home',
        component: () => import('@/views/Home.vue')
      },
      {
        path: 'cars',
        name: 'Cars',
        component: () => import('@/views/Cars.vue')
      },
      {
        path: 'cars/:id',
        name: 'CarDetail',
        component: () => import('@/views/CarDetail.vue')
      },
      {
        path: 'booking/:id',
        name: 'Booking',
        component: () => import('@/views/Booking.vue')
      },
      {
        path: 'my-bookings',
        name: 'MyBookings',
        component: () => import('@/views/MyBookings.vue')
      },
      {
        path: 'my-orders',
        name: 'MyOrders',
        component: () => import('@/views/MyOrders.vue')
      },
      {
        path: 'profile',
        name: 'Profile',
        component: () => import('@/views/Profile.vue')
      },
      {
        path: 'messages',
        name: 'Messages',
        component: () => import('@/views/Messages.vue')
      }
    ]
  },
  {
    path: '/admin',
    component: () => import('@/layouts/AdminLayout.vue'),
    meta: { requiresAuth: true, requiresAdmin: true },
    children: [
      {
        path: '',
        name: 'Dashboard',
        component: () => import('@/views/admin/Dashboard.vue')
      },
      {
        path: 'users',
        name: 'AdminUsers',
        component: () => import('@/views/admin/Users.vue')
      },
      {
        path: 'cars',
        name: 'AdminCars',
        component: () => import('@/views/admin/Cars.vue')
      },
      {
        path: 'stores',
        name: 'AdminStores',
        component: () => import('@/views/admin/Stores.vue')
      },
      {
        path: 'cities',
        name: 'AdminCities',
        component: () => import('@/views/admin/Cities.vue')
      },
      {
        path: 'bookings',
        name: 'AdminBookings',
        component: () => import('@/views/admin/Bookings.vue')
      },
      {
        path: 'orders',
        name: 'AdminOrders',
        component: () => import('@/views/admin/Orders.vue')
      },
      {
        path: 'reviews',
        name: 'AdminReviews',
        component: () => import('@/views/admin/Reviews.vue')
      },
      {
        path: 'maintenance',
        name: 'AdminMaintenance',
        component: () => import('@/views/admin/Maintenance.vue')
      },
      {
        path: 'promos',
        name: 'AdminPromos',
        component: () => import('@/views/admin/Promos.vue')
      },
      {
        path: 'pricing',
        name: 'AdminPricing',
        component: () => import('@/views/admin/Pricing.vue')
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, from, next) => {
  const userStore = useUserStore()

  if (to.meta.requiresAuth && !userStore.isLoggedIn) {
    next('/login')
  } else if (to.meta.requiresAdmin && !userStore.isAdmin) {
    next('/')
  } else if ((to.path === '/login' || to.path === '/register') && userStore.isLoggedIn) {
    next(userStore.isAdmin ? '/admin' : '/')
  } else {
    next()
  }
})

export default router
