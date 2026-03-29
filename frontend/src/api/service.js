import request from '@/utils/request'

export function getNoticeList(params) {
  return request({
    url: '/notices',
    method: 'get',
    params
  })
}

export function getNoticeDetail(id) {
  return request({
    url: `/notice/${id}`,
    method: 'get'
  })
}

export function readNotice(id) {
  return request({
    url: `/notice/read/${id}`,
    method: 'post'
  })
}

export function createRepair(data) {
  return request({
    url: '/repair/create',
    method: 'post',
    data
  })
}

export function getRepairList(params) {
  return request({
    url: '/repair/list',
    method: 'get',
    params
  })
}

export function createVisitor(data) {
  return request({
    url: '/visitor/create',
    method: 'post',
    data
  })
}

export function getVisitorList(params) {
  return request({
    url: '/visitor/list',
    method: 'get',
    params
  })
}

export function getMyParking() {
  return request({
    url: '/parking/my',
    method: 'get'
  })
}

export function bindCar(data) {
  return request({
    url: '/parking/bind',
    method: 'post',
    data
  })
}

export function getPropertyFeeList(params) {
  return request({
    url: '/property/list',
    method: 'get',
    params
  })
}

export function payPropertyFee(feeId, payload = {}) {
  return request({
    url: '/finance/pay',
    method: 'post',
    data: {
      business_id: feeId,
      business_type: 2,
      pay_type: payload.pay_type || 'password',
      password: payload.password || '',
      face_image_url: payload.face_image_url || ''
    }
  })
}

export function getStoreList() {
  return request({
    url: '/stores',
    method: 'get'
  })
}
