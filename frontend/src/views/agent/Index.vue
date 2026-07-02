<template>
  <div class="agent-page-layout">
    <Navbar />

    <div class="agent-main-container">
      <!-- Left Sidebar: Conversations list -->
      <aside :class="['sidebar', { collapsed: isSidebarCollapsed }]">
        <div class="sidebar-top-bar">
          <span class="sidebar-brand">
            <el-icon class="brand-icon"><ChatLineRound /></el-icon>
            智能社区管家
          </span>
          <el-tooltip content="收起侧边栏" placement="bottom">
            <el-button class="sidebar-toggle-btn" link @click="isSidebarCollapsed = true">
              <el-icon size="18"><Fold /></el-icon>
            </el-button>
          </el-tooltip>
        </div>

        <div class="sidebar-header">
          <el-button type="primary" class="new-chat-btn" @click="handleCreateSession">
            <el-icon><Plus /></el-icon>
            新建对话
          </el-button>
        </div>

        <div class="sessions-list-container">
          <div v-if="sessions.length === 0" class="empty-sessions">
            暂无对话历史
          </div>
          <div
            v-for="item in sessions"
            :key="item.id"
            :class="['session-item', { active: activeSessionId === item.id }]"
            @click="handleSelectSession(item.id)"
          >
            <div class="session-info">
              <div class="session-title-row">
                <span class="session-title">{{ item.title || '新对话' }}</span>
                <el-popconfirm
                  title="确定删除此对话吗？"
                  confirm-button-text="确定"
                  cancel-button-text="取消"
                  @confirm.stop="handleDeleteSession(item.id)"
                >
                  <template #reference>
                    <el-button
                      type="danger"
                      link
                      class="delete-session-btn"
                      @click.stop
                    >
                      <el-icon><Delete /></el-icon>
                    </el-button>
                  </template>
                </el-popconfirm>
              </div>
              <div class="session-summary">
                {{ item.summary || '开启与智慧社区助理的对话吧...' }}
              </div>
              <div class="session-time">
                {{ formatTime(item.updated_at) }}
              </div>
            </div>
          </div>
        </div>
      </aside>

      <!-- Right Chat Window -->
      <main class="chat-window">
        <header class="chat-header">
          <div class="header-left-group">
            <el-tooltip content="展开侧边栏" placement="bottom" v-if="isSidebarCollapsed">
              <el-button class="sidebar-toggle-btn sidebar-open-btn" link @click="isSidebarCollapsed = false">
                <el-icon size="18"><Expand /></el-icon>
              </el-button>
            </el-tooltip>
            <div class="chat-title-info">
              <h2>{{ activeSession?.title || '新对话' }}</h2>
            </div>
          </div>
          <div class="header-meta">
            <span class="model-pill">{{ currentModelLabel }}</span>
            <span class="status-text">{{ isStreaming ? '响应中' : '已连接' }}</span>
          </div>
        </header>

        <!-- Messages scroll container -->
        <div ref="msgContainer" class="messages-container">
          <div v-if="messages.length === 0" class="welcome-screen">
            <div class="welcome-avatar">
              <span class="welcome-mark">AI</span>
            </div>
            <h3>有什么我可以帮您的？</h3>
            <p class="welcome-subtitle">更快地处理社区公告、下单支付、物业报修与日常问答。</p>
            <div class="quick-starts">
              <button
                v-for="qs in quickStarts"
                :key="qs.title"
                class="qs-card"
                @click="handleQuickStart(qs.prompt)"
              >
                <div class="qs-card-top">
                  <div class="qs-title">{{ qs.title }}</div>
                  <span class="qs-arrow">+</span>
                </div>
                <div class="qs-desc">{{ qs.desc }}</div>
              </button>
            </div>
          </div>

          <div
            v-for="(msg, idx) in messages"
            :key="msg.id || idx"
            :class="['message-row', msg.role]"
          >
            <div class="avatar-col">
              <div v-if="msg.role === 'assistant'" class="avatar bot">
                <span class="assistant-mark">AI</span>
              </div>
              <div v-else class="avatar user">
                我
              </div>
            </div>
            <div class="bubble-col">
              <div class="bubble-info" v-if="msg.role === 'assistant'">
                智能管家
              </div>
              <div class="bubble-info" v-else>
                业主
              </div>
              
              <!-- If assistant message is empty and streaming, show typing indicator dots -->
              <div v-if="msg.role === 'assistant' && msg.content === '' && isStreaming && !msg.proposed_action" class="bubble typing-bubble">
                <span class="dot"></span>
                <span class="dot"></span>
                <span class="dot"></span>
              </div>

              <!-- If action approval is proposed -->
              <div v-else-if="msg.proposed_action" class="bubble approval-card-bubble">
                <div class="approval-card">
                  <div class="approval-header">
                    <el-icon color="#e6a23c" size="18"><Warning /></el-icon>
                    <h4>需要您的操作确认</h4>
                  </div>
                  
                  <div class="approval-body">
                    <div v-if="msg.proposed_action.action_type === 'create_order'" class="action-details">
                      <p class="action-desc">AI 计划为您在商城中订购商品：</p>
                      <div class="detail-row">
                        <span class="label">商品ID:</span>
                        <span class="val">{{ msg.proposed_action.payload.product_id }}</span>
                      </div>
                      <div class="detail-row">
                        <span class="label">购买数量:</span>
                        <span class="val">{{ msg.proposed_action.payload.quantity }} 件</span>
                      </div>
                    </div>

                    <div v-else-if="msg.proposed_action.action_type === 'pay_order'" class="action-details">
                      <p class="action-desc">AI 计划对以下商城订单发起余额扣款支付：</p>
                      <div class="detail-row">
                        <span class="label">订单ID:</span>
                        <span class="val">{{ msg.proposed_action.payload.order_id }}</span>
                      </div>
                    </div>

                    <div v-else-if="msg.proposed_action.action_type === 'submit_repair'" class="action-details">
                      <p class="action-desc">AI 计划为您提交物业服务单：</p>
                      <div class="detail-row">
                        <span class="label">工单类别:</span>
                        <span class="val">{{ msg.proposed_action.payload.type === 'repair' ? '报修' : '投诉' }}</span>
                      </div>
                      <div class="detail-row">
                        <span class="label">故障分类:</span>
                        <span class="val">{{ msg.proposed_action.payload.category }}</span>
                      </div>
                      <div class="detail-row">
                        <span class="label">具体描述:</span>
                        <span class="val">{{ msg.proposed_action.payload.description }}</span>
                      </div>
                    </div>
                  </div>

                  <div class="approval-actions" v-if="!msg.action_resolved">
                    <el-button size="small" type="danger" plain @click="handleRejectAction(msg)">拒绝</el-button>
                    <el-button size="small" type="primary" :loading="msg.action_submitting" @click="handleApproveAction(msg)">确认同意</el-button>
                  </div>
                  
                  <div class="approval-resolved-status" v-else>
                    <div v-if="msg.action_resolved === 'approved' && msg.proposed_action.action_type === 'create_order' && getCreatedOrderId(msg)" class="pay-after-approval-row">
                      <el-tag type="success" size="small">已同意授权并执行</el-tag>
                      <el-tag v-if="isOrderPaidOrResolved(getCreatedOrderId(msg))" type="success" size="small" style="margin-left: 8px;">已完成支付</el-tag>
                      <el-button v-else type="warning" size="small" class="pay-btn-margin" @click="handlePayForCreatedOrder(getCreatedOrderId(msg), msg)">去支付</el-button>
                    </div>
                    <el-tag v-else :type="msg.action_resolved === 'approved' ? 'success' : 'danger'" size="small">
                      {{ msg.action_resolved === 'approved' ? '已同意授权并执行' : '已拒绝授权' }}
                    </el-tag>
                  </div>
                </div>
              </div>

              <!-- Otherwise render normal markdown -->
              <div v-else-if="msg.content" class="bubble" v-html="msg.rendered_content || ''"></div>

              <!-- Tool status indicator -->
              <div v-if="msg.tool_status" class="tool-calling-status">
                <el-icon class="is-loading"><Loading /></el-icon>
                <span>{{ msg.tool_status }}</span>
              </div>
            </div>
          </div>
        </div>

        <!-- Chat Input Footer -->
        <footer class="chat-footer">
          <div class="chat-mode-row">
            <span class="chat-mode-label">模式</span>
            <div class="chat-mode-switch">
              <button
                v-for="mode in chatModes"
                :key="mode.value"
                type="button"
                :class="['chat-mode-btn', { active: chatMode === mode.value }]"
                @click="setChatMode(mode.value)"
              >
                {{ mode.label }}
              </button>
            </div>
            <div class="chat-mode-hint">{{ activeModeHint }}</div>
          </div>

          <div class="input-form">
            <el-input
              v-model="inputMsg"
              type="textarea"
              :autosize="{ minRows: 1, maxRows: 6 }"
              resize="none"
              placeholder="有问题，尽管问..."
              class="message-input"
              @keydown.enter.exact.prevent="handleSend"
            />
            
            <div class="input-right-actions">
              <span class="input-model-pill">{{ currentModelLabel }}</span>
              <el-button
                type="primary"
                class="send-btn"
                :disabled="!inputMsg.trim() || isStreaming"
                @click="handleSend"
              >
                <el-icon><Top /></el-icon>
              </el-button>
            </div>
          </div>
        </footer>
      </main>
    </div>

    <!-- 支付验证弹窗 -->
    <PayAuthDialog
      v-model="showPayAuth"
      title="订单支付验证"
      :face-registered="Boolean(userStore.userInfo?.face_registered)"
      :loading="paySubmitting"
      @confirm="submitOrderPay"
    />
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onBeforeUnmount, nextTick } from 'vue'
import { ElMessage } from 'element-plus'
import {
  Plus,
  Delete,
  Loading,
  Warning,
  Top,
  Fold,
  Expand,
  ChatLineRound
} from '@element-plus/icons-vue'
import Navbar from '@/components/layout/Navbar.vue'
import PayAuthDialog from '@/components/payment/PayAuthDialog.vue'
import { useUserStore } from '@/stores/user'
import {
  getConversations,
  createConversation,
  deleteConversation,
  getChatHistory,
  chatStream,
  approveAction,
  rejectAction
} from '@/api/chat'
import { payOrder, getOrderDetail } from '@/api/order'
import dayjs from 'dayjs'

// Quick start suggestion prompts
const quickStarts = [
  { title: '查询公告', desc: '获取社区最新通知与资讯公告', prompt: '查询最新社区公告' },
  { title: '推荐好物', desc: '帮我看看商城有哪些方便面或可乐？', prompt: '商城有哪些方便面或可乐？' },
  { title: '报修登记', desc: '物业报修：家里卫生间水龙头漏水了', prompt: '物业报修：家里卫生间水龙头漏水了' },
  { title: '订单支付', desc: '对特定商城订单号发起余额付款', prompt: '对订单号 20261001 进行付款支付' }
]

const chatModes = [
  { value: 'auto', label: '均衡', hint: '默认优先快模型，复杂问题自动切深度分析。' },
  { value: 'fast', label: '极速', hint: '更适合日常问答、查询和轻量操作。' },
  { value: 'deep', label: '高级', hint: '更适合总结、分析、报表和复杂推理。' }
]

// State variables
const isSidebarCollapsed = ref(false)
const sessions = ref([])
const activeSessionId = ref('')
const inputMsg = ref('')
const messages = ref([])
const isStreaming = ref(false)

const userStore = useUserStore()
const showPayAuth = ref(false)
const paySubmitting = ref(false)
const pendingPayMsg = ref(null)
const activeStreamingMessage = ref('')
const msgContainer = ref(null)
const payingOrderId = ref(null)
const payingOrderMsg = ref(null)
const chatMode = ref(localStorage.getItem('agent-chat-mode') || 'auto')

let pendingAssistantChunk = ''
let assistantRenderFrame = 0
let scrollFrame = 0

// Active session helper
const activeSession = computed(() => {
  return sessions.value.find(s => s.id === activeSessionId.value)
})

const activeModeHint = computed(() => {
  return chatModes.find(item => item.value === chatMode.value)?.hint || chatModes[0].hint
})

const currentModelLabel = computed(() => {
  const item = chatModes.find(m => m.value === chatMode.value)
  return item ? item.label : '均衡'
})

// Format time
const formatTime = (timeStr) => {
  if (!timeStr) return ''
  return dayjs(timeStr).format('YYYY-MM-DD HH:mm')
}

const setChatMode = (mode) => {
  chatMode.value = mode
  localStorage.setItem('agent-chat-mode', mode)
}

// Fetch all conversations
const fetchSessions = async () => {
  try {
    const list = await getConversations()
    sessions.value = Array.isArray(list) ? list : (list?.list || [])
    if (sessions.value.length > 0 && !activeSessionId.value) {
      handleSelectSession(sessions.value[0].id)
    } else if (sessions.value.length === 0) {
      // Create session dynamically if none exist
      await handleCreateSession()
    }
  } catch (err) {
    console.error('Failed to get conversations:', err)
    ElMessage.error('获取对话列表失败: ' + err.message)
  }
}

// Create new session
const handleCreateSession = async () => {
  try {
    const res = await createConversation({ title: '新对话' })
    const newSession = {
      id: res.id,
      title: '新对话',
      summary: '',
      updated_at: new Date().toISOString()
    }
    sessions.value.unshift(newSession)
    handleSelectSession(res.id)
  } catch (err) {
    ElMessage.error('创建新对话失败: ' + err.message)
  }
}

// Delete session
const handleDeleteSession = async (id) => {
  try {
    await deleteConversation(id)
    sessions.value = sessions.value.filter(s => s.id !== id)
    if (activeSessionId.value === id) {
      activeSessionId.value = ''
      messages.value = []
      if (sessions.value.length > 0) {
        handleSelectSession(sessions.value[0].id)
      } else {
        await handleCreateSession()
      }
    }
    ElMessage.success('对话已成功删除')
  } catch (err) {
    ElMessage.error('删除对话失败: ' + err.message)
  }
}

// Select session and load history
const handleSelectSession = async (id, clearMessages = true) => {
  activeSessionId.value = id
  if (clearMessages) {
    messages.value = []
  }
  try {
    const res = await getChatHistory(id)
    messages.value = (res?.list || []).map(normalizeMessage)
    
    // Fetch statuses of all approved order creations
    messages.value.forEach(msg => {
      if (msg.action_resolved === 'approved' && msg.proposed_action?.action_type === 'create_order') {
        const orderId = getCreatedOrderId(msg)
        if (orderId) {
          fetchOrderStatus(orderId)
        }
      }
    })
    
    nextTick(scheduleScrollToBottom)
  } catch (err) {
    ElMessage.error('加载聊天记录失败: ' + err.message)
  }
}

// Trigger quick start chip
const handleQuickStart = (text) => {
  inputMsg.value = text
}

const normalizeMessage = (msg) => {
  const normalized = {
    ...msg,
    content: msg.content || '',
    rendered_content: msg.content ? renderMarkdown(msg.content) : ''
  }

  if (normalized.event_type === 'approval_required' && normalized.event_payload) {
    try {
      normalized.proposed_action = JSON.parse(normalized.event_payload)
    } catch (e) {
      console.error('Failed to parse event_payload:', e)
    }
  }

  return normalized
}

// Custom Markdown parser function
const renderMarkdown = (content) => {
  if (!content) return ''
  // 1. Escape HTML
  let html = content
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')

  // 2. Code blocks: ```go ... ```
  html = html.replace(/```([a-zA-Z0-9]+)?\n([\s\S]*?)```/g, (match, lang, code) => {
    const language = lang || 'code';
    const escapedCode = encodeURIComponent(code.trim());
    return `<div class="gpt-code-wrapper">
      <div class="gpt-code-header">
        <span class="gpt-code-lang">${language}</span>
        <button class="copy-code-btn" onclick="navigator.clipboard.writeText(decodeURIComponent('${escapedCode}')).then(() => { event.target.innerText = '已复制'; setTimeout(() => event.target.innerText = '复制代码', 2000); })">复制代码</button>
      </div>
      <pre class="code-block"><code>${code.trim()}</code></pre>
    </div>`
  })

  // 3. Bold text: **text**
  html = html.replace(/\*\*([\s\S]*?)\*\*/g, '<strong>$1</strong>')

  // 4. List Items: - item
  html = html.replace(/^\s*[-*]\s+(.+)$/gm, '<li>$1</li>')
  html = html.replace(/((?:<li>.*<\/li>\s*)+)/g, '<ul>$1</ul>')

  // 5. Line breaks
  html = html.replace(/\n/g, '<br>')

  return html
}

// Scroll chat panel to bottom
const scrollToBottom = () => {
  if (msgContainer.value) {
    msgContainer.value.scrollTop = msgContainer.value.scrollHeight
  }
}

const scheduleScrollToBottom = () => {
  if (scrollFrame) return
  scrollFrame = requestAnimationFrame(() => {
    scrollFrame = 0
    scrollToBottom()
  })
}

const renderMessageContent = (msg) => {
  msg.rendered_content = msg.content ? renderMarkdown(msg.content) : ''
}

const flushAssistantChunks = (botMsgIdx) => {
  if (!pendingAssistantChunk) return
  const botMsg = messages.value[botMsgIdx]
  if (!botMsg) {
    pendingAssistantChunk = ''
    return
  }

  activeStreamingMessage.value += pendingAssistantChunk
  pendingAssistantChunk = ''
  botMsg.content = activeStreamingMessage.value
  renderMessageContent(botMsg)
  scheduleScrollToBottom()
}

const scheduleAssistantChunkFlush = (botMsgIdx) => {
  if (assistantRenderFrame) return
  assistantRenderFrame = requestAnimationFrame(() => {
    assistantRenderFrame = 0
    flushAssistantChunks(botMsgIdx)
  })
}

// Send Message stream handler
const handleSend = async () => {
  const query = inputMsg.value.trim()
  if (!query || isStreaming.value) return

  // Validate active session
  if (!activeSessionId.value) {
    ElMessage.warning('未选中有效会话')
    return
  }

  // 1. Add user message
  const userMsgId = 'u-' + Date.now()
  messages.value.push({
    id: userMsgId,
    role: 'user',
    content: query,
    created_at: new Date().toISOString()
  })

  inputMsg.value = ''
  isStreaming.value = true
  activeStreamingMessage.value = ''
  nextTick(scrollToBottom)

  // 2. Prepare payload
  const payload = {
    conversation_id: activeSessionId.value,
    message: query,
    mode: chatMode.value,
    pay_type: '',
    payment_password: '',
    face_image_url: ''
  }

  // 3. Append assistant message placeholder
  const botMsgIdx = messages.value.length
  messages.value.push({
    role: 'assistant',
    content: '',
    rendered_content: '',
    tool_status: '智能管家正在思考中...',
    created_at: new Date().toISOString()
  })

  // 4. Call chatStream
  chatStream(
    payload,
    (event) => {
      // Event Callback — guard against stale index after session reload
      const botMsg = messages.value[botMsgIdx]
      if (!botMsg) return

      if (event.type === 'message_delta') {
        if (botMsg.tool_status === '智能管家正在思考中...') {
          botMsg.tool_status = ''
        }
        pendingAssistantChunk += event.data.chunk
        scheduleAssistantChunkFlush(botMsgIdx)
      } else if (event.type === 'tool_call_start') {
        let toolText = '智能管家正在处理业务...'
        if (event.data.tool === 'list_products') {
          toolText = '正在查询商城商品列表...'
        } else if (event.data.tool === 'query_notices') {
          toolText = '正在检索社区公告通知...'
        } else if (event.data.tool === 'create_order') {
          toolText = '正在生成商品订单...'
        } else if (event.data.tool === 'pay_order') {
          toolText = '正在发起订单余额扣款支付...'
        } else if (event.data.tool === 'submit_repair') {
          toolText = '正在提交物业报修单...'
        }
        botMsg.tool_status = toolText
      } else if (event.type === 'tool_call_end') {
        botMsg.tool_status = ''
      } else if (event.type === 'approval_required') {
        flushAssistantChunks(botMsgIdx)
        isStreaming.value = false
        botMsg.proposed_action = event.data
      }
      scheduleScrollToBottom()
    },
    () => {
      flushAssistantChunks(botMsgIdx)
      isStreaming.value = false
      fetchSessions()
    },
    (err) => {
      flushAssistantChunks(botMsgIdx)
      isStreaming.value = false
      const botMsg = messages.value[botMsgIdx]
      if (botMsg) {
        botMsg.content = '⚠️ 发送错误: ' + err.message
        renderMessageContent(botMsg)
      }
      ElMessage.error('智能管家响应异常: ' + err.message)
      scheduleScrollToBottom()
    }
  )
}

const orderStatuses = ref({})

const isOrderPaidOrResolved = (orderId) => {
  if (orderId === null || orderId === undefined) return false
  const status = orderStatuses.value[orderId]
  return status !== undefined && status !== 0
}

const fetchOrderStatus = async (orderId) => {
  if (!orderId || orderStatuses.value[orderId] !== undefined) return
  try {
    const res = await getOrderDetail(orderId)
    if (res && res.status !== undefined) {
      orderStatuses.value[orderId] = res.status
    }
  } catch (err) {
    console.error('Failed to fetch status for order', orderId, err)
  }
}

const getCreatedOrderId = (msg) => {
  if (!msg.result_payload) return null
  try {
    const parsed = JSON.parse(msg.result_payload)
    return parsed.order_id
  } catch (e) {
    return null
  }
}

const handlePayForCreatedOrder = (orderId, msg) => {
  payingOrderId.value = orderId
  payingOrderMsg.value = msg
  showPayAuth.value = true
}

const handleApproveAction = async (msg) => {
  if (msg.action_submitting) return
  
  const actionId = msg.proposed_action.action_id
  const actionType = msg.proposed_action.action_type
  
  if (actionType === 'pay_order') {
    pendingPayMsg.value = msg
    showPayAuth.value = true
    return
  }

  msg.action_submitting = true
  try {
    const res = await approveAction(activeSessionId.value, actionId, {})
    // Update the msg in-place so the card transitions without disappearing
    msg.action_resolved = 'approved'
    // For create_order, store result_payload for the "去支付" button
    if (actionType === 'create_order' && res?.order_id) {
      msg.result_payload = JSON.stringify({ order_id: res.order_id })
      orderStatuses.value[res.order_id] = 0
    }
    ElMessage.success('操作已成功授权执行')
    // Refresh session list titles, but don't reload messages (would cause card flicker)
    fetchSessions()
  } catch (err) {
    ElMessage.error('接口请求异常: ' + err.message)
  } finally {
    msg.action_submitting = false
  }
}

const submitOrderPay = async (authPayload) => {
  if (payingOrderId.value) {
    paySubmitting.value = true
    const currentOrderId = payingOrderId.value
    try {
      const res = await payOrder(currentOrderId, {
        pay_type: authPayload.pay_type,
        password: authPayload.password || '',
        face_image_url: authPayload.face_image_url || '',
        return_url: window.location.origin + '/payment/result'
      })
      
      showPayAuth.value = false
      
      if (authPayload.pay_type === 'alipay' && res && res.pay_url) {
        ElMessage.success("正在为您跳转至支付宝收银台...")
        setTimeout(() => {
          window.location.href = res.pay_url
        }, 800)
        return
      }

      orderStatuses.value[currentOrderId] = 1 // Locally mark as paid
      ElMessage.success('订单支付成功')
    } catch (err) {
      // 支付请求异常时，二次确认订单实际状态（解决网络超时但后端已成功扣款的场景）
      try {
        const detail = await getOrderDetail(currentOrderId)
        if (detail && detail.status !== undefined && detail.status !== 0) {
          // 订单已支付成功，前端对齐状态
          showPayAuth.value = false
          orderStatuses.value[currentOrderId] = detail.status
          ElMessage.success('订单支付成功')
          return
        }
      } catch (_) {
        // 二次确认也失败，展示原始错误
      }
      ElMessage.error('支付失败: ' + err.message)
    } finally {
      paySubmitting.value = false
      payingOrderId.value = null
      payingOrderMsg.value = null
    }
    return
  }

  if (!pendingPayMsg.value) return
  const msg = pendingPayMsg.value
  const actionId = msg.proposed_action.action_id

  paySubmitting.value = true
  msg.action_submitting = true
  try {
    const res = await approveAction(activeSessionId.value, actionId, {
      pay_type: authPayload.pay_type,
      payment_password: authPayload.password || '',
      face_image_url: authPayload.face_image_url || '',
      return_url: window.location.origin + '/payment/result'
    })
    
    // res is the data object containing { pay_url: "..." }
    msg.action_resolved = 'approved'
    showPayAuth.value = false
    pendingPayMsg.value = null
    ElMessage.success('操作已成功授权执行')
    
    if (authPayload.pay_type === 'alipay' && res && res.pay_url) {
      ElMessage.success("正在为您跳转至支付宝收银台...")
      setTimeout(() => {
        window.location.href = res.pay_url
      }, 800)
      return
    }

    await handleSelectSession(activeSessionId.value, false)
  } catch (err) {
    // 对 pay_order 类型的操作，二次确认订单实际状态
    if (msg.proposed_action?.action_type === 'pay_order' && msg.proposed_action?.payload?.order_id) {
      try {
        const detail = await getOrderDetail(msg.proposed_action.payload.order_id)
        if (detail && detail.status !== undefined && detail.status !== 0) {
          msg.action_resolved = 'approved'
          showPayAuth.value = false
          pendingPayMsg.value = null
          ElMessage.success('订单支付成功')
          await handleSelectSession(activeSessionId.value, false)
          return
        }
      } catch (_) {}
    }
    ElMessage.error('接口请求异常: ' + err.message)
  } finally {
    paySubmitting.value = false
    msg.action_submitting = false
  }
}

const handleRejectAction = async (msg) => {
  const actionId = msg.proposed_action.action_id
  try {
    await rejectAction(activeSessionId.value, actionId)
    msg.action_resolved = 'rejected'
    ElMessage.info('已成功取消授权')
    await handleSelectSession(activeSessionId.value, false)
  } catch (err) {
    ElMessage.error('操作异常: ' + err.message)
  }
}

// Initial mounting lifecycle
onMounted(async () => {
  await fetchSessions()
  try {
    await userStore.fetchUserInfo()
  } catch (err) {
    console.error('Failed to fetch user info:', err)
  }
})

onBeforeUnmount(() => {
  if (assistantRenderFrame) {
    cancelAnimationFrame(assistantRenderFrame)
  }
  if (scrollFrame) {
    cancelAnimationFrame(scrollFrame)
  }
})
</script>

<style scoped>
/* Main Page Layout */
.agent-page-layout {
  display: flex;
  flex-direction: column;
  height: 100vh;
  overflow: hidden;
  background: #f7f7f8;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
}

.agent-main-container {
  display: flex;
  flex: 1;
  overflow: hidden;
  height: calc(100vh - 60px);
}

/* Sidebar Collapsible Styling - ChatGPT Light Style */
.sidebar {
  width: 260px;
  min-width: 260px;
  background: #171717;
  border-right: 1px solid rgba(255, 255, 255, 0.06);
  display: flex;
  flex-direction: column;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  position: relative;
  z-index: 10;
  overflow: hidden;
}

.sidebar.collapsed {
  width: 0;
  min-width: 0;
  border-right: none;
  transform: translateX(-260px);
  opacity: 0;
}

/* Sidebar Top Bar & Brand */
.sidebar-top-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 16px 8px 16px;
}

.sidebar-brand {
  font-size: 15px;
  font-weight: 600;
  color: #fafafa;
  display: flex;
  align-items: center;
  gap: 8px;
}

.brand-icon {
  color: #10a37f;
  font-size: 18px;
}

.sidebar-toggle-btn {
  color: #a1a1aa;
  padding: 6px;
  border-radius: 6px;
  transition: all 0.2s ease;
  cursor: pointer;
  background: transparent;
  border: none;
  display: flex;
  align-items: center;
  justify-content: center;
}

.sidebar-toggle-btn:hover {
  background: rgba(255, 255, 255, 0.08);
  color: #ffffff;
}

.sidebar-header {
  padding: 8px 16px 16px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
}

.new-chat-btn {
  width: 100%;
  border-radius: 8px;
  font-weight: 600;
  color: #f4f4f5;
  background: #2a2a2a;
  border: 1px solid rgba(255, 255, 255, 0.08);
  box-shadow: none;
  transition: all 0.2s ease;
}

.new-chat-btn:hover {
  background: #343434;
  border-color: rgba(255, 255, 255, 0.14);
  color: #ffffff;
}

.sessions-list-container {
  flex: 1;
  overflow-y: auto;
  padding: 12px 8px;
}

.empty-sessions {
  text-align: center;
  color: #a1a1aa;
  margin-top: 40px;
  font-size: 13px;
}

/* Session Items */
.session-item {
  padding: 10px 12px;
  border-radius: 8px;
  margin-bottom: 4px;
  cursor: pointer;
  background: transparent;
  transition: all 0.2s ease;
  border: 1px solid transparent;
}

.session-item:hover {
  background: rgba(255, 255, 255, 0.06);
}

.session-item.active {
  background: rgba(255, 255, 255, 0.08);
  border-color: rgba(255, 255, 255, 0.08);
}

.session-info {
  display: flex;
  flex-direction: column;
}

.session-title-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.session-title {
  color: #fafafa;
  font-weight: 600;
  font-size: 13.5px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  flex: 1;
}

.delete-session-btn {
  opacity: 0;
  transition: opacity 0.2s;
  padding: 2px;
  color: #a1a1aa !important;
}

.session-item:hover .delete-session-btn {
  opacity: 1;
}

.session-summary {
  color: #a1a1aa;
  font-size: 12px;
  margin: 4px 0 2px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.session-time {
  color: #737373;
  font-size: 10px;
  text-align: right;
}

/* Right Chat Window */
.chat-window {
  flex: 1;
  display: flex;
  flex-direction: column;
  background: #f7f7f8;
  position: relative;
  overflow: hidden;
}

/* Chat Header */
.chat-header {
  height: 60px;
  padding: 0 24px;
  border-bottom: 1px solid rgba(15, 23, 42, 0.06);
  display: flex;
  align-items: center;
  background: rgba(247, 247, 248, 0.85);
  justify-content: space-between;
  backdrop-filter: blur(10px);
}

.header-left-group {
  display: flex;
  align-items: center;
  gap: 16px;
}

.sidebar-open-btn {
  margin-right: -4px;
}

.chat-title-info {
  display: flex;
  align-items: center;
  gap: 10px;
}

.chat-title-info h2 {
  color: #0d0d0d;
  font-size: 15px;
  font-weight: 600;
  margin: 0;
}

.header-meta {
  display: flex;
  align-items: center;
  gap: 10px;
}

.status-text {
  color: #6b7280;
  font-size: 12px;
}

.model-pill {
  display: inline-flex;
  align-items: center;
  height: 28px;
  padding: 0 12px;
  border-radius: 999px;
  background: #ffffff;
  border: 1px solid rgba(15, 23, 42, 0.08);
  color: #111827;
  font-size: 12px;
  font-weight: 600;
}

/* Messages area - ChatGPT Centered Style */
.messages-container {
  flex: 1;
  overflow-y: auto;
  padding: 28px 0 172px;
  display: flex;
  flex-direction: column;
  gap: 28px;
  align-items: center;
}

/* Welcome screen style */
.welcome-screen {
  max-width: 768px;
  width: 90%;
  margin: 72px auto 40px;
  text-align: center;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.welcome-avatar {
  width: 56px;
  height: 56px;
  margin-bottom: 20px;
  display: grid;
  place-items: center;
  border-radius: 18px;
  background: #111827;
  color: #ffffff;
  box-shadow: 0 16px 32px rgba(15, 23, 42, 0.12);
}

.welcome-mark {
  font-size: 18px;
  font-weight: 700;
  letter-spacing: 0;
}

.welcome-screen h3 {
  color: #0d0d0d;
  font-size: 32px;
  font-weight: 600;
  margin: 0 0 10px;
}

.welcome-subtitle {
  margin: 0 0 28px;
  color: #6b7280;
  font-size: 14px;
  line-height: 1.6;
}

.quick-starts {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
  width: 100%;
  max-width: 768px;
}

.qs-card {
  background: rgba(255, 255, 255, 0.84);
  border: 1px solid rgba(15, 23, 42, 0.08);
  padding: 16px 18px;
  border-radius: 14px;
  text-align: left;
  cursor: pointer;
  transition: all 0.2s ease;
  display: flex;
  flex-direction: column;
  gap: 4px;
  box-shadow: 0 4px 14px rgba(15, 23, 42, 0.03);
}

.qs-card:hover {
  background: #ffffff;
  border-color: rgba(17, 24, 39, 0.12);
  transform: translateY(-1px);
}

.qs-card-top {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.qs-arrow {
  color: #9ca3af;
  font-size: 18px;
  line-height: 1;
}

.qs-title {
  font-size: 14px;
  font-weight: 600;
  color: #1f2937;
}

.qs-desc {
  font-size: 12.5px;
  color: #6b7280;
  line-height: 1.4;
}

/* Chat rows aligned to ChatGPT standard layout */
.message-row {
  width: 90%;
  max-width: 768px;
  display: flex;
  gap: 14px;
}

.message-row.user {
  justify-content: flex-end;
  padding-left: 15%;
}

.message-row.assistant {
  justify-content: flex-start;
  padding-right: 15%;
}

.avatar-col {
  flex-shrink: 0;
}

.avatar {
  width: 28px;
  height: 28px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 700;
  font-size: 11px;
}

.avatar.user {
  display: none; /* Hide user avatar to match clean ChatGPT modern layout */
}

.avatar.bot {
  background: #111827;
  color: #ffffff;
}

.assistant-mark {
  font-size: 11px;
  font-weight: 700;
}

.bubble-col {
  display: flex;
  flex-direction: column;
  gap: 6px;
  flex: 1;
}

.bubble-info {
  font-size: 12px;
  font-weight: 600;
  color: #6b7280;
  margin-bottom: 2px;
}

.message-row.user .bubble-info {
  display: none; /* Hide sender label for user to keep it clean */
}

.bubble {
  color: #0d0d0d;
  line-height: 1.7;
  font-size: 15px;
  word-break: break-word;
}

.message-row.user .bubble {
  background: #ffffff;
  border: 1px solid rgba(15, 23, 42, 0.08);
  padding: 10px 18px;
  border-radius: 20px;
  max-width: fit-content;
  align-self: flex-end;
  box-shadow: 0 2px 10px rgba(15, 23, 42, 0.03);
}

.message-row.assistant .bubble {
  background: transparent;
  padding: 0;
  border-radius: 0;
}

/* ChatGPT Code Block Theme Styling */
.gpt-code-wrapper {
  margin: 16px 0;
  border-radius: 8px;
  overflow: hidden;
  border: 1px solid #343541;
  background: #1e1e2e;
}

.gpt-code-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: #2f2f3d;
  padding: 8px 16px;
  font-family: monospace;
  font-size: 12px;
  color: #d1d5db;
  border-bottom: 1px solid #1e1e2e;
}

.gpt-code-lang {
  text-transform: lowercase;
  font-weight: 600;
}

.copy-code-btn {
  background: transparent;
  border: none;
  color: #a1a1aa;
  cursor: pointer;
  font-size: 12px;
  transition: color 0.15s;
}

.copy-code-btn:hover {
  color: #ffffff;
}

.code-block {
  margin: 0;
  padding: 16px;
  overflow-x: auto;
  background: #1e1e2e;
}

.code-block code {
  font-family: SFMono-Regular, Consolas, "Liberation Mono", Menlo, monospace;
  font-size: 13.5px;
  color: #f8f8f2;
  line-height: 1.5;
  background: transparent !important;
  padding: 0 !important;
}

/* Inline code styling */
.bubble :deep(code:not(.code-block code)) {
  font-family: SFMono-Regular, Consolas, monospace;
  background: rgba(0, 0, 0, 0.05);
  padding: 2px 6px;
  border-radius: 4px;
  font-size: 13.5px;
  color: #c2410c;
}

/* Markdown typography improvements */
.bubble :deep(strong) {
  font-weight: 700;
  color: #000000;
}

.bubble :deep(ul), .bubble :deep(ol) {
  margin: 8px 0;
  padding-left: 24px;
}

.bubble :deep(li) {
  margin: 6px 0;
}

/* Bouncing typing animation */
.typing-bubble {
  display: flex;
  gap: 4px;
  align-items: center;
  padding: 12px 20px;
}

.typing-bubble .dot {
  width: 6px;
  height: 6px;
  background: #9ca3af;
  border-radius: 50%;
  animation: bounce 1.3s infinite;
}

.typing-bubble .dot:nth-child(2) {
  animation-delay: 0.15s;
}

.typing-bubble .dot:nth-child(3) {
  animation-delay: 0.3s;
}

@keyframes bounce {
  0%, 60%, 100% { transform: translateY(0); }
  30% { transform: translateY(-4px); }
}

/* Chat Input Footer styling - floating ChatGPT style */
.chat-footer {
  padding: 12px 0 24px 0;
  background: linear-gradient(180deg, rgba(247, 247, 248, 0) 0%, #f7f7f8 38%);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  z-index: 5;
}

.chat-mode-row {
  width: 90%;
  max-width: 768px;
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 12px;
  flex-wrap: wrap;
  padding: 10px 14px;
  border-radius: 14px;
  background: rgba(255, 255, 255, 0.9);
  border: 1px solid rgba(15, 23, 42, 0.08);
  box-shadow: 0 12px 30px rgba(15, 23, 42, 0.05);
  backdrop-filter: blur(8px);
}

.chat-mode-label {
  font-size: 13px;
  font-weight: 600;
  color: #111827;
  white-space: nowrap;
}

.chat-mode-switch {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 4px;
  border-radius: 999px;
  background: #f3f4f6;
}

.chat-mode-btn {
  border: none;
  background: transparent;
  color: #6b7280;
  font-size: 13px;
  font-weight: 600;
  padding: 7px 14px;
  border-radius: 999px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.chat-mode-btn.active {
  background: #111827;
  color: #ffffff;
  box-shadow: 0 6px 14px rgba(17, 24, 39, 0.18);
}

.chat-mode-hint {
  font-size: 12px;
  color: #6b7280;
  flex: 1;
  min-width: 220px;
}

@media (max-width: 768px) {
  .chat-mode-row {
    gap: 10px;
    padding: 10px 12px;
  }

  .chat-mode-hint {
    min-width: 100%;
  }
}

/* ChatGPT Input container box */
.input-form {
  max-width: 768px;
  width: 90%;
  background: rgba(255, 255, 255, 0.96);
  border: 1px solid rgba(15, 23, 42, 0.08);
  border-radius: 24px;
  padding: 12px 14px 12px 18px;
  box-shadow: 0 18px 40px rgba(15, 23, 42, 0.08);
  display: flex;
  align-items: flex-end;
  gap: 12px;
  transition: border-color 0.2s, box-shadow 0.2s;
}

.input-form:focus-within {
  border-color: rgba(17, 24, 39, 0.16);
  box-shadow: 0 20px 45px rgba(15, 23, 42, 0.1);
}

.send-btn {
  height: 32px;
  width: 32px;
  border-radius: 50%;
  padding: 0 !important;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #000000;
  border: none;
  color: #ffffff;
  transition: all 0.2s ease;
  cursor: pointer;
}

.send-btn:hover:not(:disabled) {
  background: #27272a;
}

.send-btn:disabled {
  background: #e4e4e7 !important;
  color: #a1a1aa !important;
  cursor: not-allowed;
  border: none !important;
}

.send-btn :deep(i) {
  font-size: 16px;
  font-weight: bold;
}

/* Tool calling progress */
.tool-calling-status {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #4b5563;
  font-size: 13px;
  background: #f3f4f6;
  padding: 8px 14px;
  border-radius: 12px;
  align-self: flex-start;
  margin-top: 4px;
  border: 1px solid #e5e7eb;
}

.tool-calling-status .is-loading {
  animation: rotating 2s linear infinite;
}

@keyframes rotating {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

/* Approval card bubble */
.approval-card-bubble {
  background: #fffbeb !important;
  border: 1px solid #fef3c7 !important;
  border-radius: 16px !important;
  max-width: 420px;
}

.approval-card {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.approval-header {
  display: flex;
  align-items: center;
  gap: 8px;
}

.approval-header h4 {
  margin: 0;
  font-size: 14px;
  font-weight: 600;
  color: #d97706;
}

.approval-body {
  font-size: 13px;
  color: #3f3f46;
}

.action-desc {
  margin: 0 0 8px;
  font-weight: 500;
}

.action-details {
  background: #ffffff;
  border: 1px solid #f3f4f6;
  border-radius: 8px;
  padding: 12px;
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.detail-row {
  display: flex;
  justify-content: space-between;
}

.detail-row .label {
  color: #71717a;
}

.detail-row .val {
  font-weight: 600;
  color: #18181b;
}

.approval-actions {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  border-top: 1px solid #fef3c7;
  padding-top: 10px;
}

.approval-resolved-status {
  display: flex;
  justify-content: flex-end;
  border-top: 1px solid #fef3c7;
  padding-top: 10px;
}

.pay-after-approval-row {
  display: flex;
  align-items: center;
  gap: 10px;
}

.pay-btn-margin {
  margin-left: 5px;
}

.password-input-row {
  margin-top: 8px;
}

/* Input Box Action Buttons */
.input-action-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border-radius: 50%;
  border: none;
  background: transparent;
  color: #52525b;
  cursor: pointer;
  transition: all 0.2s ease;
  padding: 0;
}

.input-action-btn:hover {
  background-color: #f4f4f5;
  color: #18181b;
}

.input-action-btn :deep(i) {
  font-size: 18px;
}

.plus-btn {
  color: #71717a;
}

.mic-btn {
  color: #18181b;
}

/* Model Select trigger button inside input box */
.input-right-actions {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-left: auto;
  padding-bottom: 2px;
}

.input-model-pill {
  display: flex;
  align-items: center;
  padding: 6px 12px;
  background-color: #f4f4f5;
  border-radius: 999px;
  font-size: 13px;
  font-weight: 600;
  color: #374151;
}

.message-input {
  flex: 1;
}

.message-input :deep(.el-textarea__inner) {
  background: transparent !important;
  box-shadow: none !important;
  border: none !important;
  color: #0d0d0d !important;
  font-family: inherit;
  font-size: 15px;
  min-height: 24px !important;
  line-height: 1.6;
  padding: 0 !important;
  resize: none;
}
</style>
