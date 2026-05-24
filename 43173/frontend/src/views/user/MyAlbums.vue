<template>
  <div class="my-albums">
    <div class="page-header">
      <h2>我的专辑</h2>
      <el-button type="primary" @click="createAlbum">
        <el-icon><Plus /></el-icon>
        创建专辑
      </el-button>
    </div>
    
    <div class="albums-list" v-loading="loading">
      <div 
        v-for="album in albums" 
        :key="album.id"
        class="album-item"
        @click="editAlbum(album)"
      >
        <el-image :src="album.cover_url" class="cover" fit="cover">
          <template #error>
            <div class="image-placeholder">
              <el-icon :size="20"><Headset /></el-icon>
            </div>
          </template>
        </el-image>
        <div class="info">
          <div class="title text-ellipsis">{{ album.title }}</div>
          <div class="meta">
            <span>{{ album.work_count || 0 }} 首作品</span>
            <span>{{ formatDate(album.release_date) }}</span>
          </div>
        </div>
      </div>
      
      <el-empty v-if="albums.length === 0 && !loading" description="暂无专辑" />
    </div>
    
    <el-dialog v-model="dialog.visible" :title="dialog.isEdit ? '编辑专辑' : '创建专辑'" width="500px">
      <el-form :model="dialog.form" label-width="80px">
        <el-form-item label="专辑名称">
          <el-input v-model="dialog.form.title" />
        </el-form-item>
        <el-form-item label="专辑封面">
          <el-upload
            class="cover-uploader"
            :auto-upload="false"
            :show-file-list="false"
            :on-change="handleCoverChange"
            accept=".jpg,.jpeg,.png,.gif,.webp"
          >
            <el-image 
              v-if="dialog.form.cover_url" 
              :src="dialog.form.cover_url" 
              class="cover-preview"
              fit="cover"
            />
            <div v-else class="cover-placeholder">
              <el-icon :size="30"><Plus /></el-icon>
            </div>
          </el-upload>
        </el-form-item>
        <el-form-item label="发行日期">
          <el-date-picker v-model="dialog.form.release_date" type="date" />
        </el-form-item>
        <el-form-item label="简介">
          <el-input 
            v-model="dialog.form.description" 
            type="textarea" 
            :rows="3"
          />
        </el-form-item>
        <el-form-item label="添加作品">
          <el-select 
            v-model="dialog.form.work_ids" 
            multiple 
            filterable 
            placeholder="选择要添加的作品"
            style="width: 100%;"
          >
            <el-option 
              v-for="work in works" 
              :key="work.id" 
              :label="work.title" 
              :value="work.id" 
            />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialog.visible = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="saveAlbum">
          {{ dialog.isEdit ? '保存' : '创建' }}
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { workApi } from '@/api/work'
import { useUserStore } from '@/stores/user'

const userStore = useUserStore()
const loading = ref(false)
const saving = ref(false)
const albums = ref<any[]>([])
const works = ref<any[]>([])

const dialog = reactive({
  visible: false,
  isEdit: false,
  albumId: 0,
  form: {
    title: '',
    cover_url: '',
    release_date: '',
    description: '',
    work_ids: [] as number[]
  }
})

onMounted(() => {
  loadAlbums()
  loadWorks()
})

async function loadAlbums() {
  if (!userStore.user?.artist_info) return
  
  loading.value = true
  try {
    const res = await workApi.getAlbums(userStore.user.artist_info.id)
    albums.value = res as unknown as any[]
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

async function loadWorks() {
  if (!userStore.user?.artist_info) return
  
  try {
    const res = await workApi.getByArtist(userStore.user.artist_info.id, {
      page: 1,
      page_size: 100,
      status: 2
    })
    works.value = res.list
  } catch (e) {
    console.error(e)
  }
}

function createAlbum() {
  dialog.isEdit = false
  dialog.albumId = 0
  dialog.form = {
    title: '',
    cover_url: '',
    release_date: '',
    description: '',
    work_ids: []
  }
  dialog.visible = true
}

function editAlbum(album: any) {
  dialog.isEdit = true
  dialog.albumId = album.id
  dialog.form = {
    title: album.title,
    cover_url: album.cover_url,
    release_date: album.release_date,
    description: album.description,
    work_ids: album.works?.map((w: any) => w.id) || []
  }
  dialog.visible = true
}

function handleCoverChange(file: any) {
  if (file.raw) {
    const reader = new FileReader()
    reader.onload = (e) => {
      dialog.form.cover_url = e.target?.result as string
    }
    reader.readAsDataURL(file.raw)
  }
}

async function saveAlbum() {
  if (!dialog.form.title) {
    ElMessage.warning('请输入专辑名称')
    return
  }
  
  saving.value = true
  try {
    if (dialog.isEdit) {
      await workApi.updateAlbum(dialog.albumId, dialog.form)
      ElMessage.success('保存成功')
    } else {
      await workApi.createAlbum(dialog.form)
      ElMessage.success('创建成功')
    }
    dialog.visible = false
    loadAlbums()
  } catch (e) {
    console.error(e)
  } finally {
    saving.value = false
  }
}

function formatDate(date: string) {
  if (!date) return '-'
  return date.split('T')[0]
}
</script>

<style scoped lang="scss">
.my-albums {
  .albums-list {
    display: flex;
    flex-wrap: wrap;
    gap: 20px;
    
    .album-item {
      width: 200px;
      cursor: pointer;
      
      &:hover {
        .cover {
          transform: scale(1.02);
        }
        
        .title {
          color: var(--primary-color);
        }
      }
      
      .cover {
        width: 200px;
        height: 200px;
        border-radius: 8px;
        overflow: hidden;
        margin-bottom: 12px;
        transition: transform 0.3s;
      }
      
      .info {
        .title {
          font-weight: 500;
          margin-bottom: 4px;
          transition: color 0.3s;
        }
        
        .meta {
          font-size: 13px;
          color: var(--text-light);
          display: flex;
          gap: 8px;
        }
      }
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
  
  .image-placeholder {
    width: 100%;
    height: 100%;
    background: #f5f7fa;
    display: flex;
    align-items: center;
    justify-content: center;
    color: #c0c4cc;
  }
}
</style>
