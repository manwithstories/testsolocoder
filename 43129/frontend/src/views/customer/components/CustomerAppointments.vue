<template>
  <el-table :data="appointments" stripe v-loading="loading">
    <el-table-column prop="id" label="ID" width="80" />
    <el-table-column label="服务">
      <template #default="{ row }">
        {{ row.service?.name }}
      </template>
    </el-table-column>
    <el-table-column label="技师">
      <template #default="{ row }">
        {{ row.technician?.name }}
      </template>
    </el-table-column>
    <el-table-column label="时间">
      <template #default="{ row }">
        {{ formatDate(row.appointment_date) }} {{ row.start_time }}-{{ row.end_time }}
      </template>
    </el-table-column>
    <el-table-column prop="status" label="状态">
      <template #default="{ row }">
        <el-tag :type="getStatusType(row.status)">{{ getStatusText(row.status) }}</el-tag>
      </template>
    </el-table-column>
  </el-table>

  <div class="pagination">
    <el-pagination
      v-model:current-page="page"
      v-model:page-size="pageSize"
      :total="total"
      :page-sizes="[10, 20, 50]"
      layout="total, prev, pager, next"
      @current-change="fetchList"
      @size-change="fetchList"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { getAppointments } from '@/api/appointment'
import dayjs from 'dayjs'

const props = defineProps<{
  customerId: number
}>()

const loading = ref(false)
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)
const appointments = ref<any[]>([])

const formatDate = (date: string) => dayjs(date).format('YYYY-MM-DD')

const getStatusType = (status: string) => {
  const types: Record<string, string> = {
    pending: 'info',
    confirmed: 'primary',
    paid: 'success',
    completed: 'success',
    cancelled: 'danger'
  }
  return types[status] || 'info'
}

const getStatusText = (status: string) => {
  const texts: Record<string, string> = {
    pending: '待确认',
    confirmed: '已确认',
    paid: '已支付',
    completed: '已完成',
    cancelled: '已取消'
  }
  return texts[status] || status
}

const fetchList = async () => {
  loading.value = true
  try {
    const res = await getAppointments({
      page: page.value,
      page_size: pageSize.value,
      customer_id: props.customerId
    })
    appointments.value = res.data.list
    total.value = res.data.total
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

watch(() => props.customerId, () => {
  page.value = 1
  fetchList()
})

onMounted(fetchList)
</script>

<style scoped lang="scss">
.pagination {
  margin-top: 16px;
  display: flex;
  justify-content: flex-end;
}
</style>
