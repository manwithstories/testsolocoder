<template>
  <div class="upload-work">
    <div class="page-header">
      <h2>上传作品</h2>
    </div>
    
    <el-form 
      ref="formRef"
      :model="form" 
      :rules="rules"
      label-width="100px"
      style="max-width: 600px;"
    >
      <el-form-item label="音频文件" prop="audioFile">
        <el-upload
          ref="uploadRef"
          class="audio-uploader"
          drag
          :auto-upload="false"
          :limit="1"
          :file-list="fileList"
          :before-upload="beforeAudioUpload"
          :on-change="handleFileChange"
          :on-remove="handleFileRemove"
          accept=".mp3,.wav,.flac,.aac,.ogg,.m4a"
        >
          <el-icon class="el-icon--upload"><UploadFilled /></el-icon>
          <div class="el-upload__text">
            将音频文件拖到此处，或<em>点击上传</em>
          </div>
          <template #tip>
            <div class="el-upload__tip">
              支持 MP3、WAV、FLAC、AAC、OGG、M4A 格式，文件大小不超过 50MB
            </div>
          </template>
        </el-upload>
      </el-form-item>
      
      <el-form-item label="作品标题" prop="title">
        <el-input v-model="form.title" placeholder="请输入作品标题" />
      </el-form-item>
      
      <el-form-item label="艺术家" prop="artist_name">
        <el-input v-model="form.artist_name" placeholder="请输入艺术家名称" />
      </el-form-item>
      
      <el-form-item label="封面">
        <el-upload
          class="cover-uploader"
          :auto-upload="false"
          :show-file-list="false"
          :on-change="handleCoverChange"
          accept=".jpg,.jpeg,.png,.gif,.webp"
        >
          <el-image 
            v-if="form.cover_url" 
            :src="form.cover_url" 
            class="cover-preview"
            fit="cover"
          />
          <div v-else class="cover-placeholder">
            <el-icon :size="40"><Plus /></el-icon>
          </div>
        </el-upload>
      </el-form-item>
      
      <el-form-item label="风格">
        <el-select v-model="form.genre" placeholder="选择风格" clearable>
          <el-option label="流行" value="流行" />
          <el-option label="摇滚" value="摇滚" />
          <el-option label="电子" value="电子" />
          <el-option label="民谣" value="民谣" />
          <el-option label="说唱" value="说唱" />
          <el-option label="R&B" value="R&B" />
          <el-option label="古典" value="古典" />
          <el-option label="爵士" value="爵士" />
        </el-select>
      </el-form-item>
      
      <el-form-item label="标签">
        <el-select
          v-model="form.tags"
          multiple
          filterable
          allow-create
          default-first-option
          placeholder="输入标签，回车添加"
        >
          <el-option 
            v-for="tag in tags" 
            :key="tag.id" 
            :label="tag.name" 
            :value="tag.name" 
          />
        </el-select>
      </el-form-item>
      
      <el-form-item label="描述">
        <el-input 
          v-model="form.description" 
          type="textarea" 
          :rows="4"
          placeholder="作品描述"
        />
      </el-form-item>
      
      <el-form-item label="歌词">
        <el-input 
          v-model="form.lyrics" 
          type="textarea" 
          :rows="6"
          placeholder="输入歌词"
        />
      </el-form-item>
      
      <el-form-item label="版权信息">
        <el-collapse>
          <el-collapse-item title="版权信息" name="1">
            <el-form-item label="版权类型" label-width="100px">
              <el-select v-model="form.copyright_type" placeholder="选择版权类型" clearable>
                <el-option label="原创" value="original" />
                <el-option label="翻唱" value="cover" />
                <el-option label="改编" value="adaptation" />
              </el-select>
            </el-form-item>
            <el-form-item label="版权方" label-width="100px">
              <el-input v-model="form.copyright_owner" placeholder="版权方名称" />
            </el-form-item>
            <el-form-item label="授权类型" label-width="100px">
              <el-select v-model="form.license_type" placeholder="选择授权类型" clearable>
                <el-option label="独家授权" value="exclusive" />
                <el-option label="非独家授权" value="non-exclusive" />
              </el-select>
            </el-form-item>
          </el-collapse-item>
        </el-collapse>
      </el-form-item>
      
      <el-form-item>
        <el-button type="primary" :loading="uploading" @click="submit">
          提交
        </el-button>
        <el-button @click="reset">重置</el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, type FormInstance, type FormRules, type UploadFile } from 'element-plus'
import { workApi } from '@/api/work'
import { useUserStore } from '@/stores/user'
import type { Tag } from '@/api/work'

const router = useRouter()
const userStore = useUserStore()
const formRef = ref<FormInstance>()
const uploadRef = ref()
const uploading = ref(false)
const fileList = ref<UploadFile[]>([])
const audioFile = ref<File | null>(null)
const coverFile = ref<File | null>(null)
const tags = ref<Tag[]>([])

const form = reactive({
  title: '',
  artist_name: '',
  cover_url: '',
  genre: '',
  tags: [] as string[],
  description: '',
  lyrics: '',
  copyright_type: '',
  copyright_owner: '',
  license_type: ''
})

const rules: FormRules = {
  title: [{ required: true, message: '请输入作品标题', trigger: 'blur' }],
  artist_name: [{ required: true, message: '请输入艺术家名称', trigger: 'blur' }]
}

onMounted(() => {
  loadTags()
})

async function loadTags() {
  try {
    const res = await workApi.listTags()
    tags.value = res as unknown as Tag[]
  } catch (e) {
    console.error(e)
  }
}

function beforeAudioUpload(file: File) {
  const isAudio = ['audio/mpeg', 'audio/wav', 'audio/flac', 'audio/aac', 'audio/ogg', 'audio/mp4'].includes(file.type)
  const isLt50M = file.size / 1024 / 1024 < 50
  
  if (!isAudio) {
    ElMessage.error('只能上传音频文件!')
    return false
  }
  if (!isLt50M) {
    ElMessage.error('音频文件大小不能超过 50MB!')
    return false
  }
  return true
}

function handleFileChange(file: UploadFile) {
  if (file.raw) {
    audioFile.value = file.raw
  }
}

function handleFileRemove() {
  audioFile.value = null
}

function handleCoverChange(file: UploadFile) {
  if (file.raw) {
    coverFile.value = file.raw
    const reader = new FileReader()
    reader.onload = (e) => {
      form.cover_url = e.target?.result as string
    }
    reader.readAsDataURL(file.raw)
  }
}

async function submit() {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    
    if (!audioFile.value) {
      ElMessage.warning('请上传音频文件')
      return
    }
    
    uploading.value = true
    try {
      const formData = new FormData()
      formData.append('audio', audioFile.value)
      formData.append('title', form.title)
      formData.append('artist_name', form.artist_name)
      if (form.genre) formData.append('genre', form.genre)
      if (form.description) formData.append('description', form.description)
      if (form.lyrics) formData.append('lyrics', form.lyrics)
      if (coverFile.value) formData.append('cover', coverFile.value)
      if (form.tags.length) formData.append('tags', JSON.stringify(form.tags))
      
      await workApi.upload(formData)
      ElMessage.success('上传成功')
      router.push('/user/works')
    } catch (e) {
      console.error(e)
    } finally {
      uploading.value = false
    }
  })
}

function reset() {
  formRef.value?.resetFields()
  fileList.value = []
  audioFile.value = null
  coverFile.value = null
  form.cover_url = ''
}
</script>

<style scoped lang="scss">
.upload-work {
  .audio-uploader {
    :deep(.el-upload) {
      width: 100%;
    }
    
    :deep(.el-upload-dragger) {
      width: 100%;
    }
  }
  
  .cover-uploader {
    .cover-preview {
      width: 100px;
      height: 100px;
      border-radius: 4px;
      overflow: hidden;
    }
    
    .cover-placeholder {
      width: 100px;
      height: 100px;
      background: #f5f7fa;
      border: 1px dashed var(--border-color);
      border-radius: 4px;
      display: flex;
      align-items: center;
      justify-content: center;
      color: #c0c4cc;
    }
  }
}
</style>
