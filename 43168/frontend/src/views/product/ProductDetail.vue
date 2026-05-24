<template>
  <div class="product-detail">
    <el-page-header :icon="ArrowLeft" @back="goBack">
      <template #content>
        <span>产品详情</span>
      </template>
    </el-page-header>

    <div v-loading="loading" class="detail-wrapper">
      <el-row :gutter="20">
        <el-col :xs="24" :md="12">
          <el-card shadow="never" class="image-card">
            <el-carousel
              v-if="imageList.length > 0"
              height="420px"
              indicator-position="outside"
              arrow="always"
            >
              <el-carousel-item v-for="(img, idx) in imageList" :key="idx">
                <el-image :src="img" fit="contain" style="width: 100%; height: 100%" />
              </el-carousel-item>
            </el-carousel>
            <el-empty v-else description="暂无图片" />
          </el-card>
        </el-col>

        <el-col :xs="24" :md="12">
          <el-card shadow="never" class="info-card">
            <div class="product-title">
              <h2>{{ product.name }}</h2>
              <el-tag v-if="product.isHot" type="danger" effect="dark">热门</el-tag>
              <el-tag :type="product.status === 1 ? 'success' : 'info'">
                {{ product.status === 1 ? '上架中' : '已下架' }}
              </el-tag>
            </div>
            <div class="product-price">
              <span class="label">参考价</span>
              <span class="value">¥{{ formatPrice(product.price) }}</span>
            </div>
            <el-descriptions :column="1" border size="default" class="info-desc">
              <el-descriptions-item label="SKU">{{ product.sku || '-' }}</el-descriptions-item>
              <el-descriptions-item label="分类">{{ product.category || '-' }}</el-descriptions-item>
              <el-descriptions-item label="库存">{{ product.stock }}</el-descriptions-item>
              <el-descriptions-item label="厂商">{{ product.manufacturerName || '-' }}</el-descriptions-item>
              <el-descriptions-item label="创建时间">{{ product.createdAt }}</el-descriptions-item>
            </el-descriptions>

            <el-divider>自定义选项</el-divider>

            <div v-if="product.options && product.options.length > 0" class="options">
              <div
                v-for="opt in product.options"
                :key="opt.id || opt.type"
                class="option-group"
              >
                <div class="option-label">
                  {{ optionTypeLabel(opt.type) }} · {{ opt.name }}
                  <span v-if="opt.required" class="required">*</span>
                </div>
                <el-radio-group v-model="selectedOptions[opt.type]" size="small">
                  <el-radio-button
                    v-for="v in opt.values"
                    :key="v.value"
                    :value="v.value"
                  >
                    {{ v.value }}
                    <span v-if="v.priceAdjustment" class="price-adjust">
                      ({{ v.priceAdjustment > 0 ? '+' : '' }}¥{{ formatPrice(v.priceAdjustment) }})
                    </span>
                  </el-radio-button>
                </el-radio-group>
              </div>
            </div>
            <el-empty v-else description="该产品暂无自定义选项" :image-size="80" />

            <div class="actions">
              <el-button
                type="primary"
                size="large"
                :icon="ChatDotRound"
                @click="handleInquire"
              >
                立即询价
              </el-button>
            </div>
          </el-card>
        </el-col>
      </el-row>

      <el-card shadow="never" class="desc-card">
        <template #header>产品描述</template>
        <div class="description">{{ product.description || '暂无描述' }}</div>
      </el-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { ArrowLeft, ChatDotRound } from '@element-plus/icons-vue'
import { getProduct, inquireProduct } from '@/api/product'
import type { Product, OptionType } from '@/types'

const route = useRoute()
const router = useRouter()

const loading = ref(false)
const product = ref<Product>({
  id: 0,
  name: '',
  sku: '',
  price: 0,
  stock: 0,
  status: 0,
  manufacturerId: 0,
  createdAt: '',
  updatedAt: ''
})

const selectedOptions = reactive<Record<OptionType, string>>({
  size: '',
  material: '',
  color: ''
})

const imageList = computed(() => {
  if (product.value.images && product.value.images.length > 0) {
    return product.value.images.map((i) => i.url)
  }
  return product.value.imageUrl ? [product.value.imageUrl] : []
})

async function fetchDetail() {
  loading.value = true
  try {
    const id = route.params.id as string
    const data = await getProduct(id)
    product.value = data
  } catch (err) {
    console.error(err)
  } finally {
    loading.value = false
  }
}

function optionTypeLabel(type: OptionType) {
  return { size: '尺寸', material: '材质', color: '颜色' }[type] || type
}

async function handleInquire() {
  const missing = product.value.options
    ?.filter((o) => o.required && !selectedOptions[o.type])
    .map((o) => optionTypeLabel(o.type))
  if (missing && missing.length > 0) {
    ElMessage.warning(`请选择${missing.join('、')}`)
    return
  }
  try {
    const result = await inquireProduct(product.value.id, {
      remark: JSON.stringify(selectedOptions)
    })
    ElMessage.success('已发起询价')
    if (result?.orderId) {
      router.push(`/orders/${result.orderId}`)
    } else {
      goBack()
    }
  } catch (err) {
    console.error(err)
  }
}

function goBack() {
  router.back()
}

function formatPrice(price: number) {
  return (price || 0).toFixed(2)
}

onMounted(fetchDetail)
</script>

<style lang="scss" scoped>
.product-detail {
  .detail-wrapper {
    margin-top: 16px;
  }

  .image-card {
    margin-bottom: 16px;
  }

  .info-card {
    margin-bottom: 16px;
  }

  .product-title {
    display: flex;
    align-items: center;
    gap: 8px;
    margin-bottom: 12px;

    h2 {
      margin: 0;
      font-size: 22px;
      color: #303133;
    }
  }

  .product-price {
    display: flex;
    align-items: baseline;
    gap: 12px;
    margin-bottom: 16px;

    .label {
      color: #909399;
      font-size: 14px;
    }

    .value {
      color: #f56c6c;
      font-size: 28px;
      font-weight: 600;
    }
  }

  .info-desc {
    margin-bottom: 16px;
  }

  .options {
    .option-group {
      margin-bottom: 14px;

      .option-label {
        font-size: 14px;
        color: #606266;
        margin-bottom: 6px;

        .required {
          color: #f56c6c;
          margin-left: 4px;
        }
      }

      .price-adjust {
        font-size: 12px;
        color: #f56c6c;
      }
    }
  }

  .actions {
    margin-top: 20px;
  }

  .desc-card {
    margin-top: 16px;

    .description {
      line-height: 1.7;
      color: #606266;
      white-space: pre-wrap;
    }
  }
}
</style>
