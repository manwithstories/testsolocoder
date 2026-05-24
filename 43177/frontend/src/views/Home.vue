<template>
  <div class="home-page">
    <el-container>
      <el-header>
        <div class="header-content">
          <div class="logo">
            <h1>维修服务平台</h1>
          </div>
          <div class="nav">
            <router-link to="/home">首页</router-link>
            <router-link to="/technicians">技师</router-link>
            <router-link v-if="userStore.isCustomer" to="/orders">我的工单</router-link>
            <router-link v-if="userStore.isCustomer" to="/orders/create">创建工单</router-link>
            <router-link v-if="userStore.isTechnician" to="/orders">接单中心</router-link>
            <router-link v-if="userStore.isTechnician" to="/parts">配件管理</router-link>
            <router-link v-if="userStore.isLoggedIn" to="/profile">个人中心</router-link>
            <router-link v-if="userStore.isLoggedIn" to="/finance">财务</router-link>
          </div>
          <div class="user-actions">
            <template v-if="userStore.isLoggedIn">
              <el-dropdown>
                <span class="user-info">
                  <el-avatar :size="32" :src="userStore.userInfo?.avatar">
                    {{ userStore.userInfo?.username?.charAt(0) }}
                  </el-avatar>
                  <span class="username">{{ userStore.userInfo?.username }}</span>
                </span>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item @click="router.push('/profile')">个人中心</el-dropdown-item>
                    <el-dropdown-item @click="handleLogout">退出登录</el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
            </template>
            <template v-else>
              <el-button type="primary" @click="router.push('/login')">登录</el-button>
              <el-button @click="router.push('/register')">注册</el-button>
            </template>
          </div>
        </div>
      </el-header>

      <el-main>
        <div class="banner">
          <el-carousel height="300px">
            <el-carousel-item>
              <div class="carousel-item item1">
                <div class="carousel-text">
                  <h2>专业家电维修</h2>
                  <p>快速响应，专业上门服务</p>
                  <el-button type="primary" size="large" @click="router.push('/orders/create')">
                    立即下单
                  </el-button>
                </div>
              </div>
            </el-carousel-item>
            <el-carousel-item>
              <div class="carousel-item item2">
                <div class="carousel-text">
                  <h2>数码设备维修</h2>
                  <p>手机、电脑、平板一站式服务</p>
                  <el-button type="primary" size="large" @click="router.push('/orders/create')">
                    立即下单
                  </el-button>
                </div>
              </div>
            </el-carousel-item>
            <el-carousel-item>
              <div class="carousel-item item3">
                <div class="carousel-text">
                  <h2>汽车维修保养</h2>
                  <p>专业技师，品质保障</p>
                  <el-button type="primary" size="large" @click="router.push('/orders/create')">
                    立即下单
                  </el-button>
                </div>
              </div>
            </el-carousel-item>
          </el-carousel>
        </div>

        <div class="categories">
          <h3>服务分类</h3>
          <div class="category-list">
            <div
              v-for="category in categories"
              :key="category.id"
              class="category-card"
              @click="handleCategoryClick(category)"
            >
              <div class="category-icon">{{ category.icon }}</div>
              <div class="category-name">{{ category.name }}</div>
            </div>
          </div>
        </div>

        <div class="service-items">
          <h3>热门服务</h3>
          <div class="service-grid">
            <el-card v-for="item in serviceItems" :key="item.id" class="service-card" shadow="hover">
              <template #header>
                <div class="service-header">
                  <span class="service-name">{{ item.name }}</span>
                  <span class="service-category">{{ item.category?.name }}</span>
                </div>
              </template>
              <p class="service-desc">{{ item.description }}</p>
              <div class="service-footer">
                <span class="service-price">¥{{ item.min_price }} - ¥{{ item.max_price }}</span>
                <span class="service-time">约{{ item.estimated_time }}分钟</span>
              </div>
              <el-button type="primary" class="service-btn" @click="handleCreateOrder(item)">
                立即预约
              </el-button>
            </el-card>
          </div>
        </div>

        <div class="top-technicians">
          <h3>优秀技师</h3>
          <div class="technician-list">
            <el-card v-for="tech in technicians" :key="tech.id" class="technician-card" shadow="hover">
              <div class="tech-info">
                <el-avatar :size="60" :src="tech.avatar">
                  {{ tech.username?.charAt(0) }}
                </el-avatar>
                <div class="tech-detail">
                  <div class="tech-name">{{ tech.real_name || tech.username }}</div>
                  <div class="tech-specialty">{{ tech.specialty }}</div>
                  <el-rate :model-value="tech.rating" disabled size="small" />
                  <div class="tech-stats">
                    <span>完成{{ tech.completed_orders }}单</span>
                    <span>{{ tech.experience_years }}年经验</span>
                  </div>
                </div>
              </div>
            </el-card>
          </div>
        </div>
      </el-main>

      <el-footer>
        <div class="footer-content">
          <p>© 2024 维修服务平台 版权所有</p>
        </div>
      </el-footer>
    </el-container>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/store/user'
import { categoryApi, serviceItemApi } from '@/api/category'
import { userApi } from '@/api/user'
import type { Category, ServiceItem } from '@/types'

const router = useRouter()
const userStore = useUserStore()

const categories = ref<Category[]>([])
const serviceItems = ref<ServiceItem[]>([])
const technicians = ref<any[]>([])

onMounted(async () => {
  await loadCategories()
  await loadServiceItems()
  await loadTechnicians()
})

async function loadCategories() {
  try {
    const res = await categoryApi.getCategories()
    categories.value = res.data || []
  } catch (error) {
    console.error('Failed to load categories:', error)
  }
}

async function loadServiceItems() {
  try {
    const res = await serviceItemApi.getServiceItems()
    serviceItems.value = res.data || []
  } catch (error) {
    console.error('Failed to load service items:', error)
  }
}

async function loadTechnicians() {
  try {
    const res = await userApi.getTechnicians({ page: 1, page_size: 6 })
    technicians.value = res.data?.list || []
  } catch (error) {
    console.error('Failed to load technicians:', error)
  }
}

function handleCategoryClick(category: Category) {
  router.push({ path: '/orders/create', query: { category_id: category.id } })
}

function handleCreateOrder(item: ServiceItem) {
  if (!userStore.isLoggedIn) {
    ElMessage.warning('请先登录')
    router.push('/login')
    return
  }
  router.push({ path: '/orders/create', query: { service_item_id: item.id } })
}

function handleLogout() {
  userStore.logout()
  ElMessage.success('已退出登录')
  router.push('/home')
}
</script>

<style scoped>
.home-page {
  min-height: 100vh;
  background-color: #f5f7fa;
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  height: 100%;
  max-width: 1200px;
  margin: 0 auto;
}

.logo h1 {
  margin: 0;
  color: #409eff;
  font-size: 24px;
}

.nav {
  display: flex;
  gap: 20px;
}

.nav a {
  color: #606266;
  text-decoration: none;
  font-size: 14px;
}

.nav a:hover {
  color: #409eff;
}

.user-actions {
  display: flex;
  align-items: center;
  gap: 10px;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
}

.username {
  font-size: 14px;
  color: #606266;
}

.banner {
  margin-bottom: 30px;
}

.carousel-item {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
  border-radius: 8px;
  overflow: hidden;
}

.item1 {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.item2 {
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
}

.item3 {
  background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
}

.carousel-text {
  text-align: center;
  color: #fff;
}

.carousel-text h2 {
  font-size: 36px;
  margin-bottom: 10px;
}

.carousel-text p {
  font-size: 18px;
  margin-bottom: 20px;
}

.categories, .service-items, .top-technicians {
  max-width: 1200px;
  margin: 0 auto 40px;
  padding: 0 20px;
}

.categories h3, .service-items h3, .top-technicians h3 {
  margin-bottom: 20px;
  color: #303133;
}

.category-list {
  display: flex;
  gap: 20px;
  flex-wrap: wrap;
}

.category-card {
  flex: 1;
  min-width: 200px;
  padding: 30px;
  background: #fff;
  border-radius: 8px;
  text-align: center;
  cursor: pointer;
  transition: all 0.3s;
}

.category-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.1);
}

.category-icon {
  font-size: 48px;
  margin-bottom: 10px;
}

.category-name {
  font-size: 16px;
  color: #303133;
}

.service-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 20px;
}

.service-card {
  transition: all 0.3s;
}

.service-card:hover {
  transform: translateY(-5px);
}

.service-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.service-name {
  font-weight: 600;
  color: #303133;
}

.service-category {
  font-size: 12px;
  color: #909399;
}

.service-desc {
  color: #606266;
  font-size: 14px;
  margin-bottom: 15px;
  min-height: 40px;
}

.service-footer {
  display: flex;
  justify-content: space-between;
  margin-bottom: 15px;
}

.service-price {
  color: #f56c6c;
  font-weight: 600;
}

.service-time {
  color: #909399;
  font-size: 12px;
}

.service-btn {
  width: 100%;
}

.technician-list {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 20px;
}

.tech-info {
  display: flex;
  gap: 15px;
}

.tech-detail {
  flex: 1;
}

.tech-name {
  font-weight: 600;
  color: #303133;
  margin-bottom: 5px;
}

.tech-specialty {
  color: #909399;
  font-size: 12px;
  margin-bottom: 5px;
}

.tech-stats {
  display: flex;
  gap: 15px;
  font-size: 12px;
  color: #606266;
  margin-top: 10px;
}

.footer-content {
  text-align: center;
  color: #909399;
  font-size: 14px;
  padding: 20px;
}
</style>
