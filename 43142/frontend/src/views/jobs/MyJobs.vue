<template>
  <MainLayout>
    <div class="page-container">
      <div class="page-header">
        <div class="page-title">我的职位</div>
        <div class="page-subtitle">管理您发布的所有职位</div>
      </div>

      <div class="action-bar mb-20">
        <el-button type="primary" @click="goCreate">
          <el-icon><Plus /></el-icon>
          发布职位
        </el-button>
        <el-button @click="handleBulkImport">
          <el-icon><Upload /></el-icon>
          批量导入
        </el-button>
        <el-button type="danger" :disabled="selectedJobs.length === 0" @click="handleBulkDelete">
          <el-icon><Delete /></el-icon>
          批量删除
        </el-button>
      </div>

      <el-card>
        <el-table
          v-loading="loading"
          :data="jobs"
          @selection-change="handleSelectionChange"
          style="width: 100%"
        >
          <el-table-column type="selection" width="55" />
          <el-table-column prop="title" label="职位名称" min-width="200" />
          <el-table-column prop="location" label="工作地点" width="120" />
          <el-table-column label="薪资" width="150">
            <template #default="{ row }">
              {{ formatSalary(row) }}
            </template>
          </el-table-column>
          <el-table-column label="状态" width="100">
            <template #default="{ row }">
              <el-tag :type="getStatusType(row.status)">
                {{ getStatusText(row.status) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="view_count" label="浏览量" width="80" />
          <el-table-column prop="apply_count" label="投递数" width="80" />
          <el-table-column label="创建时间" width="180">
            <template #default="{ row }">
              {{ formatDate(row.created_at) }}
            </template>
          </el-table-column>
          <el-table-column label="操作" width="280" fixed="right">
            <template #default="{ row }">
              <el-button type="primary" link @click="goEdit(row.id)">编辑</el-button>
              <el-button
                v-if="row.status === 'draft'"
                type="success"
                link
                @click="handlePublish(row.id)"
              >
                发布
              </el-button>
              <el-button
                v-if="row.status === 'open'"
                type="warning"
                link
                @click="handleClose(row.id)"
              >
                关闭
              </el-button>
              <el-button type="primary" link @click="viewApplications(row.id)">
                投递简历
              </el-button>
              <el-button type="danger" link @click="handleDelete(row.id)">删除</el-button>
            </template>
          </el-table-column>
        </el-table>

        <div class="pagination-container">
          <el-pagination
            v-model:current-page="pagination.page"
            v-model:page-size="pagination.pageSize"
            :total="total"
            :page-sizes="[10, 20, 50]"
            layout="total, sizes, prev, pager, next"
            @current-change="fetchJobs"
            @size-change="fetchJobs"
          />
        </div>
      </el-card>

      <el-dialog
        v-model="importDialogVisible"
        title="批量导入职位"
        width="900px"
        :close-on-click-modal="false"
      >
        <div class="mb-20">
          <el-button type="primary" @click="addImportRow">
            <el-icon><Plus /></el-icon>
            添加一行
          </el-button>
          <el-button @click="downloadImportTemplate">
            <el-icon><Download /></el-icon>
            下载导入模板
          </el-button>
        </div>
        <el-table :data="importJobs" border style="width: 100%">
          <el-table-column label="#" width="50">
            <template #default="{ $index }">{{ $index + 1 }}</template>
          </el-table-column>
          <el-table-column label="职位名称" min-width="150">
            <template #default="{ row }">
              <el-input v-model="row.title" placeholder="必填" />
            </template>
          </el-table-column>
          <el-table-column label="工作地点" min-width="120">
            <template #default="{ row }">
              <el-input v-model="row.location" placeholder="必填" />
            </template>
          </el-table-column>
          <el-table-column label="最低薪资(K)" width="110">
            <template #default="{ row }">
              <el-input-number v-model="row.salary_min_k" :min="0" :step="1" controls-position="right" />
            </template>
          </el-table-column>
          <el-table-column label="最高薪资(K)" width="110">
            <template #default="{ row }">
              <el-input-number v-model="row.salary_max_k" :min="0" :step="1" controls-position="right" />
            </template>
          </el-table-column>
          <el-table-column label="工作类型" width="110">
            <template #default="{ row }">
              <el-select v-model="row.job_type" placeholder="选择">
                <el-option label="全职" value="full-time" />
                <el-option label="兼职" value="part-time" />
                <el-option label="实习" value="intern" />
              </el-select>
            </template>
          </el-table-column>
          <el-table-column label="经验要求" width="110">
            <template #default="{ row }">
              <el-select v-model="row.experience" placeholder="选择">
                <el-option label="不限" value="" />
                <el-option label="应届" value="fresh" />
                <el-option label="1-3年" value="1-3" />
                <el-option label="3-5年" value="3-5" />
                <el-option label="5年以上" value="5+" />
              </el-select>
            </template>
          </el-table-column>
          <el-table-column label="技能要求" min-width="120">
            <template #default="{ row }">
              <el-input v-model="row.skills" placeholder="如：Java,Go" />
            </template>
          </el-table-column>
          <el-table-column label="职位描述" min-width="150">
            <template #default="{ row }">
              <el-input v-model="row.description" type="textarea" :rows="1" />
            </template>
          </el-table-column>
          <el-table-column label="操作" width="70" fixed="right">
            <template #default="{ $index }">
              <el-button type="danger" link @click="removeImportRow($index)">删除</el-button>
            </template>
          </el-table-column>
        </el-table>
        <template #footer>
          <el-button @click="importDialogVisible = false">取消</el-button>
          <el-button type="primary" :loading="importing" @click="doBulkImport">
            确认导入
          </el-button>
        </template>
      </el-dialog>
    </div>
  </MainLayout>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import MainLayout from '@/layouts/MainLayout.vue'
import { listMyJobs, publishJob, closeJob, deleteJob, bulkDeleteJobs, bulkImportJobs } from '@/api/job'
import type { Job } from '@/types'
import dayjs from 'dayjs'

interface ImportJobItem {
  title: string
  location: string
  salary_min_k: number | undefined
  salary_max_k: number | undefined
  job_type: string
  experience: string
  skills: string
  description: string
}

function createEmptyImportJob(): ImportJobItem {
  return {
    title: '',
    location: '',
    salary_min_k: undefined,
    salary_max_k: undefined,
    job_type: 'full-time',
    experience: '',
    skills: '',
    description: ''
  }
}

const router = useRouter()
const loading = ref(false)
const jobs = ref<Job[]>([])
const total = ref(0)
const selectedJobs = ref<Job[]>([])

const pagination = reactive({
  page: 1,
  pageSize: 10
})

const importDialogVisible = ref(false)
const importing = ref(false)
const importJobs = ref<ImportJobItem[]>([createEmptyImportJob()])

async function fetchJobs() {
  loading.value = true
  try {
    const res = await listMyJobs({
      page: pagination.page,
      page_size: pagination.pageSize
    })
    if (res.data) {
      jobs.value = res.data.items
      total.value = res.data.pagination.total
    }
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

function handleSelectionChange(selection: Job[]) {
  selectedJobs.value = selection
}

function goCreate() {
  router.push('/my/jobs/create')
}

function goEdit(id: number) {
  router.push(`/my/jobs/${id}/edit`)
}

function viewApplications(id: number) {
  router.push(`/company/applications?jobId=${id}`)
}

async function handlePublish(id: number) {
  try {
    await publishJob(id)
    ElMessage.success('发布成功')
    fetchJobs()
  } catch (e) {
    // error handled
  }
}

async function handleClose(id: number) {
  try {
    await closeJob(id)
    ElMessage.success('已关闭')
    fetchJobs()
  } catch (e) {
    // error handled
  }
}

async function handleDelete(id: number) {
  ElMessageBox.confirm('确定要删除该职位吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      await deleteJob(id)
      ElMessage.success('删除成功')
      fetchJobs()
    } catch (e) {
      // error handled
    }
  }).catch(() => {})
}

async function handleBulkDelete() {
  if (selectedJobs.value.length === 0) return

  ElMessageBox.confirm(`确定要删除选中的 ${selectedJobs.value.length} 个职位吗？`, '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      await bulkDeleteJobs(selectedJobs.value.map(j => j.id))
      ElMessage.success('批量删除成功')
      fetchJobs()
    } catch (e) {
      // error handled
    }
  }).catch(() => {})
}

function handleBulkImport() {
  importDialogVisible.value = true
}

function addImportRow() {
  importJobs.value.push(createEmptyImportJob())
}

function removeImportRow(index: number) {
  if (importJobs.value.length > 1) {
    importJobs.value.splice(index, 1)
  }
}

async function doBulkImport() {
  const validJobs = importJobs.value.filter(item => item.title.trim() && item.location.trim())
  if (validJobs.length === 0) {
    ElMessage.warning('请至少填写一条有效的职位数据')
    return
  }

  const payload = validJobs.map(item => ({
    title: item.title,
    description: item.description,
    location: item.location,
    job_type: item.job_type,
    experience: item.experience,
    skills: item.skills,
    salary_min: item.salary_min_k ? item.salary_min_k * 1000 : 0,
    salary_max: item.salary_max_k ? item.salary_max_k * 1000 : 0,
    requirements: ''
  }))

  importing.value = true
  try {
    const res = await bulkImportJobs(payload)
    const data = res.data as { success_count: number; errors: string[] }
    if (data && data.errors && data.errors.length > 0) {
      ElMessage.warning(`导入 ${data.success_count} 条成功，${data.errors.length} 条失败：${data.errors.join('; ')}`)
    } else {
      ElMessage.success(`成功导入 ${data?.success_count ?? validJobs.length} 条职位`)
    }
    importDialogVisible.value = false
    importJobs.value = [createEmptyImportJob()]
    fetchJobs()
  } catch (e) {
    // error handled by interceptor
  } finally {
    importing.value = false
  }
}

function downloadImportTemplate() {
  const headers = ['职位名称', '工作地点', '最低薪资(K)', '最高薪资(K)', '工作类型', '经验要求', '技能要求', '职位描述']
  const sampleRow = ['高级Go工程师', '北京', '20', '40', 'full-time', '3-5', 'Go,微服务', '负责后端服务开发']
  const csvContent = '\ufeff' + [headers.join(','), sampleRow.join(',')].join('\n')
  const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' })
  const url = URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = url
  link.download = '批量导入模板.csv'
  link.click()
  URL.revokeObjectURL(url)
}

function formatSalary(job: Job) {
  if (job.salary_min && job.salary_max) {
    return `${(job.salary_min / 1000).toFixed(0)}K-${(job.salary_max / 1000).toFixed(0)}K`
  }
  return '面议'
}

function getStatusType(status: string) {
  const types: Record<string, string> = {
    'open': 'success',
    'closed': 'info',
    'paused': 'warning',
    'draft': 'info'
  }
  return types[status] || ''
}

function getStatusText(status: string) {
  const texts: Record<string, string> = {
    'open': '招聘中',
    'closed': '已关闭',
    'paused': '已暂停',
    'draft': '草稿'
  }
  return texts[status] || status
}

function formatDate(date: string) {
  return dayjs(date).format('YYYY-MM-DD HH:mm')
}

onMounted(() => {
  fetchJobs()
})
</script>

<style scoped>
.action-bar {
  display: flex;
  gap: 10px;
}

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: center;
}
</style>
