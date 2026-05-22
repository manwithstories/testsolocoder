<template>
  <div class="property-list">
    <div class="page-header">
      <h2 class="page-title">房源管理</h2>
      <el-button type="primary" @click="goToCreate">
        <el-icon><Plus /></el-icon>发布房源
      </el-button>
    </div>

    <div class="search-bar card">
      <el-input
        v-model="searchForm.keyword"
        placeholder="搜索房源名称/小区"
        clearable
        style="width: 200px"
        @keyup.enter="handleSearch"
      />
      <el-select v-model="searchForm.region" placeholder="区域" clearable style="width: 120px">
        <el-option label="浦东" value="浦东" />
        <el-option label="徐汇" value="徐汇" />
        <el-option label="静安" value="静安" />
        <el-option label="长宁" value="长宁" />
      </el-select>
      <el-select v-model="searchForm.layout" placeholder="户型" clearable style="width: 120px">
        <el-option label="一室" value="一室" />
        <el-option label="两室" value="两室" />
        <el-option label="三室" value="三室" />
        <el-option label="四室" value="四室" />
      </el-select>
      <el-input v-model="searchForm.minRent" placeholder="最低租金" clearable style="width: 120px" />
      <el-input v-model="searchForm.maxRent" placeholder="最高租金" clearable style="width: 120px" />
      <el-select v-model="searchForm.status" placeholder="状态" clearable style="width: 120px">
        <el-option label="上架" :value="1" />
        <el-option label="已租" :value="2" />
        <el-option label="下架" :value="0" />
      </el-select>
      <el-button type="primary" @click="handleSearch">搜索</el-button>
      <el-button @click="handleReset">重置</el-button>
    </div>

    <div class="card">
      <el-table :data="properties" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column label="房源信息" min-width="200">
          <template #default="{ row }">
            <div class="property-info">
              <el-image
                :src="getImageUrl(row)"
                style="width: 60px; height: 60px; border-radius: 4px"
                fit="cover"
              />
              <div>
                <div class="title">{{ row.title }}</div>
                <div class="sub">{{ row.community }} · {{ row.layout }}</div>
              </div>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="area" label="面积" width="100">
          <template #default="{ row }">{{ row.area }}㎡</template>
        </el-table-column>
        <el-table-column prop="rent" label="租金" width="120">
          <template #default="{ row }">¥{{ row.rent }}/月</template>
        </el-table-column>
        <el-table-column prop="region" label="区域" width="100" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="goToDetail(row.id)">详情</el-button>
            <el-button link type="primary" @click="goToEdit(row.id)">编辑</el-button>
            <el-button
              link
              :type="row.status === 1 ? 'warning' : 'success'"
              @click="toggleStatus(row)"
            >
              {{ row.status === 1 ? '下架' : '上架' }}
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        class="pagination"
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.pageSize"
        :total="pagination.total"
        :page-sizes="[10, 20, 50]"
        layout="total, sizes, prev, pager, next, jumper"
        @current-change="loadData"
        @size-change="loadData"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { Property } from '@/types'
import { getProperties, updatePropertyStatus } from '@/api/property'

const router = useRouter()

const loading = ref(false)
const properties = ref<Property[]>([])

const searchForm = reactive({
  keyword: '',
  region: '',
  layout: '',
  minRent: '',
  maxRent: '',
  status: ''
})

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

onMounted(() => {
  loadData()
})

async function loadData() {
  loading.value = true
  try {
    const params: any = {
      page: pagination.page,
      pageSize: pagination.pageSize
    }
    if (searchForm.keyword) params.keyword = searchForm.keyword
    if (searchForm.region) params.region = searchForm.region
    if (searchForm.layout) params.layout = searchForm.layout
    if (searchForm.minRent) params.minRent = parseFloat(searchForm.minRent)
    if (searchForm.maxRent) params.maxRent = parseFloat(searchForm.maxRent)
    if (searchForm.status !== '') params.status = searchForm.status

    const res = await getProperties(params)
    properties.value = res.data.list
    pagination.total = res.data.total
  } catch (error) {
    console.error('Failed to load properties:', error)
  } finally {
    loading.value = false
  }
}

function handleSearch() {
  pagination.page = 1
  loadData()
}

function handleReset() {
  searchForm.keyword = ''
  searchForm.region = ''
  searchForm.layout = ''
  searchForm.minRent = ''
  searchForm.maxRent = ''
  searchForm.status = ''
  handleSearch()
}

function getImageUrl(row: Property) {
  if (row.images && row.images.length > 0) {
    return row.images[0].url
  }
  return 'https://picsum.photos/100/100'
}

function getStatusType(status: number) {
  switch (status) {
    case 1: return 'success'
    case 2: return 'warning'
    default: return 'info'
  }
}

function getStatusText(status: number) {
  switch (status) {
    case 1: return '上架'
    case 2: return '已租'
    default: return '下架'
  }
}

function goToCreate() {
  router.push('/properties/create')
}

function goToEdit(id: number) {
  router.push(`/properties/${id}/edit`)
}

function goToDetail(id: number) {
  router.push(`/properties/${id}`)
}

async function toggleStatus(row: Property) {
  const newStatus = row.status === 1 ? 0 : 1
  const action = newStatus === 1 ? '上架' : '下架'
  
  try {
    await ElMessageBox.confirm(`确定要${action}该房源吗？`, '提示', {
      type: 'warning'
    })
    
    await updatePropertyStatus(row.id, newStatus)
    ElMessage.success(`${action}成功`)
    loadData()
  } catch (error) {
    if (error !== 'cancel') {
      console.error(error)
    }
  }
}
</script>

<style scoped>
.property-list {
  padding: 0;
}

.search-bar {
  display: flex;
  gap: 10px;
  margin-bottom: 20px;
  flex-wrap: wrap;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

.property-info {
  display: flex;
  gap: 10px;
  align-items: center;
}

.property-info .title {
  font-weight: 500;
  margin-bottom: 4px;
}

.property-info .sub {
  color: #909399;
  font-size: 12px;
}
</style>
