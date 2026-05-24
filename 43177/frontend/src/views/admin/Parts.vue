<template>
  <div class="admin-parts-page">
    <div class="page-header">
      <h2 class="page-title">配件管理</h2>
      <el-button type="primary" @click="showCreateDialog = true">
        新增配件
      </el-button>
    </div>

    <el-card>
      <el-table :data="parts" style="width: 100%">
        <el-table-column prop="name" label="名称" width="150" />
        <el-table-column prop="code" label="编码" width="150" />
        <el-table-column prop="category" label="分类" width="120" />
        <el-table-column prop="price" label="价格" width="100">
          <template #default="{ row }">¥{{ row.price }}</template>
        </el-table-column>
        <el-table-column prop="stock" label="库存" width="100">
          <template #default="{ row }">
            <span :class="{ 'low-stock': row.stock <= row.min_stock }">{{ row.stock }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="min_stock" label="最低库存" width="100" />
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status ? 'success' : 'danger'">
              {{ row.status ? '正常' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150">
          <template #default="{ row }">
            <el-button size="small" @click="editPart(row)">编辑</el-button>
            <el-button size="small" type="danger" @click="deletePart(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div v-if="lowStockCount > 0" class="warning">
        <el-icon><Warning /></el-icon>
        当前有 {{ lowStockCount }} 种配件库存不足
      </div>
    </el-card>

    <el-dialog v-model="showCreateDialog" :title="editingPart ? '编辑配件' : '新增配件'" width="500px">
      <el-form :model="partForm" label-width="80px">
        <el-form-item label="名称">
          <el-input v-model="partForm.name" />
        </el-form-item>
        <el-form-item label="编码">
          <el-input v-model="partForm.code" />
        </el-form-item>
        <el-form-item label="分类">
          <el-input v-model="partForm.category" />
        </el-form-item>
        <el-form-item label="价格">
          <el-input-number v-model="partForm.price" :min="0" :precision="2" />
        </el-form-item>
        <el-form-item label="库存">
          <el-input-number v-model="partForm.stock" :min="0" />
        </el-form-item>
        <el-form-item label="最低库存">
          <el-input-number v-model="partForm.min_stock" :min="0" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showCreateDialog = false">取消</el-button>
        <el-button type="primary" @click="savePart" :loading="saving">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Warning } from '@element-plus/icons-vue'
import { partApi } from '@/api/part'
import type { Part } from '@/types'

const parts = ref<Part[]>([])
const lowStockCount = ref(0)
const showCreateDialog = ref(false)
const editingPart = ref<Part | null>(null)
const saving = ref(false)

const partForm = reactive({
  name: '',
  code: '',
  category: '',
  price: 0,
  stock: 0,
  min_stock: 10
})

onMounted(() => {
  loadParts()
})

async function loadParts() {
  try {
    const res = await partApi.getParts({ page: 1, page_size: 100 })
    parts.value = res.data?.list || []
    lowStockCount.value = (res.data as any)?.low_stock_count || 0
  } catch (error) {
    console.error('Failed to load parts:', error)
  }
}

function editPart(part: Part) {
  editingPart.value = part
  partForm.name = part.name
  partForm.code = part.code
  partForm.category = part.category
  partForm.price = part.price
  partForm.stock = part.stock
  partForm.min_stock = part.min_stock
  showCreateDialog.value = true
}

async function savePart() {
  if (!partForm.name || !partForm.code) {
    ElMessage.error('请填写名称和编码')
    return
  }

  saving.value = true
  try {
    if (editingPart.value) {
      await partApi.updatePart(editingPart.value.id, partForm)
    } else {
      await partApi.createPart(partForm)
    }
    ElMessage.success('保存成功')
    showCreateDialog.value = false
    resetForm()
    loadParts()
  } catch (error) {
    console.error('Failed to save part:', error)
  } finally {
    saving.value = false
  }
}

function resetForm() {
  editingPart.value = null
  partForm.name = ''
  partForm.code = ''
  partForm.category = ''
  partForm.price = 0
  partForm.stock = 0
  partForm.min_stock = 10
}

async function deletePart(part: Part) {
  try {
    await ElMessageBox.confirm(`确定要删除配件"${part.name}"吗？`, '确认', { type: 'warning' })
    await partApi.deletePart(part.id)
    ElMessage.success('删除成功')
    loadParts()
  } catch (error: any) {
    if (error !== 'cancel') {
      console.error('Failed to delete part:', error)
    }
  }
}
</script>

<style scoped>
.admin-parts-page {
  padding: 0;
}

.low-stock {
  color: #f56c6c;
  font-weight: 600;
}

.warning {
  margin-top: 20px;
  padding: 15px;
  background-color: #fdf6ec;
  color: #e6a23c;
  border-radius: 4px;
  display: flex;
  align-items: center;
  gap: 10px;
}
</style>
