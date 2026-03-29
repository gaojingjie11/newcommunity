<template>
  <div class="admin-child-page">
    <Navbar />
    <div class="container">
      <div class="page-header">
        <h1 class="page-title">商品管理</h1>
        <el-button type="primary" @click="openModal()">+ 发布商品</el-button>
      </div>

      <!-- Search Filters -->
      <el-card class="search-card mb-4">
        <el-form :inline="true" :model="filters" class="demo-form-inline">
          <el-form-item label="商品名称">
            <el-input v-model="filters.name" placeholder="输入名称搜索" clearable @keyup.enter="handleSearch" />
          </el-form-item>
          <el-form-item label="分类">
            <el-select v-model="filters.category_id" placeholder="所有分类" clearable style="width: 150px">
              <el-option label="所有分类" :value="0" />
              <el-option v-for="c in categories" :key="c.id" :label="c.name" :value="c.id" />
            </el-select>
          </el-form-item>
           <el-form-item>
            <el-checkbox v-model="filters.is_promotion" label="仅看促销" border />
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="handleSearch">搜索</el-button>
            <el-button @click="resetFilters">重置</el-button>
          </el-form-item>
        </el-form>
      </el-card>

      <div class="table-container card">
         <el-table :data="products" stripe border style="width: 100%" v-loading="loadingData">
          <el-table-column prop="id" label="ID" width="80" />
          <el-table-column label="图片" width="100">
             <template #default="scope">
                <img :src="scope.row.image_url" class="product-thumb" alt="无图" 
                     style="width: 50px; height: 50px; object-fit: cover; border-radius: 4px;" 
                     v-if="scope.row.image_url"/>
                <span v-else>无图</span>
             </template>
          </el-table-column>
          <el-table-column prop="name" label="名称" show-overflow-tooltip/>
          <el-table-column prop="price" label="促销价" width="120">
             <template #default="scope">
                 <span style="color:red; font-weight:bold;">¥{{ scope.row.price }}</span>
             </template>
          </el-table-column>
          <el-table-column prop="original_price" label="原价" width="120">
             <template #default="scope">
                 <span v-if="scope.row.original_price" style="text-decoration: line-through; color: #999;">
                     ¥{{ scope.row.original_price }}
                 </span>
                 <span v-else>-</span>
             </template>
          </el-table-column>
          <el-table-column prop="stock" label="库存" width="100" />
          <el-table-column prop="category_name" label="分类" width="120" />
          <el-table-column prop="status" label="状态" width="100">
             <template #default="scope">
               <el-tag :type="scope.row.status === 1 ? 'success' : 'warning'">
                 {{ scope.row.status === 1 ? '上架' : '下架' }}
               </el-tag>
             </template>
          </el-table-column>
          <el-table-column label="操作" width="180" fixed="right">
            <template #default="scope">
              <el-button size="small" type="primary" link @click="openModal(scope.row)">编辑</el-button>
              <el-popconfirm
                title="确定删除该商品吗？"
                @confirm="handleDelete(scope.row.id)"
              >
                <template #reference>
                  <el-button size="small" type="danger" link>删除</el-button>
                </template>
              </el-popconfirm>
            </template>
          </el-table-column>
        </el-table>

         <div class="pagination-container">
          <el-pagination
            v-model:current-page="page"
            v-model:page-size="size"
            :page-sizes="[10, 20, 50, 100]"
            layout="total, sizes, prev, pager, next, jumper"
            :total="total"
            @size-change="handleSearch"
            @current-change="fetchProducts"
          />
        </div>
      </div>

      <!-- Edit Dialog -->
      <el-dialog
        v-model="showModal"
        :title="isEdit ? '编辑商品' : '发布商品'"
        width="600px"
      >
        <el-form ref="formRef" :model="form" label-width="80px">
          <el-form-item label="商品名称" required>
            <el-input v-model="form.name" />
          </el-form-item>
          
          <el-row :gutter="20">
             <el-col :span="12">
               <el-form-item label="原价" required>
                  <el-input-number v-model="form.original_price" :min="0" :precision="2" style="width: 100%" placeholder="市场价/原价"/>
               </el-form-item>
            </el-col>
            <el-col :span="12">
               <el-form-item label="促销价" required>
                  <el-input-number v-model="form.price" :min="0" :precision="2" style="width: 100%" placeholder="实际售卖价格"/>
               </el-form-item>
            </el-col>
          </el-row>
          
          <el-row :gutter="20">
             <el-col :span="12">
               <el-form-item label="库存" required>
                  <el-input-number v-model="form.stock" :min="0" :step="1" style="width: 100%" />
               </el-form-item>
            </el-col>
             <el-col :span="12">
                <el-form-item label="分类">
                  <el-select v-model="form.category_id" placeholder="Select" style="width: 100%">
                    <el-option label="请选择分类" :value="0" />
                    <el-option v-for="c in categories" :key="c.id" :label="c.name" :value="c.id" />
                  </el-select>
              </el-form-item>
            </el-col>
          </el-row>

           <el-form-item label="描述">
             <el-input v-model="form.description" type="textarea" :rows="3" />
          </el-form-item>
          
           <el-form-item label="商品图片">
               <!-- Simple file input for now until upload component is refined -->
              <input type="file" @change="handleUpload" accept="image/*" />
               <div v-if="uploading">上传中...</div>
               <div v-if="form.image_url" style="margin-top: 10px">
                   <img :src="form.image_url" style="max-height: 100px" />
               </div>
          </el-form-item>

           <el-form-item>
             <el-checkbox v-model="form.status" :true-label="1" :false-label="0">上架销售</el-checkbox>
          </el-form-item>
        </el-form>
        
        <template #footer>
          <span class="dialog-footer">
            <el-button @click="showModal = false">取消</el-button>
            <el-button type="primary" @click="handleSubmit" :loading="uploading">
              保存
            </el-button>
          </span>
        </template>
      </el-dialog>

    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, reactive } from 'vue'
import Navbar from '@/components/layout/Navbar.vue'
import { getProductList, getCategories } from '@/api/product'
import { createProduct, updateProduct, deleteProduct } from '@/api/admin'
import request from '@/utils/request'
import { ElMessage } from 'element-plus'

const products = ref([])
const categories = ref([])
const showModal = ref(false)
const isEdit = ref(false)
const uploading = ref(false)
const loadingData = ref(false)
const total = ref(0)
const page = ref(1)
const size = ref(10)

const form = ref({
  id: 0,
  name: '',
  price: 0,
  original_price: 0,
  stock: 0,
  category_id: 0,
  description: '',
  image_url: '',
  status: 1
})

const filters = reactive({
    name: '',
    category_id: 0,
    is_promotion: false
})

const handleUpload = async (event) => {
  const file = event.target.files[0]
  if (!file) return
  
  const formData = new FormData()
  formData.append('file', file)
  
  uploading.value = true
  try {
    const res = await request({
        url: '/upload',
        method: 'post',
        data: formData,
        headers: { 'Content-Type': 'multipart/form-data' }
    })
    form.value.image_url = res.url
  } catch (e) {
    const msg = e?.response?.data?.msg || e?.message || '上传失败'
    ElMessage.error(msg)
  } finally {
    uploading.value = false
    // Clear input value to allow re-uploading same file
    event.target.value = ''
  }
}

// ... (previous handleUpload code) ...

const handleSearch = () => {
    page.value = 1
    fetchProducts()
}

const fetchProducts = async () => {
  loadingData.value = true
  try {
    const params = {
        name: filters.name,
        category_id: filters.category_id || undefined,
        is_promotion: filters.is_promotion || undefined,
        page: page.value,
        size: size.value
    }
    const res = await getProductList(params) 
    products.value = res.list || []
    total.value = res.total || 0
  } catch (error) {
    console.error(error)
  } finally {
      loadingData.value = false
  }
}

const resetFilters = () => {
    filters.name = ''
    filters.category_id = 0
    filters.is_promotion = false
    handleSearch()
}

const fetchCategories = async () => {
    try {
        const res = await getCategories()
        categories.value = res || []
    } catch (e) {
        console.error(e)
    }
}

const openModal = (product = null) => {
  isEdit.value = !!product
  
  // Reset form completely
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
        status: product.status
    }
  } else {
    // Ensure image_url is empty for new product
    form.value = { 
        id: 0, 
        name: '', 
        price: 0, 
        original_price: 0, 
        stock: 999, 
        category_id: categories.value[0]?.id || 0, 
        description: '', 
        image_url: '', 
        status: 1 
    }
  }
  showModal.value = true
}

const handleSubmit = async () => {
    // Basic validation
    if (!form.value.name) return ElMessage.error('请输入商品名称')
    
    // Price validation
    if (form.value.original_price > 0 && form.value.price > form.value.original_price) {
        return ElMessage.error('促销价不能高于原价')
    }
  
    try {
      form.value.category_id = Number(form.value.category_id)
      
      if (isEdit.value) {
        await updateProduct(form.value)
        ElMessage.success('更新成功')
      } else {
        await createProduct(form.value)
        ElMessage.success('创建成功')
      }
      showModal.value = false
      fetchProducts()
    } catch (error) {
      ElMessage.error('操作失败: ' + (error.response?.data?.msg || error.message))
    }
}

const handleDelete = async (id) => {
  try {
    await deleteProduct(id)
    ElMessage.success('删除成功')
    fetchProducts()
  } catch (error) {
    ElMessage.error('删除失败')
  }
}

onMounted(() => {
  fetchCategories()
  fetchProducts()
})
</script>

<style scoped>
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

.search-card {
    margin-bottom: 20px;
}

.table-container {
  padding: 20px;
}

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}
</style>
