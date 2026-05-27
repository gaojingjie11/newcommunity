import axios from 'axios'

export function sendChat(data) {
  return axios({
    baseURL: '',
    url: '/agent/chat',
    method: 'post',
    data: {
      user_id: data.user_id,
      message: data.content || data.message || ''
    },
    headers: {
      Authorization: localStorage.getItem('token') ? `Bearer ${localStorage.getItem('token')}` : ''
    },
    timeout: 60000
  }).then((res) => res.data?.data || res.data)
}

export function getChatHistory(params) {
  void params
  return Promise.resolve({ list: [] })
}
