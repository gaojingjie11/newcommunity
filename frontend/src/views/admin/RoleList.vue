<template>
  <div class="admin-child-page">
    <Navbar />
    <div class="container custom-container">
      <!-- 顶部条：返回与新增 -->
      <div class="top-bar">
        <div class="back-btn" @click="$router.push('/admin')">
          <el-icon class="back-icon"><ArrowLeft /></el-icon>
          <span>返回管理后台</span>
        </div>
        <div class="header-actions">
          <button class="action-btn btn-primary" @click="openModal()">
            + 添加角色
          </button>
        </div>
      </div>

      <!-- 角色列表表格 -->
      <div class="table-wrapper">
        <el-table :data="roles" class="custom-table" style="width: 100%" v-loading="loading">
          <el-table-column prop="id" label="ID" width="80" />
          <el-table-column prop="name" label="角色名称" width="180" />
          <el-table-column prop="code" label="角色标识" width="180">
            <template #default="{ row }">
              <el-tag size="small" type="info" class="custom-tag">{{ row.code }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="remark" label="描述/备注" min-width="200" show-overflow-tooltip />

          <el-table-column label="操作" width="300" fixed="right" align="center">
            <template #default="{ row }">
              <div class="row-actions">
                <button class="action-btn btn-sm btn-outline" @click="openModal(row)">
                  编辑
                </button>
                <button class="action-btn btn-sm btn-primary" @click="openPermissionsModal(row)">
                  权限配置
                </button>
                <button
                  v-if="row.code !== 'admin'"
                  class="action-btn btn-sm btn-danger-ghost"
                  @click="handleDelete(row)"
                >
                  删除
                </button>
              </div>
            </template>
          </el-table-column>
        </el-table>
      </div>

      <!-- 添加/编辑角色弹窗 -->
      <el-dialog
        v-model="showModal"
        :title="isEdit ? '编辑角色' : '添加角色'"
        width="500px"
        class="premium-dialog"
        destroy-on-close
      >
        <el-form :model="form" :rules="rules" ref="roleFormRef" label-width="80px" @submit.prevent="handleSubmit">
          <el-form-item label="角色名称" prop="name" required>
            <el-input v-model="form.name" placeholder="请输入角色名称，如：商家管理员" />
          </el-form-item>
          <el-form-item label="角色标识" prop="code" required>
            <el-input v-model="form.code" :disabled="isEdit" placeholder="请输入唯一标识英文，如：store" />
          </el-form-item>
          <el-form-item label="描述/备注" prop="remark">
            <el-input v-model="form.remark" type="textarea" :rows="3" placeholder="请输入备注描述" />
          </el-form-item>
        </el-form>
        <template #footer>
          <div class="dialog-footer" style="display: flex; justify-content: flex-end; gap: 12px;">
            <el-button @click="closeModal">取消</el-button>
            <el-button type="primary" @click="handleSubmit" :loading="submitting">保存</el-button>
          </div>
        </template>
      </el-dialog>

      <!-- 权限配置弹窗 -->
      <el-dialog
        v-model="showPermissionsModal"
        :title="`权限配置 - ${currentRole?.name || ''}`"
        width="800px"
        class="premium-dialog"
        destroy-on-close
      >
        <div class="perm-modal-content" v-loading="loadingPerms">
          <div class="tip-banner">
            <el-icon class="tip-icon"><InfoFilled /></el-icon>
            <span>勾选权限分配后，角色对应的用户将在重新登录后获取对应的功能操作视界。</span>
          </div>

          <div class="perm-groups">
            <div v-for="(group, moduleKey) in groupedPermissions" :key="moduleKey" class="perm-group-card">
              <div class="perm-group-header">
                <h4>{{ getModuleName(moduleKey) }}</h4>
                <el-checkbox
                  :indeterminate="isIndeterminate(moduleKey)"
                  v-model="groupAllChecked[moduleKey]"
                  @change="(val) => handleCheckAllChange(moduleKey, val)"
                >
                  全选
                </el-checkbox>
              </div>
              <div class="perm-group-body">
                <el-checkbox-group v-model="checkedPermissions" @change="handleCheckedChange">
                  <el-checkbox 
                    v-for="p in group" 
                    :key="p" 
                    :value="p"
                    class="perm-checkbox"
                  >
                    <span class="perm-label">{{ getPermissionLabel(p) }}</span>
                  </el-checkbox>
                </el-checkbox-group>
              </div>
            </div>
          </div>
        </div>
        <template #footer>
          <div class="dialog-footer" style="display: flex; justify-content: flex-end; gap: 12px;">
            <el-button @click="showPermissionsModal = false">取消</el-button>
            <el-button type="primary" @click="handleSavePermissions" :loading="savingPerms">保存权限</el-button>
          </div>
        </template>
      </el-dialog>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import Navbar from '@/components/layout/Navbar.vue'
import { ArrowLeft, InfoFilled } from '@element-plus/icons-vue'
import {
  getRoles,
  createRole,
  updateRole,
  deleteRole,
  getAllPermissions,
  getRolePermissions,
  bindRolePermissions
} from '@/api/admin'
import { ElMessage, ElMessageBox } from 'element-plus'

const roles = ref([])
const loading = ref(false)
const submitting = ref(false)

const showModal = ref(false)
const isEdit = ref(false)
const roleFormRef = ref(null)

const form = ref({
  id: 0,
  name: '',
  code: '',
  remark: ''
})

const rules = {
  name: [{ required: true, message: '请输入角色名称', trigger: 'blur' }],
  code: [{ required: true, message: '请输入角色唯一标识', trigger: 'blur' }]
}

// Permissions Configuration
const showPermissionsModal = ref(false)
const loadingPerms = ref(false)
const savingPerms = ref(false)
const currentRole = ref(null)

const allPermissions = ref([])
const checkedPermissions = ref([])
const groupAllChecked = ref({})

// Friendly module name mapping
const moduleNames = {
  'community:fee': '物业费管理',
  'community:notice': '社区公告管理',
  'community:parking': '车位管理',
  'community:visitor': '访客审核管理',
  'mall:category': '商品分类管理',
  'mall:order': '订单履约管理',
  'mall:product': '商品主库库管',
  'mall:store_product': '门店商品上架',
  'mall:store': '合作商户门店',
  'rbac:menu': '系统菜单权限',
  'rbac:permission': '底层权限点',
  'rbac:role': '角色配置中心',
  'rbac:user': '后台账号管理',
  'statistics:ai_report': 'AI 报表诊断',
  'statistics:community': '社区概览统计',
  'statistics:order': '订单营收分析',
  'statistics:product': '商品排行看板',
  'statistics:workorder': '工单效率分析',
  'workorder:complaint': '居民投诉管理',
  'workorder:repair': '报事报修管理'
}

// Friendly permission labels
const permissionLabels = {
  'community:fee:create': '创建物业费账单',
  'community:fee:list': '查询物业费账单',
  'community:fee:payment_list': '查询物业缴费记录',
  'community:notice:create': '发布社区公告',
  'community:notice:delete': '删除社区公告',
  'community:notice:list': '查询公告管理列表',
  'community:notice:views': '查看公告浏览量与状态',
  'community:parking:assign': '分配与绑定车位',
  'community:parking:create': '创建车位资源',
  'community:parking:list': '查询车位列表',
  'community:parking:statistics': '查看车位统计大屏',
  'community:visitor:audit': '审核访客登记与放行',
  'community:visitor:list': '查询访客进出记录',
  'mall:category:create': '创建商品分类',
  'mall:category:delete': '删除商品分类',
  'mall:category:update': '更新商品分类',
  'mall:order:cancel': '管理员订单作废取消',
  'mall:order:list': '查询商城订单',
  'mall:order:ship': '处理发货/配送',
  'mall:product:create': '创建/发布商品',
  'mall:product:delete': '物理删除商品',
  'mall:product:list': '查询商品列表',
  'mall:product:update': '更新编辑商品',
  'mall:store_product:bind': '绑定商品到门店',
  'mall:store_product:list': '查询门店绑定商品',
  'mall:store_product:status': '上下架门店专属商品',
  'mall:store_product:stock': '管理维护门店商品库存',
  'mall:store_product:unbind': '解除门店商品绑定',
  'mall:store:create': '创建全新门店',
  'mall:store:delete': '彻底删除门店',
  'mall:store:list': '查询绑定门店列表',
  'mall:store:list_all': '查看全部系统门店',
  'mall:store:update': '更新门店信息',
  'rbac:menu:list': '查看系统菜单项列表',
  'rbac:permission:list': '查询系统底层权限点',
  'rbac:role:bind_menu': '为角色绑定对应菜单项',
  'rbac:role:bind_permission': '为角色授权绑定权限',
  'rbac:role:create': '创建系统新角色',
  'rbac:role:delete': '物理删除已有角色',
  'rbac:role:get_permissions': '查询单角色已绑权限',
  'rbac:role:list': '查询所有系统角色',
  'rbac:role:update': '更新角色信息与编码',
  'rbac:user:assign_role': '为用户分配单角色',
  'rbac:user:assign_roles': '分配角色',
  'rbac:user:freeze': '冻结或解冻用户账号',
  'rbac:user:get_roles': '查询用户的关联角色',
  'rbac:user:list': '查询后台管理人员列表',
  'rbac:user:update_balance': '人工增减调整用户余额',
  'statistics:ai_report:generate': '手动生成AI分析报告',
  'statistics:ai_report:read': '查看AI生成的分析报告',
  'statistics:community:overview': '查看社区整体运营概览',
  'statistics:order:summary': '查看商城订单汇总分析',
  'statistics:product:sales_rank': '查看商品销量排行统计',
  'statistics:product:view_rank': '查看商品浏览排行分析',
  'statistics:workorder:summary': '查看报修投诉工单分析',
  'workorder:complaint:list': '查询居民投诉列表',
  'workorder:complaint:process': '回复并处理投诉事项',
  'workorder:repair:list': '查询报事报修列表',
  'workorder:repair:process': '处理派单与确认维修'
}

const getModuleName = (key) => {
  return moduleNames[key] || key.toUpperCase() + ' 模块'
}

const getPermissionLabel = (code) => {
  return permissionLabels[code] || '基本操作'
}

// Group permissions by prefix
const groupedPermissions = computed(() => {
  const groups = {}
  allPermissions.value.forEach(p => {
    const parts = p.split(':')
    if (parts.length >= 2) {
      const moduleKey = `${parts[0]}:${parts[1]}`
      if (!groups[moduleKey]) {
        groups[moduleKey] = []
      }
      groups[moduleKey].push(p)
    } else {
      const moduleKey = 'other'
      if (!groups[moduleKey]) {
        groups[moduleKey] = []
      }
      groups[moduleKey].push(p)
    }
  })
  return groups
})

const fetchRoles = async () => {
  loading.value = true
  try {
    const res = await getRoles()
    roles.value = res || []
  } catch (e) {
    console.error(e)
    ElMessage.error('获取角色列表失败')
  } finally {
    loading.value = false
  }
}

const openModal = (role = null) => {
  isEdit.value = !!role
  if (role) {
    form.value = { ...role }
  } else {
    form.value = {
      id: 0,
      name: '',
      code: '',
      remark: ''
    }
  }
  showModal.value = true
}

const closeModal = () => {
  showModal.value = false
}

const handleSubmit = async () => {
  if (!roleFormRef.value) return
  await roleFormRef.value.validate(async (valid) => {
    if (valid) {
      submitting.value = true
      try {
        if (isEdit.value) {
          await updateRole(form.value)
          ElMessage.success('更新角色成功')
        } else {
          await createRole(form.value)
          ElMessage.success('创建角色成功')
        }
        closeModal()
        fetchRoles()
      } catch (e) {
        console.error(e)
        ElMessage.error(e.response?.data?.message || '保存失败')
      } finally {
        submitting.value = false
      }
    }
  })
}

const handleDelete = async (role) => {
  try {
    await ElMessageBox.confirm(
      `确定删除角色「${role.name}」吗？删除后绑定该角色的用户将被降权。`,
      '删除确认',
      {
        confirmButtonText: '确定删除',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    await deleteRole(role.id)
    ElMessage.success('删除角色成功')
    fetchRoles()
  } catch (e) {
    if (e !== 'cancel') {
      console.error(e)
      ElMessage.error('删除角色失败')
    }
  }
}

// Permissions handlers
const openPermissionsModal = async (role) => {
  currentRole.value = role
  showPermissionsModal.value = true
  loadingPerms.value = true
  try {
    // 1. Fetch all permissions in system
    const allRes = await getAllPermissions()
    allPermissions.value = allRes || []

    // 2. Fetch current role bound permissions
    const boundRes = await getRolePermissions(role.id)
    checkedPermissions.value = boundRes || []

    // 3. Initialize indeterminate status
    initGroupCheckStatus()
  } catch (e) {
    console.error(e)
    ElMessage.error('拉取权限配置失败')
  } finally {
    loadingPerms.value = false
  }
}

const initGroupCheckStatus = () => {
  Object.keys(groupedPermissions.value).forEach(moduleKey => {
    const groupPerms = groupedPermissions.value[moduleKey]
    const checkedCount = groupPerms.filter(p => checkedPermissions.value.includes(p)).length
    groupAllChecked.value[moduleKey] = checkedCount === groupPerms.length
  })
}

const isIndeterminate = (moduleKey) => {
  const groupPerms = groupedPermissions.value[moduleKey]
  const checkedCount = groupPerms.filter(p => checkedPermissions.value.includes(p)).length
  return checkedCount > 0 && checkedCount < groupPerms.length
}

const handleCheckAllChange = (moduleKey, val) => {
  const groupPerms = groupedPermissions.value[moduleKey]
  if (val) {
    // Add all of group to checked
    groupPerms.forEach(p => {
      if (!checkedPermissions.value.includes(p)) {
        checkedPermissions.value.push(p)
      }
    })
  } else {
    // Remove all of group from checked
    checkedPermissions.value = checkedPermissions.value.filter(p => !groupPerms.includes(p))
  }
  initGroupCheckStatus()
}

const handleCheckedChange = () => {
  initGroupCheckStatus()
}

const handleSavePermissions = async () => {
  savingPerms.value = true
  try {
    await bindRolePermissions(currentRole.value.id, checkedPermissions.value)
    ElMessage.success('配置角色权限成功')
    showPermissionsModal.value = false
  } catch (e) {
    console.error(e)
    ElMessage.error('保存权限配置失败')
  } finally {
    savingPerms.value = false
  }
}

onMounted(() => {
  fetchRoles()
})
</script>

<style scoped>
.admin-child-page { min-height: 100vh; background-color: #f8f9fa; padding-bottom: 80px; }
.custom-container { max-width: 1280px; margin: 0 auto; }

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

.table-wrapper { background: #ffffff; border-radius: 16px; padding: 24px 32px; box-shadow: 0 4px 20px rgba(0, 0, 0, 0.02); }

:deep(.custom-table) { --el-table-border-color: transparent; border-radius: 8px; overflow: hidden; }
:deep(.custom-table th.el-table__cell) { font-weight: 600; font-size: 14px; padding: 18px 0; border-bottom: 1px solid #ebeef5; background: #fbfcfd; color: #606266; }
:deep(.custom-table td.el-table__cell) { padding: 20px 0; border-bottom: 1px dashed #f0f2f5; font-size: 14px; }
:deep(.custom-table::before) { display: none; }

.row-actions { display: flex; gap: 8px; justify-content: center; }

.custom-tag { border-radius: 12px; font-weight: bold; padding: 2px 10px; }

.action-btn { padding: 10px 24px; border-radius: 20px; font-size: 14px; font-weight: 600; cursor: pointer; transition: all 0.3s; border: 1px solid transparent; display: inline-flex; align-items: center; justify-content: center; }
.btn-sm { padding: 6px 16px; font-size: 13px; }

.btn-primary { background: #2d597b; color: #ffffff; box-shadow: 0 4px 12px rgba(45, 89, 123, 0.2); }
.btn-primary:hover:not(:disabled) { background: #1f435d; transform: translateY(-2px); box-shadow: 0 6px 16px rgba(45, 89, 123, 0.3); }

.btn-outline { background: #ffffff; color: #2d597b; border-color: #2d597b; }
.btn-outline:hover { background: #f0f7ff; transform: translateY(-1px); }

.btn-danger-ghost { background: transparent; color: #f56c6c; border-color: #fbc4c4; }
.btn-danger-ghost:hover { background: #fef0f0; color: #e4393c; transform: translateY(-1px); }

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

/* Permissions configure styles */
.perm-modal-content {
  max-height: 550px;
  overflow-y: auto;
}

.tip-banner {
  background: #fdf6ec;
  color: #e6a23c;
  padding: 12px 16px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 24px;
  font-size: 13px;
  border: 1px solid #f5dab1;
}

.tip-icon {
  font-size: 16px;
}

.perm-groups {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.perm-group-card {
  background: #ffffff;
  border: 1px solid #e4e7ed;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 2px 12px rgba(0,0,0,0.01);
}

.perm-group-header {
  background: #f8fafc;
  padding: 14px 20px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-bottom: 1px solid #ebeef5;
}

.perm-group-header h4 {
  margin: 0;
  font-size: 15px;
  color: #334155;
  font-weight: 600;
}

.perm-group-body {
  padding: 20px;
}

.perm-checkbox {
  width: 48%;
  margin-right: 2%;
  margin-bottom: 16px;
  display: inline-flex;
  align-items: center;
  box-sizing: border-box;
}

.perm-label {
  font-weight: 600;
  color: #1e293b;
  margin-right: 6px;
}

.perm-code {
  font-size: 11px;
  color: #94a3b8;
  background: #f1f5f9;
  padding: 2px 6px;
  border-radius: 4px;
  font-family: monospace;
}
</style>
