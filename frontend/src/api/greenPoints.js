import request from '@/utils/request'

export function uploadGarbageImage(formData) {
  return request({
    url: '/green-points/upload-garbage',
    method: 'post',
    timeout: 120000,
    data: formData,
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
}

export function getGreenPointsLeaderboard(params) {
  return request({
    url: '/green-points/leaderboard',
    method: 'get',
    params
  })
}
