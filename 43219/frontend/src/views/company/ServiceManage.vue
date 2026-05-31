<template>
  <AppLayout>
    <div class="page">
      <div class="row" style="margin-bottom:16px">
        <h2 style="margin:0">公司服务管理</h2>
        <el-button type="primary" @click="openCreate">新增服务</el-button>
      </div>
      <el-table :data="list" border stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="名称" />
        <el-table-column prop="category" label="分类" width="120" />
        <el-table-column label="价格" width="160">
          <template #default="{ row }">¥{{ row.min_price }} - ¥{{ row.max_price }}</template>
        </el-table-column>
        <el-table-column prop="duration" label="时长(分钟)" width="120" />
        <el-table-column prop="skills" label="技能" />
        <el-table-column label="操作" width="160">
          <template #default="{ row }">
            <el-button size="small" @click="openEdit(row)">编辑</el-button>
            <el-button size="small" type="danger" @click="remove(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <el-dialog v-model="dlgVisible" :title="editing ? '编辑服务' : '新增服务'">
      <el-form :model="form" :rules="rules" ref="formRef" label-width="80px">
        <el-form-item label="名称" prop="name"><el-input v-model="form.name" /></el-form-item>
        <el-form-item label="分类" prop="category"><el-input v-model="form.category" /></el-form-item>
        <el-form-item label="描述"><el-input v-model="form.desc" type="textarea" :rows="3" /></el-form-item>
        <el-form-item label="最低价格" prop="min_price">
          <el-input-number v-model="form.min_price" :min="0" :precision="2" />
        </el-form-item>
        <el-form-item label="最高价格" prop="max_price">
          <el-input-number v-model="form.max_price" :min="0" :precision="2" />
        </el-form-item>
        <el-form-item label="时长(分钟)" prop="duration">
          <el-input-number v-model="form.duration" :min="15" :step="15" />
        </el-form-item>
        <el-form-item label="技能"><el-input v-model="form.skills" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dlgVisible=false">取消</el-button>
        <el-button type="primary" @click="submit">确定</el-button>
      </template>
    </el-dialog>
  </AppLayout>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import AppLayout from '../../components/AppLayout.vue'
import { myServices, createService, updateService, deleteService, type ServiceItem } from '../../api/service'

const list = ref<ServiceItem[]>([])
const dlgVisible = ref(false)
const editing = ref<ServiceItem | null>(null)
const formRef = ref<FormInstance>()
const form = reactive({ name: '', category: '', desc: '', min_price: 0, max_price: 0, duration: 60, skills: '' })
const rules: FormRules = {
  name: [{ required: true, message: '必填', trigger: 'blur' }],
  category: [{ required: true, message: '必填', trigger: 'blur' }],
  min_price: [{ required: true, type: 'number', message: '必填', trigger: 'change' }],
  max_price: [{ required: true, type: 'number', message: '必填', trigger: 'change' }],
  duration: [{ required: true, type: 'number', message: '必填', trigger: 'change' }],
}

async function load() {
  const res = await myServices()
  list.value = (res.data as any).data || []
}

function openCreate() {
  editing.value = null
  Object.assign(form, { name: '', category: '', desc: '', min_price: 0, max_price: 0, duration: 60, skills: '' })
  dlgVisible.value = true
}
function openEdit(row: ServiceItem) {
  editing.value = row
  Object.assign(form, {
    name: row.name, category: row.category, desc: row.desc,
    min_price: row.min_price, max_price: row.max_price, duration: row.duration, skills: row.skills,
  })
  dlgVisible.value = true
}

async function submit() {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return
  if (editing.value) {
    await updateService(editing.value.id, form)
    ElMessage.success('已更新')
  } else {
    await createService(form)
    ElMessage.success('已创建')
  }
  dlgVisible.value = false
  load()
}

async function remove(row: ServiceItem) {
  await ElMessageBox.confirm(`确认删除 "${row.name}"?`, '提示', { type: 'warning' })
  await deleteService(row.id)
  ElMessage.success('已删除')
  load()
}

onMounted(load)
</script>
