<template>
  <div class="property-page">
    <Navbar />

    <div class="container custom-container">
      <div class="page-header">
        <h1 class="page-title highlight-title">物业费缴纳</h1>
      </div>

      <div class="premium-alert">
        <el-icon class="alert-icon"><InfoFilled /></el-icon>
        <div class="alert-content">
          当前绿色积分 <strong>{{ userStore.userInfo.green_points || 0 }}</strong>，账户余额 <strong>￥{{ formatAmount(userStore.userInfo.balance || 0) }}</strong>。
          支付时将按 <strong>{{ GREEN_POINTS_PER_YUAN }} 积分 = 1 元</strong> 自动优先扣除积分，不足部分再扣余额。
        </div>
      </div>

      <div class="fee-list">
        <div v-for="fee in fees" :key="fee.id" class="fee-card">
          <div class="fee-header">
            <div class="header-left">
              <div class="fee-month">{{ fee.month }}账单</div>
              <div class="fee-amount"><span class="currency">￥</span>{{ formatAmount(fee.amount) }}</div>
            </div>
            <div class="header-right">
               <span class="status-badge" :class="fee.status === 1 ? 'is-paid' : 'is-unpaid'">
                 {{ fee.status === 1 ? '已缴清' : '待缴费' }}
               </span>
            </div>
          </div>

          <div class="fee-body">
            <div class="mix-detail" v-if="fee.status === 0">
              <div class="detail-col">
                <span class="label">预计抵扣积分</span>
                <span class="value text-green">-{{ getPreview(fee.amount).points }}</span>
              </div>
              <div class="detail-col">
                <span class="label">预计扣余额</span>
                <span class="value">￥{{ formatAmount(getPreview(fee.amount).balance) }}</span>
              </div>
            </div>

            <div class="paid-detail" v-if="fee.status === 1">
              <div class="detail-col">
                <span class="label">使用积分</span>
                <span class="value">{{ fee.used_points || 0 }}</span>
              </div>
              <div class="detail-col">
                <span class="label">实扣余额</span>
                <span class="value">￥{{ formatAmount(fee.used_balance || 0) }}</span>
              </div>
              <div class="detail-row" v-if="fee.pay_time">
                支付时间：{{ formatDate(fee.pay_time) }}
              </div>
            </div>
          </div>

          <div class="fee-footer" v-if="fee.status === 0">
            <button class="btn-pay" :disabled="payingId === fee.id" @click="handlePay(fee)">
              {{ payingId === fee.id ? '支付中...' : '立即缴费' }}
            </button>
          </div>
        </div>
      </div>

      <div class="pagination-container" v-if="total > 0">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[5, 10, 20]"
          layout="total, sizes, prev, pager, next, jumper"
          :total="total"
          @size-change="fetchFees"
          @current-change="fetchFees"
          class="custom-pagination"
        />
      </div>
    </div>

    <PayAuthDialog
      v-model="showPayAuth"
      title="物业费支付验证"
      :face-registered="Boolean(userStore.userInfo?.face_registered)"
      :loading="paySubmitting"
      @confirm="submitFeePay"
    />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import dayjs from 'dayjs'
import { ElMessage, ElMessageBox } from 'element-plus'
import Navbar from '@/components/layout/Navbar.vue'
import PayAuthDialog from '@/components/payment/PayAuthDialog.vue'
import { getPropertyFeeList, payPropertyFee } from '@/api/service'
import { useUserStore } from '@/stores/user'
import { GREEN_POINTS_PER_YUAN, getMixedPaymentPreview } from '@/utils/payment'
import { InfoFilled } from '@element-plus/icons-vue'

const userStore = useUserStore()
const fees = ref([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)
const payingId = ref(0)
const showPayAuth = ref(false)
const paySubmitting = ref(false)
const pendingFee = ref(null)

function formatAmount(value) {
  return Number(value || 0).toFixed(2)
}

function formatDate(value) {
  return dayjs(value).format('YYYY-MM-DD HH:mm')
}

function getPreview(amount) {
  return getMixedPaymentPreview(amount, userStore.userInfo.green_points)
}

async function fetchFees() {
  const res = await getPropertyFeeList({
    page: currentPage.value,
    size: pageSize.value
  })
  fees.value = res.list || []
  total.value = res.total || 0
}

async function handlePay(fee) {
  try {
    const preview = getPreview(fee.amount)
    await ElMessageBox.confirm(
      `本次将按 ${GREEN_POINTS_PER_YUAN} 积分=1元优先抵扣 ${preview.points} 积分，余额扣除 ￥${formatAmount(preview.balance)}，确认支付吗？`,
      '混合支付确认',
      { type: 'warning', confirmButtonText: '确认缴费', cancelButtonText: '暂不缴费' }
    )

    pendingFee.value = fee
    showPayAuth.value = true
  } catch (error) {
    if (error === 'cancel' || error === 'close') {
      return
    }
    ElMessage.error(error.response?.data?.msg || error.message || '支付失败')
  }
}

async function submitFeePay(authPayload) {
  if (!pendingFee.value) return

  paySubmitting.value = true
  payingId.value = pendingFee.value.id
  try {
    const res = await payPropertyFee(pendingFee.value.id, authPayload)
    const paymentResult = res?.payment_result || res
    ElMessage.success(`支付成功，使用积分 ${paymentResult.used_points}，余额 ￥${formatAmount(paymentResult.used_balance)}`)
    showPayAuth.value = false
    pendingFee.value = null
    await Promise.all([fetchFees(), userStore.fetchUserInfo()])
  } catch (error) {
    ElMessage.error(error.response?.data?.msg || error.message || '支付失败')
  } finally {
    paySubmitting.value = false
    payingId.value = 0
  }
}

onMounted(async () => {
  await Promise.all([fetchFees(), userStore.fetchUserInfo()])
})
</script>

<style scoped>
.property-page { min-height: 100vh; background-color: #f8f9fa; padding-bottom: 80px; }
.custom-container { max-width: 900px; margin: 0 auto; }
.page-header { padding: 32px 0 24px; }
.highlight-title { display: inline-block; position: relative; font-size: 32px; color: #2c3e50; font-weight: 700; margin: 0; z-index: 1; }
.highlight-title::after { content: ''; position: absolute; bottom: 4px; left: -5%; width: 110%; height: 14px; background-color: #2d597b; opacity: 0.15; border-radius: 6px; z-index: -1; }

.premium-alert { display: flex; align-items: flex-start; background: #f0f7ff; border: 1px solid #cce3f6; border-radius: 8px; padding: 16px 20px; margin-bottom: 32px; color: #2d597b; }
.alert-icon { font-size: 20px; margin-right: 12px; margin-top: 2px; }
.alert-content { font-size: 14px; line-height: 1.6; }
.alert-content strong { color: #e4393c; font-size: 15px; margin: 0 2px; }

.fee-list { display: flex; flex-direction: column; gap: 20px; }

.fee-card { background: #ffffff; border-radius: 16px; padding: 24px 32px; box-shadow: 0 4px 20px rgba(0,0,0,0.03); border: 1px solid transparent; transition: all 0.3s; }
.fee-card:hover { transform: translateY(-2px); box-shadow: 0 8px 24px rgba(0,0,0,0.06); border-color: rgba(45,89,123,0.1); }

.fee-header { display: flex; justify-content: space-between; align-items: flex-start; border-bottom: 1px dashed #ebeef5; padding-bottom: 20px; margin-bottom: 20px; }
.fee-month { font-size: 20px; font-weight: 700; color: #2c3e50; margin-bottom: 8px; }
.fee-amount { font-size: 24px; font-weight: 800; color: #e4393c; font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; }
.currency { font-size: 16px; margin-right: 2px; }

.status-badge { padding: 6px 16px; border-radius: 20px; font-size: 13px; font-weight: bold; }
.is-unpaid { background: #fff7ed; color: #d97706; }
.is-paid { background: #f0fdf4; color: #166534; }

.mix-detail { display: flex; gap: 40px; background: #fbfcfd; padding: 16px 24px; border-radius: 8px; border: 1px solid #f0f2f5; }
.paid-detail { display: flex; flex-wrap: wrap; gap: 24px; background: #f0fdf4; padding: 16px 24px; border-radius: 8px; border: 1px solid #dcfce7; }
.detail-col { display: flex; flex-direction: column; gap: 6px; }
.label { font-size: 13px; color: #909399; }
.value { font-size: 16px; font-weight: 600; color: #303133; }
.text-green { color: #00b894; }
.detail-row { width: 100%; font-size: 13px; color: #166534; margin-top: 8px; padding-top: 12px; border-top: 1px dashed #bbf7d0; }

.fee-footer { margin-top: 24px; text-align: right; }
.btn-pay { background: #2d597b; color: #fff; border: none; border-radius: 20px; padding: 10px 32px; font-size: 15px; font-weight: bold; cursor: pointer; transition: all 0.3s; box-shadow: 0 4px 12px rgba(45,89,123,0.2); }
.btn-pay:hover:not(:disabled) { background: #1f435d; transform: translateY(-1px); }
.btn-pay:disabled { background: #c0c4cc; cursor: not-allowed; }

.pagination-container { display: flex; justify-content: center; margin-top: 40px; }
:deep(.custom-pagination .el-pager li.is-active) { background-color: #2d597b; color: #fff; border-radius: 4px; }
</style>