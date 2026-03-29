<template>
  <div class="admin-child-page">
    <Navbar />
    <div class="container custom-container">
      
      <!-- 新增：顶部返回导航 -->
      <div class="page-nav">
        <div class="back-btn" @click="$router.push('/admin')">
          <el-icon class="back-icon"><ArrowLeft /></el-icon> 
          <span>返回管理后台</span>
        </div>
      </div>

      <!-- 高级标题与操作区 -->
      <div class="page-header">
        
        <button class="action-btn btn-primary" @click="openModal()">
          <el-icon style="margin-right: 4px;"><Plus /></el-icon> 发布公告
        </button>
      </div>

      <!-- 深度美化的表格容器 -->
      <div class="table-wrapper">
        <el-table :data="notices" style="width: 100%" class="custom-table" :empty-text="'暂无公告记录'">
          <el-table-column prop="title" label="公告标题" min-width="180" show-overflow-tooltip>
            <template #default="{ row }">
              <strong class="text-primary">{{ row.title }}</strong>
            </template>
          </el-table-column>
          
          <el-table-column prop="content" label="内容摘要" min-width="250" show-overflow-tooltip>
            <template #default="{ row }">
              <span class="text-secondary">{{ row.content }}</span>
            </template>
          </el-table-column>
          
          <el-table-column label="发布时间" width="180">
            <template #default="{ row }">
              <div class="time-cell">
                <el-icon><Clock /></el-icon>
                <span>{{ formatDate(row.created_at) }}</span>
              </div>
            </template>
          </el-table-column>
          
          <el-table-column label="操作" width="120" fixed="right" align="center">
            <template #default="{ row }">
              <button class="action-btn btn-sm btn-danger-ghost" @click="handleDelete(row.id)">删除</button>
            </template>
          </el-table-column>
        </el-table>

        <div class="pagination-container" v-if="total > 0">
            <el-pagination
                v-model:current-page="currentPage"
                v-model:page-size="pageSize"
                :page-sizes="[10, 20, 50]"
                layout="total, sizes, prev, pager, next, jumper"
                :total="total"
                @size-change="handleSizeChange"
                @current-change="handleCurrentChange"
                class="custom-pagination"
            />
        </div>
      </div>

      <!-- 重构的高级弹窗 -->
      <el-dialog 
        v-model="showModal" 
        title="发布新公告" 
        width="520px" 
        class="premium-dialog"
        :close-on-click-modal="false"
      >
        <el-form @submit.prevent="handleSubmit" class="premium-form" label-position="top">
          <el-form-item label="公告标题" required>
            <el-input 
              v-model="form.title" 
              placeholder="请输入醒目的公告标题" 
              class="custom-input"
            />
          </el-form-item>
          <el-form-item label="公告正文" required>
            <el-input 
              v-model="form.content" 
              type="textarea" 
              :rows="5" 
              placeholder="请输入公告详细内容..." 
              class="custom-textarea"
            />
          </el-form-item>
        </el-form>
        
        <template #footer>
          <div class="modal-actions">
            <button type="button" class="action-btn btn-default" @click="closeModal">取消</button>
            <button type="button" class="action-btn btn-primary" @click="handleSubmit">确认发布</button>
          </div>
        </template>
      </el-dialog>

    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import Navbar from '@/components/layout/Navbar.vue'
import { getNoticeList } from '@/api/service'
import { createNotice, deleteNotice } from '@/api/admin'
import dayjs from 'dayjs'
import { ElMessage, ElMessageBox } from 'element-plus'
// 引入图标
import { ArrowLeft, Clock, Plus } from '@element-plus/icons-vue'

const router = useRouter()
const notices = ref([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)
const showModal = ref(false)
const form = ref({ title: '', content: '' })

const formatDate = (date) => dayjs(date).format('YYYY-MM-DD HH:mm')

const fetchNotices = async () => {
  try {
    const res = await getNoticeList({
        page: currentPage.value,
        size: pageSize.value
    })
    // 兼容处理：如果是列表（首页）则 list=res, 如果是分页则 list=res.list
    if (res.list) {
        notices.value = res.list
        total.value = res.total
    } else {
        notices.value = res
        total.value = res.length
    }
  } catch (error) {
    console.error(error)
  }
}

const handleSizeChange = (val) => {
    pageSize.value = val
    fetchNotices()
}

const handleCurrentChange = (val) => {
    currentPage.value = val
    fetchNotices()
}

const openModal = () => {
    form.value = { title: '', content: '' }
    showModal.value = true
}

const closeModal = () => showModal.value = false

const handleSubmit = async () => {
    if (!form.value.title.trim() || !form.value.content.trim()) {
        ElMessage.warning('标题和内容均不能为空')
        return
    }
    try {
        await createNotice(form.value)
        ElMessage.success('发布成功')
        closeModal()
        fetchNotices()
    } catch (e) {
        ElMessage.error('发布失败')
    }
}

const handleDelete = async (id) => {
    try {
        await ElMessageBox.confirm('删除后不可恢复，确定要删除这条公告吗?', '删除确认', {
            confirmButtonText: '确定删除',
            cancelButtonText: '暂不删除',
            confirmButtonClass: 'el-button--danger',
            type: 'warning'
        })
        await deleteNotice(id)
        ElMessage.success('删除成功')
        fetchNotices()
    } catch (e) {
        if (e !== 'cancel') {
            ElMessage.error('删除失败')
        }
    }
}

onMounted(fetchNotices)
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

/* 核心表格容器 */
.table-wrapper { background: #ffffff; border-radius: 16px; padding: 24px 32px; box-shadow: 0 4px 20px rgba(0, 0, 0, 0.02); }

/* Element Plus 表格深度定制去后台感 */
:deep(.custom-table) { --el-table-border-color: transparent; --el-table-header-bg-color: #fbfcfd; --el-table-header-text-color: #606266; border-radius: 8px; overflow: hidden; }
:deep(.custom-table th.el-table__cell) { font-weight: 600; font-size: 14px; padding: 18px 0; border-bottom: 1px solid #ebeef5; }
:deep(.custom-table td.el-table__cell) { padding: 20px 0; border-bottom: 1px dashed #f0f2f5; font-size: 14px; }
:deep(.custom-table::before) { display: none; }

/* 单元格内容排版 */
.text-primary { color: #2c3e50; font-size: 15px; }
.text-secondary { color: #606266; line-height: 1.6; }
.time-cell { display: flex; align-items: center; gap: 6px; color: #909399; }

/* 定制化按钮 */
.action-btn { padding: 10px 24px; border-radius: 20px; font-size: 14px; font-weight: 600; cursor: pointer; transition: all 0.3s; border: 1px solid transparent; display: inline-flex; align-items: center; justify-content: center; margin-left: 1000px;}
.btn-sm { padding: 6px 16px; font-size: 13px; }

.btn-primary { background: #2d597b; color: #ffffff; box-shadow: 0 4px 12px rgba(45, 89, 123, 0.2); }
.btn-primary:hover { background: #1f435d; transform: translateY(-2px); box-shadow: 0 6px 16px rgba(45, 89, 123, 0.3); }

.btn-danger-ghost { background: transparent; color: #f56c6c; border-color: #fbc4c4; }
.btn-danger-ghost:hover { background: #fef0f0; color: #e4393c; }

.btn-default { background: #ffffff; color: #606266; border-color: #dcdfe6; }
.btn-default:hover { color: #2d597b; border-color: #2d597b; background: #f0f7ff; }

/* 分页器 */
.pagination-container { display: flex; justify-content: flex-end; margin-top: 32px; padding-top: 16px; }
:deep(.custom-pagination .el-pager li.is-active) { background-color: #2d597b !important; color: #fff; border-radius: 4px; }
:deep(.custom-pagination .el-pager li:hover) { color: #2d597b; }

/* ================= 弹窗高级美化 ================= */
:deep(.premium-dialog) { border-radius: 16px; box-shadow: 0 16px 48px rgba(0,0,0,0.1); overflow: hidden; }
:deep(.premium-dialog .el-dialog__header) { margin-right: 0; padding: 24px 32px 20px; border-bottom: 1px solid #f0f2f5; }
:deep(.premium-dialog .el-dialog__title) { font-weight: 700; color: #2c3e50; font-size: 18px; }
:deep(.premium-dialog .el-dialog__body) { padding: 24px 32px; }
:deep(.premium-dialog .el-dialog__footer) { padding: 16px 32px 24px; border-top: 1px solid #f0f2f5; }

/* 深度定制表单输入框 */
.premium-form :deep(.el-form-item__label) { font-weight: 600; color: #303133; padding-bottom: 8px; }
:deep(.custom-input .el-input__wrapper),
:deep(.custom-textarea .el-textarea__inner) {
  box-shadow: 0 0 0 1px #dcdfe6 inset; border-radius: 8px; background: #fbfcfd; transition: all 0.3s;
}
:deep(.custom-input .el-input__wrapper) { padding: 8px 12px; }
:deep(.custom-textarea .el-textarea__inner) { padding: 12px; font-family: inherit; resize: none; }
:deep(.custom-input .el-input__wrapper.is-focus),
:deep(.custom-textarea .el-textarea__inner:focus) {
  box-shadow: 0 0 0 2px rgba(45, 89, 123, 0.2) inset !important; background: #ffffff;
}

.modal-actions { display: flex; justify-content: flex-end; gap: 12px; }
</style>