<template>
  <div class="data-screen">
    <dv-full-screen-container>
      <!-- ================== 顶部 Header ================== -->
      <div class="screen-header">
        <dv-decoration-8 style="width: 300px; height: 50px" />
        <div class="header-center">
          <dv-decoration-5 style="width: 500px; height: 40px" />
          <div class="title-text">智慧社区数据中枢大屏</div>
        </div>
        <dv-decoration-8 :reverse="true" style="width: 300px; height: 50px" />
        <div class="time-text">{{ currentTime }}</div>
        <div class="back-btn" @click="goBack">
          <el-icon><HomeFilled /></el-icon> 首页
        </div>
      </div>

      <!-- ================== 主体 Body ================== -->
      <div class="screen-body">
        <!-- ================== 左侧栏 ================== -->
        <div class="column-side">
          <dv-border-box-11 class="tech-card box-h-30" title="实时核心指标">
            <div class="metrics-grid">
              <div class="metric-item">
                <div class="m-icon" style="color: #00f2fe">
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
          </dv-border-box-11>

          <dv-border-box-13 class="tech-card box-h-35 mt-15">
            <div class="chart-title">
              <span class="indicator"></span>工单问题分类占比
            </div>
            <div ref="pieChartRef" class="chart-container"></div>
          </dv-border-box-13>

          <dv-border-box-13 class="tech-card box-h-35 mt-15">
            <div class="chart-title">
              <span class="indicator"></span>7日营收趋势分析
            </div>
            <div ref="lineChartRef" class="chart-container"></div>
          </dv-border-box-13>
        </div>

        <!-- ================== 中间栏 ================== -->
        <div class="column-center">
          <div class="center-map-box" ref="centerMapBoxRef">
            <div ref="threeContainerRef" class="three-canvas-container"></div>
            <div
              v-if="viewMode === 'floors'"
              class="floor-back-btn"
              @click="switchToBuilding"
            >
              ← 返回建筑总览
            </div>

            <!-- 悬浮窗 -->
            <div
              v-show="hoverTooltip.visible"
              class="floor-tooltip"
              :style="{
                left: hoverTooltip.x + 'px',
                top: hoverTooltip.y + 'px',
              }"
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
            <dv-border-box-10 style="pointer-events: none">
              <div class="map-mask"></div>
              <div class="map-title-overlay">
                {{
                  viewMode === "floors"
                    ? "楼层 3D 分布"
                    : "社区 3D 态势感知模型"
                }}
              </div>

              <div class="center-data">
                <div class="c-item">
                  <div class="c-label">当月社区总营收 (元)</div>
                  <div class="c-num">
                    <dv-digital-flop
                      :config="flopIncomeConfig"
                      style="width: 240px; height: 50px"
                    />
                  </div>
                </div>
                <div class="c-item">
                  <div class="c-label">活跃用户基数 (人)</div>
                  <div class="c-num">
                    <dv-digital-flop
                      :config="flopUserConfig"
                      style="width: 200px; height: 50px"
                    />
                  </div>
                </div>
              </div>

              <div class="radar-scan"></div>
            </dv-border-box-10>
          </div>
        </div>

        <!-- ================== 右侧栏 ================== -->
        <div class="column-side">
          <dv-border-box-13 class="tech-card box-h-35">
            <div class="chart-title flex-between">
              <div class="title-left">
                <span class="indicator"></span>社区环保积分先锋榜
              </div>
              <el-radio-group
                v-model="rankingView"
                size="small"
                class="dark-radio"
              >
                <el-radio-button label="datav">动态展示</el-radio-button>
                <el-radio-button label="table">经典表格</el-radio-button>
              </el-radio-group>
            </div>
            <div class="ranking-wrap">
              <dv-scroll-ranking-board
                v-if="rankingView === 'datav' && rankingBoardConfig.data.length"
                :config="rankingBoardConfig"
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
                  width="85"
                  align="right"
                >
                  <template #default="scope">
                    <strong style="color: #00f2fe">{{
                      scope.row.points
                    }}</strong>
                  </template>
                </el-table-column>
              </el-table>

              <div v-else class="empty-data">暂无排行数据</div>
            </div>
          </dv-border-box-13>

          <dv-border-box-13 class="tech-card box-h-35 mt-15">
            <div class="chart-title">
              <span class="indicator"></span>社区各模块收入构成
            </div>
            <div ref="barChartRef" class="chart-container"></div>
          </dv-border-box-13>

          <dv-border-box-11
            class="tech-card box-h-30 mt-15"
            title="AI 智能预警与诊断"
          >
            <div class="ai-report-wrap">
              <div class="ai-status">
                <div class="status-left">
                  <div class="status-dot pulse"></div>
                  <span>AI 守护引擎运行中</span>
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
                  {{
                    aiReport?.report || "正在实时诊断社区运营数据，请稍候..."
                  }}
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
          </dv-border-box-11>
        </div>
      </div>

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
    </dv-full-screen-container>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, reactive, nextTick, computed } from "vue";
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
} from "@element-plus/icons-vue";
import { useUserStore } from "@/stores/user";
import * as THREE from "three";
import { GLTFLoader } from "three/addons/loaders/GLTFLoader.js";
import { DRACOLoader } from "three/addons/loaders/DRACOLoader.js";
import { OrbitControls } from "three/addons/controls/OrbitControls.js";

const router = useRouter();
const userStore = useUserStore();
const isAdmin = ref(userStore.userInfo?.role === "admin");

const viewMode = ref("building");

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
  if (!aiReport.value?.report)
    return '<div style="color:#a0cfff;">暂无详细数据</div>';
  let html = aiReport.value.report;
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
let threeObserver = null;
let threeInited = false;

let threeModelMeshes = []; // 建筑总览所有mesh
let floorModelMeshes = []; // 楼层所有mesh
let hoveredMesh = null; // 当前总览模式下高亮的单个mesh
let hoveredFloorMesh = null; // 当前楼层模式下高亮的单个mesh

const threeRaycaster = new THREE.Raycaster();
const threeMouse = new THREE.Vector2();

const initThreeJS = () => {
  const container = threeContainerRef.value;
  if (!container) return;

  const rect = container.getBoundingClientRect();
  const width = rect.width || container.offsetWidth || 600;
  const height = rect.height || container.offsetHeight || 400;

  threeScene = new THREE.Scene();
  threeScene.background = null;

  threeCamera = new THREE.PerspectiveCamera(45, width / height, 0.1, 1000);
  // 👉 初始默认：总览模式（近）
  threeCamera.position.set(5, 2.5, 4);

  threeRenderer = new THREE.WebGLRenderer({ antialias: true, alpha: true });
  threeRenderer.setSize(width, height);
  threeRenderer.setPixelRatio(window.devicePixelRatio);
  threeRenderer.outputColorSpace = THREE.SRGBColorSpace;
  threeRenderer.shadowMap.enabled = true;
  container.appendChild(threeRenderer.domElement);

  threeControls = new OrbitControls(threeCamera, threeRenderer.domElement);
  threeControls.enableDamping = true;
  threeControls.dampingFactor = 0.05;
  threeControls.autoRotate = false;
  threeControls.enablePan = false;
  // 👉 初始总览模式距离范围
  threeControls.minDistance = 2;
  threeControls.maxDistance = 8;

  // 右视角锁定
  threeControls.minAzimuthAngle = Math.PI / 3;
  threeControls.maxAzimuthAngle = (Math.PI * 2) / 3;
  threeControls.minPolarAngle = Math.PI / 6;
  threeControls.maxPolarAngle = Math.PI / 2.5;

  // 灯光
  const ambientLight = new THREE.AmbientLight(0x0a1628, 2.5);
  threeScene.add(ambientLight);
  const dirLight = new THREE.DirectionalLight(0x00f2fe, 3);
  dirLight.position.set(5, 10, 7);
  dirLight.castShadow = true;
  threeScene.add(dirLight);
  const fillLight = new THREE.DirectionalLight(0x4facfe, 1.5);
  fillLight.position.set(-5, 2, -5);
  threeScene.add(fillLight);
  const groundLight = new THREE.PointLight(0x00f2fe, 1, 20);
  groundLight.position.set(0, -2, 0);
  threeScene.add(groundLight);

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
    "https://communitysvc.xyz/community/building/build_compressed.glb",
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
            child.material = child.material.clone();
            child.material.envMapIntensity = 1.5;
            // 保存每个mesh原始自发光色，用于离开时恢复
            child.userData.originalEmissive =
              child.material.emissive?.getHex() || 0x000000;
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
              const mat = child.material.clone();
              mat.emissive = new THREE.Color(0x003366);
              mat.emissiveIntensity = 0.3 + i * 0.05;
              child.material = mat;
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
    const rect = container.getBoundingClientRect();
    threeMouse.x = ((e.clientX - rect.left) / rect.width) * 2 - 1;
    threeMouse.y = -((e.clientY - rect.top) / rect.height) * 2 + 1;
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
        if (hoveredMesh.userData.originalEmissive !== undefined) {
          hoveredMesh.material.emissive.setHex(
            hoveredMesh.userData.originalEmissive,
          );
        }
        hoveredMesh = null;
        document.body.style.cursor = "default";
        hoverTooltip.visible = false;
      }

      if (targetObj && hoveredMesh !== targetObj) {
        hoveredMesh = targetObj;
        hoveredMesh.userData.originalEmissive =
          hoveredMesh.material.emissive.getHex();
        hoveredMesh.material.emissive.setHex(0x00f2fe);
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
            child.material.emissive.setHex(0x00f2fe);
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
  const animate = () => {
    threeAnimFrameId = requestAnimationFrame(animate);
    threeControls.update();
    threeRenderer.render(threeScene, threeCamera);
  };
  animate();
};

const handleThreeResize = () => {
  if (!threeRenderer || !threeCamera) return;
  const w = threeContainerRef.value.clientWidth;
  const h = threeContainerRef.value.clientHeight;
  threeCamera.aspect = w / h;
  threeCamera.updateProjectionMatrix();
  threeRenderer.setSize(w, h);
};

// DataV 翻牌器
const flopIncomeConfig = reactive({
  number: [0],
  content: "¥ {nt}",
  style: { fontSize: 36, fill: "#00f2fe", fontWeight: "bold" },
});
const flopUserConfig = reactive({
  number: [0],
  style: { fontSize: 36, fill: "#00f2fe", fontWeight: "bold" },
});
const rankingBoardConfig = reactive({
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
    flopIncomeConfig.number = [parseFloat(stats.value.monthIncome || 0)];
    flopUserConfig.number = [stats.value.totalUsers || 0];
    if (leaderboardRes?.list) {
      leaderboardList.value = leaderboardRes.list;
      rankingBoardConfig.data = leaderboardRes.list.map((i) => ({
        name: i.nickname || i.username || `用户${i.user_id}`,
        value: i.points || 0,
      }));
    }
    if (isAdmin.value) aiReport.value = await getAIReport();
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
    color: ["#85a5ff", "#0A82A4", "#3498db4c", "#2e86c1", "#3498db", "#5dade2"],
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
        itemStyle: { borderColor: "#050a15", borderWidth: 2 },
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
      data: ["物业费", "停车费", "商城消费", "场馆预约"],
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
        data: stats.value.costStructure || [3200, 1800, 4500, 1200],
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
});
</script>

<style scoped>
@import "../../assets/styles/data.css";
/* 悬浮窗样式 */
.floor-tooltip {
  position: fixed;
  z-index: 9999;
  background: rgba(5, 15, 30, 0.85);
  border: 1px solid #00f2fe;
  border-radius: 6px;
  padding: 12px 16px;
  color: #fff;
  pointer-events: none;
  box-shadow: 0 0 15px rgba(0, 242, 254, 0.3);
  backdrop-filter: blur(4px);
  display: flex;
  flex-direction: column;
  gap: 8px;
  min-width: 140px;
}
.floor-tooltip .tt-title {
  font-size: 15px;
  font-weight: bold;
  color: #00f2fe;
  padding-bottom: 6px;
  border-bottom: 1px dashed rgba(0, 242, 254, 0.3);
  margin-bottom: 4px;
  text-shadow: 0 0 5px rgba(0, 242, 254, 0.5);
}
.floor-tooltip .tt-item {
  font-size: 13px;
  color: #a0cfff;
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.floor-tooltip .tt-item i {
  font-style: normal;
  margin-right: 4px;
}
.floor-tooltip .tt-item .hl {
  font-weight: bold;
  color: #fff;
  font-size: 15px;
  margin-left: 10px;
}

.data-screen {
  width: 100vw;
  height: 100vh;
  background: #050a15;
  color: #fff;
  overflow: hidden;
  font-family:
    -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue",
    Arial, sans-serif;
}
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
  position: absolute;
  left: 50%;
  transform: translateX(-50%);
  top: 5px;
  display: flex;
  flex-direction: column;
  align-items: center;
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

.tech-card {
  position: relative;
  box-sizing: border-box;
  padding: 20px 15px 15px;
}
.mt-15 {
  margin-top: 15px;
}
.box-h-30 {
  height: calc(30% - 10px);
}
.box-h-35 {
  height: calc(35% - 10px);
}

.chart-title {
  position: absolute;
  top: 10px;
  left: 20px;
  font-size: 16px;
  font-weight: bold;
  color: #fff;
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
  background: #00f2fe;
  margin-right: 8px;
  box-shadow: 0 0 8px #00f2fe;
}
.chart-container {
  width: 100%;
  height: 100%;
  padding-top: 25px;
  box-sizing: border-box;
}

.metrics-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 8px;
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
  color: #fff;
  font-family: Arial, sans-serif;
}
.num-small {
  font-size: 16px;
}

.center-map-box {
  width: 100%;
  height: 100%;
  position: relative;
  padding: 5px;
}
.three-canvas-container {
  position: absolute;
  top: 10px;
  left: 10px;
  right: 10px;
  bottom: 10px;
  border-radius: 8px;
  overflow: hidden;
  z-index: 1;
  pointer-events: auto;
}
.three-canvas-container canvas {
  display: block;
  width: 100% !important;
  height: 100% !important;
}
.map-mask {
  position: absolute;
  top: 10px;
  left: 10px;
  right: 10px;
  bottom: 10px;
  background: radial-gradient(
    circle,
    rgba(5, 10, 21, 0) 0%,
    rgba(5, 10, 21, 0.8) 100%
  );
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
  pointer-events: none;
}
.center-data {
  position: absolute;
  bottom: 80px;
  width: 100%;
  display: flex;
  justify-content: space-evenly;
  z-index: 5;
  pointer-events: none;
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

.radar-scan {
  position: absolute;
  top: 50%;
  left: 50%;
  width: 400px;
  height: 400px;
  margin: -200px 0 0 -200px;
  border-radius: 50%;
  border: 1px dashed rgba(0, 242, 254, 0.2);
  background: conic-gradient(
    from 0deg,
    transparent 70%,
    rgba(0, 242, 254, 0.3) 100%
  );
  animation: scan 4s linear infinite;
  pointer-events: none;
  z-index: -1;
}
@keyframes scan {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}

.ranking-wrap {
  position: absolute;
  top: 55px;
  left: 15px;
  right: 15px;
  bottom: 15px;
  overflow: hidden;
}
.floor-back-btn {
  position: absolute;
  top: 40px;
  left: 30px;
  z-index: 20;
  color: #00f2fe;
  cursor: pointer;
  font-size: 15px;
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 14px;
  border: 1px solid rgba(0, 242, 254, 0.4);
  border-radius: 4px;
  transition: all 0.3s;
  background: rgba(0, 242, 254, 0.1);
  backdrop-filter: blur(4px);
  pointer-events: auto;
}
.floor-back-btn:hover {
  background: rgba(0, 242, 254, 0.3);
  box-shadow: 0 0 10px rgba(0, 242, 254, 0.6);
  color: #fff;
}
.empty-data {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100%;
  color: #606266;
  font-size: 14px;
}

.ai-report-wrap {
  height: 100%;
  display: flex;
  flex-direction: column;
  padding-top: 25px;
  overflow: hidden;
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
  background: #4facfe;
  border-radius: 50%;
}
.pulse {
  box-shadow: 0 0 0 0 rgba(58, 67, 194, 0.7);
  animation: pulsing 1.5s infinite;
}
@keyframes pulsing {
  0% {
    transform: scale(0.95);
    box-shadow: 0 0 0 0 rgba(103, 194, 58, 0.7);
  }
  70% {
    transform: scale(1);
    box-shadow: 0 0 0 6px rgba(103, 194, 58, 0);
  }
  100% {
    transform: scale(0.95);
    box-shadow: 0 0 0 0 rgba(103, 194, 58, 0);
  }
}
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
.tag.danger {
  color: #b2d6eb;
  border: 1px solid rgba(131, 108, 245, 0.3);
}
.tag.success {
  color: #b2d6eb;
  border: 1px solid rgba(131, 108, 245, 0.3);
}

.custom-scrollbar::-webkit-scrollbar {
  width: 4px;
}
.custom-scrollbar::-webkit-scrollbar-track {
  background: rgba(0, 0, 0, 0.1);
}
.custom-scrollbar::-webkit-scrollbar-thumb {
  background: rgba(0, 242, 254, 0.3);
  border-radius: 2px;
}
.custom-scrollbar::-webkit-scrollbar-thumb:hover {
  background: rgba(0, 242, 254, 0.6);
}

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
:deep(
  .dark-radio .el-radio-button__original-radio:checked + .el-radio-button__inner
) {
  background: rgba(0, 242, 254, 0.2);
  color: #00f2fe;
  border-color: #00f2fe;
  text-shadow: 0 0 8px rgba(0, 242, 254, 0.6);
}
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
.rank-1 {
  background: #e6a23c;
  color: #fff;
  box-shadow: 0 0 8px #e6a23c;
}
.rank-2 {
  background: #a0cfff;
  color: #fff;
  box-shadow: 0 0 8px #a0cfff;
}
.rank-3 {
  background: #e0a370;
  color: #fff;
  box-shadow: 0 0 8px #e0a370;
}
</style>

<style>
:global(.dark-theme-dialog) {
  background-color: #0b172a !important;
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
:global(.dark-theme-dialog .el-dialog__headerbtn .el-dialog__close) {
  color: #00f2fe !important;
  font-size: 18px;
}
:global(.dark-theme-dialog .el-dialog__headerbtn:hover .el-dialog__close) {
  color: #fff !important;
  text-shadow: 0 0 8px #00f2fe;
}
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
  color: #e6a23c;
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
</style>
