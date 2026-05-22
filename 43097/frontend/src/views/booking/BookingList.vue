<template>
  <div class="booking-container">
    <el-card shadow="hover">
      <div class="header-bar">
        <div class="filter-bar">
          <el-select
            v-model="filterStatus"
            placeholder="选择状态"
            clearable
            style="width: 150px"
            @change="fetchList"
          >
            <el-option label="待确认" :value="BookingStatus.PENDING" />
            <el-option label="已确认" :value="BookingStatus.CONFIRMED" />
            <el-option label="已取消" :value="BookingStatus.CANCELLED" />
            <el-option label="已完成" :value="BookingStatus.CHECKED_OUT" />
          </el-select>
          <el-date-picker
            v-model="dateRange"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            value-format="YYYY-MM-DD"
            style="width: 280px"
            @change="handleDateChange"
          />
          <el-input
            v-model="searchKeyword"
            placeholder="搜索预订号/客人姓名/电话"
            clearable
            :prefix-icon="Search"
            style="width: 240px"
            @keyup.enter="handleSearch"
          />
          <el-button type="primary" :icon="Search" @click="handleSearch">搜索</el-button>
          <el-button :icon="Refresh" @click="resetFilter">重置</el-button>
        </div>
        <el-button type="primary" :icon="Plus" @click="handleAdd">添加预订</el-button>
      </div>

      <el-table :data="tableData" v-loading="loading" border stripe>
        <el-table-column prop="bookingNo" label="预订号" width="140" />
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
        <el-table-column prop="guestName" label="客人姓名" width="100" />
        <el-table-column prop="guestPhone" label="电话" width="130" />
        <el-table-column prop="checkInDate" label="入住日期" width="120" />
        <el-table-column prop="checkOutDate" label="离店日期" width="120" />
        <el-table-column prop="days" label="天数" width="80">
          <template #default="{ row }">
            {{ row.days }}晚
          </template>
        </el-table-column>
        <el-table-column prop="totalPrice" label="总价" width="120">
          <template #default="{ row }">
            <span class="price-text">¥{{ row.totalPrice.toFixed(2) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="240" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link :icon="View" @click="handleView(row)">详情</el-button>
            <el-button
              v-if="row.status === BookingStatus.PENDING"
              type="success"
              link
              :icon="Check"
              @click="handleConfirm(row)"
            >
              确认
            </el-button>
            <el-button
              v-if="row.status === BookingStatus.PENDING || row.status === BookingStatus.CONFIRMED"
              type="danger"
              link
              :icon="Close"
              @click="handleCancel(row)"
            >
              取消
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-bar">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="pagination.total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="fetchList"
          @current-change="fetchList"
        />
      </div>
    </el-card>

    <el-dialog
      v-model="dialogVisible"
      :title="isEdit ? '编辑预订' : '添加预订'"
      width="700px"
      destroy-on-close
    >
      <el-form
        ref="formRef"
        :model="formData"
        :rules="formRules"
        label-width="100px"
        label-position="right"
      >
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="客人姓名" prop="guestName">
              <el-input v-model="formData.guestName" placeholder="请输入客人姓名" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="联系电话" prop="guestPhone">
              <el-input v-model="formData.guestPhone" placeholder="请输入联系电话" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="身份证号" prop="guestIdCard">
              <el-input v-model="formData.guestIdCard" placeholder="请输入身份证号" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="预订来源" prop="source">
              <el-select v-model="formData.source" placeholder="请选择预订来源" style="width: 100%">
                <el-option label="前台" value="front_desk" />
                <el-option label="电话" value="phone" />
                <el-option label="微信" value="wechat" />
                <el-option label="美团" value="meituan" />
                <el-option label="携程" value="ctrip" />
                <el-option label="其他" value="other" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="房型" prop="roomTypeId">
              <el-select
                v-model="formData.roomTypeId"
                placeholder="请选择房型"
                style="width: 100%"
                @change="handleRoomTypeChange"
              >
                <el-option
                  v-for="type in roomTypeOptions"
                  :key="type.id"
                  :label="`${type.name} (¥${type.price}/晚)`"
                  :value="type.id"
                />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="房间" prop="roomId">
              <el-select
                v-model="formData.roomId"
                placeholder="请选择房间"
                style="width: 100%"
                :disabled="!formData.roomTypeId"
              >
                <el-option
                  v-for="room in availableRooms"
                  :key="room.id"
                  :label="room.roomNumber"
                  :value="room.id"
                />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="入住日期" prop="checkInDate">
              <el-date-picker
                v-model="formData.checkInDate"
                type="date"
                placeholder="选择入住日期"
                value-format="YYYY-MM-DD"
                style="width: 100%"
                :disabled-date="disabledCheckInDate"
                @change="calculatePrice"
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="离店日期" prop="checkOutDate">
              <el-date-picker
                v-model="formData.checkOutDate"
                type="date"
                placeholder="选择离店日期"
                value-format="YYYY-MM-DD"
                style="width: 100%"
                :disabled-date="disabledCheckOutDate"
                @change="calculatePrice"
              />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="8">
            <el-form-item label="成人" prop="adults">
              <el-input-number
                v-model="formData.adults"
                :min="1"
                :max="10"
                style="width: 100%"
              />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="儿童" prop="children">
              <el-input-number
                v-model="formData.children"
                :min="0"
                :max="10"
                style="width: 100%"
              />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="押金" prop="deposit">
              <el-input-number
                v-model="formData.deposit"
                :min="0"
                :precision="2"
                style="width: 100%"
              />
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="备注" prop="remark">
          <el-input
            v-model="formData.remark"
            type="textarea"
            :rows="2"
            placeholder="请输入备注"
          />
        </el-form-item>
        <el-divider content-position="left">价格信息</el-divider>
        <el-descriptions :column="3" border size="small">
          <el-descriptions-item label="入住天数">
            {{ formData.days || 0 }}晚
          </el-descriptions-item>
          <el-descriptions-item label="每晚价格">
            ¥{{ selectedRoomType?.price?.toFixed(2) || '0.00' }}
          </el-descriptions-item>
          <el-descriptions-item label="总价">
            <span class="price-text">¥{{ formData.totalPrice?.toFixed(2) || '0.00' }}</span>
          </el-descriptions-item>
        </el-descriptions>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitLoading" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="detailVisible" title="预订详情" width="700px">
      <el-descriptions :column="2" border v-if="currentRow">
        <el-descriptions-item label="预订号">{{ currentRow.bookingNo }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="getStatusType(currentRow.status)">
            {{ getStatusText(currentRow.status) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="客人姓名">{{ currentRow.guestName }}</el-descriptions-item>
        <el-descriptions-item label="联系电话">{{ currentRow.guestPhone }}</el-descriptions-item>
        <el-descriptions-item label="身份证号">{{ currentRow.guestIdCard || '-' }}</el-descriptions-item>
        <el-descriptions-item label="预订来源">{{ currentRow.source }}</el-descriptions-item>
        <el-descriptions-item label="房型">{{ currentRow.room?.roomType?.name || '-' }}</el-descriptions-item>
        <el-descriptions-item label="房间">{{ currentRow.room?.roomNo || '-' }}</el-descriptions-item>
        <el-descriptions-item label="入住日期">{{ currentRow.checkInDate }}</el-descriptions-item>
        <el-descriptions-item label="离店日期">{{ currentRow.checkOutDate }}</el-descriptions-item>
        <el-descriptions-item label="入住天数">{{ currentRow.days }}晚</el-descriptions-item>
        <el-descriptions-item label="入住人数">成人{{ currentRow.adults || 1 }}人，儿童{{ currentRow.children || 0 }}人</el-descriptions-item>
        <el-descriptions-item label="总价">
          <span class="price-text">¥{{ currentRow.totalPrice.toFixed(2) }}</span>
        </el-descriptions-item>
        <el-descriptions-item label="押金">¥{{ currentRow.deposit.toFixed(2) }}</el-descriptions-item>
        <el-descriptions-item label="备注" :span="2">{{ currentRow.remark || '-' }}</el-descriptions-item>
        <el-descriptions-item label="创建时间" :span="2">{{ currentRow.createdAt }}</el-descriptions-item>
      </el-descriptions>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import { Search, Refresh, Plus, View, Check, Close } from '@element-plus/icons-vue'
import dayjs from 'dayjs'
import {
  getBookingList,
  createBooking,
  confirmBooking,
  cancelBooking
} from '@/api/booking'
import { getAllRoomTypes, getAvailableRooms } from '@/api/room'
import { BookingStatus, type Booking, type RoomType, type Room, type PageParams } from '@/types'

const loading = ref(false)
const submitLoading = ref(false)
const searchKeyword = ref('')
const filterStatus = ref<BookingStatus | ''>('')
const dateRange = ref<string[]>([])
const tableData = ref<Booking[]>([])
const roomTypeOptions = ref<RoomType[]>([])
const availableRooms = ref<Room[]>([])
const dialogVisible = ref(false)
const detailVisible = ref(false)
const isEdit = ref(false)
const currentRow = ref<Booking | null>(null)
const formRef = ref<FormInstance>()

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

const formData = reactive<Partial<Booking>>({
  guestName: '',
  guestPhone: '',
  guestIdCard: '',
  roomTypeId: undefined,
  roomId: undefined,
  checkInDate: '',
  checkOutDate: '',
  days: 0,
  adults: 1,
  children: 0,
  totalPrice: 0,
  deposit: 0,
  source: 'front_desk',
  remarks: ''
})

const formRules: FormRules = {
  guestName: [{ required: true, message: '请输入客人姓名', trigger: 'blur' }],
  guestPhone: [{ required: true, message: '请输入联系电话', trigger: 'blur' }],
  roomTypeId: [{ required: true, message: '请选择房型', trigger: 'change' }],
  checkInDate: [{ required: true, message: '请选择入住日期', trigger: 'change' }],
  checkOutDate: [{ required: true, message: '请选择离店日期', trigger: 'change' }]
}

const selectedRoomType = computed(() => {
  if (!formData.roomTypeId) return null
  return roomTypeOptions.value.find((t) => t.id === formData.roomTypeId)
})

const getStatusType = (status: BookingStatus) => {
  const typeMap: Record<BookingStatus, string> = {
    [BookingStatus.PENDING]: 'warning',
    [BookingStatus.CONFIRMED]: 'primary',
    [BookingStatus.CHECKED_IN]: 'success',
    [BookingStatus.CHECKED_OUT]: 'info',
    [BookingStatus.CANCELLED]: 'danger',
    [BookingStatus.EXPIRED]: 'info'
  }
  return typeMap[status] || 'info'
}

const getStatusText = (status: BookingStatus) => {
  const textMap: Record<BookingStatus, string> = {
    [BookingStatus.PENDING]: '待确认',
    [BookingStatus.CONFIRMED]: '已确认',
    [BookingStatus.CHECKED_IN]: '已入住',
    [BookingStatus.CHECKED_OUT]: '已完成',
    [BookingStatus.CANCELLED]: '已取消',
    [BookingStatus.EXPIRED]: '已过期'
  }
  return textMap[status] || '未知'
}

const disabledCheckInDate = (date: Date) => {
  return date && dayjs(date).isBefore(dayjs().subtract(1, 'day'), 'day')
}

const disabledCheckOutDate = (date: Date) => {
  if (!formData.checkInDate) return false
  return date && dayjs(date).isBefore(formData.checkInDate, 'day')
}

const fetchRoomTypes = async () => {
  try {
    const res = await getAllRoomTypes()
    roomTypeOptions.value = res
  } catch (error) {
    console.error('Failed to fetch room types:', error)
  }
}

const fetchList = async () => {
  loading.value = true
  try {
    const params: PageParams & { status?: BookingStatus; startDate?: string; endDate?: string; keyword?: string } = {
      page: pagination.page,
      pageSize: pagination.pageSize
    }
    if (filterStatus.value) {
      params.status = filterStatus.value
    }
    if (dateRange.value && dateRange.value.length === 2) {
      params.startDate = dateRange.value[0]
      params.endDate = dateRange.value[1]
    }
    if (searchKeyword.value) {
      params.keyword = searchKeyword.value
    }
    const res = await getBookingList(params)
    tableData.value = res.list
    pagination.total = res.total
  } catch (error) {
    console.error('Failed to fetch booking list:', error)
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  pagination.page = 1
  fetchList()
}

const handleDateChange = () => {
  pagination.page = 1
  fetchList()
}

const resetFilter = () => {
  searchKeyword.value = ''
  filterStatus.value = ''
  dateRange.value = []
  pagination.page = 1
  fetchList()
}

const handleAdd = () => {
  isEdit.value = false
  Object.assign(formData, {
    guestName: '',
    guestPhone: '',
    guestIdCard: '',
    roomTypeId: undefined,
    roomId: undefined,
    checkInDate: '',
    checkOutDate: '',
    days: 0,
    adults: 1,
    children: 0,
    totalPrice: 0,
    deposit: 0,
    source: 'front_desk',
    remarks: ''
  })
  availableRooms.value = []
  dialogVisible.value = true
}

const handleView = (row: Booking) => {
  currentRow.value = row
  detailVisible.value = true
}

const handleConfirm = (row: Booking) => {
  ElMessageBox.confirm(`确定要确认预订"${row.bookingNo}"吗？`, '确认预订', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  })
    .then(async () => {
      try {
        await confirmBooking(row.id)
        ElMessage.success('确认成功')
        fetchList()
      } catch (error) {
        console.error('Failed to confirm booking:', error)
      }
    })
    .catch(() => {})
}

const handleCancel = (row: Booking) => {
  ElMessageBox.prompt('请输入取消原因（可选）', '取消预订', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning',
    inputPlaceholder: '请输入取消原因'
  })
    .then(async ({ value }) => {
      try {
        await cancelBooking(row.id, value)
        ElMessage.success('取消成功')
        fetchList()
      } catch (error) {
        console.error('Failed to cancel booking:', error)
      }
    })
    .catch(() => {})
}

const handleRoomTypeChange = async () => {
  formData.roomId = undefined
  if (formData.roomTypeId && formData.checkInDate && formData.checkOutDate) {
    try {
      const res = await getAvailableRooms(formData.checkInDate, formData.checkOutDate, formData.roomTypeId)
      availableRooms.value = res
    } catch (error) {
      console.error('Failed to fetch available rooms:', error)
    }
  }
}

const calculatePrice = () => {
  if (!formData.checkInDate || !formData.checkOutDate || !selectedRoomType.value) {
    formData.days = 0
    formData.totalPrice = 0
    return
  }

  const checkIn = dayjs(formData.checkInDate)
  const checkOut = dayjs(formData.checkOutDate)
  const days = checkOut.diff(checkIn, 'day')

  if (days <= 0) {
    ElMessage.warning('离店日期必须晚于入住日期')
    formData.days = 0
    formData.totalPrice = 0
    return
  }

  formData.days = days
  formData.totalPrice = days * selectedRoomType.value.price
}

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (valid) {
      submitLoading.value = true
      try {
        await createBooking({
          ...formData,
          createdBy: 1
        } as Omit<Booking, 'id' | 'bookingNo' | 'createdAt' | 'updatedAt'>)
        ElMessage.success('添加成功')
        dialogVisible.value = false
        fetchList()
      } catch (error) {
        console.error('Failed to submit booking:', error)
      } finally {
        submitLoading.value = false
      }
    }
  })
}

onMounted(() => {
  fetchRoomTypes()
  fetchList()
})
</script>

<style scoped lang="scss">
.booking-container {
  .header-bar {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;
    flex-wrap: wrap;
    gap: 12px;

    .filter-bar {
      display: flex;
      gap: 12px;
      flex-wrap: wrap;
    }
  }

  .pagination-bar {
    display: flex;
    justify-content: flex-end;
    margin-top: 20px;
  }

  .price-text {
    color: #f56c6c;
    font-weight: 600;
  }
}
</style>
