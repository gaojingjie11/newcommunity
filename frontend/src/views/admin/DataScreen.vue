<template>
  <div class="data-screen">
    <!-- 全屏容器，使用 DataV 提供的缩放适配方案 -->
    <dv-full-screen-container>
      
      <!-- ================== 顶部 Header ================== -->
      <div class="screen-header">
        <dv-decoration-8 style="width:300px;height:50px;" />
        <div class="header-center">
          <dv-decoration-5 style="width:500px;height:40px;" />
          <div class="title-text">智慧社区数据中枢大屏</div>
        </div>
        <dv-decoration-8 :reverse="true" style="width:300px;height:50px;" />
        <div class="time-text">{{ currentTime }}</div>
        <div class="back-btn" @click="goBack">
          <el-icon><HomeFilled /></el-icon> 首页
        </div>
      </div>

      <!-- ================== 主体 Body ================== -->
      <div class="screen-body">
        
        <!-- ================== 左侧栏 ================== -->
        <div class="column-side">
          <!-- 1. 核心运行指标 -->
          <dv-border-box-11 class="tech-card box-h-30" title="实时核心指标">
            <div class="metrics-grid">
              <div class="metric-item">
                <div class="m-icon" style="color: #00f2fe;"><el-icon><User /></el-icon></div>
                <div class="m-info">
                  <span class="m-label">总注册用户</span>
                  <span class="m-value">{{ stats.totalUsers || 0 }}</span>
                </div>
              </div>
              <div class="metric-item">
                <div class="m-icon" style="color: #4facfe;"><el-icon><ShoppingBag /></el-icon></div>
                <div class="m-info">
                  <span class="m-label">今日新增订单</span>
                  <span class="m-value">{{ stats.todayOrders || 0 }}</span>
                </div>
              </div>
              <div class="metric-item">
                <div class="m-icon" style="color: #409EFF;"><el-icon><Van /></el-icon></div>
                <div class="m-info">
                  <span class="m-label">车位占用率</span>
                  <span class="m-value">{{ stats.parkingRate || '0%' }}</span>
                </div>
              </div>
              <div class="metric-item">
                <div class="m-icon" style="color: #79bbff;"><el-icon><Money /></el-icon></div>
                <div class="m-info">
                  <span class="m-label">本月累计营收</span>
                  <span class="m-value num-small">¥{{ formatAmount(stats.monthIncome) }}</span>
                </div>
              </div>
            </div>
          </dv-border-box-11>

          <!-- 2. 工单报修分类饼图 -->
          <dv-border-box-13 class="tech-card box-h-35 mt-15">
             <div class="chart-title"><span class="indicator"></span>工单问题分类占比</div>
             <div ref="pieChartRef" class="chart-container"></div>
          </dv-border-box-13>

          <!-- 3. 营收趋势折线图 -->
          <dv-border-box-13 class="tech-card box-h-35 mt-15">
             <div class="chart-title"><span class="indicator"></span>7日营收趋势分析</div>
             <div ref="lineChartRef" class="chart-container"></div>
          </dv-border-box-13>
        </div>

        <!-- ================== 中间栏 ================== -->
        <div class="column-center">
           <div class="center-map-box">
             <dv-border-box-10>
                 <!-- 中间全景图 -->
                 <div class="map-bg"></div>
                 <div class="map-mask"></div>
                 <div class="map-title-overlay">社区 3D 态势感知模型</div>
                 
                 <!-- 核心悬浮数据：数字翻牌器 -->
                 <div class="center-data">
                    <div class="c-item">
                       <div class="c-label">当月社区总营收 (元)</div>
                       <div class="c-num">
                         <dv-digital-flop :config="flopIncomeConfig" style="width:240px;height:50px;" />
                       </div>
                    </div>
                    <div class="c-item">
                       <div class="c-label">活跃用户基数 (人)</div>
                       <div class="c-num">
                         <dv-digital-flop :config="flopUserConfig" style="width:200px;height:50px;" />
                       </div>
                    </div>
                 </div>
                 
                 <!-- 模拟的扫描雷达光效 -->
                 <div class="radar-scan"></div>
             </dv-border-box-10>
           </div>
        </div>

        <!-- ================== 右侧栏 ================== -->
        <div class="column-side">
           <!-- 1. 绿色积分排行榜 -->
           <dv-border-box-13 class="tech-card box-h-35">
             <div class="chart-title flex-between">
                <div class="title-left"><span class="indicator"></span>社区环保积分先锋榜</div>
                <!-- 切换视图组件 -->
                <el-radio-group v-model="rankingView" size="small" class="dark-radio">
                  <el-radio-button label="datav">动态展示</el-radio-button>
                  <el-radio-button label="table">经典表格</el-radio-button>
                </el-radio-group>
             </div>
             <div class="ranking-wrap">
               <!-- DataV 轮播排名榜 -->
               <dv-scroll-ranking-board 
                 v-if="rankingView === 'datav' && rankingBoardConfig.data.length" 
                 :config="rankingBoardConfig" 
                 style="width:100%;height:100%" 
               />
               
               <!-- 经典表格视图 -->
               <el-table
                 v-else-if="rankingView === 'table' && leaderboardList.length"
                 :data="leaderboardList"
                 class="dark-theme-table custom-scrollbar"
                 height="100%"
               >
                 <el-table-column label="排名" width="60" align="center">
                   <template #default="scope">
                     <span class="rank-badge" :class="'rank-' + (scope.$index + 1)">
                       {{ scope.$index + 1 }}
                     </span>
                   </template>
                 </el-table-column>
                 <el-table-column label="社区之星" show-overflow-tooltip>
                   <template #default="scope">
                     {{ scope.row.nickname || scope.row.username || `用户${scope.row.user_id}` }}
                   </template>
                 </el-table-column>
                 <el-table-column prop="points" label="环保积分" width="85" align="right">
                   <template #default="scope">
                     <strong style="color: #00f2fe;">{{ scope.row.points }}</strong>
                   </template>
                 </el-table-column>
               </el-table>

               <div v-else class="empty-data">暂无排行数据</div>
             </div>
           </dv-border-box-13>

           <!-- 2. 费用缴费构成 (增加的展示图表) -->
           <dv-border-box-13 class="tech-card box-h-35 mt-15">
             <div class="chart-title"><span class="indicator"></span>社区各模块收入构成</div>
             <div ref="barChartRef" class="chart-container"></div>
           </dv-border-box-13>

           <!-- 3. AI 运营诊断大屏播报 -->
           <dv-border-box-11 class="tech-card box-h-30 mt-15" title="AI 智能预警与诊断">
             <div class="ai-report-wrap">
               <div class="ai-status">
                  <div class="status-left">
                    <div class="status-dot pulse"></div>
                    <span>AI 守护引擎运行中</span>
                  </div>
                  <!-- 新增展开按钮 -->
                  <div class="status-right expand-btn" @click="dialogVisible = true">
                    <el-icon><FullScreen /></el-icon> 展开
                  </div>
               </div>
               <div class="report-content custom-scrollbar">
                 <p class="ai-text line-clamp">{{ aiReport?.report || '正在实时诊断社区运营数据，请稍候...' }}</p>
                 
                 <!-- 如果有具体诊断数据，展示几个关键tag -->
                 <div class="ai-tags" v-if="aiReport">
                    <span class="tag danger" v-if="aiReport.repair_pending_count > 0">
                      待办报修: {{ aiReport.repair_pending_count }}
                    </span>
                    <span class="tag success" v-if="aiReport.visitor_new_count">
                      新增访客: {{ aiReport.visitor_new_count }}
                    </span>
                 </div>
               </div>
             </div>
           </dv-border-box-11>
        </div>

      </div>

      <!-- AI 诊断内容展开对话框：采用 parsedReport 渲染高亮格式 -->
      <el-dialog
        v-model="dialogVisible"
        title="AI 智能预警与诊断深度报告"
        width="680px"
        class="dark-theme-dialog"
        append-to-body
        destroy-on-close
      >
        <div class="dialog-report-content custom-scrollbar" v-html="parsedReport"></div>
      </el-dialog>

    </dv-full-screen-container>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, reactive, nextTick, computed } from 'vue'
import * as echarts from 'echarts'
import dayjs from 'dayjs'
import { useRouter } from 'vue-router'
import { getDashboardStats, getAIReport } from '@/api/admin'
import { getGreenPointsLeaderboard } from '@/api/greenPoints'
import { HomeFilled, User, ShoppingBag, Van, Money, FullScreen } from '@element-plus/icons-vue'
import { useUserStore } from '@/stores/user'

const router = useRouter()
const userStore = useUserStore()
const isAdmin = ref(userStore.userInfo?.role === 'admin')

const currentTime = ref(dayjs().format('YYYY-MM-DD HH:mm:ss'))
let clockTimer = null

// 核心数据状态
const stats = ref({})
const aiReport = ref(null)

// 排行榜控制
const rankingView = ref('datav')
const leaderboardList = ref([])

// 对话框控制变量
const dialogVisible = ref(false)

// ================== AI 报告文本美化解析 ==================
const parsedReport = computed(() => {
  if (!aiReport.value?.report) return '<div style="color:#a0cfff;">暂无详细数据</div>'
  
  let html = aiReport.value.report
  
  // 1. 替换标题 (## / ###)
  html = html.replace(/### (.*)/g, '<h4 class="md-title">$1</h4>')
  html = html.replace(/## (.*)/g, '<h3 class="md-title">$1</h3>')
  
  // 2. 替换加粗文本 (**文字**) - 用醒目颜色突出
  html = html.replace(/\*\*(.*?)\*\*/g, '<span class="md-bold">$1</span>')
  
  // 3. 替换列表项 (- 文字)
  html = html.replace(/^- (.*)/gm, '<div class="md-list-item"><span class="md-dot">•</span> <span class="md-text">$1</span></div>')
  
  // 4. 处理多余的换行，保留正常的段落间隔
  html = html.replace(/\n/g, '<br/>')
  html = html.replace(/<\/h3><br\/>/g, '</h3>')
  html = html.replace(/<\/h4><br\/>/g, '</h4>')
  html = html.replace(/<\/div><br\/>/g, '</div>')
  
  return html
})
// =======================================================


// ECharts 实例与 DOM Ref
const pieChartRef = ref(null)
const lineChartRef = ref(null)
const barChartRef = ref(null)
let pieChart = null
let lineChart = null
let barChart = null

// DataV 翻牌器配置
const flopIncomeConfig = reactive({
  number: [0],
  content: '¥ {nt}',
  style: { fontSize: 36, fill: '#00f2fe', fontWeight: 'bold' }
})

const flopUserConfig = reactive({
  number: [0],
  style: { fontSize: 36, fill: '#00f2fe', fontWeight: 'bold' } // 修改为统一的蓝色
})

// DataV 轮播排行配置
const rankingBoardConfig = reactive({
  data: [],
  rowNum: 6,
  waitTime: 3000,
  carousel: 'single',
  unit: '分'
})

// 初始化与拉取数据
const fetchAllData = async () => {
  try {
    const [dashboardRes, leaderboardRes] = await Promise.all([
      getDashboardStats(),
      getGreenPointsLeaderboard({ limit: 15 })
    ])
    
    stats.value = dashboardRes || {}
    
    // 更新翻牌器数据 (触发深度监听)
    flopIncomeConfig.number = [parseFloat(stats.value.monthIncome || 0)]
    flopIncomeConfig.number = [...flopIncomeConfig.number]
    
    flopUserConfig.number = [stats.value.totalUsers || 0]
    flopUserConfig.number = [...flopUserConfig.number]

    // 更新排行榜数据
    if (leaderboardRes && leaderboardRes.list) {
      leaderboardList.value = leaderboardRes.list
      rankingBoardConfig.data = leaderboardRes.list.map(item => ({
        name: item.nickname || item.username || `用户${item.user_id}`,
        value: item.points || 0
      }))
    }

    // AI 诊断
    if (isAdmin.value) {
      aiReport.value = await getAIReport()
    }

    await nextTick()
    renderCharts()
  } catch (error) {
    console.error('获取大屏数据失败:', error)
  }
}

// 渲染所有图表
const renderCharts = () => {
  // 1. 饼图：报修类型
  if (pieChart) pieChart.dispose()
  pieChart = echarts.init(pieChartRef.value)
  
  // 巧妙处理饼图数据过多拥挤的问题：取前 4 名，剩下的归为"其他"
  let rawPieData = (stats.value.repairStats && stats.value.repairStats.length > 0) 
      ? JSON.parse(JSON.stringify(stats.value.repairStats)) 
      : [{ name: '暂无数据', value: 0 }]

  let pieData = rawPieData
  if (rawPieData.length > 5 && rawPieData[0].name !== '暂无数据') {
    rawPieData.sort((a, b) => b.value - a.value)
    const topData = rawPieData.slice(0, 4)
    const othersValue = rawPieData.slice(4).reduce((sum, item) => sum + item.value, 0)
    pieData = [...topData, { name: '其他', value: othersValue }]
  }
      
  pieChart.setOption({
    // 使用暗色系蓝色调，避免刺眼
    color: ['#85a5ff', '#0A82A4', '#3498db4c', '#2e86c1', '#3498db', '#5dade2'],
    tooltip: { trigger: 'item', backgroundColor: 'rgba(0,0,0,0.7)', textStyle: { color: '#fff' } },
    legend: { bottom: '0%', itemWidth: 10, itemHeight: 10, textStyle: { color: '#a0cfff' } },
    series: [{
      type: 'pie',
      radius: ['40%', '60%'],
      center: ['50%', '42%'],
      avoidLabelOverlap: false,
      itemStyle: { borderColor: '#050a15', borderWidth: 2 },
      label: { show: false },
      data: pieData
    }]
  })

  // 2. 折线图：营收趋势
  if (lineChart) lineChart.dispose()
  lineChart = echarts.init(lineChartRef.value)
  lineChart.setOption({
    tooltip: { trigger: 'axis', backgroundColor: 'rgba(0,0,0,0.7)', textStyle: { color: '#fff' } },
    grid: { top: '15%', left: '3%', right: '4%', bottom: '5%', containLabel: true },
    xAxis: {
      type: 'category',
      boundaryGap: false,
      data: stats.value.incomeDates || [],
      axisLabel: { color: '#a0cfff' },
      axisLine: { lineStyle: { color: 'rgba(255,255,255,0.1)' } }
    },
    yAxis: {
      type: 'value',
      splitLine: { lineStyle: { color: 'rgba(255,255,255,0.05)', type: 'dashed' } },
      axisLabel: { color: '#a0cfff' }
    },
    series: [{
      data: stats.value.incomeTrend || [],
      type: 'line',
      smooth: true,
      symbol: 'none',
      itemStyle: { color: '#00f2fe' },
      lineStyle: { width: 3, shadowColor: 'rgba(0, 242, 254, 0.5)', shadowBlur: 10 },
      areaStyle: {
        color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
          { offset: 0, color: 'rgba(0, 242, 254, 0.4)' },
          { offset: 1, color: 'rgba(0, 242, 254, 0.0)' }
        ])
      }
    }]
  })

  // 3. 柱状图：费用构成分析 (增加的视觉展示数据)
  if (barChart) barChart.dispose()
  barChart = echarts.init(barChartRef.value)
  // 如果后台没传costStructure，使用Mock数据保证视觉效果
  const barData = stats.value.costStructure || [3200, 1800, 4500, 1200]
  barChart.setOption({
    tooltip: { trigger: 'axis', backgroundColor: 'rgba(0,0,0,0.7)', textStyle: { color: '#fff' } },
    grid: { top: '15%', left: '3%', right: '4%', bottom: '5%', containLabel: true },
    xAxis: {
      type: 'category',
      data: ['物业费', '停车费', '商城消费', '场馆预约'],
      axisLabel: { color: '#a0cfff' },
      axisLine: { lineStyle: { color: 'rgba(255,255,255,0.1)' } }
    },
    yAxis: {
      type: 'value',
      splitLine: { lineStyle: { color: 'rgba(255,255,255,0.05)', type: 'dashed' } },
      axisLabel: { color: '#a0cfff' }
    },
    series: [{
      type: 'bar',
      barWidth: '35%',
      itemStyle: {
        color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
          { offset: 0, color: '#4facfe' },
          { offset: 1, color: '#00f2fe' }
        ]),
        borderRadius: [4, 4, 0, 0]
      },
      data: barData
    }]
  })
}

// 工具函数
const formatAmount = (val) => Number(val || 0).toFixed(2)

const handleResize = () => {
  pieChart?.resize()
  lineChart?.resize()
  barChart?.resize()
}

const goBack = () => {
  router.push('/home')
}

onMounted(() => {
  clockTimer = setInterval(() => {
    currentTime.value = dayjs().format('YYYY-MM-DD HH:mm:ss')
  }, 1000)
  
  fetchAllData()
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  clearInterval(clockTimer)
  window.removeEventListener('resize', handleResize)
  pieChart?.dispose()
  lineChart?.dispose()
  barChart?.dispose()
})
</script>

<style scoped>
/* ================== 根级布局 ================== */
.data-screen {
  width: 100vw;
  height: 100vh;
  background-color: #050a15; /* 极深蓝黑底色 */
  color: #ffffff;
  overflow: hidden;
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
}

/* ================== 顶部 Header ================== */
.screen-header {
  height: 70px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 20px;
  position: relative;
  z-index: 10;
}
.header-center {
  display: flex;
  flex-direction: column;
  align-items: center;
  position: absolute;
  left: 50%;
  transform: translateX(-50%);
  top: 5px;
}
.title-text {
  font-size: 32px;
  font-weight: 900;
  color: #00f2fe;
  margin-top: -38px;
  letter-spacing: 2px;
  text-shadow: 0 0 15px rgba(0, 242, 254, 0.6);
}
.time-text {
  position: absolute;
  right: 340px;
  top: 20px;
  color: #a0cfff;
  font-size: 16px;
  font-family: monospace;
}
.back-btn {
  position: absolute;
  left: 340px;
  top: 20px;
  color: #00f2fe;
  cursor: pointer;
  font-size: 16px;
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 4px 12px;
  border: 1px solid rgba(0, 242, 254, 0.3);
  border-radius: 4px;
  transition: all 0.3s;
  background: rgba(0, 242, 254, 0.05);
}
.back-btn:hover {
  background: rgba(0, 242, 254, 0.2);
  box-shadow: 0 0 10px rgba(0, 242, 254, 0.5);
}

/* ================== 主体 Body 栅格 ================== */
.screen-body {
  height: calc(100vh - 70px);
  padding: 10px 20px 20px;
  display: flex;
  gap: 20px;
  box-sizing: border-box;
}
.column-side {
  width: 26%;
  display: flex;
  flex-direction: column;
}
.column-center {
  width: 48%;
  display: flex;
  flex-direction: column;
}

/* ================== 卡片通用样式 ================== */
.tech-card {
  position: relative;
  box-sizing: border-box;
  padding: 20px 15px 15px; /* 避开 datav 边框遮挡 */
}
.mt-15 { margin-top: 15px; }
/* 使用 calc 精确计算高度，减去边距误差，保证总和正好 100%，防止 Flexbox 互相挤压 */
.box-h-30 { height: calc(30% - 10px); }
.box-h-35 { height: calc(35% - 10px); }

.chart-title {
  position: absolute;
  top: 10px;
  left: 20px;
  font-size: 16px;
  font-weight: bold;
  color: #ffffff;
  display: flex;
  align-items: center;
  z-index: 5;
}
.flex-between {
  width: calc(100% - 40px);
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.title-left {
  display: flex;
  align-items: center;
}
.indicator {
  display: inline-block;
  width: 4px;
  height: 16px;
  background-color: #00f2fe;
  margin-right: 8px;
  box-shadow: 0 0 8px #00f2fe;
}
.chart-container {
  width: 100%;
  height: 100%;
  padding-top: 25px;
  box-sizing: border-box;
}

/* ================== 左侧指标卡 ================== */
.metrics-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 8px; /* 间距缩小 */
  height: 100%;
  padding-top: 30px;
}
.metric-item {
  background: rgba(0, 242, 254, 0.03);
  border: 1px solid rgba(0, 242, 254, 0.1);
  border-radius: 6px;
  display: flex;
  align-items: center;
  padding: 0 10px;
  gap: 8px;
  transition: transform 0.3s;
}
.metric-item:hover {
  background: rgba(0, 242, 254, 0.1);
  transform: translateY(-2px);
}
.m-icon {
  font-size: 22px; 
  padding: 6px; 
  background: rgba(255, 255, 255, 0.05);
  border-radius: 8px;
  display: flex;
}
.m-info {
  display: flex;
  flex-direction: column;
  justify-content: center;
}
.m-label {
  font-size: 12px;
  color: #a0cfff;
  margin-bottom: 2px;
}
.m-value {
  font-size: 18px; 
  font-weight: 900;
  color: #ffffff;
  font-family: Arial, sans-serif;
}
.num-small {
  font-size: 16px;
}

/* ================== 中间全景图与翻牌器 ================== */
.center-map-box {
  width: 100%;
  height: 100%;
  position: relative;
  padding: 5px;
}
.map-bg {
  position: absolute;
  top: 10px; left: 10px; right: 10px; bottom: 10px;
  background-image: url('@/assets/images/city.png');
  background-size: cover;
  background-position: center;
  opacity: 0.6; /* 调暗背景以便凸显数据 */
  border-radius: 8px;
}
.map-mask {
  position: absolute;
  top: 10px; left: 10px; right: 10px; bottom: 10px;
  background: radial-gradient(circle, rgba(5,10,21,0) 0%, rgba(5,10,21,0.8) 100%);
  pointer-events: none;
}
.map-title-overlay {
  position: absolute;
  top: 40px;
  width: 100%;
  text-align: center;
  font-size: 22px;
  font-weight: bold;
  color: #e0f2fe;
  letter-spacing: 4px;
  text-shadow: 0 0 15px #00f2fe;
  z-index: 2;
}

.center-data {
  position: absolute;
  bottom: 80px;
  width: 100%;
  display: flex;
  justify-content: space-evenly;
  z-index: 5;
}
.c-item {
  background: rgba(0, 15, 30, 0.6);
  border: 1px solid rgba(0, 242, 254, 0.3);
  box-shadow: 0 0 20px rgba(0, 242, 254, 0.1) inset;
  padding: 15px 30px;
  border-radius: 8px;
  display: flex;
  flex-direction: column;
  align-items: center;
  backdrop-filter: blur(4px);
}
.c-label {
  font-size: 16px;
  color: #a0cfff;
  margin-bottom: 8px;
}

/* 模拟雷达扫描动画 */
.radar-scan {
  position: absolute;
  top: 50%;
  left: 50%;
  width: 400px;
  height: 400px;
  margin-top: -200px;
  margin-left: -200px;
  border-radius: 50%;
  border: 1px dashed rgba(0, 242, 254, 0.2);
  background: conic-gradient(from 0deg, transparent 70%, rgba(0, 242, 254, 0.3) 100%);
  animation: scan 4s linear infinite;
  pointer-events: none;
  z-index: 1;
}
@keyframes scan {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

/* ================== 右侧特有内容 ================== */
.ranking-wrap {
  position: absolute;
  top: 55px; /* 固定位置，避开顶部的标题和单选按钮 */
  left: 15px;
  right: 15px;
  bottom: 15px;
  overflow: hidden;
}
.empty-data {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100%;
  color: #606266;
  font-size: 14px;
}

/* AI 诊断框 */
.ai-report-wrap {
  height: 100%;
  display: flex;
  flex-direction: column;
  padding-top: 25px;
  overflow: hidden; /* 保证内部文字绝不会撑开父盒子 */
}
.ai-status {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 6px;
}
.status-left {
  display: flex;
  align-items: center;
  gap: 5px;
  font-size: 14px;
  color: #4facfe;
  padding-top: 12px;
  
}
.status-dot {
  width: 8px;
  height: 8px;
  background-color: #4facfe;
  border-radius: 50%;
}
.pulse {
  box-shadow: 0 0 0 0 rgba(58, 67, 194, 0.7);
  animation: pulsing 1.5s infinite;
}
@keyframes pulsing {
  0% { transform: scale(0.95); box-shadow: 0 0 0 0 rgba(103, 194, 58, 0.7); }
  70% { transform: scale(1); box-shadow: 0 0 0 6px rgba(103, 194, 58, 0); }
  100% { transform: scale(0.95); box-shadow: 0 0 0 0 rgba(103, 194, 58, 0); }
}

/* 展开按钮样式 */
.expand-btn {
  cursor: pointer;
  color: #4facfe;
  font-size: 13px;
  display: flex;
  align-items: center;
  gap: 4px;
  transition: all 0.3s;
}
.expand-btn:hover {
  color: #00f2fe;
  text-shadow: 0 0 8px rgba(0, 242, 254, 0.6);
}

.report-content {
  flex: 1;
  background: rgba(0, 0, 0, 0.2);
  border: 1px solid rgba(255, 255, 255, 0.05);
  border-radius: 4px;
  padding: 12px;
  overflow-y: auto;
}
.ai-text {
  font-size: 14px;
  line-height: 1.6;
  color: #dcdfe6;
  white-space: pre-wrap;
  margin: 0 0 12px 0;
}
.line-clamp {
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.ai-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}
.tag {
  font-size: 12px;
  padding: 4px 8px;
  border-radius: 4px;
  background: rgba(255, 255, 255, 0.05);
}
.tag.danger { color: #b2d6eb; border: 1px solid rgba(131, 108, 245, 0.3); }
.tag.success { color: #b2d6eb; border: 1px solid rgba(131, 108, 245, 0.3); }

/* ================== Element Plus 弹窗暗黑样式深度覆盖 ================== */
/* 必须使用 :global 才能穿透到挂载于 body 的 el-dialog 上 */
:global(.dark-theme-dialog) {
  background-color: #0b172a !important; /* 深空蓝背景 */
  border: 1px solid rgba(0, 242, 254, 0.4) !important;
  box-shadow: 0 0 30px rgba(0, 242, 254, 0.15) !important;
  border-radius: 8px !important;
}
:global(.dark-theme-dialog .el-dialog__header) {
  border-bottom: 1px solid rgba(0, 242, 254, 0.2) !important;
  margin-right: 0 !important;
  padding-bottom: 15px !important;
}
:global(.dark-theme-dialog .el-dialog__title) {
  color: #00f2fe !important;
  font-weight: bold !important;
  letter-spacing: 1px !important;
  text-shadow: 0 0 10px rgba(0, 242, 254, 0.4);
}
:global(.dark-theme-dialog .el-dialog__body) {
  color: #dcdfe6 !important;
  padding: 20px 25px !important;
}
/* 对话框关闭按钮重置 */
:global(.dark-theme-dialog .el-dialog__headerbtn .el-dialog__close) {
  color: #00f2fe !important;
  font-size: 18px;
}
:global(.dark-theme-dialog .el-dialog__headerbtn:hover .el-dialog__close) {
  color: #ffffff !important;
  text-shadow: 0 0 8px #00f2fe;
}

/* 对话框内部容器 */
.dialog-report-content {
  max-height: 60vh;
  overflow-y: auto;
  line-height: 1.8;
  font-size: 15px;
  padding: 20px;
  background: rgba(0, 242, 254, 0.05);
  border-radius: 6px;
  border: 1px solid rgba(0, 242, 254, 0.1);
  color: #e4e7ed;
}

/* ================== AI 报告内部 Markdown 解析样式 ================== */
:global(.md-title) {
  color: #00f2fe;
  margin: 18px 0 10px 0;
  font-size: 17px;
  font-weight: bold;
  border-bottom: 1px dashed rgba(0, 242, 254, 0.2);
  padding-bottom: 6px;
}
:global(.md-title:first-child) {
  margin-top: 0;
}
:global(.md-bold) {
  color: #e6a23c; /* 醒目的橙黄色，凸显"立即行动"等词汇 */
  font-weight: bold;
}
:global(.md-list-item) {
  margin-bottom: 8px;
  padding-left: 20px;
  position: relative;
}
:global(.md-dot) {
  position: absolute;
  left: 5px;
  color: #00f2fe;
}
:global(.md-text) {
  color: #a0cfff;
}

/* 滚动条暗黑优化 */
.custom-scrollbar::-webkit-scrollbar { width: 4px; }
.custom-scrollbar::-webkit-scrollbar-track { background: rgba(0,0,0,0.1); }
.custom-scrollbar::-webkit-scrollbar-thumb { background: rgba(0, 242, 254, 0.3); border-radius: 2px; }
.custom-scrollbar::-webkit-scrollbar-thumb:hover { background: rgba(0, 242, 254, 0.6); }

/* DataV 深色模式覆盖覆盖 */
:deep(.dv-scroll-ranking-board .ranking-info .rank) { color: #00f2fe; }
:deep(.dv-scroll-ranking-board .ranking-info .info-name) { color: #fff; }

/* ================== 单选按钮组 暗黑样式 ================== */
:deep(.dark-radio .el-radio-button__inner) {
  background: rgba(0, 242, 254, 0.05);
  border: 1px solid rgba(0, 242, 254, 0.3);
  color: #a0cfff;
  border-radius: 4px;
  margin: 0 4px;
  box-shadow: none !important;
  padding: 4px 10px;
  font-size: 12px;
}
:deep(.dark-radio .el-radio-button__original-radio:checked + .el-radio-button__inner) {
  background: rgba(0, 242, 254, 0.2);
  color: #00f2fe;
  border-color: #00f2fe;
  text-shadow: 0 0 8px rgba(0, 242, 254, 0.6);
}

/* ================== 表格视图 暗黑样式 ================== */
:deep(.dark-theme-table) {
  background-color: transparent !important;
  color: #dcdfe6;
  --el-table-border-color: rgba(0, 242, 254, 0.1);
  --el-table-header-bg-color: rgba(0, 242, 254, 0.05);
  --el-table-header-text-color: #00f2fe;
  --el-table-tr-bg-color: transparent;
  --el-table-row-hover-bg-color: rgba(0, 242, 254, 0.1);
}
:deep(.dark-theme-table th.el-table__cell) {
  background-color: rgba(0, 242, 254, 0.05) !important;
  border-bottom: 1px solid rgba(0, 242, 254, 0.2);
  font-weight: bold;
}
:deep(.dark-theme-table td.el-table__cell) {
  border-bottom: 1px dashed rgba(0, 242, 254, 0.1);
}
:deep(.dark-theme-table::before) {
  display: none;
}

/* 前三名徽章效果 */
.rank-badge {
  display: inline-block;
  width: 22px;
  height: 22px;
  line-height: 22px;
  text-align: center;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.1);
  color: #fff;
  font-size: 12px;
  font-weight: bold;
}
.rank-1 { background: #e6a23c; color: #fff; box-shadow: 0 0 8px #e6a23c; }
.rank-2 { background: #a0cfff; color: #fff; box-shadow: 0 0 8px #a0cfff; }
.rank-3 { background: #e0a370; color: #fff; box-shadow: 0 0 8px #e0a370; }
</style>