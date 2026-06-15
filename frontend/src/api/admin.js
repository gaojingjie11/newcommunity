import request from '@/utils/request'

export function getDashboardStats() {
  return request({
    url: '/statistics/community/overview',
    method: 'get'
  })
}

export function getAIReport() {
  return request({
    url: '/statistics/ai-report/latest',
    method: 'get'
  })
}

export function generateAIReport() {
  return request({
    url: '/statistics/ai-report/generate',
    method: 'post',
    timeout: 120000
  })
}

export function getAIReportList(params) {
  return request({
    url: '/statistics/ai-report/list',
    method: 'get',
    params
  })
}

export function getAIReportDetail(id) {
  return request({
    url: `/statistics/ai-report/${id}`,
    method: 'get'
  })
}

export function getUserList(params) {
  return request({
    url: '/admin/users',
    method: 'get',
    params
  })
}

export function freezeUser(data) {
  return request({
    url: '/admin/users/freeze',
    method: 'post',
    data
  })
}


export function getMenuList() {
  return request({
    url: '/admin/menus',
    method: 'get'
  })
}

export function getAdminVisitorList(params) {
  return request({
    url: '/admin/community/visitors',
    method: 'get',
    params
  })
}

export function auditVisitor(visitorId, data) {
  return request({
    url: `/admin/community/visitors/${visitorId}/audit`,
    method: 'post',
    data
  })
}

export function getAdminWorkorderList(params) {
  return request({
    url: '/admin/workorders',
    method: 'get',
    params
  })
}

export function processWorkorder(workorderId, data) {
  return request({
    url: `/admin/workorders/${workorderId}/process`,
    method: 'post',
    data
  })
}

export function getAdminProductList(params) {
  return request({
    url: '/admin/mall/products',
    method: 'get',
    params
  })
}

export function createProduct(data) {
  return request({
    url: '/admin/mall/products',
    method: 'post',
    data
  })
}

export function updateProduct(data) {
  return request({
    url: `/admin/mall/products/${data.id}`,
    method: 'put',
    data
  })
}

export function deleteProduct(id) {
  return request({
    url: `/admin/mall/products/${id}`,
    method: 'delete'
  })
}

export function getAdminOrderList(params) {
  return request({
    url: '/admin/mall/orders',
    method: 'get',
    params
  }).then((res) => {
    const list = res?.list || res || [];
    const mapItems = (items) => (items || []).map(item => ({
      ...item,
      product: {
        id: item.product_id,
        name: item.product_name,
        image_url: item.product_image
      }
    }));
    
    if (res && res.list) {
      return {
        ...res,
        list: list.map(order => ({
          ...order,
          items: mapItems(order.items)
        }))
      };
    }
    
    return list.map(order => ({
      ...order,
      items: mapItems(order.items)
    }));
  })
}

export function shipOrder(orderId, data) {
  return request({
    url: `/admin/mall/orders/${orderId}/ship`,
    method: 'post',
    data
  })
}

export function createStore(data) {
  return request({
    url: '/admin/mall/stores',
    method: 'post',
    data
  })
}

export function updateStore(data) {
  return request({
    url: `/admin/mall/stores/${data.id}`,
    method: 'put',
    data
  })
}

export function deleteStore(id) {
  return request({
    url: `/admin/mall/stores/${id}`,
    method: 'delete'
  })
}

export function createNotice(data) {
  return request({
    url: '/admin/community/notices',
    method: 'post',
    data
  })
}

export function deleteNotice(id) {
  return request({
    url: `/admin/community/notices/${id}`,
    method: 'delete'
  })
}

export function assignRole(data) {
  return request({
    url: '/admin/users/assign-role',
    method: 'post',
    data
  })
}

export function getUserRoles(userId) {
  return request({
    url: '/admin/users/roles',
    method: 'get',
    params: { user_id: userId }
  })
}

export function updateUserBalance(data) {
  return request({
    url: '/admin/users/update-balance',
    method: 'post',
    data
  })
}

export function getAdminParkingList(params) {
  return request({
    url: '/admin/community/parking-spaces',
    method: 'get',
    params
  })
}

export function getParkingStats() {
  return request({
    url: '/admin/community/parking-spaces/statistics',
    method: 'get'
  })
}

export function assignParking(parkingId, data) {
  return request({
    url: `/admin/community/parking-spaces/${parkingId}/assign`,
    method: 'post',
    data
  })
}

export function createParking(data) {
  return request({
    url: '/admin/community/parking-spaces',
    method: 'post',
    data
  })
}

export function createPropertyFee(data) {
  return request({
    url: '/admin/community/property-fees',
    method: 'post',
    data
  })
}

export function getAdminPropertyFeeList(params) {
  return request({
    url: '/admin/community/property-fees',
    method: 'get',
    params
  })
}

export function getStoreProducts(storeId) {
  return request({
    url: `/admin/mall/store-products/${storeId}`,
    method: 'get'
  })
}

export function bindStoreProduct(data) {
  return request({
    url: '/admin/mall/store-products',
    method: 'post',
    data
  })
}

export function unbindStoreProduct(data) {
  return request({
    url: '/admin/mall/store-products',
    method: 'delete',
    data
  })
}

export function updateStoreProductStatus(data) {
  return request({
    url: '/admin/mall/store-products/status',
    method: 'put',
    data
  })
}

export function updateStoreProductStock(data) {
  return request({
    url: '/admin/mall/store-products/stock',
    method: 'put',
    data
  })
}

export function getAdminStores(params) {
  return request({
    url: '/admin/mall/stores',
    method: 'get',
    params
  })
}

export function getUserStores(userId) {
  return request({
    url: `/admin/users/${userId}/stores`,
    method: 'get'
  })
}

export function getRoles() {
  return request({
    url: '/admin/roles',
    method: 'get'
  })
}

export function createRole(data) {
  return request({
    url: '/admin/roles',
    method: 'post',
    data
  })
}

export function updateRole(data) {
  return request({
    url: '/admin/roles',
    method: 'put',
    data
  })
}

export function deleteRole(id) {
  return request({
    url: `/admin/roles?id=${id}`,
    method: 'delete'
  })
}

export function getRolePermissions(roleId) {
  return request({
    url: `/admin/roles/${roleId}/permissions`,
    method: 'get'
  })
}

export function bindRolePermissions(roleId, permissions) {
  return request({
    url: `/admin/roles/${roleId}/permissions`,
    method: 'post',
    data: { permissions }
  })
}

export function getAllPermissions() {
  return request({
    url: '/admin/permissions',
    method: 'get'
  })
}
