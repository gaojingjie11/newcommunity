<template>
  <div class="download-page">
    <!-- Import Premium Typography -->
    <style>
      @import url('https://fonts.googleapis.com/css2?family=Plus+Jakarta+Sans:wght@300;400;500;600;700;800&display=swap');
    </style>

    <!-- Ambient Glowing Highlights -->
    <div class="ambient-glow glow-1"></div>
    <div class="ambient-glow glow-2"></div>
    <div class="ambient-glow glow-3"></div>

    <!-- Tech Grid Overlay -->
    <div class="grid-overlay"></div>

    <div class="page-container">
      <!-- Top Navigation Header -->
      <header class="header">
        <div class="brand">
          <div class="logo-icon">
            <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
              <path d="M12 2L2 7L12 12L22 7L12 2Z" fill="url(#logo-grad-1)"/>
              <path d="M2 17L12 22L22 17" stroke="url(#logo-grad-2)" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
              <path d="M2 12L12 17L22 12" stroke="url(#logo-grad-3)" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
              <defs>
                <linearGradient id="logo-grad-1" x1="2" y1="2" x2="22" y2="12" gradientUnits="userSpaceOnUse">
                  <stop stop-color="#00f2fe"/>
                  <stop offset="1" stop-color="#4facfe"/>
                </linearGradient>
                <linearGradient id="logo-grad-2" x1="2" y1="17" x2="22" y2="22" gradientUnits="userSpaceOnUse">
                  <stop stop-color="#00f2fe"/>
                  <stop offset="1" stop-color="#0000ff"/>
                </linearGradient>
                <linearGradient id="logo-grad-3" x1="2" y1="12" x2="22" y2="17" gradientUnits="userSpaceOnUse">
                  <stop stop-color="#00f2fe"/>
                  <stop offset="1" stop-color="#4facfe"/>
                </linearGradient>
              </defs>
            </svg>
          </div>
          <span class="brand-name">智享社区</span>
        </div>
        <div class="header-badge">
          <span class="status-dot"></span>
          <span>系统在线 · 已安全加密</span>
        </div>
      </header>

      <!-- Hero Showcase Section -->
      <main class="hero-section">
        <div class="hero-grid">
          <!-- Text Info Column -->
          <div class="info-content">
            <div class="tagline-badge">
              <span class="badge-accent">NEW</span>
              <span class="badge-text">v2.0.0 智慧重构版</span>
            </div>
            
            <h1 class="hero-title">
              智享生活 <span class="gradient-text">触手可及</span>
            </h1>
            
            <p class="hero-desc">
              深度融合 AI 人工智能与智慧社区物联网。一个客户端集成在线缴费、报修审批、智能车牌变更及商城购物，更搭载全局 AI 助手，为您量身定制每日社区简报与快讯。
            </p>

            <!-- Specification Cards -->
            <div class="specs-wrapper">
              <div class="spec-item">
                <span class="spec-label">文件大小</span>
                <span class="spec-value">19.0 MB</span>
              </div>
              <div class="spec-item">
                <span class="spec-label">支持平台</span>
                <span class="spec-value">Android 5.0+</span>
              </div>
              <div class="spec-item">
                <span class="spec-label">检测状态</span>
                <span class="spec-value text-green">🛡️ 官方绿标无毒</span>
              </div>
            </div>

            <!-- Download Panel -->
            <div class="download-panel">
              <div class="download-actions">
                <button class="btn-glow-action" @click="handleDownload">
                  <span class="btn-glow-content">
                    <svg class="download-icon" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                      <path d="M12 15V3M12 15L8 11M12 15L16 11M2 17V20C2 20.5304 2.21071 21.0391 2.58579 21.4142C2.96086 21.7893 3.46957 22 4 22H20C20.5304 22 21.0391 21.7893 21.4142 21.4142C21.7893 21.0391 22 20.5304 22 20V17" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"/>
                    </svg>
                    立即下载 Android 安装包
                  </span>
                </button>
              </div>

              <!-- QR Code Sub-Card -->
              <div class="qr-card">
                <div class="qr-container">
                  <!-- Futuristic QR Design with SVGs -->
                  <svg viewBox="0 0 100 100" class="qr-svg">
                    <rect x="10" y="10" width="25" height="25" fill="none" stroke="#3b82f6" stroke-width="3"/>
                    <rect x="15" y="15" width="15" height="15" fill="#3b82f6"/>
                    <rect x="65" y="10" width="25" height="25" fill="none" stroke="#3b82f6" stroke-width="3"/>
                    <rect x="70" y="15" width="15" height="15" fill="#3b82f6"/>
                    <rect x="10" y="65" width="25" height="25" fill="none" stroke="#3b82f6" stroke-width="3"/>
                    <rect x="15" y="70" width="15" height="15" fill="#3b82f6"/>
                    
                    <!-- Random bits -->
                    <rect x="42" y="15" width="5" height="5" fill="#10b981"/>
                    <rect x="50" y="25" width="10" height="5" fill="#cbd5e1"/>
                    <rect x="45" y="45" width="10" height="10" fill="#3b82f6"/>
                    <rect x="65" y="45" width="5" height="5" fill="#10b981"/>
                    <rect x="75" y="55" width="10" height="10" fill="#3b82f6"/>
                    <rect x="50" y="75" width="5" height="15" fill="#cbd5e1"/>
                    <rect x="80" y="75" width="10" height="5" fill="#10b981"/>
                  </svg>
                  <div class="qr-scan-line"></div>
                </div>
                <div class="qr-text">
                  <p class="qr-title">手机扫码下载</p>
                  <p class="qr-desc">支持任意扫码工具</p>
                </div>
              </div>
            </div>

            <!-- Version Checksum -->
            <div class="checksum-box">
              <span class="checksum-title">安装包直链 (CDN加速):</span>
              <code class="checksum-code">{{ downloadUrl }}</code>
            </div>
          </div>

          <!-- Phone Showcase Column (Isometric Float style) -->
          <div class="showcase-isometric">
            <div class="isometric-wrapper">
              <!-- Floating Card 1: Main Dashboard -->
              <div class="float-card card-main">
                <div class="card-glow"></div>
                <div class="card-header-app">
                  <div class="app-mini-icon">🤖</div>
                  <div>
                    <div class="app-mini-name">AI 生活管家</div>
                    <div class="app-mini-state">在线提供业务支持</div>
                  </div>
                </div>
                <div class="card-body-app">
                  <div class="ai-speech-bubble">
                    “您好，已成功为您绑定车牌号 <strong>粤A88888</strong>，当前绑定车位：A区-1025。”
                  </div>
                  <div class="app-mini-chart">
                    <div class="chart-bar" style="height: 40%"></div>
                    <div class="chart-bar" style="height: 70%"></div>
                    <div class="chart-bar active" style="height: 95%"></div>
                    <div class="chart-bar" style="height: 55%"></div>
                  </div>
                </div>
              </div>

              <!-- Floating Card 2: Quick Features -->
              <div class="float-card card-features">
                <div class="card-glow"></div>
                <div class="quick-header">快捷缴费</div>
                <div class="fee-item">
                  <span>物业管理费 (7月)</span>
                  <span class="fee-price">¥320.00</span>
                </div>
                <div class="fee-action-btn">一键安全支付</div>
              </div>

              <!-- Mock Phone Outline background glow -->
              <div class="phone-backlight"></div>
            </div>
          </div>
        </div>
      </main>

      <!-- Advanced Features Section -->
      <section class="section-features">
        <h2 class="grid-title">重新定义社区服务</h2>
        <div class="features-grid">
          <!-- Feature Card 1 -->
          <div class="feature-item-card">
            <div class="feat-icon-box bg-blue">
              <span class="feat-icon">💳</span>
            </div>
            <h3>聚合支付</h3>
            <p>融合物业费、商超订单多场景，内置积分抵扣算法与余额安全核算，支付得心应手。</p>
          </div>

          <!-- Feature Card 2 -->
          <div class="feature-item-card">
            <div class="feat-icon-box bg-purple">
              <span class="feat-icon">🚗</span>
            </div>
            <h3>车位绑定</h3>
            <p>支持在线自助绑定或修改车牌，符合一省简称一字母标准格式车牌校验，极速审批。</p>
          </div>

          <!-- Feature Card 3 -->
          <div class="feature-item-card">
            <div class="feat-icon-box bg-green">
              <span class="feat-icon">🤖</span>
            </div>
            <h3>AI 智能日报</h3>
            <p>由大语言模型驱动，每日智能汇总社区公告与账单异动，为您提供精炼的摘要卡片。</p>
          </div>

          <!-- Feature Card 4 -->
          <div class="feature-item-card">
            <div class="feat-icon-box bg-orange">
              <span class="feat-icon">🛒</span>
            </div>
            <h3>社区商超</h3>
            <p>打通本地便利商超与自提门店，支持查看门店列表及商品规格挑选，即买即达。</p>
          </div>
        </div>
      </section>

      <!-- Step Stepper Guides -->
      <section class="section-guides">
        <h2 class="grid-title">极速安装指引</h2>
        <div class="steps-timeline">
          <div class="timeline-node">
            <div class="node-circle">1</div>
            <div class="node-content">
              <h4>获取安装包</h4>
              <p>点击上方下载按钮，系统将自动分配最快节点下载 APK 包。</p>
            </div>
          </div>
          <div class="timeline-node">
            <div class="node-circle">2</div>
            <div class="node-content">
              <h4>信任应用源</h4>
              <p>若弹出未知源风险拦截，请点击“信任此来源”并允许应用安装。</p>
            </div>
          </div>
          <div class="timeline-node">
            <div class="node-circle">3</div>
            <div class="node-content">
              <h4>配置与使用</h4>
              <p>登录您的注册手机号即可自动获取您的业主身份与绑定信息。</p>
            </div>
          </div>
        </div>
      </section>

      <!-- Premium Footer -->
      <footer class="footer">
        <p>© 智享生活社区 · 智慧云分发管理系统</p>
        <p class="footer-sub">安全等级：TLS 1.3 | SHA256 证书校验已启用</p>
      </footer>
    </div>
  </div>
</template>

<script setup>
import { ElMessage } from 'element-plus'
import { ref } from 'vue'

const downloadUrl = import.meta.env.VITE_APP_DOWNLOAD_URL || 'https://communitysvc.xyz/smart-community/app/__UNI__036BB94__20260706193449.apk'

function handleDownload() {
  if (!downloadUrl || downloadUrl === '#') {
    ElMessage.warning('下载地址未配置，请联系管理员')
    return
  }
  window.open(downloadUrl, '_blank', 'noopener')
}
</script>

<style scoped>
/* Reset and Theme Definition */
.download-page {
  position: relative;
  min-height: 100vh;
  background-color: #07090e;
  color: #e2e8f0;
  font-family: 'Plus Jakarta Sans', -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
  padding: 60px 24px;
  box-sizing: border-box;
  overflow: hidden;
}

/* Background grid overlay */
.grid-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-image: 
    linear-gradient(rgba(255, 255, 255, 0.015) 1px, transparent 1px),
    linear-gradient(90deg, rgba(255, 255, 255, 0.015) 1px, transparent 1px);
  background-size: 40px 40px;
  background-position: center;
  z-index: 1;
  pointer-events: none;
}

/* Ambient Glowing Background Elements */
.ambient-glow {
  position: absolute;
  border-radius: 50%;
  filter: blur(140px);
  opacity: 0.22;
  pointer-events: none;
  z-index: 0;
}

.glow-1 {
  width: 500px;
  height: 500px;
  background: radial-gradient(circle, #3b82f6 0%, transparent 80%);
  top: -10%;
  left: -5%;
}

.glow-2 {
  width: 600px;
  height: 600px;
  background: radial-gradient(circle, #10b981 0%, transparent 80%);
  bottom: -15%;
  right: -5%;
}

.glow-3 {
  width: 400px;
  height: 400px;
  background: radial-gradient(circle, #8b5cf6 0%, transparent 80%);
  top: 35%;
  left: 55%;
}

.page-container {
  position: relative;
  max-width: 1100px;
  margin: 0 auto;
  z-index: 2;
}

/* Header Section */
.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 60px;
}

.brand {
  display: flex;
  align-items: center;
  gap: 12px;
}

.logo-icon {
  width: 32px;
  height: 32px;
}

.brand-name {
  font-size: 20px;
  font-weight: 800;
  letter-spacing: 0.5px;
  background: linear-gradient(135deg, #fff 0%, #cbd5e1 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.header-badge {
  display: flex;
  align-items: center;
  gap: 8px;
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.06);
  padding: 6px 14px;
  border-radius: 99px;
  font-size: 12px;
  color: #94a3b8;
}

.status-dot {
  width: 6px;
  height: 6px;
  background-color: #10b981;
  border-radius: 50%;
  box-shadow: 0 0 10px #10b981;
}

/* Main Hero Layout */
.hero-section {
  background: rgba(10, 15, 30, 0.5);
  border: 1px solid rgba(255, 255, 255, 0.05);
  border-radius: 28px;
  padding: 50px;
  backdrop-filter: blur(20px);
  box-shadow: 0 30px 60px rgba(0, 0, 0, 0.4), inset 0 1px 0 rgba(255, 255, 255, 0.05);
  margin-bottom: 50px;
}

.hero-grid {
  display: grid;
  grid-template-columns: 1fr 400px;
  gap: 60px;
  align-items: center;
}

/* Info Column Content */
.tagline-badge {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  background: rgba(59, 130, 246, 0.1);
  border: 1px solid rgba(59, 130, 246, 0.18);
  border-radius: 99px;
  padding: 4px 12px;
  margin-bottom: 24px;
}

.badge-accent {
  font-size: 10px;
  font-weight: 800;
  color: #fff;
  background: #3b82f6;
  padding: 2px 6px;
  border-radius: 99px;
}

.badge-text {
  font-size: 11px;
  font-weight: 700;
  color: #60a5fa;
  letter-spacing: 0.5px;
}

.hero-title {
  font-size: 46px;
  font-weight: 800;
  line-height: 1.15;
  margin: 0 0 20px;
  letter-spacing: -1px;
}

.gradient-text {
  background: linear-gradient(135deg, #00f2fe 0%, #4facfe 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.hero-desc {
  font-size: 16px;
  line-height: 1.75;
  color: #94a3b8;
  margin: 0 0 32px;
}

/* Specifications Grid */
.specs-wrapper {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 16px;
  background: rgba(255, 255, 255, 0.02);
  border: 1px solid rgba(255, 255, 255, 0.04);
  border-radius: 16px;
  padding: 16px 20px;
  margin-bottom: 36px;
}

.spec-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.spec-label {
  font-size: 11px;
  color: #64748b;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.spec-value {
  font-size: 14px;
  font-weight: 700;
  color: #cbd5e1;
}

.text-green {
  color: #10b981;
}

/* Download Actions Panel */
.download-panel {
  display: flex;
  gap: 24px;
  align-items: center;
  margin-bottom: 24px;
}

.download-actions {
  flex-grow: 1;
}

.btn-glow-action {
  position: relative;
  width: 100%;
  padding: 3px;
  border-radius: 16px;
  background: linear-gradient(90deg, #00c6ff, #0072ff, #00c6ff);
  background-size: 200% auto;
  border: none;
  cursor: pointer;
  outline: none;
  transition: all 0.3s ease;
  box-shadow: 0 10px 30px rgba(0, 114, 255, 0.3);
}

.btn-glow-action:hover {
  background-position: right center;
  transform: translateY(-2px);
  box-shadow: 0 12px 35px rgba(0, 114, 255, 0.45);
}

.btn-glow-content {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  background: #0d1527;
  color: #fff;
  font-size: 15px;
  font-weight: 700;
  padding: 14px 20px;
  border-radius: 13px;
  transition: background 0.3s ease;
}

.btn-glow-action:hover .btn-glow-content {
  background: transparent;
}

.download-icon {
  width: 18px;
  height: 18px;
}

/* Futuristic QR Code Card */
.qr-card {
  display: flex;
  align-items: center;
  gap: 16px;
  background: rgba(255, 255, 255, 0.02);
  border: 1px solid rgba(255, 255, 255, 0.05);
  border-radius: 16px;
  padding: 12px 20px;
  width: 200px;
  flex-shrink: 0;
}

.qr-container {
  position: relative;
  width: 50px;
  height: 50px;
  background: rgba(255, 255, 255, 0.04);
  border-radius: 8px;
  padding: 4px;
  box-sizing: border-box;
  overflow: hidden;
}

.qr-svg {
  width: 100%;
  height: 100%;
}

.qr-scan-line {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 2px;
  background: #10b981;
  box-shadow: 0 0 8px #10b981;
  animation: scan 2s linear infinite;
}

@keyframes scan {
  0% { top: 0%; }
  50% { top: 100%; }
  100% { top: 0%; }
}

.qr-text {
  display: flex;
  flex-direction: column;
}

.qr-title {
  font-size: 12px;
  font-weight: 700;
  color: #fff;
  margin: 0;
}

.qr-desc {
  font-size: 10px;
  color: #64748b;
  margin: 2px 0 0;
}

/* Checksum box */
.checksum-box {
  background: rgba(0, 0, 0, 0.15);
  border: 1px solid rgba(255, 255, 255, 0.03);
  border-radius: 12px;
  padding: 10px 16px;
  font-size: 11px;
}

.checksum-title {
  color: #475569;
  font-weight: 600;
  margin-right: 6px;
}

.checksum-code {
  color: #64748b;
  word-break: break-all;
  font-family: monospace;
}

/* Isometric Showcase Side Column */
.showcase-isometric {
  position: relative;
  display: flex;
  justify-content: center;
  align-items: center;
  height: 400px;
}

.isometric-wrapper {
  position: relative;
  width: 100%;
  height: 100%;
  transform: perspective(1000px) rotateY(-15deg) rotateX(10deg);
  transform-style: preserve-3d;
}

.float-card {
  position: absolute;
  background: rgba(17, 24, 39, 0.75);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 20px;
  backdrop-filter: blur(16px);
  padding: 20px;
  box-sizing: border-box;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.4), inset 0 1px 0 rgba(255, 255, 255, 0.1);
  animation: float 6s ease-in-out infinite;
}

.card-glow {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  border-radius: 20px;
  background: linear-gradient(135deg, rgba(59, 130, 246, 0.1), transparent 60%);
  pointer-events: none;
}

.card-main {
  width: 280px;
  top: 10%;
  left: 10%;
  z-index: 2;
}

.card-features {
  width: 220px;
  bottom: 10%;
  right: 5%;
  z-index: 3;
  animation-delay: -3s;
}

@keyframes float {
  0% { transform: translateY(0px); }
  50% { transform: translateY(-12px); }
  100% { transform: translateY(0px); }
}

/* Inside Mockup Cards design */
.card-header-app {
  display: flex;
  align-items: center;
  gap: 12px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
  padding-bottom: 12px;
  margin-bottom: 12px;
}

.app-mini-icon {
  font-size: 24px;
}

.app-mini-name {
  font-size: 13px;
  font-weight: 700;
  color: #fff;
}

.app-mini-state {
  font-size: 10px;
  color: #10b981;
}

.ai-speech-bubble {
  font-size: 11px;
  line-height: 1.5;
  color: #94a3b8;
  background: rgba(255, 255, 255, 0.02);
  border-radius: 8px;
  padding: 8px 12px;
  margin-bottom: 12px;
}

.app-mini-chart {
  display: flex;
  justify-content: space-between;
  align-items: flex-end;
  height: 40px;
  padding: 0 10px;
}

.chart-bar {
  width: 32px;
  background: rgba(255, 255, 255, 0.06);
  border-radius: 4px;
}

.chart-bar.active {
  background: linear-gradient(180deg, #60a5fa, #3b82f6);
}

.quick-header {
  font-size: 13px;
  font-weight: 700;
  color: #fff;
  margin-bottom: 12px;
}

.fee-item {
  display: flex;
  justify-content: space-between;
  font-size: 11px;
  color: #94a3b8;
  margin-bottom: 16px;
}

.fee-price {
  font-weight: 700;
  color: #fb7185;
}

.fee-action-btn {
  font-size: 11px;
  font-weight: 700;
  color: #fff;
  background: #3b82f6;
  text-align: center;
  padding: 8px 0;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(59, 130, 246, 0.3);
}

/* Backlight Glow for Isometric Phone */
.phone-backlight {
  position: absolute;
  top: 20%;
  left: 20%;
  width: 250px;
  height: 250px;
  background: radial-gradient(circle, rgba(59, 130, 246, 0.15) 0%, transparent 70%);
  z-index: 1;
  pointer-events: none;
}

/* Premium Section features */
.section-features {
  margin-top: 60px;
  margin-bottom: 60px;
}

.grid-title {
  font-size: 24px;
  font-weight: 800;
  color: #fff;
  text-align: center;
  margin-bottom: 36px;
}

.features-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 20px;
}

.feature-item-card {
  background: rgba(255, 255, 255, 0.02);
  border: 1px solid rgba(255, 255, 255, 0.04);
  border-radius: 20px;
  padding: 24px;
  transition: all 0.3s ease;
}

.feature-item-card:hover {
  transform: translateY(-4px);
  background: rgba(255, 255, 255, 0.04);
  border-color: rgba(255, 255, 255, 0.08);
  box-shadow: 0 15px 30px rgba(0, 0, 0, 0.3);
}

.feat-icon-box {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 48px;
  height: 48px;
  border-radius: 12px;
  margin-bottom: 20px;
}

.bg-blue { background: rgba(59, 130, 246, 0.1); }
.bg-purple { background: rgba(139, 92, 246, 0.1); }
.bg-green { background: rgba(16, 185, 129, 0.1); }
.bg-orange { background: rgba(249, 115, 22, 0.1); }

.feat-icon {
  font-size: 22px;
}

.feature-item-card h3 {
  font-size: 16px;
  font-weight: 700;
  color: #fff;
  margin: 0 0 10px;
}

.feature-item-card p {
  font-size: 13px;
  line-height: 1.6;
  color: #64748b;
  margin: 0;
}

/* Timeline/Guides Section */
.section-guides {
  margin-top: 60px;
  margin-bottom: 60px;
}

.steps-timeline {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 30px;
  position: relative;
}

.timeline-node {
  position: relative;
  display: flex;
  flex-direction: column;
  background: rgba(255, 255, 255, 0.01);
  border: 1px solid rgba(255, 255, 255, 0.03);
  border-radius: 20px;
  padding: 24px;
}

.node-circle {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: rgba(59, 130, 246, 0.1);
  border: 1px solid rgba(59, 130, 246, 0.3);
  color: #3b82f6;
  font-size: 14px;
  font-weight: 800;
  margin-bottom: 16px;
}

.timeline-node h4 {
  font-size: 15px;
  font-weight: 700;
  color: #fff;
  margin: 0 0 8px;
}

.timeline-node p {
  font-size: 13px;
  line-height: 1.6;
  color: #64748b;
  margin: 0;
}

/* Footer note */
.footer {
  margin-top: 80px;
  text-align: center;
  border-top: 1px solid rgba(255, 255, 255, 0.04);
  padding-top: 24px;
}

.footer p {
  font-size: 12px;
  color: #475569;
  margin: 0;
}

.footer-sub {
  font-size: 10px;
  color: #334155 !important;
  margin-top: 6px !important;
}

/* Media Queries for Responsive Design */
@media (max-width: 960px) {
  .hero-grid {
    grid-template-columns: 1fr;
    gap: 40px;
  }

  .showcase-isometric {
    height: 350px;
    margin-top: 20px;
  }

  .features-grid {
    grid-template-columns: repeat(2, 1fr);
  }

  .steps-timeline {
    grid-template-columns: 1fr;
    gap: 20px;
  }

  .hero-section {
    padding: 30px;
  }

  .hero-title {
    font-size: 34px;
  }
}

@media (max-width: 600px) {
  .download-panel {
    flex-direction: column;
    align-items: stretch;
    gap: 16px;
  }

  .qr-card {
    width: auto;
  }

  .features-grid {
    grid-template-columns: 1fr;
  }

  .specs-wrapper {
    grid-template-columns: 1fr;
    gap: 12px;
  }
}
</style>
