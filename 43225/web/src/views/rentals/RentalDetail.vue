<template>
  <div class="rental-detail" v-if="rental">
    <el-page-header @back="$router.back()" content="订单详情" class="page-header" />

    <el-row :gutter="16">
      <el-col :xs="24" :md="16">
        <div class="card-container">
          <h3>订单信息</h3>
          <el-descriptions :column="2" border>
            <el-descriptions-item label="订单号">{{ rental.id }}</el-descriptions-item>
            <el-descriptions-item label="状态">
              <el-tag :type="getStatusType(rental.status)">
                {{ getStatusText(rental.status) }}
              </el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="船只">{{ rental.ship?.name }}</el-descriptions-item>
            <el-descriptions-item label="租赁类型">{{ getRentalTypeText(rental.rental_type) }}</el-descriptions-item>
            <el-descriptions-item label="开始时间">{{ formatDateTime(rental.start_date) }}</el-descriptions-item>
            <el-descriptions-item label="结束时间">{{ formatDateTime(rental.end_date) }}</el-descriptions-item>
            <el-descriptions-item label="起始地">{{ rental.start_location || '-' }}</el-descriptions-item>
            <el-descriptions-item label="目的地">{{ rental.end_location || '-' }}</el-descriptions-item>
          </el-descriptions>
        </div>

        <div class="card-container" style="margin-top: 16px">
          <h3>租客信息</h3>
          <el-descriptions :column="2" border>
            <el-descriptions-item label="姓名">{{ rental.tenant?.full_name }}</el-descriptions-item>
            <el-descriptions-item label="邮箱">{{ rental.tenant?.email }}</el-descriptions-item>
            <el-descriptions-item label="电话">{{ rental.tenant?.phone || '-' }}</el-descriptions-item>
            <el-descriptions-item label="紧急联系人">{{ rental.emergency_contact }}</el-descriptions-item>
            <el-descriptions-item label="紧急电话">{{ rental.emergency_phone }}</el-descriptions-item>
            <el-descriptions-item label="乘客人数">{{ rental.passenger_count }}人</el-descriptions-item>
          </el-descriptions>
        </div>

        <div class="card-container" style="margin-top: 16px" v-if="rental.notes">
          <h3>备注</h3>
          <p>{{ rental.notes }}</p>
        </div>
      </el-col>

      <el-col :xs="24" :md="8">
        <div class="card-container sticky">
          <h3>费用明细</h3>
          <div class="price-list">
            <div class="price-item">
              <span>基础费用</span>
              <span>{{ rental.currency }} {{ rental.base_amount }}</span>
            </div>
            <div class="price-item">
              <span>保险费用</span>
              <span>{{ rental.currency }} {{ rental.insurance_amount }}</span>
            </div>
            <div class="price-item">
              <span>平台服务费</span>
              <span>{{ rental.currency }} {{ rental.platform_fee }}</span>
            </div>
            <el-divider />
            <div class="price-item total">
              <span>总计</span>
              <span>{{ rental.currency }} {{ rental.total_amount }}</span>
            </div>
            <div class="price-item">
              <span>押金</span>
              <span>{{ rental.currency }} {{ rental.deposit_amount }}</span>
            </div>
          </div>

          <el-divider />

          <div class="status-actions">
            <el-button
              v-if="rental.status === 'pending'"
              type="success"
              style="width: 100%"
              @click="handleUpdateStatus('confirmed')"
            >
              确认订单
            </el-button>
            <el-button
              v-if="rental.status === 'confirmed'"
              type="primary"
              style="width: 100%"
              @click="handleUpdateStatus('active')"
            >
              开始租赁
            </el-button>
            <el-button
              v-if="rental.status === 'active'"
              type="success"
              style="width: 100%"
              @click="handleUpdateStatus('completed')"
            >
              完成订单
            </el-button>
            <el-button
              v-if="['pending', 'confirmed'].includes(rental.status)"
              type="danger"
              style="width: 100%"
              @click="handleCancel"
            >
              取消订单
            </el-button>
          </div>
        </div>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getRentalApi, updateRentalStatusApi, cancelRentalApi } from '@/api/rental'
import type { Rental } from '@/types/rental'
import dayjs from 'dayjs'

const route = useRoute()
const router = useRouter()
const rental = ref<Rental | null>(null)

const formatDateTime = (date: string) => dayjs(date).format('YYYY-MM-DD HH:mm')

const getStatusType = (status: string) => {
  const map: Record<string, string> = {
    pending: 'warning',
    confirmed: 'primary',
    active: 'success',
    completed: 'info',
    cancelled: 'danger',
    refunded: 'info'
  }
  return map[status] || 'info'
}

const getStatusText = (status: string) => {
  const map: Record<string, string> = {
    pending: '待确认',
    confirmed: '已确认',
    active: '进行中',
    completed: '已完成',
    cancelled: '已取消',
    refunded: '已退款'
  }
  return map[status] || status
}

const getRentalTypeText = (type: string) => {
  const map: Record<string, string> = {
    daily: '按天租赁',
    hourly: '按小时租赁',
    voyage: '航程租赁'
  }
  return map[type] || type
}

const handleUpdateStatus = async (status: string) => {
  try {
    await ElMessageBox.confirm('确定要更新订单状态吗？', '提示', {
      type: 'warning'
    })
    await updateRentalStatusApi(rental.value!.id, { status: status as any })
    ElMessage.success('状态更新成功')
    fetchRental()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('状态更新失败')
    }
  }
}

const handleCancel = async () => {
  try {
    await ElMessageBox.confirm('确定要取消该订单吗？', '提示', {
      type: 'warning'
    })
    await cancelRentalApi(rental.value!.id)
    ElMessage.success('订单已取消')
    router.back()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('取消失败')
    }
  }
}

const fetchRental = async () => {
  const id = route.params.id as string
  try {
    const res: any = await getRentalApi(id)
    rental.value = res.data
  } catch (error) {
    console.error('Failed to fetch rental:', error)
  }
}

onMounted(fetchRental)
</script>

<style lang="scss" scoped>
.rental-detail {
  .page-header {
    margin-bottom: 16px;
    padding: 12px 20px;
    background: #fff;
    border-radius: 8px;
  }

  .card-container {
    background: #fff;
    border-radius: 8px;
    padding: 20px;

    h3 {
      margin-bottom: 16px;
    }
  }

  .sticky {
    position: sticky;
    top: 24px;
  }

  .price-list {
    .price-item {
      display: flex;
      justify-content: space-between;
      margin-bottom: 12px;
      font-size: 14px;

      &.total {
        font-size: 18px;
        font-weight: 600;
        color: #1890ff;
      }
    }
  }

  .status-actions {
    display: flex;
    flex-direction: column;
    gap: 8px;
  }
}
</style>
