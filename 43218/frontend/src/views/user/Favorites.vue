<template>
  <div class="favorites-page">
    <div class="page-header">
      <h2 class="page-title">我的收藏</h2>
    </div>

    <div class="product-grid" v-loading="loading">
      <div
        v-for="item in favorites"
        :key="item.id"
        class="product-card"
        @click="goToProduct(item.productId)"
      >
        <div class="product-image">
          <img :src="getFirstImage(item.product?.images)" :alt="item.product?.title" />
          <span class="condition-tag">{{ item.product?.condition }}</span>
        </div>
        <div class="product-info">
          <h3 class="product-title text-ellipsis">{{ item.product?.title }}</h3>
          <p class="product-brand">{{ item.product?.brand }} {{ item.product?.model }}</p>
          <div class="product-price">
            <span class="price-text">¥{{ item.product?.price.toFixed(2) }}</span>
          </div>
          <div class="product-actions">
            <el-button size="small" type="danger" @click.stop="removeFavorite(item.id)">
              取消收藏
            </el-button>
          </div>
        </div>
      </div>

      <div class="empty-state" v-if="!loading && favorites.length === 0">
        <el-empty description="暂无收藏" />
      </div>
    </div>

    <div class="pagination-wrapper" v-if="total > 0">
      <el-pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.pageSize"
        :page-sizes="[12, 24, 48]"
        :total="total"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="fetchFavorites"
        @current-change="fetchFavorites"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { productApi } from '@/api'
import type { Favorite } from '@/types'

const router = useRouter()

const loading = ref(false)
const favorites = ref<Favorite[]>([])
const total = ref(0)

const pagination = reactive({
  page: 1,
  pageSize: 12
})

function getFirstImage(images: string | undefined): string {
  if (!images) return 'https://picsum.photos/300/300'
  try {
    const arr = JSON.parse(images)
    return Array.isArray(arr) && arr.length > 0 ? arr[0] : 'https://picsum.photos/300/300'
  } catch {
    return 'https://picsum.photos/300/300'
  }
}

function goToProduct(id: number | undefined) {
  if (id) router.push(`/products/${id}`)
}

async function removeFavorite(id: number) {
  try {
    await ElMessageBox.confirm('确定取消收藏？', '提示', {
      confirmButtonText: '确认',
      cancelButtonText: '取消'
    })
    const favorite = favorites.value.find(f => f.id === id)
    if (favorite?.productId) {
      await productApi.toggleFavorite(favorite.productId)
      ElMessage.success('已取消收藏')
      fetchFavorites()
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Failed to remove favorite:', error)
    }
  }
}

async function fetchFavorites() {
  loading.value = true
  try {
    const res = await productApi.getFavorites({
      page: pagination.page,
      pageSize: pagination.pageSize
    })
    favorites.value = res.data
    total.value = res.pagination.total
  } catch (error) {
    console.error('Failed to fetch favorites:', error)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchFavorites()
})
</script>

<style lang="scss" scoped>
.favorites-page {
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
        margin-bottom: 12px;

        .price-text {
          font-size: 18px;
        }
      }

      .product-actions {
        text-align: center;
      }
    }
  }

  .pagination-wrapper {
    display: flex;
    justify-content: center;
    padding: 20px 0;
  }
}
</style>
