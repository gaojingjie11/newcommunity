<template>
  <div class="visitor-page">
    <Navbar />
    
    <div class="container custom-container">
      <div class="page-header">
        <h1 class="page-title highlight-title">访客登记</h1>
      </div>
      
      <div class="visitor-form-wrapper">
        <div class="form-header">
          <h3>预约新访客</h3>
          <p>请准确填写信息，方便门岗快速放行</p>
        </div>
        <form @submit.prevent="submitVisitor" class="premium-form">
          <div class="form-row">
            <div class="form-group half">
              <label>访客姓名</label>
              <input v-model="form.name" class="custom-input" placeholder="请输入真实姓名" required />
            </div>
            <div class="form-group half">
              <label>联系电话</label>
              <input v-model="form.mobile" type="tel" class="custom-input" placeholder="用于接收通行码" required />
            </div>
          </div>
          
          <div class="form-row">
            <div class="form-group half">
              <label>来访原因</label>
              <input v-model="form.reason" class="custom-input" placeholder="如: 探亲访友、维修、快递等" required />
            </div>
            <div class="form-group half">
              <label>预计抵达时间</label>
              <input v-model="form.visit_time" type="datetime-local" class="custom-input" required />
            </div>
          </div>
          
          <div class="form-actions">
            <button type="submit" class="btn-submit" :disabled="loading">
              {{ loading ? '正在提交...' : '确认预约' }}
            </button>
          </div>
        </form>
      </div>
      
      <div class="records-section">
        <div class="section-title-wrap">
          <span class="indicator"></span>
          <h3>我的预约记录</h3>
        </div>

        <div class="visitor-list" v-if="visitors.length > 0">
          <div class="visitor-card" v-for="visitor in visitors" :key="visitor.id">
            <div class="card-top">
              <div class="visitor-name">{{ visitor.name }}</div>
              <span class="status-badge" :class="getStatusClass(visitor.status)">
                {{ getStatusText(visitor.status) }}
              </span>
            </div>
            <div class="card-body">
              <div class="info-line">
                <span class="label">电话</span><span class="value">{{ visitor.mobile }}</span>
              </div>
              <div class="info-line">
                <span class="label">原因</span><span class="value">{{ visitor.reason }}</span>
              </div>
              <div class="info-line">
                <span class="label">时间</span><span class="value">{{ formatDate(visitor.visit_time) }}</span>
              </div>
            </div>
            <div class="card-footer" v-if="visitor.audit_remark">
              <span class="remark-label">审核回复：</span>{{ visitor.audit_remark }}
            </div>
          </div>
        </div>
        
        <el-empty v-else description="暂无访客记录" image-size="120" />

        <div class="pagination-container" v-if="total > 0">
          <el-pagination
            v-model:current-page="currentPage"
            v-model:page-size="pageSize"
            :page-sizes="[6, 12, 24]"
            layout="total, sizes, prev, pager, next, jumper"
            :total="total"
            @size-change="handleSizeChange"
            @current-change="handleCurrentChange"
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
import { createVisitor, getVisitorList } from '@/api/service'
import dayjs from 'dayjs'
import { ElMessage } from 'element-plus'

const form = ref({
  name: '',
  mobile: '',
  reason: '',
  visit_time: ''
})

const loading = ref(false)
const visitors = ref([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(12) // 配合两列网格，设置为12

const formatDate = (date) => {
  return dayjs(date).format('YYYY-MM-DD HH:mm')
}

const getStatusText = (status) => {
  const map = { 0: '待审核', 1: '已通过', 2: '已拒绝' }
  return map[status] || '未知'
}

const getStatusClass = (status) => {
  const map = { 0: 'status-pending', 1: 'status-pass', 2: 'status-reject' }
  return map[status] || ''
}

const submitVisitor = async () => {
  loading.value = true
  try {
    const submitData = {
      ...form.value,
      visitor_name: form.value.name,
      visitor_phone: form.value.mobile,
      visit_time: dayjs(form.value.visit_time).format('YYYY-MM-DD HH:mm:ss')
    }
    
    await createVisitor(submitData)
    ElMessage.success('提交成功，等待物业审核！')
    form.value = { name: '', mobile: '', reason: '', visit_time: '' }
    await fetchVisitors()
  } catch (error) {
    ElMessage.error('提交失败: ' + (error.response?.data?.msg || error.message))
  } finally {
    loading.value = false
  }
}

const fetchVisitors = async () => {
  try {
    const res = await getVisitorList({
        page: currentPage.value,
        size: pageSize.value
    })
    visitors.value = res.list
    total.value = res.total
  } catch (error) {
    console.error('获取访客记录失败:', error)
  }
}

const handleSizeChange = (val) => {
  pageSize.value = val
  fetchVisitors()
}

const handleCurrentChange = (val) => {
  currentPage.value = val
  fetchVisitors()
}

onMounted(() => {
  fetchVisitors()
})
</script>

<style scoped>
.visitor-page { min-height: 100vh; background-color: #f8f9fa; padding-bottom: 80px; }
.custom-container { max-width: 1000px; margin: 0 auto; }
.page-header { padding: 32px 0 24px; }
.highlight-title { display: inline-block; position: relative; font-size: 32px; color: #2c3e50; font-weight: 700; margin: 0; z-index: 1; }
.highlight-title::after { content: ''; position: absolute; bottom: 4px; left: -5%; width: 110%; height: 14px; background-color: #2d597b; opacity: 0.15; border-radius: 6px; z-index: -1; }

.visitor-form-wrapper { background: #ffffff; border-radius: 16px; padding: 40px; box-shadow: 0 4px 20px rgba(0,0,0,0.03); margin-bottom: 48px; }
.form-header { margin-bottom: 32px; border-bottom: 1px solid #f0f2f5; padding-bottom: 16px; }
.form-header h3 { margin: 0 0 8px 0; font-size: 20px; color: #2c3e50; }
.form-header p { margin: 0; font-size: 14px; color: #909399; }

.premium-form { display: flex; flex-direction: column; gap: 24px; }
.form-row { display: flex; gap: 24px; }
.form-group { display: flex; flex-direction: column; gap: 8px; }
.half { flex: 1; }

.form-group label { font-size: 14px; font-weight: 600; color: #303133; }
.custom-input { width: 100%; padding: 14px 16px; border: 1px solid #dcdfe6; border-radius: 8px; font-size: 15px; outline: none; background: #fafbfc; transition: all 0.3s; font-family: inherit; }
.custom-input:focus { border-color: #2d597b; background: #ffffff; box-shadow: 0 0 0 3px rgba(45,89,123,0.1); }

.form-actions { text-align: right; margin-top: 8px; }
.btn-submit { padding: 12px 40px; background: #2d597b; color: #ffffff; border: none; border-radius: 8px; font-size: 16px; font-weight: bold; cursor: pointer; transition: all 0.3s; }
.btn-submit:hover:not(:disabled) { background: #1f435d; transform: translateY(-2px); box-shadow: 0 6px 16px rgba(45,89,123,0.2); }
.btn-submit:disabled { background: #a4b0be; cursor: not-allowed; }

.records-section { margin-top: 20px; }
.section-title-wrap { display: flex; align-items: center; margin-bottom: 24px; }
.indicator { width: 4px; height: 18px; background: #2d597b; border-radius: 2px; margin-right: 12px; }
.section-title-wrap h3 { margin: 0; font-size: 20px; color: #2c3e50; }

.visitor-list { display: grid; grid-template-columns: repeat(2, 1fr); gap: 24px; }

.visitor-card { background: #ffffff; border-radius: 12px; padding: 24px; border: 1px solid #ebeef5; transition: all 0.3s; display: flex; flex-direction: column; }
.visitor-card:hover { transform: translateY(-4px); box-shadow: 0 8px 24px rgba(0,0,0,0.06); border-color: rgba(45,89,123,0.1); }

.card-top { display: flex; justify-content: space-between; align-items: center; margin-bottom: 20px; padding-bottom: 16px; border-bottom: 1px dashed #f0f2f5; }
.visitor-name { font-size: 20px; font-weight: 700; color: #2d597b; letter-spacing: 1px; }

.status-badge { padding: 4px 12px; border-radius: 20px; font-size: 12px; font-weight: bold; }
.status-pending { background: #fff7ed; color: #d97706; }
.status-pass { background: #f0fdf4; color: #166534; }
.status-reject { background: #fef2f2; color: #991b1b; }

.card-body { display: flex; flex-direction: column; gap: 12px; flex: 1; }
.info-line { display: flex; align-items: flex-start; font-size: 14px; }
.info-line .label { color: #909399; width: 60px; flex-shrink: 0; }
.info-line .value { color: #303133; font-weight: 500; }

.card-footer { margin-top: 16px; padding: 12px 16px; background: #f8f9fa; border-radius: 6px; font-size: 13px; color: #606266; border-left: 3px solid #2d597b; }
.remark-label { color: #2d597b; font-weight: bold; }

.pagination-container { display: flex; justify-content: center; margin-top: 32px; grid-column: 1 / -1; }
:deep(.custom-pagination .el-pager li.is-active) { background-color: #2d597b; color: #fff; border-radius: 4px; }

@media (max-width: 768px) {
  .form-row { flex-direction: column; gap: 16px; }
  .visitor-list { grid-template-columns: 1fr; }
}
</style>