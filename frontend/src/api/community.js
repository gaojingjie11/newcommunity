import request from '@/utils/request'

export function getCommunityMessages(params) {
  return request({
    url: '/community/messages',
    method: 'get',
    params
  })
}

export function sendCommunityMessage(data) {
  return request({
    url: '/community/message',
    method: 'post',
    data
  })
}
