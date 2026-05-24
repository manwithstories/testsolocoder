<template>
  <div class="page-container">
    <div class="page-header">
      <h1 class="page-title">{{ isEdit ? '编辑问卷' : '创建问卷' }}</h1>
      <div>
        <el-button :icon="View" @click="showPreview = true">预览</el-button>
        <el-button type="primary" :icon="Check" :loading="saving" @click="handleSave">保存</el-button>
      </div>
    </div>

    <el-row :gutter="20">
      <el-col :span="16">
        <div class="card">
          <div class="card-body">
            <el-form :model="surveyForm" label-width="100px">
              <el-form-item label="问卷标题" required>
                <el-input v-model="surveyForm.title" placeholder="请输入问卷标题" maxlength="200" show-word-limit />
              </el-form-item>
              <el-form-item label="问卷描述">
                <el-input v-model="surveyForm.description" type="textarea" :rows="3" placeholder="请输入问卷描述" />
              </el-form-item>
              <el-form-item label="封面图片">
                <el-input v-model="surveyForm.cover_image" placeholder="封面图片URL" />
              </el-form-item>
              <el-form-item label="分类">
                <el-select v-model="surveyForm.category" placeholder="请选择分类" style="width: 100%;">
                  <el-option label="市场调研" value="market_research" />
                  <el-option label="学术研究" value="academic" />
                  <el-option label="活动投票" value="event_voting" />
                  <el-option label="其他" value="other" />
                </el-select>
              </el-form-item>
              <el-form-item label="标签">
                <el-input v-model="surveyForm.tags" placeholder="用逗号分隔多个标签" />
              </el-form-item>
            </el-form>
          </div>
        </div>

        <div class="card">
          <div class="card-header" style="display: flex; justify-content: space-between; align-items: center;">
            <span style="font-weight: 600;">题目列表</span>
            <el-button type="primary" size="small" :icon="Plus" @click="addQuestion">添加题目</el-button>
          </div>
          <div class="card-body">
            <draggable
              v-model="questions"
              item-key="id"
              handle=".drag-handle"
              animation="200"
              ghost-class="sortable-ghost"
              chosen-class="sortable-chosen"
              @change="handleQuestionReorder"
            >
              <template #item="{ element }">
                <div class="question-item" v-if="element.id">
                  <div class="question-header">
                    <div style="display: flex; align-items: center; gap: 8px;">
                      <span class="drag-handle"><el-icon><Rank /></el-icon></span>
                      <span>Q{{ questions.indexOf(element) + 1 }}.</span>
                      <el-tag size="small">{{ getQuestionTypeText(element.type) }}</el-tag>
                    </div>
                    <div>
                      <el-button size="small" :icon="Edit" @click="editQuestion(element)" />
                      <el-button size="small" type="danger" :icon="Delete" @click="removeQuestion(element)" />
                    </div>
                  </div>
                  <div class="question-title">{{ element.title || '未命名题目' }}</div>
                  <div v-if="element.options && element.options.length > 0" style="margin-top: 8px; color: #909399; font-size: 13px;">
                    选项: {{ element.options.map((o: any) => o.text).join(' | ') }}
                  </div>
                  <div v-if="element.is_required" style="margin-top: 4px;">
                    <el-tag size="small" type="danger">必填</el-tag>
                  </div>
                </div>
              </template>
            </draggable>

            <el-empty v-if="questions.length === 0" description="暂无题目，点击上方按钮添加" />
          </div>
        </div>
      </el-col>

      <el-col :span="8">
        <div class="card">
          <div class="card-header" style="font-weight: 600;">问卷设置</div>
          <div class="card-body">
            <el-form :model="surveySettings" label-width="100px">
              <el-form-item label="开始时间">
                <el-date-picker v-model="surveySettings.start_time" type="datetime" placeholder="选择开始时间" />
              </el-form-item>
              <el-form-item label="结束时间">
                <el-date-picker v-model="surveySettings.end_time" type="datetime" placeholder="选择结束时间" />
              </el-form-item>
              <el-form-item label="匿名填写">
                <el-switch v-model="surveySettings.anonymous" />
              </el-form-item>
              <el-form-item label="密码保护">
                <el-input v-model="surveySettings.password" type="password" placeholder="留空表示不设密码" show-password />
              </el-form-item>
              <el-form-item label="最大答卷数">
                <el-input-number v-model="surveySettings.max_responses" :min="0" :step="1" />
                <span style="color: #909399; font-size: 12px; margin-left: 8px;">0表示不限</span>
              </el-form-item>
              <el-form-item label="每人限答">
                <el-input-number v-model="surveySettings.max_per_user" :min="0" :step="1" />
                <span style="color: #909399; font-size: 12px; margin-left: 8px;">0表示不限</span>
              </el-form-item>
              <el-form-item label="需要登录">
                <el-switch v-model="surveySettings.requires_login" />
              </el-form-item>
              <el-form-item label="允许续填">
                <el-switch v-model="surveySettings.allow_resume" />
              </el-form-item>
            </el-form>
          </div>
        </div>
      </el-col>
    </el-row>

    <el-dialog v-model="showQuestionDialog" title="编辑题目" width="600px" :close-on-click-modal="false">
      <el-form :model="currentQuestion" label-width="80px">
        <el-form-item label="题目" required>
          <el-input v-model="currentQuestion.title" placeholder="请输入题目" />
        </el-form-item>
        <el-form-item label="题型" required>
          <el-select v-model="currentQuestion.type" style="width: 100%;" @change="handleTypeChange">
            <el-option label="单选题" value="single_choice" />
            <el-option label="多选题" value="multi_choice" />
            <el-option label="填空题" value="fill_in" />
            <el-option label="评分题" value="rating" />
            <el-option label="排序题" value="ranking" />
            <el-option label="矩阵题" value="matrix" />
          </el-select>
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="currentQuestion.description" type="textarea" :rows="2" placeholder="题目描述/说明" />
        </el-form-item>
        <el-form-item label="必填">
          <el-switch v-model="currentQuestion.is_required" />
        </el-form-item>

        <template v-if="['single_choice', 'multi_choice', 'ranking', 'matrix'].includes(currentQuestion.type)">
          <el-form-item label="选项">
            <div style="width: 100%;">
              <div v-for="(opt, index) in currentQuestion.options" :key="index" style="display: flex; gap: 8px; margin-bottom: 8px;">
                <el-input v-model="opt.text" :placeholder="`选项 ${index + 1}`" style="flex: 1;" />
                <el-checkbox v-model="opt.is_other">其他</el-checkbox>
                <el-button type="danger" :icon="Delete" circle @click="removeOption(index)" />
              </div>
              <el-button size="small" :icon="Plus" @click="addOption">添加选项</el-button>
            </div>
          </el-form-item>
        </template>

        <template v-if="currentQuestion.type === 'rating'">
          <el-form-item label="最小值">
            <el-input-number v-model="currentQuestion.min_value" :min="0" />
          </el-form-item>
          <el-form-item label="最大值">
            <el-input-number v-model="currentQuestion.max_value" :min="1" />
          </el-form-item>
        </template>

        <el-form-item label="占位符">
          <el-input v-model="currentQuestion.placeholder" placeholder="输入框占位符" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showQuestionDialog = false">取消</el-button>
        <el-button type="primary" @click="saveQuestion">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="showPreview" title="问卷预览" width="800px">
      <div class="preview-container">
        <h2 class="preview-title">{{ surveyForm.title || '问卷标题' }}</h2>
        <p class="preview-description">{{ surveyForm.description }}</p>
        <div v-for="(q, index) in questions.filter((q: any) => q.id)" :key="q.id" style="margin-bottom: 24px;">
          <h3>Q{{ index + 1 }}. {{ q.title }}</h3>
          <el-tag size="small" style="margin-left: 8px;">{{ getQuestionTypeText(q.type) }}</el-tag>
          <div v-if="q.is_required" style="color: #f56c6c;">*必填</div>
          <div v-if="q.description" style="color: #909399; margin: 8px 0;">{{ q.description }}</div>
          <div v-if="q.options && q.options.length > 0" style="margin-top: 12px;">
            <el-radio v-for="opt in q.options" :key="opt.id || opt._id" :label="opt.text" disabled>
              {{ opt.text }}
              <el-tag v-if="opt.is_other" size="small" type="info" style="margin-left: 8px;">其他</el-tag>
            </el-radio>
          </div>
          <el-input v-if="q.type === 'fill_in'" disabled placeholder="请输入答案" />
          <el-rate v-if="q.type === 'rating'" disabled :max="q.max_value || 5" />
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import {
  View, Check, Plus, Edit, Delete, Rank
} from '@element-plus/icons-vue'
import { surveyApi } from '@/api/survey'
import { questionApi } from '@/api/question'
import type { Survey, Question } from '@/types'
import draggable from 'vuedraggable'

const route = useRoute()
const router = useRouter()

const isEdit = computed(() => !!route.params.id)
const surveyId = computed(() => Number(route.params.id) || 0)

const saving = ref(false)
const showPreview = ref(false)
const showQuestionDialog = ref(false)
const questions = ref<any[]>([])
let tempIdCounter = 0

const surveyForm = reactive({
  title: '',
  description: '',
  cover_image: '',
  category: '',
  tags: ''
})

const surveySettings = reactive({
  start_time: undefined as string | undefined,
  end_time: undefined as string | undefined,
  anonymous: false,
  password: '',
  max_responses: 0,
  max_per_user: 1,
  requires_login: false,
  allow_resume: true
})

const currentQuestion = reactive<any>({
  id: undefined,
  title: '',
  type: 'single_choice',
  is_required: true,
  description: '',
  placeholder: '',
  min_value: 1,
  max_value: 5,
  options: []
})

const getQuestionTypeText = (type: string) => {
  const map: Record<string, string> = {
    single_choice: '单选题',
    multi_choice: '多选题',
    fill_in: '填空题',
    rating: '评分题',
    ranking: '排序题',
    matrix: '矩阵题'
  }
  return map[type] || type
}

const loadSurvey = async () => {
  if (!isEdit.value) return

  try {
    const survey = await surveyApi.getById(surveyId.value)
    Object.assign(surveyForm, {
      title: survey.title,
      description: survey.description,
      cover_image: survey.cover_image,
      category: survey.category,
      tags: survey.tags
    })
    Object.assign(surveySettings, {
      start_time: survey.start_time,
      end_time: survey.end_time,
      anonymous: survey.anonymous,
      password: '',
      max_responses: survey.max_responses,
      max_per_user: survey.max_per_user,
      requires_login: survey.requires_login,
      allow_resume: survey.allow_resume
    })

    if (survey.questions) {
      questions.value = survey.questions.map(q => ({
        ...q,
        options: q.options || []
      }))
    }
  } catch (error) {
    console.error('Failed to load survey')
  }
}

const addQuestion = () => {
  tempIdCounter++
  questions.value.push({
    _id: `temp_${tempIdCounter}`,
    id: undefined,
    title: '',
    type: 'single_choice',
    is_required: true,
    description: '',
    placeholder: '',
    min_value: 1,
    max_value: 5,
    options: [{ text: '选项1', is_other: false }, { text: '选项2', is_other: false }]
  })
}

const editQuestion = (q: any) => {
  Object.assign(currentQuestion, {
    ...q,
    options: q.options ? q.options.map((o: any) => ({ ...o })) : []
  })
  showQuestionDialog.value = true
}

const removeQuestion = async (q: any) => {
  try {
    if (q.id) {
      await questionApi.remove(q.id)
    }
    questions.value = questions.value.filter(item => item.id !== q.id && item._id !== q._id)
    ElMessage.success('题目已删除')
  } catch (e: any) {
    ElMessage.error(e.message || '删除失败')
  }
}

const handleQuestionReorder = async () => {
  for (let i = 0; i < questions.value.length; i++) {
    const q = questions.value[i]
    if (q.id) {
      try {
        await questionApi.reorder(q.id, i)
      } catch (e) {
        console.error('Failed to reorder question')
      }
    }
  }
}

const addOption = () => {
  currentQuestion.options.push({ text: '', is_other: false })
}

const removeOption = (index: number) => {
  currentQuestion.options.splice(index, 1)
}

const handleTypeChange = () => {
  if (['single_choice', 'multi_choice', 'ranking', 'matrix'].includes(currentQuestion.type)) {
    if (!currentQuestion.options || currentQuestion.options.length === 0) {
      currentQuestion.options = [{ text: '选项1', is_other: false }, { text: '选项2', is_other: false }]
    }
  }
}

const saveQuestion = () => {
  if (!currentQuestion.title) {
    ElMessage.warning('请输入题目')
    return
  }
  if (['single_choice', 'multi_choice', 'ranking', 'matrix'].includes(currentQuestion.type) &&
      (!currentQuestion.options || currentQuestion.options.length < 2)) {
    ElMessage.warning('至少需要2个选项')
    return
  }

  const index = questions.value.findIndex((q: any) =>
    (q.id && q.id === currentQuestion.id) || (q._id && q._id === currentQuestion._id)
  )

  if (index >= 0) {
    questions.value[index] = { ...currentQuestion }
  }

  showQuestionDialog.value = false
  ElMessage.success('题目已更新')
}

const handleSave = async () => {
  if (!surveyForm.title) {
    ElMessage.warning('请输入问卷标题')
    return
  }

  saving.value = true
  try {
    const surveyData = {
      ...surveyForm,
      ...surveySettings
    }

    let savedSurvey: Survey
    if (isEdit.value) {
      await surveyApi.update(surveyId.value, surveyData)
      savedSurvey = await surveyApi.getById(surveyId.value)
    } else {
      savedSurvey = await surveyApi.create(surveyData)
    }

    const questionData = questions.value
      .filter((q: any) => q.title)
      .map((q: any, index: number) => ({
        ...q,
        order_index: index,
        options: q.options || []
      }))

    if (questionData.length > 0) {
      await questionApi.batchCreate(savedSurvey.id, { questions: questionData })
    }

    ElMessage.success('保存成功')
    if (!isEdit.value) {
      router.push(`/surveys/${savedSurvey.id}/edit`)
    }
  } catch (error: any) {
    ElMessage.error(error.message || '保存失败')
  } finally {
    saving.value = false
  }
}

onMounted(loadSurvey)
</script>

<style scoped>
.question-item {
  background: #fff;
  border-radius: 8px;
  margin-bottom: 16px;
  padding: 16px;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.06);
}
</style>
