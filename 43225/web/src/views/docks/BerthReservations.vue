<template>
  <div class="berth-reservations">
    <div class="page-header">
      <h2 class="page-title">泊位预约</h2>
    </div>

    <div class="card-container">
      <el-form :model="form" label-width="100px" inline>
        <el-form-item label="选择码头">
          <el-select v-model="form.dock_id" placeholder="请选择码头" style="width: 200px" @change="fetchBerths">
            <el-option v-for="dock in docks" :key="dock.id" :label="dock.name" :value="dock.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="泊位类型">
          <el-select v-model="form.berth_type" placeholder="全部" clearable style="width: 120px">
            <el-option label="小型" value="small" />
            <el-option label="中型" value="medium" />
            <el-option label="大型" value="large" />
          </el-select>
        </el-form-item>
      </el-form>

      <el-table :data="berths" v-loading="loading" style="width: 100%; margin-top: 16px">
        <el-table-column prop="number" label="泊位号" width="100" />
        <el-table-column prop="berth_type" label="类型" width="100">
          <template #default="{ row }">
            {{ getBerthTypeText(row.berth_type) }}
          </template>
        </el-table-column>
        <el-table-column label="尺寸" width="120">
          <template #default="{ row }">
            {{ row.max_length }}m × {{ row.max_width }}m
          </template>
        </el-table-column>
        <el-table-column label="设施" width="150">
          <template #default="{ row }">
            <el-tag v-if="row.has_water" size="small" type="success">水</el-tag>
            <el-tag v-if="row.has_electric" size="small" type="primary">电</el-tag>
            <el-tag v-if="row.has_internet" size="small" type="warning">网络</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="价格" width="180">
          <template #default="{ row }">
            ¥{{ row.daily_rate }}/天 或 ¥{{ row.hourly_rate }}/小时
          </template>
        </el-table-column>
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 'available' ? 'success' : 'warning'" size="small">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="100" fixed="right">
          <template #default="{ row }">
            <el-button
              type="primary"
              size="small"
              :disabled="row.status !== 'available'"
              @click="showBookingDialog(row)"
            >
              预约
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <el-dialog v-model="showDialog" title="预约泊位" width="500px">
      <el-form :model="bookingForm" label-width="100px">
        <el-form-item label="泊位号">
          <el-input :model-value="selectedBerth?.number" disabled />
        </el-form-item>
        <el-form-item label="开始时间">
          <el-date-picker
            v-model="bookingForm.start_time"
            type="datetime"
            placeholder="选择开始时间"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="结束时间">
          <el-date-picker
            v-model="bookingForm.end_time"
            type="datetime"
            placeholder="选择结束时间"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="bookingForm.notes" type="textarea" :rows="2" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showDialog = false">取消</el-button>
        <el-button type="primary" :loading="bookingLoading" @click="handleBooking">确认预约</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { getDocksApi, getBerthsApi, createReservationApi } from '@/api/berth'
import type { Dock, Berth } from '@/types/berth'

const loading = ref(false)
const bookingLoading = ref(false)
const docks = ref<Dock[]>([])
const berths = ref<Berth[]>([])
const showDialog = ref(false)
const selectedBerth = ref<Berth | null>(null)

const form = reactive({
  dock_id: '',
  berth_type: ''
})

const bookingForm = reactive({
  start_time: '',
  end_time: '',
  notes: ''
})

const getBerthTypeText = (type: string) => {
  const map: Record<string, string> = {
    small: '小型',
    medium: '中型',
    large: '大型'
  }
  return map[type] || type
}

const getStatusText = (status: string) => {
  const map: Record<string, string> = {
    available: '可用',
    occupied: '已占用',
    reserved: '已预约',
    maintenance: '维护中'
  }
  return map[status] || status
}

const fetchDocks = async () => {
  try {
    const res: any = await getDocksApi({ page_size: 100 })
    docks.value = res.data || []
  } catch (error) {
    console.error('Failed to fetch docks:', error)
  }
}

const fetchBerths = async () => {
  if (!form.dock_id) {
    berths.value = []
    return
  }
  loading.value = true
  try {
    const params: any = { dock_id: form.dock_id }
    if (form.berth_type) params.berth_type = form.berth_type
    const res: any = await getBerthsApi(params)
    berths.value = res.data || []
  } catch (error) {
    console.error('Failed to fetch berths:', error)
  } finally {
    loading.value = false
  }
}

const showBookingDialog = (berth: Berth) => {
  selectedBerth.value = berth
  bookingForm.start_time = ''
  bookingForm.end_time = ''
  bookingForm.notes = ''
  showDialog.value = true
}

const handleBooking = async () => {
  if (!bookingForm.start_time || !bookingForm.end_time) {
    ElMessage.warning('请选择预约时间')
    return
  }
  bookingLoading.value = true
  try {
    await createReservationApi({
      berth_id: selectedBerth.value!.id,
      start_time: bookingForm.start_time,
      end_time: bookingForm.end_time,
      notes: bookingForm.notes
    })
    ElMessage.success('预约成功')
    showDialog.value = false
    fetchBerths()
  } catch (error) {
    ElMessage.error('预约失败')
  } finally {
    bookingLoading.value = false
  }
}

onMounted(fetchDocks)
</script>

<style lang="scss" scoped>
.berth-reservations {
  .card-container {
    background: #fff;
    border-radius: 8px;
    padding: 20px;
  }
}
</style>
