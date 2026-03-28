<template>
  <div class="admin-page">
    <Navbar />

    <div class="container">
      <h1 class="page-title">管理后台</h1>

      <div class="admin-menu">
        <div
          v-for="item in visibleCards"
          :key="item.path"
          class="admin-card card"
        >
          <div class="admin-icon-wrap">
            <el-icon class="admin-icon">
              <component :is="item.icon" />
            </el-icon>
          </div>
          <h3>{{ item.name }}</h3>
          <p>{{ item.desc }}</p>
          <el-button :type="item.buttonType || 'primary'" @click="$router.push(item.path)">进入</el-button>
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
    roles: ['admin'],
    buttonType: 'success'
  },
  {
    path: '/data',
    name: '智能数据大屏',
    desc: '查看营收趋势、排行榜与运营大盘',
    icon: Monitor,
    roles: ['admin', 'property'],
    buttonType: 'success'
  }
]

const visibleCards = computed(() => cards.filter((item) => hasPermission(item.roles)))
</script>

<style scoped>
.admin-page {
  min-height: 100vh;
  padding-bottom: var(--spacing-xl);
}

.admin-menu {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
  gap: var(--spacing-lg);
  margin-top: var(--spacing-xl);
}

.admin-card {
  text-align: center;
  padding: var(--spacing-xl);
  transition: all var(--transition-base);
}

.admin-card:hover {
  transform: translateY(-4px) scale(1.01);
}

.admin-icon-wrap {
  width: 64px;
  height: 64px;
  border-radius: 20px;
  margin: 0 auto var(--spacing-md);
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, rgba(0, 180, 137, 0.12), rgba(0, 180, 137, 0.2));
}

.admin-icon {
  font-size: 30px;
  color: var(--primary-color);
}

.admin-card h3 {
  margin-bottom: var(--spacing-sm);
}

.admin-card p {
  color: var(--text-secondary);
  margin-bottom: var(--spacing-lg);
}
</style>
