<template>
  <div class="download-page">
    <section class="hero">
      <div class="hero-top">
        <div class="logo-card">
          <img v-if="showLogo" :src="logoUrl" alt="智享生活图标" class="app-logo" @error="onLogoError">
          <div v-else class="logo-fallback">智享生活</div>
        </div>

        <div class="hero-content">
          <p class="badge">智享生活 App</p>
          <h1>社区生活，一键在手</h1>
          <p class="desc">
            智享生活 App 聚合社区常用功能，覆盖缴费、服务、购物与消息互动，
            一个应用即可完成日常高频操作。
          </p>
          <p class="desc mt-8">
            本页面用于软件介绍与下载安装。
            可将该链接直接分享给用户作为统一下载入口。
          </p>

          <div class="actions">
            <button class="btn-primary" @click="handleDownload">立即下载</button>
          </div>

          <p class="meta">当前下载地址：{{ downloadUrl }}</p>

          <div class="quick-tags">
            <span>物业缴费</span>
            <span>报修访客</span>
            <span>积分抵扣</span>
            <span>商城购物</span>
            <span>订单管理</span>
            <span>AI 助手</span>
          </div>
        </div>
      </div>
    </section>

    <section class="panel">
      <h2>软件功能</h2>
      <div class="feature-grid">
        <article class="feature-card">
          <h3>在线缴费</h3>
          <p>支持物业费、订单支付和积分抵扣，常用缴费场景集中处理。</p>
        </article>
        <article class="feature-card">
          <h3>社区服务</h3>
          <p>报修、访客、停车、公告一站式办理，提升社区事务处理效率。</p>
        </article>
        <article class="feature-card">
          <h3>商城与订单</h3>
          <p>支持商品浏览下单、订单查询和状态跟踪，购物流程更清晰。</p>
        </article>
        <article class="feature-card">
          <h3>消息互动</h3>
          <p>支持社区消息沟通，及时获取通知并参与日常交流。</p>
        </article>
        <article class="feature-card">
          <h3>AI 智能助手</h3>
          <p>提供智能问答与建议推荐，帮助用户更快完成常见业务操作。</p>
        </article>
      </div>
    </section>

    <section class="panel">
      <h2>下载与安装说明</h2>
      <ul>
        <li>点击“立即下载”后将在新窗口打开安装包地址</li>
        <li>安卓设备请允许浏览器安装应用（仅首次需要）</li>
        <li>如下载失败，请切换网络后重试或联系管理员</li>
        <li>后续版本升级可复用同一下载页与链接进行发布</li>
      </ul>
    </section>

    <section class="footer-note">
      <p>© 智享生活 · 软件发布页</p>
    </section>
  </div>
</template>

<script setup>
import { ElMessage } from 'element-plus'
import { ref } from 'vue'

const downloadUrl = import.meta.env.VITE_APP_DOWNLOAD_URL || 'https://communitysvc.xyz/community/app/base.apk'
const logoUrl = '/images/logo.png'
const showLogo = ref(true)

function onLogoError() {
  showLogo.value = false
}

function handleDownload() {
  if (!downloadUrl || downloadUrl === '#') {
    ElMessage.warning('下载地址未配置，请联系管理员')
    return
  }
  window.open(downloadUrl, '_blank', 'noopener')
}
</script>

<style scoped>
.download-page {
  min-height: 100vh;
  padding: 56px 20px;
  background:
    radial-gradient(circle at 10% 10%, rgba(38, 140, 255, 0.12), transparent 36%),
    radial-gradient(circle at 90% 90%, rgba(0, 195, 132, 0.12), transparent 40%),
    #f7f9fc;
}

.hero,
.panel {
  max-width: 940px;
  margin: 0 auto;
  background: #fff;
  border: 1px solid #e6edf7;
  border-radius: 16px;
  padding: 28px;
  box-shadow: 0 12px 32px rgba(16, 24, 40, 0.08);
}

.panel {
  margin-top: 20px;
}

.hero-top {
  display: grid;
  grid-template-columns: 260px minmax(0, 1fr);
  gap: 24px;
  align-items: center;
}

.logo-card {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 260px;
  border-radius: 14px;
  border: 1px solid #d8e3f5;
  background: linear-gradient(180deg, #ffffff, #f4f7fb);
  padding: 16px;
}

.app-logo {
  width: 100%;
  max-width: 220px;
  height: auto;
  object-fit: contain;
}

.logo-fallback {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 180px;
  height: 180px;
  border-radius: 20px;
  color: #0f172a;
  font-weight: 700;
  letter-spacing: 2px;
  background: linear-gradient(135deg, #dbeafe, #dcfce7);
}

.badge {
  display: inline-block;
  margin: 0 0 12px;
  padding: 4px 10px;
  border-radius: 999px;
  background: #e8f3ff;
  color: #2563eb;
  font-size: 12px;
  font-weight: 600;
}

h1 {
  margin: 0 0 12px;
  font-size: 34px;
  color: #0f172a;
}

.desc {
  margin: 0;
  color: #475569;
  line-height: 1.7;
}

.mt-8 {
  margin-top: 8px;
}

.actions {
  display: flex;
  gap: 12px;
  margin-top: 20px;
  flex-wrap: wrap;
}

.btn-primary {
  border-radius: 12px;
  padding: 10px 18px;
  font-size: 14px;
  text-decoration: none;
  cursor: pointer;
}

.btn-primary {
  border: 0;
  color: #fff;
  background: linear-gradient(135deg, #1677ff, #00b578);
}

.meta {
  margin-top: 12px;
  color: #64748b;
  font-size: 13px;
  word-break: break-all;
}

.quick-tags {
  margin-top: 18px;
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.quick-tags span {
  font-size: 12px;
  color: #0f172a;
  background: #f1f5f9;
  border: 1px solid #dbeafe;
  border-radius: 999px;
  padding: 4px 10px;
}

h2 {
  margin: 0 0 12px;
  font-size: 20px;
  color: #0f172a;
}

h3 {
  margin: 0 0 8px;
  color: #0f172a;
  font-size: 16px;
}

ul {
  margin: 0;
  padding-left: 20px;
  color: #334155;
  line-height: 1.9;
}

.feature-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px;
}

.feature-card {
  border: 1px solid #e2e8f0;
  border-radius: 12px;
  padding: 14px;
  background: linear-gradient(180deg, #ffffff, #f8fafc);
}

.feature-card p {
  margin: 0;
  color: #475569;
  line-height: 1.7;
  font-size: 14px;
}

.footer-note {
  max-width: 940px;
  margin: 18px auto 0;
  text-align: center;
  color: #94a3b8;
  font-size: 13px;
}

@media (max-width: 860px) {
  .hero-top {
    grid-template-columns: 1fr;
    gap: 16px;
  }

  .logo-card {
    min-height: 190px;
  }

  .app-logo {
    max-width: 170px;
  }

  .feature-grid {
    grid-template-columns: 1fr;
  }

  h1 {
    font-size: 28px;
  }
}
</style>
