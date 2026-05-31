<template>
  <div class="translation-workbench">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>翻译工作台</span>
          <div>
            <el-select v-model="selectedProject" placeholder="选择项目" style="width: 300px" @change="loadSegments">
              <el-option
                v-for="p in myProjects"
                :key="p.id"
                :label="p.title"
                :value="p.id"
              />
            </el-select>
          </div>
        </div>
      </template>

      <div v-if="!selectedProject" class="empty-tip">
        <el-empty description="请选择一个进行中的项目开始翻译" />
      </div>

      <div v-else class="workbench-content">
        <div class="segment-list">
          <div
            v-for="(segment, index) in segments"
            :key="segment.id"
            class="segment-item"
            :class="{ active: activeSegment === index }"
            @click="activeSegment = index"
          >
            <div class="segment-header">
              <span class="segment-index">#{{ index + 1 }}</span>
              <el-tag v-if="segment.is_memory" type="success" size="small">记忆</el-tag>
              <el-tag v-else-if="segment.status === 'translated'" type="success" size="small">已译</el-tag>
              <el-tag v-else type="info" size="small">待译</el-tag>
            </div>
            <div class="segment-source">{{ segment.source_text }}</div>
          </div>
        </div>

        <div class="translation-area" v-if="activeSegment >= 0 && segments[activeSegment]">
          <div class="source-panel">
            <div class="panel-header">原文</div>
            <div class="source-text">{{ currentSegment?.source_text }}</div>
          </div>

          <div class="memory-panel" v-if="memorySuggestions.length > 0">
            <div class="panel-header">翻译记忆</div>
            <div
              v-for="(sug, i) in memorySuggestions"
              :key="i"
              class="memory-item"
              @click="applySuggestion(sug)"
            >
              <div class="memory-source">{{ sug.source_text }}</div>
              <div class="memory-target">{{ sug.translated_text }}</div>
              <div class="memory-meta">
                <span>相似度: {{ (sug.similarity * 100).toFixed(0) }}%</span>
                <el-button type="primary" link size="small">使用</el-button>
              </div>
            </div>
          </div>

          <div class="target-panel">
            <div class="panel-header">译文</div>
            <el-input
              v-model="targetText"
              type="textarea"
              :rows="6"
              placeholder="在此输入译文..."
              @input="loadMemorySuggestions"
            />
            <div class="action-bar">
              <el-button type="primary" @click="saveTranslation" :loading="saving">
                保存并下一段
              </el-button>
              <el-button @click="nextSegment">跳过</el-button>
            </div>
          </div>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/stores/user'
import { listProjects } from '@/api/project'
import { listSegments, updateSegmentTranslation, getMemorySuggestions } from '@/api/document'

const userStore = useUserStore()

const myProjects = ref<any[]>([])
const selectedProject = ref<number | null>(null)
const segments = ref<any[]>([])
const activeSegment = ref(0)
const targetText = ref('')
const memorySuggestions = ref<any[]>([])
const saving = ref(false)

const currentSegment = computed(() => segments.value[activeSegment.value])

async function loadProjects() {
  try {
    const res = await listProjects({ status: 'in_progress' })
    if (Array.isArray(res)) {
      myProjects.value = res
    } else {
      myProjects.value = res?.list || []
    }
  } catch (e) {
    console.error(e)
  }
}

async function loadSegments() {
  if (!selectedProject.value) return
  try {
    const res = await listSegments(selectedProject.value)
    segments.value = res || []
    activeSegment.value = 0
    targetText.value = currentSegment.value?.translated_text || ''
    loadMemorySuggestions()
  } catch (e) {
    console.error(e)
  }
}

async function loadMemorySuggestions() {
  if (!currentSegment.value) return
  const project = myProjects.value.find(p => p.id === selectedProject.value)
  if (!project) return

  try {
    const res = await getMemorySuggestions(
      currentSegment.value.source_text,
      project.source_lang,
      project.target_lang
    )
    memorySuggestions.value = res || []
  } catch (e) {
    console.error(e)
  }
}

function applySuggestion(sug: any) {
  targetText.value = sug.translated_text
}

async function saveTranslation() {
  if (!currentSegment.value || !targetText.value.trim()) {
    ElMessage.warning('请输入译文')
    return
  }

  saving.value = true
  try {
    await updateSegmentTranslation(currentSegment.value.id, targetText.value)
    ElMessage.success('保存成功')
    nextSegment()
  } catch (e: any) {
    ElMessage.error(e.message || '保存失败')
  } finally {
    saving.value = false
  }
}

function nextSegment() {
  if (activeSegment.value < segments.value.length - 1) {
    activeSegment.value++
    targetText.value = currentSegment.value?.translated_text || ''
    memorySuggestions.value = []
  } else {
    ElMessage.success('已完成所有片段')
  }
}

watch(currentSegment, (seg) => {
  if (seg) {
    targetText.value = seg.translated_text || ''
  }
})

onMounted(loadProjects)
</script>

<style lang="scss" scoped>
.translation-workbench {
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .empty-tip {
    padding: 60px 0;
  }

  .workbench-content {
    display: flex;
    gap: 16px;
    height: 600px;
  }

  .segment-list {
    width: 300px;
    overflow-y: auto;
    border-right: 1px solid #e4e7ed;
    padding-right: 8px;
  }

  .segment-item {
    padding: 12px;
    border: 1px solid #e4e7ed;
    border-radius: 4px;
    margin-bottom: 8px;
    cursor: pointer;
    transition: all 0.2s;

    &.active {
      border-color: #409EFF;
      background: #ecf5ff;
    }

    .segment-header {
      display: flex;
      gap: 8px;
      align-items: center;
      margin-bottom: 4px;

      .segment-index {
        font-weight: bold;
        color: #909399;
      }
    }

    .segment-source {
      font-size: 13px;
      color: #606266;
      overflow: hidden;
      text-overflow: ellipsis;
      display: -webkit-box;
      -webkit-line-clamp: 2;
      -webkit-box-orient: vertical;
    }
  }

  .translation-area {
    flex: 1;
    display: flex;
    flex-direction: column;
    gap: 16px;
    overflow-y: auto;
  }

  .panel-header {
    font-weight: 500;
    color: #303133;
    margin-bottom: 8px;
  }

  .source-panel {
    padding: 12px;
    background: #f5f7fa;
    border-radius: 4px;

    .source-text {
      font-size: 14px;
      line-height: 1.6;
      white-space: pre-wrap;
    }
  }

  .memory-panel {
    .memory-item {
      padding: 12px;
      border: 1px solid #e4e7ed;
      border-radius: 4px;
      margin-bottom: 8px;
      cursor: pointer;

      &:hover {
        border-color: #409EFF;
      }

      .memory-source {
        font-size: 13px;
        color: #909399;
        margin-bottom: 4px;
      }

      .memory-target {
        font-size: 14px;
        color: #303133;
        margin-bottom: 8px;
      }

      .memory-meta {
        display: flex;
        justify-content: space-between;
        align-items: center;
        font-size: 12px;
        color: #909399;
      }
    }
  }

  .target-panel {
    flex: 1;
    display: flex;
    flex-direction: column;

    .action-bar {
      margin-top: 12px;
      display: flex;
      gap: 8px;
    }
  }
}
</style>
