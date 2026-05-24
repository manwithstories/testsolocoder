<template>
  <div class="playlist-detail" v-loading="loading">
    <div v-if="playlist" class="playlist-content">
      <div class="playlist-header">
        <el-image :src="playlist.cover_url" class="cover" fit="cover">
          <template #error>
            <div class="image-placeholder">
              <el-icon :size="64"><List /></el-icon>
            </div>
          </template>
        </el-image>
        
        <div class="playlist-info">
          <h1 class="title">{{ playlist.name }}</h1>
          <div class="creator" @click="goToUser">
            <el-avatar :size="32" :src="playlist.user?.avatar">
              {{ playlist.user?.nickname?.charAt(0) || 'U' }}
            </el-avatar>
            <span>{{ playlist.user?.nickname }}</span>
          </div>
          
          <div class="stats">
            <span>{{ playlist.work_count }} 首歌曲</span>
            <span>{{ formatCount(playlist.play_count) }} 播放</span>
            <span>{{ formatCount(playlist.like_count) }} 收藏</span>
          </div>
          
          <p class="description" v-if="playlist.description">{{ playlist.description }}</p>
          
          <div class="actions" v-if="userStore.isLoggedIn">
            <el-button type="primary" size="large" @click="playAll">
              <el-icon><VideoPlay /></el-icon>
              播放全部
            </el-button>
            <el-button size="large" @click="toggleLike">
              <el-icon><Star /></el-icon>
              {{ isLiked ? '已收藏' : '收藏' }}
            </el-button>
          </div>
        </div>
      </div>
      
      <div class="playlist-body">
        <h3>歌曲列表</h3>
        <div class="works-list" v-loading="worksLoading">
          <div 
            v-for="(work, index) in works" 
            :key="work.id"
            class="work-item"
            @click="goToWork(work.id)"
          >
            <div class="index">{{ index + 1 }}</div>
            <el-image :src="work.cover_url" class="cover" fit="cover">
              <template #error>
                <div class="image-placeholder">
                  <el-icon :size="16"><Headset /></el-icon>
                </div>
              </template>
            </el-image>
            <div class="info">
              <div class="title text-ellipsis">{{ work.title }}</div>
              <div class="artist text-ellipsis">{{ work.artist_name }}</div>
            </div>
            <div class="actions" @click.stop>
              <el-button text v-if="isOwner" @click="removeWork(work.id)">
                <el-icon><Delete /></el-icon>
              </el-button>
            </div>
          </div>
          
          <el-empty v-if="works.length === 0 && !worksLoading" description="暂无歌曲" />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { communityApi } from '@/api/community'
import { useUserStore } from '@/stores/user'
import type { Playlist } from '@/api/community'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const loading = ref(false)
const worksLoading = ref(false)
const playlist = ref<Playlist | null>(null)
const works = ref<any[]>([])
const isLiked = ref(false)

const isOwner = computed(() => playlist.value?.user_id === userStore.user?.id)

onMounted(async () => {
  await loadPlaylist()
  loadWorks()
})

async function loadPlaylist() {
  loading.value = true
  try {
    const id = parseInt(route.params.id as string)
    playlist.value = await communityApi.getPlaylistById(id)
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

async function loadWorks() {
  if (!playlist.value) return
  worksLoading.value = true
  try {
    works.value = playlist.value.works || []
  } catch (e) {
    console.error(e)
  } finally {
    worksLoading.value = false
  }
}

function playAll() {
  ElMessage.info('播放全部功能开发中')
}

function toggleLike() {
  isLiked.value = !isLiked.value
  ElMessage.info(isLiked.value ? '已收藏' : '已取消收藏')
}

async function removeWork(workId: number) {
  if (!playlist.value) return
  
  try {
    await ElMessageBox.confirm('确定要从歌单中移除这首歌曲吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    await communityApi.removeWorkFromPlaylist(playlist.value.id, workId)
    ElMessage.success('已移除')
    loadPlaylist()
  } catch (e) {
    if (e !== 'cancel') {
      console.error(e)
    }
  }
}

function goToWork(id: number) {
  router.push(`/work/${id}`)
}

function goToUser() {
  if (playlist.value?.user_id) {
    router.push(`/artist/${playlist.value.user_id}`)
  }
}

function formatCount(count: number) {
  if (count >= 10000) {
    return (count / 10000).toFixed(1) + 'w'
  }
  return count
}
</script>

<style scoped lang="scss">
.playlist-detail {
  .playlist-content {
    .playlist-header {
      display: flex;
      gap: 30px;
      margin-bottom: 30px;
      
      .cover {
        width: 240px;
        height: 240px;
        border-radius: 8px;
        overflow: hidden;
        flex-shrink: 0;
      }
      
      .playlist-info {
        flex: 1;
        display: flex;
        flex-direction: column;
        
        .title {
          font-size: 28px;
          font-weight: 600;
          margin: 0 0 12px 0;
        }
        
        .creator {
          display: flex;
          align-items: center;
          gap: 8px;
          margin-bottom: 16px;
          cursor: pointer;
          
          &:hover {
            color: var(--primary-color);
          }
        }
        
        .stats {
          display: flex;
          gap: 20px;
          color: var(--text-light);
          font-size: 14px;
          margin-bottom: 16px;
        }
        
        .description {
          line-height: 1.6;
          margin-bottom: 20px;
        }
        
        .actions {
          display: flex;
          gap: 12px;
          margin-top: auto;
        }
      }
    }
    
    .playlist-body {
      .works-list {
        .work-item {
          display: flex;
          align-items: center;
          gap: 12px;
          padding: 12px;
          border-radius: 4px;
          cursor: pointer;
          transition: background 0.3s;
          
          &:hover {
            background: rgba(64, 158, 255, 0.05);
          }
          
          .index {
            width: 30px;
            text-align: center;
            color: var(--text-light);
          }
          
          .cover {
            width: 50px;
            height: 50px;
            border-radius: 4px;
            overflow: hidden;
          }
          
          .info {
            flex: 1;
            min-width: 0;
            
            .title {
              font-size: 14px;
              margin-bottom: 2px;
            }
            
            .artist {
              font-size: 12px;
              color: var(--text-light);
            }
          }
          
          .actions {
            display: none;
          }
          
          &:hover .actions {
            display: block;
          }
        }
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
