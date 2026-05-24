<template>
  <div class="page-container">
    <div class="card">
      <div class="page-header">
        <h2>我的商品</h2>
        <el-button type="primary" :icon="Plus" @click="$router.push('/seller/products/create')">
          发布商品
        </el-button>
      </div>
      <div class="filter-bar">
        <el-radio-group v-model="status" size="default" @change="loadProducts">
          <el-radio-button label="">全部</el-radio-button>
          <el-radio-button label="draft">草稿</el-radio-button>
          <el-radio-button label="on_sale">在售</el-radio-button>
          <el-radio-button label="sold">已售</el-radio-button>
          <el-radio-button label="removed">已下架</el-radio-button>
        </el-radio-group>
      </div>
    </div>

    <div v-if="loading" class="loading-wrapper">
      <el-icon class="is-loading" :size="32"><Loading /></el-icon>
    </div>

    <el-table
      v-else-if="products.length > 0"
      :data="products"
      style="width: 100%"
      stripe
    >
      <el-table-column label="商品" min-width="200">
        <template #default="{ row }">
          <div class="product-cell">
            <img v-if="row.images?.[0]" :src="row.images[0].image_url" />
            <div class="product-info">
              <span class="product-title text-ellipsis">{{ row.title }}</span>
              <span class="product-brand">{{ row.brand_name || '其他品牌' }}</span>
            </div>
          </div>
        </template>
      </el-table-column>
      <el-table-column label="价格" width="120">
        <template #default="{ row }">
          <span class="price-text">¥{{ row.price.toFixed(2) }}</span>
        </template>
      </el-table-column>
      <el-table-column label="库存" width="80" prop="stock" />
      <el-table-column label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="getStatusType(row.status)">{{ getStatusLabel(row.status) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="浏览量" width="80" prop="views" />
      <el-table-column label="认证" width="80">
        <template #default="{ row }">
          <span v-if="row.is_authenticated" class="auth-badge genuine">已认证</span>
          <span v-else>未认证</span>
        </template>
      </el-table-column>
      <el-table-column label="创建时间" width="180">
        <template #default="{ row }">
          {{ formatDate(row.created_at) }}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="200" fixed="right">
        <template #default="{ row }">
          <el-button
            v-if="row.status === 'draft'"
            type="success"
            size="small"
            @click="handlePublish(row)"
          >
            上架
          </el-button>
          <el-button
            v-if="row.status === 'on_sale'"
            type="warning"
            size="small"
            @click="handleRemove(row)"
          >
            下架
          </el-button>
          <el-button
            type="primary"
            size="small"
            @click="$router.push(`/seller/products/${row.id}/edit`)"
          >
            编辑
          </el-button>
          <el-button
            type="danger"
            size="small"
            @click="handleDelete(row)"
          >
            删除
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <div v-else class="empty-state">
      <el-icon :size="64"><Box /></el-icon>
      <p>暂无商品</p>
    </div>

    <div v-if="!loading && total > 0" class="pagination-wrapper">
      <el-pagination
        v-model:current-page="page"
        v-model:page-size="pageSize"
        :total="total"
        layout="total, prev, pager, next"
        @current-change="loadProducts"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { productApi } from '@/api/product'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { Product } from '@/types'
import dayjs from 'dayjs'
import { Plus, Loading, Box } from '@element-plus/icons-vue'

const products = ref<Product[]>([])
const loading = ref(false)
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)
const status = ref('')

const loadProducts = async () => {
  loading.value = true
  try {
    const res = await productApi.listMyProducts({
      page: page.value,
      page_size: pageSize.value,
      status: status.value || undefined
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

const getStatusType = (status: string) => {
  const types: Record<string, string> = {
    draft: 'info',
    on_sale: 'success',
    sold: 'warning',
    removed: 'danger'
  }
  return types[status] || 'info'
}

const getStatusLabel = (status: string) => {
  const labels: Record<string, string> = {
    draft: '草稿',
    on_sale: '在售',
    sold: '已售',
    removed: '已下架'
  }
  return labels[status] || status
}

const formatDate = (date: string) => dayjs(date).format('YYYY-MM-DD HH:mm:ss')

const handlePublish = async (product: Product) => {
  try {
    await ElMessageBox.confirm('确认将此商品上架吗？', '上架确认', {
      confirmButtonText: '确认',
      cancelButtonText: '取消',
      type: 'warning'
    })

    const res = await productApi.updateProductStatus(product.id, 'on_sale')
    if (res.code === 200) {
      ElMessage.success('上架成功')
      loadProducts()
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Publish error:', error)
    }
  }
}

const handleRemove = async (product: Product) => {
  try {
    await ElMessageBox.confirm('确认将此商品下架吗？', '下架确认', {
      confirmButtonText: '确认',
      cancelButtonText: '取消',
      type: 'warning'
    })

    const res = await productApi.updateProductStatus(product.id, 'removed')
    if (res.code === 200) {
      ElMessage.success('下架成功')
      loadProducts()
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Remove error:', error)
    }
  }
}

const handleDelete = async (product: Product) => {
  try {
    await ElMessageBox.confirm('确认删除此商品吗？此操作不可恢复。', '删除确认', {
      confirmButtonText: '确认',
      cancelButtonText: '取消',
      type: 'warning'
    })

    const res = await productApi.deleteProduct(product.id)
    if (res.code === 200) {
      ElMessage.success('删除成功')
      loadProducts()
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Delete error:', error)
    }
  }
}

onMounted(() => {
  loadProducts()
})
</script>

<style lang="scss" scoped>
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
  
  h2 {
    font-size: 20px;
    font-weight: 600;
  }
}

.filter-bar {
  display: flex;
  gap: 12px;
  margin-top: 16px;
}

.product-cell {
  display: flex;
  gap: 12px;
  align-items: center;
  
  img {
    width: 60px;
    height: 60px;
    object-fit: cover;
    border-radius: 4px;
  }
  
  .product-info {
    display: flex;
    flex-direction: column;
    gap: 4px;
    
    .product-title {
      font-size: 14px;
      max-width: 200px;
    }
    
    .product-brand {
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

.auth-badge {
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 500;
  
  &.genuine {
    background: #e7f5e7;
    color: var(--success-color);
  }
}
</style>
