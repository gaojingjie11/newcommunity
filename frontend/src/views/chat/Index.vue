<template>
  <div class="chat-page">
    <Navbar />
    
    <div class="container chat-container">
      <div class="chat-card">
        <!-- 头部 -->
        <div class="chat-header">
          <div class="header-left">
            <span class="header-icon">🤖</span>
            <div class="header-title-box">
              <h2>智慧社区助手</h2>
              <p>支持通知总结、报修创建、商品下单与支付</p>
            </div>
          </div>
          <div class="header-right">
            <span class="status-dot"></span>
            <span class="status-text">在线</span>
          </div>
        </div>

        <!-- 聊天记录区 -->
        <div ref="historyRef" class="chat-history custom-scrollbar">
          <div
            v-for="(msg, index) in messages"
            :key="index"
            class="message-item"
            :class="msg.role"
          >
            <div class="avatar">
              <span v-if="msg.role === 'user'">👤</span>
              <span v-else-if="msg.role === 'assistant'">🤖</span>
              <span v-else>⚠️</span>
            </div>
            <div class="content">
              <div class="bubble-wrapper">
                <div class="bubble">{{ msg.content }}</div>
              </div>
              <div class="time">{{ msg.time }}</div>
            </div>
          </div>

          <!-- 加载动画 -->
          <div v-if="loading" class="message-item assistant">
            <div class="avatar">🤖</div>
            <div class="content">
              <div class="bubble-wrapper">
                <div class="bubble loading-bubble">
                  <span>.</span><span>.</span><span>.</span>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- 输入区 -->
        <div class="chat-input-area">
          <div class="input-wrapper">
            <el-input
              v-model="inputContent"
              type="textarea"
              :rows="3"
              placeholder="请输入您的问题，例如：'帮我总结一下最新的社区公告'..."
              class="custom-textarea"
              @keyup.enter.ctrl="handleSend"
              resize="none"
            />
          </div>

          <div class="actions">
            <span class="tip">快捷发送: <strong>Ctrl + Enter</strong></span>
            <button 
              class="btn-send" 
              :class="{ 'is-loading': loading, 'is-disabled': !inputContent.trim() && !loading }" 
              :disabled="(!inputContent.trim() && !loading) || loading" 
              @click="handleSend"
            >
              {{ loading ? '发送中...' : '发 送' }}
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { nextTick, onMounted, ref } from 'vue'
import dayjs from 'dayjs'
import { ElMessage, ElMessageBox } from 'element-plus'
import Navbar from '@/components/layout/Navbar.vue'
import { getChatHistory, sendChat } from '@/api/chat'

const buildGreetingMessage = () => ({
  role: 'assistant',
  content: '您好，我是智慧社区助手。您可以让我帮您总结通知、创建报修、搜索商品、下单和支付。',
  time: dayjs().format('HH:mm')
})

const formatMessageTime = (value) => {
  const parsed = dayjs(value)
  return parsed.isValid() ? parsed.format('HH:mm') : dayjs().format('HH:mm')
}

const messages = ref([buildGreetingMessage()])
const inputContent = ref('')
const loading = ref(false)
const historyRef = ref(null)

const isPayIntent = (text) => ['支付', '付款', '结算', '确认支付'].some((kw) => text.includes(kw))

const scrollToBottom = () => {
  nextTick(() => {
    if (historyRef.value) {
      historyRef.value.scrollTop = historyRef.value.scrollHeight
    }
  })
}

const loadHistory = async () => {
  try {
    const res = await getChatHistory({ limit: 100 })
    const list = Array.isArray(res?.list) ? res.list : []
    if (list.length === 0) {
      messages.value = [buildGreetingMessage()]
      scrollToBottom()
      return
    }

    messages.value = list.map((item) => ({
      role: item.role,
      content: item.content,
      time: formatMessageTime(item.created_at)
    }))
    scrollToBottom()
  } catch {
    messages.value = [buildGreetingMessage()]
    ElMessage.warning('聊天记录加载失败，已展示默认欢迎语')
    scrollToBottom()
  }
}

const handleSend = async () => {
  const content = inputContent.value.trim()
  if (!content || loading.value) {
    return
  }

  let paymentPassword = ''
  if (isPayIntent(content)) {
    try {
      const { value } = await ElMessageBox.prompt(
        '请输入登录密码完成支付',
        '安全支付验证',
        {
          inputType: 'password',
          inputPlaceholder: '登录密码',
          confirmButtonText: '确认支付',
          cancelButtonText: '取消'
        }
      )
      paymentPassword = (value || '').trim()
      if (!paymentPassword) {
        ElMessage.warning('未输入支付密码，已取消本次支付请求')
        return
      }
    } catch {
      return
    }
  }

  messages.value.push({
    role: 'user',
    content,
    time: dayjs().format('HH:mm')
  })
  inputContent.value = ''
  loading.value = true
  scrollToBottom()

  try {
    const res = await sendChat({
      content,
      payment_password: paymentPassword
    })

    const reply = (res?.reply || '').trim()
    if (!reply) {
      throw new Error('empty AI response')
    }
    messages.value.push({
      role: 'assistant',
      content: reply,
      time: dayjs().format('HH:mm')
    })
  } catch (error) {
    ElMessage.error('AI 响应失败，请稍后重试')
    messages.value.push({
      role: 'system',
      content: `生成失败: ${error?.message || '网络错误'}`,
      time: dayjs().format('HH:mm')
    })
  } finally {
    loading.value = false
    scrollToBottom()
  }
}

onMounted(() => {
  loadHistory()
})
</script>

<style scoped>
.chat-page {
  min-height: 100vh;
  padding-bottom: 40px;
  background-color: #f5f7fa; /* 灰底色衬托主卡片 */
}

/* ★ 加宽的聊天容器 ★ */
.chat-container {
  max-width: 1100px; /* 原来是800px，现在加宽 */
  margin: 30px auto;
}

.chat-card {
  height: 82vh; /* 稍微增加一点高度 */
  display: flex;
  flex-direction: column;
  background: #ffffff;
  border-radius: 16px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.06); /* 高级软阴影 */
  overflow: hidden;
  border: 1px solid rgba(0, 0, 0, 0.02);
}

/* 头部样式 */
.chat-header {
  padding: 20px 32px;
  border-bottom: 1px solid #f0f2f5;
  background: #ffffff;
  display: flex;
  justify-content: space-between;
  align-items: center;
  z-index: 10;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 16px;
}

.header-icon {
  font-size: 36px;
  background: #f0f7ff;
  border-radius: 50%;
  width: 56px;
  height: 56px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.header-title-box h2 {
  margin: 0 0 4px 0;
  font-size: 20px;
  font-weight: 600;
  color: #2c3e50;
}

.header-title-box p {
  margin: 0;
  color: #8c939d;
  font-size: 13px;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 8px;
  background: #f6ffed;
  padding: 6px 16px;
  border-radius: 20px;
  border: 1px solid #b7eb8f;
}

.status-dot {
  width: 8px;
  height: 8px;
  background-color: #52c41a;
  border-radius: 50%;
}

.status-text {
  font-size: 13px;
  color: #52c41a;
  font-weight: 500;
}

/* 聊天历史区 */
.chat-history {
  flex: 1;
  padding: 32px;
  overflow-y: auto;
  background: #f8f9fa; /* 浅灰色背景，让气泡更突出 */
  scroll-behavior: smooth;
}

/* 优化滚动条 */
.custom-scrollbar::-webkit-scrollbar {
  width: 6px;
}
.custom-scrollbar::-webkit-scrollbar-track {
  background: transparent;
}
.custom-scrollbar::-webkit-scrollbar-thumb {
  background: #dcdfe6;
  border-radius: 3px;
}
.custom-scrollbar::-webkit-scrollbar-thumb:hover {
  background: #c0c4cc;
}

/* 气泡行布局 */
.message-item {
  display: flex;
  margin-bottom: 28px;
  align-items: flex-start;
}

.message-item.user {
  flex-direction: row-reverse;
}

.avatar {
  width: 44px;
  height: 44px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20px;
  flex-shrink: 0;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

.message-item.user .avatar {
  margin-left: 16px;
  background: #e6f1fc;
}

.message-item.assistant .avatar {
  margin-right: 16px;
  background: #ffffff;
}

.message-item.system .avatar {
  margin-right: 16px;
  background: #fff0f0;
}

.content {
  max-width: 75%; /* 加宽了容器，气泡占比也可以相应增加 */
  display: flex;
  flex-direction: column;
}

.message-item.user .content {
  align-items: flex-end;
}

.message-item.assistant .content,
.message-item.system .content {
  align-items: flex-start;
}

/* 对话气泡样式 */
.bubble {
  padding: 14px 20px;
  border-radius: 16px;
  font-size: 15px;
  line-height: 1.6;
  white-space: pre-wrap;
  word-break: break-word;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.02);
}

.message-item.user .bubble {
  background: #2d597b; /* 统一的商务蓝 */
  color: #ffffff;
  border-top-right-radius: 4px; /* 尖角朝向右侧 */
}

.message-item.assistant .bubble {
  background: #ffffff;
  color: #2c3e50;
  border: 1px solid #ebeef5;
  border-top-left-radius: 4px; /* 尖角朝向左侧 */
}

.message-item.system .bubble {
  background: #fff0f0;
  color: #f56c6c;
  border: 1px solid #fbc4c4;
  border-top-left-radius: 4px;
}

.time {
  font-size: 12px;
  color: #b2bec3;
  margin-top: 8px;
}

/* 输入区域 */
.chat-input-area {
  padding: 24px 32px;
  background: #ffffff;
  border-top: 1px solid #f0f2f5;
}

.input-wrapper {
  background: #f8f9fa;
  border: 1px solid #ebeef5;
  border-radius: 12px;
  padding: 4px;
  transition: all 0.3s;
}

.input-wrapper:focus-within {
  border-color: #2d597b;
  background: #ffffff;
  box-shadow: 0 0 0 3px rgba(45, 89, 123, 0.1);
}

:deep(.custom-textarea .el-textarea__inner) {
  background: transparent;
  box-shadow: none !important;
  border: none !important;
  color: #2c3e50;
  font-size: 15px;
  padding: 12px;
  line-height: 1.5;
}

:deep(.custom-textarea .el-textarea__inner::placeholder) {
  color: #a4b0be;
}

.actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 16px;
}

.tip {
  font-size: 13px;
  color: #909399;
}

.tip strong {
  color: #606266;
  background: #f4f4f5;
  padding: 4px 8px;
  border-radius: 4px;
  border: 1px solid #e9e9eb;
}

/* 发送按钮定制 */
.btn-send {
  background: #2d597b;
  color: #ffffff;
  border: none;
  border-radius: 20px; /* 胶囊按钮 */
  padding: 10px 32px;
  font-size: 15px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s;
  box-shadow: 0 4px 12px rgba(45, 89, 123, 0.2);
}

.btn-send:hover:not(.is-disabled) {
  background: #1f435d;
  transform: translateY(-2px);
  box-shadow: 0 6px 16px rgba(45, 89, 123, 0.3);
}

.btn-send.is-disabled {
  background: #c0c4cc;
  box-shadow: none;
  cursor: not-allowed;
}

/* 加载动画点点 */
.loading-bubble {
  padding: 12px 24px;
}

.loading-bubble span {
  display: inline-block;
  animation: bounce 1.4s infinite ease-in-out both;
  margin: 0 3px;
  font-size: 24px;
  line-height: 10px;
  font-weight: bold;
  color: #a4b0be;
}

.loading-bubble span:nth-child(1) {
  animation-delay: -0.32s;
}

.loading-bubble span:nth-child(2) {
  animation-delay: -0.16s;
}

@keyframes bounce {
  0%, 80%, 100% {
    transform: scale(0);
    opacity: 0.5;
  }
  40% {
    transform: scale(1);
    opacity: 1;
  }
}
</style>