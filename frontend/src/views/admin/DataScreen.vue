<template>
  <div class="data-screen" :class="{ 'pure-active': pureMode }">
    <dv-full-screen-container>
      <!-- ===== 全屏 Three.js 底层背景 (始终铺满整个屏幕) ===== -->
      <div ref="threeContainerRef" class="three-fullscreen"></div>

      <!-- ===== 中间叠加层：雷达/数字/遮罩 (pointer-events: none) ===== -->
      <div class="center-overlay">
        <div class="hud-container">
          <!-- 4个角落的科技装饰线 -->
          <div class="hud-corner top-left"></div>
          <div class="hud-corner top-right"></div>
          <div class="hud-corner bottom-left"></div>
          <div class="hud-corner bottom-right"></div>

          <div class="map-mask"></div>

          <!-- 顶部 3D 态势感知 HUD 标题
          <div class="hud-title-container">
            <div class="hud-title-line"></div>
            <div class="hud-title-badge">
              <span class="pulse-dot-cyan"></span>
              <span class="hud-title-text">{{ viewMode === "floors" ? "楼层 3D 空间分布" : "社区 3D 态势感知模型" }}</span>
            </div>
            <div class="hud-title-sub">REAL-TIME SPATIAL DIAGNOSTIC HUD</div>
          </div> -->

          <!-- 中间浮动的营收和活跃用户数据 -->
          <div class="center-data">
            <div class="glass-metric-card">
              <div class="c-label-container">
                <span class="c-label-icon">
                  <el-icon><Money /></el-icon>
                </span>
                <span class="c-label-text">当月社区总营收 (元)</span>
              </div>
              <div class="c-num">
                <dv-digital-flop
                  :config="flopIncomeConfig"
                  style="width: 200px; height: 50px"
                />
              </div>
              <div class="c-meta-info">
                <span class="c-trend-up">↑ 12.4%</span>
                <span class="c-trend-desc">较上月同期</span>
              </div>
            </div>

            <div class="glass-metric-card">
              <div class="c-label-container">
                <span class="c-label-icon">
                  <el-icon><User /></el-icon>
                </span>
                <span class="c-label-text">活跃用户基数 (人)</span>
              </div>
              <div class="c-num">
                <dv-digital-flop
                  :config="flopUserConfig"
                  style="width: 200px; height: 50px"
                />
              </div>
              <div class="c-meta-info">
                <span class="c-pulse-indicator"></span>
                <span class="c-trend-desc"
                  >当前在线:
                  {{ Math.floor(stats.totalUsers * 0.15) || 120 }}人</span
                >
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 楼层返回按钮 (需要 pointer-events) -->
      <div
        v-if="viewMode === 'floors'"
        class="floor-back-btn"
        @click="switchToBuilding"
      >
        ← 返回建筑总览
      </div>

      <!-- 悬浮 Tooltip -->
      <div
        v-show="hoverTooltip.visible"
        class="floor-tooltip"
        :style="{ left: hoverTooltip.x + 'px', top: hoverTooltip.y + 'px' }"
      >
        <div class="tt-title">{{ hoverTooltip.title }}</div>
        <div class="tt-item">
          住户数：<span class="hl">{{ hoverTooltip.residents }}</span> 人
        </div>
        <div class="tt-item">
          出勤率：<span class="hl">{{ hoverTooltip.attendance }}</span>
        </div>
        <div class="tt-item">
          报修数：<span class="hl">{{ hoverTooltip.repairs }}</span> 件
        </div>
      </div>

      <!-- ===== 顶部 Header (浮动，纯净模式上滑隐藏) ===== -->
      <div class="screen-header">
        <dv-decoration-8 style="width: 300px; height: 50px" />
        <div class="header-center">
          <div class="title-text" style="margin-top: 10px">
            智慧社区数据中枢大屏
          </div>
        </div>
        <dv-decoration-8 :reverse="true" style="width: 300px; height: 50px" />
        <div class="time-text">{{ currentTime }}</div>
        <div class="back-btn" @click="goBack">
          <el-icon><HomeFilled /></el-icon> 首页
        </div>
      </div>

      <!-- ===== 左侧抽屉面板 (纯净模式向左滑出) ===== -->
      <div class="drawer-panel drawer-left">
        <!-- 卡牌 1：实时核心指标 -->
        <div class="glass-card box-h-30">
          <div class="glass-card-header">
            <div class="glass-card-accent"></div>
            <span class="glass-card-title">实时核心指标</span>
            <div class="glass-card-badge-dot"></div>
          </div>
          <div class="metrics-grid">
            <div class="metric-item">
              <div class="m-icon" style="color: #4facfe">
                <el-icon><User /></el-icon>
              </div>
              <div class="m-info">
                <span class="m-label">总注册用户</span>
                <span class="m-value">{{ stats.totalUsers || 0 }}</span>
              </div>
            </div>
            <div class="metric-item">
              <div class="m-icon" style="color: #4facfe">
                <el-icon><ShoppingBag /></el-icon>
              </div>
              <div class="m-info">
                <span class="m-label">今日新增订单</span>
                <span class="m-value">{{ stats.todayOrders || 0 }}</span>
              </div>
            </div>
            <div class="metric-item">
              <div class="m-icon" style="color: #409eff">
                <el-icon><Van /></el-icon>
              </div>
              <div class="m-info">
                <span class="m-label">车位占用率</span>
                <span class="m-value">{{ stats.parkingRate || "0%" }}</span>
              </div>
            </div>
            <div class="metric-item">
              <div class="m-icon" style="color: #79bbff">
                <el-icon><Money /></el-icon>
              </div>
              <div class="m-info">
                <span class="m-label">本月累计营收</span>
                <span class="m-value num-small"
                  >¥{{ formatAmount(stats.monthIncome) }}</span
                >
              </div>
            </div>
          </div>
        </div>

        <!-- 卡牌 2：饲图 -->
        <div class="glass-card box-h-35">
          <div class="glass-card-header">
            <div class="glass-card-accent"></div>
            <span class="glass-card-title">工单问题分类占比</span>
          </div>
          <div
            ref="pieChartRef"
            class="chart-container chart-container-glass"
          ></div>
        </div>

        <!-- 卡牌 3：折线图 -->
        <div class="glass-card box-h-35">
          <div class="glass-card-header">
            <div class="glass-card-accent"></div>
            <span class="glass-card-title">7日营收趋势分析</span>
          </div>
          <div
            ref="lineChartRef"
            class="chart-container chart-container-glass"
          ></div>
        </div>
      </div>

      <!-- ===== 右侧抽屉面板 (纯净模式向右滑出) ===== -->
      <div class="drawer-panel drawer-right">
        <!-- 卡牌 1：积分排行 -->
        <div class="glass-card box-h-35">
          <div class="glass-card-header">
            <div class="glass-card-accent"></div>
            <span class="glass-card-title">社区环保积分先锋榜</span>
            <el-radio-group
              v-model="rankingView"
              size="small"
              class="dark-radio"
              style="margin-left: auto"
            >
              <el-radio-button label="datav">动态展示</el-radio-button>
              <el-radio-button label="table">经典表格</el-radio-button>
            </el-radio-group>
          </div>
          <div class="ranking-wrap ranking-wrap-glass">
            <dv-scroll-ranking-board
              v-if="rankingView === 'datav' && leaderboardList.length"
              :config="rankingBoardConfig"
              :key="leaderboardList.length"
              style="width: 100%; height: 100%"
            />

            <el-table
              v-else-if="rankingView === 'table' && leaderboardList.length"
              :data="leaderboardList"
              class="dark-theme-table custom-scrollbar"
              height="100%"
            >
              <el-table-column label="排名" width="60" align="center">
                <template #default="scope">
                  <span
                    class="rank-badge"
                    :class="'rank-' + (scope.$index + 1)"
                  >
                    {{ scope.$index + 1 }}
                  </span>
                </template>
              </el-table-column>
              <el-table-column label="社区之星" show-overflow-tooltip>
                <template #default="scope">
                  {{
                    scope.row.nickname ||
                    scope.row.username ||
                    `用户${scope.row.user_id}`
                  }}
                </template>
              </el-table-column>
              <el-table-column
                prop="points"
                label="环保积分"
                width="110"
                align="center"
              >
                <template #default="scope">
                  <strong style="color: #00f2fe">{{ scope.row.points }}</strong>
                </template>
              </el-table-column>
            </el-table>

            <div v-else class="empty-data">暂无排行数据</div>
          </div>
        </div>

        <!-- 卡牌 2：柱状图 -->
        <div class="glass-card box-h-35">
          <div class="glass-card-header">
            <div class="glass-card-accent"></div>
            <span class="glass-card-title">社区各模块收入构成</span>
          </div>
          <div
            ref="barChartRef"
            class="chart-container chart-container-glass"
          ></div>
        </div>

        <!-- 卡牌 3： AI 诊断 -->
        <div class="glass-card box-h-30">
          <div class="glass-card-header">
            <div
              class="glass-card-accent"
              style="background: linear-gradient(to right, #4facfe, #00f2fe)"
            ></div>
            <span class="glass-card-title">AI 智能预警与诊断</span>
            <div
              class="status-left"
              style="margin-left: auto; padding-top: 0; font-size: 12px"
            >
              <div class="status-dot pulse"></div>
              <span style="color: #4facfe; font-size: 11px">运行中</span>
            </div>
          </div>
          <div class="ai-report-wrap">
            <div class="ai-status">
              <div class="status-left" style="padding-top: 0">
                <span style="font-size: 13px; color: #a0cfff"
                  >实时诊断社区运营数据</span
                >
              </div>
              <div
                class="status-right expand-btn"
                @click="dialogVisible = true"
              >
                <el-icon><FullScreen /></el-icon> 展开
              </div>
            </div>
            <div class="report-content custom-scrollbar">
              <p class="ai-text line-clamp">
                {{ aiReport?.report_markdown || "正在实时诊断社区运营数据，请稍候..." }}
              </p>
              <div class="ai-tags" v-if="aiReport">
                <span
                  class="tag danger"
                  v-if="aiReport.repair_pending_count > 0"
                >
                  待办报修: {{ aiReport.repair_pending_count }}
                </span>
                <span class="tag success" v-if="aiReport.visitor_new_count">
                  新增访客: {{ aiReport.visitor_new_count }}
                </span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- ===== 纯净模式切换按钮 (始终可见) ===== -->
      <div class="screen-controls">
        <div class="control-btn" @click="togglePureMode">
          <el-icon>
            <Expand v-if="pureMode" />
            <Fold v-else />
          </el-icon>
          <span>{{ pureMode ? "返回大屏" : "纯净模式" }}</span>
        </div>
      </div>

      <!-- ===== 灯光调节控制面板 (只在纯净模式下有 CSS 过渡展示) ===== -->
      <div class="light-control-panel glass-card" :class="{ 'panel-active': pureMode }">
        <div class="glass-card-header">
          <div class="glass-card-accent"></div>
          <span class="glass-card-title">三维光源调节中枢</span>
          <div class="glass-card-badge-dot"></div>
        </div>
        <div class="panel-content">
          <el-tabs v-model="activeLightTab" class="light-tabs">
            <el-tab-pane label="环境光" name="ambient">
              <div class="control-item">
                <span class="control-label">光源强度</span>
                <el-slider v-model="lightConfig.ambient.intensity" :min="0" :max="10" :step="0.1" />
              </div>
              <div class="control-item">
                <span class="control-label">光源颜色</span>
                <el-color-picker v-model="lightConfig.ambient.color" color-format="hex" />
              </div>
            </el-tab-pane>
            <el-tab-pane label="主平行光" name="dir">
              <div class="control-item">
                <span class="control-label">光源强度</span>
                <el-slider v-model="lightConfig.dir.intensity" :min="0" :max="10" :step="0.1" />
              </div>
              <div class="control-item">
                <span class="control-label">光源颜色</span>
                <el-color-picker v-model="lightConfig.dir.color" color-format="hex" />
              </div>
              <div class="control-item">
                <span class="control-label">位置 X</span>
                <el-slider v-model="lightConfig.dir.x" :min="-20" :max="20" :step="0.5" />
              </div>
              <div class="control-item">
                <span class="control-label">位置 Y</span>
                <el-slider v-model="lightConfig.dir.y" :min="-20" :max="20" :step="0.5" />
              </div>
              <div class="control-item">
                <span class="control-label">位置 Z</span>
                <el-slider v-model="lightConfig.dir.z" :min="-20" :max="20" :step="0.5" />
              </div>
            </el-tab-pane>
            <el-tab-pane label="辅助平行" name="fill">
              <div class="control-item">
                <span class="control-label">光源强度</span>
                <el-slider v-model="lightConfig.fill.intensity" :min="0" :max="10" :step="0.1" />
              </div>
              <div class="control-item">
                <span class="control-label">光源颜色</span>
                <el-color-picker v-model="lightConfig.fill.color" color-format="hex" />
              </div>
              <div class="control-item">
                <span class="control-label">位置 X</span>
                <el-slider v-model="lightConfig.fill.x" :min="-20" :max="20" :step="0.5" />
              </div>
              <div class="control-item">
                <span class="control-label">位置 Y</span>
                <el-slider v-model="lightConfig.fill.y" :min="-20" :max="20" :step="0.5" />
              </div>
              <div class="control-item">
                <span class="control-label">位置 Z</span>
                <el-slider v-model="lightConfig.fill.z" :min="-20" :max="20" :step="0.5" />
              </div>
            </el-tab-pane>
            <el-tab-pane label="底部光源" name="ground">
              <div class="control-item">
                <span class="control-label">光源强度</span>
                <el-slider v-model="lightConfig.ground.intensity" :min="0" :max="10" :step="0.1" />
              </div>
              <div class="control-item">
                <span class="control-label">光源颜色</span>
                <el-color-picker v-model="lightConfig.ground.color" color-format="hex" />
              </div>
              <div class="control-item">
                <span class="control-label">照射距离</span>
                <el-slider v-model="lightConfig.ground.distance" :min="0" :max="100" :step="1" />
              </div>
            </el-tab-pane>
          </el-tabs>
        </div>
      </div>
    </dv-full-screen-container>

    <!-- AI 报告弹窗 -->
    <el-dialog
      v-model="dialogVisible"
      title="AI 智能预警与诊断深度报告"
      width="680px"
      class="dark-theme-dialog"
      append-to-body
      destroy-on-close
    >
      <div
        class="dialog-report-content custom-scrollbar"
        v-html="parsedReport"
      ></div>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, reactive, nextTick, computed, watch } from "vue";
import * as echarts from "echarts";
import dayjs from "dayjs";
import { useRouter } from "vue-router";
import { getDashboardStats, getAIReport } from "@/api/admin";
import { getGreenPointsLeaderboard } from "@/api/greenPoints";
import {
  HomeFilled,
  User,
  ShoppingBag,
  Van,
  Money,
  FullScreen,
  Expand,
  Fold,
} from "@element-plus/icons-vue";
import { useUserStore } from "@/stores/user";
import * as THREE from "three";
import { GLTFLoader } from "three/addons/loaders/GLTFLoader.js";
import { DRACOLoader } from "three/addons/loaders/DRACOLoader.js";
import { OrbitControls } from "three/addons/controls/OrbitControls.js";

import { hasPermission } from "@/utils/permission";

const router = useRouter();
const userStore = useUserStore();
const canReadReport = computed(() => hasPermission('statistics:ai_report:read'));

const viewMode = ref("building");

// ===== 纯净模式状态 =====
const pureMode = ref(false);
const togglePureMode = () => {
  pureMode.value = !pureMode.value;
  // 过渡完成后刷新图表尺寸
  setTimeout(() => {
    handleResize();
  }, 650);
};

const hoverTooltip = reactive({
  visible: false,
  x: 0,
  y: 0,
  title: "",
  residents: 0,
  attendance: "0%",
  repairs: 0,
});

// 👉 严格区分两个模式的相机初始位置+距离范围
const switchToFloors = () => {
  viewMode.value = "floors";
  const bGroup = threeScene.getObjectByName("buildingGroup");
  const fGroup = threeScene.getObjectByName("floorsGroup");
  if (bGroup) bGroup.visible = false;
  if (fGroup) fGroup.visible = true;

  // 详情模式：相机拉远、看全楼层
  threeCamera.position.set(12, 5, 6);
  threeControls.minDistance = 6;
  threeControls.maxDistance = 18;
  threeControls.update();
  document.body.style.cursor = "default";
};

const switchToBuilding = () => {
  viewMode.value = "building";
  const bGroup = threeScene.getObjectByName("buildingGroup");
  const fGroup = threeScene.getObjectByName("floorsGroup");
  if (bGroup) bGroup.visible = true;
  if (fGroup) fGroup.visible = false;

  // 总览模式：相机拉近、紧凑视角
  threeCamera.position.set(8, 3.2, 2.2);
  threeControls.minDistance = 2;
  threeControls.maxDistance = 8;
  threeControls.update();
};

const currentTime = ref(dayjs().format("YYYY-MM-DD HH:mm:ss"));
let clockTimer = null;

const stats = ref({});
const aiReport = ref(null);
const rankingView = ref("datav");
const leaderboardList = ref([]);
const dialogVisible = ref(false);

const parsedReport = computed(() => {
  if (!aiReport.value?.report_markdown)
    return '<div style="color:#a0cfff;">暂无详细数据</div>';
  let html = aiReport.value.report_markdown;
  html = html.replace(/### (.*)/g, '<h4 class="md-title">$1</h4>');
  html = html.replace(/## (.*)/g, '<h3 class="md-title">$1</h3>');
  html = html.replace(/\*\*(.*?)\*\*/g, '<span class="md-bold">$1</span>');
  html = html.replace(
    /^- (.*)/gm,
    '<div class="md-list-item"><span class="md-dot">•</span> <span class="md-text">$1</span></div>',
  );
  html = html.replace(/\n/g, "<br/>");
  html = html.replace(/<\/h3><br\/>/g, "</h3>");
  html = html.replace(/<\/h4><br\/>/g, "</h4>");
  html = html.replace(/<\/div><br\/>/g, "</div>");
  return html;
});

const pieChartRef = ref(null);
const lineChartRef = ref(null);
const barChartRef = ref(null);
let pieChart = null;
let lineChart = null;
let barChart = null;

const threeContainerRef = ref(null);
let threeRenderer = null;
let threeScene = null;
let threeCamera = null;
let threeControls = null;
let threeAnimFrameId = null;
let threeInited = false;
let containerRect = null;

// Light references for dynamic adjustment
let ambientLight = null;
let dirLight = null;
let fillLight = null;
let groundLight = null;

const activeLightTab = ref("ambient");

const lightConfig = reactive({
  ambient: {
    color: "#0a1628",
    intensity: 2.5,
  },
  dir: {
    color: "#7fdfff",
    intensity: 2.5,
    x: 5,
    y: 10,
    z: 7,
  },
  fill: {
    color: "#4facfe",
    intensity: 1.5,
    x: -5,
    y: 2,
    z: -5,
  },
  ground: {
    color: "#00f2fe",
    intensity: 1.0,
    distance: 20,
  },
});

// Watch lightConfig changes and update Three.js scene dynamically
watch(
  () => lightConfig.ambient.color,
  (val) => {
    if (ambientLight) ambientLight.color.set(val);
  }
);
watch(
  () => lightConfig.ambient.intensity,
  (val) => {
    if (ambientLight) ambientLight.intensity = val;
  }
);

watch(
  () => lightConfig.dir.color,
  (val) => {
    if (dirLight) dirLight.color.set(val);
  }
);
watch(
  () => lightConfig.dir.intensity,
  (val) => {
    if (dirLight) dirLight.intensity = val;
  }
);
watch(
  () => [lightConfig.dir.x, lightConfig.dir.y, lightConfig.dir.z],
  ([x, y, z]) => {
    if (dirLight) dirLight.position.set(x, y, z);
  }
);

watch(
  () => lightConfig.fill.color,
  (val) => {
    if (fillLight) fillLight.color.set(val);
  }
);
watch(
  () => lightConfig.fill.intensity,
  (val) => {
    if (fillLight) fillLight.intensity = val;
  }
);
watch(
  () => [lightConfig.fill.x, lightConfig.fill.y, lightConfig.fill.z],
  ([x, y, z]) => {
    if (fillLight) fillLight.position.set(x, y, z);
  }
);

watch(
  () => lightConfig.ground.color,
  (val) => {
    if (groundLight) groundLight.color.set(val);
  }
);
watch(
  () => lightConfig.ground.intensity,
  (val) => {
    if (groundLight) groundLight.intensity = val;
  }
);
watch(
  () => lightConfig.ground.distance,
  (val) => {
    if (groundLight) groundLight.distance = val;
  }
);

let threeModelMeshes = []; // 建筑总览所有mesh
let floorModelMeshes = []; // 楼层所有mesh
let hoveredMesh = null; // 当前总览模式下高亮的单个mesh
let hoveredFloorMesh = null; // 当前楼层模式下高亮的单个mesh
let starsMaterial = null; // 星空材质引用

const threeRaycaster = new THREE.Raycaster();
const threeMouse = new THREE.Vector2();

// 无限网格生成函数 (模仿 Drei 的 Grid 组件)
const createInfiniteGrid = () => {
  const geometry = new THREE.PlaneGeometry(2, 2, 1, 1);
  const material = new THREE.ShaderMaterial({
    side: THREE.DoubleSide,
    transparent: true,
    uniforms: {
      uCellSize: { value: 0.3 },
      uCellThickness: { value: 0.6 },
      uCellColor: { value: new THREE.Color("#6f6f6f") },
      uSectionSize: { value: 1.5 },
      uSectionThickness: { value: 1.5 },
      uSectionColor: { value: new THREE.Color("#00f2fe") },
      uFadeDistance: { value: 30.0 },
    },
    vertexShader: `
      varying vec3 worldPosition;
      uniform float uFadeDistance;
      void main() {
        vec3 pos = vec3(position.x, 0.0, position.y) * uFadeDistance;
        pos.xz += cameraPosition.xz;
        worldPosition = pos;
        gl_Position = projectionMatrix * modelViewMatrix * vec4(pos, 1.0);
      }
    `,
    fragmentShader: `
      varying vec3 worldPosition;
      uniform float uCellSize;
      uniform float uCellThickness;
      uniform vec3 uCellColor;
      uniform float uSectionSize;
      uniform float uSectionThickness;
      uniform vec3 uSectionColor;
      uniform float uFadeDistance;

      float getGrid(float size, float thickness) {
        vec2 r = worldPosition.xz / size;
        vec2 grid = abs(fract(r - 0.5) - 0.5) / (fwidth(r) * thickness);
        float line = min(grid.x, grid.y);
        return 1.0 - min(line, 1.0);
      }

      void main() {
        float cell = getGrid(uCellSize, uCellThickness);
        float section = getGrid(uSectionSize, uSectionThickness);
        vec3 color = mix(uCellColor, uSectionColor, section);
        float alpha = max(cell, section);
        float dist = distance(cameraPosition.xz, worldPosition.xz);
        float fade = 1.0 - min(dist / uFadeDistance, 1.0);
        gl_FragColor = vec4(color, alpha * fade);
        if (gl_FragColor.a < 0.01) discard;
      }
    `,
  });

  const mesh = new THREE.Mesh(geometry, material);
  mesh.frustumCulled = false;
  mesh.position.y = -1.0;
  return mesh;
};

// 粒子星空生成函数 (模仿 Drei 的 Stars 组件)
const createStars = (
  radius = 100,
  depth = 50,
  count = 1000,
  saturation = 0,
  factor = 8,
  fade = true,
) => {
  const genStar = (r) => {
    const spherical = new THREE.Spherical(
      r,
      Math.acos(1 - Math.random() * 2),
      Math.random() * 2 * Math.PI,
    );
    return new THREE.Vector3().setFromSpherical(spherical);
  };

  const positions = [];
  const colors = [];
  const sizes = [];
  const colorObj = new THREE.Color();
  let r = radius + depth;
  const increment = depth / count;
  for (let i = 0; i < count; i++) {
    r -= increment * Math.random();
    const starPos = genStar(r);
    positions.push(starPos.x, starPos.y, starPos.z);
    colorObj.setHSL(i / count, saturation, 0.9);
    colors.push(colorObj.r, colorObj.g, colorObj.b);
    sizes.push((0.5 + 0.5 * Math.random()) * factor);
  }

  const geometry = new THREE.BufferGeometry();
  geometry.setAttribute(
    "position",
    new THREE.Float32BufferAttribute(positions, 3),
  );
  geometry.setAttribute("color", new THREE.Float32BufferAttribute(colors, 3));
  geometry.setAttribute("size", new THREE.Float32BufferAttribute(sizes, 1));

  const mat = new THREE.ShaderMaterial({
    uniforms: {
      time: { value: 0.0 },
      fade: { value: fade ? 1.0 : 0.0 },
    },
    vertexShader: `
      uniform float time;
      attribute float size;
      varying vec3 vColor;
      void main() {
        vColor = color;
        vec4 mvPosition = modelViewMatrix * vec4(position, 0.5);
        gl_PointSize = size * (30.0 / -mvPosition.z) * (3.0 + sin(time + 100.0));
        gl_Position = projectionMatrix * mvPosition;
      }
    `,
    fragmentShader: `
      uniform float fade;
      varying vec3 vColor;
      void main() {
        float opacity = 1.0;
        if (fade == 1.0) {
          float d = distance(gl_PointCoord, vec2(0.5, 0.5));
          opacity = 1.0 / (1.0 + exp(16.0 * (d - 0.25)));
        }
        gl_FragColor = vec4(vColor, opacity);
        #include <tonemapping_fragment>
        #include <colorspace_fragment>
      }
    `,
    blending: THREE.AdditiveBlending,
    depthWrite: false,
    transparent: true,
    vertexColors: true,
  });

  const stars = new THREE.Points(geometry, mat);
  return { mesh: stars, material: mat };
};

const initThreeJS = () => {
  const container = threeContainerRef.value;
  if (!container) return;

  containerRect = container.getBoundingClientRect();
  const width =
    containerRect.width || container.offsetWidth || window.innerWidth;
  const height =
    containerRect.height || container.offsetHeight || window.innerHeight;

  threeScene = new THREE.Scene();
  // 设置 3D 场景背景色为深灰色 #26282a
  threeScene.background = new THREE.Color("#26282a");

  threeCamera = new THREE.PerspectiveCamera(45, width / height, 0.1, 1000);
  // 👉 初始默认：总览模式（近）
  threeCamera.position.set(5, 2.5, 4);

  threeRenderer = new THREE.WebGLRenderer({ antialias: true, alpha: false });
  threeRenderer.setSize(width, height);
  threeRenderer.setPixelRatio(window.devicePixelRatio);
  threeRenderer.outputColorSpace = THREE.SRGBColorSpace;
  threeRenderer.shadowMap.enabled = true;
  container.appendChild(threeRenderer.domElement);

  threeControls = new OrbitControls(threeCamera, threeRenderer.domElement);
  threeControls.enableDamping = true;
  threeControls.dampingFactor = 0.05;
  threeControls.autoRotate = false;
  threeControls.enablePan = true;
  // 👉 初始总览模式距离范围
  threeControls.minDistance = 2;
  threeControls.maxDistance = 8;

  // 右视角锁定
  //threeControls.minAzimuthAngle = Math.PI / 3;
  //threeControls.maxAzimuthAngle = (Math.PI * 2) / 3;
  threeControls.minPolarAngle = Math.PI / 4;
  threeControls.maxPolarAngle = Math.PI / 2.5;

  // 灯光
  ambientLight = new THREE.AmbientLight(lightConfig.ambient.color, lightConfig.ambient.intensity);
  threeScene.add(ambientLight);

  dirLight = new THREE.DirectionalLight(lightConfig.dir.color, lightConfig.dir.intensity);
  dirLight.position.set(lightConfig.dir.x, lightConfig.dir.y, lightConfig.dir.z);
  dirLight.castShadow = true;
  threeScene.add(dirLight);

  fillLight = new THREE.DirectionalLight(lightConfig.fill.color, lightConfig.fill.intensity);
  fillLight.position.set(lightConfig.fill.x, lightConfig.fill.y, lightConfig.fill.z);
  threeScene.add(fillLight);

  groundLight = new THREE.PointLight(lightConfig.ground.color, lightConfig.ground.intensity, lightConfig.ground.distance);
  groundLight.position.set(0, -2, 0);
  threeScene.add(groundLight);

  // 添加无限网格
  const grid = createInfiniteGrid();
  threeScene.add(grid);

  // 添加粒子星空背景
  const starsData = createStars(100, 50, 1000, 0, 8, true);
  threeScene.add(starsData.mesh);
  starsMaterial = starsData.material;

  // 分组
  const buildingGroup = new THREE.Group();
  buildingGroup.name = "buildingGroup";
  const floorsGroup = new THREE.Group();
  floorsGroup.name = "floorsGroup";
  floorsGroup.visible = false;
  threeScene.add(buildingGroup);
  threeScene.add(floorsGroup);

  const loader = new GLTFLoader();
  const dracoLoader = new DRACOLoader();
  dracoLoader.setDecoderPath("https://www.gstatic.com/draco/v1/decoders/");
  loader.setDRACOLoader(dracoLoader);

  // 加载建筑总览模型
  loader.load(
    // "https://communitysvc.xyz/community/building/build_compressed.glb",//修改
    "/modelNew.glb",
    (gltf) => {
      const model = gltf.scene;
      const box = new THREE.Box3().setFromObject(model);
      const center = box.getCenter(new THREE.Vector3());
      const size = box.getSize(new THREE.Vector3());
      const maxDim = Math.max(size.x, size.y, size.z);
      const scale = 3 / maxDim;
      model.scale.setScalar(scale);
      model.position.sub(center.multiplyScalar(scale));

      model.traverse((child) => {
        if (child.isMesh) {
          child.castShadow = true;
          child.receiveShadow = true;
          if (child.material) {
            child.material.envMapIntensity = 1.5;
          }
          threeModelMeshes.push(child);
        }
      });
      buildingGroup.add(model);
    },
    undefined,
    (err) => console.error("建模加载失败:", err),
  );

  // 加载楼层模型
  const FLOOR_COUNT = 5;
  const FLOOR_SPACING = 0.9;
  for (let i = 0; i < FLOOR_COUNT; i++) {
    loader.load(
      "/楼层.glb",
      (gltf) => {
        const model = gltf.scene;
        const box = new THREE.Box3().setFromObject(model);
        const maxDim = Math.max(
          box.getSize(new THREE.Vector3()).x,
          box.getSize(new THREE.Vector3()).z,
        );
        const scale = 3.5 / maxDim;
        model.scale.setScalar(scale);
        model.position.y = -box.min.y * scale;

        const group = new THREE.Group();
        group.add(model);
        group.position.y = i * FLOOR_SPACING - 1;
        group.userData.floorIndex = i;

        model.traverse((child) => {
          if (child.isMesh) {
            child.castShadow = true;
            child.receiveShadow = true;
            if (child.material) {
              child.material.emissive = new THREE.Color(0x003366);
              child.material.emissiveIntensity = 0.3 + i * 0.05;//修改蓝色
              // child.material.emissive.setHex(0x3fa9f5);
              // child.material.emissiveIntensity = 0.08; //变浅
              child.userData.originalEmissive = 0x003366;
            }
            floorModelMeshes.push(child);
          }
        });
        floorsGroup.add(group);
      },
      undefined,
      (err) => console.error(`楼层 ${i} 加载失败:`, err),
    );
  }

  // 初始视角
  threeCamera.lookAt(0, 0, 0);
  threeControls.update();

  // 射线检测工具
  const getIntersect = (e, meshes) => {
    if (!containerRect) return [];
    threeMouse.x =
      ((e.clientX - containerRect.left) / containerRect.width) * 2 - 1;
    threeMouse.y =
      -((e.clientY - containerRect.top) / containerRect.height) * 2 + 1;
    threeRaycaster.setFromCamera(threeMouse, threeCamera);
    return threeRaycaster.intersectObjects(meshes);
  };

  // 点击切换模式
  const onThreeClick = (e) => {
    if (viewMode.value === "building") {
      const hits = getIntersect(e, threeModelMeshes);
      if (hits.length > 0) switchToFloors();
    }
  };

  // 👉 辅助函数：向上查找到顶级组
  const getTopParent = (obj, stopName) => {
    let p = obj;
    while (
      p.parent &&
      p.parent.name !== stopName &&
      p.parent.type !== "Scene"
    ) {
      p = p.parent;
    }
    return p;
  };

  // 👉 核心修复：总览模式单建筑高亮，详情模式单楼层高亮
  const onThreeMouseMove = (e) => {
    // 总览模式：单个建筑高亮
    if (viewMode.value === "building") {
      const hits = getIntersect(e, threeModelMeshes);
      const targetObj = hits.length > 0 ? hits[0].object : null;

      if (hoveredMesh && hoveredMesh !== targetObj) {
        if (hoveredMesh.userData.originalMaterial) {
          hoveredMesh.material = hoveredMesh.userData.originalMaterial;
        }
        hoveredMesh = null;
        document.body.style.cursor = "default";
        hoverTooltip.visible = false;
      }

      if (targetObj && hoveredMesh !== targetObj) {
        hoveredMesh = targetObj;
        hoveredMesh.userData.originalMaterial = hoveredMesh.material;
        hoveredMesh.material = hoveredMesh.material.clone();
        hoveredMesh.material.emissive.setHex(0x00f2fe); //修改
        document.body.style.cursor = "pointer";
        hoverTooltip.visible = true;
      }

      if (hoverTooltip.visible && hoveredMesh) {
        hoverTooltip.x = e.clientX + 15;
        hoverTooltip.y = e.clientY + 15;
        hoverTooltip.title = "社区建筑数据";
        hoverTooltip.residents = 100 + Math.floor(Math.random() * 50);
        hoverTooltip.attendance = 85 + Math.floor(Math.random() * 10) + "%";
        hoverTooltip.repairs = Math.floor(Math.random() * 3);
      }
    }

    // 楼层模式：整个楼层高亮
    if (viewMode.value === "floors") {
      const hits = getIntersect(e, floorModelMeshes);
      const targetGroup =
        hits.length > 0 ? getTopParent(hits[0].object, "floorsGroup") : null;

      if (hoveredFloorMesh && hoveredFloorMesh !== targetGroup) {
        hoveredFloorMesh.traverse((child) => {
          if (
            child.isMesh &&
            child.material &&
            child.userData.originalEmissive !== undefined
          ) {
            child.material.emissive.setHex(child.userData.originalEmissive);
          }
        });
        hoveredFloorMesh = null;
        document.body.style.cursor = "default";
        hoverTooltip.visible = false;
      }

      if (targetGroup && hoveredFloorMesh !== targetGroup) {
        hoveredFloorMesh = targetGroup;
        hoveredFloorMesh.traverse((child) => {
          if (child.isMesh && child.material) {
            child.material.emissive.setHex(0x00f2fe); //修改
          }
        });
        document.body.style.cursor = "pointer";
        hoverTooltip.visible = true;
      }

      if (hoverTooltip.visible && hoveredFloorMesh) {
        hoverTooltip.x = e.clientX + 15;
        hoverTooltip.y = e.clientY + 15;
        const floorIndex = hoveredFloorMesh.userData.floorIndex || 0;
        hoverTooltip.title = `第 ${floorIndex + 1} 层 运行数据`;
        hoverTooltip.residents = 120 + floorIndex * 15;
        hoverTooltip.attendance = 92 + floorIndex + "%";
        hoverTooltip.repairs = Math.floor(Math.random() * 5);
      }
    }
  };

  container.addEventListener("click", onThreeClick);
  container.addEventListener("mousemove", onThreeMouseMove);
  container._click = onThreeClick;
  container._move = onThreeMouseMove;

  // 渲染循环
  const threeClock = new THREE.Clock();
  const animate = () => {
    threeAnimFrameId = requestAnimationFrame(animate);
    threeControls.update();
    if (starsMaterial) {
      starsMaterial.uniforms.time.value = threeClock.getElapsedTime() * 2.0;
    }
    threeRenderer.render(threeScene, threeCamera);
  };
  animate();
};

const handleThreeResize = () => {
  if (!threeRenderer || !threeCamera || !threeContainerRef.value) return;
  const w = threeContainerRef.value.clientWidth;
  const h = threeContainerRef.value.clientHeight;
  threeCamera.aspect = w / h;
  threeCamera.updateProjectionMatrix();
  threeRenderer.setSize(w, h);
  if (threeContainerRef.value) {
    containerRect = threeContainerRef.value.getBoundingClientRect();
  }
};

// DataV 翻牌器
const flopIncomeConfig = ref({
  number: [0],
  content: "¥ {nt}",
  style: { fontSize: 36, fill: "#00f2fe", fontWeight: "bold" },
});
const flopUserConfig = ref({
  number: [0],
  style: { fontSize: 36, fill: "#00f2fe", fontWeight: "bold" },
});
const rankingBoardConfig = ref({
  data: [],
  rowNum: 6,
  waitTime: 3000,
  carousel: "single",
  unit: "分",
});

// 数据拉取
const fetchAllData = async () => {
  try {
    const [dashboardRes, leaderboardRes] = await Promise.all([
      getDashboardStats(),
      getGreenPointsLeaderboard({ limit: 15 }),
    ]);
    stats.value = dashboardRes || {};
    flopIncomeConfig.value = {
      ...flopIncomeConfig.value,
      number: [parseFloat(stats.value.monthIncome || 0)],
    };
    flopUserConfig.value = {
      ...flopUserConfig.value,
      number: [stats.value.totalUsers || 0],
    };
    if (leaderboardRes?.list) {
      leaderboardList.value = leaderboardRes.list;
      rankingBoardConfig.value = {
        ...rankingBoardConfig.value,
        data: leaderboardRes.list.map((i) => ({
          name: i.nickname || i.username || `用户${i.user_id}`,
          value: i.points || 0,
        })),
      };
    }
    if (canReadReport.value) {
      const res = await getAIReport();
      aiReport.value = res.report;
    }
    await nextTick();
    renderCharts();
  } catch (e) {
    console.error("数据加载失败", e);
  }
};

// 图表渲染
const renderCharts = () => {
  // 饼图
  if (pieChart) pieChart.dispose();
  pieChart = echarts.init(pieChartRef.value);
  let rawPieData = stats.value.repairStats?.length
    ? JSON.parse(JSON.stringify(stats.value.repairStats))
    : [{ name: "暂无数据", value: 0 }];
  let pieData = rawPieData;
  if (rawPieData.length > 5 && rawPieData[0].name !== "暂无数据") {
    rawPieData.sort((a, b) => b.value - a.value);
    pieData = [
      ...rawPieData.slice(0, 4),
      {
        name: "其他",
        value: rawPieData.slice(4).reduce((s, i) => s + i.value, 0),
      },
    ];
  }
  pieChart.setOption({
    color: ["#00f2fe", "#4facfe", "#79bbff", "#2e86c1", "#3498db", "#5dade2"],
    tooltip: {
      trigger: "item",
      backgroundColor: "rgba(0,0,0,0.7)",
      textStyle: { color: "#fff" },
    },
    legend: {
      bottom: "0%",
      itemWidth: 10,
      itemHeight: 10,
      textStyle: { color: "#a0cfff" },
    },
    series: [
      {
        type: "pie",
        radius: ["40%", "60%"],
        center: ["50%", "42%"],
        itemStyle: {
          borderColor: "#050a15",
          borderWidth: 2,
          borderRadius: 6,
        },
        label: { show: false },
        data: pieData,
      },
    ],
  });

  // 折线图
  if (lineChart) lineChart.dispose();
  lineChart = echarts.init(lineChartRef.value);
  lineChart.setOption({
    tooltip: {
      trigger: "axis",
      backgroundColor: "rgba(0,0,0,0.7)",
      textStyle: { color: "#fff" },
    },
    grid: {
      top: "15%",
      left: "3%",
      right: "4%",
      bottom: "5%",
      containLabel: true,
    },
    xAxis: {
      type: "category",
      data: stats.value.incomeDates || [],
      axisLabel: { color: "#a0cfff" },
      axisLine: { lineStyle: { color: "rgba(255,255,255,0.1)" } },
    },
    yAxis: {
      type: "value",
      splitLine: {
        lineStyle: { color: "rgba(255,255,255,0.05)", type: "dashed" },
      },
      axisLabel: { color: "#a0cfff" },
    },
    series: [
      {
        data: stats.value.incomeTrend || [],
        type: "line",
        smooth: true,
        symbol: "none",
        lineStyle: { width: 3, color: "#00f2fe" },
        areaStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: "rgba(0,242,254,0.4)" },
            { offset: 1, color: "rgba(0,242,254,0)" },
          ]),
        },
      },
    ],
  });

  // 柱状图
  if (barChart) barChart.dispose();
  barChart = echarts.init(barChartRef.value);
  barChart.setOption({
    tooltip: {
      trigger: "axis",
      backgroundColor: "rgba(0,0,0,0.7)",
      textStyle: { color: "#fff" },
    },
    grid: {
      top: "15%",
      left: "3%",
      right: "4%",
      bottom: "5%",
      containLabel: true,
    },
    xAxis: {
      type: "category",
      data: ["物业费", "停车费", "商城消费"],
      axisLabel: { color: "#a0cfff" },
      axisLine: { lineStyle: { color: "rgba(255,255,255,0.1)" } },
    },
    yAxis: {
      type: "value",
      splitLine: {
        lineStyle: { color: "rgba(255,255,255,0.05)", type: "dashed" },
      },
      axisLabel: { color: "#a0cfff" },
    },
    series: [
      {
        type: "bar",
        barWidth: "35%",
        itemStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: "#4facfe" },
            { offset: 1, color: "#00f2fe" },
          ]),
          borderRadius: [4, 4, 0, 0],
        },
        data: stats.value.costStructure || [3200, 1800, 4500],
      },
    ],
  });
};

// 工具函数
const formatAmount = (val) => Number(val || 0).toFixed(2);
const handleResize = () => {
  pieChart?.resize();
  lineChart?.resize();
  barChart?.resize();
};
const goBack = () => router.push("/home");

onMounted(() => {
  document.documentElement.classList.add('data-screen-active');
  document.body.classList.add('data-screen-active');
  clockTimer = setInterval(
    () => (currentTime.value = dayjs().format("YYYY-MM-DD HH:mm:ss")),
    1000,
  );
  fetchAllData();
  setTimeout(() => initThreeJS(), 800);
  window.addEventListener("resize", handleResize);
  window.addEventListener("resize", handleThreeResize);
});

onUnmounted(() => {
  document.documentElement.classList.remove('data-screen-active');
  document.body.classList.remove('data-screen-active');
  clearInterval(clockTimer);
  window.removeEventListener("resize", handleResize);
  window.removeEventListener("resize", handleThreeResize);
  pieChart?.dispose();
  lineChart?.dispose();
  barChart?.dispose();
  if (threeAnimFrameId) cancelAnimationFrame(threeAnimFrameId);
  if (threeControls) threeControls.dispose();
  if (threeRenderer) {
    const c = threeContainerRef.value;
    if (c) {
      c.removeEventListener("click", c._click);
      c.removeEventListener("mousemove", c._move);
    }
    threeRenderer.dispose();
    threeRenderer.domElement.remove();
  }
  starsMaterial = null;
});
</script>

<style scoped>
@import "../../assets/styles/data.css";

/* ===================================================
   全屏 Three.js 容器：始终铺满整个 dv-full-screen-container
   =================================================== */
.three-fullscreen {
  position: absolute;
  inset: 0;
  z-index: 0;
  overflow: hidden;
}
.three-fullscreen canvas {
  display: block;
  width: 100% !important;
  height: 100% !important;
}

/* ===================================================
   顶部 Header：绝对浮动，带渐变遮罩
   =================================================== */
.screen-header {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 70px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 20px;
  z-index: 10;
  background: linear-gradient(
    to bottom,
    rgba(5, 10, 21, 0.92) 0%,
    rgba(5, 10, 21, 0.1) 100%
  );
  /* 纯净模式上滑动画 */
  transition:
    opacity 0.5s cubic-bezier(0.4, 0, 0.2, 1),
    transform 0.5s cubic-bezier(0.4, 0, 0.2, 1);
}

/* ===================================================
   抽屉面板通用样式
   =================================================== */
.drawer-panel {
  position: absolute;
  top: 70px;
  width: 26%;
  height: calc(100% - 250px);
  display: flex;
  flex-direction: column;
  /* gap 替代 mt-15，避免额外高度溢出 */
  gap: 8px;
  padding: 8px 10px 10px;
  z-index: 5;
  box-sizing: border-box;
  overflow: hidden;
  /* 抽屉滑动动画 */
  transition:
    transform 0.6s cubic-bezier(0.4, 0, 0.2, 1),
    opacity 0.6s ease;
}

.drawer-left {
  left: 0;
  background: linear-gradient(
    to right,
    rgba(5, 10, 21, 0.1) 65%,
    rgba(5, 10, 21, 0) 100%
  );
}

.drawer-right {
  right: 0;
  background: linear-gradient(
    to left,
    rgba(5, 10, 21, 0.1) 65%,
    rgba(5, 10, 21, 0) 100%
  );
}

/* 使用 flex 比例替代固定百分比高度，让卡片自动均分可用空间 */
.drawer-panel .box-h-30 {
  flex: 30 1 0;
  height: 0 !important; /* flex-basis=0, flex-grow 控制实际高度 */
  min-height: 0;
}
.drawer-panel .box-h-35 {
  flex: 35 1 0;
  height: 0 !important;
  min-height: 0;
}
/* 消除 mt-15 多余 margin，由 gap 统一控制间距 */
.drawer-panel .mt-15 {
  margin-top: 0 !important;
}

/* ===================================================
   中间叠加层：全屏，pointer-events: none
   =================================================== */
.center-overlay {
  position: absolute;
  inset: 0;
  z-index: 2;
  pointer-events: none;
  transition: opacity 0.5s ease;
}

/* ===================================================
   纯净模式激活：左右抽屉像门一样滑出屏幕
   =================================================== */
.data-screen.pure-active .drawer-left {
  transform: translateX(-110%);
  opacity: 0;
  pointer-events: none;
}
.data-screen.pure-active .drawer-right {
  transform: translateX(110%);
  opacity: 0;
  pointer-events: none;
}
.data-screen.pure-active .screen-header {
  opacity: 0;
  transform: translateY(-100%);
  pointer-events: none;
}
.data-screen.pure-active .center-overlay {
  opacity: 0;
  pointer-events: none;
}
.data-screen.pure-active .screen-controls {
  top: 36px;
  left: 36px;
  bottom: auto;
  transform: none;
}

/* ===================================================
   纯净模式切换按钮（始终可见，始终居中底部）
   =================================================== */
.screen-controls {
  position: absolute;
  bottom: 960px;
  left: 68%;
  transform: translateX(-50%);
  z-index: 20;
  pointer-events: auto;
  display: flex;
  gap: 14px;
}
.control-btn {
  color: #00f2fe;
  cursor: pointer;
  font-size: 14px;
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 9px 22px;
  border: 1px solid rgba(0, 242, 254, 0.35);
  border-radius: 4px;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  background: rgba(5, 10, 21, 0.88);
  backdrop-filter: blur(10px);
  box-shadow:
    0 0 16px rgba(0, 242, 254, 0.18),
    inset 0 0 30px rgba(0, 242, 254, 0.03);
  user-select: none;
  font-weight: bold;
  letter-spacing: 1px;
}
.control-btn:hover {
  background: rgba(0, 242, 254, 0.18);
  box-shadow: 0 0 24px rgba(0, 242, 254, 0.55);
  transform: translateY(-3px);
  color: #ffffff;
  border-color: #00f2fe;
}
.control-btn:active {
  transform: translateY(-1px);
}

/* ===================================================
   楼层返回按钮
   =================================================== */
.floor-back-btn {
  position: absolute;
  top: 82px;
  left: 32%;
  transform: translateX(-50%);
  z-index: 10;
  background: rgba(5, 10, 21, 0.82);
  border: 1px solid rgba(0, 242, 254, 0.4);
  border-radius: 4px;
  padding: 6px 18px;
  color: #00f2fe;
  cursor: pointer;
  font-size: 14px;
  backdrop-filter: blur(6px);
  transition: all 0.3s;
  pointer-events: auto;
}
.floor-back-btn:hover {
  background: rgba(0, 242, 254, 0.2);
  box-shadow: 0 0 12px rgba(0, 242, 254, 0.4);
}

/* ===================================================
   其余样式兼容（data.css 中已有，这里补充覆盖）
   =================================================== */
.data-screen {
  position: fixed;
  width: 100%;
  height: 100%;
  top: 0;
  left: 0;
  background: #26282a;
  overflow: hidden;
  z-index: 999;
}

:global(html.data-screen-active),
:global(body.data-screen-active) {
  margin: 0 !important;
  padding: 0 !important;
  width: 100% !important;
  height: 100% !important;
  overflow: hidden !important;
  background-color: #26282a !important; /* 统一底色为大屏深灰，防止加载或边缘出现白色缝隙 */
}

:global(.data-screen-active #app) {
  max-width: none !important;
  margin: 0 !important;
  padding: 0 !important;
  width: 100% !important;
  height: 100% !important;
  background-color: #26282a !important;
}

/* ===================================================
   毛玻璃卡片：仿照 sc-datav Demo0 Card 样式
   =================================================== */
.glass-card {
  position: relative;
  display: flex;
  flex-direction: column;
  color: #ffffff;
  pointer-events: auto;
  /* 毛玻璃核心 */
  backdrop-filter: blur(14px);
  -webkit-backdrop-filter: blur(20px);
  background: rgba(10, 20, 42, 0.05);
  /* 细边框 */
  border: 1px solid rgba(141, 141, 141, 0.22);
  border-radius: 6px;
  box-sizing: border-box;
  overflow: hidden;
  /* 顶部左上角青色光 */
  box-shadow:
    0 4px 24px rgba(0, 0, 0, 0.35),
    inset 0 1px 0 rgba(255, 255, 255, 0.06),
    0 0 0 0.5px rgba(0, 242, 254, 0.08);
  transition:
    box-shadow 0.3s ease,
    border-color 0.3s ease;
}
.glass-card:hover {
  border-color: rgba(0, 242, 254, 0.28);
  box-shadow:
    0 6px 32px rgba(0, 0, 0, 0.45),
    inset 0 1px 0 rgba(255, 255, 255, 0.08),
    0 0 20px rgba(0, 242, 254, 0.08);
}

/* 内层渐变光泽（模拟 Demo0 纺理图的叠加效果） */
.glass-card::before {
  content: "";
  position: absolute;
  inset: 0;
  background: linear-gradient(
    135deg,
    rgba(0, 242, 254, 0.04) 0%,
    rgba(79, 172, 254, 0.02) 40%,
    rgba(0, 0, 0, 0) 70%
  );
  pointer-events: none;
  z-index: 0;
}

/* 顶部带青色渐变描边线 */
.glass-card::after {
  content: "";
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 1px;
  background: linear-gradient(
    to right,
    rgba(0, 242, 254, 0.7) 0%,
    rgba(79, 172, 254, 0.4) 50%,
    rgba(0, 242, 254, 0) 100%
  );
  z-index: 1;
}

/* 卡牌内部内容层需在 z-index=1 以上 */
.glass-card > * {
  position: relative;
  z-index: 1;
}

/* 卡牌标题栏 */
.glass-card-header {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 14px 8px;
  border-bottom: 1px solid rgba(141, 141, 141, 0.12);
  flex-shrink: 0;
}

/* 左侧彬色竖线（对标 Demo0 CardTitle 的左边设计） */
.glass-card-accent {
  width: 3px;
  height: 16px;
  border-radius: 2px;
  background: linear-gradient(to bottom, #00f2fe, rgba(0, 242, 254, 0.3));
  box-shadow: 0 0 8px rgba(0, 242, 254, 0.6);
  flex-shrink: 0;
}

.glass-card-title {
  font-size: 14px;
  font-weight: 600;
  color: #e8f4ff;
  letter-spacing: 0.5px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.glass-card-badge-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: #00f2fe;
  box-shadow: 0 0 6px #00f2fe;
  margin-left: auto;
  animation: dotPulse 2s ease-in-out infinite;
}
@keyframes dotPulse {
  0%,
  100% {
    opacity: 1;
    box-shadow: 0 0 6px #00f2fe;
  }
  50% {
    opacity: 0.4;
    box-shadow: 0 0 2px #00f2fe;
  }
}

/* glass-card 内的 metrics-grid 上方需要 padding */
.glass-card .metrics-grid {
  padding: 8px 10px 6px;
  flex: 1;
  min-height: 0;
  height: auto !important; /* 覆盖 data.css 中的 height: 100%，防止 flex 溢出 */
}

/* 图表容器在 glass-card 里需要补充展开 */
.chart-container-glass {
  flex: 1;
  min-height: 0;
  height: auto !important; /* 覆盖 data.css 中的 height: 100%，防止 flex 溢出 */
  padding: 4px 6px 6px;
}

/* ranking-wrap 在 glass-card 里的内嵌调整 */
.ranking-wrap-glass {
  position: absolute !important;
  top: 46px;
  left: 10px;
  right: 10px;
  bottom: 10px;
}

/* AI 卡片内 ai-report-wrap 补充 */
.glass-card .ai-report-wrap {
  padding-top: 6px;
  flex: 1;
  min-height: 0;
  height: auto !important; /* 覆盖 data.css 中的 height: 100%，防止 flex 溢出 */
}

/* ===================================================
   HUD 科技感知核心容器
   =================================================== */
.hud-container {
  position: absolute;
  inset: 12px;
  border: 1px solid rgba(0, 242, 254, 0.08);
  box-shadow: inset 0 0 32px rgba(0, 242, 254, 0.03);
  pointer-events: none;
  border-radius: 8px;
  overflow: hidden;
  z-index: 1;
}

/* 4角高科技定位角标 */
.hud-corner {
  position: absolute;
  width: 14px;
  height: 14px;
  border-color: #00f2fe;
  border-style: solid;
  opacity: 0.75;
  filter: drop-shadow(0 0 4px rgba(0, 242, 254, 0.5));
  pointer-events: none;
  z-index: 2;
}
.hud-corner.top-left {
  top: 12px;
  left: 12px;
  border-width: 2px 0 0 2px;
}
.hud-corner.top-right {
  top: 12px;
  right: 12px;
  border-width: 2px 2px 0 0;
}
.hud-corner.bottom-left {
  bottom: 12px;
  left: 12px;
  border-width: 0 0 2px 2px;
}
.hud-corner.bottom-right {
  bottom: 12px;
  right: 12px;
  border-width: 0 2px 2px 0;
}

/* HUD 顶部悬浮科技标题 */
.hud-title-container {
  position: absolute;
  top: 30px;
  left: 50%;
  transform: translateX(-50%);
  display: flex;
  flex-direction: column;
  align-items: center;
  z-index: 5;
  pointer-events: none;
}
.hud-title-line {
  width: 140px;
  height: 1px;
  background: linear-gradient(to right, transparent, #00f2fe, transparent);
  margin-bottom: 8px;
}
.hud-title-badge {
  background: rgba(10, 20, 42, 0.65);
  border: 1px solid rgba(0, 242, 254, 0.4);
  backdrop-filter: blur(10px);
  border-radius: 20px;
  padding: 6px 20px;
  display: flex;
  align-items: center;
  gap: 10px;
  box-shadow:
    0 4px 15px rgba(0, 0, 0, 0.25),
    inset 0 0 10px rgba(0, 242, 254, 0.15);
}
.hud-title-text {
  font-size: 18px;
  font-weight: bold;
  color: #ffffff;
  letter-spacing: 3px;
  text-shadow: 0 0 10px rgba(0, 242, 254, 0.5);
}
.hud-title-sub {
  font-size: 10px;
  color: rgba(160, 223, 255, 0.55);
  margin-top: 5px;
  letter-spacing: 2px;
  font-family: monospace, sans-serif;
}
.pulse-dot-cyan {
  width: 6px;
  height: 6px;
  background-color: #00f2fe;
  border-radius: 50%;
  box-shadow: 0 0 0 0 rgba(0, 242, 254, 0.7);
  animation: pulsing-cyan 1.8s infinite;
}
@keyframes pulsing-cyan {
  0% {
    transform: scale(0.95);
    box-shadow: 0 0 0 0 rgba(0, 242, 254, 0.7);
  }
  70% {
    transform: scale(1.1);
    box-shadow: 0 0 0 8px rgba(0, 242, 254, 0);
  }
  100% {
    transform: scale(0.95);
    box-shadow: 0 0 0 0 rgba(0, 242, 254, 0);
  }
}

/* ===================================================
   HUD 底部浮动数据卡片（毛玻璃拟态）
   =================================================== */
.hud-container .center-data {
  position: absolute;
  bottom: 180px; /* 整体上移，留出底部的主视野和控制按钮 */
  width: 100%;
  display: flex;
  justify-content: center;
  gap: 20px;
  z-index: 5;
  pointer-events: none;
}
.glass-metric-card {
  background: rgba(10, 20, 42, 0.65);
  border: 1px solid rgba(0, 242, 254, 0.22);
  backdrop-filter: blur(14px);
  -webkit-backdrop-filter: blur(14px);
  border-radius: 12px;
  padding: 16px 28px;
  display: flex;
  flex-direction: column;
  align-items: center;
  box-shadow:
    0 10px 30px rgba(0, 0, 0, 0.4),
    inset 0 1px 0 rgba(255, 255, 255, 0.05),
    inset 0 0 15px rgba(0, 242, 254, 0.08);
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  position: relative;
  pointer-events: auto; /* 关键：恢复鼠标交互 */
  cursor: pointer;
}
.glass-metric-card::after {
  content: "";
  position: absolute;
  bottom: 0;
  left: 10%;
  width: 80%;
  height: 2px;
  background: linear-gradient(to right, transparent, #00f2fe, transparent);
  opacity: 0.7;
}
.glass-metric-card:hover {
  transform: translateY(-5px);
  border-color: rgba(0, 242, 254, 0.5);
  box-shadow:
    0 15px 40px rgba(0, 242, 254, 0.15),
    inset 0 1px 0 rgba(255, 255, 255, 0.1),
    inset 0 0 20px rgba(0, 242, 254, 0.15);
}

.c-label-container {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 6px;
}
.c-label-icon {
  font-size: 16px;
  color: #00f2fe;
  text-shadow: 0 0 8px rgba(0, 242, 254, 0.5);
  display: flex;
  align-items: center;
}
.c-label-text {
  font-size: 14px;
  font-weight: 500;
  color: #a0cfff;
  letter-spacing: 1px;
}

.c-meta-info {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-top: 4px;
  font-size: 12px;
}
.c-trend-up {
  color: #00ffc4;
  font-weight: bold;
  text-shadow: 0 0 6px rgba(0, 255, 196, 0.4);
}
.c-trend-desc {
  color: rgba(160, 207, 255, 0.7);
}
.c-pulse-indicator {
  width: 6px;
  height: 6px;
  background-color: #38ef7d;
  border-radius: 50%;
  display: inline-block;
  box-shadow: 0 0 6px #38ef7d;
}

/* ===================================================
   雷达扫描效果美化优化
   =================================================== */
.hud-container .radar-scan {
  position: absolute;
  top: 50%;
  left: 50%;
  width: 340px;
  height: 340px;
  margin-top: -170px;
  margin-left: -170px;
  border-radius: 50%;
  border: 1.5px dashed rgba(0, 242, 254, 0.15);
  background: conic-gradient(
    from 0deg,
    transparent 60%,
    rgba(0, 242, 254, 0.12) 100%
  );
  animation: scan 6s linear infinite;
  pointer-events: none;
  z-index: 0;
  opacity: 0.65;
}

/* ===================================================
   三维光源调节中枢控制面板样式 (玻璃拟态 + 暗色定制)
   =================================================== */
.light-control-panel {
  position: absolute;
  top: 90px;
  right: 30px;
  width: 320px;
  z-index: 100;
  padding: 16px;
  transition: all 0.6s cubic-bezier(0.4, 0, 0.2, 1);
  transform: translateX(120%);
  opacity: 0;
  pointer-events: none;
  background: rgba(10, 20, 42, 0.65) !important;
}

.light-control-panel.panel-active {
  transform: translateX(0);
  opacity: 1;
  pointer-events: auto;
}

.panel-content {
  margin-top: 12px;
}

/* 深度自定义 Element Plus Tabs 样式 */
:deep(.light-tabs) {
  --el-tabs-header-height: 36px;
}

:deep(.light-tabs .el-tabs__item) {
  color: #a0cfff !important;
  font-size: 13px !important;
  font-weight: 500;
  padding: 0 12px !important;
  transition: all 0.3s;
}

:deep(.light-tabs .el-tabs__item.is-active) {
  color: #00f2fe !important;
  font-weight: bold;
  text-shadow: 0 0 8px rgba(0, 242, 254, 0.5);
}

:deep(.light-tabs .el-tabs__active-bar) {
  background-color: #00f2fe !important;
  box-shadow: 0 0 8px #00f2fe;
}

:deep(.light-tabs .el-tabs__nav-wrap::after) {
  background-color: rgba(141, 141, 141, 0.12) !important;
}

/* 控制项间距 */
.control-item {
  display: flex;
  flex-direction: column;
  gap: 6px;
  margin-bottom: 14px;
}

.control-item:last-child {
  margin-bottom: 4px;
}

.control-label {
  font-size: 12px;
  color: #a0cfff;
  opacity: 0.85;
}

/* 深度自定义 Element Plus Slider 样式 */
:deep(.el-slider) {
  --el-slider-main-bg-color: #00f2fe;
  --el-slider-runway-bg-color: rgba(255, 255, 255, 0.15);
  --el-slider-stop-bg-color: transparent;
  --el-slider-button-size: 12px;
  --el-slider-button-wrapper-size: 30px;
}

:deep(.el-slider__button) {
  border: 2px solid #00f2fe !important;
  background-color: #050a15 !important;
  box-shadow: 0 0 6px #00f2fe;
}

/* 深度自定义 Element Plus ColorPicker 样式 */
:deep(.el-color-picker) {
  width: 100%;
}

:deep(.el-color-picker__trigger) {
  width: 100% !important;
  height: 32px !important;
  border: 1px solid rgba(0, 242, 254, 0.35) !important;
  background: rgba(5, 10, 21, 0.88) !important;
  border-radius: 4px !important;
  padding: 3px !important;
  box-sizing: border-box !important;
}

:deep(.el-color-picker__color) {
  border-radius: 2px !important;
  border: none !important;
}
</style>