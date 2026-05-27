import request from '@/utils/request'

export function register(data) {
  return request({
    url: '/users/register',
    method: 'post',
    data
  })
}

export function login(data) {
  return request({
    url: '/users/login',
    method: 'post',
    data
  })
}

export function logout() {
  return request({
    url: '/users/logout',
    method: 'post'
  })
}

export function sendCode(data) {
  return request({
    url: '/users/sms-code',
    method: 'post',
    data
  })
}

export function loginByCode(data) {
  return request({
    url: '/users/login-code',
    method: 'post',
    data
  })
}

export function sendPasswordResetCode(data) {
  return request({
    url: '/users/password-reset/code',
    method: 'post',
    data
  })
}

export function forgetPassword(data) {
  return request({
    url: '/users/password-reset',
    method: 'post',
    data
  })
}

export function getUserInfo() {
  return request({
    url: '/users/me',
    method: 'get'
  })
}

export function updateUserInfo(data) {
  return request({
    url: '/users/me',
    method: 'put',
    data
  })
}

export function changePassword(data) {
  return request({
    url: '/users/me/password',
    method: 'put',
    data
  })
}
