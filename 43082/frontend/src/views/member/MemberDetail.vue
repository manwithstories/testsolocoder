<template>
  <div class="member-detail">
    <el-row :gutter="20">
      <el-col :span="8">
        <el-card>
          <template #header>
            <span style="font-weight: 600">基本信息</span>
          </template>
          <div class="profile-section" v-if="member">
            <div class="avatar-section">
              <el-avatar :size="100">
                <el-icon size="50"><UserFilled /></el-icon>
              </el-avatar>
              <h3>{{ member.name }}</h3>
              <el-tag :type="member.status === 1 ? 'success' : 'info'">
                {{ member.status === 1 ? '正常' : '冻结' }}
              </el-tag>
            </div>
            <el-divider />
            <div class="info-list">
              <div class="info-item">
                <span class="label">手机号：</span>
                <span class="value">{{ member.phone }}</span>
              </div>
              <div class="info-item">
                <span class="label">邮箱：</span>
                <span class="value">{{ member.email || '-' }}</span>
              </div>
              <div class="info-item">
                <span class="label">性别：</span>
                <span class="value">{{ member.gender || '-' }}</span>
              </div>
              <div class="info-item">
                <span class="label">生日：</span>
                <span class="value">{{ member.birthday || '-' }}</span>
              </div>
              <div class="info-item">
                <span class="label">注册时间：</span>
                <span class="value">{{ formatDate(member.created_at) }}</span>
              </div>
            </div>
          </div>
        </el-card>

        <el-card style="margin-top: 20px" v-if="membership">
          <template #header>
            <span style="font-weight: 600">会员卡信息</span>
            <el-button type="primary" size="small" @click="openRenewDialog" style="float: right">
              续费/升级
            </el-button>
          </template>
          <div class="membership-info">
            <div class="membership-type">
              <el-tag :type="getMembershipTagType(membership.type)" size="large">
                {{ getMembershipTypeName(membership.type) }}
              </el-tag>
            </div>
            <div class="info-list" style="margin-top: 16px">
              <div class="info-item">
                <span class="label">开始日期：</span>
                <span class="value">{{ membership.start_date }}</span>
              </div>
              <div class="info-item">
                <span class="label">到期日期：</span>
                <span class="value" :class="{ 'expiring': isExpiringSoon }">
                  {{ membership.end_date }}
                  <el-tag v-if="isExpiringSoon" type="danger" size="small" style="margin-left: 8px">
                    即将到期
                  </el-tag>
                </span>
              </div>
              <div class="info-item">
                <span class="label">剩余天数：</span>
                <span class="value" :class="{ 'text-danger': daysRemaining < 7 }">
                  {{ daysRemaining }} 天
                </span>
              </div>
              <div class="info-item">
                <span class="label">价格：</span>
                <span class="value">¥{{ membership.price }}</span>
              </div>
            </div>
          </div>
        </el-card>

        <el-card style="margin-top: 20px" v-else>
          <template #header>
            <span style="font-weight: 600">会员卡信息</span>
            <el-button type="primary" size="small" @click="openNewMembershipDialog" style="float: right">
              开通会员卡
            </el-button>
          </template>
          <div class="no-membership">
            <el-icon size="48" style="color: #909399"><Warning /></el-icon>
            <p>暂无有效会员卡</p>
          </div>
        </el-card>
      </el-col>

      <el-col :span="16">
        <el-card>
          <el-tabs v-model="activeTab">
            <el-tab-pane label="预约记录" name="bookings">
              <el-table :data="bookings" style="width: 100%" v-loading="loadingBookings">
                <el-table-column prop="schedule.course.name" label="课程名称" />
                <el-table-column prop="schedule.course.coach.name" label="教练" />
                <el-table-column prop="schedule.start_time" label="课程时间">
                  <template #default="{ row }">
                    {{ formatDateTime(row.schedule?.start_time) }}
                  </template>
                </el-table-column>
                <el-table-column prop="booking_time" label="预约时间">
                  <template #default="{ row }">
                    {{ formatDateTime(row.booking_time) }}
                  </template>
                </el-table-column>
                <el-table-column prop="status" label="状态">
                  <template #default="{ row }">
                    <el-tag :type="getBookingStatusType(row.status)">
                      {{ getBookingStatusName(row.status) }}
                    </el-tag>
                  </template>
                </el-table-column>
                <el-table-column label="操作" width="120">
                  <template #default="{ row }">
                    <el-button
                      v-if="row.status === 1"
                      type="danger"
                      link
                      size="small"
                      @click="cancelBooking(row.id)"
                    >
                      取消预约
                    </el-button>
                  </template>
                </el-table-column>
              </el-table>
            </el-tab-pane>
            <el-tab-pane label="签到记录" name="checkins">
              <el-table :data="checkIns" style="width: 100%" v-loading="loadingCheckIns">
                <el-table-column prop="check_in_time" label="签到时间">
                  <template #default="{ row }">
                    {{ formatDateTime(row.check_in_time) }}
                  </template>
                </el-table-column>
                <el-table-column prop="schedule.course.name" label="课程名称">
                  <template #default="{ row }">
                    {{ row.schedule?.course?.name || '自由训练' }}
                  </template>
                </el-table-column>
                <el-table-column prop="check_type" label="签到类型">
                  <template #default="{ row }">
                    {{ row.check_type === 1 ? '正常签到' : '补签' }}
                  </template>
                </el-table-column>
                <el-table-column prop="remark" label="备注" />
              </el-table>
            </el-tab-pane>
          </el-tabs>
        </el-card>
      </el-col>
    </el-row>

    <el-dialog v-model="renewDialogVisible" title="续费/升级" width="400px">
      <el-form label-width="100px">
        <el-form-item label="会员卡类型">
          <el-select v-model="renewForm.type" style="width: 100%">
            <el-option label="月卡" value="monthly" />
            <el-option label="季卡" value="quarter" />
            <el-option label="年卡" value="yearly" />
          </el-select>
        </el-form-item>
        <el-form-item v-if="!membership" label="价格">
          <el-input-number v-model="renewForm.price" :min="0" style="width: 100%" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="renewDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleRenew" :loading="submitting">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { memberApi, membershipApi, bookingApi, checkInApi } from '@/api/member'
import type { Member, Membership, Booking, CheckIn } from '@/types'
import dayjs from 'dayjs'

const route = useRoute()
const memberId = computed(() => Number(route.params.id))

const member = ref<Member | null>(null)
const membership = ref<Membership | null>(null)
const bookings = ref<Booking[]>([])
const checkIns = ref<CheckIn[]>([])
const activeTab = ref('bookings')
const loadingBookings = ref(false)
const loadingCheckIns = ref(false)

const renewDialogVisible = ref(false)
const submitting = ref(false)
const renewForm = ref({
  type: 'monthly',
  price: 299
})

const daysRemaining = computed(() => {
  if (!membership.value) return 0
  return dayjs(membership.value.end_date).diff(dayjs(), 'day')
})

const isExpiringSoon = computed(() => daysRemaining.value > 0 && daysRemaining.value <= 7)

const formatDate = (date: string) => dayjs(date).format('YYYY-MM-DD')
const formatDateTime = (date: string) => date ? dayjs(date).format('YYYY-MM-DD HH:mm') : '-'

const getMembershipTypeName = (type: string) => {
  const map: Record<string, string> = { monthly: '月卡', quarter: '季卡', yearly: '年卡' }
  return map[type] || type
}

const getMembershipTagType = (type: string) => {
  const map: Record<string, string> = { monthly: '', quarter: 'warning', yearly: 'success' }
  return map[type] || ''
}

const getBookingStatusName = (status: number) => {
  const map: Record<number, string> = { 1: '已预约', 2: '已取消', 3: '已签到', 4: '未到场' }
  return map[status] || '未知'
}

const getBookingStatusType = (status: number) => {
  const map: Record<number, string> = { 1: 'primary', 2: 'info', 3: 'success', 4: 'danger' }
  return map[status] || ''
}

const loadMemberDetail = async () => {
  try {
    const res = await memberApi.getById(memberId.value)
    member.value = res.data
    membership.value = res.data.membership || null
  } catch (error) {
    console.error(error)
  }
}

const loadBookings = async () => {
  try {
    loadingBookings.value = true
    const res = await bookingApi.getByMember(memberId.value, { page_size: 100 })
    bookings.value = res.data
  } catch (error) {
    console.error(error)
  } finally {
    loadingBookings.value = false
  }
}

const loadCheckIns = async () => {
  try {
    loadingCheckIns.value = true
    const res = await checkInApi.getByMember(memberId.value, { page_size: 100 })
    checkIns.value = res.data
  } catch (error) {
    console.error(error)
  } finally {
    loadingCheckIns.value = false
  }
}

const openRenewDialog = () => {
  renewForm.value.type = membership.value?.type || 'monthly'
  renewDialogVisible.value = true
}

const openNewMembershipDialog = () => {
  renewForm.value.type = 'monthly'
  renewForm.value.price = 299
  renewDialogVisible.value = true
}

const handleRenew = async () => {
  try {
    submitting.value = true
    if (membership.value) {
      await membershipApi.renew(membership.value.id, renewForm.value.type as any)
    } else {
      await membershipApi.upgrade(memberId.value, {
        new_type: renewForm.value.type as any,
        price: renewForm.value.price
      })
    }
    ElMessage.success('操作成功')
    renewDialogVisible.value = false
    loadMemberDetail()
  } catch (error) {
    console.error(error)
  } finally {
    submitting.value = false
  }
}

const cancelBooking = async (id: number) => {
  try {
    await ElMessageBox.confirm('确定要取消该预约吗？开课前2小时内无法取消', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await bookingApi.cancel(id)
    ElMessage.success('取消成功')
    loadBookings()
  } catch (error) {
    if (error !== 'cancel') {
      console.error(error)
    }
  }
}

onMounted(() => {
  loadMemberDetail()
  loadBookings()
  loadCheckIns()
})
</script>

<style scoped>
.profile-section {
  text-align: center;
}

.avatar-section {
  margin-bottom: 16px;
}

.avatar-section h3 {
  margin: 12px 0 8px 0;
  font-size: 20px;
}

.info-list {
  text-align: left;
}

.info-item {
  display: flex;
  padding: 8px 0;
  border-bottom: 1px solid #f0f0f0;
}

.info-item:last-child {
  border-bottom: none;
}

.info-item .label {
  width: 100px;
  color: #909399;
}

.info-item .value {
  flex: 1;
  color: #303133;
}

.info-item .value.expiring {
  color: #f56c6c;
}

.text-danger {
  color: #f56c6c;
  font-weight: 600;
}

.membership-type {
  text-align: center;
}

.no-membership {
  text-align: center;
  padding: 40px 0;
  color: #909399;
}

.no-membership p {
  margin-top: 12px;
}
</style>
