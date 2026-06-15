import request from '@/utils/request'

export function getConversations() {
  return request({
    url: '/agent/conversations',
    method: 'get'
  })
}

export function createConversation(data) {
  return request({
    url: '/agent/conversations',
    method: 'post',
    data
  })
}

export function deleteConversation(id) {
  return request({
    url: `/agent/conversations/${id}`,
    method: 'delete'
  })
}

export function getChatHistory(id) {
  return request({
    url: `/agent/conversations/${id}/history`,
    method: 'get'
  })
}

export function chatStream(data, onChunk, onDone, onError) {
  const token = localStorage.getItem('token');
  fetch('/api/agent/chat/stream', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': token ? `Bearer ${token}` : ''
    },
    body: JSON.stringify({
      conversation_id: data.conversation_id,
      message: data.message,
      mode: data.mode || 'auto',
      pay_type: data.pay_type || '',
      payment_password: data.payment_password || '',
      face_image_url: data.face_image_url || ''
    })
  })
  .then(async (response) => {
    if (!response.ok) {
      const text = await response.text();
      let errorMsg = '请求失败';
      try {
        const parsed = JSON.parse(text);
        errorMsg = parsed.message || parsed.msg || errorMsg;
      } catch (e) {
        errorMsg = text || errorMsg;
      }
      throw new Error(errorMsg);
    }
    const reader = response.body.getReader();
    const decoder = new TextDecoder('utf-8');
    let buffer = '';

    while (true) {
      const { value, done } = await reader.read();
      if (done) break;

      buffer += decoder.decode(value, { stream: true });
      const lines = buffer.split('\n');
      buffer = lines.pop(); // keep partial line in buffer

      for (const line of lines) {
        const cleaned = line.trim();
        if (!cleaned) continue;

        if (cleaned.startsWith('data: ')) {
          const content = cleaned.slice(6);
          if (content === '[DONE]') {
            onDone();
            return;
          }
          if (content.startsWith('[ERROR]')) {
            onError(new Error(content.slice(8)));
            return;
          }
          try {
            const parsed = JSON.parse(content);
            if (parsed && parsed.type !== undefined) {
              onChunk(parsed);
            } else if (parsed && parsed.chunk !== undefined) {
              onChunk({
                type: 'message_delta',
                data: { chunk: parsed.chunk }
              });
            }
          } catch (e) {
            console.error('Error parsing stream chunk:', e);
          }
        }
      }
    }
    onDone();
  })
  .catch((err) => {
    onError(err);
  });
}

export function approveAction(conversationId, actionId, data) {
  return request({
    url: `/agent/sessions/${conversationId}/actions/${actionId}/approve`,
    method: 'post',
    data
  })
}

export function rejectAction(conversationId, actionId) {
  return request({
    url: `/agent/sessions/${conversationId}/actions/${actionId}/reject`,
    method: 'post'
  })
}
