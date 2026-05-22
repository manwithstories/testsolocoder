<template>
  <div class="task-list">
    <div class="page-header">
      <el-select v-model="statusFilter" placeholder="状态" clearable style="width: 120px" @change="fetchTasks">
        <el-option label="待处理" value="pending" />
        <el-option label="进行中" value="in_progress" />
        <el-option label="已完成" value="completed" />
        <el-option label="已取消" value="cancelled" />
      </el-select>
      <el-button type="primary" :icon="Plus" @click="showCreateDialog = true">
        添加任务
      </el-button>
      <el-button :icon="Document" @click="showTemplateDialog = true">
        使用模板
      </el-button>
    </div>

    <el-table :data="tasks" v-loading="loading" stripe>
      <el-table-column label="任务名称" min-width="200">
        <template #default="{ row }">
          <el-checkbox
            :model-value="row.status === 'completed'"
            @change="toggleTaskStatus(row)"
          />
          <span :style="{ textDecoration: row.status === 'completed' ? 'line-through' : 'none' }">
            {{ row.title }}
          </span>
        </template>
      </el-table-column>
      <el-table-column prop="category" label="分类" width="100" />
      <el-table-column prop="assignee" label="负责人" width="120" />
      <el-table-column label="截止日期" width="120">
        <template #default="{ row }">
          <span :class="{ overdue: isOverdue(row) }">
            {{ formatDate(row.due_date) }}
          </span>
        </template>
      </el-table-column>
      <el-table-column label="优先级" width="100">
        <template #default="{ row }">
          <el-tag :type="priorityType(row.priority)" size="small">{{ priorityText(row.priority) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="150" fixed="right">
        <template #default="{ row }">
          <el-button type="primary" link @click="editTask(row)">编辑</el-button>
          <el-button type="danger" link @click="deleteTask(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="showCreateDialog" :title="editingTask ? '编辑任务' : '添加任务'" width="500px">
      <el-form ref="taskForm" :model="taskForm" :rules="taskRules" label-width="100px">
        <el-form-item label="任务名称" prop="title">
          <el-input v-model="taskForm.title" />
        </el-form-item>
        <el-form-item label="分类">
          <el-select v-model="taskForm.category" style="width: 100%">
            <el-option v-for="cat in categories" :key="cat" :label="cat" :value="cat" />
          </el-select>
        </el-form-item>
        <el-form-item label="负责人">
          <el-input v-model="taskForm.assignee" />
        </el-form-item>
        <el-form-item label="截止日期">
          <el-date-picker
            v-model="taskForm.due_date"
            type="date"
            style="width: 100%"
            value-format="YYYY-MM-DD"
          />
        </el-form-item>
        <el-form-item label="优先级">
          <el-select v-model="taskForm.priority" style="width: 100%">
            <el-option label="高" value="high" />
            <el-option label="中" value="medium" />
            <el-option label="低" value="low" />
          </el-select>
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="taskForm.description" type="textarea" :rows="3" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showCreateDialog = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="saveTask">保存</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="showTemplateDialog" title="应用任务模板" width="400px">
      <el-form label-width="80px">
        <el-form-item label="选择模板">
          <el-select v-model="selectedTemplate" style="width: 100%">
            <el-option
              v-for="template in templates"
              :key="template.id"
              :label="template.name"
              :value="template.id"
            />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showTemplateDialog = false">取消</el-button>
        <el-button type="primary" :loading="applying" @click="applyTemplate">应用</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { taskApi } from '@/api/task'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import { Plus, Document } from '@element-plus/icons-vue'
import dayjs from 'dayjs'
import type { Task, TaskTemplate } from '@/types'

const props = defineProps<{
  weddingId: number
}>()

const loading = ref(false)
const saving = ref(false)
const applying = ref(false)
const tasks = ref<Task[]>([])
const templates = ref<TaskTemplate[]>([])
const categories = ref<string[]>([])
const statusFilter = ref('')
const showCreateDialog = ref(false)
const showTemplateDialog = ref(false)
const editingTask = ref<Task | null>(null)
const selectedTemplate = ref<number | null>(null)

const taskForm = reactive({
  title: '',
  category: '',
  assignee: '',
  due_date: '',
  priority: 'medium',
  description: ''
})

const taskRules: FormRules = {
  title: [{ required: true, message: '请输入任务名称', trigger: 'blur' }]
}

const taskFormRef = ref<FormInstance>()

function formatDate(date?: string) {
  return date ? dayjs(date).format('YYYY-MM-DD') : '-'
}

function priorityType(priority: string) {
  const types: Record<string, string> = { high: 'danger', medium: 'warning', low: 'success' }
  return types[priority] || 'info'
}

function priorityText(priority: string) {
  const texts: Record<string, string> = { high: '高', medium: '中', low: '低' }
  return texts[priority] || priority
}

function isOverdue(task: Task) {
  return task.status !== 'completed' && task.due_date && dayjs(task.due_date).isBefore(dayjs())
}

async function fetchCategories() {
  try {
    const res = await taskApi.getCategories(props.weddingId)
    categories.value = res.data
  } catch (error) {
    console.error('Failed to fetch categories:', error)
  }
}

async function fetchTasks() {
  loading.value = true
  try {
    const params: any = {}
    if (statusFilter.value) {
      params.status = statusFilter.value
    }
    const res = await taskApi.getList(props.weddingId, params)
    tasks.value = res.data
  } catch (error) {
    console.error('Failed to fetch tasks:', error)
  } finally {
    loading.value = false
  }
}

async function fetchTemplates() {
  try {
    const res = await taskApi.getTemplates()
    templates.value = res.data
  } catch (error) {
    console.error('Failed to fetch templates:', error)
  }
}

function editTask(task: Task) {
  editingTask.value = task
  Object.assign(taskForm, {
    title: task.title,
    category: task.category || '',
    assignee: task.assignee || '',
    due_date: task.due_date ? task.due_date.split('T')[0] : '',
    priority: task.priority,
    description: task.description || ''
  })
  showCreateDialog.value = true
}

async function saveTask() {
  if (!taskFormRef.value) return
  
  await taskFormRef.value.validate(async (valid) => {
    if (valid) {
      saving.value = true
      try {
        if (editingTask.value) {
          await taskApi.update(props.weddingId, editingTask.value.id, taskForm)
          ElMessage.success('更新成功')
        } else {
          await taskApi.create(props.weddingId, taskForm)
          ElMessage.success('创建成功')
        }
        showCreateDialog.value = false
        fetchTasks()
      } catch (error: any) {
        ElMessage.error(error.message || '保存失败')
      } finally {
        saving.value = false
      }
    }
  })
}

async function toggleTaskStatus(task: Task) {
  const newStatus = task.status === 'completed' ? 'pending' : 'completed'
  try {
    await taskApi.updateStatus(props.weddingId, task.id, newStatus)
    task.status = newStatus
    ElMessage.success('状态更新成功')
  } catch (error) {
    console.error('Failed to update task status:', error)
  }
}

async function deleteTask(task: Task) {
  try {
    await ElMessageBox.confirm(`确定要删除任务"${task.title}"吗？`, '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    await taskApi.delete(props.weddingId, task.id)
    ElMessage.success('删除成功')
    fetchTasks()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Failed to delete task:', error)
    }
  }
}

async function applyTemplate() {
  if (!selectedTemplate.value) {
    ElMessage.warning('请选择一个模板')
    return
  }
  
  applying.value = true
  try {
    await taskApi.applyTemplate(props.weddingId, selectedTemplate.value)
    ElMessage.success('模板应用成功')
    showTemplateDialog.value = false
    fetchTasks()
  } catch (error: any) {
    ElMessage.error(error.message || '应用失败')
  } finally {
    applying.value = false
  }
}

onMounted(() => {
  fetchCategories()
  fetchTasks()
  fetchTemplates()
})
</script>

<style scoped>
.task-list {
  padding: 0;
}

.page-header {
  display: flex;
  gap: 16px;
  margin-bottom: 20px;
}

.overdue {
  color: #F56C6C;
}
</style>
