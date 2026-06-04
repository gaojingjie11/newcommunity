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
          <div class="balance-info">
            <span class="summary-label">当前账户余额</span>
            <div class="balance-amount-row">
              <span class="currency">¥</span>
              <span class="amount">{{ (userInfo.balance || 0).toFixed(2) }}</span>
            </div>
          </div>
        </div>
        <div class="summary-actions">
          <el-button type="primary" size="large" class="recharge-btn" @click="showRechargeDialog = true">
            <el-icon class="btn-icon"><Plus /></el-icon>充值余额
          </el-button>
        </div>
      </div>

      <!-- 余额充值弹窗 -->
      <el-dialog
        v-model="showRechargeDialog"
        title="账户余额充值"
        width="460px"
        align-center
        class="recharge-dialog"
        :before-close="closeRechargeDialog"
      >
        <el-form :model="rechargeForm" label-position="top">
          <el-form-item label="充值金额 (元)">
            <el-input-number 
              v-model="rechargeForm.amount" 
              :min="0.01" 
              :precision="2" 
              :step="10" 
              style="width: 100%"
            />
          </el-form-item>
          <el-form-item label="支付方式">
            <el-radio-group v-model="rechargeForm.payType" class="pay-type-group">
              <el-radio-button label="alipay">支付宝 (沙箱)</el-radio-button>
              <el-radio-button label="mock">模拟充值</el-radio-button>
            </el-radio-group>
          </el-form-item>
        </el-form>
        <template #footer>
          <div class="dialog-footer">
            <el-button @click="closeRechargeDialog">取消</el-button>
            <el-button type="primary" :loading="rechargeLoading" @click="handleRechargeSubmit">确认充值</el-button>
          </div>
        </template>
      </el-dialog>

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
              <span class="type-tag" :class="getTypeClass(scope.row)">
                {{ getTypeLabel(scope.row) }}
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
          
          <el-table-column prop="remark" label="支付方式 / 关联单号" min-width="220" show-overflow-tooltip>
            <template #default="scope">
              <div class="memo-content">
                <span class="memo-main" style="font-weight: 600; color: #2c3e50;">
                  {{ getPaymentMethodLabel(scope.row) }}
                </span>
                <span v-if="scope.row.biz_id" class="memo-sub">单号: {{ scope.row.biz_id }}</span>
                <span v-else-if="scope.row.id" class="memo-sub">流水号: {{ scope.row.id }}</span>
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
import { useRouter, useRoute } from 'vue-router'
import Navbar from '@/components/layout/Navbar.vue'
import { getTransactionList, getWalletBalance, recharge } from '@/api/finance'
import dayjs from 'dayjs'
import { ElMessage } from 'element-plus'
// 引入必需的图标
import { ArrowLeft, Wallet, Plus } from '@element-plus/icons-vue'

const router = useRouter()
const route = useRoute()
const transactions = ref([])
const userInfo = ref({ balance: 0 })
const loading = ref(false)
const total = ref(0)
const page = ref(1)
const size = ref(10)

const showRechargeDialog = ref(false)
const rechargeLoading = ref(false)
const rechargeForm = ref({
  amount: 50.00,
  payType: 'alipay'
})

const closeRechargeDialog = () => {
  showRechargeDialog.value = false
  rechargeForm.value.amount = 50.00
  rechargeForm.value.payType = 'alipay'
}

const handleRechargeSubmit = async () => {
  if (rechargeForm.value.amount <= 0) {
    ElMessage.warning('请输入大于0的充值金额')
    return
  }
  rechargeLoading.value = true
  try {
    const res = await recharge(rechargeForm.value.amount, rechargeForm.value.payType)
    if (res && (res.code === 0 || !res.code)) {
      if (rechargeForm.value.payType === 'alipay' && res.pay_url) {
        ElMessage.success('正在跳转到支付宝支付收银台...')
        window.location.href = res.pay_url
      } else {
        ElMessage.success('充值成功！')
        closeRechargeDialog()
        fetchTransactions()
        fetchUser()
      }
    } else {
      ElMessage.error(res?.message || '充值失败')
    }
  } catch (error) {
    console.error(error)
    ElMessage.error('充值请求失败，请稍后重试')
  } finally {
    rechargeLoading.value = false
  }
}

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
        const res = await getWalletBalance()
        userInfo.value = { balance: Number(res?.balance || 0) }
    } catch (e) {}
}

const formatDate = (dateStr) => {
    if(!dateStr) return ''
    return dayjs(dateStr).format('YYYY-MM-DD HH:mm:ss')
}

const getTypeLabel = (row) => {
    const map = {
        'order_pay': '商城购买',
        'transfer': '用户转账',
        'recharge': '个人充值',
        'order_refund': '订单退款',
        'property_fee': '物业缴费'
    }
    return map[row.biz_type] || map[row.type] || '其他交易'
}

const getTypeClass = (row) => {
    const map = {
        'order_pay': 'type-1',
        'transfer': 'type-2',
        'recharge': 'type-3',
        'order_refund': 'type-4',
        'property_fee': 'type-5'
    }
    return map[row.biz_type] || 'type-default'
}

const getPaymentMethodLabel = (row) => {
    const remark = row.remark || ''
    if (remark.includes('支付宝') || (row.biz_id && row.biz_id.startsWith('RECH_') && !row.biz_id.includes('mock'))) {
        return '支付宝支付'
    }
    if (remark.includes('模拟充值') || (row.biz_id && row.biz_id.includes('mock'))) {
        return '系统模拟充值'
    }
    if (remark.includes('积分支付')) {
        return '积分支付'
    }
    if (remark.includes('积分+钱包')) {
        return '积分+钱包支付'
    }
    if (remark.includes('钱包支付') || row.biz_type === 'property_fee' || remark.includes('订单支付') || remark.includes('转账')) {
        return '钱包支付'
    }
    return remark || '钱包支付'
}

onMounted(() => {
    fetchTransactions()
    fetchUser()
    if (route.query.recharge === 'true') {
        showRechargeDialog.value = true
    }
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

.balance-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.summary-label {
  font-size: 14px;
  font-weight: 500;
  color: #7f8c8d;
}

.balance-amount-row {
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

.recharge-btn {
  background-color: #2d597b;
  border-color: #2d597b;
  font-weight: 600;
  border-radius: 8px;
  padding: 12px 24px;
  transition: all 0.3s ease;
  color: #fff;
  display: flex;
  align-items: center;
}

.recharge-btn:hover {
  background-color: #1e3f5a;
  border-color: #1e3f5a;
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(45, 89, 123, 0.2);
}

.btn-icon {
  margin-right: 6px;
}

.pay-type-group {
  display: flex;
  gap: 12px;
  width: 100%;
}

:deep(.pay-type-group .el-radio-button) {
  flex: 1;
}

:deep(.pay-type-group .el-radio-button__inner) {
  border-radius: 8px !important;
  border: 1px solid #dcdfe6 !important;
  padding: 12px 20px;
  width: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
}

:deep(.pay-type-group .el-radio-button__orig-radio:checked + .el-radio-button__inner) {
  background-color: #e8f4ff;
  border-color: #409eff !important;
  color: #409eff;
  box-shadow: none;
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
.type-1 { background: #fff7ed; color: #d97706; } /* 商城订单 */
.type-2 { background: #eef6ff; color: #2563eb; } /* 用户转账 */
.type-3 { background: #f0fdf4; color: #166534; } /* 账户充值 */
.type-4 { background: #fdf6f6; color: #e4393c; } /* 订单退款 */
.type-5 { background: #f3f4ff; color: #4f46e5; } /* 物业费 */
.type-default { background: #f3f4f6; color: #4b5563; }

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
