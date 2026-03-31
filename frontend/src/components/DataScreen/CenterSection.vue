<template>
  <div class="column-center">
    <div class="center-map-box">
      <div ref="threeContainerRef" class="three-canvas-container"></div>
      <dv-border-box-10>
        <div class="map-mask"></div>
        <div class="map-title-overlay">社区 3D 态势感知模型</div>
        <div class="center-data">
          <div class="c-item">
            <div class="c-label">当月社区总营收 (元)</div>
            <dv-digital-flop :config="flopIncomeConfig" style="width:240px;height:50px;" />
          </div>
          <div class="c-item">
            <div class="c-label">活跃用户基数 (人)</div>
            <dv-digital-flop :config="flopUserConfig" style="width:200px;height:50px;" />
          </div>
        </div>
        <div class="radar-scan"></div>
      </dv-border-box-10>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, watch, onMounted, onUnmounted } from 'vue'
import * as THREE from 'three'
import { GLTFLoader } from 'three/addons/loaders/GLTFLoader.js'
import { OrbitControls } from 'three/addons/controls/OrbitControls.js'

const props = defineProps(['stats'])
const threeContainerRef = ref(null)

const flopIncomeConfig = reactive({ number: [0], content: '¥ {nt}', style: { fontSize: 36, fill: '#00f2fe' }})
const flopUserConfig = reactive({ number: [0], style: { fontSize: 36, fill: '#00f2fe' }})

watch(() => props.stats, (newVal) => {
  flopIncomeConfig.number = [parseFloat(newVal.monthIncome || 0)]
  flopUserConfig.number = [newVal.totalUsers || 0]
}, { deep: true })

// 实现 initThreeJS() 逻辑 (代码同原文件)

onMounted(() => setTimeout(initThreeJS, 800))
</script>

<style scoped>
.column-center { width: 48%; display: flex; flex-direction: column; }
.center-map-box { width: 100%; height: 100%; position: relative; }
.three-canvas-container { position: absolute; inset: 10px; z-index: 1; }
@import "@/assets/styles/data.css";
</style>