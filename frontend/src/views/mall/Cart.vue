<template>
  <div class="cart-page">
    <Navbar />

    <div class="container custom-container">
      <div class="page-header">
        <h1 class="page-title highlight-title">我的购物车</h1>
      </div>

      <div v-if="cartItems.length > 0" class="cart-content">
        <!-- 购物车商品列表 -->
        <div class="cart-list">
          <div class="cart-item" v-for="item in cartItems" :key="item.id">
            <!-- 1. 商品图 -->
            <div
              class="item-image-wrapper"
              @click="$router.push(`/product/${item.product.id}`)"
            >
              <img
                :src="
                  item.product.image_url || 'https://via.placeholder.com/100'
                "
                class="item-image"
              />
            </div>

            <!-- 2. 商品信息 -->
            <div class="item-info">
              <div
                class="item-name"
                @click="$router.push(`/product/${item.product.id}`)"
              >
                {{ item.product.name }}
              </div>
              <div class="item-price">
                <span class="price-currency">¥</span
                ><span class="price-amount">{{ item.product.price }}</span>
              </div>
            </div>

            <!-- 3. 数量控制器 (定制高级样式) -->
            <div class="quantity-control">
              <button
                class="qty-btn"
                :disabled="item.quantity <= 1"
                @click="changeQuantity(item, -1)"
              >
                <el-icon><Minus /></el-icon>
              </button>
              <input
                type="number"
                class="qty-input custom-hide-arrows"
                v-model.number="item.quantity"
                @change="manualUpdate(item)"
                min="1"
              />
              <button class="qty-btn" @click="changeQuantity(item, 1)">
                <el-icon><Plus /></el-icon>
              </button>
            </div>

            <!-- 4. 小计 (高级感细节：展示单行总价) -->
            <div class="item-subtotal">
              <span class="price-currency">¥</span
              ><span class="price-amount">{{
                (item.product.price * item.quantity).toFixed(2)
              }}</span>
            </div>

            <!-- 5. 删除操作 -->
            <div class="item-actions">
              <button
                class="btn-delete"
                @click="removeItem(item.id)"
                title="移除商品"
              >
                <el-icon><Delete /></el-icon>
              </button>
            </div>
          </div>
        </div>

        <!-- ★ 底部结算栏 (贴底通栏悬浮) ★ -->
        <div class="cart-footer">
          <div class="cart-footer-inner">
            <div class="footer-left">
              <span class="total-count"
                >共
                <span class="highlight-count">{{ cartItems.length }}</span>
                件商品</span
              >
            </div>
            <div class="footer-right">
              <div class="total">
                <span class="total-label">合计金额：</span>
                <span class="total-price">
                  <span class="price-currency">¥</span>{{ totalPrice }}
                </span>
              </div>
              <button class="btn-checkout" @click="checkout">去结算</button>
            </div>
          </div>
        </div>
      </div>

      <!-- 优雅的空状态 -->
      <div class="empty-state-wrapper" v-else>
        <el-empty
          description="购物车空空如也，快去挑选心仪的商品吧"
          image-size="160"
        >
          <button class="btn-go-mall" @click="$router.push('/mall')">
            去逛逛
          </button>
        </el-empty>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from "vue";
import { useRouter } from "vue-router";
import Navbar from "@/components/layout/Navbar.vue";
import { getCartList, deleteCartItem, updateCartQuantity } from "@/api/order";
import { useCartStore } from "@/stores/cart";
import { ElMessage } from "element-plus";
// 引入精致的图标
import { Minus, Plus, Delete } from "@element-plus/icons-vue";

const router = useRouter();
const cartStore = useCartStore();
const cartItems = ref([]);

const totalPrice = computed(() => {
  return cartItems.value
    .reduce((sum, item) => {
      return sum + item.product.price * item.quantity;
    }, 0)
    .toFixed(2);
});

const fetchCart = async () => {
  try {
    cartItems.value = await getCartList();
    cartStore.fetchCart(); // 同步 store 数量
  } catch (error) {
    console.error("获取购物车失败:", error);
  }
};

const changeQuantity = async (item, delta) => {
  const newQty = item.quantity + delta;
  if (newQty < 1) return;

  // 乐观更新
  item.quantity = newQty;
  try {
    await updateCartQuantity(item.id, newQty);
  } catch (e) {
    item.quantity -= delta; // 回滚
    console.error(e);
  }
};

const manualUpdate = async (item) => {
  if (item.quantity < 1) item.quantity = 1;
  try {
    await updateCartQuantity(item.id, item.quantity);
  } catch (e) {
    console.error(e);
  }
};

const removeItem = async (id) => {
  try {
    await deleteCartItem(id);
    fetchCart();
    ElMessage.success("已移出购物车");
  } catch (error) {
    ElMessage.error("删除失败");
  }
};

const checkout = () => {
  router.push("/order/create");
};

onMounted(() => {
  fetchCart();
});
</script>

<style scoped>
/* 全局背景 */
.cart-page {
  min-height: 100vh;
  background-color: #f8f9fa;
  /* 增加底部内边距，防止商品列表底部被全宽悬浮的结算栏挡住 */
  padding-bottom: 120px;
}

.custom-container {
  max-width: 1100px;
  margin: 0 auto;
}

/* 顶部标题区域 */
.page-header {
  padding: 32px 0 24px;
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
  content: "";
  position: absolute;
  bottom: 4px;
  left: -5%;
  width: 110%;
  height: 12px;
  background-color: #2d597b; /* 商务蓝 */
  opacity: 0.2;
  border-radius: 4px;
  z-index: -1;
  transition: all 0.3s ease;
}

/* 购物车内容区 */
.cart-content {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

..cart-list {
  display: flex;
  flex-direction: column;
  gap: 0;
}

/* 单个商品卡片 */
.cart-item {
  display: flex;
  align-items: center;
  padding: 24px;
  background: #ffffff;
  border-radius: 0;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.03);
  transition: all 0.3s cubic-bezier(0.25, 0.8, 0.25, 1);
  border: 1px solid transparent;
}

.cart-item:hover {
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.06);
  border-color: rgba(45, 89, 123, 0.1);
}

.item-image-wrapper {
  width: 100px;
  height: 100px;
  background: #f8f9fa;
  border-radius: 8px;
  overflow: hidden;
  margin-right: 24px;
  cursor: pointer;
  flex-shrink: 0;
}

.item-image {
  width: 100%;
  height: 100%;
  object-fit: contain;
  mix-blend-mode: multiply;
  transition: transform 0.3s;
}

.item-image-wrapper:hover .item-image {
  transform: scale(1.05);
}

.item-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: center;
  min-width: 200px;
}

.item-name {
  font-size: 16px;
  font-weight: 600;
  color: #303133;
  margin-bottom: 12px;
  cursor: pointer;
  transition: color 0.2s;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.item-name:hover {
  color: #2d597b;
}

.item-price {
  color: #909399;
}

.price-currency {
  font-size: 14px;
  margin-right: 2px;
}

.price-amount {
  font-size: 18px;
  font-weight: bold;
  font-family:
    -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
}

/* 数量控制器定制 */
.quantity-control {
  display: flex;
  align-items: center;
  border: 1px solid #dcdfe6;
  border-radius: 6px;
  overflow: hidden;
  margin: 0 40px;
  background: #ffffff;
}

.qty-btn {
  width: 36px;
  height: 32px;
  background: #f5f7fa;
  border: none;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #606266;
  cursor: pointer;
  transition: all 0.2s;
}

.qty-btn:hover:not(:disabled) {
  background: #e4e7ed;
  color: #2d597b;
}

.qty-btn:disabled {
  color: #c0c4cc;
  cursor: not-allowed;
  background: #fafbfc;
}

.qty-input {
  width: 48px;
  height: 32px;
  text-align: center;
  border: none;
  border-left: 1px solid #dcdfe6;
  border-right: 1px solid #dcdfe6;
  font-size: 14px;
  color: #303133;
  font-weight: 600;
  outline: none;
}

.custom-hide-arrows::-webkit-outer-spin-button,
.custom-hide-arrows::-webkit-inner-spin-button {
  -webkit-appearance: none;
  margin: 0;
}
.custom-hide-arrows[type="number"] {
  appearance: textfield;
}

/* 小计金额 */
.item-subtotal {
  width: 120px;
  text-align: right;
  color: #e4393c;
  margin-right: 32px;
}

/* 删除按钮 */
.item-actions {
  width: 40px;
  text-align: right;
}

.btn-delete {
  background: transparent;
  border: none;
  color: #a4b0be;
  font-size: 20px;
  cursor: pointer;
  transition: all 0.2s;
  padding: 8px;
  border-radius: 50%;
}

.btn-delete:hover {
  color: #e74c3c;
  background: #fdf6f6;
}

/* ★ 底部结算栏 (全宽贴底悬浮) ★ */
.cart-footer {
  position: fixed;
  bottom: 0; /* 贴紧底部 */
  left: 0;
  width: 100%; /* 占满全宽 */
  z-index: 999;

  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  box-shadow: 0 -4px 20px rgba(0, 0, 0, 0.05); /* 向上投射的阴影 */
  border-top: 1px solid rgba(0, 0, 0, 0.05); /* 顶部细微边框 */
}

/* 内容限制宽度并居中，保持和上方商品列表对齐 */
.cart-footer-inner {
  max-width: 1100px;
  margin: 0 auto;
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 24px;
}

.footer-left {
  color: #606266;
  font-size: 15px;
}

.highlight-count {
  color: #2d597b;
  font-weight: bold;
  font-size: 16px;
  margin: 0 4px;
}

.footer-right {
  display: flex;
  align-items: center;
  gap: 32px;
}

.total {
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
  font-family:
    -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
  letter-spacing: -1px;
}

.btn-checkout {
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

.btn-checkout:hover {
  background: #1f435d;
  transform: translateY(-2px);
  box-shadow: 0 6px 16px rgba(45, 89, 123, 0.4);
}

/* 空状态 */
.empty-state-wrapper {
  background: #ffffff;
  border-radius: 12px;
  padding: 60px 0;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.03);
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

/* 响应式适配 */
@media (max-width: 768px) {
  .cart-item {
    flex-wrap: wrap;
    position: relative;
    padding: 16px;
  }
  .item-image-wrapper {
    width: 80px;
    height: 80px;
    margin-right: 16px;
  }
  .item-info {
    min-width: 0;
  }
  .quantity-control {
    margin: 16px 0 0 0;
  }
  .item-subtotal {
    display: none;
  }
  .item-actions {
    position: absolute;
    top: 16px;
    right: 16px;
  }
  .cart-footer-inner {
    flex-direction: column;
    gap: 16px;
    padding: 16px;
  }
  .footer-right {
    width: 100%;
    justify-content: space-between;
    gap: 16px;
  }
  .btn-checkout {
    padding: 12px 24px;
  }
}
</style>
