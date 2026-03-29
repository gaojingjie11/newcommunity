<template>
  <div class="admin-child-page">
    <Navbar />
    <div class="container custom-container">
      
      <!-- 极简顶部区：只保留返回导航，去除了大标题 -->
      <div class="top-bar">
        <div class="back-btn" @click="$router.push('/admin')">
          <el-icon class="back-icon"><ArrowLeft /></el-icon> 
          <span>返回管理后台</span>
        </div>
      </div>

      <!-- 深度美化的表格容器 -->
      <div class="table-wrapper">
        <el-table :data="visitors" class="custom-table" style="width: 100%" v-loading="loadingData" :empty-text="'暂无待处理的访客记录'">
          <el-table-column prop="name" label="访客姓名" width="140" align="center">
            <template #default="{ row }">
              <strong class="text-primary">{{ row.name }}</strong>
            </template>
          </el-table-column>
          
          <el-table-column prop="mobile" label="联系电话" width="160" align="center">
            <template #default="{ row }">
              <span class="mobile-text">{{ row.mobile }}</span>
            </template>
          </el-table-column>
          
          <el-table-column prop="reason" label="来访事由" min-width="200" show-overflow-tooltip>
            <template #default="{ row }">
              <span class="text-secondary">{{ row.reason }}</span>
            </template>
          </el-table-column>
          
          <el-table-column prop="visit_time" label="预计来访时间" width="200">
              <template #default="{ row }">
                <div class="time-cell">
                  <el-icon><Clock /></el-icon>
                  <span>{{ formatDate(row.visit_time) }}</span>
                </div>
              </template>
          </el-table-column>
          
          <el-table-column prop="created_at" label="提交申请时间" width="200">
              <template #default="{ row }">
                <div class="time-cell">
                  <el-icon><Clock /></el-icon>
                  <span>{{ formatDate(row.created_at) }}</span>
                </div>
              </template>
          </el-table-column>
          
          <el-table-column prop="status" label="审核状态" width="120" align="center">
             <template #default="{ row }">
               <span class="status-badge" :class="getStatusClass(row.status)">
                 {{ getStatusText(row.status) }}
               </span>
             </template>
          </el-table-column>
          
          <el-table-column label="审核操作" width="180" fixed="right" align="center">
             <template #default="{ row }">
               <div class="row-actions" v-if="row.status === 0">
                 <button class="action-btn btn-sm btn-success-ghost" @click="handleAudit(row, 1)">同意</button>
                 <button class="action-btn btn-sm btn-danger-ghost" @click="handleAudit(row, 2)">拒绝</button>
               </div>
               <span v-else class="remark-text" :title="row.audit_remark">
                 {{ row.audit_remark || '系统已处理' }}
               </span>
             </template>
          </el-table-column>
        </el-table>

        <!-- 分页器 -->
        <div class="pagination-container" v-if="total > 0">
          <el-pagination
            v-model:current-page="page"
            v-model:page-size="size"
            :page-sizes="[10, 20, 50]"
            layout="total, sizes, prev, pager, next, jumper"
            :total="total"
            @size-change="fetchList"
            @current-change="fetchList"
            class="custom-pagination"
          />
        </div>
      </div>

    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import Navbar from '@/components/layout/Navbar.vue'
import { getAdminVisitorList, auditVisitor } from '@/api/admin'
import dayjs from 'dayjs'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ArrowLeft, Clock } from '@element-plus/icons-vue'

const visitors = ref([])
const loadingData = ref(false)
const total = ref(0)
const page = ref(1)
const size = ref(10)

const formatDate = (date) => dayjs(date).format('YYYY-MM-DD HH:mm')

const getStatusText = (s) => {
    const map = { 0: '待审核', 1: '已通过', 2: '已驳回' }
    return map[s] || s
}

const getStatusClass = (s) => {
    const map = { 0: 'is-pending', 1: 'is-pass', 2: 'is-reject' }
    return map[s] || 'is-pending'
}

const fetchList = async () => {
  loadingData.value = true
  try {
    const res = await getAdminVisitorList({ page: page.value, size: size.value })
    if (res.list) {
        visitors.value = res.list
        total.value = res.total
    } else {
        visitors.value = res // Fallback
        total.value = res.length
    }
  } catch (error) {
    console.error(error)
  } finally {
      loadingData.value = false
  }
}

const handleAudit = async (item, status) => {
    if (status === 1) {
        // Approve
         try {
             await auditVisitor({
                id: item.id,
                status: 1,
                audit_remark: '同意通行'
            })
            ElMessage.success('访客申请已通过')
            fetchList()
         } catch(e) {
             ElMessage.error('操作失败')
         }
    } else {
        // Reject
        ElMessageBox.prompt('请输入拒绝该访客通行的原因', '拒绝审核', {
            confirmButtonText: '确认驳回',
            cancelButtonText: '取消',
            inputPattern: /.+/,
            inputErrorMessage: '原因不能为空',
            customClass: 'premium-msg-box' // 方便如果你之后全局覆盖弹窗样式
        }).then(async ({ value }) => {
            try {
                await auditVisitor({
                    id: item.id,
                    status: 2,
                    audit_remark: value
                })
                ElMessage.success('已驳回该访客申请')
                fetchList()
            } catch (e) {
                ElMessage.error('操作失败')
            }
        }).catch(() => {
            // Cancelled
        })
    }
}

onMounted(() => {
    fetchList()
})
</script>

<style scoped>
/* 全局页面底色与容器 */
.admin-child-page { min-height: 100vh; background-color: #f8f9fa; padding-bottom: 80px; }
.custom-container { max-width: 1280px; margin: 0 auto; }

/* 极简顶部区：只保留返回导航 */
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
.text-primary { color: #2c3e50; font-size: 15px; }
.text-secondary { color: #606266; line-height: 1.6; }
.mobile-text { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; font-size: 15px; color: #2d597b; font-weight: 600; }
.time-cell { display: flex; align-items: center; justify-content: center; gap: 6px; color: #909399; }
.remark-text { font-size: 13px; color: #a4b0be; background: #f4f4f5; padding: 4px 12px; border-radius: 6px; display: inline-block; max-width: 120px; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }

/* 状态标签定制 */
.status-badge { padding: 6px 14px; border-radius: 20px; font-size: 12px; font-weight: bold; display: inline-block; }
.is-pending { background: #fff7ed; color: #d97706; border: 1px solid #fed7aa; }
.is-pass { background: #f0fdf4; color: #166534; border: 1px solid #bbf7d0; }
.is-reject { background: #fef0f0; color: #e4393c; border: 1px solid #fbc4c4; }

/* 列表操作区 */
.row-actions { display: flex; gap: 12px; justify-content: center; }

/* 定制化按钮 */
.action-btn { padding: 8px 16px; border-radius: 20px; font-size: 13px; font-weight: 600; cursor: pointer; transition: all 0.3s; border: 1px solid transparent; display: inline-flex; align-items: center; justify-content: center; }
.btn-sm { padding: 6px 16px; font-size: 13px; }

.btn-success-ghost { background: transparent; color: #00b894; border-color: #a7e9d9; }
.btn-success-ghost:hover { background: #f0fdf4; color: #166534; transform: translateY(-1px); }

.btn-danger-ghost { background: transparent; color: #f56c6c; border-color: #fbc4c4; }
.btn-danger-ghost:hover { background: #fef0f0; color: #e4393c; transform: translateY(-1px); }

/* 分页器 */
.pagination-container { display: flex; justify-content: flex-end; margin-top: 32px; padding-top: 16px; }
:deep(.custom-pagination .el-pager li.is-active) { background-color: #2d597b !important; color: #fff; border-radius: 4px; }
:deep(.custom-pagination .el-pager li:hover) { color: #2d597b; }

@media (max-width: 768px) {
  .top-bar { flex-direction: column; align-items: flex-start; gap: 16px; }
}
</style>