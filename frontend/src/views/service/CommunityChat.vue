<template>
  <div class="community-chat-page">
    <Navbar />

    <div class="container">
      <h1 class="page-title">社区群聊</h1>

      <div class="chat-card card">
        <div class="chat-list" ref="listRef" v-loading="loading">
          <el-empty v-if="!loading && messages.length === 0" description="暂无消息，来发第一条吧" />

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
            placeholder="说点什么..."
            @keydown.ctrl.enter.prevent="handleSend"
          />
          <div class="actions">
            <span class="hint">Ctrl + Enter 发送</span>
            <el-button type="primary" :loading="sending" @click="handleSend">发送</el-button>
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
  padding-bottom: 40px;
}

.chat-card {
  display: flex;
  flex-direction: column;
  min-height: 70vh;
  overflow: hidden;
}

.chat-list {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
}

.message-item {
  display: flex;
  align-items: flex-start;
  gap: 10px;
  margin-bottom: 14px;
}

.message-item.is-self {
  flex-direction: row-reverse;
}

.avatar {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  object-fit: cover;
  background: #f5f5f5;
}

.bubble-wrap {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  align-items: flex-start;
}

.message-item.is-self .bubble-wrap {
  align-items: flex-end;
}

.meta {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 12px;
  color: var(--text-secondary);
  margin-bottom: 4px;
}

.message-item.is-self .meta {
  justify-content: flex-end;
}

.name {
  font-weight: 600;
}

.bubble {
  display: inline-block;
  max-width: 75%;
  background: #f2f4f7;
  border-radius: 10px;
  border-top-left-radius: 4px;
  padding: 10px 12px;
  line-height: 1.6;
  word-break: break-word;
}

.message-item.is-self .bubble {
  background: linear-gradient(135deg, #3a8bff, #66b1ff);
  color: #fff;
  border-top-left-radius: 10px;
  border-top-right-radius: 4px;
}

.composer {
  border-top: 1px solid var(--border-color);
  padding: 12px 16px 16px;
  background: #fff;
}

.actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 10px;
}

.hint {
  font-size: 12px;
  color: var(--text-secondary);
}
</style>
