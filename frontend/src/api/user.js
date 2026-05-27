import request from '@/utils/request'

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

export function registerFace(data) {
    return request({
        url: '/users/me/face',
        method: 'post',
        data
    })
}
