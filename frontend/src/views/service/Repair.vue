<template>
  <div class="repair-page">
    <Navbar />
    
    <div class="container custom-container">
      <div class="page-header">
        <h1 class="page-title highlight-title">报修与投诉</h1>
      </div>
      
      <div class="repair-form-wrapper">
        <div class="form-header">
          <h3>提交新工单</h3>
          <p>物业中心将尽快为您安排处理</p>
        </div>
        <form @submit.prevent="submitRepair" class="premium-form">
          <div class="form-row">
            <div class="form-group half">
              <label>服务类型</label>
              <select v-model="form.type" class="custom-input">
                <option :value="1">设备报修</option>
                <option :value="2">建议投诉</option>
              </select>
            </div>
            
            <div class="form-group half">
              <label>问题分类</label>
              <input v-model="form.category" class="custom-input" placeholder="如: 水电、门窗、环境等" required />
            </div>
          </div>
          
          <div class="form-group">
            <label>详细描述</label>
            <textarea v-model="form.content" class="custom-input textarea-input" rows="4" placeholder="请详细描述您遇到的问题，以便我们更好地定位和解决..." required></textarea>
          </div>
          
          <div class="form-actions">
            <button type="submit" class="btn-submit" :disabled="loading">
              {{ loading ? '正在提交...' : '确认提交' }}
            </button>
          </div>
        </form>
      </div>
      
      <div class="records-section">
        <div class="section-title-wrap">
          <span class="indicator"></span>
          <h3>我的历史记录</h3>
        </div>

        <div class="repair-list" v-if="repairs.length > 0">
          <div class="record-card" v-for="item in repairs" :key="item.id">
            <div class="record-header">
              <div class="header-left">
                <span class="type-badge" :class="item.type === 1 ? 'is-repair' : 'is-complain'">
                  {{ item.type === 1 ? '报修' : '投诉' }}
                </span>
                <span class="category-name">{{ item.category }}</span>
              </div>
              <span class="status-badge" :class="getStatusClass(item.status)">
                {{ getStatusText(item.status) }}
              </span>
            </div>
            
            <div class="record-body">
              <p class="content-text">{{ item.content }}</p>
              <div class="result-box" v-if="item.result">
                <span class="result-label">处理结果：</span>
                <span class="result-text">{{ item.result }}</span>
              </div>
            </div>
            
            <div class="record-footer">
              提交时间: {{ formatDate(item.created_at) }}
            </div>
          </div>
        </div>
        
        <el-empty v-else description="暂无历史记录" image-size="120" />

        <div class="pagination-container" v-if="total > 0">
          <el-pagination
            v-model:current-page="currentPage"
            v-model:page-size="pageSize"
            :page-sizes="[5, 10, 20]"
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
import { createRepair, getRepairList } from '@/api/service'
import dayjs from 'dayjs'
import { ElMessage } from 'element-plus'

const form = ref({
  type: 1,
  category: '',
  content: ''
})

const loading = ref(false)
const repairs = ref([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)

const formatDate = (date) => {
  return dayjs(date).format('YYYY-MM-DD HH:mm')
}

const getStatusText = (status) => {
  const map = { 0: '待处理', 1: '处理中', 2: '已完成' }
  return map[status] || '未知'
}

const getStatusClass = (status) => {
  const map = { 0: 'status-pending', 1: 'status-processing', 2: 'status-done' }
  return map[status] || ''
}

const submitRepair = async () => {
  loading.value = true
  try {
    await createRepair(form.value)
    ElMessage.success('提交成功！')
    form.value = { type: 1, category: '', content: '' }
    await fetchRepairs()
  } catch (error) {
    ElMessage.error('提交失败')
  } finally {
    loading.value = false
  }
}

const fetchRepairs = async () => {
  try {
    const res = await getRepairList({
        page: currentPage.value,
        size: pageSize.value
    })
    repairs.value = res.list
    total.value = res.total
  } catch (error) {
    console.error('获取记录失败:', error)
  }
}

const handleSizeChange = (val) => {
  pageSize.value = val
  fetchRepairs()
}

const handleCurrentChange = (val) => {
  currentPage.value = val
  fetchRepairs()
}

onMounted(() => {
  fetchRepairs()
})
</script>

<style scoped>
.repair-page { min-height: 100vh; background-color: #f8f9fa; padding-bottom: 80px; }
.custom-container { max-width: 900px; margin: 0 auto; }
.page-header { padding: 32px 0 24px; }
.highlight-title { display: inline-block; position: relative; font-size: 32px; color: #2c3e50; font-weight: 700; margin: 0; z-index: 1; }
.highlight-title::after { content: ''; position: absolute; bottom: 4px; left: -5%; width: 110%; height: 14px; background-color: #2d597b; opacity: 0.15; border-radius: 6px; z-index: -1; }

.repair-form-wrapper { background: #ffffff; border-radius: 16px; padding: 40px; box-shadow: 0 4px 20px rgba(0,0,0,0.03); margin-bottom: 48px; }
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
.textarea-input { resize: vertical; min-height: 120px; }

.form-actions { text-align: right; margin-top: 16px; }
.btn-submit { padding: 12px 40px; background: #2d597b; color: #ffffff; border: none; border-radius: 8px; font-size: 16px; font-weight: bold; cursor: pointer; transition: all 0.3s; }
.btn-submit:hover:not(:disabled) { background: #1f435d; transform: translateY(-2px); box-shadow: 0 6px 16px rgba(45,89,123,0.2); }
.btn-submit:disabled { background: #a4b0be; cursor: not-allowed; }

.records-section { margin-top: 20px; }
.section-title-wrap { display: flex; align-items: center; margin-bottom: 24px; }
.indicator { width: 4px; height: 18px; background: #2d597b; border-radius: 2px; margin-right: 12px; }
.section-title-wrap h3 { margin: 0; font-size: 20px; color: #2c3e50; }

.repair-list { display: flex; flex-direction: column; gap: 20px; }
.record-card { background: #ffffff; border-radius: 12px; padding: 24px; border: 1px solid #ebeef5; transition: all 0.3s; }
.record-card:hover { transform: translateX(4px); box-shadow: 0 4px 16px rgba(0,0,0,0.05); border-color: rgba(45,89,123,0.2); }

.record-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px; }
.header-left { display: flex; align-items: center; gap: 12px; }
.type-badge { padding: 4px 10px; border-radius: 4px; font-size: 12px; font-weight: bold; }
.is-repair { background: #f0f7ff; color: #2d597b; border: 1px solid #cce3f6; }
.is-complain { background: #fdf6f6; color: #e4393c; border: 1px solid #fbc4c4; }
.category-name { font-size: 16px; font-weight: 600; color: #303133; }

.status-badge { padding: 4px 12px; border-radius: 20px; font-size: 12px; font-weight: bold; }
.status-pending { background: #fff7ed; color: #d97706; }
.status-processing { background: #f0f7ff; color: #0984e3; }
.status-done { background: #f0fdf4; color: #166534; }

.record-body { margin-bottom: 16px; }
.content-text { margin: 0 0 16px 0; font-size: 15px; color: #606266; line-height: 1.6; }
.result-box { background: #f8f9fa; padding: 16px; border-radius: 8px; border-left: 4px solid #00b894; }
.result-label { font-size: 13px; color: #909399; margin-bottom: 4px; display: block; }
.result-text { font-size: 14px; color: #303133; font-weight: 500; }

.record-footer { font-size: 13px; color: #a4b0be; text-align: right; }

.pagination-container { display: flex; justify-content: center; margin-top: 32px; }
:deep(.custom-pagination .el-pager li.is-active) { background-color: #2d597b; color: #fff; border-radius: 4px; }
</style>