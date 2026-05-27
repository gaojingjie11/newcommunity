import axios from 'axios'

const request = axios.create({
  baseURL: '/api',
  timeout: 10000
})

function createBusinessError(res) {
  const code = Number(res?.code || 0)
  const message = res?.message || res?.msg || '请求失败'
  const error = new Error(message)
  error.response = {
    status: code === 401 ? 401 : 200,
    data: { ...res, msg: res?.msg || message }
  }
  error.code = String(code || '')
  return error
}

function clearAuth() {
  localStorage.removeItem('token')
  localStorage.removeItem('userInfo')
}

function isAuthRequest(config = {}) {
  const url = String(config.url || '')
  return (
    url.includes('/users/login') ||
    url.includes('/users/register') ||
    url.includes('/users/password-reset')
  )
}

function redirectToLogin() {
  const current = `${window.location.pathname}${window.location.search}${window.location.hash}`
  if (window.location.pathname === '/login') return

  const loginUrl = current && current !== '/'
    ? `/login?redirect=${encodeURIComponent(current)}`
    : '/login'
  window.location.replace(loginUrl)
}

request.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => Promise.reject(error)
)

request.interceptors.response.use(
  (response) => {
    const res = response.data
    if (res.code === 0 || res.code === 200) {
      return res.data
    }

    if (res.code === 401) {
      clearAuth()
      if (!isAuthRequest(response.config)) {
        redirectToLogin()
      }
      return Promise.reject(createBusinessError(res))
    }

    return Promise.reject(createBusinessError(res))
  },
  (error) => {
    const data = error.response?.data
    const message = data?.message || data?.msg
    if (message) {
      error.message = message
      error.response.data = { ...data, msg: data.msg || message }
    }
    if (error.response?.status === 401) {
      clearAuth()
      if (!isAuthRequest(error.config)) {
        redirectToLogin()
      }
    }
    return Promise.reject(error)
  }
)

export default request
