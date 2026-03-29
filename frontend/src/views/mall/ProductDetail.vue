<template>
  <div class="detail-page">
    <Navbar />
    
    <!-- ★ 重点修改：动态改变外层容器的宽度限制 ★ -->
    <div class="container custom-container" :class="{ 'is-wide': showComments }">
      <!-- 顶部返回导航 -->
      <div class="page-nav">
        <div class="back-btn" @click="$router.push('/mall')">
          <el-icon class="back-icon"><ArrowLeft /></el-icon> 
          <span>返回商城</span>
        </div>
      </div>

      <!-- 商品详情卡片 -->
      <div class="product-detail" v-if="product">
        
        <!-- 左侧主图橱窗 -->
        <div class="detail-image-showcase">
          <img :src="product.image_url || 'https://via.placeholder.com/400'" :alt="product.name" />
        </div>
        
        <!-- 中间信息面板 (★ 重点修改：锁定宽度，防止被挤压 ★) -->
        <div class="detail-info">
          <div class="info-header">
            <h1 class="detail-title">{{ product.name }}</h1>
            <p class="detail-desc">{{ product.description }}</p>
          </div>
          
          <!-- 价格展示区 -->
          <div class="detail-price-box">
            <div class="price-main">
              <span class="price-currency">¥</span>
              <span class="price-current">{{ product.price }}</span>
              <span class="price-original" v-if="product.original_price">¥{{ product.original_price }}</span>
            </div>
            <div class="price-tags" v-if="product.is_promotion">
              <span class="tag-minimal">限时特惠</span>
            </div>
          </div>
          
          <!-- 数据元信息 -->
          <div class="detail-meta">
            <div class="meta-item">
              <span class="meta-label">累计销量</span>
              <span class="meta-value">{{ product.sales }} <span class="meta-unit">件</span></span>
            </div>
            <div class="meta-divider"></div>
            <div class="meta-item">
              <span class="meta-label">当前库存</span>
              <span class="meta-value">{{ product.stock }} <span class="meta-unit">件</span></span>
            </div>
          </div>

          <!-- 商品评价展开触发器 -->
          <div class="review-trigger-box" @click="showComments = !showComments">
            <div class="trigger-left">
              <span class="trigger-title">商品评价</span>
              <span class="trigger-count" v-if="total">{{ total }}条</span>
            </div>
            <div class="trigger-right">
              <span class="trigger-text">{{ showComments ? '收起面板' : '点击查看' }}</span>
              <el-icon class="trigger-icon" :class="{ 'is-rotated': showComments }"><ArrowRight /></el-icon>
            </div>
          </div>
          
          <!-- 操作按钮区 -->
          <div class="detail-actions">
            <button class="action-btn btn-cart" @click="addToCart(false)">
              加入购物车
            </button>
            <button 
                class="action-btn" 
                :class="isFavorite ? 'btn-fav-active' : 'btn-fav'" 
                @click="toggleFavorite"
            >
              <el-icon class="fav-icon"><StarFilled v-if="isFavorite" /><Star v-else /></el-icon>
              {{ isFavorite ? '已收藏' : '加入收藏' }}
            </button>
          </div>
        </div>

        <!-- 右侧平移展开的评价面板 -->
        <div class="comments-side-panel" :class="{ 'is-expanded': showComments }">
          <div class="panel-inner">
            <div class="panel-header">
              <h3 class="panel-title">评价详情</h3>
              <el-icon class="close-panel-icon" @click="showComments = false"><Close /></el-icon>
            </div>
            
            <div class="panel-content custom-scrollbar">
              <!-- 评价输入框区 -->
              <div class="comment-form-box">
                  <el-input 
                    v-model="commentContent" 
                    type="textarea" 
                    :rows="3" 
                    placeholder="分享一下使用体验..." 
                    class="custom-textarea"
                  />
                  <div class="form-footer">
                      <div class="rating-wrapper">
                        <!-- 坚持黄色星星 -->
                        <el-rate v-model="commentRating" :colors="['#f39c12', '#f39c12', '#f39c12']" />
                      </div>
                      <button class="submit-btn" @click="submitComment">发表评价</button>
                  </div>
              </div>

              <!-- 评价列表 -->
              <div class="comment-list">
                  <div v-for="item in comments" :key="item.id" class="comment-item">
                      <div class="comment-avatar">
                          <img :src="item.user?.avatar || 'https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png'" />
                      </div>
                      <div class="comment-body">
                          <div class="comment-header">
                              <span class="username">{{ item.user?.real_name || item.user?.username || '匿名用户' }}</span>
                              <span class="time">{{ formatDate(item.created_at) }}</span>
                          </div>
                          <!-- 坚持黄色星星 -->
                          <el-rate :model-value="item.rating" disabled size="small" :colors="['#f39c12', '#f39c12', '#f39c12']" class="readonly-rate" />
                          <div class="comment-content">{{ item.content }}</div>
                      </div>
                  </div>
                  
                  <div class="empty-comment" v-if="comments.length === 0">
                    <el-empty description="暂无评价，来抢沙发吧！" image-size="80" />
                  </div>
              </div>

              <!-- 分页 -->
              <div class="pagination-container" v-if="total > 0">
                    <el-pagination
                      v-model:current-page="page"
                      v-model:page-size="size"
                      :total="total"
                      layout="prev, pager, next"
                      @current-change="fetchComments"
                      class="custom-pagination"
                      small
                    />
              </div>
            </div>
          </div>
        </div>

      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import Navbar from '@/components/layout/Navbar.vue'
import { getProductDetail } from '@/api/product'
import { addToCart as addToCartApi } from '@/api/order'
import { addFavorite, deleteFavorite, checkFavorite } from '@/api/favorite'
import { useCartStore } from '@/stores/cart'
import { ElMessage } from 'element-plus'
import { ArrowLeft, Star, StarFilled, ArrowRight, Close } from '@element-plus/icons-vue' 

const route = useRoute()
const router = useRouter()
const cartStore = useCartStore()
const product = ref(null)
const isFavorite = ref(false)

const showComments = ref(false) 

const addToCart = async (silent) => {
  if (typeof silent !== 'boolean') silent = false
  try {
    await addToCartApi({
      product_id: product.value.id,
      quantity: 1
    })
    if (!silent) {
      ElMessage.success('已成功添加到购物车！')
    }
    cartStore.fetchCart()
  } catch (error) {
    if (error.response?.status === 401) {
       ElMessage.warning('请先登录后再进行操作')
       router.push('/login')
    } else {
       ElMessage.error('添加失败: ' + error.message)
    }
  }
}

const toggleFavorite = async () => {
    try {
        if (isFavorite.value) {
            await deleteFavorite({ product_id: product.value.id })
            isFavorite.value = false
            ElMessage.success('已取消收藏')
        } else {
            await addFavorite({ product_id: product.value.id })
            isFavorite.value = true
            ElMessage.success('收藏成功')
        }
    } catch (e) {
        if (e.response?.status === 401) {
            router.push('/login')
        } else {
            ElMessage.error('操作失败')
        }
    }
}

const checkFavStatus = async () => {
    try {
        const res = await checkFavorite(product.value.id)
        isFavorite.value = res.is_favorite
    } catch (e) {
        // ignore
    }
}

// --- Comment Logic ---
import { getCommentList, createComment } from '@/api/comment'
import dayjs from 'dayjs'

const comments = ref([])
const total = ref(0)
const page = ref(1)
const size = ref(10)
const commentContent = ref('')
const commentRating = ref(5)

const fetchComments = async () => {
    try {
        const res = await getCommentList({ 
            product_id: route.params.id,
            page: page.value,
            size: size.value 
        })
        comments.value = res.list || []
        total.value = res.total
    } catch (e) {
        console.error(e)
    }
}

const submitComment = async () => {
    if (!commentContent.value.trim()) {
        ElMessage.warning('请输入评价内容')
        return
    }
    try {
        await createComment({
            product_id: Number(route.params.id),
            content: commentContent.value,
            rating: commentRating.value
        })
        ElMessage.success('评价发表成功')
        commentContent.value = ''
        fetchComments()
        showComments.value = true
    } catch (e) {
         if (e.response?.status === 401) {
            router.push('/login')
         } else {
             ElMessage.error(e.response?.data?.msg || '评价失败')
         }
    }
}

const formatDate = (date) => dayjs(date).format('YYYY-MM-DD HH:mm')

onMounted(async () => {
  try {
    product.value = await getProductDetail(route.params.id)
    checkFavStatus()
    fetchComments()
  } catch (error) {
    console.error('获取商品详情失败:', error)
  }
})
</script>

<style scoped>
/* 全局页面背景 */
.detail-page {
  min-height: 100vh;
  background-color: #f8f9fa;
  padding-bottom: 80px;
}

/* ★ 核心动态容器：默认窄，展开评论时变宽 ★ */
.custom-container {
  max-width: 980px; /* 默认宽度，让详情页更聚焦，不空旷 */
  transition: max-width 0.5s cubic-bezier(0.25, 0.8, 0.25, 1);
  margin: 0 auto;
}

.custom-container.is-wide {
  max-width: 1400px; /* 展开评论区时，释放空间让右侧滑出 */
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

/* 商品主卡片区 */
.product-detail {
  display: flex;
  align-items: stretch;
  background: #ffffff;
  border-radius: 12px;
  padding: 32px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.04);
  margin-bottom: 24px;
  width: fit-content; /* 让内部元素撑开盒子 */
  margin: 0 auto;     /* 保持卡片居中 */
  transition: all 0.5s cubic-bezier(0.25, 0.8, 0.25, 1);
  overflow: hidden;
}

/* 左侧：图片橱窗 */
.detail-image-showcase {
  width: 420px;
  flex-shrink: 0; /* 绝对不允许图片变窄 */
  background: #f8f9fa;
  border-radius: 8px;
  overflow: hidden;
  display: flex;
  align-items: center;
  justify-content: center;
  aspect-ratio: 1 / 1;
}

.detail-image-showcase img {
  width: 100%;
  height: 100%;
  object-fit: contain;
  mix-blend-mode: multiply;
}

/* ★ 中间信息面板 ★ */
.detail-info {
  width: 440px;     /* 固定具体宽度 */
  flex-shrink: 0;   /* 绝对不允许信息区被压缩！解决你截图中的痛点 */
  display: flex;
  flex-direction: column;
  padding: 0 0 0 40px;
}

.info-header {
  margin-bottom: 24px;
}

.detail-title {
  font-size: 26px;
  font-weight: 600;
  color: #2c3e50;
  line-height: 1.4;
  margin-bottom: 12px;
}

.detail-desc {
  font-size: 14px;
  color: #7f8c8d;
  line-height: 1.6;
}

/* 价格区域 */
.detail-price-box {
  background: #fcfdfe;
  border: 1px solid #f0f2f5;
  border-radius: 8px;
  padding: 16px 20px;
  margin-bottom: 24px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.price-main {
  display: flex;
  align-items: baseline;
}

.price-currency {
  font-size: 16px;
  color: #e4393c;
  font-weight: bold;
  margin-right: 2px;
}

.price-current {
  font-size: 32px;
  font-weight: 700;
  color: #e4393c;
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
}

.price-original {
  font-size: 14px;
  color: #a4b0be;
  text-decoration: line-through;
  margin-left: 10px;
}

.tag-minimal {
  border: 1px solid #2d597b;
  color: #2d597b; 
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 12px;
}

/* 销量库存数据 */
.detail-meta {
  display: flex;
  align-items: center;
  gap: 32px;
  margin-bottom: 20px;
  padding: 0 10px;
}

.meta-divider {
  width: 1px;
  height: 16px;
  background-color: #ebeef5;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 8px;
}

.meta-label {
  font-size: 13px;
  color: #909399;
}

.meta-value {
  font-size: 16px;
  font-weight: 600;
  color: #303133;
}

.meta-unit {
  font-size: 12px;
  font-weight: normal;
  color: #909399;
}

/* 高级触发器条 */
.review-trigger-box {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: #f8f9fa;
  padding: 16px 20px;
  border-radius: 8px;
  cursor: pointer;
  margin-bottom: 32px;
  border: 1px solid transparent;
  transition: all 0.3s ease;
}

.review-trigger-box:hover {
  background: #f2f5f8;
  border-color: #2d597b;
}

.trigger-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.trigger-title {
  font-size: 15px;
  font-weight: 600;
  color: #303133;
}

.trigger-count {
  background: #2d597b;
  color: white;
  font-size: 12px;
  padding: 2px 8px;
  border-radius: 12px;
}

.trigger-right {
  display: flex;
  align-items: center;
  gap: 4px;
  color: #8c939d;
  font-size: 14px;
}

.trigger-icon {
  transition: transform 0.3s;
}
.trigger-icon.is-rotated {
  transform: rotate(180deg);
}

/* 底部操作按钮 */
.detail-actions {
  display: flex;
  gap: 16px;
  margin-top: auto;
}

.action-btn {
  height: 48px;
  border-radius: 6px; 
  font-size: 16px;
  font-weight: 600;
  border: none;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
}

.btn-cart {
  flex: 2;
  background: #2d597b;
  color: white;
}

.btn-cart:hover {
  background: #1f435d;
}

.btn-fav {
  flex: 1;
  background: #ffffff;
  color: #606266;
  border: 1px solid #dcdfe6;
}

.btn-fav:hover {
  border-color: #2d597b;
  color: #2d597b;
}

.btn-fav-active {
  flex: 1;
  background: #f4f7f9;
  color: #2d597b;
  border: 1px solid #2d597b;
}

.fav-icon {
  margin-right: 6px;
  font-size: 18px;
  color: inherit;
}

/* ================= 右侧隐藏评价面板 ================= */
.comments-side-panel {
  max-width: 0;      
  opacity: 0;        
  visibility: hidden;
  overflow: hidden;
  transition: max-width 0.5s cubic-bezier(0.25, 0.8, 0.25, 1), 
              opacity 0.4s ease, 
              margin-left 0.4s ease, 
              padding-left 0.4s ease,
              visibility 0.5s;
  display: flex;
  flex-direction: column;
}

/* 展开状态 */
.comments-side-panel.is-expanded {
  max-width: 420px; /* 固定的评论区宽度 */
  opacity: 1;
  visibility: visible;
  margin-left: 40px;
  padding-left: 32px;
  border-left: 1px solid #f0f2f5;
}

.panel-inner {
  width: 387px; 
  height: 100%;
  max-height: 580px; /* 和左侧区域对齐的高度 */
  display: flex;
  flex-direction: column;
}

.panel-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-bottom: 16px;
  border-bottom: 1px solid #f0f2f5;
  margin-bottom: 16px;
}

.panel-title {
  font-size: 16px;
  font-weight: 600;
  color: #2c3e50;
  margin: 0;
  border-left: 4px solid #2d597b;
  padding-left: 8px;
}

.close-panel-icon {
  font-size: 20px;
  color: #909399;
  cursor: pointer;
  transition: color 0.3s;
}
.close-panel-icon:hover {
  color: #e4393c;
}

/* 面板内部定制滚动条 */
.panel-content {
  flex: 1;
  overflow-y: auto;
  padding-right: 8px;
}

.custom-scrollbar::-webkit-scrollbar {
  width: 6px;
}
.custom-scrollbar::-webkit-scrollbar-track {
  background: transparent;
}
.custom-scrollbar::-webkit-scrollbar-thumb {
  background: #dcdfe6;
  border-radius: 3px;
}
.custom-scrollbar::-webkit-scrollbar-thumb:hover {
  background: #c0c4cc;
}

/* 表单内部 */
.comment-form-box {
  margin-bottom: 24px;
}

:deep(.custom-textarea .el-textarea__inner) {
  box-shadow: none;
  border: 1px solid #dcdfe6;
  border-radius: 6px;
  transition: all 0.3s;
}

:deep(.custom-textarea .el-textarea__inner:focus) {
  border-color: #2d597b; 
  box-shadow: 0 0 0 2px rgba(45, 89, 123, 0.1);
}

.form-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 12px;
}

.submit-btn {
  padding: 6px 16px;
  border-radius: 4px;
  background: #2d597b; 
  color: white;
  border: none;
  font-size: 13px;
  cursor: pointer;
  transition: all 0.2s;
}

.submit-btn:hover {
  background: #1f435d;
}

/* 评价列表内部 */
.comment-item {
  display: flex;
  gap: 12px;
  padding: 16px 0;
  border-bottom: 1px solid #f2f3f5;
}

.comment-item:last-child {
  border-bottom: none;
}

.comment-avatar img {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  object-fit: cover;
}

.comment-body {
  flex: 1;
}

.comment-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 4px;
}

.username {
  font-weight: 600;
  font-size: 13px;
  color: #303133;
}

.time {
  color: #a4b0be;
  font-size: 12px;
}

.readonly-rate {
  margin-bottom: 6px;
}

.comment-content {
  color: #606266;
  line-height: 1.5;
  font-size: 13px;
}

.empty-comment {
  padding: 40px 0;
}

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: center;
}

/* 分页器颜色 */
:deep(.custom-pagination .el-pager li.is-active) {
  background-color: #2d597b;
  color: #fff;
  border-radius: 4px;
}
:deep(.custom-pagination .el-pager li:hover) {
  color: #2d597b;
}

/* 适配移动端 */
@media (max-width: 900px) {
  .product-detail {
    flex-direction: column;
    padding: 20px;
  }
  .detail-image-showcase {
    width: 100%;
  }
  .detail-info {
    width: 100%;
    padding: 20px 0 0 0;
  }
  .comments-side-panel.is-expanded {
    max-width: 100%;
    margin-left: 0;
    padding-left: 0;
    border-left: none;
    border-top: 1px solid #f0f2f5;
    margin-top: 24px;
    padding-top: 24px;
  }
  .panel-inner {
    width: 100%;
    max-height: none;
  }
}
</style>