<template>
  <div class="production-list">
    <el-card shadow="never" class="search-card">
      <el-form :inline="true" :model="searchForm" @submit.prevent>
        <el-form-item label="订单号">
          <el-input v-model="searchForm.orderNo" placeholder="请输入订单号" clearable style="width: 200px" />
        </el-form-item>
        <el-form-item label="生产阶段">
          <el-select v-model="searchForm.currentStage" placeholder="请选择阶段" clearable style="width: 200px">
            <el-option v-for="item in stageOptions" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :icon="Search" @click="handleSearch">搜索</el-button>
          <el-button :icon="Refresh" @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card shadow="never">
      <template #header>
        <span>生产进度列表</span>
      </template>

      <el-table :data="list" v-loading="loading" stripe>
        <el-table-column prop="id" label="编号" width="80" />
        <el-table-column prop="orderNo" label="订单号" width="180" />
        <el-table-column prop="productName" label="产品" min-width="160" show-overflow-tooltip />
        <el-table-column prop="quantity" label="数量" width="80" />
        <el-table-column label="当前阶段" width="120">
          <template #default="{ row }">
            <el-tag :type="getStageType(row.currentStage)">{{ getStageLabel(row.currentStage) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="进度" min-width="200">
          <template #default="{ row }">
            <el-progress :percentage="row.progress" :status="getProgressStatus(row)" />
          </template>
        </el-table-column>
        <el-table-column prop="estimatedDate" label="预计完成" width="120" />
        <el-table-column label="操作" width="160" fixed="right">
          <template #default="{ row }">
            <el-dropdown @command="(val) => handleUpdateStage(row, val)" :disabled="row.currentStage >= 5">
              <el-button type="primary" link>
                更新阶段<el-icon class="el-icon--right"><ArrowDown /></el-icon>
              </el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item v-for="s in stageOptions" :key="s.value" :value="s.value" :disabled="s.value <= row.currentStage">
                    {{ s.label }}
                  </el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-wrapper">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="pagination.total"
          layout="total, sizes, prev, pager, next, jumper"
          background
          @size-change="fetchList"
          @current-change="fetchList"
        />
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Search, Refresh, ArrowDown } from '@element-plus/icons-vue'
import {
  listProduction,
  updateProductionStage,
  type ProductionProgress,
  type ListProductionParams
} from '@/api/production'

const loading = ref(false)
const list = ref<ProductionProgress[]>([])

const searchForm = reactive<ListProductionParams>({
  page: 1,
  pageSize: 10,
  orderNo: '',
  currentStage: undefined as unknown as number
})

const pagination = reactive({ page: 1, pageSize: 10, total: 0 })

const stageOptions = [
  { label: '待生产', value: 0 },
  { label: '备料中', value: 1 },
  { label: '生产中', value: 2 },
  { label: '质检中', value: 3 },
  { label: '包装中', value: 4 },
  { label: '已完成', value: 5 }
]

function getStageLabel(stage: number) {
  return stageOptions.find((i) => i.value === stage)?.label ?? '未知'
}
function getStageType(stage: number) {
  const map: Record<number, string> = {
    0: 'info',
    1: 'warning',
    2: 'primary',
    3: '',
    4: '',
    5: 'success'
  }
  return map[stage] ?? ''
}
function getProgressStatus(row: ProductionProgress) {
  if (row.progress >= 100) return 'success'
  if (row.progress >= 60) return ''
  return 'warning'
}

async function fetchList() {
  loading.value = true
  try {
    const res = await listProduction({
      orderNo: searchForm.orderNo || undefined,
      currentStage: searchForm.currentStage ?? undefined,
      page: pagination.page,
      pageSize: pagination.pageSize
    })
    list.value = (res as any)?.list ?? []
    pagination.total = (res as any)?.total ?? 0
  } finally {
    loading.value = false
  }
}

function handleSearch() {
  pagination.page = 1
  fetchList()
}

function handleReset() {
  searchForm.orderNo = ''
  searchForm.currentStage = undefined as unknown as number
  pagination.page = 1
  fetchList()
}

async function handleUpdateStage(row: ProductionProgress, stage: number) {
  await updateProductionStage(row.id, stage)
  ElMessage.success('阶段更新成功')
  fetchList()
}

onMounted(fetchList)
</script>

<style lang="scss" scoped>
.production-list {
  .search-card {
    margin-bottom: 16px;
  }
  .pagination-wrapper {
    display: flex;
    justify-content: flex-end;
    padding-top: 16px;
  }
}
</style>
