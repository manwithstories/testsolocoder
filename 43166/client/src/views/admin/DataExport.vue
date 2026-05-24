<template>
  <div class="data-export">
    <div class="page-header">
      <h2 class="page-title">数据导出</h2>
    </div>

    <el-row :gutter="20">
      <el-col :span="12">
        <div class="card">
          <h3 class="font-weight-600 mb-16">新建导出任务</h3>
          <el-form :model="exportForm" label-width="100px">
            <el-form-item label="导出类型">
              <el-select v-model="exportForm.type" style="width: 100%">
                <el-option label="申请列表" value="applications" />
                <el-option label="费用明细" value="fees" />
                <el-option label="专员工作记录" value="agents" />
              </el-select>
            </el-form-item>
            <el-form-item label="时间范围">
              <el-date-picker
                v-model="exportForm.dateRange"
                type="daterange"
                range-separator="至"
                start-placeholder="开始日期"
                end-placeholder="结束日期"
                value-format="YYYY-MM-DD"
              />
            </el-form-item>
            <el-form-item label="状态筛选">
              <el-select v-model="exportForm.conditions.status" placeholder="全部" clearable style="width: 100%">
                <el-option label="全部" value="" />
                <el-option label="草稿" value="draft" />
                <el-option label="待审核" value="pending_review" />
                <el-option label="处理中" value="processing" />
                <el-option label="已完成" value="completed" />
                <el-option label="已驳回" value="rejected" />
                <el-option label="待支付" value="pending" />
                <el-option label="已支付" value="paid" />
              </el-select>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" :loading="exporting" @click="handleExport">
                开始导出
              </el-button>
            </el-form-item>
          </el-form>
        </div>
      </el-col>
      <el-col :span="12">
        <div class="card">
          <h3 class="font-weight-600 mb-16">导出历史</h3>
          <el-table :data="exportTasks" style="width: 100%" v-loading="loading">
            <el-table-column prop="type" label="类型" width="120">
              <template #default="{ row }">
                {{ getTypeText(row.type) }}
              </template>
            </el-table-column>
            <el-table-column prop="fileName" label="文件名" />
            <el-table-column prop="status" label="状态" width="100">
              <template #default="{ row }">
                <el-tag :type="getStatusType(row.status)">
                  {{ getStatusText(row.status) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="createdAt" label="创建时间" width="160">
              <template #default="{ row }">
                {{ formatTime(row.createdAt) }}
              </template>
            </el-table-column>
            <el-table-column label="操作" width="100">
              <template #default="{ row }">
                <el-button
                  v-if="row.status === 'completed'"
                  type="primary"
                  link
                  @click="downloadExport(row)"
                >
                  下载
                </el-button>
              </template>
            </el-table-column>
          </el-table>
        </div>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { exportApi } from '@/api/statistics'
import { ExportTask } from '@/types'
import dayjs from 'dayjs'

const loading = ref(false)
const exporting = ref(false)
const exportTasks = ref<ExportTask[]>([])

const exportForm = reactive({
  type: 'applications',
  dateRange: [] as string[],
  conditions: {
    status: ''
  }
})

const fetchTasks = async () => {
  loading.value = true
  try {
    const res = await exportApi.getTasks({ page: 1, pageSize: 10 })
    exportTasks.value = res?.list || []
  } catch (error) {
    console.error('获取导出任务失败:', error)
  } finally {
    loading.value = false
  }
}

const handleExport = async () => {
  exporting.value = true
  try {
    const params: any = {
      type: exportForm.type
    }
    if (exportForm.dateRange.length === 2) {
      params.startDate = exportForm.dateRange[0]
      params.endDate = exportForm.dateRange[1]
    }
    if (exportForm.conditions.status) {
      params.conditions = exportForm.conditions
    }

    await exportApi.createTask(params)
    ElMessage.success('导出任务已创建，处理完成后可下载')
    fetchTasks()
  } catch (error: any) {
    ElMessage.error(error.message || '导出失败')
  } finally {
    exporting.value = false
  }
}

const downloadExport = async (task: ExportTask) => {
  try {
    const res = await exportApi.download(task.id)
    const blob = new Blob([res as any])
    const url = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = task.fileName
    link.click()
    window.URL.revokeObjectURL(url)
  } catch (error) {
    ElMessage.error('下载失败')
  }
}

const getTypeText = (type: string) => {
  const map: Record<string, string> = {
    applications: '申请列表',
    fees: '费用明细',
    agents: '专员记录'
  }
  return map[type] || type
}

const getStatusType = (status: string) => {
  const map: Record<string, string> = {
    pending: 'warning',
    completed: 'success',
    failed: 'danger'
  }
  return map[status] || 'info'
}

const getStatusText = (status: string) => {
  const map: Record<string, string> = {
    pending: '处理中',
    completed: '已完成',
    failed: '失败'
  }
  return map[status] || status
}

const formatTime = (time: string) => {
  return dayjs(time).format('YYYY-MM-DD HH:mm')
}

onMounted(fetchTasks)
</script>
