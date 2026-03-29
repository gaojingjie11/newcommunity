<template>
  <div class="admin-child-page">
    <Navbar />
    <div class="container custom-container">
      <!-- 顶部返回导航 -->
      <div class="page-nav">
        <div class="back-btn" @click="$router.push('/admin')">
          <el-icon class="back-icon"><ArrowLeft /></el-icon>
          <span>返回管理后台</span>
        </div>
      </div>

      <!-- 高级标题与操作区 -->
      <div class="page-header">
        <h1 class="page-title highlight-title">商品管理</h1>
        <button class="action-btn btn-primary" @click="openModal()">
          <el-icon style="margin-right: 4px"><Plus /></el-icon> 发布新商品
        </button>
      </div>

      <!-- 沉浸式搜索过滤卡片 -->
      <div class="search-card">
        <el-form
          :inline="true"
          :model="filters"
          class="premium-search-form"
          @submit.prevent
        >
          <el-form-item label="商品名称">
            <el-input
              v-model="filters.name"
              placeholder="输入名称搜索"
              clearable
              @keyup.enter="handleSearch"
              class="custom-input search-input"
            >
              <template #prefix
                ><el-icon><Search /></el-icon
              ></template>
            </el-input>
          </el-form-item>

          <el-form-item label="商品分类">
            <el-select
              v-model="filters.category_id"
              placeholder="所有分类"
              clearable
              class="custom-select"
              style="width: 160px"
            >
              <el-option label="全部类目" :value="0" />
              <el-option
                v-for="c in categories"
                :key="c.id"
                :label="c.name"
                :value="c.id"
              />
            </el-select>
          </el-form-item>

          <el-form-item class="checkbox-item">
            <el-checkbox
              v-model="filters.is_promotion"
              label="仅看促销"
              class="custom-checkbox"
            />
          </el-form-item>

          <el-form-item class="form-actions">
            <button
              type="button"
              class="action-btn btn-primary btn-sm"
              @click="handleSearch"
            >
              搜索
            </button>
            <button
              type="button"
              class="action-btn btn-default btn-sm"
              @click="resetFilters"
            >
              重置
            </button>
          </el-form-item>
        </el-form>
      </div>

      <!-- 深度美化的表格容器 -->
      <div class="table-wrapper">
        <el-table
          :data="products"
          class="custom-table"
          style="width: 100%"
          v-loading="loadingData"
          :empty-text="'暂无符合条件的商品'"
        >
          <el-table-column prop="id" label="ID" width="80" align="center" />

          <el-table-column label="商品图片" width="100" align="center">
            <template #default="scope">
              <div class="product-thumb-wrapper">
                <img
                  :src="scope.row.image_url"
                  class="product-thumb"
                  alt="商品"
                  v-if="scope.row.image_url"
                />
                <span class="no-img" v-else>无图</span>
              </div>
            </template>
          </el-table-column>

          <el-table-column
            prop="name"
            label="商品名称"
            min-width="200"
            show-overflow-tooltip
          >
            <template #default="{ row }">
              <strong class="text-primary">{{ row.name }}</strong>
            </template>
          </el-table-column>

          <el-table-column
            prop="price"
            label="促销价 / 现价"
            width="140"
            align="right"
          >
            <template #default="scope">
              <span class="price-text"
                >¥{{ Number(scope.row.price).toFixed(2) }}</span
              >
            </template>
          </el-table-column>

          <el-table-column
            prop="original_price"
            label="市场原价"
            width="120"
            align="right"
          >
            <template #default="scope">
              <span v-if="scope.row.original_price" class="original-price-text">
                ¥{{ Number(scope.row.original_price).toFixed(2) }}
              </span>
              <span v-else class="text-secondary">-</span>
            </template>
          </el-table-column>

          <el-table-column
            prop="stock"
            label="库存"
            width="100"
            align="center"
          />

          <el-table-column
            prop="category_name"
            label="所属分类"
            width="120"
            align="center"
          >
            <template #default="{ row }">
              <span class="category-tag">{{
                row.category_name || "未分类"
              }}</span>
            </template>
          </el-table-column>

          <el-table-column
            prop="status"
            label="上架状态"
            width="100"
            align="center"
          >
            <template #default="scope">
              <span
                class="status-badge"
                :class="scope.row.status === 1 ? 'is-active' : 'is-inactive'"
              >
                {{ scope.row.status === 1 ? "售卖中" : "已下架" }}
              </span>
            </template>
          </el-table-column>

          <el-table-column
            label="操作"
            width="160"
            fixed="right"
            align="center"
          >
            <template #default="scope">
              <div class="row-actions">
                <button
                  class="action-btn btn-sm btn-outline"
                  @click="openModal(scope.row)"
                >
                  编辑
                </button>
                <el-popconfirm
                  title="确定将该商品从库中彻底删除吗？"
                  confirm-button-text="确定删除"
                  cancel-button-text="取消"
                  confirm-button-type="danger"
                  @confirm="handleDelete(scope.row.id)"
                >
                  <template #reference>
                    <button class="action-btn btn-sm btn-danger-ghost">
                      删除
                    </button>
                  </template>
                </el-popconfirm>
              </div>
            </template>
          </el-table-column>
        </el-table>

        <!-- 分页器 -->
        <div class="pagination-container" v-if="total > 0">
          <el-pagination
            v-model:current-page="page"
            v-model:page-size="size"
            :page-sizes="[10, 20, 50, 100]"
            layout="total, sizes, prev, pager, next, jumper"
            :total="total"
            @size-change="handleSearch"
            @current-change="fetchProducts"
            class="custom-pagination"
          />
        </div>
      </div>

      <!-- 高级双栏编辑弹窗 -->
      <el-dialog
        v-model="showModal"
        :title="isEdit ? '编辑商品信息' : '发布全新商品'"
        width="880px"
        class="premium-dialog"
        :close-on-click-modal="false"
      >
        <el-form
          ref="formRef"
          :model="form"
          label-position="top"
          class="premium-form"
        >
          <div class="form-grid">
            <!-- 左侧：基础信息 -->
            <div class="form-col-left">
              <el-form-item label="商品名称" required>
                <el-input
                  v-model="form.name"
                  placeholder="请输入精准的商品名称"
                  class="custom-input"
                />
              </el-form-item>

              <el-row :gutter="20">
                <el-col :span="12">
                  <el-form-item label="商品分类">
                    <el-select
                      v-model="form.category_id"
                      placeholder="请选择商品分类"
                      class="custom-select"
                      style="width: 100%"
                    >
                      <el-option label="未分类" :value="0" />
                      <el-option
                        v-for="c in categories"
                        :key="c.id"
                        :label="c.name"
                        :value="c.id"
                      />
                    </el-select>
                  </el-form-item>
                </el-col>
                <el-col :span="12">
                  <el-form-item label="库存数量" required>
                    <el-input-number
                      v-model="form.stock"
                      :min="0"
                      :step="1"
                      class="custom-input-number"
                      style="width: 100%"
                    />
                  </el-form-item>
                </el-col>
              </el-row>

              <el-row :gutter="20">
                <el-col :span="12">
                  <el-form-item label="市场原价 (¥)" required>
                    <el-input-number
                      v-model="form.original_price"
                      :min="0"
                      :precision="2"
                      class="custom-input-number"
                      placeholder="原价"
                      style="width: 100%"
                    />
                  </el-form-item>
                </el-col>
                <el-col :span="12">
                  <el-form-item label="实际售卖价 (¥)" required>
                    <el-input-number
                      v-model="form.price"
                      :min="0"
                      :precision="2"
                      class="custom-input-number highlight-number"
                      placeholder="现价"
                      style="width: 100%"
                    />
                  </el-form-item>
                </el-col>
              </el-row>

              <el-form-item class="mb-0 mt-2">
                <el-checkbox
                  v-model="form.status"
                  :true-label="1"
                  :false-label="0"
                  class="custom-checkbox bold-check"
                >
                  立即上架并在商城中展示
                </el-checkbox>
              </el-form-item>
            </div>

            <!-- 右侧：媒体与描述 -->
            <div class="form-col-right">
              <el-form-item label="商品主图">
                <div
                  class="custom-upload-box compact-upload"
                  @click="triggerFileUpload"
                >
                  <input
                    ref="fileInput"
                    type="file"
                    @change="handleUpload"
                    accept="image/*"
                    style="display: none"
                  />

                  <div v-if="uploading" class="upload-status">
                    <el-icon class="is-loading"><Loading /></el-icon> 上传中...
                  </div>
                  <div v-else-if="form.image_url" class="upload-preview">
                    <img :src="form.image_url" alt="preview" />
                    <div class="preview-overlay"><span>更换图片</span></div>
                  </div>
                  <div v-else class="upload-placeholder">
                    <el-icon class="upload-icon"><Upload /></el-icon>
                    <span>点击上传主图</span>
                  </div>
                </div>
              </el-form-item>

              <el-form-item label="商品图文描述" class="mb-0">
                <el-input
                  v-model="form.description"
                  type="textarea"
                  :rows="4"
                  class="custom-textarea"
                  placeholder="介绍一下这个商品的卖点和特色..."
                />
              </el-form-item>
            </div>
          </div>
        </el-form>

        <template #footer>
          <div class="modal-actions">
            <button
              type="button"
              class="action-btn btn-default"
              @click="showModal = false"
            >
              取消放弃
            </button>
            <button
              type="button"
              class="action-btn btn-primary"
              @click="handleSubmit"
              :disabled="uploading"
            >
              {{ uploading ? "处理中..." : "确认保存" }}
            </button>
          </div>
        </template>
      </el-dialog>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, reactive } from "vue";
import Navbar from "@/components/layout/Navbar.vue";
import { getProductList, getCategories } from "@/api/product";
import { createProduct, updateProduct, deleteProduct } from "@/api/admin";
import request from "@/utils/request";
import { ElMessage } from "element-plus";
import {
  ArrowLeft,
  Plus,
  Search,
  Upload,
  Loading,
} from "@element-plus/icons-vue";

const products = ref([]);
const categories = ref([]);
const showModal = ref(false);
const isEdit = ref(false);
const uploading = ref(false);
const loadingData = ref(false);
const total = ref(0);
const page = ref(1);
const size = ref(10);
const fileInput = ref(null);

const form = ref({
  id: 0,
  name: "",
  price: 0,
  original_price: 0,
  stock: 0,
  category_id: 0,
  description: "",
  image_url: "",
  status: 1,
});

const filters = reactive({
  name: "",
  category_id: 0,
  is_promotion: false,
});

const triggerFileUpload = () => {
  if (fileInput.value) {
    fileInput.value.click();
  }
};

const handleUpload = async (event) => {
  const file = event.target.files[0];
  if (!file) return;

  const formData = new FormData();
  formData.append("file", file);

  uploading.value = true;
  try {
    const res = await request({
      url: "/upload",
      method: "post",
      data: formData,
      headers: { "Content-Type": "multipart/form-data" },
    });
    form.value.image_url = res.url;
    ElMessage.success("图片上传成功");
  } catch (e) {
    const msg = e?.response?.data?.msg || e?.message || "上传失败";
    ElMessage.error(msg);
  } finally {
    uploading.value = false;
    event.target.value = ""; // Allow re-uploading same file
  }
};

const handleSearch = () => {
  page.value = 1;
  fetchProducts();
};

const fetchProducts = async () => {
  loadingData.value = true;
  try {
    const params = {
      name: filters.name,
      category_id: filters.category_id || undefined,
      is_promotion: filters.is_promotion || undefined,
      page: page.value,
      size: size.value,
    };
    const res = await getProductList(params);
    products.value = res.list || [];
    total.value = res.total || 0;
  } catch (error) {
    console.error(error);
  } finally {
    loadingData.value = false;
  }
};

const resetFilters = () => {
  filters.name = "";
  filters.category_id = 0;
  filters.is_promotion = false;
  handleSearch();
};

const fetchCategories = async () => {
  try {
    const res = await getCategories();
    categories.value = res || [];
  } catch (e) {
    console.error(e);
  }
};

const openModal = (product = null) => {
  isEdit.value = !!product;

  if (product) {
    form.value = {
      id: product.id,
      name: product.name,
      price: product.price,
      original_price: product.original_price || product.price,
      stock: product.stock,
      category_id: product.category_id || 0,
      description: product.description,
      image_url: product.image_url,
      status: product.status,
    };
  } else {
    form.value = {
      id: 0,
      name: "",
      price: 0,
      original_price: 0,
      stock: 999,
      category_id: categories.value[0]?.id || 0,
      description: "",
      image_url: "",
      status: 1,
    };
  }
  showModal.value = true;
};

const handleSubmit = async () => {
  if (!form.value.name) return ElMessage.error("请输入商品名称");
  if (
    form.value.original_price > 0 &&
    form.value.price > form.value.original_price
  ) {
    return ElMessage.error("实际售价不能高于市场原价");
  }

  try {
    form.value.category_id = Number(form.value.category_id);

    if (isEdit.value) {
      await updateProduct(form.value);
      ElMessage.success("商品信息已更新");
    } else {
      await createProduct(form.value);
      ElMessage.success("商品发布成功");
    }
    showModal.value = false;
    fetchProducts();
  } catch (error) {
    ElMessage.error(
      "操作失败: " + (error.response?.data?.msg || error.message),
    );
  }
};

const handleDelete = async (id) => {
  try {
    await deleteProduct(id);
    ElMessage.success("该商品已被永久删除");
    fetchProducts();
  } catch (error) {
    ElMessage.error("删除失败");
  }
};

onMounted(() => {
  fetchCategories();
  fetchProducts();
});
</script>

<style scoped>
/* 全局页面底色与容器 */
.admin-child-page {
  min-height: 100vh;
  background-color: #f8f9fa;
  padding-bottom: 80px;
}
.custom-container {
  max-width: 1280px;
  margin: 0 auto;
}

/* 顶部返回导航 */
.page-nav {
  padding: 24px 0 16px;
}
.back-btn {
  display: inline-flex;
  align-items: center;
  color: #606266;
  font-size: 15px;
  cursor: pointer;
  transition: color 0.3s;
  padding: 8px 16px 8px 0;
}
.back-btn:hover {
  color: #2d597b;
}
.back-icon {
  margin-right: 6px;
  font-size: 16px;
}

/* 统一的高光标题与顶部动作区 */
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 0 24px;
}
.page-title {
  display: inline-block;
  position: relative;
  font-size: 32px;
  color: #2c3e50;
  font-weight: 700;
  margin: 0;
  z-index: 1;
}
.page-title::after {
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
}

/* ★ 沉浸式搜索卡片 ★ */
.search-card {
  background: #ffffff;
  border-radius: 16px;
  padding: 24px 32px 4px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.02);
  margin-bottom: 24px;
}
.premium-search-form {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 8px;
}
.premium-search-form :deep(.el-form-item) {
  margin-bottom: 20px;
  margin-right: 16px;
}
.premium-search-form :deep(.el-form-item__label) {
  font-weight: 600;
  color: #303133;
}

:deep(.custom-input .el-input__wrapper),
:deep(.custom-select .el-select__wrapper) {
  box-shadow: 0 0 0 1px #dcdfe6 inset;
  border-radius: 8px;
  background: #fbfcfd;
  transition: all 0.3s;
  padding: 4px 12px;
}
:deep(.custom-input .el-input__wrapper.is-focus),
:deep(.custom-select .el-select__wrapper.is-focus) {
  box-shadow: 0 0 0 2px rgba(45, 89, 123, 0.2) inset !important;
  background: #ffffff;
}
:deep(.custom-checkbox .el-checkbox__label) {
  font-weight: 600;
  color: #303133;
}
:deep(.custom-checkbox.is-checked .el-checkbox__label) {
  color: #2d597b;
}
:deep(.custom-checkbox .el-checkbox__inner) {
  border-radius: 4px;
}
:deep(.custom-checkbox.is-checked .el-checkbox__inner) {
  background-color: #2d597b;
  border-color: #2d597b;
}

.checkbox-item {
  display: flex;
  align-items: center;
  margin-top: 0;
}
.form-actions {
  margin-left: auto;
  margin-top: 32px;
}

/* ★ 核心表格容器 ★ */
.table-wrapper {
  background: #ffffff;
  border-radius: 16px;
  padding: 24px 32px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.02);
}

/* Element Plus 表格定制 */
:deep(.custom-table) {
  --el-table-border-color: transparent;
  border-radius: 8px;
  overflow: hidden;
}
:deep(.custom-table th.el-table__cell) {
  font-weight: 600;
  font-size: 14px;
  padding: 18px 0;
  border-bottom: 1px solid #ebeef5;
  background: #fbfcfd;
}
:deep(.custom-table td.el-table__cell) {
  padding: 16px 0;
  border-bottom: 1px dashed #f0f2f5;
}
:deep(.custom-table::before) {
  display: none;
}

/* 单元格排版 */
.text-primary {
  color: #2c3e50;
  font-size: 15px;
}
.text-secondary {
  color: #a4b0be;
}

.product-thumb-wrapper {
  width: 56px;
  height: 56px;
  margin: 0 auto;
  background: #fbfcfd;
  border: 1px solid #ebeef5;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
}
.product-thumb {
  width: 100%;
  height: 100%;
  object-fit: contain;
  mix-blend-mode: multiply;
}
.no-img {
  font-size: 12px;
  color: #a4b0be;
}

.price-text {
  color: #e4393c;
  font-weight: 800;
  font-family:
    -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
  font-size: 16px;
}
.original-price-text {
  color: #a4b0be;
  text-decoration: line-through;
  font-size: 13px;
}

.category-tag {
  background: #f4f4f5;
  color: #606266;
  padding: 4px 10px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 500;
}

.status-badge {
  padding: 4px 12px;
  border-radius: 20px;
  font-size: 12px;
  font-weight: bold;
}
.is-active {
  background: #f0fdf4;
  color: #166534;
  border: 1px solid #bbf7d0;
}
.is-inactive {
  background: #fff7ed;
  color: #d97706;
  border: 1px solid #fed7aa;
}

.row-actions {
  display: flex;
  gap: 8px;
  justify-content: center;
}

/* 定制按钮 */
.action-btn {
  padding: 10px 24px;
  border-radius: 20px;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s;
  border: 1px solid transparent;
  display: inline-flex;
  align-items: center;
  justify-content: center;
}
.btn-sm {
  padding: 6px 16px;
  font-size: 13px;
}

.btn-primary {
  background: #2d597b;
  color: #ffffff;
  box-shadow: 0 4px 12px rgba(45, 89, 123, 0.2);
}
.btn-primary:hover:not(:disabled) {
  background: #1f435d;
  transform: translateY(-1px);
  box-shadow: 0 6px 16px rgba(45, 89, 123, 0.3);
}

.btn-outline {
  background: #ffffff;
  color: #2d597b;
  border-color: #2d597b;
}
.btn-outline:hover:not(:disabled) {
  background: #f0f7ff;
  transform: translateY(-1px);
}

.btn-danger-ghost {
  background: transparent;
  color: #f56c6c;
  border-color: #fbc4c4;
}
.btn-danger-ghost:hover {
  background: #fef0f0;
  color: #e4393c;
  transform: translateY(-1px);
}

.btn-default {
  background: #ffffff;
  color: #606266;
  border-color: #dcdfe6;
}
.btn-default:hover {
  color: #2d597b;
  border-color: #2d597b;
  background: #f0f7ff;
}

.action-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

/* 分页器 */
.pagination-container {
  display: flex;
  justify-content: flex-end;
  margin-top: 32px;
  padding-top: 16px;
  border-top: 1px solid #f0f2f5;
}
:deep(.custom-pagination .el-pager li.is-active) {
  background-color: #2d597b !important;
  color: #fff;
  border-radius: 4px;
}
:deep(.custom-pagination .el-pager li:hover) {
  color: #2d597b;
}

/* ================= 高级弹窗与表单美化 (双栏重构) ================= */
:deep(.premium-dialog) {
  border-radius: 16px;
  box-shadow: 0 16px 48px rgba(0, 0, 0, 0.1);
  overflow: hidden;
}
:deep(.premium-dialog .el-dialog__header) {
  margin-right: 0;
  padding: 24px 32px 20px;
  border-bottom: 1px solid #f0f2f5;
}
:deep(.premium-dialog .el-dialog__title) {
  font-weight: 700;
  color: #2c3e50;
  font-size: 18px;
  border-left: 4px solid #2d597b;
  padding-left: 10px;
}
:deep(.premium-dialog .el-dialog__body) {
  padding: 32px;
}
:deep(.premium-dialog .el-dialog__footer) {
  padding: 16px 32px 24px;
  border-top: 1px solid #f0f2f5;
  background: #fafbfc;
}

/* 网格双栏布局 */
.form-grid {
  display: flex;
  gap: 32px;
}
.form-col-left {
  flex: 1.2;
  display: flex;
  flex-direction: column;
}
.form-col-right {
  flex: 1;
  background: #fbfcfd;
  border: 1px solid #ebeef5;
  border-radius: 12px;
  padding: 20px;
  display: flex;
  flex-direction: column;
}

.premium-form :deep(.el-form-item__label) {
  font-weight: 600;
  color: #303133;
  padding-bottom: 6px;
}

:deep(.custom-input-number .el-input__wrapper) {
  box-shadow: 0 0 0 1px #dcdfe6 inset;
  border-radius: 8px;
  background: #fbfcfd;
  transition: all 0.3s;
}
:deep(.custom-input-number .el-input__wrapper.is-focus) {
  box-shadow: 0 0 0 2px rgba(45, 89, 123, 0.2) inset !important;
  background: #ffffff;
}
:deep(.highlight-number .el-input__inner) {
  color: #e4393c;
  font-weight: 800;
  font-size: 16px;
}

:deep(.custom-textarea .el-textarea__inner) {
  box-shadow: 0 0 0 1px #dcdfe6 inset;
  border-radius: 8px;
  background: #ffffff;
  transition: all 0.3s;
  padding: 12px;
  font-family: inherit;
  resize: none;
}
:deep(.custom-textarea .el-textarea__inner:focus) {
  box-shadow: 0 0 0 2px rgba(45, 89, 123, 0.2) inset !important;
  background: #ffffff;
}

.bold-check {
  margin-top: 8px;
}
.bold-check :deep(.el-checkbox__label) {
  font-weight: 600;
  color: #2d597b;
}
.mb-0 {
  margin-bottom: 0 !important;
}
.mt-2 {
  margin-top: 8px !important;
}

/* 极简紧凑型高级上传控件 */
.custom-upload-box.compact-upload {
  width: 100%;
  border: 2px dashed #dcdfe6;
  border-radius: 10px;
  background: #ffffff;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 24px 20px;
  cursor: pointer;
  transition: all 0.3s;
  position: relative;
  overflow: hidden;
  min-height: 160px;
}
.custom-upload-box.compact-upload:hover {
  border-color: #2d597b;
  background: #f0f7ff;
}

.upload-placeholder {
  display: flex;
  flex-direction: column;
  align-items: center;
  color: #606266;
  font-size: 13px;
  font-weight: 500;
}
.compact-upload .upload-icon {
  font-size: 32px;
  color: #a4b0be;
  margin-bottom: 8px;
  transition: color 0.3s;
}
.custom-upload-box.compact-upload:hover .upload-icon {
  color: #2d597b;
}

.upload-preview {
  width: 100%;
  display: flex;
  justify-content: center;
}
.upload-preview img {
  max-height: 120px;
  border-radius: 8px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}
.preview-overlay {
  position: absolute;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 600;
  opacity: 0;
  transition: opacity 0.3s;
}
.custom-upload-box.compact-upload:hover .preview-overlay {
  opacity: 1;
}

.upload-status {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #2d597b;
  font-weight: 600;
}
.is-loading {
  animation: rotating 2s linear infinite;
}

.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

@media (max-width: 768px) {
  .premium-search-form {
    flex-direction: column;
    align-items: stretch;
  }
  .premium-search-form :deep(.el-form-item) {
    margin-right: 0;
  }
  .form-actions {
    margin-left: 0;
    justify-content: flex-end;
  }
  .form-grid {
    flex-direction: column;
    gap: 24px;
  }
}
</style>
