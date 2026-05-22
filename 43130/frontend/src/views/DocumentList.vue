<template>
  <div class="document-list">
    <div class="page-header">
      <el-select v-model="categoryFilter" placeholder="分类" clearable style="width: 140px" @change="fetchDocuments">
        <el-option v-for="cat in categories" :key="cat" :label="cat" :value="cat" />
      </el-select>
      <el-upload
        :show-file-list="false"
        :before-upload="handleUpload"
        accept=".pdf,.doc,.docx,.xls,.xlsx,.jpg,.jpeg,.png"
      >
        <el-button type="primary" :icon="Upload">上传文档</el-button>
      </el-upload>
    </div>

    <el-table :data="documents" v-loading="loading" stripe>
      <el-table-column label="文件名" min-width="200">
        <template #default="{ row }">
          <el-icon style="margin-right: 8px"><Document /></el-icon>
          <span>{{ row.file_name }}</span>
        </template>
      </el-table-column>
      <el-table-column prop="category" label="分类" width="100" />
      <el-table-column label="大小" width="100">
        <template #default="{ row }">
          {{ formatFileSize(row.file_size) }}
        </template>
      </el-table-column>
      <el-table-column label="版本" width="80">
        <template #default="{ row }">
          <el-tag size="small" type="info">v{{ row.version }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="上传时间" width="180">
        <template #default="{ row }">
          {{ formatDate(row.created_at) }}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="260" fixed="right">
        <template #default="{ row }">
          <el-button type="primary" link @click="downloadDocument(row)">下载</el-button>
          <el-button type="info" link @click="showVersions(row)">历史版本</el-button>
          <el-button type="success" link @click="uploadNewVersion(row)">上传新版本</el-button>
          <el-button type="danger" link @click="deleteDocument(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="showVersionDialog" title="上传新版本" width="400px">
      <el-upload
        :show-file-list="false"
        :before-upload="handleNewVersion"
        accept=".pdf,.doc,.docx,.xls,.xlsx,.jpg,.jpeg,.png"
      >
        <el-button type="primary">选择文件上传</el-button>
      </el-upload>
      <div style="margin-top: 12px; color: #909399; font-size: 12px;">
        当前文件：{{ currentDocument?.file_name }}<br />
        当前版本：v{{ currentDocument?.version }}
      </div>
    </el-dialog>

    <el-dialog v-model="showHistoryDialog" title="历史版本" width="700px">
      <el-empty v-if="versions.length === 0" description="暂无历史版本" />
      <el-table v-else :data="versions" stripe>
        <el-table-column label="版本" width="100">
          <template #default="{ row }">
            <el-tag size="small">v{{ row.version }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="file_name" label="文件名" min-width="180" />
        <el-table-column label="大小" width="100">
          <template #default="{ row }">
            {{ formatFileSize(row.file_size) }}
          </template>
        </el-table-column>
        <el-table-column label="上传时间" width="180">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="120">
          <template #default="{ row }">
            <el-button type="primary" link @click="downloadHistoryVersion(row)">下载</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { documentApi } from '@/api/document'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Upload, Document as DocumentIcon } from '@element-plus/icons-vue'
import dayjs from 'dayjs'
import type { Document as DocType } from '@/types'

const props = defineProps<{
  weddingId: number
}>()

const loading = ref(false)
const documents = ref<DocType[]>([])
const categories = ref<string[]>([])
const categoryFilter = ref('')
const showVersionDialog = ref(false)
const showHistoryDialog = ref(false)
const currentDocument = ref<DocType | null>(null)
const versions = ref<DocType[]>([])

const Document = DocumentIcon

function formatFileSize(bytes: number) {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

function formatDate(date: string) {
  return dayjs(date).format('YYYY-MM-DD HH:mm')
}

async function fetchCategories() {
  try {
    const res = await documentApi.getCategories(props.weddingId)
    categories.value = res.data
  } catch (error) {
    console.error('Failed to fetch categories:', error)
  }
}

async function fetchDocuments() {
  loading.value = true
  try {
    const params: any = {}
    if (categoryFilter.value) {
      params.category = categoryFilter.value
    }
    const res = await documentApi.getList(props.weddingId, params)
    documents.value = res.data
  } catch (error) {
    console.error('Failed to fetch documents:', error)
  } finally {
    loading.value = false
  }
}

async function handleUpload(file: File) {
  try {
    await documentApi.upload(props.weddingId, file, {})
    ElMessage.success('上传成功')
    fetchDocuments()
  } catch (error: any) {
    ElMessage.error(error.message || '上传失败')
  }
  return false
}

async function downloadDocument(doc: DocType) {
  try {
    const res = await documentApi.download(props.weddingId, doc.id) as any
    const blob = new Blob([res])
    const url = URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = doc.file_name
    link.click()
    URL.revokeObjectURL(url)
  } catch (error) {
    ElMessage.error('下载失败')
  }
}

function uploadNewVersion(doc: DocType) {
  currentDocument.value = doc
  showVersionDialog.value = true
}

async function showVersions(doc: DocType) {
  currentDocument.value = doc
  showHistoryDialog.value = true
  try {
    const res = await documentApi.getById(props.weddingId, doc.id)
    versions.value = res.data.versions || []
  } catch (error) {
    console.error('Failed to fetch versions:', error)
    ElMessage.error('获取版本列表失败')
  }
}

async function downloadHistoryVersion(ver: DocType) {
  try {
    const res = await documentApi.download(props.weddingId, ver.id) as any
    const blob = new Blob([res])
    const url = URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = ver.file_name
    link.click()
    URL.revokeObjectURL(url)
  } catch (error) {
    ElMessage.error('下载失败')
  }
}

async function handleNewVersion(file: File) {
  if (!currentDocument.value) return
  
  try {
    await documentApi.uploadVersion(props.weddingId, currentDocument.value.id, file)
    ElMessage.success('新版本上传成功')
    showVersionDialog.value = false
    fetchDocuments()
  } catch (error: any) {
    ElMessage.error(error.message || '上传失败')
  }
  return false
}

async function deleteDocument(doc: DocType) {
  try {
    await ElMessageBox.confirm(`确定要删除文档"${doc.file_name}"吗？这将删除所有历史版本！`, '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    await documentApi.delete(props.weddingId, doc.id)
    ElMessage.success('删除成功')
    fetchDocuments()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Failed to delete document:', error)
    }
  }
}

onMounted(() => {
  fetchCategories()
  fetchDocuments()
})
</script>

<style scoped>
.document-list {
  padding: 0;
}

.page-header {
  display: flex;
  gap: 16px;
  margin-bottom: 20px;
}
</style>
