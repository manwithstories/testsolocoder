<template>
  <div class="home">
    <div class="banner">
      <h1>探索自由，畅享出行</h1>
      <p>精选车辆，灵活租赁，随时随地启程</p>
    </div>

    <div class="container">
      <div class="search-section">
        <el-form :inline="true" :model="searchForm" class="search-form">
          <el-form-item label="取车城市">
            <el-select v-model="searchForm.cityId" placeholder="选择城市" style="width: 150px">
              <el-option
                v-for="city in cities"
                :key="city.id"
                :label="city.name"
                :value="city.id"
              />
            </el-select>
          </el-form-item>
          <el-form-item label="取车门店">
            <el-select v-model="searchForm.storeId" placeholder="选择门店" style="width: 200px">
              <el-option
                v-for="store in stores"
                :key="store.id"
                :label="store.name"
                :value="store.id"
              />
            </el-select>
          </el-form-item>
          <el-form-item label="取车日期">
            <el-date-picker
              v-model="searchForm.pickupDate"
              type="datetime"
              placeholder="选择取车时间"
              style="width: 200px"
            />
          </el-form-item>
          <el-form-item label="还车日期">
            <el-date-picker
              v-model="searchForm.returnDate"
              type="datetime"
              placeholder="选择还车时间"
              style="width: 200px"
            />
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="searchCars">搜索车辆</el-button>
          </el-form-item>
        </el-form>
      </div>

      <div class="section-title">
        <h2>热门车辆</h2>
      </div>

      <div class="card-grid">
        <div
          v-for="car in cars"
          :key="car.id"
          class="car-card"
          @click="$router.push(`/cars/${car.id}`)"
        >
          <el-image
            class="car-image"
            :src="car.images?.[0]?.url || 'https://picsum.photos/400/300?random=' + car.id"
            fit="cover"
          />
          <div class="car-info">
            <div class="car-title">{{ car.brand }} {{ car.model }}</div>
            <div class="car-meta">
              <span>{{ car.seats }}座</span>
              <span>{{ car.transmission }}</span>
              <span>{{ car.fuel_type }}</span>
            </div>
            <div style="margin-top: 12px; display: flex; justify-content: space-between; align-items: center;">
              <div class="car-price">
                ¥{{ car.daily_rent }}<span class="car-price-unit">/天</span>
              </div>
              <el-rate :model-value="car.rating" disabled size="small" />
            </div>
          </div>
        </div>
      </div>

      <el-empty v-if="cars.length === 0 && !loading" description="暂无车辆" />
      <div v-if="loading" style="text-align: center; padding: 40px;">
        <el-icon class="is-loading" :size="32"><Loading /></el-icon>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { Loading } from '@element-plus/icons-vue'
import { carApi, storeApi } from '@/api'
import type { Car, City, Store } from '@/types'

const router = useRouter()
const route = useRoute()

const cars = ref<Car[]>([])
const cities = ref<City[]>([])
const stores = ref<Store[]>([])
const loading = ref(false)

const searchForm = reactive({
  cityId: undefined as number | undefined,
  storeId: undefined as number | undefined,
  pickupDate: new Date(),
  returnDate: new Date(Date.now() + 24 * 60 * 60 * 1000)
})

onMounted(() => {
  loadCities()
  loadCars()

  if (route.query.keyword) {
    loadCars(route.query.keyword as string)
  }
})

const loadCities = async () => {
  try {
    const res = await storeApi.getCities()
    cities.value = res.data
  } catch {
    // ignore
  }
}

const loadCars = async (keyword?: string) => {
  loading.value = true
  try {
    const res = await carApi.getAvailableCars({ page: 1, page_size: 20 })
    cars.value = res.data.items
  } catch {
    // ignore
  } finally {
    loading.value = false
  }
}

const searchCars = () => {
  router.push({
    path: '/cars',
    query: {
      store_id: searchForm.storeId,
      keyword: ''
    }
  })
}
</script>

<style scoped>
.home {
  min-height: calc(100vh - 120px);
}

.banner {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: #fff;
  text-align: center;
  padding: 60px 20px;
  margin-bottom: 30px;
}

.banner h1 {
  font-size: 36px;
  margin-bottom: 12px;
}

.banner p {
  font-size: 18px;
  opacity: 0.9;
}

.search-section {
  background: #fff;
  padding: 24px;
  border-radius: 8px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
  margin-bottom: 30px;
}

.search-form {
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
}

.section-title {
  margin-bottom: 20px;
}

.section-title h2 {
  font-size: 24px;
  color: #303133;
}
</style>
