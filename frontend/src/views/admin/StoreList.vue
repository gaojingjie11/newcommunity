<template>
  <div class="admin-child-page">
    <Navbar />
    <div class="container custom-container">
      <!-- 极简顶部区：只保留返回导航与操作按钮，去除了大标题 -->
      <div class="top-bar">
        <div class="back-btn" @click="$router.push('/admin')">
          <el-icon class="back-icon"><ArrowLeft /></el-icon>
          <span>返回管理后台</span>
        </div>
        <div class="header-actions">
          <button class="action-btn btn-primary" @click="openModal()">
            + 添加门店
          </button>
        </div>
      </div>

      <div class="table-wrapper">
        <el-table :data="stores" class="custom-table" style="width: 100%">
          <el-table-column prop="id" label="ID" width="80" />
          <el-table-column
            prop="name"
            label="名称"
            min-width="150"
            show-overflow-tooltip
          />
          <el-table-column
            prop="address"
            label="地址"
            min-width="200"
            show-overflow-tooltip
          />
          <el-table-column prop="phone" label="电话" width="120" />
          <el-table-column prop="business_hours" label="营业时间" width="150" />

          <el-table-column label="操作" width="280" fixed="right" align="center">
            <template #default="{ row }">
              <div class="row-actions">
                <button class="action-btn btn-sm btn-outline" @click="openModal(row)">
                  编辑
                </button>
                <button class="action-btn btn-sm btn-primary" @click="openStoreProducts(row)">
                  商品管理
                </button>
                <button
                  class="action-btn btn-sm btn-danger-ghost"
                  @click="handleDelete(row.id)"
                >
                  删除
                </button>
              </div>
            </template>
          </el-table-column>
        </el-table>
      </div>

      <!-- 门店基本信息弹窗 -->
      <el-dialog
        v-model="showModal"
        :title="isEdit ? '编辑门店' : '添加门店'"
        width="500px"
        class="premium-dialog"
        destroy-on-close
      >
        <el-form :model="form" label-width="80px" @submit.prevent="handleSubmit">
          <el-form-item label="门店名称" required>
            <el-input v-model="form.name" placeholder="请输入门店名称" />
          </el-form-item>
          <el-form-item label="地址" required>
            <el-input v-model="form.address" placeholder="请输入地址" />
          </el-form-item>
          <el-form-item label="电话" required>
            <el-input v-model="form.phone" placeholder="请输入电话" />
          </el-form-item>
          <el-form-item label="区域">
            <el-input v-model="form.region" placeholder="如: A区" />
          </el-form-item>
          <el-form-item label="营业时间">
            <el-input v-model="form.business_hours" placeholder="如: 09:00 - 22:00" />
          </el-form-item>
        </el-form>
        <template #footer>
          <div class="dialog-footer" style="display: flex; justify-content: flex-end; gap: 12px;">
            <el-button @click="closeModal">取消</el-button>
            <el-button type="primary" @click="handleSubmit">保存</el-button>
          </div>
        </template>
      </el-dialog>

      <!-- 门店商品管理弹窗 -->
      <el-dialog
        v-model="showStoreProductsDialog"
        :title="`门店商品管理 - ${currentStore?.name || ''}`"
        width="1000px"
        class="premium-dialog"
        destroy-on-close
      >
        <div class="dialog-action-bar" style="margin-bottom: 20px; display: flex; justify-content: space-between; align-items: center;">
          <span style="font-size: 14px; color: #606266;">
            提示：在此为门店分配商品，并维护各门店独立的库存与上架状态。
          </span>
          <el-button type="primary" size="default" @click="openBindDialog" :icon="Plus">
            绑定新商品
          </el-button>
        </div>

        <el-table :data="storeProducts" class="custom-table" style="width: 100%" v-loading="loadingProducts">
          <el-table-column prop="product_id" label="商品ID" width="90" align="center" />
          <el-table-column label="商品图片" width="100" align="center">
            <template #default="{ row }">
              <div class="product-thumb-wrapper">
                <img 
                  :src="row.product?.image_url || DEFAULT_PRODUCT_IMAGE" 
                  class="product-thumb"
                  alt="商品"
                />
              </div>
            </template>
          </el-table-column>
          <el-table-column label="商品名称" min-width="150" show-overflow-tooltip>
            <template #default="{ row }">
              <span style="font-weight: 600;">{{ row.product?.name || '未知商品' }}</span>
            </template>
          </el-table-column>
          <el-table-column label="分类" width="120" align="center">
            <template #default="{ row }">
              <span>{{ row.product?.category_name || '-' }}</span>
            </template>
          </el-table-column>
          <el-table-column label="单价" width="100" align="center">
            <template #default="{ row }">
              <span style="color: #e4393c; font-weight: bold;">¥{{ row.product?.price }}</span>
            </template>
          </el-table-column>
          <el-table-column label="门店库存" width="150" align="center">
            <template #default="{ row }">
              <el-input-number 
                v-model="row.stock" 
                :min="0" 
                size="small" 
                style="width: 110px;"
                @change="(val) => handleUpdateStock(row, val)" 
              />
            </template>
          </el-table-column>
          <el-table-column label="上架状态" width="120" align="center">
            <template #default="{ row }">
              <el-switch 
                v-model="row.status" 
                :active-value="1" 
                :inactive-value="0"
                active-text="上架"
                inactive-text="下架"
                inline-prompt
                @change="(val) => handleToggleStatus(row, val)" 
              />
            </template>
          </el-table-column>
          <el-table-column label="操作" width="100" align="center" fixed="right">
            <template #default="{ row }">
              <el-button type="danger" size="small" @click="handleUnbind(row)">解绑</el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-dialog>

      <!-- 绑定商品二级弹窗 -->
      <el-dialog
        v-model="showBindDialog"
        title="绑定新商品到门店"
        width="450px"
        class="premium-dialog"
        append-to-body
      >
        <el-form :model="bindForm" label-width="80px">
          <el-form-item label="选择商品" required>
            <el-select 
              v-model="bindForm.product_id" 
              filterable 
              placeholder="请输入商品名称搜索..."
              style="width: 100%"
            >
              <el-option
                v-for="p in availableProducts"
                :key="p.id"
                :label="`${p.name} (单价: ¥${p.price} | 系统库存: ${p.stock})`"
                :value="p.id"
              />
            </el-select>
          </el-form-item>
          <el-form-item label="分配库存" required>
            <el-input-number 
              v-model="bindForm.stock" 
              :min="0" 
              style="width: 100%"
            />
          </el-form-item>
        </el-form>
        <template #footer>
          <div class="dialog-footer" style="display: flex; justify-content: flex-end; gap: 12px;">
            <el-button @click="showBindDialog = false">取消</el-button>
            <el-button type="primary" @click="submitBind" :loading="bindingLoading">确认绑定</el-button>
          </div>
        </template>
      </el-dialog>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from "vue";
import Navbar from "@/components/layout/Navbar.vue";
import request from "@/utils/request";
import { DEFAULT_PRODUCT_IMAGE } from "@/utils/constants";
import { 
  createStore, 
  updateStore, 
  deleteStore,
  getStoreProducts,
  bindStoreProduct,
  unbindStoreProduct,
  updateStoreProductStatus,
  updateStoreProductStock
} from "@/api/admin";
import { getProductList } from "@/api/product";
import { ElMessage, ElMessageBox } from "element-plus";
import { Plus, ArrowLeft } from "@element-plus/icons-vue";

// 临时补充 getStores，因为 api/service 里好像没暴露给admin用
const getStores = () => {
  return request({ url: "/mall/stores", method: "get" });
};

const stores = ref([]);
const showModal = ref(false);
const isEdit = ref(false);
const form = ref({
  id: 0,
  name: "",
  address: "",
  phone: "",
  region: "",
  business_hours: "",
});

const showStoreProductsDialog = ref(false);
const loadingProducts = ref(false);
const storeProducts = ref([]);
const currentStore = ref(null);

const systemProducts = ref([]);
const showBindDialog = ref(false);
const bindingLoading = ref(false);
const bindForm = ref({
  product_id: "",
  stock: 0,
});

const fetchStores = async () => {
  try {
    const res = await getStores();
    stores.value = res?.list || res || [];
  } catch (e) {
    console.error(e);
  }
};

const openModal = (store = null) => {
  isEdit.value = !!store;
  if (store) {
    form.value = { ...store };
  } else {
    form.value = {
      id: 0,
      name: "",
      address: "",
      phone: "",
      region: "",
      business_hours: "",
    };
  }
  showModal.value = true;
};

const closeModal = () => (showModal.value = false);

const handleSubmit = async () => {
  try {
    if (isEdit.value) {
      await updateStore(form.value);
    } else {
      await createStore(form.value);
    }
    ElMessage.success("保存成功");
    closeModal();
    fetchStores();
  } catch (e) {
    ElMessage.error("操作失败");
  }
};

const handleDelete = async (id) => {
  try {
    await ElMessageBox.confirm("确定删除?", "删除确认", {
      confirmButtonText: "删除",
      cancelButtonText: "取消",
      type: "warning",
    });
    await deleteStore(id);
    ElMessage.success("删除成功");
    fetchStores();
  } catch (e) {
    if (e !== "cancel") {
      ElMessage.error("删除失败");
    }
  }
};

const openStoreProducts = async (store) => {
  currentStore.value = store;
  showStoreProductsDialog.value = true;
  await loadSystemProducts();
  await fetchStoreProducts();
};

const fetchStoreProducts = async () => {
  if (!currentStore.value) return;
  loadingProducts.value = true;
  try {
    const res = await getStoreProducts(currentStore.value.id);
    const rawList = res?.list || res || [];
    storeProducts.value = rawList.map(item => {
      const prod = systemProducts.value.find(p => p.id === item.product_id);
      return {
        ...item,
        product: prod || {
          name: item.product_name || '未知商品',
          price: item.price,
          image_url: '',
          category_name: '-'
        }
      };
    });
  } catch (e) {
    console.error(e);
    ElMessage.error("获取门店商品列表失败");
  } finally {
    loadingProducts.value = false;
  }
};

const loadSystemProducts = async () => {
  try {
    const res = await getProductList({ page: 1, size: 1000 });
    systemProducts.value = res?.list || res || [];
  } catch (e) {
    console.error(e);
  }
};

const availableProducts = computed(() => {
  const boundIds = new Set(storeProducts.value.map((sp) => sp.product_id));
  return systemProducts.value.filter((p) => !boundIds.has(p.id));
});

const openBindDialog = () => {
  bindForm.value = {
    product_id: "",
    stock: 10,
  };
  showBindDialog.value = true;
};

const submitBind = async () => {
  if (!bindForm.value.product_id) {
    ElMessage.warning("请选择要绑定的商品");
    return;
  }
  bindingLoading.value = true;
  try {
    await bindStoreProduct({
      store_id: currentStore.value.id,
      product_id: bindForm.value.product_id,
      stock: bindForm.value.stock,
    });
    ElMessage.success("绑定商品成功");
    showBindDialog.value = false;
    await fetchStoreProducts();
  } catch (e) {
    console.error(e);
    ElMessage.error(e.response?.data?.message || "绑定商品失败");
  } finally {
    bindingLoading.value = false;
  }
};

const handleUpdateStock = async (row, val) => {
  try {
    await updateStoreProductStock({
      store_id: row.store_id,
      product_id: row.product_id,
      stock: val,
    });
    ElMessage.success("修改库存成功");
  } catch (e) {
    console.error(e);
    ElMessage.error("修改库存失败");
    await fetchStoreProducts();
  }
};

const handleToggleStatus = async (row, val) => {
  try {
    await updateStoreProductStatus({
      store_id: row.store_id,
      product_id: row.product_id,
      status: val,
    });
    ElMessage.success(val === 1 ? "商品已上架" : "商品已下架");
  } catch (e) {
    console.error(e);
    ElMessage.error("操作失败");
    await fetchStoreProducts();
  }
};

const handleUnbind = async (row) => {
  try {
    await ElMessageBox.confirm(`确定解除商品「${row.product?.name || '未知商品'}」的门店绑定吗？`, "解除绑定", {
      confirmButtonText: "解除",
      cancelButtonText: "取消",
      type: "warning",
    });
    await unbindStoreProduct({
      store_id: row.store_id,
      product_id: row.product_id,
    });
    ElMessage.success("解绑成功");
    await fetchStoreProducts();
  } catch (e) {
    if (e !== "cancel") {
      console.error(e);
      ElMessage.error("解绑失败");
    }
  }
};

onMounted(fetchStores);
</script>

<style scoped>
/* 全局页面底色与容器 */
.admin-child-page { min-height: 100vh; background-color: #f8f9fa; padding-bottom: 80px; }
.custom-container { max-width: 1280px; margin: 0 auto; }

/* 极简顶部区：只保留返回和操作 */
.top-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 32px 0 24px;
}

.back-btn {
  display: inline-flex; align-items: center; color: #606266; font-size: 16px; font-weight: 600;
  cursor: pointer; transition: color 0.3s; padding: 8px 16px 8px 0;
}
.back-btn:hover { color: #2d597b; }
.back-icon { margin-right: 6px; font-size: 18px; }

.header-actions { display: flex; align-items: center; gap: 12px; }

/* 核心表格容器 */
.table-wrapper { background: #ffffff; border-radius: 16px; padding: 24px 32px; box-shadow: 0 4px 20px rgba(0, 0, 0, 0.02); }

/* Element Plus 表格定制去后台感 */
:deep(.custom-table) { --el-table-border-color: transparent; border-radius: 8px; overflow: hidden; }
:deep(.custom-table th.el-table__cell) { font-weight: 600; font-size: 14px; padding: 18px 0; border-bottom: 1px solid #ebeef5; background: #fbfcfd; color: #606266; }
:deep(.custom-table td.el-table__cell) { padding: 20px 0; border-bottom: 1px dashed #f0f2f5; font-size: 14px; }
:deep(.custom-table::before) { display: none; }

/* 列表操作区 */
.row-actions { display: flex; gap: 8px; justify-content: center; }

/* 定制化按钮 */
.action-btn { padding: 10px 24px; border-radius: 20px; font-size: 14px; font-weight: 600; cursor: pointer; transition: all 0.3s; border: 1px solid transparent; display: inline-flex; align-items: center; justify-content: center; }
.btn-sm { padding: 6px 16px; font-size: 13px; }

.btn-primary { background: #2d597b; color: #ffffff; box-shadow: 0 4px 12px rgba(45, 89, 123, 0.2); }
.btn-primary:hover:not(:disabled) { background: #1f435d; transform: translateY(-2px); box-shadow: 0 6px 16px rgba(45, 89, 123, 0.3); }

.btn-outline { background: #ffffff; color: #2d597b; border-color: #2d597b; }
.btn-outline:hover { background: #f0f7ff; transform: translateY(-1px); }

.btn-danger-ghost { background: transparent; color: #f56c6c; border-color: #fbc4c4; }
.btn-danger-ghost:hover { background: #fef0f0; color: #e4393c; transform: translateY(-1px); }

/* Modal Styles - Reuse */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 2000;
}
.modal {
  padding: 24px;
  width: 400px;
  max-width: 90%;
}
.form-group {
  margin-bottom: 16px;
  display: flex;
  flex-direction: column;
}
.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  margin-top: 24px;
}

/* Premium Dialog Styling */
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
  padding: 24px 32px 12px;
}
:deep(.premium-dialog .el-dialog__footer) {
  padding: 16px 32px 24px;
  border-top: 1px solid #f0f2f5;
  background: #fafbfc;
}

.product-thumb-wrapper {
  width: 50px;
  height: 50px;
  margin: 0 auto;
  background: #fbfcfd;
  border: 1px solid #ebeef5;
  border-radius: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
}
.product-thumb {
  width: 100%;
  height: 100%;
  object-fit: cover;
}
</style>
