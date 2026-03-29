<template>
  <div class="transaction-page">
    <Navbar />
    
    <div class="container custom-container">
      <!-- 顶部返回导航 -->
      <div class="page-nav">
        <div class="back-btn" @click="$router.push('/profile')">
          <el-icon class="back-icon"><ArrowLeft /></el-icon> 
          <span>返回个人中心</span>
        </div>
      </div>

      
      <!-- 高级资产摘要卡片 -->
      <div class="balance-summary-card">
        <div class="summary-left">
          <div class="icon-wrapper">
            <el-icon><Wallet /></el-icon>
          </div>
          <span class="summary-label">当前账户余额</span>
        </div>
        <div class="summary-right">
          <span class="currency">¥</span>
          <span class="amount">{{ (userInfo.balance || 0).toFixed(2) }}</span>
        </div>
      </div>

      <!-- 深度定制的流水表格 -->
      <div class="table-wrapper">
        <el-table 
          :data="transactions" 
          style="width: 100%" 
          v-loading="loading"
          class="custom-table"
          :empty-text="'暂无账单流水记录'"
        >
          <el-table-column prop="created_at" label="交易时间" min-width="180">
            <template #default="scope">
              <span class="time-text">{{ formatDate(scope.row.created_at) }}</span>
            </template>
          </el-table-column>
          
          <el-table-column prop="type" label="交易类型" min-width="120">
            <template #default="scope">
              <span class="type-tag" :class="`type-${scope.row.type}`">
                {{ getTypeLabel(scope.row.type) }}
              </span>
            </template>
          </el-table-column>
          
          <el-table-column prop="amount" label="交易金额" min-width="150" align="right">
            <template #default="scope">
              <span class="amount-text" :class="scope.row.amount > 0 ? 'is-income' : 'is-expense'">
                {{ scope.row.amount > 0 ? '+' : '' }}{{ scope.row.amount.toFixed(2) }}
              </span>
            </template>
          </el-table-column>
          
          <el-table-column prop="related_id" label="单号 / 备注" min-width="220" show-overflow-tooltip>
            <template #default="scope">
              <div class="memo-content">
                <span class="memo-main">相关ID: {{ scope.row.related_id }}</span>
                <span v-if="scope.row.id" class="memo-sub">流水号: {{ scope.row.id }}</span>
              </div>
            </template>
          </el-table-column>
          
          <!-- 空状态插槽 -->
          <template #empty>
            <el-empty description="暂无账单流水记录" image-size="120" />
          </template>
        </el-table>
        
        <!-- 分页器 -->
        <div class="pagination-container" v-if="total > 0">
          <el-pagination
            v-model:current-page="page"
            v-model:page-size="size"
            :page-sizes="[10, 20, 50]"
            layout="total, sizes, prev, pager, next, jumper"
            :total="total"
            @size-change="fetchTransactions"
            @current-change="fetchTransactions"
            class="custom-pagination"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import Navbar from '@/components/layout/Navbar.vue'
import { getTransactionList } from '@/api/finance'
import { getUserInfo } from '@/api/user'
import dayjs from 'dayjs'
// 引入必需的图标
import { ArrowLeft, Wallet } from '@element-plus/icons-vue'

const router = useRouter()
const transactions = ref([])
const userInfo = ref({ balance: 0 })
const loading = ref(false)
const total = ref(0)
const page = ref(1)
const size = ref(10)

const fetchTransactions = async () => {
    loading.value = true
    try {
        const res = await getTransactionList({ page: page.value, size: size.value })
        if (res.list) {
            transactions.value = res.list
            total.value = res.total
        } else if (Array.isArray(res)) {
            transactions.value = res
            total.value = res.length
        }
    } catch (e) {
        console.error(e)
    } finally {
        loading.value = false
    }
}

const fetchUser = async () => {
    try {
        const res = await getUserInfo()
        userInfo.value = res
    } catch (e) {}
}

const formatDate = (dateStr) => {
    if(!dateStr) return ''
    return dayjs(dateStr).format('YYYY-MM-DD HH:mm:ss')
}

const getTypeLabel = (type) => {
    const map = {
        1: '商城订单',
        2: '物业费',
        3: '账户充值',
        4: '系统转账'
    }
    return map[type] || '其他交易'
}

onMounted(() => {
    fetchTransactions()
    fetchUser()
})
</script>

<style scoped>
/* 全局页面背景设定 */
.transaction-page {
  min-height: 100vh;
  background-color: #f8f9fa;
  padding-bottom: 80px;
}

/* 大气加宽容器 */
.custom-container {
  max-width: 1000px;
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

/* ★ 高级资产摘要卡片 ★ */
.balance-summary-card {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: linear-gradient(135deg, #ffffff 0%, #f0f7ff 100%);
  border: 1px solid #e1f0ff;
  border-radius: 16px;
  padding: 28px 40px;
  margin-bottom: 32px;
  box-shadow: 0 4px 20px rgba(45, 89, 123, 0.05);
}

.summary-left {
  display: flex;
  align-items: center;
  gap: 16px;
}

.icon-wrapper {
  width: 48px;
  height: 48px;
  background: #2d597b;
  color: #ffffff;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  box-shadow: 0 4px 12px rgba(45, 89, 123, 0.2);
}

.summary-label {
  font-size: 18px;
  font-weight: 600;
  color: #2c3e50;
}

.summary-right {
  display: flex;
  align-items: baseline;
  color: #2d597b;
}

.currency {
  font-size: 20px;
  font-weight: 600;
  margin-right: 4px;
}

.amount {
  font-size: 40px;
  font-weight: 800;
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
  letter-spacing: -1px;
}

/* ★ 表格容器与深度美化 ★ */
.table-wrapper {
  background: #ffffff;
  border-radius: 16px;
  padding: 24px 32px;
  box-shadow: 0 2px 16px rgba(0, 0, 0, 0.03);
}

/* 覆盖 Element Plus 的默认表格样式，去后台感 */
:deep(.custom-table) {
  --el-table-border-color: transparent;
  --el-table-header-bg-color: #fbfcfd;
  --el-table-header-text-color: #606266;
}

:deep(.custom-table th.el-table__cell) {
  background-color: #fbfcfd;
  color: #606266;
  font-weight: 600;
  font-size: 14px;
  padding: 16px 0;
  border-bottom: 1px solid #ebeef5;
}

:deep(.custom-table td.el-table__cell) {
  padding: 20px 0;
  border-bottom: 1px dashed #f0f2f5;
}

/* 去除最底部的实体线 */
:deep(.custom-table::before) {
  display: none;
}

/* 表格内元素排版 */
.time-text {
  color: #606266;
  font-size: 14px;
}

/* 定制类型标签 */
.type-tag {
  display: inline-block;
  padding: 4px 12px;
  border-radius: 6px;
  font-size: 12px;
  font-weight: 600;
}
.type-1 { background: #fff7ed; color: #d97706; } /* 商城订单 - 橙色 */
.type-2 { background: #f0f7ff; color: #0984e3; } /* 物业费 - 蓝色 */
.type-3 { background: #f0fdf4; color: #166534; } /* 充值 - 绿色 */
.type-4 { background: #fdf6f6; color: #e4393c; } /* 转账 - 红色 */

/* 交易金额字体与颜色 */
.amount-text {
  font-size: 18px;
  font-weight: 700;
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
}
.is-income { color: #00b894; } /* 收入用绿色 */
.is-expense { color: #2c3e50; } /* 支出用沉稳的深色 */

/* 备注信息 */
.memo-content {
  display: flex;
  flex-direction: column;
  gap: 4px;
}
.memo-main {
  font-size: 14px;
  color: #303133;
}
.memo-sub {
  font-size: 12px;
  color: #a4b0be;
}

/* 定制分页器 */
.pagination-container {
  display: flex;
  justify-content: flex-end;
  margin-top: 32px;
  padding-top: 16px;
}

:deep(.custom-pagination .el-pager li.is-active) {
  background-color: #2d597b;
  color: #fff;
  border-radius: 4px;
}
:deep(.custom-pagination .el-pager li:hover) {
  color: #2d597b;
}

/* 响应式适配 */
@media (max-width: 768px) {
  .balance-summary-card {
    flex-direction: column;
    align-items: flex-start;
    gap: 16px;
    padding: 24px;
  }
  .table-wrapper {
    padding: 16px;
  }
}
</style>