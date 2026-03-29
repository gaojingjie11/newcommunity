<template>
  <div class="home-page">
    <Navbar />

    <!-- ★ 重构：智能波浪式 & 滚动避让弹幕区 ★ -->
    <div class="barrage-container" v-show="showBarrage">
      <div
        v-for="(item, index) in currentBarrages"
        :key="item.uniqueId"
        class="barrage-item"
        :style="{
          top: item.top,
          animationDelay: item.delay,
          animationDuration: item.duration
        }"
        @click="goToNotice(item.id)"
      >
        <!-- 新增：用来控制随滚动左右散开的包装层 -->
        <div class="scatter-wrapper" :style="getScatterStyle(index)">
          <div class="barrage-content">
            <el-icon class="barrage-icon"><Bell /></el-icon>
            <span class="barrage-tag">社区资讯</span>
            <span class="barrage-text">{{ item.title }}</span>
          </div>
        </div>
      </div>
    </div>

    <!-- 右下角弹幕控制悬浮按钮 -->
    <div class="barrage-toggle" @click="showBarrage = !showBarrage" :class="{ 'is-closed': !showBarrage }">
      <el-icon class="toggle-icon">
        <ChatLineSquare v-if="showBarrage" />
        <Close v-else />
      </el-icon>
      <span class="toggle-text">{{ showBarrage ? '关闭弹幕公告' : '开启弹幕公告' }}</span>
    </div>

    <div class="banner">
      <div
        v-for="(img, index) in bannerImages"
        :key="index"
        class="banner-slide"
        :class="{ active: currentBanner === index }"
        :style="{ backgroundImage: `url(${img})` }"
      ></div>
      <div class="banner-overlay"></div>

      <div class="container banner-content">
        <h1 class="banner-title">欢迎来到智慧社区</h1>
        <p class="banner-subtitle">让生活更便捷、更智能、更美好</p>
      </div>

      <div class="banner-dots">
        <span
          v-for="(_, index) in bannerImages"
          :key="index"
          class="dot"
          :class="{ active: currentBanner === index }"
          @click="currentBanner = index"
        ></span>
      </div>
    </div>

    <div class="container">
      <div class="service-section">
        <div class="section-header-center">
          <h2 class="section-title-main highlight-title">我们的服务</h2>
          <p class="section-subtitle">全方位社区资讯服务，守护您的美好生活</p>
        </div>

        <div class="quick-menu">
          <template v-for="(item, index) in menuItems" :key="index">
            <div
              class="quick-item"
              @click="$router.push(item.path)"
              v-if="
                !item.requireAdmin ||
                userStore.userInfo?.real_name === '系统管理员'
              "
            >
              <div class="quick-icon-wrap">
                <img
                  :src="item.image"
                  :alt="item.name"
                  class="quick-icon-img"
                />
              </div>
              <div class="quick-text">{{ item.name }}</div>
              <div class="quick-desc">{{ item.desc }}</div>
            </div>
          </template>
        </div>
      </div>
    </div>

    <!-- 社区风采/竖向长图海报展示区 -->
    <div class="container">
      <div class="feature-section">
        <div class="section-header-center">
          <h2 class="section-title-main highlight-title">共建美好社区</h2>
          <p class="section-subtitle">打造智慧、绿色、和谐的现代化居住新标杆</p>
        </div>

        <div class="feature-grid">
          <div class="feature-card" v-for="(item, index) in featureCards" :key="index">
            <div class="feature-img-wrapper">
              <img :src="item.img" :alt="item.title" />
              <!-- 渐变遮罩与文案 -->
              <div class="feature-overlay">
                <div class="overlay-content">
                  <h3>{{ item.title }}</h3>
                  <p>{{ item.desc }}</p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 横屏生活与环保积分展示区（单行交替、图片破框错层排版） -->
    <div class="container">
      <div class="lifestyle-section">
        <div class="section-header-center">
          <h2 class="section-title-main highlight-title">乐享绿居生活</h2>
          <p class="section-subtitle">倡导低碳环保新风尚，体验便捷有温度的社区生活</p>
        </div>

        <!-- 垂直排列的错层卡片 -->
        <div class="lifestyle-list">
          <div class="lifestyle-card" v-for="(item, index) in lifestyleCards" :key="index">
            <div class="lifestyle-img-wrapper">
              <img :src="item.img" :alt="item.title" />
              <div class="image-overlay"></div>
            </div>
            <div class="lifestyle-content">
              <h3>{{ item.title }}</h3>
              <p>{{ item.desc }}</p>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div class="container">
      <div class="news-section">
        <div class="section-header-flex">
          <div class="header-left">
            <span class="hot-tag">NEW</span>
            <h2 class="section-title">最新公告</h2>
          </div>
          <div class="header-right" @click="$router.push('/service/notice')">
            <span class="more-text">查看全部 →</span>
          </div>
        </div>

        <div class="notice-grid">
          <div
            class="notice-card"
            v-for="notice in notices"
            :key="notice.id"
            @click="goToNotice(notice.id)"
          >
            <div class="notice-tag">智慧社区</div>
            <h3 class="notice-title">{{ notice.title }}</h3>
            <p class="notice-excerpt">
              为了更好地服务社区居民，我们持续更新相关动态，请各位业主及时查阅...
            </p>
            <div class="notice-meta">
              <span class="publisher">👤 {{ notice.publisher }}</span>
              <span class="date">🕒 {{ formatDate(notice.created_at) }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from "vue";
import { useRouter } from "vue-router";
import { useUserStore } from "@/stores/user"; 
import Navbar from "@/components/layout/Navbar.vue";
import { getNoticeList } from "@/api/service";
import dayjs from "dayjs";
import { ChatLineSquare, Close, Bell } from '@element-plus/icons-vue';

// 轮播图图片
import banner1 from "@/assets/images/首页1.png";
import banner2 from "@/assets/images/首页2.png";
import banner3 from "@/assets/images/首页3.png";
import banner4 from "@/assets/images/首页4.png";

// 服务区图标
import cartImg from "@/assets/images/购物车.png";
import noticeImg from "@/assets/images/公告.png";
import repairImg from "@/assets/images/维修.png";
import visitorImg from "@/assets/images/需求登记.png";
import dataImg from "@/assets/images/数据大屏 (1).png";

// 竖向长图图片
import featureImg1 from "@/assets/images/长图1.png";
import featureImg2 from "@/assets/images/长图2.png";
import featureImg3 from "@/assets/images/长图3.png";

// 横向生活图片
import lifeImg1 from "@/assets/images/生活1.png";
import lifeImg2 from "@/assets/images/生活2.png";
import garbageImg from "@/assets/images/垃圾分类.png";

const router = useRouter();
const userStore = useUserStore(); 

// ★ 弹幕状态与波浪式管理机制 ★
const showBarrage = ref(true);
const originalNotices = ref([]); // 存储拉取到的真实数据
const currentBarrages = ref([]); // 当前正在飞行的那一波数据
let waveTimer = null;

// ★ 滚动监听相关 ★
const scrollY = ref(0);

const handleScroll = () => {
  scrollY.value = window.scrollY;
};

// 核心逻辑：根据滚动距离动态计算每个弹幕的偏移和透明度
const getScatterStyle = (index) => {
  // 滚动距离在 50px 以内不散开
  if (scrollY.value <= 50) {
    return {
      transform: 'translateX(0)',
      opacity: 1,
      transition: 'transform 0.3s ease-out, opacity 0.3s ease-out'
    };
  }
  
  // 超过 50px 开始计算散开进度 (250px 时完全隐身)
  const progress = Math.min((scrollY.value - 50) / 250, 1); 
  // 偶数向左跑 (-1)，奇数向右跑 (1)
  const direction = index % 2 === 0 ? -1 : 1; 
  // 最大偏移量可以达到 80vw
  const moveX = direction * progress * 80; 
  
  return {
    transform: `translateX(${moveX}vw)`,
    opacity: 1 - progress,
    transition: 'transform 0.1s linear, opacity 0.1s linear' // 滚动时需要跟手，动画给快一点
  };
};

// 核心逻辑：发射一波弹幕
const playWave = () => {
  if (!originalNotices.value.length || !showBarrage.value) {
    // 如果没数据或关闭了弹幕，5秒后再重试探测
    waveTimer = setTimeout(playWave, 5000);
    return;
  }

  let maxTime = 0;
  
  // 把源数据包装为这一波要飞的弹幕
  currentBarrages.value = originalNotices.value.map((n, index) => {
    const delay = Math.random() * 3; // 随机 0 ~ 3 秒内起飞
    const duration = 12 + Math.random() * 15; // 随机 12 ~ 18 秒飞完
    const top = 12 + Math.random() * 30; // 高度分布在 12% ~ 42% 之间
    
    // 记录这波弹幕需要飞多久才能彻底清空屏幕
    if (delay + duration > maxTime) {
      maxTime = delay + duration;
    }

    return {
      ...n,
      uniqueId: `b-${n.id}-${Date.now()}-${index}`, // 保证 key 唯一
      top: `${top}%`,
      delay: `${delay}s`,
      duration: `${duration}s`
    };
  });

  // 等这一波全部飞完，额外静默休息 4 秒钟，再发下一波
  waveTimer = setTimeout(() => {
    currentBarrages.value = []; // 清空 DOM 释放资源
    setTimeout(playWave, 500);  // 等 0.5 秒重新组装下一波
  }, (maxTime + 4) * 1000);
};


// 菜单配置
const menuItems = [
  { name: "社区商城", image: cartImg, desc: "优质好物，送货上门", path: "/mall", requireAdmin: false },
  { name: "公告通知", image: noticeImg, desc: "重要信息，一手掌握", path: "/service/notice", requireAdmin: false },
  { name: "报修投诉", image: repairImg, desc: "报修反馈，贴心管家", path: "/service/repair", requireAdmin: false },
  { name: "访客登记", image: visitorImg, desc: "便捷访客，安全无忧", path: "/service/visitor", requireAdmin: false },
  { name: "数据大屏", image: dataImg, desc: "社区运行，实时监测", path: "/data", requireAdmin: true }
];

const featureCards = [
  { img: featureImg1, title: '智能安防体系', desc: '24小时无死角巡控，人脸识别无感通行，科技守护您的每一天。' },
  { img: featureImg2, title: '绿色生态家园', desc: '打造充满绿意的社区微环境，让您每次呼吸都清新自然。' },
  { img: featureImg3, title: '尊享管家服务', desc: '极速报修响应，在线物业缴费，让您的生活更加从容不迫。' }
];

const lifestyleCards = [
  { img: lifeImg1, title: '多彩社区活动', desc: '丰富多元的文体活动，拉近邻里距离，让社区生活充满欢声笑语与温暖人情。' },
  { img: lifeImg2, title: '全龄友好空间', desc: '精心打造适宜全年龄段的休闲娱乐空间，让老人舒心漫步，让孩子尽情奔跑开心。' },
  { img: garbageImg, title: '环保分类赚积分', desc: 'AI 智能极速识别垃圾分类。您的每一次环保微行动都能转化为绿色积分，直接抵现购物与物业费！' }
];

const bannerImages = [banner1, banner2, banner3, banner4];
const currentBanner = ref(0);
const notices = ref([]);
let bannerTimer = null;

const startBannerLoop = () => {
  bannerTimer = setInterval(() => {
    currentBanner.value = (currentBanner.value + 1) % bannerImages.length;
  }, 4500);
};

const formatDate = (date) => dayjs(date).format("YYYY-MM-DD");
const goToNotice = (id) => router.push(`/service/notice?id=${id}`);

onMounted(async () => {
  window.addEventListener("scroll", handleScroll); // 监听滚动
  startBannerLoop();
  
  try {
    const list = await getNoticeList();
    // 底部卡片区只保留最新3条
    notices.value = list.slice(0, 3);
    
    // 初始化弹幕源数据，并打开发射器
    originalNotices.value = list;
    playWave();
  } catch (error) {
    console.error(error);
  }
});

onUnmounted(() => {
  window.removeEventListener("scroll", handleScroll);
  if (bannerTimer) clearInterval(bannerTimer);
  if (waveTimer) clearTimeout(waveTimer);
});
</script>

<style scoped>
@import "@/assets/styles/home.css";
</style>