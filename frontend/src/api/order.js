import request from '@/utils/request'

function appendQueryParams(url, params = {}) {
  const base = url || ''
  const [path, query = ''] = base.split('?')
  const search = new URLSearchParams(query)
  Object.entries(params).forEach(([key, value]) => {
    if (value !== undefined && value !== null && value !== '') {
      search.set(key, String(value))
    }
  })
  const queryString = search.toString()
  return queryString ? `${path}?${queryString}` : path
}

export function addToCart(data) {
  return request({
    url: '/mall/cart/items',
    method: 'post',
    data
  })
}

export function getCartList() {
  return request({
    url: '/mall/cart/items',
    method: 'get'
  }).then((res) => {
    const list = res?.list || res || [];
    return list.map(item => ({
      id: item.id,
      user_id: item.user_id,
      product_id: item.product_id,
      quantity: item.quantity,
      created_at: item.created_at,
      product: {
        id: item.product_id,
        name: item.product_name,
        price: item.product_price,
        image_url: item.product_image
      }
    }));
  })
}

export function deleteCartItem(id) {
  return request({
    url: `/mall/cart/items/${id}`,
    method: 'delete'
  })
}

export function updateCartQuantity(id, quantity) {
  return request({
    url: `/mall/cart/items/${id}`,
    method: 'put',
    data: { quantity }
  })
}

export function createOrder(data) {
  return request({
    url: '/mall/orders',
    method: 'post',
    data
  })
}

export function getAvailableStores(cartIds) {
  return request({
    url: '/mall/orders/available-stores',
    method: 'get',
    params: {
      cart_ids: (cartIds || []).join(',')
    }
  })
}

export function getOrderList(params) {
  return request({
    url: '/mall/orders',
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

export function payOrder(orderId, data = {}) {
  const idempotencyKey = data.idempotency_key || data.idempotencyKey ||
    `order-pay-${orderId}-${Date.now()}-${Math.random().toString(36).slice(2, 10)}`
  const returnUrl = appendQueryParams(data.return_url || data.returnUrl || '', {
    order_id: orderId
  })

  return request({
    url: `/mall/orders/${orderId}/pay`,
    method: 'post',
    data: {
      pay_type: data.pay_type || 'password',
      password: data.password || '',
      face_image_url: data.face_image_url || '',
      idempotency_key: idempotencyKey,
      return_url: returnUrl
    }
  })
}

export function cancelOrder(orderId) {
  return request({
    url: `/mall/orders/${orderId}/cancel`,
    method: 'post'
  })
}

export function receiveOrder(orderId) {
  return request({
    url: `/mall/orders/${orderId}/receive`,
    method: 'post'
  })
}

export function getOrderDetail(id) {
  return request({
    url: `/mall/orders/${id}`,
    method: 'get'
  }).then((res) => {
    if (!res) return res;
    return {
      ...res,
      store: {
        id: res.store_id,
        name: res.store_name,
        address: res.store_address,
        phone: res.store_phone
      },
      items: (res.items || []).map(item => ({
        ...item,
        product: {
          id: item.product_id,
          name: item.product_name,
          image_url: item.product_image
        }
      }))
    };
  })
}
