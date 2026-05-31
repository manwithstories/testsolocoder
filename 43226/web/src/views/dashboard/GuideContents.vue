<template>
  <div class="dashboard-guide-contents">
    <div class="card-shadow p-20 mb-20">
      <div class="flex-between mb-20">
        <h2 class="page-title">导览内容管理</h2>
        <div class="action-buttons">
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon> 新增导览内容
          </el-button>
        </div>
      </div>

      <el-form :inline="true" :model="query" class="mb-20">
        <el-form-item>
          <el-input
            v-model="query.keyword"
            placeholder="搜索内容"
            clearable
            @keyup.enter="fetchList"
          />
        </el-form-item>
        <el-form-item label="语言">
          <el-select v-model="query.language" clearable style="width: 150px" @change="fetchList">
            <el-option label="中文" value="zh" />
            <el-option label="英文" value="en" />
            <el-option label="日文" value="ja" />
          </el-select>
        </el-form-item>
        <el-form-item label="关联藏品">
          <el-select v-model="query.collection_id" clearable style="width: 200px" @change="fetchList">
            <el-option v-for="c in collections" :key="c.id" :label="c.name" :value="c.id" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="fetchList">查询</el-button>
          <el-button @click="resetQuery">重置</el-button>
        </el-form-item>
      </el-form>

      <el-table :data="list" v-loading="loading" border>
        <el-table-column type="index" width="60" label="#" />
        <el-table-column label="关联藏品" width="150">
          <template #default="{ row }">
            <div class="flex-center gap-10">
              <img
                v-if="row.collection?.image_url"
                :src="row.collection.image_url"
                style="width: 40px; height: 40px; object-fit: cover; border-radius: 4px;"
              />
              <span>{{ row.collection?.name || '-' }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="language" label="语言" width="100">
          <template #default="{ row }">{{ languageText(row.language) }}</template>
        </el-table-column>
        <el-table-column prop="content" label="讲解内容" min-width="200" show-overflow-tooltip>
          <template #default="{ row }">
            <span v-if="row.content">{{ row.content.substring(0, 50) }}{{ row.content.length > 50 ? '...' : '' }}</span>
            <span v-else class="text-gray">暂无内容</span>
          </template>
        </el-table-column>
        <el-table-column label="音频" width="100">
          <template #default="{ row }">
            <el-tag v-if="row.audio_url" type="success" size="small" effect="light">
              <el-icon size="12"><VideoPlay /></el-icon> 有音频
            </el-tag>
            <el-tag v-else type="info" size="small" effect="light">无音频</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="sort_order" label="排序" width="80" />
        <el-table-column prop="created_at" label="创建时间" width="160">
          <template #default="{ row }">{{ formatDateTime(row.created_at) }}</template>
        </el-table-column>
        <el-table-column label="操作" width="180" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" size="small" link @click="handlePreview(row)">预览</el-button>
            <el-button type="primary" size="small" link @click="handleEdit(row)">编辑</el-button>
            <el-button type="danger" size="small" link @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination mt-20">
        <el-pagination
          v-model:current-page="query.page"
          v-model:page-size="query.page_size"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          :page-sizes="[10, 20, 50, 100]"
          @size-change="fetchList"
          @current-change="fetchList"
        />
      </div>
    </div>

    <el-dialog v-model="showDialog" :title="isEdit ? '编辑导览内容' : '新增导览内容'" width="700px">
      <el-form :model="form" :rules="rules" ref="formRef" label-width="100px">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="关联藏品" prop="collection_id">
              <el-select v-model="form.collection_id" style="width: 100%" filterable>
                <el-option v-for="c in collections" :key="c.id" :label="c.name" :value="c.id">
                  <span>{{ c.name }}</span>
                  <span style="float: right; color: #8492a6; font-size: 13px;">{{ c.code }}</span>
                </el-option>
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="语言" prop="language">
              <el-select v-model="form.language" style="width: 100%">
                <el-option label="中文" value="zh" />
                <el-option label="英文" value="en" />
                <el-option label="日文" value="ja" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="排序" prop="sort_order">
              <el-input-number v-model="form.sort_order" :min="0" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="关联展览">
              <el-select v-model="form.exhibition_id" clearable style="width: 100%" filterable>
                <el-option v-for="e in exhibitions" :key="e.id" :label="e.title" :value="e.id" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="音频文件">
          <el-upload
            class="audio-uploader"
            action="/api/v1/uploads/audio"
            :headers="{ Authorization: `Bearer ${userStore.token}` }"
            :show-file-list="false"
            :accept="'audio/*'"
            :on-success="handleAudioUploadSuccess"
            :before-upload="beforeAudioUpload"
          >
            <div v-if="form.audio_url" class="audio-preview">
              <el-icon size="24"><Microphone /></el-icon>
              <span>已上传音频</span>
              <audio :src="form.audio_url" controls style="width: 200px; margin-left: 10px;"></audio>
            </div>
            <div v-else class="upload-placeholder">
              <el-icon size="24"><UploadFilled /></el-icon>
              <div>点击上传音频</div>
              <div class="el-upload__tip">支持 mp3, wav, ogg 格式</div>
            </div>
          </el-upload>
          <el-button v-if="form.audio_url" type="danger" text @click="removeAudio">删除音频</el-button>
        </el-form-item>
        <el-form-item label="讲解内容" prop="content">
          <el-input v-model="form.content" type="textarea" :rows="8" placeholder="请输入讲解内容..." />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showDialog = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="showPreview" title="导览内容预览" width="600px">
      <div class="preview-content">
        <div class="preview-header mb-20">
          <div class="flex-between mb-10">
            <h3>{{ previewData.collection?.name }}</h3>
            <el-tag :type="languageType(previewData.language)">{{ languageText(previewData.language) }}</el-tag>
          </div>
          <div class="text-gray text-sm">编号：{{ previewData.collection?.code }}</div>
        </div>
        <div v-if="previewData.audio_url" class="preview-audio mb-20">
          <audio :src="previewData.audio_url" controls style="width: 100%;"></audio>
        </div>
        <div class="preview-text">
          <h4 class="mb-10">讲解内容</h4>
          <div class="content-text">{{ previewData.content || '暂无内容' }}</div>
        </div>
      </div>
      <template #footer>
        <el-button @click="showPreview = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules, type UploadRawFile } from 'element-plus'
import { Plus, VideoPlay, Microphone, UploadFilled } from '@element-plus/icons-vue'
import * as guideApi from '@/api/guide'
import * as collectionApi from '@/api/collection'
import * as exhibitionApi from '@/api/exhibition'
import type { GuideContent, Collection, Exhibition } from '@/types'
import { useUserStore } from '@/stores/user'
import dayjs from 'dayjs'

const userStore = useUserStore()
const loading = ref(false)
const submitting = ref(false)
const list = ref<GuideContent[]>([])
const total = ref(0)
const collections = ref<Collection[]>([])
const exhibitions = ref<Exhibition[]>([])
const showDialog = ref(false)
const showPreview = ref(false)
const isEdit = ref(false)
const formRef = ref<FormInstance>()

const previewData = ref<GuideContent>({
  id: 0,
  collection_id: 0,
  language: 'zh',
  content: '',
  audio_url: '',
  sort_order: 0,
  created_at: ''
})

const query = reactive({
  page: 1,
  page_size: 10,
  keyword: '',
  collection_id: undefined as number | undefined,
  language: ''
})

const form = reactive({
  id: 0,
  collection_id: 0,
  exhibition_id: undefined as number | undefined,
  language: 'zh',
  content: '',
  audio_url: '',
  sort_order: 0
})

const rules: FormRules = {
  collection_id: [{ required: true, message: '请选择关联藏品', trigger: 'change' }],
  language: [{ required: true, message: '请选择语言', trigger: 'change' }],
  content: [{ required: true, message: '请输入讲解内容', trigger: 'blur' }]
}

const languageText = (lang: string) => {
  const map: Record<string, string> = {
    zh: '中文',
    en: '英文',
    ja: '日文'
  }
  return map[lang] || lang
}

const languageType = (lang: string) => {
  const map: Record<string, string> = {
    zh: 'success',
    en: 'primary',
    ja: 'warning'
  }
  return map[lang] || 'info'
}

const formatDateTime = (date: string) => {
  return dayjs(date).format('YYYY-MM-DD HH:mm')
}

const beforeAudioUpload = (file: UploadRawFile) => {
  const isAudio = file.type.startsWith('audio/')
  const isLt10M = file.size / 1024 / 1024 < 10
  if (!isAudio) {
    ElMessage.error('只能上传音频文件！')
  }
  if (!isLt10M) {
    ElMessage.error('音频大小不能超过 10MB！')
  }
  return isAudio && isLt10M
}

const handleAudioUploadSuccess = (res: any) => {
  form.audio_url = res.data.url
  ElMessage.success('音频上传成功')
}

const removeAudio = () => {
  ElMessageBox.confirm('确定要删除该音频吗？', '提示', { type: 'warning' }).then(() => {
    form.audio_url = ''
    ElMessage.success('删除成功')
  }).catch(() => {})
}

const fetchList = async () => {
  try {
    loading.value = true
    const params: any = {
      page: query.page,
      page_size: query.page_size
    }
    if (query.collection_id) params.collection_id = query.collection_id
    if (query.language) params.language = query.language
    const res = await guideApi.listGuideContents(params)
    list.value = res.data
    total.value = res.data.length
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

const fetchCollections = async () => {
  try {
    const res = await collectionApi.listCollections({ page: 1, page_size: 1000, status: 'active' })
    collections.value = res.data.list
  } catch (e) {
    console.error(e)
  }
}

const fetchExhibitions = async () => {
  try {
    const res = await exhibitionApi.listExhibitions({ page: 1, page_size: 1000, status: 'published' })
    exhibitions.value = res.data.list
  } catch (e) {
    console.error(e)
  }
}

const resetQuery = () => {
  query.page = 1
  query.keyword = ''
  query.collection_id = undefined
  query.language = ''
  fetchList()
}

const handleAdd = () => {
  isEdit.value = false
  Object.assign(form, {
    id: 0,
    collection_id: collections.value[0]?.id || 0,
    exhibition_id: undefined,
    language: 'zh',
    content: '',
    audio_url: '',
    sort_order: 0
  })
  showDialog.value = true
}

const handleEdit = (row: GuideContent) => {
  isEdit.value = true
  Object.assign(form, {
    id: row.id,
    collection_id: row.collection_id,
    exhibition_id: row.exhibition_id,
    language: row.language,
    content: row.content,
    audio_url: row.audio_url || '',
    sort_order: row.sort_order
  })
  showDialog.value = true
}

const handlePreview = (row: GuideContent) => {
  previewData.value = row
  showPreview.value = true
}

const handleDelete = (row: GuideContent) => {
  ElMessageBox.confirm('确定要删除该导览内容吗？', '提示', { type: 'warning' }).then(async () => {
    await guideApi.deleteGuideContent(row.id)
    ElMessage.success('删除成功')
    fetchList()
  }).catch(() => {})
}

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate()
  try {
    submitting.value = true
    const data = {
      collection_id: form.collection_id,
      exhibition_id: form.exhibition_id,
      language: form.language,
      content: form.content,
      audio_url: form.audio_url || undefined,
      sort_order: form.sort_order
    }
    if (isEdit.value) {
      await guideApi.updateGuideContent(form.id, data)
      ElMessage.success('更新成功')
    } else {
      await guideApi.createGuideContent(data)
      ElMessage.success('创建成功')
    }
    showDialog.value = false
    fetchList()
  } catch (e) {
    console.error(e)
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  fetchCollections()
  fetchExhibitions()
  fetchList()
})
</script>

<style scoped lang="scss">
.dashboard-guide-contents {
  .page-title {
    margin: 0;
    font-size: 20px;
  }

  .audio-uploader {
    width: 100%;

    :deep(.el-upload) {
      width: 100%;
    }

    .audio-preview {
      display: flex;
      align-items: center;
      padding: 15px;
      border: 1px dashed #d9d9d9;
      border-radius: 6px;
      background: #fafafa;
    }

    .upload-placeholder {
      text-align: center;
      padding: 30px;
      border: 1px dashed #d9d9d9;
      border-radius: 6px;
      color: #8c939d;
      background: #fafafa;
      cursor: pointer;
      transition: border-color 0.2s;

      &:hover {
        border-color: #409eff;
      }

      .el-upload__tip {
        font-size: 12px;
        color: #999;
        margin-top: 5px;
      }
    }
  }

  .preview-content {
    .preview-header {
      h3 {
        margin: 0;
        font-size: 18px;
      }

      .text-gray {
        color: #909399;
      }

      .text-sm {
        font-size: 12px;
      }
    }

    .preview-audio {
      padding: 15px;
      background: #f5f7fa;
      border-radius: 6px;
    }

    .preview-text {
      h4 {
        margin: 0;
        font-size: 14px;
        color: #303133;
      }

      .content-text {
        line-height: 1.8;
        color: #606266;
        white-space: pre-wrap;
      }
    }
  }

  .text-gray {
    color: #909399;
  }
}
</style>
