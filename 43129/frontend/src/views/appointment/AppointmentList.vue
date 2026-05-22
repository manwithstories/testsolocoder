<template>
  <div class="page-container">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>预约管理</span>
          <div class="header-actions">
            <el-select v-model="filterStatus" placeholder="状态" clearable style="width: 120px">
              <el-option label="待确认" value="pending" />
              <el-option label="已确认" value="confirmed" />
              <el-option label="已支付" value="paid" />
              <el-option label="已完成" value="completed" />
              <el-option label="已取消" value="cancelled" />
            </el-select>
            <el-date-picker
              v-model="filterDate"
              type="date"
              placeholder="选择日期"
              clearable
              style="width: 160px"
              value-format="YYYY-MM-DD"
            />
            <el-button :icon="Search" @click="fetchList">查询</el-button>
            <el-button type="primary" :icon="Plus" @click="$router.push('/appointments/create')">
              新建预约
            </el-button>
          </div>
        </div>
      </template>

      <el-table :data="appointments" stripe v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column label="顾客">
          <template #default="{ row }">
            {{ row.customer?.name || row.customer?.user?.phone }}
          </template>
        </el-table-column>
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
        <el-table-column label="日期" width="120">
          <template #default="{ row }">
            {{ formatDate(row.appointment_date) }}
          </template>
        </el-table-column>
        <el-table-column label="时间" width="140">
          <template #default="{ row }">
            {{ row.start_time }} - {{ row.end_time }}
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">{{ getStatusText(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button
              v-if="row.status === 'confirmed' || row.status === 'paid'"
              type="success"
              link
              @click="handleComplete(row)"
            >
              完成
            </el-button>
            <el-button
              v-if="row.status === 'confirmed' || row.status === 'pending'"
              type="warning"
              link
              @click="handleCancel(row)"
            >
              取消
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination">
        <el-pagination
          v-model:current-page="page"
          v-model:page-size="pageSize"
          :total="total"
          :page-sizes="[10, 20, 50]"
          layout="total, sizes, prev, pager, next, jumper"
          @current-change="fetchList"
          @size-change="fetchList"
        />
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getAppointments, completeAppointment, cancelAppointment } from '@/api/appointment'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search, Plus } from '@element-plus/icons-vue'
import dayjs from 'dayjs'
import type { Appointment } from '@/types'

const loading = ref(false)
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)
const appointments = ref<Appointment[]>([])
const filterStatus = ref('')
const filterDate = ref('')

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
    const params: any = {
      page: page.value,
      page_size: pageSize.value
    }
    if (filterStatus.value) {
      params.status = filterStatus.value
    }
    if (filterDate.value) {
      params.start_date = filterDate.value
      params.end_date = filterDate.value
    }
    const res = await getAppointments(params)
    appointments.value = res.data.list
    total.value = res.data.total
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

const handleComplete = async (row: Appointment) => {
  try {
    await ElMessageBox.confirm('确定标记为已完成吗？', '提示', { type: 'warning' })
    await completeAppointment(row.id)
    ElMessage.success('操作成功')
    fetchList()
  } catch (e) {
    if (e !== 'cancel') console.error(e)
  }
}

const handleCancel = async (row: Appointment) => {
  try {
    const { value } = await ElMessageBox.prompt('请输入取消原因', '取消预约', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      inputType: 'textarea'
    })
    await cancelAppointment({ id: row.id, cancel_reason: value })
    ElMessage.success('取消成功')
    fetchList()
  } catch (e) {
    if (e !== 'cancel') console.error(e)
  }
}

onMounted(fetchList)
</script>

<style scoped lang="scss">
.page-container {
  .card-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
  }

  .header-actions {
    display: flex;
    gap: 10px;
  }

  .pagination {
    margin-top: 20px;
    display: flex;
    justify-content: flex-end;
  }
}
</style>
