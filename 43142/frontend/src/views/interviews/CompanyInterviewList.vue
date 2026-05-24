<template>
  <MainLayout>
    <div class="page-container">
      <div class="page-header">
        <div class="page-title">面试管理</div>
        <div class="page-subtitle">管理所有面试安排</div>
      </div>

      <div class="action-bar mb-20">
        <el-button type="primary" @click="goSchedule">
          <el-icon><Plus /></el-icon>
          安排面试
        </el-button>
      </div>

      <el-card>
        <el-table v-loading="loading" :data="interviews" style="width: 100%">
          <el-table-column label="求职者" min-width="150">
            <template #default="{ row }">
              <div class="applicant-info">
                <el-avatar :size="40">{{ getApplicantName(row).charAt(0) }}</el-avatar>
                <div>
                  <div class="name">{{ getApplicantName(row) }}</div>
                </div>
              </div>
            </template>
          </el-table-column>
          <el-table-column label="应聘职位" min-width="150">
            <template #default="{ row }">
              {{ row.application?.job?.title }}
            </template>
          </el-table-column>
          <el-table-column label="面试时间" width="180">
            <template #default="{ row }">
              {{ formatDateTime(row.scheduled_at) }}
            </template>
          </el-table-column>
          <el-table-column label="时长" width="100">
            <template #default="{ row }">
              {{ row.duration }}分钟
            </template>
          </el-table-column>
          <el-table-column label="面试官" width="120">
            <template #default="{ row }">
              {{ row.interviewer || '-' }}
            </template>
          </el-table-column>
          <el-table-column label="状态" width="120">
            <template #default="{ row }">
              <el-tag :type="getStatusType(row.status)">
                {{ getStatusText(row.status) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="操作" width="200" fixed="right">
            <template #default="{ row }">
              <el-dropdown trigger="click" @command="(cmd: string) => handleAction(row.id, cmd)">
                <el-button type="primary" link>
                  操作
                  <el-icon><ArrowDown /></el-icon>
                </el-button>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item command="complete">完成面试</el-dropdown-item>
                    <el-dropdown-item command="cancel">取消面试</el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
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
            @current-change="fetchInterviews"
            @size-change="fetchInterviews"
          />
        </div>
      </el-card>

      <el-dialog v-model="completeDialogVisible" title="完成面试" width="500px">
        <el-form :model="completeForm" label-width="80px">
          <el-form-item label="评分">
            <el-rate v-model="completeForm.rating" />
          </el-form-item>
          <el-form-item label="面试反馈">
            <el-input
              v-model="completeForm.feedback"
              type="textarea"
              :rows="4"
              placeholder="请输入面试反馈"
            />
          </el-form-item>
        </el-form>
        <template #footer>
          <el-button @click="completeDialogVisible = false">取消</el-button>
          <el-button type="primary" :loading="completing" @click="confirmComplete">确认</el-button>
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
import { listCompanyInterviews, completeInterview, cancelInterview } from '@/api/interview'
import type { Interview } from '@/types'
import dayjs from 'dayjs'

const router = useRouter()
const loading = ref(false)
const interviews = ref<Interview[]>([])
const total = ref(0)
const currentInterviewId = ref<number | null>(null)
const completeDialogVisible = ref(false)
const completing = ref(false)

const completeForm = reactive({
  rating: 0,
  feedback: ''
})

const pagination = reactive({
  page: 1,
  pageSize: 10
})

async function fetchInterviews() {
  loading.value = true
  try {
    const res = await listCompanyInterviews({
      page: pagination.page,
      page_size: pagination.pageSize
    })
    if (res.data) {
      interviews.value = res.data.items
      total.value = res.data.pagination.total
    }
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

function goSchedule() {
  router.push('/company/interviews/schedule')
}

function handleAction(id: number, action: string) {
  if (action === 'complete') {
    currentInterviewId.value = id
    completeForm.rating = 0
    completeForm.feedback = ''
    completeDialogVisible.value = true
  } else if (action === 'cancel') {
    ElMessageBox.confirm('确定要取消该面试吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }).then(async () => {
      try {
        await cancelInterview(id)
        ElMessage.success('已取消面试')
        fetchInterviews()
      } catch (e) {
        // error handled
      }
    }).catch(() => {})
  }
}

async function confirmComplete() {
  if (currentInterviewId.value === null) return

  completing.value = true
  try {
    await completeInterview(currentInterviewId.value, completeForm.feedback, completeForm.rating)
    ElMessage.success('面试已完成')
    completeDialogVisible.value = false
    fetchInterviews()
  } catch (e) {
    // error handled
  } finally {
    completing.value = false
  }
}

function getApplicantName(row: Interview) {
  return row.application?.applicant?.profile?.full_name || row.application?.applicant?.email || '未知'
}

function getStatusType(status: string) {
  const types: Record<string, string> = {
    'pending': 'warning',
    'accepted': 'success',
    'rejected': 'danger',
    'completed': 'success',
    'cancelled': 'info'
  }
  return types[status] || ''
}

function getStatusText(status: string) {
  const texts: Record<string, string> = {
    'pending': '待确认',
    'accepted': '已接受',
    'rejected': '已拒绝',
    'completed': '已完成',
    'cancelled': '已取消'
  }
  return texts[status] || status
}

function formatDateTime(date: string) {
  return dayjs(date).format('YYYY-MM-DD HH:mm')
}

onMounted(() => {
  fetchInterviews()
})
</script>

<style scoped>
.applicant-info {
  display: flex;
  align-items: center;
  gap: 10px;
}

.applicant-info .name {
  font-weight: 500;
}

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
