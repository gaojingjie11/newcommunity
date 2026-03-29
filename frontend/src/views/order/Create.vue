<template>
  <div class="create-order-page">
    <Navbar />

    <div class="container custom-container">
      <!-- 顶部返回导航 -->
      <div class="page-nav">
        <div class="back-btn" @click="$router.push('/cart')">
          <el-icon class="back-icon"><ArrowLeft /></el-icon> 
          <span>返回购物车</span>
        </div>
      </div>

      <div class="page-header">
        <h1 class="page-title highlight-title">确认订单</h1>
      </div>

      <!-- 积分提示横幅 (替代原有 alert，更高级) -->
      <div class="premium-alert">
        <el-icon class="alert-icon"><InfoFilled /></el-icon>
        <div class="alert-content">
          当前积分 <span class="highlight-text">{{ userStore.userInfo.green_points || 0 }}</span>，
          余额 <span class="highlight-text">¥{{ formatAmount(userStore.userInfo.balance || 0) }}</span>。
          订单支付时将按 <strong>{{ GREEN_POINTS_PER_YUAN }} 积分 = 1 元</strong> 自动优先抵扣。
        </div>
      </div>

      <div class="order-content">
        <!-- 1. 门店选择卡片 -->
        <div class="order-card">
          <div class="section-header">
            <el-icon class="header-icon"><Location /></el-icon> 选择服务门店
          </div>
          <div class="store-selector">
            <el-select 
              v-model="selectedStoreId" 
              placeholder="请选择提货/服务门店" 
              class="custom-select"
              size="large"
            >
              <el-option 
                v-for="store in storeList" 
                :key="store.id" 
                :label="`${store.name} (${store.address})`" 
                :value="store.id" 
              />
            </el-select>
          </div>
        </div>

        <!-- 2. 商品清单卡片 -->
        <div class="order-card">
          <div class="section-header">
            <el-icon class="header-icon"><Goods /></el-icon> 商品清单
          </div>
          <div class="order-items">
            <div class="item" v-for="item in cartItems" :key="item.id">
              <div class="item-img-wrapper">
                <img :src="item.product.image_url" class="item-img" />
              </div>
              <div class="item-info">
                <div class="item-name">{{ item.product.name }}</div>
                <div class="item-price-calc">
                  <span class="unit-price">¥{{ formatAmount(item.product.price) }}</span>
                  <span class="item-qty">x {{ item.quantity }}</span>
                </div>
              </div>
              <div class="item-total">
                <span class="currency">¥</span>{{ formatAmount(item.product.price * item.quantity) }}
              </div>
            </div>
          </div>
        </div>

        <!-- 3. 支付预估卡片 -->
        <div class="order-card">
          <div class="section-header">
            <el-icon class="header-icon"><Wallet /></el-icon> 支付预估
          </div>
          <div class="preview-box">
            <div class="preview-line">
              <span class="line-label">商品总额</span>
              <span class="line-value price-text">¥{{ totalPrice }}</span>
            </div>
            <div class="preview-line deduct-line">
              <span class="line-label">预计抵扣积分</span>
              <span class="line-value">- {{ paymentPreview.points }} 积分</span>
            </div>
            <div class="preview-line total-line">
              <span class="line-label">预计扣除余额</span>
              <span class="line-value final-price">
                <span class="currency">¥</span>{{ formatAmount(paymentPreview.balance) }}
              </span>
            </div>
          </div>
        </div>
      </div>

      <!-- ★ 底部悬浮提交栏 ★ -->
      <div class="footer-bar">
        <div class="footer-inner">
          <div class="total-info">
            <span class="count-text">共 <span class="highlight-count">{{ totalCount }}</span> 件商品</span>
            <div class="total-price-wrap">
              <span class="total-label">实付余额：</span>
              <span class="total-price">
                <span class="currency">¥</span>{{ formatAmount(paymentPreview.balance) }}
              </span>
            </div>
          </div>
          <button 
            class="btn-submit" 
            :class="{ 'is-loading': submitting }" 
            :disabled="submitting" 
            @click="submitOrder"
          >
            {{ submitting ? '提交中...' : '提交订单' }}
          </button>
        </div>
      </div>

    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import Navbar from '@/components/layout/Navbar.vue'
import { useUserStore } from '@/stores/user'
import { getCartList, createOrder } from '@/api/order'
import { getStoreList } from '@/api/service'
import { GREEN_POINTS_PER_YUAN, getMixedPaymentPreview } from '@/utils/payment'
// 引入图标
import { ArrowLeft, InfoFilled, Location, Goods, Wallet } from '@element-plus/icons-vue'

const router = useRouter()
const userStore = useUserStore()
const cartItems = ref([])
const storeList = ref([])
const selectedStoreId = ref()
const submitting = ref(false)

const totalCount = computed(() => cartItems.value.reduce((sum, item) => sum + item.quantity, 0))
const totalPriceNumber = computed(() => cartItems.value.reduce((sum, item) => sum + item.product.price * item.quantity, 0))
const totalPrice = computed(() => formatAmount(totalPriceNumber.value))
const paymentPreview = computed(() => getMixedPaymentPreview(totalPriceNumber.value, userStore.userInfo.green_points))

function formatAmount(value) {
  return Number(value || 0).toFixed(2)
}

async function loadData() {
  const [cartRes, storeRes] = await Promise.all([getCartList(), getStoreList(), userStore.fetchUserInfo()])
  cartItems.value = cartRes || []
  storeList.value = storeRes || []
  if (storeList.value.length > 0) {
    selectedStoreId.value = storeList.value[0].id
  }
  if (cartItems.value.length === 0) {
    ElMessage.warning('购物车为空')
    router.push('/mall')
  }
}

async function submitOrder() {
  if (!selectedStoreId.value) {
    ElMessage.warning('请选择提货门店')
    return
  }

  submitting.value = true
  try {
    const items = cartItems.value.map((item) => ({
      cart_id: item.id,
      quantity: item.quantity
    }))

    await createOrder({
      store_id: selectedStoreId.value,
      items
    })
    ElMessage.success('订单创建成功')
    router.push('/order')
  } catch (error) {
    ElMessage.error(error.response?.data?.msg || error.message || '创建订单失败')
  } finally {
    submitting.value = false
  }
}

onMounted(loadData)
</script>

<style scoped>
/* 全局背景 */
.create-order-page {
  min-height: 100vh;
  background-color: #f8f9fa;
  padding-bottom: 140px; 
}

.custom-container {
  max-width: 900px; /* 订单确认页不需要太宽，窄一点更聚拢、阅读体验更好 */
  margin: 0 auto;
}

/* 顶部返回导航 */
.page-nav {
  padding: 24px 0 16px;
}

.back-btn {
  display: inline-flex;
  align-items: center;
  color: #606266;
  font-size: 15px;
  cursor: pointer;
  transition: color 0.3s;
  padding: 8px 16px 8px 0;
}

.back-btn:hover {
  color: #2d597b;
}

.back-icon {
  margin-right: 6px;
  font-size: 16px;
}

/* 标题 */
.page-header {
  margin-bottom: 24px;
}

.highlight-title {
  display: inline-block;
  position: relative;
  font-size: 28px;
  color: #2c3e50;
  font-weight: 600;
  margin: 0;
  z-index: 1;
}

.highlight-title::after {
  content: '';
  position: absolute;
  bottom: 4px;
  left: -5%;
  width: 110%;
  height: 12px;
  background-color: #2d597b; 
  opacity: 0.2;
  border-radius: 4px;
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

.highlight-text {
  font-weight: bold;
  font-size: 15px;
  color: #e4393c;
  margin: 0 2px;
}

/* 卡片统一规范 */
.order-content {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.order-card {
  background: #ffffff;
  border-radius: 12px;
  padding: 24px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.03);
}

.section-header {
  display: flex;
  align-items: center;
  font-size: 18px;
  font-weight: 600;
  color: #303133;
  margin-bottom: 20px;
  border-bottom: 1px solid #f0f2f5;
  padding-bottom: 16px;
}

.header-icon {
  color: #2d597b;
  margin-right: 8px;
  font-size: 20px;
}

/* 门店选择器深度调整 */
.store-selector {
  width: 100%;
}
:deep(.custom-select) {
  width: 100%;
}
:deep(.custom-select .el-input__wrapper) {
  box-shadow: 0 0 0 1px #dcdfe6 inset;
  border-radius: 8px;
  padding: 4px 12px;
}
:deep(.custom-select .el-input__wrapper.is-focus) {
  box-shadow: 0 0 0 1px #2d597b inset !important;
}

/* 商品清单 */
.order-items {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.item {
  display: flex;
  align-items: center;
}

.item-img-wrapper {
  width: 72px;
  height: 72px;
  background: #f8f9fa;
  border-radius: 8px;
  overflow: hidden;
  margin-right: 16px;
  flex-shrink: 0;
}

.item-img {
  width: 100%;
  height: 100%;
  object-fit: contain;
  mix-blend-mode: multiply;
}

.item-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.item-name {
  font-size: 15px;
  font-weight: 600;
  color: #303133;
  margin-bottom: 8px;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.item-price-calc {
  font-size: 13px;
  color: #909399;
}

.unit-price {
  margin-right: 8px;
}

.item-total {
  font-size: 18px;
  font-weight: bold;
  color: #303133;
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
  margin-left: 20px;
}

.currency {
  font-size: 13px;
  margin-right: 2px;
  font-weight: normal;
}

/* 支付明细账单 */
.preview-box {
  background: #fcfdfe;
  border: 1px solid #ebeef5;
  border-radius: 8px;
  padding: 20px;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.preview-line {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.line-label {
  font-size: 14px;
  color: #606266;
}

.line-value {
  font-size: 15px;
  font-weight: 600;
  color: #303133;
}

.deduct-line .line-value {
  color: #00b894; /* 抵扣用清新的绿色提示，让人感觉赚到了 */
}

.total-line {
  margin-top: 8px;
  padding-top: 20px;
  border-top: 1px dashed #dcdfe6;
}

.total-line .line-label {
  font-size: 16px;
  font-weight: bold;
  color: #303133;
}

.final-price {
  font-size: 24px;
  font-weight: 800;
  color: #e4393c;
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
}

/* ★ 底部结算悬浮栏 (复刻购物车样式) ★ */
.footer-bar {
  position: fixed; 
  bottom: 0;     
  left: 0;       
  width: 100%;   
  z-index: 999;  
  
  background: rgba(255, 255, 255, 0.95); 
  backdrop-filter: blur(10px);           
  box-shadow: 0 -4px 20px rgba(0, 0, 0, 0.05); 
  border-top: 1px solid rgba(0, 0, 0, 0.05);   
}

.footer-inner {
  max-width: 900px; /* 匹配上面容器的宽度 */
  margin: 0 auto;
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 24px;
}

.total-info {
  display: flex;
  align-items: baseline;
  gap: 24px;
}

.count-text {
  color: #606266;
  font-size: 15px;
}

.highlight-count {
  color: #2d597b;
  font-weight: bold;
  font-size: 16px;
  margin: 0 2px;
}

.total-price-wrap {
  display: flex;
  align-items: baseline;
}

.total-label {
  font-size: 15px;
  color: #303133;
}

.total-price {
  color: #e4393c;
  font-size: 32px;
  font-weight: 800;
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
  letter-spacing: -1px;
}

.btn-submit {
  background: #2d597b;
  color: #ffffff;
  border: none;
  border-radius: 24px;
  padding: 12px 36px;
  font-size: 16px;
  font-weight: bold;
  cursor: pointer;
  transition: all 0.3s;
  box-shadow: 0 4px 12px rgba(45, 89, 123, 0.3);
}

.btn-submit:hover:not(:disabled) {
  background: #1f435d;
  transform: translateY(-2px);
  box-shadow: 0 6px 16px rgba(45, 89, 123, 0.4);
}

.btn-submit:disabled {
  background: #a4b0be;
  box-shadow: none;
  cursor: not-allowed;
  transform: none;
}

/* 响应式 */
@media (max-width: 768px) {
  .footer-inner {
    flex-direction: column;
    gap: 16px;
    padding: 16px;
  }
  .total-info {
    width: 100%;
    justify-content: space-between;
  }
  .btn-submit {
    width: 100%;
  }
}
</style>