<template>
  <div class="container" style="padding-top: 20px;">
    <el-card v-loading="loading">
      <template v-if="booking">
        <h2>预订详情</h2>
        <el-descriptions :column="2" border>
          <el-descriptions-item label="预订号">{{ booking.booking_no }}</el-descriptions-item>
          <el-descriptions-item label="状态">
            <el-tag :type="getStatusType(booking.status)">{{ getStatusText(booking.status) }}</el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="车辆">{{ booking.car?.brand }} {{ booking.car?.model }}</el-descriptions-item>
          <el-descriptions-item label="车牌号">{{ booking.car?.license_plate }}</el-descriptions-item>
          <el-descriptions-item label="取车门店">{{ booking.pickup_store?.name }}</el-descriptions-item>
          <el-descriptions-item label="还车门店">{{ booking.return_store?.name }}</el-descriptions-item>
          <el-descriptions-item label="取车时间">{{ booking.pickup_time }}</el-descriptions-item>
          <el-descriptions-item label="还车时间">{{ booking.return_time }}</el-descriptions-item>
          <el-descriptions-item label="租赁天数">{{ booking.total_days }}天</el-descriptions-item>
          <el-descriptions-item label="基础价格">¥{{ booking.base_price }}</el-descriptions-item>
          <el-descriptions-item label="优惠金额" v-if="booking.discount > 0">-¥{{ booking.discount }}</el-descriptions-item>
          <el-descriptions-item label="实付金额">¥{{ booking.final_price }}</el-descriptions-item>
          <el-descriptions-item label="押金" v-if="booking.deposit > 0">¥{{ booking.deposit }}</el-descriptions-item>
        </el-descriptions>

        <div style="margin-top: 20px;">
          <el-button
            v-if="booking.status === 'pending' || booking.status === 'confirmed'"
            type="danger"
            @click="handleCancel"
          >
            取消预订
          </el-button>
        </div>
      </template>
      <el-empty v-else description="预订不存在" />
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { bookingApi } from '@/api'
import type { Booking } from '@/types'

const route = useRoute()
const loading = ref(false)
const booking = ref<Booking | null>(null)

const getStatusType = (status: string) => {
  const types: Record<string, string> = {
    pending: 'warning',
    confirmed: 'primary',
    completed: 'success',
    cancelled: 'info',
    no_show: 'danger'
  }
  return types[status] || ''
}

const getStatusText = (status: string) => {
  const texts: Record<string, string> = {
    pending: '待确认',
    confirmed: '已确认',
    completed: '已完成',
    cancelled: '已取消',
    no_show: '未到店'
  }
  return texts[status] || status
}

const loadBooking = async () => {
  loading.value = true
  try {
    const id = route.params.id as string
    const res = await bookingApi.getBookingById(parseInt(id))
    booking.value = res.data
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}

const handleCancel = async () => {
  try {
    await ElMessageBox.confirm('确定要取消此预订吗？', '取消预订', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await bookingApi.cancelBooking(booking.value!.id, '用户取消')
    ElMessage.success('预订已取消')
    loadBooking()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error.message || '取消失败')
    }
  }
}

onMounted(() => {
  loadBooking()
})
</script>
