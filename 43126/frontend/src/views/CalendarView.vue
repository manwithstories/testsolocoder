<template>
  <div class="calendar-view">
    <el-card class="filter-card">
      <el-form :inline="true" :model="filterForm">
        <el-form-item label="视图">
          <el-radio-group v-model="viewType">
            <el-radio-button value="week">周视图</el-radio-button>
            <el-radio-button value="month">月视图</el-radio-button>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="日期">
          <el-date-picker
            v-model="currentDate"
            :type="viewType === 'week' ? 'week' : 'month'"
            placeholder="选择日期"
            @change="loadCalendar"
          />
        </el-form-item>
        <el-form-item label="楼层">
          <el-select v-model="filterForm.floor" placeholder="全部楼层" clearable @change="loadCalendar">
            <el-option v-for="floor in floors" :key="floor" :label="floor" :value="floor" />
          </el-select>
        </el-form-item>
        <el-form-item label="会议室">
          <el-select v-model="filterForm.room_id" placeholder="全部会议室" clearable @change="loadCalendar">
            <el-option v-for="room in rooms" :key="room.id" :label="room.name" :value="room.id" />
          </el-select>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card style="margin-top: 20px; border-radius: 8px">
      <div class="calendar-header">
        <el-button :icon="ArrowLeft" circle @click="prevPeriod" />
        <span class="period-label">{{ periodLabel }}</span>
        <el-button :icon="ArrowRight" circle @click="nextPeriod" />
        <el-button type="primary" @click="goToday">今天</el-button>
      </div>

      <div v-if="viewType === 'week'" class="week-view">
        <div class="time-column">
          <div class="time-header">时间</div>
          <div v-for="hour in timeSlots" :key="hour" class="time-slot">{{ hour }}</div>
        </div>
        <div v-for="day in weekDays" :key="day.date" class="day-column">
          <div class="day-header" :class="{ today: isToday(day.date) }">
            <div class="day-name">{{ day.name }}</div>
            <div class="day-date">{{ day.date.slice(5) }}</div>
          </div>
          <div class="day-slots">
            <div
              v-for="hour in timeSlots"
              :key="hour"
              class="slot"
              :class="{ 'slot-booked': isSlotBooked(day.date, hour) }"
              @click="onSlotClick(day.date, hour)"
            >
              <div v-if="getBookingForSlot(day.date, hour)" class="booking-item" :style="{ background: getBookingForSlot(day.date, hour)?.color }">
                {{ getBookingForSlot(day.date, hour)?.title }}
              </div>
            </div>
          </div>
        </div>
      </div>

      <div v-else class="month-view">
        <div class="month-grid">
          <div v-for="dayName in ['周一', '周二', '周三', '周四', '周五', '周六', '周日']" :key="dayName" class="month-header">
            {{ dayName }}
          </div>
          <div
            v-for="(day, index) in monthDays"
            :key="index"
            class="month-day"
            :class="{ other: day.otherMonth, today: isToday(day.date) }"
          >
            <div class="day-number">{{ day.date.slice(8) }}</div>
            <div class="day-bookings">
              <div
                v-for="booking in getBookingsForDate(day.date)"
                :key="booking.id"
                class="mini-booking"
                :style="{ background: booking.color }"
              >
                {{ booking.title }}
              </div>
            </div>
          </div>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, watch } from 'vue'
import { ArrowLeft, ArrowRight } from '@element-plus/icons-vue'
import { api } from '@/api'
import dayjs from 'dayjs'
import isoWeek from 'dayjs/plugin/isoWeek'

dayjs.extend(isoWeek)

const viewType = ref<'week' | 'month'>('week')
const currentDate = ref(dayjs().toDate())
const calendarData = ref<any[]>([])
const rooms = ref<any[]>([])
const floors = ref<string[]>([])

const filterForm = reactive({
  floor: '',
  room_id: undefined as number | undefined
})

const timeSlots = computed(() => {
  const slots: string[] = []
  if (filterForm.room_id) {
    const room = rooms.value.find(r => r.id === filterForm.room_id)
    if (room) {
      const startHour = parseInt(room.available_start?.split(':')[0] || '8')
      const endHour = parseInt(room.available_end?.split(':')[0] || '22')
      for (let h = startHour; h <= endHour; h++) {
        slots.push(`${h.toString().padStart(2, '0')}:00`)
      }
      return slots
    }
  }
  for (let h = 8; h <= 22; h++) {
    slots.push(`${h.toString().padStart(2, '0')}:00`)
  }
  return slots
})

const weekDays = computed(() => {
  const start = dayjs(currentDate.value).startOf('isoWeek')
  const days = []
  for (let i = 0; i < 7; i++) {
    const d = start.add(i, 'day')
    days.push({
      date: d.format('YYYY-MM-DD'),
      name: ['周一', '周二', '周三', '周四', '周五', '周六', '周日'][i]
    })
  }
  return days
})

const monthDays = computed(() => {
  const start = dayjs(currentDate.value).startOf('month')
  const end = dayjs(currentDate.value).endOf('month')
  const days: any[] = []

  const startWeekday = (start.isoWeekday() % 7)
  for (let i = startWeekday - 1; i >= 0; i--) {
    const d = start.subtract(i + 1, 'day')
    days.push({ date: d.format('YYYY-MM-DD'), otherMonth: true })
  }

  for (let d = start; d.isBefore(end) || d.isSame(end, 'day'); d = d.add(1, 'day')) {
    days.push({ date: d.format('YYYY-MM-DD'), otherMonth: false })
  }

  const remaining = 42 - days.length
  for (let i = 1; i <= remaining; i++) {
    const d = end.add(i, 'day')
    days.push({ date: d.format('YYYY-MM-DD'), otherMonth: true })
  }

  return days
})

const periodLabel = computed(() => {
  if (viewType.value === 'week') {
    const start = dayjs(currentDate.value).startOf('isoWeek')
    const end = start.add(6, 'day')
    return `${start.format('YYYY年MM月DD日')} - ${end.format('MM月DD日')}`
  }
  return dayjs(currentDate.value).format('YYYY年MM月')
})

onMounted(() => {
  loadRooms()
  loadFloors()
  loadCalendar()
})

watch(viewType, () => {
  loadCalendar()
})

async function loadRooms() {
  try {
    const res: any = await api.listAllRooms()
    rooms.value = res.data || []
  } catch (e) {
    console.error(e)
  }
}

async function loadFloors() {
  try {
    const res: any = await api.getFloors()
    floors.value = res.data || []
  } catch (e) {
    console.error(e)
  }
}

async function loadCalendar() {
  try {
    let res: any
    if (viewType.value === 'week') {
      res = await api.getWeekCalendar({
        date: dayjs(currentDate.value).format('YYYY-MM-DD'),
        room_id: filterForm.room_id,
        floor: filterForm.floor
      })
    } else {
      res = await api.getMonthCalendar({
        date: dayjs(currentDate.value).format('YYYY-MM'),
        room_id: filterForm.room_id,
        floor: filterForm.floor
      })
    }
    calendarData.value = res.data || []
  } catch (e) {
    console.error(e)
  }
}

function prevPeriod() {
  if (viewType.value === 'week') {
    currentDate.value = dayjs(currentDate.value).subtract(1, 'week').toDate()
  } else {
    currentDate.value = dayjs(currentDate.value).subtract(1, 'month').toDate()
  }
  loadCalendar()
}

function nextPeriod() {
  if (viewType.value === 'week') {
    currentDate.value = dayjs(currentDate.value).add(1, 'week').toDate()
  } else {
    currentDate.value = dayjs(currentDate.value).add(1, 'month').toDate()
  }
  loadCalendar()
}

function goToday() {
  currentDate.value = dayjs().toDate()
  loadCalendar()
}

function isToday(date: string) {
  return dayjs(date).isSame(dayjs(), 'day')
}

function getBookingsForDate(date: string) {
  return calendarData.value.filter(b => dayjs(b.start).isSame(date, 'day'))
}

function isSlotBooked(date: string, hour: string) {
  const h = parseInt(hour)
  return calendarData.value.some(b => {
    const startH = dayjs(b.start).hour()
    const endH = dayjs(b.end).hour()
    return dayjs(b.start).isSame(date, 'day') && h >= startH && h < endH
  })
}

function getBookingForSlot(date: string, hour: string) {
  const h = parseInt(hour)
  return calendarData.value.find(b => {
    const startH = dayjs(b.start).hour()
    return dayjs(b.start).isSame(date, 'day') && h === startH
  })
}

function onSlotClick(date: string, hour: string) {}
</script>

<style scoped>
.filter-card {
  border-radius: 8px;
}

.calendar-header {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 20px;
}

.period-label {
  font-size: 16px;
  font-weight: 600;
  color: #303133;
}

.week-view {
  display: flex;
  border: 1px solid #ebeef5;
  border-radius: 4px;
  overflow: hidden;
}

.time-column {
  width: 60px;
  flex-shrink: 0;
  background: #f5f7fa;
}

.time-header, .day-header {
  height: 40px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  background: #f5f7fa;
  border-bottom: 1px solid #ebeef5;
  font-size: 12px;
  color: #606266;
}

.day-header.today {
  background: #409EFF;
  color: #fff;
}

.day-name {
  font-weight: 500;
}

.day-date {
  font-size: 11px;
}

.time-slot {
  height: 50px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  color: #909399;
  border-bottom: 1px solid #ebeef5;
}

.day-column {
  flex: 1;
  border-left: 1px solid #ebeef5;
}

.day-slots {
  position: relative;
}

.slot {
  height: 50px;
  border-bottom: 1px solid #ebeef5;
  cursor: pointer;
  position: relative;
}

.slot:hover {
  background: #f0f9ff;
}

.slot-booked {
  background: #ecf5ff;
}

.booking-item {
  position: absolute;
  left: 2px;
  right: 2px;
  top: 2px;
  bottom: 2px;
  border-radius: 4px;
  padding: 2px 6px;
  font-size: 11px;
  color: #fff;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.month-view {
  border: 1px solid #ebeef5;
  border-radius: 4px;
  overflow: hidden;
}

.month-grid {
  display: grid;
  grid-template-columns: repeat(7, 1fr);
}

.month-header {
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f5f7fa;
  border-bottom: 1px solid #ebeef5;
  font-size: 13px;
  color: #606266;
  font-weight: 500;
}

.month-day {
  min-height: 100px;
  border-right: 1px solid #ebeef5;
  border-bottom: 1px solid #ebeef5;
  padding: 4px;
  cursor: pointer;
}

.month-day:hover {
  background: #f0f9ff;
}

.month-day.other {
  background: #fafafa;
  color: #c0c4cc;
}

.month-day.today .day-number {
  background: #409EFF;
  color: #fff;
  border-radius: 50%;
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.day-number {
  font-size: 13px;
  margin-bottom: 4px;
}

.day-bookings {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.mini-booking {
  padding: 2px 4px;
  border-radius: 2px;
  font-size: 10px;
  color: #fff;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
</style>
