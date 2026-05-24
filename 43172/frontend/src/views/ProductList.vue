<template>
  <div class="page-container">
    <div class="card">
      <div class="filter-bar">
        <el-input
          v-model="keyword"
          placeholder="搜索商品名称、品牌..."
          :prefix-icon="Search"
          clearable
          style="width: 300px"
          @keyup.enter="loadProducts"
        />
        <el-select v-model="category" placeholder="全部分类" clearable style="width: 150px">
          <el-option
            v-for="cat in CATEGORY_OPTIONS"
            :key="cat.value"
            :label="cat.label"
            :value="cat.value"
          />
        </el-select>
        <el-select v-model="sortBy" placeholder="排序方式" style="width: 150px">
          <el-option label="最新发布" value="newest" />
          <el-option label="价格从低到高" value="price_asc" />
          <el-option label="价格从高到低" value="price_desc" />
          <el-option label="最多浏览" value="views" />
        </el-select>
        <el-button type="primary" :icon="Search" @click="loadProducts">搜索</el-button>
      </div>
    </div>

    <div v-if="loading" class="loading-wrapper">
      <el-icon class="is-loading" :size="32"><Loading /></el-icon>
    </div>

    <el-row v-else :gutter="16">
      <el-col :span="6" v-for="product in products" :key="product.id">
        <div class="product-card" @click="$router.push(`/products/${product.id}`)">
          <div class="product-image">
            <img v-if="product.images?.[0]" :src="product.images[0].image_url" :alt="product.title" />
            <div v-else class="no-image">暂无图片</div>
            <div v-if="product.is_authenticated" class="auth-badge genuine">已认证</div>
          </div>
          <div class="product-info">
            <h3 class="product-title text-ellipsis">{{ product.title }}</h3>
            <p class="product-brand">{{ product.brand_name || '其他品牌' }}</p>
            <p class="product-price">¥{{ product.price.toFixed(2) }}</p>
            <p class="product-views">{{ product.views }} 次浏览</p>
          </div>
        </div>
      </el-col>
    </el-row>

    <div v-if="!loading && products.length === 0" class="empty-state">
      <el-icon :size="64"><Box /></el-icon>
      <p>暂无商品</p>
    </div>

    <div v-if="!loading && total > 0" class="pagination-wrapper">
      <el-pagination
        v-model:current-page="page"
        v-model:page-size="pageSize"
        :page-sizes="[12, 24, 48]"
        :total="total"
        layout="total, sizes, prev, pager, next, jumper"
        @current-change="loadProducts"
        @size-change="loadProducts"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { productApi } from '@/api/product'
import type { Product } from '@/types'
import { CATEGORY_OPTIONS } from '@/types'
import { Search, Loading, Box } from '@element-plus/icons-vue'

const products = ref<Product[]>([])
const loading = ref(false)
const total = ref(0)
const page = ref(1)
const pageSize = ref(12)
const keyword = ref('')
const category = ref('')
const sortBy = ref('newest')

const loadProducts = async () => {
  loading.value = true
  try {
    const res = await productApi.listProducts({
      page: page.value,
      page_size: pageSize.value,
      category: category.value || undefined,
      status: 'on_sale',
      sort_by: sortBy.value,
      keyword: keyword.value || undefined
    })
    if (res.code === 200) {
      const data = res.data as any
      products.value = data?.list || []
      total.value = data?.total || 0
    }
  } catch (error) {
    console.error('Load products error:', error)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadProducts()
})
</script>

<style lang="scss" scoped>
.filter-bar {
  display: flex;
  gap: 12px;
  align-items: center;
}

.product-card {
  background: #fff;
  border-radius: 8px;
  overflow: hidden;
  cursor: pointer;
  transition: all 0.3s;
  margin-bottom: 16px;
  border: 1px solid var(--border-color);
  
  &:hover {
    transform: translateY(-4px);
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  }
  
  .product-image {
    position: relative;
    height: 240px;
    background: #f5f5f5;
    
    img {
      width: 100%;
      height: 100%;
      object-fit: cover;
    }
    
    .no-image {
      display: flex;
      align-items: center;
      justify-content: center;
      height: 100%;
      color: var(--text-light);
    }
    
    .auth-badge {
      position: absolute;
      top: 8px;
      right: 8px;
    }
  }
  
  .product-info {
    padding: 16px;
    
    .product-title {
      font-size: 16px;
      font-weight: 500;
      margin-bottom: 8px;
    }
    
    .product-brand {
      font-size: 13px;
      color: var(--text-light);
      margin-bottom: 8px;
    }
    
    .product-price {
      font-size: 20px;
      font-weight: 600;
      color: var(--danger-color);
      margin-bottom: 4px;
    }
    
    .product-views {
      font-size: 12px;
      color: var(--text-light);
    }
  }
}

.pagination-wrapper {
  display: flex;
  justify-content: center;
  margin-top: 24px;
}
</style>
