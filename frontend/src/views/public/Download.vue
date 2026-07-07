<template>
  <div class="download-page">
    <!-- Main Background Decorative Elements -->
    <div class="bg-glow bg-glow-primary"></div>
    <div class="bg-glow bg-glow-secondary"></div>
    <div class="bg-glow bg-glow-tertiary"></div>

    <div class="container">
      <!-- Header App Showcase -->
      <section class="hero-section glass-panel">
        <div class="hero-grid">
          <!-- Logo & Mockup Presentation -->
          <div class="showcase-card">
            <div class="phone-mockup">
              <div class="phone-notch"></div>
              <div class="phone-screen">
                <div class="mock-header">
                  <span class="mock-time">10:54</span>
                  <div class="mock-status-icons">📶 🔋</div>
                </div>
                <div class="mock-app-content">
                  <div class="mock-logo-wrap">
                    <img v-if="showLogo" :src="logoUrl" alt="Logo" class="mock-logo" @error="onLogoError">
                    <div v-else class="mock-logo-fallback">🤖</div>
                  </div>
                  <div class="mock-app-title">智享社区</div>
                  <div class="mock-app-subtitle">Smart Community</div>
                  <div class="mock-widget">
                    <div class="mock-widget-title">智能管家已就绪</div>
                    <div class="mock-widget-desc">“今天已为您总结 3 条社区公告”</div>
                  </div>
                  <div class="mock-menu-grid">
                    <span class="mock-menu-item">缴费</span>
                    <span class="mock-menu-item">报修</span>
                    <span class="mock-menu-item">访客</span>
                    <span class="mock-menu-item">商城</span>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- Description & Action Area -->
          <div class="info-card">
            <span class="app-badge">Official Release v2.0.0</span>
            <h1 class="main-title">智享生活 App</h1>
            <p class="main-desc">
              重塑智慧社区体验。一键畅享在线物业缴费、报修投诉办理、访客登记以及社区便利商超购物，更有强大的 AI 智能生活管家，实时为您总结社区资讯，打理日常业务。
            </p>

            <!-- Download Details Card -->
            <div class="download-info-box">
              <div class="info-row">
                <span class="info-label">应用大小</span>
                <span class="info-val">19.0 MB</span>
              </div>
              <div class="info-row">
                <span class="info-label">适用平台</span>
                <span class="info-val">Android 5.0+</span>
              </div>
              <div class="info-row">
                <span class="info-label">更新时间</span>
                <span class="info-val">今日最新发布</span>
              </div>
            </div>

            <!-- Download Button -->
            <div class="action-row">
              <button class="btn btn-download" @click="handleDownload">
                <span class="icon">📥</span>
                <span class="btn-text">立即下载 (Android APK)</span>
              </button>
            </div>

            <p class="meta-url">下载源: <code>{{ displayUrlShort }}</code></p>
          </div>
        </div>
      </section>

      <!-- Features Showcase Section -->
      <section class="features-section">
        <h2 class="section-title">核心功能一览</h2>
        <div class="features-grid">
          <div class="feature-card glass-panel-light">
            <div class="feature-icon icon-blue">💳</div>
            <h3>在线缴费</h3>
            <p>物业费、停车费快捷账单，支持混合积分抵扣与账户余额秒级支付，账单明细清晰完整。</p>
          </div>

          <div class="feature-card glass-panel-light">
            <div class="feature-icon icon-green">🛠️</div>
            <h3>物业报修</h3>
            <p>水电设施、公共区域损坏一键拍照提报，处理进度实时通知，物业响应更迅速。</p>
          </div>

          <div class="feature-card glass-panel-light">
            <div class="feature-icon icon-purple">🚗</div>
            <h3>智能停车</h3>
            <p>名下车位快捷绑定，支持线上输入车牌号随时修改绑定车辆，出行缴费轻松自如。</p>
          </div>

          <div class="feature-card glass-panel-light">
            <div class="feature-icon icon-orange">🛒</div>
            <h3>社区商超</h3>
            <p>周边便利商品一应俱全，购物车一键结算，支持挑选就近服务门店，自提配送随心选。</p>
          </div>

          <div class="feature-card glass-panel-light">
            <div class="feature-icon icon-teal">🤖</div>
            <h3>AI 日报与助手</h3>
            <p>智能 AI 管家全天候在线，随时提炼总结社区公告；支持语音指令与自动业务导航。</p>
          </div>
        </div>
      </section>

      <!-- Guides & Instructions -->
      <section class="guides-section glass-panel">
        <h2 class="section-title">安装引导说明</h2>
        <div class="guides-grid">
          <div class="guide-step">
            <span class="step-num">01</span>
            <h4>点击下载</h4>
            <p>点击上方的“立即下载”按钮，浏览器会自动下载 APK 安装包。</p>
          </div>
          <div class="guide-step">
            <span class="step-num">02</span>
            <h4>授权安装</h4>
            <p>如果系统提示“未知来源安全风险”，请选择允许在此浏览器中进行应用安装（仅首次需要）。</p>
          </div>
          <div class="guide-step">
            <span class="step-num">03</span>
            <h4>登录畅享</h4>
            <p>安装完成后，使用您在社区绑定的手机号登录，即可体验多端联动的生活服务。</p>
          </div>
        </div>
      </section>

      <!-- Footer Note -->
      <footer class="footer-note">
        <p>© 智享生活社区 · App 官方分发平台</p>
      </footer>
    </div>
  </div>
</template>

<script setup>
import { ElMessage } from 'element-plus'
import { ref, computed } from 'vue'

const downloadUrl = import.meta.env.VITE_APP_DOWNLOAD_URL || 'https://communitysvc.xyz/smart-community/app/__UNI__036BB94__20260706193449.apk'
const logoUrl = '/images/logo.png'
const showLogo = ref(true)

const displayUrlShort = computed(() => {
  try {
    const url = new URL(downloadUrl)
    return `${url.origin}${url.pathname.slice(0, 20)}.../${url.pathname.split('/').pop()}`
  } catch (e) {
    return downloadUrl
  }
})

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
  position: relative;
  min-height: 100vh;
  padding: 80px 24px;
  background-color: #0b1329;
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
  color: #f1f5f9;
  overflow: hidden;
}

/* Background Glowing Orbs */
.bg-glow {
  position: absolute;
  border-radius: 50%;
  filter: blur(120px);
  z-index: 0;
  opacity: 0.15;
  pointer-events: none;
}

.bg-glow-primary {
  width: 500px;
  height: 500px;
  background: #3b82f6;
  top: -10%;
  left: -10%;
}

.bg-glow-secondary {
  width: 600px;
  height: 600px;
  background: #10b981;
  bottom: -20%;
  right: -10%;
}

.bg-glow-tertiary {
  width: 400px;
  height: 400px;
  background: #8b5cf6;
  top: 40%;
  left: 60%;
}

.container {
  position: relative;
  max-width: 1000px;
  margin: 0 auto;
  z-index: 1;
}

/* Glassmorphism Styles */
.glass-panel {
  background: rgba(17, 24, 39, 0.6);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 24px;
  padding: 40px;
  backdrop-filter: blur(12px);
  box-shadow: 0 20px 50px rgba(0, 0, 0, 0.3);
}

.glass-panel-light {
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.05);
  border-radius: 20px;
  padding: 24px;
  backdrop-filter: blur(8px);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

/* Hero Section Layout */
.hero-grid {
  display: grid;
  grid-template-columns: 320px 1fr;
  gap: 48px;
  align-items: center;
}

/* Phone Mockup Presentation */
.showcase-card {
  display: flex;
  justify-content: center;
  align-items: center;
}

.phone-mockup {
  position: relative;
  width: 250px;
  height: 500px;
  background: #1e293b;
  border: 8px solid #334155;
  border-radius: 36px;
  box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.5), 0 0 40px rgba(59, 130, 246, 0.2);
  overflow: hidden;
}

.phone-notch {
  position: absolute;
  top: 0;
  left: 50%;
  transform: translateX(-50%);
  width: 110px;
  height: 18px;
  background: #334155;
  border-bottom-left-radius: 12px;
  border-bottom-right-radius: 12px;
  z-index: 10;
}

.phone-screen {
  position: relative;
  width: 100%;
  height: 100%;
  background: linear-gradient(135deg, #0f172a 0%, #1e1b4b 100%);
  padding: 16px;
  box-sizing: border-box;
}

.mock-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 10px;
  color: #64748b;
  margin-top: 4px;
  margin-bottom: 30px;
}

.mock-app-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
}

.mock-logo-wrap {
  width: 76px;
  height: 76px;
  background: linear-gradient(135deg, #3b82f6, #10b981);
  border-radius: 18px;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 8px 20px rgba(59, 130, 246, 0.3);
  margin-bottom: 14px;
}

.mock-logo {
  width: 54px;
  height: 54px;
  object-fit: contain;
}

.mock-logo-fallback {
  font-size: 32px;
}

.mock-app-title {
  font-size: 18px;
  font-weight: 700;
  color: #fff;
  letter-spacing: 0.5px;
}

.mock-app-subtitle {
  font-size: 9px;
  color: #3b82f6;
  text-transform: uppercase;
  letter-spacing: 1px;
  margin-top: 2px;
  margin-bottom: 24px;
}

.mock-widget {
  width: 100%;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 12px;
  padding: 12px;
  box-sizing: border-box;
  text-align: left;
  margin-bottom: 20px;
}

.mock-widget-title {
  font-size: 11px;
  font-weight: 600;
  color: #10b981;
}

.mock-widget-desc {
  font-size: 9px;
  color: #94a3b8;
  margin-top: 4px;
}

.mock-menu-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 8px;
  width: 100%;
}

.mock-menu-item {
  font-size: 9px;
  color: #cbd5e1;
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.05);
  border-radius: 8px;
  padding: 6px 0;
  text-align: center;
}

/* Info Side Panel Layout */
.info-card {
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.app-badge {
  align-self: flex-start;
  font-size: 11px;
  font-weight: 700;
  color: #10b981;
  background: rgba(16, 185, 129, 0.1);
  border: 1px solid rgba(16, 185, 129, 0.2);
  padding: 6px 14px;
  border-radius: 99px;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  margin-bottom: 16px;
}

.main-title {
  font-size: 38px;
  font-weight: 800;
  margin: 0 0 16px;
  background: linear-gradient(135deg, #ffffff 0%, #cbd5e1 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.main-desc {
  font-size: 15px;
  line-height: 1.7;
  color: #94a3b8;
  margin: 0 0 24px;
}

.download-info-box {
  background: rgba(255, 255, 255, 0.02);
  border: 1px solid rgba(255, 255, 255, 0.04);
  border-radius: 16px;
  padding: 16px 20px;
  margin-bottom: 24px;
}

.info-row {
  display: flex;
  justify-content: space-between;
  padding: 6px 0;
  font-size: 13px;
}

.info-row:not(:last-child) {
  border-bottom: 1px solid rgba(255, 255, 255, 0.04);
}

.info-label {
  color: #64748b;
}

.info-val {
  font-weight: 600;
  color: #cbd5e1;
}

/* Interactive Download Buttons */
.action-row {
  margin-bottom: 14px;
}

.btn {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
  width: 100%;
  padding: 16px 28px;
  font-size: 15px;
  font-weight: 700;
  border-radius: 14px;
  border: none;
  cursor: pointer;
  transition: all 0.25s ease;
}

.btn-download {
  color: #ffffff;
  background: linear-gradient(135deg, #3b82f6 0%, #1d4ed8 100%);
  box-shadow: 0 10px 25px rgba(59, 130, 246, 0.35), inset 0 1px 0 rgba(255, 255, 255, 0.2);
}

.btn-download:hover {
  transform: translateY(-2px);
  background: linear-gradient(135deg, #2563eb 0%, #1e40af 100%);
  box-shadow: 0 12px 30px rgba(59, 130, 246, 0.45), inset 0 1px 0 rgba(255, 255, 255, 0.3);
}

.btn-download:active {
  transform: translateY(0);
}

.meta-url {
  font-size: 11px;
  color: #475569;
  margin: 0;
}

.meta-url code {
  color: #64748b;
  background: rgba(255, 255, 255, 0.02);
  padding: 2px 6px;
  border-radius: 4px;
}

/* Features grid section */
.features-section {
  margin-top: 48px;
}

.section-title {
  font-size: 22px;
  font-weight: 800;
  color: #f1f5f9;
  text-align: center;
  margin: 0 0 28px;
}

.features-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 20px;
}

.feature-card:hover {
  transform: translateY(-4px);
  background: rgba(255, 255, 255, 0.05);
  border-color: rgba(255, 255, 255, 0.1);
  box-shadow: 0 12px 24px rgba(0, 0, 0, 0.2);
}

.feature-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 44px;
  height: 44px;
  border-radius: 12px;
  font-size: 20px;
  margin-bottom: 16px;
}

.icon-blue { background: rgba(59, 130, 246, 0.15); }
.icon-green { background: rgba(16, 185, 129, 0.15); }
.icon-purple { background: rgba(139, 92, 246, 0.15); }
.icon-orange { background: rgba(249, 115, 22, 0.15); }
.icon-teal { background: rgba(20, 184, 166, 0.15); }

.feature-card h3 {
  font-size: 16px;
  font-weight: 700;
  margin: 0 0 8px;
  color: #fff;
}

.feature-card p {
  font-size: 13px;
  line-height: 1.6;
  color: #64748b;
  margin: 0;
}

/* Install Guide Section */
.guides-section {
  margin-top: 48px;
}

.guides-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 32px;
}

.guide-step {
  display: flex;
  flex-direction: column;
}

.step-num {
  font-size: 32px;
  font-weight: 900;
  color: rgba(59, 130, 246, 0.2);
  line-height: 1;
  margin-bottom: 12px;
}

.guide-step h4 {
  font-size: 15px;
  font-weight: 700;
  color: #fff;
  margin: 0 0 6px;
}

.guide-step p {
  font-size: 13px;
  line-height: 1.6;
  color: #64748b;
  margin: 0;
}

/* Footer Section */
.footer-note {
  margin-top: 60px;
  text-align: center;
  font-size: 12px;
  color: #475569;
}

/* Responsive adjustments */
@media (max-width: 800px) {
  .hero-grid {
    grid-template-columns: 1fr;
    gap: 36px;
  }

  .glass-panel {
    padding: 24px;
  }

  .guides-grid {
    grid-template-columns: 1fr;
    gap: 24px;
  }

  .main-title {
    font-size: 30px;
  }
}
</style>
