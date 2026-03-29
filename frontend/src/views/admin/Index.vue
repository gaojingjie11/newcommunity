<template>
  <div class="admin-page">
    <Navbar />

    <div class="container custom-container">
      

      <div class="admin-menu">
        <div
          v-for="item in visibleCards"
          :key="item.path"
          class="admin-card"
          @click="$router.push(item.path)"
        >
          <div class="admin-icon-wrap">
            <el-icon class="admin-icon">
              <component :is="item.icon" />
            </el-icon>
          </div>
          <div class="admin-content">
            <h3>{{ item.name }}</h3>
            <p>{{ item.desc }}</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import Navbar from '@/components/layout/Navbar.vue'
import { useUserStore } from '@/stores/user'
import {
  User,
  Goods,
  List,
  Shop,
  Bell,
  Tools,
  UserFilled,
  Van,
  Wallet,
  DataAnalysis,
  Monitor
} from '@element-plus/icons-vue'

const userStore = useUserStore()
const role = computed(() => userStore.userInfo.role || 'user')

function hasPermission(allowedRoles) {
  return allowedRoles.includes(role.value)
}

const cards = [
  {
    path: '/admin/users',
    name: '用户管理',
    desc: '冻结、解冻、分配角色、调整余额',
    icon: User,
    roles: ['admin']
  },
  {
    path: '/admin/products',
    name: '商品管理',
    desc: '维护商城商品与营销活动',
    icon: Goods,
    roles: ['admin', 'store']
  },
  {
    path: '/admin/orders',
    name: '订单管理',
    desc: '查看订单状态与发货进度',
    icon: List,
    roles: ['admin', 'store']
  },
  {
    path: '/admin/stores',
    name: '门店管理',
    desc: '维护社区门店与商品绑定',
    icon: Shop,
    roles: ['admin', 'store']
  },
  {
    path: '/admin/notices',
    name: '公告管理',
    desc: '发布和删除社区公告',
    icon: Bell,
    roles: ['admin', 'property']
  },
  {
    path: '/admin/repairs',
    name: '报修管理',
    desc: '处理工单并跟踪完结情况',
    icon: Tools,
    roles: ['admin', 'property']
  },
  {
    path: '/admin/visitors',
    name: '访客管理',
    desc: '审核访客登记与通行状态',
    icon: UserFilled,
    roles: ['admin', 'property']
  },
  {
    path: '/admin/parking',
    name: '车位管理',
    desc: '新增车位并维护绑定关系',
    icon: Van,
    roles: ['admin', 'property']
  },
  {
    path: '/admin/property-fee',
    name: '物业费管理',
    desc: '创建账单并查看混合支付明细',
    icon: Wallet,
    roles: ['admin', 'property']
  },
  {
    path: '/admin/ai-report',
    name: 'AI 报表中心',
    desc: '手动生成社区分析报表并支持列表刷新',
    icon: DataAnalysis,
    roles: ['admin']
  },
  {
    path: '/data',
    name: '智能数据大屏',
    desc: '查看营收趋势、排行榜与运营大盘',
    icon: Monitor,
    roles: ['admin', 'property']
  }
]

const visibleCards = computed(() => cards.filter((item) => hasPermission(item.roles)))
</script>

<style scoped>
/* 全局背景 */
.admin-page {
  min-height: 100vh;
  background-color: #f8f9fa;
  padding-bottom: 80px;
}

/* 加宽容器，更符合管理后台的气质 */
.custom-container {
  max-width: 1200px;
  margin: 0 auto;
}

/* 统一的高光标题 */
.page-header {
  padding: 32px 0 24px;
}

.highlight-title {
  display: inline-block;
  position: relative;
  font-size: 32px;
  color: #2c3e50;
  font-weight: 700;
  margin: 0;
  z-index: 1;
}

.highlight-title::after {
  content: '';
  position: absolute;
  bottom: 4px;
  left: -5%;
  width: 110%;
  height: 14px;
  background-color: #2d597b; 
  opacity: 0.15;
  border-radius: 6px;
  z-index: -1;
  transition: all 0.3s ease;
}

/* ★ 网格布局：强制一行四个 ★ */
.admin-menu {
  margin-top: 50px;
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 24px;
  margin-top: 16px;
}

/* 高级模块卡片 */
.admin-card {
  background: #ffffff;
  border-radius: 16px;
  padding: 36px 24px;
  text-align: center;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.03);
  transition: all 0.4s cubic-bezier(0.175, 0.885, 0.32, 1.275);
  border: 1px solid transparent;
  display: flex;
  flex-direction: column;
  align-items: center;
  cursor: pointer;
}

.admin-card:hover {
  transform: translateY(-6px);
  box-shadow: 0 12px 32px rgba(45, 89, 123, 0.08);
  border-color: rgba(45, 89, 123, 0.1);
}

/* 卡片内部图标 */
.admin-icon-wrap {
  width: 72px;
  height: 72px;
  border-radius: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f0f7ff;
  color: #2d597b;
  font-size: 32px;
  margin-bottom: 24px;
  transition: all 0.4s ease;
}

/* 悬停联动：图标颜色翻转反差 */
.admin-card:hover .admin-icon-wrap {
  background: #2d597b;
  color: #ffffff;
  transform: scale(1.1) rotate(-5deg); /* 微小的旋转和放大，视觉更灵动 */
  box-shadow: 0 8px 16px rgba(45, 89, 123, 0.2);
}

/* 内容区域 */
.admin-content {
  flex: 1; 
  display: flex;
  flex-direction: column;
}

.admin-card h3 {
  font-size: 18px;
  font-weight: 700;
  color: #2c3e50;
  margin: 0 0 12px 0;
  transition: color 0.3s;
}

.admin-card:hover h3 {
  color: #2d597b;
}

.admin-card p {
  font-size: 14px;
  color: #8c939d;
  line-height: 1.6;
  margin: 0;
}

/* 响应式：在屏幕变窄时逐级减少列数 */
@media (max-width: 1024px) {
  .admin-menu {
    grid-template-columns: repeat(3, 1fr);
  }
}

@media (max-width: 768px) {
  .admin-menu {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 480px) {
  .admin-menu {
    grid-template-columns: 1fr;
  }
}
</style>