<template>
  <div class="column-side">
    <dv-border-box-13 class="tech-card box-h-35">
      <div class="chart-title flex-between">
        <span>社区环保积分先锋榜</span>
        <el-radio-group v-model="rankingView" size="small" class="dark-radio">
          <el-radio-button label="datav">动态</el-radio-button>
          <el-radio-button label="table">表格</el-radio-button>
        </el-radio-group>
      </div>
      <div class="ranking-wrap">
        <dv-scroll-ranking-board v-if="rankingView === 'datav'" :config="rankingBoardConfig" />
        <el-table v-else :data="leaderboardList" class="dark-theme-table">
            </el-table>
      </div>
    </dv-border-box-13>

    <dv-border-box-13 class="tech-card box-h-35 mt-15">
       <div ref="barChartRef" class="chart-container"></div>
    </dv-border-box-13>

    <dv-border-box-11 class="tech-card box-h-30 mt-15" title="AI 智能预警">
      <div class="ai-status">
        <span>AI 守护引擎运行中</span>
        <div class="expand-btn" @click="dialogVisible = true"><el-icon><FullScreen /></el-icon> 展开</div>
      </div>
      <p class="ai-text">{{ aiReport?.report }}</p>
    </dv-border-box-11>

    <el-dialog v-model="dialogVisible" title="AI 深度报告" class="dark-theme-dialog">
       <div v-html="parsedReport"></div>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { getGreenPointsLeaderboard } from '@/api/greenPoints'

const props = defineProps(['stats', 'aiReport'])
const rankingView = ref('datav')
const dialogVisible = ref(false)
const rankingBoardConfig = reactive({ data: [], unit: '分' })

// 排行榜数据获取与 Chart 渲染逻辑
// 实现 parsedReport 的计算属性

onMounted(() => {
  // 获取排行榜数据等
})
</script>

<style scoped>
@import "@/assets/styles/data.css";
/* ... 其他样式同原文件 ... */
</style>