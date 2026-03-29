<template>
  <div class="order-page">
    <Navbar />

    <div class="container custom-container">
      <div class="page-header">
        <h1 class="page-title highlight-title">我的订单</h1>
      </div>

      <!-- 高级提示横幅 -->
      <div class="premium-alert">
        <el-icon class="alert-icon"><InfoFilled /></el-icon>
        <div class="alert-content">
          混合支付已启用（<strong>{{ GREEN_POINTS_PER_YUAN }} 积分 = 1 元</strong>）。系统会优先扣除绿色积分，不足部分再扣余额。
        </div>
      </div>

      <!-- 深度美化的分类 Tabs -->
      <div class="tabs-container">
        <el-tabs v-model="currentTab" @tab-click="handleTabClick" class="custom-tabs">
          <el-tab-pane label="全部订单" name="all" />
          <el-tab-pane label="待支付" name="0" />
          <el-tab-pane label="已支付" name="1" />
          <el-tab-pane label="已发货" name="2" />
          <el-tab-pane label="已完成" name="3" />
        </el-tabs>
      </div>

      <div class="order-list" v-loading="loading">
        <template v-if="orders.length > 0">
          <div v-for="order in orders" :key="order.id" class="order-card" @click="$router.push(`/order/${order.id}`)">
            
            <!-- 订单头部 -->
            <div class="order-header">
              <div class="header-left">
                <span class="order-no">订单号：{{ order.order_no }}</span>
                <span class="order-time">{{ formatDate(order.created_at) }}</span>
              </div>
              <div class="header-right">
                <el-tag :type="getStatusType(order.status)" class="custom-tag" effect="light" round>
                  {{ getStatusText(order.status) }}
                </el-tag>
              </div>
            </div>

            <!-- 订单商品 (新增画中画背景) -->
            <div class="order-products">
              <div class="product-item" v-for="item in order.items" :key="item.id">
                <div class="thumb-wrapper">
                  <img :src="item.product?.image_url || 'https://via.placeholder.com/80'" class="thumb" />
                </div>
                <div class="product-info">
                  <div class="product-name">{{ item.product?.name || '未知商品' }}</div>
                  <div class="product-meta">
                    <span class="price"><span class="currency">￥</span>{{ formatAmount(item.price) }}</span>
                    <span class="qty">x{{ item.quantity }}</span>
                  </div>
                </div>
              </div>
            </div>

            <!-- 支付明细与操作底栏 -->
            <div class="order-footer">
              <div class="payment-info">
                <div class="payment-details">
                  <span v-if="order.used_points" class="detail-item point-item">积分抵扣: -{{ order.used_points }}</span>
                  <span v-if="order.used_balance" class="detail-item balance-item">余额支付: ￥{{ formatAmount(order.used_balance) }}</span>
                </div>
                <div class="total-amount">
                  共 {{ order.items?.length || 0 }} 件商品，总额：<span class="highlight-price"><span class="currency">￥</span>{{ formatAmount(order.total_amount) }}</span>
                </div>
              </div>
              
              <div class="order-actions">
                <!-- 待支付状态操作 -->
                <template v-if="order.status === 0">
                  <button
                    class="action-btn btn-cancel"
                    :disabled="loadingMap[order.id]"
                    @click.stop="cancelCurrentOrder(order.id)"
                  >
                    取消订单
                  </button>
                  <button
                    class="action-btn btn-primary"
                    :class="{ 'is-loading': loadingMap[order.id] }"
                    :disabled="loadingMap[order.id]"
                    @click.stop="payCurrentOrder(order)"
                  >
                    立即支付
                  </button>
                </template>
                
                <!-- 已发货状态操作 -->
                <template v-if="order.status === 2">
                  <button
                    class="action-btn btn-success"
                    :class="{ 'is-loading': loadingMap[order.id] }"
                    :disabled="loadingMap[order.id]"
                    @click.stop="confirmReceipt(order.id)"
                  >
                    确认收货
                  </button>
                </template>

                <!-- 其他状态：查看详情 (用作保底按钮) -->
                <button
                  v-if="order.status === 1 || order.status === 3 || order.status === 40"
                  class="action-btn btn-default"
                  @click.stop="$router.push(`/order/${order.id}`)"
                >
                  查看详情
                </button>
              </div>
            </div>
            
          </div>

          <!-- 定制分页器 -->
          <div class="pagination-container">
            <el-pagination
              v-model:current-page="page"
              v-model:page-size="size"
              :page-sizes="[10, 20, 50]"
              layout="total, sizes, prev, pager, next, jumper"
              :total="total"
              @size-change="fetchOrders"
              @current-change="fetchOrders"
              class="custom-pagination"
            />
          </div>
        </template>

        <!-- 空状态 -->
        <div class="empty-wrapper" v-else>
          <el-empty description="暂时没有相关订单记录" image-size="160">
            <button class="btn-go-mall" @click="$router.push('/mall')">去挑选好物</button>
          </el-empty>
        </div>
      </div>
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
import { ref, onMounted, reactive, watch } from 'vue'
import dayjs from 'dayjs'
import { ElMessage, ElMessageBox } from 'element-plus'
import Navbar from '@/components/layout/Navbar.vue'
import PayAuthDialog from '@/components/payment/PayAuthDialog.vue'
import { getOrderList, payOrder, cancelOrder, receiveOrder } from '@/api/order'
import { useUserStore } from '@/stores/user'
import { GREEN_POINTS_PER_YUAN, getMixedPaymentPreview } from '@/utils/payment'
import { InfoFilled } from '@element-plus/icons-vue'

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
    1: 'primary',
    2: 'info',
    3: 'success',
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
      { 
        type: 'warning',
        confirmButtonText: '确认支付',
        cancelButtonText: '暂不支付'
      }
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
    await ElMessageBox.confirm('确认取消该订单吗？取消后无法恢复。', '取消确认', { 
      type: 'warning',
      confirmButtonText: '确认取消',
      cancelButtonText: '暂不取消',
      confirmButtonClass: 'el-button--danger'
    })
    await cancelOrder(orderId)
    ElMessage.success('订单已取消')
    await fetchOrders()
  }).catch(() => {})
}

function confirmReceipt(orderId) {
  withLoading(orderId, async () => {
    await ElMessageBox.confirm('请确认您已收到商品，确认收货后交易将完成。', '收货确认', { 
      type: 'success',
      confirmButtonText: '确认收货'
    })
    await receiveOrder(orderId)
    ElMessage.success('确认收货成功')
    await fetchOrders()
  }).catch(() => {})
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
/* 全局背景 */
.order-page {
  min-height: 100vh;
  background-color: #f8f9fa;
  padding-bottom: 80px;
}

/* ★ 大气加宽容器 ★ */
.custom-container {
  max-width: 1200px; /* 从1000px大幅拓宽，更显大气 */
  margin: 0 auto;
}

/* 顶部标题区域 */
.page-header {
  padding: 36px 0 28px;
}

.highlight-title {
  display: inline-block;
  position: relative;
  font-size: 32px; /* 加大标题 */
  color: #2c3e50;
  font-weight: 700;
  margin: 0;
  z-index: 1;
}

.highlight-title::after {
  content: '';
  position: absolute;
  bottom: 4px;
  left: -5%;
  width: 110%;
  height: 14px;
  background-color: #2d597b; 
  opacity: 0.15;
  border-radius: 6px;
  z-index: -1;
  transition: all 0.3s ease;
}

/* 高级提示横幅 */
.premium-alert {
  display: flex;
  align-items: flex-start;
  background: #f0f7ff;
  border: 1px solid #cce3f6;
  border-radius: 8px;
  padding: 16px 20px;
  margin-bottom: 24px;
  color: #2d597b;
}

.alert-icon {
  font-size: 20px;
  margin-right: 12px;
  margin-top: 2px;
  color: #2d597b;
}

.alert-content {
  font-size: 14px;
  line-height: 1.6;
}

/* 深度定制 Tabs 样式 */
.tabs-container {
  margin-bottom: 24px;
  background: #ffffff;
  border-radius: 12px;
  padding: 0 24px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.02);
}

:deep(.custom-tabs .el-tabs__nav-wrap::after) {
  height: 1px;
  background-color: #f0f2f5;
}

:deep(.custom-tabs .el-tabs__active-bar) {
  background-color: #2d597b;
  height: 3px;
  border-radius: 2px;
}

:deep(.custom-tabs .el-tabs__item) {
  font-size: 16px;
  color: #606266;
  height: 60px; /* 加高Tab栏 */
  line-height: 60px;
  transition: all 0.3s;
}

:deep(.custom-tabs .el-tabs__item:hover) {
  color: #2d597b;
}

:deep(.custom-tabs .el-tabs__item.is-active) {
  color: #2d597b;
  font-weight: 600;
  font-size: 17px; /* 激活时字体稍微变大 */
}

/* 订单列表容器 */
.order-list {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

/* ★ 独立重构的纯白悬浮订单卡片 ★ */
.order-card {
  background: #ffffff;
  border-radius: 12px;
  padding: 0 32px; /* 加大两侧留白，呼吸感增强 */
  box-shadow: 0 2px 16px rgba(0, 0, 0, 0.03);
  transition: all 0.3s ease;
  border: 1px solid transparent;
  cursor: pointer;
}

.order-card:hover {
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.06);
  border-color: rgba(45, 89, 123, 0.1);
  transform: translateY(-2px);
}

/* 订单头部 */
.order-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px 0;
  border-bottom: 1px dashed #ebeef5;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 16px;
}

.order-no {
  font-size: 15px; /* 加大字号 */
  font-weight: 600;
  color: #2c3e50;
}

.order-time {
  font-size: 14px;
  color: #a4b0be;
}

:deep(.custom-tag) {
  border: none;
  font-weight: 600;
  padding: 0 16px;
  height: 28px;
  font-size: 13px;
}

/* ★ 订单商品区 (新增画中画背景底色) ★ */
.order-products {
  display: flex;
  flex-direction: column;
  gap: 16px;
  padding: 20px 24px;
  background: #fbfcfd; /* 极浅的专属底色 */
  border-radius: 8px;
  margin: 20px 0;
}

.product-item {
  display: flex;
  align-items: center;
  gap: 20px;
}

.thumb-wrapper {
  width: 80px; /* 大图展示 */
  height: 80px;
  background: #ffffff; /* 图片底色改白，与外层浅灰形成对比 */
  border-radius: 8px;
  overflow: hidden;
  flex-shrink: 0;
  border: 1px solid #f0f2f5;
}

.thumb {
  width: 100%;
  height: 100%;
  object-fit: contain;
  mix-blend-mode: multiply;
}

.product-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.product-name {
  font-size: 16px;
  font-weight: 600;
  color: #303133;
  margin-bottom: 10px;
  line-height: 1.5;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.product-meta {
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-size: 14px;
  color: #909399;
}

.currency {
  font-size: 13px;
  margin-right: 2px;
}

.price {
  color: #303133;
  font-weight: 600;
  font-size: 16px;
}

.qty {
  font-size: 15px;
}

/* 支付信息与底部操作区 */
.order-footer {
  display: flex;
  flex-direction: column;
  padding: 0 0 28px; /* 移除上部padding，因为上方有了 margin */
  gap: 16px;
}

.payment-info {
  display: flex;
  justify-content: space-between;
  align-items: flex-end;
}

.payment-details {
  display: flex;
  gap: 16px;
  font-size: 14px;
}

.point-item {
  color: #00b894; 
  background: rgba(0, 184, 148, 0.1);
  padding: 4px 10px;
  border-radius: 6px;
}

.balance-item {
  color: #606266;
}

.total-amount {
  font-size: 15px;
  color: #606266;
}

.highlight-price {
  font-size: 24px; /* 加大总金额 */
  font-weight: 800;
  color: #e4393c; 
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
}

.order-actions {
  display: flex;
  justify-content: flex-end;
  gap: 16px;
}

/* 定制操作按钮 (放量设计) */
.action-btn {
  padding: 10px 24px; /* 加大点击区域 */
  border-radius: 20px; 
  font-size: 15px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s;
  border: 1px solid transparent;
  background: #ffffff;
}

.btn-primary {
  background: #2d597b;
  color: #ffffff;
  box-shadow: 0 4px 12px rgba(45, 89, 123, 0.2);
}
.btn-primary:hover:not(:disabled) {
  background: #1f435d;
  transform: translateY(-1px);
  box-shadow: 0 6px 16px rgba(45, 89, 123, 0.3);
}

.btn-success {
  background: #00b894;
  color: #ffffff;
  box-shadow: 0 4px 12px rgba(0, 184, 148, 0.2);
}
.btn-success:hover:not(:disabled) {
  background: #00997a;
  transform: translateY(-1px);
  box-shadow: 0 6px 16px rgba(0, 184, 148, 0.3);
}

.btn-cancel {
  border-color: #dcdfe6;
  color: #909399;
}
.btn-cancel:hover:not(:disabled) {
  color: #e4393c;
  border-color: #fbc4c4;
  background: #fef0f0;
}

.btn-default {
  border-color: #dcdfe6;
  color: #606266;
}
.btn-default:hover {
  color: #2d597b;
  border-color: #2d597b;
}

.action-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

/* 定制分页器 */
.pagination-container {
  display: flex;
  justify-content: center;
  margin-top: 32px;
}

:deep(.custom-pagination .el-pager li.is-active) {
  background-color: #2d597b;
  color: #fff;
  border-radius: 4px;
}
:deep(.custom-pagination .el-pager li:hover) {
  color: #2d597b;
}

/* 空状态 */
.empty-wrapper {
  background: #ffffff;
  border-radius: 12px;
  padding: 80px 0;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.02);
}

.btn-go-mall {
  margin-top: 16px;
  background: #ffffff;
  color: #2d597b;
  border: 1px solid #2d597b;
  padding: 10px 32px;
  border-radius: 20px;
  font-size: 15px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s;
}

.btn-go-mall:hover {
  background: #2d597b;
  color: #ffffff;
}

/* 响应式 */
@media (max-width: 768px) {
  .custom-container {
    padding: 0 16px;
  }
  .order-card {
    padding: 0 20px;
  }
  .order-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 12px;
  }
  .header-right {
    align-self: flex-end;
  }
  .payment-info {
    flex-direction: column;
    align-items: flex-end;
    gap: 8px;
  }
  .order-actions {
    width: 100%;
    justify-content: flex-end;
  }
}
</style>