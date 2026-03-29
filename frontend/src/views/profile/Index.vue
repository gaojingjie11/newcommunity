<template>
  <div class="profile-page">
    <Navbar />

    <div class="container custom-container">
      <!-- 顶部个人信息主卡片 -->
      <div class="profile-card">
        <div class="profile-header">
          <div class="avatar-wrapper" @click="openEditProfile">
            <img
              :src="userInfo.avatar || defaultAvatar"
              class="avatar-img"
              alt="avatar"
            />
            <div class="avatar-overlay">
              <el-icon><Edit /></el-icon>
              <span>修改</span>
            </div>
          </div>

          <div class="user-info">
            <h2 class="user-name">
              {{ userInfo.real_name || userInfo.username }}
            </h2>
            <p class="user-mobile">{{ userInfo.mobile || "未绑定手机号" }}</p>
            <div class="user-tags">
              <span class="custom-tag tag-role">{{ roleText }}</span>
              <span
                class="custom-tag"
                :class="faceRegistered ? 'tag-success' : 'tag-warning'"
              >
                <el-icon class="tag-icon"
                  ><Camera v-if="faceRegistered" /><Warning v-else
                /></el-icon>
                {{ faceRegistered ? "已录入人脸" : "未录入人脸" }}
              </span>
            </div>
          </div>

          <div class="header-actions">
            <button class="action-btn btn-outline" @click="openEditProfile">
              编辑资料
            </button>
            <button class="action-btn btn-primary" @click="openFaceDialog">
              {{ faceRegistered ? "重新录入人脸" : "录入人脸" }}
            </button>
          </div>
        </div>

        <!-- 资产与状态统计 -->
        <div class="profile-stats">
          <div class="stat-item stat-balance">
            <div class="stat-label">账户余额 (元)</div>
            <div class="stat-value">￥{{ formatAmount(userInfo.balance) }}</div>
          </div>
          <div class="stat-item stat-points">
            <div class="stat-label">绿色积分</div>
            <div class="stat-value green">{{ userInfo.green_points || 0 }}</div>
          </div>
          <div class="stat-item stat-status">
            <div class="stat-label">当前状态</div>
            <div class="stat-value">
              <span
                class="status-indicator"
                :class="userInfo.status === 1 ? 'is-active' : 'is-frozen'"
              ></span>
              <span class="status-text">{{
                userInfo.status === 1 ? "账号正常" : "账号冻结"
              }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- 功能菜单矩阵 -->
      <div class="menu-section">
        <h3 class="section-title">我的服务</h3>
        <div class="menu-grid">
          <div
            v-for="item in menuItems"
            :key="item.title"
            class="menu-card"
            @click="item.action"
          >
            <div
              class="menu-icon-wrap"
              :class="{ 'is-danger': item.title === '退出' }"
            >
              <el-icon class="menu-icon"><component :is="item.icon" /></el-icon>
            </div>
            <div class="menu-text-wrap">
              <div
                class="menu-title"
                :class="{ 'text-danger': item.title === '退出' }"
              >
                {{ item.title }}
              </div>
              <span class="menu-desc">{{ item.desc }}</span>
            </div>
            <el-icon class="menu-arrow"><ArrowRight /></el-icon>
          </div>
        </div>
      </div>
    </div>

    <!-- 个人设置弹窗 -->
    <el-dialog
      v-model="showEditDialog"
      title="个人设置"
      width="560px"
      class="premium-dialog"
      :close-on-click-modal="false"
    >
      <el-tabs v-model="activeTab" class="custom-tabs">
        <el-tab-pane label="基本信息" name="info">
          <el-form :model="editForm" label-width="80px" class="premium-form">
            <el-form-item label="头像设置">
              <div class="avatar-uploader" @click="triggerFileUpload">
                <img
                  v-if="editForm.avatar"
                  :src="editForm.avatar"
                  class="upload-avatar"
                  alt="upload-avatar"
                />
                <div v-else class="upload-placeholder">
                  <el-icon><Plus /></el-icon>
                  <span>点击上传</span>
                </div>
                <input
                  ref="fileInput"
                  type="file"
                  accept="image/*"
                  style="display: none"
                  @change="handleUpload"
                />
              </div>
            </el-form-item>
            <el-form-item label="真实姓名">
              <el-input
                v-model="editForm.real_name"
                placeholder="请输入真实姓名"
                class="custom-input"
              />
            </el-form-item>
            <el-form-item label="显示昵称">
              <el-input
                v-model="editForm.username"
                placeholder="请输入昵称"
                class="custom-input"
              />
            </el-form-item>
            <el-form-item label="性别选择">
              <el-radio-group v-model="editForm.gender">
                <el-radio :label="1">先生</el-radio>
                <el-radio :label="2">女士</el-radio>
              </el-radio-group>
            </el-form-item>
            <el-form-item label="年龄">
              <el-input-number
                v-model="editForm.age"
                :min="1"
                :max="120"
                controls-position="right"
              />
            </el-form-item>
            <el-form-item label="电子邮箱">
              <el-input
                v-model="editForm.email"
                placeholder="请输入邮箱地址"
                class="custom-input"
              />
            </el-form-item>
            <el-form-item class="form-actions">
              <el-button
                type="primary"
                :loading="loading"
                @click="submitInfo"
                class="btn-submit"
                >保存修改</el-button
              >
            </el-form-item>
          </el-form>
        </el-tab-pane>

        <el-tab-pane label="修改密码" name="pwd">
          <el-form :model="pwdForm" label-width="80px" class="premium-form">
            <el-form-item label="原密码">
              <el-input
                v-model="pwdForm.old_password"
                type="password"
                show-password
                placeholder="请输入当前使用的密码"
                class="custom-input"
              />
            </el-form-item>
            <el-form-item label="新密码">
              <el-input
                v-model="pwdForm.new_password"
                type="password"
                show-password
                placeholder="请输入新密码"
                class="custom-input"
              />
            </el-form-item>
            <el-form-item class="form-actions">
              <el-button
                type="danger"
                :loading="loading"
                @click="submitPwd"
                class="btn-submit-danger"
                >确认修改并重新登录</el-button
              >
            </el-form-item>
          </el-form>
        </el-tab-pane>
      </el-tabs>
    </el-dialog>

    <!-- 人脸录入弹窗 -->
    <el-dialog
      v-model="showFaceDialog"
      title="人脸录入扫描"
      width="600px"
      :close-on-click-modal="false"
      @close="closeFaceDialog"
      class="premium-dialog"
    >
      <div class="face-panel">
        <div
          class="premium-alert"
          :class="faceRegistered ? 'is-warning' : 'is-info'"
        >
          <el-icon class="alert-icon"><InfoFilled /></el-icon>
          <div class="alert-content">
            {{
              faceRegistered
                ? "当前账号已录入人脸，重新录入将覆盖原底图。"
                : "请先开启摄像头并正对屏幕，确保光线充足。"
            }}
          </div>
        </div>

        <div class="scanner-container">
          <video
            ref="faceVideoRef"
            class="face-video"
            autoplay
            playsinline
            muted
            v-show="!facePreview"
          />
          <img
            v-if="facePreview"
            :src="facePreview"
            class="face-preview"
            alt="face-preview"
          />
          <div class="scanner-overlay" v-if="faceCameraReady && !facePreview">
            <div class="scan-line"></div>
          </div>
        </div>

        <p v-if="faceError" class="face-error">
          <el-icon><Warning /></el-icon> {{ faceError }}
        </p>

        <div class="face-actions-bar">
          <button
            class="action-btn btn-outline"
            @click="startFaceCamera"
            v-if="!faceCameraReady && !facePreview"
          >
            开启摄像头
          </button>
          <button
            class="action-btn btn-primary"
            @click="captureFace"
            v-if="faceCameraReady && !facePreview"
          >
            点击抓拍
          </button>
          <button
            class="action-btn btn-outline"
            @click="
              resetFaceState;
              startFaceCamera();
            "
            v-if="facePreview"
          >
            重新拍摄
          </button>
          <button
            class="action-btn btn-success"
            :class="{ 'is-loading': enrolling }"
            @click="saveFace"
            v-if="facePreview"
          >
            {{ enrolling ? "正在上传分析..." : "确认并保存" }}
          </button>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { computed, onBeforeUnmount, onMounted, ref, watch } from "vue";
import { useRouter } from "vue-router";
import { ElMessage, ElMessageBox } from "element-plus";
// 补全用到的图标
import {
  ShoppingBag,
  Trophy,
  Star,
  Wallet,
  Service,
  Setting,
  SwitchButton,
  Edit,
  Camera,
  Warning,
  ArrowRight,
  Plus,
  InfoFilled,
} from "@element-plus/icons-vue";
import Navbar from "@/components/layout/Navbar.vue";
import request from "@/utils/request";
import { useUserStore } from "@/stores/user";
import { updateUserInfo, changePassword, registerFace } from "@/api/user";

const router = useRouter();
const userStore = useUserStore();
const userInfo = computed(() => userStore.userInfo || {});
const defaultAvatar =
  "https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png";
const faceRegistered = computed(() => Boolean(userInfo.value.face_registered));

const showEditDialog = ref(false);
const activeTab = ref("info");
const loading = ref(false);
const fileInput = ref(null);
const editForm = ref({});
const pwdForm = ref({ old_password: "", new_password: "" });

const showFaceDialog = ref(false);
const faceVideoRef = ref(null);
const faceCameraReady = ref(false);
const facePreview = ref("");
const faceBlob = ref(null);
const faceError = ref("");
const enrolling = ref(false);
let faceStream = null;

const roleText = computed(() => {
  const map = {
    admin: "系统管理员",
    store: "商户",
    property: "物业人员",
    user: "居民",
  };
  return map[userInfo.value.role] || "居民";
});

const menuItems = computed(() => {
  const items = [
    {
      title: "订单",
      desc: "我的订单记录",
      icon: ShoppingBag,
      action: () => router.push("/order"),
    },
    {
      title: "积分",
      desc: "绿色积分中心",
      icon: Trophy,
      action: () => router.push("/service/green-points"),
    },
    {
      title: "收藏",
      desc: "我的商品收藏",
      icon: Star,
      action: () => router.push("/user/favorites"),
    },
    {
      title: "账单",
      desc: "钱包交易流水",
      icon: Wallet,
      action: () => router.push("/user/transactions"),
    },
    {
      title: "服务",
      desc: "社区便民服务",
      icon: Service,
      action: () => router.push("/service"),
    },
  ];

  if (["admin", "store", "property"].includes(userInfo.value.role)) {
    items.push({
      title: "后台",
      desc: "前往管理后台",
      icon: Setting,
      action: () => router.push("/admin"),
    });
  }

  items.push({
    title: "退出",
    desc: "注销当前账号",
    icon: SwitchButton,
    action: handleLogout,
  });

  return items;
});

function formatAmount(value) {
  return Number(value || 0).toFixed(2);
}

function openEditProfile() {
  editForm.value = { ...userInfo.value };
  pwdForm.value = { old_password: "", new_password: "" };
  activeTab.value = "info";
  showEditDialog.value = true;
}

function triggerFileUpload() {
  fileInput.value?.click();
}

async function handleUpload(event) {
  const file = event.target.files?.[0];
  if (!file) return;

  const formData = new FormData();
  formData.append("file", file);

  try {
    const res = await request({
      url: "/upload",
      method: "post",
      data: formData,
      headers: { "Content-Type": "multipart/form-data" },
    });
    editForm.value.avatar = res.url;
  } catch (error) {
    ElMessage.error("上传失败");
  }
}

async function submitInfo() {
  loading.value = true;
  try {
    await updateUserInfo(editForm.value);
    await userStore.fetchUserInfo();
    ElMessage.success("保存成功");
    showEditDialog.value = false;
  } catch (error) {
    ElMessage.error(error.response?.data?.msg || error.message || "保存失败");
  } finally {
    loading.value = false;
  }
}

async function submitPwd() {
  if (!pwdForm.value.old_password || !pwdForm.value.new_password) {
    ElMessage.warning("请输入完整密码");
    return;
  }

  loading.value = true;
  try {
    await changePassword(pwdForm.value);
    ElMessage.success("密码已修改，请重新登录");
    await userStore.logout();
    router.push("/login");
  } catch (error) {
    ElMessage.error(error.response?.data?.msg || error.message || "修改失败");
  } finally {
    loading.value = false;
  }
}

function handleLogout() {
  ElMessageBox.confirm("确认退出当前登录状态吗？", "注销提示", {
    type: "warning",
    confirmButtonText: "确认退出",
    cancelButtonText: "暂不退出",
    confirmButtonClass: "el-button--danger",
  })
    .then(async () => {
      await userStore.logout();
      router.push("/home");
    })
    .catch(() => {});
}

function resetFaceState() {
  facePreview.value = "";
  faceBlob.value = null;
  faceError.value = "";
  faceCameraReady.value = false;
}

function openFaceDialog() {
  showFaceDialog.value = true;
}

async function startFaceCamera() {
  faceError.value = "";
  facePreview.value = "";
  if (faceStream) return;

  try {
    faceStream = await navigator.mediaDevices.getUserMedia({
      video: { facingMode: "user" },
      audio: false,
    });

    if (!faceVideoRef.value) return;
    faceVideoRef.value.srcObject = faceStream;
    await faceVideoRef.value.play();
    faceCameraReady.value = true;
  } catch (error) {
    faceError.value = "摄像头打开失败，请检查浏览器权限设置或设备连接";
    stopFaceCamera();
  }
}

function stopFaceCamera() {
  if (faceStream) {
    faceStream.getTracks().forEach((track) => track.stop());
    faceStream = null;
  }
  if (faceVideoRef.value) {
    faceVideoRef.value.srcObject = null;
  }
  faceCameraReady.value = false;
}

function captureFace() {
  if (!faceVideoRef.value || !faceCameraReady.value) {
    ElMessage.warning("请先开启摄像头");
    return;
  }

  const video = faceVideoRef.value;
  const canvas = document.createElement("canvas");
  canvas.width = video.videoWidth || 640;
  canvas.height = video.videoHeight || 480;

  const ctx = canvas.getContext("2d");
  if (!ctx) {
    ElMessage.error("抓拍失败，请重试");
    return;
  }

  ctx.drawImage(video, 0, 0, canvas.width, canvas.height);
  facePreview.value = canvas.toDataURL("image/jpeg", 0.9);
  canvas.toBlob(
    (blob) => {
      faceBlob.value = blob;
    },
    "image/jpeg",
    0.9,
  );
}

async function saveFace() {
  if (!faceBlob.value) {
    ElMessage.warning("请先抓拍人脸照片");
    return;
  }

  enrolling.value = true;
  try {
    const formData = new FormData();
    formData.append("file", faceBlob.value, `face-register-${Date.now()}.jpg`);

    const uploadRes = await request({
      url: "/upload",
      method: "post",
      data: formData,
      headers: { "Content-Type": "multipart/form-data" },
    });

    await registerFace({ face_image_url: uploadRes.url });
    await userStore.fetchUserInfo();
    ElMessage.success("人脸特征录入成功");
    closeFaceDialog();
  } catch (error) {
    ElMessage.error(
      error.response?.data?.msg || error.message || "人脸录入失败",
    );
  } finally {
    enrolling.value = false;
  }
}

function closeFaceDialog() {
  stopFaceCamera();
  resetFaceState();
  showFaceDialog.value = false;
}

watch(showFaceDialog, (visible) => {
  if (visible) {
    resetFaceState();
    return;
  }
  stopFaceCamera();
});

onBeforeUnmount(() => {
  stopFaceCamera();
});

onMounted(() => {
  userStore.fetchUserInfo();
});
</script>

<style scoped>
.profile-page {
  min-height: 100vh;
  background-color: #f8f9fa;
  padding-bottom: 80px;
}

.custom-container {
  max-width: 1000px;
  margin: 0 auto;
}

.page-header {
  padding: 32px 0 24px;
}

.highlight-title {
  display: inline-block;
  position: relative;
  font-size: 32px;
  color: #2c3e50;
  font-weight: 700;
  margin: 0;
  z-index: 1;
}

.highlight-title::after {
  content: "";
  position: absolute;
  bottom: 4px;
  left: -5%;
  width: 110%;
  height: 14px;
  background-color: #2d597b;
  opacity: 0.15;
  border-radius: 6px;
  z-index: -1;
  transition: all 0.3s ease;
}

/* ================= 1. 顶部个人信息主卡片 ================= */
.profile-card {
  background: #ffffff;
  border-radius: 16px;
  padding: 40px;
  box-shadow: 0 4px 24px rgba(0, 0, 0, 0.03);
  margin-bottom: 40px;
  margin-top: 40px;
}

.profile-header {
  display: flex;
  align-items: center;
  gap: 32px;
  margin-bottom: 32px;
  padding-bottom: 32px;
  border-bottom: 1px dashed #ebeef5;
}

.avatar-wrapper {
  position: relative;
  width: 100px;
  height: 100px;
  border-radius: 50%;
  overflow: hidden;
  cursor: pointer;
  border: 4px solid #ffffff;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);
}

.avatar-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.avatar-overlay {
  position: absolute;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: #ffffff;
  font-size: 13px;
  opacity: 0;
  transition: opacity 0.3s;
}

.avatar-wrapper:hover .avatar-overlay {
  opacity: 1;
}

.user-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.user-name {
  margin: 0;
  font-size: 26px;
  font-weight: 700;
  color: #2c3e50;
}

.user-mobile {
  margin: 0;
  font-size: 15px;
  color: #909399;
}

.user-tags {
  display: flex;
  gap: 12px;
  margin-top: 4px;
}

.custom-tag {
  display: inline-flex;
  align-items: center;
  padding: 4px 12px;
  border-radius: 6px;
  font-size: 12px;
  font-weight: 600;
}

.tag-icon {
  margin-right: 4px;
  font-size: 14px;
}
.tag-role {
  background: #f0f7ff;
  color: #2d597b;
  border: 1px solid #cce3f6;
}
.tag-success {
  background: #f0fdf4;
  color: #166534;
  border: 1px solid #bbf7d0;
}
.tag-warning {
  background: #fff7ed;
  color: #9a3412;
  border: 1px solid #fed7aa;
}

.header-actions {
  display: flex;
  gap: 16px;
}

/* 按钮通用规范 */
.action-btn {
  padding: 10px 24px;
  border-radius: 20px;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s;
  border: 1px solid transparent;
}

.btn-primary {
  background: #2d597b;
  color: #ffffff;
  box-shadow: 0 4px 12px rgba(45, 89, 123, 0.2);
}
.btn-primary:hover {
  background: #1f435d;
  transform: translateY(-2px);
  box-shadow: 0 6px 16px rgba(45, 89, 123, 0.3);
}

.btn-success {
  background: #2d597b;
  color: #ffffff;
  box-shadow: 0 4px 12px rgba(45, 89, 123, 0.2);
}
.btn-success:hover {
  background: #1f435d;
  transform: translateY(-2px);
  box-shadow: 0 6px 16px rgba(45, 89, 123, 0.3);
}

.btn-outline {
  background: #ffffff;
  color: #2d597b;
  border-color: #2d597b;
}
.btn-outline:hover {
  background: #f0f7ff;
}

/* 资产与状态统计 */
.profile-stats {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 24px;
}

.stat-item {
  padding: 24px;
  border-radius: 12px;
  display: flex;
  flex-direction: column;
  gap: 8px;
  border: 1px solid transparent;
}

.stat-balance {
  background: #fdf6f6;
  border-color: #fbc4c4;
}
.stat-balance .stat-value {
  color: #e4393c;
}

.stat-points {
  background: #f0fdf4;
  border-color: #bbf7d0;
}
.stat-points .stat-value {
  color: #00b894;
}

.stat-status {
  background: #f8f9fa;
  border-color: #ebeef5;
}
.stat-status .stat-value {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 18px;
  color: #303133;
}

.stat-label {
  font-size: 14px;
  color: #606266;
}

.stat-value {
  font-size: 28px;
  font-weight: 800;
  font-family:
    -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
}

.status-indicator {
  width: 10px;
  height: 10px;
  border-radius: 50%;
}
.status-indicator.is-active {
  background: #52c41a;
  box-shadow: 0 0 8px rgba(82, 196, 26, 0.5);
}
.status-indicator.is-frozen {
  background: #f56c6c;
}

/* ================= 2. 功能菜单矩阵 ================= */
.menu-section {
  margin-top: 20px;
}

.section-title {
  font-size: 20px;
  font-weight: 700;
  color: #2c3e50;
  margin-bottom: 24px;
  padding-left: 12px;
  border-left: 4px solid #2d597b;
}

.menu-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 20px;
}

.menu-card {
  display: flex;
  align-items: center;
  background: #ffffff;
  padding: 24px;
  border-radius: 16px;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.02);
  border: 1px solid transparent;
}

.menu-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 24px rgba(45, 89, 123, 0.08);
  border-color: rgba(45, 89, 123, 0.1);
}

.menu-icon-wrap {
  width: 56px;
  height: 56px;
  border-radius: 12px;
  background: #f0f7ff;
  color: #2d597b;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 28px;
  margin-right: 20px;
  transition: all 0.3s;
}

.menu-icon-wrap.is-danger {
  background: #fdf6f6;
  color: #e4393c;
}
.menu-card:hover .menu-icon-wrap {
  background: #2d597b;
  color: #ffffff;
}
.menu-card:hover .menu-icon-wrap.is-danger {
  background: #e4393c;
  color: #ffffff;
}

.menu-text-wrap {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.menu-title {
  font-size: 16px;
  font-weight: 600;
  color: #303133;
}
.menu-title.text-danger {
  color: #e4393c;
}

.menu-desc {
  font-size: 13px;
  color: #909399;
}

.menu-arrow {
  color: #dcdfe6;
  font-size: 20px;
  transition: transform 0.3s;
}
.menu-card:hover .menu-arrow {
  transform: translateX(4px);
  color: #2d597b;
}

/* ================= 3. 弹窗深度美化 ================= */
:deep(.premium-dialog .el-dialog__header) {
  border-bottom: 1px solid #f0f2f5;
  padding-bottom: 16px;
  margin-right: 0;
}
:deep(.premium-dialog .el-dialog__title) {
  font-weight: 700;
  color: #2c3e50;
}

.premium-form {
  padding: 20px 0 0;
}

/* 头像上传区 */
.avatar-uploader {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  border: 2px dashed #dcdfe6;
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
  cursor: pointer;
  transition: all 0.3s;
  background: #fafbfc;
}
.avatar-uploader:hover {
  border-color: #2d597b;
  background: #f0f7ff;
  color: #2d597b;
}
.upload-avatar {
  width: 100%;
  height: 100%;
  object-fit: cover;
}
.upload-placeholder {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
  color: #a4b0be;
  font-size: 12px;
}

/* 深度定制表单输入框 */
:deep(.custom-input .el-input__wrapper) {
  box-shadow: 0 0 0 1px #dcdfe6 inset;
  border-radius: 8px;
  padding: 4px 12px;
  background: #fbfcfd;
  transition: all 0.3s;
}
:deep(.custom-input .el-input__wrapper.is-focus) {
  box-shadow: 0 0 0 2px rgba(45, 89, 123, 0.2) inset !important;
  background: #ffffff;
}

.form-actions {
  margin-top: 32px;
  text-align: right;
}
.btn-submit {
  padding: 10px 32px;
  border-radius: 20px;
  font-weight: 600;
  background: #2d597b;
  border: none;
}
.btn-submit-danger {
  padding: 10px 32px;
  border-radius: 20px;
  font-weight: 600;
}

/* 人脸录入专属样式 */
.face-panel {
  display: flex;
  flex-direction: column;
  gap: 20px;
  padding-top: 10px;
}

.premium-alert {
  display: flex;
  align-items: flex-start;
  border-radius: 8px;
  padding: 12px 16px;
  font-size: 14px;
  line-height: 1.5;
}
.premium-alert.is-warning {
  background: #fff7ed;
  border: 1px solid #fed7aa;
  color: #9a3412;
}
.premium-alert.is-info {
  background: #f0f7ff;
  border: 1px solid #cce3f6;
  color: #2d597b;
}
.alert-icon {
  font-size: 18px;
  margin-right: 8px;
  margin-top: 2px;
}

/* 扫描仪效果容器 */
.scanner-container {
  position: relative;
  width: 100%;
  aspect-ratio: 4 / 3;
  border-radius: 16px;
  overflow: hidden;
  background: #1e293b;
  box-shadow: inset 0 0 20px rgba(0, 0, 0, 0.5);
  border: 2px solid #e2e8f0;
}

.face-video,
.face-preview {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transform: scaleX(-1); /* 镜像翻转，符合照镜子习惯 */
}
.face-preview {
  transform: none;
} /* 拍出来的照片不需要镜像 */

.scanner-overlay {
  position: absolute;
  inset: 0;
  background: rgba(45, 89, 123, 0.1);
  border: 2px solid rgba(45, 89, 123, 0.5);
  border-radius: 16px;
  pointer-events: none;
}
.scan-line {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 4px;
  background: #00b894;
  box-shadow: 0 0 10px #00b894;
  opacity: 0.8;
  animation: scan 2.5s infinite linear;
}
@keyframes scan {
  0% {
    top: 0%;
  }
  50% {
    top: 100%;
  }
  100% {
    top: 0%;
  }
}

.face-error {
  color: #e4393c;
  font-size: 14px;
  display: flex;
  align-items: center;
  gap: 6px;
  margin: 0;
}

.face-actions-bar {
  display: flex;
  justify-content: center;
  gap: 16px;
  margin-top: 10px;
}

/* 响应式 */
@media (max-width: 768px) {
  .profile-header {
    flex-direction: column;
    text-align: center;
  }
  .header-actions {
    margin-left: 0;
    justify-content: center;
  }
  .user-tags {
    justify-content: center;
  }
  .profile-stats {
    grid-template-columns: 1fr;
  }
  .menu-grid {
    grid-template-columns: 1fr;
  }
}
</style>
