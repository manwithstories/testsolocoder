<template>
  <div class="appointment-list">
    <div class="page-header">
      <h2 class="page-title">看房预约</h2>
    </div>

    <div class="search-bar card">
      <el-select v-model="statusFilter" placeholder="状态" clearable style="width: 150px">
        <el-option label="待确认" :value="1" />
        <el-option label="已确认" :value="2" />
        <el-option label="已完成" :value="3" />
        <el-option label="已取消" :value="0" />
      </el-select>
      <el-button type="primary" @click="loadData">搜索</el-button>
    </div>

    <div class="card">
      <el-table :data="appointments" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column label="租户" width="150">
          <template #default="{ row }">
            {{ row.tenant?.name }} ({{ row.tenant?.phone }})
          </template>
        </el-table-column>
        <el-table-column label="房源" min-width="150">
          <template #default="{ row }">
            {{ row.property?.title }}
          </template>
        </el-table-column>
        <el-table-column prop="visitTime" label="看房时间" width="180">
          <template #default="{ row }">
            {{ formatDate(row.visitTime) }}
          </template>
        </el-table-column>
        <el-table-column prop="remark" label="备注" min-width="150" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">{{ getStatusText(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="180" fixed="right">
          <template #default="{ row }">
            <el-button link type="success" @click="updateStatus(row, 2)" v-if="row.status === 1">确认</el-button>
            <el-button link type="primary" @click="updateStatus(row, 3)" v-if="row.status === 2">完成</el-button>
            <el-button link type="danger" @click="updateStatus(row, 0)">取消</el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        class="pagination"
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.pageSize"
        :total="pagination.total"
        :page-sizes="[10, 20, 50]"
        layout="total, sizes, prev, pager, next, jumper"
        @current-change="loadData"
        @size-change="loadData"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import type { Appointment } from '@/types'
import { getAppointments, updateAppointmentStatus } from '@/api/tenant'
import dayjs from 'dayjs'

const loading = ref(false)
const appointments = ref<Appointment[]>([])
const statusFilter = ref<number | ''>('')

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

onMounted(() => {
  loadData()
})

async function loadData() {
  loading.value = true
  try {
    const res = await getAppointments({
      page: pagination.page,
      pageSize: pagination.pageSize,
      status: statusFilter.value || undefined
    })
    appointments.value = res.data.list
    pagination.total = res.data.total
  } catch (error) {
    console.error('Failed to load appointments:', error)
  } finally {
    loading.value = false
  }
}

function formatDate(date: string) {
  return dayjs(date).format('YYYY-MM-DD HH:mm')
}

function getStatusType(status: number) {
  switch (status) {
    case 1: return 'warning'
    case 2: return 'primary'
    case 3: return 'success'
    default: return 'info'
  }
}

function getStatusText(status: number) {
  switch (status) {
    case 1: return '待确认'
    case 2: return '已确认'
    case 3: return '已完成'
    default: return '已取消'
  }
}

async function updateStatus(row: Appointment, status: number) {
  try {
    await updateAppointmentStatus(row.id, status)
    ElMessage.success('操作成功')
    loadData()
  } catch (error) {
    console.error(error)
  }
}
</script>

<style scoped>
.appointment-list {
  padding: 0;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}
</style>
