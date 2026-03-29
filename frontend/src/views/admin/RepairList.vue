<template>
  <div class="admin-child-page">
    <Navbar />
    <div class="container">
 <div class="page-header">
        <div class="page-nav">
        <div class="back-btn" @click="$router.push('/admin')">
          <el-icon class="back-icon"><ArrowLeft /></el-icon> 
          <span></span>
        </div>
      </div>
        
      </div>

      <div class="table-container card" style="margin-top: 0px;">
        <el-table :data="list" style="width: 100%" stripe border>
          <el-table-column label="提交人" min-width="120">
            <template #default="{ row }">
              {{ row.user ? (row.user.real_name || row.user.username) : row.user_id }}
            </template>
          </el-table-column>
          
          <el-table-column label="电话" width="120">
            <template #default="{ row }">
              {{ row.user ? row.user.mobile : '--' }}
            </template>
          </el-table-column>
          
          <el-table-column prop="content" label="内容" min-width="200" show-overflow-tooltip />
          
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
          
          <el-table-column label="操作" width="180" fixed="right">
            <template #default="{ row }">
              <button 
                v-if="row.status !== 2" 
                class="btn btn-sm btn-primary" 
                @click="openProcess(row)"
              >
                {{ row.status === 0 ? '开始处理' : '完成处理' }}
              </button>
              <div v-else>
                  <el-tooltip :content="row.result" placement="top" v-if="row.result">
                     <span class="text-truncate" style="display:inline-block; max-width: 150px;">{{ row.result }}</span>
                  </el-tooltip>
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
          <h3>处理报修</h3>
          <p class="mb-4">报修内容: {{ currentItem?.content }}</p>
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
import { getAdminRepairList, processRepair } from '@/api/admin'
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

const fetchList = async () => {
    try {
        const res = await getAdminRepairList({
            page: currentPage.value,
            size: pageSize.value
        })
        list.value = res.list
        total.value = res.total
    } catch (e) {
        console.error(e)
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
        await processRepair({
            id: currentItem.value.id,
            feedback: processForm.value.result,
            status: processForm.value.status 
        })
        ElMessage.success('处理成功')
        closeModal()
        fetchList()
    } catch(e) {
        ElMessage.error('提交失败: ' + (e.response?.data?.msg || '未知错误'))
    }
}

onMounted(fetchList)
</script>

<style scoped>
/* Reuse styles */
.page-nav { padding: 24px 0 16px; }
.back-btn {
  display: inline-flex; align-items: center; color: #606266; font-size: 15px;
  cursor: pointer; transition: color 0.3s; padding: 8px 16px 8px 0;
}
.back-icon { margin-right: 6px; font-size: 16px; }
.admin-child-page { min-height: 100vh; padding-bottom: var(--spacing-xl); }
.page-header { display: flex; justify-content: space-before; margin-bottom: var(--spacing-lg); }
.table { width: 100%; border-collapse: collapse; }
.table th, .table td { padding: 12px; border-bottom: 1px solid #eee; text-align: left; }
.content-cell { max-width: 250px; }
.mb-4 { margin-bottom: 16px; }

.modal-overlay { position: fixed; top: 0; left: 0; right: 0; bottom: 0; background: rgba(0,0,0,0.5); display: flex; justify-content: center; align-items: center; z-index: 2000; }
.modal { padding: 24px; width: 400px; max-width: 90%; background: #fff; z-index: 2001; }
.form-group { margin-bottom: 16px; display: flex; flex-direction: column; }
.textarea { height: 100px; resize: vertical; }
.modal-actions { display: flex; justify-content: flex-end; gap: 12px; margin-top: 24px; }
.pagination-container { display: flex; justify-content: flex-end; padding-top: 20px; }

</style>
