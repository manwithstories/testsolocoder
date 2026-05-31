<template>
  <div class="dashboard-exhibitions">
    <div class="card-shadow p-20 mb-20">
      <div class="flex-between mb-20">
        <h2 class="page-title">展览管理</h2>
        <el-button type="primary" @click="handleAdd">
          <el-icon><Plus /></el-icon> 新增展览
        </el-button>
      </div>

      <el-form :inline="true" :model="query" class="mb-20">
        <el-form-item>
          <el-input v-model="query.keyword" placeholder="搜索展览名称" clearable @keyup.enter="fetchList" />
        </el-form-item>
        <el-form-item>
          <el-select v-model="query.status" clearable placeholder="状态" @change="fetchList">
            <el-option label="草稿" value="draft" />
            <el-option label="已发布" value="published" />
            <el-option label="已结束" value="closed" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="fetchList">查询</el-button>
          <el-button @click="resetQuery">重置</el-button>
        </el-form-item>
      </el-form>

      <el-table :data="list" v-loading="loading" border>
        <el-table-column label="封面" width="120">
          <template #default="{ row }">
            <img :src="row.image_url || '/placeholder.svg'" style="width: 100px; height: 60px; object-fit: cover; border-radius: 4px;" />
          </template>
        </el-table-column>
        <el-table-column prop="title" label="展览名称" min-width="180" show-overflow-tooltip />
        <el-table-column label="时间" width="220">
          <template #default="{ row }">
            {{ formatDate(row.start_date) }} - {{ formatDate(row.end_date) }}
          </template>
        </el-table-column>
        <el-table-column prop="location" label="地点" width="120" />
        <el-table-column prop="hall_number" label="展厅" width="80" />
        <el-table-column label="票价" width="80">
          <template #default="{ row }">
            {{ row.ticket_price > 0 ? '¥' + row.ticket_price : '免费' }}
          </template>
        </el-table-column>
        <el-table-column prop="view_count" label="浏览量" width="80" />
        <el-table-column label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="row.status === 'published' ? 'success' : row.status === 'closed' ? 'info' : 'warning'" size="small">
              {{ row.status === 'published' ? '已发布' : row.status === 'closed' ? '已结束' : '草稿' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="220" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" size="small" link @click="$router.push(`/dashboard/exhibitions/${row.id}/time-slots`)">
              时段
            </el-button>
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

    <el-dialog v-model="showDialog" :title="isEdit ? '编辑展览' : '新增展览'" width="800px">
      <el-form :model="form" :rules="rules" ref="formRef" label-width="100px">
        <el-form-item label="展览标题" prop="title">
          <el-input v-model="form.title" />
        </el-form-item>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="开始日期" prop="start_date">
              <el-date-picker v-model="form.start_date" type="date" value-format="YYYY-MM-DD" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="结束日期" prop="end_date">
              <el-date-picker v-model="form.end_date" type="date" value-format="YYYY-MM-DD" style="width: 100%" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="8">
            <el-form-item label="地点" prop="location">
              <el-input v-model="form.location" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="展厅">
              <el-input v-model="form.hall_number" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="票价" prop="ticket_price">
              <el-input-number v-model="form.ticket_price" :min="0" :precision="2" style="width: 100%" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="最大人数">
              <el-input-number v-model="form.max_visitors" :min="1" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="状态" prop="status">
              <el-select v-model="form.status" style="width: 100%">
                <el-option label="草稿" value="draft" />
                <el-option label="已发布" value="published" />
                <el-option label="已结束" value="closed" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="虚拟展厅">
          <el-switch v-model="form.is_virtual" />
        </el-form-item>
        <el-form-item v-if="form.is_virtual" label="虚拟展厅链接">
          <el-input v-model="form.virtual_url" />
        </el-form-item>
        <el-form-item label="选择藏品">
          <el-select v-model="form.collection_ids" multiple filterable style="width: 100%">
            <el-option v-for="c in allCollections" :key="c.id" :label="c.name" :value="c.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="封面图片">
          <el-upload
            action="/api/v1/uploads/image"
            :headers="{ Authorization: `Bearer ${userStore.token}` }"
            :show-file-list="false"
            :on-success="handleUploadSuccess"
          >
            <img v-if="form.image_url" :src="form.image_url" style="width: 150px; height: 100px; object-fit: cover; border-radius: 4px;" />
            <el-icon v-else size="40"><Plus /></el-icon>
          </el-upload>
        </el-form-item>
        <el-form-item label="展览描述" prop="description">
          <el-input v-model="form.description" type="textarea" :rows="4" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showDialog = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import * as exhibitionApi from '@/api/exhibition'
import * as collectionApi from '@/api/collection'
import type { Exhibition, ExhibitionQuery, Collection } from '@/types'
import { useUserStore } from '@/stores/user'
import dayjs from 'dayjs'

const userStore = useUserStore()
const loading = ref(false)
const submitting = ref(false)
const list = ref<Exhibition[]>([])
const total = ref(0)
const allCollections = ref<Collection[]>([])
const showDialog = ref(false)
const isEdit = ref(false)
const editId = ref(0)
const formRef = ref<FormInstance>()

const query = reactive<ExhibitionQuery>({
  page: 1,
  page_size: 10,
  keyword: '',
  status: ''
})

const form = reactive({
  title: '',
  description: '',
  start_date: '',
  end_date: '',
  location: '',
  hall_number: '',
  ticket_price: 0,
  max_visitors: 50,
  image_url: '',
  status: 'draft',
  is_virtual: false,
  virtual_url: '',
  collection_ids: [] as number[]
})

const rules: FormRules = {
  title: [{ required: true, message: '请输入展览标题', trigger: 'blur' }],
  start_date: [{ required: true, message: '请选择开始日期', trigger: 'change' }],
  end_date: [{ required: true, message: '请选择结束日期', trigger: 'change' }],
  location: [{ required: true, message: '请输入地点', trigger: 'blur' }],
  ticket_price: [{ required: true, message: '请输入票价', trigger: 'blur' }],
  status: [{ required: true, message: '请选择状态', trigger: 'change' }],
  description: [{ required: true, message: '请输入展览描述', trigger: 'blur' }]
}

const formatDate = (date: string) => {
  return dayjs(date).format('YYYY-MM-DD')
}

const resetQuery = () => {
  query.keyword = ''
  query.status = ''
  query.page = 1
  fetchList()
}

const fetchList = async () => {
  try {
    loading.value = true
    const res = await exhibitionApi.listExhibitions(query)
    list.value = res.data.list
    total.value = res.data.total
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

const fetchCollections = async () => {
  try {
    const res = await collectionApi.listCollections({ page: 1, page_size: 1000 })
    allCollections.value = res.data.list
  } catch (e) {
    console.error(e)
  }
}

const handleAdd = () => {
  isEdit.value = false
  editId.value = 0
  Object.assign(form, {
    title: '',
    description: '',
    start_date: '',
    end_date: '',
    location: '',
    hall_number: '',
    ticket_price: 0,
    max_visitors: 50,
    image_url: '',
    status: 'draft',
    is_virtual: false,
    virtual_url: '',
    collection_ids: []
  })
  showDialog.value = true
}

const handleEdit = async (row: Exhibition) => {
  isEdit.value = true
  editId.value = row.id
  const detail = await exhibitionApi.getExhibition(row.id)
  const collections = await exhibitionApi.getExhibitionCollections(row.id)
  Object.assign(form, {
    title: detail.data.title,
    description: detail.data.description,
    start_date: dayjs(detail.data.start_date).format('YYYY-MM-DD'),
    end_date: dayjs(detail.data.end_date).format('YYYY-MM-DD'),
    location: detail.data.location,
    hall_number: detail.data.hall_number,
    ticket_price: detail.data.ticket_price,
    max_visitors: detail.data.max_visitors,
    image_url: detail.data.image_url,
    status: detail.data.status,
    is_virtual: detail.data.is_virtual,
    virtual_url: detail.data.virtual_url,
    collection_ids: collections.data.map(c => c.id)
  })
  showDialog.value = true
}

const handleDelete = (row: Exhibition) => {
  ElMessageBox.confirm('确定要删除该展览吗？', '提示', {
    type: 'warning'
  }).then(async () => {
    await exhibitionApi.deleteExhibition(row.id)
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
    if (isEdit.value) {
      await exhibitionApi.updateExhibition(editId.value, form)
      ElMessage.success('更新成功')
    } else {
      await exhibitionApi.createExhibition(form)
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

onMounted(() => {
  fetchList()
  fetchCollections()
})
</script>
