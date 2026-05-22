<template>
  <div class="container" style="padding-top: 20px;">
    <div class="page-header">
      <h1 class="page-title">车辆列表</h1>
    </div>

    <div class="filter-bar">
      <el-select v-model="filters.brand" placeholder="品牌" clearable style="width: 150px">
        <el-option v-for="brand in brands" :key="brand" :label="brand" :value="brand" />
      </el-select>
      <el-select v-model="filters.status" placeholder="状态" clearable style="width: 120px">
        <el-option label="可用" value="available" />
        <el-option label="出租中" value="rented" />
        <el-option label="维护中" value="maintenance" />
        <el-option label="已停用" value="disabled" />
      </el-select>
      <el-input
        v-model="filters.keyword"
        placeholder="搜索品牌/型号/车牌号"
        clearable
        style="width: 250px"
        @keyup.enter="loadCars"
      >
        <template #prefix>
          <el-icon><Search /></el-icon>
        </template>
      </el-input>
      <el-button type="primary" @click="loadCars">搜索</el-button>
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
            <div>
              <span class="car-price">¥{{ car.daily_rent }}</span>
              <span class="car-price-unit">/天</span>
            </div>
            <el-tag :type="getStatusType(car.status)" size="small">
              {{ getStatusText(car.status) }}
            </el-tag>
          </div>
        </div>
      </div>
    </div>

    <div style="text-align: center; margin-top: 30px;">
      <el-pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.pageSize"
        :total="pagination.total"
        :page-sizes="[12, 24, 48]"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="loadCars"
        @current-change="loadCars"
      />
    </div>

    <el-empty v-if="cars.length === 0 && !loading" description="暂无车辆" />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { Search } from '@element-plus/icons-vue'
import { carApi } from '@/api'
import type { Car } from '@/types'

const route = useRoute()

const cars = ref<Car[]>([])
const loading = ref(false)
const brands = ref<string[]>(['奥迪', '宝马', '奔驰', '丰田', '本田', '大众', '特斯拉', '比亚迪'])

const filters = reactive({
  keyword: route.query.keyword as string || '',
  brand: '',
  status: ''
})

const pagination = reactive({
  page: 1,
  pageSize: 12,
  total: 0
})

onMounted(() => {
  loadCars()
})

const loadCars = async () => {
  loading.value = true
  try {
    const res = await carApi.getCars({
      page: pagination.page,
      page_size: pagination.pageSize,
      keyword: filters.keyword,
      brand: filters.brand,
      status: filters.status
    })
    cars.value = res.data.items
    pagination.total = res.data.total
  } catch {
    // ignore
  } finally {
    loading.value = false
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
