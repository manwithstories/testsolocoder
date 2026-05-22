<template>
  <div class="ticket-list">
    <div class="page-header">
      <h2 class="page-title">票务管理</h2>
      <el-button type="primary" @click="dialogVisible = true" v-if="userStore.isAdmin">
        <el-icon><Plus /></el-icon>
        添加票型
      </el-button>
    </div>

    <el-card>
      <div class="search-bar">
        <el-select v-model="search.activityId" placeholder="选择活动" clearable style="width: 240px">
          <el-option v-for="act in activities" :key="act.id" :label="act.title" :value="act.id" />
        </el-select>
        <el-select v-model="search.status" placeholder="票型状态" clearable style="width: 140px">
          <el-option label="在售" value="on_sale" />
          <el-option label="售罄" value="sold_out" />
          <el-option label="下架" value="off_sale" />
        </el-select>
        <el-button type="primary" @click="loadData">搜索</el-button>
      </div>

      <el-table :data="list" style="width: 100%" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="activity.title" label="所属活动" min-width="180" />
        <el-table-column prop="name" label="票型名称" />
        <el-table-column prop="type" label="类型" width="100">
          <template #default="{ row }">
            <el-tag :type="getTypeColor(row.type)">{{ getTypeText(row.type) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="price" label="价格(元)" width="100">
          <template #default="{ row }">¥{{ row.price }}</template>
        </el-table-column>
        <el-table-column prop="stock" label="库存" width="100" />
        <el-table-column prop="soldCount" label="已售" width="100" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">{{ getStatusText(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="180" v-if="userStore.isAdmin">
          <template #default="{ row }">
            <div class="table-actions">
              <el-button size="small" @click="handleEdit(row)">编辑</el-button>
              <el-button size="small" type="danger" @click="handleDelete(row.id)">删除</el-button>
            </div>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.pageSize"
        :total="pagination.total"
        layout="total, prev, pager, next, jumper"
        @current-change="loadData"
        style="margin-top: 20px; justify-content: flex-end"
      />
    </el-card>

    <el-dialog v-model="dialogVisible" :title="editMode ? '编辑票型' : '添加票型'" width="500px">
      <el-form :model="form" :rules="rules" ref="formRef" label-width="100px">
        <el-form-item label="所属活动" prop="activityId">
          <el-select v-model="form.activityId" placeholder="选择活动" style="width: 100%">
            <el-option v-for="act in activities" :key="act.id" :label="act.title" :value="act.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="票型名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入票型名称" />
        </el-form-item>
        <el-form-item label="票型类型" prop="type">
          <el-select v-model="form.type" placeholder="选择票型类型" style="width: 100%">
            <el-option label="普通票" value="normal" />
            <el-option label="VIP票" value="vip" />
            <el-option label="早鸟票" value="early_bird" />
          </el-select>
        </el-form-item>
        <el-form-item label="价格" prop="price">
          <el-input-number v-model="form.price" :min="0" :precision="2" style="width: 100%" />
        </el-form-item>
        <el-form-item label="库存" prop="stock">
          <el-input-number v-model="form.stock" :min="0" style="width: 100%" />
        </el-form-item>
        <el-form-item label="状态" prop="status" v-if="editMode">
          <el-select v-model="form.status" placeholder="选择状态" style="width: 100%">
            <el-option label="在售" value="on_sale" />
            <el-option label="售罄" value="sold_out" />
            <el-option label="下架" value="off_sale" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import { getTicketTypeList, createTicketType, updateTicketType, deleteTicketType } from '@/api/ticket'
import { getActivityList } from '@/api/activity'
import { useUserStore } from '@/store/user'

const userStore = useUserStore()

const loading = ref(false)
const submitting = ref(false)
const list = ref<any[]>([])
const activities = ref<any[]>([])
const dialogVisible = ref(false)
const editMode = ref(false)
const formRef = ref<FormInstance>()

const search = reactive({
  activityId: '',
  status: ''
})

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

const form = reactive({
  id: 0,
  activityId: 0,
  name: '',
  type: 'normal',
  price: 0,
  stock: 0,
  status: 'on_sale'
})

const rules: FormRules = {
  activityId: [{ required: true, message: '请选择活动', trigger: 'change' }],
  name: [{ required: true, message: '请输入票型名称', trigger: 'blur' }],
  type: [{ required: true, message: '请选择票型类型', trigger: 'change' }],
  price: [{ required: true, message: '请输入价格', trigger: 'blur' }],
  stock: [{ required: true, message: '请输入库存', trigger: 'blur' }]
}

const loadActivities = async () => {
  try {
    const res = await getActivityList({ page: 1, pageSize: 100, status: 'published' })
    activities.value = res.list
  } catch (error) {
    console.error(error)
  }
}

const loadData = async () => {
  try {
    loading.value = true
    const res = await getTicketTypeList({
      page: pagination.page,
      pageSize: pagination.pageSize,
      activityId: search.activityId,
      status: search.status
    })
    list.value = res.list
    pagination.total = res.total
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}

const handleEdit = (row: any) => {
  editMode.value = true
  Object.assign(form, row)
  dialogVisible.value = true
}

const handleDelete = async (id: number) => {
  try {
    await ElMessageBox.confirm('确定要删除这个票型吗？', '提示', { type: 'warning' })
    await deleteTicketType(id)
    ElMessage.success('删除成功')
    loadData()
  } catch (error) {
    if (error !== 'cancel') {
      console.error(error)
    }
  }
}

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    try {
      submitting.value = true
      if (editMode.value) {
        await updateTicketType(form.id, form)
        ElMessage.success('更新成功')
      } else {
        await createTicketType(form as any)
        ElMessage.success('创建成功')
      }
      dialogVisible.value = false
      loadData()
    } catch (error) {
      console.error(error)
    } finally {
      submitting.value = false
    }
  })
}

const getTypeColor = (type: string) => {
  const map: Record<string, string> = { normal: '', vip: 'warning', early_bird: 'success' }
  return map[type] || ''
}

const getTypeText = (type: string) => {
  const map: Record<string, string> = { normal: '普通票', vip: 'VIP票', early_bird: '早鸟票' }
  return map[type] || type
}

const getStatusType = (status: string) => {
  const map: Record<string, string> = { on_sale: 'success', sold_out: 'info', off_sale: 'danger' }
  return map[status] || 'info'
}

const getStatusText = (status: string) => {
  const map: Record<string, string> = { on_sale: '在售', sold_out: '售罄', off_sale: '下架' }
  return map[status] || status
}

onMounted(() => {
  loadActivities()
  loadData()
})
</script>
