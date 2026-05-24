<template>
  <el-card>
    <template #header>
      <div class="header">
        <span>设备列表</span>
        <el-button type="primary" @click="showSearch = !showSearch">
          <el-icon><Search /></el-icon>
          搜索
        </el-button>
      </div>
    </template>

    <el-form :inline="true" v-show="showSearch" @submit.prevent="handleSearch">
      <el-form-item label="关键字">
        <el-input v-model="search.keyword" placeholder="设备名称/品牌" clearable />
      </el-form-item>
      <el-form-item label="区域">
        <el-input v-model="search.region" placeholder="区域" clearable />
      </el-form-item>
      <el-form-item label="品牌">
        <el-input v-model="search.brand" placeholder="品牌" clearable />
      </el-form-item>
      <el-form-item label="租赁时段">
        <el-date-picker
          v-model="search.dateRange"
          type="daterange"
          range-separator="至"
          start-placeholder="开始日期"
          end-placeholder="结束日期"
          value-format="YYYY-MM-DD"
        />
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="handleSearch">搜索</el-button>
        <el-button @click="resetSearch">重置</el-button>
      </el-form-item>
    </el-form>

    <el-table :data="drones" v-loading="loading">
      <el-table-column prop="name" label="设备名称" min-width="150" />
      <el-table-column label="品牌/型号" min-width="150">
        <template #default="{ row }">
          {{ row.brand }} {{ row.model }}
        </template>
      </el-table-column>
      <el-table-column prop="region" label="区域" width="100" />
      <el-table-column prop="price_per_day" label="日租金" width="100">
        <template #default="{ row }">¥{{ row.price_per_day }}</template>
      </el-table-column>
      <el-table-column label="评分" width="120">
        <template #default="{ row }">
          <el-rate :model-value="row.rating" disabled size="small" />
        </template>
      </el-table-column>
      <el-table-column label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="statusTagType(row.status)">{{ statusText(row.status) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="120" fixed="right">
        <template #default="{ row }">
          <el-button type="primary" link @click="$router.push(`/drone/${row.id}`)">详情</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-pagination
      style="margin-top: 20px; justify-content: flex-end"
      layout="prev, pager, next, total"
      :total="total"
      :page-size="pageSize"
      v-model:current-page="currentPage"
      @current-change="fetchDrones"
    />
  </el-card>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import request from '@/utils/request'

const loading = ref(false)
const showSearch = ref(false)
const drones = ref<Drone[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)

const search = reactive({
  keyword: '',
  region: '',
  brand: '',
  dateRange: [] as string[]
})

onMounted(() => {
  fetchDrones()
})

async function fetchDrones() {
  loading.value = true
  try {
    const params: any = {
      page: currentPage.value,
      page_size: pageSize.value,
      keyword: search.keyword,
      region: search.region
    }
    
    if (search.dateRange && search.dateRange.length === 2) {
      params.start_date = search.dateRange[0]
      params.end_date = search.dateRange[1]
      const res: any = await request.get('/drones/search', { params })
      drones.value = res.data.list || []
      total.value = res.data.total || 0
    } else {
      params.status = 'online'
      if (search.brand) {
        params.brand = search.brand
      }
      const res: any = await request.get('/drones', { params })
      drones.value = res.data.list || []
      total.value = res.data.total || 0
    }
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

function handleSearch() {
  currentPage.value = 1
  fetchDrones()
}

function resetSearch() {
  search.keyword = ''
  search.region = ''
  search.brand = ''
  search.dateRange = []
  currentPage.value = 1
  fetchDrones()
}

function statusText(status: string) {
  const map: Record<string, string> = {
    offline: '已下架', online: '可租', rented: '已租', maintenance: '维护中'
  }
  return map[status] || status
}

function statusTagType(status: string) {
  const map: Record<string, string> = {
    offline: 'info', online: 'success', rented: 'warning', maintenance: 'danger'
  }
  return map[status] || ''
}
</script>

<style scoped>
.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
