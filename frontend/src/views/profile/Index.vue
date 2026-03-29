<template>
  <div class="profile-page">
    <Navbar />

    <div class="container">
      <h1 class="page-title">个人中心</h1>

      <el-card class="profile-card">
        <div class="profile-header">
          <div class="avatar-container" @click="openEditProfile">
            <img :src="userInfo.avatar || defaultAvatar" class="avatar-img" alt="avatar" />
            <div class="avatar-overlay">编辑资料</div>
          </div>

          <div class="user-info">
            <h2>{{ userInfo.real_name || userInfo.username }}</h2>
            <p>{{ userInfo.mobile }}</p>
            <div class="user-tags">
              <el-tag size="small">{{ roleText }}</el-tag>
              <el-tag size="small" :type="faceRegistered ? 'success' : 'warning'">
                {{ faceRegistered ? '已录入人脸' : '未录入人脸' }}
              </el-tag>
            </div>
          </div>

          <div class="header-actions">
            <el-button type="primary" plain @click="openEditProfile">编辑资料</el-button>
            <el-button type="success" plain @click="openFaceDialog">
              {{ faceRegistered ? '重新录入人脸' : '录入人脸' }}
            </el-button>
          </div>
        </div>

        <div class="profile-stats">
          <div class="stat-item">
            <div class="stat-label">余额</div>
            <div class="stat-value">￥{{ formatAmount(userInfo.balance) }}</div>
          </div>
          <div class="stat-item">
            <div class="stat-label">绿色积分</div>
            <div class="stat-value green">{{ userInfo.green_points || 0 }}</div>
          </div>
          <div class="stat-item">
            <div class="stat-label">状态</div>
            <el-tag :type="userInfo.status === 1 ? 'success' : 'danger'">
              {{ userInfo.status === 1 ? '正常' : '冻结' }}
            </el-tag>
          </div>
        </div>
      </el-card>

      <div class="menu-grid">
        <el-card
          v-for="item in menuItems"
          :key="item.title"
          shadow="hover"
          class="menu-item"
          @click="item.action"
        >
          <div class="menu-icon-wrap">
            <el-icon class="menu-icon"><component :is="item.icon" /></el-icon>
          </div>
          <div class="menu-title">{{ item.title }}</div>
          <span>{{ item.desc }}</span>
        </el-card>
      </div>
    </div>

    <el-dialog v-model="showEditDialog" title="个人设置" width="520px">
      <el-tabs v-model="activeTab">
        <el-tab-pane label="基本信息" name="info">
          <el-form :model="editForm" label-width="80px">
            <el-form-item label="头像">
              <div class="avatar-uploader" @click="triggerFileUpload">
                <img v-if="editForm.avatar" :src="editForm.avatar" class="upload-avatar" alt="upload-avatar" />
                <span v-else>上传</span>
                <input ref="fileInput" type="file" accept="image/*" style="display: none" @change="handleUpload" />
              </div>
            </el-form-item>
            <el-form-item label="姓名">
              <el-input v-model="editForm.real_name" />
            </el-form-item>
            <el-form-item label="昵称">
              <el-input v-model="editForm.username" />
            </el-form-item>
            <el-form-item label="性别">
              <el-radio-group v-model="editForm.gender">
                <el-radio :label="1">男</el-radio>
                <el-radio :label="2">女</el-radio>
              </el-radio-group>
            </el-form-item>
            <el-form-item label="年龄">
              <el-input-number v-model="editForm.age" :min="1" :max="120" />
            </el-form-item>
            <el-form-item label="邮箱">
              <el-input v-model="editForm.email" />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" :loading="loading" @click="submitInfo">保存信息</el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>

        <el-tab-pane label="修改密码" name="pwd">
          <el-form :model="pwdForm" label-width="80px">
            <el-form-item label="旧密码">
              <el-input v-model="pwdForm.old_password" type="password" show-password />
            </el-form-item>
            <el-form-item label="新密码">
              <el-input v-model="pwdForm.new_password" type="password" show-password />
            </el-form-item>
            <el-form-item>
              <el-button type="danger" :loading="loading" @click="submitPwd">确认修改</el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>
      </el-tabs>
    </el-dialog>

    <el-dialog
      v-model="showFaceDialog"
      title="人脸录入"
      width="560px"
      :close-on-click-modal="false"
      @close="closeFaceDialog"
    >
      <div class="face-panel">
        <el-alert
          :type="faceRegistered ? 'warning' : 'info'"
          :closable="false"
          show-icon
          class="face-alert"
        >
          <template #title>
            {{ faceRegistered ? '当前账号已录入人脸，重新录入将覆盖原底图。' : '请先开启摄像头并抓拍清晰正脸照片。' }}
          </template>
        </el-alert>

        <video ref="faceVideoRef" class="face-video" autoplay playsinline muted />
        <img v-if="facePreview" :src="facePreview" class="face-preview" alt="face-preview" />

        <p v-if="faceError" class="face-error">{{ faceError }}</p>

        <div class="face-actions">
          <el-button @click="startFaceCamera">开启摄像头</el-button>
          <el-button type="primary" @click="captureFace" :disabled="!faceCameraReady">抓拍</el-button>
          <el-button type="success" :loading="enrolling" @click="saveFace">拍照并保存</el-button>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ShoppingBag, Trophy, Star, Wallet, Service, Setting, SwitchButton } from '@element-plus/icons-vue'
import Navbar from '@/components/layout/Navbar.vue'
import request from '@/utils/request'
import { useUserStore } from '@/stores/user'
import { updateUserInfo, changePassword, registerFace } from '@/api/user'

const router = useRouter()
const userStore = useUserStore()
const userInfo = computed(() => userStore.userInfo || {})
const defaultAvatar = 'https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png'
const faceRegistered = computed(() => Boolean(userInfo.value.face_registered))

const showEditDialog = ref(false)
const activeTab = ref('info')
const loading = ref(false)
const fileInput = ref(null)
const editForm = ref({})
const pwdForm = ref({ old_password: '', new_password: '' })

const showFaceDialog = ref(false)
const faceVideoRef = ref(null)
const faceCameraReady = ref(false)
const facePreview = ref('')
const faceBlob = ref(null)
const faceError = ref('')
const enrolling = ref(false)
let faceStream = null

const roleText = computed(() => {
  const map = { admin: '系统管理员', store: '商户', property: '物业人员', user: '居民' }
  return map[userInfo.value.role] || '居民'
})

const menuItems = computed(() => {
  const items = [
    {
      title: '订单',
      desc: '我的订单',
      icon: ShoppingBag,
      action: () => router.push('/order')
    },
    {
      title: '积分',
      desc: '绿色积分中心',
      icon: Trophy,
      action: () => router.push('/service/green-points')
    },
    {
      title: '收藏',
      desc: '我的收藏',
      icon: Star,
      action: () => router.push('/user/favorites')
    },
    {
      title: '账单',
      desc: '交易流水',
      icon: Wallet,
      action: () => router.push('/user/transactions')
    },
    {
      title: '服务',
      desc: '社区服务',
      icon: Service,
      action: () => router.push('/service')
    }
  ]

  if (['admin', 'store', 'property'].includes(userInfo.value.role)) {
    items.push({
      title: '后台',
      desc: '管理后台',
      icon: Setting,
      action: () => router.push('/admin')
    })
  }

  items.push({
    title: '退出',
    desc: '退出登录',
    icon: SwitchButton,
    action: handleLogout
  })

  return items
})

function formatAmount(value) {
  return Number(value || 0).toFixed(2)
}

function openEditProfile() {
  editForm.value = { ...userInfo.value }
  pwdForm.value = { old_password: '', new_password: '' }
  activeTab.value = 'info'
  showEditDialog.value = true
}

function triggerFileUpload() {
  fileInput.value?.click()
}

async function handleUpload(event) {
  const file = event.target.files?.[0]
  if (!file) return

  const formData = new FormData()
  formData.append('file', file)

  try {
    const res = await request({
      url: '/upload',
      method: 'post',
      data: formData,
      headers: { 'Content-Type': 'multipart/form-data' }
    })
    editForm.value.avatar = res.url
  } catch (error) {
    ElMessage.error('上传失败')
  }
}

async function submitInfo() {
  loading.value = true
  try {
    await updateUserInfo(editForm.value)
    await userStore.fetchUserInfo()
    ElMessage.success('保存成功')
    showEditDialog.value = false
  } catch (error) {
    ElMessage.error(error.response?.data?.msg || error.message || '保存失败')
  } finally {
    loading.value = false
  }
}

async function submitPwd() {
  if (!pwdForm.value.old_password || !pwdForm.value.new_password) {
    ElMessage.warning('请输入完整密码')
    return
  }

  loading.value = true
  try {
    await changePassword(pwdForm.value)
    ElMessage.success('密码已修改，请重新登录')
    await userStore.logout()
    router.push('/login')
  } catch (error) {
    ElMessage.error(error.response?.data?.msg || error.message || '修改失败')
  } finally {
    loading.value = false
  }
}

function handleLogout() {
  ElMessageBox.confirm('确认退出登录吗？', '提示', { type: 'warning' })
    .then(async () => {
      await userStore.logout()
      router.push('/home')
    })
    .catch(() => {})
}

function resetFaceState() {
  facePreview.value = ''
  faceBlob.value = null
  faceError.value = ''
  faceCameraReady.value = false
}

function openFaceDialog() {
  showFaceDialog.value = true
}

async function startFaceCamera() {
  faceError.value = ''
  if (faceStream) return

  try {
    faceStream = await navigator.mediaDevices.getUserMedia({
      video: { facingMode: 'user' },
      audio: false
    })

    if (!faceVideoRef.value) return
    faceVideoRef.value.srcObject = faceStream
    await faceVideoRef.value.play()
    faceCameraReady.value = true
  } catch (error) {
    faceError.value = '摄像头打开失败，请检查浏览器权限设置'
    stopFaceCamera()
  }
}

function stopFaceCamera() {
  if (faceStream) {
    faceStream.getTracks().forEach((track) => track.stop())
    faceStream = null
  }
  if (faceVideoRef.value) {
    faceVideoRef.value.srcObject = null
  }
  faceCameraReady.value = false
}

function captureFace() {
  if (!faceVideoRef.value || !faceCameraReady.value) {
    ElMessage.warning('请先开启摄像头')
    return
  }

  const video = faceVideoRef.value
  const canvas = document.createElement('canvas')
  canvas.width = video.videoWidth || 640
  canvas.height = video.videoHeight || 480

  const ctx = canvas.getContext('2d')
  if (!ctx) {
    ElMessage.error('抓拍失败，请重试')
    return
  }

  ctx.drawImage(video, 0, 0, canvas.width, canvas.height)
  facePreview.value = canvas.toDataURL('image/jpeg', 0.9)
  canvas.toBlob((blob) => {
    faceBlob.value = blob
  }, 'image/jpeg', 0.9)
}

async function saveFace() {
  if (!faceBlob.value) {
    ElMessage.warning('请先抓拍人脸照片')
    return
  }

  enrolling.value = true
  try {
    const formData = new FormData()
    formData.append('file', faceBlob.value, `face-register-${Date.now()}.jpg`)

    const uploadRes = await request({
      url: '/upload',
      method: 'post',
      data: formData,
      headers: { 'Content-Type': 'multipart/form-data' }
    })

    await registerFace({ face_image_url: uploadRes.url })
    await userStore.fetchUserInfo()
    ElMessage.success('人脸录入成功')
    closeFaceDialog()
  } catch (error) {
    ElMessage.error(error.response?.data?.msg || error.message || '人脸录入失败')
  } finally {
    enrolling.value = false
  }
}

function closeFaceDialog() {
  stopFaceCamera()
  resetFaceState()
  showFaceDialog.value = false
}

watch(showFaceDialog, (visible) => {
  if (visible) {
    resetFaceState()
    return
  }
  stopFaceCamera()
})

onBeforeUnmount(() => {
  stopFaceCamera()
})

onMounted(() => {
  userStore.fetchUserInfo()
})
</script>

<style scoped>
.profile-page {
  min-height: 100vh;
  padding-bottom: 40px;
}

.profile-card {
  margin-bottom: 20px;
}

.profile-header {
  display: flex;
  align-items: center;
  gap: 20px;
}

.avatar-container {
  position: relative;
  width: 88px;
  height: 88px;
  border-radius: 50%;
  overflow: hidden;
  cursor: pointer;
}

.avatar-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.avatar-overlay {
  position: absolute;
  inset: auto 0 0;
  text-align: center;
  padding: 4px 0;
  font-size: 12px;
  color: #fff;
  background: rgba(0, 0, 0, 0.55);
}

.user-info {
  flex: 1;
}

.user-tags {
  display: flex;
  gap: 8px;
}

.header-actions {
  margin-left: auto;
  display: flex;
  gap: 10px;
}

.profile-stats {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 16px;
  margin-top: 20px;
}

.stat-item {
  padding: 16px;
  border-radius: 12px;
  background: var(--bg-gray);
  text-align: center;
}

.stat-label {
  color: var(--text-secondary);
  margin-bottom: 8px;
}

.stat-value {
  font-size: 24px;
  font-weight: 700;
  color: var(--primary-color);
}

.stat-value.green {
  color: #2e7d32;
}

.menu-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(160px, 1fr));
  gap: 15px;
}

.menu-item {
  cursor: pointer;
  text-align: center;
  transition: all var(--transition-base);
}

.menu-item:hover {
  transform: translateY(-4px);
}

.menu-icon-wrap {
  width: 58px;
  height: 58px;
  border-radius: 16px;
  margin: 0 auto 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, rgba(0, 180, 137, 0.12), rgba(0, 180, 137, 0.2));
}

.menu-icon {
  font-size: 30px;
  color: var(--primary-color);
}

.menu-title {
  font-size: 18px;
  font-weight: 700;
  margin-bottom: 6px;
}

.avatar-uploader {
  width: 100px;
  height: 100px;
  border: 1px dashed #d9d9d9;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  overflow: hidden;
}

.upload-avatar {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.face-panel {
  display: grid;
  gap: 10px;
}

.face-alert {
  margin-bottom: 8px;
}

.face-video,
.face-preview {
  width: 100%;
  border-radius: 10px;
  background: #f3f4f6;
  aspect-ratio: 4 / 3;
  object-fit: cover;
}

.face-error {
  color: #d93025;
  margin: 0;
}

.face-actions {
  display: flex;
  gap: 10px;
}

@media (max-width: 768px) {
  .profile-header,
  .profile-stats {
    display: grid;
    grid-template-columns: 1fr;
  }

  .header-actions {
    margin-left: 0;
  }
}
</style>
