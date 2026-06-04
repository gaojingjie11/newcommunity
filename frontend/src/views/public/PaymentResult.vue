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
        
        <h2 class="result-title">{{ success ? '支付成功' : '支付确认中' }}</h2>
        <p class="result-subtitle">{{ success ? '感谢您的支付，资金已安全入账' : '正在确认您的支付状态，请稍后' }}</p>
        
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
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import Navbar from '@/components/layout/Navbar.vue'
import { SuccessFilled, CircleCloseFilled } from '@element-plus/icons-vue'

const route = useRoute()
const router = useRouter()

const loading = ref(false)
const success = ref(true)

const outTradeNo = computed(() => route.query.out_trade_no || '')
const tradeNo = computed(() => route.query.trade_no || '')
const amount = computed(() => route.query.total_amount || '')

const isRecharge = computed(() => {
  return outTradeNo.value.startsWith('RECH_')
})

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
