<template>
  <div class="order-list">
    <SearchBar
      :fields="searchFields"
      v-model="searchForm"
      @search="handleSearch"
      @reset="handleReset"
    />

    <el-card class="list-card" shadow="never">
      <div class="card-header">
        <el-button :icon="Refresh" @click="fetchList">刷新</el-button>
      </div>

      <el-table
        v-loading="loading"
        :data="tableData"
        stripe
        border
        style="width: 100%"
      >
        <el-table-column label="产品" min-width="220">
          <template #default="{ row }">
            <div class="product-cell">
              <el-image
                :src="row.productImage"
                fit="cover"
                style="width: 48px; height: 48px; border-radius: 4px"
              >
                <template #error>
                  <div class="image-placeholder">
                    <el-icon :size="16"><Picture /></el-icon>
                  </div>
                </template>
              </el-image>
              <div class="product-info">
                <div class="product-name">{{ row.productName }}</div>
                <div class="product-meta">数量：{{ row.quantity }}</div>
              </div>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="orderNo" label="订单号" width="180" show-overflow-tooltip />
        <el-table-column label="金额" width="130" align="right">
          <template #default="{ row }">¥{{ formatPrice(row.totalPrice) }}</template>
        </el-table-column>
        <el-table-column label="状态" width="110" align="center">
          <template #default="{ row }">
            <el-tag :type="statusTag(row.status).color || 'info'">
              {{ statusTag(row.status).label }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="manufacturerName" label="厂商" width="140" show-overflow-tooltip />
        <el-table-column prop="ownerName" label="业主" width="120" show-overflow-tooltip />
        <el-table-column prop="createdAt" label="创建时间" width="180" />
        <el-table-column label="操作" width="140" align="center" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="goDetail(row)">详情</el-button>
            <el-button
              v-if="canQuote(row)"
              link
              type="primary"
              @click="handleQuote(row)"
            >
              报价
            </el-button>
            <el-button
              v-if="canCancel(row)"
              link
              type="danger"
              @click="handleCancel(row)"
            >
              取消
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <Pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.pageSize"
        :total="pagination.total"
        @change="handlePageChange"
      />
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Refresh, Picture } from '@element-plus/icons-vue'
import SearchBar from '@/components/SearchBar.vue'
import Pagination from '@/components/Pagination.vue'
import { listOrders, cancelOrder, quoteOrder } from '@/api/order'
import { useUserStore } from '@/stores/user'
import { OrderStatus, OrderStatusMap } from '@/types'
import type { Order } from '@/types'

const router = useRouter()
const userStore = useUserStore()

const loading = ref(false)
const tableData = ref<Order[]>([])
const searchForm = reactive<Record<string, unknown>>({})

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

const searchFields = [
  { prop: 'orderNo', label: '订单号', type: 'input' as const, placeholder: '请输入订单号' },
  {
    prop: 'status',
    label: '状态',
    type: 'select' as const,
    options: Object.entries(OrderStatusMap).map(([k, v]) => ({
      label: v.label,
      value: Number(k)
    }))
  },
  { prop: 'dateRange', label: '时间范围', type: 'daterange' as const }
]

async function fetchList() {
  loading.value = true
  try {
    const dateRange = searchForm.dateRange as [string, string] | undefined
    const data = await listOrders({
      page: pagination.page,
      pageSize: pagination.pageSize,
      orderNo: searchForm.orderNo as string | undefined,
      status: searchForm.status as number | undefined,
      startTime: dateRange?.[0],
      endTime: dateRange?.[1]
    })
    tableData.value = data.list
    pagination.total = data.total
  } catch (err) {
    console.error(err)
  } finally {
    loading.value = false
  }
}

function handleSearch(form: Record<string, unknown>) {
  Object.assign(searchForm, form)
  pagination.page = 1
  fetchList()
}

function handleReset() {
  Object.keys(searchForm).forEach((k) => delete searchForm[k])
  pagination.page = 1
  fetchList()
}

function handlePageChange(page: number, pageSize: number) {
  pagination.page = page
  pagination.pageSize = pageSize
  fetchList()
}

function goDetail(row: Order) {
  router.push(`/orders/${row.id}`)
}

function statusTag(status: number) {
  return OrderStatusMap[status] || { label: '未知', color: 'info' }
}

const canQuote = (row: Order) =>
  userStore.userRole === 'manufacturer' && row.status === OrderStatus.PENDING

const canCancel = (row: Order) =>
  [OrderStatus.PENDING, OrderStatus.QUOTED].includes(row.status as any)

async function handleQuote(row: Order) {
  try {
    const { value } = await ElMessageBox.prompt(
      `请输入订单「${row.orderNo}」的报价金额`,
      '订单报价',
      {
        confirmButtonText: '提交',
        cancelButtonText: '取消',
        inputPattern: /^\d+(\.\d{1,2})?$/,
        inputErrorMessage: '请输入正确的金额'
      }
    )
    await quoteOrder(row.id, { unitPrice: Number(value) })
    ElMessage.success('报价成功')
    fetchList()
  } catch {
    // 取消
  }
}

async function handleCancel(row: Order) {
  try {
    await ElMessageBox.confirm(`确定取消订单「${row.orderNo}」吗？`, '提示', {
      type: 'warning'
    })
    await cancelOrder(row.id)
    ElMessage.success('已取消')
    fetchList()
  } catch {
    // 取消
  }
}

function formatPrice(price: number) {
  return (price || 0).toFixed(2)
}

onMounted(fetchList)
</script>

<style lang="scss" scoped>
.order-list {
  .card-header {
    margin-bottom: 16px;
  }

  .product-cell {
    display: flex;
    align-items: center;
    gap: 10px;

    .image-placeholder {
      width: 48px;
      height: 48px;
      background-color: #f5f7fa;
      display: flex;
      align-items: center;
      justify-content: center;
      color: #c0c4cc;
      border-radius: 4px;
    }

    .product-info {
      .product-name {
        color: #303133;
        font-weight: 500;
      }

      .product-meta {
        color: #909399;
        font-size: 12px;
        margin-top: 4px;
      }
    }
  }
}
</style>
