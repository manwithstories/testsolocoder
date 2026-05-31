<template>
  <div class="dashboard-guide-schedules">
    <div class="card-shadow p-20">
      <div class="flex-between mb-20">
        <h2 class="page-title">导览排班</h2>
        <div class="action-buttons">
          <el-radio-group v-model="viewMode" size="default" class="mr-20">
            <el-radio-button value="calendar">
              <el-icon><Calendar /></el-icon> 日历视图
            </el-radio-button>
            <el-radio-button value="list">
              <el-icon><List /></el-icon> 列表视图
            </el-radio-button>
          </el-radio-group>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon> 新增时段
          </el-button>
        </div>
      </div>

      <div v-if="viewMode === 'calendar'" class="calendar-view">
        <div class="calendar-header flex-between mb-20">
          <div class="flex-center gap-20">
            <el-button :icon="ArrowLeft" circle @click="prevMonth" />
            <span class="calendar-title">{{ currentMonth }}</span>
            <el-button :icon="ArrowRight" circle @click="nextMonth" />
          </div>
          <el-button @click="goToToday">今天</el-button>
        </div>

        <div class="calendar-grid">
          <div class="calendar-weekdays">
            <div v-for="day in weekDays" :key="day" class="weekday">{{ day }}</div>
          </div>
          <div class="calendar-days">
            <div
              v-for="(day, index) in calendarDays"
              :key="index"
              class="calendar-day"
              :class="{
                'other-month': !day.isCurrentMonth,
                'today': day.isToday,
                'has-schedule': day.schedules.length > 0
              }"
              @click="selectDate(day.date)"
            >
              <div class="day-number">{{ day.day }}</div>
              <div v-if="day.schedules.length > 0" class="day-schedules">
                <div
                  v-for="(schedule, sIdx) in day.schedules.slice(0, 2)"
                  :key="sIdx"
                  class="day-schedule"
                  :class="{ 'unavailable': !schedule.is_available }"
                >
                  {{ schedule.start_time }}-{{ schedule.end_time }}
                </div>
                <div v-if="day.schedules.length > 2" class="more-schedules">
                  +{{ day.schedules.length - 2 }}
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div v-else class="list-view">
        <el-form :inline="true" :model="query" class="mb-20">
          <el-form-item label="日期范围">
            <el-date-picker
              v-model="query.start_date"
              type="date"
              placeholder="开始日期"
              value-format="YYYY-MM-DD"
              @change="fetchList"
            />
            <span class="mx-10">至</span>
            <el-date-picker
              v-model="query.end_date"
              type="date"
              placeholder="结束日期"
              value-format="YYYY-MM-DD"
              @change="fetchList"
            />
          </el-form-item>
          <el-form-item label="状态">
            <el-select v-model="query.is_available" clearable style="width: 150px" @change="fetchList">
              <el-option label="可预约" :value="true" />
              <el-option label="不可预约" :value="false" />
            </el-select>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="fetchList">查询</el-button>
            <el-button @click="resetQuery">重置</el-button>
          </el-form-item>
        </el-form>

        <el-table :data="filteredList" v-loading="loading" border>
          <el-table-column prop="id" label="ID" width="60" />
          <el-table-column prop="date" label="日期" width="120">
            <template #default="{ row }">{{ formatDate(row.date) }}</template>
          </el-table-column>
          <el-table-column prop="start_time" label="开始时间" width="100" />
          <el-table-column prop="end_time" label="结束时间" width="100" />
          <el-table-column label="时长" width="80">
            <template #default="{ row }">{{ calculateDuration(row.start_time, row.end_time) }}</template>
          </el-table-column>
          <el-table-column prop="is_available" label="状态" width="100">
            <template #default="{ row }">
              <el-tag :type="row.is_available ? 'success' : 'info'" size="small">
                {{ row.is_available ? '可预约' : '不可预约' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="reservation_id" label="预约状态" width="100">
            <template #default="{ row }">
              <el-tag v-if="row.reservation_id" type="warning" size="small">已预约</el-tag>
              <el-tag v-else type="success" size="small">空闲</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="created_at" label="创建时间" width="160">
            <template #default="{ row }">{{ formatDateTime(row.created_at) }}</template>
          </el-table-column>
          <el-table-column label="操作" width="150" fixed="right">
            <template #default="{ row }">
              <el-button type="primary" size="small" link @click="handleEdit(row)">编辑</el-button>
              <el-button type="danger" size="small" link @click="handleDelete(row)">删除</el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </div>

    <el-dialog v-model="showDialog" :title="isEdit ? '编辑时段' : '新增时段'" width="500px">
      <el-form :model="form" :rules="rules" ref="formRef" label-width="100px">
        <el-form-item label="日期" prop="date">
          <el-date-picker
            v-model="form.date"
            type="date"
            placeholder="选择日期"
            value-format="YYYY-MM-DD"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="开始时间" prop="start_time">
          <el-time-picker
            v-model="form.start_time"
            placeholder="选择开始时间"
            value-format="HH:mm"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="结束时间" prop="end_time">
          <el-time-picker
            v-model="form.end_time"
            placeholder="选择结束时间"
            value-format="HH:mm"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="是否可用" prop="is_available">
          <el-switch v-model="form.is_available" active-text="可预约" inactive-text="不可预约" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showDialog = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="showDayDetail" :title="selectedDateTitle" width="600px">
      <div class="day-detail-header mb-20">
        <h3>{{ selectedDateTitle }}</h3>
      </div>
      <el-table :data="selectedDateSchedules" v-loading="loading" border>
        <el-table-column prop="start_time" label="开始时间" width="100" />
        <el-table-column prop="end_time" label="结束时间" width="100" />
        <el-table-column label="时长" width="80">
          <template #default="{ row }">{{ calculateDuration(row.start_time, row.end_time) }}</template>
        </el-table-column>
        <el-table-column prop="is_available" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.is_available ? 'success' : 'info'" size="small">
              {{ row.is_available ? '可预约' : '不可预约' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="预约状态" width="100">
          <template #default="{ row }">
            <el-tag v-if="row.reservation_id" type="warning" size="small">已预约</el-tag>
            <el-tag v-else type="success" size="small">空闲</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" size="small" link @click="handleEdit(row); showDayDetail = false;">编辑</el-button>
            <el-button type="danger" size="small" link @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      <template #footer>
        <el-button @click="showDayDetail = false">关闭</el-button>
        <el-button type="primary" @click="handleAddForDate">新增时段</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import { Calendar, List, Plus, ArrowLeft, ArrowRight } from '@element-plus/icons-vue'
import * as guideApi from '@/api/guide'
import type { GuideSchedule } from '@/types'
import dayjs from 'dayjs'

const loading = ref(false)
const submitting = ref(false)
const list = ref<GuideSchedule[]>([])
const showDialog = ref(false)
const showDayDetail = ref(false)
const isEdit = ref(false)
const formRef = ref<FormInstance>()
const viewMode = ref<'calendar' | 'list'>('calendar')
const selectedDate = ref<Date>(new Date())
const currentDate = ref<dayjs.Dayjs>(dayjs())

const weekDays = ['日', '一', '二', '三', '四', '五', '六']

const query = reactive({
  start_date: '',
  end_date: '',
  is_available: undefined as boolean | undefined
})

const form = reactive({
  id: 0,
  date: '',
  start_time: '',
  end_time: '',
  is_available: true
})

const rules: FormRules = {
  date: [{ required: true, message: '请选择日期', trigger: 'change' }],
  start_time: [{ required: true, message: '请选择开始时间', trigger: 'change' }],
  end_time: [{ required: true, message: '请选择结束时间', trigger: 'change' }]
}

const currentMonth = computed(() => {
  return currentDate.value.format('YYYY年MM月')
})

const selectedDateTitle = computed(() => {
  return dayjs(selectedDate.value).format('YYYY年MM月DD日')
})

const calendarDays = computed(() => {
  const year = currentDate.value.year()
  const month = currentDate.value.month()
  const firstDay = dayjs(`${year}-${month + 1}-01`)
  const startDay = firstDay.startOf('week')
  const days = []

  for (let i = 0; i < 42; i++) {
    const day = startDay.add(i, 'day')
    const dateStr = day.format('YYYY-MM-DD')
    const schedules = list.value.filter(s => s.date === dateStr)

    days.push({
      date: day.toDate(),
      day: day.date(),
      isCurrentMonth: day.month() === month,
      isToday: day.isSame(dayjs(), 'day'),
      schedules
    })
  }

  return days
})

const filteredList = computed(() => {
  let result = [...list.value]

  if (query.start_date) {
    result = result.filter(s => dayjs(s.date).isAfter(dayjs(query.start_date).subtract(1, 'day')))
  }
  if (query.end_date) {
    result = result.filter(s => dayjs(s.date).isBefore(dayjs(query.end_date).add(1, 'day')))
  }
  if (query.is_available !== undefined) {
    result = result.filter(s => s.is_available === query.is_available)
  }

  return result.sort((a, b) => {
    const dateCompare = dayjs(a.date).valueOf() - dayjs(b.date).valueOf()
    if (dateCompare !== 0) return dateCompare
    return a.start_time.localeCompare(b.start_time)
  })
})

const selectedDateSchedules = computed(() => {
  const dateStr = dayjs(selectedDate.value).format('YYYY-MM-DD')
  return list.value.filter(s => s.date === dateStr).sort((a, b) => a.start_time.localeCompare(b.start_time))
})

const formatDate = (date: string) => {
  return dayjs(date).format('YYYY-MM-DD')
}

const formatDateTime = (date: string) => {
  return dayjs(date).format('YYYY-MM-DD HH:mm')
}

const calculateDuration = (start: string, end: string) => {
  const [startH, startM] = start.split(':').map(Number)
  const [endH, endM] = end.split(':').map(Number)
  const minutes = (endH * 60 + endM) - (startH * 60 + startM)
  const hours = Math.floor(minutes / 60)
  const mins = minutes % 60
  return hours > 0 ? `${hours}小时${mins > 0 ? mins + '分' : ''}` : `${mins}分钟`
}

const prevMonth = () => {
  currentDate.value = currentDate.value.subtract(1, 'month')
}

const nextMonth = () => {
  currentDate.value = currentDate.value.add(1, 'month')
}

const goToToday = () => {
  currentDate.value = dayjs()
}

const selectDate = (date: Date) => {
  selectedDate.value = date
  showDayDetail.value = true
}

const fetchList = async () => {
  try {
    loading.value = true
    const params: { start_date?: string; end_date?: string } = {}
    if (viewMode.value === 'calendar') {
      const year = currentDate.value.year()
      const month = currentDate.value.month()
      const firstDay = dayjs(`${year}-${month + 1}-01`).startOf('week').format('YYYY-MM-DD')
      const lastDay = dayjs(`${year}-${month + 1}-01`).endOf('month').endOf('week').format('YYYY-MM-DD')
      params.start_date = firstDay
      params.end_date = lastDay
    }
    const res = await guideApi.listGuideSchedules(params)
    list.value = res.data
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

const resetQuery = () => {
  query.start_date = ''
  query.end_date = ''
  query.is_available = undefined
  fetchList()
}

const handleAdd = () => {
  isEdit.value = false
  Object.assign(form, {
    id: 0,
    date: dayjs().format('YYYY-MM-DD'),
    start_time: '09:00',
    end_time: '10:00',
    is_available: true
  })
  showDialog.value = true
}

const handleAddForDate = () => {
  isEdit.value = false
  Object.assign(form, {
    id: 0,
    date: dayjs(selectedDate.value).format('YYYY-MM-DD'),
    start_time: '09:00',
    end_time: '10:00',
    is_available: true
  })
  showDialog.value = true
}

const handleEdit = (row: GuideSchedule) => {
  isEdit.value = true
  Object.assign(form, {
    id: row.id,
    date: row.date,
    start_time: row.start_time,
    end_time: row.end_time,
    is_available: row.is_available
  })
  showDialog.value = true
}

const handleDelete = (row: GuideSchedule) => {
  ElMessageBox.confirm('确定要删除该时段吗？', '提示', { type: 'warning' }).then(async () => {
    await guideApi.deleteGuideSchedule(row.id)
    ElMessage.success('删除成功')
    fetchList()
  }).catch(() => {})
}

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate()
  try {
    submitting.value = true
    if (isEdit.value) {
      await guideApi.updateGuideSchedule(form.id, {
        date: form.date,
        start_time: form.start_time,
        end_time: form.end_time,
        is_available: form.is_available
      })
      ElMessage.success('更新成功')
    } else {
      await guideApi.createGuideSchedule({
        date: form.date,
        start_time: form.start_time,
        end_time: form.end_time,
        is_available: form.is_available
      })
      ElMessage.success('创建成功')
    }
    showDialog.value = false
    fetchList()
  } catch (e) {
    console.error(e)
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  fetchList()
})
</script>

<style scoped lang="scss">
.dashboard-guide-schedules {
  .page-title {
    margin: 0;
    font-size: 20px;
  }

  .calendar-view {
    .calendar-header {
      .calendar-title {
        font-size: 18px;
        font-weight: 600;
      }
    }

    .calendar-grid {
      border: 1px solid #e4e7ed;
      border-radius: 4px;
      overflow: hidden;

      .calendar-weekdays {
        display: grid;
        grid-template-columns: repeat(7, 1fr);
        background: #f5f7fa;

        .weekday {
          padding: 12px;
          text-align: center;
          font-weight: 500;
          border-right: 1px solid #e4e7ed;

          &:last-child {
            border-right: none;
          }
        }
      }

      .calendar-days {
        display: grid;
        grid-template-columns: repeat(7, 1fr);

        .calendar-day {
          min-height: 100px;
          padding: 8px;
          border-right: 1px solid #e4e7ed;
          border-bottom: 1px solid #e4e7ed;
          cursor: pointer;
          transition: background-color 0.2s;

          &:nth-child(7n) {
            border-right: none;
          }

          &:hover {
            background: #f5f7fa;
          }

          &.other-month {
            opacity: 0.4;
          }

          &.today {
            background: #ecf5ff;

            .day-number {
              background: #409eff;
              color: #fff;
              border-radius: 50%;
              width: 28px;
              height: 28px;
              display: flex;
              align-items: center;
              justify-content: center;
            }
          }

          &.has-schedule {
            background: #f0f9ff;
          }

          .day-number {
            font-size: 14px;
            margin-bottom: 4px;
          }

          .day-schedules {
            .day-schedule {
              font-size: 12px;
              padding: 2px 4px;
              margin-bottom: 2px;
              background: #67c23a;
              color: #fff;
              border-radius: 2px;
              white-space: nowrap;
              overflow: hidden;
              text-overflow: ellipsis;

              &.unavailable {
                background: #909399;
              }
            }

            .more-schedules {
              font-size: 12px;
              color: #909399;
              padding-left: 4px;
            }
          }
        }
      }
    }
  }

  .day-detail-header {
    h3 {
      margin: 0;
      font-size: 16px;
    }
  }
}
</style>
