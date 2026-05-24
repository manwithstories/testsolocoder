<template>
  <div class="order-detail">
    <el-page-header :icon="ArrowLeft" @back="goBack">
      <template #content>
        <span>订单详情</span>
        <el-tag
          v-if="order.orderNo"
          :type="statusTag(order.status).color || 'info'"
          style="margin-left: 12px"
        >
          {{ statusTag(order.status).label }}
        </el-tag>
      </template>
    </el-page-header>

    <div v-loading="loading" class="detail-wrapper">
      <el-row :gutter="20">
        <el-col :xs="24" :md="16">
          <el-card shadow="never" class="base-card">
            <template #header>基本信息</template>
            <el-descriptions :column="2" border>
              <el-descriptions-item label="订单号">{{ order.orderNo }}</el-descriptions-item>
              <el-descriptions-item label="创建时间">{{ order.createdAt }}</el-descriptions-item>
              <el-descriptions-item label="厂商">{{ order.manufacturerName || '-' }}</el-descriptions-item>
              <el-descriptions-item label="业主">{{ order.ownerName || '-' }}</el-descriptions-item>
              <el-descriptions-item label="设计师">
                {{ order.designerName || '-' }}
              </el-descriptions-item>
              <el-descriptions-item label="更新时间">{{ order.updatedAt }}</el-descriptions-item>
              <el-descriptions-item label="备注" :span="2">{{ order.remark || '-' }}</el-descriptions-item>
            </el-descriptions>
          </el-card>

          <el-card shadow="never" class="items-card">
            <template #header>订单项</template>
            <el-table :data="order.items || []" border stripe>
              <el-table-column label="产品" min-width="200">
                <template #default="{ row }">
                  <div class="item-cell">
                    <el-image
                      :src="row.productImage"
                      fit="cover"
                      style="width: 48px; height: 48px; border-radius: 4px"
                    />
                    <span class="item-name">{{ row.productName }}</span>
                  </div>
                </template>
              </el-table-column>
              <el-table-column label="单价" width="120" align="right">
                <template #default="{ row }">¥{{ formatPrice(row.unitPrice) }}</template>
              </el-table-column>
              <el-table-column prop="quantity" label="数量" width="80" align="center" />
              <el-table-column label="小计" width="120" align="right">
                <template #default="{ row }">¥{{ formatPrice(row.totalPrice) }}</template>
              </el-table-column>
              <el-table-column label="自定义选项" min-width="200">
                <template #default="{ row }">
                  <template v-if="row.options && row.options.length > 0">
                    <el-tag
                      v-for="(opt, idx) in row.options"
                      :key="idx"
                      class="opt-tag"
                      size="small"
                    >
                      {{ optTypeLabel(opt.type) }}：{{ opt.value }}
                    </el-tag>
                  </template>
                  <span v-else class="no-opt">-</span>
                </template>
              </el-table-column>
            </el-table>

            <div class="summary">
              <div>合计：</div>
              <div class="total">¥{{ formatPrice(order.totalPrice) }}</div>
            </div>
          </el-card>
        </el-col>

        <el-col :xs="24" :md="8">
          <el-card shadow="never" class="timeline-card">
            <template #header>状态流转</template>
            <el-timeline>
              <el-timeline-item
                v-for="(log, idx) in (order.statusLogs || []).slice().reverse()"
                :key="log.id || idx"
                :timestamp="log.createdAt"
                :color="timelineColor(log.status)"
              >
                <div class="timeline-item">
                  <div class="timeline-status">
                    {{ statusTag(log.status).label }}
                  </div>
                  <div class="timeline-operator">
                    {{ log.operatorName || '-' }}
                  </div>
                  <div v-if="log.remark" class="timeline-remark">
                    {{ log.remark }}
                  </div>
                </div>
              </el-timeline-item>
              <el-timeline-item v-if="!order.statusLogs || order.statusLogs.length === 0" color="#909399">
                暂无流转记录
              </el-timeline-item>
            </el-timeline>
          </el-card>

          <el-card shadow="never" class="action-card">
            <template #header>操作</template>
            <div class="action-list">
              <el-button
                v-if="canQuote"
                type="primary"
                :loading="actioning"
                @click="handleQuote"
              >
                报价
              </el-button>
              <el-button
                v-if="canConfirm"
                type="success"
                :loading="actioning"
                @click="handleConfirm"
              >
                确认订单
              </el-button>
              <el-button
                v-if="canProduce"
                :loading="actioning"
                @click="handleProduce"
              >
                开始生产
              </el-button>
              <el-button
                v-if="canShip"
                :loading="actioning"
                @click="handleShip"
              >
                发货
              </el-button>
              <el-button
                v-if="canComplete"
                type="success"
                :loading="actioning"
                @click="handleComplete"
              >
                完成订单
              </el-button>
              <el-button
                v-if="canCancel"
                type="danger"
                :loading="actioning"
                @click="handleCancel"
              >
                取消订单
              </el-button>
              <el-empty
                v-if="!canQuote && !canConfirm && !canProduce && !canShip && !canComplete && !canCancel"
                description="当前状态无可执行操作"
                :image-size="60"
              />
            </div>
          </el-card>
        </el-col>
      </el-row>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ArrowLeft } from '@element-plus/icons-vue'
import {
  getOrder,
  quoteOrder,
  confirmOrder,
  cancelOrder,
  startProduceOrder,
  shipOrder,
  completeOrder
} from '@/api/order'
import { OrderStatus, OrderStatusMap } from '@/types'
import type { Order, OptionType } from '@/types'
import { useUserStore } from '@/stores/user'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const loading = ref(false)
const actioning = ref(false)
const order = ref<Order>({
  id: 0,
  orderNo: '',
  productId: 0,
  quantity: 0,
  unitPrice: 0,
  totalPrice: 0,
  status: 0,
  manufacturerId: 0,
  createdAt: '',
  updatedAt: ''
})

async function fetchDetail() {
  loading.value = true
  try {
    const data = await getOrder(route.params.id as string)
    order.value = data
  } catch (err) {
    console.error(err)
  } finally {
    loading.value = false
  }
}

const canQuote = computed(
  () => userStore.userRole === 'manufacturer' && order.value.status === OrderStatus.PENDING
)
const canConfirm = computed(
  () => userStore.userRole === 'owner' && order.value.status === OrderStatus.QUOTED
)
const canProduce = computed(
  () => userStore.userRole === 'manufacturer' && order.value.status === OrderStatus.CONFIRMED
)
const canShip = computed(
  () => userStore.userRole === 'manufacturer' && order.value.status === OrderStatus.PRODUCING
)
const canComplete = computed(
  () => userStore.userRole === 'owner' && order.value.status === OrderStatus.SHIPPED
)
const canCancel = computed(
  () =>
    [OrderStatus.PENDING, OrderStatus.QUOTED].includes(order.value.status as any) &&
    (userStore.userRole === 'manufacturer' || userStore.userRole === 'owner')
)

function statusTag(status: number) {
  return OrderStatusMap[status] || { label: '未知', color: 'info' }
}

function timelineColor(status: number) {
  const map: Record<number, string> = {
    0: '#e6a23c',
    1: '#909399',
    2: '#409eff',
    3: '#409eff',
    4: '#409eff',
    5: '#67c23a',
    9: '#f56c6c'
  }
  return map[status] || '#909399'
}

function optTypeLabel(type: OptionType) {
  return ({ size: '尺寸', material: '材质', color: '颜色' } as Record<OptionType, string>)[type] || type
}

async function handleQuote() {
  try {
    const { value } = await ElMessageBox.prompt('请输入报价金额', '订单报价', {
      confirmButtonText: '提交',
      inputPattern: /^\d+(\.\d{1,2})?$/,
      inputErrorMessage: '请输入正确的金额'
    })
    actioning.value = true
    await quoteOrder(order.value.id, { unitPrice: Number(value) })
    ElMessage.success('报价成功')
    fetchDetail()
  } catch {
    // 取消
  } finally {
    actioning.value = false
  }
}

async function handleConfirm() {
  try {
    await ElMessageBox.confirm('确定确认该订单吗？确认后将进入生产阶段。', '提示', {
      type: 'warning'
    })
    actioning.value = true
    await confirmOrder(order.value.id)
    ElMessage.success('订单已确认')
    fetchDetail()
  } catch {
    // 取消
  } finally {
    actioning.value = false
  }
}

async function handleProduce() {
  try {
    actioning.value = true
    await startProduceOrder(order.value.id)
    ElMessage.success('已开始生产')
    fetchDetail()
  } catch (err) {
    console.error(err)
  } finally {
    actioning.value = false
  }
}

async function handleShip() {
  try {
    actioning.value = true
    await shipOrder(order.value.id)
    ElMessage.success('已发货')
    fetchDetail()
  } catch (err) {
    console.error(err)
  } finally {
    actioning.value = false
  }
}

async function handleComplete() {
  try {
    await ElMessageBox.confirm('确定完成该订单吗？', '提示', { type: 'warning' })
    actioning.value = true
    await completeOrder(order.value.id)
    ElMessage.success('订单已完成')
    fetchDetail()
  } catch {
    // 取消
  } finally {
    actioning.value = false
  }
}

async function handleCancel() {
  try {
    const { value } = await ElMessageBox.prompt('请输入取消原因', '取消订单', {
      confirmButtonText: '提交'
    })
    actioning.value = true
    await cancelOrder(order.value.id, { remark: value })
    ElMessage.success('订单已取消')
    fetchDetail()
  } catch {
    // 取消
  } finally {
    actioning.value = false
  }
}

function goBack() {
  router.back()
}

function formatPrice(price: number) {
  return (price || 0).toFixed(2)
}

onMounted(fetchDetail)
</script>

<style lang="scss" scoped>
.order-detail {
  .detail-wrapper {
    margin-top: 16px;
  }

  .base-card,
  .items-card,
  .timeline-card,
  .action-card {
    margin-bottom: 16px;
  }

  .items-card {
    .item-cell {
      display: flex;
      align-items: center;
      gap: 10px;

      .item-name {
        color: #303133;
      }
    }

    .opt-tag {
      margin-right: 6px;
      margin-bottom: 4px;
    }

    .no-opt {
      color: #c0c4cc;
    }
  }

  .summary {
    display: flex;
    justify-content: flex-end;
    align-items: center;
    gap: 12px;
    margin-top: 12px;
    font-size: 16px;
    color: #606266;

    .total {
      font-size: 22px;
      font-weight: 600;
      color: #f56c6c;
    }
  }

  .timeline-card {
    .timeline-item {
      .timeline-status {
        font-size: 14px;
        font-weight: 500;
        color: #303133;
      }

      .timeline-operator {
        font-size: 12px;
        color: #909399;
        margin-top: 2px;
      }

      .timeline-remark {
        font-size: 12px;
        color: #606266;
        margin-top: 4px;
        background-color: #f5f7fa;
        padding: 4px 8px;
        border-radius: 4px;
      }
    }
  }

  .action-card {
    .action-list {
      display: flex;
      flex-direction: column;
      gap: 10px;

      .el-button {
        width: 100%;
      }
    }
  }
}
</style>
