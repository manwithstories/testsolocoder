<template>
  <div class="stats-report-page">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>统计报表</span>
          <el-date-picker
            v-model="dateRange"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            value-format="YYYY-MM-DD"
            @change="fetchReports"
          />
        </div>
      </template>

      <el-tabs v-model="activeTab">
        <el-tab-pane label="用户活动" name="activity">
          <div class="toolbar">
            <el-button type="primary" @click="exportActivityReport">
              <el-icon><Download /></el-icon>
              导出报表
            </el-button>
          </div>

          <el-table :data="activityReports" style="width: 100%">
            <el-table-column prop="date" label="日期" width="120" />
            <el-table-column prop="newQuestions" label="新增问题" width="120" />
            <el-table-column prop="newAnswers" label="新增回答" width="120" />
            <el-table-column prop="newUsers" label="新增用户" width="120" />
            <el-table-column prop="newComments" label="新增评论" width="120" />
          </el-table>
        </el-tab-pane>

        <el-tab-pane label="审核统计" name="audit">
          <div class="toolbar">
            <el-button type="primary" @click="exportAuditReport">
              <el-icon><Download /></el-icon>
              导出报表
            </el-button>
          </div>

          <el-table :data="auditReports" style="width: 100%">
            <el-table-column prop="date" label="日期" width="120" />
            <el-table-column prop="reviewedCount" label="审核数量" width="120" />
            <el-table-column prop="approvedCount" label="通过数量" width="120" />
            <el-table-column prop="rejectedCount" label="拒绝数量" width="120" />
          </el-table>
        </el-tab-pane>
      </el-tabs>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { statsApi } from '@/api'
import type { ActivityReport, AuditReport } from '@/types'
import dayjs from 'dayjs'

const activityReports = ref<ActivityReport[]>([])
const auditReports = ref<AuditReport[]>([])
const activeTab = ref('activity')
const dateRange = ref<string[]>([
  dayjs().subtract(30, 'day').format('YYYY-MM-DD'),
  dayjs().format('YYYY-MM-DD')
])

const fetchReports = async () => {
  if (!dateRange.value || dateRange.value.length !== 2) return

  try {
    const [activityRes, auditRes] = await Promise.all([
      statsApi.getActivityReport({
        startDate: dateRange.value[0],
        endDate: dateRange.value[1]
      }),
      statsApi.getAuditReport({
        startDate: dateRange.value[0],
        endDate: dateRange.value[1]
      })
    ])
    activityReports.value = activityRes.data || []
    auditReports.value = auditRes.data || []
  } catch (e) {
    console.error(e)
  }
}

const exportActivityReport = async () => {
  if (!dateRange.value || dateRange.value.length !== 2) return
  try {
    const res = await statsApi.exportActivityReport({
      startDate: dateRange.value[0],
      endDate: dateRange.value[1]
    })
    downloadFile(res.data as Blob, 'activity_report.csv')
  } catch (e) {
    console.error(e)
  }
}

const exportAuditReport = async () => {
  if (!dateRange.value || dateRange.value.length !== 2) return
  try {
    const res = await statsApi.exportAuditReport({
      startDate: dateRange.value[0],
      endDate: dateRange.value[1]
    })
    downloadFile(res.data as Blob, 'audit_report.csv')
  } catch (e) {
    console.error(e)
  }
}

const downloadFile = (blob: Blob, filename: string) => {
  const url = window.URL.createObjectURL(new Blob([blob]))
  const link = document.createElement('a')
  link.href = url
  link.setAttribute('download', filename)
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
}

onMounted(() => {
  fetchReports()
})
</script>

<style scoped lang="scss">
.stats-report-page {
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .toolbar {
    margin-bottom: 16px;
  }
}
</style>
