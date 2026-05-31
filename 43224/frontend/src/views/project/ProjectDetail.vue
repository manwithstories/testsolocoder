<template>
  <div class="project-detail">
    <el-page-header @back="$router.back()">
      <template #content>项目详情</template>
    </el-page-header>

    <el-card v-loading="loading" style="margin-top: 20px">
      <el-descriptions :column="2" border>
        <el-descriptions-item label="项目名称">{{ project?.title }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="getStatusType(project?.status)">{{ getStatusText(project?.status) }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="语言对">{{ project?.source_lang }} - {{ project?.target_lang }}</el-descriptions-item>
        <el-descriptions-item label="字数">{{ project?.word_count }}</el-descriptions-item>
        <el-descriptions-item label="金额">{{ project?.total_amount?.toFixed(2) }} 元</el-descriptions-item>
        <el-descriptions-item label="紧急程度">{{ getUrgencyText(project?.urgency) }}</el-descriptions-item>
        <el-descriptions-item label="截止日期">{{ formatDate(project?.deadline) }}</el-descriptions-item>
        <el-descriptions-item label="客户">{{ project?.client?.username }}</el-descriptions-item>
        <el-descriptions-item label="项目经理">{{ project?.pm?.username || '未指派' }}</el-descriptions-item>
        <el-descriptions-item label="译者">{{ project?.translator?.username || '未指派' }}</el-descriptions-item>
        <el-descriptions-item label="专业领域" :span="2">
          <el-tag v-for="tag in project?.expertise_tags" :key="tag.id" style="margin-right: 8px">
            {{ tag.name }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="项目描述" :span="2">{{ project?.description }}</el-descriptions-item>
      </el-descriptions>
    </el-card>

    <el-tabs v-model="activeTab" style="margin-top: 20px">
      <el-tab-pane label="文档管理" name="documents">
        <div class="section-header">
          <span>源文档</span>
          <el-upload
            :action="`/api/projects/${route.params.id}/documents`"
            :headers="{ Authorization: `Bearer ${userStore.token}` }"
            :data="{ is_source: 'true' }"
            :show-file-list="false"
            :on-success="handleUploadSuccess"
          >
            <el-button type="primary"><el-icon><Upload /></el-icon>上传源文档</el-button>
          </el-upload>
        </div>
        <el-table :data="sourceDocuments" v-loading="docLoading">
          <el-table-column prop="file_name" label="文件名" />
          <el-table-column prop="file_type" label="类型" width="100" />
          <el-table-column prop="file_size" label="大小" width="120">
            <template #default="{ row }">{{ formatFileSize(row.file_size) }}</template>
          </el-table-column>
          <el-table-column prop="version" label="版本" width="80" />
          <el-table-column label="操作" width="200">
            <template #default="{ row }">
              <el-button type="primary" link @click="handleDownload(row)">下载</el-button>
              <el-button type="primary" link @click="handleExtract(row)">提取片段</el-button>
            </template>
          </el-table-column>
        </el-table>

        <div class="section-header" style="margin-top: 20px">
          <span>译文文档</span>
          <el-upload
            :action="`/api/projects/${route.params.id}/documents`"
            :headers="{ Authorization: `Bearer ${userStore.token}` }"
            :data="{ is_source: 'false' }"
            :show-file-list="false"
            :on-success="handleUploadSuccess"
          >
            <el-button type="primary"><el-icon><Upload /></el-icon>上传译文</el-button>
          </el-upload>
        </div>
        <el-table :data="translatedDocuments" v-loading="docLoading">
          <el-table-column prop="file_name" label="文件名" />
          <el-table-column prop="file_type" label="类型" width="100" />
          <el-table-column prop="file_size" label="大小" width="120">
            <template #default="{ row }">{{ formatFileSize(row.file_size) }}</template>
          </el-table-column>
          <el-table-column prop="version" label="版本" width="80" />
          <el-table-column label="操作" width="100">
            <template #default="{ row }">
              <el-button type="primary" link @click="handleDownload(row)">下载</el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-tab-pane>

      <el-tab-pane label="翻译片段" name="segments">
        <el-table :data="segments" v-loading="segLoading" max-height="500">
          <el-table-column prop="source_text" label="原文" min-width="200">
            <template #default="{ row }">
              <div class="segment-text">{{ row.source_text }}</div>
            </template>
          </el-table-column>
          <el-table-column prop="translated_text" label="译文" min-width="200">
            <template #default="{ row }">
              <div class="segment-text">{{ row.translated_text || '未翻译' }}</div>
            </template>
          </el-table-column>
          <el-table-column prop="status" label="状态" width="100">
            <template #default="{ row }">
              <el-tag v-if="row.is_memory" type="success" size="small">记忆</el-tag>
              <el-tag v-else :type="row.status === 'translated' ? 'success' : 'info'" size="small">
                {{ row.status === 'translated' ? '已翻译' : '待翻译' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="相似度" width="100">
            <template #default="{ row }">{{ row.memory_match ? (row.memory_match * 100).toFixed(0) + '%' : '-' }}</template>
          </el-table-column>
        </el-table>
      </el-tab-pane>

      <el-tab-pane label="审核记录" name="reviews">
        <el-table :data="reviewTasks" v-loading="reviewLoading">
          <el-table-column prop="round" label="轮次" width="80" />
          <el-table-column prop="reviewer.username" label="审核人" width="120" />
          <el-table-column prop="comment" label="审核意见" />
          <el-table-column prop="suggestion" label="建议修改" />
          <el-table-column prop="status" label="状态" width="100">
            <template #default="{ row }">
              <el-tag :type="row.status === 'approved' ? 'success' : row.status === 'rejected' ? 'danger' : 'warning'" size="small">
                {{ row.status === 'approved' ? '通过' : row.status === 'rejected' ? '不通过' : '待处理' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="created_at" label="时间" width="160">
            <template #default="{ row }">{{ formatDate(row.created_at) }}</template>
          </el-table-column>
        </el-table>
      </el-tab-pane>

      <el-tab-pane label="评论交流" name="comments">
        <div class="comment-list">
          <div v-for="comment in project?.comments || []" :key="comment.id" class="comment-item">
            <el-avatar :size="36">{{ comment.user?.username?.charAt(0)?.toUpperCase() }}</el-avatar>
            <div class="comment-content">
              <div class="comment-header">
                <span class="comment-user">{{ comment.user?.username }}</span>
                <span class="comment-time">{{ formatDate(comment.created_at) }}</span>
              </div>
              <div class="comment-text">{{ comment.content }}</div>
            </div>
          </div>
        </div>
        <div class="comment-input">
          <el-input
            v-model="newComment"
            type="textarea"
            :rows="2"
            placeholder="输入评论..."
            @keyup.enter.ctrl="handleAddComment"
          />
          <el-button type="primary" :loading="commentLoading" @click="handleAddComment" style="margin-top: 8px">
            发送 (Ctrl+Enter)
          </el-button>
        </div>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/stores/user'
import { getProject, addProjectComment } from '@/api/project'
import { listDocuments, downloadDocument, extractSegments, listSegments } from '@/api/document'
import { listReviewTasks } from '@/api/statistics'
import dayjs from 'dayjs'
import { saveAs } from 'file-saver'

const route = useRoute()
const userStore = useUserStore()

const activeTab = ref('documents')
const loading = ref(false)
const project = ref<any>(null)

const sourceDocuments = ref<any[]>([])
const translatedDocuments = ref<any[]>([])
const docLoading = ref(false)

const segments = ref<any[]>([])
const segLoading = ref(false)

const reviewTasks = ref<any[]>([])
const reviewLoading = ref(false)

const newComment = ref('')
const commentLoading = ref(false)

async function loadData() {
  loading.value = true
  try {
    const id = route.params.id as string
    project.value = await getProject(Number(id))
    loadDocuments()
    loadSegments()
    loadReviews()
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

async function loadDocuments() {
  docLoading.value = true
  try {
    const id = route.params.id as string
    const allDocs = await listDocuments(Number(id)) as any[]
    sourceDocuments.value = allDocs.filter(d => d.is_source)
    translatedDocuments.value = allDocs.filter(d => !d.is_source)
  } catch (e) {
    console.error(e)
  } finally {
    docLoading.value = false
  }
}

async function loadSegments() {
  segLoading.value = true
  try {
    const id = route.params.id as string
    const res = await listSegments(Number(id))
    segments.value = res || []
  } catch (e) {
    console.error(e)
  } finally {
    segLoading.value = false
  }
}

async function loadReviews() {
  reviewLoading.value = true
  try {
    const id = route.params.id as string
    const res = await listReviewTasks({ project_id: id })
    if (Array.isArray(res)) {
      reviewTasks.value = res
    } else {
      reviewTasks.value = (res as any)?.list || []
    }
  } catch (e) {
    console.error(e)
  } finally {
    reviewLoading.value = false
  }
}

function handleUploadSuccess() {
  ElMessage.success('上传成功')
  loadDocuments()
}

async function handleDownload(row: any) {
  try {
    const blob = await downloadDocument(row.id) as any
    saveAs(blob, row.file_name)
  } catch (e) {
    console.error(e)
  }
}

async function handleExtract(row: any) {
  try {
    await extractSegments(row.id)
    ElMessage.success('提取完成')
    loadSegments()
  } catch (e: any) {
    ElMessage.error(e.message || '提取失败')
  }
}

async function handleAddComment() {
  if (!newComment.value.trim()) return
  commentLoading.value = true
  try {
    const id = route.params.id as string
    await addProjectComment(Number(id), newComment.value)
    newComment.value = ''
    loadData()
    ElMessage.success('评论已发送')
  } catch (e: any) {
    ElMessage.error(e.message || '发送失败')
  } finally {
    commentLoading.value = false
  }
}

function getStatusType(status: string) {
  const map: Record<string, string> = {
    pending: 'warning', approved: 'primary', assigned: 'info',
    in_progress: '', review: 'warning', completed: 'success', cancelled: 'danger'
  }
  return map[status] || ''
}

function getStatusText(status: string) {
  const map: Record<string, string> = {
    pending: '待审核', approved: '已通过', assigned: '已分配',
    in_progress: '进行中', review: '审核中', completed: '已完成', cancelled: '已取消'
  }
  return map[status] || status
}

function getUrgencyText(urgency: string) {
  const map: Record<string, string> = { normal: '普通', fast: '加急', urgent: '特急' }
  return map[urgency] || urgency
}

function formatDate(date: string) {
  return dayjs(date).format('YYYY-MM-DD HH:mm')
}

function formatFileSize(bytes: number) {
  if (bytes < 1024) return bytes + ' B'
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB'
  return (bytes / 1024 / 1024).toFixed(1) + ' MB'
}

onMounted(loadData)
</script>

<style lang="scss" scoped>
.project-detail {
  .section-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 12px;

    span {
      font-size: 16px;
      font-weight: 500;
    }
  }

  .segment-text {
    max-height: 100px;
    overflow-y: auto;
    white-space: pre-wrap;
  }

  .comment-list {
    max-height: 400px;
    overflow-y: auto;
    padding: 16px;
    background: #f5f7fa;
    border-radius: 4px;
  }

  .comment-item {
    display: flex;
    gap: 12px;
    margin-bottom: 16px;
  }

  .comment-content {
    flex: 1;
  }

  .comment-header {
    display: flex;
    gap: 12px;
    margin-bottom: 4px;

    .comment-user {
      font-weight: 500;
      color: #303133;
    }

    .comment-time {
      color: #909399;
      font-size: 12px;
    }
  }

  .comment-text {
    color: #606266;
  }

  .comment-input {
    margin-top: 16px;
  }
}
</style>
