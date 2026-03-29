<template>
  <div class="service-page">
    <Navbar />

    <div class="container">
      <!-- 统一的高光标题样式 -->
      <div class="page-header">
        <h1 class="page-title highlight-title">社区服务</h1>
      </div>

      <div class="service-menu">
        <div
          v-for="item in serviceCards"
          :key="item.path"
          class="service-card card"
          @click="$router.push(item.path)"
        >
          <div class="service-icon-wrap">
            <el-icon class="service-icon">
              <component :is="item.icon" />
            </el-icon>
          </div>
          <div class="service-name">{{ item.name }}</div>
          <div class="service-desc">{{ item.desc }}</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import Navbar from '@/components/layout/Navbar.vue'
import { Bell, Tools, UserFilled, Van, Wallet, Trophy, ChatDotRound } from '@element-plus/icons-vue'

const serviceCards = [
  {
    path: '/service/notice',
    name: '公告通知',
    desc: '查看社区最新公告',
    icon: Bell
  },
  {
    path: '/service/repair',
    name: '报修投诉',
    desc: '提交报修并跟踪处理进度',
    icon: Tools
  },
  {
    path: '/service/visitor',
    name: '访客登记',
    desc: '录入访客信息并等待审核',
    icon: UserFilled
  },
  {
    path: '/service/parking',
    name: '车位管理',
    desc: '查看车位绑定信息',
    icon: Van
  },
  {
    path: '/service/property',
    name: '物业费',
    desc: '支持积分优先抵扣的混合支付',
    icon: Wallet
  },
  {
    path: '/service/green-points',
    name: '绿色积分',
    desc: '垃圾分类识别、积分发放与排行榜',
    icon: Trophy
  },
  {
    path: '/service/community-chat',
    name: '社区群聊',
    desc: '全社区用户公共讨论区',
    icon: ChatDotRound
  }
]
</script>

<style scoped>
.service-page {
  min-height: 100vh;
  padding-bottom: var(--spacing-xl);
}

/* ================= 统一标题样式 ================= */
.page-header {
  padding: 32px 0 24px;
}

.highlight-title {
  display: inline-block;
  position: relative;
  font-size: 28px;
  color: #2c3e50;
  font-weight: 600;
  margin: 0;
  z-index: 1;
}

.highlight-title::after {
  content: '';
  position: absolute;
  bottom: 4px;
  left: -5%;
  width: 110%;
  height: 12px;
  background-color: #2d597b; /* 统一的商务蓝 */
  opacity: 0.2;
  border-radius: 4px;
  z-index: -1;
  transition: all 0.3s ease;
}

/* ================= 原有卡片样式保留 ================= */
.service-menu {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
  gap: var(--spacing-lg);
  margin-top: var(--spacing-xl);
}

.service-card {
  text-align: center;
  padding: var(--spacing-xl);
  cursor: pointer;
  transition: all var(--transition-base);
}

.service-card:hover {
  transform: translateY(-4px) scale(1.01);
}

.service-icon-wrap {
  width: 64px;
  height: 64px;
  border-radius: 20px;
  margin: 0 auto var(--spacing-md);
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, rgba(0, 180, 137, 0.12), rgba(0, 180, 137, 0.2));
}

.service-icon {
  font-size: 30px;
  color: var(--primary-color);
}

.service-name {
  font-size: var(--font-size-xl);
  font-weight: 600;
  margin-bottom: var(--spacing-sm);
}

.service-desc {
  color: var(--text-secondary);
}
</style>