<template>
  <div class="service-list-page">
    <header class="page-header-inner">
      <div class="container flex-between">
        <router-link to="/" class="logo">
          <el-icon :size="24"><ShoppingCart /></el-icon>
          <span>二手交易平台</span>
        </router-link>
        <div class="header-actions">
          <router-link to="/">首页</router-link>
          <router-link to="/products">商品市场</router-link>
          <router-link to="/services">维修服务</router-link>
        </div>
      </div>
    </header>

    <div class="container main-content">
      <div class="search-filter-bar card">
        <div class="filter-row">
          <el-input
            v-model="searchForm.keyword"
            placeholder="搜索维修服务"
            clearable
            style="width: 300px"
            @keyup.enter="handleSearch"
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
        </div>
        <div class="filter-row">
          <span class="filter-label">服务类型：</span>
          <el-radio-group v-model="searchForm.serviceType" @change="handleSearch">
            <el-radio-button value="">全部</el-radio-button>
            <el-radio-button
              v-for="type in serviceTypes"
              :key="type"
              :value="type"
            >
              {{ type }}
            </el-radio-button>
          </el-radio-group>
        </div>
        <div class="filter-row">
          <span class="filter-label">排序：</span>
          <el-select v-model="searchForm.sortBy" @change="handleSearch">
            <el-option label="最新发布" value="newest" />
            <el-option label="价格从低到高" value="price_asc" />
            <el-option label="价格从高到低" value="price_desc" />
            <el-option label="评分优先" value="rating" />
            <el-option label="销量优先" value="orders" />
          </el-select>
        </div>
      </div>

      <div class="service-grid" v-loading="loading">
        <div
          v-for="service in services"
          :key="service.id"
          class="service-card"
          @click="goToDetail(service.id)"
        >
          <div class="service-icon">
            <el-icon :size="48"><Tools /></el-icon>
          </div>
          <h3 class="service-title">{{ service.title }}</h3>
          <p class="service-type">{{ service.serviceType }}</p>
          <p class="service-description text-ellipsis">{{ service.description }}</p>
          <div class="service-price">
            <span class="price-text">¥{{ service.price.toFixed(2) }}</span>
            <span v-if="service.minPrice" class="price-range">
              ¥{{ service.minPrice }}-¥{{ service.maxPrice }}
            </span>
          </div>
          <div class="service-meta">
            <span>预计 {{ service.estimatedDays }} 天</span>
            <span><el-icon><Star /></el-icon> {{ service.rating.toFixed(1) }}</span>
            <span>已接 {{ service.orderCount }} 单</span>
          </div>
          <div class="technician-info">
            <el-avatar :size="24" :src="service.technician?.avatar">
              {{ service.technician?.nickname?.charAt(0) }}
            </el-avatar>
            <span>{{ service.technician?.nickname || service.technician?.username }}</span>
          </div>
        </div>
      </div>

      <div class="pagination-wrapper" v-if="total > 0">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.pageSize"
          :page-sizes="[12, 24, 48]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="fetchServices"
          @current-change="fetchServices"
        />
      </div>

      <div class="empty-state" v-if="!loading && services.length === 0">
        <el-empty description="暂无维修服务" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { repairApi } from '@/api'
import { ServiceTypes } from '@/types'
import type { RepairService } from '@/types'

const router = useRouter()

const loading = ref(false)
const services = ref<RepairService[]>([])
const total = ref(0)
const serviceTypes = ServiceTypes

const searchForm = reactive({
  keyword: '',
  serviceType: '',
  sortBy: 'newest'
})

const pagination = reactive({
  page: 1,
  pageSize: 12
})

function goToDetail(id: number) {
  router.push(`/services/${id}`)
}

async function handleSearch() {
  pagination.page = 1
  await fetchServices()
}

async function fetchServices() {
  loading.value = true
  try {
    const res = await repairApi.getServiceList({
      page: pagination.page,
      pageSize: pagination.pageSize,
      keyword: searchForm.keyword || undefined,
      serviceType: searchForm.serviceType || undefined,
      sortBy: searchForm.sortBy
    })
    services.value = res.data
    total.value = res.pagination.total
  } catch (error) {
    console.error('Failed to fetch services:', error)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchServices()
})
</script>

<style lang="scss" scoped>
.service-list-page {
  min-height: 100vh;
  background: #f5f7fa;
}

.page-header-inner {
  background: #fff;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);

  .logo {
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 16px;
    font-weight: 600;
    color: var(--primary-color);
  }

  .header-actions {
    display: flex;
    gap: 24px;

    a {
      color: var(--text-color);
      padding: 20px 0;

      &.router-link-active {
        color: var(--primary-color);
      }
    }
  }
}

.main-content {
  padding: 20px;
}

.search-filter-bar {
  margin-bottom: 20px;

  .filter-row {
    display: flex;
    align-items: center;
    gap: 12px;
    margin-bottom: 16px;

    &:last-child {
      margin-bottom: 0;
    }

    .filter-label {
      color: var(--text-light-color);
      min-width: 80px;
    }
  }
}

.service-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 20px;
  margin-bottom: 20px;
}

.service-card {
  background: #fff;
  border-radius: 8px;
  padding: 24px;
  cursor: pointer;
  transition: all 0.3s;

  &:hover {
    transform: translateY(-4px);
    box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
  }

  .service-icon {
    color: var(--primary-color);
    text-align: center;
    margin-bottom: 16px;
  }

  .service-title {
    font-size: 18px;
    text-align: center;
    margin-bottom: 8px;
  }

  .service-type {
    text-align: center;
    color: var(--primary-color);
    font-size: 13px;
    margin-bottom: 8px;
  }

  .service-description {
    font-size: 13px;
    color: var(--text-lighter-color);
    margin-bottom: 16px;
    height: 36px;
  }

  .service-price {
    text-align: center;
    margin-bottom: 16px;

    .price-text {
      font-size: 24px;
    }

    .price-range {
      display: block;
      font-size: 12px;
      color: var(--text-lighter-color);
    }
  }

  .service-meta {
    display: flex;
    justify-content: center;
    gap: 12px;
    font-size: 12px;
    color: var(--text-lighter-color);
    margin-bottom: 16px;

    span {
      display: flex;
      align-items: center;
      gap: 4px;
    }
  }

  .technician-info {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 8px;
    font-size: 13px;
    padding-top: 12px;
    border-top: 1px solid #f0f0f0;
  }
}

.pagination-wrapper {
  display: flex;
  justify-content: center;
  padding: 20px 0;
}
</style>
