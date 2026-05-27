<template>
  <div class="agent-page">
    <Navbar />

    <main class="container">
      <section class="page-header">
        <h1 class="page-title">Agent 智能服务</h1>
        <p class="page-desc">当前为微服务改造阶段预留入口，后续将统一接入 gateway-service 与 agent-service。</p>
      </section>

      <section class="agent-grid">
        <article v-for="item in agents" :key="item.name" class="agent-card">
          <div class="agent-name">{{ item.name }}</div>
          <p>{{ item.desc }}</p>
          <span class="agent-status">占位接口已规划</span>
        </article>
      </section>
    </main>
  </div>
</template>

<script setup>
import Navbar from '@/components/layout/Navbar.vue'

const agents = [
  { name: '社区客服 Agent', desc: '对应 POST /agent/chat，处理咨询和服务引导。' },
  { name: '报修派单 Agent', desc: '对应 POST /agent/repair-classify，输出分类、紧急度和建议部门。' },
  { name: '投诉风险 Agent', desc: '对应 POST /agent/complaint-risk，输出风险等级和处理建议。' },
  { name: '推荐 Agent', desc: '对应 POST /agent/recommend，输出商品或社区服务推荐。' }
]
</script>

<style scoped>
.agent-page {
  min-height: 100vh;
  background: #f8f9fa;
  padding-bottom: 80px;
}

.page-header {
  padding: 32px 0 24px;
}

.page-title {
  margin: 0 0 10px;
  color: #243447;
  font-size: 30px;
  font-weight: 700;
}

.page-desc {
  margin: 0;
  color: #6b7280;
  font-size: 15px;
}

.agent-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
  gap: 18px;
}

.agent-card {
  min-height: 150px;
  padding: 22px;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  background: #fff;
  box-shadow: 0 6px 20px rgba(15, 23, 42, 0.06);
}

.agent-name {
  margin-bottom: 10px;
  color: #1f2937;
  font-size: 18px;
  font-weight: 700;
}

.agent-card p {
  min-height: 48px;
  margin: 0 0 16px;
  color: #4b5563;
  line-height: 1.6;
}

.agent-status {
  display: inline-flex;
  align-items: center;
  height: 28px;
  padding: 0 10px;
  border-radius: 6px;
  background: #ecfdf5;
  color: #047857;
  font-size: 13px;
  font-weight: 600;
}
</style>
