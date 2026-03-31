<template>
  <div class="column-side">
    <dv-border-box-11 class="tech-card box-h-30" title="实时核心指标">
      <div class="metrics-grid">
        <div class="metric-item" v-for="item in metricList" :key="item.label">
          <div class="m-icon" :style="{color: item.color}"><el-icon><component :is="item.icon"/></el-icon></div>
          <div class="m-info">
            <span class="m-label">{{ item.label }}</span>
            <span class="m-value">{{ item.value }}</span>
          </div>
        </div>
      </div>
    </dv-border-box-11>

    <dv-border-box-13 class="tech-card box-h-35 mt-15">
      <div class="chart-title"><span class="indicator"></span>工单问题分类占比</div>
      <div ref="pieChartRef" class="chart-container"></div>
    </dv-border-box-13>

    <dv-border-box-13 class="tech-card box-h-35 mt-15">
      <div class="chart-title"><span class="indicator"></span>7日营收趋势分析</div>
      <div ref="lineChartRef" class="chart-container"></div>
    </dv-border-box-13>
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue'
import * as echarts from 'echarts'
import { User, ShoppingBag, Van, Money } from '@element-plus/icons-vue'

const props = defineProps(['stats'])
const pieChartRef = ref(null)
const lineChartRef = ref(null)

const metricList = computed(() => [
  { label: '总注册用户', value: props.stats.totalUsers || 0, icon: User, color: '#00f2fe' },
  { label: '今日新增订单', value: props.stats.todayOrders || 0, icon: ShoppingBag, color: '#4facfe' },
  { label: '车位占用率', value: props.stats.parkingRate || '0%', icon: Van, color: '#409EFF' },
  { label: '月营收', value: `¥${props.stats.monthIncome || 0}`, icon: Money, color: '#79bbff' }
])

const initCharts = () => {
  // 此处实现具体的 echarts.init 和 setOption 逻辑 (代码同原文件)
}

watch(() => props.stats, () => initCharts(), { deep: true })
onMounted(() => initCharts())
</script>

<style scoped>
@import "@/assets/styles/data.css";
.column-side { width: 26%; display: flex; flex-direction: column; }
.tech-card { position: relative; padding: 20px 15px 15px; }
.box-h-30 { height: calc(30% - 10px); }
.box-h-35 { height: calc(35% - 10px); }
.mt-15 { margin-top: 15px; }
.metrics-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 8px; padding-top: 30px; }
/* ... 其他样式同原文件 ... */
</style>