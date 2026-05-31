<template>
  <div class="product-detail-page">
    <div class="container">
      <el-breadcrumb class="breadcrumb" separator="/">
        <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
        <el-breadcrumb-item :to="{ path: '/products' }">商品市场</el-breadcrumb-item>
        <el-breadcrumb-item>{{ product?.title }}</el-breadcrumb-item>
      </el-breadcrumb>

      <div class="product-content card" v-loading="loading">
        <div class="product-gallery">
          <el-image
            :src="currentImage"
            :preview-src-list="imageList"
            fit="cover"
            class="main-image"
          />
          <div class="thumbnail-list">
            <div
              v-for="(img, index) in imageList"
              :key="index"
              class="thumbnail"
              :class="{ active: currentImageIndex === index }"
              @click="currentImageIndex = index"
            >
              <img :src="img" />
            </div>
          </div>
        </div>

        <div class="product-info">
          <h1 class="product-title">{{ product?.title }}</h1>
          <div class="product-meta">
            <span class="status-tag status-success" v-if="product?.status === 3">在售</span>
            <span class="status-tag status-pending" v-else>状态：{{ getStatusText(product?.status) }}</span>
            <span class="view-count">浏览 {{ product?.viewCount }} 次</span>
          </div>

          <div class="price-section">
            <span class="current-price">¥{{ product?.price.toFixed(2) }}</span>
            <span class="original-price" v-if="product?.originalPrice">
              原价 ¥{{ product?.originalPrice.toFixed(2) }}
            </span>
          </div>

          <div class="product-specs">
            <div class="spec-item">
              <span class="spec-label">分类：</span>
              <span>{{ product?.category }}</span>
            </div>
            <div class="spec-item">
              <span class="spec-label">品牌：</span>
              <span>{{ product?.brand }}</span>
            </div>
            <div class="spec-item">
              <span class="spec-label">型号：</span>
              <span>{{ product?.model }}</span>
            </div>
            <div class="spec-item">
              <span class="spec-label">成色：</span>
              <span>{{ product?.condition }}</span>
            </div>
            <div class="spec-item">
              <span class="spec-label">质保：</span>
              <span>{{ product?.warrantyDays }} 天</span>
            </div>
          </div>

          <div class="seller-info">
            <el-avatar :size="40" :src="product?.seller?.avatar">
              {{ product?.seller?.nickname?.charAt(0) }}
            </el-avatar>
            <div class="seller-detail">
              <div class="seller-name">{{ product?.seller?.nickname || product?.seller?.username }}</div>
              <div class="seller-credit">
                <span>信用分：{{ product?.seller?.creditScore }}</span>
                <el-rate :model-value="getRating(product?.seller?.creditScore)" disabled size="small" />
              </div>
            </div>
          </div>

          <div class="action-buttons">
            <el-button
              type="primary"
              size="large"
              :disabled="!canBuy"
              @click="showOrderDialog = true"
            >
              立即购买
            </el-button>
            <el-button
              size="large"
              :disabled="!userStore.isLoggedIn"
              @click="handleToggleFavorite"
            >
              <el-icon><Star /></el-icon>
              {{ isFavorited ? '已收藏' : '收藏' }}
            </el-button>
            <el-button
              size="large"
              :disabled="!userStore.isLoggedIn"
              @click="showNegotiationDialog = true"
            >
              议价
            </el-button>
          </div>
        </div>
      </div>

      <div class="product-description card">
        <h3>商品描述</h3>
        <p>{{ product?.description }}</p>
      </div>
    </div>

    <el-dialog v-model="showOrderDialog" title="确认订单" width="500px">
      <el-form :model="orderForm" label-width="100px">
        <el-form-item label="商品名称">
          <span>{{ product?.title }}</span>
        </el-form-item>
        <el-form-item label="商品价格">
          <span class="price-text">¥{{ product?.price.toFixed(2) }}</span>
        </el-form-item>
        <el-form-item label="议价价格" v-if="negotiatedPrice > 0">
          <span class="price-text">¥{{ negotiatedPrice.toFixed(2) }}</span>
        </el-form-item>
        <el-form-item label="收货人" prop="receiverName">
          <el-input v-model="orderForm.receiverName" placeholder="请输入收货人姓名" />
        </el-form-item>
        <el-form-item label="联系电话" prop="receiverPhone">
          <el-input v-model="orderForm.receiverPhone" placeholder="请输入联系电话" />
        </el-form-item>
        <el-form-item label="收货地址" prop="receiverAddress">
          <el-input
            v-model="orderForm.receiverAddress"
            type="textarea"
            :rows="2"
            placeholder="请输入收货地址"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showOrderDialog = false">取消</el-button>
        <el-button type="primary" :loading="submittingOrder" @click="submitOrder">
          提交订单
        </el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="showNegotiationDialog" title="议价" width="400px">
      <el-form :model="negotiationForm" label-width="80px">
        <el-form-item label="商品价格">
          <span>¥{{ product?.price.toFixed(2) }}</span>
        </el-form-item>
        <el-form-item label="期望价格" prop="offeredPrice">
          <el-input-number
            v-model="negotiationForm.offeredPrice"
            :min="1"
            :max="product?.price || 0"
            :precision="2"
          />
        </el-form-item>
        <el-form-item label="留言">
          <el-input
            v-model="negotiationForm.message"
            type="textarea"
            :rows="2"
            placeholder="请输入议价说明（选填）"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showNegotiationDialog = false">取消</el-button>
        <el-button type="primary" :loading="submittingNegotiation" @click="submitNegotiation">
          提交议价
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, reactive, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/store/user'
import { productApi, orderApi } from '@/api'
import { ProductStatusText } from '@/types'
import type { Product } from '@/types'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const loading = ref(false)
const product = ref<Product | null>(null)
const currentImageIndex = ref(0)
const isFavorited = ref(false)
const showOrderDialog = ref(false)
const showNegotiationDialog = ref(false)
const submittingOrder = ref(false)
const submittingNegotiation = ref(false)
const negotiatedPrice = ref(0)

const orderForm = reactive({
  receiverName: '',
  receiverPhone: '',
  receiverAddress: ''
})

const negotiationForm = reactive({
  offeredPrice: 0,
  message: ''
})

const imageList = computed(() => {
  if (!product.value?.images) return ['https://picsum.photos/600/600']
  try {
    const arr = JSON.parse(product.value.images)
    return Array.isArray(arr) && arr.length > 0 ? arr : ['https://picsum.photos/600/600']
  } catch {
    return ['https://picsum.photos/600/600']
  }
})

const currentImage = computed(() => imageList.value[currentImageIndex.value] || 'https://picsum.photos/600/600')

const canBuy = computed(() => {
  return userStore.isLoggedIn &&
    product.value?.status === 3 &&
    product.value?.sellerId !== userStore.userInfo?.id
})

function getStatusText(status?: number): string {
  return ProductStatusText[status || 0] || '未知'
}

function getRating(creditScore?: number): number {
  if (!creditScore) return 3
  if (creditScore >= 90) return 5
  if (creditScore >= 80) return 4
  if (creditScore >= 70) return 3
  if (creditScore >= 60) return 2
  return 1
}

async function fetchProduct() {
  const id = parseInt(route.params.id as string)
  if (!id) return

  loading.value = true
  try {
    const res = await productApi.getById(id)
    product.value = res.data
  } catch (error) {
    console.error('Failed to fetch product:', error)
    ElMessage.error('获取商品详情失败')
  } finally {
    loading.value = false
  }
}

async function handleToggleFavorite() {
  if (!product.value) return
  try {
    const res = await productApi.toggleFavorite(product.value.id)
    isFavorited.value = res.data.isFavorited
    ElMessage.success(isFavorited.value ? '已收藏' : '已取消收藏')
  } catch (error) {
    console.error('Failed to toggle favorite:', error)
  }
}

async function submitOrder() {
  if (!orderForm.receiverName || !orderForm.receiverPhone || !orderForm.receiverAddress) {
    ElMessage.warning('请填写完整的收货信息')
    return
  }

  if (!product.value) return

  submittingOrder.value = true
  try {
    const res = await orderApi.create({
      productId: product.value.id,
      receiverName: orderForm.receiverName,
      receiverPhone: orderForm.receiverPhone,
      receiverAddress: orderForm.receiverAddress,
      negotiatedPrice: negotiatedPrice.value
    })
    ElMessage.success('订单创建成功')
    showOrderDialog.value = false
    router.push(`/user/orders`)
  } catch (error: any) {
    ElMessage.error(error.message || '创建订单失败')
  } finally {
    submittingOrder.value = false
  }
}

async function submitNegotiation() {
  if (!negotiationForm.offeredPrice) {
    ElMessage.warning('请输入期望价格')
    return
  }

  if (!product.value) return

  submittingNegotiation.value = true
  try {
    const res = await orderApi.create({
      productId: product.value.id,
      receiverName: userStore.userInfo?.nickname || userStore.userInfo?.username || '',
      receiverPhone: userStore.userInfo?.phone || '',
      receiverAddress: '',
      negotiatedPrice: negotiationForm.offeredPrice
    })

    await orderApi.negotiate({
      orderNo: res.data.orderNo,
      offeredPrice: negotiationForm.offeredPrice,
      message: negotiationForm.message
    })

    ElMessage.success('议价已提交')
    showNegotiationDialog.value = false
    negotiatedPrice.value = negotiationForm.offeredPrice
  } catch (error: any) {
    ElMessage.error(error.message || '提交议价失败')
  } finally {
    submittingNegotiation.value = false
  }
}

onMounted(() => {
  fetchProduct()
})
</script>

<style lang="scss" scoped>
.product-detail-page {
  min-height: 100vh;
  background: #f5f7fa;
  padding: 20px 0;
}

.breadcrumb {
  margin-bottom: 20px;
}

.product-content {
  display: flex;
  gap: 40px;
  margin-bottom: 20px;
}

.product-gallery {
  width: 400px;
  flex-shrink: 0;

  .main-image {
    width: 400px;
    height: 400px;
    border-radius: 8px;
    overflow: hidden;
    margin-bottom: 16px;
  }

  .thumbnail-list {
    display: flex;
    gap: 12px;

    .thumbnail {
      width: 70px;
      height: 70px;
      border-radius: 4px;
      overflow: hidden;
      cursor: pointer;
      border: 2px solid transparent;

      img {
        width: 100%;
        height: 100%;
        object-fit: cover;
      }

      &.active {
        border-color: var(--primary-color);
      }
    }
  }
}

.product-info {
  flex: 1;

  .product-title {
    font-size: 24px;
    margin-bottom: 16px;
  }

  .product-meta {
    display: flex;
    align-items: center;
    gap: 16px;
    margin-bottom: 20px;

    .view-count {
      color: var(--text-lighter-color);
    }
  }

  .price-section {
    background: #fff5f5;
    padding: 20px;
    border-radius: 8px;
    margin-bottom: 20px;
    display: flex;
    align-items: baseline;
    gap: 16px;

    .current-price {
      font-size: 32px;
      color: var(--danger-color);
      font-weight: 600;
    }

    .original-price {
      font-size: 16px;
      color: var(--text-lighter-color);
      text-decoration: line-through;
    }
  }

  .product-specs {
    margin-bottom: 20px;

    .spec-item {
      display: flex;
      margin-bottom: 12px;
      font-size: 14px;

      .spec-label {
        color: var(--text-lighter-color);
        width: 80px;
      }
    }
  }

  .seller-info {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 16px;
    background: #f5f7fa;
    border-radius: 8px;
    margin-bottom: 20px;

    .seller-detail {
      .seller-name {
        font-weight: 500;
        margin-bottom: 4px;
      }

      .seller-credit {
        display: flex;
        align-items: center;
        gap: 8px;
        font-size: 12px;
        color: var(--text-lighter-color);
      }
    }
  }

  .action-buttons {
    display: flex;
    gap: 12px;

    .el-button {
      flex: 1;
    }
  }
}

.product-description {
  h3 {
    margin-bottom: 16px;
  }

  p {
    line-height: 1.8;
    color: var(--text-light-color);
    white-space: pre-wrap;
  }
}
</style>
