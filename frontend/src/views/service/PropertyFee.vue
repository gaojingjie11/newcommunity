<template>
  <div class="property-page">
    <Navbar />

    <div class="container">
      <h1 class="page-title">物业费缴纳</h1>

      <el-alert type="success" show-icon :closable="false" class="summary-alert">
        <template #title>
          当前绿色积分 {{ userStore.userInfo.green_points || 0 }}，账户余额 ￥{{ formatAmount(userStore.userInfo.balance || 0) }}。
          支付时将按 {{ GREEN_POINTS_PER_YUAN }} 积分 = 1 元自动优先扣除积分，不足部分再扣余额。
        </template>
      </el-alert>

      <div class="fee-list">
        <el-card v-for="fee in fees" :key="fee.id" class="fee-card">
          <div class="fee-header">
            <div>
              <div class="fee-month">{{ fee.month }}</div>
              <div class="fee-meta">应缴金额 ￥{{ formatAmount(fee.amount) }}</div>
            </div>
            <el-tag :type="fee.status === 1 ? 'success' : 'warning'">
              {{ fee.status === 1 ? '已缴费' : '待缴费' }}
            </el-tag>
          </div>

          <div class="mix-detail">
            <span>预计抵扣积分：{{ getPreview(fee.amount).points }}</span>
            <span>预计扣余额：￥{{ formatAmount(getPreview(fee.amount).balance) }}</span>
          </div>

          <div class="paid-detail" v-if="fee.status === 1">
            <span>实际使用积分：{{ fee.used_points || 0 }}</span>
            <span>实际使用余额：￥{{ formatAmount(fee.used_balance || 0) }}</span>
            <span v-if="fee.pay_time">支付时间：{{ formatDate(fee.pay_time) }}</span>
          </div>

          <div class="actions" v-if="fee.status === 0">
            <el-button type="primary" :loading="payingId === fee.id" @click="handlePay(fee)">
              立即缴费
            </el-button>
          </div>
        </el-card>
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
        />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import dayjs from 'dayjs'
import { ElMessage, ElMessageBox } from 'element-plus'
import Navbar from '@/components/layout/Navbar.vue'
import { getPropertyFeeList, payPropertyFee } from '@/api/service'
import { useUserStore } from '@/stores/user'
import { GREEN_POINTS_PER_YUAN, getMixedPaymentPreview } from '@/utils/payment'

const userStore = useUserStore()
const fees = ref([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)
const payingId = ref(0)

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
      { type: 'warning' }
    )

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

    const password = (value || '').trim()
    if (!password) {
      ElMessage.warning('未输入登录密码，已取消支付')
      return
    }

    payingId.value = fee.id
    const result = await payPropertyFee(fee.id, password)
    ElMessage.success(`支付成功，使用积分 ${result.used_points}，余额 ￥${formatAmount(result.used_balance)}`)
    await Promise.all([fetchFees(), userStore.fetchUserInfo()])
  } catch (error) {
    if (error === 'cancel' || error === 'close') {
      return
    }
    ElMessage.error(error.response?.data?.msg || error.message || '支付失败')
  } finally {
    payingId.value = 0
  }
}

onMounted(async () => {
  await Promise.all([fetchFees(), userStore.fetchUserInfo()])
})
</script>

<style scoped>
.property-page {
  min-height: 100vh;
  padding-bottom: 40px;
}

.summary-alert {
  margin-bottom: 20px;
}

.fee-list {
  display: grid;
  gap: 16px;
}

.fee-card {
  border-radius: 16px;
}

.fee-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.fee-month {
  font-size: 20px;
  font-weight: 600;
}

.fee-meta,
.mix-detail,
.paid-detail {
  display: flex;
  gap: 16px;
  flex-wrap: wrap;
  margin-top: 12px;
  color: var(--text-secondary);
}

.actions {
  margin-top: 16px;
}

.pagination-container {
  display: flex;
  justify-content: flex-end;
  margin-top: 24px;
}
</style>
