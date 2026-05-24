<template>
  <div class="my-tickets">
    <div class="page-header">
      <h2>我的票</h2>
    </div>
    
    <el-tabs v-model="activeTab">
      <el-tab-pane label="待使用" name="upcoming" />
      <el-tab-pane label="已使用" name="used" />
      <el-tab-pane label="已取消" name="cancelled" />
    </el-tabs>
    
    <div class="tickets-list" v-loading="loading">
      <div 
        v-for="ticket in tickets" 
        :key="ticket.id"
        class="ticket-item"
      >
        <div class="ticket-info">
          <div class="event-title">{{ ticket.event?.title }}</div>
          <div class="event-meta">
            <el-icon><Location /></el-icon>
            <span>{{ ticket.event?.venue }}</span>
          </div>
          <div class="event-meta">
            <el-icon><Clock /></el-icon>
            <span>{{ formatTime(ticket.event?.start_time) }}</span>
          </div>
        </div>
        <div class="ticket-right">
          <div class="price">¥{{ ticket.price.toFixed(2) }}</div>
          <el-tag :type="getTicketTypeTag(ticket.status)" size="small">
            {{ getTicketStatusText(ticket.status) }}
          </el-tag>
          <div class="seat" v-if="ticket.seat_code">
            {{ ticket.seat_code }}
          </div>
        </div>
      </div>
      
      <el-empty v-if="tickets.length === 0 && !loading" description="暂无票" />
    </div>
    
    <div class="pagination" v-if="total > 0">
      <el-pagination
        v-model:current-page="page"
        v-model:page-size="pageSize"
        :total="total"
        layout="total, prev, pager, next"
        @current-change="loadTickets"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { eventApi } from '@/api/event'

const loading = ref(false)
const activeTab = ref('upcoming')
const tickets = ref<any[]>([])
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)

onMounted(() => {
  loadTickets()
})

async function loadTickets() {
  loading.value = true
  try {
    const res = await eventApi.getMyTickets({
      page: page.value,
      page_size: pageSize.value
    })
    tickets.value = res.list
    total.value = res.total
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

function formatTime(time: string) {
  if (!time) return '-'
  return time.replace('T', ' ').substring(0, 16)
}

function getTicketTypeTag(status: number): 'primary' | 'success' | 'warning' | 'info' | 'danger' {
  const tags: Record<number, 'primary' | 'success' | 'warning' | 'info' | 'danger'> = {
    0: 'warning',
    1: 'success',
    2: 'info',
    3: 'danger'
  }
  return tags[status] || 'info'
}

function getTicketStatusText(status: number) {
  const texts: Record<number, string> = {
    0: '待使用',
    1: '已使用',
    2: '已取消',
    3: '已退款'
  }
  return texts[status] || '未知'
}
</script>

<style scoped lang="scss">
.my-tickets {
  .tickets-list {
    .ticket-item {
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: 16px;
      border: 1px solid var(--border-color);
      border-radius: 8px;
      margin-bottom: 12px;
      
      .ticket-info {
        flex: 1;
        
        .event-title {
          font-weight: 500;
          font-size: 16px;
          margin-bottom: 8px;
        }
        
        .event-meta {
          display: flex;
          align-items: center;
          gap: 4px;
          font-size: 13px;
          color: var(--text-light);
          margin-bottom: 4px;
        }
      }
      
      .ticket-right {
        text-align: right;
        
        .price {
          font-size: 18px;
          font-weight: 600;
          color: var(--primary-color);
          margin-bottom: 8px;
        }
        
        .seat {
          margin-top: 8px;
          font-size: 13px;
          color: var(--text-light);
        }
      }
    }
  }
  
  .pagination {
    display: flex;
    justify-content: center;
    margin-top: 20px;
  }
}
</style>
