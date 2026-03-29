<template>
  <div class="admin-child-page">
    <Navbar />
    <div class="container">
      <div class="page-header">
        <div class="page-nav">
          <div class="back-btn" @click="$router.push('/admin')">
            <el-icon class="back-icon"><ArrowLeft /></el-icon>
            <span></span>
          </div>
        </div>
        <button
          class="btn btn-primary"
          style="margin-left: 1000px; margin-top: 40px"
          @click="openModal()"
        >
          + 添加门店
        </button>
      </div>

      <div class="table-container card">
        <el-table :data="stores" style="width: 100%" stripe border>
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

          <el-table-column label="操作" width="150" fixed="right">
            <template #default="{ row }">
              <button class="btn btn-sm btn-info" @click="openModal(row)">
                编辑
              </button>
              <button
                class="btn btn-sm btn-danger ml-2"
                @click="handleDelete(row.id)"
              >
                删除
              </button>
            </template>
          </el-table-column>
        </el-table>
      </div>

      <!-- 弹窗 -->
      <div class="modal-overlay" v-if="showModal">
        <div class="modal card">
          <h3>{{ isEdit ? "编辑门店" : "添加门店" }}</h3>
          <form @submit.prevent="handleSubmit">
            <div class="form-group" style="margin-top: 20px">
              <label>门店名称</label>
              <input v-model="form.name" class="input" required />
            </div>
            <div class="form-group">
              <label>地址</label>
              <input v-model="form.address" class="input" required />
            </div>
            <div class="form-group">
              <label>电话</label>
              <input v-model="form.phone" class="input" required />
            </div>
            <div class="form-group">
              <label>区域</label>
              <input
                v-model="form.region"
                class="input"
                placeholder="如: A区"
              />
            </div>
            <div class="form-group">
              <label>营业时间</label>
              <input
                v-model="form.business_hours"
                class="input"
                placeholder="如: 09:00 - 22:00"
              />
            </div>

            <div class="modal-actions">
              <button
                type="button"
                class="btn btn-secondary"
                @click="closeModal"
              >
                取消
              </button>
              <button type="submit" class="btn btn-primary">保存</button>
            </div>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from "vue";
import Navbar from "@/components/layout/Navbar.vue";
import request from "@/utils/request"; // 没有封装store list API，这里直接调或者补充API
import { createStore, updateStore, deleteStore } from "@/api/admin";
import { ElMessage, ElMessageBox } from "element-plus";

// 临时补充 getStores，因为 api/service 里好像没暴露给admin用
const getStores = () => {
  return request({ url: "/stores", method: "get" });
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

const fetchStores = async () => {
  try {
    const res = await getStores();
    stores.value = res;
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

onMounted(fetchStores);
</script>

<style scoped>
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
.back-icon {
  margin-right: 6px;
  font-size: 16px;
}
.admin-child-page {
  min-height: 100vh;
  padding-bottom: var(--spacing-xl);
}
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--spacing-lg);
}
.table-container {
  min-height: 600px;
}
.table {
  width: 100%;
  border-collapse: collapse;
}
.table th,
.table td {
  padding: 12px;
  border-bottom: 1px solid #eee;
  text-align: left;
}
.ml-2 {
  margin-left: 8px;
}

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
</style>
