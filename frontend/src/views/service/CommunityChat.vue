<template>
  <div class="community-chat-page">
    <Navbar />

    <div class="container custom-container">
      <div class="page-header">
        <h1 class="page-title highlight-title">社区群聊</h1>
      </div>

      <div class="chat-card">
        <div class="chat-list custom-scrollbar" ref="listRef" v-loading="loading">
          <el-empty v-if="!loading && messages.length === 0" description="暂无消息，来发第一条吧" image-size="120" />

          <div
            v-for="item in messages"
            :key="item.id"
            class="message-item"
            :class="{ 'is-self': isSelf(item) }"
          >
            <img
              class="avatar"
              :src="item.user?.avatar || defaultAvatar"
              alt="avatar"
            />
            <div class="bubble-wrap">
              <div class="meta">
                <span class="name">{{ item.user?.username || '用户' }}</span>
                <span class="time">{{ formatTime(item.created_at) }}</span>
              </div>
              <div class="bubble">{{ item.content }}</div>
            </div>
          </div>
        </div>

        <div class="composer">
          <el-input
            v-model="draft"
            type="textarea"
            :rows="3"
            maxlength="1000"
            show-word-limit
            placeholder="参与社区讨论，说点什么..."
            class="custom-textarea"
            @keydown.ctrl.enter.prevent="handleSend"
          />
          <div class="actions">
            <span class="hint">快捷发送: <strong>Ctrl + Enter</strong></span>
            <button class="btn-send" :class="{ 'is-loading': sending }" :disabled="!draft.trim() || sending" @click="handleSend">
              {{ sending ? '发送中...' : '发 送' }}
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, nextTick, onMounted, onUnmounted, ref } from 'vue'
import dayjs from 'dayjs'
import { ElMessage } from 'element-plus'
import Navbar from '@/components/layout/Navbar.vue'
import { getCommunityMessages, sendCommunityMessage } from '@/api/community'
import { useUserStore } from '@/stores/user'

const defaultAvatar = 'https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png'
const userStore = useUserStore()
const messages = ref([])
const draft = ref('')
const loading = ref(false)
const sending = ref(false)
const listRef = ref(null)
let pollTimer = null
const currentUserID = computed(() => Number(userStore.userInfo?.id || 0))

function isSelf(item) {
  return currentUserID.value > 0 && Number(item?.user_id || 0) === currentUserID.value
}

function formatTime(value) {
  return dayjs(value).format('MM-DD HH:mm:ss')
}

function scrollToBottom() {
  nextTick(() => {
    if (!listRef.value) return
    listRef.value.scrollTop = listRef.value.scrollHeight
  })
}

async function fetchMessages() {
  loading.value = true
  try {
    const res = await getCommunityMessages({ page: 1, size: 100 })
    const list = Array.isArray(res?.list) ? res.list.slice().reverse() : []
    messages.value = list
    scrollToBottom()
  } catch (error) {
    ElMessage.error(error.response?.data?.msg || error.message || '加载消息失败')
  } finally {
    loading.value = false
  }
}

async function handleSend() {
  const content = draft.value.trim()
  if (!content || sending.value) return

  sending.value = true
  try {
    await sendCommunityMessage({ content })
    draft.value = ''
    await fetchMessages()
  } catch (error) {
    ElMessage.error(error.response?.data?.msg || error.message || '发送失败')
  } finally {
    sending.value = false
  }
}

onMounted(async () => {
  if (!currentUserID.value) {
    await userStore.fetchUserInfo()
  }
  await fetchMessages()
  pollTimer = setInterval(fetchMessages, 5000)
})

onUnmounted(() => {
  if (pollTimer) {
    clearInterval(pollTimer)
    pollTimer = null
  }
})
</script>

<style scoped>
.community-chat-page {
  min-height: 100vh;
  background-color: #f8f9fa;
  padding-bottom: 60px;
}

.custom-container {
  max-width: 1000px;
  margin: 0 auto;
}

.page-header { padding: 32px 0 24px; }

.highlight-title {
  display: inline-block; position: relative; font-size: 32px; color: #2c3e50; font-weight: 700; z-index: 1; margin: 0;
}
.highlight-title::after {
  content: ''; position: absolute; bottom: 4px; left: -5%; width: 110%; height: 14px;
  background-color: #2d597b; opacity: 0.15; border-radius: 6px; z-index: -1;
}

.chat-card {
  display: flex;
  flex-direction: column;
  height: 75vh;
  background: #ffffff;
  border-radius: 16px;
  box-shadow: 0 4px 24px rgba(0, 0, 0, 0.04);
  overflow: hidden;
  border: 1px solid rgba(0,0,0,0.02);
}

.chat-list {
  flex: 1;
  overflow-y: auto;
  padding: 32px;
  background: #fbfcfd;
}

.custom-scrollbar::-webkit-scrollbar { width: 6px; }
.custom-scrollbar::-webkit-scrollbar-track { background: transparent; }
.custom-scrollbar::-webkit-scrollbar-thumb { background: #dcdfe6; border-radius: 3px; }
.custom-scrollbar::-webkit-scrollbar-thumb:hover { background: #c0c4cc; }

.message-item {
  display: flex;
  align-items: flex-start;
  gap: 16px;
  margin-bottom: 24px;
}

.message-item.is-self {
  flex-direction: row-reverse;
}

.avatar {
  width: 44px;
  height: 44px;
  border-radius: 50%;
  object-fit: cover;
  box-shadow: 0 2px 8px rgba(0,0,0,0.05);
}

.bubble-wrap {
  max-width: 70%;
  display: flex;
  flex-direction: column;
}

.message-item.is-self .bubble-wrap {
  align-items: flex-end;
}

.meta {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
  color: #a4b0be;
  margin-bottom: 6px;
}

.name { font-weight: 600; color: #606266; }

.bubble {
  display: inline-block;
  padding: 12px 18px;
  font-size: 15px;
  line-height: 1.6;
  word-break: break-word;
  background: #ffffff;
  color: #2c3e50;
  border-radius: 12px;
  border-top-left-radius: 2px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.02);
  border: 1px solid #ebeef5;
}

.message-item.is-self .bubble {
  background: #2d597b;
  color: #ffffff;
  border: none;
  border-top-left-radius: 12px;
  border-top-right-radius: 2px;
}

.composer {
  padding: 24px 32px;
  background: #ffffff;
  border-top: 1px solid #f0f2f5;
}

:deep(.custom-textarea .el-textarea__inner) {
  background: #f8f9fa;
  border: 1px solid #ebeef5;
  border-radius: 8px;
  padding: 12px;
  font-size: 15px;
  box-shadow: none;
  transition: all 0.3s;
}

:deep(.custom-textarea .el-textarea__inner:focus) {
  background: #ffffff;
  border-color: #2d597b;
  box-shadow: 0 0 0 3px rgba(45, 89, 123, 0.1);
}

.actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 16px;
}

.hint { font-size: 13px; color: #a4b0be; }
.hint strong { background: #f4f4f5; padding: 2px 6px; border-radius: 4px; color: #606266; }

.btn-send {
  background: #2d597b;
  color: #ffffff;
  border: none;
  border-radius: 20px;
  padding: 8px 32px;
  font-size: 15px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s;
  box-shadow: 0 4px 12px rgba(45, 89, 123, 0.2);
}

.btn-send:hover:not(:disabled) {
  background: #1f435d;
  transform: translateY(-2px);
}

.btn-send:disabled {
  background: #c0c4cc;
  box-shadow: none;
  cursor: not-allowed;
}
</style>