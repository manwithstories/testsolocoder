<template>
  <div class="dashboard-reservations">
    <div class="card-shadow p-20">
      <h2 class="page-title mb-20">预约管理</h2>

      <el-form :inline="true" :model="query" class="mb-20">
        <el-form-item label="展览">
          <el-select v-model="query.exhibition_id" clearable style="width: 200px" @change="fetchList">
            <el-option v-for="e in exhibitions" :key="e.id" :label="e.title" :value="e.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="query.status" clearable style="width: 150px" @change="fetchList">
            <el-option label="待确认" value="pending" />
            <el-option label="已确认" value="confirmed" />
            <el-option label="已取消" value="cancelled" />
            <el-option label="已完成" value="completed" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="fetchList">查询</el-button>
        </el-form-item>
      </el-form>

      <el-table :data="list" v-loading="loading" border>
        <el-table-column prop="id" label="ID" width="60" />
        <el-table-column prop="user.nickname" label="用户" width="100" />
        <el-table-column prop="exhibition.title" label="展览" min-width="150" show-overflow-tooltip />
        <el-table-column label="时段" width="180">
          <template #default="{ row }">
            {{ formatDate(row.time_slot?.date || row.created_at) }} {{ row.time_slot?.start_time }} - {{ row.time_slot?.end_time }}
          </template>
        </el-table-column>
        <el-table-column prop="visitor_count" label="人数" width="80" />
        <el-table-column label="导览类型" width="100">
          <template #default="{ row }">
            {{ guideTypeText(row.guide_type) }}
          </template>
        </el-table-column>
        <el-table-column prop="total_price" label="金额" width="100">
          <template #default="{ row }">¥{{ row.total_price }}</template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="statusType(row.status)" size="small">
              {{ statusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="qr_code" label="预约码" width="140" />
        <el-table-column prop="created_at" label="创建时间" width="160">
          <template #default="{ row }">{{ formatDateTime(row.created_at) }}</template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button v-if="row.status === 'pending'" type="success" size="small" link @click="handleConfirm(row)">确认</el-button>
            <el-button v-if="row.status === 'confirmed'" type="primary" size="small" link @click="handleCheckIn(row)">签到</el-button>
            <el-button v-if="row.status === 'pending' || row.status === 'confirmed'" type="warning" size="small" link @click="handleCancel(row)">取消</el-button>
            <el-button type="primary" size="small" link @click="handleDetail(row)">详情</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination mt-20">
        <el-pagination
          v-model:current-page="query.page"
          v-model:page-size="query.page_size"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          :page-sizes="[10, 20, 50, 100]"
          @size-change="fetchList"
          @current-change="fetchList"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import * as reservationApi from '@/api/reservation'
import * as exhibitionApi from '@/api/exhibition'
import type { Reservation, Exhibition } from '@/types'
import dayjs from 'dayjs'

const loading = ref(false)
const list = ref<Reservation[]>([])
const total = ref(0)
const exhibitions = ref<Exhibition[]>([])

const query = reactive({
  page: 1,
  page_size: 10,
  exhibition_id: undefined as number | undefined,
  status: ''
})

const statusType = (status: string) => {
  if (status === 'confirmed') return 'success'
  if (status === 'cancelled') return 'info'
  if (status === 'completed') return 'success'
  return 'warning'
}

const statusText = (status: string) => {
  const map: Record<string, string> = {
    pending: '待确认',
    confirmed: '已确认',
    cancelled: '已取消',
    completed: '已完成'
  }
  return map[status] || status
}

const guideTypeText = (type: string) => {
  const map: Record<string, string> = {
    standard: '标准参观',
    audio: '语音导览',
    human: '人工讲解'
  }
  return map[type] || type
}

const formatDate = (date: string) => {
  return dayjs(date).format('YYYY-MM-DD')
}

const formatDateTime = (date: string) => {
  return dayjs(date).format('YYYY-MM-DD HH:mm')
}

const fetchExhibitions = async () => {
  try {
    const res = await exhibitionApi.listExhibitions({ page: 1, page_size: 1000 })
    exhibitions.value = res.data.list
  } catch (e) {
    console.error(e)
  }
}

const fetchList = async () => {
  if (!query.exhibition_id && exhibitions.value.length > 0) {
    query.exhibition_id = exhibitions.value[0]?.id
  }
  if (!query.exhibition_id) return
  try {
    loading.value = true
    const res = await reservationApi.listExhibitionReservations(query.exhibition_id, {
      page: query.page,
      page_size: query.page_size,
      status: query.status
    })
    list.value = res.data.list
    total.value = res.data.total
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

const handleConfirm = async (row: Reservation) => {
  ElMessageBox.confirm('确定确认该预约吗？', '提示', { type: 'warning' }).then(async () => {
    await reservationApi.confirmReservation(row.id)
    ElMessage.success('确认成功')
    fetchList()
  }).catch(() => {})
}

const handleCancel = async (row: Reservation) => {
  ElMessageBox.confirm('确定取消该预约吗？', '提示', { type: 'warning' }).then(async () => {
    await reservationApi.cancelReservation(row.id, '管理员取消')
    ElMessage.success('取消成功')
    fetchList()
  }).catch(() => {})
}

const handleCheckIn = async (row: Reservation) => {
  ElMessageBox.confirm('确定为该预约签到吗？', '提示', { type: 'warning' }).then(async () => {
    await reservationApi.checkIn(row.qr_code)
    ElMessage.success('签到成功')
    fetchList()
  }).catch(() => {})
}

const handleDetail = (row: Reservation) => {
  console.log('Detail:', row)
}

onMounted(() => {
  fetchExhibitions()
})
</script>
