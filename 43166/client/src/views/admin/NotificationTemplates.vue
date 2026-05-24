<template>
  <div class="notification-templates">
    <div class="page-header flex-between">
      <h2 class="page-title">通知模板管理</h2>
      <el-button type="primary" @click="showCreateDialog = true">
        <el-icon><Plus /></el-icon>
        添加模板
      </el-button>
    </div>

    <div class="card">
      <el-table :data="templates" style="width: 100%" v-loading="loading">
        <el-table-column prop="code" label="模板编码" width="150" />
        <el-table-column prop="name" label="模板名称" />
        <el-table-column prop="type" label="类型" width="100">
          <template #default="{ row }">
            {{ getTypeText(row.type) }}
          </template>
        </el-table-column>
        <el-table-column prop="title" label="标题" />
        <el-table-column prop="isActive" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.isActive ? 'success' : 'info'">
              {{ row.isActive ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="editTemplate(row)">编辑</el-button>
            <el-button type="danger" link @click="deleteTemplate(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <el-dialog v-model="showCreateDialog" :title="isEdit ? '编辑模板' : '添加模板'" width="600px">
      <el-form ref="templateFormRef" :model="templateForm" :rules="templateRules" label-width="100px">
        <el-form-item label="模板编码" prop="code">
          <el-input v-model="templateForm.code" :disabled="isEdit" />
        </el-form-item>
        <el-form-item label="模板名称" prop="name">
          <el-input v-model="templateForm.name" />
        </el-form-item>
        <el-form-item label="类型" prop="type">
          <el-select v-model="templateForm.type" style="width: 100%">
            <el-option label="系统通知" value="system" />
            <el-option label="邮件通知" value="email" />
            <el-option label="短信通知" value="sms" />
          </el-select>
        </el-form-item>
        <el-form-item label="标题" prop="title">
          <el-input v-model="templateForm.title" />
        </el-form-item>
        <el-form-item label="内容" prop="content">
          <el-input
            v-model="templateForm.content"
            type="textarea"
            :rows="4"
            placeholder="使用 {{变量名}} 表示动态变量"
          />
        </el-form-item>
        <el-form-item label="变量列表">
          <el-input
            v-model="templateForm.variables"
            placeholder="多个变量用逗号分隔，如: company_name,status"
          />
        </el-form-item>
        <el-form-item label="启用状态">
          <el-switch v-model="templateForm.isActive" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showCreateDialog = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="handleSaveTemplate">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import { notificationApi } from '@/api/notification'
import { NotificationTemplate } from '@/types'

const loading = ref(false)
const templates = ref<NotificationTemplate[]>([])
const showCreateDialog = ref(false)
const isEdit = ref(false)
const saving = ref(false)
const templateFormRef = ref<FormInstance>()

const templateForm = reactive({
  id: 0,
  code: '',
  name: '',
  type: 'system',
  title: '',
  content: '',
  variables: '',
  isActive: true
})

const templateRules: FormRules = {
  code: [{ required: true, message: '请输入模板编码', trigger: 'blur' }],
  name: [{ required: true, message: '请输入模板名称', trigger: 'blur' }],
  type: [{ required: true, message: '请选择类型', trigger: 'change' }],
  title: [{ required: true, message: '请输入标题', trigger: 'blur' }],
  content: [{ required: true, message: '请输入内容', trigger: 'blur' }]
}

const fetchData = async () => {
  loading.value = true
  try {
    const res = await notificationApi.getTemplates()
    templates.value = res || []
  } catch (error) {
    console.error('获取模板列表失败:', error)
  } finally {
    loading.value = false
  }
}

const getTypeText = (type: string) => {
  const map: Record<string, string> = {
    system: '系统',
    email: '邮件',
    sms: '短信'
  }
  return map[type] || type
}

const editTemplate = (row: NotificationTemplate) => {
  isEdit.value = true
  Object.assign(templateForm, row)
  showCreateDialog.value = true
}

const deleteTemplate = async (row: NotificationTemplate) => {
  try {
    await ElMessageBox.confirm(`确认删除模板 ${row.name}？`, '提示', {
      confirmButtonText: '确认',
      cancelButtonText: '取消',
      type: 'warning'
    })

    await notificationApi.deleteTemplate(row.id)
    ElMessage.success('删除成功')
    fetchData()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error.message || '删除失败')
    }
  }
}

const handleSaveTemplate = async () => {
  if (!templateFormRef.value) return

  await templateFormRef.value.validate(async (valid) => {
    if (valid) {
      saving.value = true
      try {
        if (isEdit.value) {
          await notificationApi.updateTemplate(templateForm.id, templateForm)
        } else {
          await notificationApi.createTemplate(templateForm as any)
        }
        ElMessage.success('保存成功')
        showCreateDialog.value = false
        fetchData()
      } catch (error: any) {
        ElMessage.error(error.message || '保存失败')
      } finally {
        saving.value = false
      }
    }
  })
}

onMounted(fetchData)
</script>
