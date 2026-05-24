<template>
  <div class="my-playlists">
    <div class="page-header">
      <h2>我的歌单</h2>
      <el-button type="primary" @click="createPlaylist">
        <el-icon><Plus /></el-icon>
        创建歌单
      </el-button>
    </div>
    
    <div class="playlists-list" v-loading="loading">
      <div 
        v-for="playlist in playlists" 
        :key="playlist.id"
        class="playlist-item"
        @click="goToDetail(playlist)"
      >
        <el-image :src="playlist.cover_url" class="cover" fit="cover">
          <template #error>
            <div class="image-placeholder">
              <el-icon :size="20"><Headset /></el-icon>
            </div>
          </template>
        </el-image>
        <div class="info">
          <div class="title text-ellipsis">{{ playlist.name }}</div>
          <div class="meta">
            <span>{{ playlist.work_count || 0 }} 首</span>
            <span>{{ formatCount(playlist.play_count) }} 播放</span>
          </div>
        </div>
        <div class="actions" @click.stop>
          <el-button text @click="editPlaylist(playlist)">
            <el-icon><Edit /></el-icon>
          </el-button>
          <el-button text type="danger" @click="deletePlaylist(playlist)">
            <el-icon><Delete /></el-icon>
          </el-button>
        </div>
      </div>
      
      <el-empty v-if="playlists.length === 0 && !loading" description="暂无歌单" />
    </div>
    
    <el-dialog v-model="dialog.visible" :title="dialog.isEdit ? '编辑歌单' : '创建歌单'" width="400px">
      <el-form :model="dialog.form" label-width="80px">
        <el-form-item label="歌单名称">
          <el-input v-model="dialog.form.name" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input 
            v-model="dialog.form.description" 
            type="textarea" 
            :rows="3"
          />
        </el-form-item>
        <el-form-item label="公开">
          <el-switch v-model="dialog.form.is_public" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialog.visible = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="savePlaylist">
          {{ dialog.isEdit ? '保存' : '创建' }}
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { communityApi } from '@/api/community'

const router = useRouter()
const loading = ref(false)
const saving = ref(false)
const playlists = ref<any[]>([])

const dialog = reactive({
  visible: false,
  isEdit: false,
  playlistId: 0,
  form: {
    name: '',
    description: '',
    is_public: true
  }
})

onMounted(() => {
  loadPlaylists()
})

async function loadPlaylists() {
  loading.value = true
  try {
    const res = await communityApi.getMyPlaylists()
    playlists.value = res as unknown as any[]
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

function createPlaylist() {
  dialog.isEdit = false
  dialog.playlistId = 0
  dialog.form = {
    name: '',
    description: '',
    is_public: true
  }
  dialog.visible = true
}

function editPlaylist(playlist: any) {
  dialog.isEdit = true
  dialog.playlistId = playlist.id
  dialog.form = {
    name: playlist.name,
    description: playlist.description,
    is_public: playlist.is_public
  }
  dialog.visible = true
}

async function savePlaylist() {
  if (!dialog.form.name) {
    ElMessage.warning('请输入歌单名称')
    return
  }
  
  saving.value = true
  try {
    if (dialog.isEdit) {
      await communityApi.updatePlaylist(dialog.playlistId, dialog.form)
      ElMessage.success('保存成功')
    } else {
      await communityApi.createPlaylist(dialog.form)
      ElMessage.success('创建成功')
    }
    dialog.visible = false
    loadPlaylists()
  } catch (e) {
    console.error(e)
  } finally {
    saving.value = false
  }
}

async function deletePlaylist(playlist: any) {
  try {
    await ElMessageBox.confirm('确定要删除该歌单吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    await communityApi.deletePlaylist(playlist.id)
    ElMessage.success('删除成功')
    loadPlaylists()
  } catch (e) {
    if (e !== 'cancel') {
      console.error(e)
    }
  }
}

function goToDetail(playlist: any) {
  router.push(`/playlists/${playlist.id}`)
}

function formatCount(count: number) {
  if (count >= 10000) {
    return (count / 10000).toFixed(1) + 'w'
  }
  return count
}
</script>

<style scoped lang="scss">
.my-playlists {
  .playlists-list {
    display: flex;
    flex-wrap: wrap;
    gap: 20px;
    
    .playlist-item {
      width: 200px;
      cursor: pointer;
      position: relative;
      
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
      
      .actions {
        position: absolute;
        top: 8px;
        right: 8px;
        display: none;
        gap: 4px;
      }
      
      &:hover .actions {
        display: flex;
      }
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
