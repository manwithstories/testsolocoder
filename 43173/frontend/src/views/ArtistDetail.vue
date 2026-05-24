<template>
  <div class="artist-detail" v-loading="loading">
    <div v-if="artist" class="artist-content">
      <div class="artist-header">
        <el-avatar :size="120" :src="artist.avatar">
          {{ artist.nickname?.charAt(0) || 'A' }}
        </el-avatar>
        
        <div class="artist-info">
          <h1 class="name">{{ artist.nickname }}</h1>
          <div class="desc">{{ artist.artist_info?.artist_name || '独立音乐人' }}</div>
          
          <div class="stats">
            <div class="stat-item">
              <div class="value">{{ formatCount(artist.artist_info?.total_followers || 0) }}</div>
              <div class="label">粉丝</div>
            </div>
            <div class="stat-item">
              <div class="value">{{ formatCount(artist.artist_info?.total_works || 0) }}</div>
              <div class="label">作品</div>
            </div>
            <div class="stat-item">
              <div class="value">{{ formatCount(artist.artist_info?.total_plays || 0) }}</div>
              <div class="label">播放量</div>
            </div>
          </div>
          
          <div class="actions" v-if="userStore.isLoggedIn && artist.id !== userStore.user?.id">
            <el-button 
              :type="isFollowing ? 'default' : 'primary'" 
              size="large"
              @click="toggleFollow"
            >
              <el-icon><UserFilled /></el-icon>
              {{ isFollowing ? '已关注' : '关注' }}
            </el-button>
          </div>
        </div>
      </div>
      
      <div class="artist-body">
        <el-tabs v-model="activeTab">
          <el-tab-pane label="作品" name="works">
            <div class="works-grid" v-loading="worksLoading">
              <div 
                v-for="work in works" 
                :key="work.id"
                class="work-item"
                @click="goToWork(work.id)"
              >
                <el-image :src="work.cover_url" class="cover" fit="cover">
                  <template #error>
                    <div class="image-placeholder">
                      <el-icon :size="24"><Headset /></el-icon>
                    </div>
                  </template>
                </el-image>
                <div class="info">
                  <div class="title text-ellipsis">{{ work.title }}</div>
                  <div class="meta">
                    <el-icon><Headset /></el-icon>
                    {{ formatCount(work.play_count) }}
                  </div>
                </div>
              </div>
              
              <el-empty v-if="works.length === 0 && !worksLoading" description="暂无作品" />
            </div>
          </el-tab-pane>
          
          <el-tab-pane label="演出" name="events">
            <div class="events-list" v-loading="eventsLoading">
              <div 
                v-for="event in events" 
                :key="event.id"
                class="event-item"
                @click="goToEvent(event.id)"
              >
                <el-image :src="event.cover_url" class="cover" fit="cover">
                  <template #error>
                    <div class="image-placeholder">
                      <el-icon :size="24"><Calendar /></el-icon>
                    </div>
                  </template>
                </el-image>
                <div class="info">
                  <div class="title text-ellipsis">{{ event.title }}</div>
                  <div class="meta">
                    <el-icon><Location /></el-icon>
                    {{ event.venue }}
                  </div>
                  <div class="meta">
                    <el-icon><Clock /></el-icon>
                    {{ formatTime(event.start_time) }}
                  </div>
                </div>
              </div>
              
              <el-empty v-if="events.length === 0 && !eventsLoading" description="暂无演出" />
            </div>
          </el-tab-pane>
          
          <el-tab-pane label="简介" name="intro">
            <div class="intro" v-if="artist.bio || artist.artist_info?.genre">
              <p v-if="artist.bio">{{ artist.bio }}</p>
              <div class="info-list" v-if="artist.artist_info">
                <div class="info-item" v-if="artist.artist_info.genre">
                  <span class="label">风格</span>
                  <span class="value">{{ artist.artist_info.genre }}</span>
                </div>
                <div class="info-item" v-if="artist.artist_info.label">
                  <span class="label">厂牌</span>
                  <span class="value">{{ artist.artist_info.label }}</span>
                </div>
                <div class="info-item" v-if="artist.artist_info.website">
                  <span class="label">网站</span>
                  <span class="value">{{ artist.artist_info.website }}</span>
                </div>
              </div>
            </div>
            <el-empty v-else description="暂无简介" />
          </el-tab-pane>
        </el-tabs>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { userApi } from '@/api/auth'
import { workApi } from '@/api/work'
import { eventApi } from '@/api/event'
import { communityApi } from '@/api/community'
import { useUserStore } from '@/stores/user'
import dayjs from 'dayjs'
import type { UserInfo } from '@/api/auth'
import type { Work } from '@/api/work'
import type { Event } from '@/api/event'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const loading = ref(false)
const artist = ref<UserInfo | null>(null)
const activeTab = ref('works')
const isFollowing = ref(false)

const worksLoading = ref(false)
const works = ref<Work[]>([])

const eventsLoading = ref(false)
const events = ref<Event[]>([])

onMounted(async () => {
  await loadArtist()
  if (userStore.isLoggedIn && artist.value) {
    checkFollowStatus()
  }
  loadWorks()
  loadEvents()
})

async function loadArtist() {
  loading.value = true
  try {
    const id = parseInt(route.params.id as string)
    artist.value = await userApi.getById(id)
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

async function checkFollowStatus() {
  if (!artist.value) return
  try {
    const res = await communityApi.isFollowing(artist.value.id)
    isFollowing.value = res.is_following
  } catch (e) {
    console.error(e)
  }
}

async function loadWorks() {
  if (!artist.value?.artist_info) return
  worksLoading.value = true
  try {
    const res = await workApi.getByArtist(artist.value.artist_info.id, { page: 1, page_size: 20 })
    works.value = res.list
  } catch (e) {
    console.error(e)
  } finally {
    worksLoading.value = false
  }
}

async function loadEvents() {
  if (!artist.value?.artist_info) return
  eventsLoading.value = true
  try {
    const res = await eventApi.list({ 
      page: 1, 
      page_size: 20, 
      artist_id: artist.value.artist_info.id 
    })
    events.value = res.list
  } catch (e) {
    console.error(e)
  } finally {
    eventsLoading.value = false
  }
}

async function toggleFollow() {
  if (!artist.value || !userStore.isLoggedIn) {
    ElMessage.warning('请先登录')
    router.push('/login')
    return
  }
  
  try {
    if (isFollowing.value) {
      await communityApi.unfollow(artist.value.id)
      isFollowing.value = false
      ElMessage.success('已取消关注')
    } else {
      await communityApi.follow(artist.value.id)
      isFollowing.value = true
      ElMessage.success('关注成功')
    }
  } catch (e) {
    console.error(e)
  }
}

function goToWork(id: number) {
  router.push(`/work/${id}`)
}

function goToEvent(id: number) {
  router.push(`/event/${id}`)
}

function formatCount(count: number) {
  if (count >= 10000) {
    return (count / 10000).toFixed(1) + 'w'
  }
  return count
}

function formatTime(time: string) {
  return dayjs(time).format('YYYY-MM-DD HH:mm')
}
</script>

<style scoped lang="scss">
.artist-detail {
  .artist-content {
    .artist-header {
      display: flex;
      gap: 30px;
      margin-bottom: 30px;
      padding-bottom: 30px;
      border-bottom: 1px solid var(--border-color);
      
      .artist-info {
        flex: 1;
        
        .name {
          font-size: 28px;
          font-weight: 600;
          margin: 0 0 8px 0;
        }
        
        .desc {
          font-size: 16px;
          color: var(--text-light);
          margin-bottom: 16px;
        }
        
        .stats {
          display: flex;
          gap: 30px;
          margin-bottom: 20px;
          
          .stat-item {
            .value {
              font-size: 24px;
              font-weight: 600;
              color: var(--primary-color);
            }
            
            .label {
              font-size: 13px;
              color: var(--text-light);
            }
          }
        }
        
        .actions {
          display: flex;
          gap: 12px;
        }
      }
    }
    
    .artist-body {
      .works-grid {
        display: grid;
        grid-template-columns: repeat(6, 1fr);
        gap: 20px;
        
        .work-item {
          cursor: pointer;
          
          &:hover {
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
            padding: 8px 0;
            
            .title {
              font-size: 14px;
              margin-bottom: 4px;
            }
            
            .meta {
              font-size: 12px;
              color: var(--text-light);
              display: flex;
              align-items: center;
              gap: 4px;
            }
          }
        }
      }
      
      .events-list {
        display: grid;
        grid-template-columns: repeat(3, 1fr);
        gap: 20px;
        
        .event-item {
          cursor: pointer;
          
          &:hover {
            .cover {
              box-shadow: 0 8px 24px rgba(0, 0, 0, 0.15);
            }
          }
          
          .cover {
            width: 100%;
            aspect-ratio: 16/9;
            border-radius: 8px;
            overflow: hidden;
            transition: box-shadow 0.3s;
          }
          
          .info {
            padding: 12px 0;
            
            .title {
              font-weight: 500;
              margin-bottom: 8px;
            }
            
            .meta {
              font-size: 13px;
              color: var(--text-light);
              display: flex;
              align-items: center;
              gap: 4px;
              margin-bottom: 4px;
            }
          }
        }
      }
      
      .intro {
        line-height: 1.8;
        
        .info-list {
          margin-top: 20px;
          
          .info-item {
            display: flex;
            padding: 8px 0;
            border-bottom: 1px solid var(--border-color);
            
            .label {
              width: 80px;
              color: var(--text-light);
            }
            
            .value {
              flex: 1;
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
  }
}
</style>
