<template>
  <div class="admin-child-page">
    <Navbar />
    <div class="container custom-container">
      
      <!-- 极简顶部区：包含返回导航与操作按钮，去除了大标题 -->
      <div class="top-bar">
        <div class="back-btn" @click="$router.push('/admin')">
          <el-icon class="back-icon"><ArrowLeft /></el-icon> 
          <span>返回管理后台</span>
        </div>
        <div class="header-actions">
          <button class="action-btn btn-primary" @click="openCreateModal">
            <el-icon style="margin-right: 4px;"><Plus /></el-icon> 新增账单
          </button>
        </div>
      </div>

      <!-- 深度美化的表格容器 -->
      <div class="table-wrapper">
        <el-table :data="list" class="custom-table" style="width: 100%" v-loading="loadingData" :empty-text="'暂无物业费账单记录'">
          <el-table-column prop="id" label="账单 ID" width="100" align="center" />
          
          <el-table-column prop="user_id" label="绑定用户 ID" min-width="120" align="center">
            <template #default="{ row }">
              <span class="user-id-text">{{ row.user_id }}</span>
            </template>
          </el-table-column>
          
          <el-table-column prop="month" label="账单月份" min-width="140">
            <template #default="{ row }">
              <div class="icon-text-cell">
                <el-icon><Calendar /></el-icon>
                <strong>{{ row.month }}</strong>
              </div>
            </template>
          </el-table-column>
          
          <el-table-column prop="amount" label="应缴金额" min-width="140" align="right">
             <template #default="{ row }">
               <span class="price-text"><span class="currency">¥</span>{{ Number(row.amount).toFixed(2) }}</span>
             </template>
          </el-table-column>
          
          <el-table-column prop="status" label="缴费状态" width="140" align="center">
             <template #default="{ row }">
                <span class="status-badge" :class="row.status === 1 ? 'is-paid' : 'is-unpaid'">
                  {{ row.status === 1 ? '已缴清' : '待缴费' }}
                </span>
             </template>
          </el-table-column>
          
          <el-table-column prop="pay_time" label="支付完成时间" min-width="200">
            <template #default="{ row }">
              <div class="time-cell" v-if="row.pay_time">
                <el-icon><Clock /></el-icon>
                <span>{{ formatDate(row.pay_time) }}</span>
              </div>
              <span v-else class="text-secondary">-</span>
            </template>
          </el-table-column>
        </el-table>

        <!-- 分页器 -->
        <div class="pagination-container" v-if="total > 0">
          <el-pagination
            v-model:current-page="page"
            v-model:page-size="size"
            :page-sizes="[10, 20, 50, 100]"
            layout="total, sizes, prev, pager, next, jumper"
            :total="total"
            @size-change="fetchData"
            @current-change="fetchData"
            class="custom-pagination"
          />
        </div>
      </div>

      <!-- 高级新增弹窗 -->
      <el-dialog
        v-model="showModal"
        title="下发新物业费账单"
        width="480px"
        class="premium-dialog"
        :close-on-click-modal="false"
      >
        <el-form label-position="top" class="premium-form">
          <el-form-item label="用户 ID" required>
             <el-input 
               v-model.number="form.user_id" 
               type="number" 
               placeholder="请输入需要下发账单的用户 ID" 
               class="custom-input"
             />
          </el-form-item>
          
          <el-form-item label="账单月份 (YYYY-MM)" required>
             <el-input 
               v-model="form.month" 
               placeholder="例如: 2024-01" 
               class="custom-input"
             />
          </el-form-item>
          
          <el-form-item label="应缴金额 (¥)" required>
             <el-input-number 
               v-model="form.amount" 
               :min="0" 
               :precision="2" 
               style="width: 100%" 
               class="custom-input-number highlight-number"
             />
          </el-form-item>
        </el-form>
        
        <template #footer>
          <div class="modal-actions">
            <button class="action-btn btn-default" @click="showModal = false">取消放弃</button>
            <button class="action-btn btn-primary" @click="saveCreate" :disabled="loading">
              {{ loading ? '正在生成...' : '确认下发' }}
            </button>
          </div>
        </template>
      </el-dialog>

    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, reactive } from 'vue'
import Navbar from '@/components/layout/Navbar.vue'
import { getAdminPropertyFeeList, createPropertyFee } from '@/api/admin'
import dayjs from 'dayjs'
import { ElMessage } from 'element-plus'
import { ArrowLeft, Plus, Clock, Calendar } from '@element-plus/icons-vue'

const list = ref([])
const total = ref(0)
const page = ref(1)
const size = ref(20)
const loadingData = ref(false)

const formatDate = (date) => dayjs(date).format('YYYY-MM-DD HH:mm')

const fetchData = async () => {
  loadingData.value = true
  try {
    const res = await getAdminPropertyFeeList({ page: page.value, size: size.value })
    list.value = res.list || []
    total.value = res.total
  } catch (error) {
    console.error(error)
  } finally {
    loadingData.value = false
  }
}

// Modal
const showModal = ref(false)
const loading = ref(false)
const form = reactive({
  user_id: '',
  month: dayjs().format('YYYY-MM'),
  amount: 200
})

const openCreateModal = () => {
  form.user_id = ''
  form.month = dayjs().format('YYYY-MM')
  form.amount = 200
  showModal.value = true
}

const saveCreate = async () => {
  if (!form.user_id || !form.month || form.amount <= 0) {
    ElMessage.warning('请填写完整的账单信息')
    return
  }
  loading.value = true
  try {
    await createPropertyFee({
      user_id: Number(form.user_id),
      month: form.month,
      amount: Number(form.amount)
    })
    ElMessage.success('物业费账单下发成功')
    showModal.value = false
    fetchData()
  } catch (error) {
    ElMessage.error('创建失败: ' + (error.response?.data?.message || '未知错误'))
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchData()
})
</script>

<style scoped>
/* 全局页面底色与容器 */
.admin-child-page { min-height: 100vh; background-color: #f8f9fa; padding-bottom: 80px; }
.custom-container { max-width: 1200px; margin: 0 auto; }

/* 极简顶部区：只保留返回和操作按钮 */
.top-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 32px 0 24px;
}

.back-btn {
  display: inline-flex; align-items: center; color: #606266; font-size: 16px; font-weight: 600;
  cursor: pointer; transition: color 0.3s; padding: 8px 16px 8px 0;
}
.back-btn:hover { color: #2d597b; }
.back-icon { margin-right: 6px; font-size: 18px; }

/* 核心表格容器 */
.table-wrapper { background: #ffffff; border-radius: 16px; padding: 24px 32px; box-shadow: 0 4px 20px rgba(0, 0, 0, 0.02); }

/* Element Plus 表格定制去后台感 */
:deep(.custom-table) { --el-table-border-color: transparent; border-radius: 8px; overflow: hidden; }
:deep(.custom-table th.el-table__cell) { font-weight: 600; font-size: 14px; padding: 18px 0; border-bottom: 1px solid #ebeef5; background: #fbfcfd; color: #606266; }
:deep(.custom-table td.el-table__cell) { padding: 20px 0; border-bottom: 1px dashed #f0f2f5; font-size: 14px; }
:deep(.custom-table::before) { display: none; }

/* 单元格排版精修 */
.user-id-text { font-family: monospace; font-size: 15px; color: #606266; background: #f4f4f5; padding: 2px 8px; border-radius: 4px; }
.icon-text-cell { display: flex; align-items: center; gap: 8px; color: #2c3e50; }
.time-cell { display: flex; align-items: center; gap: 6px; color: #909399; }
.text-secondary { color: #dcdfe6; }

/* 账单金额 */
.price-text { color: #e4393c; font-weight: 800; font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; font-size: 18px; }
.currency { font-size: 14px; margin-right: 2px; font-weight: 600; }

/* 状态标签定制 */
.status-badge { padding: 6px 14px; border-radius: 20px; font-size: 12px; font-weight: bold; display: inline-block; }
.is-unpaid { background: #fff7ed; color: #d97706; border: 1px solid #fed7aa; }
.is-paid { background: #f0fdf4; color: #166534; border: 1px solid #bbf7d0; }

/* 定制化按钮 */
.action-btn { padding: 10px 24px; border-radius: 20px; font-size: 14px; font-weight: 600; cursor: pointer; transition: all 0.3s; border: 1px solid transparent; display: inline-flex; align-items: center; justify-content: center; }
.btn-primary { background: #2d597b; color: #ffffff; box-shadow: 0 4px 12px rgba(45, 89, 123, 0.2); }
.btn-primary:hover:not(:disabled) { background: #1f435d; transform: translateY(-2px); box-shadow: 0 6px 16px rgba(45, 89, 123, 0.3); }
.btn-default { background: #ffffff; color: #606266; border-color: #dcdfe6; }
.btn-default:hover { color: #2d597b; border-color: #2d597b; background: #f0f7ff; }
.action-btn:disabled { opacity: 0.6; cursor: not-allowed; }

/* 分页器 */
.pagination-container { display: flex; justify-content: flex-end; margin-top: 32px; padding-top: 16px; }
:deep(.custom-pagination .el-pager li.is-active) { background-color: #2d597b !important; color: #fff; border-radius: 4px; }
:deep(.custom-pagination .el-pager li:hover) { color: #2d597b; }

/* ================= 弹窗高级美化 ================= */
:deep(.premium-dialog) { border-radius: 16px; box-shadow: 0 16px 48px rgba(0,0,0,0.1); overflow: hidden; }
:deep(.premium-dialog .el-dialog__header) { margin-right: 0; padding: 24px 32px 20px; border-bottom: 1px solid #f0f2f5; }
:deep(.premium-dialog .el-dialog__title) { font-weight: 700; color: #2c3e50; font-size: 18px; border-left: 4px solid #2d597b; padding-left: 10px; }
:deep(.premium-dialog .el-dialog__body) { padding: 24px 32px 12px; }
:deep(.premium-dialog .el-dialog__footer) { padding: 16px 32px 24px; border-top: 1px solid #f0f2f5; background: #fafbfc; }

/* 深度定制表单输入框 */
.premium-form :deep(.el-form-item__label) { font-weight: 600; color: #303133; padding-bottom: 6px; }
:deep(.custom-input .el-input__wrapper) {
  box-shadow: 0 0 0 1px #dcdfe6 inset; border-radius: 8px; background: #fbfcfd; transition: all 0.3s; padding: 8px 12px;
}
:deep(.custom-input .el-input__wrapper.is-focus) {
  box-shadow: 0 0 0 2px rgba(45, 89, 123, 0.2) inset !important; background: #ffffff;
}

:deep(.custom-input-number .el-input__wrapper) {
  box-shadow: 0 0 0 1px #dcdfe6 inset; border-radius: 8px; background: #fbfcfd; transition: all 0.3s;
}
:deep(.custom-input-number .el-input__wrapper.is-focus) {
  box-shadow: 0 0 0 2px rgba(45, 89, 123, 0.2) inset !important; background: #ffffff;
}
:deep(.highlight-number .el-input__inner) { color: #e4393c; font-weight: 800; font-size: 16px; }

.modal-actions { display: flex; justify-content: flex-end; gap: 12px; }
</style>