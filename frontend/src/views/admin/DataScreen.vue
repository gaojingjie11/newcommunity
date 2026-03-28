<template>
  <div class="data-screen">
    <Navbar />

    <div class="container">
      <div class="screen-header">
        <div>
          <h1 class="page-title">社区智能分析大屏</h1>
          <p class="sub-title">更新时间：{{ currentTime }}</p>
        </div>
        <div class="header-actions">
          <el-button v-if="isAdmin" @click="router.push('/admin/ai-report')">AI 报表中心</el-button>
          <el-button type="primary" @click="refreshAll">刷新数据</el-button>
        </div>
      </div>

      <div class="stats-grid">
        <el-card class="metric-card">
          <span class="metric-label">总用户数</span>
          <strong class="metric-value">{{ stats.totalUsers || 0 }}</strong>
        </el-card>
        <el-card class="metric-card">
          <span class="metric-label">今日订单</span>
          <strong class="metric-value">{{ stats.todayOrders || 0 }}</strong>
        </el-card>
        <el-card class="metric-card">
          <span class="metric-label">月营收</span>
          <strong class="metric-value">¥{{ formatAmount(stats.monthIncome) }}</strong>
        </el-card>
        <el-card class="metric-card">
          <span class="metric-label">车位占用率</span>
          <strong class="metric-value">{{ stats.parkingRate || '0%' }}</strong>
        </el-card>
      </div>

      <div class="content-grid">
        <el-card>
          <template #header>7 日营收趋势</template>
          <div ref="lineChartRef" class="chart-box"></div>
        </el-card>

        <el-card>
          <template #header>报修分类分布</template>
          <div ref="pieChartRef" class="chart-box"></div>
        </el-card>

        <el-card>
          <template #header>
            <div class="ranking-header">
              <span>累计绿色积分排行榜</span>
              <el-radio-group v-model="rankingView" size="small">
                <el-radio-button label="datav">DataV</el-radio-button>
                <el-radio-button label="table">表格</el-radio-button>
              </el-radio-group>
            </div>
          </template>

          <dv-scroll-ranking-board
            v-if="rankingView === 'datav'"
            class="ranking-board"
            :config="rankingBoardConfig"
          />

          <el-table v-else :data="leaderboard" stripe height="320">
            <el-table-column prop="rank" label="#" width="60" />
            <el-table-column label="用户">
              <template #default="{ row }">
                <div class="table-user">
                  <el-avatar :src="row.avatar" />
                  <span>{{ row.nickname || row.username }}</span>
                </div>
              </template>
            </el-table-column>
            <el-table-column prop="points" label="积分" width="100" />
          </el-table>
        </el-card>

        <el-card v-if="isAdmin">
          <template #header>AI 社区分析报表</template>
          <div class="report-summary" v-if="aiReport">
            <div class="summary-item">近 7 日新增报修：{{ aiReport.repair_new_count }}</div>
            <div class="summary-item">未处理报修：{{ aiReport.repair_pending_count }}</div>
            <div class="summary-item">新增访客：{{ aiReport.visitor_new_count }}</div>
            <div class="summary-item">物业缴费金额：¥{{ formatAmount(aiReport.property_paid_amount) }}</div>
          </div>
          <el-scrollbar height="220px">
            <div class="report-text">{{ aiReport?.report || '暂无数据' }}</div>
          </el-scrollbar>
        </el-card>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted, onUnmounted, ref } from 'vue'
import * as echarts from 'echarts'
import dayjs from 'dayjs'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import Navbar from '@/components/layout/Navbar.vue'
import { getDashboardStats, getAIReport } from '@/api/admin'
import { getGreenPointsLeaderboard } from '@/api/greenPoints'
import { useUserStore } from '@/stores/user'

const router = useRouter()
const userStore = useUserStore()
const currentTime = ref(dayjs().format('YYYY-MM-DD HH:mm:ss'))
const stats = ref({})
const aiReport = ref(null)
const leaderboard = ref([])
const rankingView = ref('datav')
const isAdmin = computed(() => userStore.userInfo.role === 'admin')

const lineChartRef = ref(null)
const pieChartRef = ref(null)
let lineChart
let pieChart
let timer

const rankingBoardConfig = computed(() => ({
  rowNum: 5,
  waitTime: 3500,
  carousel: 'single',
  unit: '分',
  data: leaderboard.value.map((item) => ({
    name: buildRankName(item),
    value: item.points || 0
  }))
}))

function formatAmount(value) {
  return Number(value || 0).toFixed(2)
}

function escapeHtml(text) {
  return String(text || '')
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
    .replace(/'/g, '&#39;')
}

function buildRankName(item) {
  const nickname = escapeHtml(item.nickname || item.username || `用户${item.user_id}`)
  const avatarUrl = escapeHtml(item.avatar || '')

  if (!avatarUrl) {
    return `<span style="display:flex;align-items:center;gap:8px;"><span>${nickname}</span></span>`
  }
  return `<span style="display:flex;align-items:center;gap:8px;"><img src="${avatarUrl}" style="width:24px;height:24px;border-radius:50%;object-fit:cover;" /><span>${nickname}</span></span>`
}

function renderCharts() {
  if (lineChart) lineChart.dispose()
  if (pieChart) pieChart.dispose()

  lineChart = echarts.init(lineChartRef.value)
  pieChart = echarts.init(pieChartRef.value)

  lineChart.setOption({
    tooltip: { trigger: 'axis' },
    xAxis: {
      type: 'category',
      data: stats.value.incomeDates || []
    },
    yAxis: { type: 'value' },
    series: [{
      data: stats.value.incomeTrend || [],
      type: 'line',
      smooth: true,
      areaStyle: {}
    }]
  })

  pieChart.setOption({
    tooltip: { trigger: 'item' },
    series: [{
      type: 'pie',
      radius: ['40%', '70%'],
      data: stats.value.repairStats || []
    }]
  })
}

async function refreshAll() {
  try {
    const [dashboardRes, leaderboardRes] = await Promise.all([
      getDashboardStats(),
      getGreenPointsLeaderboard({ limit: 10 })
    ])

    stats.value = dashboardRes || {}
    leaderboard.value = leaderboardRes.list || []

    if (isAdmin.value) {
      aiReport.value = await getAIReport()
    } else {
      aiReport.value = null
    }
    renderCharts()
  } catch (error) {
    ElMessage.error(error.response?.data?.msg || error.message || '刷新数据失败')
  }
}

function handleResize() {
  lineChart?.resize()
  pieChart?.resize()
}

onMounted(async () => {
  timer = setInterval(() => {
    currentTime.value = dayjs().format('YYYY-MM-DD HH:mm:ss')
  }, 1000)
  await refreshAll()
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  clearInterval(timer)
  window.removeEventListener('resize', handleResize)
  lineChart?.dispose()
  pieChart?.dispose()
})
</script>

<style scoped>
.data-screen {
  min-height: 100vh;
  padding-bottom: 40px;
  background: linear-gradient(180deg, #f7fbff 0%, #eef5ec 100%);
}

.screen-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-end;
  margin-bottom: 24px;
}

.header-actions {
  display: flex;
  gap: 12px;
}

.sub-title {
  color: var(--text-secondary);
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
  margin-bottom: 24px;
}

.metric-card {
  border-radius: 16px;
}

.metric-label {
  display: block;
  color: var(--text-secondary);
  margin-bottom: 8px;
}

.metric-value {
  font-size: 28px;
  color: #1f7a4d;
}

.content-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 20px;
}

.chart-box {
  height: 320px;
}

.ranking-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 12px;
}

.ranking-board {
  height: 320px;
}

.table-user {
  display: flex;
  align-items: center;
  gap: 12px;
}

.report-summary {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 8px;
  margin-bottom: 12px;
}

.summary-item {
  padding: 10px 12px;
  border-radius: 10px;
  background: #f5f9f3;
}

.report-text {
  line-height: 1.8;
  white-space: pre-wrap;
}

@media (max-width: 960px) {
  .stats-grid,
  .content-grid {
    grid-template-columns: 1fr;
  }

  .screen-header {
    align-items: flex-start;
    flex-direction: column;
    gap: 12px;
  }
}
</style>
