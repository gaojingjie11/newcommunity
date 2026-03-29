<template>
  <div class="mall-page">
    <Navbar />

    <div class="container">
      <div class="page-header">
        <h1 class="page-title highlight-title">社区商城</h1>
        <p class="page-subtitle">甄选社区好物，即刻送货上门</p>
      </div>

      <div class="search-panel">
        <div class="search-console">
          <div class="search-inputs">
            <el-input
              v-model="searchKeyword"
              placeholder="输入商品名称探索好物..."
              class="custom-input"
              @keyup.enter="searchProducts"
              clearable
            >
              <template #prefix
                ><el-icon><Search /></el-icon
              ></template>
            </el-input>

            <div class="divider"></div>

            <el-select
              v-model="selectedCategory"
              placeholder="全部分类"
              clearable
              class="custom-select"
              @change="handleSearch"
            >
              <el-option :value="0" label="所有分类" />
              <el-option
                v-for="c in categories"
                :key="c.id"
                :value="c.id"
                :label="c.name"
              />
            </el-select>
          </div>
        </div>

        <el-button type="primary" class="search-btn" @click="handleSearch">
          <el-icon class="btn-icon"><Search /></el-icon>
        </el-button>
      </div>

      <div class="product-grid" v-loading="loading">
        <el-row :gutter="24">
          <el-col
            :xs="12"
            :sm="8"
            :md="6"
            :lg="6"
            :xl="4"
            v-for="product in products"
            :key="product.id"
          >
            <div class="product-card" @click="goToDetail(product.id)">
              <div class="product-image">
                <img
                  :src="product.image_url || 'https://via.placeholder.com/200'"
                  :alt="product.name"
                />
                <div class="product-badge" v-if="product.is_promotion">
                  限时特惠
                </div>
              </div>
              <div class="product-info">
                <div class="product-name" :title="product.name">
                  {{ product.name }}
                </div>
                <div class="product-desc" :title="product.description">
                  {{ product.description }}
                </div>
                <div class="product-footer">
                  <div class="product-price">
                    <span class="price-currency">¥</span>
                    <span class="price-current">{{ product.price }}</span>
                    <span class="price-original" v-if="product.original_price"
                      >¥{{ product.original_price }}</span
                    >
                  </div>
                  <div
                    class="cart-btn-wrapper"
                    @click.stop="addToCart(product)"
                  >
                    <el-button
                      type="primary"
                      size="small"
                      :icon="ShoppingCart"
                      circle
                      class="cart-btn"
                    />
                  </div>
                </div>
              </div>
            </div>
          </el-col>
        </el-row>
      </div>

      <div class="empty-state" v-if="!loading && products.length === 0">
        <el-empty
          description="未找到相关商品，换个关键词试试吧"
          image-size="200"
        />
      </div>

      <div class="pagination-container" v-if="total > 0">
        <el-pagination
          v-model:current-page="page"
          v-model:page-size="size"
          :page-sizes="[12, 24, 48, 96]"
          layout="total, sizes, prev, pager, next, jumper"
          :total="total"
          @size-change="handleSearch"
          @current-change="handleSearch"
          class="custom-pagination"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from "vue";
import { useRouter } from "vue-router";
import Navbar from "@/components/layout/Navbar.vue";
import { getProductList, getCategories } from "@/api/product";
import { addToCart as addToCartApi } from "@/api/order";
import { useCartStore } from "@/stores/cart";
import { ElMessage } from "element-plus";
import { Search, ShoppingCart } from "@element-plus/icons-vue";

const router = useRouter();
const cartStore = useCartStore();

const products = ref([]);
const categories = ref([]);
const searchKeyword = ref("");
const selectedCategory = ref(0);
const page = ref(1);
const size = ref(12);
const total = ref(0);
const loading = ref(false);

const handleSearch = () => {
  fetchData();
};

const searchProducts = () => {
  page.value = 1;
  fetchData();
};

const fetchData = async () => {
  loading.value = true;
  try {
    const data = await getProductList({
      name: searchKeyword.value,
      category_id: selectedCategory.value || undefined,
      status: 1,
      page: page.value,
      size: size.value,
    });
    products.value = data.list || [];
    total.value = data.total;
  } catch (error) {
    console.error("搜索失败:", error);
    ElMessage.error("加载商品失败");
  } finally {
    loading.value = false;
  }
};

const fetchCategories = async () => {
  try {
    const res = await getCategories();
    categories.value = res || [];
  } catch (e) {
    console.error(e);
  }
};

const goToDetail = (id) => {
  router.push(`/product/${id}`);
};

const addToCart = async (product) => {
  try {
    await addToCartApi({
      product_id: product.id,
      quantity: 1,
    });
    ElMessage.success("已成功加入购物车！");
    cartStore.fetchCart();
  } catch (error) {
    if (error.response?.status === 401) {
      ElMessage.warning("请先登录后再进行操作");
      router.push("/login");
    } else {
      ElMessage.error("添加失败，请稍后重试");
    }
  }
};

onMounted(async () => {
  await fetchCategories();
  await fetchData();
});
</script>

<style scoped>
/* 确保背景色能够衬托出白色卡片的高级感 */
.mall-page {
  min-height: 100vh;
  background-color: #f5f7fa;
  padding-bottom: 60px;
}

/* 头部与高光标题 */
.page-header {
  padding: 40px 0 20px;
  text-align: center;
}

.highlight-title {
  display: inline-block;
  position: relative;
  font-size: 32px;
  color: #2d3436;
  font-weight: 700;
  margin-bottom: 10px;
  z-index: 1;
}

.highlight-title::after {
  content: "";
  position: absolute;
  bottom: 4px;
  left: -5%;
  width: 110%;
  height: 14px;
  background-color: #0d347c;
  opacity: 0.2;
  border-radius: 6px;
  z-index: -1;
  transition: all 0.3s ease;
}

.highlight-title:hover::after {
  opacity: 0.35;
  bottom: 6px;
}

.page-subtitle {
  font-size: 15px;
  color: #8c939d;
  letter-spacing: 1px;
}

/* 高级搜索控制台 */
.search-panel {
  display: flex;
  align-items: stretch;
  justify-content: center;
  gap: 16px;
  max-width: 860px;
  margin: 0 auto 40px;
}

.search-console {
  display: flex;
  align-items: center;
  flex: 1;
  background: #ffffff;
  padding: 10px 24px;
  border-radius: 50px; /* 大圆角胶囊形 */
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.04);
  transition: box-shadow 0.3s ease;
  min-height: 52px;
}

.search-console:hover {
  box-shadow: 0 12px 32px rgba(9, 132, 227, 0.08); /* 悬浮泛出微蓝光 */
}

.search-inputs {
  display: flex;
  align-items: center;
  width: 100%;
}

.divider {
  width: 1px;
  height: 24px;
  background-color: #e4e7ed;
  margin: 0 16px;
}

/* 深度覆盖 Element Plus 输入框样式，去除默认边框 */
:deep(.custom-input .el-input__wrapper),
:deep(.custom-select .el-select__wrapper) {
  box-shadow: none !important;
  background-color: transparent;
  padding: 0;
}

:deep(.custom-input input) {
  font-size: 15px;
  color: #303133;
}

:deep(.custom-select) {
  width: 160px;
}

.search-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  height: 52px;
  min-width: 90px;
  border-radius: 50px;
  padding: 0 20px;
  font-size: 15px;
  font-weight: bold;
  background: #2d597b;
  border: none;
}

.search-btn:hover {
  background: #2d597b;
  transform: translateY(-2px);
}

.btn-icon {
  width: 18px;
}

/* 商品卡片网格 */
.product-grid {
  margin-bottom: 40px;
}

.product-card {
  cursor: pointer;
  background: #ffffff;
  border-radius: 12px;
  overflow: hidden;
  transition: all 0.4s cubic-bezier(0.175, 0.885, 0.32, 1.275); /* 加入弹簧过渡曲线 */
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.03);
  margin-bottom: 24px;
  border: 1px solid transparent;
}

.product-card:hover {
  transform: translateY(-8px);
  /* 悬停时边框变微蓝，散发蓝色阴影，高级感拉满 */
  border-color: rgba(9, 132, 227, 0.1);
  box-shadow: 0 16px 32px rgba(9, 132, 227, 0.12);
}

.product-image {
  position: relative;
  width: 100%;
  aspect-ratio: 1 / 1; /* 完美的正方形比例 */
  overflow: hidden;
  background: #f8f9fa;
}

.product-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.5s ease;
}

.product-card:hover .product-image img {
  transform: scale(1.05); /* 图片缓慢放大 */
}

/* 促销标签 - 采用类似热点 NEW 的极简设计 */
.product-badge {
  position: absolute;
  top: 12px;
  left: 12px;
  background: linear-gradient(135deg, #ff7675 0%, #d63031 100%);
  color: white;
  padding: 4px 10px;
  border-radius: 6px;
  font-size: 12px;
  font-weight: bold;
  letter-spacing: 1px;
  box-shadow: 0 2px 8px rgba(214, 48, 49, 0.3);
}

.product-info {
  padding: 16px;
}

.product-name {
  font-size: 16px;
  font-weight: 600;
  color: #2d3436;
  margin-bottom: 6px;
  line-height: 1.4;
  /* 两行省略号 */
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  transition: color 0.3s;
}

.product-card:hover .product-name {
  color: #2d597b;
}

.product-desc {
  font-size: 13px;
  color: #a4b0be;
  margin-bottom: 16px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.product-footer {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
}

.product-price {
  display: flex;
  align-items: baseline;
}

.price-currency {
  font-size: 14px;
  color: #e74c3c;
  font-weight: bold;
  margin-right: 2px;
}

.price-current {
  font-size: 22px;
  font-weight: 800;
  color: #e74c3c;
  font-family:
    -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
}

.price-original {
  font-size: 12px;
  color: #b2bec3;
  text-decoration: line-through;
  margin-left: 8px;
}

/* 购物车按钮动效 */
.cart-btn-wrapper {
  transition: transform 0.2s;
}

.cart-btn-wrapper:hover {
  transform: scale(1.15);
}

:deep(.cart-btn) {
  background-color: #f1f8ff;
  border-color: #d2e9ff;
  color: #2d597b;
}

.cart-btn-wrapper:hover :deep(.cart-btn) {
  background-color: #2d597b;
  color: #ffffff;
  border-color: #2d597b;
}

/* 空状态和分页 */
.empty-state {
  padding: 80px 0;
  background: #ffffff;
  border-radius: 16px;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.02);
}

.pagination-container {
  display: flex;
  justify-content: center;
  margin-top: 40px;
}

/* 覆盖 Element Plus 分页样式，使其更柔和 */
:deep(.custom-pagination .el-pager li.is-active) {
  background-color: #2d597b;
  color: #fff;
  font-weight: bold;
}
:deep(.custom-pagination .el-pager li) {
  border-radius: 6px;
  background: #ffffff;
  margin: 0 4px;
}
</style>
