<template>
  <div class="playlists">
    <div class="page-header">
      <h1 class="page-title">歌单</h1>
      <el-button type="primary" v-if="userStore.isLoggedIn" @click="showCreateDialog = true">
        <el-icon><Plus /></el-icon>
        创建歌单
      </el-button>
    </div>
    
    <div class="playlists-grid" v-loading="loading">
      <div 
        v-for="playlist in playlists" 
        :key="playlist.id"
        class="playlist-item"
        @click="goToDetail(playlist.id)"
      >
        <el-image :src="playlist.cover_url" class="cover" fit="cover">
          <template #error>
            <div class="image-placeholder">
              <el-icon :size="40"><List /></el-icon>
            </div>
          </template>
        </el-image>
        <div class="info">
          <div class="title text-ellipsis">{{ playlist.name }}</div>
          <div class="meta">
            <span>{{ playlist.work_count }} 首歌曲</span>
            <span>{{ formatCount(playlist.play_count) }} 播放</span>
          </div>
        </div>
      </div>
      
      <el-empty v-if="playlists.length === 0 && !loading" description="暂无歌单" />
    </div>
    
    <div class="pagination" v-if="total > 0">
      <el-pagination
        v-model:current-page="page"
        v-model:page-size="pageSize"
        :total="total"
        :page-sizes="[12, 24, 48]"
        layout="total, sizes, prev, pager, next, jumper"
        @current-change="loadPlaylists"
        @size-change="loadPlaylists"
      />
    </div>
    
    <el-dialog v-model="showCreateDialog" title="创建歌单" width="400px">
      <el-form :model="createForm" label-width="80px">
        <el-form-item label="名称">
          <el-input v-model="createForm.name" placeholder="请输入歌单名称" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input 
            v-model="createForm.description" 
            type="textarea" 
            :rows="3"
            placeholder="请输入歌单描述" 
          />
        </el-form-item>
        <el-form-item label="公开">
          <el-switch v-model="createForm.is_public" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showCreateDialog = false">取消</el-button>
        <el-button type="primary" :loading="creating" @click="createPlaylist">
          创建
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { communityApi } from '@/api/community'
import { useUserStore } from '@/stores/user'

const router = useRouter()
const userStore = useUserStore()

const loading = ref(false)
const playlists = ref<any[]>([])
const page = ref(1)
const pageSize = ref(12)
const total = ref(0)

const showCreateDialog = ref(false)
const creating = ref(false)
const createForm = reactive({
  name: '',
  description: '',
  is_public: true
})

onMounted(() => {
  loadPlaylists()
})

async function loadPlaylists() {
  loading.value = true
  try {
    const res = await communityApi.listPlaylists({
      page: page.value,
      page_size: pageSize.value
    })
    playlists.value = res.list
    total.value = res.total
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

async function createPlaylist() {
  if (!createForm.name.trim()) {
    ElMessage.warning('请输入歌单名称')
    return
  }
  
  creating.value = true
  try {
    await communityApi.createPlaylist(createForm)
    ElMessage.success('创建成功')
    showCreateDialog.value = false
    createForm.name = ''
    createForm.description = ''
    createForm.is_public = true
    loadPlaylists()
  } catch (e) {
    console.error(e)
  } finally {
    creating.value = false
  }
}

function formatCount(count: number) {
  if (count >= 10000) {
    return (count / 10000).toFixed(1) + 'w'
  }
  return count
}

function goToDetail(id: number) {
  router.push(`/playlist/${id}`)
}
</script>

<style scoped lang="scss">
.playlists {
  .playlists-grid {
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    gap: 24px;
    
    .playlist-item {
      cursor: pointer;
      transition: transform 0.3s;
      
      &:hover {
        transform: translateY(-5px);
        
        .cover {
          box-shadow: 0 8px 24px rgba(0, 0, 0, 0.15);
        }
      }
      
      .cover {
        width: 100%;
        aspect-ratio: 1;
        border-radius: 8px;
        overflow: hidden;
        transition: box-shadow 0.3s;
      }
      
      .info {
        padding: 12px 4px;
        
        .title {
          font-weight: 500;
          margin-bottom: 4px;
        }
        
        .meta {
          font-size: 13px;
          color: var(--text-light);
          display: flex;
          gap: 12px;
        }
      }
    }
  }
  
  .pagination {
    display: flex;
    justify-content: center;
    margin-top: 30px;
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
