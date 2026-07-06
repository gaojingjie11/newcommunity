<template>
  <div class="agent-page-layout">
    <Navbar />

    <div class="agent-shell">
      <aside :class="['sidebar', { collapsed: isSidebarCollapsed }]">
        <div class="sidebar-top">
          <div class="sidebar-brand">
            <span class="brand-mark">AI</span>
            <span class="brand-text">智能管家</span>
          </div>
          <el-button class="sidebar-icon-btn" link @click="isSidebarCollapsed = true">
            <el-icon size="18"><Fold /></el-icon>
          </el-button>
        </div>

        <div class="sidebar-actions">
          <el-button class="new-chat-btn" @click="handleCreateSession">
            <el-icon><EditPen /></el-icon>
            <span>新聊天</span>
          </el-button>
        </div>

        <div class="sidebar-section-title">对话</div>

        <div class="sessions-list">
          <div v-if="sessions.length === 0" class="empty-sessions">暂无对话历史</div>
          <button
            v-for="item in sessions"
            :key="item.id"
            type="button"
            :class="['session-item', { active: activeSessionId === item.id }]"
            @click="handleSelectSession(item.id)"
          >
            <div class="session-copy">
              <div class="session-title">{{ item.title || '新对话' }}</div>
              <div class="session-summary">{{ item.summary || '开始新的对话' }}</div>
            </div>

            <el-popconfirm
              title="确定删除此对话吗？"
              confirm-button-text="确定"
              cancel-button-text="取消"
              @confirm.stop="handleDeleteSession(item.id)"
            >
              <template #reference>
                <el-button
                  link
                  class="delete-session-btn"
                  @click.stop
                >
                  <el-icon><Delete /></el-icon>
                </el-button>
              </template>
            </el-popconfirm>
          </button>
        </div>

        <div class="sidebar-footer">
          <div class="sidebar-user">
            <span class="sidebar-user-avatar">{{ userInitial }}</span>
            <div class="sidebar-user-meta">
              <div class="sidebar-user-name">{{ userStore.userInfo?.nickname || userStore.userInfo?.mobile || '当前用户' }}</div>
              <div class="sidebar-user-plan">社区服务</div>
            </div>
          </div>
        </div>
      </aside>

      <main class="chat-stage">
        <header class="chat-header">
          <div class="chat-header-left">
            <el-button
              v-if="isSidebarCollapsed"
              class="sidebar-icon-btn"
              link
              @click="isSidebarCollapsed = false"
            >
              <el-icon size="18"><Expand /></el-icon>
            </el-button>
            <div class="chat-header-title">
              <h2>{{ activeSession?.title || '智能管家' }}</h2>
              <p>{{ isStreaming ? '正在生成回复' : '在线' }}</p>
            </div>
          </div>

          <div class="chat-header-right">
            <span class="header-model-pill">{{ currentModelLabel }}</span>
          </div>
        </header>

        <div ref="msgContainer" class="messages-container">
          <div v-if="messages.length === 0" class="empty-state">
            <div class="empty-state-badge">AI</div>
            <h1>今天想让智能管家帮您处理什么？</h1>
            <p>查公告、买东西、报修、跟进订单，都可以直接说。</p>

            <div class="suggestion-grid">
              <button
                v-for="qs in quickStarts"
                :key="qs.title"
                type="button"
                class="suggestion-card"
                @click="handleQuickStart(qs.prompt)"
              >
                <div class="suggestion-title">{{ qs.title }}</div>
                <div class="suggestion-desc">{{ qs.desc }}</div>
              </button>
            </div>
          </div>

          <div
            v-for="(msg, idx) in messages"
            :key="msg.id || idx"
            :class="['message-row', msg.role]"
          >
            <template v-if="msg.role === 'assistant'">
              <div class="assistant-avatar">AI</div>
              <div class="message-body">
                <div class="message-name">智能管家</div>

                <div
                  v-if="msg.content === '' && isStreaming && !msg.proposed_action"
                  class="assistant-message typing-bubble"
                >
                  <span class="dot"></span>
                  <span class="dot"></span>
                  <span class="dot"></span>
                </div>

                <div v-else-if="msg.proposed_action" class="assistant-message approval-card-bubble">
                  <div class="approval-card">
                    <div class="approval-header">
                      <el-icon color="#d97706" size="18"><Warning /></el-icon>
                      <h4>需要您的确认</h4>
                    </div>

                    <div class="approval-body">
                      <div v-if="msg.proposed_action.action_type === 'create_order'" class="action-details">
                        <p class="action-desc">AI 准备为您创建商城订单。</p>
                        <div class="detail-row">
                          <span class="label">商品ID</span>
                          <span class="val">{{ msg.proposed_action.payload.product_id }}</span>
                        </div>
                        <div class="detail-row">
                          <span class="label">数量</span>
                          <span class="val">{{ msg.proposed_action.payload.quantity }} 件</span>
                        </div>
                      </div>

                      <div v-else-if="msg.proposed_action.action_type === 'pay_order'" class="action-details">
                        <p class="action-desc">AI 准备对以下订单发起支付。</p>
                        <div class="detail-row">
                          <span class="label">订单ID</span>
                          <span class="val">{{ msg.proposed_action.payload.order_id }}</span>
                        </div>
                      </div>

                      <div v-else-if="msg.proposed_action.action_type === 'submit_repair'" class="action-details">
                        <p class="action-desc">AI 准备为您提交物业服务单。</p>
                        <div class="detail-row">
                          <span class="label">工单类别</span>
                          <span class="val">{{ msg.proposed_action.payload.type === 'repair' ? '报修' : '投诉' }}</span>
                        </div>
                        <div class="detail-row">
                          <span class="label">故障分类</span>
                          <span class="val">{{ msg.proposed_action.payload.category }}</span>
                        </div>
                        <div class="detail-row">
                          <span class="label">描述</span>
                          <span class="val">{{ msg.proposed_action.payload.description }}</span>
                        </div>
                      </div>
                    </div>

                    <div class="approval-actions" v-if="!msg.action_resolved">
                      <el-button size="small" plain @click="handleRejectAction(msg)">拒绝</el-button>
                      <el-button size="small" type="primary" :loading="msg.action_submitting" @click="handleApproveAction(msg)">
                        确认
                      </el-button>
                    </div>

                    <div class="approval-resolved-status" v-else>
                      <div
                        v-if="msg.action_resolved === 'approved' && msg.proposed_action.action_type === 'create_order' && getCreatedOrderId(msg)"
                        class="pay-after-approval-row"
                      >
                        <el-tag type="success" size="small">已执行</el-tag>
                        <el-tag
                          v-if="isOrderPaidOrResolved(getCreatedOrderId(msg))"
                          type="success"
                          size="small"
                        >
                          已支付
                        </el-tag>
                        <el-button
                          v-else
                          size="small"
                          type="warning"
                          @click="handlePayForCreatedOrder(getCreatedOrderId(msg), msg)"
                        >
                          去支付
                        </el-button>
                      </div>

                      <el-tag v-else :type="msg.action_resolved === 'approved' ? 'success' : 'danger'" size="small">
                        {{ msg.action_resolved === 'approved' ? '已同意授权' : '已拒绝授权' }}
                      </el-tag>
                    </div>
                  </div>
                </div>

                <div
                  v-else-if="msg.content"
                  class="assistant-message"
                  v-html="msg.rendered_content || ''"
                ></div>

                <div v-if="msg.tool_status" class="tool-calling-status">
                  <el-icon class="is-loading"><Loading /></el-icon>
                  <span>{{ msg.tool_status }}</span>
                </div>
              </div>
            </template>

            <template v-else>
              <div class="message-body user-body">
                <div class="user-message" v-html="msg.rendered_content || ''"></div>
              </div>
            </template>
          </div>
        </div>

        <footer class="composer-wrap">
          <div class="composer">
            <el-input
              v-model="inputMsg"
              type="textarea"
              :autosize="{ minRows: 1, maxRows: 8 }"
              resize="none"
              placeholder="给智能管家发消息"
              class="message-input"
              @keydown.enter.exact.prevent="handleSend"
            />

            <div class="composer-footer">
              <el-dropdown trigger="click" @command="handleModeCommand" popper-class="agent-mode-popper">
                <button type="button" class="mode-trigger">
                  <span>{{ currentModeMenuLabel }}</span>
                  <el-icon><ArrowDown /></el-icon>
                </button>
                <template #dropdown>
                  <el-dropdown-menu class="mode-menu">
                    <el-dropdown-item
                      v-for="mode in chatModes"
                      :key="mode.value"
                      :command="mode.value"
                      :class="{ 'is-active': chatMode === mode.value }"
                    >
                      <div class="mode-menu-row">
                        <div class="mode-menu-copy">
                          <span class="mode-menu-title">{{ mode.menuLabel || mode.label }}</span>
                          <span class="mode-menu-desc">{{ mode.hint }}</span>
                        </div>
                        <el-icon v-if="chatMode === mode.value" class="mode-check"><Check /></el-icon>
                      </div>
                    </el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>

              <div class="composer-actions">
                <span class="composer-hint">{{ currentModeHint }}</span>
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
          </div>
        </footer>
      </main>
    </div>

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
import { ref, computed, onMounted, onBeforeUnmount, nextTick } from "vue";
import { ElMessage } from "element-plus";
import {
  Delete,
  Loading,
  Warning,
  Top,
  Fold,
  Expand,
  ArrowDown,
  Check,
  EditPen,
} from "@element-plus/icons-vue";
import Navbar from "@/components/layout/Navbar.vue";
import PayAuthDialog from "@/components/payment/PayAuthDialog.vue";
import { useUserStore } from "@/stores/user";
import {
  getConversations,
  createConversation,
  deleteConversation,
  getChatHistory,
  chatStream,
  approveAction,
  rejectAction,
} from "@/api/chat";
import { payOrder, getOrderDetail } from "@/api/order";

const quickStarts = [
  { title: "查询公告", desc: "查看社区最近通知和公告", prompt: "查询最新社区公告" },
  { title: "推荐商品", desc: "看看商城里有什么日用品和饮料", prompt: "帮我推荐一些便利店商品" },
  { title: "物业报修", desc: "提交漏水、断电等物业服务需求", prompt: "家里卫生间水龙头漏水了，帮我报修" },
  { title: "订单支付", desc: "跟进待支付订单并继续完成支付", prompt: "帮我查看待支付订单" },
];

const chatModes = [
  {
    value: "fast",
    label: "极速",
    menuLabel: "极速",
    hint: "适合快速问答、轻量查询和日常对话。",
  },
  {
    value: "auto",
    label: "均衡",
    menuLabel: "智能",
    hint: "默认模式，速度和分析能力更均衡。",
  },
  {
    value: "deep",
    label: "高级",
    menuLabel: "高级",
    hint: "适合总结、分析、报表和复杂推理。",
  },
];

const isSidebarCollapsed = ref(false);
const sessions = ref([]);
const activeSessionId = ref("");
const inputMsg = ref("");
const messages = ref([]);
const isStreaming = ref(false);

const userStore = useUserStore();
const showPayAuth = ref(false);
const paySubmitting = ref(false);
const pendingPayMsg = ref(null);
const activeStreamingMessage = ref("");
const msgContainer = ref(null);
const payingOrderId = ref(null);
const payingOrderMsg = ref(null);
const chatMode = ref(localStorage.getItem("agent-chat-mode") || "auto");

let pendingAssistantChunk = "";
let assistantRenderFrame = 0;
let scrollFrame = 0;

const activeSession = computed(() => sessions.value.find((s) => s.id === activeSessionId.value));

const currentModelLabel = computed(() => {
  const item = chatModes.find((mode) => mode.value === chatMode.value);
  return item ? item.label : "均衡";
});

const currentModeMenuLabel = computed(() => {
  const item = chatModes.find((mode) => mode.value === chatMode.value);
  return item ? item.menuLabel || item.label : "智能";
});

const currentModeHint = computed(() => {
  const item = chatModes.find((mode) => mode.value === chatMode.value);
  return item ? item.hint : chatModes[1].hint;
});

const userInitial = computed(() => {
  const source = userStore.userInfo?.nickname || userStore.userInfo?.mobile || "AI";
  return String(source).slice(0, 1).toUpperCase();
});

const setChatMode = (mode) => {
  chatMode.value = mode;
  localStorage.setItem("agent-chat-mode", mode);
};

const handleModeCommand = (mode) => {
  setChatMode(mode);
};

const fetchSessions = async () => {
  try {
    const list = await getConversations();
    const rawSessions = Array.isArray(list) ? list : list?.list || [];
    sessions.value = rawSessions
      .map((item) => ({
        ...item,
        id: String(item?.id || "").trim(),
      }))
      .filter((item) => item.id);
    if (sessions.value.length > 0 && !activeSessionId.value) {
      handleSelectSession(sessions.value[0].id);
    } else if (sessions.value.length === 0) {
      await handleCreateSession();
    }
  } catch (err) {
    ElMessage.error("获取对话列表失败: " + err.message);
  }
};

const handleCreateSession = async () => {
  try {
    const res = await createConversation({ title: "新对话" });
    sessions.value.unshift({
      id: res.id,
      title: "新对话",
      summary: "",
      updated_at: new Date().toISOString(),
    });
    handleSelectSession(res.id);
  } catch (err) {
    ElMessage.error("创建新对话失败: " + err.message);
  }
};

const handleDeleteSession = async (id) => {
  if (!id || !String(id).trim()) {
    ElMessage.warning("该聊天记录编号异常，正在刷新列表");
    await fetchSessions();
    return;
  }

  try {
    await deleteConversation(id);
    sessions.value = sessions.value.filter((session) => session.id !== id);
    if (activeSessionId.value === id) {
      activeSessionId.value = "";
      messages.value = [];
      if (sessions.value.length > 0) {
        handleSelectSession(sessions.value[0].id);
      } else {
        await handleCreateSession();
      }
    }
    ElMessage.success("对话已成功删除");
  } catch (err) {
    ElMessage.error("删除对话失败: " + err.message);
  }
};

const handleSelectSession = async (id, clearMessages = true) => {
  if (!id || !String(id).trim()) {
    ElMessage.warning("该聊天记录编号异常，正在刷新列表");
    await fetchSessions();
    return;
  }

  activeSessionId.value = id;
  if (clearMessages) {
    messages.value = [];
  }

  try {
    const res = await getChatHistory(id);
    messages.value = (res?.list || []).map(normalizeMessage);

    messages.value.forEach((msg) => {
      if (msg.action_resolved === "approved" && msg.proposed_action?.action_type === "create_order") {
        const orderId = getCreatedOrderId(msg);
        if (orderId) {
          fetchOrderStatus(orderId);
        }
      }
    });

    nextTick(scheduleScrollToBottom);
  } catch (err) {
    ElMessage.error("加载聊天记录失败: " + err.message);
  }
};

const handleQuickStart = (text) => {
  inputMsg.value = text;
};

const normalizeMessage = (msg) => {
  const normalized = {
    ...msg,
    content: msg.content || "",
    rendered_content: msg.content ? renderMarkdown(msg.content) : "",
  };

  if (normalized.event_type === "approval_required" && normalized.event_payload) {
    try {
      normalized.proposed_action = JSON.parse(normalized.event_payload);
    } catch (error) {
      console.error("Failed to parse event_payload:", error);
    }
  }

  return normalized;
};

const renderMarkdown = (content) => {
  if (!content) return "";

  let html = content
    .replace(/&/g, "&amp;")
    .replace(/</g, "&lt;")
    .replace(/>/g, "&gt;");

  html = html.replace(/```([a-zA-Z0-9]+)?\n([\s\S]*?)```/g, (match, lang, code) => {
    const language = lang || "code";
    const escapedCode = encodeURIComponent(code.trim());
    return `<div class="gpt-code-wrapper">
      <div class="gpt-code-header">
        <span class="gpt-code-lang">${language}</span>
        <button class="copy-code-btn" onclick="navigator.clipboard.writeText(decodeURIComponent('${escapedCode}')).then(() => { event.target.innerText = '已复制'; setTimeout(() => event.target.innerText = '复制代码', 2000); })">复制代码</button>
      </div>
      <pre class="code-block"><code>${code.trim()}</code></pre>
    </div>`;
  });

  html = html.replace(/\*\*([\s\S]*?)\*\*/g, "<strong>$1</strong>");
  html = html.replace(/^\s*[-*]\s+(.+)$/gm, "<li>$1</li>");
  html = html.replace(/((?:<li>.*<\/li>\s*)+)/g, "<ul>$1</ul>");
  html = html.replace(/\n/g, "<br>");

  return html;
};

const scrollToBottom = () => {
  if (msgContainer.value) {
    msgContainer.value.scrollTop = msgContainer.value.scrollHeight;
  }
};

const scheduleScrollToBottom = () => {
  if (scrollFrame) return;
  scrollFrame = requestAnimationFrame(() => {
    scrollFrame = 0;
    scrollToBottom();
  });
};

const renderMessageContent = (msg) => {
  msg.rendered_content = msg.content ? renderMarkdown(msg.content) : "";
};

const flushAssistantChunks = (botMsgIdx) => {
  if (!pendingAssistantChunk) return;

  const botMsg = messages.value[botMsgIdx];
  if (!botMsg) {
    pendingAssistantChunk = "";
    return;
  }

  activeStreamingMessage.value += pendingAssistantChunk;
  pendingAssistantChunk = "";
  botMsg.content = activeStreamingMessage.value;
  renderMessageContent(botMsg);
  scheduleScrollToBottom();
};

const scheduleAssistantChunkFlush = (botMsgIdx) => {
  if (assistantRenderFrame) return;
  assistantRenderFrame = requestAnimationFrame(() => {
    assistantRenderFrame = 0;
    flushAssistantChunks(botMsgIdx);
  });
};

const handleSend = async () => {
  const query = inputMsg.value.trim();
  if (!query || isStreaming.value) return;

  if (!activeSessionId.value) {
    ElMessage.warning("未选中有效会话");
    return;
  }

  const userMsgId = "u-" + Date.now();
  messages.value.push({
    id: userMsgId,
    role: "user",
    content: query,
    rendered_content: renderMarkdown(query),
    created_at: new Date().toISOString(),
  });

  inputMsg.value = "";
  isStreaming.value = true;
  activeStreamingMessage.value = "";
  pendingAssistantChunk = "";
  nextTick(scheduleScrollToBottom);

  const payload = {
    conversation_id: activeSessionId.value,
    message: query,
    mode: chatMode.value,
    pay_type: "",
    payment_password: "",
    face_image_url: "",
  };

  const botMsgIdx = messages.value.length;
  messages.value.push({
    role: "assistant",
    content: "",
    rendered_content: "",
    tool_status: "智能管家正在思考中...",
    created_at: new Date().toISOString(),
  });

  chatStream(
    payload,
    (event) => {
      const botMsg = messages.value[botMsgIdx];
      if (!botMsg) return;

      if (event.type === "message_delta") {
        if (botMsg.tool_status === "智能管家正在思考中...") {
          botMsg.tool_status = "";
        }
        pendingAssistantChunk += event.data.chunk;
        scheduleAssistantChunkFlush(botMsgIdx);
      } else if (event.type === "tool_call_start") {
        let toolText = "智能管家正在处理业务...";
        if (event.data.tool === "list_products") {
          toolText = "正在查询商城商品列表...";
        } else if (event.data.tool === "query_notices") {
          toolText = "正在检索社区公告通知...";
        } else if (event.data.tool === "create_order") {
          toolText = "正在生成商品订单...";
        } else if (event.data.tool === "pay_order") {
          toolText = "正在发起订单支付...";
        } else if (event.data.tool === "submit_repair") {
          toolText = "正在提交物业报修单...";
        }
        botMsg.tool_status = toolText;
      } else if (event.type === "tool_call_end") {
        botMsg.tool_status = "";
      } else if (event.type === "approval_required") {
        flushAssistantChunks(botMsgIdx);
        isStreaming.value = false;
        botMsg.proposed_action = event.data;
      }
      scheduleScrollToBottom();
    },
    () => {
      flushAssistantChunks(botMsgIdx);
      isStreaming.value = false;
      fetchSessions();
    },
    (err) => {
      flushAssistantChunks(botMsgIdx);
      isStreaming.value = false;
      const botMsg = messages.value[botMsgIdx];
      if (botMsg) {
        botMsg.content = "⚠️ 发送错误: " + err.message;
        renderMessageContent(botMsg);
      }
      ElMessage.error("智能管家响应异常: " + err.message);
      scheduleScrollToBottom();
    }
  );
};

const orderStatuses = ref({});

const isOrderPaidOrResolved = (orderId) => {
  if (orderId === null || orderId === undefined) return false;
  const status = orderStatuses.value[orderId];
  return status !== undefined && status !== 0;
};

const fetchOrderStatus = async (orderId) => {
  if (!orderId || orderStatuses.value[orderId] !== undefined) return;
  try {
    const res = await getOrderDetail(orderId);
    if (res && res.status !== undefined) {
      orderStatuses.value[orderId] = res.status;
    }
  } catch (err) {
    console.error("Failed to fetch status for order", orderId, err);
  }
};

const getCreatedOrderId = (msg) => {
  if (!msg.result_payload) return null;
  try {
    const parsed = JSON.parse(msg.result_payload);
    return parsed.order_id;
  } catch (error) {
    return null;
  }
};

const handlePayForCreatedOrder = (orderId, msg) => {
  payingOrderId.value = orderId;
  payingOrderMsg.value = msg;
  showPayAuth.value = true;
};

const handleApproveAction = async (msg) => {
  if (msg.action_submitting) return;

  const actionId = msg.proposed_action.action_id;
  const actionType = msg.proposed_action.action_type;

  if (actionType === "pay_order") {
    pendingPayMsg.value = msg;
    showPayAuth.value = true;
    return;
  }

  msg.action_submitting = true;
  try {
    const res = await approveAction(activeSessionId.value, actionId, {});
    msg.action_resolved = "approved";
    if (actionType === "create_order" && res?.order_id) {
      msg.result_payload = JSON.stringify({ order_id: res.order_id });
      orderStatuses.value[res.order_id] = 0;
    }
    ElMessage.success("操作已成功授权执行");
    fetchSessions();
  } catch (err) {
    ElMessage.error("接口请求异常: " + err.message);
  } finally {
    msg.action_submitting = false;
  }
};

const submitOrderPay = async (authPayload) => {
  if (payingOrderId.value) {
    paySubmitting.value = true;
    const currentOrderId = payingOrderId.value;
    try {
      const res = await payOrder(currentOrderId, {
        pay_type: authPayload.pay_type,
        password: authPayload.password || "",
        face_image_url: authPayload.face_image_url || "",
        return_url: window.location.origin + "/payment/result",
      });

      showPayAuth.value = false;

      if (authPayload.pay_type === "alipay" && res?.pay_url) {
        ElMessage.success("正在为您跳转至支付宝收银台...");
        setTimeout(() => {
          window.location.href = res.pay_url;
        }, 800);
        return;
      }

      orderStatuses.value[currentOrderId] = 1;
      ElMessage.success("订单支付成功");
    } catch (err) {
      try {
        const detail = await getOrderDetail(currentOrderId);
        if (detail && detail.status !== undefined && detail.status !== 0) {
          showPayAuth.value = false;
          orderStatuses.value[currentOrderId] = detail.status;
          ElMessage.success("订单支付成功");
          return;
        }
      } catch (_) {}
      ElMessage.error("支付失败: " + err.message);
    } finally {
      paySubmitting.value = false;
      payingOrderId.value = null;
      payingOrderMsg.value = null;
    }
    return;
  }

  if (!pendingPayMsg.value) return;

  const msg = pendingPayMsg.value;
  const actionId = msg.proposed_action.action_id;

  paySubmitting.value = true;
  msg.action_submitting = true;
  try {
    const res = await approveAction(activeSessionId.value, actionId, {
      pay_type: authPayload.pay_type,
      payment_password: authPayload.password || "",
      face_image_url: authPayload.face_image_url || "",
      return_url: window.location.origin + "/payment/result",
    });

    msg.action_resolved = "approved";
    showPayAuth.value = false;
    pendingPayMsg.value = null;
    ElMessage.success("操作已成功授权执行");

    if (authPayload.pay_type === "alipay" && res?.pay_url) {
      ElMessage.success("正在为您跳转至支付宝收银台...");
      setTimeout(() => {
        window.location.href = res.pay_url;
      }, 800);
      return;
    }

    await handleSelectSession(activeSessionId.value, false);
  } catch (err) {
    if (msg.proposed_action?.action_type === "pay_order" && msg.proposed_action?.payload?.order_id) {
      try {
        const detail = await getOrderDetail(msg.proposed_action.payload.order_id);
        if (detail && detail.status !== undefined && detail.status !== 0) {
          msg.action_resolved = "approved";
          showPayAuth.value = false;
          pendingPayMsg.value = null;
          ElMessage.success("订单支付成功");
          await handleSelectSession(activeSessionId.value, false);
          return;
        }
      } catch (_) {}
    }
    ElMessage.error("接口请求异常: " + err.message);
  } finally {
    paySubmitting.value = false;
    msg.action_submitting = false;
  }
};

const handleRejectAction = async (msg) => {
  const actionId = msg.proposed_action.action_id;
  try {
    await rejectAction(activeSessionId.value, actionId);
    msg.action_resolved = "rejected";
    ElMessage.info("已成功取消授权");
    await handleSelectSession(activeSessionId.value, false);
  } catch (err) {
    ElMessage.error("操作异常: " + err.message);
  }
};

onMounted(async () => {
  await fetchSessions();
  try {
    await userStore.fetchUserInfo();
  } catch (err) {
    console.error("Failed to fetch user info:", err);
  }
});

onBeforeUnmount(() => {
  if (assistantRenderFrame) {
    cancelAnimationFrame(assistantRenderFrame);
  }
  if (scrollFrame) {
    cancelAnimationFrame(scrollFrame);
  }
});
</script>

<style scoped>
.agent-page-layout {
  display: flex;
  flex-direction: column;
  height: 100vh;
  overflow: hidden;
  background: #ffffff;
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
}

.agent-shell {
  display: flex;
  flex: 1;
  min-height: 0;
}

.sidebar {
  width: 300px;
  min-width: 300px;
  display: flex;
  flex-direction: column;
  background: #f7f7f5;
  border-right: 1px solid #ecebe8;
  transition: width 0.22s ease, min-width 0.22s ease, opacity 0.22s ease;
}

.sidebar.collapsed {
  width: 0;
  min-width: 0;
  opacity: 0;
  overflow: hidden;
  border-right: 0;
}

.sidebar-top {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 18px 18px 12px;
}

.sidebar-brand {
  display: flex;
  align-items: center;
  gap: 10px;
}

.brand-mark {
  width: 32px;
  height: 32px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border-radius: 11px;
  background: #111827;
  color: #ffffff;
  font-size: 13px;
  font-weight: 700;
}

.brand-text {
  font-size: 16px;
  font-weight: 600;
  color: #101828;
}

.sidebar-icon-btn {
  color: #6b7280;
  border-radius: 10px;
}

.sidebar-icon-btn:hover {
  background: #ecebe8;
  color: #111827;
}

.sidebar-actions {
  padding: 0 18px 14px;
}

.new-chat-btn {
  width: 100%;
  height: 42px;
  justify-content: flex-start;
  gap: 8px;
  border-radius: 14px;
  border: 1px solid #e4e3df;
  background: #ffffff;
  color: #111827;
  font-weight: 600;
  box-shadow: none;
}

.new-chat-btn:hover {
  background: #fbfbfa;
  border-color: #d8d6d1;
}

.sidebar-section-title {
  padding: 0 18px 8px;
  font-size: 12px;
  font-weight: 600;
  color: #8a8f98;
  text-transform: uppercase;
  letter-spacing: 0.04em;
}

.sessions-list {
  flex: 1;
  min-height: 0;
  overflow-y: auto;
  padding: 0 12px 16px;
}

.empty-sessions {
  padding: 22px 10px;
  color: #8a8f98;
  font-size: 13px;
  text-align: center;
}

.session-item {
  width: 100%;
  display: flex;
  align-items: flex-start;
  gap: 10px;
  padding: 12px;
  border: 0;
  border-radius: 14px;
  background: transparent;
  text-align: left;
  cursor: pointer;
  transition: background 0.18s ease;
}

.session-item:hover,
.session-item.active {
  background: #ecebe8;
}

.session-copy {
  flex: 1;
  min-width: 0;
}

.session-title {
  color: #111827;
  font-size: 14px;
  font-weight: 600;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.session-summary {
  margin-top: 4px;
  color: #7a7f87;
  font-size: 12px;
  line-height: 1.4;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.delete-session-btn {
  opacity: 0;
  color: #9ca3af !important;
}

.session-item:hover .delete-session-btn,
.session-item.active .delete-session-btn {
  opacity: 1;
}

.sidebar-footer {
  padding: 14px 18px 18px;
  border-top: 1px solid #ecebe8;
}

.sidebar-user {
  display: flex;
  align-items: center;
  gap: 10px;
}

.sidebar-user-avatar {
  width: 34px;
  height: 34px;
  border-radius: 50%;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  background: #d4b36f;
  color: #ffffff;
  font-size: 13px;
  font-weight: 700;
}

.sidebar-user-name {
  color: #111827;
  font-size: 14px;
  font-weight: 600;
}

.sidebar-user-plan {
  color: #8a8f98;
  font-size: 12px;
}

.chat-stage {
  position: relative;
  flex: 1;
  min-width: 0;
  min-height: 0;
  display: flex;
  flex-direction: column;
  background: #ffffff;
}

.chat-header {
  height: 60px;
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24px;
  background: rgba(255, 255, 255, 0.92);
  border-bottom: 1px solid #f1f1ef;
  backdrop-filter: blur(10px);
}

.chat-header-left,
.chat-header-right {
  display: flex;
  align-items: center;
  gap: 12px;
}

.chat-header-title h2 {
  margin: 0;
  font-size: 15px;
  font-weight: 600;
  color: #111827;
}

.chat-header-title p {
  margin: 2px 0 0;
  font-size: 12px;
  color: #8a8f98;
}

.header-model-pill {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 54px;
  height: 28px;
  padding: 0 12px;
  border-radius: 999px;
  background: #f5f5f3;
  color: #4b5563;
  font-size: 12px;
  font-weight: 600;
}

.messages-container {
  flex: 1;
  min-height: 0;
  overflow-y: auto;
  padding: 24px 0 172px;
}

.empty-state {
  max-width: 840px;
  margin: 72px auto 0;
  padding: 0 32px;
  text-align: center;
}

.empty-state-badge {
  width: 58px;
  height: 58px;
  margin: 0 auto 22px;
  display: grid;
  place-items: center;
  border-radius: 18px;
  background: #111827;
  color: #ffffff;
  font-size: 18px;
  font-weight: 700;
}

.empty-state h1 {
  margin: 0;
  color: #111827;
  font-size: 34px;
  font-weight: 600;
  line-height: 1.2;
}

.empty-state p {
  margin: 12px 0 28px;
  color: #6b7280;
  font-size: 15px;
  line-height: 1.6;
}

.suggestion-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px;
}

.suggestion-card {
  min-height: 96px;
  padding: 16px 18px;
  border: 1px solid #ecebe8;
  border-radius: 18px;
  background: #fafaf9;
  text-align: left;
  cursor: pointer;
  transition: border-color 0.18s ease, background 0.18s ease, transform 0.18s ease;
}

.suggestion-card:hover {
  background: #ffffff;
  border-color: #d8d6d1;
  transform: translateY(-1px);
}

.suggestion-title {
  color: #111827;
  font-size: 15px;
  font-weight: 600;
}

.suggestion-desc {
  margin-top: 6px;
  color: #6b7280;
  font-size: 13px;
  line-height: 1.5;
}

.message-row {
  max-width: 840px;
  margin: 0 auto;
  padding: 0 32px 24px;
}

.message-row.user {
  display: flex;
  justify-content: flex-end;
}

.message-row.assistant {
  display: flex;
  align-items: flex-start;
  gap: 14px;
}

.assistant-avatar {
  width: 32px;
  height: 32px;
  flex-shrink: 0;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border-radius: 11px;
  background: #111827;
  color: #ffffff;
  font-size: 12px;
  font-weight: 700;
  margin-top: 2px;
}

.message-body {
  flex: 1;
  min-width: 0;
}

.message-name {
  margin-bottom: 8px;
  color: #111827;
  font-size: 13px;
  font-weight: 600;
}

.assistant-message {
  color: #1f2937;
  font-size: 15px;
  line-height: 1.8;
  word-break: break-word;
}

.user-body {
  flex: none;
  max-width: min(72%, 620px);
}

.user-message {
  display: inline-block;
  padding: 14px 18px;
  border-radius: 26px;
  background: #f3f4f6;
  color: #111827;
  font-size: 15px;
  line-height: 1.6;
  word-break: break-word;
}

.typing-bubble {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  padding: 10px 4px;
}

.typing-bubble .dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: #9ca3af;
  animation: bounce 1.25s infinite;
}

.typing-bubble .dot:nth-child(2) {
  animation-delay: 0.15s;
}

.typing-bubble .dot:nth-child(3) {
  animation-delay: 0.3s;
}

@keyframes bounce {
  0%, 60%, 100% { transform: translateY(0); }
  30% { transform: translateY(-3px); }
}

.assistant-message :deep(code:not(.code-block code)),
.user-message :deep(code:not(.code-block code)) {
  font-family: SFMono-Regular, Consolas, monospace;
  background: rgba(15, 23, 42, 0.06);
  color: #c2410c;
  padding: 2px 6px;
  border-radius: 6px;
  font-size: 13px;
}

.assistant-message :deep(strong),
.user-message :deep(strong) {
  color: #111827;
  font-weight: 700;
}

.assistant-message :deep(ul),
.assistant-message :deep(ol),
.user-message :deep(ul),
.user-message :deep(ol) {
  margin: 8px 0;
  padding-left: 24px;
}

.assistant-message :deep(li),
.user-message :deep(li) {
  margin: 6px 0;
}

.gpt-code-wrapper {
  margin: 16px 0;
  overflow: hidden;
  border-radius: 12px;
  border: 1px solid #2f3540;
  background: #111827;
}

.gpt-code-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 14px;
  background: #1f2937;
  border-bottom: 1px solid #111827;
  color: #cbd5e1;
  font-size: 12px;
}

.gpt-code-lang {
  font-weight: 600;
}

.copy-code-btn {
  border: 0;
  background: transparent;
  color: #cbd5e1;
  cursor: pointer;
  font-size: 12px;
}

.code-block {
  margin: 0;
  padding: 16px;
  overflow-x: auto;
  background: #111827;
}

.code-block code {
  color: #f8fafc;
  font-family: SFMono-Regular, Consolas, "Liberation Mono", Menlo, monospace;
  font-size: 13px;
  line-height: 1.6;
}

.tool-calling-status {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  margin-top: 12px;
  padding: 8px 12px;
  border-radius: 999px;
  background: #f5f5f3;
  color: #4b5563;
  font-size: 13px;
}

.tool-calling-status .is-loading {
  animation: rotating 2s linear infinite;
}

@keyframes rotating {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

.approval-card-bubble {
  display: inline-block;
  max-width: 520px;
  margin-top: 4px;
  padding: 16px;
  border: 1px solid #fde68a;
  border-radius: 18px;
  background: #fffbeb;
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
  color: #b45309;
  font-size: 14px;
  font-weight: 600;
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
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding: 12px;
  border-radius: 12px;
  background: rgba(255, 255, 255, 0.85);
  border: 1px solid #f3f4f6;
}

.detail-row {
  display: flex;
  justify-content: space-between;
  gap: 16px;
}

.detail-row .label {
  color: #71717a;
}

.detail-row .val {
  color: #18181b;
  font-weight: 600;
  text-align: right;
}

.approval-actions,
.approval-resolved-status {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  padding-top: 10px;
  border-top: 1px solid #fde68a;
}

.pay-after-approval-row {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}

.composer-wrap {
  position: absolute;
  left: 0;
  right: 0;
  bottom: 0;
  padding: 20px 24px 28px;
  background: linear-gradient(180deg, rgba(255, 255, 255, 0) 0%, #ffffff 34%);
}

.composer {
  max-width: 840px;
  margin: 0 auto;
  padding: 14px 16px 12px;
  border: 1px solid #e7e5e4;
  border-radius: 28px;
  background: #ffffff;
  box-shadow: 0 12px 32px rgba(15, 23, 42, 0.08);
}

.message-input :deep(.el-textarea__inner) {
  border: 0 !important;
  box-shadow: none !important;
  padding: 0 !important;
  background: transparent !important;
  color: #111827 !important;
  font-size: 16px;
  line-height: 1.6;
  resize: none;
}

.composer-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-top: 10px;
}

.mode-trigger {
  height: 34px;
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 0 14px;
  border: 0;
  border-radius: 999px;
  background: #f3f4f6;
  color: #374151;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  transition: background 0.18s ease, color 0.18s ease;
}

.mode-trigger:hover {
  background: #e5e7eb;
  color: #111827;
}

.composer-actions {
  display: flex;
  align-items: center;
  gap: 12px;
  min-width: 0;
}

.composer-hint {
  color: #8a8f98;
  font-size: 12px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.send-btn {
  width: 38px;
  height: 38px;
  border-radius: 50%;
  border: 0;
  padding: 0 !important;
  background: #111827;
  color: #ffffff;
}

.send-btn:hover:not(:disabled) {
  background: #1f2937;
}

.send-btn:disabled {
  background: #e5e7eb !important;
  color: #9ca3af !important;
}

@media (max-width: 1100px) {
  .sidebar {
    width: 272px;
    min-width: 272px;
  }
}

@media (max-width: 900px) {
  .sidebar {
    position: absolute;
    top: 0;
    bottom: 0;
    left: 0;
    z-index: 20;
    box-shadow: 0 18px 40px rgba(15, 23, 42, 0.14);
  }

  .empty-state,
  .message-row {
    padding-left: 20px;
    padding-right: 20px;
  }

  .composer-wrap {
    padding-left: 16px;
    padding-right: 16px;
  }

  .composer-footer {
    flex-direction: column;
    align-items: stretch;
  }

  .composer-actions {
    justify-content: space-between;
  }

  .composer-hint {
    white-space: normal;
  }
}

@media (max-width: 640px) {
  .chat-header {
    padding: 0 16px;
  }

  .suggestion-grid {
    grid-template-columns: 1fr;
  }

  .empty-state h1 {
    font-size: 28px;
  }

  .message-row.assistant {
    gap: 10px;
  }

  .assistant-avatar {
    width: 28px;
    height: 28px;
    border-radius: 10px;
    font-size: 11px;
  }

  .user-body {
    max-width: 86%;
  }
}
</style>

<style>
.agent-mode-popper {
  border: 1px solid #ecebe8 !important;
  border-radius: 18px !important;
  padding: 6px !important;
  background: #ffffff !important;
  box-shadow: 0 18px 40px rgba(15, 23, 42, 0.12) !important;
}

.mode-menu {
  padding: 0 !important;
}

.mode-menu .el-dropdown-menu__item {
  padding: 0 !important;
  border-radius: 12px;
  margin: 2px 0;
  white-space: normal !important;
}

.mode-menu-row {
  width: 280px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding: 12px 14px;
}

.mode-menu-copy {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.mode-menu-title {
  color: #111827;
  font-size: 14px;
  font-weight: 600;
}

.mode-menu-desc {
  color: #6b7280;
  font-size: 12px;
  line-height: 1.45;
}

.mode-menu .el-dropdown-menu__item:hover,
.mode-menu .el-dropdown-menu__item:focus,
.mode-menu .el-dropdown-menu__item.is-active {
  background: #f7f7f5 !important;
}

.mode-check {
  color: #111827;
  font-size: 14px;
  font-weight: 700;
}
</style>
