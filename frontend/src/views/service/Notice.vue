<template>
  <div class="notice-page">
    <Navbar />
    
    <div class="container custom-container">
      <div class="page-header">
        <h1 class="page-title highlight-title">公告通知</h1>
      </div>
      
      <div class="notice-list">
        <div class="notice-card" v-for="notice in notices" :key="notice.id">
          <div class="notice-main">
            <h3 class="notice-title">{{ notice.title }}</h3>
            <p class="notice-content">{{ notice.content }}</p>
            <div class="notice-footer">
              <div class="notice-meta">
                <span class="publisher"><el-icon><Avatar /></el-icon> {{ notice.publisher }}</span>
                <span class="date"><el-icon><Calendar /></el-icon> {{ formatDate(notice.created_at) }}</span>
              </div>
              <span class="views"><el-icon><View /></el-icon> {{ notice.view_count }} 浏览</span>
            </div>
          </div>
        </div>
      </div>

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
</template>

<script setup>
import { ref, onMounted } from 'vue'
import Navbar from '@/components/layout/Navbar.vue'
import { getNoticeList } from '@/api/service'
import dayjs from 'dayjs'
import { Avatar, Calendar, View } from '@element-plus/icons-vue'

const notices = ref([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)

const formatDate = (date) => {
  return dayjs(date).format('YYYY-MM-DD HH:mm')
}

const fetchList = async () => {
  try {
    const res = await getNoticeList({
        page: currentPage.value,
        size: pageSize.value
    })
    if (res.list) {
        notices.value = res.list
        total.value = res.total
    } else {
        notices.value = res
        total.value = res.length
    }
  } catch (error) {
    console.error('获取公告失败:', error)
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

onMounted(() => {
  fetchList()
})
</script>

<style scoped>
.notice-page { min-height: 100vh; background-color: #f8f9fa; padding-bottom: 80px; }
.custom-container { max-width: 1000px; margin: 0 auto; }
.page-header { padding: 32px 0 24px; }
.highlight-title { display: inline-block; position: relative; font-size: 32px; color: #2c3e50; font-weight: 700; margin: 0; z-index: 1; }
.highlight-title::after { content: ''; position: absolute; bottom: 4px; left: -5%; width: 110%; height: 14px; background-color: #2d597b; opacity: 0.15; border-radius: 6px; z-index: -1; }

.notice-list { display: flex; flex-direction: column; gap: 20px; }

.notice-card {
  padding: 32px; background: #ffffff; border-radius: 16px; box-shadow: 0 2px 12px rgba(0,0,0,0.03);
  transition: all 0.3s ease; border: 1px solid transparent; cursor: pointer;
}
.notice-card:hover { transform: translateY(-4px); box-shadow: 0 12px 30px rgba(0,0,0,0.06); border-color: rgba(45, 89, 123, 0.1); }

.notice-title { font-size: 22px; font-weight: 700; color: #2c3e50; margin: 0 0 16px 0; transition: color 0.3s; }
.notice-card:hover .notice-title { color: #2d597b; }

.notice-content { font-size: 15px; color: #606266; line-height: 1.8; margin-bottom: 24px; display: -webkit-box; -webkit-line-clamp: 3; -webkit-box-orient: vertical; overflow: hidden; }

.notice-footer { display: flex; justify-content: space-between; align-items: center; padding-top: 20px; border-top: 1px dashed #ebeef5; }

.notice-meta { display: flex; gap: 24px; }
.publisher, .date { display: flex; align-items: center; gap: 6px; font-size: 14px; color: #909399; }
.publisher { color: #2d597b; font-weight: 500; background: #f0f7ff; padding: 4px 12px; border-radius: 20px; }

.views { display: flex; align-items: center; gap: 4px; font-size: 13px; color: #a4b0be; }

.pagination-container { display: flex; justify-content: center; margin-top: 40px; }
:deep(.custom-pagination .el-pager li.is-active) { background-color: #2d597b; color: #fff; border-radius: 4px; }
:deep(.custom-pagination .el-pager li:hover) { color: #2d597b; }
</style>