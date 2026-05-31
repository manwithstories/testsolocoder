<template>
  <div class="negotiations-page">
    <div class="page-header">
      <h2 class="page-title">议价管理</h2>
    </div>

    <el-card>
      <el-table :data="negotiations" v-loading="loading" stripe>
        <el-table-column prop="orderNo" label="订单号" width="180" />
        <el-table-column label="商品" min-width="200">
          <template #default="{ row }">
            <div class="product-cell">
              <img :src="getFirstImage(row.product?.images)" class="product-thumb" />
              <div>
                <div class="product-title text-ellipsis">{{ row.product?.title }}</div>
                <div class="product-price">原价：¥{{ row.originalPrice?.toFixed(2) }}</div>
              </div>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="offeredPrice" label="买家出价" width="120">
          <template #default="{ row }">
            <span class="price-text">¥{{ row.offeredPrice?.toFixed(2) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="counterPrice" label="还价" width="120">
          <template #default="{ row }">
            <span v-if="row.counterPrice">¥{{ row.counterPrice.toFixed(2) }}</span>
            <span v-else class="text-lighter">-</span>
          </template>
        </el-table-column>
        <el-table-column prop="buyerMessage" label="买家留言" min-width="150" show-overflow-tooltip />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="createdAt" label="创建时间" width="180">
          <template #default="{ row }">
            {{ formatTime(row.createdAt) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <template v-if="row.status === 1">
              <el-button size="small" type="success" link @click="acceptNegotiation(row)">
                接受
              </el-button>
              <el-button size="small" type="warning" link @click="showCounterDialog(row)">
                还价
              </el-button>
              <el-button size="small" type="danger" link @click="rejectNegotiation(row)">
                拒绝
              </el-button>
            </template>
            <span v-else class="text-lighter">-</span>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-wrapper">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.pageSize"
          :page-sizes="[10, 20, 50]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="fetchNegotiations"
          @current-change="fetchNegotiations"
        />
      </div>
    </el-card>

    <el-dialog v-model="showCounter" title="还价" width="400px">
      <el-form :model="counterForm" label-width="80px">
        <el-form-item label="还价金额">
          <el-input-number
            v-model="counterForm.price"
            :min="0"
            :precision="2"
            placeholder="请输入还价金额"
          />
        </el-form-item>
        <el-form-item label="说明">
          <el-input
            v-model="counterForm.message"
            type="textarea"
            :rows="3"
            placeholder="请输入说明（选填）"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showCounter = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="submitCounter">
          提交
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { orderApi } from '@/api'
import type { Order } from '@/types'

const loading = ref(false)
const negotiations = ref<any[]>([])
const total = ref(0)
const showCounter = ref(false)
const currentOrder = ref<Order | null>(null)
const submitting = ref(false)

const pagination = reactive({
  page: 1,
  pageSize: 10
})

const counterForm = reactive({
  price: 0,
  message: ''
})

function getFirstImage(images: string | undefined): string {
  if (!images) return 'https://picsum.photos/50/50'
  try {
    const arr = JSON.parse(images)
    return Array.isArray(arr) && arr.length > 0 ? arr[0] : 'https://picsum.photos/50/50'
  } catch {
    return 'https://picsum.photos/50/50'
  }
}

function getStatusText(status: number): string {
  const textMap: Record<number, string> = {
    1: '待处理',
    2: '已接受',
    3: '已拒绝',
    4: '已过期'
  }
  return textMap[status] || '未知'
}

function getStatusType(status: number): string {
  const typeMap: Record<number, string> = {
    1: 'warning',
    2: 'success',
    3: 'danger',
    4: 'info'
  }
  return typeMap[status] || 'info'
}

function formatTime(time: string): string {
  return new Date(time).toLocaleString()
}

async function fetchNegotiations() {
  loading.value = true
  try {
    const res = await orderApi.getSellerOrders({
      page: pagination.page,
      pageSize: pagination.pageSize,
      status: 9
    })
    negotiations.value = res.data
      .filter((order: Order) => order.negotiations && order.negotiations.length > 0)
      .map((order: Order) => ({
        ...order.negotiations![0],
        orderNo: order.orderNo,
        originalPrice: order.originalPrice,
        product: order.product
      }))
    total.value = negotiations.value.length
  } catch (error) {
    console.error('Failed to fetch negotiations:', error)
  } finally {
    loading.value = false
  }
}

async function acceptNegotiation(negotiation: any) {
  try {
    await ElMessageBox.confirm('确定接受该议价？', '提示', {
      confirmButtonText: '确认',
      cancelButtonText: '取消'
    })
    await orderApi.handleNegotiation({
      orderNo: negotiation.orderNo,
      accepted: true,
      counterPrice: 0,
      message: ''
    })
    ElMessage.success('已接受议价')
    fetchNegotiations()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Failed to accept:', error)
    }
  }
}

function showCounterDialog(negotiation: any) {
  currentOrder.value = negotiation
  counterForm.price = negotiation.offeredPrice
  counterForm.message = ''
  showCounter.value = true
}

async function submitCounter() {
  if (!currentOrder.value || !counterForm.price) {
    ElMessage.warning('请输入还价金额')
    return
  }

  submitting.value = true
  try {
    await orderApi.handleNegotiation({
      orderNo: currentOrder.value.orderNo,
      accepted: false,
      counterPrice: counterForm.price,
      message: counterForm.message
    })
    ElMessage.success('已提交还价')
    showCounter.value = false
    fetchNegotiations()
  } catch (error: any) {
    ElMessage.error(error.message || '提交失败')
  } finally {
    submitting.value = false
  }
}

async function rejectNegotiation(negotiation: any) {
  try {
    await ElMessageBox.confirm('确定拒绝该议价？', '提示', {
      confirmButtonText: '确认',
      cancelButtonText: '取消'
    })
    await orderApi.handleNegotiation({
      orderNo: negotiation.orderNo,
      accepted: false,
      counterPrice: 0,
      message: ''
    })
    ElMessage.success('已拒绝议价')
    fetchNegotiations()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Failed to reject:', error)
    }
  }
}

onMounted(() => {
  fetchNegotiations()
})
</script>

<style lang="scss" scoped>
.negotiations-page {
  .product-cell {
    display: flex;
    align-items: center;
    gap: 12px;

    .product-thumb {
      width: 50px;
      height: 50px;
      border-radius: 4px;
      object-fit: cover;
    }

    .product-title {
      max-width: 180px;
      margin-bottom: 4px;
    }

    .product-price {
      font-size: 12px;
      color: var(--text-lighter-color);
    }
  }

  .pagination-wrapper {
    display: flex;
    justify-content: center;
    padding: 20px 0 0;
  }

  .text-lighter {
    color: var(--text-lighter-color);
  }
}
</style>
