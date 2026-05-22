<template>
  <div class="checkin">
    <div class="page-header">
      <h2 class="page-title">签到管理</h2>
    </div>

    <el-row :gutter="20">
      <el-col :span="8">
        <el-card>
          <template #header>扫码签到</template>
          <div class="checkin-box">
            <el-input
              v-model="qrCode"
              placeholder="输入或扫描签到码"
              size="large"
              @keyup.enter="handleCheckIn"
            >
              <template #append>
                <el-button type="primary" :loading="checking" @click="handleCheckIn">签到</el-button>
              </template>
            </el-input>
            <div class="scan-tip">
              <el-icon size="48" color="#909399"><Camera /></el-icon>
              <p>请扫描门票二维码进行签到</p>
            </div>
          </div>
        </el-card>

        <el-card style="margin-top: 20px">
          <template #header>签到统计</template>
          <el-row :gutter="20">
            <el-col :span="12">
              <div class="stat-box">
                <p class="stat-num">{{ checkInStats.total }}</p>
                <p class="stat-label">总票数</p>
              </div>
            </el-col>
            <el-col :span="12">
              <div class="stat-box checked">
                <p class="stat-num">{{ checkInStats.checked }}</p>
                <p class="stat-label">已签到</p>
              </div>
            </el-col>
          </el-row>
        </el-card>
      </el-col>

      <el-col :span="16">
        <el-card>
          <template #header>签到记录</template>
          <div class="search-bar">
            <el-select v-model="search.activityId" placeholder="选择活动" clearable style="width: 200px">
              <el-option v-for="act in activities" :key="act.id" :label="act.title" :value="act.id" />
            </el-select>
            <el-select v-model="search.checkedIn" placeholder="签到状态" clearable style="width: 140px">
              <el-option label="已签到" :value="true" />
              <el-option label="未签到" :value="false" />
            </el-select>
            <el-input v-model="search.keyword" placeholder="订单号" clearable style="width: 200px" />
            <el-button type="primary" @click="loadData">搜索</el-button>
          </div>

          <el-table :data="list" style="width: 100%" v-loading="loading">
            <el-table-column prop="qrCode" label="签到码" width="180" />
            <el-table-column prop="activity.title" label="活动名称" min-width="180" />
            <el-table-column prop="order.orderNo" label="订单号" width="180" />
            <el-table-column label="票型" width="120">
              <template #default="{ row }">
                {{ getTicketTypeName(row.ticketTypeID) }}
              </template>
            </el-table-column>
            <el-table-column label="状态" width="100">
              <template #default="{ row }">
                <el-tag :type="row.checkedIn ? 'success' : 'info'">
                  {{ row.checkedIn ? '已签到' : '未签到' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="checkedAt" label="签到时间" width="180" :formatter="formatDate" />
          </el-table>

          <el-pagination
            v-model:current-page="pagination.page"
            v-model:page-size="pagination.pageSize"
            :total="pagination.total"
            layout="total, prev, pager, next, jumper"
            @current-change="loadData"
            style="margin-top: 20px; justify-content: flex-end"
          />
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { checkIn, getCheckInList, getCheckInStatistics } from '@/api/checkin'
import { getActivityList } from '@/api/activity'
import dayjs from 'dayjs'

const checking = ref(false)
const loading = ref(false)
const qrCode = ref('')
const list = ref<any[]>([])
const activities = ref<any[]>([])

const search = reactive({
  activityId: undefined as number | undefined,
  checkedIn: undefined as boolean | undefined,
  keyword: ''
})

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

const checkInStats = reactive({
  total: 0,
  checked: 0
})

const loadActivities = async () => {
  try {
    const res = await getActivityList({ page: 1, pageSize: 100 })
    activities.value = res.list
  } catch (error) {
    console.error(error)
  }
}

const loadData = async () => {
  try {
    loading.value = true
    const res = await getCheckInList({
      page: pagination.page,
      pageSize: pagination.pageSize,
      activityId: search.activityId,
      checkedIn: search.checkedIn,
      keyword: search.keyword
    })
    list.value = res.list
    pagination.total = res.total
    loadStats()
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}

const loadStats = async () => {
  try {
    const res = await getCheckInStatistics(search.activityId)
    checkInStats.total = res.total
    checkInStats.checked = res.checked
  } catch (error) {
    console.error(error)
  }
}

const handleCheckIn = async () => {
  if (!qrCode.value) {
    ElMessage.warning('请输入签到码')
    return
  }
  try {
    checking.value = true
    await checkIn(qrCode.value)
    ElMessage.success('签到成功')
    qrCode.value = ''
    loadData()
  } catch (error) {
    console.error(error)
  } finally {
    checking.value = false
  }
}

const getTicketTypeName = (id: number) => {
  for (const act of activities.value) {
    const ticket = act.ticketTypes?.find((t: any) => t.id === id)
    if (ticket) return ticket.name
  }
  return '-'
}

const formatDate = (_row: any, _column: any, value: string) => value ? dayjs(value).format('YYYY-MM-DD HH:mm:ss') : '-'

onMounted(() => {
  loadActivities()
  loadData()
})
</script>

<style scoped lang="scss">
.checkin-box {
  text-align: center;
}

.scan-tip {
  margin-top: 40px;
  color: #909399;

  p {
    margin-top: 10px;
  }
}

.stat-box {
  text-align: center;
  padding: 20px;
  background: #f5f7fa;
  border-radius: 8px;

  &.checked {
    background: #f0f9eb;
  }

  .stat-num {
    font-size: 32px;
    font-weight: 600;
    margin: 0;
    color: #303133;
  }

  .stat-label {
    margin: 8px 0 0;
    color: #909399;
  }
}
</style>
