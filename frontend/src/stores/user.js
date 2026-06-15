import { defineStore } from 'pinia'
import {
  login as apiLogin,
  loginByCode as apiLoginByCode,
  register as apiRegister,
  logout as apiLogout,
  getUserInfo
} from '@/api/auth'
import { getWalletBalance } from '@/api/finance'

export const useUserStore = defineStore('user', {
  state: () => ({
    token: localStorage.getItem('token') || '',
    userInfo: JSON.parse(localStorage.getItem('userInfo') || '{}'),
    permissions: JSON.parse(localStorage.getItem('permissions') || '[]'),
    isLoggedIn: !!localStorage.getItem('token'),
    isInfoFetched: false
  }),

  actions: {
    async login(data) {
      const res = await apiLogin(data)
      this.token = res.token
      this.userInfo = res.user_info
      this.permissions = res.user_info?.permissions || []
      this.isLoggedIn = true
      this.isInfoFetched = true
      localStorage.setItem('token', res.token)
      localStorage.setItem('userInfo', JSON.stringify(res.user_info))
      localStorage.setItem('permissions', JSON.stringify(this.permissions))
      await this.refreshWalletBalance()
      return res
    },

    async loginByCode(data) {
      const res = await apiLoginByCode(data)
      this.token = res.token
      this.userInfo = res.user_info
      this.permissions = res.user_info?.permissions || []
      this.isLoggedIn = true
      this.isInfoFetched = true
      localStorage.setItem('token', res.token)
      localStorage.setItem('userInfo', JSON.stringify(res.user_info))
      localStorage.setItem('permissions', JSON.stringify(this.permissions))
      await this.refreshWalletBalance()
      return res
    },

    async register(data) {
      return apiRegister(data)
    },

    async fetchUserInfo() {
      try {
        const userInfo = await getUserInfo()
        this.userInfo = userInfo || {}
        this.permissions = userInfo?.permissions || []
        this.isInfoFetched = true
        await this.refreshWalletBalance()
        localStorage.setItem('userInfo', JSON.stringify(this.userInfo))
        localStorage.setItem('permissions', JSON.stringify(this.permissions))
        return this.userInfo
      } catch (error) {
        console.error('fetchUserInfo failed', error)
        return null
      }
    },

    async refreshWalletBalance() {
      if (!this.token) return
      try {
        const wallet = await getWalletBalance()
        if (wallet && wallet.balance !== undefined) {
          this.userInfo = {
            ...this.userInfo,
            balance: wallet.balance
          }
          localStorage.setItem('userInfo', JSON.stringify(this.userInfo))
        }
      } catch (error) {
        console.warn('refreshWalletBalance failed', error)
      }
    },

    async logout() {
      try {
        await apiLogout()
      } catch (error) {
        console.warn('logout request failed', error)
      } finally {
        this.token = ''
        this.userInfo = {}
        this.permissions = []
        this.isLoggedIn = false
        this.isInfoFetched = false
        localStorage.removeItem('token')
        localStorage.removeItem('userInfo')
        localStorage.removeItem('permissions')
      }
    }
  }
})
