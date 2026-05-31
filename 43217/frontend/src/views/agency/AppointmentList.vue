<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { ElTable, ElTableColumn, ElButton, ElInput, ElSelect, ElOption, ElPagination, ElTag, ElCard, ElDatePicker, ElMessage, ElDialog, ElForm, ElFormItem } from 'element-plus'
import { Search, Refresh, Calendar, Edit, Check } from '@element-plus/icons-vue'
import { getAgencyAppointments, completeAppointment } from '@/api/appointment'
import type { Appointment } from '@/types'
import dayjs from 'dayjs'

const loading = ref(false)
const appointments = ref<Appointment[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)
const searchKeyword = ref('')
const statusFilter = ref<number | undefined>(undefined)
const dateRange = ref<string[]>([])

const dialogVisible = ref(false)
const currentAppointment = ref<Appointment | null>(null)

const statusMap: Record<number, { text: string; type: string }> = {
  0: { text: '待确认', type: 'warning' },
  1: { text: '已确认', type: 'primary' },
  2: { text: '已完成', type: 'success' },
  3: { text: '已取消', type: 'danger' },
  4: { text: '已过期', type: 'info' }
}

const fetchAppointments = async () => {
  loading.value = true
  try {
    const response = await getAgencyAppointments({ page: page.value, page_size: pageSize.value })
    appointments.value = response.items
    total.value = response.total
  } catch (error) {
    console.error('Failed to fetch appointments:', error)
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  page.value = 1
  fetchAppointments()
}

const handleComplete = async (row: Appointment) => {
  try {
    await completeAppointment(row.id)
    ElMessage.success('预约已完成')
    fetchAppointments()
  } catch (error) {
    console.error('Failed to complete appointment:', error)
  }
}

const handlePageChange = (newPage: number) => {
  page.value = newPage
  fetchAppointments()
}

const handleSizeChange = (newSize: number) => {
  pageSize.value = newSize
  page.value = 1
  fetchAppointments()
}

onMounted(() => {
  fetchAppointments()
})
</script>

<template>
  <div class="appointment-list">
    <ElCard class="filter-card">
      <div class="filter-container">
        <ElInput
          v-model="searchKeyword"
          placeholder="搜索员工姓名/预约号"
          :prefix-icon="Search"
          clearable
          class="search-input"
          @keyup.enter="handleSearch"
        />
        <ElSelect
          v-model="statusFilter"
          placeholder="预约状态"
          clearable
          class="status-select"
        >
          <ElOption label="待确认" :value="0" />
          <ElOption label="已确认" :value="1" />
          <ElOption label="已完成" :value="2" />
          <ElOption label="已取消" :value="3" />
          <ElOption label="已过期" :value="4" />
        </ElSelect>
        <ElDatePicker
          v-model="dateRange"
          type="daterange"
          range-separator="至"
          start-placeholder="开始日期"
          end-placeholder="结束日期"
          class="date-picker"
        />
        <ElButton type="primary" :icon="Search" @click="handleSearch">
          搜索
        </ElButton>
        <ElButton :icon="Refresh" @click="fetchAppointments">
          刷新
        </ElButton>
      </div>
    </ElCard>

    <ElCard>
      <ElTable :data="appointments" v-loading="loading" border stripe>
        <ElTableColumn prop="appointment_no" label="预约号" width="180" />
        <ElTableColumn prop="employee.real_name" label="员工姓名" width="100" />
        <ElTableColumn prop="package.name" label="体检套餐" width="150" />
        <ElTableColumn label="预约日期" width="120">
          <template #default="{ row }">
            {{ dayjs(row.appointment_date).format('YYYY-MM-DD') }}
          </template>
        </ElTableColumn>
        <ElTableColumn label="预约时段" width="120">
          <template #default="{ row }">
            {{ row.start_time }} - {{ row.end_time }}
          </template>
        </ElTableColumn>
        <ElTableColumn label="状态" width="100">
          <template #default="{ row }">
            <ElTag :type="statusMap[row.status]?.type">
              {{ statusMap[row.status]?.text }}
            </ElTag>
          </template>
        </ElTableColumn>
        <ElTableColumn label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <ElButton 
              v-if="row.status === 1" 
              type="success" 
              link 
              :icon="Check"
              @click="handleComplete(row)"
            >
              完成
            </ElButton>
            <ElButton type="primary" link :icon="Edit">
              详情
            </ElButton>
          </template>
        </ElTableColumn>
      </ElTable>

      <div class="pagination-container">
        <ElPagination
          v-model:current-page="page"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handlePageChange"
        />
      </div>
    </ElCard>
  </div>
</template>

<style scoped lang="scss">
.appointment-list {
  .filter-card {
    margin-bottom: 20px;

    .filter-container {
      display: flex;
      gap: 10px;
      align-items: center;
      flex-wrap: wrap;

      .search-input {
        width: 200px;
      }

      .status-select {
        width: 150px;
      }

      .date-picker {
        width: 260px;
      }
    }
  }

  .pagination-container {
    display: flex;
    justify-content: flex-end;
    margin-top: 20px;
  }
}
</style>
