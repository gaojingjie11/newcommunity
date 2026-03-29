<template>
  <div class="order-page">
    <Navbar />

    <div class="container">
      <h1 class="page-title">我的订单</h1>

      <el-alert type="success" show-icon :closable="false" class="tip-alert">
        <template #title>
          混合支付已启用（{{ GREEN_POINTS_PER_YUAN }} 积分 = 1 元）。系统会优先扣除绿色积分，不足部分再扣余额。
        </template>
      </el-alert>

      <div class="tabs-container">
        <el-tabs v-model="currentTab" @tab-click="handleTabClick">
          <el-tab-pane label="全部" name="all" />
          <el-tab-pane label="待支付" name="0" />
          <el-tab-pane label="已支付" name="1" />
          <el-tab-pane label="已发货" name="2" />
          <el-tab-pane label="已完成" name="3" />
        </el-tabs>
      </div>

      <div class="order-list" v-loading="loading">
        <template v-if="orders.length > 0">
          <el-card v-for="order in orders" :key="order.id" shadow="hover" class="order-item">
            <div class="order-header">
              <div>
                <div class="order-no">订单号：{{ order.order_no }}</div>
                <div class="order-time">{{ formatDate(order.created_at) }}</div>
              </div>
              <el-tag :type="getStatusType(order.status)">{{ getStatusText(order.status) }}</el-tag>
            </div>

            <div class="order-products" @click="$router.push(`/order/${order.id}`)">
              <div class="product-item" v-for="item in order.items" :key="item.id">
                <el-image :src="item.product?.image_url || 'https://via.placeholder.com/60'" fit="cover" class="thumb" />
                <div class="product-info">
                  <div class="product-name">{{ item.product?.name || '商品' }}</div>
                  <div class="product-meta">
                    <span>￥{{ formatAmount(item.price) }}</span>
                    <span>x{{ item.quantity }}</span>
                  </div>
                </div>
              </div>
            </div>

            <div class="payment-info">
              <span>总额 ￥{{ formatAmount(order.total_amount) }}</span>
              <span v-if="order.used_points">积分抵扣 {{ order.used_points }}</span>
              <span v-if="order.used_balance">余额支付 ￥{{ formatAmount(order.used_balance) }}</span>
            </div>

            <div class="order-footer">
              <div>共 {{ order.items.length }} 件商品</div>
              <div class="order-actions">
                <el-button
                  v-if="order.status === 0"
                  type="danger"
                  size="small"
                  :loading="loadingMap[order.id]"
                  @click.stop="payCurrentOrder(order)"
                >
                  立即支付
                </el-button>
                <el-button
                  v-if="order.status === 0"
                  size="small"
                  :loading="loadingMap[order.id]"
                  @click.stop="cancelCurrentOrder(order.id)"
                >
                  取消订单
                </el-button>
                <el-button
                  v-if="order.status === 2"
                  type="success"
                  size="small"
                  :loading="loadingMap[order.id]"
                  @click.stop="confirmReceipt(order.id)"
                >
                  确认收货
                </el-button>
              </div>
            </div>
          </el-card>

          <div class="pagination-container">
            <el-pagination
              v-model:current-page="page"
              v-model:page-size="size"
              :page-sizes="[10, 20, 50]"
              layout="total, sizes, prev, pager, next, jumper"
              :total="total"
              @size-change="fetchOrders"
              @current-change="fetchOrders"
            />
          </div>
        </template>

        <el-empty v-else description="暂无相关订单">
          <el-button type="primary" @click="$router.push('/mall')">去逛逛</el-button>
        </el-empty>
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
import { ref, onMounted, reactive, watch } from 'vue'
import dayjs from 'dayjs'
import { ElMessage, ElMessageBox } from 'element-plus'
import Navbar from '@/components/layout/Navbar.vue'
import PayAuthDialog from '@/components/payment/PayAuthDialog.vue'
import { getOrderList, payOrder, cancelOrder, receiveOrder } from '@/api/order'
import { useUserStore } from '@/stores/user'
import { GREEN_POINTS_PER_YUAN, getMixedPaymentPreview } from '@/utils/payment'

const userStore = useUserStore()
const orders = ref([])
const currentTab = ref('all')
const loading = ref(false)
const total = ref(0)
const page = ref(1)
const size = ref(10)
const loadingMap = reactive({})
const fetchToken = ref(0)

const showPayAuth = ref(false)
const paySubmitting = ref(false)
const pendingOrder = ref(null)

function formatAmount(value) {
  return Number(value || 0).toFixed(2)
}

function formatDate(date) {
  return dayjs(date).format('YYYY-MM-DD HH:mm')
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

function getStatusType(status) {
  return {
    0: 'warning',
    1: 'success',
    2: 'info',
    3: 'primary',
    40: 'danger'
  }[status] || 'info'
}

function getPaymentPreview(amount) {
  const currentPoints = Number(userStore.userInfo?.green_points || 0)
  return getMixedPaymentPreview(amount, currentPoints)
}

async function withLoading(id, fn) {
  if (loadingMap[id]) return
  loadingMap[id] = true
  try {
    await fn()
  } finally {
    loadingMap[id] = false
  }
}

function handleTabClick() {
  page.value = 1
}

async function payCurrentOrder(order) {
  try {
    const preview = getPaymentPreview(order.total_amount)
    await ElMessageBox.confirm(
      `本次将优先抵扣 ${preview.points} 积分，余额支付 ￥${formatAmount(preview.balance)}，确认继续吗？`,
      '订单支付确认',
      { type: 'warning' }
    )

    pendingOrder.value = order
    showPayAuth.value = true
  } catch (error) {
    if (error === 'cancel' || error === 'close') {
      return
    }
    ElMessage.error(error?.response?.data?.msg || error?.message || '支付失败')
  }
}

async function submitOrderPay(authPayload) {
  if (!pendingOrder.value) return
  const current = pendingOrder.value

  paySubmitting.value = true
  try {
    await withLoading(current.id, async () => {
      const res = await payOrder({
        order_id: current.id,
        business_type: 1,
        ...authPayload
      })

      const paymentResult = res?.payment_result || res
      ElMessage.success(`支付成功，使用积分 ${paymentResult.used_points}，余额 ￥${formatAmount(paymentResult.used_balance)}`)
      showPayAuth.value = false
      pendingOrder.value = null
      await Promise.all([fetchOrders(), userStore.fetchUserInfo()])
    })
  } catch (error) {
    ElMessage.error(error?.response?.data?.msg || error?.message || '支付失败')
  } finally {
    paySubmitting.value = false
  }
}

function cancelCurrentOrder(orderId) {
  withLoading(orderId, async () => {
    await ElMessageBox.confirm('确认取消该订单吗？', '提示', { type: 'warning' })
    await cancelOrder(orderId)
    ElMessage.success('订单已取消')
    await fetchOrders()
  }).catch(() => {})
}

function confirmReceipt(orderId) {
  withLoading(orderId, async () => {
    await receiveOrder(orderId)
    ElMessage.success('确认收货成功')
    await fetchOrders()
  })
}

async function fetchOrders() {
  const token = ++fetchToken.value
  loading.value = true
  try {
    const status = currentTab.value === 'all' ? undefined : Number(currentTab.value)
    const res = await getOrderList({
      status,
      page: page.value,
      size: size.value
    })
    if (token !== fetchToken.value) return
    orders.value = res.list || []
    total.value = res.total || 0
  } catch (error) {
    if (token !== fetchToken.value) return
    ElMessage.error(error.response?.data?.msg || error.message || '获取订单失败')
  } finally {
    if (token === fetchToken.value) {
      loading.value = false
    }
  }
}

onMounted(async () => {
  await Promise.all([fetchOrders(), userStore.fetchUserInfo()])
})

watch(currentTab, () => {
  page.value = 1
  fetchOrders()
})
</script>

<style scoped>
.order-page {
  min-height: 100vh;
  padding-bottom: 40px;
}

.tip-alert,
.tabs-container {
  margin-bottom: 20px;
}

.order-list {
  display: grid;
  gap: 16px;
}

.order-header,
.order-footer,
.payment-info,
.product-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.order-products {
  display: grid;
  gap: 12px;
  margin: 16px 0;
  cursor: pointer;
}

.product-item {
  display: flex;
  gap: 12px;
}

.thumb {
  width: 60px;
  height: 60px;
  border-radius: 8px;
}

.product-info {
  flex: 1;
}

.payment-info {
  flex-wrap: wrap;
  gap: 12px;
  color: var(--text-secondary);
  margin-bottom: 12px;
}

.order-actions {
  display: flex;
  gap: 8px;
}

.pagination-container {
  display: flex;
  justify-content: flex-end;
}
</style>
