<template>
  <div class="page-container">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>服务项目</span>
          <el-button type="primary" :icon="Plus" @click="addVisible = true">添加服务</el-button>
        </div>
      </template>

      <el-table :data="services" stripe v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="服务名称" />
        <el-table-column prop="category" label="分类" width="100" />
        <el-table-column prop="price" label="价格" width="100">
          <template #default="{ row }">¥{{ row.price?.toFixed(2) }}</template>
        </el-table-column>
        <el-table-column prop="duration" label="时长(分钟)" width="120" />
        <el-table-column label="套餐" width="100">
          <template #default="{ row }">
            <el-tag v-if="row.is_package" type="warning">
              套餐({{ row.package_count }}次)
            </el-tag>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column label="动态定价" width="100">
          <template #default="{ row }">
            <el-tag v-if="row.dynamic_pricing" type="success">已开启</el-tag>
            <span v-else>未开启</span>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="160" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleEdit(row)">编辑</el-button>
            <el-button type="danger" link @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination">
        <el-pagination
          v-model:current-page="page"
          v-model:page-size="pageSize"
          :total="total"
          :page-sizes="[10, 20, 50]"
          layout="total, sizes, prev, pager, next, jumper"
          @current-change="fetchList"
          @size-change="fetchList"
        />
      </div>
    </el-card>

    <el-dialog v-model="addVisible" :title="editing ? '编辑服务' : '添加服务'" width="600px">
      <el-form :model="form" label-width="100px">
        <el-form-item label="服务名称" required>
          <el-input v-model="form.name" />
        </el-form-item>
        <el-form-item label="分类" required>
          <el-select v-model="form.category">
            <el-option label="剪发" value="剪发" />
            <el-option label="染发" value="染发" />
            <el-option label="烫发" value="烫发" />
            <el-option label="美容" value="美容" />
            <el-option label="SPA" value="SPA" />
            <el-option label="其他" value="其他" />
          </el-select>
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="form.description" type="textarea" />
        </el-form-item>
        <el-form-item label="价格" required>
          <el-input-number v-model="form.price" :min="0" :precision="2" />
        </el-form-item>
        <el-form-item label="时长(分钟)" required>
          <el-input-number v-model="form.duration" :min="15" :step="15" />
        </el-form-item>
        <el-form-item label="所需技能">
          <el-input v-model="form.required_skill" />
        </el-form-item>
        <el-form-item label="是否套餐">
          <el-switch v-model="form.is_package" />
        </el-form-item>
        <el-form-item v-if="form.is_package" label="套餐次数">
          <el-input-number v-model="form.package_count" :min="1" />
        </el-form-item>
        <el-form-item label="动态定价">
          <el-switch v-model="form.dynamic_pricing" />
        </el-form-item>
        <el-form-item v-if="form.dynamic_pricing" label="周末价格">
          <el-input-number v-model="form.weekend_price" :min="0" :precision="2" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="addVisible = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="handleSave">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { getServices, createService, updateService, deleteService } from '@/api/service'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import type { Service } from '@/types'

const loading = ref(false)
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)
const services = ref<Service[]>([])
const addVisible = ref(false)
const editing = ref(false)
const editingId = ref<number | null>(null)
const saving = ref(false)

const form = reactive({
  name: '',
  category: '',
  description: '',
  price: 0,
  duration: 60,
  required_skill: '',
  is_package: false,
  package_count: 0,
  dynamic_pricing: false,
  weekend_price: 0
})

const fetchList = async () => {
  loading.value = true
  try {
    const res = await getServices({ page: page.value, page_size: pageSize.value })
    services.value = res.data.list
    total.value = res.data.total
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

const handleEdit = (row: Service) => {
  editing.value = true
  editingId.value = row.id
  Object.assign(form, {
    name: row.name,
    category: row.category,
    description: row.description,
    price: row.price,
    duration: row.duration,
    required_skill: row.required_skill,
    is_package: row.is_package,
    package_count: row.package_count,
    dynamic_pricing: row.dynamic_pricing,
    weekend_price: row.weekend_price
  })
  addVisible.value = true
}

const handleDelete = async (row: Service) => {
  try {
    await ElMessageBox.confirm(`确定删除服务"${row.name}"吗？`, '提示', {
      type: 'warning'
    })
    await deleteService(row.id)
    ElMessage.success('删除成功')
    fetchList()
  } catch (e) {
    if (e !== 'cancel') console.error(e)
  }
}

const handleSave = async () => {
  if (!form.name || !form.category) {
    ElMessage.warning('请填写必填项')
    return
  }

  saving.value = true
  try {
    if (editing.value && editingId.value) {
      await updateService(editingId.value, form)
      ElMessage.success('更新成功')
    } else {
      await createService(form)
      ElMessage.success('添加成功')
    }
    addVisible.value = false
    fetchList()
  } catch (e: any) {
    ElMessage.error(e.message || '保存失败')
  } finally {
    saving.value = false
  }
}

onMounted(fetchList)
</script>

<style scoped lang="scss">
.page-container {
  .card-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
  }

  .pagination {
    margin-top: 20px;
    display: flex;
    justify-content: flex-end;
  }
}
</style>
