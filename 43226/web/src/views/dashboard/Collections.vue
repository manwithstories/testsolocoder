<template>
  <div class="dashboard-collections">
    <div class="card-shadow p-20 mb-20">
      <div class="flex-between mb-20">
        <h2 class="page-title">藏品管理</h2>
        <div class="action-buttons">
          <el-button type="success" @click="showImport = true">
            <el-icon><Upload /></el-icon> 批量导入
          </el-button>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon> 新增藏品
          </el-button>
        </div>
      </div>

      <el-form :inline="true" :model="query" class="mb-20">
        <el-form-item>
          <el-input
            v-model="query.keyword"
            placeholder="搜索名称、编号"
            clearable
            @keyup.enter="fetchList"
          />
        </el-form-item>
        <el-form-item>
          <el-select v-model="query.category_id" clearable placeholder="分类" @change="fetchList">
            <el-option v-for="cat in categories" :key="cat.id" :label="cat.name" :value="cat.id" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-select v-model="query.status" clearable placeholder="状态" @change="fetchList">
            <el-option label="展出中" value="active" />
            <el-option label="未展出" value="inactive" />
            <el-option label="修复中" value="repair" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="fetchList">查询</el-button>
          <el-button @click="resetQuery">重置</el-button>
        </el-form-item>
      </el-form>

      <el-table :data="list" v-loading="loading" border>
        <el-table-column type="index" width="60" label="#" />
        <el-table-column label="图片" width="80">
          <template #default="{ row }">
            <img :src="row.image_url || '/placeholder.svg'" style="width: 60px; height: 60px; object-fit: cover; border-radius: 4px;" />
          </template>
        </el-table-column>
        <el-table-column prop="name" label="名称" min-width="150" show-overflow-tooltip />
        <el-table-column prop="code" label="编号" width="120" />
        <el-table-column prop="category.name" label="分类" width="100" />
        <el-table-column prop="era" label="年代" width="100" />
        <el-table-column prop="material" label="材质" width="100" />
        <el-table-column prop="view_count" label="浏览量" width="80" />
        <el-table-column prop="status" label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="row.status === 'active' ? 'success' : row.status === 'repair' ? 'warning' : 'info'" size="small">
              {{ row.status === 'active' ? '展出中' : row.status === 'repair' ? '修复中' : '未展出' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" size="small" link @click="handleEdit(row)">编辑</el-button>
            <el-button type="danger" size="small" link @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination mt-20">
        <el-pagination
          v-model:current-page="query.page"
          v-model:page-size="query.page_size"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          :page-sizes="[10, 20, 50, 100]"
          @size-change="fetchList"
          @current-change="fetchList"
        />
      </div>
    </div>

    <el-dialog v-model="showDialog" :title="isEdit ? '编辑藏品' : '新增藏品'" width="700px">
      <el-form :model="form" :rules="rules" ref="formRef" label-width="100px">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="藏品名称" prop="name">
              <el-input v-model="form.name" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="藏品编号" prop="code">
              <el-input v-model="form.code" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="分类" prop="category_id">
              <el-select v-model="form.category_id" style="width: 100%">
                <el-option v-for="cat in categories" :key="cat.id" :label="cat.name" :value="cat.id" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="状态" prop="status">
              <el-select v-model="form.status" style="width: 100%">
                <el-option label="展出中" value="active" />
                <el-option label="未展出" value="inactive" />
                <el-option label="修复中" value="repair" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="年代">
              <el-input v-model="form.era" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="材质">
              <el-input v-model="form.material" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="尺寸">
              <el-input v-model="form.size" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="来源">
              <el-input v-model="form.source" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="保存状态">
              <el-input v-model="form.condition" />
            </el-form-item>
          </el-col>
          <el-col :span="24">
            <el-form-item label="标签">
              <el-select v-model="selectedTags" multiple filterable allow-create style="width: 100%">
                <el-option v-for="tag in tags" :key="tag.id" :label="tag.name" :value="tag.name" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="24">
            <el-form-item label="图片">
              <el-upload
                action="/api/v1/uploads/image"
                :headers="{ Authorization: `Bearer ${userStore.token}` }"
                :show-file-list="false"
                :on-success="handleUploadSuccess"
              >
                <img v-if="form.image_url" :src="form.image_url" style="width: 120px; height: 120px; object-fit: cover; border-radius: 4px;" />
                <el-icon v-else size="40"><Plus /></el-icon>
              </el-upload>
            </el-form-item>
          </el-col>
          <el-col :span="24">
            <el-form-item label="描述">
              <el-input v-model="form.description" type="textarea" :rows="4" />
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      <template #footer>
        <el-button @click="showDialog = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="showImport" title="批量导入藏品" width="500px">
      <el-alert type="info" show-icon class="mb-20">
        请上传JSON格式的藏品数据数组，包含name、code、category_id等必填字段
      </el-alert>
      <el-upload
        drag
        action=""
        :auto-upload="false"
        :on-change="handleImportFile"
        accept=".json"
      >
        <el-icon class="el-icon--upload"><UploadFilled /></el-icon>
        <div class="el-upload__text">将JSON文件拖到此处，或<em>点击上传</em></div>
      </el-upload>
      <template #footer>
        <el-button @click="showImport = false">取消</el-button>
        <el-button type="primary" :loading="importing" @click="handleImport">导入</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules, type UploadFile } from 'element-plus'
import * as collectionApi from '@/api/collection'
import type { Collection, CollectionQuery, CollectionCategory, CollectionTag } from '@/types'
import { useUserStore } from '@/stores/user'

const userStore = useUserStore()
const loading = ref(false)
const submitting = ref(false)
const importing = ref(false)
const list = ref<Collection[]>([])
const total = ref(0)
const categories = ref<CollectionCategory[]>([])
const tags = ref<CollectionTag[]>([])
const showDialog = ref(false)
const showImport = ref(false)
const isEdit = ref(false)
const formRef = ref<FormInstance>()
const importData = ref<any[]>([])
const selectedTags = ref<string[]>([])

const query = reactive<CollectionQuery>({
  page: 1,
  page_size: 10,
  keyword: '',
  category_id: undefined,
  status: '',
  sort_by: 'created_at',
  sort_order: 'desc'
})

const form = reactive({
  id: 0,
  name: '',
  code: '',
  category_id: 0,
  era: '',
  material: '',
  size: '',
  source: '',
  condition: '',
  description: '',
  image_url: '',
  status: 'active'
})

const rules: FormRules = {
  name: [{ required: true, message: '请输入藏品名称', trigger: 'blur' }],
  code: [{ required: true, message: '请输入藏品编号', trigger: 'blur' }],
  category_id: [{ required: true, message: '请选择分类', trigger: 'change' }],
  status: [{ required: true, message: '请选择状态', trigger: 'change' }]
}

const resetQuery = () => {
  query.keyword = ''
  query.category_id = undefined
  query.status = ''
  query.page = 1
  fetchList()
}

const fetchList = async () => {
  try {
    loading.value = true
    const res = await collectionApi.listCollections(query)
    list.value = res.data.list
    total.value = res.data.total
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

const fetchCategories = async () => {
  try {
    const res = await collectionApi.listCategories()
    categories.value = res.data
  } catch (e) {
    console.error(e)
  }
}

const fetchTags = async () => {
  try {
    const res = await collectionApi.listTags()
    tags.value = res.data
  } catch (e) {
    console.error(e)
  }
}

const handleAdd = () => {
  isEdit.value = false
  Object.assign(form, {
    id: 0,
    name: '',
    code: '',
    category_id: categories.value[0]?.id || 0,
    era: '',
    material: '',
    size: '',
    source: '',
    condition: '',
    description: '',
    image_url: '',
    status: 'active'
  })
  selectedTags.value = []
  showDialog.value = true
}

const handleEdit = (row: Collection) => {
  isEdit.value = true
  Object.assign(form, {
    id: row.id,
    name: row.name,
    code: row.code,
    category_id: row.category_id,
    era: row.era,
    material: row.material,
    size: row.size,
    source: row.source,
    condition: row.condition,
    description: row.description,
    image_url: row.image_url,
    status: row.status
  })
  selectedTags.value = row.tags ? row.tags.split(',').filter(t => t) : []
  showDialog.value = true
}

const handleDelete = (row: Collection) => {
  ElMessageBox.confirm('确定要删除该藏品吗？', '提示', {
    type: 'warning'
  }).then(async () => {
    await collectionApi.deleteCollection(row.id)
    ElMessage.success('删除成功')
    fetchList()
  }).catch(() => {})
}

const handleUploadSuccess = (res: any) => {
  form.image_url = res.data.url
}

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate()
  try {
    submitting.value = true
    const data = { ...form, tags: selectedTags.value.join(',') }
    if (isEdit.value) {
      await collectionApi.updateCollection(form.id, data)
      ElMessage.success('更新成功')
    } else {
      await collectionApi.createCollection(data)
      ElMessage.success('创建成功')
    }
    showDialog.value = false
    fetchList()
  } catch (e) {
    console.error(e)
  } finally {
    submitting.value = false
  }
}

const handleImportFile = (file: UploadFile) => {
  const reader = new FileReader()
  reader.onload = (e) => {
    try {
      importData.value = JSON.parse(e.target?.result as string)
      ElMessage.success(`读取到 ${importData.value.length} 条数据`)
    } catch (e) {
      ElMessage.error('JSON格式错误')
    }
  }
  reader.readAsText(file.raw as File)
}

const handleImport = async () => {
  if (importData.value.length === 0) {
    ElMessage.warning('请先选择导入文件')
    return
  }
  try {
    importing.value = true
    await collectionApi.batchImportCollections(importData.value)
    ElMessage.success('导入成功')
    showImport.value = false
    fetchList()
  } catch (e) {
    console.error(e)
  } finally {
    importing.value = false
  }
}

onMounted(() => {
  fetchCategories()
  fetchTags()
  fetchList()
})
</script>

<style scoped lang="scss">
.dashboard-collections {
  .page-title {
    margin: 0;
    font-size: 20px;
  }
}
</style>
