<template>
  <div class="collection-detail page-container" v-loading="loading">
    <div v-if="collection" class="detail-content">
      <div class="main-section card-shadow">
        <div class="image-gallery">
          <img :src="collection.image_url || '/placeholder.svg'" :alt="collection.name" class="main-image" />
        </div>
        <div class="detail-info">
          <div class="flex-between">
            <h1>{{ collection.name }}</h1>
            <el-tag type="success" v-if="collection.status === 'active'">展出中</el-tag>
            <el-tag type="info" v-else-if="collection.status === 'inactive'">未展出</el-tag>
            <el-tag type="warning" v-else>修复中</el-tag>
          </div>
          <div class="code-row">
            <span class="label">藏品编号:</span>
            <span class="value">{{ collection.code }}</span>
          </div>
          <div class="info-grid">
            <div class="info-item">
              <span class="label">分类</span>
              <span class="value">{{ collection.category?.name }}</span>
            </div>
            <div class="info-item">
              <span class="label">年代</span>
              <span class="value">{{ collection.era || '-' }}</span>
            </div>
            <div class="info-item">
              <span class="label">材质</span>
              <span class="value">{{ collection.material || '-' }}</span>
            </div>
            <div class="info-item">
              <span class="label">尺寸</span>
              <span class="value">{{ collection.size || '-' }}</span>
            </div>
            <div class="info-item">
              <span class="label">来源</span>
              <span class="value">{{ collection.source || '-' }}</span>
            </div>
            <div class="info-item">
              <span class="label">保存状态</span>
              <span class="value">{{ collection.condition || '-' }}</span>
            </div>
          </div>
          <div class="tags-row" v-if="collection.tags">
            <el-tag
              v-for="tag in collection.tags.split(',').filter(t => t)"
              :key="tag"
              style="margin-right: 8px; margin-bottom: 8px;"
            >
              {{ tag }}
            </el-tag>
          </div>
          <div class="description-section">
            <h3>藏品介绍</h3>
            <p>{{ collection.description }}</p>
          </div>
          <div class="action-buttons">
            <el-button type="primary" @click="showGuide = true">
              <el-icon><Microphone /></el-icon> 查看导览
            </el-button>
            <el-button v-if="userStore.isLoggedIn" @click="applyResearch">
              <el-icon><Document /></el-icon> 申请研究使用
            </el-button>
          </div>
          <div class="stats-row">
            <span><el-icon><View /></el-icon> {{ collection.view_count }} 次浏览</span>
          </div>
        </div>
      </div>

      <div class="guide-section card-shadow p-20 mt-20" v-if="guideContents.length > 0">
        <h2>导览讲解</h2>
        <div class="lang-selector">
          <el-radio-group v-model="currentLang" @change="fetchGuideContents">
            <el-radio-button value="zh">中文</el-radio-button>
            <el-radio-button value="en">English</el-radio-button>
            <el-radio-button value="ja">日本語</el-radio-button>
          </el-radio-group>
        </div>
        <div v-for="content in filteredGuideContents" :key="content.id" class="guide-content">
          <h4>{{ content.collection?.name }}</h4>
          <p>{{ content.content }}</p>
          <div v-if="content.audio_url" class="audio-player">
            <el-icon size="20"><VideoPlay /></el-icon>
            <audio :src="content.audio_url" controls></audio>
          </div>
        </div>
      </div>
    </div>

    <el-dialog v-model="showResearchDialog" title="申请学术研究使用" width="500px">
      <el-form :model="researchForm" label-width="100px">
        <el-form-item label="藏品名称">
          <span>{{ collection?.name }}</span>
        </el-form-item>
        <el-form-item label="研究机构" prop="institution">
          <el-input v-model="researchForm.institution" placeholder="请输入您的研究机构" />
        </el-form-item>
        <el-form-item label="研究用途" prop="purpose">
          <el-input
            v-model="researchForm.purpose"
            type="textarea"
            :rows="4"
            placeholder="请详细描述您的研究用途"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showResearchDialog = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="submitResearch">提交申请</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import * as collectionApi from '@/api/collection'
import * as guideApi from '@/api/guide'
import type { Collection, GuideContent } from '@/types'
import { useUserStore } from '@/stores/user'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()
const id = Number(route.params.id)

const loading = ref(false)
const submitting = ref(false)
const collection = ref<Collection | null>(null)
const guideContents = ref<GuideContent[]>([])
const currentLang = ref('zh')
const showResearchDialog = ref(false)
const showGuide = ref(false)

const researchForm = reactive({
  institution: '',
  purpose: ''
})

const filteredGuideContents = computed(() => {
  return guideContents.value.filter(g => g.language === currentLang.value)
})

const fetchDetail = async () => {
  try {
    loading.value = true
    const res = await collectionApi.getCollection(id)
    collection.value = res.data
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

const fetchGuideContents = async () => {
  try {
    const res = await guideApi.listGuideContents({
      collection_id: id,
      language: currentLang.value
    })
    guideContents.value = res.data
  } catch (e) {
    console.error(e)
  }
}

const applyResearch = () => {
  showResearchDialog.value = true
}

const submitResearch = async () => {
  if (!researchForm.institution || !researchForm.purpose) {
    ElMessage.warning('请填写完整信息')
    return
  }
  try {
    submitting.value = true
    await guideApi.createResearchApplication({
      collection_id: id,
      institution: researchForm.institution,
      purpose: researchForm.purpose
    })
    ElMessage.success('申请提交成功，请等待管理员审批')
    showResearchDialog.value = false
  } catch (e) {
    console.error(e)
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  fetchDetail()
  fetchGuideContents()
})
</script>

<style scoped lang="scss">
.collection-detail {
  max-width: 1400px;
  margin: 0 auto;

  .main-section {
    display: flex;
    gap: 40px;
    padding: 30px;
    border-radius: 12px;

    .image-gallery {
      width: 45%;
      flex-shrink: 0;

      .main-image {
        width: 100%;
        height: 450px;
        object-fit: cover;
        border-radius: 8px;
      }
    }

    .detail-info {
      flex: 1;

      h1 {
        font-size: 30px;
        margin: 0 0 8px;
      }

      .code-row {
        margin-bottom: 20px;
        color: #909399;

        .label {
          margin-right: 8px;
        }

        .value {
          font-family: monospace;
        }
      }

      .info-grid {
        display: grid;
        grid-template-columns: 1fr 1fr;
        gap: 16px;
        margin-bottom: 20px;

        .info-item {
          display: flex;
          flex-direction: column;
          gap: 4px;

          .label {
            font-size: 13px;
            color: #909399;
          }

          .value {
            font-size: 15px;
          }
        }
      }

      .tags-row {
        margin-bottom: 20px;
      }

      .description-section {
        margin-bottom: 24px;

        h3 {
          font-size: 18px;
          margin-bottom: 12px;
        }

        p {
          line-height: 1.8;
          color: #606266;
          text-indent: 2em;
        }
      }

      .action-buttons {
        display: flex;
        gap: 12px;
        margin-bottom: 16px;
      }

      .stats-row {
        color: #909399;
        font-size: 13px;
      }
    }
  }

  .guide-section {
    border-radius: 12px;

    h2 {
      font-size: 22px;
      margin-bottom: 16px;
      padding-bottom: 12px;
      border-bottom: 1px solid #ebeef5;
    }

    .lang-selector {
      margin-bottom: 20px;
    }

    .guide-content {
      padding: 16px;
      background: #f5f7fa;
      border-radius: 8px;
      margin-bottom: 16px;

      h4 {
        margin-bottom: 8px;
        color: #303133;
      }

      p {
        line-height: 1.8;
        color: #606266;
        margin-bottom: 12px;
      }

      .audio-player {
        display: flex;
        align-items: center;
        gap: 12px;

        audio {
          flex: 1;
        }
      }
    }
  }
}

@media (max-width: 768px) {
  .main-section {
    flex-direction: column;

    .image-gallery {
      width: 100% !important;
    }
  }
}
</style>
