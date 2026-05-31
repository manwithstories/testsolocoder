<template>
  <div class="product-list-page">
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
            placeholder="搜索商品名称、品牌、型号"
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
          <span class="filter-label">分类：</span>
          <el-radio-group v-model="searchForm.category" @change="handleSearch">
            <el-radio-button value="">全部</el-radio-button>
            <el-radio-button
              v-for="cat in categories"
              :key="cat"
              :value="cat"
            >
              {{ cat }}
            </el-radio-button>
          </el-radio-group>
        </div>
        <div class="filter-row">
          <span class="filter-label">成色：</span>
          <el-radio-group v-model="searchForm.condition" @change="handleSearch">
            <el-radio-button value="">全部</el-radio-button>
            <el-radio-button
              v-for="cond in conditions"
              :key="cond"
              :value="cond"
            >
              {{ cond }}
            </el-radio-button>
          </el-radio-group>
        </div>
        <div class="filter-row">
          <span class="filter-label">价格：</span>
          <el-input-number
            v-model="searchForm.minPrice"
            :min="0"
            placeholder="最低"
            @change="handleSearch"
          />
          <span class="price-separator">-</span>
          <el-input-number
            v-model="searchForm.maxPrice"
            :min="0"
            placeholder="最高"
            @change="handleSearch"
          />
          <el-select v-model="searchForm.sortBy" placeholder="排序" @change="handleSearch" style="margin-left: 20px">
            <el-option label="最新发布" value="newest" />
            <el-option label="价格从低到高" value="price_asc" />
            <el-option label="价格从高到低" value="price_desc" />
            <el-option label="销量优先" value="sold" />
          </el-select>
        </div>
      </div>

      <div class="product-grid" v-loading="loading">
        <div
          v-for="product in products"
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

      <div class="pagination-wrapper" v-if="total > 0">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.pageSize"
          :page-sizes="[12, 24, 48]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="fetchProducts"
          @current-change="fetchProducts"
        />
      </div>

      <div class="empty-state" v-if="!loading && products.length === 0">
        <el-empty description="暂无商品" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { productApi } from '@/api'
import { ProductCategories, ProductConditions } from '@/types'
import type { Product } from '@/types'

const route = useRoute()
const router = useRouter()

const loading = ref(false)
const products = ref<Product[]>([])
const total = ref(0)
const categories = ProductCategories
const conditions = ProductConditions

const searchForm = reactive({
  keyword: route.query.keyword as string || '',
  category: route.query.category as string || '',
  condition: '',
  minPrice: undefined as number | undefined,
  maxPrice: undefined as number | undefined,
  sortBy: 'newest'
})

const pagination = reactive({
  page: 1,
  pageSize: 12
})

function getFirstImage(images: string): string {
  if (!images) return 'https://picsum.photos/300/300'
  try {
    const arr = JSON.parse(images)
    return Array.isArray(arr) && arr.length > 0 ? arr[0] : 'https://picsum.photos/300/300'
  } catch {
    return 'https://picsum.photos/300/300'
  }
}

function goToDetail(id: number) {
  router.push(`/products/${id}`)
}

async function handleSearch() {
  pagination.page = 1
  await fetchProducts()
}

async function fetchProducts() {
  loading.value = true
  try {
    const res = await productApi.getList({
      page: pagination.page,
      pageSize: pagination.pageSize,
      keyword: searchForm.keyword || undefined,
      category: searchForm.category || undefined,
      condition: searchForm.condition || undefined,
      minPrice: searchForm.minPrice,
      maxPrice: searchForm.maxPrice,
      sortBy: searchForm.sortBy
    })
    products.value = res.data
    total.value = res.pagination.total
  } catch (error) {
    console.error('Failed to fetch products:', error)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchProducts()
})
</script>

<style lang="scss" scoped>
.product-list-page {
  min-height: 100vh;
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
      min-width: 60px;
    }

    .price-separator {
      margin: 0 8px;
      color: var(--text-lighter-color);
    }
  }
}

.product-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 20px;
  margin-bottom: 20px;
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

.pagination-wrapper {
  display: flex;
  justify-content: center;
  padding: 20px 0;
}
</style>
