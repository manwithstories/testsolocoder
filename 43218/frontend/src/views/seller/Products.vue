<template>
  <div class="seller-products-page">
    <div class="page-header">
      <h2 class="page-title">商品管理</h2>
      <el-button type="primary" @click="$router.push('/seller/products/create')">
        <el-icon><Plus /></el-icon>
        发布商品
      </el-button>
    </div>

    <div class="search-bar">
      <el-input
        v-model="searchForm.keyword"
        placeholder="搜索商品名称"
        clearable
        style="width: 240px"
      />
      <el-select v-model="searchForm.status" placeholder="状态" clearable style="width: 140px">
        <el-option label="全部" :value="" />
        <el-option label="待审核" :value="0" />
        <el-option label="在售" :value="3" />
        <el-option label="已售出" :value="4" />
        <el-option label="已下架" :value="5" />
      </el-select>
      <el-button type="primary" @click="fetchProducts">搜索</el-button>
    </div>

    <el-card>
      <el-table :data="products" v-loading="loading" stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="title" label="商品名称" min-width="200">
          <template #default="{ row }">
            <div class="product-cell">
              <img :src="getFirstImage(row.images)" class="product-thumb" />
              <span class="product-title text-ellipsis">{{ row.title }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="category" label="分类" width="100" />
        <el-table-column prop="condition" label="成色" width="100" />
        <el-table-column prop="price" label="价格" width="100">
          <template #default="{ row }">
            <span class="price-text">¥{{ row.price.toFixed(2) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="viewCount" label="浏览" width="80" />
        <el-table-column prop="soldCount" label="销量" width="80" />
        <el-table-column prop="createdAt" label="发布时间" width="180">
          <template #default="{ row }">
            {{ formatTime(row.createdAt) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button size="small" link @click="editProduct(row.id)">编辑</el-button>
            <el-button
              v-if="row.status === 3"
              size="small"
              link
              type="warning"
              @click="offShelf(row)"
            >
              下架
            </el-button>
            <el-button size="small" link type="danger" @click="deleteProduct(row)">
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-wrapper">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.pageSize"
          :page-sizes="[10, 20, 50]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="fetchProducts"
          @current-change="fetchProducts"
        />
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { productApi } from '@/api'
import { ProductStatus, ProductStatusText } from '@/types'
import type { Product } from '@/types'

const router = useRouter()

const loading = ref(false)
const products = ref<Product[]>([])
const total = ref(0)

const searchForm = reactive({
  keyword: '',
  status: ''
})

const pagination = reactive({
  page: 1,
  pageSize: 10
})

function getFirstImage(images: string): string {
  if (!images) return 'https://picsum.photos/60/60'
  try {
    const arr = JSON.parse(images)
    return Array.isArray(arr) && arr.length > 0 ? arr[0] : 'https://picsum.photos/60/60'
  } catch {
    return 'https://picsum.photos/60/60'
  }
}

function getStatusText(status: number): string {
  return ProductStatusText[status] || '未知'
}

function getStatusType(status: number): string {
  const typeMap: Record<number, string> = {
    0: 'warning',
    1: 'success',
    2: 'danger',
    3: 'success',
    4: 'info',
    5: 'info'
  }
  return typeMap[status] || 'info'
}

function formatTime(time: string): string {
  return new Date(time).toLocaleString()
}

function editProduct(id: number) {
  router.push(`/products/${id}`)
}

async function offShelf(product: Product) {
  try {
    await ElMessageBox.confirm('确定下架该商品？', '提示', {
      confirmButtonText: '确认',
      cancelButtonText: '取消'
    })
    await productApi.offShelf(product.id)
    ElMessage.success('商品已下架')
    fetchProducts()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Failed to off shelf:', error)
    }
  }
}

async function deleteProduct(product: Product) {
  try {
    await ElMessageBox.confirm('确定删除该商品？此操作不可恢复', '提示', {
      confirmButtonText: '确认',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await productApi.remove(product.id)
    ElMessage.success('商品已删除')
    fetchProducts()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Failed to delete:', error)
    }
  }
}

async function fetchProducts() {
  loading.value = true
  try {
    const res = await productApi.getMyProducts({
      page: pagination.page,
      pageSize: pagination.pageSize,
      status: searchForm.status || undefined
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
.seller-products-page {
  .search-bar {
    display: flex;
    gap: 12px;
    margin-bottom: 20px;
  }

  .product-cell {
    display: flex;
    align-items: center;
    gap: 12px;

    .product-thumb {
      width: 50px;
      height: 50px;
      border-radius: 4px;
      object-fit: cover;
    }

    .product-title {
      max-width: 180px;
    }
  }

  .pagination-wrapper {
    display: flex;
    justify-content: center;
    padding: 20px 0 0;
  }
}
</style>
