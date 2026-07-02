<template>
  <div class="payment-result-page">
    <Navbar />
    
    <div class="result-container">
      <div class="result-card" v-loading="loading">
        <div class="status-icon success" v-if="success">
          <el-icon><SuccessFilled /></el-icon>
        </div>
        <div class="status-icon error" v-else>
          <el-icon><CircleCloseFilled /></el-icon>
        </div>
        
        <h2 class="result-title">{{ statusText }}</h2>
        <p class="result-subtitle">{{ subtitleText }}</p>
        
        <div class="detail-list">
          <div class="detail-item" v-if="outTradeNo">
            <span class="label">交易单号</span>
            <span class="value">{{ outTradeNo }}</span>
          </div>
          <div class="detail-item" v-if="amount">
            <span class="label">交易金额</span>
            <span class="value amount-value">¥{{ parseFloat(amount).toFixed(2) }}</span>
          </div>
          <div class="detail-item" v-if="tradeNo">
            <span class="label">支付宝交易号</span>
            <span class="value">{{ tradeNo }}</span>
          </div>
          <div class="detail-item">
            <span class="label">交易类型</span>
            <span class="value">{{ isRecharge ? '余额充值' : '订单支付' }}</span>
          </div>
          <div class="detail-item" v-if="!success && checked">
            <span class="label">当前状态</span>
            <span class="value">等待系统确认</span>
          </div>
        </div>
        
        <div class="actions">
          <button class="action-btn main-btn" @click="goNext">
            {{ isRecharge ? '查看钱包账单' : '查看我的订单' }}
          </button>
          <button class="action-btn sub-btn" @click="goHome">返回首页</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import Navbar from '@/components/layout/Navbar.vue'
import { SuccessFilled, CircleCloseFilled } from '@element-plus/icons-vue'
import { getTransactionList } from '@/api/finance'
import { getOrderDetail } from '@/api/order'

const route = useRoute()
const router = useRouter()

const loading = ref(false)
const success = ref(false)
const checked = ref(false)
const statusText = ref('支付确认中')
const subtitleText = ref('正在确认您的支付状态，请稍后')
const pollIntervalMs = 2500
const maxPollAttempts = 8
let pollTimer = null
let pollAttempts = 0

const outTradeNo = computed(() => route.query.out_trade_no || '')
const tradeNo = computed(() => route.query.trade_no || '')
const amount = computed(() => route.query.total_amount || '')
const orderId = computed(() => Number(route.query.order_id || 0))

const isRecharge = computed(() => {
  return outTradeNo.value.startsWith('RECH_')
})

const stopPolling = () => {
  if (pollTimer) {
    window.clearTimeout(pollTimer)
    pollTimer = null
  }
}

const setPendingState = (title = '支付确认中', subtitle = '正在确认您的支付状态，请稍后') => {
  success.value = false
  statusText.value = title
  subtitleText.value = subtitle
}

const setSuccessState = (title = '支付成功', subtitle = '感谢您的支付，资金已安全入账') => {
  success.value = true
  checked.value = true
  statusText.value = title
  subtitleText.value = subtitle
  stopPolling()
}

const scheduleRetry = () => {
  if (pollAttempts >= maxPollAttempts) {
    checked.value = true
    setPendingState('支付确认中', '系统还在和支付结果对账，稍后刷新本页或前往业务页面查看最新状态')
    return
  }
  pollAttempts += 1
  pollTimer = window.setTimeout(() => {
    verifyPaymentResult()
  }, pollIntervalMs)
}

const confirmRechargeResult = async () => {
  const txResp = await getTransactionList({ page: 1, size: 100 })
  const list = Array.isArray(txResp?.list) ? txResp.list : Array.isArray(txResp) ? txResp : []
  const matched = list.find((item) => item?.biz_id === outTradeNo.value)
  if (matched && Number(matched.amount) > 0) {
    setSuccessState('充值成功', '充值金额已进入账户余额，可在钱包账单中查看')
    return true
  }
  return false
}

const confirmOrderResult = async () => {
  if (!orderId.value) {
    return false
  }
  const detail = await getOrderDetail(orderId.value)
  const status = Number(detail?.status)
  if (status >= 1 && status !== 40) {
    setSuccessState('支付成功', '订单已完成支付，稍后可以在我的订单中查看')
    return true
  }
  if (status === 40) {
    checked.value = true
    setPendingState('支付未完成', '订单已取消或已过期，请返回订单页重新发起支付')
    stopPolling()
    return true
  }
  return false
}

const verifyPaymentResult = async () => {
  if (!outTradeNo.value) {
    checked.value = true
    setPendingState('缺少交易信息', '未获取到有效的支付回跳参数，请返回业务页面查看最新状态')
    return
  }

  loading.value = true
  try {
    const confirmed = isRecharge.value
      ? await confirmRechargeResult()
      : await confirmOrderResult()

    if (!confirmed) {
      setPendingState()
      scheduleRetry()
    }
  } catch (error) {
    if (pollAttempts >= maxPollAttempts) {
      checked.value = true
      setPendingState('支付确认中', error?.message || '暂时无法确认支付结果，请稍后重试')
      return
    }
    scheduleRetry()
  } finally {
    loading.value = false
  }
}

const goNext = () => {
  if (isRecharge.value) {
    router.push('/user/transactions')
  } else {
    router.push('/order')
  }
}

const goHome = () => {
  router.push('/home')
}

onMounted(() => {
  verifyPaymentResult()
})

onUnmounted(() => {
  stopPolling()
})
</script>

<style scoped>
.payment-result-page {
  min-height: 100vh;
  background-color: #f8f9fa;
  padding-bottom: 60px;
}

.result-container {
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 80px 20px;
}

.result-card {
  background: #ffffff;
  border-radius: 20px;
  width: 100%;
  max-width: 500px;
  padding: 40px;
  text-align: center;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.05);
  border: 1px solid #ebeef5;
}

.status-icon {
  font-size: 72px;
  margin-bottom: 24px;
  display: inline-flex;
  justify-content: center;
  align-items: center;
}

.status-icon.success {
  color: #67c23a;
}

.status-icon.error {
  color: #f56c6c;
}

.result-title {
  font-size: 24px;
  color: #2c3e50;
  font-weight: 700;
  margin: 0 0 8px 0;
}

.result-subtitle {
  font-size: 14px;
  color: #909399;
  margin: 0 0 32px 0;
}

.detail-list {
  background: #fcfcfd;
  border: 1px dashed #e4e7ed;
  border-radius: 12px;
  padding: 20px;
  margin-bottom: 32px;
  text-align: left;
}

.detail-item {
  display: flex;
  justify-content: space-between;
  margin-bottom: 12px;
  font-size: 14px;
}

.detail-item:last-child {
  margin-bottom: 0;
}

.detail-item .label {
  color: #909399;
}

.detail-item .value {
  color: #303133;
  font-weight: 500;
}

.detail-item .amount-value {
  color: #2d597b;
  font-weight: 700;
  font-size: 16px;
}

.actions {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.action-btn {
  width: 100%;
  padding: 14px 0;
  border-radius: 10px;
  font-size: 15px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  border: 1px solid transparent;
  outline: none;
  display: block;
  box-sizing: border-box;
  margin: 0;
}

.main-btn {
  background-color: #2d597b;
  color: #ffffff;
  border-color: #2d597b;
}

.main-btn:hover {
  background-color: #1e3f5a;
  border-color: #1e3f5a;
}

.sub-btn {
  background-color: #ffffff;
  color: #606266;
  border: 1px solid #dcdfe6;
}

.sub-btn:hover {
  color: #2d597b;
  border-color: #c6e2ff;
  background-color: #ecf5ff;
}
</style>
