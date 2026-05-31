<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>文档管理</span>
        </div>
      </template>

      <el-alert
        title="请先在项目详情页面管理文档"
        type="info"
        show-icon
        :closable="false"
      />

      <el-table :data="documents" v-loading="loading" style="margin-top: 20px">
        <el-table-column prop="project_id" label="项目ID" width="100" />
        <el-table-column prop="file_name" label="文件名" />
        <el-table-column prop="file_type" label="类型" width="100" />
        <el-table-column prop="file_size" label="大小" width="120">
          <template #default="{ row }">{{ formatFileSize(row.file_size) }}</template>
        </el-table-column>
        <el-table-column prop="version" label="版本" width="80" />
        <el-table-column label="来源" width="100">
          <template #default="{ row }">
            <el-tag :type="row.is_source ? 'primary' : 'success'" size="small">
              {{ row.is_source ? '原文' : '译文' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="120">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleDownload(row)">下载</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { downloadDocument } from '@/api/document'
import { saveAs } from 'file-saver'

const documents = ref<any[]>([])
const loading = ref(false)

function formatFileSize(bytes: number) {
  if (bytes < 1024) return bytes + ' B'
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB'
  return (bytes / 1024 / 1024).toFixed(1) + ' MB'
}

async function handleDownload(row: any) {
  try {
    const blob = await downloadDocument(row.id) as any
    saveAs(blob, row.file_name)
  } catch (e) {
    console.error(e)
  }
}

onMounted(() => {})
</script>

<style lang="scss" scoped>
.page-container {
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }
}
</style>
