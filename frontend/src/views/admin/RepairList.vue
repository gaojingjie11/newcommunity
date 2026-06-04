<template>
  <div class="admin-child-page">
    <Navbar />
    <div class="container custom-container">
      <div class="top-bar">
        <div class="back-btn" @click="$router.push('/admin')">
          <el-icon class="back-icon"><ArrowLeft /></el-icon> 
          <span>返回管理后台</span>
        </div>
      </div>

      <div class="table-wrapper">
        <el-table :data="list" class="custom-table" style="width: 100%">
          <el-table-column label="提交人" min-width="120">
            <template #default="{ row }">
              {{ displayUser(row) }}
            </template>
          </el-table-column>
          
          <el-table-column label="电话" width="120">
            <template #default="{ row }">
              {{ row.user_mobile || row.user?.mobile || '--' }}
            </template>
          </el-table-column>
          
          <el-table-column label="类型" width="120">
            <template #default="{ row }">
              <span class="type-tag" :class="row.type === 'complaint' ? 'is-complaint' : 'is-repair'">
                {{ row.type === 'complaint' ? '投诉' : '报修' }}
              </span>
            </template>
          </el-table-column>

          <el-table-column label="事项类型" width="120">
            <template #default="{ row }">
              {{ row.category || '--' }}
            </template>
          </el-table-column>

          <el-table-column label="内容" min-width="200" show-overflow-tooltip>
            <template #default="{ row }">
              {{ row.description || row.content || '--' }}
            </template>
          </el-table-column>
          
          <el-table-column label="状态" width="100" align="center">
            <template #default="{ row }">
              <span class="tag" :class="getStatusClass(row.status)">
                {{ getStatusText(row.status) }}
              </span>
            </template>
          </el-table-column>
          
          <el-table-column label="提交时间" width="160">
            <template #default="{ row }">
              {{ formatDate(row.created_at) }}
            </template>
          </el-table-column>
          
          <el-table-column label="操作" width="180" fixed="right" align="center">
            <template #default="{ row }">
              <div class="row-actions">
                <button 
                  v-if="row.status !== 2" 
                  class="action-btn btn-sm btn-primary" 
                  @click="openProcess(row)"
                >
                  {{ row.status === 0 ? '开始处理' : '完成处理' }}
                </button>
                <div v-else>
                    <el-tooltip :content="row.result" placement="top" v-if="row.result">
                       <span class="text-truncate" style="display:inline-block; max-width: 150px; color: #606266;">{{ row.result }}</span>
                    </el-tooltip>
                </div>
              </div>
            </template>
          </el-table-column>
        </el-table>

        <div class="pagination-container mt-4">
            <el-pagination
                v-model:current-page="currentPage"
                v-model:page-size="pageSize"
                :page-sizes="[10, 20, 50]"
                layout="total, sizes, prev, pager, next, jumper"
                :total="total"
                @size-change="handleSizeChange"
                @current-change="handleCurrentChange"
            />
        </div>
      </div>


      <div class="modal-overlay" v-if="showModal">
        <div class="modal card">
          <h3>处理{{ currentItem?.type === 'complaint' ? '投诉' : '报修' }}</h3>
          <p class="mb-4">工单内容: {{ currentItem?.description || currentItem?.content || '--' }}</p>
          <form @submit.prevent="handleSubmit">
            <div class="form-group">
                <label>更新状态</label>
                <div class="radio-group" style="display:flex; gap:16px; margin-bottom:12px;">
                    <label style="display:inline-flex; align-items:center; cursor:pointer;">
                        <input type="radio" v-model="processForm.status" :value="1" :disabled="currentItem.status > 1"> 处理中
                    </label>
                    <label style="display:inline-flex; align-items:center; cursor:pointer;">
                        <input type="radio" v-model="processForm.status" :value="2"> 已完成
                    </label>
                </div>
            </div>

            <div class="form-group">
              <label>处理结果/反馈</label>
              <textarea v-model="processForm.result" class="input textarea" required placeholder="请输入处理结果..."></textarea>
            </div>
             <div class="modal-actions">
              <button type="button" class="btn btn-secondary" @click="closeModal">取消</button>
              <button type="submit" class="btn btn-primary">提交</button>
            </div>
          </form>
        </div>
      </div>

    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import Navbar from '@/components/layout/Navbar.vue'
import { getAdminWorkorderList, processWorkorder } from '@/api/admin'
import dayjs from 'dayjs'
import { ElMessage } from 'element-plus'

const list = ref([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)
const showModal = ref(false)
const currentItem = ref(null)
const processForm = ref({ result: '', status: 1 })

const formatDate = (date) => dayjs(date).format('YYYY-MM-DD HH:mm')
const getStatusText = (s) => {
    if(s === 0) return '待处理'
    if(s === 1) return '处理中'
    return '已完成'
}
const getStatusClass = (s) => {
    if(s === 0) return 'tag-warning'
    if(s === 1) return 'tag-primary' // processing color
    return 'tag-success'
}
const displayUser = (row) => {
    return row.user_name || row.user?.real_name || row.user?.username || `用户 #${row.user_id}`
}
const fetchList = async () => {
    try {
        const params = {
            page: currentPage.value,
            size: pageSize.value
        }
        const res = await getAdminWorkorderList(params)
        list.value = res.list || []
        total.value = res.total || 0
    } catch (e) {
        console.error(e)
        ElMessage.error(e.response?.data?.msg || e.message || '获取工单失败')
    }
}

const handleSizeChange = (val) => {
    pageSize.value = val
    fetchList()
}

const handleCurrentChange = (val) => {
    currentPage.value = val
    fetchList()
}

const openProcess = (item) => {
    currentItem.value = item
    processForm.value.result = item.result || ''
    // If pending (0), default to Processing (1). If Processing (1), default to Completed (2)
    processForm.value.status = item.status === 0 ? 1 : 2
    showModal.value = true
}
const closeModal = () => showModal.value = false

const handleSubmit = async () => {
    try {
        const payload = {
            result: processForm.value.result,
            status: processForm.value.status
        }
        await processWorkorder(currentItem.value.id, payload)
        ElMessage.success('处理成功')
        closeModal()
        fetchList()
    } catch(e) {
        ElMessage.error('提交失败: ' + (e.response?.data?.msg || e.message || '未知错误'))
    }
}

onMounted(fetchList)
</script>

<style scoped>
/* Reuse styles */
/* 全局页面底色与容器 */
.admin-child-page { min-height: 100vh; background-color: #f8f9fa; padding-bottom: 80px; }
.custom-container { max-width: 1280px; margin: 0 auto; }

/* 极简顶部区：只保留返回 */
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

.mb-4 { margin-bottom: 16px; }

.modal-overlay { position: fixed; top: 0; left: 0; right: 0; bottom: 0; background: rgba(0,0,0,0.5); display: flex; justify-content: center; align-items: center; z-index: 2000; }
.modal { padding: 24px; width: 400px; max-width: 90%; background: #fff; z-index: 2001; }
.form-group { margin-bottom: 16px; display: flex; flex-direction: column; }
.textarea { height: 100px; resize: vertical; }
.modal-actions { display: flex; justify-content: flex-end; gap: 12px; margin-top: 24px; }
.pagination-container { display: flex; justify-content: flex-end; padding-top: 20px; }
.type-tag { display: inline-block; padding: 4px 10px; border-radius: 4px; font-size: 12px; font-weight: 700; }
.type-tag.is-repair { background: #f0f7ff; color: #2d597b; border: 1px solid #cce3f6; }
.type-tag.is-complaint { background: #fdf6f6; color: #e4393c; border: 1px solid #fbc4c4; }

</style>
