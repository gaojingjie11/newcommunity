
import request from '@/utils/request'

export function getCommentList(params) {
    return request({
        url: '/mall/comments',
        method: 'get',
        params
    })
}

export function createComment(data) {
    return request({
        url: '/mall/comments',
        method: 'post',
        data
    })
}
