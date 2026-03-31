<template>
  <div class="login-page">
    <!-- 背景轮播层（和登录页相同） -->
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
        <h2 class="login-title">创建账号</h2>
        <p class="login-subtitle">加入智慧社区</p>
        
        <form @submit.prevent="handleRegister" class="login-form">
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
            <label>真实姓名</label>
            <input 
              v-model="form.real_name" 
              type="text" 
              class="input" 
              placeholder="请输入真实姓名"
              required
            />
          </div>
          
          <div class="form-group">
            <label>密码</label>
            <input 
              v-model="form.password" 
              type="password" 
              class="input" 
              placeholder="请输入密码（至少6位）"
              required
              minlength="6"
            />
          </div>
          
          <div class="form-footer">
            <router-link to="/login" class="link">已有账号？登录</router-link>
          </div>
          
          <button type="submit" class="btn btn-primary btn-lg" :disabled="loading">
            <span v-if="!loading">注册</span>
            <span v-else class="loading"></span>
          </button>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { ElMessage } from 'element-plus'

const router = useRouter()
const userStore = useUserStore()

const form = ref({
  mobile: '',
  real_name: '',
  password: '',
  username: '',
  age: 0,
  gender: 0
})

const loading = ref(false)

// 轮播背景（和登录页一致）
const bgImages = [
  "https://communitysvc.xyz/community/show/denglu.png",
  "https://communitysvc.xyz/community/show/denglu2.png",
  "https://communitysvc.xyz/community/show/denglu3.png",
];
const currentBg = ref(0);
let bgTimer = null;

onMounted(() => {
  bgTimer = setInterval(() => {
    currentBg.value = (currentBg.value + 1) % bgImages.length;
  }, 7000);
});

onUnmounted(() => {
  if (bgTimer) clearInterval(bgTimer);
});

const handleRegister = async () => {
  if (!form.value.username) {
    form.value.username = 'user_' + form.value.mobile.slice(-6)
  }
  
  loading.value = true
  try {
    await userStore.register(form.value)
    ElMessage.success('注册成功！请登录')
    router.push('/login')
  } catch (error) {
    ElMessage.error(error.response?.data?.message || '注册失败')
  } finally {
    loading.value = false
  }
}
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

/* 背景轮播 */
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
  transition: opacity 1.5s ease;
  will-change: opacity;
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

/* 毛玻璃半透明卡片 */
.login-card {
  background: rgba(255, 255, 255, 0.25);
  backdrop-filter: blur(12px);
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
  color: #fff;
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

/* 输入框半透明 */
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

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}
</style>