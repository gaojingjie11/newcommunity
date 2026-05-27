import request from '@/utils/request'

export function addFavorite(productId) {
    return request({
        url: '/mall/favorites',
        method: 'post',
        data: { product_id: productId }
    })
}

export function deleteFavorite(productId) {
    return request({
        url: `/mall/favorites/${productId}`,
        method: 'delete'
    })
}

export function getFavoriteList() {
    return request({
        url: '/mall/favorites',
        method: 'get'
    })
}

export function checkFavorite(productId) {
    return request({
        url: `/mall/favorites/check/${productId}`,
        method: 'get'
    })
}
