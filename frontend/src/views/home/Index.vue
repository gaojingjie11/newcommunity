<template>
  <div class="home-page">
    <Navbar />

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
import { useUserStore } from "@/stores/user"; // 引入 userStore 获取权限信息
import Navbar from "@/components/layout/Navbar.vue";
import { getNoticeList } from "@/api/service";
import dayjs from "dayjs";

// 轮播图图片
import banner1 from "@/assets/images/首页1.png";
import banner2 from "@/assets/images/首页2.png";
import banner3 from "@/assets/images/首页3.png";
import banner4 from "@/assets/images/首页4.png";

// 服务区图标替换为本地图片
import cartImg from "@/assets/images/购物车.png";
import noticeImg from "@/assets/images/公告.png";
import repairImg from "@/assets/images/维修.png";
import visitorImg from "@/assets/images/需求登记.png";
import dataImg from "@/assets/images/数据大屏 (1).png";

const router = useRouter();
const userStore = useUserStore(); // 初始化 userStore

// 菜单配置：加入 requireAdmin 标识需要管理员权限
const menuItems = [
  {
    name: "社区商城",
    image: cartImg,
    desc: "优质好物，送货上门",
    path: "/mall",
    requireAdmin: false,
  },
  {
    name: "公告通知",
    image: noticeImg,
    desc: "重要信息，一手掌握",
    path: "/service/notice",
    requireAdmin: false,
  },
  {
    name: "报修投诉",
    image: repairImg,
    desc: "报修反馈，贴心管家",
    path: "/service/repair",
    requireAdmin: false,
  },
  {
    name: "访客登记",
    image: visitorImg,
    desc: "便捷访客，安全无忧",
    path: "/service/visitor",
    requireAdmin: false,
  },
  {
    name: "数据大屏",
    image: dataImg,
    desc: "社区运行，实时监测",
    path: "/data",
    requireAdmin: true,
  }, // 只有系统管理员可见
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
  startBannerLoop();
  try {
    const list = await getNoticeList();
    notices.value = list.slice(0, 3);
  } catch (error) {
    console.error(error);
  }
});

onUnmounted(() => {
  if (bannerTimer) clearInterval(bannerTimer);
});
</script>

<style scoped>
.home-page {
  background-color: var(--bg-color); /* 适配最新的全局变量 */
}

/* Banner 区域 */
.banner {
  height: 420px;
  position: relative;
  overflow: hidden;
  margin-bottom: 50px;
}

.banner-slide {
  position: absolute;
  width: 100%;
  height: 100%;
  background-size: cover;
  background-position: center;
  opacity: 0;
  transition: opacity 1.2s ease;
}

.banner-slide.active {
  opacity: 1;
}

.banner-overlay {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;

  z-index: 2;
}

.banner-content {
  position: relative;
  z-index: 3;
  color: white;
  text-align: center;
  padding-top: 130px;
}

.banner-title {
  font-size: 52px;
  font-weight: 800;
  text-shadow: 0 4px 10px rgba(0, 0, 0, 0.2);
  margin-bottom: 15px;
}
.banner-subtitle {
  font-size: 20px;
}

.banner-dots {
  position: absolute;
  bottom: 20px;
  left: 50%;
  transform: translateX(-50%);
  z-index: 3;
  display: flex;
  gap: 8px;
}
.dot {
  width: 24px;
  height: 4px;
  background: rgba(255, 255, 255, 0.4);
  cursor: pointer;
  transition: all 0.3s;
}
.dot.active {
  background: var(--primary-color);
  width: 30px;
}

/* ================== 服务区 ================== */
.service-section {
  /* ★ 增加了距离下方公告区的边距 ★ */
  margin-bottom: 80px;
}

.section-header-center {
  text-align: center;
  margin-bottom: 40px;
}

/* 高光笔样式标题 */
.highlight-title {
  display: inline-block;
  position: relative;
  font-size: 28px;
  color: var(--text-primary);
  font-weight: 600;
  margin-bottom: 15px;
  z-index: 1;
}

.highlight-title::after {
  content: "";
  position: absolute;
  bottom: 4px;
  left: -5%;
  width: 110%;
  height: 12px;
  background-color: #0d347c;
  opacity: 0.25;
  border-radius: 4px;
  z-index: -1;
  transition: all 0.3s ease;
}

.highlight-title:hover::after {
  opacity: 0.4;
  bottom: 6px;
}

.section-subtitle {
  font-size: 14px;
  color: var(--text-secondary);
}

.quick-menu {
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
  gap: 20px;
}

.quick-item {
  flex: 1;
  min-width: 200px;
  max-width: 220px;
  background: var(--bg-white);
  border-radius: var(--border-radius);
  padding: 30px 20px;
  text-align: center;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: var(--shadow-sm);
  border: 1px solid transparent;
  border-bottom: 3px solid transparent;
}

.quick-item:hover {
  transform: translateY(-5px);
  box-shadow: var(--shadow-md);
  border-bottom-color: var(--primary-color);
}

.quick-icon-wrap {
  width: 70px;
  height: 70px;
  background: var(--bg-gray);
  border-radius: 50%;
  margin: 0 auto 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: 0.3s;
}

.quick-icon-img {
  width: 36px;
  height: 36px;
  object-fit: contain;
  transition: transform 0.3s ease;
}

.quick-item:hover .quick-icon-wrap {
  background: rgba(9, 132, 227, 0.1);
}

.quick-item:hover .quick-icon-img {
  transform: scale(1.15);
}

.quick-text {
  font-size: 18px;
  font-weight: bold;
  color: var(--text-primary);
  margin-bottom: 8px;
}

.quick-desc {
  font-size: 13px;
  color: var(--text-secondary);
}

/* ================== 公告区 ================== */
.news-section {
  margin-bottom: 60px;
}

.section-header-flex {
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-left: 5px solid var(--primary-color);
  padding-left: 15px;
  margin-bottom: 30px;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 8px;
}

.hot-tag {
  background: var(--primary-color);
  color: white;
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: bold;
}

.section-title {
  font-size: 22px;
  font-weight: bold;
  color: var(--text-primary);
}

.more-text {
  font-size: 14px;
  color: var(--text-secondary);
  border: 1px solid var(--border-color);
  padding: 4px 12px;
  border-radius: 16px;
  transition: all 0.3s;
  cursor: pointer;
}

.more-text:hover {
  background: var(--primary-color);
  color: white;
  border-color: var(--primary-color);
}

.notice-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 20px;
}

.notice-card {
  background: var(--bg-white);
  border-radius: var(--border-radius-sm);
  padding: 24px;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: var(--shadow-sm);
  display: flex;
  flex-direction: column;
}

.notice-card:hover {
  box-shadow: var(--shadow-md);
  transform: translateY(-4px);
}

.notice-tag {
  background: rgba(9, 132, 227, 0.1);
  color: var(--primary-color);
  font-weight: 600;
  padding: 4px 12px;
  border-radius: 20px;
  font-size: 12px;
  margin-bottom: 15px;
  width: fit-content;
}

.notice-title {
  font-size: 18px;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 12px;
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  transition: color 0.3s;
}

.notice-card:hover .notice-title {
  color: var(--primary-color);
}

.notice-excerpt {
  font-size: 14px;
  color: var(--text-secondary);
  line-height: 1.6;
  margin-bottom: 24px;
  flex-grow: 1;
}

.notice-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 13px;
  color: var(--text-light);
  border-top: 1px solid var(--border-color);
  padding-top: 16px;
}

@media (max-width: 992px) {
  .notice-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 768px) {
  .banner-title {
    font-size: 32px;
  }
  .notice-grid {
    grid-template-columns: 1fr;
  }
}
</style>
