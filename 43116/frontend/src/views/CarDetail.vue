<template>
  <div class="container" style="padding-top: 20px;">
    <div v-if="car" class="car-detail">
      <div class="detail-header">
        <el-carousel :interval="4000" type="card" height="400px">
          <el-carousel-item v-for="(image, index) in car.images" :key="image.id">
            <img :src="image.url" class="carousel-image" />
          </el-carousel-item>
          <el-carousel-item v-if="!car.images || car.images.length === 0">
            <img :src="'https://picsum.photos/800/400?random=' + car.id" class="carousel-image" />
          </el-carousel-item>
        </el-carousel>
      </div>

      <div class="detail-body">
        <div class="detail-main">
          <h1 class="car-title">{{ car.brand }} {{ car.model }} {{ car.year }}</h1>
          <div class="car-meta">
            <el-tag type="primary" size="large">{{ car.seats }}座</el-tag>
            <el-tag size="large">{{ car.transmission }}</el-tag>
            <el-tag size="large">{{ car.fuel_type }}</el-tag>
            <el-tag :type="getStatusType(car.status)" size="large">
              {{ getStatusText(car.status) }}
            </el-tag>
          </div>
          <div class="price-info">
            <span class="price">¥{{ car.daily_rent }}</span>
            <span class="price-unit">/天</span>
            <span v-if="car.deposit > 0" class="deposit">押金: ¥{{ car.deposit }}</span>
          </div>

          <div class="detail-section">
            <div class="detail-label">车辆信息</div>
            <div class="detail-content">
              <p>车牌号: {{ car.license_plate || '暂无' }}</p>
              <p>颜色: {{ car.color || '暂无' }}</p>
              <p>里程: {{ car.mileage }} km</p>
              <p>评分: <el-rate :model-value="car.rating" disabled size="large" /> ({{ car.review_count }}条评价)</p>
            </div>
          </div>

          <div v-if="car.description" class="detail-section">
            <div class="detail-label">车辆描述</div>
            <div class="detail-content">{{ car.description }}</div>
          </div>

          <div v-if="car.features" class="detail-section">
            <div class="detail-label">车辆特点</div>
            <div class="detail-content">{{ car.features }}</div>
          </div>

          <div class="detail-section">
            <div class="detail-label">门店信息</div>
            <div class="detail-content" v-if="car.store">
              <p>{{ car.store.name }}</p>
              <p>{{ car.store.city?.name }} · {{ car.store.address }}</p>
              <p>电话: {{ car.store.phone }}</p>
              <p>营业时间: {{ car.store.business_hours }}</p>
            </div>
          </div>
        </div>

        <div class="detail-sidebar">
          <div class="booking-card">
            <h3>立即预订</h3>
            <el-form :model="bookingForm" label-position="top">
              <el-form-item label="取车门店">
                <el-select v-model="bookingForm.pickupStoreId" placeholder="选择门店">
                  <el-option
                    v-for="store in stores"
                    :key="store.id"
                    :label="store.name"
                    :value="store.id"
                  />
                </el-select>
              </el-form-item>
              <el-form-item label="还车门店">
                <el-select v-model="bookingForm.returnStoreId" placeholder="选择门店">
                  <el-option
                    v-for="store in stores"
                    :key="store.id"
                    :label="store.name"
                    :value="store.id"
                  />
                </el-select>
              </el-form-item>
              <el-form-item label="取车时间">
                <el-date-picker
                  v-model="bookingForm.pickupTime"
                  type="datetime"
                  placeholder="选择取车时间"
                  style="width: 100%"
                />
              </el-form-item>
              <el-form-item label="还车时间">
                <el-date-picker
                  v-model="bookingForm.returnTime"
                  type="datetime"
                  placeholder="选择还车时间"
                  style="width: 100%"
                />
              </el-form-item>
              <el-form-item label="优惠码">
                <el-input v-model="bookingForm.promoCode" placeholder="输入优惠码" />
              </el-form-item>
            </el-form>

            <div v-if="priceInfo" class="price-summary">
              <div class="price-row">
                <span>基础价格</span>
                <span>¥{{ priceInfo.base_price.toFixed(2) }}</span>
              </div>
              <div class="price-row" v-if="priceInfo.discount > 0">
                <span>优惠</span>
                <span class="discount">-¥{{ priceInfo.discount.toFixed(2) }}</span>
              </div>
              <div class="price-row total">
                <span>总计 ({{ priceInfo.total_days }}天)</span>
                <span>¥{{ priceInfo.final_price.toFixed(2) }}</span>
              </div>
            </div>

            <el-button
              type="primary"
              size="large"
              style="width: 100%; margin-top: 16px;"
              :disabled="car.status !== 'available'"
              :loading="booking"
              @click="handleBooking"
            >
              {{ car.status === 'available' ? '立即预订' : '暂不可预订' }}
            </el-button>
          </div>
        </div>
      </div>

      <div class="reviews-section">
        <h2>用户评价 ({{ car.review_count }}条)</h2>
        <review-list :car-id="car.id" />
      </div>
    </div>

    <el-empty v-if="!car && !loading" description="车辆不存在" />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { carApi, storeApi, bookingApi } from '@/api'
import type { Car, Store, PriceCalculation } from '@/types'
import ReviewList from '@/components/ReviewList.vue'

const route = useRoute()
const router = useRouter()

const car = ref<Car | null>(null)
const stores = ref<Store[]>([])
const loading = ref(false)
const booking = ref(false)
const priceInfo = ref<PriceCalculation | null>(null)

const bookingForm = reactive({
  pickupStoreId: undefined as number | undefined,
  returnStoreId: undefined as number | undefined,
  pickupTime: new Date(Date.now() + 24 * 60 * 60 * 1000),
  returnTime: new Date(Date.now() + 3 * 24 * 60 * 60 * 1000),
  promoCode: ''
})

onMounted(() => {
  loadCar()
})

const loadCar = async () => {
  loading.value = true
  try {
    const id = parseInt(route.params.id as string)
    const res = await carApi.getCarById(id)
    car.value = res.data

    if (car.value?.store) {
      bookingForm.pickupStoreId = car.value.store.id
      bookingForm.returnStoreId = car.value.store.id
    }

    loadStores()
    calculatePrice()
  } catch {
    // ignore
  } finally {
    loading.value = false
  }
}

const loadStores = async () => {
  if (!car.value?.store?.city_id) return
  try {
    const res = await storeApi.getStoresByCity(car.value.store.city_id)
    stores.value = res.data
  } catch {
    // ignore
  }
}

const calculatePrice = async () => {
  if (!car.value || !bookingForm.pickupTime || !bookingForm.returnTime) return
  try {
    const res = await bookingApi.calculatePrice({
      car_id: car.value.id,
      pickup_time: new Date(bookingForm.pickupTime).toISOString(),
      return_time: new Date(bookingForm.returnTime).toISOString(),
      promo_code: bookingForm.promoCode
    })
    priceInfo.value = res.data
  } catch {
    priceInfo.value = null
  }
}

watch(() => [bookingForm.pickupTime, bookingForm.returnTime, bookingForm.promoCode], calculatePrice)

const handleBooking = async () => {
  if (!car.value) return
  if (!bookingForm.pickupStoreId || !bookingForm.returnStoreId) {
    ElMessage.error('请选择取车和还车门店')
    return
  }

  booking.value = true
  try {
    const res = await bookingApi.createBooking({
      car_id: car.value.id,
      pickup_store_id: bookingForm.pickupStoreId!,
      return_store_id: bookingForm.returnStoreId!,
      pickup_time: new Date(bookingForm.pickupTime).toISOString(),
      return_time: new Date(bookingForm.returnTime).toISOString(),
      promo_code: bookingForm.promoCode || undefined
    })
    ElMessage.success('预订成功！')
    router.push(`/my-bookings`)
  } catch (err: any) {
    ElMessage.error(err.message || '预订失败')
  } finally {
    booking.value = false
  }
}

const getStatusType = (status: string) => {
  const map: Record<string, string> = {
    available: 'success',
    rented: 'primary',
    maintenance: 'warning',
    disabled: 'danger'
  }
  return map[status] || 'info'
}

const getStatusText = (status: string) => {
  const map: Record<string, string> = {
    available: '可用',
    rented: '出租中',
    maintenance: '维护中',
    disabled: '已停用'
  }
  return map[status] || status
}
</script>

<style scoped>
.car-detail {
  background: #fff;
  border-radius: 8px;
  overflow: hidden;
}

.detail-header {
  width: 100%;
}

.carousel-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.detail-body {
  display: flex;
  gap: 30px;
  padding: 30px;
}

.detail-main {
  flex: 1;
}

.car-title {
  font-size: 28px;
  font-weight: 600;
  margin-bottom: 16px;
}

.car-meta {
  display: flex;
  gap: 10px;
  margin-bottom: 20px;
}

.price-info {
  margin-bottom: 20px;
}

.price {
  font-size: 32px;
  font-weight: 700;
  color: #409eff;
}

.price-unit {
  font-size: 16px;
  color: #909399;
}

.deposit {
  margin-left: 20px;
  color: #909399;
}

.detail-section {
  margin-bottom: 24px;
}

.detail-label {
  font-weight: 600;
  color: #606266;
  margin-bottom: 8px;
  font-size: 16px;
}

.detail-content {
  color: #303133;
  line-height: 1.8;
}

.detail-sidebar {
  width: 350px;
}

.booking-card {
  background: #f5f7fa;
  border-radius: 8px;
  padding: 24px;
  position: sticky;
  top: 20px;
}

.booking-card h3 {
  margin-bottom: 20px;
  font-size: 18px;
}

.price-summary {
  background: #fff;
  border-radius: 6px;
  padding: 16px;
  margin-top: 16px;
}

.price-row {
  display: flex;
  justify-content: space-between;
  padding: 8px 0;
}

.price-row.total {
  border-top: 1px solid #ebeef5;
  margin-top: 8px;
  padding-top: 12px;
  font-weight: 600;
  font-size: 16px;
}

.discount {
  color: #67c23a;
}

.reviews-section {
  padding: 30px;
  border-top: 1px solid #ebeef5;
}

.reviews-section h2 {
  margin-bottom: 20px;
}
</style>
