<template>
  <div class="project-list">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>项目列表</span>
          <el-button
            v-if="userStore.hasRole(['client', 'admin'])"
            type="primary"
            @click="$router.push('/projects/create')"
          >
            <el-icon><Plus /></el-icon>新建项目
          </el-button>
        </div>
      </template>

      <div class="filter-bar">
        <el-select v-model="filters.status" placeholder="状态" clearable style="width: 140px">
          <el-option label="待审核" value="pending" />
          <el-option label="已通过" value="approved" />
          <el-option label="已分配" value="assigned" />
          <el-option label="进行中" value="in_progress" />
          <el-option label="审核中" value="review" />
          <el-option label="已完成" value="completed" />
          <el-option label="已取消" value="cancelled" />
        </el-select>
        <el-input
          v-model="filters.keyword"
          placeholder="搜索项目名称"
          clearable
          style="width: 200px"
          @keyup.enter="loadData"
        >
          <template #prefix><el-icon><Search /></el-icon></template>
        </el-input>
        <el-button type="primary" @click="loadData">
          <el-icon><Search /></el-icon>搜索
        </el-button>
      </div>

      <el-table :data="projects" stripe v-loading="loading">
        <el-table-column prop="title" label="项目名称" min-width="150" />
        <el-table-column prop="source_lang" label="源语言" width="100" />
        <el-table-column prop="target_lang" label="目标语言" width="100" />
        <el-table-column prop="word_count" label="字数" width="100" />
        <el-table-column prop="total_amount" label="金额(元)" width="120">
          <template #default="{ row }">{{ row.total_amount?.toFixed(2) }}</template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="120">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">{{ getStatusText(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="deadline" label="截止日期" width="160">
          <template #default="{ row }">{{ formatDate(row.deadline) }}</template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="viewDetail(row)">详情</el-button>
            <el-button
              v-if="row.status === 'pending' && userStore.hasRole(['pm', 'admin'])"
              type="success"
              link
              @click="handleApprove(row)"
            >通过</el-button>
            <el-button
              v-if="row.status === 'approved' && userStore.hasRole(['pm', 'admin'])"
              type="warning"
              link
              @click="showAssignDialog(row)"
            >派单</el-button>
            <el-button
              v-if="row.status === 'in_progress' && row.translator_id === userStore.userInfo?.id"
              type="primary"
              link
              @click="handleSubmit(row)"
            >提交审核</el-button>
            <el-button
              v-if="row.review_status === 'approved' && row.status === 'review' && userStore.hasRole(['pm', 'admin'])"
              type="success"
              link
              @click="handleComplete(row)"
            >完成</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.pageSize"
          :total="pagination.total"
          :page-sizes="[10, 20, 50]"
          layout="total, sizes, prev, pager, next, jumper"
          @change="loadData"
        />
      </div>
    </el-card>

    <el-dialog v-model="assignDialog.visible" title="分配译者" width="600px">
      <div class="recommend-section">
        <el-button type="primary" @click="loadRecommendations" :loading="recLoading">
          推荐译者
        </el-button>
        <el-button @click="handleAutoAssign" :loading="autoLoading">
          自动派单
        </el-button>
      </div>
      <el-table :data="translators" height="300" v-loading="recLoading">
        <el-table-column prop="username" label="用户名" />
        <el-table-column prop="real_name" label="姓名" />
        <el-table-column prop="rating" label="评分" width="80" />
        <el-table-column prop="completed_count" label="完成数" width="80" />
        <el-table-column label="匹配度" width="100">
          <template #default="{ row }">
            <span v-if="row.score">{{ row.score }}分</span>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="100">
          <template #default="{ row }">
            <el-button type="primary" link @click="assignTranslator(row)">选择</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useUserStore } from '@/stores/user'
import {
  listProjects, approveProject, assignTranslator as apiAssign,
  submitForReview, completeProject, recommendTranslators, autoAssignTranslator
} from '@/api/project'
import dayjs from 'dayjs'

const router = useRouter()
const userStore = useUserStore()

const projects = ref<any[]>([])
const loading = ref(false)
const filters = reactive({ status: '', keyword: '' })
const pagination = reactive({ page: 1, pageSize: 10, total: 0 })

const assignDialog = reactive({ visible: false, projectId: 0 })
const translators = ref<any[]>([])
const recLoading = ref(false)
const autoLoading = ref(false)

async function loadData() {
  loading.value = true
  try {
    const res = await listProjects({
      status: filters.status,
      keyword: filters.keyword,
      page: pagination.page,
      page_size: pagination.pageSize
    })
    if (res) {
      if (Array.isArray(res)) {
        projects.value = res
      } else {
        projects.value = res.list || []
        pagination.total = res.total || 0
      }
    }
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

function viewDetail(row: any) {
  router.push(`/projects/${row.id}`)
}

async function handleApprove(row: any) {
  try {
    await ElMessageBox.confirm('确定通过该项目吗？', '提示', { type: 'warning' })
    await approveProject(row.id)
    ElMessage.success('已通过')
    loadData()
  } catch (_) {}
}

function showAssignDialog(row: any) {
  assignDialog.visible = true
  assignDialog.projectId = row.id
  translators.value = []
}

async function loadRecommendations() {
  recLoading.value = true
  try {
    const res = await recommendTranslators(assignDialog.projectId)
    translators.value = res || []
  } catch (e) {
    console.error(e)
  } finally {
    recLoading.value = false
  }
}

async function handleAutoAssign() {
  autoLoading.value = true
  try {
    const res = await autoAssignTranslator(assignDialog.projectId)
    ElMessage.success(`已自动分配: ${res.translator.username}`)
    assignDialog.visible = false
    loadData()
  } catch (e) {
    console.error(e)
  } finally {
    autoLoading.value = false
  }
}

async function assignTranslator(row: any) {
  try {
    await ElMessageBox.confirm(`确定分配给 ${row.username} 吗？`, '提示')
    await apiAssign(assignDialog.projectId, row.User?.id || row.id)
    ElMessage.success('分配成功')
    assignDialog.visible = false
    loadData()
  } catch (_) {}
}

async function handleSubmit(row: any) {
  try {
    await ElMessageBox.confirm('确定提交审核吗？', '提示')
    await submitForReview(row.id)
    ElMessage.success('已提交审核')
    loadData()
  } catch (_) {}
}

async function handleComplete(row: any) {
  try {
    await ElMessageBox.confirm('确定完成该项目吗？', '提示')
    await completeProject(row.id)
    ElMessage.success('项目已完成')
    loadData()
  } catch (_) {}
}

function getStatusType(status: string) {
  const map: Record<string, string> = {
    pending: 'warning', approved: 'primary', assigned: 'info',
    in_progress: '', review: 'warning', completed: 'success', cancelled: 'danger'
  }
  return map[status] || ''
}

function getStatusText(status: string) {
  const map: Record<string, string> = {
    pending: '待审核', approved: '已通过', assigned: '已分配',
    in_progress: '进行中', review: '审核中', completed: '已完成', cancelled: '已取消'
  }
  return map[status] || status
}

function formatDate(date: string) {
  return dayjs(date).format('YYYY-MM-DD HH:mm')
}

onMounted(loadData)
</script>

<style lang="scss" scoped>
.project-list {
  .filter-bar {
    display: flex;
    gap: 12px;
    margin-bottom: 16px;
  }

  .pagination {
    margin-top: 16px;
    display: flex;
    justify-content: flex-end;
  }

  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .recommend-section {
    margin-bottom: 16px;
    display: flex;
    gap: 12px;
  }
}
</style>
