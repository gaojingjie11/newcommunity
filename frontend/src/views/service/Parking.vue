<template>
  <div class="parking-page">
    <Navbar />
    
    <div class="container custom-container">
      <div class="page-header">
        <h1 class="page-title highlight-title">车位管理</h1>
      </div>
      
      <div v-if="parkingList && parkingList.length > 0" class="parking-list">
          <div class="parking-card" v-for="item in parkingList" :key="item.id">
            <div class="card-header">
              <h3>我的车位 <span class="parking-no">{{item.parking_no}}</span></h3>
              <span class="status-tag" :class="item.status === 1 ? 'is-occupied' : 'is-free'">
                {{ item.status === 1 ? '已占用' : '空闲中' }}
              </span>
            </div>
            
            <div class="info-grid">
              <div class="info-item">
                <span class="info-label">绑定车牌</span>
                <span class="plate-badge" v-if="item.car_plate">{{ item.car_plate }}</span>
                <span class="info-value text-muted" v-else>未绑定</span>
              </div>
            </div>
            
            <div class="bind-form">
              <div class="form-title">{{ item.car_plate ? '更新车牌信息' : '绑定车辆' }}</div>
              <div class="form-row">
                <div class="input-wrapper">
                  <input 
                    v-model="item.editCarPlate" 
                    class="custom-input" 
                    :placeholder="item.car_plate ? '输入新车牌号' : '请输入车牌号'"
                  />
                </div>
                <button class="btn-action" :class="{'is-loading': loading}" :disabled="loading" @click="handleBindCar(item)">
                  {{ loading ? '处理中...' : '确认更新' }}
                </button>
              </div>
            </div>
          </div>
      </div>
      
      <div class="empty-wrapper" v-else>
        <el-empty description="系统未查到您的车位信息，请联系物业分配" image-size="160" />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import Navbar from '@/components/layout/Navbar.vue'
import { getMyParking, bindCar } from '@/api/service'
import { ElMessage } from 'element-plus'

const parkingList = ref([])
const loading = ref(false)

const fetchParking = async () => {
  try {
    const data = await getMyParking()
    if (Array.isArray(data)) {
        parkingList.value = data.map(p => ({
            ...p,
            editCarPlate: p.car_plate 
        }))
    } else if (data) {
        parkingList.value = [{...data, editCarPlate: data.car_plate}]
    } else {
        parkingList.value = []
    }
  } catch (error) {
    console.log('未分配车位或获取失败', error)
  }
}

const handleBindCar = async (item) => {
  if (!item.editCarPlate || !item.editCarPlate.trim()) {
    ElMessage.warning('请输入车牌号')
    return
  }
  
  loading.value = true
  try {
    await bindCar({ 
        parking_id: item.id, 
        car_plate: item.editCarPlate 
    })
    ElMessage.success('更新成功！')
    await fetchParking()
  } catch (error) {
    ElMessage.error('操作失败：' + (error.response?.data?.msg || error.message))
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchParking()
})
</script>

<style scoped>
.parking-page { min-height: 100vh; background-color: #f8f9fa; padding-bottom: 80px; }
.custom-container { max-width: 800px; margin: 0 auto; }
.page-header { padding: 32px 0 24px; }
.highlight-title { display: inline-block; position: relative; font-size: 32px; color: #2c3e50; font-weight: 700; margin: 0; z-index: 1; }
.highlight-title::after { content: ''; position: absolute; bottom: 4px; left: -5%; width: 110%; height: 14px; background-color: #2d597b; opacity: 0.15; border-radius: 6px; z-index: -1; }

.parking-list { display: flex; flex-direction: column; gap: 24px; }

.parking-card {
  background: #ffffff; border-radius: 16px; padding: 32px; box-shadow: 0 4px 20px rgba(0,0,0,0.03);
}

.card-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 24px; padding-bottom: 20px; border-bottom: 1px dashed #ebeef5; }
.card-header h3 { margin: 0; font-size: 20px; color: #2c3e50; }
.parking-no { color: #2d597b; font-weight: 800; font-size: 22px; margin-left: 8px; font-family: monospace; }

.status-tag { padding: 6px 16px; border-radius: 20px; font-size: 13px; font-weight: bold; }
.is-occupied { background: #f0fdf4; color: #166534; border: 1px solid #dcfce7; }
.is-free { background: #fff7ed; color: #9a3412; border: 1px solid #ffedd5; }

.info-grid { background: #fbfcfd; padding: 24px; border-radius: 12px; margin-bottom: 32px; }
.info-item { display: flex; align-items: center; justify-content: space-between; }
.info-label { color: #606266; font-size: 15px; }
.text-muted { color: #a4b0be; font-style: italic; }

/* 拟物车牌样式 */
.plate-badge {
  background: #1e40af; color: #ffffff; padding: 8px 20px; border-radius: 6px;
  font-size: 18px; font-weight: bold; letter-spacing: 2px; box-shadow: inset 0 0 0 2px rgba(255,255,255,0.2);
  border: 1px solid #1e3a8a; font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
}

.form-title { font-size: 16px; font-weight: 600; color: #303133; margin-bottom: 16px; }
.form-row { display: flex; gap: 16px; }
.input-wrapper { flex: 1; }
.custom-input {
  width: 100%; padding: 12px 16px; border: 1px solid #dcdfe6; border-radius: 8px; font-size: 15px;
  outline: none; transition: all 0.3s; background: #fafbfc;
}
.custom-input:focus { border-color: #2d597b; background: #ffffff; box-shadow: 0 0 0 3px rgba(45,89,123,0.1); }

.btn-action {
  padding: 0 32px; background: #2d597b; color: #ffffff; border: none; border-radius: 8px;
  font-size: 15px; font-weight: bold; cursor: pointer; transition: all 0.3s;
}
.btn-action:hover:not(:disabled) { background: #1f435d; }
.btn-action:disabled { background: #a4b0be; cursor: not-allowed; }

.empty-wrapper { background: #ffffff; border-radius: 12px; padding: 80px 0; box-shadow: 0 2px 12px rgba(0,0,0,0.02); }
</style>