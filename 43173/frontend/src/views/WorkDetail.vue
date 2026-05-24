<template>
  <div class="work-detail" v-loading="loading">
    <div v-if="work" class="work-content">
      <div class="work-header">
        <el-image :src="work.cover_url" class="cover" fit="cover">
          <template #error>
            <div class="image-placeholder">
              <el-icon :size="64"><Headset /></el-icon>
            </div>
          </template>
        </el-image>
        
        <div class="work-info">
          <h1 class="title">{{ work.title }}</h1>
          <div class="artist" @click="goToArtist">
            <el-avatar :size="32" :src="work.user?.avatar">
              {{ work.artist_name?.charAt(0) }}
            </el-avatar>
            <span>{{ work.artist_name }}</span>
          </div>
          
          <div class="meta">
            <span><el-icon><Headset /></el-icon> {{ formatCount(work.play_count) }} 播放</span>
            <span><el-icon><Star /></el-icon> {{ formatCount(work.like_count) }} 收藏</span>
            <span><el-icon><ChatDotRound /></el-icon> {{ formatCount(work.comment_count) }} 评论</span>
            <span><el-icon><Share /></el-icon> {{ formatCount(work.share_count) }} 分享</span>
          </div>
          
          <div class="tags" v-if="work.tags?.length">
            <el-tag 
              v-for="tag in work.tags" 
              :key="tag.id" 
              type="info" 
              effect="plain"
              style="margin-right: 8px;"
            >
              {{ tag.name }}
            </el-tag>
          </div>
          
          <div class="actions">
            <el-button type="primary" size="large" @click="play">
              <el-icon><VideoPlay /></el-icon>
              播放
            </el-button>
            <el-button size="large" @click="toggleLike">
              <el-icon><Star /></el-icon>
              {{ isLiked ? '已收藏' : '收藏' }}
            </el-button>
            <el-button size="large" @click="showPlaylistDialog = true">
              <el-icon><Plus /></el-icon>
              添加到歌单
            </el-button>
            <el-button size="large" @click="followArtist">
              <el-icon><UserFilled /></el-icon>
              {{ isFollowing ? '已关注' : '关注' }}
            </el-button>
          </div>
        </div>
      </div>
      
      <div class="work-body">
        <el-tabs v-model="activeTab">
          <el-tab-pane label="简介" name="intro">
            <div class="intro" v-if="work.description">
              {{ work.description }}
            </div>
            <el-empty v-else description="暂无简介" />
            
            <div class="info-list" v-if="work.copyright">
              <h3>版权信息</h3>
              <div class="info-item">
                <span class="label">版权类型</span>
                <span class="value">{{ work.copyright.copyright_type }}</span>
              </div>
              <div class="info-item">
                <span class="label">版权方</span>
                <span class="value">{{ work.copyright.owner }}</span>
              </div>
              <div class="info-item">
                <span class="label">授权类型</span>
                <span class="value">{{ work.copyright.license_type }}</span>
              </div>
            </div>
          </el-tab-pane>
          
          <el-tab-pane label="歌词" name="lyrics">
            <div class="lyrics" v-if="work.lyrics">
              <pre>{{ work.lyrics }}</pre>
            </div>
            <el-empty v-else description="暂无歌词" />
          </el-tab-pane>
          
          <el-tab-pane :label="`评论 (${work.comment_count})`" name="comments">
            <CommentList :work-id="work.id" />
          </el-tab-pane>
        </el-tabs>
      </div>
    </div>
    
    <el-dialog v-model="showPlaylistDialog" title="添加到歌单" width="400px">
      <div class="playlist-selector">
        <div 
          v-for="playlist in playlists" 
          :key="playlist.id"
          class="playlist-item"
          @click="addToPlaylist(playlist.id)"
        >
          {{ playlist.title }}
        </div>
        <el-empty v-if="playlists.length === 0" description="暂无歌单，请先创建" />
      </div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { workApi } from '@/api/work'
import { communityApi } from '@/api/community'
import { useUserStore } from '@/stores/user'
import CommentList from '@/components/CommentList.vue'
import type { Work } from '@/api/work'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const loading = ref(false)
const work = ref<Work | null>(null)
const activeTab = ref('intro')
const isLiked = ref(false)
const isFollowing = ref(false)
const showPlaylistDialog = ref(false)
const playlists = ref<any[]>([])

onMounted(async () => {
  await loadWork()
  if (userStore.isLoggedIn) {
    checkLikeStatus()
    checkFollowStatus()
    loadPlaylists()
  }
})

async function loadWork() {
  loading.value = true
  try {
    const id = parseInt(route.params.id as string)
    work.value = await workApi.getById(id)
    recordPlay()
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

async function recordPlay() {
  if (!work.value || !userStore.isLoggedIn) return
  try {
    await workApi.recordPlay(work.value.id, 0)
  } catch (e) {
    console.error(e)
  }
}

async function checkLikeStatus() {
}

async function checkFollowStatus() {
  if (!work.value) return
  try {
    const res = await communityApi.isFollowing(work.value.user_id)
    isFollowing.value = res.is_following
  } catch (e) {
    console.error(e)
  }
}

async function loadPlaylists() {
  try {
    const res = await communityApi.listPlaylists({ 
      user_id: userStore.user?.id,
      page: 1, 
      page_size: 50 
    })
    playlists.value = res.list
  } catch (e) {
    console.error(e)
  }
}

function play() {
  ElMessage.info('开始播放')
}

async function toggleLike() {
  ElMessage.info(isLiked.value ? '已取消收藏' : '已收藏')
  isLiked.value = !isLiked.value
}

async function followArtist() {
  if (!work.value || !userStore.isLoggedIn) {
    ElMessage.warning('请先登录')
    router.push('/login')
    return
  }
  
  try {
    if (isFollowing.value) {
      await communityApi.unfollow(work.value.user_id)
      isFollowing.value = false
      ElMessage.success('已取消关注')
    } else {
      await communityApi.follow(work.value.user_id)
      isFollowing.value = true
      ElMessage.success('关注成功')
    }
  } catch (e) {
    console.error(e)
  }
}

async function addToPlaylist(playlistId: number) {
  if (!work.value) return
  try {
    await communityApi.addWorkToPlaylist(playlistId, work.value.id)
    ElMessage.success('已添加到歌单')
    showPlaylistDialog.value = false
  } catch (e) {
    console.error(e)
  }
}

function goToArtist() {
  if (work.value) {
    router.push(`/artist/${work.value.user_id}`)
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
.work-detail {
  .work-content {
    .work-header {
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
      
      .work-info {
        flex: 1;
        display: flex;
        flex-direction: column;
        
        .title {
          font-size: 28px;
          font-weight: 600;
          margin: 0 0 12px 0;
        }
        
        .artist {
          display: flex;
          align-items: center;
          gap: 8px;
          margin-bottom: 16px;
          cursor: pointer;
          
          &:hover {
            color: var(--primary-color);
          }
        }
        
        .meta {
          display: flex;
          gap: 20px;
          color: var(--text-light);
          font-size: 14px;
          margin-bottom: 16px;
          
          span {
            display: flex;
            align-items: center;
            gap: 4px;
          }
        }
        
        .tags {
          margin-bottom: 20px;
        }
        
        .actions {
          display: flex;
          gap: 12px;
          margin-top: auto;
        }
      }
    }
    
    .work-body {
      .intro {
        white-space: pre-wrap;
        line-height: 1.8;
      }
      
      .info-list {
        margin-top: 30px;
        
        h3 {
          margin-bottom: 16px;
        }
        
        .info-item {
          display: flex;
          padding: 8px 0;
          border-bottom: 1px solid var(--border-color);
          
          .label {
            width: 100px;
            color: var(--text-light);
          }
          
          .value {
            flex: 1;
          }
        }
      }
      
      .lyrics {
        pre {
          white-space: pre-wrap;
          line-height: 2;
          font-size: 15px;
        }
      }
    }
  }
  
  .playlist-selector {
    max-height: 300px;
    overflow-y: auto;
    
    .playlist-item {
      padding: 12px;
      border-radius: 4px;
      cursor: pointer;
      transition: background 0.3s;
      
      &:hover {
        background: rgba(64, 158, 255, 0.1);
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
