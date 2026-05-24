<template>
  <div class="events">
    <div class="page-header">
      <h1 class="page-title">演出</h1>
      <div class="search-bar">
        <el-input
          v-model="keyword"
          placeholder="搜索演出"
          clearable
          :prefix-icon="Search"
          @keyup.enter="search"
        />
        <el-select v-model="city" placeholder="城市" clearable class="city-select">
          <el-option label="北京" value="北京" />
          <el-option label="上海" value="上海" />
          <el-option label="广州" value="广州" />
          <el-option label="深圳" value="深圳" />
          <el-option label="杭州" value="杭州" />
          <el-option label="成都" value="成都" />
        </el-select>
        <el-button type="primary" @click="search">搜索</el-button>
      </div>
    </div>
    
    <div class="events-grid" v-loading="loading">
      <div 
        v-for="event in events" 
        :key="event.id"
        class="event-item"
        @click="goToDetail(event.id)"
      >
        <el-image :src="event.cover_url" class="cover" fit="cover">
          <template #error>
            <div class="image-placeholder">
              <el-icon :size="40"><Calendar /></el-icon>
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
          <div class="price">
            <span class="amount">¥{{ event.ticket_price }}</span>
            <span class="status" :class="{ sold: event.sold_tickets >= event.total_tickets }">
              {{ event.sold_tickets >= event.total_tickets ? '已售罄' : `剩余 ${event.total_tickets - event.sold_tickets}` }}
            </span>
          </div>
        </div>
      </div>
      
      <el-empty v-if="events.length === 0 && !loading" description="暂无演出" />
    </div>
    
    <div class="pagination" v-if="total > 0">
      <el-pagination
        v-model:current-page="page"
        v-model:page-size="pageSize"
        :total="total"
        :page-sizes="[12, 24, 48]"
        layout="total, sizes, prev, pager, next, jumper"
        @current-change="loadEvents"
        @size-change="loadEvents"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Search } from '@element-plus/icons-vue'
import { eventApi } from '@/api/event'
import type { Event } from '@/api/event'
import dayjs from 'dayjs'

const router = useRouter()
const loading = ref(false)
const keyword = ref('')
const city = ref('')
const events = ref<Event[]>([])
const page = ref(1)
const pageSize = ref(12)
const total = ref(0)

onMounted(() => {
  loadEvents()
})

async function loadEvents() {
  loading.value = true
  try {
    const res = await eventApi.list({
      page: page.value,
      page_size: pageSize.value,
      keyword: keyword.value,
      city: city.value,
      status: 1
    })
    events.value = res.list
    total.value = res.total
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

function search() {
  page.value = 1
  loadEvents()
}

function formatTime(time: string) {
  return dayjs(time).format('YYYY-MM-DD HH:mm')
}

function goToDetail(id: number) {
  router.push(`/event/${id}`)
}
</script>

<style scoped lang="scss">
.events {
  .search-bar {
    display: flex;
    gap: 12px;
    
    .city-select {
      width: 150px;
    }
  }
  
  .events-grid {
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    gap: 24px;
    
    .event-item {
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
        aspect-ratio: 16/9;
        border-radius: 8px;
        overflow: hidden;
        transition: box-shadow 0.3s;
      }
      
      .info {
        padding: 16px 4px;
        
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
        
        .price {
          display: flex;
          justify-content: space-between;
          align-items: center;
          margin-top: 12px;
          
          .amount {
            font-size: 18px;
            font-weight: 600;
            color: var(--danger-color);
          }
          
          .status {
            font-size: 12px;
            color: var(--success-color);
            
            &.sold {
              color: var(--danger-color);
            }
          }
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
