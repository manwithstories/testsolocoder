<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>质量审核</span>
        </div>
      </template>

      <div class="filter-bar">
        <el-select v-model="status" placeholder="审核状态" clearable style="width: 140px" @change="loadData">
          <el-option label="待处理" value="pending" />
          <el-option label="已通过" value="approved" />
          <el-option label="不通过" value="rejected" />
        </el-select>
      </div>

      <el-table :data="tasks" stripe v-loading="loading">
        <el-table-column prop="project.title" label="项目" min-width="150" />
        <el-table-column prop="round" label="轮次" width="80" />
        <el-table-column prop="comment" label="审核意见" min-width="200">
          <template #default="{ row }">{{ row.comment || '-' }}</template>
        </el-table-column>
        <el-table-column prop="suggestion" label="建议修改" min-width="200">
          <template #default="{ row }">{{ row.suggestion || '-' }}</template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 'approved' ? 'success' : row.status === 'rejected' ? 'danger' : 'warning'" size="small">
              {{ row.status === 'approved' ? '通过' : row.status === 'rejected' ? '不通过' : '待处理' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="160">
          <template #default="{ row }">{{ formatDate(row.created_at) }}</template>
        </el-table-column>
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button
              v-if="row.status === 'pending'"
              type="primary"
              link
              @click="handleProcess(row)"
            >处理</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="processDialog.visible" title="处理审核" width="600px">
      <el-form :model="processForm" label-width="80px">
        <el-form-item label="审核结果">
          <el-radio-group v-model="processForm.status">
            <el-radio value="approved">通过</el-radio>
            <el-radio value="rejected">不通过</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="审核意见">
          <el-input v-model="processForm.comment" type="textarea" :rows="3" placeholder="请输入审核意见" />
        </el-form-item>
        <el-form-item v-if="processForm.status === 'rejected'" label="建议修改">
          <el-input v-model="processForm.suggestion" type="textarea" :rows="3" placeholder="请输入修改建议" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="processDialog.visible = false">取消</el-button>
        <el-button type="primary" @click="submitProcess">提交</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { listReviewTasks, processReview } from '@/api/statistics'
import dayjs from 'dayjs'

const tasks = ref<any[]>([])
const loading = ref(false)
const status = ref('')

const processDialog = reactive({ visible: false, taskId: 0 })
const processForm = reactive({ status: 'approved', comment: '', suggestion: '' })

async function loadData() {
  loading.value = true
  try {
    const res = await listReviewTasks({ status: status.value }) as any
    if (Array.isArray(res)) {
      tasks.value = res
    } else {
      tasks.value = res?.list || []
    }
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

function handleProcess(row: any) {
  processDialog.visible = true
  processDialog.taskId = row.id
  processForm.status = 'approved'
  processForm.comment = ''
  processForm.suggestion = ''
}

async function submitProcess() {
  try {
    await processReview(processDialog.taskId, processForm)
    ElMessage.success('处理成功')
    processDialog.visible = false
    loadData()
  } catch (e: any) {
    ElMessage.error(e.message || '处理失败')
  }
}

function formatDate(date: string) {
  return dayjs(date).format('YYYY-MM-DD HH:mm')
}

onMounted(loadData)
</script>

<style lang="scss" scoped>
.page-container {
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .filter-bar {
    display: flex;
    gap: 12px;
    margin-bottom: 16px;
  }
}
</style>
