<template>
  <div class="favorite-page">
    <Navbar />
    
    <div class="container custom-container">
      <!-- 顶部返回导航 -->
      <div class="page-nav">
        <div class="back-btn" @click="$router.push('/profile')">
          <el-icon class="back-icon"><ArrowLeft /></el-icon> 
          <span>返回个人中心</span>
        </div>
      </div>

      
      <div class="product-grid" v-if="list.length > 0">
        <div class="product-card" v-for="item in list" :key="item.id" @click="goDetail(item.product.id)">
          
          <!-- 商品图橱窗 -->
          <div class="card-img-wrapper">
            <img :src="item.product.image_url || 'https://via.placeholder.com/200'" :alt="item.product.name" class="product-img" />
          </div>
          
          <!-- 商品信息 -->
          <div class="card-body">
            <h3 class="product-name" :title="item.product.name">{{ item.product.name }}</h3>
            
            <div class="card-footer">
              <div class="price-row">
                <span class="price"><span class="price-currency">¥</span>{{ item.product.price }}</span>
                <span class="original-price" v-if="item.product.original_price">¥{{ item.product.original_price }}</span>
              </div>
              
              <!-- 取消收藏按钮 -->
              <button class="btn-cancel-fav" @click.stop="removeFavorite(item.product.id)" title="取消收藏">
                <el-icon><Delete /></el-icon>
                <span>取消</span>
              </button>
            </div>
          </div>
        </div>
      </div>
      
      <!-- 优雅的空状态 -->
      <div v-else class="empty-wrapper">
        <el-empty description="您的收藏夹空空如也，快去发现好物吧" image-size="160">
          <button class="btn-go-mall" @click="$router.push('/mall')">去商城逛逛</button>
        </el-empty>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import Navbar from '@/components/layout/Navbar.vue'
import { getFavoriteList, deleteFavorite } from '@/api/favorite'
import { ElMessage } from 'element-plus'
import { ArrowLeft, Delete } from '@element-plus/icons-vue'

const router = useRouter()
const list = ref([])

const fetchList = async () => {
    try {
        const res = await getFavoriteList()
        list.value = res.list || []
    } catch (e) {
        console.error(e)
    }
}

const removeFavorite = async (pid) => {
    try {
        await deleteFavorite({ product_id: pid })
        ElMessage.success('已取消收藏')
        fetchList() // 刷新列表
    } catch (e) {
        ElMessage.error('操作失败')
    }
}

const goDetail = (id) => {
    router.push(`/product/${id}`)
}

onMounted(fetchList)
</script>

<style scoped>
/* 全局页面背景设定 */
.favorite-page {
  min-height: 100vh;
  background-color: #f8f9fa;
  padding-bottom: 80px;
}

/* 大气加宽容器 */
.custom-container {
  max-width: 1200px;
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

/* 统一的高光标题 */
.page-header {
  padding: 16px 0 32px;
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

/* 商品卡片网格 */
.product-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(260px, 1fr));
  gap: 24px;
}

.product-card {
  background: #ffffff;
  border-radius: 12px;
  overflow: hidden;
  transition: all 0.4s cubic-bezier(0.175, 0.885, 0.32, 1.275);
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.03);
  cursor: pointer;
  border: 1px solid transparent;
  display: flex;
  flex-direction: column;
}

.product-card:hover {
  transform: translateY(-6px);
  border-color: rgba(45, 89, 123, 0.1);
  box-shadow: 0 12px 24px rgba(45, 89, 123, 0.08);
}

/* 图片画中画效果 */
.card-img-wrapper {
  width: 100%;
  aspect-ratio: 1 / 1;
  background: #fbfcfd;
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
  border-bottom: 1px solid #f0f2f5;
}

.product-img {
  width: 100%;
  height: 100%;
  object-fit: contain;
  mix-blend-mode: multiply;
  transition: transform 0.5s ease;
}

.product-card:hover .product-img {
  transform: scale(1.05);
}

.card-body {
  padding: 20px;
  display: flex;
  flex-direction: column;
  flex: 1;
}

.product-name {
  font-size: 16px;
  font-weight: 600;
  color: #2d3436;
  margin: 0 0 16px 0;
  line-height: 1.5;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  transition: color 0.3s;
  flex: 1;
}

.product-card:hover .product-name {
  color: #2d597b;
}

.card-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-top: auto;
}

.price-row {
  display: flex;
  align-items: baseline;
}

.price {
  color: #e4393c;
  font-weight: 800;
  font-size: 22px;
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
}

.price-currency {
  font-size: 14px;
  margin-right: 2px;
}

.original-price {
  color: #b2bec3;
  text-decoration: line-through;
  font-size: 13px;
  margin-left: 8px;
  font-weight: normal;
}

/* 高级感取消收藏按钮 */
.btn-cancel-fav {
  display: flex;
  align-items: center;
  gap: 4px;
  background: transparent;
  color: #a4b0be;
  border: 1px solid #ebeef5;
  padding: 6px 12px;
  border-radius: 20px;
  font-size: 13px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.btn-cancel-fav:hover {
  background: #fef0f0;
  color: #e4393c;
  border-color: #fbc4c4;
}

/* 空状态 */
.empty-wrapper {
  background: #ffffff;
  border-radius: 16px;
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
  box-shadow: 0 4px 12px rgba(45, 89, 123, 0.2);
}
</style>