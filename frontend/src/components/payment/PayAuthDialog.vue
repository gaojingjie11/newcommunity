<template>
  <el-dialog
    :model-value="modelValue"
    :title="title"
    width="520px"
    :close-on-click-modal="false"
    @close="handleClose"
    @update:model-value="emit('update:modelValue', $event)"
  >
    <el-tabs v-model="activeTab" @tab-change="handleTabChange">
      <el-tab-pane label="密码支付" name="password">
        <el-form label-width="88px">
          <el-form-item label="登录密码">
            <el-input
              v-model="password"
              type="password"
              show-password
              placeholder="请输入登录密码"
              @keyup.enter="submit"
            />
          </el-form-item>
        </el-form>
      </el-tab-pane>

      <el-tab-pane label="刷脸支付" name="face">
        <div class="face-panel">
          <video ref="videoRef" class="face-video" autoplay playsinline muted />
          <img v-if="facePreview" :src="facePreview" class="face-preview" alt="face preview" />

          <p class="face-tip" v-if="!faceRegistered">
            当前账号未登记人脸，请先到个人中心录入人脸。
          </p>
          <p class="face-tip" v-else-if="cameraError">
            {{ cameraError }}
          </p>
          <p class="face-tip" v-else>
            请正对摄像头，确保光线充足后点击“抓拍”。
          </p>

          <div class="face-actions">
            <el-button @click="startCamera" :disabled="!faceRegistered">开启摄像头</el-button>
            <el-button type="primary" @click="captureFace" :disabled="!faceRegistered || !cameraReady">抓拍</el-button>
          </div>
        </div>
      </el-tab-pane>
    </el-tabs>

    <template #footer>
      <el-button @click="handleClose">取消</el-button>
      <el-button type="primary" :loading="processing || loading" @click="submit">
        {{ activeTab === 'face' ? '确认刷脸' : '确认支付' }}
      </el-button>
    </template>
  </el-dialog>
</template>

<script setup>
import { onBeforeUnmount, ref, watch } from 'vue'
import { ElMessage } from 'element-plus'
import request from '@/utils/request'

const props = defineProps({
  modelValue: { type: Boolean, default: false },
  title: { type: String, default: '支付验证' },
  faceRegistered: { type: Boolean, default: false },
  loading: { type: Boolean, default: false }
})

const emit = defineEmits(['update:modelValue', 'confirm'])

const activeTab = ref('password')
const password = ref('')
const videoRef = ref(null)
const cameraReady = ref(false)
const cameraError = ref('')
const facePreview = ref('')
const faceBlob = ref(null)
const processing = ref(false)
let mediaStream = null

watch(
  () => props.modelValue,
  (visible) => {
    if (visible) {
      resetDialog()
    } else {
      stopCamera()
    }
  }
)

onBeforeUnmount(() => {
  stopCamera()
})

function resetDialog() {
  activeTab.value = 'password'
  password.value = ''
  facePreview.value = ''
  faceBlob.value = null
  cameraError.value = ''
  cameraReady.value = false
}

function handleClose() {
  stopCamera()
  emit('update:modelValue', false)
}

function handleTabChange(name) {
  if (name === 'face') {
    if (!props.faceRegistered) {
      ElMessage.warning('当前账号未登记人脸，请先到个人中心录入人脸')
      activeTab.value = 'password'
      return
    }
    startCamera()
    return
  }
  stopCamera()
}

async function startCamera() {
  if (!props.faceRegistered || mediaStream) {
    return
  }
  cameraError.value = ''
  cameraReady.value = false

  try {
    mediaStream = await navigator.mediaDevices.getUserMedia({
      video: { facingMode: 'user' },
      audio: false
    })
    if (!videoRef.value) return
    videoRef.value.srcObject = mediaStream
    await videoRef.value.play()
    cameraReady.value = true
  } catch (error) {
    cameraError.value = '摄像头打开失败，请检查浏览器权限'
    stopCamera()
  }
}

function stopCamera() {
  if (mediaStream) {
    mediaStream.getTracks().forEach((track) => track.stop())
    mediaStream = null
  }
  if (videoRef.value) {
    videoRef.value.srcObject = null
  }
  cameraReady.value = false
}

function captureFace() {
  if (!videoRef.value || !cameraReady.value) {
    ElMessage.warning('请先开启摄像头')
    return
  }
  const video = videoRef.value
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

async function uploadCapturedFace() {
  if (!faceBlob.value) {
    throw new Error('请先抓拍人脸照片')
  }
  const formData = new FormData()
  formData.append('file', faceBlob.value, `face-pay-${Date.now()}.jpg`)

  const res = await request({
    url: '/upload',
    method: 'post',
    data: formData,
    headers: { 'Content-Type': 'multipart/form-data' }
  })
  return res?.url || ''
}

async function submit() {
  if (activeTab.value === 'password') {
    const pwd = password.value.trim()
    if (!pwd) {
      ElMessage.warning('请输入登录密码')
      return
    }
    emit('confirm', { pay_type: 'password', password: pwd })
    return
  }

  if (!props.faceRegistered) {
    ElMessage.warning('当前账号未登记人脸，请先到个人中心录入人脸')
    return
  }

  processing.value = true
  try {
    const faceImageURL = await uploadCapturedFace()
    emit('confirm', { pay_type: 'face', face_image_url: faceImageURL })
  } catch (error) {
    ElMessage.error(error.response?.data?.msg || error.message || '刷脸支付准备失败')
  } finally {
    processing.value = false
  }
}
</script>

<style scoped>
.face-panel {
  display: grid;
  gap: 10px;
}

.face-video,
.face-preview {
  width: 100%;
  border-radius: 10px;
  background: #f3f4f6;
  aspect-ratio: 4 / 3;
  object-fit: cover;
}

.face-tip {
  margin: 0;
  color: var(--text-secondary);
  font-size: 13px;
}

.face-actions {
  display: flex;
  gap: 10px;
}
</style>
