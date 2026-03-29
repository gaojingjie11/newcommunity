import request from '@/utils/request'

export function getUserInfo() {
    return request({
        url: '/user/info',
        method: 'get'
    })
}

export function updateUserInfo(data) {
    return request({
        url: '/user/update',
        method: 'post',
        data
    })
}

export function changePassword(data) {
    return request({
        url: '/user/change_password',
        method: 'post',
        data
    })
}

export function registerFace(data) {
    return request({
        url: '/user/face/register',
        method: 'post',
        data
    })
}
