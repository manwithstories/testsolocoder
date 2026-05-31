<template>
  <div class="project-create">
    <el-card>
      <template #header>
        <span>创建翻译项目</span>
      </template>

      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="100px"
        style="max-width: 700px"
      >
        <el-form-item label="项目名称" prop="title">
          <el-input v-model="form.title" placeholder="请输入项目名称" maxlength="200" />
        </el-form-item>

        <el-form-item label="项目描述" prop="description">
          <el-input
            v-model="form.description"
            type="textarea"
            :rows="3"
            placeholder="请输入项目描述"
          />
        </el-form-item>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="源语言" prop="source_lang">
              <el-select v-model="form.source_lang" placeholder="请选择" style="width: 100%">
                <el-option v-for="lp in languagePairs" :key="lp.id" :label="lp.display_name || lp.source_lang + ' - ' + lp.target_lang" :value="lp.source_lang" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="目标语言" prop="target_lang">
              <el-select v-model="form.target_lang" placeholder="请选择" style="width: 100%">
                <el-option v-for="lp in languagePairs" :key="lp.id" :label="lp.display_name || lp.target_lang" :value="lp.target_lang" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>

        <el-form-item label="专业领域" prop="expertise_tag_ids">
          <el-select
            v-model="form.expertise_tag_ids"
            multiple
            placeholder="请选择专业领域"
            style="width: 100%"
          >
            <el-option v-for="tag in expertiseTags" :key="tag.id" :label="tag.name" :value="tag.id" />
          </el-select>
        </el-form-item>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="字数" prop="word_count">
              <el-input-number
                v-model="form.word_count"
                :min="1"
                :max="1000000"
                style="width: 100%"
                @change="calculateFee"
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="紧急程度" prop="urgency">
              <el-select v-model="form.urgency" placeholder="请选择" style="width: 100%" @change="calculateFee">
                <el-option label="普通" value="normal" />
                <el-option label="加急" value="fast" />
                <el-option label="特急" value="urgent" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>

        <el-form-item label="截止日期" prop="deadline">
          <el-date-picker
            v-model="form.deadline"
            type="datetime"
            placeholder="选择截止日期"
            style="width: 100%"
            value-format="YYYY-MM-DDTHH:mm:ssZ"
          />
        </el-form-item>

        <el-form-item label="费用预估">
          <div class="fee-preview" v-if="feeResult">
            <el-descriptions :column="3" border>
              <el-descriptions-item label="基础费用">{{ feeResult.base_amount?.toFixed(2) }} 元</el-descriptions-item>
              <el-descriptions-item label="加急费用">{{ feeResult.urgency_fee?.toFixed(2) }} 元</el-descriptions-item>
              <el-descriptions-item label="难度费用">{{ feeResult.difficulty_fee?.toFixed(2) }} 元</el-descriptions-item>
              <el-descriptions-item label="总计" :span="3">
                <span class="total-fee">{{ feeResult.total_amount?.toFixed(2) }} 元</span>
              </el-descriptions-item>
            </el-descriptions>
          </div>
          <span v-else class="fee-hint">填写完字数和紧急程度后自动计算</span>
        </el-form-item>

        <el-form-item>
          <el-button type="primary" :loading="loading" @click="handleSubmit">
            提交项目
          </el-button>
          <el-button @click="$router.back()">取消</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { createProject } from '@/api/project'
import { listLanguagePairs, listExpertiseTags } from '@/api/project'
import { calculateFee as apiCalculateFee } from '@/api/statistics'

const router = useRouter()
const formRef = ref<FormInstance>()
const loading = ref(false)

const languagePairs = ref<any[]>([])
const expertiseTags = ref<any[]>([])
const feeResult = ref<any>(null)

const form = reactive({
  title: '',
  description: '',
  source_lang: '',
  target_lang: '',
  expertise_tag_ids: [] as number[],
  word_count: 0,
  urgency: 'normal',
  deadline: ''
})

const rules: FormRules = {
  title: [{ required: true, message: '请输入项目名称', trigger: 'blur' }],
  source_lang: [{ required: true, message: '请选择源语言', trigger: 'change' }],
  target_lang: [{ required: true, message: '请选择目标语言', trigger: 'change' }],
  word_count: [{ required: true, message: '请输入字数', trigger: 'blur' }],
  urgency: [{ required: true, message: '请选择紧急程度', trigger: 'change' }],
  deadline: [{ required: true, message: '请选择截止日期', trigger: 'change' }]
}

async function loadOptions() {
  try {
    const [pairs, tags] = await Promise.all([
      listLanguagePairs(),
      listExpertiseTags()
    ])
    languagePairs.value = pairs || []
    expertiseTags.value = tags || []
  } catch (e) {
    console.error(e)
  }
}

async function calculateFee() {
  if (!form.word_count || !form.source_lang || !form.target_lang) {
    feeResult.value = null
    return
  }
  try {
    const res = await apiCalculateFee({
      word_count: form.word_count,
      source_lang: form.source_lang,
      target_lang: form.target_lang,
      urgency: form.urgency,
      expertise_tag_ids: form.expertise_tag_ids
    })
    feeResult.value = res
  } catch (e) {
    console.error(e)
  }
}

async function handleSubmit() {
  if (!formRef.value) return
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return

  loading.value = true
  try {
    const data = { ...form }
    if (data.deadline) {
      data.deadline = new Date(data.deadline).toISOString()
    }
    const res = await createProject(data)
    ElMessage.success('项目创建成功')
    router.push(`/projects/${res.project?.id || res.id}`)
  } catch (e: any) {
    ElMessage.error(e.message || '创建失败')
  } finally {
    loading.value = false
  }
}

onMounted(loadOptions)
</script>

<style lang="scss" scoped>
.project-create {
  .fee-preview {
    width: 100%;
  }

  .total-fee {
    font-size: 18px;
    font-weight: bold;
    color: #f56c6c;
  }

  .fee-hint {
    color: #909399;
    font-size: 13px;
  }
}
</style>
