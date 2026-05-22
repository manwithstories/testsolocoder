<template>
  <div class="booking-create-container">
    <el-card shadow="hover">
      <template #header>
        <div class="card-header">
        <span class="card-title">创建预订</span>
        <el-button :icon="ArrowLeft" @click="goBack">返回列表</el-button>
      </div>
    </template>

    <el-steps :active="currentStep" finish-status="success" simple>
      <el-step title="选择房型" />
      <el-step title="选择房间" />
      <el-step title="填写信息" />
      <el-step title="确认预订" />
    </el-steps>

    <div class="step-content">
      <div v-show="currentStep === 0" class="step-1">
        <h3 class="step-title">请选择房型</h3>
        <div class="room-type-grid">
          <div
            v-for="type in roomTypeList"
            :key="type.id"
            class="room-type-card"
            :class="{ active: selectedRoomType?.id === type.id }"
            @click="selectRoomType(type)"
          >
            <div class="room-type-header">
              <span class="room-type-name">{{ type.name }}</span>
              <span class="room-type-price">¥{{ type.basePrice }}/晚</span>
            </div>
            <div class="room-type-info">
              <span><el-icon><House /></el-icon> {{ type.bedCount }}张床</span>
              <span><el-icon><User /></el-icon> 最多{{ type.maxGuests }}人</span>
            </div>
            <div class="room-type-facilities">
              <el-tag
                v-for="(facility, index) in type.facilities.slice(0, 4)"
                :key="index"
                size="small"
                type="info"
              >
                {{ facility }}
              </el-tag>
            </div>
          </div>
        </div>
      </div>

      <div v-show="currentStep === 1" class="step-2">
        <h3 class="step-title">请选择房间</h3>
        <div class="date-selector">
          <el-form :inline="true" label-width="80px">
            <el-form-item label="入住日期">
              <el-date-picker
                v-model="formData.checkInDate"
                type="date"
                placeholder="选择入住日期"
                value-format="YYYY-MM-DD"
                :disabled-date="disabledCheckInDate"
                @change="fetchAvailableRooms"
              />
            </el-form-item>
            <el-form-item label="离店日期">
              <el-date-picker
                v-model="formData.checkOutDate"
                type="date"
                placeholder="选择离店日期"
                value-format="YYYY-MM-DD"
                :disabled-date="disabledCheckOutDate"
                @change="fetchAvailableRooms"
              />
            </el-form-item>
          </el-form>
        </div>
        <div class="room-grid">
          <div
            v-for="room in availableRooms"
            :key="room.id"
            class="room-card"
            :class="{ active: selectedRoom?.id === room.id }"
            @click="selectRoom(room)"
          >
            <div class="room-number">{{ room.roomNumber }}</div>
            <div class="room-floor">{{ room.floor }}楼</div>
            <div class="room-status">
              <el-tag size="small" type="success">空闲</el-tag>
            </div>
          </div>
          <div v-if="!availableRooms.length && formData.checkInDate && formData.checkOutDate" class="no-rooms">
            <el-empty description="该时间段内没有可用房间" />
          </div>
        </div>
      </div>

      <div v-show="currentStep === 2" class="step-3">
        <h3 class="step-title">请填写客人信息</h3>
        <el-form
          ref="formRef"
          :model="formData"
          :rules="formRules"
          label-width="120px"
          class="guest-form"
        >
          <el-row :gutter="24">
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
          <el-row :gutter="24">
            <el-col :span="12">
              <el-form-item label="身份证号" prop="guestIdCard">
                <el-input v-model="formData.guestIdCard" placeholder="请输入身份证号" />
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="会员手机号">
                <el-input v-model="formData.memberPhone" placeholder="输入会员手机号（可选）" />
              </el-form-item>
            </el-col>
          </el-row>
          <el-row :gutter="24">
            <el-col :span="8">
              <el-form-item label="成人数量" prop="adults">
                <el-input-number
                  v-model="formData.adults"
                  :min="1"
                  :max="10"
                  style="width: 100%"
                />
              </el-form-item>
            </el-col>
            <el-col :span="8">
              <el-form-item label="儿童数量" prop="children">
                <el-input-number
                  v-model="formData.children"
                  :min="0"
                  :max="10"
                  style="width: 100%"
                />
              </el-form-item>
            </el-col>
            <el-col :span="8">
              <el-form-item label="预订来源" prop="source">
                <el-select v-model="formData.source" style="width: 100%">
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
          <el-row :gutter="24">
            <el-col :span="12">
              <el-form-item label="押金金额" prop="deposit">
                <el-input-number
                  v-model="formData.deposit"
                  :min="0"
                  :precision="2"
                  style="width: 100%"
                  placeholder="请输入押金金额"
                />
              </el-form-item>
            </el-col>
          </el-row>
          <el-form-item label="备注" prop="remark">
            <el-input
              v-model="formData.remark"
              type="textarea"
              :rows="3"
              placeholder="请输入备注"
            />
          </el-form-item>
        </el-form>
      </div>

      <div v-show="currentStep === 3" class="step-4">
        <h3 class="step-title">请确认预订信息</h3>
        <el-descriptions :column="2" border class="confirm-descriptions">
          <el-descriptions-item label="房型">{{ selectedRoomType?.name }}</el-descriptions-item>
          <el-descriptions-item label="房间号">{{ selectedRoom?.roomNumber }}</el-descriptions-item>
          <el-descriptions-item label="入住日期">{{ formData.checkInDate }}</el-descriptions-item>
          <el-descriptions-item label="离店日期">{{ formData.checkOutDate }}</el-descriptions-item>
          <el-descriptions-item label="入住天数">{{ formData.nights }}晚</el-descriptions-item>
          <el-descriptions-item label="入住人数">成人{{ formData.adults }}人，儿童{{ formData.children }}人</el-descriptions-item>
          <el-descriptions-item label="客人姓名">{{ formData.guestName }}</el-descriptions-item>
          <el-descriptions-item label="联系电话">{{ formData.guestPhone }}</el-descriptions-item>
          <el-descriptions-item label="身份证号">{{ formData.guestIdCard || '-' }}</el-descriptions-item>
          <el-descriptions-item label="预订来源">{{ formData.source }}</el-descriptions-item>
        </el-descriptions>

        <el-card class="price-card">
          <h4>价格明细</h4>
          <div class="price-row">
            <span>原价</span>
            <span>¥{{ originalPrice.toFixed(2) }}</span>
          </div>
          <div class="price-row">
            <span>会员折扣</span>
            <span class="discount">-¥{{ discountAmount.toFixed(2) }}</span>
          </div>
          <el-divider />
          <div class="price-row total">
            <span>应付金额</span>
            <span class="total-price">¥{{ finalPrice.toFixed(2) }}</span>
          </div>
        </el-card>
      </div>
    </div>

    <div class="step-actions">
      <el-button v-if="currentStep > 0" @click="prevStep">上一步</el-button>
      <el-button
        v-if="currentStep < 3"
        type="primary"
        :disabled="!canProceed"
        @click="nextStep"
      >
        下一步
      </el-button>
      <el-button
        v-if="currentStep === 3"
        type="success"
        :loading="submitLoading"
        @click="handleSubmit"
      >
        确认提交
      </el-button>
    </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { ArrowLeft, House, User, DataLine } from '@element-plus/icons-vue'
import dayjs from 'dayjs'
import { getAllRoomTypes, getAvailableRooms } from '@/api/room'
import { createBooking } from '@/api/booking'
import { type RoomType, type Room, type Booking } from '@/types'

const router = useRouter()
const formRef = ref<FormInstance>()
const submitLoading = ref(false)
const currentStep = ref(0)
const roomTypeList = ref<RoomType[]>([])
const availableRooms = ref<Room[]>([])
const selectedRoomType = ref<RoomType | null>(null)
const selectedRoom = ref<Room | null>(null)

const formData = reactive({
  guestName: '',
  guestPhone: '',
  guestIdCard: '',
  memberPhone: '',
  checkInDate: '',
  checkOutDate: '',
  nights: 0,
  adults: 1,
  children: 0,
  deposit: 0,
  source: 'front_desk',
  remark: ''
})

const formRules: FormRules = {
  guestName: [{ required: true, message: '请输入客人姓名', trigger: 'blur' }],
  guestPhone: [{ required: true, message: '请输入联系电话', trigger: 'blur' }],
  checkInDate: [{ required: true, message: '请选择入住日期', trigger: 'change' }],
  checkOutDate: [{ required: true, message: '请选择离店日期', trigger: 'change' }]
}

const canProceed = computed(() => {
  switch (currentStep.value) {
    case 0:
      return !!selectedRoomType.value
    case 1:
      return !!selectedRoom.value && formData.checkInDate && formData.checkOutDate
    case 2:
      return formData.guestName && formData.guestPhone
    default:
      return true
  }
})

const originalPrice = computed(() => {
  if (!selectedRoomType.value || !formData.nights) return 0
  return selectedRoomType.value.basePrice * formData.nights
})

const discountAmount = computed(() => {
  return 0
})

const finalPrice = computed(() => {
  return originalPrice.value - discountAmount.value
})

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
    roomTypeList.value = res
  } catch (error) {
    console.error('Failed to fetch room types:', error)
  }
}

const fetchAvailableRooms = async () => {
  if (!selectedRoomType.value || !formData.checkInDate || !formData.checkOutDate) {
    availableRooms.value = []
    return
  }

  const checkIn = dayjs(formData.checkInDate)
  const checkOut = dayjs(formData.checkOutDate)
  const nights = checkOut.diff(checkIn, 'day')

  if (nights <= 0) {
    ElMessage.warning('离店日期必须晚于入住日期')
    formData.nights = 0
    return
  }

  formData.nights = nights

  try {
    const res = await getAvailableRooms(formData.checkInDate, formData.checkOutDate, selectedRoomType.value.id)
    availableRooms.value = res
  } catch (error) {
    console.error('Failed to fetch available rooms:', error)
  }
}

const selectRoomType = (type: RoomType) => {
  selectedRoomType.value = type
  selectedRoom.value = null
  availableRooms.value = []
}

const selectRoom = (room: Room) => {
  selectedRoom.value = room
}

const nextStep = () => {
  if (currentStep.value === 2) {
    if (!formRef.value) return
    formRef.value.validate((valid) => {
      if (valid) {
        currentStep.value++
      }
    })
  } else {
    currentStep.value++
  }
}

const prevStep = () => {
  currentStep.value--
}

const goBack = () => {
  router.push('/booking/list')
}

const handleSubmit = async () => {
  submitLoading.value = true
  try {
    await createBooking({
      ...formData,
      roomTypeId: selectedRoomType.value!.id,
      roomId: selectedRoom.value!.id,
      totalPrice: finalPrice.value,
      createdBy: 1
    } as Omit<Booking, 'id' | 'bookingNo' | 'createdAt' | 'updatedAt'>)
    ElMessage.success('预订创建成功')
    router.push('/booking/list')
  } catch (error) {
    console.error('Failed to create booking:', error)
  } finally {
    submitLoading.value = false
  }
}

onMounted(() => {
  fetchRoomTypes()
})
</script>

<style scoped lang="scss">
.booking-create-container {
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;

    .card-title {
      font-size: 18px;
      font-weight: 600;
      color: #303133;
    }
  }

  .step-content {
    margin-top: 32px;

    .step-title {
      font-size: 16px;
      font-weight: 600;
      margin-bottom: 24px;
      color: #303133;
    }
  }

  .room-type-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
    gap: 16px;

    .room-type-card {
      padding: 20px;
      border: 2px solid #ebeef5;
      border-radius: 8px;
      cursor: pointer;
      transition: all 0.3s;

      &:hover {
        border-color: #409eff;
        box-shadow: 0 4px 12px rgba(64, 158, 255, 0.15);
      }

      &.active {
        border-color: #409eff;
        background: #ecf5ff;
      }

      .room-type-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 12px;

        .room-type-name {
          font-size: 16px;
          font-weight: 600;
          color: #303133;
        }

        .room-type-price {
          font-size: 18px;
          font-weight: 600;
          color: #f56c6c;
        }
      }

      .room-type-info {
        display: flex;
        gap: 16px;
        margin-bottom: 12px;
        color: #606266;
        font-size: 13px;

        span {
          display: flex;
          align-items: center;
          gap: 4px;
        }
      }

      .room-type-facilities {
        display: flex;
        flex-wrap: wrap;
        gap: 6px;
      }
    }
  }

  .date-selector {
    margin-bottom: 24px;
  }

  .room-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(120px, 1fr));
    gap: 12px;

    .room-card {
      padding: 20px 12px;
      border: 2px solid #ebeef5;
      border-radius: 8px;
      text-align: center;
      cursor: pointer;
      transition: all 0.3s;

      &:hover {
        border-color: #409eff;
        box-shadow: 0 4px 12px rgba(64, 158, 255, 0.15);
      }

      &.active {
        border-color: #409eff;
        background: #ecf5ff;
      }

      .room-number {
        font-size: 20px;
        font-weight: 600;
        color: #303133;
        margin-bottom: 8px;
      }

      .room-floor {
        font-size: 13px;
        color: #909399;
        margin-bottom: 8px;
      }
    }

    .no-rooms {
      grid-column: 1 / -1;
      padding: 40px 0;
    }
  }

  .guest-form {
    max-width: 800px;
  }

  .confirm-descriptions {
    margin-bottom: 24px;
  }

  .price-card {
    max-width: 500px;
    margin: 0 auto;

    h4 {
      margin: 0 0 16px 0;
      font-size: 16px;
      font-weight: 600;
    }

    .price-row {
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: 8px 0;
      font-size: 14px;
      color: #606266;

      &.total {
        font-size: 16px;
        font-weight: 600;
        color: #303133;

        .total-price {
          font-size: 24px;
          color: #f56c6c;
        }
      }

      .discount {
        color: #67c23a;
      }
    }
  }

  .step-actions {
    display: flex;
    justify-content: center;
    gap: 16px;
    margin-top: 32px;
    padding-top: 24px;
    border-top: 1px solid #ebeef5;
  }
}
</style>
