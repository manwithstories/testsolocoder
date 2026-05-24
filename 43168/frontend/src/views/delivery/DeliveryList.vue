<template>
  <div class="delivery-list">
    <el-card shadow="never" class="search-card">
      <el-form :inline="true" :model="searchForm" @submit.prevent>
        <el-form-item label="订单号">
          <el-input v-model="searchForm.orderNo" placeholder="请输入订单号" clearable style="width: 200px" />
        </el-form-item>
        <el-form-item label="联系电话">
          <el-input v-model="searchForm.contactPhone" placeholder="请输入联系电话" clearable style="width: 200px" />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="请选择状态" clearable style="width: 200px">
            <el-option v-for="item in statusOptions" :key="item.value" :label="item.label" :value="item.value" />
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
        <div class="card-header">
          <span>配送安装列表</span>
          <el-button type="primary" :icon="Plus" @click="handleCreate">预约安装</el-button>
        </div>
      </template>

      <el-table :data="list" v-loading="loading" stripe>
        <el-table-column prop="orderNo" label="订单号" width="180" />
        <el-table-column prop="contactName" label="联系人" width="100" />
        <el-table-column prop="contactPhone" label="联系电话" width="140" />
        <el-table-column prop="address" label="地址" min-width="200" show-overflow-tooltip />
        <el-table-column prop="deliveryDate" label="配送日期" width="120" />
        <el-table-column prop="timeSlot" label="时间段" width="120" />
        <el-table-column prop="installer" label="安装师傅" width="120" show-overflow-tooltip />
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">{{ getStatusLabel(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="handleEdit(row)">编辑</el-button>
            <el-popconfirm title="确认删除？" @confirm="handleDelete(row)">
              <template #reference>
                <el-button link type="danger">删除</el-button>
              </template>
            </el-popconfirm>
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
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Search, Refresh, Plus } from '@element-plus/icons-vue'
import { listDeliveries, deleteDelivery, type Delivery, type ListDeliveryParams } from '@/api/delivery'

const router = useRouter()
const loading = ref(false)
const list = ref<Delivery[]>([])

const searchForm = reactive<ListDeliveryParams>({
  page: 1,
  pageSize: 10,
  orderNo: '',
  contactPhone: '',
  status: undefined as unknown as number
})

const pagination = reactive({ page: 1, pageSize: 10, total: 0 })

const statusOptions = [
  { label: '待预约', value: 0 },
  { label: '已预约', value: 1 },
  { label: '配送中', value: 2 },
  { label: '安装中', value: 3 },
  { label: '已完成', value: 4 },
  { label: '已取消', value: 5 }
]

function getStatusLabel(status: number) {
  return statusOptions.find((i) => i.value === status)?.label ?? '未知'
}
function getStatusType(status: number) {
  const map: Record<number, string> = {
    0: 'info',
    1: 'warning',
    2: 'primary',
    3: 'primary',
    4: 'success',
    5: 'danger'
  }
  return map[status] ?? ''
}

async function fetchList() {
  loading.value = true
  try {
    const res = await listDeliveries({
      orderNo: searchForm.orderNo || undefined,
      contactPhone: searchForm.contactPhone || undefined,
      status: searchForm.status ?? undefined,
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
  searchForm.contactPhone = ''
  searchForm.status = undefined as unknown as number
  pagination.page = 1
  fetchList()
}

function handleCreate() {
  router.push('/deliveries/new')
}

function handleEdit(row: Delivery) {
  router.push(`/deliveries/${row.id}/edit`)
}

async function handleDelete(row: Delivery) {
  await deleteDelivery(row.id)
  ElMessage.success('删除成功')
  fetchList()
}

onMounted(fetchList)
</script>

<style lang="scss" scoped>
.delivery-list {
  .search-card {
    margin-bottom: 16px;
  }
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }
  .pagination-wrapper {
    display: flex;
    justify-content: flex-end;
    padding-top: 16px;
  }
}
</style>
