<template>
  <div class="floor-detail-page">
    <!-- 面包屑导航 -->
    <div class="breadcrumb-bar">
      <span class="breadcrumb-item" @click="goBack">
        <svg viewBox="0 0 24 24" width="16" height="16" fill="currentColor">
          <path d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6"/>
        </svg>
        数据大屏
      </span>
      <span class="breadcrumb-sep">›</span>
      <span class="breadcrumb-item active">楼层详情</span>
    </div>

    <!-- 标题 -->
    <div class="page-title">
      <div class="title-glow">社区楼栋 · 三维楼层分布</div>
      <div class="title-sub">点击楼层可查看详细信息 · 共 {{ FLOOR_COUNT }} 层</div>
    </div>

    <!-- Three.js 画布容器 -->
    <div ref="canvasRef" class="three-canvas"></div>

    <!-- 楼层信息浮窗 -->
    <transition name="fade">
      <div v-if="hoveredFloor !== null" class="floor-tooltip" :style="tooltipStyle">
        <div class="tooltip-title">第 {{ hoveredFloor + 1 }} 层</div>
        <div class="tooltip-row"><span>住户数</span><span>{{ floorData[hoveredFloor].residents }} 户</span></div>
        <div class="tooltip-row"><span>出勤率</span><span>{{ floorData[hoveredFloor].rate }}%</span></div>
        <div class="tooltip-row"><span>报修</span><span>{{ floorData[hoveredFloor].repairs }} 件</span></div>
      </div>
    </transition>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, reactive } from 'vue'
import { useRouter } from 'vue-router'
import * as THREE from 'three'
import { GLTFLoader } from 'three/addons/loaders/GLTFLoader.js'
import { OrbitControls } from 'three/addons/controls/OrbitControls.js'

const router = useRouter()
const canvasRef = ref(null)

const FLOOR_COUNT = 5
const FLOOR_SPACING = 2.8   // 楼层间距（Three.js 单位）

// 模拟各楼层数据
const floorData = [
  { residents: 32, rate: 91, repairs: 2 },
  { residents: 28, rate: 85, repairs: 1 },
  { residents: 35, rate: 94, repairs: 0 },
  { residents: 30, rate: 88, repairs: 3 },
  { residents: 25, rate: 79, repairs: 1 },
]

const hoveredFloor = ref(null)
const tooltipStyle = reactive({ left: '0px', top: '0px' })

let renderer, scene, camera, controls, animFrameId
let floorMeshGroups = []  // 每一层的 mesh group
const raycaster = new THREE.Raycaster()
const mouse = new THREE.Vector2()

const goBack = () => router.push('/data')

const initThree = () => {
  const container = canvasRef.value
  if (!container) return

  const w = container.clientWidth
  const h = container.clientHeight

  // 场景
  scene = new THREE.Scene()
  scene.background = new THREE.Color(0x050a15)
  scene.fog = new THREE.FogExp2(0x050a15, 0.04)

  // 相机
  camera = new THREE.PerspectiveCamera(45, w / h, 0.1, 500)
  camera.position.set(12, FLOOR_SPACING * FLOOR_COUNT * 0.5, 14)
  camera.lookAt(0, FLOOR_SPACING * FLOOR_COUNT * 0.3, 0)

  // 渲染器
  renderer = new THREE.WebGLRenderer({ antialias: true })
  renderer.setSize(w, h)
  renderer.setPixelRatio(window.devicePixelRatio)
  renderer.outputColorSpace = THREE.SRGBColorSpace
  renderer.shadowMap.enabled = true
  renderer.shadowMap.type = THREE.PCFSoftShadowMap
  container.appendChild(renderer.domElement)

  // 轨道控制器
  controls = new OrbitControls(camera, renderer.domElement)
  controls.enableDamping = true
  controls.dampingFactor = 0.06
  controls.autoRotate = false
  controls.enablePan = true
  controls.minDistance = 5
  controls.maxDistance = 60

  // 灯光
  const ambient = new THREE.AmbientLight(0x102040, 3)
  scene.add(ambient)

  const sunLight = new THREE.DirectionalLight(0x00f2fe, 4)
  sunLight.position.set(10, 20, 10)
  sunLight.castShadow = true
  sunLight.shadow.mapSize.width = 2048
  sunLight.shadow.mapSize.height = 2048
  scene.add(sunLight)

  const fillLight = new THREE.DirectionalLight(0x4facfe, 2)
  fillLight.position.set(-8, 5, -8)
  scene.add(fillLight)

  // 地面网格光效
  const gridHelper = new THREE.GridHelper(40, 40, 0x00f2fe, 0x0a1a2e)
  gridHelper.position.y = -0.3
  scene.add(gridHelper)

  // 加载 5 个楼层模型，上下叠放
  const loader = new GLTFLoader()
  let loadedCount = 0

  for (let i = 0; i < FLOOR_COUNT; i++) {
    loader.load(
      '/楼层.glb',
      (gltf) => {
        const model = gltf.scene

        // 计算单次边界盒用于缩放
        const box = new THREE.Box3().setFromObject(model)
        const size = box.getSize(new THREE.Vector3())
        const center = box.getCenter(new THREE.Vector3())
        const maxDim = Math.max(size.x, size.z)
        const scale = 4 / maxDim
        model.scale.setScalar(scale)
        // 使模型底部 y=0
        model.position.y = -box.min.y * scale

        // 整层上移到第 i 层高度
        const group = new THREE.Group()
        group.add(model)
        group.position.y = i * FLOOR_SPACING
        group.userData.floorIndex = i

        // 材质：科技蓝色系，每层轻微色差
        model.traverse((child) => {
          if (child.isMesh) {
            child.castShadow = true
            child.receiveShadow = true
            // 保留原材质，叠加蓝色自发光
            if (child.material) {
              const mat = child.material.clone()
              mat.emissive = new THREE.Color(0x003366)
              mat.emissiveIntensity = 0.3 + i * 0.06
              child.material = mat
            }
          }
        })

        scene.add(group)
        floorMeshGroups.push(group)

        loadedCount++
        // 所有楼层加载完成后自动居中视角
        if (loadedCount === FLOOR_COUNT) {
          const totalH = (FLOOR_COUNT - 1) * FLOOR_SPACING
          controls.target.set(0, totalH * 0.4, 0)
          controls.update()
        }
      },
      undefined,
      (err) => console.error(`楼层 ${i + 1} 加载失败:`, err)
    )
  }

  // 渲染循环
  const animate = () => {
    animFrameId = requestAnimationFrame(animate)
    controls.update()
    renderer.render(scene, camera)
  }
  animate()
}

// 鼠标移动 → 悬停检测
const onMouseMove = (e) => {
  const container = canvasRef.value
  if (!container || !renderer) return
  const rect = container.getBoundingClientRect()
  mouse.x = ((e.clientX - rect.left) / rect.width) * 2 - 1
  mouse.y = -((e.clientY - rect.top) / rect.height) * 2 + 1

  raycaster.setFromCamera(mouse, camera)
  const allMeshes = []
  floorMeshGroups.forEach(g => g.traverse(c => { if (c.isMesh) allMeshes.push(c) }))
  const hits = raycaster.intersectObjects(allMeshes)
  if (hits.length > 0) {
    // 找到属于哪一层
    let obj = hits[0].object
    while (obj.parent && !obj.parent.userData.floorIndex !== undefined && obj.parent !== scene) {
      obj = obj.parent
    }
    const group = floorMeshGroups.find(g => g === obj || g.getObjectById(hits[0].object.id))
    if (group) {
      hoveredFloor.value = group.userData.floorIndex
      tooltipStyle.left = (e.clientX + 16) + 'px'
      tooltipStyle.top = (e.clientY - 40) + 'px'
    }
    document.body.style.cursor = 'pointer'
  } else {
    hoveredFloor.value = null
    document.body.style.cursor = 'default'
  }
}

const onResize = () => {
  const container = canvasRef.value
  if (!container || !renderer || !camera) return
  const w = container.clientWidth
  const h = container.clientHeight
  camera.aspect = w / h
  camera.updateProjectionMatrix()
  renderer.setSize(w, h)
}

onMounted(() => {
  setTimeout(initThree, 100)
  window.addEventListener('resize', onResize)
  window.addEventListener('mousemove', onMouseMove)
})

onUnmounted(() => {
  window.removeEventListener('resize', onResize)
  window.removeEventListener('mousemove', onMouseMove)
  document.body.style.cursor = 'default'
  if (animFrameId) cancelAnimationFrame(animFrameId)
  if (controls) controls.dispose()
  if (renderer) {
    renderer.dispose()
    renderer.domElement.remove()
  }
})
</script>

<style scoped>
.floor-detail-page {
  width: 100vw;
  height: 100vh;
  background: #050a15;
  color: #fff;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
}

/* ====== 面包屑 ====== */
.breadcrumb-bar {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 14px 28px;
  background: rgba(0, 15, 30, 0.9);
  border-bottom: 1px solid rgba(0, 242, 254, 0.15);
  z-index: 10;
  flex-shrink: 0;
}
.breadcrumb-item {
  display: flex;
  align-items: center;
  gap: 5px;
  font-size: 14px;
  color: #4facfe;
  cursor: pointer;
  transition: color 0.2s;
}
.breadcrumb-item:hover { color: #00f2fe; }
.breadcrumb-item.active {
  color: #e0f2fe;
  cursor: default;
  font-weight: bold;
}
.breadcrumb-sep {
  color: rgba(255, 255, 255, 0.3);
  font-size: 16px;
}

/* ====== 标题 ====== */
.page-title {
  text-align: center;
  padding: 18px 0 10px;
  flex-shrink: 0;
}
.title-glow {
  font-size: 26px;
  font-weight: 900;
  color: #00f2fe;
  letter-spacing: 3px;
  text-shadow: 0 0 20px rgba(0, 242, 254, 0.7);
}
.title-sub {
  margin-top: 6px;
  font-size: 13px;
  color: #4facfe;
  opacity: 0.8;
}

/* ====== 画布 ====== */
.three-canvas {
  flex: 1;
  width: 100%;
  min-height: 0;
  position: relative;
}
.three-canvas canvas {
  display: block;
  width: 100% !important;
  height: 100% !important;
}

/* ====== 悬停提示 ====== */
.floor-tooltip {
  position: fixed;
  background: rgba(0, 12, 25, 0.92);
  border: 1px solid rgba(0, 242, 254, 0.4);
  border-radius: 8px;
  padding: 12px 18px;
  min-width: 160px;
  pointer-events: none;
  z-index: 999;
  backdrop-filter: blur(6px);
  box-shadow: 0 0 20px rgba(0, 242, 254, 0.2);
}
.tooltip-title {
  font-size: 15px;
  font-weight: bold;
  color: #00f2fe;
  margin-bottom: 8px;
  letter-spacing: 1px;
}
.tooltip-row {
  display: flex;
  justify-content: space-between;
  font-size: 13px;
  color: #a0cfff;
  padding: 2px 0;
}
.tooltip-row span:last-child {
  color: #e0f2fe;
  font-weight: 500;
}

/* 过渡动画 */
.fade-enter-active, .fade-leave-active { transition: opacity 0.15s; }
.fade-enter-from, .fade-leave-to { opacity: 0; }
</style>
