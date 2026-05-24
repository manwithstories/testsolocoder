<template>
  <el-card>
    <template #header>我的手表</template>
    <el-table :data="list">
      <el-table-column prop="brand" label="品牌" />
      <el-table-column prop="model" label="型号" />
      <el-table-column prop="price" label="价格" />
      <el-table-column prop="status" label="状态" />
      <el-table-column label="操作">
        <template #default="{ row }">
          <el-button size="small" @click="edit(row)">编辑</el-button>
          <el-button size="small" type="danger" @click="remove(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="showEdit" title="编辑手表" width="600px">
      <el-form :model="editForm" label-width="100px">
        <el-form-item label="品牌"><el-input v-model="editForm.brand" /></el-form-item>
        <el-form-item label="型号"><el-input v-model="editForm.model" /></el-form-item>
        <el-form-item label="参考号"><el-input v-model="editForm.reference_no" /></el-form-item>
        <el-form-item label="年份"><el-input-number v-model="editForm.year" :min="1900" :max="2100" /></el-form-item>
        <el-form-item label="机芯"><el-input v-model="editForm.movement" /></el-form-item>
        <el-form-item label="表径(mm)"><el-input-number v-model="editForm.case_size_mm" :precision="1" /></el-form-item>
        <el-form-item label="表壳"><el-input v-model="editForm.case_material" /></el-form-item>
        <el-form-item label="表盘"><el-input v-model="editForm.dial_color" /></el-form-item>
        <el-form-item label="表带"><el-input v-model="editForm.bracelet" /></el-form-item>
        <el-form-item label="成色"><el-input v-model="editForm.condition" /></el-form-item>
        <el-form-item label="描述"><el-input v-model="editForm.description" type="textarea" :rows="3" /></el-form-item>
        <el-form-item label="价格"><el-input-number v-model="editForm.price" :min="0" :precision="2" /></el-form-item>
        <el-form-item label="状态">
          <el-select v-model="editForm.status">
            <el-option label="在售" value="on_sale" />
            <el-option label="已预留" value="reserved" />
            <el-option label="已售出" value="sold" />
            <el-option label="下架" value="offline" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showEdit = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="submitEdit">保存</el-button>
      </template>
    </el-dialog>
  </el-card>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '@/utils/request'
import type { Watch } from '@/types'

const list = ref<Watch[]>([])
const showEdit = ref(false)
const saving = ref(false)
const editingId = ref<number | null>(null)
const editForm = reactive({
  brand: '', model: '', reference_no: '', year: 2024, movement: '', case_size_mm: 40,
  case_material: '', dial_color: '', bracelet: '', condition: '', description: '', price: 0, status: 'on_sale'
})

onMounted(async () => {
  const res: any = await request.get('/watches')
  list.value = res.list || []
})

function edit(row: Watch) {
  editingId.value = row.id
  editForm.brand = row.brand
  editForm.model = row.model
  editForm.reference_no = row.reference_no || ''
  editForm.year = row.year
  editForm.movement = row.movement || ''
  editForm.case_size_mm = row.case_size_mm || 40
  editForm.case_material = row.case_material || ''
  editForm.dial_color = row.dial_color || ''
  editForm.bracelet = row.bracelet || ''
  editForm.condition = row.condition || ''
  editForm.description = row.description || ''
  editForm.price = row.price
  editForm.status = row.status
  showEdit.value = true
}

async function submitEdit() {
  if (!editingId.value) return
  saving.value = true
  try {
    await request.put(`/watches/${editingId.value}`, editForm)
    ElMessage.success('已保存')
    showEdit.value = false
    const res: any = await request.get('/watches')
    list.value = res.list || []
  } finally {
    saving.value = false
  }
}

async function remove(row: Watch) {
  await ElMessageBox.confirm('确定删除?', '提示')
  await request.delete(`/watches/${row.id}`)
  ElMessage.success('已删除')
  list.value = list.value.filter((x) => x.id !== row.id)
}
</script>
