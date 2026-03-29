<template>
  <div class="admin-child-page">
    <Navbar />
    <div class="container custom-container">
      
      <div class="page-nav">
        <div class="back-btn" @click="$router.push('/admin')">
          <el-icon class="back-icon"><ArrowLeft /></el-icon> 
          <span>返回管理后台</span>
        </div>
      </div>

      <div class="page-header">
        
        <div class="header-actions" style="margin-left: 900px;">
          <el-input 
            v-model="searchUserId" 
            placeholder="按用户ID搜索订单..." 
            class="search-input"
            clearable 
            @clear="handleSearch"
            @keyup.enter="handleSearch"
          >
            <template #append>
              <el-button @click="handleSearch">
                <el-icon><Search /></el-icon>
              </el-button>
            </template>
          </el-input>
        </div>
      </div>

      <div class="table-wrapper">
        <el-table 
          :data="orders" 
          class="custom-table" 
          v-loading="loadingData"
          :empty-text="'暂无相关订单记录'"
        >
          <el-table-column type="expand">
            <template #default="scope">
              <div class="expand-detail-container">
                <div class="expand-title">
                  <el-icon><Goods /></el-icon> 商品明细清单
                </div>
                <el-table :data="scope.row.items" class="inner-table" size="small">
                  <el-table-column label="商品预览" width="100" align="center">
                    <template #default="itemScope">
                      <div class="inner-thumb-wrapper">
                        <img :src="itemScope.row.product?.image_url" class="inner-thumb" />
                      </div>
                    </template>
                  </el-table-column>
                  <el-table-column prop="product.name" label="商品名称" min-width="200" />
                  <el-table-column label="销售单价" width="120">
                    <template #default="itemScope">¥{{ itemScope.row.price }}</template>
                  </el-table-column>
                  <el-table-column prop="quantity" label="购买数量" width="100" align="center" />
                  <el-table-column label="小计" width="120" align="right">
                    <template #default="itemScope">
                      <strong class="text-primary">¥{{ (itemScope.row.price * itemScope.row.quantity).toFixed(2) }}</strong>
                    </template>
                  </el-table-column>
                </el-table>
              </div>
            </template>
          </el-table-column>

          <el-table-column prop="order_no" label="流水订单号" min-width="180">
            <template #default="{ row }">
              <code class="order-no-code">{{ row.order_no }}</code>
            </template>
          </el-table-column>
          
          <el-table-column prop="store.name" label="服务门店" min-width="150">
            <template #default="{ row }">
              <span class="store-text">{{ row.store?.name || '社区总店' }}</span>
            </template>
          </el-table-column>

          <el-table-column prop="user_id" label="买家ID" width="100" align="center" />

          <el-table-column label="实付金额" width="140" align="right">
            <template #default="{ row }">
              <span class="amount-text">¥{{ row.total_amount }}</span>
            </template>
          </el-table-column>

          <el-table-column label="交易状态" width="120" align="center">
            <template #default="scope">
              <span class="status-badge" :class="`status-${scope.row.status}`">
                {{ getStatusText(scope.row.status) }}
              </span>
            </template>
          </el-table-column>

          <el-table-column label="下单时间" width="180">
            <template #default="scope">
              <div class="time-cell">
                <el-icon><Clock /></el-icon>
                <span>{{ formatDate(scope.row.created_at) }}</span>
              </div>
            </template>
          </el-table-column>

          <el-table-column label="管理操作" width="120" fixed="right" align="center">
            <template #default="scope">
              <button 
                v-if="scope.row.status === 1" 
                class="action-btn btn-primary btn-sm"
                @click="handleShip(scope.row)"
                :disabled="loadingMap[scope.row.id]"
              >
                <el-icon v-if="loadingMap[scope.row.id]" class="is-loading"><Loading /></el-icon>
                立即发货
              </button>
              <span v-else class="text-secondary">-</span>
            </template>
          </el-table-column>
        </el-table>

        <div class="pagination-container" v-if="total > 0">
          <el-pagination
            v-model:current-page="page"
            v-model:page-size="size"
            :page-sizes="[10, 20, 50, 100]"
            layout="total, sizes, prev, pager, next, jumper"
            :total="total"
            @size-change="fetchOrders"
            @current-change="fetchOrders"
            class="custom-pagination"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, reactive } from 'vue'
import { useRouter } from 'vue-router'
import Navbar from '@/components/layout/Navbar.vue'
import { getAdminOrderList, shipOrder } from '@/api/admin'
import dayjs from 'dayjs'
import { ElMessage } from 'element-plus'
// 引入图标
import { ArrowLeft, Clock, Goods, Search, Loading } from '@element-plus/icons-vue'

const router = useRouter()
const orders = ref([])
const loadingMap = reactive({})
const loadingData = ref(false)
const total = ref(0)
const page = ref(1)
const size = ref(10)
const searchUserId = ref('')

const formatDate = (date) => dayjs(date).format('YYYY-MM-DD HH:mm')

const getStatusText = (status) => {
  const map = { 0: '待支付', 1: '待发货', 2: '已发货', 3: '已完成', 40: '已取消' }
  return map[status] || status
}

const fetchOrders = async () => {
  loadingData.value = true
  try {
    const params = { page: page.value, size: size.value }
    if (searchUserId.value) {
        params.user_id = searchUserId.value
    }
    const res = await getAdminOrderList(params)
    if (res.list) {
        orders.value = res.list
        total.value = res.total
    } else {
        orders.value = res
        total.value = res.length
    }
  } catch (error) {
    console.error(error)
  } finally {
      loadingData.value = false
  }
}

const handleShip = async (order) => {
  if (loadingMap[order.id]) return
  loadingMap[order.id] = true
  
  try {
    await shipOrder({ id: order.id })
    ElMessage.success('发货指令已下达，更新状态成功')
    fetchOrders()
  } catch (error) {
    ElMessage.error('发货操作失败')
  } finally {
    loadingMap[order.id] = false
  }
}

const handleSearch = () => {
    page.value = 1
    fetchOrders()
}
 
onMounted(() => {
  fetchOrders()
})
</script>

<style scoped>
/* 全局页面底色与容器 */
.admin-child-page { min-height: 100vh; background-color: #f8f9fa; padding-bottom: 80px; }
.custom-container { max-width: 1200px; margin: 0 auto; }

/* 顶部返回导航 */
.page-nav { padding: 24px 0 16px; }
.back-btn {
  display: inline-flex; align-items: center; color: #606266; font-size: 15px;
  cursor: pointer; transition: color 0.3s; padding: 8px 16px 8px 0;
}
.back-btn:hover { color: #2d597b; }
.back-icon { margin-right: 6px; font-size: 16px; }

/* 统一的高光标题与顶部动作区 */
.page-header { display: flex; justify-content: space-between; align-items: center; padding: 16px 0 32px; }
.page-title { display: inline-block; position: relative; font-size: 32px; color: #2c3e50; font-weight: 700; margin: 0; z-index: 1; }
.page-title::after { content: ''; position: absolute; bottom: 4px; left: -5%; width: 110%; height: 14px; background-color: #2d597b; opacity: 0.15; border-radius: 6px; z-index: -1; transition: all 0.3s ease; }

/* 搜索框深度美化 */
:deep(.search-input .el-input__wrapper) {
  border-radius: 20px 0 0 20px; box-shadow: 0 2px 10px rgba(0,0,0,0.02); background: #ffffff;
}
:deep(.search-input .el-input-group__append) {
  border-radius: 0 20px 20px 0; background: #2d597b; color: white; border: none;
}

/* 核心表格容器 */
.table-wrapper { background: #ffffff; border-radius: 16px; padding: 24px 32px; box-shadow: 0 4px 20px rgba(0, 0, 0, 0.02); }

/* Element Plus 表格定制 */
:deep(.custom-table) { --el-table-border-color: transparent; border-radius: 8px; overflow: hidden; }
:deep(.custom-table th.el-table__cell) { font-weight: 600; font-size: 14px; padding: 18px 0; border-bottom: 1px solid #ebeef5; background: #fbfcfd; }
:deep(.custom-table td.el-table__cell) { padding: 20px 0; border-bottom: 1px dashed #f0f2f5; }
:deep(.custom-table::before) { display: none; }

/* 单元格排版 */
.order-no-code { font-family: monospace; background: #f4f4f5; padding: 2px 6px; border-radius: 4px; color: #606266; }
.store-text { font-weight: 600; color: #2d597b; }
.amount-text { font-size: 18px; font-weight: 800; color: #e4393c; font-family: sans-serif; }
.time-cell { display: flex; align-items: center; gap: 6px; color: #909399; }
.text-secondary { color: #dcdfe6; }

/* 状态标签定制 */
.status-badge { padding: 4px 12px; border-radius: 20px; font-size: 12px; font-weight: bold; }
.status-0 { background: #fef0f0; color: #f56c6c; } /* 待支付 */
.status-1 { background: #fff7ed; color: #d97706; } /* 待发货 - 橘黄提示 */
.status-2 { background: #f0f7ff; color: #2d597b; } /* 已发货 - 品牌蓝 */
.status-3 { background: #f0fdf4; color: #00b894; } /* 已完成 */
.status-40 { background: #f4f4f5; color: #909399; } /* 已取消 */

/* ★ 展开行深度美化 ★ */
.expand-detail-container { padding: 20px 40px; background: #fbfcfd; border-radius: 12px; margin: 10px 0; }
.expand-title { font-size: 14px; font-weight: 700; color: #2c3e50; margin-bottom: 16px; display: flex; align-items: center; gap: 8px; }

:deep(.inner-table) { background: transparent !important; }
:deep(.inner-table tr) { background: transparent !important; }
.inner-thumb-wrapper { width: 50px; height: 50px; border-radius: 6px; overflow: hidden; border: 1px solid #ebeef5; background: #fff; }
.inner-thumb { width: 100%; height: 100%; object-fit: contain; mix-blend-mode: multiply; }

/* 定制按钮 */
.action-btn { padding: 8px 20px; border-radius: 20px; font-size: 13px; font-weight: 600; cursor: pointer; transition: all 0.3s; border: none; }
.btn-primary { background: #2d597b; color: #ffffff; box-shadow: 0 4px 12px rgba(45, 89, 123, 0.2); }
.btn-primary:hover:not(:disabled) { background: #1f435d; transform: translateY(-1px); }
.action-btn:disabled { opacity: 0.6; cursor: not-allowed; }

/* 分页器 */
.pagination-container { display: flex; justify-content: flex-end; margin-top: 32px; padding-top: 16px; border-top: 1px solid #f0f2f5; }
:deep(.custom-pagination .el-pager li.is-active) { background-color: #2d597b !important; color: #fff; border-radius: 4px; }
</style>