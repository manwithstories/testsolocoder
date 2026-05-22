<template>
  <div class="bookings">
    <el-card class="filter-card">
      <el-form :inline="true" :model="filterForm">
        <el-form-item label="会议室">
          <el-select v-model="filterForm.room_id" placeholder="全部会议室" clearable @change="loadBookings">
            <el-option v-for="room in rooms" :key="room.id" :label="room.name" :value="room.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="filterForm.status" placeholder="全部状态" clearable @change="loadBookings">
            <el-option label="待确认" :value="0" />
            <el-option label="已确认" :value="1" />
            <el-option label="已取消" :value="2" />
            <el-option label="已完成" :value="3" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadBookings">
            <el-icon><Search /></el-icon>
            搜索
          </el-button>
        </el-form-item>
        <el-form-item>
          <el-button type="success" @click="showCreateDialog">
            <el-icon><Plus /></el-icon>
            新建预订
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card style="margin-top: 20px; border-radius: 8px">
      <el-table :data="bookings" style="width: 100%" v-loading="loading">
        <el-table-column prop="title" label="会议名称" min-width="150" />
        <el-table-column prop="room.name" label="会议室" width="120" />
        <el-table-column label="时间" width="200">
          <template #default="{ row }">
            {{ formatDateTime(row.start_time) }}<br />
            <span style="color: #909399; font-size: 12px">{{ formatTime(row.start_time) }} - {{ formatTime(row.end_time) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="user.real_name" label="预订人" width="100" />
        <el-table-column prop="total_price" label="费用" width="100">
          <template #default="{ row }">
            ¥{{ row.total_price }}
          </template>
        </el-table-column>
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">{{ getStatusText(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="240">
          <template #default="{ row }">
            <el-button type="primary" size="small" @click="showDetail(row)">详情</el-button>
            <el-button v-if="row.status === 1 && canCancel(row)" type="warning" size="small" @click="cancelBooking(row)">取消</el-button>
            <el-button v-if="row.status === 1 && canCancel(row)" size="small" @click="showRescheduleDialog(row)">改期</el-button>
            <el-button v-if="isAdmin && row.status === 0" type="success" size="small" @click="approveBooking(row)">审核</el-button>
            <el-button v-if="isAdmin && row.status === 1 && isPast(row)" type="success" size="small" @click="completeBooking(row)">完成</el-button>
          </template>
        </el-table-column>
      </el-table>
      <el-pagination
        v-if="total > 0"
        style="margin-top: 16px; justify-content: flex-end"
        layout="total, prev, pager, next"
        :total="total"
        :current-page="filterForm.page"
        :page-size="filterForm.page_size"
        @current-change="handlePageChange"
      />
    </el-card>

    <el-dialog v-model="createDialogVisible" title="新建预订" width="600px">
      <el-form ref="createFormRef" :model="createForm" :rules="createRules" label-width="100px">
        <el-form-item label="会议名称" prop="title">
          <el-input v-model="createForm.title" placeholder="请输入会议名称" />
        </el-form-item>
        <el-form-item label="会议室" prop="room_id">
          <el-select v-model="createForm.room_id" placeholder="选择会议室" style="width: 100%" @change="onRoomChange">
            <el-option v-for="room in rooms" :key="room.id" :label="`${room.name} (${room.capacity}人, ¥${room.price_per_hour}/h)`" :value="room.id" />
          </el-select>
        </el-form-item>
        <el-alert
          v-if="selectedRoom"
          :title="`可用时段: ${selectedRoom.available_start || '08:00'} - ${selectedRoom.available_end || '22:00'}`"
          type="info"
          :closable="false"
          style="margin-bottom: 18px"
        />
        <el-form-item label="开始时间" prop="start_time">
          <el-date-picker
            v-model="createForm.start_time"
            type="datetime"
            placeholder="选择开始时间"
            format="YYYY-MM-DD HH:mm"
            value-format="YYYY-MM-DDTHH:mm:ss"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="结束时间" prop="end_time">
          <el-date-picker
            v-model="createForm.end_time"
            type="datetime"
            placeholder="选择结束时间"
            format="YYYY-MM-DD HH:mm"
            value-format="YYYY-MM-DDTHH:mm:ss"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="重复类型">
          <el-select v-model="createForm.recurrence_type" placeholder="不重复" style="width: 100%">
            <el-option label="不重复" value="none" />
            <el-option label="每天" value="daily" />
            <el-option label="每周" value="weekly" />
            <el-option label="每两周" value="biweekly" />
            <el-option label="每月" value="monthly" />
          </el-select>
        </el-form-item>
        <el-form-item v-if="createForm.recurrence_type && createForm.recurrence_type !== 'none'" label="重复结束">
          <el-date-picker
            v-model="createForm.recurrence_end"
            type="date"
            placeholder="选择重复结束日期"
            format="YYYY-MM-DD"
            value-format="YYYY-MM-DDTHH:mm:ss"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="createForm.description" type="textarea" :rows="2" placeholder="请输入会议描述" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="createDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="createBooking">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="detailDialogVisible" title="预订详情" width="500px">
      <el-descriptions :column="1" border v-if="currentBooking">
        <el-descriptions-item label="会议名称">{{ currentBooking.title }}</el-descriptions-item>
        <el-descriptions-item label="会议室">{{ currentBooking.room?.name }}</el-descriptions-item>
        <el-descriptions-item label="预订人">{{ currentBooking.user?.real_name || currentBooking.user?.username }}</el-descriptions-item>
        <el-descriptions-item label="时间">
          {{ formatDateTime(currentBooking.start_time) }} - {{ formatTime(currentBooking.end_time) }}
        </el-descriptions-item>
        <el-descriptions-item label="费用">¥{{ currentBooking.total_price }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="getStatusType(currentBooking.status)">{{ getStatusText(currentBooking.status) }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item v-if="currentBooking.description" label="描述">{{ currentBooking.description }}</el-descriptions-item>
      </el-descriptions>
      <template #footer>
        <el-button @click="detailDialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="rescheduleDialogVisible" title="改期" width="500px">
      <el-form :model="rescheduleForm" label-width="100px">
        <el-form-item label="新开始时间">
          <el-date-picker
            v-model="rescheduleForm.start_time"
            type="datetime"
            placeholder="选择新的开始时间"
            format="YYYY-MM-DD HH:mm"
            value-format="YYYY-MM-DDTHH:mm:ss"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="新结束时间">
          <el-date-picker
            v-model="rescheduleForm.end_time"
            type="datetime"
            placeholder="选择新的结束时间"
            format="YYYY-MM-DD HH:mm"
            value-format="YYYY-MM-DDTHH:mm:ss"
            style="width: 100%"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="rescheduleDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="rescheduleBooking">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import { api } from '@/api'
import { useUserStore } from '@/stores/user'
import dayjs from 'dayjs'

const userStore = useUserStore()
const isAdmin = computed(() => userStore.isAdmin)
const userId = computed(() => userStore.userInfo?.id)

const loading = ref(false)
const submitting = ref(false)
const bookings = ref<any[]>([])
const rooms = ref<any[]>([])
const total = ref(0)
const createDialogVisible = ref(false)
const detailDialogVisible = ref(false)
const rescheduleDialogVisible = ref(false)
const currentBooking = ref<any>(null)
const createFormRef = ref<FormInstance>()

const filterForm = reactive({
  page: 1,
  page_size: 10,
  room_id: undefined as number | undefined,
  status: undefined as number | undefined
})

const createForm = reactive({
  title: '',
  room_id: undefined as number | undefined,
  start_time: '',
  end_time: '',
  recurrence_type: 'none',
  recurrence_end: '',
  description: ''
})

const createRules: FormRules = {
  title: [{ required: true, message: '请输入会议名称', trigger: 'blur' }],
  room_id: [{ required: true, message: '请选择会议室', trigger: 'change' }],
  start_time: [{ required: true, message: '请选择开始时间', trigger: 'change' }],
  end_time: [{ required: true, message: '请选择结束时间', trigger: 'change' }]
}

const rescheduleForm = reactive({
  start_time: '',
  end_time: ''
})

const selectedRoom = computed(() => {
  if (!createForm.room_id) return null
  return rooms.value.find(r => r.id === createForm.room_id)
})

onMounted(() => {
  loadRooms()
  loadBookings()
})

async function loadRooms() {
  try {
    const res: any = await api.listAllRooms()
    rooms.value = res.data || []
  } catch (e) {
    console.error(e)
  }
}

async function loadBookings() {
  loading.value = true
  try {
    const res: any = await api.getBookings(filterForm)
    bookings.value = res.data?.bookings || []
    total.value = res.data?.total || 0
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

function handlePageChange(page: number) {
  filterForm.page = page
  loadBookings()
}

function showCreateDialog() {
  Object.assign(createForm, {
    title: '',
    room_id: undefined,
    start_time: '',
    end_time: '',
    recurrence_type: 'none',
    recurrence_end: '',
    description: ''
  })
  createDialogVisible.value = true
}

function onRoomChange() {}

async function createBooking() {
  if (!createFormRef.value) return
  await createFormRef.value.validate(async (valid) => {
    if (valid) {
      submitting.value = true
      try {
        const data: any = { ...createForm }
        if (data.recurrence_type === 'none') {
          delete data.recurrence_end
        }
        await api.createBooking(data)
        ElMessage.success('预订成功')
        createDialogVisible.value = false
        loadBookings()
      } catch (e: any) {
        console.error(e)
      } finally {
        submitting.value = false
      }
    }
  })
}

function showDetail(row: any) {
  currentBooking.value = row
  detailDialogVisible.value = true
}

async function cancelBooking(row: any) {
  try {
    const { value } = await ElMessageBox.prompt('请输入取消原因', '取消预订', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      inputPlaceholder: '取消原因'
    })
    await api.cancelBooking(row.id, value || '')
    ElMessage.success('取消成功')
    loadBookings()
  } catch (e: any) {
    if (e !== 'cancel') console.error(e)
  }
}

function showRescheduleDialog(row: any) {
  currentBooking.value = row
  rescheduleForm.start_time = row.start_time
  rescheduleForm.end_time = row.end_time
  rescheduleDialogVisible.value = true
}

async function rescheduleBooking() {
  submitting.value = true
  try {
    await api.rescheduleBooking(currentBooking.value.id, rescheduleForm)
    ElMessage.success('改期成功')
    rescheduleDialogVisible.value = false
    loadBookings()
  } catch (e: any) {
    console.error(e)
  } finally {
    submitting.value = false
  }
}

async function approveBooking(row: any) {
  try {
    await api.approveBooking(row.id)
    ElMessage.success('审核通过')
    loadBookings()
  } catch (e: any) {
    console.error(e)
  }
}

async function completeBooking(row: any) {
  try {
    await api.completeBooking(row.id)
    ElMessage.success('已标记完成')
    loadBookings()
  } catch (e: any) {
    console.error(e)
  }
}

function canCancel(row: any) {
  return dayjs(row.start_time).isAfter(dayjs())
}

function isPast(row: any) {
  return dayjs(row.end_time).isBefore(dayjs())
}

function formatDateTime(date: string) {
  return dayjs(date).format('YYYY-MM-DD HH:mm')
}

function formatTime(date: string) {
  return dayjs(date).format('HH:mm')
}

function getStatusType(status: number) {
  const map: Record<number, string> = { 0: 'warning', 1: 'success', 2: 'info', 3: 'primary' }
  return map[status] || 'info'
}

function getStatusText(status: number) {
  const map: Record<number, string> = { 0: '待确认', 1: '已确认', 2: '已取消', 3: '已完成' }
  return map[status] || '未知'
}
</script>

<style scoped>
.filter-card {
  border-radius: 8px;
}
</style>
