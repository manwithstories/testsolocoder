<template>
  <div class="home-container">
    <section class="hero-section">
      <div class="hero-content">
        <h1>奢侈品二手交易与鉴定平台</h1>
        <p>专业鉴定·安全交易·品质保障</p>
        <div class="hero-buttons">
          <el-button type="primary" size="large" @click="$router.push('/products')">
            浏览商品
          </el-button>
          <el-button size="large" @click="$router.push('/register')">
            立即注册
          </el-button>
        </div>
      </div>
    </section>

    <section class="features-section">
      <div class="container">
        <el-row :gutter="20">
          <el-col :span="6">
            <div class="feature-card">
              <div class="feature-icon">🛡️</div>
              <h3>专业鉴定</h3>
              <p>资深鉴定师把关，正品保障</p>
            </div>
          </el-col>
          <el-col :span="6">
            <div class="feature-card">
              <div class="feature-icon">💰</div>
              <h3>安全交易</h3>
              <p>担保交易，资金安全有保障</p>
            </div>
          </el-col>
          <el-col :span="6">
            <div class="feature-card">
              <div class="feature-icon">📦</div>
              <h3>品质严选</h3>
              <p>每件商品都经过严格筛选</p>
            </div>
          </el-col>
          <el-col :span="6">
            <div class="feature-card">
              <div class="feature-icon">⭐</div>
              <h3>信用评价</h3>
              <p>双向评价体系，交易更放心</p>
            </div>
          </el-col>
        </el-row>
      </div>
    </section>

    <section class="categories-section">
      <div class="container">
        <h2 class="section-title">热门分类</h2>
        <el-row :gutter="16">
          <el-col :span="8" v-for="cat in categories" :key="cat.value">
            <div class="category-card" @click="$router.push(`/products?category=${cat.value}`)">
              <div class="category-icon">{{ cat.icon }}</div>
              <span>{{ cat.label }}</span>
            </div>
          </el-col>
        </el-row>
      </div>
    </section>

    <section class="hot-products-section">
      <div class="container">
        <div class="section-header">
          <h2 class="section-title">热门商品</h2>
          <el-link type="primary" @click="$router.push('/products')">查看更多 →</el-link>
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
              </div>
            </div>
          </el-col>
        </el-row>
        <div v-if="!loading && products.length === 0" class="empty-state">
          <el-icon :size="64"><Box /></el-icon>
          <p>暂无商品</p>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { productApi } from '@/api/product'
import type { Product } from '@/types'
import { Loading, Box } from '@element-plus/icons-vue'
import { CATEGORY_OPTIONS } from '@/types'

const products = ref<Product[]>([])
const loading = ref(false)

const categories = [
  { label: '包包', value: 'bag', icon: '👜' },
  { label: '首饰', value: 'jewelry', icon: '💎' },
  { label: '手表', value: 'watch', icon: '⌚' },
  { label: '服装', value: 'clothing', icon: '👗' },
  { label: '鞋履', value: 'shoes', icon: '👠' },
  { label: '其他', value: 'other', icon: '✨' }
]

const loadProducts = async () => {
  loading.value = true
  try {
    const res = await productApi.listProducts({
      page: 1,
      page_size: 8,
      status: 'on_sale',
      sort_by: 'views'
    })
    if (res.code === 200) {
      products.value = (res.data as any)?.list || []
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
.home-container {
  min-height: 100vh;
}

.hero-section {
  background: linear-gradient(135deg, #1a1a1a 0%, #2d2d2d 100%);
  padding: 80px 20px;
  text-align: center;
  color: #fff;
  
  .hero-content {
    max-width: 800px;
    margin: 0 auto;
    
    h1 {
      font-size: 42px;
      font-weight: 700;
      margin-bottom: 16px;
    }
    
    p {
      font-size: 18px;
      opacity: 0.8;
      margin-bottom: 32px;
    }
    
    .hero-buttons {
      display: flex;
      gap: 16px;
      justify-content: center;
    }
  }
}

.features-section {
  background: #fff;
  padding: 60px 20px;
  
  .container {
    max-width: 1200px;
    margin: 0 auto;
  }
  
  .feature-card {
    text-align: center;
    padding: 24px;
    
    .feature-icon {
      font-size: 48px;
      margin-bottom: 16px;
    }
    
    h3 {
      font-size: 18px;
      font-weight: 600;
      margin-bottom: 8px;
    }
    
    p {
      color: var(--text-secondary);
      font-size: 14px;
    }
  }
}

.categories-section {
  background: #f8f9fa;
  padding: 40px 20px;
  
  .container {
    max-width: 1200px;
    margin: 0 auto;
  }
  
  .section-title {
    font-size: 24px;
    font-weight: 600;
    margin-bottom: 24px;
    text-align: center;
  }
  
  .category-card {
    background: #fff;
    border-radius: 8px;
    padding: 24px;
    text-align: center;
    cursor: pointer;
    transition: all 0.3s;
    margin-bottom: 16px;
    
    &:hover {
      transform: translateY(-4px);
      box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
    }
    
    .category-icon {
      font-size: 36px;
      margin-bottom: 8px;
    }
    
    span {
      font-size: 16px;
      font-weight: 500;
    }
  }
}

.hot-products-section {
  background: #fff;
  padding: 60px 20px;
  
  .container {
    max-width: 1200px;
    margin: 0 auto;
  }
  
  .section-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 24px;
  }
  
  .section-title {
    font-size: 24px;
    font-weight: 600;
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
      height: 200px;
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
      padding: 12px;
      
      .product-title {
        font-size: 14px;
        font-weight: 500;
        margin-bottom: 4px;
      }
      
      .product-brand {
        font-size: 12px;
        color: var(--text-light);
        margin-bottom: 4px;
      }
      
      .product-price {
        font-size: 16px;
        font-weight: 600;
        color: var(--danger-color);
      }
    }
  }
}
</style>
