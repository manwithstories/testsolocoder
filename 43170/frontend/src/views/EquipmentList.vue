<template>
  <div class="page-container">
    <h1 class="page-title">设备列表</h1>

    <div class="card search-filter">
      <el-form :model="searchParams" inline>
        <el-form-item label="分类">
          <el-select v-model="searchParams.category" placeholder="全部分类" clearable style="width: 150px">
            <el-option
              v-for="cat in categories"
              :key="cat"
              :label="cat"
              :value="cat"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="品牌">
          <el-select v-model="searchParams.brand" placeholder="全部品牌" clearable style="width: 150px">
            <el-option
              v-for="brand in brands"
              :key="brand"
              :label="brand"
              :value="brand"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="价格区间">
          <el-input-number
            v-model="searchParams.minPrice"
            :min="0"
            placeholder="最低价"
            style="width: 120px"
          />
          <span style="margin: 0 10px">-</span>
          <el-input-number
            v-model="searchParams.maxPrice"
            :min="0"
            placeholder="最高价"
            style="width: 120px"
          />
        </el-form-item>
        <el-form-item label="开始日期">
          <el-date-picker
            v-model="searchParams.startDate"
            type="date"
            placeholder="选择日期"
            style="width: 150px"
            value-format="YYYY-MM-DD"
          />
        </el-form-item>
        <el-form-item label="结束日期">
          <el-date-picker
            v-model="searchParams.endDate"
            type="date"
            placeholder="选择日期"
            style="width: 150px"
            value-format="YYYY-MM-DD"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <div v-loading="loading" class="equipment-list">
      <el-row :gutter="20">
        <el-col
          v-for="equipment in equipments"
          :key="equipment.id"
          :xs="24"
          :sm="12"
          :md="8"
          :lg="6"
        >
          <div class="card equipment-card" @click="goToDetail(equipment.id)">
            <el-carousel
              v-if="equipment.images && equipment.images.length > 0"
              height="200px"
              :indicator-position="'none'"
            >
              <el-carousel-item v-for="img in equipment.images" :key="img.id">
                <img
                  :src="`/uploads/${img.imageUrl}`"
                  :alt="equipment.name"
                  class="equipment-image"
                />
              </el-carousel-item>
            </el-carousel>
            <div v-else class="equipment-image-placeholder">
              <el-icon :size="48"><Camera /></el-icon>
            </div>
            <div class="equipment-info">
              <div class="equipment-name">{{ equipment.name }}</div>
              <div class="equipment-brand">{{ equipment.brand }} {{ equipment.model }}</div>
              <div class="equipment-rating">
                <el-rate :model-value="equipment.rating" disabled size="small" />
                <span>({{ equipment.reviewCount }})</span>
              </div>
              <div class="equipment-price">¥{{ equipment.dailyRent }}/天</div>
              <div class="equipment-deposit">押金: ¥{{ equipment.deposit }}</div>
              <div class="equipment-status">
                <el-tag
                  :type="equipment.status === 'available' ? 'success' : 'info'"
                  size="small"
                >
                  {{ getStatusText(equipment.status) }}
                </el-tag>
              </div>
            </div>
          </div>
        </el-col>
      </el-row>

      <div v-if="!loading && equipments.length === 0" class="empty-state">
        <el-icon><Box /></el-icon>
        <p>暂无符合条件的设备</p>
      </div>
    </div>

    <div v-if="total > 0" class="pagination">
      <el-pagination
        v-model:current-page="searchParams.page"
        v-model:page-size="searchParams.pageSize"
        :page-sizes="[8, 16, 24, 32]"
        :total="total"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="handleSearch"
        @current-change="handleSearch"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { searchApi, equipmentApi } from '@/api/equipment'
import type { Equipment, SearchRequest } from '@/types'

const router = useRouter()
const loading = ref(false)
const equipments = ref<Equipment[]>([])
const categories = ref<string[]>([])
const brands = ref<string[]>([])
const total = ref(0)

const searchParams = reactive<SearchRequest>({
  category: '',
  brand: '',
  minPrice: undefined,
  maxPrice: undefined,
  startDate: '',
  endDate: '',
  page: 1,
  pageSize: 8,
  sortBy: 'created_at',
  sortOrder: 'desc'
})

onMounted(async () => {
  loadCategories()
  loadBrands()
  handleSearch()
})

async function loadCategories() {
  try {
    const response = await equipmentApi.getCategories()
    categories.value = response.data
  } catch (error) {
    console.error('Failed to load categories:', error)
  }
}

async function loadBrands() {
  try {
    const response = await equipmentApi.getBrands()
    brands.value = response.data
  } catch (error) {
    console.error('Failed to load brands:', error)
  }
}

async function handleSearch() {
  loading.value = true
  try {
    const response = await searchApi.searchEquipments(searchParams)
    equipments.value = response.data
    total.value = response.total
  } catch (error) {
    console.error('Search failed:', error)
  } finally {
    loading.value = false
  }
}

function handleReset() {
  searchParams.category = ''
  searchParams.brand = ''
  searchParams.minPrice = undefined
  searchParams.maxPrice = undefined
  searchParams.startDate = ''
  searchParams.endDate = ''
  searchParams.page = 1
  handleSearch()
}

function goToDetail(id: number) {
  router.push(`/equipments/${id}`)
}

function getStatusText(status: string) {
  const statusMap: Record<string, string> = {
    available: '可出租',
    rented: '已出租',
    maintenance: '维护中'
  }
  return statusMap[status] || status
}
</script>

<style scoped>
.equipment-image {
  width: 100%;
  height: 200px;
  object-fit: cover;
}

.equipment-image-placeholder {
  width: 100%;
  height: 200px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f5f7fa;
  color: #c0c4cc;
}

.equipment-brand {
  font-size: 14px;
  color: #909399;
  margin-bottom: 8px;
}

.equipment-deposit {
  font-size: 12px;
  color: #909399;
  margin-top: 4px;
}

.equipment-status {
  margin-top: 8px;
}
</style>
