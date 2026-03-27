<template>
  <div class="mall-page">
    <Navbar />
    
    <div class="container">
      <div class="page-header">
        <h1 class="page-title">社区商城</h1>
        
        <div class="search-bar">
          <el-input 
            v-model="searchKeyword" 
            placeholder="搜索商品..."
            class="search-input"
            @keyup.enter="searchProducts"
            clearable
          >
             <template #prefix><el-icon><Search /></el-icon></template>
          </el-input>
          
           <el-select 
              v-model="selectedCategory" 
              placeholder="选择分类" 
              clearable 
              style="width: 180px; margin-left: 10px;"
              @change="handleSearch"
            >
              <el-option :value="0" label="所有分类" />
              <el-option v-for="c in categories" :key="c.id" :value="c.id" :label="c.name" />
           </el-select>
           
          <el-button type="primary" class="ml-2" @click="handleSearch">搜索</el-button>
        </div>
      </div>
      
      <div class="product-grid" v-loading="loading">
        <el-row :gutter="20">
          <el-col :xs="12" :sm="8" :md="6" :lg="6" :xl="4" v-for="product in products" :key="product.id">
            <div class="product-card" @click="goToDetail(product.id)">
               <div class="product-image">
                <img :src="product.image_url || 'https://via.placeholder.com/200'" :alt="product.name" />
                <div class="product-badge" v-if="product.is_promotion">促销</div>
              </div>
              <div class="product-info">
                <div class="product-name" :title="product.name">{{ product.name }}</div>
                <div class="product-desc" :title="product.description">{{ product.description }}</div>
                <div class="product-footer">
                  <div class="product-price">
                    <span class="price-current">¥{{ product.price }}</span>
                    <span class="price-original" v-if="product.original_price">¥{{ product.original_price }}</span>
                  </div>
                  <el-button type="primary" size="small" :icon="'ShoppingCart'" circle @click.stop="addToCart(product)" />
                </div>
              </div>
            </div>
          </el-col>
        </el-row>
      </div>
      
      <div class="empty-state" v-if="!loading && products.length === 0">
        <el-empty description="暂无商品" />
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
          />
        </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import Navbar from '@/components/layout/Navbar.vue'
import { getProductList, getCategories } from '@/api/product'
import { addToCart as addToCartApi } from '@/api/order'
import { useCartStore } from '@/stores/cart'
import { ElMessage } from 'element-plus'
import { Search, ShoppingCart } from '@element-plus/icons-vue'

const router = useRouter()
const cartStore = useCartStore()

const products = ref([])
const categories = ref([])
const searchKeyword = ref('')
const selectedCategory = ref(0)
const page = ref(1)
const size = ref(12) // Grid prefers multiple of 3/4
const total = ref(0)
const loading = ref(false)

const handleSearch = () => {
    // If searching (keyword change or category change), might want to reset page, but usually Keep page if just refreshing
    // Here we treat any explicit action as maybe "Refetch", but if category changed, reset page to 1
    // For simplicity, let's keep page unless manually reset. 
    // Wait, if filter changes, should reset to page 1.
    // 'handleSearch' is called by button/change event.
    // But pagination calls specific events.
    // Let's make a generic fetch and a reset-fetch
    fetchData()
}

// Check if called from pagination or filter
// Actually simpler: 
// On Filter Change -> page=1, fetch
// On Page Change -> fetch

const searchProducts = () => {
    page.value = 1
    fetchData()
}

const fetchData = async () => {
  loading.value = true
  try {
    const data = await getProductList({
      name: searchKeyword.value,
      category_id: selectedCategory.value || undefined,
      status: 1,
      page: page.value,
      size: size.value
    })
    products.value = data.list || []
    total.value = data.total
  } catch (error) {
    console.error('搜索失败:', error)
    ElMessage.error('加载商品失败')
  } finally {
      loading.value = false
  }
}

const fetchCategories = async () => {
    try {
        const res = await getCategories()
        categories.value = res || []
    } catch (e) {
        console.error(e)
    }
}

const goToDetail = (id) => {
  router.push(`/product/${id}`)
}

const addToCart = async (product) => {
  try {
    await addToCartApi({
      product_id: product.id,
      quantity: 1
    })
    ElMessage.success('已添加到购物车！')
    cartStore.fetchCart()
  } catch (error) {
    if (error.response?.status === 401) {
      ElMessage.warning('请先登录')
      router.push('/login')
    } else {
      ElMessage.error('添加失败')
    }
  }
}

onMounted(async () => {
  await fetchCategories()
  await fetchData()
})
</script>

<style scoped>
.mall-page {
  min-height: 100vh;
  padding-bottom: var(--spacing-xl);
}

.page-header {
  padding: var(--spacing-xl) 0;
}

.search-bar {
  display: flex;
  align-items: center;
  max-width: 800px;
  margin-top: var(--spacing-lg);
}

.search-input {
    width: 300px;
}

.product-grid {
  margin-bottom: 30px;
}

.product-card {
  cursor: pointer;
  background: white;
  border-radius: 8px;
  overflow: hidden;
  transition: all 0.3s;
  box-shadow: 0 2px 12px 0 rgba(0,0,0,0.05);
  margin-bottom: 20px;
  border: 1px solid var(--border-color-lighter);
}

.product-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 4px 16px 0 rgba(0,0,0,0.1);
}

.product-image {
  position: relative;
  width: 100%;
  padding-bottom: 100%; /* Square Aspect Ratio */
  /* OR fixed height */
  height: 0;
  overflow: hidden;
  background: #f5f7fa;
}

.product-image img {
  position: absolute;
  top: 0; left: 0;
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.product-badge {
  position: absolute;
  top: 10px;
  right: 10px;
  background: #f56c6c;
  color: white;
  padding: 2px 8px;
  border-radius: 12px;
  font-size: 12px;
}

.product-info {
  padding: 12px;
}

.product-name {
  font-size: 16px;
  font-weight: 600;
  color: #303133;
  margin-bottom: 6px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.product-desc {
  font-size: 12px;
  color: #909399;
  margin-bottom: 10px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.product-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.price-current {
  font-size: 18px;
  font-weight: 700;
  color: #f56c6c;
}

.price-original {
  font-size: 12px;
  color: #c0c4cc;
  text-decoration: line-through;
  margin-left: 6px;
}

.empty-state {
    padding: 50px;
    text-align: center;
}

.pagination-container {
    display: flex;
    justify-content: center;
    margin-top: 30px;
}

.ml-2 { margin-left: 8px; }
</style>
