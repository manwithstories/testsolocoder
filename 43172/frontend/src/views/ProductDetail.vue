<template>
  <div class="page-container">
    <div v-if="loading" class="loading-wrapper">
      <el-icon class="is-loading" :size="32"><Loading /></el-icon>
    </div>

    <template v-else-if="product">
      <el-row :gutter="24">
        <el-col :span="12">
          <div class="product-gallery">
            <el-carousel height="400px">
              <el-carousel-item v-for="(image, index) in product.images" :key="index">
                <img :src="image.image_url" :alt="product.title" />
              </el-carousel-item>
              <el-carousel-item v-if="!product.images || product.images.length === 0">
                <div class="no-image">暂无图片</div>
              </el-carousel-item>
            </el-carousel>
          </div>
        </el-col>
        <el-col :span="12">
          <div class="product-detail">
            <div class="product-header">
              <h1 class="product-title">{{ product.title }}</h1>
              <div class="product-badges">
                <span v-if="product.is_authenticated" class="auth-badge genuine">
                  ✅ 已认证
                </span>
                <span v-if="product.brand_name" class="brand-badge">
                  {{ product.brand_name }}
                </span>
              </div>
            </div>
            
            <div class="product-price">
              <span class="current-price">¥{{ product.price.toFixed(2) }}</span>
              <span v-if="product.original_price" class="original-price">
                ¥{{ product.original_price.toFixed(2) }}
              </span>
            </div>

            <el-descriptions :column="2" border class="product-info">
              <el-descriptions-item label="分类">
                {{ getCategoryLabel(product.category) }}
              </el-descriptions-item>
              <el-descriptions-item label="成色">
                {{ product.condition || '未填写' }}
              </el-descriptions-item>
              <el-descriptions-item label="颜色">
                {{ product.color || '未填写' }}
              </el-descriptions-item>
              <el-descriptions-item label="尺寸">
                {{ product.size || '未填写' }}
              </el-descriptions-item>
              <el-descriptions-item label="材质">
                {{ product.material || '未填写' }}
              </el-descriptions-item>
              <el-descriptions-item label="库存">
                <span :class="{ 'out-of-stock': product.stock <= 0 }">
                  {{ product.stock }} 件
                </span>
              </el-descriptions-item>
              <el-descriptions-item label="浏览量" :span="2">
                {{ product.views }} 次
              </el-descriptions-item>
            </el-descriptions>

            <div class="product-actions" v-if="userStore.isLoggedIn && userStore.userRole === 'buyer'">
              <el-checkbox v-model="needAuth">需要鉴定服务</el-checkbox>
              <el-button
                type="primary"
                size="large"
                :disabled="product.stock <= 0"
                :loading="buying"
                @click="handleBuy"
              >
                {{ product.stock <= 0 ? '暂时缺货' : '立即购买' }}
              </el-button>
            </div>

            <div class="seller-info" v-if="product.seller">
              <el-divider />
              <div class="seller-profile">
                <el-avatar :size="48" :src="product.seller.avatar">
                  {{ product.seller.username?.charAt(0)?.toUpperCase() }}
                </el-avatar>
                <div class="seller-details">
                  <h4>{{ product.seller.username }}</h4>
                  <p>信用评分: {{ product.seller.credit_score }}</p>
                </div>
              </div>
            </div>
          </div>
        </el-col>
      </el-row>

      <div class="card product-description">
        <h3>商品描述</h3>
        <p>{{ product.description }}</p>
      </div>

      <div v-if="authentication" class="card auth-section">
        <h3>鉴定记录</h3>
        <el-descriptions :column="2" border>
          <el-descriptions-item label="鉴定状态">
            <el-tag :type="getAuthStatusType(authentication.status)">
              {{ getAuthStatusLabel(authentication.status) }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="鉴定结果">
            <el-tag v-if="authentication.result" :type="getAuthResultType(authentication.result)">
              {{ getAuthResultLabel(authentication.result) }}
            </el-tag>
            <span v-else>待鉴定</span>
          </el-descriptions-item>
          <el-descriptions-item label="鉴定师" :span="2">
            {{ authentication.authenticator?.username || '待分配' }}
          </el-descriptions-item>
          <el-descriptions-item v-if="authentication.report_content" label="鉴定报告" :span="2">
            {{ authentication.report_content }}
          </el-descriptions-item>
        </el-descriptions>
        <div v-if="authentication.report_file" class="report-download">
          <el-button type="primary" :icon="Download" @click="downloadReport">
            下载鉴定报告
          </el-button>
        </div>
      </div>
    </template>

    <div v-else class="empty-state">
      <el-icon :size="64"><Box /></el-icon>
      <p>商品不存在</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { productApi } from '@/api/product'
import { orderApi } from '@/api/order'
import { authServiceApi } from '@/api/authentication'
import { useUserStore } from '@/stores/user'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { Product, Authentication } from '@/types'
import { AUTH_STATUS_OPTIONS, AUTH_RESULT_OPTIONS, CATEGORY_OPTIONS } from '@/types'
import { Loading, Box, Download } from '@element-plus/icons-vue'

const route = useRoute()
const userStore = useUserStore()

const product = ref<Product | null>(null)
const authentication = ref<Authentication | null>(null)
const loading = ref(false)
const buying = ref(false)
const needAuth = ref(false)

const productId = Number(route.params.id)

const loadProduct = async () => {
  loading.value = true
  try {
    const res = await productApi.getProduct(productId)
    if (res.code === 200 && res.data) {
      product.value = res.data
      
      if (product.value.is_authenticated) {
        loadAuthentication()
      }
    }
  } catch (error) {
    console.error('Load product error:', error)
  } finally {
    loading.value = false
  }
}

const loadAuthentication = async () => {
  try {
    const res = await authServiceApi.getAuthenticationByOrder(productId)
    if (res.code === 200 && res.data) {
      authentication.value = res.data
    }
  } catch (error) {
    console.error('Load authentication error:', error)
  }
}

const getCategoryLabel = (category: string) => {
  const cat = CATEGORY_OPTIONS.find(c => c.value === category)
  return cat?.label || category
}

const getAuthStatusType = (status: string) => {
  const opt = AUTH_STATUS_OPTIONS.find(o => o.value === status)
  return opt?.type || 'info'
}

const getAuthStatusLabel = (status: string) => {
  const opt = AUTH_STATUS_OPTIONS.find(o => o.value === status)
  return opt?.label || status
}

const getAuthResultType = (result: string) => {
  const opt = AUTH_RESULT_OPTIONS.find(o => o.value === result)
  return opt?.type || 'info'
}

const getAuthResultLabel = (result: string) => {
  const opt = AUTH_RESULT_OPTIONS.find(o => o.value === result)
  return opt?.label || result
}

const handleBuy = async () => {
  if (!product.value) return

  try {
    await ElMessageBox.confirm(
      `确认购买 "${product.value.title}" 吗？${needAuth.value ? '将同时申请鉴定服务。' : ''}`,
      '确认购买',
      {
        confirmButtonText: '确认',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    buying.value = true
    const res = await orderApi.createOrder({
      product_id: product.value.id,
      shipping_address: userStore.user?.address || '请填写收货地址',
      need_auth: needAuth.value
    })

    if (res.code === 201 || res.code === 200) {
      ElMessage.success('下单成功，请前往支付')
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Buy error:', error)
    }
  } finally {
    buying.value = false
  }
}

const downloadReport = () => {
  if (authentication.value?.report_file) {
    window.open(`/api/v1/authentications/${authentication.value.id}/report/download`)
  }
}

onMounted(() => {
  loadProduct()
})
</script>

<style lang="scss" scoped>
.product-gallery {
  background: #fff;
  border-radius: 8px;
  overflow: hidden;
  
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
    background: #f5f5f5;
    color: var(--text-light);
  }
}

.product-detail {
  .product-header {
    margin-bottom: 16px;
    
    .product-title {
      font-size: 24px;
      font-weight: 600;
      margin-bottom: 8px;
    }
    
    .product-badges {
      display: flex;
      gap: 8px;
      
      .auth-badge {
        padding: 4px 8px;
        border-radius: 4px;
        font-size: 12px;
        font-weight: 500;
        
        &.genuine {
          background: #e7f5e7;
          color: var(--success-color);
        }
      }
      
      .brand-badge {
        padding: 4px 8px;
        border-radius: 4px;
        font-size: 12px;
        background: #f0f0f0;
        color: var(--text-secondary);
      }
    }
  }
  
  .product-price {
    margin-bottom: 16px;
    padding: 16px;
    background: #fff5f5;
    border-radius: 8px;
    
    .current-price {
      font-size: 28px;
      font-weight: 600;
      color: var(--danger-color);
    }
    
    .original-price {
      font-size: 16px;
      color: var(--text-light);
      text-decoration: line-through;
      margin-left: 12px;
    }
  }
  
  .product-info {
    margin-bottom: 16px;
  }
  
  .out-of-stock {
    color: var(--danger-color);
  }
  
  .product-actions {
    display: flex;
    gap: 16px;
    align-items: center;
    margin-bottom: 16px;
    
    .el-button {
      flex: 1;
    }
  }
  
  .seller-info {
    .seller-profile {
      display: flex;
      gap: 12px;
      align-items: center;
      
      .seller-details {
        h4 {
          font-size: 16px;
          font-weight: 500;
          margin-bottom: 4px;
        }
        
        p {
          font-size: 14px;
          color: var(--text-secondary);
        }
      }
    }
  }
}

.product-description {
  margin-top: 24px;
  
  h3 {
    font-size: 18px;
    font-weight: 600;
    margin-bottom: 12px;
  }
  
  p {
    white-space: pre-wrap;
    line-height: 1.8;
    color: var(--text-secondary);
  }
}

.auth-section {
  margin-top: 24px;
  
  h3 {
    font-size: 18px;
    font-weight: 600;
    margin-bottom: 12px;
  }
  
  .report-download {
    margin-top: 16px;
    text-align: right;
  }
}
</style>
