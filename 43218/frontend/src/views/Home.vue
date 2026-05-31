<template>
  <div class="home-page">
    <header class="home-header">
      <div class="container flex-between">
        <router-link to="/" class="logo">
          <el-icon :size="28"><ShoppingCart /></el-icon>
          <span>二手交易平台</span>
        </router-link>
        <el-menu mode="horizontal" :default-active="activeMenu" class="nav-menu" @select="handleMenuSelect">
          <el-menu-item index="home">首页</el-menu-item>
          <el-menu-item index="products">商品市场</el-menu-item>
          <el-menu-item index="services">维修服务</el-menu-item>
        </el-menu>
        <div class="header-right">
          <el-input v-model="searchKeyword" placeholder="搜索商品或服务" clearable style="width: 240px" @keyup.enter="handleSearch">
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
          <template v-if="userStore.isLoggedIn">
            <el-dropdown @command="handleUserCommand">
              <span class="user-info">
                <el-avatar :size="32" :src="userStore.userInfo?.avatar">
                  {{ userStore.userInfo?.nickname?.charAt(0) || 'U' }}
                </el-avatar>
                <span class="username">{{ userStore.userInfo?.nickname || userStore.userInfo?.username }}</span>
                <el-icon><ArrowDown /></el-icon>
              </span>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="profile">个人中心</el-dropdown-item>
                  <el-dropdown-item v-if="userStore.isSeller" command="seller">卖家中心</el-dropdown-item>
                  <el-dropdown-item v-if="userStore.isTechnician" command="technician">技师中心</el-dropdown-item>
                  <el-dropdown-item v-if="userStore.isAdmin" command="admin">管理后台</el-dropdown-item>
                  <el-dropdown-item divided command="logout">退出登录</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </template>
          <template v-else>
            <el-button type="primary" @click="$router.push('/login')">登录</el-button>
            <el-button @click="$router.push('/register')">注册</el-button>
          </template>
        </div>
      </div>
    </header>

    <section class="banner">
      <div class="container">
        <div class="banner-content">
          <h1>二手电子产品交易与维修服务平台</h1>
          <p>买卖二手电子产品，专业维修服务，安全有保障</p>
          <div class="banner-actions">
            <el-button type="primary" size="large" @click="$router.push('/products')">浏览商品</el-button>
            <el-button size="large" @click="$router.push('/services')">维修服务</el-button>
          </div>
        </div>
      </div>
    </section>

    <section class="categories-section">
      <div class="container">
        <h2 class="section-title">商品分类</h2>
        <div class="category-grid">
          <div
            v-for="category in categories"
            :key="category"
            class="category-item"
            @click="goToCategory(category)"
          >
            <el-icon :size="32"><component :is="getCategoryIcon(category)" /></el-icon>
            <span>{{ category }}</span>
          </div>
        </div>
      </div>
    </section>

    <section class="hot-products">
      <div class="container">
        <div class="flex-between">
          <h2 class="section-title">热门商品</h2>
          <el-link type="primary" @click="$router.push('/products')">查看更多 →</el-link>
        </div>
        <div class="product-grid" v-loading="loading">
          <div
            v-for="product in hotProducts"
            :key="product.id"
            class="product-card"
            @click="goToDetail(product.id)"
          >
            <div class="product-image">
              <img :src="getFirstImage(product.images)" :alt="product.title" />
              <span class="condition-tag">{{ product.condition }}</span>
            </div>
            <div class="product-info">
              <h3 class="product-title text-ellipsis">{{ product.title }}</h3>
              <p class="product-brand">{{ product.brand }} {{ product.model }}</p>
              <div class="product-price">
                <span class="price-text">¥{{ product.price.toFixed(2) }}</span>
                <span v-if="product.originalPrice" class="original-price">¥{{ product.originalPrice.toFixed(2) }}</span>
              </div>
              <div class="product-meta">
                <span><el-icon><View /></el-icon> {{ product.viewCount }}</span>
                <span><el-icon><Star /></el-icon> {{ product.favoriteCount }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>

    <section class="services-section">
      <div class="container">
        <div class="flex-between">
          <h2 class="section-title">维修服务</h2>
          <el-link type="primary" @click="$router.push('/services')">查看更多 →</el-link>
        </div>
        <div class="service-grid" v-loading="servicesLoading">
          <div
            v-for="service in hotServices"
            :key="service.id"
            class="service-card"
            @click="$router.push(`/services/${service.id}`)"
          >
            <div class="service-icon">
              <el-icon :size="40"><Tools /></el-icon>
            </div>
            <h3 class="service-title">{{ service.title }}</h3>
            <p class="service-type">{{ service.serviceType }}</p>
            <div class="service-price">
              <span class="price-text">¥{{ service.price.toFixed(2) }}</span>
              <span v-if="service.minPrice" class="price-range">¥{{ service.minPrice }}-¥{{ service.maxPrice }}</span>
            </div>
            <div class="service-meta">
              <span>预计 {{ service.estimatedDays }} 天</span>
              <span><el-icon><Star /></el-icon> {{ service.rating.toFixed(1) }}</span>
            </div>
          </div>
        </div>
      </div>
    </section>

    <footer class="home-footer">
      <div class="container">
        <div class="footer-content">
          <div class="footer-section">
            <h4>关于我们</h4>
            <p>二手电子产品交易与维修服务平台</p>
          </div>
          <div class="footer-section">
            <h4>帮助中心</h4>
            <p>常见问题 | 联系我们</p>
          </div>
          <div class="footer-section">
            <h4>服务条款</h4>
            <p>用户协议 | 隐私政策</p>
          </div>
        </div>
        <div class="footer-bottom">
          <p>© 2024 二手交易平台 版权所有</p>
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/store/user'
import { productApi, repairApi } from '@/api'
import { ProductCategories } from '@/types'
import type { Product, RepairService } from '@/types'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

const searchKeyword = ref('')
const loading = ref(false)
const servicesLoading = ref(false)
const hotProducts = ref<Product[]>([])
const hotServices = ref<RepairService[]>([])

const categories = ProductCategories

const activeMenu = computed(() => {
  if (route.path.startsWith('/products')) return 'products'
  if (route.path.startsWith('/services')) return 'services'
  return 'home'
})

function handleMenuSelect(index: string) {
  switch (index) {
    case 'home':
      router.push('/')
      break
    case 'products':
      router.push('/products')
      break
    case 'services':
      router.push('/services')
      break
  }
}

function handleSearch() {
  if (searchKeyword.value) {
    router.push({ path: '/products', query: { keyword: searchKeyword.value } })
  }
}

function handleUserCommand(command: string) {
  switch (command) {
    case 'profile':
      router.push('/user/profile')
      break
    case 'seller':
      router.push('/seller/products')
      break
    case 'technician':
      router.push('/technician/services')
      break
    case 'admin':
      router.push('/admin/dashboard')
      break
    case 'logout':
      userStore.logout()
      ElMessage.success('已退出登录')
      router.push('/')
      break
  }
}

function getCategoryIcon(category: string): string {
  const iconMap: Record<string, string> = {
    '手机': 'Iphone',
    '电脑': 'Monitor',
    '相机': 'Camera',
    '耳机': 'Headset',
    '平板': 'Tablet',
    '智能手表': 'Watch',
    '游戏机': 'Gamepad',
    '其他数码': 'Box'
  }
  return iconMap[category] || 'Box'
}

function getFirstImage(images: string): string {
  if (!images) return 'https://picsum.photos/300/300'
  try {
    const arr = JSON.parse(images)
    return Array.isArray(arr) && arr.length > 0 ? arr[0] : 'https://picsum.photos/300/300'
  } catch {
    return 'https://picsum.photos/300/300'
  }
}

function goToCategory(category: string) {
  router.push({ path: '/products', query: { category } })
}

function goToDetail(id: number) {
  router.push(`/products/${id}`)
}

async function fetchHotProducts() {
  loading.value = true
  try {
    const res = await productApi.getHot(8)
    hotProducts.value = res.data
  } catch (error) {
    console.error('Failed to fetch hot products:', error)
  } finally {
    loading.value = false
  }
}

async function fetchHotServices() {
  servicesLoading.value = true
  try {
    const res = await repairApi.getServiceList({ page: 1, pageSize: 8, sortBy: 'orders' })
    hotServices.value = res.data
  } catch (error) {
    console.error('Failed to fetch hot services:', error)
  } finally {
    servicesLoading.value = false
  }
}

onMounted(() => {
  fetchHotProducts()
  fetchHotServices()
})
</script>

<style lang="scss" scoped>
.home-page {
  min-height: 100vh;
}

.home-header {
  background: #fff;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
  position: sticky;
  top: 0;
  z-index: 100;

  .logo {
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 18px;
    font-weight: 600;
    color: var(--primary-color);
  }

  .nav-menu {
    border: none;
    flex: 1;
    margin-left: 40px;
  }

  .header-right {
    display: flex;
    align-items: center;
    gap: 16px;
  }

  .user-info {
    display: flex;
    align-items: center;
    gap: 8px;
    cursor: pointer;

    .username {
      max-width: 100px;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
    }
  }
}

.banner {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 60px 0;
  color: #fff;

  .banner-content {
    text-align: center;

    h1 {
      font-size: 36px;
      margin-bottom: 16px;
    }

    p {
      font-size: 18px;
      margin-bottom: 24px;
      opacity: 0.9;
    }

    .banner-actions {
      display: flex;
      gap: 16px;
      justify-content: center;
    }
  }
}

.section-title {
  font-size: 20px;
  font-weight: 600;
  margin-bottom: 20px;
}

.categories-section {
  padding: 40px 0;
  background: #fff;

  .category-grid {
    display: grid;
    grid-template-columns: repeat(8, 1fr);
    gap: 20px;
  }

  .category-item {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 8px;
    padding: 20px;
    border-radius: 8px;
    cursor: pointer;
    transition: all 0.3s;

    &:hover {
      background: #f5f7fa;
      transform: translateY(-2px);
    }
  }
}

.hot-products {
  padding: 40px 0;
  background: #f5f7fa;

  .product-grid {
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    gap: 20px;
  }

  .product-card {
    background: #fff;
    border-radius: 8px;
    overflow: hidden;
    cursor: pointer;
    transition: all 0.3s;

    &:hover {
      transform: translateY(-4px);
      box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
    }

    .product-image {
      position: relative;
      height: 200px;
      overflow: hidden;

      img {
        width: 100%;
        height: 100%;
        object-fit: cover;
      }

      .condition-tag {
        position: absolute;
        top: 10px;
        left: 10px;
        background: rgba(0, 0, 0, 0.6);
        color: #fff;
        padding: 4px 8px;
        border-radius: 4px;
        font-size: 12px;
      }
    }

    .product-info {
      padding: 12px;

      .product-title {
        font-size: 14px;
        margin-bottom: 4px;
      }

      .product-brand {
        color: var(--text-lighter-color);
        font-size: 12px;
        margin-bottom: 8px;
      }

      .product-price {
        display: flex;
        align-items: baseline;
        gap: 8px;
        margin-bottom: 8px;

        .price-text {
          font-size: 18px;
        }

        .original-price {
          font-size: 12px;
          color: var(--text-lighter-color);
          text-decoration: line-through;
        }
      }

      .product-meta {
        display: flex;
        gap: 16px;
        font-size: 12px;
        color: var(--text-lighter-color);

        span {
          display: flex;
          align-items: center;
          gap: 4px;
        }
      }
    }
  }
}

.services-section {
  padding: 40px 0;
  background: #fff;

  .service-grid {
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    gap: 20px;
  }

  .service-card {
    background: #f5f7fa;
    border-radius: 8px;
    padding: 24px;
    cursor: pointer;
    transition: all 0.3s;
    text-align: center;

    &:hover {
      background: #e8f3ff;
      transform: translateY(-4px);
    }

    .service-icon {
      color: var(--primary-color);
      margin-bottom: 12px;
    }

    .service-title {
      font-size: 16px;
      margin-bottom: 8px;
    }

    .service-type {
      font-size: 12px;
      color: var(--text-lighter-color);
      margin-bottom: 12px;
    }

    .service-price {
      margin-bottom: 12px;

      .price-text {
        font-size: 20px;
      }

      .price-range {
        font-size: 12px;
        color: var(--text-lighter-color);
      }
    }

    .service-meta {
      display: flex;
      justify-content: center;
      gap: 16px;
      font-size: 12px;
      color: var(--text-lighter-color);

      span {
        display: flex;
        align-items: center;
        gap: 4px;
      }
    }
  }
}

.home-footer {
  background: #303133;
  color: #fff;
  padding: 40px 0 20px;

  .footer-content {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 40px;
    margin-bottom: 30px;

    h4 {
      margin-bottom: 12px;
    }

    p {
      color: rgba(255, 255, 255, 0.7);
      font-size: 13px;
    }
  }

  .footer-bottom {
    border-top: 1px solid rgba(255, 255, 255, 0.1);
    padding-top: 20px;
    text-align: center;
    color: rgba(255, 255, 255, 0.5);
    font-size: 13px;
  }
}
</style>
