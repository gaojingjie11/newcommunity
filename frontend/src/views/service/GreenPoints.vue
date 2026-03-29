<template>
  <div class="green-points-page">
    <Navbar />

    <div class="container custom-container">
      <div class="page-header">
        <h1 class="page-title highlight-title">绿色积分中心</h1>
      </div>

      <div class="hero-banner">
        <div class="hero-left">
          <h2>AI 智能分类，环保赚积分</h2>
          <p class="hero-desc">
            上传垃圾分类图片，AI 自动识别并发放积分。积分按 <strong>{{ GREEN_POINTS_PER_YUAN }} 积分 = 1 元</strong> 比例，可在支付商城订单和物业费时优先抵扣。
          </p>
        </div>
        <div class="hero-stats">
          <div class="stat-box">
            <span class="stat-label">当前积分余额</span>
            <span class="stat-value text-green">{{ userStore.userInfo.green_points || 0 }}</span>
          </div>
          <div class="stat-box">
            <span class="stat-label">账户余额</span>
            <span class="stat-value">¥{{ formatAmount(userStore.userInfo.balance || 0) }}</span>
          </div>
        </div>
      </div>

      <div class="content-grid">
        <div class="card premium-card">
          <div class="card-header">
            <span class="header-indicator"></span> AI 垃圾分类识别
          </div>
          
          <div class="upload-wrapper">
            <el-upload
              drag
              :auto-upload="false"
              :show-file-list="false"
              accept="image/*"
              :on-change="handleFileChange"
              class="custom-upload"
            >
              <el-icon class="upload-icon"><UploadFilled /></el-icon>
              <div class="el-upload__text">拖拽图片到此处，或 <em>点击选择文件</em></div>
              <template #tip>
                <div class="upload-tip">建议上传清晰的垃圾分类照片，支持 jpg/png/webp</div>
              </template>
            </el-upload>
          </div>

          <div v-if="previewUrl" class="preview-box">
            <img :src="previewUrl" alt="preview" class="preview-image" />
          </div>

          <div class="actions">
            <button class="btn-action btn-success" :disabled="!selectedFile || uploading" @click="submitGarbageImage">
              {{ uploading ? 'AI识别中...' : '开始识别并领取积分' }}
            </button>
          </div>

          <el-result
            v-if="recognitionResult"
            icon="success"
            :title="`本次奖励 ${recognitionResult.points} 积分`"
            :sub-title="recognitionResult.reason"
            class="custom-result"
          >
            <template #extra>
              <div class="result-extra">当前积分余额：<strong>{{ recognitionResult.green_points }}</strong></div>
            </template>
          </el-result>
        </div>

        <div class="card premium-card">
          <div class="card-header">
            <span class="header-indicator"></span> 累计积分排行榜
          </div>
          
          <el-table :data="leaderboard" style="width: 100%" class="custom-table" :header-cell-style="{background:'#f8f9fa', color:'#606266'}">
            <el-table-column label="排名" width="80" align="center">
              <template #default="{ row }">
                <span class="rank-badge" :class="`rank-${row.rank}`">#{{ row.rank }}</span>
              </template>
            </el-table-column>
            <el-table-column label="环保卫士">
              <template #default="{ row }">
                <div class="user-cell">
                  <img :src="row.avatar || 'https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png'" class="table-avatar" />
                  <span class="user-name">{{ row.nickname || row.username }}</span>
                </div>
              </template>
            </el-table-column>
            <el-table-column prop="points" label="总积分" width="100" align="right">
              <template #default="{ row }">
                <strong class="points-text">{{ row.points }}</strong>
              </template>
            </el-table-column>
          </el-table>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { UploadFilled } from '@element-plus/icons-vue'
import Navbar from '@/components/layout/Navbar.vue'
import { uploadGarbageImage, getGreenPointsLeaderboard } from '@/api/greenPoints'
import { useUserStore } from '@/stores/user'
import { GREEN_POINTS_PER_YUAN } from '@/utils/payment'

const userStore = useUserStore()
const selectedFile = ref(null)
const previewUrl = ref('')
const uploading = ref(false)
const recognitionResult = ref(null)
const leaderboard = ref([])

function formatAmount(value) {
  return Number(value || 0).toFixed(2)
}

function handleFileChange(file) {
  selectedFile.value = file.raw
  recognitionResult.value = null
  if (file.raw) {
    previewUrl.value = URL.createObjectURL(file.raw)
  }
}

async function submitGarbageImage() {
  if (!selectedFile.value) {
    ElMessage.warning('请先选择图片')
    return
  }

  const formData = new FormData()
  formData.append('file', selectedFile.value)

  uploading.value = true
  try {
    const res = await uploadGarbageImage(formData)
    recognitionResult.value = res
    ElMessage.success(`识别成功，奖励 ${res.points} 积分`)
    await Promise.all([fetchLeaderboard(), userStore.fetchUserInfo()])
  } catch (error) {
    ElMessage.error(error.response?.data?.msg || error.message || '识别失败')
  } finally {
    uploading.value = false
  }
}

async function fetchLeaderboard() {
  const res = await getGreenPointsLeaderboard({ limit: 10 })
  leaderboard.value = res.list || []
}

onMounted(async () => {
  await Promise.all([fetchLeaderboard(), userStore.fetchUserInfo()])
})
</script>

<style scoped>
.green-points-page { min-height: 100vh; background-color: #f8f9fa; padding-bottom: 60px; }
.custom-container { max-width: 1100px; margin: 0 auto; }
.page-header { padding: 32px 0 24px; }
.highlight-title { display: inline-block; position: relative; font-size: 32px; color: #2c3e50; font-weight: 700; margin: 0; z-index: 1; }
.highlight-title::after { content: ''; position: absolute; bottom: 4px; left: -5%; width: 110%; height: 14px; background-color: #2d597b; opacity: 0.15; border-radius: 6px; z-index: -1; }

.hero-banner {
  background: linear-gradient(135deg, #2d597b 0%, #2d597b 100%);
  border-radius: 16px; padding: 40px; color: #ffffff; display: flex; justify-content: space-between; align-items: center; margin-bottom: 32px; box-shadow: 0 8px 24px rgba(0, 184, 148, 0.2);
}
.hero-left h2 { margin: 0 0 12px 0; font-size: 26px; font-weight: 700; }
.hero-desc { margin: 0; font-size: 15px; opacity: 0.9; line-height: 1.6; max-width: 600px; }

.hero-stats { display: flex; gap: 20px; }
.stat-box { background: rgba(255, 255, 255, 0.95); border-radius: 12px; padding: 20px 32px; min-width: 160px; box-shadow: 0 4px 12px rgba(0,0,0,0.05); }
.stat-label { display: block; font-size: 13px; color: #606266; margin-bottom: 8px; }
.stat-value { font-size: 28px; font-weight: 800; color: #303133; }
.text-green { color: #2d597b; }

.content-grid { display: grid; grid-template-columns: 1.2fr 1fr; gap: 24px; }
.premium-card { background: #ffffff; border-radius: 16px; padding: 32px; box-shadow: 0 2px 16px rgba(0,0,0,0.03); }
.card-header { font-size: 18px; font-weight: 700; color: #2c3e50; margin-bottom: 24px; display: flex; align-items: center; }
.header-indicator { display: inline-block; width: 4px; height: 18px; background: #2d597b; border-radius: 2px; margin-right: 12px; }

.upload-wrapper { margin-bottom: 24px; }
:deep(.custom-upload .el-upload-dragger) { background: #fbfcfd; border: 1px dashed #dcdfe6; border-radius: 12px; padding: 40px 0; transition: all 0.3s; }
:deep(.custom-upload .el-upload-dragger:hover) { border-color: #2d597b; background: #f0fdf4; }
.upload-icon { font-size: 48px; color: #a4b0be; margin-bottom: 16px; }
.upload-tip { font-size: 13px; color: #a4b0be; margin-top: 12px; }
:deep(.el-upload__text em) { color: #2d597b; font-style: normal; font-weight: bold; }

.preview-box { text-align: center; margin-bottom: 24px; }
.preview-image { max-width: 100%; max-height: 280px; border-radius: 8px; box-shadow: 0 4px 12px rgba(0,0,0,0.1); border: 4px solid #fff; }

.actions { text-align: center; margin-bottom: 24px; }
.btn-action { padding: 12px 40px; border-radius: 24px; font-size: 16px; font-weight: bold; border: none; cursor: pointer; transition: all 0.3s; }
.btn-success { background: #2d597b; color: #fff; box-shadow: 0 4px 12px rgba(0,184,148,0.2); }
.btn-success:hover:not(:disabled) { background: #2d597b; transform: translateY(-2px); box-shadow: 0 6px 16px rgba(0,184,148,0.3); }
.btn-action:disabled { background: #dcdfe6; box-shadow: none; cursor: not-allowed; }

.custom-result { padding: 20px; background: #f0fdf4; border-radius: 12px; }
.result-extra { color: #166534; font-size: 14px; background: #dcfce7; padding: 6px 16px; border-radius: 20px; display: inline-block; }

/* Table styles */
:deep(.custom-table) { border-radius: 8px; overflow: hidden; }
:deep(.el-table td.el-table__cell) { border-bottom: 1px dashed #ebeef5; padding: 16px 0; }

.rank-badge { display: inline-block; padding: 2px 10px; border-radius: 12px; font-size: 13px; font-weight: bold; background: #f4f4f5; color: #909399; }
.rank-1 { background: #fffbe6; color: #e6a23c; font-size: 15px; }
.rank-2 { background: #f0f9eb; color: #67c23a; font-size: 14px; }
.rank-3 { background: #fdf6f6; color: #f56c6c; font-size: 14px; }

.user-cell { display: flex; align-items: center; gap: 12px; }
.table-avatar { width: 36px; height: 36px; border-radius: 50%; object-fit: cover; }
.user-name { font-weight: 600; color: #303133; }
.points-text { color: #2d597b; font-size: 16px; font-weight: 800; }

@media (max-width: 900px) {
  .hero-banner, .content-grid { flex-direction: column; display: flex; }
  .hero-stats { width: 100%; flex-wrap: wrap; }
  .stat-box { flex: 1; }
}
</style>