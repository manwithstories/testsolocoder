<template>
  <Layout>
    <div class="show-detail" v-if="show">
      <div class="show-header">
        <img :src="show.poster || 'https://picsum.photos/600/400'" class="show-poster" />
        <div class="show-info">
          <h1>{{ show.name }}</h1>
          <div class="meta">
            <p><strong>艺人：</strong>{{ show.artist }}</p>
            <p><strong>场馆：</strong>{{ show.venue }}</p>
            <p><strong>地址：</strong>{{ show.address }}</p>
            <p><strong>主办方：</strong>{{ show.organizer }}</p>
          </div>
          <p class="description">{{ show.description }}</p>
        </div>
        </div>

      <div class="sessions-section">
        <h3>选择场次</h3>
        <div class="sessions-list">
          <div
            v-for="session in sessions"
            :key="session.id"
            class="session-item"
            :class="{ active: selectedSession?.id === session.id }"
            @click="selectSession(session)"
          >
            <div class="session-time">
              <p class="date">{{ formatDate(session.start_time) }}</p>
              <p class="time">{{ formatTime(session.start_time) }} - {{ formatTime(session.end_time) }}</p>
            </div>
            <div class="session-info">
              <el-tag type="success" v-if="session.total_seats > 0">
                剩余 {{ session.total_seats - session.sold_seats }} 张
              </el-tag>
            </div>
            <el-button type="primary" size="small" @click.stop="goToSeatSelect(session)">
              选座购票
            </el-button>
          </div>
        </div>
      </div>
    </div>
  </Layout>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import dayjs from 'dayjs'
import { showApi } from '@/api'
import Layout from '@/components/Layout.vue'
import type { Show, Session } from '@/types'
import { useUserStore } from '@/store'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const show = ref<Show | null>(null)
const sessions = ref<Session[]>([])
const selectedSession = ref<Session | null>(null)

async function fetchShow() {
  const id = Number(route.params.id as string)
  try {
    show.value = await showApi.get(id)
  } catch (err) {
    console.error(err)
  }
}

async function fetchSessions() {
  const id = Number(route.params.id as string)
  try {
    sessions.value = await showApi.getSessions(id)
  } catch (err) {
    console.error(err)
  }
}

function selectSession(session: Session) {
  selectedSession.value = session
}

function goToSeatSelect(session: Session) {
  if (!userStore.isLoggedIn) {
    ElMessage.warning('请先登录')
    router.push('/login')
    return
  }
  router.push(`/seat-select/${session.id}`)
}

function formatDate(date: string) {
  return dayjs(date).format('YYYY-MM-DD')
}

function formatTime(date: string) {
  return dayjs(date).format('HH:mm')
}

onMounted(() => {
  fetchShow()
  fetchSessions()
})
</script>

<style lang="scss" scoped>
.show-detail {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

.show-header {
  display: flex;
  gap: 30px;
  background: white;
  padding: 24px;
  border-radius: 12px;
  margin-bottom: 24px;

  .show-poster {
    width: 400px;
    height: 300px;
    object-fit: cover;
    border-radius: 8px;
  }

  .show-info {
    flex: 1;

    h1 {
      margin: 0 0 16px 0;
      font-size: 28px;
    }

    .meta p {
      margin: 8px 0;
      color: #666;
    }

    .description {
      margin-top: 16px 0 0 0;
      color: #555;
      line-height: 1.8;
    }
  }
}

.sessions-section {
  background: white;
  padding: 24px;
  border-radius: 12px;

  h3 {
    margin: 0 0 16px 0;
  }
}

.sessions-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.session-item {
  display: flex;
  align-items: center;
  padding: 16px;
  border: 1px solid #e4e7ed;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;

  &:hover, &.active {
    border-color: #409eff;
    background: #ecf5ff;
  }

  .session-time {
    flex: 1;

    .date {
      font-size: 16px;
      font-weight: 600;
      margin: 0 0 4px 0;
    }

    .time {
      margin: 0;
      color: #666;
    }
  }

  .session-info {
    margin-right: 16px;
  }
}
</style>
