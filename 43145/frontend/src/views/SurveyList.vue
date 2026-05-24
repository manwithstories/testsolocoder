<template>
  <div class="page-container">
    <div class="page-header">
      <h1 class="page-title">问卷列表</h1>
      <el-button type="primary" :icon="Plus" @click="handleCreate">创建问卷</el-button>
    </div>

    <div class="card">
      <div class="card-body">
        <el-form :inline="true" :model="filterForm" class="filter-form">
          <el-form-item label="状态">
            <el-select v-model="filterForm.status" placeholder="全部" clearable @change="loadSurveys">
              <el-option label="草稿" :value="1" />
              <el-option label="已发布" :value="2" />
              <el-option label="已关闭" :value="3" />
            </el-select>
          </el-form-item>
          <el-form-item label="分类">
            <el-select v-model="filterForm.category" placeholder="全部" clearable @change="loadSurveys">
              <el-option label="市场调研" value="market_research" />
              <el-option label="学术研究" value="academic" />
              <el-option label="活动投票" value="event_voting" />
              <el-option label="其他" value="other" />
            </el-select>
          </el-form-item>
          <el-form-item label="关键词">
            <el-input v-model="filterForm.keyword" placeholder="搜索问卷标题" clearable @input="handleSearch" />
          </el-form-item>
        </el-form>
      </div>
    </div>

    <div v-loading="loading" class="survey-grid">
      <el-row :gutter="20">
        <el-col v-for="survey in surveys" :key="survey.id" :xs="24" :sm="12" :md="8" :lg="6">
          <div class="survey-card" @click="handleEdit(survey.id)">
            <div class="survey-card-cover">
              {{ survey.title.charAt(0) }}
            </div>
            <div class="survey-card-content">
              <div class="survey-card-title">{{ survey.title }}</div>
              <div class="survey-card-meta">
                <span :class="['status-tag', getStatusClass(survey.status)]">
                  {{ getStatusText(survey.status) }}
                </span>
                <span>{{ survey.response_count }} 份答卷</span>
              </div>
              <div class="survey-card-actions" @click.stop>
                <el-button size="small" :icon="Edit" @click="handleEdit(survey.id)">编辑</el-button>
                <el-button size="small" :icon="DataLine" @click="handleStatistics(survey.id)">统计</el-button>
                <el-dropdown trigger="click" @command="(cmd: string) => handleCommand(cmd, survey)">
                  <el-button size="small" :icon="MoreFilled" />
                  <template #dropdown>
                    <el-dropdown-menu>
                      <el-dropdown-item command="copy" :icon="Document">复制</el-dropdown-item>
                      <el-dropdown-item command="publish" :icon="Promotion" v-if="survey.status === 1">发布</el-dropdown-item>
                      <el-dropdown-item command="close" :icon="SwitchButton" v-if="survey.status === 2">关闭</el-dropdown-item>
                      <el-dropdown-item command="delete" :icon="Delete" divided style="color: #f56c6c;">删除</el-dropdown-item>
                    </el-dropdown-menu>
                  </template>
                </el-dropdown>
              </div>
            </div>
          </div>
        </el-col>
      </el-row>
    </div>

    <el-empty v-if="!loading && surveys.length === 0" description="暂无问卷" />

    <el-pagination
      v-if="total > 0"
      class="pagination"
      v-model:current-page="filterForm.page"
      v-model:page-size="filterForm.page_size"
      :total="total"
      :page-sizes="[10, 20, 50, 100]"
      layout="total, sizes, prev, pager, next, jumper"
      @size-change="loadSurveys"
      @current-change="loadSurveys"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Plus, Edit, DataLine, MoreFilled, Document, Promotion, SwitchButton, Delete
} from '@element-plus/icons-vue'
import { surveyApi } from '@/api/survey'
import type { Survey } from '@/types'

const router = useRouter()

const loading = ref(false)
const surveys = ref<Survey[]>([])
const total = ref(0)

const filterForm = reactive({
  page: 1,
  page_size: 12,
  status: undefined as number | undefined,
  category: '',
  keyword: ''
})

const loadSurveys = async () => {
  loading.value = true
  try {
    const res = await surveyApi.list(filterForm)
    surveys.value = res.items
    total.value = res.total
  } catch (error) {
    console.error('Failed to load surveys')
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  filterForm.page = 1
  loadSurveys()
}

const handleCreate = () => {
  router.push('/surveys/create')
}

const handleEdit = (id: number) => {
  router.push(`/surveys/${id}/edit`)
}

const handleStatistics = (id: number) => {
  router.push(`/surveys/${id}/statistics`)
}

const handleCommand = async (command: string, survey: Survey) => {
  switch (command) {
    case 'copy':
      try {
        await surveyApi.copy(survey.id)
        ElMessage.success('问卷复制成功')
        loadSurveys()
      } catch (e: any) {
        ElMessage.error(e.message || '复制失败')
      }
      break
    case 'publish':
      try {
        await surveyApi.publish(survey.id)
        ElMessage.success('问卷已发布')
        loadSurveys()
      } catch (e: any) {
        ElMessage.error(e.message || '发布失败')
      }
      break
    case 'close':
      try {
        await ElMessageBox.confirm('确定要关闭此问卷吗？关闭后无法再收集答卷。', '确认关闭', {
          type: 'warning'
        })
        await surveyApi.close(survey.id)
        ElMessage.success('问卷已关闭')
        loadSurveys()
      } catch (e: any) {
        if (e !== 'cancel') {
          ElMessage.error(e.message || '关闭失败')
        }
      }
      break
    case 'delete':
      try {
        await ElMessageBox.confirm('确定要删除此问卷吗？此操作不可恢复。', '确认删除', {
          type: 'warning'
        })
        await surveyApi.remove(survey.id)
        ElMessage.success('问卷已删除')
        loadSurveys()
      } catch (e: any) {
        if (e !== 'cancel') {
          ElMessage.error(e.message || '删除失败')
        }
      }
      break
  }
}

const getStatusText = (status: number) => {
  const map: Record<number, string> = {
    1: '草稿',
    2: '已发布',
    3: '已关闭'
  }
  return map[status] || '未知'
}

const getStatusClass = (status: number) => {
  const map: Record<number, string> = {
    1: 'status-draft',
    2: 'status-published',
    3: 'status-closed'
  }
  return map[status] || ''
}

onMounted(loadSurveys)
</script>

<style scoped>
.survey-grid {
  margin-bottom: 20px;
}

.pagination {
  justify-content: center;
  margin-top: 20px;
}

.filter-form {
  margin-bottom: 0;
}

.survey-card-actions {
  margin-top: 12px;
  display: flex;
  gap: 8px;
}
</style>
