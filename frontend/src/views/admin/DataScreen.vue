<template>
  <div class="data-screen">
    <Navbar />

    <div class="container custom-container">
      <div class="screen-header">
        <div class="header-left">
          <h1 class="page-title highlight-title">社区智能分析大盘</h1>
          <p class="sub-title">
            <el-icon><Clock /></el-icon> 数据最后更新于：{{ currentTime }}
          </p>
        </div>
        <div class="header-actions">
          <button v-if="isAdmin" class="action-btn btn-outline" @click="router.push('/admin/ai-report')">
            <el-icon><DataAnalysis /></el-icon> AI 报表中心
          </button>
          <button class="action-btn btn-primary" @click="refreshAll">
            <el-icon><Refresh /></el-icon> 刷新数据
          </button>
        </div>
      </div>

      <!-- 核心指标卡片区 -->
      <div class="stats-grid">
        <div class="metric-card">
          <div class="metric-icon-wrap icon-user">
            <el-icon><User /></el-icon>
          </div>
          <div class="metric-content">
            <span class="metric-label">总注册用户数</span>
            <strong class="metric-value">{{ stats.totalUsers || 0 }}</strong>
          </div>
        </div>
        
        <div class="metric-card">
          <div class="metric-icon-wrap icon-order">
            <el-icon><ShoppingBag /></el-icon>
          </div>
          <div class="metric-content">
            <span class="metric-label">今日新增订单</span>
            <strong class="metric-value">{{ stats.todayOrders || 0 }}</strong>
          </div>
        </div>
        
        <div class="metric-card">
          <div class="metric-icon-wrap icon-income">
            <el-icon><Money /></el-icon>
          </div>
          <div class="metric-content">
            <span class="metric-label">本月累计营收</span>
            <strong class="metric-value highlight-income"><span class="currency">¥</span>{{ formatAmount(stats.monthIncome) }}</strong>
          </div>
        </div>
        
        <div class="metric-card">
          <div class="metric-icon-wrap icon-parking">
            <el-icon><Van /></el-icon>
          </div>
          <div class="metric-content">
            <span class="metric-label">当前车位占用率</span>
            <strong class="metric-value">{{ stats.parkingRate || '0%' }}</strong>
          </div>
        </div>
      </div>

      <!-- 图表与列表区 -->
      <div class="content-grid">
        <!-- 营收趋势图 -->
        <div class="premium-card">
          <div class="card-header">
            <span class="header-indicator"></span> 7 日营收趋势分析
          </div>
          <div class="card-body">
            <div ref="lineChartRef" class="chart-box"></div>
          </div>
        </div>

        <!-- 报修分类饼图 -->
        <div class="premium-card">
          <div class="card-header">
            <span class="header-indicator"></span> 工单问题分类占比
          </div>
          <div class="card-body">
            <div ref="pieChartRef" class="chart-box"></div>
          </div>
        </div>

        <!-- 积分排行榜 -->
        <div class="premium-card">
          <div class="card-header ranking-header">
            <div class="header-title">
              <span class="header-indicator"></span> 累计绿色积分排行榜
            </div>
            <el-radio-group v-model="rankingView" size="small" class="custom-radio-group">
              <el-radio-button label="datav">动态展示</el-radio-button>
              <el-radio-button label="table">经典表格</el-radio-button>
            </el-radio-group>
          </div>
          <div class="card-body">
            <dv-scroll-ranking-board
              v-if="rankingView === 'datav'"
              class="ranking-board"
              :config="rankingBoardConfig"
            />

            <el-table v-else :data="leaderboard" class="custom-table" height="320">
              <el-table-column prop="rank" label="排名" width="80" align="center">
                <template #default="{ row }">
                  <span class="rank-badge" :class="`rank-${row.rank}`">#{{ row.rank }}</span>
                </template>
              </el-table-column>
              <el-table-column label="社区之星">
                <template #default="{ row }">
                  <div class="table-user">
                    <img :src="row.avatar || 'https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png'" class="user-avatar" />
                    <span class="user-name">{{ row.nickname || row.username }}</span>
                  </div>
                </template>
              </el-table-column>
              <el-table-column prop="points" label="环保积分" width="120" align="right">
                <template #default="{ row }">
                  <strong class="points-text">{{ row.points }}</strong>
                </template>
              </el-table-column>
            </el-table>
          </div>
        </div>

        <!-- AI 社区报表 -->
        <div class="premium-card" v-if="isAdmin">
          <div class="card-header">
            <span class="header-indicator"></span> AI 社区运营诊断报表
          </div>
          <div class="card-body">
            <div class="report-summary" v-if="aiReport">
              <div class="summary-item repair">近 7 日新增报修：<strong>{{ aiReport.repair_new_count }}</strong> 单</div>
              <div class="summary-item pending">未处理报修：<strong>{{ aiReport.repair_pending_count }}</strong> 单</div>
              <div class="summary-item visitor">新增访客：<strong>{{ aiReport.visitor_new_count }}</strong> 人</div>
              <div class="summary-item income">物业收缴：<strong>¥{{ formatAmount(aiReport.property_paid_amount) }}</strong></div>
            </div>
            <div class="report-content-wrap custom-scrollbar">
              <div class="report-text">{{ aiReport?.report || '正在等待 AI 生成数据诊断...' }}</div>
            </div>
          </div>
        </div>
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
// 引入必需的图标
import { Clock, DataAnalysis, Refresh, User, ShoppingBag, Money, Van } from '@element-plus/icons-vue'

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
    tooltip: { trigger: 'axis', backgroundColor: 'rgba(255,255,255,0.9)', borderColor: '#ebeef5', textStyle: { color: '#303133' } },
    grid: { left: '3%', right: '4%', bottom: '3%', containLabel: true },
    xAxis: {
      type: 'category',
      boundaryGap: false,
      data: stats.value.incomeDates || [],
      axisLine: { lineStyle: { color: '#dcdfe6' } },
      axisLabel: { color: '#606266' }
    },
    yAxis: { 
      type: 'value',
      splitLine: { lineStyle: { color: '#ebeef5', type: 'dashed' } },
      axisLabel: { color: '#606266' }
    },
    series: [{
      data: stats.value.incomeTrend || [],
      type: 'line',
      smooth: true,
      symbol: 'circle',
      symbolSize: 8,
      itemStyle: { color: '#2d597b' },
      lineStyle: { color: '#2d597b', width: 3 },
      areaStyle: {
        color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
          { offset: 0, color: 'rgba(45, 89, 123, 0.4)' },
          { offset: 1, color: 'rgba(45, 89, 123, 0.05)' }
        ])
      }
    }]
  })

  pieChart.setOption({
    tooltip: { trigger: 'item', backgroundColor: 'rgba(255,255,255,0.9)', borderColor: '#ebeef5', textStyle: { color: '#303133' } },
    legend: { bottom: '0%', left: 'center', itemWidth: 10, itemHeight: 10, textStyle: { color: '#606266' } },
    color: ['#2d597b', '#00b894', '#f39c12', '#e74c3c', '#8e44ad', '#0984e3'],
    series: [{
      type: 'pie',
      radius: ['45%', '70%'],
      center: ['50%', '45%'],
      avoidLabelOverlap: false,
      itemStyle: { borderRadius: 6, borderColor: '#fff', borderWidth: 2 },
      label: { show: false, position: 'center' },
      emphasis: {
        label: { show: true, fontSize: 18, fontWeight: 'bold' }
      },
      labelLine: { show: false },
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
/* 全局页面底色与容器 */
.data-screen {
  min-height: 100vh;
  background-color: #f4f7f9; /* 浅灰蓝色，衬托纯白卡片 */
  padding-bottom: 80px;
}

.custom-container {
  max-width: 1400px; /* 大屏视野需要更宽的容器 */
  margin: 0 auto;
}

/* 头部信息区 */
.screen-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-end;
  padding: 32px 0 24px;
  margin-bottom: 12px;
}

.header-left {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.highlight-title {
  display: inline-block;
  position: relative;
  font-size: 32px;
  color: #2c3e50;
  font-weight: 800;
  margin: 0;
  z-index: 1;
  letter-spacing: 1px;
}

.highlight-title::after {
  content: '';
  position: absolute;
  bottom: 4px;
  left: -2%;
  width: 104%;
  height: 14px;
  background-color: #2d597b; 
  opacity: 0.15;
  border-radius: 6px;
  z-index: -1;
}

.sub-title {
  color: #8c939d;
  font-size: 14px;
  margin: 0;
  display: flex;
  align-items: center;
  gap: 6px;
}

.header-actions {
  display: flex;
  gap: 16px;
}

/* 定制操作按钮 */
.action-btn {
  padding: 10px 24px;
  border-radius: 20px;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s;
  border: 1px solid transparent;
  display: inline-flex;
  align-items: center;
  gap: 6px;
}

.btn-primary { background: #2d597b; color: #ffffff; box-shadow: 0 4px 12px rgba(45, 89, 123, 0.2); }
.btn-primary:hover { background: #1f435d; transform: translateY(-2px); box-shadow: 0 6px 16px rgba(45, 89, 123, 0.3); }

.btn-outline { background: #ffffff; color: #2d597b; border-color: #2d597b; }
.btn-outline:hover { background: #f0f7ff; transform: translateY(-2px); }


/* ================= 指标卡片矩阵 ================= */
.stats-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 24px;
  margin-bottom: 32px;
}

.metric-card {
  background: #ffffff;
  border-radius: 16px;
  padding: 24px;
  display: flex;
  align-items: center;
  gap: 20px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.02);
  transition: all 0.4s cubic-bezier(0.175, 0.885, 0.32, 1.275);
  border: 1px solid transparent;
}

.metric-card:hover {
  transform: translateY(-6px);
  box-shadow: 0 12px 24px rgba(45, 89, 123, 0.08);
  border-color: rgba(45, 89, 123, 0.1);
}

.metric-icon-wrap {
  width: 64px;
  height: 64px;
  border-radius: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 32px;
}

/* 各模块独立颜色 */
.icon-user { background: #f0f7ff; color: #409eff; }
.icon-order { background: #fdf6f6; color: #e4393c; }
.icon-income { background: #f0fdf4; color: #00b894; }
.icon-parking { background: #fff7ed; color: #e6a23c; }

.metric-content {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.metric-label {
  color: #909399;
  font-size: 14px;
  font-weight: 500;
}

.metric-value {
  font-size: 32px;
  font-weight: 800;
  color: #2c3e50;
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
  letter-spacing: -0.5px;
}

.highlight-income {
  color: #e4393c;
}
.currency {
  font-size: 20px;
  margin-right: 4px;
  font-weight: 600;
}


/* ================= 图表与列表区 ================= */
.content-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 24px;
}

.premium-card {
  background: #ffffff;
  border-radius: 16px;
  box-shadow: 0 2px 16px rgba(0, 0, 0, 0.02);
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.card-header {
  padding: 20px 24px;
  border-bottom: 1px solid #f0f2f5;
  font-size: 18px;
  font-weight: 700;
  color: #2c3e50;
  display: flex;
  align-items: center;
}

.header-indicator {
  display: inline-block;
  width: 4px;
  height: 18px;
  background: #2d597b;
  border-radius: 2px;
  margin-right: 12px;
}

.card-body {
  padding: 24px;
  flex: 1;
}

.chart-box {
  height: 340px;
  width: 100%;
}

/* 积分排行榜定制 */
.ranking-header {
  justify-content: space-between;
}
.header-title {
  display: flex;
  align-items: center;
}

:deep(.custom-radio-group .el-radio-button__inner) {
  border-radius: 6px !important;
  margin: 0 4px;
  border: 1px solid #dcdfe6;
  box-shadow: none !important;
}
:deep(.custom-radio-group .el-radio-button__original-radio:checked + .el-radio-button__inner) {
  background-color: #2d597b;
  border-color: #2d597b;
  color: #ffffff;
}

.ranking-board {
  height: 340px;
  width: 100%;
}

/* 定制排行榜表格 */
:deep(.custom-table) {
  --el-table-header-bg-color: #fbfcfd;
  --el-table-border-color: transparent;
}
:deep(.custom-table th.el-table__cell) { font-weight: 600; padding: 12px 0; border-bottom: 1px solid #ebeef5; }
:deep(.custom-table td.el-table__cell) { padding: 12px 0; border-bottom: 1px dashed #f0f2f5; }
:deep(.custom-table::before) { display: none; }

.rank-badge {
  display: inline-block;
  padding: 2px 10px;
  border-radius: 12px;
  font-size: 13px;
  font-weight: bold;
  background: #f4f4f5;
  color: #909399;
}
.rank-1 { background: #fffbe6; color: #e6a23c; font-size: 14px; }
.rank-2 { background: #f0f9eb; color: #67c23a; font-size: 13px; }
.rank-3 { background: #fdf6f6; color: #f56c6c; font-size: 13px; }

.table-user {
  display: flex;
  align-items: center;
  gap: 12px;
}
.user-avatar {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  object-fit: cover;
  border: 1px solid #ebeef5;
}
.user-name {
  font-weight: 600;
  color: #303133;
}
.points-text {
  color: #497db7;
  font-size: 16px;
  font-weight: 800;
}

/* AI 报表区域 */
.report-summary {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16px;
  margin-bottom: 24px;
}

.summary-item {
  padding: 16px;
  border-radius: 12px;
  font-size: 14px;
  color: #606266;
  border: 1px solid transparent;
}
.summary-item strong { font-size: 18px; margin: 0 4px; }

.summary-item.repair { background: #fdf6f6; border-color: #fbc4c4; color: #e4393c; }
.summary-item.pending { background: #fff7ed; border-color: #fed7aa; color: #d97706; }
.summary-item.visitor { background: #f0f7ff; border-color: #cce3f6; color: #2d597b; }
.summary-item.income { background: #f0fdf4; border-color: #bbf7d0; color: #166534; }

.report-content-wrap {
  height: 240px;
  overflow-y: auto;
  background: #fafbfc;
  border-radius: 12px;
  border: 1px solid #ebeef5;
  padding: 20px;
}

.report-text {
  line-height: 1.8;
  color: #475569;
  font-size: 15px;
  white-space: pre-wrap;
}

.custom-scrollbar::-webkit-scrollbar { width: 6px; }
.custom-scrollbar::-webkit-scrollbar-track { background: transparent; }
.custom-scrollbar::-webkit-scrollbar-thumb { background: #dcdfe6; border-radius: 3px; }
.custom-scrollbar::-webkit-scrollbar-thumb:hover { background: #c0c4cc; }

/* 响应式 */
@media (max-width: 1024px) {
  .stats-grid, .content-grid { grid-template-columns: 1fr; }
  .screen-header { flex-direction: column; align-items: flex-start; gap: 16px; }
}
</style>