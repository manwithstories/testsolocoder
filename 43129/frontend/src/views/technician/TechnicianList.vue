<template>
  <div class="page-container">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>技师管理</span>
          <el-button type="primary" :icon="Plus" @click="addVisible = true">添加技师</el-button>
        </div>
      </template>

      <el-table :data="technicians" stripe v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column label="头像" width="80">
          <template #default="{ row }">
            <el-avatar :size="40" :src="row.avatar">
              {{ row.name?.[0] || 'T' }}
            </el-avatar>
          </template>
        </el-table-column>
        <el-table-column prop="name" label="姓名" />
        <el-table-column prop="title" label="职称" />
        <el-table-column prop="specialties" label="擅长项目" show-overflow-tooltip />
        <el-table-column label="评分" width="120">
          <template #default="{ row }">
            <el-rate :model-value="row.rating" disabled size="small" />
            <span style="margin-left: 8px">{{ row.rating }}</span>
          </template>
        </el-table-column>
        <el-table-column label="工作时间">
          <template #default="{ row }">
            {{ row.work_start_time }} - {{ row.work_end_time }}
          </template>
        </el-table-column>
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'">
              {{ row.status === 1 ? '在职' : '离职' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="160" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="$router.push(`/technicians/${row.id}`)">详情</el-button>
            <el-button type="warning" link @click="handleEdit(row)">编辑</el-button>
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

    <el-dialog v-model="addVisible" :title="editing ? '编辑技师' : '添加技师'" width="600px">
      <el-form :model="form" label-width="100px">
        <el-form-item label="姓名" required>
          <el-input v-model="form.name" />
        </el-form-item>
        <el-form-item label="职称">
          <el-input v-model="form.title" />
        </el-form-item>
        <el-form-item label="擅长项目">
          <el-input v-model="form.specialties" type="textarea" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="form.description" type="textarea" />
        </el-form-item>
        <el-form-item label="上班时间">
          <el-time-picker
            v-model="form.work_start_time"
            format="HH:mm"
            value-format="HH:mm"
            placeholder="上班时间"
          />
        </el-form-item>
        <el-form-item label="下班时间">
          <el-time-picker
            v-model="form.work_end_time"
            format="HH:mm"
            value-format="HH:mm"
            placeholder="下班时间"
          />
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
import { getTechnicians, createTechnician, updateTechnician } from '@/api/technician'
import { ElMessage } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import type { Technician } from '@/types'

const loading = ref(false)
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)
const technicians = ref<Technician[]>([])
const addVisible = ref(false)
const editing = ref(false)
const editingId = ref<number | null>(null)
const saving = ref(false)

const form = reactive({
  name: '',
  title: '',
  specialties: '',
  description: '',
  work_start_time: '09:00',
  work_end_time: '21:00'
})

const fetchList = async () => {
  loading.value = true
  try {
    const res = await getTechnicians({ page: page.value, page_size: pageSize.value })
    technicians.value = res.data.list
    total.value = res.data.total
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

const handleEdit = (row: Technician) => {
  editing.value = true
  editingId.value = row.id
  Object.assign(form, {
    name: row.name,
    title: row.title,
    specialties: row.specialties,
    description: row.description,
    work_start_time: row.work_start_time,
    work_end_time: row.work_end_time
  })
  addVisible.value = true
}

const handleSave = async () => {
  if (!form.name) {
    ElMessage.warning('请输入姓名')
    return
  }

  saving.value = true
  try {
    if (editing.value && editingId.value) {
      await updateTechnician(editingId.value, form)
      ElMessage.success('更新成功')
    } else {
      await createTechnician(form)
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
