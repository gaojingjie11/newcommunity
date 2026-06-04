import request from '@/utils/request'

export function recharge(amount, payType = 'alipay') {
  return request({
    url: '/mall/wallet/recharge',
    method: 'post',
    data: { amount: Number(amount), pay_type: payType }
  })
}

export function transfer(data) {
  return request({
    url: '/mall/wallet/transfer',
    method: 'post',
    data: {
      ...data,
      idempotency_key: data.idempotency_key || `transfer-${Date.now()}-${Math.random().toString(16).slice(2)}`
    }
  })
}

export function getTransactionList(params) {
  return request({
    url: '/mall/wallet/transactions',
    method: 'get',
    params
  })
}

export function getWalletBalance() {
  return request({
    url: '/mall/wallet/balance',
    method: 'get'
  })
}

export function getPropertyFeeList(params) {
  return request({
    url: '/community/property-fees',
    method: 'get',
    params
  })
}
