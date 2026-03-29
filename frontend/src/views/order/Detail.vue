<template>
  <div class="order-detail-page">
    <Navbar />

    <div class="container custom-container">
      <!-- 顶部返回导航 -->
      <div class="page-nav">
        <div class="back-btn" @click="$router.back()">
          <el-icon class="back-icon"><ArrowLeft /></el-icon>
          <span>返回列表</span>
        </div>
      </div>

      

      <div v-loading="loading" class="detail-content">
        <template v-if="order">
          <!-- 1. 核心状态横幅 -->
          <div class="status-banner">
            <div class="status-left">
              <div class="status-text" :class="`status-${order.status}`">
                {{ getStatusText(order.status) }}
              </div>
              <div class="status-meta">
                <span class="order-no">订单号：{{ order.order_no }}</span>
                <span class="order-time"
                  >创建时间：{{ formatDate(order.created_at) }}</span
                >
              </div>
            </div>
            <div class="status-right">
              <div class="order-total-highlight">
                <span class="label">应付总额</span>
                <span class="amount"
                  ><span class="currency">￥</span
                  >{{ formatAmount(order.total_amount) }}</span
                >
              </div>
              <div class="status-actions">
                <button
                  v-if="order.status === 0"
                  class="action-btn btn-cancel"
                  @click="cancelCurrentOrder"
                >
                  取消订单
                </button>
                <button
                  v-if="order.status === 0"
                  class="action-btn btn-primary"
                  @click="payCurrentOrder"
                >
                  立即支付
                </button>
                <button
                  v-if="order.status === 2"
                  class="action-btn btn-success"
                  @click="confirmReceipt"
                >
                  确认收货
                </button>
              </div>
            </div>
          </div>

          <!-- 2. 两栏信息网格布局 (支付信息 & 用户信息) -->
          <div class="info-grid">
            <!-- 支付信息卡片 -->
            <div class="detail-card">
              <div class="card-header">
                <el-icon class="header-icon"><Wallet /></el-icon> 支付明细
              </div>
              <div class="card-body">
                <div class="info-row">
                  <span class="row-label">订单总额</span>
                  <span class="row-value price-text"
                    >￥{{ formatAmount(order.total_amount) }}</span
                  >
                </div>
                <div class="info-row">
                  <span class="row-label">积分抵扣</span>
                  <span class="row-value point-text"
                    >- {{ order.used_points || 0 }} 积分</span
                  >
                </div>
                <div class="info-row">
                  <span class="row-label">余额支付</span>
                  <span class="row-value"
                    >￥{{ formatAmount(order.used_balance || 0) }}</span
                  >
                </div>
                <div class="info-row" v-if="order.paid_at">
                  <span class="row-label">支付时间</span>
                  <span class="row-value">{{ formatDate(order.paid_at) }}</span>
                </div>
              </div>
            </div>

            <!-- 用户与服务信息卡片 -->
            <div class="detail-card">
              <div class="card-header">
                <el-icon class="header-icon"><User /></el-icon> 服务信息
              </div>
              <div class="card-body">
                <div class="info-row">
                  <span class="row-label">服务门店</span>
                  <span class="row-value">{{
                    order.store?.name || "默认门店"
                  }}</span>
                </div>
                <div class="info-row">
                  <span class="row-label">联系人</span>
                  <span class="row-value">{{
                    order.sys_user?.real_name || order.sys_user?.username || "-"
                  }}</span>
                </div>
                <div class="info-row">
                  <span class="row-label">联系电话</span>
                  <span class="row-value">{{
                    order.sys_user?.mobile || "-"
                  }}</span>
                </div>
              </div>
            </div>
          </div>

          <!-- 3. 商品清单卡片 -->
          <div class="detail-card product-card">
            <div class="card-header">
              <el-icon class="header-icon"><Goods /></el-icon> 商品清单
            </div>
            <div class="product-list">
              <div
                v-for="item in order.items"
                :key="item.id"
                class="product-item"
                @click="$router.push(`/product/${item.product_id}`)"
              >
                <div class="thumb-wrapper">
                  <img
                    :src="
                      item.product?.image_url ||
                      'https://via.placeholder.com/80'
                    "
                    class="thumb"
                  />
                </div>
                <div class="product-info">
                  <div class="product-name">
                    {{ item.product?.name || "未知商品" }}
                  </div>
                  <div class="product-meta">
                    <span class="price"
                      ><span class="currency">￥</span
                      >{{ formatAmount(item.price) }}</span
                    >
                    <span class="qty">x{{ item.quantity }}</span>
                  </div>
                </div>
                <div class="product-subtotal">
                  <span class="currency">￥</span
                  >{{ formatAmount(item.price * item.quantity) }}
                </div>
              </div>
            </div>
          </div>
        </template>

        <!-- 空状态 -->
        <div class="empty-wrapper" v-else>
          <el-empty description="未找到订单信息，可能已被删除" image-size="160">
            <button class="btn-go-mall" @click="$router.push('/mall')">
              去商城看看
            </button>
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
import { ref, onMounted } from "vue";
import { useRoute, useRouter } from "vue-router";
import dayjs from "dayjs";
import { ElMessage, ElMessageBox } from "element-plus";
import Navbar from "@/components/layout/Navbar.vue";
import PayAuthDialog from "@/components/payment/PayAuthDialog.vue";
import {
  getOrderDetail,
  payOrder,
  cancelOrder,
  receiveOrder,
} from "@/api/order";
import { useUserStore } from "@/stores/user";
import { GREEN_POINTS_PER_YUAN, getMixedPaymentPreview } from "@/utils/payment";
// 引入图标
import { ArrowLeft, Wallet, User, Goods } from "@element-plus/icons-vue";

const route = useRoute();
const router = useRouter();
const userStore = useUserStore();
const order = ref(null);
const loading = ref(false);
const showPayAuth = ref(false);
const paySubmitting = ref(false);

function formatAmount(value) {
  return Number(value || 0).toFixed(2);
}

function formatDate(date) {
  return dayjs(date).format("YYYY-MM-DD HH:mm:ss");
}

function getStatusText(status) {
  return (
    {
      0: "待支付",
      1: "已支付",
      2: "已发货",
      3: "已完成",
      40: "已取消",
    }[status] || "未知"
  );
}

function getPaymentPreview(amount) {
  return getMixedPaymentPreview(amount, userStore.userInfo.green_points);
}

async function fetchDetail() {
  loading.value = true;
  try {
    order.value = await getOrderDetail(route.params.id);
  } catch (error) {
    ElMessage.error(
      error.response?.data?.msg || error.message || "获取订单详情失败",
    );
  } finally {
    loading.value = false;
  }
}

async function payCurrentOrder() {
  const preview = getPaymentPreview(order.value.total_amount);
  try {
    await ElMessageBox.confirm(
      `本次将按 ${GREEN_POINTS_PER_YUAN} 积分=1元优先抵扣 ${preview.points} 积分，余额支付 ￥${formatAmount(preview.balance)}，确认支付吗？`,
      "支付确认",
      {
        type: "warning",
        confirmButtonText: "确认支付",
        cancelButtonText: "暂不支付",
      },
    );
    showPayAuth.value = true;
  } catch (error) {
    if (error === "cancel" || error === "close") {
      return;
    }
    ElMessage.error(error.response?.data?.msg || error.message || "支付失败");
  }
}

async function submitOrderPay(authPayload) {
  if (!order.value) return;

  paySubmitting.value = true;
  try {
    const res = await payOrder({
      order_id: order.value.id,
      business_type: 1,
      ...authPayload,
    });
    const paymentResult = res?.payment_result || res;
    ElMessage.success(
      `支付成功，使用积分 ${paymentResult.used_points}，余额 ￥${formatAmount(paymentResult.used_balance)}`,
    );
    showPayAuth.value = false;
    await Promise.all([fetchDetail(), userStore.fetchUserInfo()]);
  } catch (error) {
    ElMessage.error(error.response?.data?.msg || error.message || "支付失败");
  } finally {
    paySubmitting.value = false;
  }
}

async function cancelCurrentOrder() {
  try {
    await ElMessageBox.confirm("确认取消该订单吗？取消后无法恢复。", "提示", {
      type: "warning",
      confirmButtonText: "确认取消",
      cancelButtonText: "暂不取消",
      confirmButtonClass: "el-button--danger",
    });
    await cancelOrder(order.value.id);
    ElMessage.success("订单已取消");
    await fetchDetail();
  } catch (error) {
    if (error !== "cancel" && error !== "close") {
      ElMessage.error(error.response?.data?.msg || error.message || "操作失败");
    }
  }
}

async function confirmReceipt() {
  try {
    await ElMessageBox.confirm(
      "请确认您已收到商品，确认收货后交易将完成。",
      "收货确认",
      {
        type: "success",
        confirmButtonText: "确认收货",
      },
    );
    await receiveOrder(order.value.id);
    ElMessage.success("确认收货成功");
    await fetchDetail();
  } catch (error) {
    if (error !== "cancel" && error !== "close") {
      ElMessage.error(error.response?.data?.msg || error.message || "操作失败");
    }
  }
}

onMounted(async () => {
  await Promise.all([fetchDetail(), userStore.fetchUserInfo()]);
});
</script>

<style scoped>
/* 全局背景 */
.order-detail-page {
  min-height: 100vh;
  background-color: #f8f9fa;
  padding-bottom: 80px;
}

/* ★ 大气加宽容器 ★ */
.custom-container {
  max-width: 1000px; /* 详情页采用1000px，既宽广又聚拢视线 */
  margin: 0 auto;
}

/* 顶部返回导航 */
.page-nav {
  display: flex;
  justify-content: flex-start;
  width: 100%;
  padding: 16px 0 8px;
}

.back-btn {
  display: inline-flex;
  align-items: center;
  color: #606266;
  font-size: 15px;
  cursor: pointer;
  transition: color 0.3s;
  padding: 8px 0;
  margin: 0;
}

.back-btn:hover {
  color: #2d597b;
}

.back-icon {
  margin-right: 6px;
  font-size: 16px;
}

/* 页面标题 */
.page-header {
  padding: 8px 0 16px;
}

.highlight-title {
  display: inline-block;
  position: relative;
  font-size: 32px;
  color: #2c3e50;
  font-weight: 700;
  margin: 0;
  z-index: 1;
}

.highlight-title::after {
  content: "";
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

.detail-content {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

/* ================= 1. 状态横幅 ================= */
.status-banner {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: #ffffff;
  border-radius: 12px;
  padding: 32px 40px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.03);
  border-top: 4px solid #2d597b; /* 顶部加粗品牌色边框 */
}

.status-left {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.status-text {
  font-size: 28px;
  font-weight: 800;
}

/* 状态颜色动态绑定 */
.status-0 {
  color: #e6a23c;
} /* 待支付 - 橙色 */
.status-1 {
  color: #2d597b;
} /* 已支付 - 商务蓝 */
.status-2 {
  color: #409eff;
} /* 已发货 - 浅蓝 */
.status-3 {
  color: #00b894;
} /* 已完成 - 绿色 */
.status-40 {
  color: #f56c6c;
} /* 已取消 - 红色 */

.status-meta {
  display: flex;
  flex-direction: column;
  gap: 6px;
  color: #909399;
  font-size: 14px;
}

.status-right {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 20px;
}

.order-total-highlight {
  text-align: right;
}

.order-total-highlight .label {
  font-size: 14px;
  color: #606266;
  margin-right: 8px;
}

.order-total-highlight .amount {
  font-size: 32px;
  font-weight: 800;
  color: #e4393c;
  font-family:
    -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
  letter-spacing: -1px;
}

.currency {
  font-size: 16px;
  margin-right: 2px;
}

.status-actions {
  display: flex;
  gap: 12px;
}

/* 定制操作按钮 */
.action-btn {
  padding: 10px 28px;
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
.btn-primary:hover {
  background: #1f435d;
  transform: translateY(-2px);
  box-shadow: 0 6px 16px rgba(45, 89, 123, 0.3);
}

.btn-success {
  background: #00b894;
  color: #ffffff;
  box-shadow: 0 4px 12px rgba(0, 184, 148, 0.2);
}
.btn-success:hover {
  background: #00997a;
  transform: translateY(-2px);
  box-shadow: 0 6px 16px rgba(0, 184, 148, 0.3);
}

.btn-cancel {
  border-color: #dcdfe6;
  color: #909399;
}
.btn-cancel:hover {
  color: #e4393c;
  border-color: #fbc4c4;
  background: #fef0f0;
}

/* ================= 2. 两栏信息网格 ================= */
.info-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 24px;
}

.detail-card {
  background: #ffffff;
  border-radius: 12px;
  padding: 28px 32px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.02);
}

.card-header {
  display: flex;
  align-items: center;
  font-size: 18px;
  font-weight: 600;
  color: #303133;
  margin-bottom: 24px;
  padding-bottom: 16px;
  border-bottom: 1px dashed #ebeef5;
}

.header-icon {
  color: #2d597b;
  margin-right: 8px;
  font-size: 20px;
}

.info-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 0;
}

.row-label {
  color: #909399;
  font-size: 14px;
}

.row-value {
  color: #303133;
  font-weight: 500;
  font-size: 15px;
}

.price-text {
  font-weight: 700;
  color: #303133;
}

.point-text {
  color: #00b894;
  font-weight: 600;
}

/* ================= 3. 商品清单 ================= */
.product-card {
  margin-bottom: 20px;
}

.product-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.product-item {
  display: flex;
  align-items: center;
  gap: 24px;
  padding: 16px;
  background: #fbfcfd; /* 画中画底色 */
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s ease;
  border: 1px solid transparent;
}

.product-item:hover {
  background: #ffffff;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.04);
  border-color: rgba(45, 89, 123, 0.1);
  transform: translateX(4px);
}

.thumb-wrapper {
  width: 80px;
  height: 80px;
  background: #ffffff;
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
  gap: 16px;
  font-size: 14px;
  color: #909399;
}

.product-meta .price {
  color: #303133;
  font-weight: 600;
}

.product-subtotal {
  font-size: 18px;
  font-weight: 700;
  color: #e4393c;
  min-width: 100px;
  text-align: right;
  font-family:
    -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
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
  .status-banner {
    flex-direction: column;
    align-items: flex-start;
    gap: 24px;
    padding: 24px;
  }
  .status-right {
    width: 100%;
    align-items: flex-start;
  }
  .order-total-highlight {
    text-align: left;
  }
  .info-grid {
    grid-template-columns: 1fr;
  }
  .product-item {
    flex-wrap: wrap;
    position: relative;
  }
  .product-subtotal {
    position: absolute;
    right: 16px;
    bottom: 16px;
  }
}
</style>
