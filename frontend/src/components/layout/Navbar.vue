<template>
  <div class="navbar">
    <div class="container">
      <div class="navbar-content">
        <div class="navbar-logo" @click="$router.push('/home')">
          <img :src="logoIcon" alt="智享生活" class="logo-icon-img" />
          智享生活
        </div>

        <div class="navbar-menu">
          <router-link to="/home" class="nav-link">首页</router-link>
          <router-link to="/mall" class="nav-link">商城</router-link>
          <router-link to="/service" class="nav-link">社区服务</router-link>
          <router-link to="/chat" class="nav-link" v-if="userStore.isLoggedIn"
            >AI助手</router-link
          >
          <router-link to="/order" class="nav-link" v-if="userStore.isLoggedIn"
            >我的订单</router-link
          >
        </div>

        <div class="navbar-right">
          <router-link to="/cart" class="cart-icon" v-if="userStore.isLoggedIn">
            <img :src="cartIcon" alt="购物车" class="cart-icon-img" />
            <span class="cart-badge" v-if="cartStore.cartCount > 0">
              {{ cartStore.cartCount }}
            </span>
          </router-link>

          <div v-if="userStore.isLoggedIn" class="user-menu">
            <div
              class="user-avatar"
              @click="showUserMenu = !showUserMenu"
              :title="
                userStore.userInfo.real_name || userStore.userInfo.username
              "
            >
              <img
                :src="
                  userStore.userInfo.avatar ||
                  'https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png'
                "
                style="
                  width: 100%;
                  height: 100%;
                  border-radius: 50%;
                  object-fit: cover;
                "
              />
            </div>
            <div class="user-dropdown" v-if="showUserMenu">
              <router-link to="/profile" class="dropdown-item"
                >个人中心</router-link
              >
              <router-link
                to="/admin"
                class="dropdown-item"
                v-if="
                  ['admin', 'store', 'property'].includes(
                    userStore.userInfo.role,
                  )
                "
              >
                管理后台
              </router-link>
              <div class="dropdown-item" @click="handleLogout">退出登录</div>
            </div>
          </div>

          <div v-else class="auth-buttons">
            <router-link to="/login" class="btn btn-login">登录</router-link>
            <router-link to="/register" class="btn btn-register"
              >注册</router-link
            >
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from "vue";
import { useRouter } from "vue-router";
import { useUserStore } from "@/stores/user";
import { useCartStore } from "@/stores/cart";
import cartIcon from "@/assets/images/购物车 (2).png";
import logoIcon from "@/assets/images/社区管理 (1).png";

const router = useRouter();
const userStore = useUserStore();
const cartStore = useCartStore();
const showUserMenu = ref(false);

const handleLogout = () => {
  userStore.logout();
  showUserMenu.value = false;
  router.push("/home");
};
</script>

<style scoped>
.navbar {
  /* 仿照图片使用深蓝色背景 */
  background: #0d347c;
  position: sticky;
  top: 0;
  z-index: 1000;
}

.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 16px;
}

.navbar-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 64px;
}

.navbar-logo {
  font-size: 20px;
  font-weight: 600;
  color: #ffffff;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 8px;
}

.logo-icon-img {
  width: 26px;
  height: 26px;
  object-fit: contain;
}

.navbar-menu {
  display: flex;
  gap: 32px;
}

.nav-link {
  color: rgba(255, 255, 255, 0.85);
  text-decoration: none;
  font-weight: 500;
  font-size: 16px;
  padding: 8px 0;
  transition: all 0.3s ease;
  position: relative;
}

.nav-link:hover,
.nav-link.router-link-active {
  color: #ffffff;
}

/* 激活状态的底部白条 */
.nav-link.router-link-active::after {
  content: "";
  position: absolute;
  bottom: 0;
  left: 50%;
  transform: translateX(-50%);
  width: 100%;
  height: 3px;
  background: #ffffff;
  border-radius: 2px;
}

.navbar-right {
  display: flex;
  align-items: center;
  gap: 20px;
}

.cart-icon {
  position: relative;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 34px;
  height: 34px;
  cursor: pointer;
  text-decoration: none;
  transition: transform 0.3s;
}

.cart-icon:hover {
  transform: scale(1.05);
}

.cart-icon-img {
  width: 24px;
  height: 24px;
  display: block;
}

.cart-badge {
  position: absolute;
  top: -6px;
  right: -8px;
  background: #ff4d4f;
  color: white;
  font-size: 12px;
  padding: 0 5px;
  border-radius: 10px;
  min-width: 16px;
  text-align: center;
  line-height: 16px;
}

.user-menu {
  position: relative;
}

.user-avatar {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  border: 2px solid rgba(255, 255, 255, 0.3);
  background: #fff;
  cursor: pointer;
}

.user-dropdown {
  position: absolute;
  top: 100%;
  right: 0;
  margin-top: 12px;
  background: white;
  border-radius: 4px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  min-width: 120px;
  overflow: hidden;
}

.dropdown-item {
  display: block;
  padding: 10px 16px;
  color: #333;
  font-size: 14px;
  cursor: pointer;
  text-decoration: none;
  transition: background 0.3s;
}

.dropdown-item:hover {
  background: #f5f7fa;
  color: #0d347c;
}

.auth-buttons {
  display: flex;
  gap: 12px;
}

.btn {
  padding: 6px 16px;
  border-radius: 4px;
  font-size: 14px;
  cursor: pointer;
  text-decoration: none;
  transition: all 0.3s;
}

.btn-login {
  color: #fff;
  border: 1px solid rgba(255, 255, 255, 0.6);
}

.btn-login:hover {
  background: rgba(255, 255, 255, 0.1);
}

.btn-register {
  background: #fff;
  color: #0d347c;
}

.btn-register:hover {
  background: #f0f0f0;
}
</style>
