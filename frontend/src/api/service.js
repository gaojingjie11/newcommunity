import request from '@/utils/request'

export function getNoticeList(params) {
  return request({
    url: '/community/notices',
    method: 'get',
    params
  })
}

export function getNoticeDetail(id) {
  return request({
    url: `/community/notices/${id}`,
    method: 'get'
  })
}

export function readNotice(id) {
  return request({
    url: `/community/notices/${id}/read`,
    method: 'post'
  })
}

export function createWorkorder(data) {
  return request({
    url: '/workorders',
    method: 'post',
    data
  })
}

export function getWorkorderList(params) {
  return request({
    url: '/workorders',
    method: 'get',
    params
  })
}

export function createVisitor(data) {
  return request({
    url: '/community/visitors',
    method: 'post',
    data
  })
}

export function getVisitorList(params) {
  return request({
    url: '/community/visitors',
    method: 'get',
    params
  })
}

export function getMyParking() {
  return request({
    url: '/community/parking-spaces/my',
    method: 'get'
  })
}

export function bindCar(parkingId, data) {
  return request({
    url: `/community/parking-spaces/${parkingId}/plate`,
    method: 'put',
    data
  })
}

export function getPropertyFeeList(params) {
  return request({
    url: '/community/property-fees',
    method: 'get',
    params
  })
}

export function payPropertyFee(feeId, payload = {}) {
  const idempotencyKey = payload.idempotency_key || payload.idempotencyKey ||
    `property-fee-${feeId}-${Date.now()}-${Math.random().toString(36).slice(2, 10)}`

  return request({
    url: `/community/property-fees/${feeId}/pay`,
    method: 'post',
    data: {
      pay_type: payload.pay_type || 'password',
      password: payload.password || '',
      face_image_url: payload.face_image_url || '',
      idempotency_key: idempotencyKey
    }
  })
}

export function getStoreList() {
  return request({
    url: '/mall/stores',
    method: 'get'
  })
}
