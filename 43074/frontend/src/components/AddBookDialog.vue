<template>
  <el-dialog
    v-model="visible"
    title="添加图书"
    width="600px"
    :close-on-click-modal="false"
    @close="resetForm"
  >
    <el-form ref="formRef" :model="form" :rules="rules" label-width="100px">
      <el-form-item label="ISBN" prop="isbn">
        <el-input v-model="form.isbn" placeholder="输入ISBN自动获取书籍信息" clearable>
          <template #append>
            <el-button @click="fetchBookInfo" :loading="fetching">获取信息</el-button>
          </template>
        </el-input>
        <div class="isbn-tip">支持ISBN-10和ISBN-13格式</div>
      </el-form-item>
      <el-form-item label="书名" prop="title">
        <el-input v-model="form.title" placeholder="请输入书名" />
      </el-form-item>
      <el-form-item label="作者">
        <el-input v-model="form.author" placeholder="请输入作者" />
      </el-form-item>
      <el-form-item label="出版社">
        <el-input v-model="form.publisher" placeholder="请输入出版社" />
      </el-form-item>
      <el-form-item label="总页数">
        <el-input-number v-model="form.total_pages" :min="1" :max="10000" />
      </el-form-item>
      <el-form-item label="简介">
        <el-input
          v-model="form.summary"
          type="textarea"
          :rows="3"
          placeholder="请输入书籍简介"
        />
      </el-form-item>
      <el-form-item label="标签">
        <el-select v-model="form.tag_ids" multiple filterable placeholder="选择标签">
          <el-option
            v-for="tag in tags"
            :key="tag.id"
            :label="tag.name"
            :value="tag.id"
          >
            <span :style="{ color: tag.color }">●</span>
            {{ tag.name }}
          </el-option>
        </el-select>
      </el-form-item>
      <el-form-item label="分类">
        <el-cascader
          v-model="form.category_ids"
          :options="categoryOptions"
          :props="{ multiple: true, checkStrictly: true, label: 'name', value: 'id', children: 'children' }"
          placeholder="选择分类"
          clearable
        />
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="visible = false">取消</el-button>
      <el-button type="primary" @click="submitForm" :loading="submitting">确定</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, watch, onMounted, computed } from 'vue'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { createBook, fetchBookByISBN, type CreateBookData } from '@/api/book'
import { getTags, getCategories } from '@/api/common'
import type { Tag, Category } from '@/types'

const props = defineProps<{ modelValue: boolean }>()
const emit = defineEmits<{
  (e: 'update:modelValue', value: boolean): void
  (e: 'success'): void
}>()

const visible = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val)
})

const formRef = ref<FormInstance>()
const fetching = ref(false)
const submitting = ref(false)
const tags = ref<Tag[]>([])
const categories = ref<Category[]>([])

const form = reactive<CreateBookData>({
  title: '',
  author: '',
  publisher: '',
  isbn: '',
  summary: '',
  total_pages: 0,
  tag_ids: [],
  category_ids: []
})

const rules: FormRules = {
  title: [{ required: true, message: '请输入书名', trigger: 'blur' }],
  isbn: [{ validator: validateISBN, trigger: 'blur' }]
}

const categoryOptions = computed(() => categories.value)

function validateISBN(rule: any, value: string, callback: any) {
  if (!value) return callback()
  const clean = value.replace(/[-\s]/g, '').toUpperCase()
  if (clean.length === 10 || clean.length === 13) {
    callback()
  } else {
    callback(new Error('ISBN格式不正确'))
  }
}

const fetchBookInfo = async () => {
  if (!form.isbn) {
    ElMessage.warning('请输入ISBN')
    return
  }
  fetching.value = true
  try {
    const info = await fetchBookByISBN(form.isbn)
    if (info.title) form.title = info.title
    if (info.author) form.author = info.author
    if (info.publisher) form.publisher = info.publisher
    if (info.summary) form.summary = info.summary
    if (info.pages) form.total_pages = info.pages
    ElMessage.success('获取书籍信息成功')
  } catch (e) {
    ElMessage.error('未找到该ISBN对应的书籍信息，请手动填写')
  } finally {
    fetching.value = false
  }
}

const loadTags = async () => {
  try {
    tags.value = await getTags()
  } catch (e) {}
}

const loadCategories = async () => {
  try {
    categories.value = await getCategories(false)
  } catch (e) {}
}

const submitForm = async () => {
  if (!formRef.value) return
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return

  submitting.value = true
  try {
    await createBook({
      ...form,
      total_pages: form.total_pages || undefined
    })
    ElMessage.success('添加成功')
    emit('success')
    visible.value = false
  } catch (e) {
  } finally {
    submitting.value = false
  }
}

const resetForm = () => {
  form.title = ''
  form.author = ''
  form.publisher = ''
  form.isbn = ''
  form.summary = ''
  form.total_pages = 0
  form.tag_ids = []
  form.category_ids = []
  formRef.value?.resetFields()
}

onMounted(() => {
  loadTags()
  loadCategories()
})

watch(visible, (val) => {
  if (val) {
    loadTags()
    loadCategories()
  }
})
</script>

<style scoped lang="scss">
.isbn-tip {
  font-size: 12px;
  color: #909399;
  margin-top: 4px;
}
</style>
