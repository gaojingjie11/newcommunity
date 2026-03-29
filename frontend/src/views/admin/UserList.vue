<template>
  <div class="admin-child-page">
    <Navbar />
    <div class="container custom-container">
      
      <!-- 极简顶部区：只保留返回导航与搜索框，去除了大标题 -->
      <div class="top-bar">
        <div class="back-btn" @click="$router.push('/admin')">
          <el-icon class="back-icon"><ArrowLeft /></el-icon> 
          <span>返回管理后台</span>
        </div>
        <div class="header-actions">
          <el-input 
            v-model="keyword" 
            placeholder="搜索用户名 / 手机号 / 真实姓名" 
            class="search-input"
            clearable
            @keyup.enter="handleSearch"
            @clear="handleSearch"
          >
            <template #append>
              <el-button @click="handleSearch">
                <el-icon><Search /></el-icon>
              </el-button>
            </template>
          </el-input>
        </div>
      </div>

      <!-- 深度美化的表格容器 -->
      <div class="table-wrapper">
        <el-table :data="users" class="custom-table" style="width: 100%" v-loading="loadingData" :empty-text="'暂无匹配的用户数据'">
          <el-table-column prop="id" label="用户 ID" width="100" align="center">
            <template #default="{ row }">
              <span class="user-id-text">{{ row.id }}</span>
            </template>
          </el-table-column>
          
          <el-table-column prop="username" label="系统用户名" min-width="140">
            <template #default="{ row }">
              <strong class="text-primary">{{ row.username }}</strong>
            </template>
          </el-table-column>
          
          <el-table-column prop="real_name" label="真实姓名" min-width="120">
            <template #default="scope">
              <span :class="scope.row.real_name ? 'text-primary' : 'text-secondary'">
                {{ scope.row.real_name || '未填写' }}
              </span>
            </template>
          </el-table-column>
          
          <el-table-column prop="mobile" label="联系手机号" width="140" align="center">
            <template #default="scope">
              <span class="mobile-text">{{ scope.row.mobile || '-' }}</span>
            </template>
          </el-table-column>
          
          <el-table-column prop="balance" label="账户余额" width="140" align="right">
            <template #default="scope">
              <span class="price-text"><span class="currency">¥</span>{{ Number(scope.row.balance).toFixed(2) }}</span>
            </template>
          </el-table-column>
          
          <el-table-column prop="role" label="系统角色" width="130" align="center">
             <template #default="scope">
               <span class="role-badge" :class="`role-${scope.row.role}`">
                 {{ getRoleText(scope.row.role) }}
               </span>
             </template>
          </el-table-column>
          
          <el-table-column prop="status" label="账号状态" width="100" align="center">
            <template #default="scope">
              <span class="status-badge" :class="scope.row.status === 1 ? 'is-active' : 'is-frozen'">
                {{ scope.row.status === 1 ? '正常' : '已冻结' }}
              </span>
            </template>
          </el-table-column>
          
          <el-table-column prop="created_at" label="注册时间" min-width="180">
             <template #default="scope">
               <div class="time-cell">
                 <el-icon><Clock /></el-icon>
                 <span>{{ formatDate(scope.row.created_at) }}</span>
               </div>
             </template>
          </el-table-column>
          
          <el-table-column label="管理操作" width="160" fixed="right" align="center">
            <template #default="scope">
              <div class="row-actions">
                <button class="action-btn btn-sm btn-outline" @click="openEditModal(scope.row)">编辑</button>
                <el-popconfirm
                  :title="`确定要${scope.row.status === 1 ? '冻结' : '解冻'}该用户账号吗？`"
                  :confirm-button-text="scope.row.status === 1 ? '确认冻结' : '确认解冻'"
                  cancel-button-text="取消"
                  :confirm-button-type="scope.row.status === 1 ? 'danger' : 'success'"
                  @confirm="handleFreeze(scope.row)"
                >
                  <template #reference>
                    <button 
                      class="action-btn btn-sm" 
                      :class="scope.row.status === 1 ? 'btn-danger-ghost' : 'btn-success-ghost'"
                    >
                      {{ scope.row.status === 1 ? '冻结' : '解冻' }}
                    </button>
                  </template>
                </el-popconfirm>
              </div>
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
            @size-change="handleSearch"
            @current-change="fetchUsers"
            class="custom-pagination"
          />
        </div>
      </div>

      <!-- 高级编辑用户弹窗 -->
      <el-dialog
        v-model="showModal"
        title="配置用户权限与资产"
        width="520px"
        class="premium-dialog"
        :close-on-click-modal="false"
      >
        <el-form label-position="top" class="premium-form">
          <div class="current-user-panel">
            <div class="user-avatar-placeholder">
              <el-icon><UserFilled /></el-icon>
            </div>
            <div class="user-meta">
              <span class="meta-name">{{ currentUser.username }}</span>
              <span class="meta-balance">当前余额：¥{{ Number(currentUser.balance || 0).toFixed(2) }}</span>
            </div>
          </div>

          <el-form-item label="系统角色分配" required>
             <el-select v-model="editForm.role" class="custom-select" style="width: 100%">
                <el-option label="普通居民用户" value="user" />
                <el-option label="社区物业人员" value="property" />
                <el-option label="入驻商场店主" value="store" />
                <el-option label="系统超级管理员" value="admin" />
             </el-select>
          </el-form-item>
          
          <el-form-item label="余额人工干预 (¥)">
             <el-input-number 
               v-model="editForm.balanceAmount" 
               :precision="2" 
               :step="10" 
               style="width: 100%" 
               class="custom-input-number highlight-number"
             />
             <div class="form-tip">
               <el-icon><InfoFilled /></el-icon> 填入正数代表充值加款，填入负数代表扣除余额。若不修改请保持为 0。
             </div>
          </el-form-item>
        </el-form>
        
        <template #footer>
          <div class="modal-actions">
            <button class="action-btn btn-default" @click="showModal = false">取消放弃</button>
            <button class="action-btn btn-primary" @click="saveEdit" :disabled="loading">
              {{ loading ? '保存中...' : '确认保存更改' }}
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
import { getUserList, freezeUser, assignRole, updateUserBalance } from '@/api/admin'
import dayjs from 'dayjs'
import { ElMessage } from 'element-plus'
// 引入图标
import { ArrowLeft, Search, Clock, UserFilled, InfoFilled } from '@element-plus/icons-vue'

const users = ref([])
const keyword = ref('')
const total = ref(0) 
const page = ref(1)
const size = ref(20)
const loadingData = ref(false)

const formatDate = (date) => dayjs(date).format('YYYY-MM-DD HH:mm')

const getRoleText = (role) => {
  const map = { 'admin': '管理员', 'user': '普通用户', 'property': '物业管理', 'store': '合作商户' }
  return map[role] || role
}

const handleSearch = () => {
    page.value = 1
    fetchUsers()
}

const fetchUsers = async () => {
  loadingData.value = true
  try {
    const res = await getUserList({ 
      page: page.value, 
      size: size.value, 
      keyword: keyword.value 
    })
    users.value = res.list || []
    total.value = res.total || 0
  } catch (error) {
    console.error('获取用户列表失败', error)
  } finally {
    loadingData.value = false
  }
}

const handleFreeze = async (user) => {
  try {
    const newStatus = user.status === 1 ? 0 : 1
    await freezeUser({ id: user.id, status: newStatus })
    ElMessage.success(`账号已成功${newStatus === 1 ? '解冻' : '冻结'}`)
    await fetchUsers()
  } catch (error) {
    ElMessage.error('操作失败')
  }
}

// Modal Logic
const showModal = ref(false)
const loading = ref(false)
const currentUser = ref({})
const editForm = reactive({
  role: 'user',
  balanceAmount: 0
})

const openEditModal = (user) => {
  currentUser.value = user
  editForm.role = user.role
  editForm.balanceAmount = 0
  showModal.value = true
}

const saveEdit = async () => {
  loading.value = true
  try {
    // 1. Update Role if changed
    if (editForm.role !== currentUser.value.role) {
      await assignRole({ user_id: currentUser.value.id, role_code: editForm.role })
    }
    // 2. Update Balance if entered
    if (editForm.balanceAmount !== 0) {
      await updateUserBalance({ user_id: currentUser.value.id, amount: editForm.balanceAmount })
    }
    
    ElMessage.success('用户配置修改成功')
    showModal.value = false
    fetchUsers()
  } catch (error) {
    console.error(error)
    ElMessage.error('操作失败: ' + (error.response?.data?.message || '未知错误'))
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchUsers()
})
</script>

<style scoped>
/* 全局页面底色与容器 */
.admin-child-page { min-height: 100vh; background-color: #f8f9fa; padding-bottom: 80px; }
.custom-container { max-width: 1280px; margin: 0 auto; }

/* 极简顶部区：只保留返回和搜索框 */
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

/* 搜索框深度美化 */
.header-actions { width: 360px; }
:deep(.search-input .el-input__wrapper) {
  border-radius: 20px 0 0 20px; box-shadow: 0 2px 10px rgba(0,0,0,0.02); background: #ffffff;
}
:deep(.search-input .el-input-group__append) {
  border-radius: 0 20px 20px 0; background: #2d597b; color: white; border: none; box-shadow: 0 2px 10px rgba(45,89,123,0.1);
}

/* 核心表格容器 */
.table-wrapper { background: #ffffff; border-radius: 16px; padding: 24px 32px; box-shadow: 0 4px 20px rgba(0, 0, 0, 0.02); }

/* Element Plus 表格定制去后台感 */
:deep(.custom-table) { --el-table-border-color: transparent; border-radius: 8px; overflow: hidden; }
:deep(.custom-table th.el-table__cell) { font-weight: 600; font-size: 14px; padding: 18px 0; border-bottom: 1px solid #ebeef5; background: #fbfcfd; color: #606266; }
:deep(.custom-table td.el-table__cell) { padding: 20px 0; border-bottom: 1px dashed #f0f2f5; font-size: 14px; }
:deep(.custom-table::before) { display: none; }

/* 单元格排版精修 */
.text-primary { color: #2c3e50; font-size: 15px; }
.text-secondary { color: #a4b0be; }
.user-id-text { font-family: monospace; font-size: 15px; color: #606266; background: #f4f4f5; padding: 2px 8px; border-radius: 4px; }
.mobile-text { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; font-size: 14px; color: #606266; }
.time-cell { display: flex; align-items: center; gap: 6px; color: #909399; }

/* 余额排版 */
.price-text { color: #e4393c; font-weight: 800; font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; font-size: 16px; }
.currency { font-size: 13px; margin-right: 2px; font-weight: 600; }

/* 角色标签定制 */
.role-badge { padding: 4px 10px; border-radius: 6px; font-size: 12px; font-weight: 600; display: inline-block; }
.role-admin { background: #f0f7ff; color: #2d597b; border: 1px solid #cce3f6; }
.role-store { background: #fff7ed; color: #d97706; border: 1px solid #fed7aa; }
.role-property { background: #f4f4f5; color: #606266; border: 1px solid #dcdfe6; }
.role-user { background: #f0fdf4; color: #166534; border: 1px solid #bbf7d0; }

/* 状态标签定制 */
.status-badge { padding: 6px 14px; border-radius: 20px; font-size: 12px; font-weight: bold; display: inline-block; }
.is-active { background: #f0fdf4; color: #166534; }
.is-frozen { background: #fef0f0; color: #e4393c; }

/* 列表操作区 */
.row-actions { display: flex; gap: 8px; justify-content: center; }

/* 定制化按钮 */
.action-btn { padding: 10px 24px; border-radius: 20px; font-size: 14px; font-weight: 600; cursor: pointer; transition: all 0.3s; border: 1px solid transparent; display: inline-flex; align-items: center; justify-content: center; }
.btn-sm { padding: 6px 16px; font-size: 13px; }

.btn-primary { background: #2d597b; color: #ffffff; box-shadow: 0 4px 12px rgba(45, 89, 123, 0.2); }
.btn-primary:hover:not(:disabled) { background: #1f435d; transform: translateY(-2px); box-shadow: 0 6px 16px rgba(45, 89, 123, 0.3); }

.btn-outline { background: #ffffff; color: #2d597b; border-color: #2d597b; }
.btn-outline:hover { background: #f0f7ff; transform: translateY(-1px); }

.btn-danger-ghost { background: transparent; color: #f56c6c; border-color: #fbc4c4; }
.btn-danger-ghost:hover { background: #fef0f0; color: #e4393c; transform: translateY(-1px); }

.btn-success-ghost { background: transparent; color: #00b894; border-color: #a7e9d9; }
.btn-success-ghost:hover { background: #f0fdf4; color: #166534; transform: translateY(-1px); }

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

/* 弹窗内的用户信息面板 */
.current-user-panel {
  display: flex;
  align-items: center;
  gap: 16px;
  background: #f8f9fa;
  padding: 16px 20px;
  border-radius: 12px;
  border: 1px solid #ebeef5;
  margin-bottom: 24px;
}

.user-avatar-placeholder {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  background: #e1f0ff;
  color: #2d597b;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
}

.user-meta {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.meta-name {
  font-size: 16px;
  font-weight: 700;
  color: #2c3e50;
}

.meta-balance {
  font-size: 13px;
  color: #606266;
}

/* 深度定制表单输入框 */
.premium-form :deep(.el-form-item__label) { font-weight: 600; color: #303133; padding-bottom: 6px; }

:deep(.custom-select .el-select__wrapper) {
  box-shadow: 0 0 0 1px #dcdfe6 inset; border-radius: 8px; background: #fbfcfd; transition: all 0.3s; padding: 6px 12px;
}
:deep(.custom-select .el-select__wrapper.is-focus) {
  box-shadow: 0 0 0 2px rgba(45, 89, 123, 0.2) inset !important; background: #ffffff;
}

:deep(.custom-input-number .el-input__wrapper) {
  box-shadow: 0 0 0 1px #dcdfe6 inset; border-radius: 8px; background: #fbfcfd; transition: all 0.3s;
}
:deep(.custom-input-number .el-input__wrapper.is-focus) {
  box-shadow: 0 0 0 2px rgba(45, 89, 123, 0.2) inset !important; background: #ffffff;
}
:deep(.highlight-number .el-input__inner) { color: #e4393c; font-weight: 800; font-size: 16px; }

.form-tip {
  margin-top: 8px;
  font-size: 13px;
  color: #909399;
  display: flex;
  align-items: center;
  gap: 6px;
}

.modal-actions { display: flex; justify-content: flex-end; gap: 12px; }

@media (max-width: 768px) {
  .top-bar { flex-direction: column; align-items: flex-start; gap: 16px; }
  .header-actions { width: 100%; }
}
</style>