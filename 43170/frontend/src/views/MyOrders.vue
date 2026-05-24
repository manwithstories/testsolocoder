<template>
  <div class="page-container">
    <h1 class="page-title">我的订单</h1>

    <div class="card">
      <el-tabs v-model="activeTab" @tab-change="handleTabChange">
        <el-tab-pane label="全部订单" name="all">
          <div v-loading="loading" class="order-list">
            <el-table :data="orders" stripe style="width: 100%">
              <el-table-column label="订单编号" prop="id" width="100" />
              <el-table-column label="设备" min-width="150">
                <template #default="{ row }">
                  <div class="equipment-cell" @click="goToEquipment(row.equipmentId)">
                    <el-avatar
                      v-if="row.equipment?.images && row.equipment.images.length > 0"
                      :size="48"
                      :src="`/uploads/${row.equipment.images[0].imageUrl}`"
                      shape="square"
                    />
                    <el-avatar v-else :size="48" shape="square">
                      <el-icon><Camera /></el-icon>
                    </el-avatar>
                    <span class="equipment-name">{{ row.equipment?.name || '-' }}</span>
                  </div>
                </template>
              </el-table-column>
              <el-table-column v-if="userStore.isOwner()" label="租借方" width="120">
                <template #default="{ row }">
                  {{ row.renter?.realName || row.renter?.username || '-' }}
                </template>
              </el-table-column>
              <el-table-column v-else label="出租方" width="120">
                <template #default="{ row }">
                  {{ row.owner?.realName || row.owner?.username || '-' }}
                </template>
              </el-table-column>
              <el-table-column label="租期" min-width="200">
                <template #default="{ row }">
                  {{ row.startDate }} 至 {{ row.endDate }}
                </template>
              </el-table-column>
              <el-table-column label="租金" width="100">
                <template #default="{ row }">
                  <span class="price">¥{{ row.totalRent }}</span>
                </template>
              </el-table-column>
              <el-table-column label="状态" width="100">
                <template #default="{ row }">
                  <el-tag :class="getStatusClass(row.status)">
                    {{ getStatusText(row.status) }}
                  </el-tag>
                </template>
              </el-table-column>
              <el-table-column label="操作" width="150" fixed="right">
                <template #default="{ row }">
                  <el-button type="primary" link @click="goToDetail(row.id)">查看</el-button>
                  <el-button
                    v-if="canConfirm(row)"
                    type="success"
                    link
                    @click="handleConfirm(row)"
                  >确认</el-button>
                  <el-button
                    v-if="canComplete(row)"
                    type="success"
                    link
                    @click="handleComplete(row)"
                  >完成</el-button>
                </template>
              </el-table-column>
            </el-table>

            <div v-if="!loading && orders.length === 0" class="empty-state">
              <el-icon><Document /></el-icon>
              <p>暂无订单</p>
            </div>
          </div>
        </el-tab-pane>

        <el-tab-pane v-for="tab in statusTabs" :key="tab.value" :label="tab.label" :name="tab.value">
          <div v-loading="loading" class="order-list">
            <el-table :data="orders.filter(o => o.status === tab.value)" stripe style="width: 100%">
              <el-table-column label="订单编号" prop="id" width="100" />
              <el-table-column label="设备" min-width="150">
                <template #default="{ row }">
                  <div class="equipment-cell" @click="goToEquipment(row.equipmentId)">
                    <el-avatar
                      v-if="row.equipment?.images && row.equipment.images.length > 0"
                      :size="48"
                      :src="`/uploads/${row.equipment.images[0].imageUrl}`"
                      shape="square"
                    />
                    <el-avatar v-else :size="48" shape="square">
                      <el-icon><Camera /></el-icon>
                    </el-avatar>
                    <span class="equipment-name">{{ row.equipment?.name || '-' }}</span>
                  </div>
                </template>
              </el-table-column>
              <el-table-column v-if="userStore.isOwner()" label="租借方" width="120">
                <template #default="{ row }">
                  {{ row.renter?.realName || row.renter?.username || '-' }}
                </template>
              </el-table-column>
              <el-table-column v-else label="出租方" width="120">
                <template #default="{ row }">
                  {{ row.owner?.realName || row.owner?.username || '-' }}
                </template>
              </el-table-column>
              <el-table-column label="租期" min-width="200">
                <template #default="{ row }">
                  {{ row.startDate }} 至 {{ row.endDate }}
                </template>
              </el-table-column>
              <el-table-column label="租金" width="100">
                <template #default="{ row }">
                  <span class="price">¥{{ row.totalRent }}</span>
                </template>
              </el-table-column>
              <el-table-column label="操作" width="150" fixed="right">
                <template #default="{ row }">
                  <el-button type="primary" link @click="goToDetail(row.id)">查看</el-button>
                  <el-button
                    v-if="canConfirm(row)"
                    type="success"
                    link
                    @click="handleConfirm(row)"
                  >确认</el-button>
                  <el-button
                    v-if="canComplete(row)"
                    type="success"
                    link
                    @click="handleComplete(row)"
                  >完成</el-button>
                </template>
              </el-table-column>
            </el-table>

            <div v-if="!loading && orders.filter(o => o.status === tab.value).length === 0" class="empty-state">
              <el-icon><Document /></el-icon>
              <p>暂无{{ tab.label }}订单</p>
            </div>
          </div>
        </el-tab-pane>
      </el-tabs>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { orderApi } from '@/api/order'
import { useUserStore } from '@/stores/user'
import type { Order } from '@/types'

const router = useRouter()
const userStore = useUserStore()
const loading = ref(false)
const orders = ref<Order[]>([])
const activeTab = ref('all')

const statusTabs = [
  { label: '待确认', value: 'pending' },
  { label: '已确认', value: 'confirmed' },
  { label: '租赁中', value: 'rented' },
  { label: '已完成', value: 'completed' },
  { label: '已取消', value: 'cancelled' }
]

onMounted(() => {
  loadOrders()
})

async function loadOrders() {
  loading.value = true
  try {
    const response = await orderApi.getMyOrders()
    orders.value = response.data
  } catch (error) {
    console.error('Failed to load orders:', error)
    ElMessage.error('加载订单列表失败')
  } finally {
    loading.value = false
  }
}

function handleTabChange() {
  // Tab change handler
}

function goToEquipment(id: number) {
  router.push(`/equipments/${id}`)
}

function goToDetail(id: number) {
  router.push(`/orders/${id}`)
}

function canConfirm(order: Order) {
  return order.status === 'pending' && userStore.isOwner() && order.ownerId === userStore.user?.id
}

function canComplete(order: Order) {
  return order.status === 'rented' && userStore.isOwner() && order.ownerId === userStore.user?.id
}

async function handleConfirm(order: Order) {
  try {
    await ElMessageBox.confirm(
      '确认要接受此租赁申请吗？',
      '确认订单',
      {
        confirmButtonText: '确认',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    await orderApi.confirmOrder(order.id)
    ElMessage.success('订单已确认')
    loadOrders()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Failed to confirm order:', error)
    }
  }
}

async function handleComplete(order: Order) {
  try {
    await ElMessageBox.confirm(
      '确认要完成此订单吗？设备已归还且无损坏？',
      '完成订单',
      {
        confirmButtonText: '确认',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    await orderApi.completeOrder(order.id)
    ElMessage.success('订单已完成')
    loadOrders()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Failed to complete order:', error)
    }
  }
}

function getStatusText(status: string) {
  const textMap: Record<string, string> = {
    pending: '待确认',
    confirmed: '已确认',
    rented: '租赁中',
    completed: '已完成',
    cancelled: '已取消'
  }
  return textMap[status] || status
}

function getStatusClass(status: string) {
  return `status-${status}`
}
</script>

<style scoped>
.equipment-cell {
  display: flex;
  align-items: center;
  gap: 12px;
  cursor: pointer;
}

.equipment-name {
  font-weight: 500;
}

.price {
  color: #f56c6c;
  font-weight: 500;
}

.status-pending {
  color: #e6a23c;
}

.status-confirmed {
  color: #409eff;
}

.status-rented {
  color: #67c23a;
}

.status-completed {
  color: #909399;
}

.status-cancelled {
  color: #f56c6c;
}
</style>
