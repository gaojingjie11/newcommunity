USE smart_community;

-- Roles: admin, property, store, user
INSERT INTO `sys_role` (`id`, `name`, `code`, `remark`, `created_at`) VALUES
(1, '系统管理员', 'admin', '全局管理权限', NOW()),
(2, '物业管理员', 'property', '物业与报修管理', NOW()),
(3, '门店管理员', 'store', '门店与商品管理', NOW()),
(4, '普通用户', 'user', '居民用户', NOW())
ON DUPLICATE KEY UPDATE `name`=VALUES(`name`), `remark`=VALUES(`remark`);

-- Menus
INSERT INTO `sys_menu` (`id`, `parent_id`, `name`, `path`, `component`, `sort`, `type`, `created_at`) VALUES
(1, 0, '仪表盘', '/admin/dashboard', 'views/admin/Dashboard.vue', 1, 1, NOW()),
(2, 0, '订单管理', '/admin/order', 'views/admin/order/OrderList.vue', 2, 1, NOW()),
(3, 0, '物业管理', '/admin/property', 'views/admin/property/PropertyFeeList.vue', 3, 1, NOW()),
(4, 0, 'AI报表', '/admin/ai-report', 'views/admin/AIReport.vue', 4, 1, NOW()),
(5, 0, '权限管理', '/admin/rbac', 'views/admin/rbac/RbacLayout.vue', 5, 1, NOW()),
(6, 5, '角色管理', '/admin/rbac/roles', 'views/admin/rbac/RoleList.vue', 1, 2, NOW()),
(7, 5, '用户管理', '/admin/rbac/users', 'views/admin/rbac/UserList.vue', 2, 2, NOW()),
(8, 5, '会员管理', '/admin/rbac/members', 'views/admin/rbac/MemberList.vue', 3, 2, NOW()),
(9, 0, '商品管理', '/admin/mall', 'views/admin/mall/MallLayout.vue', 6, 1, NOW()),
(10, 9, '商品列表', '/admin/mall/products', 'views/admin/mall/ProductList.vue', 1, 2, NOW()),
(11, 9, '分类管理', '/admin/mall/categories', 'views/admin/mall/CategoryList.vue', 2, 2, NOW()),
(12, 9, '营销管理', '/admin/mall/promotions', 'views/admin/mall/PromotionList.vue', 3, 2, NOW()),
(13, 0, '门店管理', '/admin/store', 'views/admin/store/StoreLayout.vue', 7, 1, NOW()),
(14, 13, '门店列表', '/admin/store/list', 'views/admin/store/StoreList.vue', 1, 2, NOW()),
(15, 13, '服务区域', '/admin/store/areas', 'views/admin/store/ServiceAreaList.vue', 2, 2, NOW()),
(16, 0, '社区管理', '/admin/community', 'views/admin/community/CommunityLayout.vue', 8, 1, NOW()),
(17, 16, '公告管理', '/admin/community/notices', 'views/admin/community/NoticeList.vue', 1, 2, NOW()),
(18, 16, '访客管理', '/admin/community/visitors', 'views/admin/community/VisitorList.vue', 2, 2, NOW()),
(19, 16, '车位管理', '/admin/community/parking', 'views/admin/community/ParkingList.vue', 3, 2, NOW()),
(20, 0, '工单管理', '/admin/workorder', 'views/admin/workorder/WorkorderLayout.vue', 9, 1, NOW()),
(21, 20, '报修投诉', '/admin/repairs', 'views/admin/RepairList.vue', 1, 2, NOW()),
(23, 0, '系统管理', '/admin/system', 'views/admin/system/SystemLayout.vue', 10, 1, NOW()),
(24, 23, '用户登录日志', '/admin/system/user-login-logs', 'views/admin/system/UserLoginLogs.vue', 1, 2, NOW()),
(25, 23, '管理员登录日志', '/admin/system/admin-login-logs', 'views/admin/system/AdminLoginLogs.vue', 2, 2, NOW())
ON DUPLICATE KEY UPDATE `name`=VALUES(`name`), `path`=VALUES(`path`), `component`=VALUES(`component`), `sort`=VALUES(`sort`);

-- Permissions (API access control points)
INSERT INTO `sys_permission` (`id`, `code`, `name`, `resource`, `method`, `path`, `type`, `status`, `created_at`, `updated_at`) VALUES
(1,  'rbac:role:create',          '创建角色',       'role',       'POST',   '/api/admin/roles',              2, 1, NOW(), NOW()),
(2,  'rbac:role:update',          '更新角色',       'role',       'PUT',    '/api/admin/roles',              2, 1, NOW(), NOW()),
(3,  'rbac:role:delete',          '删除角色',       'role',       'DELETE', '/api/admin/roles',              2, 1, NOW(), NOW()),
(4,  'rbac:role:list',            '查询角色列表',    'role',       'GET',    '/api/admin/roles',              1, 1, NOW(), NOW()),
(5,  'rbac:role:bind_menu',       '角色绑定菜单',    'role',       'POST',   '/api/admin/roles/:id/menus',    2, 1, NOW(), NOW()),
(6,  'rbac:role:bind_permission', '角色绑定权限',    'role',       'POST',   '/api/admin/roles/:id/permissions', 2, 1, NOW(), NOW()),
(7,  'rbac:role:get_permissions', '查询角色权限',    'role',       'GET',    '/api/admin/roles/:id/permissions', 1, 1, NOW(), NOW()),
(8,  'rbac:user:list',            '查询管理员列表',  'user',       'GET',    '/api/admin/users',              1, 1, NOW(), NOW()),
(9,  'rbac:user:freeze',          '冻结/解冻用户',  'user',       'POST',   '/api/admin/users/freeze',       2, 1, NOW(), NOW()),
(10, 'rbac:user:assign_role',     '分配用户角色(旧)', 'user',      'POST',   '/api/admin/users/assign-role',  2, 1, NOW(), NOW()),
(11, 'rbac:user:assign_roles',    '分配用户多角色',  'user',       'POST',   '/api/admin/users/:id/roles',    2, 1, NOW(), NOW()),
(12, 'rbac:user:get_roles',       '查询用户角色',    'user',       'GET',    '/api/admin/users/:id/roles',    1, 1, NOW(), NOW()),
(13, 'rbac:member:list',          '查询会员列表',    'member',     'GET',    '/api/admin/members',            1, 1, NOW(), NOW()),
(14, 'rbac:permission:list',      '查询权限列表',    'permission', 'GET',    '/api/admin/permissions',        1, 1, NOW(), NOW()),
(15, 'rbac:menu:list',            '查询菜单列表',    'menu',       'GET',    '/api/admin/menus',              1, 1, NOW(), NOW()),
(16, 'log:user_login:list',       '查询用户登录日志', 'log',       'GET',    '/api/admin/user-login-logs',    1, 1, NOW(), NOW()),
(17, 'log:admin_login:list',      '查询管理员登录日志', 'log',     'GET',    '/api/admin/admin-login-logs',   1, 1, NOW(), NOW()),
-- Mall service permissions
(19, 'mall:product:list',         '查询商品列表',     'mall_product',     'GET',    '/api/admin/mall/products',             1, 1, NOW(), NOW()),
(20, 'mall:category:create',      '创建商品分类',     'mall_category',    'POST',   '/api/admin/mall/categories',           2, 1, NOW(), NOW()),
(21, 'mall:category:update',      '更新商品分类',     'mall_category',    'PUT',    '/api/admin/mall/categories/:id',       2, 1, NOW(), NOW()),
(22, 'mall:category:delete',      '删除商品分类',     'mall_category',    'DELETE', '/api/admin/mall/categories/:id',       2, 1, NOW(), NOW()),
(23, 'mall:product:create',       '创建商品',        'mall_product',     'POST',   '/api/admin/mall/products',             2, 1, NOW(), NOW()),
(24, 'mall:product:update',       '更新商品',        'mall_product',     'PUT',    '/api/admin/mall/products/:id',         2, 1, NOW(), NOW()),
(25, 'mall:product:delete',       '删除商品',        'mall_product',     'DELETE', '/api/admin/mall/products/:id',         2, 1, NOW(), NOW()),
(26, 'mall:promotion:create',     '创建促销',        'mall_promotion',   'POST',   '/api/admin/mall/promotions',           2, 1, NOW(), NOW()),
(27, 'mall:promotion:update',     '更新促销',        'mall_promotion',   'PUT',    '/api/admin/mall/promotions/:id',       2, 1, NOW(), NOW()),
(28, 'mall:promotion:delete',     '删除促销',        'mall_promotion',   'DELETE', '/api/admin/mall/promotions/:id',       2, 1, NOW(), NOW()),
(29, 'mall:promotion:bind_product','促销绑定商品',    'mall_promotion',   'POST',   '/api/admin/mall/promotions/:id/products', 2, 1, NOW(), NOW()),
(30, 'mall:service_area:create',  '创建服务区域',     'mall_service_area','POST',   '/api/admin/mall/service-areas',        2, 1, NOW(), NOW()),
(31, 'mall:service_area:update',  '更新服务区域',     'mall_service_area','PUT',    '/api/admin/mall/service-areas/:id',    2, 1, NOW(), NOW()),
(32, 'mall:service_area:delete',  '删除服务区域',     'mall_service_area','DELETE', '/api/admin/mall/service-areas/:id',    2, 1, NOW(), NOW()),
(33, 'mall:store:create',         '创建门店',        'mall_store',       'POST',   '/api/admin/mall/stores',               2, 1, NOW(), NOW()),
(34, 'mall:store:update',         '更新门店',        'mall_store',       'PUT',    '/api/admin/mall/stores/:id',           2, 1, NOW(), NOW()),
(35, 'mall:store:delete',         '删除门店',        'mall_store',       'DELETE', '/api/admin/mall/stores/:id',           2, 1, NOW(), NOW()),
(36, 'mall:store_product:bind',   '绑定门店商品',     'mall_store_product','POST',  '/api/admin/mall/store-products',       2, 1, NOW(), NOW()),
(37, 'mall:store_product:unbind', '解绑门店商品',     'mall_store_product','DELETE','/api/admin/mall/store-products',       2, 1, NOW(), NOW()),
(38, 'mall:store_product:status', '上下架门店商品',   'mall_store_product','PUT',   '/api/admin/mall/store-products/status',2, 1, NOW(), NOW()),
(39, 'mall:store_product:stock',  '门店商品库存',     'mall_store_product','PUT',   '/api/admin/mall/store-products/stock', 2, 1, NOW(), NOW()),
(40, 'mall:store_product:list',   '查询门店商品',     'mall_store_product','GET',   '/api/admin/mall/store-products/:store_id', 1, 1, NOW(), NOW()),
(41, 'mall:order:list',           '查询订单',        'mall_order',       'GET',    '/api/admin/mall/orders',               1, 1, NOW(), NOW()),
(42, 'mall:order:ship',           '订单发货',        'mall_order',       'POST',   '/api/admin/mall/orders/:id/ship',      2, 1, NOW(), NOW()),
(43, 'mall:order:cancel',         '订单作废',        'mall_order',       'POST',   '/api/admin/mall/orders/:id/cancel',    2, 1, NOW(), NOW()),
-- Community service permissions
(50, 'community:notice:list',          '查询公告管理列表', 'community_notice',  'GET',    '/api/admin/community/notices',                  1, 1, NOW(), NOW()),
(51, 'community:notice:create',        '发布公告',        'community_notice',  'POST',   '/api/admin/community/notices',                  2, 1, NOW(), NOW()),
(52, 'community:notice:delete',        '删除公告',        'community_notice',  'DELETE', '/api/admin/community/notices/:id',              2, 1, NOW(), NOW()),
(53, 'community:notice:views',         '查询公告浏览状态', 'community_notice',  'GET',    '/api/admin/community/notices/:id/views',        1, 1, NOW(), NOW()),
(54, 'community:visitor:list',         '查询访客记录',     'community_visitor', 'GET',    '/api/admin/community/visitors',                 1, 1, NOW(), NOW()),
(55, 'community:visitor:audit',        '访客审核放行',     'community_visitor', 'POST',   '/api/admin/community/visitors/:id/audit',       2, 1, NOW(), NOW()),
(56, 'community:parking:list',         '查询车位列表',     'community_parking', 'GET',    '/api/admin/community/parking-spaces',           1, 1, NOW(), NOW()),
(57, 'community:parking:create',       '创建车位',        'community_parking', 'POST',   '/api/admin/community/parking-spaces',           2, 1, NOW(), NOW()),
(58, 'community:parking:assign',       '分配车位',        'community_parking', 'POST',   '/api/admin/community/parking-spaces/:id/assign',2, 1, NOW(), NOW()),
(59, 'community:parking:statistics',   '查询车位统计',     'community_parking', 'GET',    '/api/admin/community/parking-spaces/statistics',1, 1, NOW(), NOW()),
(60, 'community:fee:list',             '查询物业费',       'community_fee',     'GET',    '/api/admin/community/property-fees',            1, 1, NOW(), NOW()),
(61, 'community:fee:create',           '创建物业费',       'community_fee',     'POST',   '/api/admin/community/property-fees',            2, 1, NOW(), NOW()),
(62, 'community:fee:payment_list',     '查询缴费记录',     'community_fee',     'GET',    '/api/admin/community/property-fees/payments',   1, 1, NOW(), NOW()),
-- Workorder service permissions
(70, 'workorder:repair:list',          '查询报修投诉列表', 'workorder',           'GET',  '/api/admin/workorders',                        1, 1, NOW(), NOW()),
(71, 'workorder:repair:process',       '处理报修投诉',     'workorder',           'POST', '/api/admin/workorders/:id/process',            2, 1, NOW(), NOW()),
-- Statistics service permissions
(80, 'statistics:product:sales_rank',  '商品销售排行',     'statistics_product',  'GET',  '/api/statistics/products/sales-rank',           1, 1, NOW(), NOW()),
(81, 'statistics:product:view_rank',   '商品访客排行',     'statistics_product',  'GET',  '/api/statistics/products/view-rank',            1, 1, NOW(), NOW()),
(82, 'statistics:community:overview',  '社区运营概览',     'statistics_community','GET',  '/api/statistics/community/overview',            1, 1, NOW(), NOW()),
(83, 'statistics:order:summary',       '订单统计',        'statistics_order',    'GET',  '/api/statistics/orders',                        1, 1, NOW(), NOW()),
(84, 'statistics:workorder:summary',   '报修投诉统计',     'statistics_workorder','GET',  '/api/statistics/workorders',                    1, 1, NOW(), NOW())
ON DUPLICATE KEY UPDATE `name`=VALUES(`name`), `resource`=VALUES(`resource`), `method`=VALUES(`method`), `path`=VALUES(`path`), `status`=VALUES(`status`), `updated_at`=NOW();

-- Admin role gets all permissions
INSERT INTO `sys_role_permission` (`id`, `role_id`, `permission_id`) VALUES
(1,  1, 1),  (2,  1, 2),  (3,  1, 3),  (4,  1, 4),  (5,  1, 5),
(6,  1, 6),  (7,  1, 7),  (8,  1, 8),  (9,  1, 9),  (10, 1, 10),
(11, 1, 11), (12, 1, 12), (13, 1, 13), (14, 1, 14), (15, 1, 15),
(16, 1, 16), (17, 1, 17), (39, 1, 19),
-- Admin: mall permissions (20~43)
(40, 1, 20), (41, 1, 21), (42, 1, 22), (43, 1, 23), (44, 1, 24),
(45, 1, 25), (46, 1, 26), (47, 1, 27), (48, 1, 28), (49, 1, 29),
(50, 1, 30), (51, 1, 31), (52, 1, 32), (53, 1, 33), (54, 1, 34),
(55, 1, 35), (56, 1, 36), (57, 1, 37), (58, 1, 38), (59, 1, 39),
(60, 1, 40), (61, 1, 41), (62, 1, 42), (63, 1, 43),
-- Admin: community permissions (50~62)
(100, 1, 50), (101, 1, 51), (102, 1, 52), (103, 1, 53), (104, 1, 54),
(105, 1, 55), (106, 1, 56), (107, 1, 57), (108, 1, 58), (109, 1, 59),
(110, 1, 60), (111, 1, 61), (112, 1, 62),
-- Admin: workorder permissions
(140, 1, 70), (141, 1, 71),
-- Admin: statistics permissions (80~84)
(160, 1, 80), (161, 1, 81), (162, 1, 82), (163, 1, 83), (164, 1, 84),
-- Property role: community + workorder permissions
(20, 2, 4),  (21, 2, 13), (22, 2, 16), (23, 2, 17),
(120, 2, 50), (121, 2, 51), (122, 2, 52), (123, 2, 53), (124, 2, 54),
(125, 2, 55), (126, 2, 56), (127, 2, 57), (128, 2, 58), (129, 2, 59),
(130, 2, 60), (131, 2, 61), (132, 2, 62),
(150, 2, 70), (151, 2, 71),
-- Store role: mall management permissions (20~43)
(69, 3, 19), (70, 3, 20), (71, 3, 21), (72, 3, 22), (73, 3, 23), (74, 3, 24),
(75, 3, 25), (76, 3, 26), (77, 3, 27), (78, 3, 28), (79, 3, 29),
(80, 3, 30), (81, 3, 31), (82, 3, 32), (83, 3, 33), (84, 3, 34),
(85, 3, 35), (86, 3, 36), (87, 3, 37), (88, 3, 38), (89, 3, 39),
(90, 3, 40), (91, 3, 41), (92, 3, 42), (93, 3, 43)
ON DUPLICATE KEY UPDATE `role_id`=VALUES(`role_id`), `permission_id`=VALUES(`permission_id`);

-- Role-menu bindings: admin sees all
INSERT INTO `sys_role_menu` (`id`, `role_id`, `menu_id`) VALUES
(1, 1, 1), (2, 1, 2), (3, 1, 3), (4, 1, 4), (5, 1, 5), (6, 1, 6), (7, 1, 7),
(8, 1, 8), (9, 1, 9), (10, 1, 10), (11, 1, 11), (12, 1, 12), (13, 1, 13),
(14, 1, 14), (15, 1, 15), (16, 1, 16), (17, 1, 17), (18, 1, 18), (19, 1, 19),
(20, 1, 20), (21, 1, 21), (22, 1, 22), (23, 1, 23), (24, 1, 24), (25, 1, 25),
-- Property sees: dashboard, property, community, workorder, system logs
(30, 2, 1), (31, 2, 3), (32, 2, 16), (33, 2, 17), (34, 2, 18), (35, 2, 19),
(36, 2, 20), (37, 2, 21), (39, 2, 23), (40, 2, 24), (41, 2, 25),
-- Store sees: dashboard, order, mall, store
(50, 3, 1), (51, 3, 2), (52, 3, 9), (53, 3, 10), (54, 3, 11), (55, 3, 12),
(56, 3, 13), (57, 3, 14), (58, 3, 15)
ON DUPLICATE KEY UPDATE `role_id`=VALUES(`role_id`), `menu_id`=VALUES(`menu_id`);

-- Default admin user (password: 123456, bcrypt cost=14)
INSERT INTO `sys_user` (`id`, `username`, `password`, `real_name`, `mobile`, `age`, `gender`, `email`, `avatar`, `green_points`, `balance`, `role`, `status`, `created_at`, `updated_at`) VALUES
(1, 'admin', '$2a$14$5Kc3aGL3vFnv3LAUDOCEZOTKQWfqag5edtPgvsKuS0C2qePaD4i46', '系统管理员', '13800000001', 30, 1, 'admin@community.com', 'https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png', 5000, 8000.00, 'admin', 1, NOW(), NOW())
ON DUPLICATE KEY UPDATE `username`=VALUES(`username`), `password`=VALUES(`password`);

-- Admin user-role binding
INSERT INTO `sys_user_role` (`id`, `user_id`, `role_id`) VALUES
(1, 1, 1)
ON DUPLICATE KEY UPDATE `user_id`=VALUES(`user_id`), `role_id`=VALUES(`role_id`);
