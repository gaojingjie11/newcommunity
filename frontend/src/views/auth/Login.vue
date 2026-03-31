<template>
  <div class="login-page">
    <!-- 背景轮播层（预加载 + 淡入淡出，无闪动） -->
    <div class="bg-wrapper">
      <div
        class="bg-image"
        v-for="(img, index) in bgImages"
        :key="index"
        :style="{
          backgroundImage: `url(${img})`,
          opacity: currentBg === index ? 1 : 0
        }"
      ></div>
    </div>

    <div class="login-container">
      <div class="login-card card">
        <h2 class="login-title">智享生活社区服务平台</h2>
        <p class="login-subtitle">欢迎回来</p>

        <div class="tabs">
          <div
            class="tab-item"
            :class="{ active: loginType === 'password' }"
            @click="loginType = 'password'"
          >
            密码登录
          </div>
          <div
            class="tab-item"
            :class="{ active: loginType === 'code' }"
            @click="loginType = 'code'"
          >
            验证码登录
          </div>
        </div>

        <form
          v-if="loginType === 'password'"
          @submit.prevent="handleLogin"
          class="login-form"
        >
          <div class="form-group">
            <label>手机号</label>
            <input
              v-model="form.mobile"
              type="tel"
              class="input"
              placeholder="请输入手机号"
              required
            />
          </div>
          <div class="form-group">
            <label>密码</label>
            <input
              v-model="form.password"
              type="password"
              class="input"
              placeholder="请输入密码"
              required
            />
          </div>
          <div class="form-footer">
            <router-link to="/register" class="link"
              >还没有账号？注册</router-link
            >
          </div>
          <button
            type="submit"
            class="btn btn-primary btn-lg"
            :disabled="loading || isThrottled"
          >
            <span v-if="!loading">{{
              isThrottled ? "请稍候..." : "登录"
            }}</span>
            <span v-else class="loading"></span>
          </button>
        </form>

        <form v-else @submit.prevent="handleCodeLogin" class="login-form">
          <div class="form-group">
            <label>手机号</label>
            <input
              v-model="codeForm.mobile"
              type="tel"
              class="input"
              placeholder="请输入手机号"
              required
            />
          </div>
          <div class="form-group">
            <label>验证码</label>
            <div class="code-input-group">
              <input
                v-model="codeForm.code"
                type="text"
                class="input"
                placeholder="6位验证码"
                required
              />
              <button
                type="button"
                class="btn btn-secondary btn-code"
                :disabled="timer > 0"
                @click="sendSms"
              >
                {{ timer > 0 ? `${timer}s` : "发送验证码" }}
              </button>
            </div>
          </div>
          <div class="form-footer">
            <router-link to="/register" class="link"
              >还没有账号？注册</router-link
            >
          </div>
          <button
            type="submit"
            class="btn btn-primary btn-lg"
            :disabled="loading || isThrottled"
          >
            <span v-if="!loading">{{
              isThrottled ? "请稍候..." : "登录"
            }}</span>
            <span v-else class="loading"></span>
          </button>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from "vue";
import { useRouter } from "vue-router";
import { useUserStore } from "@/stores/user";
import { sendCode, loginByCode } from "@/api/auth";
import { ElMessage } from "element-plus";

const router = useRouter();
const userStore = useUserStore();

const loginType = ref("password");
const loading = ref(false);
const timer = ref(0);
const isThrottled = ref(false);

const form = ref({ mobile: "", password: "" });
const codeForm = ref({ mobile: "", code: "" });

// 轮播壁纸（预加载所有图片，无闪动）
const bgImages = [
  "https://communitysvc.xyz/community/show/denglu.png",
  "https://communitysvc.xyz/community/show/denglu2.png",
  "https://communitysvc.xyz/community/show/denglu3.png",
];
const currentBg = ref(0);
let bgTimer = null;

onMounted(() => {
  // 自动轮播
  bgTimer = setInterval(() => {
    currentBg.value = (currentBg.value + 1) % bgImages.length;
  }, 7000);
});

onUnmounted(() => {
  if (bgTimer) clearInterval(bgTimer);
});

// 密码登录
const handleLogin = async () => {
  if (isThrottled.value) return;
  if (!/^1[3-9]\d{9}$/.test(form.value.mobile)) {
    ElMessage.warning("请输入正确的11位手机号");
    return;
  }

  isThrottled.value = true;
  setTimeout(() => (isThrottled.value = false), 3000);

  loading.value = true;
  try {
    await userStore.login(form.value);
    ElMessage.success("登录成功！");
    router.push("/home");
  } catch (error) {
    ElMessage.error(error.response?.data?.message || "登录失败");
  } finally {
    loading.value = false;
  }
};

// 发送验证码
const sendSms = async () => {
  if (!codeForm.value.mobile) {
    ElMessage.warning("请输入手机号");
    return;
  }
  if (!/^1[3-9]\d{9}$/.test(codeForm.value.mobile)) {
    ElMessage.warning("请输入正确的11位手机号");
    return;
  }
  try {
    await sendCode({ mobile: codeForm.value.mobile });
    ElMessage.success("验证码已发送");
    timer.value = 60;
    const interval = setInterval(() => {
      timer.value--;
      if (timer.value <= 0) clearInterval(interval);
    }, 1000);
  } catch (e) {
    ElMessage.error(e.response?.data?.message || "发送失败");
  }
};

// 验证码登录
const handleCodeLogin = async () => {
  if (isThrottled.value) return;
  if (!/^1[3-9]\d{9}$/.test(codeForm.value.mobile)) {
    ElMessage.warning("请输入正确的11位手机号");
    return;
  }

  isThrottled.value = true;
  setTimeout(() => (isThrottled.value = false), 3000);

  loading.value = true;
  try {
    await userStore.loginByCode(codeForm.value);
    ElMessage.success("登录成功");
    router.push("/home");
  } catch (e) {
    ElMessage.error(e.message || e.response?.data?.msg || "登录失败");
  } finally {
    loading.value = false;
  }
};
</script>

<style scoped>
/* 页面容器 */
.login-page {
  position: relative;
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
}

/* 背景轮播容器 —— 关键：无闪动 */
.bg-wrapper {
  position: absolute;
  inset: 0;
  z-index: 0;
}
.bg-image {
  position: absolute;
  inset: 0;
  background-size: cover;
  background-position: center;
  background-repeat: no-repeat;
  transition: opacity 1.5s ease; /* 平滑淡入淡出 */
  will-change: opacity; /* 浏览器优化 */
}

/* 遮罩层 */
.login-page::before {
  content: "";
  position: absolute;
  inset: 0;
  background: rgba(128, 128, 128, 0.35);
  z-index: 1;
  pointer-events: none;
}

/* 登录内容层 */
.login-container {
  position: relative;
  z-index: 2;
  width: 100%;
  max-width: 400px;
  padding: var(--spacing-lg);
}

/* ✅ 核心：半透明毛玻璃登录卡片 */
.login-card {
  background: rgba(255, 255, 255, 0.174); /* 半透明白色 */
  backdrop-filter: blur(80px); /* 毛玻璃模糊 */
  -webkit-backdrop-filter: blur(12px);
  border: 1px solid rgba(255, 255, 255, 0.3);
  border-radius: 16px;
  padding: 32px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.15);
  animation: fadeIn 0.5s ease;
}

.login-title {
  font-size: var(--font-size-2xl);
  font-weight: 600;
  text-align: center;
  margin-bottom: var(--spacing-sm);
  color: #fff; /* 白色文字更通透 */
}

.login-subtitle {
  text-align: center;
  color: rgba(255, 255, 255, 0.85);
  margin-bottom: var(--spacing-xl);
}

.login-form {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-lg);
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-sm);
}

.form-group label {
  font-weight: 500;
  color: #fff;
}

/* 输入框也半透明 */
.input {
  background: rgba(255, 255, 255, 0.2);
  border: 1px solid rgba(255, 255, 255, 0.3);
  color: #fff;
  padding: 12px 16px;
  border-radius: 8px;
  outline: none;
}
.input::placeholder {
  color: rgba(255, 255, 255, 0.6);
}
.input:focus {
  border-color: rgba(255, 255, 255, 0.5);
  background: rgba(255, 255, 255, 0.25);
}

.form-footer {
  display: flex;
  justify-content: flex-end;
}

.link {
  color: rgba(255, 255, 255, 0.9);
  text-decoration: none;
  font-size: var(--font-size-sm);
}
.link:hover {
  text-decoration: underline;
  color: #fff;
}

.btn-lg {
  width: 100%;
}

/* 选项卡 */
.tabs {
  display: flex;
  justify-content: center;
  border-bottom: 1px solid rgba(255, 255, 255, 0.2);
  margin-bottom: 20px;
}
.tab-item {
  padding: 10px 20px;
  cursor: pointer;
  color: rgba(255, 255, 255, 0.7);
  border-bottom: 2px solid transparent;
}
.tab-item.active {
  color: #fff;
  border-bottom-color: #fff;
  font-weight: bold;
}

/* 验证码输入框 */
.code-input-group {
  display: flex;
  gap: 10px;
}
.btn-code {
  white-space: nowrap;
  width: 120px;
  background: rgba(255, 255, 255, 0.2);
  color: #fff;
  border: 1px solid rgba(255, 255, 255, 0.3);
}
.btn-code:disabled {
  background: rgba(255, 255, 255, 0.1);
  color: rgba(255, 255, 255, 0.5);
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}
</style>