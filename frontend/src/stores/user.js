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
    isLoggedIn: !!localStorage.getItem('token')
  }),

  actions: {
    async login(data) {
      const res = await apiLogin(data)
      this.token = res.token
      this.userInfo = res.user_info
      this.isLoggedIn = true
      localStorage.setItem('token', res.token)
      localStorage.setItem('userInfo', JSON.stringify(res.user_info))
      await this.refreshWalletBalance()
      return res
    },

    async loginByCode(data) {
      const res = await apiLoginByCode(data)
      this.token = res.token
      this.userInfo = res.user_info
      this.isLoggedIn = true
      localStorage.setItem('token', res.token)
      localStorage.setItem('userInfo', JSON.stringify(res.user_info))
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
        await this.refreshWalletBalance()
        localStorage.setItem('userInfo', JSON.stringify(this.userInfo))
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
        this.isLoggedIn = false
        localStorage.removeItem('token')
        localStorage.removeItem('userInfo')
      }
    }
  }
})
