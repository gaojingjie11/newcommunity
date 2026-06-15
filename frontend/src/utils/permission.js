import { useUserStore } from '@/stores/user'

export function hasPermission(perm) {
  const userStore = useUserStore()
  // System admin role bypasses all permission checks
  if (userStore.userInfo?.role === 'admin') {
    return true
  }
  return userStore.permissions ? userStore.permissions.includes(perm) : false
}
