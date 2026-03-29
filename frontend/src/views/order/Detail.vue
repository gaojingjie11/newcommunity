<template>
  <div class="order-detail-page">
    <Navbar />

    <div class="container">
      <el-page-header @back="$router.back()" content="订单详情" title="返回" class="page-header" />

      <div v-loading="loading">
        <template v-if="order">
          <el-card class="status-card">
            <div class="status-header">
              <div>
                <div class="status-text">{{ getStatusText(order.status) }}</div>
                <div class="store-name" v-if="order.store">门店：{{ order.store.name }}</div>
              </div>
              <div class="order-total">总金额 ￥{{ formatAmount(order.total_amount) }}</div>
            </div>
            <div class="status-actions">
              <el-button v-if="order.status === 0" type="danger" @click="payCurrentOrder">立即支付</el-button>
              <el-button v-if="order.status === 0" @click="cancelCurrentOrder">取消订单</el-button>
              <el-button v-if="order.status === 2" type="primary" @click="confirmReceipt">确认收货</el-button>
            </div>
          </el-card>

          <el-card class="section-card">
            <template #header>支付信息</template>
            <div class="info-row"><span>订单总额</span><strong>￥{{ formatAmount(order.total_amount) }}</strong></div>
            <div class="info-row"><span>积分抵扣</span><strong>{{ order.used_points || 0 }}</strong></div>
            <div class="info-row"><span>余额支付</span><strong>￥{{ formatAmount(order.used_balance || 0) }}</strong></div>
            <div class="info-row" v-if="order.paid_at"><span>支付时间</span><strong>{{ formatDate(order.paid_at) }}</strong></div>
          </el-card>

          <el-card class="section-card">
            <template #header>用户信息</template>
            <div class="info-row"><span>用户名</span><strong>{{ order.sys_user?.username || '-' }}</strong></div>
            <div class="info-row"><span>真实姓名</span><strong>{{ order.sys_user?.real_name || '-' }}</strong></div>
            <div class="info-row"><span>手机号</span><strong>{{ order.sys_user?.mobile || '-' }}</strong></div>
          </el-card>

          <el-card class="section-card">
            <template #header>商品信息</template>
            <div class="product-list">
              <div v-for="item in order.items" :key="item.id" class="product-item" @click="$router.push(`/product/${item.product_id}`)">
                <el-image :src="item.product?.image_url" fit="cover" class="product-img" />
                <div class="product-info">
                  <div class="product-name">{{ item.product?.name }}</div>
                  <div class="product-meta">
                    <span>￥{{ formatAmount(item.price) }}</span>
                    <span>x{{ item.quantity }}</span>
                  </div>
                </div>
              </div>
            </div>
          </el-card>
        </template>

        <el-empty v-else description="未找到订单信息" />
      </div>
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
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import dayjs from 'dayjs'
import { ElMessage, ElMessageBox } from 'element-plus'
import Navbar from '@/components/layout/Navbar.vue'
import PayAuthDialog from '@/components/payment/PayAuthDialog.vue'
import { getOrderDetail, payOrder, cancelOrder, receiveOrder } from '@/api/order'
import { useUserStore } from '@/stores/user'
import { GREEN_POINTS_PER_YUAN, getMixedPaymentPreview } from '@/utils/payment'

const route = useRoute()
const userStore = useUserStore()
const order = ref(null)
const loading = ref(false)
const showPayAuth = ref(false)
const paySubmitting = ref(false)

function formatAmount(value) {
  return Number(value || 0).toFixed(2)
}

function formatDate(date) {
  return dayjs(date).format('YYYY-MM-DD HH:mm:ss')
}

function getStatusText(status) {
  return {
    0: '待支付',
    1: '已支付',
    2: '已发货',
    3: '已完成',
    40: '已取消'
  }[status] || '未知'
}

function getPaymentPreview(amount) {
  return getMixedPaymentPreview(amount, userStore.userInfo.green_points)
}

async function fetchDetail() {
  loading.value = true
  try {
    order.value = await getOrderDetail(route.params.id)
  } catch (error) {
    ElMessage.error(error.response?.data?.msg || error.message || '获取订单详情失败')
  } finally {
    loading.value = false
  }
}

async function payCurrentOrder() {
  const preview = getPaymentPreview(order.value.total_amount)
  try {
    await ElMessageBox.confirm(
      `本次将按 ${GREEN_POINTS_PER_YUAN} 积分=1元优先抵扣 ${preview.points} 积分，余额支付 ￥${formatAmount(preview.balance)}，确认支付吗？`,
      '支付确认',
      { type: 'warning' }
    )
    showPayAuth.value = true
  } catch (error) {
    if (error === 'cancel' || error === 'close') {
      return
    }
    ElMessage.error(error.response?.data?.msg || error.message || '支付失败')
  }
}

async function submitOrderPay(authPayload) {
  if (!order.value) return

  paySubmitting.value = true
  try {
    const res = await payOrder({
      order_id: order.value.id,
      business_type: 1,
      ...authPayload
    })
    const paymentResult = res?.payment_result || res
    ElMessage.success(`支付成功，使用积分 ${paymentResult.used_points}，余额 ￥${formatAmount(paymentResult.used_balance)}`)
    showPayAuth.value = false
    await Promise.all([fetchDetail(), userStore.fetchUserInfo()])
  } catch (error) {
    ElMessage.error(error.response?.data?.msg || error.message || '支付失败')
  } finally {
    paySubmitting.value = false
  }
}

async function cancelCurrentOrder() {
  try {
    await ElMessageBox.confirm('确认取消该订单吗？', '提示', { type: 'warning' })
    await cancelOrder(order.value.id)
    ElMessage.success('订单已取消')
    await fetchDetail()
  } catch (error) {
    if (error !== 'cancel' && error !== 'close') {
      ElMessage.error(error.response?.data?.msg || error.message || '操作失败')
    }
  }
}

async function confirmReceipt() {
  try {
    await receiveOrder(order.value.id)
    ElMessage.success('确认收货成功')
    await fetchDetail()
  } catch (error) {
    ElMessage.error(error.response?.data?.msg || error.message || '操作失败')
  }
}

onMounted(async () => {
  await Promise.all([fetchDetail(), userStore.fetchUserInfo()])
})
</script>

<style scoped>
.order-detail-page {
  min-height: 100vh;
  padding-bottom: 40px;
}

.page-header,
.section-card,
.status-card {
  margin-bottom: 16px;
}

.status-header,
.info-row,
.product-meta {
  display: flex;
  justify-content: space-between;
  gap: 16px;
}

.status-text {
  font-size: 24px;
  font-weight: 700;
}

.order-total {
  font-size: 20px;
  font-weight: 700;
  color: var(--danger-color);
}

.status-actions {
  margin-top: 16px;
  display: flex;
  gap: 12px;
}

.info-row {
  margin-bottom: 10px;
}

.product-list {
  display: grid;
  gap: 12px;
}

.product-item {
  display: flex;
  gap: 16px;
  cursor: pointer;
}

.product-img {
  width: 80px;
  height: 80px;
  border-radius: 8px;
}

.product-info {
  flex: 1;
}
</style>
