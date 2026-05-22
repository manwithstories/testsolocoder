<template>
  <div class="contract-list">
    <div class="page-header">
      <h2 class="page-title">合同管理</h2>
      <el-button type="primary" @click="goToCreate">
        <el-icon><Plus /></el-icon>创建合同
      </el-button>
    </div>

    <div class="search-bar card">
      <el-select v-model="statusFilter" placeholder="状态" clearable style="width: 150px">
        <el-option label="生效中" :value="1" />
        <el-option label="已到期" :value="2" />
        <el-option label="已终止" :value="0" />
      </el-select>
      <el-button type="primary" @click="loadData">搜索</el-button>
    </div>

    <div class="card">
      <el-table :data="contracts" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column label="租户" width="150">
          <template #default="{ row }">
            {{ row.tenant?.name }} ({{ row.tenant?.phone }})
          </template>
        </el-table-column>
        <el-table-column label="房源" min-width="150">
          <template #default="{ row }">
            {{ row.property?.title }}
          </template>
        </el-table-column>
        <el-table-column label="租期" width="220">
          <template #default="{ row }">
            {{ formatDate(row.startDate) }} 至 {{ formatDate(row.endDate) }}
          </template>
        </el-table-column>
        <el-table-column prop="rent" label="租金" width="120">
          <template #default="{ row }">¥{{ row.rent }}/月</template>
        </el-table-column>
        <el-table-column prop="deposit" label="押金" width="100">
          <template #default="{ row }">¥{{ row.deposit }}</template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">{{ getStatusText(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="viewDetail(row.id)">详情</el-button>
            <el-button
              link
              type="warning"
              @click="terminateContract(row)"
              v-if="row.status === 1"
            >终止</el-button>
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
import type { Contract } from '@/types'
import { getContracts, updateContractStatus } from '@/api/tenant'
import dayjs from 'dayjs'

const router = useRouter()
const loading = ref(false)
const contracts = ref<Contract[]>([])
const statusFilter = ref<number | ''>('')

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
    const res = await getContracts({
      page: pagination.page,
      pageSize: pagination.pageSize,
      status: statusFilter.value || undefined
    })
    contracts.value = res.data.list
    pagination.total = res.data.total
  } catch (error) {
    console.error('Failed to load contracts:', error)
  } finally {
    loading.value = false
  }
}

function formatDate(date: string) {
  return dayjs(date).format('YYYY-MM-DD')
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
    case 1: return '生效中'
    case 2: return '已到期'
    default: return '已终止'
  }
}

function goToCreate() {
  router.push('/contracts/create')
}

function viewDetail(id: number) {
  router.push(`/contracts/${id}`)
}

async function terminateContract(row: Contract) {
  try {
    await ElMessageBox.confirm('确定要终止该合同吗？', '提示', { type: 'warning' })
    await updateContractStatus(row.id, 0)
    ElMessage.success('合同已终止')
    loadData()
  } catch (error) {
    if (error !== 'cancel') {
      console.error(error)
    }
  }
}
</script>

<style scoped>
.contract-list {
  padding: 0;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}
</style>
