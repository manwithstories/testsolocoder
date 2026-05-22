<template>
  <Layout>
    <div class="checkin-page">
      <div class="page-header">
        <h2>入场验票</h2>
        <p class="subtitle">扫描票券二维码或输入票号进行验票</p>
      </div>

      <div class="checkin-content">
        <el-card class="checkin-card">
          <el-steps :active="currentStep" finish-status="success" class="checkin-steps">
            <el-step title="输入票号" />
            <el-step title="验证票券" />
            <el-step title="验票完成" />
          </el-steps>

          <div class="checkin-form">
            <el-form :model="form" label-width="100px" @submit.prevent="handleCheckin">
              <el-form-item label="票券编号">
                <el-input
                  v-model="form.ticket_no"
                  placeholder="请扫描或输入票券编号"
                  size="large"
                  @keyup.enter="handleCheckin"
                  ref="ticketInput"
                >
                  <template #append>
                    <el-button type="primary" size="large" @click="handleCheckin" :loading="checking">
                      验票
                    </el-button>
                  </template>
                </el-input>
              </el-form-item>
            </el-form>

            <div class="scan-hint">
              <el-icon><Camera /></el-icon>
              <span>使用扫码枪扫描二维码，或手动输入票号</span>
            </div>
          </div>

          <div class="checkin-result" v-if="checkinResult">
            <el-alert
              :title="checkinResult.title"
              :type="checkinResult.type"
              :description="checkinResult.description"
              show-icon
              :closable="false"
            >
              <template #default v-if="checkinResult.data">
                <div class="result-details">
                  <p><strong>座位信息：</strong>{{ checkinResult.data.seat_info }}</p>
                  <p><strong>观演人：</strong>{{ checkinResult.data.real_name }}</p>
                  <p><strong>验票时间：</strong>{{ checkinResult.data.checkin_time }}</p>
                </div>
              </template>
            </el-alert>
          </div>
        </el-card>

        <el-card class="history-card">
          <template #header>
            <div class="card-header">
              <span>今日验票记录</span>
              <el-button type="primary" size="small" @click="fetchHistory">刷新</el-button>
            </div>
          </template>

          <el-table :data="historyList" style="width: 100%" v-if="historyList.length > 0">
            <el-table-column prop="ticket_no" label="票号" width="200" />
            <el-table-column prop="seat_info" label="座位" />
            <el-table-column prop="real_name" label="观演人" />
            <el-table-column label="状态" width="100">
              <template #default="{ row }">
                <el-tag v-if="row.status === 'success'" type="success">成功</el-tag>
                <el-tag v-else type="danger">失败</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="time" label="时间" width="180" />
          </el-table>

          <el-empty v-else description="暂无验票记录" />
        </el-card>
      </div>
    </div>
  </Layout>
</template>

<script setup lang="ts">
import { ref, onMounted, nextTick } from 'vue'
import { ElMessage } from 'element-plus'
import { Camera } from '@element-plus/icons-vue'
import dayjs from 'dayjs'
import { checkinApi } from '@/api'
import Layout from '@/components/Layout.vue'

interface CheckinForm {
  ticket_no: string
}

interface CheckinResult {
  title: string
  type: 'success' | 'error' | 'warning' | 'info'
  description: string
  data?: any
}

interface HistoryItem {
  ticket_no: string
  seat_info: string
  real_name: string
  status: string
  time: string
}

const form = ref<CheckinForm>({
  ticket_no: ''
})
const currentStep = ref(0)
const checking = ref(false)
const checkinResult = ref<CheckinResult | null>(null)
const historyList = ref<HistoryItem[]>([])
const ticketInput = ref()

async function handleCheckin() {
  if (!form.value.ticket_no.trim()) {
    ElMessage.warning('请输入票券编号')
    return
  }

  checking.value = true
  currentStep.value = 1
  checkinResult.value = null

  try {
    const res = await checkinApi.checkin({
      ticket_no: form.value.ticket_no.trim()
    })

    currentStep.value = 2
    checkinResult.value = {
      title: '验票成功',
      type: 'success',
      description: '欢迎入场，请妥善保管好您的票券',
      data: {
        seat_info: res.ticket?.seat_info || '',
        real_name: res.ticket?.real_name || '',
        checkin_time: dayjs().format('YYYY-MM-DD HH:mm:ss')
      }
    }

    addToHistory(form.value.ticket_no.trim(), res.ticket?.seat_info || '', res.ticket?.real_name || '', 'success')

    setTimeout(() => {
      form.value.ticket_no = ''
      currentStep.value = 0
      nextTick(() => {
        ticketInput.value?.focus()
      })
    }, 3000)

  } catch (err: any) {
    currentStep.value = 0
    checkinResult.value = {
      title: '验票失败',
      type: 'error',
      description: err?.message || '票券验证失败，请检查票号是否正确'
    }

    addToHistory(form.value.ticket_no.trim(), '', '', 'error')
  } finally {
    checking.value = false
  }
}

function addToHistory(ticketNo: string, seatInfo: string, realName: string, status: string) {
  historyList.value.unshift({
    ticket_no: ticketNo,
    seat_info: seatInfo,
    real_name: realName,
    status: status,
    time: dayjs().format('YYYY-MM-DD HH:mm:ss')
  })

  if (historyList.value.length > 50) {
    historyList.value = historyList.value.slice(0, 50)
  }
}

function fetchHistory() {
  ElMessage.success('记录已刷新')
}

onMounted(() => {
  nextTick(() => {
    ticketInput.value?.focus()
  })
})
</script>

<style lang="scss" scoped>
.checkin-page {
  max-width: 1000px;
  margin: 0 auto;
  padding: 20px;

  .page-header {
    text-align: center;
    margin-bottom: 30px;

    h2 {
      margin: 0 0 8px 0;
      font-size: 28px;
    }

    .subtitle {
      color: #999;
      margin: 0;
    }
  }

  .checkin-content {
    display: flex;
    flex-direction: column;
    gap: 24px;
  }

  .checkin-card {
    .checkin-steps {
      max-width: 500px;
      margin: 0 auto 40px;
    }

    .checkin-form {
      max-width: 600px;
      margin: 0 auto;

      .scan-hint {
        display: flex;
        align-items: center;
        justify-content: center;
        gap: 8px;
        margin-top: 20px;
        color: #999;
        font-size: 14px;
      }
    }

    .checkin-result {
      max-width: 600px;
      margin: 30px auto 0;

      .result-details {
        margin-top: 12px;
        padding-top: 12px;
        border-top: 1px solid rgba(0, 0, 0, 0.1);

        p {
          margin: 6px 0;
        }
      }
    }
  }

  .history-card {
    .card-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
    }
  }
}
</style>
