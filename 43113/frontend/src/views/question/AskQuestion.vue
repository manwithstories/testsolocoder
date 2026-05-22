<template>
  <div class="ask-question-page">
    <el-card>
      <template #header>
        <h2>发布问题</h2>
      </template>

      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-position="top"
      >
        <el-form-item label="问题标题" prop="title">
          <el-input
            v-model="form.title"
            placeholder="请输入问题标题"
            maxlength="200"
            show-word-limit
          />
        </el-form-item>

        <el-form-item label="问题分类" prop="categoryId">
          <el-select v-model="form.categoryId" placeholder="请选择分类" style="width: 100%">
            <el-option
              v-for="category in categories"
              :key="category.id"
              :label="category.name"
              :value="category.id"
            >
              <span>{{ category.icon }} {{ category.name }}</span>
            </el-option>
          </el-select>
        </el-form-item>

        <el-form-item label="标签" prop="tagIds">
          <el-select
            v-model="form.tagIds"
            multiple
            filterable
            allow-create
            default-first-option
            placeholder="选择或输入标签"
            style="width: 100%"
          >
            <el-option
              v-for="tag in tags"
              :key="tag.id"
              :label="tag.name"
              :value="tag.id"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="悬赏积分" prop="rewardPoints">
          <el-input-number
            v-model="form.rewardPoints"
            :min="0"
            :max="userStore.userInfo?.points || 0"
            :step="10"
            controls-position="right"
          />
          <div class="tip">当前积分: {{ userStore.userInfo?.points || 0 }}</div>
        </el-form-item>

        <el-form-item label="问题内容" prop="content">
          <div class="editor-wrapper">
            <div class="editor-toolbar">
              <el-button-group>
                <el-button size="small" @click="insertMarkdown('**', '**')">加粗</el-button>
                <el-button size="small" @click="insertMarkdown('*', '*')">斜体</el-button>
                <el-button size="small" @click="insertMarkdown('`', '`')">代码</el-button>
                <el-button size="small" @click="insertMarkdown('\n```', '```\n')">代码块</el-button>
                <el-button size="small" @click="insertMarkdown('[', '](url)')">链接</el-button>
                <el-button size="small" @click="insertMarkdown('\n- ', '')">列表</el-button>
              </el-button-group>
            </div>
            <el-input
              v-model="form.content"
              type="textarea"
              :rows="10"
              placeholder="请输入问题内容，支持 Markdown 格式..."
            />
          </div>
          <div class="content-preview" v-if="form.content">
            <h4>预览</h4>
            <div v-html="renderMarkdown(form.content)"></div>
          </div>
        </el-form-item>

        <el-form-item>
          <el-button type="primary" size="large" :loading="submitting" @click="handleSubmit">
            发布问题
          </el-button>
          <el-button size="large" @click="goBack">取消</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { questionApi } from '@/api'
import type { Category, Tag } from '@/types'
import type { FormInstance, FormRules } from 'element-plus'
import MarkdownIt from 'markdown-it'

const router = useRouter()
const userStore = useUserStore()

const md = new MarkdownIt()

const formRef = ref<FormInstance>()
const submitting = ref(false)
const categories = ref<Category[]>([])
const tags = ref<Tag[]>([])

const form = reactive({
  title: '',
  content: '',
  categoryId: null as number | null,
  tagIds: [] as number[],
  rewardPoints: 0
})

const rules: FormRules = {
  title: [
    { required: true, message: '请输入问题标题', trigger: 'blur' },
    { min: 5, max: 200, message: '标题长度在 5 到 200 个字符', trigger: 'blur' }
  ],
  content: [
    { required: true, message: '请输入问题内容', trigger: 'blur' },
    { min: 10, message: '问题内容至少 10 个字符', trigger: 'blur' }
  ],
  categoryId: [
    { required: true, message: '请选择分类', trigger: 'change' }
  ]
}

const fetchCategories = async () => {
  try {
    const res = await questionApi.getCategories()
    categories.value = res.data || []
  } catch (e) {
    console.error(e)
  }
}

const fetchTags = async () => {
  try {
    const res = await questionApi.getTags()
    tags.value = res.data || []
  } catch (e) {
    console.error(e)
  }
}

const insertMarkdown = (prefix: string, suffix: string) => {
  form.content += prefix + suffix
}

const renderMarkdown = (content: string) => {
  return md.render(content)
}

const handleSubmit = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (valid) {
      submitting.value = true
      try {
        const res = await questionApi.createQuestion({
          title: form.title,
          content: form.content,
          categoryId: form.categoryId!,
          tagIds: form.tagIds,
          rewardPoints: form.rewardPoints
        })
        if (res.data) {
          router.push(`/questions/${res.data.id}`)
        }
      } catch (e) {
        console.error(e)
      } finally {
        submitting.value = false
      }
    }
  })
}

const goBack = () => {
  router.back()
}

onMounted(() => {
  fetchCategories()
  fetchTags()
})
</script>

<style scoped lang="scss">
.ask-question-page {
  max-width: 900px;
  margin: 0 auto;

  .editor-wrapper {
    width: 100%;

    .editor-toolbar {
      margin-bottom: 8px;
    }
  }

  .tip {
    font-size: 12px;
    color: #909399;
    margin-top: 4px;
  }

  .content-preview {
    margin-top: 16px;
    padding: 16px;
    background: #f5f7fa;
    border-radius: 4px;

    h4 {
      margin: 0 0 12px 0;
    }

    :deep(pre) {
      background: white;
      padding: 12px;
      border-radius: 4px;
      overflow-x: auto;
    }
  }
}
</style>
