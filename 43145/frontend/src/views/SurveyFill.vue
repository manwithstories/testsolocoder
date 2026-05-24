<template>
  <div class="fill-survey-container">
    <div v-if="showPasswordDialog">
      <div class="login-container">
        <div class="login-card">
          <h2 class="login-title">问卷密码验证</h2>
          <el-form :model="passwordForm" @submit.prevent>
            <el-form-item>
              <el-input
                v-model="passwordForm.password"
                type="password"
                placeholder="请输入问卷密码"
                show-password
                size="large"
                @keyup.enter="handlePasswordSubmit"
              />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" size="large" :loading="verifying" @click="handlePasswordSubmit" style="width: 100%;">
                验证
              </el-button>
            </el-form-item>
          </el-form>
        </div>
      </div>
    </div>

    <div v-else-if="!survey">
      <div class="login-container">
        <el-empty description="问卷不存在或已关闭" />
      </div>
    </div>

    <div v-else class="preview-container">
      <h2 class="preview-title">{{ survey.title }}</h2>
      <p class="preview-description">{{ survey.description }}</p>

      <div v-for="(question, index) in questions" :key="question.id" class="fill-question-card">
        <div style="margin-bottom: 12px;">
          <span style="font-size: 16px; font-weight: 500;">Q{{ index + 1 }}. {{ question.title }}</span>
          <span v-if="question.is_required" class="fill-question-required">*</span>
        </div>
        <div v-if="question.description" style="color: #909399; margin-bottom: 12px;">
          {{ question.description }}
        </div>

        <div v-if="question.type === 'single_choice'">
          <el-radio-group v-model="answers[question.id]">
            <div v-for="opt in question.options" :key="opt.id" class="fill-option"
                 :class="{ selected: answers[question.id] === opt.id }"
                 @click="answers[question.id] = opt.id; handleAutoSave()">
              <el-radio :label="opt.id" :disabled="true" />
              <span>{{ opt.text }}</span>
              <el-input
                v-if="opt.is_other && answers[question.id] === opt.id"
                v-model="otherAnswers[question.id]"
                placeholder="请输入其他"
                size="small"
                style="width: 200px; margin-left: 8px;"
                @input="handleAutoSave"
              />
            </div>
          </el-radio-group>
        </div>

        <div v-else-if="question.type === 'multi_choice'">
          <div v-for="opt in question.options" :key="opt.id" class="fill-option"
               :class="{ selected: (multiAnswers[question.id] || []).includes(opt.id) }"
               @click="toggleMultiAnswer(question.id, opt.id)">
            <el-checkbox :model-value="(multiAnswers[question.id] || []).includes(opt.id)" :disabled="true" />
            <span>{{ opt.text }}</span>
          </div>
        </div>

        <div v-else-if="question.type === 'fill_in'">
          <el-input
            v-model="textAnswers[question.id]"
            type="textarea"
            :rows="3"
            :placeholder="question.placeholder || '请输入答案'"
            @input="handleAutoSave"
          />
        </div>

        <div v-else-if="question.type === 'rating'">
          <el-rate
            v-model="ratingAnswers[question.id]"
            :max="question.max_value || 5"
            :min="question.min_value || 1"
            show-text
            @change="handleAutoSave"
          />
        </div>

        <div v-else-if="question.type === 'ranking'">
          <draggable
            v-model="rankingAnswers[question.id]"
            item-key="id"
            animation="200"
            @end="handleAutoSave"
          >
            <template #item="{ element }">
              <div class="fill-option" style="cursor: move;">
                <el-icon><Rank /></el-icon>
                <span style="margin-left: 8px;">{{ element.text }}</span>
              </div>
            </template>
          </draggable>
        </div>
      </div>

      <div style="text-align: center; margin-top: 30px;">
        <el-button type="primary" size="large" :loading="submitting" @click="handleSubmit">
          提交问卷
        </el-button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Rank } from '@element-plus/icons-vue'
import { publicApi } from '@/api/response'
import { questionApi } from '@/api/question'
import { surveyApi } from '@/api/survey'
import type { Survey, Question } from '@/types'
import draggable from 'vuedraggable'

const route = useRoute()
const router = useRouter()

const token = computed(() => route.params.token as string)
const surveyId = ref(0)

const survey = ref<Survey>()
const questions = ref<Question[]>([])
const sessionId = ref('')
const submitting = ref(false)
const verifying = ref(false)
const showPasswordDialog = ref(false)

const passwordForm = reactive({
  password: ''
})

const answers = reactive<Record<number, number>>({})
const multiAnswers = reactive<Record<number, number[]>>({})
const textAnswers = reactive<Record<number, string>>({})
const ratingAnswers = reactive<Record<number, number>>({})
const rankingAnswers = reactive<Record<number, any[]>>({})
const otherAnswers = reactive<Record<number, string>>({})

let autoSaveTimer: any = null

const loadSurvey = async () => {
  try {
    const dist = await publicApi.getSurveyByToken(token.value)
    surveyId.value = dist.survey_id

    const [surveyRes, questionsRes] = await Promise.all([
      surveyApi.getById(dist.survey_id),
      questionApi.list(dist.survey_id)
    ])

    survey.value = surveyRes
    questions.value = questionsRes

    if (surveyRes.has_password) {
      showPasswordDialog.value = true
      return
    }

    const result: any = await publicApi.startResponse(dist.survey_id, {}, token.value)
    sessionId.value = result.session_id

    for (const q of questions.value) {
      if (q.type === 'ranking' && q.options) {
        rankingAnswers[q.id] = [...q.options]
      }
      if (q.type === 'multi_choice') {
        multiAnswers[q.id] = []
      }
    }
  } catch (error: any) {
    ElMessage.error(error.message || '加载问卷失败')
  }
}

const handlePasswordSubmit = async () => {
  if (!passwordForm.password) {
    ElMessage.warning('请输入密码')
    return
  }

  verifying.value = true
  try {
    await publicApi.validateAccess(surveyId.value, { password: passwordForm.password })
    showPasswordDialog.value = false

    const result: any = await publicApi.startResponse(surveyId.value, {}, token.value)
    sessionId.value = result.session_id
  } catch (e: any) {
    ElMessage.error(e.message || '密码错误')
  } finally {
    verifying.value = false
  }
}

const toggleMultiAnswer = (questionId: number, optionId: number) => {
  if (!multiAnswers[questionId]) {
    multiAnswers[questionId] = []
  }
  const idx = multiAnswers[questionId].indexOf(optionId)
  if (idx >= 0) {
    multiAnswers[questionId].splice(idx, 1)
  } else {
    multiAnswers[questionId].push(optionId)
  }
  handleAutoSave()
}

const handleAutoSave = () => {
  if (autoSaveTimer) clearTimeout(autoSaveTimer)
  autoSaveTimer = setTimeout(async () => {
    if (!sessionId.value || !surveyId.value) return

    try {
      const answerData = buildAnswers()
      await publicApi.saveProgress(surveyId.value, {
        session_id: sessionId.value,
        answers: answerData
      })
    } catch (e) {
      console.error('Auto save failed')
    }
  }, 2000)
}

const buildAnswers = () => {
  const answerList: any[] = []

  for (const q of questions.value) {
    const answer: any = { question_id: q.id }

    if (q.type === 'single_choice' && answers[q.id]) {
      answer.option_id = answers[q.id]
      if (otherAnswers[q.id]) {
        answer.text_value = otherAnswers[q.id]
      }
    } else if (q.type === 'multi_choice' && multiAnswers[q.id]) {
      for (const optId of multiAnswers[q.id]) {
        answerList.push({ question_id: q.id, option_id: optId })
      }
      continue
    } else if (q.type === 'fill_in') {
      answer.text_value = textAnswers[q.id] || ''
    } else if (q.type === 'rating') {
      answer.numeric_value = ratingAnswers[q.id] || 0
    } else if (q.type === 'ranking' && rankingAnswers[q.id]) {
      answer.ranking_order = JSON.stringify(rankingAnswers[q.id].map((o: any) => o.id))
    }

    answerList.push(answer)
  }

  return answerList
}

const validateAnswers = () => {
  for (const q of questions.value) {
    if (q.is_required) {
      if (q.type === 'single_choice' && !answers[q.id]) {
        ElMessage.warning(`请回答第 ${questions.value.indexOf(q) + 1} 题`)
        return false
      }
      if (q.type === 'multi_choice' && (!multiAnswers[q.id] || multiAnswers[q.id].length === 0)) {
        ElMessage.warning(`请回答第 ${questions.value.indexOf(q) + 1} 题`)
        return false
      }
      if (q.type === 'fill_in' && !textAnswers[q.id]) {
        ElMessage.warning(`请回答第 ${questions.value.indexOf(q) + 1} 题`)
        return false
      }
      if (q.type === 'rating' && !ratingAnswers[q.id]) {
        ElMessage.warning(`请回答第 ${questions.value.indexOf(q) + 1} 题`)
        return false
      }
    }
  }
  return true
}

const handleSubmit = async () => {
  if (!validateAnswers()) return

  try {
    await ElMessageBox.confirm('确定要提交问卷吗？提交后无法修改。', '确认提交', {
      type: 'warning'
    })
  } catch {
    return
  }

  submitting.value = true
  try {
    const answerData = buildAnswers()
    await publicApi.submitResponse(surveyId.value, {
      session_id: sessionId.value,
      answers: answerData
    })
    ElMessage.success('问卷提交成功')
    router.push('/survey/thank-you')
  } catch (e: any) {
    ElMessage.error(e.message || '提交失败')
  } finally {
    submitting.value = false
  }
}

onMounted(loadSurvey)
</script>
