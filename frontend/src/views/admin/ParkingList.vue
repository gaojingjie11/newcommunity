<template>
  <div class="admin-child-page">
    <Navbar />
    <div class="container custom-container">
      
      <!-- 顶部返回导航 -->
      <div class="page-nav">
        <div class="back-btn" @click="$router.push('/admin')">
          <el-icon class="back-icon"><ArrowLeft /></el-icon> 
          <span>返回管理后台</span>
        </div>
      </div>

      <!-- 页面标题与操作 -->
      <div class="page-header">
        <h1 class="page-title highlight-title">车位管理</h1>
        <div class="actions">
          <button class="action-btn btn-primary" @click="openCreateModal">
            <el-icon style="margin-right: 6px;"><Plus /></el-icon>新增车位
          </button>
        </div>
      </div>

      <!-- 高级指标卡片 -->
      <div class="stats-grid">
        <div class="metric-card">
          <div class="metric-icon-wrap icon-total">
            <el-icon><Van /></el-icon>
          </div>
          <div class="metric-content">
            <span class="metric-label">总车位数量</span>
            <strong class="metric-value">{{ stats.total || 0 }}</strong>
          </div>
        </div>
        <div class="metric-card">
          <div class="metric-icon-wrap icon-used">
            <el-icon><Warning /></el-icon>
          </div>
          <div class="metric-content">
            <span class="metric-label">已占用车位</span>
            <strong class="metric-value text-warning">{{ stats.used || 0 }}</strong>
          </div>
        </div>
        <div class="metric-card">
          <div class="metric-icon-wrap icon-free">
            <el-icon><CircleCheck /></el-icon>
          </div>
          <div class="metric-content">
            <span class="metric-label">当前空闲车位</span>
            <strong class="metric-value text-success">{{ stats.free || 0 }}</strong>
          </div>
        </div>
      </div>

      <!-- 深度定制数据表格 -->
      <div class="table-wrapper">
         <el-table :data="list" class="custom-table" style="width: 100%" v-loading="loadingData" empty-text="暂无车位数据">
          <el-table-column prop="id" label="ID" width="80" align="center" />
          
          <el-table-column prop="parking_no" label="车位编号" min-width="120">
            <template #default="{ row }">
              <span class="parking-no">{{ row.parking_no }}</span>
            </template>
          </el-table-column>

          <el-table-column prop="status" label="使用状态" width="120" align="center">
             <template #default="scope">
                <span class="status-badge" :class="scope.row.status === 1 ? 'is-occupied' : 'is-free'">
                  {{ scope.row.status === 1 ? '已占用' : '空闲' }}
                </span>
             </template>
          </el-table-column>
          
          <el-table-column prop="user_id" label="绑定用户ID" width="120" align="center">
             <template #default="scope">
               <span class="text-secondary">{{ scope.row.user_id || '-' }}</span>
             </template>
          </el-table-column>
          
          <el-table-column prop="car_plate" label="绑定车牌" min-width="150">
             <template #default="scope">
               <span v-if="scope.row.car_plate" class="plate-badge">{{ scope.row.car_plate }}</span>
               <span v-else class="text-secondary">-</span>
             </template>
          </el-table-column>
          
          <el-table-column label="操作" width="160" fixed="right" align="center">
             <template #default="scope">
               <button 
                 class="action-btn btn-sm" 
                 :class="scope.row.status === 1 ? 'btn-outline' : 'btn-primary'" 
                 @click="openAssignModal(scope.row)"
               >
                  {{ scope.row.status === 1 ? '修改 / 解绑' : '分配车位' }}
               </button>
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
      
      <!-- 分配/解绑 弹窗 -->
      <el-dialog
        v-model="showModal"
        :title="`分配车位: ${currentItem.parking_no}`"
        width="450px"
        class="premium-dialog"
        :close-on-click-modal="false"
      >
        <div class="premium-alert is-info">
          <el-icon class="alert-icon"><InfoFilled /></el-icon>
          <div class="alert-content">请输入用户 ID 和车牌号进行绑定。<strong>输入 ID 为 0 则解除绑定。</strong></div>
        </div>
        
        <el-form label-position="top" class="premium-form">
          <el-form-item label="用户 ID (User ID)">
            <el-input v-model.number="form.user_id" type="number" placeholder="输入绑定的用户 ID" class="custom-input" />
          </el-form-item>
          <el-form-item label="车牌号 (Car Plate)">
            <el-input v-model="form.car_plate" placeholder="例如: 辽A66666" class="custom-input" />
          </el-form-item>
        </el-form>
        
        <template #footer>
          <div class="modal-actions">
            <button class="action-btn btn-default" @click="showModal = false">取消</button>
            <button class="action-btn btn-primary" :disabled="loading" @click="saveAssign">
              {{ loading ? '保存中...' : '确认保存' }}
            </button>
          </div>
        </template>
      </el-dialog>

      <!-- 新增车位 弹窗 -->
      <el-dialog 
        v-model="showCreateModal" 
        title="新增车位" 
        width="450px"
        class="premium-dialog"
        :close-on-click-modal="false"
      >
        <el-form :model="createForm" label-position="top" class="premium-form">
          <el-form-item label="车位编号" required>
            <el-input v-model="createForm.parking_no" placeholder="例如: A-001" class="custom-input" />
          </el-form-item>
        </el-form>
        <template #footer>
          <div class="modal-actions">
            <button class="action-btn btn-default" @click="showCreateModal = false">取消</button>
            <button class="action-btn btn-primary" :disabled="creating" @click="saveCreate">
              {{ creating ? '创建中...' : '确认创建' }}
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
import { getAdminParkingList, getParkingStats, assignParking, createParking } from '@/api/admin'
import { ElMessage } from 'element-plus'
import { ArrowLeft, Plus, Van, Warning, CircleCheck, InfoFilled } from '@element-plus/icons-vue'

const list = ref([])
const stats = ref({})
const total = ref(0)
const page = ref(1)
const size = ref(20)
const loadingData = ref(false)

const fetchData = async () => {
  loadingData.value = true
  try {
    const res = await getAdminParkingList({ page: page.value, size: size.value })
    list.value = res.list || []
    total.value = res.total
    
    // fetch stats
    const statsRes = await getParkingStats()
    stats.value = statsRes || {}
  } catch (error) {
    console.error(error)
  } finally {
    loadingData.value = false
  }
}

// Modal
const showModal = ref(false)
const currentItem = ref({})
const form = reactive({
  user_id: '',
  car_plate: ''
})
const loading = ref(false)

const openAssignModal = (item) => {
  currentItem.value = item
  form.user_id = item.user_id || ''
  form.car_plate = item.car_plate || ''
  showModal.value = true
}

const saveAssign = async () => {
  loading.value = true
  try {
    await assignParking({
      id: currentItem.value.id,
      user_id: Number(form.user_id),
      car_plate: form.car_plate
    })
    ElMessage.success('分配/解绑操作成功')
    showModal.value = false
    fetchData()
  } catch (error) {
    ElMessage.error('操作失败: ' + (error.response?.data?.message || '未知错误'))
  } finally {
    loading.value = false
  }
}

// Create Logic
const showCreateModal = ref(false)
const createForm = reactive({ parking_no: '' })
const creating = ref(false)

const openCreateModal = () => {
    createForm.parking_no = ''
    showCreateModal.value = true
}

const saveCreate = async () => {
    if (!createForm.parking_no) {
        ElMessage.warning('请输入车位号')
        return
    }
    creating.value = true
    try {
        await createParking({ parking_no: createForm.parking_no })
        ElMessage.success('新增车位成功')
        showCreateModal.value = false
        fetchData()
    } catch (e) {
        ElMessage.error(e.response?.data?.msg || '创建失败')
    } finally {
        creating.value = false
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

/* ★ 指标卡片重构 ★ */
.stats-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 24px;
  margin-bottom: 32px;
}

.metric-card {
  background: #ffffff;
  border-radius: 16px;
  padding: 24px 32px;
  display: flex;
  align-items: center;
  gap: 20px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.02);
  transition: all 0.4s ease;
  border: 1px solid transparent;
}

.metric-card:hover {
  transform: translateY(-6px);
  box-shadow: 0 12px 24px rgba(45, 89, 123, 0.08);
  border-color: rgba(45, 89, 123, 0.1);
}

.metric-icon-wrap {
  width: 56px;
  height: 56px;
  border-radius: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 28px;
}

.icon-total { background: #f0f7ff; color: #409eff; }
.icon-used { background: #fff7ed; color: #e6a23c; }
.icon-free { background: #f0fdf4; color: #00b894; }

.metric-content { display: flex; flex-direction: column; gap: 4px; }
.metric-label { color: #909399; font-size: 14px; font-weight: 500; }
.metric-value { font-size: 32px; font-weight: 800; color: #2c3e50; font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; letter-spacing: -0.5px; }
.text-warning { color: #d97706; }
.text-success { color: #166534; }

/* 核心表格容器 */
.table-wrapper { background: #ffffff; border-radius: 16px; padding: 24px 32px; box-shadow: 0 4px 20px rgba(0, 0, 0, 0.02); }

/* Element Plus 表格定制 */
:deep(.custom-table) { --el-table-border-color: transparent; border-radius: 8px; overflow: hidden; }
:deep(.custom-table th.el-table__cell) { font-weight: 600; font-size: 14px; padding: 18px 0; border-bottom: 1px solid #ebeef5; background: #fbfcfd; }
:deep(.custom-table td.el-table__cell) { padding: 20px 0; border-bottom: 1px dashed #f0f2f5; }
:deep(.custom-table::before) { display: none; }

/* 单元格排版 */
.parking-no { font-weight: 800; font-size: 16px; color: #2d597b; font-family: monospace; }
.text-secondary { color: #a4b0be; }

/* 状态标签 */
.status-badge { padding: 4px 12px; border-radius: 20px; font-size: 12px; font-weight: bold; }
.is-occupied { background: #fff7ed; color: #d97706; }
.is-free { background: #f0fdf4; color: #166534; }

/* 蓝底白字仿真车牌 */
.plate-badge {
  display: inline-block;
  background: #1e40af; color: #ffffff; padding: 4px 12px; border-radius: 4px;
  font-size: 13px; font-weight: bold; letter-spacing: 1px; box-shadow: inset 0 0 0 1px rgba(255,255,255,0.3);
  border: 1px solid #1e3a8a;
}

/* 定制按钮 */
.action-btn { padding: 10px 24px; border-radius: 20px; font-size: 13px; font-weight: 600; cursor: pointer; transition: all 0.3s; border: 1px solid transparent; display: inline-flex; align-items: center; justify-content: center; }
.btn-sm { padding: 6px 16px; font-size: 13px; }

.btn-primary { background: #2d597b; color: #ffffff; box-shadow: 0 4px 12px rgba(45, 89, 123, 0.2); }
.btn-primary:hover:not(:disabled) { background: #1f435d; transform: translateY(-1px); box-shadow: 0 6px 16px rgba(45, 89, 123, 0.3); }
.btn-outline { background: #ffffff; color: #2d597b; border-color: #2d597b; }
.btn-outline:hover:not(:disabled) { background: #f0f7ff; transform: translateY(-1px); }
.btn-default { background: #ffffff; color: #606266; border-color: #dcdfe6; }
.btn-default:hover { color: #2d597b; border-color: #2d597b; background: #f0f7ff; }

.action-btn:disabled { opacity: 0.6; cursor: not-allowed; }

/* 分页器 */
.pagination-container { display: flex; justify-content: flex-end; margin-top: 32px; padding-top: 16px; border-top: 1px solid #f0f2f5; }
:deep(.custom-pagination .el-pager li.is-active) { background-color: #2d597b !important; color: #fff; border-radius: 4px; }
:deep(.custom-pagination .el-pager li:hover) { color: #2d597b; }

/* ================= 弹窗高级美化 ================= */
:deep(.premium-dialog) { border-radius: 16px; box-shadow: 0 16px 48px rgba(0,0,0,0.1); overflow: hidden; }
:deep(.premium-dialog .el-dialog__header) { margin-right: 0; padding: 24px 32px 20px; border-bottom: 1px solid #f0f2f5; }
:deep(.premium-dialog .el-dialog__title) { font-weight: 700; color: #2c3e50; font-size: 18px; }
:deep(.premium-dialog .el-dialog__body) { padding: 24px 32px; }
:deep(.premium-dialog .el-dialog__footer) { padding: 16px 32px 24px; border-top: 1px solid #f0f2f5; }

.premium-alert { display: flex; align-items: flex-start; border-radius: 8px; padding: 12px 16px; font-size: 13px; line-height: 1.5; margin-bottom: 24px; }
.premium-alert.is-info { background: #f0f7ff; border: 1px solid #cce3f6; color: #2d597b; }
.alert-icon { font-size: 16px; margin-right: 8px; margin-top: 2px; }

/* 深度定制表单输入框 */
.premium-form :deep(.el-form-item__label) { font-weight: 600; color: #303133; padding-bottom: 8px; }
:deep(.custom-input .el-input__wrapper) {
  box-shadow: 0 0 0 1px #dcdfe6 inset; border-radius: 8px; background: #fbfcfd; transition: all 0.3s; padding: 8px 12px;
}
:deep(.custom-input .el-input__wrapper.is-focus) {
  box-shadow: 0 0 0 2px rgba(45, 89, 123, 0.2) inset !important; background: #ffffff;
}

.modal-actions { display: flex; justify-content: flex-end; gap: 12px; }

@media (max-width: 768px) {
  .stats-grid { grid-template-columns: 1fr; }
}
</style>