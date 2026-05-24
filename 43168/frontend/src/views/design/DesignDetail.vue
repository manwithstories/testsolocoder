<template>
  <div class="design-detail">
    <el-page-header @back="router.back()">
      <template #content>
        <span>方案详情</span>
      </template>
    </el-page-header>

    <el-row :gutter="16" style="margin-top: 16px">
      <el-col :xs="24" :md="16">
        <el-card shadow="never">
          <template #header>
            <div class="card-header">
              <span>{{ detail.name }}</span>
              <div>
                <el-tag :type="getStatusType(detail.status)">{{ getStatusLabel(detail.status) }}</el-tag>
              </div>
            </div>
          </template>

          <el-descriptions :column="2" border>
            <el-descriptions-item label="业主">{{ detail.ownerName || '-' }}</el-descriptions-item>
            <el-descriptions-item label="设计师">{{ detail.designerName || '-' }}</el-descriptions-item>
            <el-descriptions-item label="房型">{{ detail.houseType || '-' }}</el-descriptions-item>
            <el-descriptions-item label="面积">{{ detail.area ? `${detail.area} ㎡` : '-' }}</el-descriptions-item>
            <el-descriptions-item label="预算">{{ detail.budget ? `¥${detail.budget.toLocaleString()}` : '-' }}</el-descriptions-item>
            <el-descriptions-item label="创建时间">{{ detail.createdAt || '-' }}</el-descriptions-item>
            <el-descriptions-item label="描述" :span="2">{{ detail.description || '-' }}</el-descriptions-item>
          </el-descriptions>
        </el-card>

        <el-card shadow="never" style="margin-top: 16px">
          <template #header>
            <div class="card-header">
              <span>设计图片</span>
            </div>
          </template>
          <div class="image-grid">
            <div
              v-for="img in detail.images || []"
              :key="img.id"
              class="image-item"
              @click="handleAnnotate(img)"
            >
              <el-image :src="img.url" fit="cover" style="width: 100%; height: 100%" />
              <div class="image-actions">
                <el-button size="small" type="primary" :icon="EditPen" @click.stop="handleAnnotate(img)">批注</el-button>
              </div>
            </div>
            <el-empty v-if="!detail.images || detail.images.length === 0" description="暂无图片" />
          </div>
        </el-card>
      </el-col>

      <el-col :xs="24" :md="8">
        <el-card shadow="never">
          <template #header>
            <span>状态操作</span>
          </template>
          <div class="action-group">
            <el-button
              v-for="act in statusActions"
              :key="act.value"
              type="primary"
              plain
              style="width: 100%; margin-bottom: 8px"
              @click="handleStatusChange(act.value)"
            >
              {{ act.label }}
            </el-button>
          </div>
        </el-card>

        <el-card shadow="never" style="margin-top: 16px">
          <template #header>
            <span>批注记录</span>
          </template>
          <div class="annotation-list">
            <div v-for="item in annotations" :key="item.id" class="annotation-item">
              <div class="annotation-header">
                <el-icon><ChatDotRound /></el-icon>
                <span class="annotation-author">{{ item.author || '匿名' }}</span>
                <span class="annotation-time">{{ item.createdAt }}</span>
              </div>
              <div class="annotation-content">{{ item.content }}</div>
            </div>
            <el-empty v-if="annotations.length === 0" description="暂无批注" :image-size="80" />
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-dialog v-model="annotateVisible" title="添加批注" width="520px">
      <div class="annotate-preview">
        <el-image v-if="currentImage" :src="currentImage.url" fit="contain" style="width: 100%; height: 280px" />
      </div>
      <el-form :model="annotateForm" label-width="80px" style="margin-top: 16px">
        <el-form-item label="批注内容">
          <el-input v-model="annotateForm.content" type="textarea" :rows="4" placeholder="请输入批注内容" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="annotateVisible = false">取消</el-button>
        <el-button type="primary" :loading="annotating" @click="submitAnnotation">提交</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { EditPen, ChatDotRound } from '@element-plus/icons-vue'
import {
  getDesign,
  updateDesignStatus,
  addAnnotation,
  listAnnotations,
  type DesignProject,
  type DesignImage,
  type DesignAnnotation
} from '@/api/design'

const route = useRoute()
const router = useRouter()
const detail = ref<DesignProject>({} as DesignProject)
const annotations = ref<DesignAnnotation[]>([])

const statusOptions = [
  { label: '草稿', value: 0 },
  { label: '进行中', value: 1 },
  { label: '待审核', value: 2 },
  { label: '已通过', value: 3 },
  { label: '已驳回', value: 4 },
  { label: '已归档', value: 5 }
]

const statusActions = [
  { label: '提交审核', value: 2 },
  { label: '审核通过', value: 3 },
  { label: '驳回', value: 4 },
  { label: '归档', value: 5 }
]

function getStatusLabel(status: number) {
  return statusOptions.find((i) => i.value === status)?.label ?? '未知'
}
function getStatusType(status: number) {
  const map: Record<number, string> = {
    0: 'info',
    1: 'primary',
    2: 'warning',
    3: 'success',
    4: 'danger',
    5: ''
  }
  return map[status] ?? ''
}

async function fetchDetail() {
  const id = route.params.id
  if (!id) return
  detail.value = await getDesign(id)
}

async function fetchAnnotations() {
  const id = route.params.id
  if (!id) return
  try {
    annotations.value = await listAnnotations(id)
  } catch {
    annotations.value = []
  }
}

async function handleStatusChange(status: number) {
  const id = route.params.id
  if (!id) return
  await updateDesignStatus(id, status)
  ElMessage.success('状态更新成功')
  fetchDetail()
}

const annotateVisible = ref(false)
const annotating = ref(false)
const currentImage = ref<DesignImage | null>(null)
const annotateForm = reactive({ content: '' })

function handleAnnotate(img: DesignImage) {
  currentImage.value = img
  annotateForm.content = ''
  annotateVisible.value = true
}

async function submitAnnotation() {
  if (!annotateForm.content.trim()) {
    ElMessage.warning('请输入批注内容')
    return
  }
  const id = route.params.id
  if (!id) return
  annotating.value = true
  try {
    await addAnnotation(id, {
      imageId: currentImage.value?.id,
      content: annotateForm.content
    })
    ElMessage.success('批注提交成功')
    annotateVisible.value = false
    fetchAnnotations()
  } finally {
    annotating.value = false
  }
}

onMounted(() => {
  fetchDetail()
  fetchAnnotations()
})
</script>

<style lang="scss" scoped>
.design-detail {
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }
  .image-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(180px, 1fr));
    gap: 12px;
    .image-item {
      position: relative;
      height: 140px;
      border-radius: 4px;
      overflow: hidden;
      cursor: pointer;
      &:hover .image-actions {
        opacity: 1;
      }
      .image-actions {
        position: absolute;
        inset: 0;
        background: rgba(0, 0, 0, 0.45);
        display: flex;
        align-items: center;
        justify-content: center;
        opacity: 0;
        transition: opacity 0.2s;
      }
    }
  }
  .action-group {
    display: flex;
    flex-direction: column;
  }
  .annotation-list {
    max-height: 520px;
    overflow-y: auto;
    .annotation-item {
      padding: 12px 0;
      border-bottom: 1px solid #f0f2f5;
      &:last-child {
        border-bottom: none;
      }
      .annotation-header {
        display: flex;
        align-items: center;
        gap: 6px;
        font-size: 12px;
        color: #909399;
        .annotation-author {
          color: #303133;
          font-weight: 500;
        }
        .annotation-time {
          margin-left: auto;
        }
      }
      .annotation-content {
        margin-top: 6px;
        color: #303133;
        line-height: 1.6;
      }
    }
  }
  .annotate-preview {
    background: #f5f7fa;
    border-radius: 4px;
    padding: 8px;
  }
}
</style>
