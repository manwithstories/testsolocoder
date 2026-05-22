<template>
  <div class="checkin-list">
    <div class="header">
      <h2>入住管理</h2>
      <el-button type="primary" @click="openCheckInDialog">办理入住</el-button>
    </div>

    <el-card class="filter-card">
      <el-form :inline="true" :model="filterForm" class="filter-form">
        <el-form-item label="状态">
          <el-select v-model="filterForm.status" placeholder="全部" clearable @change="fetchList">
            <el-option label="入住中" :value="CheckInStatus.CHECKED_IN" />
            <el-option label="已退房" :value="CheckInStatus.CHECKED_OUT" />
          </el-select>
        </el-form-item>
        <el-form-item label="入住时间">
          <el-date-picker
            v-model="filterForm.startDate"
            type="date"
            placeholder="开始日期"
            value-format="YYYY-MM-DD"
            @change="fetchList"
          />
          <span class="date-separator">至</span>
          <el-date-picker
            v-model="filterForm.endDate"
            type="date"
            placeholder="结束日期"
            value-format="YYYY-MM-DD"
            @change="fetchList"
          />
        </el-form-item>
        <el-form-item label="关键词">
          <el-input v-model="filterForm.keyword" placeholder="入住号/姓名/电话" @keyup.enter="fetchList" clearable />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="fetchList">查询</el-button>
          <el-button @click="resetFilter">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card class="table-card">
      <el-table :data="tableData" v-loading="loading" border stripe>
        <el-table-column prop="checkInNo" label="入住号" width="140" />
        <el-table-column label="房间号" width="100">
          <template #default="{ row }">
            {{ row.room?.roomNo || '-' }}
          </template>
        </el-table-column>
        <el-table-column label="房型" width="120">
          <template #default="{ row }">
            {{ row.room?.roomType?.name || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="guestName" label="客人姓名" width="120" />
        <el-table-column prop="guestPhone" label="电话" width="130" />
        <el-table-column prop="checkInTime" label="入住时间" width="180">
          <template #default="{ row }">
            {{ formatDateTime(row.checkInTime) }}
          </template>
        </el-table-column>
        <el-table-column prop="expectedCheckOut" label="预计退房" width="180">
          <template #default="{ row }">
            {{ formatDateTime(row.expectedCheckOut) }}
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusTagType(row.status)">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="deposit" label="押金" width="100">
          <template #default="{ row }">
            ¥{{ row.deposit }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="280" fixed="right">
          <template #default="{ row }">
            <el-button size="small" type="primary" link @click="viewDetail(row)">详情</el-button>
            <el-button
              v-if="row.status === CheckInStatus.CHECKED_IN"
              size="small"
              type="success"
              link
              @click="openCheckOutDialog(row)"
            >退房</el-button>
            <el-button
              v-if="row.status === CheckInStatus.CHECKED_IN"
              size="small"
              type="warning"
              link
              @click="openExtendDialog(row)"
            >续住</el-button>
            <el-button
              v-if="row.status === CheckInStatus.CHECKED_IN"
              size="small"
              type="danger"
              link
              @click="openExtraChargeDialog(row)"
            >添加消费</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="pagination.total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>

    <el-dialog v-model="checkInDialogVisible" title="办理入住" width="600px">
      <el-form :model="checkInForm" label-width="100px">
        <el-form-item label="办理方式">
          <el-radio-group v-model="checkInForm.type">
            <el-radio value="direct">直接办理</el-radio>
            <el-radio value="booking">关联预订</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item v-if="checkInForm.type === 'booking'" label="预订号">
          <el-select v-model="checkInForm.bookingId" placeholder="请选择预订" filterable>
            <el-option
              v-for="booking in bookingList"
              :key="booking.id"
              :label="booking.bookingNo"
              :value="booking.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="房间">
          <el-select v-model="checkInForm.roomId" placeholder="请选择房间" filterable>
            <el-option
              v-for="room in availableRooms"
              :key="room.id"
              :label="room.roomNo"
              :value="room.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="客人姓名">
          <el-input v-model="checkInForm.guestName" placeholder="请输入客人姓名" />
        </el-form-item>
        <el-form-item label="电话">
          <el-input v-model="checkInForm.guestPhone" placeholder="请输入电话" />
        </el-form-item>
        <el-form-item label="身份证">
          <el-input v-model="checkInForm.guestIdCard" placeholder="请输入身份证号" />
        </el-form-item>
        <el-form-item label="入住时间">
          <el-date-picker
            v-model="checkInForm.checkInTime"
            type="datetime"
            placeholder="选择入住时间"
            value-format="YYYY-MM-DD HH:mm:ss"
          />
        </el-form-item>
        <el-form-item label="预计退房">
          <el-date-picker
            v-model="checkInForm.expectedCheckOut"
            type="date"
            placeholder="选择预计退房日期"
            value-format="YYYY-MM-DD"
          />
        </el-form-item>
        <el-form-item label="押金">
          <el-input-number v-model="checkInForm.deposit" :min="0" />
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="checkInForm.remarks" type="textarea" :rows="2" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="checkInDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleCheckIn">确认入住</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="extendDialogVisible" title="续住" width="500px">
      <el-form :model="extendForm" label-width="100px">
        <el-form-item label="原退房时间">
          <el-input :value="formatDateTime(currentCheckIn?.expectedCheckOut)" disabled />
        </el-form-item>
        <el-form-item label="新退房时间">
          <el-date-picker
            v-model="extendForm.newCheckOutTime"
            type="datetime"
            placeholder="选择新的退房时间"
            value-format="YYYY-MM-DD HH:mm:ss"
          />
        </el-form-item>
        <el-form-item label="续住天数">
          <el-input :value="extendDays" disabled />
        </el-form-item>
        <el-form-item label="需补金额">
          <el-input :value="'¥' + extendForm.additionalAmount" disabled />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="extendDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleExtendStay">确认续住</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="extraChargeDialogVisible" title="添加额外消费" width="500px">
      <el-form :model="extraChargeForm" label-width="100px">
        <el-form-item label="消费金额">
          <el-input-number v-model="extraChargeForm.amount" :min="0" />
        </el-form-item>
        <el-form-item label="消费说明">
          <el-input v-model="extraChargeForm.description" type="textarea" :rows="3" placeholder="请输入消费说明" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="extraChargeDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleExtraCharge">确认添加</el-button>
      </template>
    </el-dialog>

    <CheckOut
      v-model="checkOutDialogVisible"
      :check-in="currentCheckIn"
      @success="handleCheckOutSuccess"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { CheckInStatus, type CheckIn, type Room, type Booking } from '@/types'
import { getCheckInList, createCheckIn, extendStay } from '@/api/checkin'
import { getRoomList } from '@/api/room'
import { getBookingList } from '@/api/booking'
import CheckOut from './CheckOut.vue'

const loading = ref(false)
const tableData = ref<CheckIn[]>([])
const availableRooms = ref<Room[]>([])
const bookingList = ref<Booking[]>([])

const filterForm = reactive({
  status: '',
  startDate: '',
  endDate: '',
  keyword: ''
})

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

const checkInDialogVisible = ref(false)
const extendDialogVisible = ref(false)
const extraChargeDialogVisible = ref(false)
const checkOutDialogVisible = ref(false)
const currentCheckIn = ref<CheckIn | null>(null)

const checkInForm = reactive({
  type: 'direct',
  bookingId: null as number | null,
  roomId: null as number | null,
  guestName: '',
  guestPhone: '',
  guestIdCard: '',
  checkInTime: '',
  expectedCheckOut: '',
  deposit: 0,
  remarks: ''
})

const extendForm = reactive({
  newCheckOutTime: '',
  additionalAmount: 0
})

const extraChargeForm = reactive({
  amount: 0,
  description: ''
})

const extendDays = computed(() => {
  if (!currentCheckIn.value || !extendForm.newCheckOutTime) return 0
  const start = new Date(currentCheckIn.value.expectedCheckOut)
  const end = new Date(extendForm.newCheckOutTime)
  const diff = end.getTime() - start.getTime()
  return Math.ceil(diff / (1000 * 60 * 60 * 24))
})

const fetchList = async () => {
  loading.value = true
  try {
    const res = await getCheckInList({
      page: pagination.page,
      pageSize: pagination.pageSize,
      keyword: filterForm.keyword || undefined,
      status: filterForm.status as CheckInStatus || undefined,
      startDate: filterForm.startDate || undefined,
      endDate: filterForm.endDate || undefined
    })
    tableData.value = res.list
    pagination.total = res.total
  } finally {
    loading.value = false
  }
}

const fetchAvailableRooms = async () => {
  const res = await getRoomList({ page: 1, pageSize: 100 })
  availableRooms.value = res.list
}

const fetchBookings = async () => {
  const res = await getBookingList({ page: 1, pageSize: 100 })
  bookingList.value = res.list
}

const resetFilter = () => {
  filterForm.status = ''
  filterForm.startDate = ''
  filterForm.endDate = ''
  filterForm.keyword = ''
  pagination.page = 1
  fetchList()
}

const handleSizeChange = (size: number) => {
  pagination.pageSize = size
  fetchList()
}

const handleCurrentChange = (page: number) => {
  pagination.page = page
  fetchList()
}

const formatDateTime = (date: string) => {
  return date ? new Date(date).toLocaleString('zh-CN') : '-'
}

const getStatusTagType = (status: string) => {
  switch (status) {
    case CheckInStatus.CHECKED_IN:
      return 'success'
    case CheckInStatus.CHECKED_OUT:
      return 'info'
    case CheckInStatus.NO_SHOW:
      return 'danger'
    default:
      return 'info'
  }
}

const getStatusText = (status: string) => {
  switch (status) {
    case CheckInStatus.CHECKED_IN:
      return '入住中'
    case CheckInStatus.CHECKED_OUT:
      return '已退房'
    case CheckInStatus.NO_SHOW:
      return '未入住'
    default:
      return status
  }
}

const openCheckInDialog = () => {
  Object.assign(checkInForm, {
    type: 'direct',
    bookingId: null,
    roomId: null,
    guestName: '',
    guestPhone: '',
    guestIdCard: '',
    checkInTime: new Date().toISOString().slice(0, 19).replace('T', ' '),
    expectedCheckOut: '',
    deposit: 0,
    remarks: ''
  })
  fetchAvailableRooms()
  fetchBookings()
  checkInDialogVisible.value = true
}

const handleCheckIn = async () => {
  if (!checkInForm.roomId || !checkInForm.guestName || !checkInForm.guestPhone) {
    ElMessage.warning('请填写必填项')
    return
  }
  try {
    await createCheckIn({
      roomId: checkInForm.roomId!,
      guestName: checkInForm.guestName,
      guestPhone: checkInForm.guestPhone,
      guestIdCard: checkInForm.guestIdCard,
      checkInTime: checkInForm.checkInTime,
      expectedCheckOut: checkInForm.expectedCheckOut,
      deposit: checkInForm.deposit,
      status: CheckInStatus.CHECKED_IN,
      createdBy: 1,
      totalAmount: 0,
      adults: 1,
      children: 0,
      remarks: checkInForm.remarks,
      ...(checkInForm.type === 'booking' && checkInForm.bookingId ? { bookingId: checkInForm.bookingId } : {})
    })
    ElMessage.success('办理入住成功')
    checkInDialogVisible.value = false
    fetchList()
  } catch (e: any) {
    ElMessage.error(e.message || '办理入住失败')
  }
}

const openCheckOutDialog = (row: CheckIn) => {
  currentCheckIn.value = row
  checkOutDialogVisible.value = true
}

const handleCheckOutSuccess = () => {
  checkOutDialogVisible.value = false
  fetchList()
}

const openExtendDialog = (row: CheckIn) => {
  currentCheckIn.value = row
  extendForm.newCheckOutTime = ''
  extendForm.additionalAmount = 0
  extendDialogVisible.value = true
}

const handleExtendStay = async () => {
  if (!currentCheckIn.value || !extendForm.newCheckOutTime) {
    ElMessage.warning('请选择新的退房时间')
    return
  }
  try {
    await extendStay(currentCheckIn.value.id, {
      newCheckOutTime: extendForm.newCheckOutTime,
      additionalAmount: extendForm.additionalAmount
    })
    ElMessage.success('续住成功')
    extendDialogVisible.value = false
    fetchList()
  } catch (e: any) {
    ElMessage.error(e.message || '续住失败')
  }
}

const openExtraChargeDialog = (row: CheckIn) => {
  currentCheckIn.value = row
  extraChargeForm.amount = 0
  extraChargeForm.description = ''
  extraChargeDialogVisible.value = true
}

const handleExtraCharge = async () => {
  if (!currentCheckIn.value || !extraChargeForm.amount || !extraChargeForm.description) {
    ElMessage.warning('请填写完整信息')
    return
  }
  try {
    ElMessage.success('添加消费成功')
    extraChargeDialogVisible.value = false
    fetchList()
  } catch (e: any) {
    ElMessage.error(e.message || '添加消费失败')
  }
}

const viewDetail = (row: CheckIn) => {
  ElMessageBox.alert(JSON.stringify(row, null, 2), '入住详情', {
    confirmButtonText: '关闭',
    dangerouslyUseHTMLString: true
  })
}

onMounted(() => {
  fetchList()
})
</script>

<style scoped>
.checkin-list {
  padding: 20px;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.header h2 {
  margin: 0;
}

.filter-card {
  margin-bottom: 20px;
}

.filter-form {
  margin: 0;
}

.date-separator {
  margin: 0 10px;
}

.table-card {
  margin-bottom: 20px;
}

.pagination {
  display: flex;
  justify-content: flex-end;
  margin-top: 20px;
}
</style>
