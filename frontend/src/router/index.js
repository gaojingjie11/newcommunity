import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '@/stores/user'

const routes = [
  { path: '/', redirect: '/home' },
  { path: '/data', name: 'DataScreen', component: () => import('@/views/admin/DataScreen.vue'), meta: { requiresAuth: true, requiresAdmin: true, roles: ['admin', 'property'] } },
  { path: '/login', name: 'Login', component: () => import('@/views/auth/Login.vue'), meta: { hideNav: true } },
  { path: '/register', name: 'Register', component: () => import('@/views/auth/Register.vue'), meta: { hideNav: true } },
  { path: '/home', name: 'Home', component: () => import('@/views/home/Index.vue'), meta: { requiresAuth: false } },
  { path: '/mall', name: 'Mall', component: () => import('@/views/mall/Index.vue'), meta: { requiresAuth: false } },
  { path: '/product/:id', name: 'ProductDetail', component: () => import('@/views/mall/ProductDetail.vue'), meta: { requiresAuth: false } },
  { path: '/cart', name: 'Cart', component: () => import('@/views/mall/Cart.vue'), meta: { requiresAuth: true } },
  { path: '/order', name: 'Order', component: () => import('@/views/order/Index.vue'), meta: { requiresAuth: true } },
  { path: '/order/:id', name: 'OrderDetail', component: () => import('@/views/order/Detail.vue'), meta: { requiresAuth: true } },
  { path: '/order/create', name: 'CreateOrder', component: () => import('@/views/order/Create.vue'), meta: { requiresAuth: true } },
  { path: '/service', name: 'Service', component: () => import('@/views/service/Index.vue'), meta: { requiresAuth: false } },
  { path: '/service/notice', name: 'Notice', component: () => import('@/views/service/Notice.vue'), meta: { requiresAuth: false } },
  { path: '/service/repair', name: 'Repair', component: () => import('@/views/service/Repair.vue'), meta: { requiresAuth: true } },
  { path: '/service/visitor', name: 'Visitor', component: () => import('@/views/service/Visitor.vue'), meta: { requiresAuth: true } },
  { path: '/service/parking', name: 'Parking', component: () => import('@/views/service/Parking.vue'), meta: { requiresAuth: true } },
  { path: '/service/property', name: 'PropertyFee', component: () => import('@/views/service/PropertyFee.vue'), meta: { requiresAuth: true } },
  { path: '/service/green-points', name: 'GreenPoints', component: () => import('@/views/service/GreenPoints.vue'), meta: { requiresAuth: true } },
  { path: '/service/community-chat', name: 'CommunityChat', component: () => import('@/views/service/CommunityChat.vue'), meta: { requiresAuth: true } },
  { path: '/profile', name: 'Profile', component: () => import('@/views/profile/Index.vue'), meta: { requiresAuth: true } },
  { path: '/admin', name: 'Admin', component: () => import('@/views/admin/Index.vue'), meta: { requiresAuth: true, requiresAdmin: true, roles: ['admin', 'store', 'property'] } },
  { path: '/user/favorites', name: 'FavoriteList', component: () => import('@/views/user/FavoriteList.vue'), meta: { requiresAuth: true } },
  { path: '/user/transactions', name: 'TransactionList', component: () => import('@/views/user/TransactionList.vue'), meta: { requiresAuth: true } },
  { path: '/admin/users', name: 'AdminUserList', component: () => import('@/views/admin/UserList.vue'), meta: { requiresAuth: true, requiresAdmin: true, roles: ['admin'] } },
  { path: '/admin/products', name: 'AdminProductList', component: () => import('@/views/admin/ProductList.vue'), meta: { requiresAuth: true, requiresAdmin: true, roles: ['admin', 'store'] } },
  { path: '/admin/orders', name: 'AdminOrderList', component: () => import('@/views/admin/OrderList.vue'), meta: { requiresAuth: true, requiresAdmin: true, roles: ['admin', 'store'] } },
  { path: '/admin/stores', name: 'AdminStoreList', component: () => import('@/views/admin/StoreList.vue'), meta: { requiresAuth: true, requiresAdmin: true, roles: ['admin', 'store'] } },
  { path: '/admin/notices', name: 'AdminNoticeList', component: () => import('@/views/admin/NoticeList.vue'), meta: { requiresAuth: true, requiresAdmin: true, roles: ['admin', 'property'] } },
  { path: '/admin/repairs', name: 'AdminRepairList', component: () => import('@/views/admin/RepairList.vue'), meta: { requiresAuth: true, requiresAdmin: true, roles: ['admin', 'property'] } },
  { path: '/admin/visitors', name: 'AdminVisitorList', component: () => import('@/views/admin/VisitorList.vue'), meta: { requiresAuth: true, requiresAdmin: true, roles: ['admin', 'property'] } },
  { path: '/admin/parking', name: 'AdminParkingList', component: () => import('@/views/admin/ParkingList.vue'), meta: { requiresAuth: true, requiresAdmin: true, roles: ['admin', 'property'] } },
  { path: '/admin/property-fee', name: 'AdminPropertyFeeList', component: () => import('@/views/admin/PropertyFeeList.vue'), meta: { requiresAuth: true, requiresAdmin: true, roles: ['admin', 'property'] } },
  { path: '/admin/ai-report', name: 'AdminAIReport', component: () => import('@/views/admin/AIReport.vue'), meta: { requiresAuth: true, requiresAdmin: true, roles: ['admin'] } },
  { path: '/chat', name: 'Chat', component: () => import('@/views/chat/Index.vue'), meta: { requiresAuth: true } }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, from, next) => {
  const userStore = useUserStore()

  if (to.meta.requiresAuth && !userStore.isLoggedIn) {
    next('/login')
    return
  }

  if (to.meta.requiresAdmin) {
    const userRole = userStore.userInfo.role
    const allowedRoles = to.meta.roles || ['admin', 'store', 'property']
    if (!allowedRoles.includes(userRole)) {
      window.alert('No permission to access this page')
      next('/')
      return
    }
  }

  next()
})

export default router
